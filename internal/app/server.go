package app

import (
	"net/http"

	"golang.org/x/sys/windows/svc/debug"
)

type server struct {
	winlog debug.Log
	// a local logger
	// a database connection
	// your app configuration
	httpServer *http.Server
}
