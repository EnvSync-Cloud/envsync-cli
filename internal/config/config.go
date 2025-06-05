package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type AppConfig struct {
	AccessToken string `json:"access_token"`
}

var cfg AppConfig
var once sync.Once

const configDirPath = "/.config/envsync"

func New() (AppConfig, error) {
	var err error

	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			err = fmt.Errorf("failed to get user home directory: %w", err)
			return
		}

		//Get absolute path of config directory
		dirPath, err := filepath.Abs(home + configDirPath)
		if err != nil {
			err = fmt.Errorf("failed to get absolute path of config directory: %w", err)
			return
		}

		// Check if directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				err = fmt.Errorf("failed to create config directory: %w", err)
				return
			}
		}

		cfg, err = ReadConfigFile(dirPath + "/config.json")
	})

	return cfg, err
}

func (c *AppConfig) WriteConfigFile(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func ReadConfigFile(filePath string) (AppConfig, error) {
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
