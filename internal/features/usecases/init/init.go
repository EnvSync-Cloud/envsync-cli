package init

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type initCaseUse struct {
	appService services.ApplicationService
	tui        *factory.InitFactory
}

func NewInitUseCase() InitUseCase {
	tui := factory.NewInitFactory()
	appService := services.NewAppService()
	return &initCaseUse{
		appService: appService,
		tui:        tui,
	}
}

func (uc *initCaseUse) Execute(ctx context.Context, config string) error {
	// Check if the configuration file already exists
	if err := uc.checkConfigExists(config); err == nil {
		// If it already exists, return an warning
		return errors.New("configuration file already exists")
	}

	// Fetch all the applications
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return fmt.Errorf("failed to retrieve applications: %w", err)
	}

	// Open a form to collect user input for the configuration
	appID, envID, err := uc.tui.OpenInitForm(apps)
	if err != nil {
		return fmt.Errorf("failed to open configuration form: %w", err)
	}

	syncConfig := domain.SyncConfig{
		AppID:     appID,
		EnvTypeID: envID,
	}

	if err := uc.saveConfig(syncConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

func (uc *initCaseUse) checkConfigExists(configPath string) error {
	// Check if the configuration file exists at the specified path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist at path: %s", configPath)
	}
	return nil
}

func (uc *initCaseUse) saveConfig(cfg domain.SyncConfig) error {
	file, err := os.Create(constants.DefaultProjectConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
