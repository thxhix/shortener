package drivers

import (
	"context"
	"database/sql"
	"errors"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/models"
	"sync"
)

type MemoryDatabase struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func (db *MemoryDatabase) AddLink(original string, shorten string) (string, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.storage[shorten]; exists {
		return "", customErrors.ErrDuplicate
	}
	db.storage[shorten] = original
	return shorten, nil
}

func (db *MemoryDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, link := range list {
		if _, exists := db.storage[link.Hash]; exists {
			return customErrors.ErrDuplicate
		}
		db.storage[link.Hash] = link.URL
	}

	return nil
}

func (db *MemoryDatabase) GetFullLink(hash string) (string, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	value, ok := db.storage[hash]
	if ok {
		return value, nil
	}
	return "", errors.New("нет такой записи в БД")
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
		mutex:   sync.RWMutex{}, // для явности
	}, nil
}
