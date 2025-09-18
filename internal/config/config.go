package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

// Config holds application configuration parameters.
// Load order / priority: defaults <- JSON file (-c/-config/CONFIG) <- ENV <- flags.
type Config struct {
	// Address specifies the HTTP server listen address, e.g. "localhost:8080".
	Address string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`

	// BaseURL defines the base URL for shortened links, e.g. "http://localhost:8080".
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`

	// DBFileName specifies the path to the JSON file used by the file storage driver.
	DBFileName string `env:"FILE_STORAGE_PATH" envDefault:"./db.json"`

	// PostgresQL is the PostgreSQL DSN (data source name).
	// If set, the service will use PostgreSQL as the database backend.
	PostgresQL string `env:"DATABASE_DSN"`

	// SecretKey is used for signing user authentication tokens.
	SecretKey string `env:"SECRET_KEY" envDefault:"secret"`

	// DeleteWorkersCount sets the number of concurrent workers for batch link deletion.
	DeleteWorkersCount int `env:"DELETE_WORKERS_COUNT" envDefault:"10"`

	// DeleteBatchSize sets the maximum batch size for deleting user links.
	DeleteBatchSize int `env:"DELETE_BATCH_SIZE" envDefault:"1000"`

	// EnableProfiler enables the built-in pprof profiler if true.
	EnableProfiler bool `env:"ENABLE_PROFILER" envDefault:"false"`

	// ProfilerAddress specifies the address for the pprof profiler, e.g. "localhost:9090".
	ProfilerAddress string `env:"PROFILER_ADDRESS" envDefault:"localhost:9090"`

	// EnableHTTPS turns server as HTTPS protocol, if true
	EnableHTTPS bool `env:"ENABLE_HTTPS" envDefault:"false"`
}

// NewConfig loads config from a JSON file (low priority), then ENV, then flags.
// Environment variables take precedence, but values can be overridden via flags.
// Returns a Config pointer or an error if parsing fails.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cfg.parseFile(); err != nil {
		return nil, err
	}

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
	enablePprof := flag.Bool("pprof", c.EnableProfiler, "Включить pprof (профайлер)")
	enableHTTPS := flag.Bool("s", c.EnableHTTPS, "enable HTTPS mode)")

	flag.Parse()

	c.Address = *address
	c.BaseURL = *baseURL
	c.DBFileName = *dbFile
	c.PostgresQL = *postgres
	c.EnableProfiler = *enablePprof
	c.EnableHTTPS = *enableHTTPS
}
