package database

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"time"
)

type Database interface {
	AddLink(original string, shorten string) (string, error)
	GetFullLink(hash string) (string, error)
	Close() error
	PingConnection() error
}

type FileDatabase struct {
	file    *os.File
	encoder *json.Encoder
}

type LinkRow struct {
	UUID int    `json:"uuid"`
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

func NewFileDatabase(filePath string) (Database, error) {
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

func (db *FileDatabase) WriteRow(row *LinkRow) error {
	err := db.encoder.Encode(row)
	if err != nil {
		return err
	}
	return db.file.Sync()
}

func (db *FileDatabase) FindByHash(hash string) (*LinkRow, error) {
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

func (db *FileDatabase) getLastUUID() (int, error) {
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

func (db *FileDatabase) AddLink(original string, shorten string) (string, error) {
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

func (db *FileDatabase) GetFullLink(hash string) (string, error) {
	byHash, err := db.FindByHash(hash)
	if err != nil {
		return "", err
	}
	return byHash.URL, nil
}

func (p *FileDatabase) PingConnection() error {
	return nil
}

type PostgresQLDatabase struct {
	driver *sql.DB
}

func (p *PostgresQLDatabase) AddLink(original string, shorten string) (string, error) {
	//TODO implement me
	panic("i can't do anything now..")
}

func (p *PostgresQLDatabase) GetFullLink(hash string) (string, error) {
	//TODO implement me
	panic("i can't do anything now..")
}

func (p *PostgresQLDatabase) Close() error {
	return p.driver.Close()
}

func (p *PostgresQLDatabase) PingConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return p.driver.PingContext(ctx)
}

func NewPQLDatabase(params string) (Database, error) {
	fmt.Println(params)
	db, err := sql.Open("pgx", params)
	if err != nil {
		return nil, err
	}
	return &PostgresQLDatabase{
		driver: db,
	}, nil
}
