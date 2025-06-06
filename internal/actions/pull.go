package actions

import (
	"github.com/EnvSync-Cloud/envsync-cli/internal/helper"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func PullAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		if err := helper.CheckProjectConfig(); err != nil {
			return err
		}

		syncService := services.NewSyncService()
		projectCfgService := services.NewProjectConfigService()

		// Read project configuration
		projectCfg, err := projectCfgService.ReadProjectConfig()
		if err != nil {
			return err
		}

		// Fetch environment variables from cloud
		env, err := syncService.PullEnv(projectCfg.AppID, projectCfg.EnvType)
		if err != nil {
			return err
		}

		if err := helper.WriteEnv(env); err != nil {
			return err
		}

		return nil
	}
}
