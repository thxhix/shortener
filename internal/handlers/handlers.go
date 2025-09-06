package handlers

import (
	json2 "encoding/json"
	"errors"
	"github.com/thxhix/shortener/internal/config"
	custorErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/middleware"
	"github.com/thxhix/shortener/internal/models"
	urlUseCase "github.com/thxhix/shortener/internal/url"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
)

// Handler groups all HTTP handlers for the URL shortener service.
// It provides endpoints for creating, retrieving, listing and deleting links.
type Handler struct {
	config     config.Config
	URLUsecase urlUseCase.URLUseCaseInterface
}

// NewHandler creates a new Handler instance with the given configuration
// and URL use case implementation.
func NewHandler(cfg *config.Config, useCase urlUseCase.URLUseCaseInterface) *Handler {
	return &Handler{
		config:     *cfg,
		URLUsecase: useCase,
	}
}

// StoreLink It reads the raw URL from the request body, validates it,
// and returns a shortened link as plain text
// Responds with 201 Created on success or 409 Conflict if the URL already exists.
func (h *Handler) StoreLink(w http.ResponseWriter, r *http.Request) {
	targetURL, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("Error closing body: %v", err)
		}
	}()
	if err != nil || string(targetURL) == "" {
		http.Error(w, "не удалось прочитать ссылку из тела запроса", http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(string(targetURL))
	if err != nil {
		http.Error(w, "не удалось спарсить переданную ссылку", http.StatusBadRequest)
		return
	}

	var isConflict = false
	link, err := h.URLUsecase.Shorten(r.Context(), parsedURL.String())
	if err != nil {
		if errors.Is(err, custorErrors.ErrDuplicate) {
			isConflict = true
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	if !isConflict {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	_, err = io.WriteString(w, h.config.BaseURL+"/"+link)
	if err != nil {
		http.Error(w, "не удалось записать ответ", http.StatusInternalServerError)
		return
	}
}

// Redirect It looks up the full URL by the short hash and issues a 307 redirect.
// If the link was deleted, responds with 410 Gone.
// If the link does not exist, responds with 400 Bad Request.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	link, err := h.URLUsecase.GetFullURL(r.Context(), id)
	if err != nil {
		if errors.Is(err, urlUseCase.ErrLinkDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, "такой страницы нет", http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// APIStoreLink It reads a JSON payload with the original URL, validates it,
// and returns a JSON response containing the shortened URL.
// Responds with 201 Created on success or 409 Conflict if the URL already exists.
func (h *Handler) APIStoreLink(w http.ResponseWriter, r *http.Request) {
	json, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("Error closing body: %v", err)
		}
	}()
	if err != nil {
		http.Error(w, "не удалось прочитать ссылку из тела запроса", http.StatusBadRequest)
		return
	}

	fullURL := &models.FullURL{}
	err = easyjson.Unmarshal(json, fullURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var isConflict = false

	link, err := h.URLUsecase.Shorten(r.Context(), fullURL.URL)
	if err != nil {
		if errors.Is(err, custorErrors.ErrDuplicate) {
			isConflict = true
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	link = h.config.BaseURL + "/" + link

	shortURL := models.ShortURL{Result: link}
	result, err := easyjson.Marshal(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !isConflict {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "не удалось записать ответ", http.StatusInternalServerError)
		return
	}
}

// BatchStoreLink It reads a JSON array of objects with correlation_id and original_url,
// validates the input, and returns a JSON array of shortened links.
func (h *Handler) BatchStoreLink(w http.ResponseWriter, r *http.Request) {
	json, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "не удалось прочитать ссылки из тела запроса", http.StatusBadRequest)
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("Error closing body: %v", err)
		}
	}()

	var batch models.BatchShortenRequestList

	if err := easyjson.Unmarshal(json, &batch); err != nil {
		http.Error(w, "невалидный JSON", http.StatusBadRequest)
		return
	}

	if len(batch) == 0 {
		http.Error(w, "пустой batch в запросе", http.StatusBadRequest)
		return
	}

	data, err := h.URLUsecase.BatchShorten(r.Context(), batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := data.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// PingDatabase It checks the database connectivity and responds with 200 OK if successful,
// or 500 Internal Server Error otherwise.
func (h *Handler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	err := h.URLUsecase.PingDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UserList It returns all URLs belonging to the authenticated user.
// If the user has no links, responds with 204 No Content.
func (h *Handler) UserList(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if userID == "" {
		http.Error(w, "неверный данные авторизации..", http.StatusUnauthorized)
		return
	}

	links, err := h.URLUsecase.UserList(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(links) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := links.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(result)
	if err != nil {
		log.Printf("ошибка при записи ответа: %v", err)
		return
	}
}

func (h *Handler) UserDeleteRows(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if userID == "" {
		http.Error(w, "неверный данные авторизации..", http.StatusUnauthorized)
		return
	}

	var ids []string
	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	err = json2.Unmarshal(jsonBody, &ids)
	if err != nil {
		http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	go h.URLUsecase.UserDeleteRows(middleware.GetUserID(r.Context()), ids, h.config.DeleteWorkersCount, h.config.DeleteBatchSize)

	w.WriteHeader(http.StatusAccepted)
}
