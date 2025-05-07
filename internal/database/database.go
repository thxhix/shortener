package database

import (
	"github.com/thxhix/shortener/internal/config"
	drivers2 "github.com/thxhix/shortener/internal/database/drivers"
	"github.com/thxhix/shortener/internal/database/interfaces"
)

func NewDatabase(config *config.Config) (interfaces.Database, error) {
	if config.PostgresQL != "" {
		return drivers2.NewPQLDatabase(config.PostgresQL)
	}
	if config.DBFileName != "" {
		return drivers2.NewFileDatabase(config.DBFileName)
	}
	return drivers2.NewMemoryDatabase()
}
