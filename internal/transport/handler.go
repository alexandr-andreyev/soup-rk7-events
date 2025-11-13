package transport

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	// Выводим содержимое тела запроса в консоль (или лог)
	fmt.Println("Event body:", string(body))

	// Можно добавить лог в файл позже
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event received"))
}
