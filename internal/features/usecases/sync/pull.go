package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
	"github.com/joho/godotenv"
)

type pullUseCase struct {
	syncService services.SyncService
}

func NewPullUseCase() PullUseCase {
	service := services.NewSyncService()
	return &pullUseCase{
		syncService: service,
	}
}

func (uc *pullUseCase) Execute(ctx context.Context, configPath string) (SyncResponse, error) {
	// Check if the configuration file exists
	fmt.Println("checking config existance")
	if err := uc.checkConfigFileExists(configPath); err != nil {
		return SyncResponse{}, fmt.Errorf("configuration file check failed: %w", err)
	}
	fmt.Println("completed checking config existance")

	fmt.Println("reading remote env")
	// Read remote remote environment variables
	remoteEnv, err := uc.syncService.ReadRemoteEnv()
	if err != nil {
		return SyncResponse{}, fmt.Errorf("failed to read remote environment variables: %w", err)
	}
	fmt.Println("reading remote env completed")
	jsonData, _ := json.Marshal(remoteEnv)
	fmt.Println("Remote Environment Variables:")
	fmt.Println(string(jsonData))

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
		// print remote environment variables to the console
		jsonData, _ := json.Marshal(remoteEnv)
		fmt.Println("Remote Environment Variables:")
		fmt.Println(string(jsonData))

		err := uc.writeToLocalEnv(remoteEnvMap)
		if err != nil {
			return SyncResponse{}, fmt.Errorf("failed to write updated environment variables to local file: %w", err)
		}
		fmt.Println("Printed to env")
	}

	return diff, nil
}

func (uc *pullUseCase) checkConfigFileExists(configPath string) error {
	// Check if the configuration file exists at the specified path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist at path: %s", configPath)
	}

	return nil
}

func (uc *pullUseCase) calculateEnvDiff(remoteEnv, localEnv map[string]string) (SyncResponse, error) {
	added := make([]domain.EnvironmentVariable, 0)
	updated := make([]domain.EnvironmentVariable, 0)
	deleted := make([]domain.EnvironmentVariable, 0)
	conflicts := make([]domain.EnvironmentVariable, 0)
	warnings := make([]string, 0)

	// Track keys in local for deletion detection
	localKeys := make(map[string]struct{}, len(localEnv))
	for k := range localEnv {
		localKeys[k] = struct{}{}
	}

	// Check for added and updated/conflicted keys (from remote to local)
	for k, remoteVal := range remoteEnv {
		if localVal, ok := localEnv[k]; !ok {
			added = append(added, domain.EnvironmentVariable{Key: k, Value: remoteVal})
		} else if localVal != remoteVal {
			updated = append(updated, domain.EnvironmentVariable{Key: k, Value: remoteVal})
			conflicts = append(conflicts, domain.EnvironmentVariable{Key: k, Value: remoteVal})
			warnings = append(warnings, "Conflict for key '"+k+"': local='"+localVal+"' remote='"+remoteVal+"'")
		}
		delete(localKeys, k)
	}

	// Remaining keys in localKeys are deleted in remote (should be deleted from local)
	for k := range localKeys {
		deleted = append(deleted, domain.EnvironmentVariable{Key: k, Value: localEnv[k]})
	}

	return SyncResponse{
		Added:     added,
		Updated:   updated,
		Deleted:   deleted,
		Conflicts: conflicts,
		Warnings:  warnings,
	}, nil
}

func (uc *pullUseCase) writeToLocalEnv(env map[string]string) error {
	return godotenv.Write(env, ".env")
}
