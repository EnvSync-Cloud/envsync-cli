package actions

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type ProjectEnvConfig struct {
	AppID   string `toml:"app_id"`
	EnvType string `toml:"env_type"`
}

func GenConfigAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		appID := ctx.String("app-id")
		envType := ctx.String("env-type")

		config := ProjectEnvConfig{
			AppID:   appID,
			EnvType: envType,
		}

		// Check if the project config file exists
		if _, err := os.Stat("envsyncrc.toml"); err != nil {
			if os.IsNotExist(err) {
				os.Create("envsyncrc.toml")
			}
		}

		// Write the config to the file
		file, err := os.Create("envsyncrc.toml")
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
