package models

//go:generate easyjson -all models.go

//easyjson:json
type FullURL struct {
	URL string `json:"url"`
}

//easyjson:json
type ShortURL struct {
	Result string `json:"result"`
}
