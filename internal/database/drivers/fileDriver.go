package drivers

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/thxhix/shortener/internal/database/interfaces"
	"github.com/thxhix/shortener/internal/models"
	"log"
	"os"
)

var (
	// ErrUserNotFound is returned when no records exist for the given user ID.
	ErrUserNotFound = errors.New("пользователь не найден")
)

// FileDatabase implements the Database interface, and using a JSON-lines file.
// Each record is stored as a single JSON object per line.
type FileDatabase struct {
	file    *os.File
	encoder *json.Encoder
}

// NewFileDatabase creates a new FileDatabase instance for the given file path.
// If the file does not exist, it will be created.
// Returns an error if the file cannot be opened.g
func NewFileDatabase(filePath string) (interfaces.Database, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileDatabase{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// RunMigrations is a no-op for FileDatabase.
// Always returns nil.
func (db *FileDatabase) RunMigrations() error {
	return nil
}

// Close closes the underlying file used by FileDatabase.
func (db *FileDatabase) Close() error {
	return db.file.Close()
}

// WriteRow appends a new DBShortenRow to the file as a JSON object.
// The file is synced after each write.
func (db *FileDatabase) WriteRow(row *models.DBShortenRow) error {
	err := db.encoder.Encode(row)
	if err != nil {
		return err
	}
	return db.file.Sync()
}

// FindByHash scans the file and returns the first record matching the hash.
// Returns an error if not found or if JSON decoding fails.
func (db *FileDatabase) FindByHash(hash string) (*models.DBShortenRow, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var row models.DBShortenRow
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			continue
		}

		if row.Hash == hash {
			return &row, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, errors.New("запись не найдена")
}

// FindByUserID scans the file and returns all records for the given user ID.
// Returns ErrUserNotFound if no records exist.
func (db *FileDatabase) FindByUserID(userID string) (models.DBShortenRowList, error) {
	result := models.DBShortenRowList{}

	_, err := db.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var row models.DBShortenRow
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			log.Printf("ошибка чтения строки из файла: %v", err)
			continue
		}

		if row.UserID == userID {
			result = append(result, row)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrUserNotFound
	}

	return result, nil
}

// AddLink stores a single shortened link in the file.
// Returns the hash of the link or an error if writing fails.
func (db *FileDatabase) AddLink(ctx context.Context, original string, shorten string, userID string) (string, error) {
	err := db.WriteRow(&models.DBShortenRow{
		Hash:   shorten,
		URL:    original,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	return shorten, nil
}

// AddLinks stores multiple shortened links in the file.
// Returns an error if any write fails.
func (db *FileDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error {
	for _, link := range list {

		link.UserID = userID
		err := db.WriteRow(&link)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFullLink retrieves the original URL by its short hash.
// Returns an error if the hash does not exist.
func (db *FileDatabase) GetFullLink(ctx context.Context, hash string) (models.DBShortenRow, error) {
	byHash, err := db.FindByHash(hash)
	if err != nil {
		return models.DBShortenRow{}, err
	}
	return models.DBShortenRow{URL: byHash.URL}, nil
}

// GetUserFullLinks retrieves all links belonging to the specified user ID.
// Returns ErrUserNotFound if no records exist.
func (db *FileDatabase) GetUserFullLinks(ctx context.Context, userID string) (models.DBShortenRowList, error) {
	return db.FindByUserID(userID)
}

// RemoveUserLinks is not implemented for FileDatabase.
// Always returns nil.
func (db *FileDatabase) RemoveUserLinks(ctx context.Context, userID string, ids []string) error {
	return nil
}

// PingConnection always succeeds for FileDatabase.
// It returns nil to indicate the database is available.
func (db *FileDatabase) PingConnection() error {
	return nil
}

// GetDriver always returns nil for FileDatabase since it has no SQL driver.
func (db *FileDatabase) GetDriver() *sql.DB { return nil }
