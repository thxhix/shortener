package main

import (
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database"
	"github.com/thxhix/shortener/internal/meta"
	r "github.com/thxhix/shortener/internal/router"
	http "github.com/thxhix/shortener/internal/server"
	"go.uber.org/zap"
	"log"
)

// buildVersion, buildDate, and buildCommit are global variables that can be
// set at build time using go build -ldflags to include build metadata such as
// version, build date, and commit hash. If left empty, "N/A" will be displayed.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	meta.PrintMeta(buildVersion, buildDate, buildCommit)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Кривое исполнение, но пока не представляю как работают миграции в Go
	err = db.RunMigrations()
	if err != nil {
		log.Fatal(err)
	}

	zapLogger := zap.NewExample()
	defer func() {
		err := zapLogger.Sync()
		if err != nil {
			log.Fatal("Error syncing logger", zap.Error(err))
		}
	}()

	router := r.NewRouter(cfg, db, zapLogger.Sugar())

	server := http.NewServer(*cfg, *router, db, zapLogger.Sugar())
	err = server.StartPooling()
	if err != nil {
		log.Fatal(err)
	}
}
