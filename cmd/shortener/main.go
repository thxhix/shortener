// пакеты исполняемых приложений должны называться main
package main

import (
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
		baseURL := "http://" + req.Host

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
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortLink)
		r.Get("/{id}", getFullLink)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
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

func isURL(link string) bool {
	if len([]rune(link)) > 0 {
		parsedURL, err := url.Parse(link)
		return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
	}
	return false
}
