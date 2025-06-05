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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
