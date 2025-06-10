package actions

import (
	"context"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func RunAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		c := strings.Split(cmd.String("command"), " ")

		// Step1: Initialize Sync service
		s := services.NewSyncService()

		// Step2: Check sync config file exists
		if err := s.SyncConfigExist(); err != nil {
			return err
		}

		// Step3: Fetch the remote env
		remoteEnv, err := s.ReadRemoteEnv()
		if err != nil {
			return err
		}

		// Step4: Set env in terminal environment
		for _, env := range remoteEnv {
			if err := os.Setenv(env.Key, env.Value); err != nil {
				return err
			}
		}

		// Step5: Extract redactValues from remoteEnv
		var redactedValues []string
		for _, env := range remoteEnv {
			redactedValues = append(redactedValues, env.Value)
		}

		// Step6: Initialize redactor service and run redactor
		r := services.NewRedactorService(redactedValues)
		_ = r.RunRedactor(c)

		return nil
	}
}
