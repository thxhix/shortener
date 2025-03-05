package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/handlers"
)

func InitRouter(cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()
	handlers := handlers.InitHandler(cfg)

	router.Route("/", func(r chi.Router) {
		router.Post("/", handlers.StoreLink)
		router.Get("/{id}", handlers.Redirect)
	})

	return router
}
