package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	API    APIConfig    `yaml:"api"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type APIConfig struct {
	MassiveAPIKey string `yaml:"massive_api_key"`
}

// Load reads configuration from file
func Load(configPath string) (*Config, error) {
	cfg := &Config{}
	if configPath == "" {
		return nil, fmt.Errorf("config file path is required")
	}
	if err := loadFromYAML(cfg, configPath); err != nil {
		return nil, fmt.Errorf("failed to load config from YAML: %w", err)
	}
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func loadFromYAML(cfg *Config, configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Expand environment variables in YAML
	expanded := os.ExpandEnv(string(data))

	return yaml.Unmarshal([]byte(expanded), cfg)
}

func validate(cfg *Config) error {
	if cfg.API.MassiveAPIKey == "" || cfg.API.MassiveAPIKey == "your_massive_api_key_here" {
		return fmt.Errorf("Massive.com API key is required. Please set MASSIVE_API_KEY environment variable or update the config file. Get your API key at https://www.massive.com/")
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	return nil
}

// GetAddress returns the server address in host:port format
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}