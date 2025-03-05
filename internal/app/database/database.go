package database

import (
	"errors"
	"strconv"
)

type Database struct {
	storage map[string]string
}

func CreateDatabase() *Database {
	return &Database{
		storage: make(map[string]string),
	}
}

func (db *Database) AddRow(prefix string, link string) string {
	index := prefix + strconv.Itoa(len(db.storage)+1)
	db.storage[index] = link
	return index
}

func (db *Database) GetRow(link string) (string, error) {
	value, ok := db.storage[link]
	if ok {
		return value, nil
	}
	return value, errors.New("нет такой записи в БД")
}
