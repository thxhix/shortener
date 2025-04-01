package database

import (
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database/drivers"
	"github.com/thxhix/shortener/internal/app/database/interfaces"
)

func NewDatabase(config *config.Config) (interfaces.Database, error) {
	if config.PostgresQL != "" {
		return drivers.NewPQLDatabase(config.PostgresQL)
	}
	if config.DBFileName != "" {
		return drivers.NewFileDatabase(config.DBFileName)
	}
	return drivers.NewMemoryDatabase()
}
