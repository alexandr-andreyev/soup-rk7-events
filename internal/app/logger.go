package app

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// crlfWriter wraps an io.Writer and converts LF to CRLF
type crlfWriter struct {
	w io.Writer
}

func (c *crlfWriter) Write(p []byte) (n int, err error) {
	// Replace LF with CRLF for Windows
	var result []byte
	for i := 0; i < len(p); i++ {
		if p[i] == '\n' && (i == 0 || p[i-1] != '\r') {
			result = append(result, '\r', '\n')
		} else {
			result = append(result, p[i])
		}
	}
	written, err := c.w.Write(result)
	if err != nil {
		return written, err
	}
	return len(p), nil
}

// multiHandler writes to multiple slog handlers
type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func initFileLogger(logDir string) (*slog.Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Create log file with date-based name
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logPath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	// File handler with CRLF line endings
	fileHandler := slog.NewTextHandler(&crlfWriter{w: file}, opts)

	// Console handler without CRLF conversion
	consoleHandler := slog.NewTextHandler(os.Stdout, opts)

	logger := slog.New(&multiHandler{handlers: []slog.Handler{fileHandler, consoleHandler}})

	return logger, nil
}
