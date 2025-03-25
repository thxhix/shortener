package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type DatabaseInterface interface {
	AddLink(original string, shorten string) (string, error)
	GetFullLink(hash string) (string, error)
	Close() error
}

type Database struct {
	file    *os.File
	encoder *json.Encoder
}

type LinkRow struct {
	UUID int    `json:"uuid"`
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

func NewDatabase(filePath string) (DatabaseInterface, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Database{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (db *Database) Close() error {
	return db.file.Close()
}

func (db *Database) WriteRow(row *LinkRow) error {
	err := db.encoder.Encode(row)
	if err != nil {
		return err
	}
	return db.file.Sync()
}

func (db *Database) FindByHash(hash string) (*LinkRow, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		var row LinkRow
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

func (db *Database) getLastUUID() (int, error) {
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(db.file)
	var lastUUID int

	for scanner.Scan() {
		var row LinkRow
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

func (db *Database) AddLink(original string, shorten string) (string, error) {
	lastID, err := db.getLastUUID()
	if err != nil {
		return "", err
	}

	newID := lastID + 1

	err = db.WriteRow(&LinkRow{
		UUID: newID,
		Hash: shorten,
		URL:  original,
	})
	if err != nil {
		return "", err
	}
	return shorten, nil
}

func (db *Database) GetFullLink(hash string) (string, error) {
	byHash, err := db.FindByHash(hash)
	if err != nil {
		return "", err
	}
	return byHash.URL, nil
}
