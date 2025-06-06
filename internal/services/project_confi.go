package services

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/models"
)

type ProjectConfigService interface {
	ReadProjectConfig() (models.ProjectEnvConfig, error)
	WriteProjectConfig(projCfg models.ProjectEnvConfig) error
}

type projectConfig struct {
}

func NewProjectConfigService() ProjectConfigService {
	return &projectConfig{}
}

func (p *projectConfig) ReadProjectConfig() (models.ProjectEnvConfig, error) {
	var projCfg models.ProjectEnvConfig

	file, err := os.Open(constants.DefaultProjectConfig)
	if err != nil {
		return projCfg, err
	}
	defer file.Close()

	_, err = toml.NewDecoder(file).Decode(&projCfg)
	if err != nil {
		return projCfg, err
	}

	return projCfg, nil
}

func (p *projectConfig) WriteProjectConfig(projCfg models.ProjectEnvConfig) error {
	// Check if the project config file exists
	if _, err := os.Stat(constants.DefaultProjectConfig); err != nil {
		if os.IsNotExist(err) {
			os.Create(constants.DefaultProjectConfig)
		}
	}

	// Write the config to the file
	file, err := os.Create(constants.DefaultProjectConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the config to the file
	err = toml.NewEncoder(file).Encode(projCfg)
	if err != nil {
		return err
	}

	return nil
}
