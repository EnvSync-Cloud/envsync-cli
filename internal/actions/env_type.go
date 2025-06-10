package actions

import (
	"context"
	"encoding/json"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v3"
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
