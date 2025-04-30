package handlers

import (
	"errors"
	custorErrors "github.com/thxhix/shortener/internal/app/errors"
	"github.com/thxhix/shortener/internal/app/models"
	"github.com/thxhix/shortener/internal/app/usecase"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
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
	link, err := h.URLUsecase.Shorten(parsedURL.String())
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

// * Было замечание насчет названия метода, но в итоге не придумал ничего лучше. Кажется итак сойдет..
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

	link, err := h.URLUsecase.Shorten(fullURL.URL)
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

func (h *Handler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	err := h.URLUsecase.PingDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
