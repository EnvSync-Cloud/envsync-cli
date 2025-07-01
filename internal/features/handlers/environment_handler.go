package handlers

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/environment"
	"github.com/urfave/cli/v3"
)

type EnvironmentHandler struct {
	getEnvUseCase    environment.GetEnvUseCase
	switchEnvUseCase environment.SwitchEnvUseCase
}

func NewEnvironmentHandler(
	getEnvUseCase environment.GetEnvUseCase,
	switchEnvUseCase environment.SwitchEnvUseCase,
) *EnvironmentHandler {
	return &EnvironmentHandler{
		getEnvUseCase:    getEnvUseCase,
		switchEnvUseCase: switchEnvUseCase,
	}
}

func (h *EnvironmentHandler) SwitchEnvironment(ctx context.Context, cmd *cli.Command) error {
	env := domain.EnvType{}

	if err := h.switchEnvUseCase.Execute(ctx, env); err != nil {
		return err
	}

	return nil
}
