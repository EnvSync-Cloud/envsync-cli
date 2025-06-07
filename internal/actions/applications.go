package actions

import (
	"encoding/json"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func CreateApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		appService := services.NewAppService()

		res, err := appService.CreateApp(*domain.NewApplication(
			c.String("name"),
			c.String("description"),
			map[string]any{
				"cli_create": true,
			},
		))

		if err != nil {
			return err
		}

		if c.Bool("json") {
			jsonOutput, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(jsonOutput))
			return nil
		}

		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		c.App.Writer.Write([]byte("Application created successfully.\n"))
		c.App.Writer.Write([]byte("ID: " + res.ID + "\n"))
		c.App.Writer.Write([]byte("Name: " + res.Name + "\n"))
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}

func ListApplications() cli.ActionFunc {
	return func(c *cli.Context) error {
		appService := services.NewAppService()

		res, err := appService.GetAllApps()
		if err != nil {
			return err
		}

		if c.Bool("json") {
			jsonOutput, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(jsonOutput))
			return nil
		}

		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		for _, app := range res {
			c.App.Writer.Write([]byte("ID: " + app.ID + "\n"))
			c.App.Writer.Write([]byte("Name: " + app.Name + "\n"))
			c.App.Writer.Write([]byte("Description: " + app.Description + "\n"))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		}

		return nil
	}
}

func DeleteApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		appService := services.NewAppService()

		appID := c.String("id")
		if appID == "" {
			return cli.Exit("Application ID is required", 1)
		}

		app, err := appService.GetAppByID(appID)
		if err != nil {
			return err
		}

		if err := appService.DeleteApp(app); err != nil {
			return err
		}

		if c.Bool("json") {
			jsonOutput, err := json.MarshalIndent(app, "", "  ")
			if err != nil {
				return err
			}
			c.App.Writer.Write([]byte(jsonOutput))
			return nil
		}

		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		c.App.Writer.Write([]byte("Application deleted successfully.\n"))
		c.App.Writer.Write([]byte("ID: " + app.ID + "\n"))
		c.App.Writer.Write([]byte("Name: " + app.Name + "\n"))
		c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))

		return nil
	}
}
