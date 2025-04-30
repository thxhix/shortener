package config

import (
	"flag"
	"os"

	parser "github.com/thxhix/shortener/internal/app/flags"
)

const DefaultAddress = "localhost:8080"
const DefaultBaseURL = "http://localhost:8080"
const DefaultDBFileName = "./db.json"
const DefaultPostgresQL = "user=postgres password=129755 dbname=yp_go sslmode=disable" // TODO : Это для себя сделал, пока пусть будет

type Config struct {
	Address    parser.Address
	BaseURL    parser.BaseURL
	DBFileName string
	PostgresQL string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		Address: parser.Address{},
		BaseURL: parser.BaseURL{},
	}

	err := cfg.Address.Set(DefaultAddress)
	if err != nil {
		return nil, err
	}
	err = cfg.BaseURL.Set(DefaultBaseURL)
	if err != nil {
		return nil, err
	}
	//cfg.DBFileName = DefaultDBFileName
	//cfg.PostgresQL = DefaultPostgresQL

	cfg.ParseFlags()
	err = cfg.LoadEnv()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) LoadEnv() error {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		err := c.Address.Set(envAddr)
		if err != nil {
			return err
		}
	}
	if envBase := os.Getenv("BASE_URL"); envBase != "" {
		err := c.BaseURL.Set(envBase)
		if err != nil {
			return err
		}
	}
	if envFile := os.Getenv("FILE_STORAGE_PATH"); envFile != "" {
		c.DBFileName = envFile
	}
	if envFile := os.Getenv("DATABASE_DSN"); envFile != "" && c.PostgresQL == "" {
		c.PostgresQL = envFile
	}
	return nil
}

func (c *Config) ParseFlags() {
	flag.Var(&c.Address, "a", "Address (например, localhost:8080)")
	flag.Var(&c.BaseURL, "b", "Base URL (например, http://example.com:8080)")

	dbFileName := flag.String("f", DefaultDBFileName, "Путь к файлу БД (например, ./db.json)")
	postgresDSN := flag.String("d", "", "Данные для PostgresQL")

	flag.Parse()

	c.DBFileName = *dbFileName
	if c.PostgresQL == "" {
		c.PostgresQL = *postgresDSN
	}
}
