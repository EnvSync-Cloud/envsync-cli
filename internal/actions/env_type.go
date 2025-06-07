package actions

import (
	"encoding/json"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func ListEnvTypes() cli.ActionFunc {
	return func(c *cli.Context) error {
		envTypeService := services.NewEnvTypeService()

		// Get user info from the auth service
		envTypes, err := envTypeService.GetAllEnvTypes()

		if err != nil {
			return err
		}

		// If JSON flag is set, print in JSON format
		if c.Bool("json") {
			jsonOutput, err := json.MarshalIndent(envTypes, "", "  ")
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(jsonOutput))
			return nil
		}

		// Print environment types
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		for _, envType := range envTypes {
			c.App.Writer.Write([]byte("ID: " + envType.ID + "\n"))
			c.App.Writer.Write([]byte("Name: " + envType.Name + "\n"))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}
