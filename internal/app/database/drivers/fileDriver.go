package drivers

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/thxhix/shortener/internal/app/database"
	"os"
)

type FileDatabase struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileDatabase(filePath string) (database.Database, error) {
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

func (db *FileDatabase) WriteRow(row *database.LinkRow) error {
	err := db.encoder.Encode(row)
	if err != nil {
		return err
	}
	return db.file.Sync()
}

func (db *FileDatabase) FindByHash(hash string) (*database.LinkRow, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var row database.LinkRow
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

func (db *FileDatabase) getLastUUID() (int, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(db.file)
	var lastUUID int

	for scanner.Scan() {
		var row database.LinkRow
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err == nil {
			lastUUID = row.UUID
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lastUUID, nil
}

func (db *FileDatabase) AddLink(original string, shorten string) (string, error) {
	lastID, err := db.getLastUUID()
	if err != nil {
		return "", err
	}

	newID := lastID + 1

	err = db.WriteRow(&database.LinkRow{
		UUID: newID,
		Hash: shorten,
		URL:  original,
	})
	if err != nil {
		return "", err
	}
	return shorten, nil
}

func (db *FileDatabase) GetFullLink(hash string) (string, error) {
	byHash, err := db.FindByHash(hash)
	if err != nil {
		return "", err
	}
	return byHash.URL, nil
}

func (db *FileDatabase) PingConnection() error {
	return nil
}
