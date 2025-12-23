package app

import (
	"log/slog"
	"net/http"
)

type server struct {
	// a local logger
	Logger *slog.Logger
	// a database connection
	// your app configuration
	httpServer *http.Server
}
