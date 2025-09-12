package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	handle "github.com/thxhix/shortener/internal/handlers"
	"github.com/thxhix/shortener/internal/middleware"
	"github.com/thxhix/shortener/internal/url"
	"go.uber.org/zap"
)

// NewRouter creates and configures a new chi.Mux router instance.
// It sets up the URL shortening use case, attaches middleware,
// and registers all HTTP endpoints for the service:
//
//   - POST   /               → Store a short link
//
//   - GET    /{id}           → Redirect to original URL
//
//   - GET    /ping           → Ping database
//
//   - GET    /api/user/urls  → List user links
//
//   - DELETE /api/user/urls  → Delete user links
//
//   - POST   /api/shorten          → Store a short link via API
//
//   - POST   /api/shorten/batch    → Store multiple links via API
//
// The following middleware are applied to the root route group:
//   - WithLogging: request logging using zap logger
//   - CompressorMiddleware: response compression
//   - Auth: authentication based on SecretKey
func NewRouter(cfg *config.Config, db interfaces.Database, logger *zap.SugaredLogger) *chi.Mux {
	uc := url.NewURLUseCase(db, *cfg)

	router := chi.NewRouter()
	handlers := handle.NewHandler(cfg, uc)

	router.Route("/", func(r chi.Router) {
		// Кидаем на группу мидлвару с логами
		r.Use(middleware.WithLogging(logger))
		r.Use(middleware.CompressorMiddleware)
		r.Use(middleware.Auth(cfg.SecretKey))

		r.Post("/", handlers.StoreLink)
		r.Get("/{id}", handlers.Redirect)
		r.Get("/ping", handlers.PingDatabase)

		r.Route("/api", func(r chi.Router) {

			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", handlers.UserList)
				r.Delete("/urls", handlers.UserDeleteRows)
			})

			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", handlers.APIStoreLink)
				r.Post("/batch", handlers.BatchStoreLink)
			})
		})
	})

	return router
}
