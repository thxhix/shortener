package drivers

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/thxhix/shortener/internal/app/database"
	"time"
)

type PostgresQLDatabase struct {
	driver *sql.DB
}

func (db *PostgresQLDatabase) AddLink(original string, shorten string) (string, error) {
	//TODO implement me
	panic("i can't do anything now..")
}

func (db *PostgresQLDatabase) GetFullLink(hash string) (string, error) {
	//TODO implement me
	panic("i can't do anything now..")
}

func (db *PostgresQLDatabase) Close() error {
	return db.driver.Close()
}

func (db *PostgresQLDatabase) PingConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return db.driver.PingContext(ctx)
}

func NewPQLDatabase(params string) (database.Database, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	return &PostgresQLDatabase{
		driver: db,
	}, nil
}
