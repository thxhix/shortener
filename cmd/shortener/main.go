// пакеты исполняемых приложений должны называться main
package main

import (
	"fmt"
	"io"
	"net/http"
)

var fullLink string = "https://google.com"

func shortLink(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}
	req.Body.Close()

	fullLink = string(body)

	// Отправляем ответ
	baseURL := "http://" + req.Host

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	io.WriteString(w, baseURL+"/EwHXdJfB")
}

func getFullLink(w http.ResponseWriter, req *http.Request) {
	// id := req.PathValue("id")
	fmt.Println(fullLink)
	w.Header().Add("Location", fullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// функция main вызывается автоматически при запуске приложения
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", shortLink)
	mux.HandleFunc("GET /{id}", getFullLink)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
