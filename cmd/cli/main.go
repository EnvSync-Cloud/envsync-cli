package main

import (
	"log"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/actions"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "envsync-cli",
		Usage:  "Sync environment variables between local and remote environments",
		Action: actions.IndexAction(),
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "Login to EnvSync Cloud",
				Action: actions.LoginAction(),
			},
			{
				Name:   "gen-config",
				Usage:  "Generate a new configuration file",
				Action: actions.GenConfigAction(),
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
			},
			{
				Name:   "push",
				Usage:  "Push environment variables to remote environment",
				Action: actions.PushAction(),
			},
			{
				Name:   "pull",
				Usage:  "Pull environment variables from remote environment",
				Action: actions.PullAction(),
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
