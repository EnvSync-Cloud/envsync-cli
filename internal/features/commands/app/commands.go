package app

import (
	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/app"
)

// Commands returns all app-related commands
func Commands(handler *app.Handler) *cli.Command {
	return &cli.Command{
		Name:  "app",
		Usage: "Interact with your applications",
		Commands: []*cli.Command{
			CreateCommand(handler),
			DeleteCommand(handler),
			ListCommand(handler),
			SelectCommand(handler),
		},
	}
}

func CreateCommand(handler *app.Handler) *cli.Command {
	return &cli.Command{
		Name:   "create",
		Usage:  "Create a new application",
		Action: handler.Create,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Application name",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "description",
				Usage:   "Application description",
				Aliases: []string{"d"},
			},
			&cli.StringSliceFlag{
				Name:    "metadata",
				Usage:   "Application metadata in key=value format",
				Aliases: []string{"m"},
			},
		},
	}
}

func DeleteCommand(handler *app.Handler) *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "Delete an application",
		Action: handler.Delete,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "id",
				Usage:   "Application ID to delete",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Application name to delete",
				Aliases: []string{"n"},
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Force deletion without confirmation",
				Value: false,
			},
		},
	}
}

func ListCommand(handler *app.Handler) *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all applications",
		Action: handler.List,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "interactive",
				Usage: "Use interactive mode (default: true)",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "no-interactive",
				Usage: "Disable interactive mode",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format (table, compact, list)",
				Value: "table",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Maximum number of applications to display",
				Value: 0, // 0 means no limit
			},
			&cli.IntFlag{
				Name:  "offset",
				Usage: "Number of applications to skip",
				Value: 0,
			},
			&cli.StringFlag{
				Name:    "org-id",
				Usage:   "Filter by organization ID",
				Aliases: []string{"org"},
			},
		},
	}
}

func SelectCommand(handler *app.Handler) *cli.Command {
	return &cli.Command{
		Name:   "select",
		Usage:  "Interactively select an application",
		Action: handler.Select,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "set-default",
				Usage: "Set selected application as default",
				Value: false,
			},
		},
	}
}
