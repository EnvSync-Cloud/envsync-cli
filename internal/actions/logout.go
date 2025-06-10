package actions

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v3"
)

func Logout() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		authService := services.NewAuthService()

		// Call the logout function in the auth service
		if err := authService.Logout(); err != nil {
			return err
		}

		cmd.Writer.Write([]byte("Successfully logged out.\n"))
		return nil
	}
}
