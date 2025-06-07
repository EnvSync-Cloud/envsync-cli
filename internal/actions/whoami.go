package actions

import (
	"encoding/json"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func Whoami() cli.ActionFunc {
	return func(c *cli.Context) error {
		authService := services.NewAuthService()

		// Get user info from the auth service
		userInfo, err := authService.Whoami()

		if err != nil {
			return err
		}

		// If JSON flag is set, print in JSON format
		if c.Bool("json") {
			jsonOutput, err := json.MarshalIndent(userInfo, "", "  ")
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(jsonOutput))
			return nil
		}

		// Print user info
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		c.App.Writer.Write([]byte("User ID: " + userInfo.UserId + "\n"))
		c.App.Writer.Write([]byte("Email: " + userInfo.Email + "\n"))
		c.App.Writer.Write([]byte("Organization: " + userInfo.Org + "\n"))
		c.App.Writer.Write([]byte("Role: " + userInfo.Role + "\n"))
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}
