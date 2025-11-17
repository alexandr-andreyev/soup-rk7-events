package app

import (
	"log"
	"net/http"
)

// The wrapper of your app
func runApp(s server) {
	log.Println("Starting HTTP server on :7080")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
