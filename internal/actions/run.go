package actions

import (
	"os"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func RunAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		cmd := strings.Split(c.String("command"), " ")

		// Step1: Initialize Sync service
		s := services.NewSyncService()

		// Step2: Check sync config file exists
		if err := s.CheckSyncConfig(); err != nil {
			return err
		}

		// Step3: Read sync config file
		projCfg, err := s.ReadConfigData()
		if err != nil {
			return err
		}

		// Step4: Fetch the remote env
		remoteEnv, err := s.GetAllEnv(projCfg.AppID, projCfg.EnvTypeID)
		if err != nil {
			return err
		}

		// Step5: Set env in terminal environment
		for _, env := range remoteEnv {
			if err := os.Setenv(env.Key, env.Value); err != nil {
				return err
			}
		}

		// Step6: Extract redactValues from remoteEnv
		var redactedValues []string
		for _, env := range remoteEnv {
			redactedValues = append(redactedValues, env.Value)
		}

		// Step6: Initialize Redactor service and run redactor
		r := services.NewRedactorService(redactedValues)
		_ = r.RunRedactor(cmd)

		return nil
	}
}
