package services

import (
	"log/slog"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/client/subscriber"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

// EventDispatcher routes events to multiple subscribers based on filters
type EventDispatcher struct {
	subscribers []*subscriberWithFilters
}

type subscriberWithFilters struct {
	client     *subscriber.SubscriberClient
	config     *config.Subscriber
	eventTypes []string
	tableCodes map[int]bool
}

// NewEventDispatcher creates a new dispatcher from config
func NewEventDispatcher(cfg *config.Config) *EventDispatcher {
	var subscribers []*subscriberWithFilters

	for i := range cfg.Subscribers {
		sub := &cfg.Subscribers[i]

		// Skip disabled subscribers
		if !sub.Enabled {
			slog.Info("Skipping disabled subscriber", "name", sub.Name)
			continue
		}

		// Create table codes map for fast lookup
		tableCodes := make(map[int]bool)
		for _, code := range sub.TableCodes {
			tableCodes[code] = true
		}

		// Normalize event type patterns from config for consistent matching
		normalizedEventTypes := make([]string, len(sub.EventTypes))
		for j, eventType := range sub.EventTypes {
			normalizedEventTypes[j] = normalizeEventType(eventType)
		}

		subscribers = append(subscribers, &subscriberWithFilters{
			client:     subscriber.NewSubscriberClient(sub),
			config:     sub,
			eventTypes: normalizedEventTypes,
			tableCodes: tableCodes,
		})

		slog.Info("Registered subscriber",
			"name", sub.Name,
			"url", sub.URL,
			"event_types", sub.EventTypes,
			"table_codes", sub.TableCodes)
	}

	return &EventDispatcher{
		subscribers: subscribers,
	}
}

// Dispatch sends an event to all matching subscribers
// If customEventType is empty, uses event.Name for matching
func (d *EventDispatcher) Dispatch(event *models.Rk7NotifyEvent, customEventType string) error {
	if len(d.subscribers) == 0 {
		slog.Debug("No subscribers configured")
		return nil
	}

	// Use custom event type if provided, otherwise derive from RK7 event name
	var eventType string
	if customEventType != "" {
		eventType = normalizeEventType(customEventType)
	} else {
		eventType = normalizeEventType(event.Name)
	}

	// Get table code if available
	var tableCode int
	if event.Order != nil && event.Order.Table != nil && event.Order.Table.Code != "" {
		if code, err := strconv.Atoi(event.Order.Table.Code); err == nil {
			tableCode = code
		}
	}

	// Find matching subscribers
	var matchingSubscribers []*subscriberWithFilters
	for _, sub := range d.subscribers {
		if d.shouldDispatch(sub, eventType, tableCode) {
			matchingSubscribers = append(matchingSubscribers, sub)
		}
	}

	if len(matchingSubscribers) == 0 {
		slog.Debug("No matching subscribers for event",
			"event_type", eventType,
			"table_code", tableCode)
		return nil
	}

	slog.Info("Dispatching event to subscribers",
		"event_type", eventType,
		"table_code", tableCode,
		"subscriber_count", len(matchingSubscribers))

	// Convert to external event format
	externalEvent := event.ToExternalEvent()

	// Send to all matching subscribers in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, sub := range matchingSubscribers {
		wg.Add(1)
		go func(s *subscriberWithFilters) {
			defer wg.Done()

			if err := s.client.SendEvent(externalEvent); err != nil {
				slog.Error("Failed to send to subscriber",
					"subscriber", s.config.Name,
					"error", err)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	// Return first error if any occurred
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

// shouldDispatch determines if an event should be sent to a subscriber
func (d *EventDispatcher) shouldDispatch(sub *subscriberWithFilters, eventType string, tableCode int) bool {
	// Check table code filter (if specified)
	if len(sub.tableCodes) > 0 {
		if tableCode == 0 || !sub.tableCodes[tableCode] {
			slog.Debug("Event filtered by table code",
				"subscriber", sub.config.Name,
				"table_code", tableCode,
				"allowed_codes", sub.config.TableCodes)
			return false
		}
	}

	// If no event types specified, accept all events
	if len(sub.eventTypes) == 0 {
		slog.Debug("Subscriber accepts all events", "subscriber", sub.config.Name)
		return true
	}

	// Check if event type matches any pattern
	for _, pattern := range sub.eventTypes {
		if matchEventType(pattern, eventType) {
			slog.Debug("Event matched pattern",
				"subscriber", sub.config.Name,
				"pattern", pattern,
				"event_type", eventType)
			return true
		}
	}

	slog.Debug("Event filtered by event type",
		"subscriber", sub.config.Name,
		"event_type", eventType,
		"patterns", sub.eventTypes)
	return false
}

// normalizeEventType converts RK7 event name to normalized event type
// Examples: "Order Changed" -> "order.changed", "Print Receipt" -> "print.receipt"
func normalizeEventType(name string) string {
	// Convert to lowercase and replace spaces with dots
	normalized := strings.ToLower(name)
	normalized = strings.ReplaceAll(normalized, " ", ".")
	return normalized
}

// matchEventType checks if event type matches a pattern (supports wildcards)
// Examples:
//   - "order.changed" matches "order.changed"
//   - "order.changed" matches "order.*"
//   - "order.changed" matches "*"
func matchEventType(pattern, eventType string) bool {
	// Exact match
	if pattern == eventType {
		return true
	}

	// Wildcard match (e.g., "order.*" matches "order.changed")
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		if strings.HasPrefix(eventType, prefix+".") {
			return true
		}
	}

	// Match all wildcard
	if pattern == "*" {
		return true
	}

	// Path-style wildcard match (e.g., "order*" matches "order.changed")
	matched, _ := path.Match(pattern, eventType)
	return matched
}
