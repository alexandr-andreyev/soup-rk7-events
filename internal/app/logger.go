package app

import (
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

	// Wrap file with CRLF writer
	writer := &crlfWriter{w: file}

	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return logger, nil
}
