package handlers

import (
	"github.com/thxhix/shortener/internal/app/usecase"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/config"
)

type Handler struct {
	config     config.Config
	URLUsecase usecase.URLUseCaseInterface
}

func NewHandler(cfg *config.Config, useCase usecase.URLUseCaseInterface) *Handler {
	return &Handler{
		config:     *cfg,
		URLUsecase: useCase,
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

	link, err := h.URLUsecase.Shorten(parsedURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")

	_, err = io.WriteString(w, h.config.BaseURL.String()+"/"+link)
	if err != nil {
		http.Error(w, "не удалось записать ответ", http.StatusBadRequest)
		return
	}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	link, err := h.URLUsecase.GetFullURL(id)
	if err != nil {
		http.Error(w, "такой страницы нет", http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
