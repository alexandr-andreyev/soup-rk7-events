package transport

import (
	"encoding/xml"
	"io"
	"log"
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
	log.Println("Received event")

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//log.Println("Request Body:", string(body))

	var rk7NotifyEvent models.Rk7NotifyEvent
	err = xml.Unmarshal(body, &rk7NotifyEvent)
	if err != nil {
		http.Error(w, "Failed to parse XML", http.StatusBadRequest)
		return
	}

	// Отправляем OK сразу после успешного парсинга
	w.WriteHeader(http.StatusOK)

	// Обрабатываем уведомление синхронно для сохранения порядка событий
	if err := h.NotifyService.HandleNotification(&rk7NotifyEvent); err != nil {
		log.Printf("Error handling notification: %v", err)
	}
}
