package app

import (
	"net/http"
)

type server struct {
	// a local logger
	// a database connection
	// your app configuration
	httpServer *http.Server
}
