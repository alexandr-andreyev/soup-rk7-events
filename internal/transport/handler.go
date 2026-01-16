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
	NotifyService services.NotifyEventService
}

func NewHandler(notifySvc services.NotifyEventService) *Handler {
	return &Handler{
		NotifyService: notifySvc,
	}
}

func (h *Handler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	const op = "transport.HandleEvents"

	// Читаем тело запроса от кассовой системы
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//slog.Info("Request Body", "body", string(body))

	var rk7NotifyEvent models.Rk7NotifyEvent
	err = xml.Unmarshal(body, &rk7NotifyEvent)
	if err != nil {
		http.Error(w, "Failed to parse XML", http.StatusBadRequest)
		return
	}

	// Отправляем OK сразу после успешного парсинга
	w.WriteHeader(http.StatusOK)

	// Обработка события по правилам
	if err := h.NotifyService.HandleNotification(&rk7NotifyEvent); err != nil {
		slog.Error("Error handling notification", "error", err)
	}
}
