package main

import (
	"fmt"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database"
	"github.com/thxhix/shortener/internal/app/database/migrations"
	r "github.com/thxhix/shortener/internal/app/router"
	http "github.com/thxhix/shortener/internal/app/server"
	"os"
)

func main() {
	cfg := config.NewConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Кривое исполнение, но пока не представляю как работают миграции в Go
	err = migrations.Migrate(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := r.NewRouter(cfg, db)

	server := http.NewServer(*cfg, *router, db)
	err = server.StartPooling()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
