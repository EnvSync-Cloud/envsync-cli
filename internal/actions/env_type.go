package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func ListEnvTypes() cli.ActionFunc {
	return func(c context.Context, cmd *cli.Command) error {
		envTypeService := services.NewEnvTypeService()

		// Get user info from the auth service
		envTypes, err := envTypeService.GetAllEnvTypes()

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
			cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}

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
