package actions

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/helper"
	"github.com/urfave/cli/v2"
)

func PushAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := helper.CheckProjectConfig(); err != nil {
			return err
		}

		return nil
	}
}
