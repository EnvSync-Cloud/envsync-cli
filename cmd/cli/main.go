package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/EnvSync-Cloud/envsync-cli/internal/actions"
	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	"github.com/EnvSync-Cloud/envsync-cli/internal/logger"
)

func main() {
	app := &cli.Command{
		Name:                  "envsync",
		Usage:                 "sync environment variables between local and remote environments",
		Action:                actions.IndexAction(),
		Suggest:               true,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "Output in JSON format",
				Aliases:     []string{"j"},
				Value:       false,
				DefaultText: "false",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			l := logger.NewLogger()
			c := context.WithValue(ctx, constants.LoggerKey, l)
			return c, nil
		},
		After: func(ctx context.Context, cmd *cli.Command) error {
			if l, ok := ctx.Value(constants.LoggerKey).(*zap.Logger); ok && l != nil {
				l.Sync()
			}
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:     "login",
				Usage:    "Login to Envsync Cloud",
				Action:   actions.LoginAction(),
				Category: "AUTH",
			},
			{
				Name:     "whoami",
				Usage:    "Display current user information",
				Action:   actions.Whoami(),
				Category: "AUTH",
			},
			{
				Name:     "logout",
				Usage:    "Logout from Envsync Cloud",
				Action:   actions.Logout(),
				Category: "AUTH",
			},
			{
				Name:   "init",
				Usage:  "Generate a new configuration file",
				Action: actions.InitAction(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "app",
						Usage:    "Name of your app",
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
			},
			{
				Name:     "push",
				Usage:    "Push environment variables to remote environment",
				Action:   actions.PushAction(),
				Category: "SYNC",
			},
			{
				Name:     "pull",
				Usage:    "Pull environment variables from remote environment",
				Action:   actions.PullAction(),
				Category: "SYNC",
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
					&cli.StringSliceFlag{
						Name:    "redact",
						Usage:   "Values to redact from output (for testing)",
						Aliases: []string{"r"},
					},
				},
				Action: actions.RunAction(),
			},
			{
				Name:  "app",
				Usage: "Interact with your apps.",
				Commands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new app.",
						Action: actions.CreateApplication(),
					},
					{
						Name:   "delete",
						Usage:  "Delete an app.",
						Action: actions.DeleteApplication(),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Usage:    "ID of the app to delete",
								Aliases:  []string{"i"},
								Required: true,
							},
						},
					},
					{
						Name:   "list",
						Usage:  "List all apps.",
						Action: actions.ListApplications(),
					},
				},
			},
			{
				Name:  "env-type",
				Usage: "Manage environment types.",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all environment types for a specific app.",
						Action: actions.GetEnvTypesByApp(),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "app-id",
								Usage:    "ID of the app to list environment types for",
								Aliases:  []string{"a"},
								Required: true,
							},
						},
					},
					{
						Name:   "view",
						Usage:  "View details of a specific environment type.",
						Action: actions.GetEnvTypeByID(),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Usage:    "ID of the environment type to view",
								Aliases:  []string{"i"},
								Required: true,
							},
						},
					},
					{
						Name:   "switch",
						Usage:  "Switch to a different environment type for the current app.",
						Action: actions.SwitchEnvType(),
					},
				},
			},
			{
				Name:  "config",
				Usage: "Manage configuration settings.",
				Commands: []*cli.Command{
					{
						Name:   "set",
						Usage:  "Set a configuration value.",
						Action: actions.SetConfigAction(),
					},
					{
						Name:   "get",
						Usage:  "Get a configuration value.",
						Action: actions.GetConfigAction(),
					},
				},
			},
			{
				Name:  "users",
				Usage: "Manage users.",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all active users.",
						Action: actions.ListUsers(),
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output users in JSON format",
								Aliases: []string{"j"},
								Value:   false,
							},
						},
					},
				},
			},
			{
				Name:   "roles",
				Usage:  "Manage roles.",
				Action: actions.ListRoles(),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Usage:   "Output roles in JSON format",
						Aliases: []string{"j"},
						Value:   false,
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
