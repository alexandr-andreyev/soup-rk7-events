package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig   `yaml:"server"`
	Subscribers []Subscriber   `yaml:"subscribers"`
	Database    DatabaseConfig `yaml:"database"`
	Logging     LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type Subscriber struct {
	Name               string            `yaml:"name"`
	URL                string            `yaml:"url"`
	Enabled            bool              `yaml:"enabled"`
	Timeout            int               `yaml:"timeout"`
	RetryCount         int               `yaml:"retry_count"`
	RetryDelaySeconds  int               `yaml:"retry_delay_seconds"`
	TableCodes         []int             `yaml:"table_codes"`
	EventTypes         []string          `yaml:"event_types"`
	Headers            map[string]string `yaml:"headers"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type LoggingConfig struct {
	Directory string `yaml:"directory"`
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

	// Set defaults for subscribers
	for i := range cfg.Subscribers {
		if cfg.Subscribers[i].Timeout == 0 {
			cfg.Subscribers[i].Timeout = 5
		}
		if cfg.Subscribers[i].RetryCount == 0 {
			cfg.Subscribers[i].RetryCount = 3
		}
		if cfg.Subscribers[i].RetryDelaySeconds == 0 {
			cfg.Subscribers[i].RetryDelaySeconds = 1
		}
	}

	if cfg.Database.Path == "" {
		cfg.Database.Path = "orders.db"
	}
	if cfg.Logging.Directory == "" {
		cfg.Logging.Directory = "logs"
	}

	return &cfg, nil
}
