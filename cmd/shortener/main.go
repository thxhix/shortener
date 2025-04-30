package main

import (
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database"
	"github.com/thxhix/shortener/internal/app/database/migrations"
	r "github.com/thxhix/shortener/internal/app/router"
	http "github.com/thxhix/shortener/internal/app/server"
	"go.uber.org/zap"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Кривое исполнение, но пока не представляю как работают миграции в Go
	err = migrations.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	logger := zap.NewExample()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal("Error syncing logger", zap.Error(err))
		}
	}()

	router := r.NewRouter(cfg, db, *logger)

	server := http.NewServer(*cfg, *router, db)
	err = server.StartPooling()
	if err != nil {
		log.Fatal(err)
	}
}
