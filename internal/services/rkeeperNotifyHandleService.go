package services

import (
	"log/slog"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/database"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"gorm.io/gorm"
)

type NotifyEventService interface {
	HandleNotification(data *models.Rk7NotifyEvent) error
}

type RkNotifyHandleService struct {
	dispatcher *EventDispatcher
	orderRepo  *database.OrderStateRepository
}

func NewRkNotifyHandleService(dispatcher *EventDispatcher, orderRepo *database.OrderStateRepository) *RkNotifyHandleService {
	return &RkNotifyHandleService{
		dispatcher: dispatcher,
		orderRepo:  orderRepo,
	}
}

func (s *RkNotifyHandleService) HandleNotification(data *models.Rk7NotifyEvent) error {
	// Log received event
	slog.Info("Received event", "name", data.Name, "guid", data.GUID)

	// check subscribers exist

	// Process event based on type
	switch data.Name {
	case "Started":
		slog.Info("Started event", "restCode", data.RestCode)
		// Don't forward system events
		return nil

	case "New Order":
		slog.Info("New Order event", "data", data)
		// Forward to external API
		// return s.sendToExternalAPI(data)

	case "Order Changed":
		if data.Order == nil {
			slog.Info("Ignoring Order Changed event with no Order data")
			return nil
		}
		if data.Order.Kdsstate == "" {
			slog.Info("Ignoring Order Changed event with empty Kdsstate")
			return nil
		}
		return s.handleOrderStatusChange(data)

	case "Print Receipt":
		if data.Order == nil {
			slog.Info("Ignoring Print Receipt event with no Order data")
			return nil
		}
		return s.handleOrderStatusChange(data)

	default:
		slog.Info("Unhandled event", "name", data.Name)
		return nil
	}
	return nil
}

// handleOrderStatusChange checks if order status changed and sends event if needed
func (s *RkNotifyHandleService) handleOrderStatusChange(data *models.Rk7NotifyEvent) error {
	orderGUID := data.Order.GUID
	newStatus := data.Order.Kdsstate

	slog.Info("Processing order", "guid", orderGUID, "status", newStatus)

	// Try to get existing order state from DB
	existingState, err := s.orderRepo.GetByOrderGuid(orderGUID)

	statusChanged := false
	if err != nil {
		// Check if error is "record not found"
		if err == gorm.ErrRecordNotFound {
			// Order not found in DB - add new order
			slog.Info("New order detected, adding to DB", "guid", orderGUID, "status", newStatus)

			newOrderState := &models.OrderState{
				OriginalOrderId: orderGUID,
				OrderGUID:       orderGUID,
				Kdsstate:        newStatus,
			}

			if err := s.orderRepo.Create(newOrderState); err != nil {
				slog.Error("Failed to create order state in DB", "error", err)
				return err
			}

			statusChanged = true
		} else {
			// Other database error
			slog.Error("Database error while getting order state", "error", err)
			return err
		}
	} else {
		// Order found in DB - check if status changed
		if existingState.Kdsstate != newStatus {
			// Status changed - update DB
			slog.Info("Order status changed", "guid", orderGUID, "oldStatus", existingState.Kdsstate, "newStatus", newStatus)

			if err := s.orderRepo.UpdateStatus(orderGUID, newStatus); err != nil {
				slog.Error("Failed to update order state in DB", "error", err)
				return err
			}

			statusChanged = true
		} else {
			slog.Info("Order status unchanged", "guid", orderGUID, "status", newStatus)
		}
	}

	// Dispatch events based on what changed
	var dispatchErrors []error

	// 1. If status changed, dispatch OrderStatusChanged event
	if statusChanged {
		if err := s.dispatcher.Dispatch(data, "OrderStatusChanged"); err != nil {
			slog.Error("Failed to dispatch OrderStatusChanged event", "error", err)
			dispatchErrors = append(dispatchErrors, err)
		}
	}

	// 2. Always dispatch OrderChanged event (content may have changed)
	// Subscribers can filter whether they want this event type
	if err := s.dispatcher.Dispatch(data, "OrderChanged"); err != nil {
		slog.Error("Failed to dispatch OrderChanged event", "error", err)
		dispatchErrors = append(dispatchErrors, err)
	}

	// Return first error if any occurred
	if len(dispatchErrors) > 0 {
		return dispatchErrors[0]
	}

	return nil
}

func (s *RkNotifyHandleService) sendToExternalAPI(data *models.Rk7NotifyEvent) error {
	// Dispatch to all matching subscribers (uses default event.Name for matching)
	if err := s.dispatcher.Dispatch(data, ""); err != nil {
		slog.Error("Failed to dispatch event", "error", err)
		return err
	}

	slog.Info("Event dispatched successfully", "guid", data.GUID)
	return nil
}
