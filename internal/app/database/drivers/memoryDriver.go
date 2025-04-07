package drivers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/thxhix/shortener/internal/app/models"
)

type MemoryDatabase struct {
	storage map[string]string
}

func (db *MemoryDatabase) AddLink(original string, shorten string) (string, error) {
	if _, err := db.GetFullLink(shorten); err == nil {
		return "", errors.New("такая запись в БД уже есть")
	}
	db.storage[shorten] = original
	return shorten, nil
}

func (db *MemoryDatabase) AddLinks(ctx context.Context, list models.DatabaseRowList) error {
	for _, link := range list {
		if _, err := db.GetFullLink(link.Hash); err == nil {
			return err
		}
		db.storage[link.Hash] = link.URL
	}

	return nil
}

func (db *MemoryDatabase) GetFullLink(hash string) (string, error) {
	value, ok := db.storage[hash]
	if ok {
		return value, nil
	}
	return value, errors.New("нет такой записи в БД")
}

func (db *MemoryDatabase) Close() error {
	return nil
}

func (db *MemoryDatabase) PingConnection() error {
	return nil
}

func (db *MemoryDatabase) GetDriver() *sql.DB {
	return nil
}

func NewMemoryDatabase() (*MemoryDatabase, error) {
	return &MemoryDatabase{
		storage: make(map[string]string),
	}, nil
}
