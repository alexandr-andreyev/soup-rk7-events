package app

import (
	"log"
	"net/http"
	"time"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/client"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/database"
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
		log.Printf("Failed to load config: %v, using defaults", err)
		// Use default config
		cfg = &config.Config{
			Server: config.ServerConfig{
				Port:         ":7080",
				ReadTimeout:  10,
				WriteTimeout: 10,
			},
			External: config.ExternalConfig{
				URL:            "",
				TimeoutSeconds: 10,
				MaxRetries:     3,
				RetryDelay:     1,
			},
			Database: config.DatabaseConfig{
				Path: "orders.db",
			},
		}
	}

	// Initialize database
	db, err := database.InitDB(cfg.Database.Path)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return s, err
	}

	// Create repository
	orderRepo := database.NewOrderStateRepository(db)

	// Create external client
	externalClient := client.NewExternalClient(&cfg.External)

	// Create service with external client and repository
	handleEventService := services.NewRkNotifyHandleService(
		externalClient,
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
