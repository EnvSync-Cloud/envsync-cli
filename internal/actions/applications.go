package actions

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func CreateApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		return nil
	}
}

func ListApplications() cli.ActionFunc {
	return func(c *cli.Context) error {
		// Step1: Initialize application service
		as := services.NewAppService()

		// Step2: Get all applications
		apps, err := as.GetAllApps()
		if err != nil {
			return err
		}

		// Step3: Print applications
		c.App.Writer.Write([]byte("🚀 Available Applications:\n"))
		for _, app := range apps {
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
			c.App.Writer.Write([]byte(fmt.Sprintf("📛 Name: %s\n", app.Name)))
			c.App.Writer.Write([]byte(fmt.Sprintf("🆔 ID: %s\n", app.ID)))
			c.App.Writer.Write([]byte(fmt.Sprintf("📝 Description: %s\n", app.Description)))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}

func DeleteApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		return nil
	}
}
