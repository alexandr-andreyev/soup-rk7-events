package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	External ExternalConfig `yaml:"external"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type ExternalConfig struct {
	URL            string `yaml:"url"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	MaxRetries     int    `yaml:"max_retries"`
	RetryDelay     int    `yaml:"retry_delay_seconds"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults if not specified
	if cfg.Server.Port == "" {
		cfg.Server.Port = ":7080"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10
	}
	if cfg.External.TimeoutSeconds == 0 {
		cfg.External.TimeoutSeconds = 10
	}
	if cfg.External.MaxRetries == 0 {
		cfg.External.MaxRetries = 3
	}
	if cfg.External.RetryDelay == 0 {
		cfg.External.RetryDelay = 1
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "orders.db"
	}

	return &cfg, nil
}
