package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/database"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/logger"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/services"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/transport"
)

// if setup returns an error, the service doesn't start
func setup(svcName, sha1ver string) (server, error) {
	var s server

	// did we get a full SHA1?
	if len(sha1ver) == 40 {
		sha1ver = sha1ver[0:7]
	}

	if sha1ver == "" {
		sha1ver = "dev"
	}

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Warn("Failed to load config, using defaults", "error", err)
		// Use default config
		cfg = &config.Config{
			Server: config.ServerConfig{
				Port:         ":7080",
				ReadTimeout:  10,
				WriteTimeout: 10,
			},
			Subscribers: []config.Subscriber{},
			Database: config.DatabaseConfig{
				Path: "orders.db",
			},
			Logging: config.LoggingConfig{
				Directory: "logs",
			},
		}
	}

	// Initialize file logging
	if err := logger.InitFileLogger(cfg.Logging.Directory); err != nil {
		slog.Warn("Failed to initialize file logging", "error", err)
		// Continue without file logging
	}

	// Initialize database
	db, err := database.InitDB(cfg.Database.Path)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return s, err
	}

	// Create repository
	orderRepo := database.NewOrderStateRepository(db)

	// Create event dispatcher with subscribers from config
	dispatcher := services.NewEventDispatcher(cfg)

	// Create service with dispatcher and repository
	handleEventService := services.NewRkNotifyHandleService(
		dispatcher,
		orderRepo,
	)
	s.httpServer = setupHttpServer(handleEventService, cfg)
	return s, nil
}

func setupHttpServer(handleEventService services.NotifyEventService, cfg *config.Config) *http.Server {
	handler := transport.NewHandler(handleEventService)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", handler.HandleEvents)

	s := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return s
}
