package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body:", err)
	}
	defer r.Body.Close()

	// Логируем основную информацию
	log.Println("========== NEW REQUEST ==========")
	log.Println("Method:", r.Method)
	log.Println("URL:", r.URL.String())
	log.Println("RemoteAddr:", r.RemoteAddr)

	// Логируем заголовки
	log.Println("Headers:")
	for k := range r.Header {
		log.Printf("  %s: %s\n", k, r.Header.Get(k))
	}

	// Логируем тело
	log.Println("Body:")
	fmt.Println(string(body))
	log.Println("=================================")

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received\n"))
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Server started on :28080")
	if err := http.ListenAndServe(":28080", nil); err != nil {
		log.Fatal(err)
	}
}
