package app

import (
	"net/http"
)

// The wrapper of your app
func runApp(s server) {
	s.Logger.Info("Soup Event Service is started")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Logger.Error("Soup Event Service HTTP server error", "error", err)
	}
}
