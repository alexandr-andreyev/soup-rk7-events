package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/database"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/services"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/transport"
)

// if setup returns an error, the service doesn't start
func setup(svcName string, build BuildInfo) (server, error) {
	var s server

	slog.Info("Starting service", "version", build.Version, "commit", build.Commit, "buildTime", build.BuildTime)

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
				Path: "events.db",
			},
			Logging: config.LoggingConfig{
				Directory: "logs",
			},
		}
	}

	// Initialize file logging
	s.Logger, err = initFileLogger(cfg.Logging.Directory)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		return s, err
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
		s.Logger,
		dispatcher,
		orderRepo,
	)
	s.httpServer = setupHttpServer(s.Logger, handleEventService, cfg)
	return s, nil
}

func setupHttpServer(logger *slog.Logger, handleEventService services.NotifyEventService, cfg *config.Config) *http.Server {
	logger.Info("Setting up HTTP server", "port", cfg.Server.Port)
	handler := transport.NewHandler(logger, handleEventService)

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
