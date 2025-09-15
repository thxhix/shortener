package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

// ServerInterface the contract for running the HTTP server.
type ServerInterface interface {
	// StartPooling starts the HTTP server and begins listening for requests.
	// If profiling is enabled in the config, a separate pprof server is also started.
	StartPooling() error
}

// Server is the concrete implementation of the HTTP server.
// It holds the configuration, router, database and logger instances.
type Server struct {
	config   *config.Config
	router   *chi.Mux
	database interfaces.Database
	logger   *zap.SugaredLogger
}

// NewServer creates a new Server instance with the provided configuration,
// router, database, and logger. It returns a ServerInterface implementation.
func NewServer(config config.Config, router chi.Mux, db interfaces.Database, logger *zap.SugaredLogger) ServerInterface {
	return &Server{
		config:   &config,
		router:   &router,
		database: db,
		logger:   logger,
	}
}

// StartPooling starts the main HTTP API server using the configured address.
// If profiling is enabled in the configuration, a separate pprof server is
// started on the ProfilerAddress in a separate goroutine. The method blocks
// until the main HTTP server exits or encounters an error.
func (s *Server) StartPooling() error {
	s.logger.Info("* * * Запускаюсь * * *")
	s.logger.Infof("Адрес: %s", s.config.Address)
	s.logger.Infof("Base URL: %s", s.config.BaseURL)
	if s.config.EnableProfiler {
		s.logger.Infof("ProfilerAddress: %s", s.config.ProfilerAddress)
	}
	s.logger.Info("* * * * * * * * * * *")

	if s.config.EnableProfiler {
		// pprof отдельно, чтобы не мешать в API
		go func() {
			addr := s.config.ProfilerAddress
			if err := http.ListenAndServe(addr, nil); err != nil {
				s.logger.Errorf("pprof server error: %v", err)
			}
		}()
	}

	if !s.config.EnableHTTPS {
		return http.ListenAndServe(s.config.Address, s.router)
	}

	return startWithHTTPS(s)
}
