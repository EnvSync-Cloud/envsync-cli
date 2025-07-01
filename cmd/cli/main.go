package main

import (
	"context"
	"log"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/commands"
	appHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers"
	appUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	authUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/auth"
	configUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/config"
	envUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/environment"
	syncUseCase "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/sync"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/formatters"
)

func main() {
	// Initialize dependencies
	container := buildDependencyContainer()

	// Build command registry with dependencies
	registry := commands.NewCommandRegistry(
		container.AppHandler,
		container.AuthHandler,
		container.ConfigHandler,
		container.EnvironmentHandler,
		container.SyncHandler,
	)

	// Build CLI app
	app := registry.RegisterCLI()

	// Run the application
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// Container holds the handler dependencies
type Container struct {
	AppHandler         *appHandler.AppHandler
	AuthHandler        *appHandler.AuthHandler
	ConfigHandler      *appHandler.ConfigHandler
	EnvironmentHandler *appHandler.EnvironmentHandler
	SyncHandler        *appHandler.SyncHandler
}

// buildDependencyContainer creates and wires all handler dependencies
func buildDependencyContainer() *Container {
	c := &Container{}

	// Initialize formatters
	appFormatter := formatters.NewAppFormatter()
	authFormatter := formatters.NewAuthFormatter()
	configFormatter := formatters.NewConfigFormatter()

	// Initialize use cases
	createAppUseCase := appUseCases.NewCreateAppUseCase()
	deleteAppUseCase := appUseCases.NewDeleteAppUseCase()
	listAppsUseCase := appUseCases.NewListAppsUseCase()

	loginUseCase := authUseCases.NewLoginUseCase()
	logoutUseCase := authUseCases.NewLogoutUseCase()
	whoamiUseCase := authUseCases.NewWhoamiUseCase()

	setConfigUseCase := configUseCases.NewSetConfigUseCase()
	getConfigUseCase := configUseCases.NewGetConfigUseCase()
	resetConfigUseCase := configUseCases.NewResetConfigUseCase()

	getEnvironmentUseCase := envUseCases.NewGetEnvUseCase()
	switchEnvironmentUseCase := envUseCases.NewSwitchEnvUseCase()

	pullUseCase := syncUseCase.NewPullUseCase()
	pushUseCase := syncUseCase.NewPushUseCase()

	// Initialize handlers
	c.AppHandler = appHandler.NewAppHandler(
		createAppUseCase,
		deleteAppUseCase,
		listAppsUseCase,
		appFormatter,
	)

	c.AuthHandler = appHandler.NewAuthHandler(
		loginUseCase,
		logoutUseCase,
		whoamiUseCase,
		authFormatter,
	)

	c.ConfigHandler = appHandler.NewConfigHandler(
		setConfigUseCase,
		getConfigUseCase,
		resetConfigUseCase,
		configFormatter,
	)

	c.EnvironmentHandler = appHandler.NewEnvironmentHandler(
		getEnvironmentUseCase,
		switchEnvironmentUseCase,
	)

	c.SyncHandler = appHandler.NewSyncHandler(
		pullUseCase,
		pushUseCase,
	)

	return c
}
