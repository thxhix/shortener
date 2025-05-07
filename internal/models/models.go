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
type DBShortenRowList []DBShortenRow

//easyjson:json
type DBShortenRow struct {
	ID   int       `json:"id"`
	Hash string    `json:"hash"`
	URL  string    `json:"url"`
	Time time.Time `json:"time"`
}

//easyjson:json
type BatchShortenRequestList []BatchShortenRequest

//easyjson:json
type BatchShortenRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

//easyjson:json
type BatchShortenResponseList []BatchShortenResponse

type BatchShortenResponse struct {
	ID   string `json:"correlation_id"`
	Hash string `json:"short_url"`
}
