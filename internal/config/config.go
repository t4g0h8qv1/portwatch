// Package config handles loading and validating portwatch configuration
// from a YAML file.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	Target   string   `yaml:"target"`
	Ports    string   `yaml:"ports"`
	Baseline string   `yaml:"baseline"`
	Alert    Alert    `yaml:"alert"`
	Scan     ScanOpts `yaml:"scan"`
}

// Alert configures notification behaviour.
type Alert struct {
	Stdout  bool   `yaml:"stdout"`
	Webhook string `yaml:"webhook"`
}

// ScanOpts controls scanner tuning parameters.
type ScanOpts struct {
	TimeoutMs   int `yaml:"timeout_ms"`
	Concurrency int `yaml:"concurrency"`
}

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required fields are present and values are sensible.
func (c *Config) validate() error {
	if c.Target == "" {
		return fmt.Errorf("config: target must not be empty")
	}
	if c.Ports == "" {
		return fmt.Errorf("config: ports must not be empty")
	}
	if c.Scan.TimeoutMs <= 0 {
		c.Scan.TimeoutMs = 500
	}
	if c.Scan.Concurrency <= 0 {
		c.Scan.Concurrency = 100
	}
	return nil
}

// HasWebhook reports whether a webhook URL has been configured.
func (a *Alert) HasWebhook() bool {
	return a.Webhook != ""
}

// IsAlertEnabled reports whether at least one alert channel is configured.
func (a *Alert) IsAlertEnabled() bool {
	return a.Stdout || a.HasWebhook()
}
