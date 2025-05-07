package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	"go.uber.org/zap"
	"net/http"
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
	s.logger.Info("* * * * * * * * * * *")

	if err := http.ListenAndServe(s.config.Address, s.router); err != nil {
		return err
	}

	return nil
}
