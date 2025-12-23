package app

import (
	"net/http"
	"time"

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

	// read configuration
	// configure more logging
	handleEventService := services.NewRkNotifyHandleService()
	s.httpServer = setupHttpServer(handleEventService)
	return s, nil
}

func setupHttpServer(handleEventService services.NotifyEventService) *http.Server {
	handler := transport.NewHandler(handleEventService)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", handler.HandleEvents)

	s := &http.Server{
		Addr:         ":7080", // порт, куда будет постить внешняя система
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}
