package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type AppConfig struct {
	IdToken string `mapstructure:"idToken"`
	APIUrl  string `mapstructure:"apiUrl"`
}

func New() (AppConfig, error) {
	var cfg AppConfig

	configDir, err := os.UserConfigDir()
	if err != nil {
		return cfg, err
	}
	appConfigDir := filepath.Join(configDir, "envsync-cli")
	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		return cfg, err
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(appConfigDir)

	v.SetDefault("idToken", "")
	v.SetDefault("apiUrl", "")

	// Read or initialize config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			_ = v.WriteConfigAs(filepath.Join(appConfigDir, "config.yaml"))
		} else {
			return cfg, err
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *AppConfig) AddConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	appConfigDir := filepath.Join(configDir, "envsync-cli")

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(appConfigDir)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.Set("idToken", c.IdToken)
	v.Set("apiUrl", c.APIUrl)

	return v.WriteConfigAs(filepath.Join(appConfigDir, "config.yaml"))
}

func (c *AppConfig) ReadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	appConfigDir := filepath.Join(configDir, "envsync-cli")

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(appConfigDir)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(c)
}
