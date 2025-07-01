package handlers

import (
	"context"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/run"
	"github.com/urfave/cli/v3"
)

type RunHandler struct {
	redactUseCase    run.RedactUseCase
	injectEnvUseCase run.InjectEnvUseCase
}

func NewRunHandler(ruc run.RedactUseCase, iuc run.InjectEnvUseCase) *RunHandler {
	return &RunHandler{
		redactUseCase:    ruc,
		injectEnvUseCase: iuc,
	}
}

func (h *RunHandler) Run(ctx context.Context, cmd *cli.Command) error {
	c := strings.Split(cmd.String("command"), " ")

	envs, err := h.injectEnvUseCase.Execute(ctx)
	if err != nil {
		return err
	}

	var redactedValues []string
	for _, env := range envs {
		redactedValues = append(redactedValues, env)
	}

	_ = h.redactUseCase.Execute(ctx, c, redactedValues)

	return nil
}
