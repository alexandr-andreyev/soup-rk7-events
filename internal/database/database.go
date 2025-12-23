package database

import (
	"log"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes SQLite database with GORM
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.OrderState{}); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
