package database

import (
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"gorm.io/gorm"
)

type OrderStateRepository struct {
	db *gorm.DB
}

func NewOrderStateRepository(db *gorm.DB) *OrderStateRepository {
	return &OrderStateRepository{db: db}
}

// GetByOrderGuid retrieves an order state by its order GUID
func (r *OrderStateRepository) GetByOrderGuid(orderGuid string) (*models.OrderState, error) {
	var orderState models.OrderState
	err := r.db.Where("order_guid = ?", orderGuid).First(&orderState).Error
	if err != nil {
		return nil, err
	}
	return &orderState, nil
}

// Create creates a new order state record
func (r *OrderStateRepository) Create(orderState *models.OrderState) error {
	return r.db.Create(orderState).Error
}

// UpdateStatus updates the kdsstate for an existing order
func (r *OrderStateRepository) UpdateStatus(orderGuid, newStatus string) error {
	return r.db.Model(&models.OrderState{}).
		Where("order_guid = ?", orderGuid).
		Update("kdsstate", newStatus).Error
}
