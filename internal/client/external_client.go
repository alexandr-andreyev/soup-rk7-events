package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/config"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

type ExternalClient struct {
	httpClient *http.Client
	config     *config.ExternalConfig
}

func NewExternalClient(cfg *config.ExternalConfig) *ExternalClient {
	return &ExternalClient{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
		config: cfg,
	}
}

// SendEvent sends an event to the external API with retry logic
func (c *ExternalClient) SendEvent(event *models.ExternalEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 1 {
			log.Printf("Retry attempt %d/%d for event %s", attempt, c.config.MaxRetries, event.ResponseEventCommon.EventGUID)
			time.Sleep(time.Duration(c.config.RetryDelay) * time.Second)
		}

		req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			log.Printf("Attempt %d failed: %v", attempt, err)
			continue
		}

		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Check if request was successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("Event sent successfully: %s (status: %d)", event.ResponseEventCommon.EventGUID, resp.StatusCode)
			return nil
		}

		lastErr = fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		log.Printf("Attempt %d failed with status %d: %s", attempt, resp.StatusCode, string(body))

		// Don't retry on 4xx errors (client errors)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return lastErr
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", c.config.MaxRetries, lastErr)
}
