package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database"
	handle "github.com/thxhix/shortener/internal/app/handlers"
	"github.com/thxhix/shortener/internal/app/usecase"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	db := database.NewDatabase()
	uc := usecase.NewURLUseCase(db)

	router := chi.NewRouter()
	handlers := handle.NewHandler(cfg, uc)

	router.Route("/", func(r chi.Router) {
		router.Post("/", handlers.StoreLink)
		router.Get("/{id}", handlers.Redirect)
	})

	return router
}
