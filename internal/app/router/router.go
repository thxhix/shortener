package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database/interfaces"
	handle "github.com/thxhix/shortener/internal/app/handlers"
	"github.com/thxhix/shortener/internal/app/middleware"
	"github.com/thxhix/shortener/internal/app/usecase"
	"go.uber.org/zap"
)

func NewRouter(cfg *config.Config, db interfaces.Database) *chi.Mux {
	uc := usecase.NewURLUseCase(db)

	logger := zap.NewExample()
	defer logger.Sync()

	router := chi.NewRouter()
	handlers := handle.NewHandler(cfg, uc)

	router.Route("/", func(r chi.Router) {
		// Кидаем на группу мидлвару с логами
		r.Use(middleware.WithLogging(logger.Sugar()))
		r.Use(middleware.CompressorMiddleware)

		r.Post("/", handlers.StoreLink)
		r.Get("/{id}", handlers.Redirect)
		r.Get("/ping", handlers.PingDatabase)

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", handlers.APIStoreLink)
				r.Post("/batch", handlers.BatchStoreLink)
			})
		})
	})

	return router
}
