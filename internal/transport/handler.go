package transport

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"github.com/alexandr-andreyev/soup-rk7-events/internal/models"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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

	var rk7NotifyEvent models.Rk7NotifyEvent
	err = xml.Unmarshal(body, &rk7NotifyEvent)
	if err != nil {
		http.Error(w, "Failed to parse XML", http.StatusBadRequest)
		return
	}

	log.Println("Parsed Event:", rk7NotifyEvent)
	log.Println("Event Name:", rk7NotifyEvent.Name)
	log.Println("Event RestCode:", rk7NotifyEvent.RestCode)

	// Можно добавить лог в файл позже
	w.WriteHeader(http.StatusOK)
}
