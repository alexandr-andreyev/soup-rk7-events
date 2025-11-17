package services

import (
	"log"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

type NotifyEventService interface {
	HandleNotification(data *models.Rk7NotifyEvent) error
}

type RkNotifyHandleService struct{}

func NewRkNotifyHandleService() *RkNotifyHandleService {
	return &RkNotifyHandleService{}
}

func (s *RkNotifyHandleService) HandleNotification(data *models.Rk7NotifyEvent) error {
	// Здесь будет логика обработки уведомления
	log.Println("Handling notification for event:", data.Name)
	log.Println("RestCode:", data.RestCode)
	return nil
}
