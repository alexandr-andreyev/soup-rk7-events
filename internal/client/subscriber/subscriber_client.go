package subscriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

// SubscriberClient handles sending events to a single subscriber
type SubscriberClient struct {
	httpClient *http.Client
	config     *config.Subscriber
}

// NewSubscriberClient creates a new client for a subscriber
func NewSubscriberClient(cfg *config.Subscriber) *SubscriberClient {
	return &SubscriberClient{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		config: cfg,
	}
}

// SendEvent sends an event to the subscriber with retry logic
func (c *SubscriberClient) SendEvent(event *models.ExternalEvent) error {
	if !c.config.Enabled {
		slog.Debug("Subscriber disabled, skipping", "subscriber", c.config.Name)
		return nil
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	var lastErr error
	maxRetries := c.config.RetryCount
	if maxRetries == 0 {
		maxRetries = 1 // At least one attempt
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			slog.Info("Retry attempt",
				"subscriber", c.config.Name,
				"attempt", attempt,
				"maxRetries", maxRetries,
				"guid", event.ResponseEventCommon.EventGUID)
			time.Sleep(time.Duration(c.config.RetryDelaySeconds) * time.Second)
		}

		req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Set standard headers
		req.Header.Set("Content-Type", "application/json")

		// Set custom headers from config
		for key, value := range c.config.Headers {
			req.Header.Set(key, value)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			slog.Warn("Send attempt failed",
				"subscriber", c.config.Name,
				"attempt", attempt,
				"error", err)
			continue
		}

		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Check if request was successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			slog.Info("Event sent to subscriber",
				"subscriber", c.config.Name,
				"guid", event.ResponseEventCommon.EventGUID,
				"status", resp.StatusCode)
			return nil
		}

		lastErr = fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		slog.Warn("Send attempt failed with non-success status",
			"subscriber", c.config.Name,
			"attempt", attempt,
			"status", resp.StatusCode,
			"body", string(body))

		// Don't retry on 4xx errors (client errors)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return lastErr
		}
	}

	return fmt.Errorf("failed to send to subscriber %s after %d attempts: %w",
		c.config.Name, maxRetries, lastErr)
}
