package database

import (
	"errors"
)

type Database struct {
	storage map[string]string
}

func CreateDatabase() *Database {
	return &Database{
		storage: make(map[string]string),
	}
}

func (db *Database) AddLink(original string, shorten string) (string, error) {
	if _, err := db.GetFullLink(shorten); err == nil {
		return "", errors.New("такая запись в БД уже есть")
	}
	db.storage[shorten] = original
	return shorten, nil
}

func (db *Database) GetFullLink(hash string) (string, error) {
	value, ok := db.storage[hash]
	if ok {
		return value, nil
	}
	return value, errors.New("нет такой записи в БД")
}
