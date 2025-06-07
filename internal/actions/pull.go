package actions

import (
	"fmt"

	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/urfave/cli/v2"
)

func PullAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		// Step 1: Initiate sync service
		syncService := services.NewSyncService()

		// Step2:  Check if sync config available.
		// If not found throw error.
		if err := syncService.CheckSyncConfig(); err != nil {
			return err
		}

		// Step3: Read the config file and get the data
		cfg, err := syncService.ReadConfigData()
		if err != nil {
			return err
		}

		// Step4: Fetch env from remote
		remoteEnvs, err := syncService.GetAllEnv(cfg.AppID, cfg.EnvTypeID)
		if err != nil {
			return err
		}

		// Convert remote env variables to map for processing
		remoteEnvMap := make(map[string]string)
		for _, env := range remoteEnvs {
			remoteEnvMap[env.Key] = env.Value
		}

		// Step5: Calculate the diff from local env
		localEnvs, err := syncService.ReadLocalEnv()
		if err != nil {
			return err
		}

		envDiff := syncService.CalculateEnvDiff(localEnvs, remoteEnvMap)

		// Step6: Write to local env
		if envDiff.HasChanges() {
			// Write remote environment variables to local .env file
			if err := syncService.WriteLocalEnv(remoteEnvMap); err != nil {
				return err
			}

			summary := envDiff.GetSummary()
			c.App.Writer.Write([]byte("\n🎉 Environment variables synced successfully!\n"))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
			c.App.Writer.Write([]byte(fmt.Sprintf("✅ Added:   %d variables\n", summary.AddCount)))
			c.App.Writer.Write([]byte(fmt.Sprintf("🔄 Updated: %d variables\n", summary.UpdateCount)))
			c.App.Writer.Write([]byte(fmt.Sprintf("🗑️  Deleted: %d variables\n", summary.DeleteCount)))
			c.App.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"))
		} else {
			c.App.Writer.Write([]byte("\n✨ No changes detected. Environment is already in sync.\n\n"))
		}

		return nil
	}
}
