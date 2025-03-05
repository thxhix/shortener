package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/router"
)

type Server struct {
	config config.Config
}

func initServer() *Server {
	return &Server{
		config: *config.InitConfig(),
	}
}

func StartPooling() {
	params := initServer()
	router := router.InitRouter(&params.config)

	fmt.Println("* * * Запускаюсь * * *")
	fmt.Println("Адрес: " + params.config.Address.String())
	fmt.Println("Base URL: " + params.config.BaseURL.String())
	fmt.Println("* * * * * * * * * * *")

	if err := http.ListenAndServe(params.config.Address.String(), router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
