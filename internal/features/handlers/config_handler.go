package handlers

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/formatters"
)

type ConfigHandler struct {
	setUseCase   config.SetConfigUseCase
	getUseCase   config.GetConfigUseCase
	resetUseCase config.ResetConfigUseCase
	formatter    *formatters.ConfigFormatter
}

func NewConfigHandler(
	setUseCase config.SetConfigUseCase,
	getUseCase config.GetConfigUseCase,
	resetUseCase config.ResetConfigUseCase,
	formatter *formatters.ConfigFormatter,
) *ConfigHandler {
	return &ConfigHandler{
		setUseCase:   setUseCase,
		getUseCase:   getUseCase,
		resetUseCase: resetUseCase,
		formatter:    formatter,
	}
}

func (h *ConfigHandler) Set(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()

	if args.Len() < 1 {
		return h.formatter.FormatError(cmd.Writer, "No arguments provided. Usage: envsync config set key=value")
	}

	// Parse key=value pairs from arguments
	keyValuePairs := make(map[string]string)
	for i := 0; i < args.Len(); i++ {
		arg := args.Get(i)
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return h.formatter.FormatError(cmd.Writer, "Invalid format: '"+arg+"'. Expected format: key=value")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return h.formatter.FormatError(cmd.Writer, "Empty key provided in: '"+arg+"'")
		}

		keyValuePairs[key] = value
	}

	// Build request
	req := config.SetConfigRequest{
		KeyValuePairs: keyValuePairs,
		OverwriteAll:  cmd.Bool("overwrite"),
	}

	// Execute use case
	if err := h.setUseCase.Execute(ctx, req); err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	// Format success message
	return h.formatter.FormatSuccess(cmd.Writer, "Configuration updated successfully!")
}

func (h *ConfigHandler) Get(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	keys := make([]string, args.Len())
	for i := 0; i < args.Len(); i++ {
		keys[i] = strings.TrimSpace(args.Get(i))
	}

	// Build request
	req := config.GetConfigRequest{
		Keys: keys,
	}

	// Execute use case
	response, err := h.getUseCase.Execute(ctx, req)
	if err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	// Format output based on requested format
	if cmd.Bool("json") {
		return h.formatter.FormatJSON(cmd.Writer, response.Config)
	}

	// If specific keys were requested, show only those
	if len(keys) > 0 {
		for _, key := range keys {
			if value, exists := response.Values[key]; exists {
				if err := h.formatter.FormatSingleValue(cmd.Writer, key, value); err != nil {
					return err
				}
			} else {
				if err := h.formatter.FormatWarning(cmd.Writer, "Key '"+key+"' not found"); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return nil
}

func (h *ConfigHandler) Reset(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	keys := make([]string, args.Len())
	for i := 0; i < args.Len(); i++ {
		keys[i] = strings.TrimSpace(args.Get(i))
	}

	// Build request
	req := config.ResetConfigRequest{
		Keys: keys,
	}

	// Execute use case
	if err := h.resetUseCase.Execute(ctx, req); err != nil {
		return h.formatUseCaseError(cmd, err)
	}

	// Format success message
	if len(keys) == 0 {
		return h.formatter.FormatSuccess(cmd.Writer, "All configuration values reset successfully!")
	} else {
		return h.formatter.FormatSuccess(cmd.Writer, "Configuration values reset successfully!")
	}
}

// Helper methods

func (h *ConfigHandler) formatUseCaseError(cmd *cli.Command, err error) error {
	// Handle different types of use case errors
	switch e := err.(type) {
	case *config.ConfigError:
		switch e.Code {
		case config.ConfigErrorCodeValidation:
			return h.formatter.FormatError(cmd.Writer, "Validation error: "+e.Message)
		case config.ConfigErrorCodeFileSystem:
			return h.formatter.FormatError(cmd.Writer, "File system error: "+e.Message)
		case config.ConfigErrorCodePermission:
			return h.formatter.FormatError(cmd.Writer, "Permission error: "+e.Message)
		case config.ConfigErrorCodeNotFound:
			return h.formatter.FormatError(cmd.Writer, "Configuration not found: "+e.Message)
		case config.ConfigErrorCodeCorrupted:
			return h.formatter.FormatError(cmd.Writer, "Configuration corrupted: "+e.Message)
		default:
			return h.formatter.FormatError(cmd.Writer, "Service error: "+e.Message)
		}
	default:
		return h.formatter.FormatError(cmd.Writer, "Unexpected error: "+err.Error())
	}
}
