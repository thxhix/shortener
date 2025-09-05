package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

type ServerInterface interface {
	StartPooling() error
}

type Server struct {
	config   *config.Config
	router   *chi.Mux
	database interfaces.Database
	logger   *zap.SugaredLogger
}

func NewServer(config config.Config, router chi.Mux, db interfaces.Database, logger *zap.SugaredLogger) ServerInterface {
	return &Server{
		config:   &config,
		router:   &router,
		database: db,
		logger:   logger,
	}
}

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

	if err := http.ListenAndServe(s.config.Address, s.router); err != nil {
		return err
	}

	return nil
}
