package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type AppConfig struct {
	AccessToken string `json:"access_token"`
	BackendURL  string `json:"backend_url"`
}

var cfg AppConfig
var once sync.Once

const configDirPath = "/.config/envsync"

func New() AppConfig {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		filePath := filepath.Join(home, configDirPath, "config.json")

		// Ensure directory exists
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(err)
		}

		// Create file if it doesn't exist
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			file, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			file.Close()
		}

		cfg, err = ReadConfigFile()
	})

	return cfg
}

func (c *AppConfig) WriteConfigFile() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(home, configDirPath, "config.json")

	return os.WriteFile(filePath, data, 0644)
}

func ReadConfigFile() (AppConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(home, configDirPath, "config.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
