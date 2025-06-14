package actions

import (
	"context"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func InitAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		app := cmd.String("app")
		envType := cmd.String("env-type")

		// Initialize the application service
		appService := services.NewAppService()

		// Get All Apps
		apps, err := appService.GetAllApps()
		if err != nil {
			return fmt.Errorf("failed to fetch applications: %w", err)
		}

		// Check if the app exists
		var appID string
		for _, a := range apps {
			if a.Name == app {
				appID = a.ID
				break
			}
		}
		if appID == "" {
			return fmt.Errorf("application '%s' not found", app)
		}

		// Get All Environment Types
		envTypes, err := appService.ReadAppEnvTypes()
		if err != nil {
			return fmt.Errorf("failed to fetch environment types: %w", err)
		}

		// Check if the environment type exists
		var envTypeID string
		for _, et := range envTypes {
			if et.Name == envType {
				envTypeID = et.ID
				break
			}
		}
		if envTypeID == "" {
			return fmt.Errorf("environment type '%s' not found", envType)
		}

		config := domain.SyncConfig{
			AppID:     appID,
			EnvTypeID: envTypeID,
		}

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

		err = toml.NewEncoder(file).Encode(config)
		if err != nil {
			return err
		}

		fmt.Println("Project config file generated successfully.")

		return nil
	}
}
