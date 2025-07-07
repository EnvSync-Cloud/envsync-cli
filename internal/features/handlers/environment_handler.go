package handlers

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/environment"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/formatters"
)

type EnvironmentHandler struct {
	getEnvUseCase    environment.GetEnvUseCase
	switchEnvUseCase environment.SwitchEnvUseCase
	formatter        *formatters.EnvFormatter
}

func NewEnvironmentHandler(
	getEnvUseCase environment.GetEnvUseCase,
	switchEnvUseCase environment.SwitchEnvUseCase,
	formatter *formatters.EnvFormatter,
) *EnvironmentHandler {
	return &EnvironmentHandler{
		getEnvUseCase:    getEnvUseCase,
		switchEnvUseCase: switchEnvUseCase,
		formatter:        formatter,
	}
}

func (h *EnvironmentHandler) SwitchEnvironment(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("json") && (!cmd.IsSet("app-id") && cmd.IsSet("env-id")) {
		return h.formatUseCaseError(cmd, errors.New("app-id or env-id must be provided with json flag"))
	}

	env := domain.EnvType{
		AppID: cmd.String("app-id"),
		ID:    cmd.String("env-id"),
	}

	if err := h.switchEnvUseCase.Execute(ctx, env); !errors.Is(err, tea.ErrProgramKilled) && err != nil {
		return err
	}

	return nil
}

func (h *EnvironmentHandler) formatUseCaseError(cmd *cli.Command, err error) error {
	if cmd.Bool("json") {
		// If JSON output is requested, format the error as JSON
		jsonOutput := map[string]any{
			"error": err.Error(),
		}
		return h.formatter.FormatJSON(cmd.Writer, jsonOutput)
	}

	// Handle different types of use case errors
	switch e := err.(type) {
	case *environment.EnvError:
		switch e.Code {
		case environment.EnvErrorCodeValidation:
			return h.formatter.FormatError(cmd.Writer, "Validation error: "+e.Message)
		case environment.EnvErrorCodeServiceError:
			return h.formatter.FormatError(cmd.Writer, "Service error: "+e.Message)
		case environment.EnvErrorCodeNotFound:
			return h.formatter.FormatError(cmd.Writer, "Environment not found: "+e.Message)
		case environment.EnvErrorCodeCorrupted:
			return h.formatter.FormatError(cmd.Writer, "Environment data is corrupted: "+e.Message)
		case environment.EnvErrorCodePermission:
			return h.formatter.FormatError(cmd.Writer, "Permission error: "+e.Message)
		case environment.EnvErrorCodeFileSystem:
			return h.formatter.FormatError(cmd.Writer, "File system error: "+e.Message)
		default:
			return h.formatter.FormatError(cmd.Writer, "Service error: "+e.Message)
		}
	default:
		return h.formatter.FormatError(cmd.Writer, "Unexpected error: "+err.Error())
	}
}
