package interfaces

import "database/sql"

type Database interface {
	AddLink(original string, shorten string) (string, error)
	GetFullLink(hash string) (string, error)
	Close() error
	PingConnection() error
	GetDriver() *sql.DB
}
