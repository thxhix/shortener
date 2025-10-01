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

	// GetUsersCount returns the total number of unique users in the database.
	// A user is identified by their user_id value. If user_id is NULL, it is ignored.
	// The method executes a COUNT(DISTINCT user_id) query and returns the result.
	GetUsersCount(ctx context.Context) (int64, error)

	// GetURLsCount returns the total number of shortened URLs stored in the database.
	// Deleted URLs (where is_deleted = true) are excluded from the count.
	// The method executes a COUNT(*) query with the proper filter and returns the result.
	GetURLsCount(ctx context.Context) (int64, error)
}
