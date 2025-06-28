package commands

import (
	"context"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	appCommands "github.com/EnvSync-Cloud/envsync-cli/internal/features/commands/app"
	authCommands "github.com/EnvSync-Cloud/envsync-cli/internal/features/commands/auth"
	configCommands "github.com/EnvSync-Cloud/envsync-cli/internal/features/commands/config"
	appHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/app"
	authHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/auth"
	configHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/config"
	"github.com/EnvSync-Cloud/envsync-cli/internal/logger"
)

type CommandRegistry struct {
	appHandler    *appHandler.Handler
	authHandler   *authHandler.Handler
	configHandler *configHandler.Handler
}

func NewCommandRegistry(
	appHandler *appHandler.Handler,
	authHandler *authHandler.Handler,
	configHandler *configHandler.Handler,
) *CommandRegistry {
	return &CommandRegistry{
		appHandler:    appHandler,
		authHandler:   authHandler,
		configHandler: configHandler,
	}
}

func (r *CommandRegistry) RegisterCLI() *cli.Command {
	return &cli.Command{
		Name:                  "envsync",
		Usage:                 "EnvSync CLI for managing applications and configurations",
		Suggest:               true,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Usage:   "Output in JSON format",
				Aliases: []string{"j"},
				Value:   false,
			},
		},
		Before: r.beforeHook,
		After:  r.afterHook,
		Commands: []*cli.Command{
			appCommands.Commands(r.appHandler),
			authCommands.Commands(r.authHandler),
			configCommands.Commands(r.configHandler),
		},
	}
}

func (r *CommandRegistry) beforeHook(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	l := logger.NewLogger()
	return context.WithValue(ctx, constants.LoggerKey, l), nil
}

func (r *CommandRegistry) afterHook(ctx context.Context, cmd *cli.Command) error {
	if l, ok := ctx.Value(constants.LoggerKey).(*zap.Logger); ok && l != nil {
		l.Sync()
	}
	return nil
}
