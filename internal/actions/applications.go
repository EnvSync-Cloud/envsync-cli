package actions

import "github.com/urfave/cli/v2"

func CreateApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		return nil
	}
}

func ListApplications() cli.ActionFunc {
	return func(c *cli.Context) error {
		return nil
	}
}

func DeleteApplication() cli.ActionFunc {
	return func(c *cli.Context) error {
		return nil
	}
}
