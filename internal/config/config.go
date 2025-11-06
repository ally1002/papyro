package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetConfig() (string, string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	dir := filepath.Join(homedir, ".config", "papyro")
	fullPath := filepath.Join(dir, "profiles.json")

	return dir, fullPath, nil
}

func CheckAndCreateConfiguration() {
	dir, fullPath, err := GetConfig()
	if err != nil {
		log.Fatalf("could not determine user home directory: %s", err)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Error creating configuration directory %s: %v", dir, err)
		}

		data := map[string][]map[string]string{
			"profiles": {},
		}
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		err = os.WriteFile(fullPath, jsonBytes, 0600)
		if err != nil {
			log.Fatalf("Error creating configuration file %s: %v", fullPath, err)
		}
	}
}
