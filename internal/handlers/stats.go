package handlers

import (
	"github.com/mailru/easyjson"
	"net/http"
)

func (h *Handler) GetServiceStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.URLUsecase.GetServiceStats(r.Context())
	if err != nil {
		return
	}

	result, err := easyjson.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "не удалось записать ответ", http.StatusInternalServerError)
		return
	}
}
