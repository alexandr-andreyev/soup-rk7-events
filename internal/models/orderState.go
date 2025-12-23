package models

import "time"

// OrderState represents the current state of an order in the database
type OrderState struct {
	ID              uint      `gorm:"primaryKey"`
	OriginalOrderId string    `gorm:"uniqueIndex;not null"` // Unique order identifier
	OrderGUID       string    `gorm:"index"`                // Order GUID for reference
	Kdsstate        string    `gorm:"not null"`             // Current KDS state
	LastUpdated     time.Time `gorm:"autoUpdateTime"`       // Last update timestamp
	CreatedAt       time.Time `gorm:"autoCreateTime"`       // Creation timestamp
}
