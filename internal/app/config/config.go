package config

import (
	"flag"
	"fmt"
	"os"

	parser "github.com/thxhix/shortener/internal/app/flags"
)

const DefaultAddress = "localhost:8080"
const DefaultBaseURL = "http://localhost:8080"
const DefaultDBFileName = "./db.json"
const DefaultPostgresQL = "postgres://postgres:129755@localhost:5432/yp_go"

type Config struct {
	Address    parser.Address
	BaseURL    parser.BaseURL
	DBFileName string
	PostgresQL string
}

func NewConfig() *Config {
	cfg := &Config{
		Address: parser.Address{},
		BaseURL: parser.BaseURL{},
	}

	cfg.Address.Set(DefaultAddress)
	cfg.BaseURL.Set(DefaultBaseURL)
	cfg.DBFileName = DefaultDBFileName
	//cfg.PostgresQL = DefaultPostgresQL

	cfg.ParseFlags()
	cfg.LoadEnv()

	fmt.Println(cfg.PostgresQL)

	return cfg
}

func (c *Config) LoadEnv() {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		c.Address.Set(envAddr)
	}
	if envBase := os.Getenv("BASE_URL"); envBase != "" {
		c.BaseURL.Set(envBase)
	}
	if envFile := os.Getenv("FILE_STORAGE_PATH"); envFile != "" {
		c.DBFileName = envFile
	}
	if envFile := os.Getenv("DATABASE_DSN"); envFile != "" {
		c.PostgresQL = envFile
	}
}

func (c *Config) ParseFlags() {
	flag.Var(&c.Address, "a", "Address (например, localhost:8080)")
	flag.Var(&c.BaseURL, "b", "Base URL (например, http://example.com:8080)")

	dbFileName := flag.String("f", DefaultDBFileName, "Путь к файлу БД (например, ./db.json)")
	postgresDSN := flag.String("d", "", "Данные для PostgresQL")

	flag.Parse()

	c.DBFileName = *dbFileName
	c.PostgresQL = *postgresDSN
}
