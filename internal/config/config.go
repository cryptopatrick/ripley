// Package config provides configuration management for the Ripley daemon.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the Ripley daemon.
type Config struct {
	Daemon struct {
		Interval string `yaml:"interval"` // e.g. "30m", "1h"
		DBPath   string `yaml:"db_path"`
	} `yaml:"daemon"`

	Claude struct {
		Model            string `yaml:"model"`
		DefaultMaxTokens int    `yaml:"default_max_tokens"`
	} `yaml:"claude"`

	Monitoring struct {
		RollingWindow    int     `yaml:"rolling_window"`
		WarningThreshold float64 `yaml:"warning_threshold"`
	} `yaml:"monitoring"`
}

// Load reads and parses a YAML configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// LoadWithDefaults returns a Config with sensible default values.
// Use this when no config file is provided.
func LoadWithDefaults() *Config {
	cfg := &Config{}

	cfg.Daemon.Interval = "30m"
	cfg.Daemon.DBPath = "./ripley.db"

	cfg.Claude.Model = "Sonnet"
	cfg.Claude.DefaultMaxTokens = 200

	cfg.Monitoring.RollingWindow = 10
	cfg.Monitoring.WarningThreshold = 0.7

	return cfg
}

// GetInterval parses the interval string and returns a time.Duration.
func (c *Config) GetInterval() (time.Duration, error) {
	duration, err := time.ParseDuration(c.Daemon.Interval)
	if err != nil {
		return 0, fmt.Errorf("invalid interval format: %w", err)
	}
	return duration, nil
}

// validate checks that all required fields are set and valid.
func (c *Config) validate() error {
	if c.Daemon.Interval == "" {
		return fmt.Errorf("daemon.interval is required")
	}

	if _, err := time.ParseDuration(c.Daemon.Interval); err != nil {
		return fmt.Errorf("daemon.interval must be a valid duration (e.g. '30m', '1h'): %w", err)
	}

	if c.Daemon.DBPath == "" {
		return fmt.Errorf("daemon.db_path is required")
	}

	if c.Claude.Model == "" {
		return fmt.Errorf("claude.model is required")
	}

	if c.Monitoring.RollingWindow <= 0 {
		return fmt.Errorf("monitoring.rolling_window must be positive")
	}

	if c.Monitoring.WarningThreshold < 0 || c.Monitoring.WarningThreshold > 1 {
		return fmt.Errorf("monitoring.warning_threshold must be between 0 and 1")
	}

	return nil
}
