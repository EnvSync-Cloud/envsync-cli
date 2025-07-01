package commands

import (
	"context"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/EnvSync-Cloud/envsync-cli/internal/constants"
	appHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers"
	"github.com/EnvSync-Cloud/envsync-cli/internal/logger"
)

// ExecutionMode represents how the command should be executed
type ExecutionMode int

const (
	ExecutionModeJSON ExecutionMode = iota
	ExecutionModeCLI
)

type CommandRegistry struct {
	appHandler         *appHandler.AppHandler
	authHandler        *appHandler.AuthHandler
	configHandler      *appHandler.ConfigHandler
	environmentHandler *appHandler.EnvironmentHandler
	syncHandler        *appHandler.SyncHandler
}

func NewCommandRegistry(
	appHandler *appHandler.AppHandler,
	authHandler *appHandler.AuthHandler,
	configHandler *appHandler.ConfigHandler,
	environmentHandler *appHandler.EnvironmentHandler,
	syncHandler *appHandler.SyncHandler,
) *CommandRegistry {
	return &CommandRegistry{
		appHandler:         appHandler,
		authHandler:        authHandler,
		configHandler:      configHandler,
		environmentHandler: environmentHandler,
		syncHandler:        syncHandler,
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
		Action: RootCommand(),
		Commands: []*cli.Command{
			AppCommands(r.appHandler),
			AuthCommands(r.authHandler),
			ConfigCommands(r.configHandler),
			EnvironmentCommands(r.environmentHandler),
			PullCommand(r.syncHandler),
			PushCommand(r.syncHandler),
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

func RootCommand() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		cmd.Writer.Write([]byte("Welcome to EnvSync CLI!\n"))
		cmd.Writer.Write([]byte("Use 'envsync --help' to see available commands.\n"))
		cmd.Writer.Write([]byte("For more information, visit: https://envsync.cloud/docs\n"))
		cmd.Writer.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
		return nil
	}
}
