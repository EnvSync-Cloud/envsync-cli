package helper

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/models"
)

func CheckProjectConfig(filename ...string) error {
	cfg := constants.DefaultProjectConfig

	if len(filename) > 0 {
		cfg = filename[0]
	}

	if _, err := os.Stat(cfg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("project configuration file not found")
		}
		return err
	}

	return nil
}

func LoadProjectConfig(filename ...string) (models.ProjectEnvConfig, error) {
	cfg := constants.DefaultProjectConfig
	var config models.ProjectEnvConfig

	if len(filename) > 0 {
		cfg = filename[0]
	}

	if _, err := toml.DecodeFile(cfg, &config); err != nil {
		return models.ProjectEnvConfig{}, err
	}

	return models.ProjectEnvConfig{}, nil
}
