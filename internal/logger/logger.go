package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
)

// crlfWriter wraps an io.Writer and converts LF to CRLF for Windows compatibility
type crlfWriter struct {
	w io.Writer
}

func (cw *crlfWriter) Write(p []byte) (n int, err error) {
	// Replace LF with CRLF
	converted := bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	_, err = cw.w.Write(converted)
	if err != nil {
		return 0, err
	}
	// Return original length to satisfy io.Writer interface
	return len(p), nil
}

// InitFileLogger initializes file logging to the logs directory
func InitFileLogger(logsDir string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with date in filename
	logFileName := fmt.Sprintf("soup-rk7-events_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logsDir, logFileName)

	// Open log file in append mode
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Wrap file writer with CRLF converter for Windows Notepad compatibility
	fileWriter := &crlfWriter{w: logFile}

	// Set log output to both file (with CRLF) and stdout (without CRLF)
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// Create text handler with custom options
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(multiWriter, opts)

	// Set as default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("File logging initialized", "path", logFilePath)
	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
