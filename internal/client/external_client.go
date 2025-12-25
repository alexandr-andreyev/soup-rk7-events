package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

// White Server отправляет события при помощи Post-запросов:

// Body — в виде json запроса
// Headers — здесь добавлено поле Signature для проверки хэш кода содержимого body, удостоверяющее отправителя.

// Пример Headers:

// {
// "content-length":"1140",
// "signature":"+sMQMDTMoeE82ONz5ZR03ocQ+TSihyCGCW2Immrlwps=",
// "content-type":"application/json; charset=utf-8"
// }

// SendEvent sends an event to the external API with retry logic
func (c *ExternalClient) SendEvent(event *models.ExternalEvent) error {

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 1 {
			slog.Info("Retry attempt", "attempt", attempt, "maxRetries", c.config.MaxRetries, "guid", event.ResponseEventCommon.EventGUID)
			time.Sleep(time.Duration(c.config.RetryDelay) * time.Second)
		}

		req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// headers
		req.Header.Set("Content-Type", "application/json")

		// signature
		signature := c.GenerateHash(jsonData)
		req.Header.Set("Signature", signature)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			slog.Warn("Attempt failed", "attempt", attempt, "error", err)
			continue
		}

		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		slog.Info("Sending event",
			"guid", event.ResponseEventCommon.EventGUID,
			"event", string(jsonData))
		// Check if request was successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			slog.Info("Event sent successfully", "guid", event.ResponseEventCommon.EventGUID, "status", resp.StatusCode)

			return nil
		}

		lastErr = fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		slog.Warn("Attempt failed with non-success status", "attempt", attempt, "status", resp.StatusCode, "body", string(body))

		// Don't retry on 4xx errors (client errors)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return lastErr
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", c.config.MaxRetries, lastErr)
}
