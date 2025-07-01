package sync

import (
	"context"
	"fmt"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type pushUseCase struct {
	syncService services.SyncService
}

func NewPushUseCase() PushUseCase {
	service := services.NewSyncService()
	return &pushUseCase{
		syncService: service,
	}
}

func (uc *pushUseCase) Execute(ctx context.Context, configPath string) (SyncResponse, error) {
	// Check if the configuration file exists
	if err := uc.checkConfigFileExists(configPath); err != nil {
		return SyncResponse{}, fmt.Errorf("configuration file check failed: %w", err)
	}

	// Read remote environment variables
	remoteEnv, err := uc.syncService.ReadRemoteEnv()
	if err != nil {
		return SyncResponse{}, fmt.Errorf("failed to read remote environment variables: %w", err)
	}

	// Convert remote env variables to map for processing
	remoteEnvMap := make(map[string]string)
	for _, env := range remoteEnv {
		remoteEnvMap[env.Key] = env.Value
	}

	// Read local environment variables from the specified config file
	localEnv, err := uc.syncService.ReadLocalEnv()
	if err != nil {
		return SyncResponse{}, fmt.Errorf("failed to read local environment variables: %w", err)
	}

	// Calculate the differences between remote and local environment variables
	diff, err := uc.calculateEnvDiff(remoteEnvMap, localEnv)
	if err != nil {
		return SyncResponse{}, fmt.Errorf("failed to calculate environment differences: %w", err)
	}

	if len(diff.Added) > 0 || len(diff.Updated) > 0 || len(diff.Deleted) > 0 {
		envSync := &domain.EnvironmentSync{
			ToAdd:    diff.Added,
			ToUpdate: diff.Updated,
		}
		// ToDelete expects a slice of keys
		for _, v := range diff.Deleted {
			envSync.ToDelete = append(envSync.ToDelete, v.Key)
		}
		if err := uc.syncService.WriteRemoteEnv(envSync); err != nil {
			return SyncResponse{}, fmt.Errorf("failed to write remote environment variables: %w", err)
		}
	}

	return diff, nil
}

func (uc *pushUseCase) checkConfigFileExists(configPath string) error {
	// Check if the configuration file exists at the specified path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist at path: %s", configPath)
	}

	return nil
}

func (uc *pushUseCase) calculateEnvDiff(remoteEnv, localEnv map[string]string) (SyncResponse, error) {
	var (
		added     []domain.EnvironmentVariable
		updated   []domain.EnvironmentVariable
		deleted   []domain.EnvironmentVariable
		conflicts []domain.EnvironmentVariable
		warnings  []string
	)

	// Track keys in remote for deletion detection
	remoteKeys := make(map[string]struct{}, len(remoteEnv))
	for k := range remoteEnv {
		remoteKeys[k] = struct{}{}
	}

	// Check for added, updated, and conflicts
	for k, localVal := range localEnv {
		if remoteVal, ok := remoteEnv[k]; !ok {
			added = append(added, domain.EnvironmentVariable{Key: k, Value: localVal})
		} else if remoteVal != localVal {
			updated = append(updated, domain.EnvironmentVariable{Key: k, Value: localVal})
			conflicts = append(conflicts, domain.EnvironmentVariable{Key: k, Value: localVal})
			warnings = append(warnings, "Conflict for key '"+k+"': remote='"+remoteVal+"' local='"+localVal+"'")
		}
		delete(remoteKeys, k)
	}

	// Remaining keys in remoteKeys are deleted in local
	for k := range remoteKeys {
		deleted = append(deleted, domain.EnvironmentVariable{Key: k, Value: remoteEnv[k]})
	}

	return SyncResponse{
		Added:     added,
		Updated:   updated,
		Deleted:   deleted,
		Conflicts: conflicts,
		Warnings:  warnings,
	}, nil
}
