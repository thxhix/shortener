package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DBFileName string `env:"FILE_STORAGE_PATH" envDefault:"./db.json"`
	PostgresQL string `env:"DATABASE_DSN" envDefault:"user=postgres password=129755 dbname=yp_go sslmode=disable"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Переопределяем значениями из флагов, если переданы
	cfg.parseFlags()

	return cfg, nil
}

func (c *Config) parseFlags() {
	address := flag.String("a", c.Address, "Address (например, localhost:8080)")
	baseURL := flag.String("b", c.BaseURL, "Base URL (например, http://example.com:8080)")
	dbFile := flag.String("f", c.DBFileName, "Путь к файлу БД (например, ./db.json)")
	postgres := flag.String("d", c.PostgresQL, "PostgreSQL DSN")

	flag.Parse()

	c.Address = *address
	c.BaseURL = *baseURL
	c.DBFileName = *dbFile
	c.PostgresQL = *postgres
}
