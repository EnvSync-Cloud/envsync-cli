package environment

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type switchEnvUseCase struct {
	envTypeService services.EnvTypeService
	syncService    services.SyncService
	tui            *factory.EnvFactory
}

func NewSwitchEnvUseCase() SwitchEnvUseCase {
	envTypeService := services.NewEnvTypeService()
	syncService := services.NewSyncService()
	tui := factory.NewEnvFactory()

	return &switchEnvUseCase{
		envTypeService: envTypeService,
		syncService:    syncService,
		tui:            tui,
	}
}

func (uc *switchEnvUseCase) Execute(ctx context.Context, envType domain.EnvType) error {
	// Read current sync config
	syncConfig, err := uc.syncService.ReadConfigData()
	if err != nil {
		return NewFileSystemError("failed to read sync config", err)
	}

	// Fetch available environment types for the app
	envs, err := uc.envTypeService.GetEnvTypeByAppID(syncConfig.AppID)
	if err != nil {
		return NewServiceError("failed to fetch environment types", err)
	}
	if len(envs) == 0 {
		return NewNotFoundError("no environment types found for the current app", nil)
	}

	// Select environment type using TUI
	selectedEnv, err := uc.tui.SelectEnvironmentTUI(envs)
	if err != nil {
		if err == tea.ErrProgramKilled {
			return nil // User cancelled the selection
		}
		return NewServiceError("failed to select environment type", err)
	}

	// Update sync config with the selected environment type
	syncConfig.EnvTypeID = selectedEnv.ID
	if err := uc.syncService.WriteConfigData(syncConfig); err != nil {
		return NewFileSystemError("failed to update sync config with selected environment type", err)
	}

	return nil
}
