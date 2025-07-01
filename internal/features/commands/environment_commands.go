package commands

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers"
	"github.com/urfave/cli/v3"
)

func EnvironmentCommands(handlers *handlers.EnvironmentHandler) *cli.Command {
	return &cli.Command{
		Name:    "environment",
		Aliases: []string{"env"},
		Usage:   "Manage environments for applications",
		Commands: []*cli.Command{
			SwitchEnvironmentCommand(handlers),
		},
	}
}

func SwitchEnvironmentCommand(handlers *handlers.EnvironmentHandler) *cli.Command {
	return &cli.Command{
		Name:    "switch",
		Aliases: []string{"sw"},
		Usage:   "Switch to a different environment",
		Action:  handlers.SwitchEnvironment,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "app-id",
				Usage:    "ID of the application to switch environment for",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "env-id",
				Usage:    "ID of the environment to switch to",
				Required: false,
			},
		},
	}
}
