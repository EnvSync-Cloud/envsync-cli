package actions

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/models"
	"github.com/urfave/cli/v2"
)

func GenConfigAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		appID := ctx.String("app-id")
		envType := ctx.String("env-type-id")

		config := models.ProjectEnvConfig{
			AppID:   appID,
			EnvType: envType,
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
