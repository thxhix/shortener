package interfaces

import (
	"context"
	"database/sql"
	"github.com/thxhix/shortener/internal/models"
)

// Database defines the contract for all storage backends.
type Database interface {
	// RunMigrations runs initial schema migrations for the database.
	RunMigrations() error

	// AddLink stores a single shortened link and returns its hash.
	AddLink(ctx context.Context, original string, shorten string, userID string) (string, error)

	// AddLinks stores a batch of shortened links.
	AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error

	// GetFullLink retrieves the original link by its short hash.
	GetFullLink(ctx context.Context, hash string) (models.DBShortenRow, error)

	// GetUserFullLinks retrieves all links created by the given user.
	GetUserFullLinks(ctx context.Context, userID string) (models.DBShortenRowList, error)

	// RemoveUserLinks deletes links by their IDs for the given user.
	RemoveUserLinks(ctx context.Context, userID string, ids []string) error

	// Close releases resources and closes the database connection.
	Close() error

	// PingConnection checks if the database connection is alive.
	PingConnection() error

	// GetDriver returns the underlying sql.DB object for advanced use.
	GetDriver() *sql.DB
}
