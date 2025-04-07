package interfaces

import (
	"context"
	"database/sql"
	"github.com/thxhix/shortener/internal/app/models"
)

type Database interface {
	AddLink(original string, shorten string) (string, error)
	AddLinks(ctx context.Context, list models.BatchList) error
	GetFullLink(hash string) (string, error)
	Close() error
	PingConnection() error
	GetDriver() *sql.DB
}
