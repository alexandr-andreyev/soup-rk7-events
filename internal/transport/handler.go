package transport

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
	"github.com/alexandr-andreyev/soup-rk7-events/internal/services"
)

type Handler struct {
	Logger        *slog.Logger
	NotifyService services.NotifyEventService
}

func NewHandler(logger *slog.Logger, notifySvc services.NotifyEventService) *Handler {
	return &Handler{
		Logger:        logger,
		NotifyService: notifySvc,
	}
}

func (h *Handler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	const op = "transport.HandleEvents"

	logger := h.Logger.With("op", op)

	// Читаем тело запроса от кассовой системы
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(op, "Failed to read request body", "error", err)
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logger.Info("events body", "body", string(body))

	var rk7NotifyEvent models.Rk7NotifyEvent
	err = xml.Unmarshal(body, &rk7NotifyEvent)
	if err != nil {
		logger.Error("Failed to parse XML", "error", err)
		http.Error(w, "Failed to parse XML", http.StatusBadRequest)
		return
	}

	// Отправляем OK сразу после успешного парсинга
	w.WriteHeader(http.StatusOK)

	// Обработка события по правилам
	if err := h.NotifyService.HandleNotification(&rk7NotifyEvent); err != nil {
		logger.Error("Error handling notification", "error", err)
	}
}
