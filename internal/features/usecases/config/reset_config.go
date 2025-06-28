package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
)

type resetConfigUseCase struct{}

func NewResetConfigUseCase() ResetConfigUseCase {
	return &resetConfigUseCase{}
}

func (uc *resetConfigUseCase) Execute(ctx context.Context, req ResetConfigRequest) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return NewValidationError("invalid reset config request", "", err)
	}

	// Read current configuration
	cfg, err := config.ReadConfigFile()
	if err != nil {
		return NewFileSystemError("failed to read config file", err)
	}

	// Create backup before reset (optional but recommended)
	if err := uc.createBackup(cfg); err != nil {
		return NewFileSystemError("failed to create backup", err)
	}

	// Reset configuration based on request
	if len(req.Keys) == 0 {
		// Reset all configuration
		err = uc.resetAllConfig()
	} else {
		// Reset specific keys
		err = uc.resetSpecificKeys(cfg, req.Keys)
	}

	if err != nil {
		return err
	}

	return nil
}

func (uc *resetConfigUseCase) resetAllConfig() error {
	// Create empty configuration
	emptyCfg := config.AppConfig{}

	// Write empty configuration to file
	if err := emptyCfg.WriteConfigFile(); err != nil {
		return NewFileSystemError("failed to write reset config file", err)
	}

	return nil
}

func (uc *resetConfigUseCase) resetSpecificKeys(cfg config.AppConfig, keys []string) error {
	// Reset specific configuration keys
	for _, key := range keys {
		if err := uc.resetConfigKey(&cfg, key); err != nil {
			return NewValidationError("failed to reset config key", key, err)
		}
	}

	// Write updated configuration to file
	if err := cfg.WriteConfigFile(); err != nil {
		return NewFileSystemError("failed to write updated config file", err)
	}

	return nil
}

func (uc *resetConfigUseCase) resetConfigKey(cfg *config.AppConfig, key string) error {
	// Normalize key to lowercase for comparison
	normalizedKey := strings.ToLower(key)

	switch normalizedKey {
	case "access_token", "accesstoken":
		cfg.AccessToken = ""
	case "backend_url", "backendurl":
		cfg.BackendURL = ""
	default:
		return fmt.Errorf("unknown configuration key: '%s'. Valid keys are: access_token, backend_url", key)
	}

	return nil
}

func (uc *resetConfigUseCase) createBackup(cfg config.AppConfig) error {
	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFilename := fmt.Sprintf("envsync-config-backup-%s.json", timestamp)

	// Get config directory (same as main config file)
	configDir, err := uc.getConfigDirectory()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	backupPath := fmt.Sprintf("%s/%s", configDir, backupFilename)

	// Create backup content (JSON format for easier restoration)
	backupContent := fmt.Sprintf(`{
  "access_token": "%s",
  "backend_url": "%s",
  "backup_created_at": "%s",
  "backup_note": "Automatic backup created before config reset"
}`, cfg.AccessToken, cfg.BackendURL, time.Now().Format(time.RFC3339))

	// Write backup file
	if err := os.WriteFile(backupPath, []byte(backupContent), 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

func (uc *resetConfigUseCase) getConfigDirectory() (string, error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// EnvSync config directory
	configDir := fmt.Sprintf("%s/.envsync", homeDir)

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func (uc *resetConfigUseCase) validateResetPermissions() error {
	// Check if we have permission to modify the config file
	cfg, err := config.ReadConfigFile()
	if err != nil {
		// If file doesn't exist, we can create it
		return nil
	}

	// Try to write a test to verify permissions
	if err := cfg.WriteConfigFile(); err != nil {
		return NewPermissionError("insufficient permissions to reset configuration", err)
	}

	return nil
}

func (uc *resetConfigUseCase) cleanupOldBackups() error {
	// Clean up old backup files (keep only last 5 backups)
	configDir, err := uc.getConfigDirectory()
	if err != nil {
		return err
	}

	// Read directory contents
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	// Filter backup files
	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "envsync-config-backup-") {
			backupFiles = append(backupFiles, entry)
		}
	}

	// If we have more than 5 backup files, remove the oldest ones
	if len(backupFiles) > 5 {
		// Sort by modification time (oldest first)
		// Note: This is a simplified implementation
		// In a real implementation, you'd want to sort by file modification time

		filesToRemove := len(backupFiles) - 5
		for i := 0; i < filesToRemove; i++ {
			backupPath := fmt.Sprintf("%s/%s", configDir, backupFiles[i].Name())
			if err := os.Remove(backupPath); err != nil {
				// Log warning but don't fail the operation
				// In a real implementation, you'd use proper logging
				continue
			}
		}
	}

	return nil
}

func (uc *resetConfigUseCase) generateResetSummary(keys []string, isFullReset bool) string {
	if isFullReset {
		return "All configuration values have been reset"
	}

	if len(keys) == 1 {
		return fmt.Sprintf("Configuration key '%s' has been reset", keys[0])
	}

	return fmt.Sprintf("Configuration keys [%s] have been reset", strings.Join(keys, ", "))
}

func (uc *resetConfigUseCase) validateResetSafety(cfg config.AppConfig, keys []string) error {
	// Perform safety checks before reset

	// Check if user is about to reset all config while having important data
	if len(keys) == 0 && (cfg.AccessToken != "" || cfg.BackendURL != "") {
		// This would typically prompt for confirmation in an interactive context
		// For now, we'll just warn through error message
		return NewValidationError("attempting to reset all configuration with existing data", "",
			fmt.Errorf("use --force flag to confirm full configuration reset"))
	}

	return nil
}
