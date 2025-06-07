package main

import (
	"log"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/actions"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "envsync-cli",
		Usage:                "Sync environment variables between local and remote environments",
		Action:               actions.IndexAction(),
		Suggest:              true,
		EnableBashCompletion: true,
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
						Name:     "app-id",
						Usage:    "App ID of your application",
						Aliases:  []string{"a"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "env-type-id",
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
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new application.",
						Action: actions.CreateApplication(),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Name of the application",
								Aliases:  []string{"n"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Usage:    "Description of the application",
								Aliases:  []string{"d"},
								Required: false,
							},
							&cli.BoolFlag{
								Name:    "json",
								Usage:   "Output application details in JSON format",
								Aliases: []string{"j"},
								Value:   false,
							},
						},
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
					},
				},
			},
			{
				Name:  "env-type",
				Usage: "Manage environment types.",
				Subcommands: []*cli.Command{
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
