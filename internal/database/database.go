package database

import (
	"log/slog"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// InitDB initializes SQLite database with GORM
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.OrderState{}); err != nil {
		return nil, err
	}

	slog.Info("Database initialized successfully")
	return db, nil
}
