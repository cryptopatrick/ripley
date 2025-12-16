package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadWithDefaults(t *testing.T) {
	cfg := LoadWithDefaults()

	if cfg.Daemon.Interval != "30m" {
		t.Errorf("Expected default interval '30m', got '%s'", cfg.Daemon.Interval)
	}

	if cfg.Daemon.DBPath != "./ripley.db" {
		t.Errorf("Expected default db_path './ripley.db', got '%s'", cfg.Daemon.DBPath)
	}

	if cfg.Claude.Model != "Sonnet" {
		t.Errorf("Expected default model 'Sonnet', got '%s'", cfg.Claude.Model)
	}

	if cfg.Monitoring.RollingWindow != 10 {
		t.Errorf("Expected default rolling_window 10, got %d", cfg.Monitoring.RollingWindow)
	}

	if cfg.Monitoring.WarningThreshold != 0.7 {
		t.Errorf("Expected default warning_threshold 0.7, got %.2f", cfg.Monitoring.WarningThreshold)
	}
}

func TestGetInterval(t *testing.T) {
	cfg := LoadWithDefaults()

	duration, err := cfg.GetInterval()
	if err != nil {
		t.Fatalf("GetInterval() failed: %v", err)
	}

	expected := 30 * time.Minute
	if duration != expected {
		t.Errorf("Expected duration %v, got %v", expected, duration)
	}
}

func TestGetIntervalInvalid(t *testing.T) {
	cfg := &Config{}
	cfg.Daemon.Interval = "invalid"

	_, err := cfg.GetInterval()
	if err == nil {
		t.Error("Expected error for invalid interval, got nil")
	}
}

func TestLoadValidConfig(t *testing.T) {
	// Create temporary config file
	content := `
daemon:
  interval: "15m"
  db_path: "/tmp/test.db"
claude:
  model: "Haiku"
  default_max_tokens: 100
monitoring:
  rolling_window: 5
  warning_threshold: 0.8
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if cfg.Daemon.Interval != "15m" {
		t.Errorf("Expected interval '15m', got '%s'", cfg.Daemon.Interval)
	}

	if cfg.Daemon.DBPath != "/tmp/test.db" {
		t.Errorf("Expected db_path '/tmp/test.db', got '%s'", cfg.Daemon.DBPath)
	}

	if cfg.Claude.Model != "Haiku" {
		t.Errorf("Expected model 'Haiku', got '%s'", cfg.Claude.Model)
	}

	if cfg.Monitoring.RollingWindow != 5 {
		t.Errorf("Expected rolling_window 5, got %d", cfg.Monitoring.RollingWindow)
	}

	if cfg.Monitoring.WarningThreshold != 0.8 {
		t.Errorf("Expected warning_threshold 0.8, got %.2f", cfg.Monitoring.WarningThreshold)
	}
}

func TestLoadInvalidInterval(t *testing.T) {
	content := `
daemon:
  interval: "not-a-duration"
  db_path: "./test.db"
claude:
  model: "Sonnet"
monitoring:
  rolling_window: 10
  warning_threshold: 0.7
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid interval, got nil")
	}
}

func TestLoadMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "missing interval",
			content: `
daemon:
  db_path: "./test.db"
claude:
  model: "Sonnet"
monitoring:
  rolling_window: 10
  warning_threshold: 0.7
`,
		},
		{
			name: "missing db_path",
			content: `
daemon:
  interval: "30m"
claude:
  model: "Sonnet"
monitoring:
  rolling_window: 10
  warning_threshold: 0.7
`,
		},
		{
			name: "missing model",
			content: `
daemon:
  interval: "30m"
  db_path: "./test.db"
claude:
  default_max_tokens: 200
monitoring:
  rolling_window: 10
  warning_threshold: 0.7
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "config-*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			_, err = Load(tmpfile.Name())
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		expectErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Daemon: struct {
					Interval string `yaml:"interval"`
					DBPath   string `yaml:"db_path"`
				}{Interval: "30m", DBPath: "./test.db"},
				Claude: struct {
					Model            string `yaml:"model"`
					DefaultMaxTokens int    `yaml:"default_max_tokens"`
				}{Model: "Sonnet", DefaultMaxTokens: 200},
				Monitoring: struct {
					RollingWindow    int     `yaml:"rolling_window"`
					WarningThreshold float64 `yaml:"warning_threshold"`
				}{RollingWindow: 10, WarningThreshold: 0.7},
			},
			expectErr: false,
		},
		{
			name: "invalid threshold high",
			cfg: &Config{
				Daemon: struct {
					Interval string `yaml:"interval"`
					DBPath   string `yaml:"db_path"`
				}{Interval: "30m", DBPath: "./test.db"},
				Claude: struct {
					Model            string `yaml:"model"`
					DefaultMaxTokens int    `yaml:"default_max_tokens"`
				}{Model: "Sonnet", DefaultMaxTokens: 200},
				Monitoring: struct {
					RollingWindow    int     `yaml:"rolling_window"`
					WarningThreshold float64 `yaml:"warning_threshold"`
				}{RollingWindow: 10, WarningThreshold: 1.5},
			},
			expectErr: true,
		},
		{
			name: "invalid threshold low",
			cfg: &Config{
				Daemon: struct {
					Interval string `yaml:"interval"`
					DBPath   string `yaml:"db_path"`
				}{Interval: "30m", DBPath: "./test.db"},
				Claude: struct {
					Model            string `yaml:"model"`
					DefaultMaxTokens int    `yaml:"default_max_tokens"`
				}{Model: "Sonnet", DefaultMaxTokens: 200},
				Monitoring: struct {
					RollingWindow    int     `yaml:"rolling_window"`
					WarningThreshold float64 `yaml:"warning_threshold"`
				}{RollingWindow: 10, WarningThreshold: -0.1},
			},
			expectErr: true,
		},
		{
			name: "invalid rolling window",
			cfg: &Config{
				Daemon: struct {
					Interval string `yaml:"interval"`
					DBPath   string `yaml:"db_path"`
				}{Interval: "30m", DBPath: "./test.db"},
				Claude: struct {
					Model            string `yaml:"model"`
					DefaultMaxTokens int    `yaml:"default_max_tokens"`
				}{Model: "Sonnet", DefaultMaxTokens: 200},
				Monitoring: struct {
					RollingWindow    int     `yaml:"rolling_window"`
					WarningThreshold float64 `yaml:"warning_threshold"`
				}{RollingWindow: 0, WarningThreshold: 0.7},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("validate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
