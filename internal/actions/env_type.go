package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func GetEnvTypeByID() cli.ActionFunc {
	return func(c context.Context, cmd *cli.Command) error {
		envTypeService := services.NewEnvTypeService()

		envTypeID := cmd.String("id")
		if envTypeID == "" {
			return cli.Exit("Please provide an environment type ID using --id flag", 1)
		}

		envType, err := envTypeService.GetEnvTypeByID(envTypeID)
		if err != nil {
			return err
		}

		// If JSON flag is set, print in JSON format
		if cmd.Bool("json") {
			jsonOutput, err := json.MarshalIndent(envType, "", "  ")
			if err != nil {
				return err
			}
			cmd.Writer.Write([]byte(jsonOutput))
			return nil
		}

		// Print environment type details
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		cmd.Writer.Write([]byte("ID: " + envType.ID + "\n"))
		cmd.Writer.Write([]byte("Name: " + envType.Name + "\n"))
		cmd.Writer.Write([]byte("AppID: " + envType.AppID + "\n"))
		cmd.Writer.Write([]byte("IsDefault: " + fmt.Sprintf("%t", envType.IsDefault) + "\n"))
		cmd.Writer.Write([]byte("IsProtected: " + fmt.Sprintf("%t", envType.IsProtected) + "\n"))
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}

func GetEnvTypesByApp() cli.ActionFunc {
	return func(c context.Context, cmd *cli.Command) error {
		envTypeService := services.NewEnvTypeService()

		appID := cmd.String("app-id")
		if appID == "" {
			return cli.Exit("Please provide an app ID using --app-id flag", 1)
		}

		envTypes, err := envTypeService.GetEnvTypeByAppID(appID)
		if err != nil {
			return err
		}

		// If JSON flag is set, print in JSON format
		if cmd.Bool("json") {
			jsonOutput, err := json.MarshalIndent(envTypes, "", "  ")
			if err != nil {
				return err
			}
			cmd.Writer.Write([]byte(jsonOutput))
			return nil
		}

		// Print environment types
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		for _, envType := range envTypes {
			cmd.Writer.Write([]byte("ID: " + envType.ID + "\n"))
			cmd.Writer.Write([]byte("Name: " + envType.Name + "\n"))
			cmd.Writer.Write([]byte("AppID: " + envType.AppID + "\n"))
			cmd.Writer.Write([]byte("IsDefault: " + fmt.Sprintf("%t", envType.IsDefault) + "\n"))
			cmd.Writer.Write([]byte("IsProtected: " + fmt.Sprintf("%t", envType.IsProtected) + "\n"))
			cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}

func SwitchEnvType() cli.ActionFunc {
	return func(c context.Context, cmd *cli.Command) error {
		es := services.NewEnvTypeService()
		ss := services.NewSyncService()

		syncConfig, err := ss.ReadConfigData()
		if err != nil {
			return fmt.Errorf("failed to read sync config: %w", err)
		}

		envs, err := es.GetEnvTypeByAppID(syncConfig.AppID)
		if err != nil {
			return fmt.Errorf("failed to fetch environment types: %w", err)
		}

		if len(envs) == 0 {
			return cli.Exit("No environment types found for the current app.", 1)
		}

		cmd.Writer.Write([]byte("Available Environment Types:\n"))
		for i, env := range envs {
			cmd.Writer.Write([]byte(fmt.Sprintf("%d. %s (ID: %s)\n", i+1, env.Name, env.ID)))
		}
		cmd.Writer.Write([]byte("Please enter the number of the environment type you want to switch to: "))
		var choice int
		_, err = fmt.Scanf("%d", &choice)
		if err != nil || choice < 1 || choice > len(envs) {
			return cli.Exit("Invalid choice. Please enter a valid number.", 1)
		}

		selectedEnv := envs[choice-1]
		if err := ss.WriteConfigData(domain.SyncConfig{
			AppID:     syncConfig.AppID,
			EnvTypeID: selectedEnv.ID,
		}); err != nil {
			return fmt.Errorf("failed to write sync config: %w", err)
		}

		return nil
	}
}
