package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var Database = map[string]string{}

func shortLink(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	isCorrectLink := isURL(string(body))

	if isCorrectLink {
		link := writeLink(string(body))

		// Отправляем ответ
		baseURL := destinationURL.String()

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		io.WriteString(w, baseURL+"/"+link)
	} else {
		badRequest(w)
	}
}

func getFullLink(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if hasLink(id) {
		w.Header().Add("Location", Database[id])
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		badRequest(w)
	}
}

// функция main вызывается автоматически при запуске приложения
func main() {
	hostParams.Set("localhost:8080")
	destinationURL.Set("http://localhost:8080")

	parseFlags()

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortLink)
		r.Get("/{id}", getFullLink)
	})

	fmt.Println("Запускаю сервер по адресу:", hostParams.String())
	fmt.Println("Вывод по адресу:", destinationURL.String())

	if err := http.ListenAndServe(hostParams.String(), r); err != nil {
		panic(err)
	}
}

func writeLink(link string) string {
	index := "link" + strconv.Itoa(len(Database)+1)
	Database[index] = link
	return index
}

func hasLink(slug string) bool {
	_, ok := Database[slug]
	return ok
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "Bad request")
}

// Проверка валидности URL
func isURL(link string) bool {
	parsedURL, err := url.Parse(link)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}
