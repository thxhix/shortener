package config

import (
	"flag"
	"os"

	parser "github.com/thxhix/shortener/internal/app/flags"
)

const DefaulAddress = "localhost:8080"
const DefaulBaseURL = "http://localhost:8080"

type Config struct {
	Address parser.Address
	BaseURL parser.BaseURL
}

func NewConfig() *Config {
	cfg := &Config{
		Address: parser.Address{},
		BaseURL: parser.BaseURL{},
	}

	cfg.Address.Set(DefaulAddress)
	cfg.BaseURL.Set(DefaulBaseURL)

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
	flag.Var(&c.BaseURL, "b", "Base URL (например, http://example.com:8080)")

	flag.Parse()
}
