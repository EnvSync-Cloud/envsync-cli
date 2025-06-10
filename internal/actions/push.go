package actions

import (
	"context"
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v3"
)

func PushAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// Step1: Initialize sync service
		syncService := services.NewSyncService()

		// Step2: Check if the sync config available
		// If not found thorw error
		if err := syncService.CheckSyncConfig(); err != nil {
			return err
		}

		// Step3: Read the config file and get the data
		cfg, err := syncService.ReadConfigData()
		if err != nil {
			return err
		}

		// Step4: Get remote env
		remoteEnvs, err := syncService.GetAllEnv(cfg.AppID, cfg.EnvTypeID)
		if err != nil {
			return err
		}

		// Convert remote env variables to map for processing
		remoteEnvMap := make(map[string]string)
		for _, env := range remoteEnvs {
			remoteEnvMap[env.Key] = env.Value
		}

		// Step5: Get local env
		localEnvs, err := syncService.ReadLocalEnv()
		if err != nil {
			return err
		}

		// Step5: Calculate env diff
		envDiff := syncService.CalculateEnvDiff(localEnvs, remoteEnvMap)

		if envDiff.HasChanges() {
			if err := syncService.PushEnv(envDiff); err != nil {
				return err
			}

			summary := envDiff.GetSummary()
			cmd.Writer.Write([]byte("\n🎉 Environment variables synced successfully!\n"))
			cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
			cmd.Writer.Write([]byte(fmt.Sprintf("✅ Added:   %d variables\n", summary.AddCount)))
			cmd.Writer.Write([]byte(fmt.Sprintf("🔄 Updated: %d variables\n", summary.UpdateCount)))
			cmd.Writer.Write([]byte(fmt.Sprintf("🗑️  Deleted: %d variables\n", summary.DeleteCount)))
			cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"))
		} else {
			cmd.Writer.Write([]byte("\n✨ No changes detected. Environment is already in sync.\n\n"))
		}

		return nil
	}
}
