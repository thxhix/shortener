package models

import "time"

//go:generate easyjson -all models.go

//easyjson:json
type FullURL struct {
	URL string `json:"url"`
}

//easyjson:json
type ShortURL struct {
	Result string `json:"result"`
}

//easyjson:json
type DatabaseRow struct {
	ID   int       `json:"id"`
	Hash string    `json:"hash"`
	URL  string    `json:"url"`
	Time time.Time `json:"time"`
}
