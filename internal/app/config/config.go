package config

import (
	"flag"
	"os"

	parser "github.com/thxhix/shortener/internal/app/flags"
)

const DEFAULT_ADDRESS = "localhost:8080"
const DEFAULT_BASE_URL = "http://localhost:8080"

type Config struct {
	Address parser.Address
	BaseURL parser.BaseURL
}

func InitConfig() *Config {
	cfg := &Config{
		Address: parser.Address{},
		BaseURL: parser.BaseURL{},
	}

	cfg.Address.Set(DEFAULT_ADDRESS)
	cfg.BaseURL.Set(DEFAULT_BASE_URL)

	cfg.ParseFlags()
	cfg.LoadEnv()

	return cfg
}

func (c *Config) LoadEnv() {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		c.Address.Set(envAddr)
	}
	if envBase := os.Getenv("BASE_URL"); envBase != "" {
		c.BaseURL.Set(envBase)
	}
}

func (c *Config) ParseFlags() {
	flag.Var(&c.Address, "a", "Address (например, localhost:8080)")
	flag.Var(&c.BaseURL, "d", "Base URL (например, http://example.com:8080)")

	flag.Parse()
}
