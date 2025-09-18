package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type fileConfig struct {
	Address     *string `json:"server_address"`
	BaseURL     *string `json:"base_url"`
	DBFileName  *string `json:"file_storage_path"`
	PostgresQL  *string `json:"database_dsn"`
	EnableHTTPS *bool   `json:"enable_https"`

	SecretKey          *string `json:"secret_key"`
	DeleteWorkersCount *int    `json:"delete_workers_count"`
	DeleteBatchSize    *int    `json:"delete_batch_size"`
	EnableProfiler     *bool   `json:"enable_profiler"`
	ProfilerAddress    *string `json:"profiler_address"`
}

func (c *Config) parseFile() error {
	configPath := os.Getenv("CONFIG")
	if v, ok := getConfigFileValue(); ok {
		configPath = v
	}
	if configPath != "" {
		return readConfigFile(c, configPath)
	}
	return nil
}

func readConfigFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	var fc fileConfig
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&fc); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	if fc.Address != nil {
		cfg.Address = *fc.Address
	}
	if fc.BaseURL != nil {
		cfg.BaseURL = *fc.BaseURL
	}
	if fc.DBFileName != nil {
		cfg.DBFileName = *fc.DBFileName
	}
	if fc.PostgresQL != nil {
		cfg.PostgresQL = *fc.PostgresQL
	}
	if fc.EnableHTTPS != nil {
		cfg.EnableHTTPS = *fc.EnableHTTPS
	}

	if fc.SecretKey != nil {
		cfg.SecretKey = *fc.SecretKey
	}
	if fc.DeleteWorkersCount != nil {
		cfg.DeleteWorkersCount = *fc.DeleteWorkersCount
	}
	if fc.DeleteBatchSize != nil {
		cfg.DeleteBatchSize = *fc.DeleteBatchSize
	}
	if fc.EnableProfiler != nil {
		cfg.EnableProfiler = *fc.EnableProfiler
	}
	if fc.ProfilerAddress != nil {
		cfg.ProfilerAddress = *fc.ProfilerAddress
	}

	return nil
}

// getConfigFileValue scans os.Args for -c or -config flags.
// Returns the value and true if found, otherwise empty string and false.
func getConfigFileValue() (string, bool) {
	args := os.Args
	for i := 0; i < len(args); i++ {
		// short with pattern: -c {{value}}
		if args[i] == "-c" && i+1 < len(args) {
			return args[i+1], true
		}
		// short with pattern: -c={{value}}
		if strings.HasPrefix(args[i], "-c=") {
			return args[i][len("-c="):], true
		}

		// long with pattern: -config {{value}}
		if args[i] == "-config" && i+1 < len(args) {
			return args[i+1], true
		}
		// long with pattern: -config={{value}}
		if strings.HasPrefix(args[i], "-config=") {
			return args[i][len("-config="):], true
		}
	}
	return "", false
}
