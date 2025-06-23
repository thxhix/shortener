package interfaces

import (
	"context"
	"database/sql"
	"github.com/thxhix/shortener/internal/models"
)

type Database interface {
	AddLink(ctx context.Context, original string, shorten string, userID string) (string, error)
	AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error
	GetFullLink(ctx context.Context, hash string) (models.DBShortenRow, error)
	GetUserFullLinks(ctx context.Context, userID string) (models.DBShortenRowList, error)
	RemoveUserLinks(ctx context.Context, userID string, ids []string) error
	Close() error
	PingConnection() error
	GetDriver() *sql.DB
}
