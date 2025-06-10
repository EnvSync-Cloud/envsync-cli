package actions

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

func PullAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// Step 1: Initiate sync service
		syncService := services.NewSyncService()

		// Step2: Check if sync config available.
		// If not found throw error.
		if err := syncService.SyncConfigExist(); err != nil {
			return err
		}

		// Step3: Fetch env from remote
		remoteEnvs, err := syncService.ReadRemoteEnv()
		if err != nil {
			return err
		}

		// Convert remote env variables to map for processing
		remoteEnvMap := make(map[string]string)
		for _, env := range remoteEnvs {
			remoteEnvMap[env.Key] = env.Value
		}

		// Step4: Calculate the diff from local env
		localEnvs, err := syncService.ReadLocalEnv()
		if err != nil {
			return err
		}

		envDiff := syncService.CalculateEnvDiff(localEnvs, remoteEnvMap)

		// Step5: Write to local env
		if envDiff.HasChanges() {
			// Write remote environment variables to local .env file
			if err := syncService.WriteLocalEnv(remoteEnvMap); err != nil {
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
