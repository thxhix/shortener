package interfaces

import (
	"context"
	"database/sql"
	"github.com/thxhix/shortener/internal/models"
)

type Database interface {
	AddLink(original string, shorten string, userID string) (string, error)
	AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error
	GetFullLink(hash string) (string, error)
	Close() error
	PingConnection() error
	GetDriver() *sql.DB
}
