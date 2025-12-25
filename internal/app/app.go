package app

import (
	"log/slog"
	"net/http"
)

// The wrapper of your app
func runApp(s server) {
	slog.Info("Starting HTTP server on :7080")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server error", "error", err)
	}
}
