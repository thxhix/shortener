package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database"
)

var DB = database.CreateDatabase()

type Handler struct {
	config config.Config
}

func InitHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: *cfg,
	}
}

func (h *Handler) StoreLink(w http.ResponseWriter, r *http.Request) {
	targetURL, err := io.ReadAll(r.Body)
	if err != nil || string(targetURL) == "" {
		http.Error(w, "не удалось прочитать ссылку из тела запроса", http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(string(targetURL))
	if err != nil {
		http.Error(w, "не удалось спарсить переданную ссылку", http.StatusBadRequest)
		return
	}

	index := DB.AddRow("link", parsedURL.String())

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")

	io.WriteString(w, h.config.BaseURL.String()+"/"+index)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	link, err := DB.GetRow(id)
	if err != nil {
		http.Error(w, "такой страницы нет", http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
