package drivers

import (
	"context"
	"database/sql"
	"errors"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/models"
	"sync"
)

// MemoryDatabase implements the Database interface using
// an in-memory map with synchronization via RWMutex.
// Data is not persisted and will be lost when the process exits.
type MemoryDatabase struct {
	storage map[string]string
	mutex   sync.RWMutex
}

// NewMemoryDatabase creates and returns a new MemoryDatabase instance.
// The storage is initialized as an empty map.
func NewMemoryDatabase() (*MemoryDatabase, error) {
	return &MemoryDatabase{
		storage: make(map[string]string),
		mutex:   sync.RWMutex{}, // для явности
	}, nil
}

// RunMigrations is a no-op for MemoryDatabase.
// It always returns nil.
func (db *MemoryDatabase) RunMigrations() error {
	return nil
}

// AddLink stores a single shortened link in memory.
// Returns ErrDuplicate if the hash already exists.
func (db *MemoryDatabase) AddLink(ctx context.Context, original string, shorten string, userID string) (string, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.storage[shorten]; exists {
		return "", customErrors.ErrDuplicate
	}
	db.storage[shorten] = original
	return shorten, nil
}

// AddLinks stores multiple shortened links in memory.
// Returns ErrDuplicate if any hash already exists.
func (db *MemoryDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error {
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

// GetFullLink retrieves the original URL by hash from memory.
// Returns an error if the hash does not exist.
func (db *MemoryDatabase) GetFullLink(ctx context.Context, hash string) (models.DBShortenRow, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	value, ok := db.storage[hash]
	if ok {
		return models.DBShortenRow{URL: value}, nil
	}
	return models.DBShortenRow{}, errors.New("нет такой записи в БД")
}

// Close is a no-op for MemoryDatabase.
// It always returns nil.
func (db *MemoryDatabase) Close() error {
	return nil
}

// PingConnection always succeeds for MemoryDatabase.
// It returns nil to indicate the database is available.
func (db *MemoryDatabase) PingConnection() error {
	return nil
}

// GetDriver always returns nil for MemoryDatabase since it has no SQL driver.
func (db *MemoryDatabase) GetDriver() *sql.DB {
	return nil
}

// GetUserFullLinks is not implemented for MemoryDatabase.
// Always returns nil.
func (db *MemoryDatabase) GetUserFullLinks(ctx context.Context, userID string) (models.DBShortenRowList, error) {
	return nil, nil
}

// RemoveUserLinks is not implemented for MemoryDatabase.
// Always returns nil.
func (db *MemoryDatabase) RemoveUserLinks(ctx context.Context, userID string, ids []string) error {
	return nil
}
