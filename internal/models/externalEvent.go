package models

import (
	"time"

	"github.com/google/uuid"
)

// ExternalEvent represents the JSON payload sent to external API
type ExternalEvent struct {
	ResponseEventCommon ResponseEventCommon `json:"responseEventCommon"`
	Response            Response            `json:"response"`
}

type ResponseEventCommon struct {
	ObjectID                            string `json:"objectId"`
	AgentGUID                           string `json:"agentGuid"`
	EventGUID                           string `json:"eventGuid"`
	DateTimeServerReceiveEventFromAgent string `json:"dateTimeServerReceiveEventFromAgent"`
	EventType                           string `json:"eventType"`
}

type Response struct {
	OriginalOrderID string           `json:"originalOrderId"`
	OrderGUID       string           `json:"orderGuid"`
	Status          *Status          `json:"status,omitempty"`
	Substate        string           `json:"substate,omitempty"`
	RejectingReason *RejectingReason `json:"rejectingReason,omitempty"`
}

type Status struct {
	Value         string `json:"value"`
	IsBillPrinted bool   `json:"isBillPrinted"`
}

type RejectingReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ToExternalEvent converts RK7 notification to external event format
func (rk7 *Rk7NotifyEvent) ToExternalEvent() *ExternalEvent {
	event := &ExternalEvent{
		ResponseEventCommon: ResponseEventCommon{
			ObjectID:                            rk7.RestCode,
			AgentGUID:                           rk7.RestCode + "-agent",
			EventGUID:                           uuid.New().String(),
			DateTimeServerReceiveEventFromAgent: time.Now().Format(time.RFC3339),
			EventType:                           mapEventType(rk7.Name),
		},
		Response: Response{
			OriginalOrderID: "",
			OrderGUID:       "",
			Substate:        "",
		},
	}

	// Add order info if present
	if rk7.Order != nil {
		event.Response.OrderGUID = rk7.Order.GUID
		event.Response.OriginalOrderID = rk7.Order.OrderIdent

		// Map KDS state to status
		if rk7.Order.Kdsstate != "" {
			event.Response.Status = &Status{
				Value:         rk7.Order.Kdsstate,
				IsBillPrinted: false, // This should be determined from actual RK7 data
			}
		}
	}

	return event
}

// mapEventType maps RK7 event names to external event types
func mapEventType(rk7EventName string) string {
	switch rk7EventName {
	case "Order Changed":
		return "OrderStatusChanged"
	case "New Order":
		return "OrderCreated"
	case "Print Receipt":
		return "OrderCompleted"
	default:
		return rk7EventName
	}
}
