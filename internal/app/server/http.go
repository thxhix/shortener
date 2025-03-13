package server

import (
	"fmt"
	"net/http"

	"github.com/thxhix/shortener/internal/app/config"
	r "github.com/thxhix/shortener/internal/app/router"
)

type ServerInterface interface {
	StartPooling() error
}

type Server struct {
	config config.Config
}

func NewServer() ServerInterface {
	return &Server{
		config: *config.NewConfig(),
	}
}

func (s *Server) StartPooling() error {
	router := r.NewRouter(&s.config)

	fmt.Println("* * * Запускаюсь * * *")
	fmt.Println("Адрес: " + s.config.Address.String())
	fmt.Println("Base URL: " + s.config.BaseURL.String())
	fmt.Println("* * * * * * * * * * *")

	if err := http.ListenAndServe(s.config.Address.String(), router); err != nil {
		return err
	}

	return nil
}
