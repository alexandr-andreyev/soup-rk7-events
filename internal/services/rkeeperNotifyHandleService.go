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
	Logger     *slog.Logger
	dispatcher *EventDispatcher
	orderRepo  *database.OrderStateRepository
}

func NewRkNotifyHandleService(logger *slog.Logger, dispatcher *EventDispatcher, orderRepo *database.OrderStateRepository) *RkNotifyHandleService {
	return &RkNotifyHandleService{
		Logger:     logger,
		dispatcher: dispatcher,
		orderRepo:  orderRepo,
	}
}

func (s *RkNotifyHandleService) HandleNotification(data *models.Rk7NotifyEvent) error {
	const op = "RkNotifyHandleService.HandleNotification"

	logger := s.Logger.With("op", op)
	// Log received event
	logger.Info("Received RK7 Notify Event", "name", data.Name, "guid", data.GUID)

	// check subscribers exist

	// Process event based on type
	switch data.Name {
	case "Started":
		logger.Info("Ignoring Started event")
		// Don't forward system events
		return nil

	case "New Order":
		logger.Info("Processing New Order event", "data", data)
		// Forward to external API
		// return s.sendToExternalAPI(data)

	case "Order Changed":
		if data.Order == nil {
			logger.Info("Ignoring Order Changed event with no Order data", "order", data.GUID)
			return nil
		}
		if data.Order.Kdsstate == "" {
			logger.Info("Ignoring Order Changed event with empty Kdsstate", "order", data.GUID, "status", data.Order.Kdsstate)
			return nil
		}
		return s.handleOrderStatusChange(data)

	case "Print Receipt":
		if data.Order == nil {
			logger.Info("Ignoring Print Receipt event with no Order data", "order", data.GUID)
			return nil
		}
		return s.handleOrderStatusChange(data)

	default:
		logger.Info("Ignoring unhandled event type", "name", data.Name)
		logger.Debug("Ignoring event data", "data", data)
		return nil
	}
	return nil
}

// handleOrderStatusChange checks if order status changed and sends event if needed
func (s *RkNotifyHandleService) handleOrderStatusChange(data *models.Rk7NotifyEvent) error {
	const op = "RkNotifyHandleService.handleOrderStatusChange"
	slog := s.Logger.With("op", op)

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
