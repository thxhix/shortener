package config

import (
	"flag"
	"os"

	parser "github.com/thxhix/shortener/internal/app/flags"
)

const DefaulAddress = "localhost:8080"
const DefaulBaseURL = "http://localhost:8080"

type Config struct {
	Address string
	BaseURL parser.BaseURL
}

func InitConfig() *Config {
	cfg := &Config{
		Address: DefaulAddress,
		BaseURL: parser.BaseURL{},
	}

	cfg.BaseURL.Set(DefaulBaseURL)

	cfg.ParseFlags()
	cfg.LoadEnv()

	return cfg
}

func (c *Config) LoadEnv() {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		c.Address = envAddr
	}
	if envBase := os.Getenv("BASE_URL"); envBase != "" {
		c.BaseURL.Set(envBase)
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.Address, "a", c.Address, "Address (например, localhost:8080)")
	flag.Var(&c.BaseURL, "b", "Base URL (например, http://example.com:8080)")

	flag.Parse()
}
