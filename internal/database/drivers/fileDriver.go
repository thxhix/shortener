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
	ErrUserNotFound = errors.New("пользователь не найден")
)

type FileDatabase struct {
	file    *os.File
	encoder *json.Encoder
}

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

func (db *FileDatabase) Close() error {
	return db.file.Close()
}

func (db *FileDatabase) WriteRow(row *models.DBShortenRow) error {
	err := db.encoder.Encode(row)
	if err != nil {
		return err
	}
	return db.file.Sync()
}

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

func (db *FileDatabase) getLastUUID() (int, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(db.file)
	var lastUUID int

	for scanner.Scan() {
		var row models.DBShortenRow
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err == nil {
			lastUUID = row.ID
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lastUUID, nil
}

func (db *FileDatabase) AddLink(original string, shorten string, userID string) (string, error) {
	lastID, err := db.getLastUUID()
	if err != nil {
		return "", err
	}

	newID := lastID + 1

	err = db.WriteRow(&models.DBShortenRow{
		ID:     newID,
		Hash:   shorten,
		URL:    original,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	return shorten, nil
}

func (db *FileDatabase) AddLinks(ctx context.Context, list models.DBShortenRowList, userID string) error {
	for _, link := range list {
		lastID, err := db.getLastUUID()
		if err != nil {
			return err
		}

		link.ID = lastID + 1
		link.UserID = userID

		err = db.WriteRow(&link)

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *FileDatabase) GetFullLink(hash string) (models.DBShortenRow, error) {
	byHash, err := db.FindByHash(hash)
	if err != nil {
		return models.DBShortenRow{}, err
	}
	return models.DBShortenRow{URL: byHash.URL}, nil
}

func (db *FileDatabase) GetUserFullLinks(userID string) (models.DBShortenRowList, error) {
	return db.FindByUserID(userID)
}

func (db *FileDatabase) RemoveUserLinks(userID string, ids []string) error {
	return nil
}

func (db *FileDatabase) PingConnection() error {
	return nil
}

func (db *FileDatabase) GetDriver() *sql.DB { return nil }
