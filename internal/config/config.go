package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Dir      string
	FilePath string
}

func NewConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return &Config{}, err
	}

	dir := filepath.Join(configDir, "papyro")
	fullPath := filepath.Join(dir, "profiles.json")

	return &Config{
		Dir:      dir,
		FilePath: fullPath,
	}, nil
}

func (c *Config) Exists() bool {
	_, err := os.Stat(c.FilePath)
	return !os.IsNotExist(err)
}

func (c *Config) CreateIfNotExists() error {
	if c.Exists() {
		return nil
	}

	err := os.MkdirAll(c.Dir, 0755)
	if err != nil {
		return err
	}

	data := map[string][]map[string]string{
		"profiles": {},
	}
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.FilePath, jsonBytes, 0600)
}

func (c *Config) CheckAndCreateConfiguration() {
	if err := c.CreateIfNotExists(); err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}
}
