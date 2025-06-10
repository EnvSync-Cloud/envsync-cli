package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/actions"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/logger"
)

func main() {
	app := &cli.Command{
		Name:                  "envsync",
		Usage:                 "Sync environment variables between local and remote environments",
		Action:                actions.IndexAction(),
		Suggest:               true,
		EnableShellCompletion: true,
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			l := logger.NewLogger()
			c := context.WithValue(ctx, constants.LoggerKey, l)
			return c, nil
		},
		Commands: []*cli.Command{
			{
				Name:     "login",
				Usage:    "Login to EnvSync Cloud",
				Action:   actions.LoginAction(),
				Category: "Auth",
			},
			{
				Name:   "whoami",
				Usage:  "Display current user information",
				Action: actions.Whoami(),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Usage:   "Output user info in JSON format",
						Aliases: []string{"j"},
					},
				},
				Category: "Auth",
			},
			{
				Name:     "logout",
				Usage:    "Logout from EnvSync Cloud",
				Action:   actions.Logout(),
				Category: "Auth",
			},
			{
				Name:   "init",
				Usage:  "Generate a new configuration file",
				Action: actions.InitAction(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "app",
						Usage:    "Name of your application",
						Aliases:  []string{"a"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "env-type",
						Usage:    "Type of your environment",
						Aliases:  []string{"e"},
						Required: true,
					},
				},
				Category: "Config",
			},
			{
				Name:     "push",
				Usage:    "Push environment variables to remote environment",
				Action:   actions.PushAction(),
				Category: "Sync",
			},
			{
				Name:     "pull",
				Usage:    "Pull environment variables from remote environment",
				Action:   actions.PullAction(),
				Category: "Sync",
			},
			{
				Name:  "run",
				Usage: "Run with project command",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "command",
						Usage:    "Execute the command along with envsync",
						Aliases:  []string{"c"},
						Required: true,
					},
				},
				Action: actions.RunAction(),
			},
			{
				Name:  "app",
				Usage: "Interact with your applications.",
				Commands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new application.",
						Action: actions.CreateApplication(),
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output application details in JSON format",
								Aliases: []string{"j"},
								Value:   false,
							},
						},
						Category: "Application",
					},
					{
						Name:   "delete",
						Usage:  "Delete an application.",
						Action: actions.DeleteApplication(),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Usage:    "ID of the application to delete",
								Aliases:  []string{"i"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output deletion confirmation in JSON format",
								Aliases: []string{"j"},
								Value:   false,
							},
						},
						Category: "Application",
					},
					{
						Name:   "list",
						Usage:  "List all applications.",
						Action: actions.ListApplications(),
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output applications in JSON format",
								Aliases: []string{"j"},
								Value:   false,
							},
						},
						Category: "Application",
					},
				},
			},
			{
				Name:  "env-type",
				Usage: "Manage environment types.",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all environment types.",
						Action: actions.ListEnvTypes(),
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output environment types in JSON format",
								Aliases: []string{"j"},
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
