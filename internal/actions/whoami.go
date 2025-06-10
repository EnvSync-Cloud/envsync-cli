package actions

import (
	"context"
	"encoding/json"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v3"
)

func Whoami() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		authService := services.NewAuthService()

		// Get user info from the auth service
		userInfo, err := authService.Whoami()

		if err != nil {
			return err
		}

		// If JSON flag is set, print in JSON format
		if cmd.Bool("json") {
			jsonOutput, err := json.MarshalIndent(userInfo, "", "  ")
			if err != nil {
				return err
			}
			cmd.Writer.Write(append(jsonOutput, '\n'))
			return nil
		}

		// Print user info
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		cmd.Writer.Write([]byte("User ID: " + userInfo.UserId + "\n"))
		cmd.Writer.Write([]byte("Email: " + userInfo.Email + "\n"))
		cmd.Writer.Write([]byte("Organization: " + userInfo.Org + "\n"))
		cmd.Writer.Write([]byte("Role: " + userInfo.Role + "\n"))
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}
