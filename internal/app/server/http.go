package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/shortener/internal/app/database/interfaces"
	"net/http"

	"github.com/thxhix/shortener/internal/app/config"
)

type ServerInterface interface {
	StartPooling() error
}

type Server struct {
	config   *config.Config
	router   *chi.Mux
	database interfaces.Database
}

func NewServer(config config.Config, router chi.Mux, db interfaces.Database) ServerInterface {
	return &Server{
		config:   &config,
		router:   &router,
		database: db,
	}
}

func (s *Server) StartPooling() error {
	fmt.Println("* * * Запускаюсь * * *")
	fmt.Println("Адрес: " + s.config.Address.String())
	fmt.Println("Base URL: " + s.config.BaseURL.String())
	fmt.Println("* * * * * * * * * * *")

	if err := http.ListenAndServe(s.config.Address.String(), s.router); err != nil {
		return err
	}

	return nil
}
