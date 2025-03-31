package database

type Database interface {
	AddLink(original string, shorten string) (string, error)
	GetFullLink(hash string) (string, error)
	Close() error
	PingConnection() error
}

type LinkRow struct {
	UUID int    `json:"uuid"`
	Hash string `json:"hash"`
	URL  string `json:"url"`
}
