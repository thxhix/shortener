package database

import (
	"github.com/thxhix/shortener/internal/config"
	drivers2 "github.com/thxhix/shortener/internal/database/drivers"
	"github.com/thxhix/shortener/internal/database/interfaces"
)

// NewDatabase creates a new Database implementation based on the provided configuration.
//
// Priority of selection:
//  1. If PostgresQL DSN is provided, returns a PostgreSQL-backed database.
//  2. If a file path (DBFileName) is provided, returns a file-based database.
//  3. Otherwise, returns an in-memory database (not persistent).
//
// The returned value implements the Database interface. An error is returned
// if the chosen backend cannot be initialized (e.g., failed to connect to PostgreSQL).
func NewDatabase(config *config.Config) (interfaces.Database, error) {
	if config.PostgresQL != "" {
		return drivers2.NewPQLDatabase(config.PostgresQL)
	}
	if config.DBFileName != "" {
		return drivers2.NewFileDatabase(config.DBFileName)
	}
	return drivers2.NewMemoryDatabase()
}
