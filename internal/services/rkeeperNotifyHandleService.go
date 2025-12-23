package services

import (
	"log"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/client"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/database"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"gorm.io/gorm"
)

type NotifyEventService interface {
	HandleNotification(data *models.Rk7NotifyEvent) error
}

type RkNotifyHandleService struct {
	externalClient *client.ExternalClient
	orderRepo      *database.OrderStateRepository
}

func NewRkNotifyHandleService(externalClient *client.ExternalClient, orderRepo *database.OrderStateRepository) *RkNotifyHandleService {
	return &RkNotifyHandleService{
		externalClient: externalClient,
		orderRepo:      orderRepo,
	}
}

func (s *RkNotifyHandleService) HandleNotification(data *models.Rk7NotifyEvent) error {
	// Log received event
	log.Printf("Received event: %s (GUID: %s)", data.Name, data.GUID)

	// Process event based on type
	switch data.Name {
	case "Started":
		log.Println("event Name:", data.Name)
		log.Println("RestCode:", data.RestCode)
		// Don't forward system events
		return nil

	case "New Order":
		log.Println("event Name:", data.Name)
		log.Printf("Data: %+v\n", data)
		// Forward to external API
		// return s.sendToExternalAPI(data)

	case "Order Changed":
		if data.Order == nil {
			log.Println("Ignoring Order Changed event with no Order data")
			return nil
		}
		if data.Order.Kdsstate == "" {
			log.Println("Ignoring Order Changed event with empty Kdsstate")
			return nil
		}
		return s.handleOrderStatusChange(data)

	case "Print Receipt":
		if data.Order == nil {
			log.Println("Ignoring Print Receipt event with no Order data")
			return nil
		}
		return s.handleOrderStatusChange(data)

	default:
		log.Println("Unhandled event:", data.Name)
		return nil
	}
	return nil
}

// handleOrderStatusChange checks if order status changed and sends event if needed
func (s *RkNotifyHandleService) handleOrderStatusChange(data *models.Rk7NotifyEvent) error {
	orderGUID := data.Order.GUID
	newStatus := data.Order.Kdsstate

	log.Printf("Processing order %s with status: %s", orderGUID, newStatus)

	// Try to get existing order state from DB
	existingState, err := s.orderRepo.GetByOrderGuid(orderGUID)

	if err != nil {
		// Check if error is "record not found"
		if err == gorm.ErrRecordNotFound {
			// Order not found in DB - add new order
			log.Printf("New order detected: %s, adding to DB with status: %s", orderGUID, newStatus)

			newOrderState := &models.OrderState{
				OriginalOrderId: data.Order.OrderName,
				OrderGUID:       data.Order.GUID,
				Kdsstate:        newStatus,
			}

			if err := s.orderRepo.Create(newOrderState); err != nil {
				log.Printf("Failed to create order state in DB: %v", err)
				return err
			}

			// Send event to external API
			return s.sendToExternalAPI(data)
		}

		// Other database error
		log.Printf("Database error while getting order state: %v", err)
		return err
	}

	// Order found in DB - check if status changed
	if existingState.Kdsstate == newStatus {
		log.Printf("Order %s status unchanged (%s), ignoring event", orderGUID, newStatus)
		return nil
	}

	// Status changed - update DB and send event
	log.Printf("Order %s status changed: %s -> %s", orderGUID, existingState.Kdsstate, newStatus)

	if err := s.orderRepo.UpdateStatus(orderGUID, newStatus); err != nil {
		log.Printf("Failed to update order state in DB: %v", err)
		return err
	}

	// Send event to external API
	return s.sendToExternalAPI(data)
}

func (s *RkNotifyHandleService) sendToExternalAPI(data *models.Rk7NotifyEvent) error {
	// Convert to external event format
	externalEvent := data.ToExternalEvent()

	// Send to external API
	if err := s.externalClient.SendEvent(externalEvent); err != nil {
		log.Printf("Failed to send event to external API: %v", err)
		return err
	}

	log.Printf("Event forwarded successfully: %s", externalEvent.ResponseEventCommon.EventGUID)
	return nil
}
