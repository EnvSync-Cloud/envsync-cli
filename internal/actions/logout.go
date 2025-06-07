package actions

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func Logout() cli.ActionFunc {
	return func(c *cli.Context) error {
		authService := services.NewAuthService()

		// Call the logout function in the auth service
		if err := authService.Logout(); err != nil {
			return err
		}

		c.App.Writer.Write([]byte("Successfully logged out.\n"))
		return nil
	}
}
