package main

import (
	"context"
	"log"
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/features/commands"
	appHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/app"
	authHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/auth"
	configHandler "github.com/EnvSync-Cloud/envsync-cli/internal/features/handlers/config"
	appUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	authUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/auth"
	configUseCases "github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/config"
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
	AppHandler    *appHandler.Handler
	AuthHandler   *authHandler.Handler
	ConfigHandler *configHandler.Handler
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
	validateConfigUseCase := configUseCases.NewValidateConfigUseCase()
	resetConfigUseCase := configUseCases.NewResetConfigUseCase()

	// Initialize handlers
	c.AppHandler = appHandler.NewHandler(
		createAppUseCase,
		deleteAppUseCase,
		listAppsUseCase,
		appFormatter,
	)

	c.AuthHandler = authHandler.NewHandler(
		loginUseCase,
		logoutUseCase,
		whoamiUseCase,
		authFormatter,
	)

	c.ConfigHandler = configHandler.NewHandler(
		setConfigUseCase,
		getConfigUseCase,
		validateConfigUseCase,
		resetConfigUseCase,
		configFormatter,
	)

	return c
}
