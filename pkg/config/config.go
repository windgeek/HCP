package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the HCP configuration.
type Config struct {
	IdentityKeyPath string `yaml:"identity_key_path"`
	// Add more config fields here as needed
}

// DefaultConfig returns the default configuration.
func DefaultConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Config{
		IdentityKeyPath: filepath.Join(home, ".hcp", "identity.key"),
	}, nil
}

// LoadConfig loads configuration with priority:
// 1. CLI Override (passed as argument)
// 2. Env Var (HCP_KEY_PATH)
// 3. Config File (.hcp/config.yaml or ~/.hcp/config.yaml)
// 4. Default (~/.hcp/identity.key)
func LoadConfig(cliKeyPath string) (*Config, error) {
	// Start with defaults
	cfg, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	// 3. Load from Config File
	// Check local .hcp/config.yaml
	localConf := ".hcp/config.yaml"
	if _, err := os.Stat(localConf); err == nil {
		if err := loadFromFile(localConf, cfg); err != nil {
			return nil, err
		}
	} else {
		// Check home ~/.hcp/config.yaml
		home, _ := os.UserHomeDir()
		homeConf := filepath.Join(home, ".hcp", "config.yaml")
		if _, err := os.Stat(homeConf); err == nil {
			if err := loadFromFile(homeConf, cfg); err != nil {
				return nil, err
			}
		}
	}

	// 2. Env Var Override
	if envPath := os.Getenv("HCP_KEY_PATH"); envPath != "" {
		cfg.IdentityKeyPath = envPath
	}

	// 1. CLI Override
	if cliKeyPath != "" {
		cfg.IdentityKeyPath = cliKeyPath
	}

	// Ensure absolute path for key
	absPath, err := filepath.Abs(cfg.IdentityKeyPath)
	if err != nil {
		return nil, fmt.Errorf("invalid key path: %w", err)
	}
	cfg.IdentityKeyPath = absPath

	return cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

// SaveConfig saves the configuration to a file.
func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
