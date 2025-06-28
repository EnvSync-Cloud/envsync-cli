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
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/cli/formatters"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
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

// Container holds all dependencies
type Container struct {
	// Services
	AppService  services.ApplicationService
	AuthService services.AuthService

	// Use Cases
	CreateAppUseCase      appUseCases.CreateAppUseCase
	DeleteAppUseCase      appUseCases.DeleteAppUseCase
	ListAppsUseCase       appUseCases.ListAppsUseCase
	GetAppUseCase         appUseCases.GetAppUseCase
	LoginUseCase          authUseCases.LoginUseCase
	LogoutUseCase         authUseCases.LogoutUseCase
	WhoamiUseCase         authUseCases.WhoamiUseCase
	SetConfigUseCase      configUseCases.SetConfigUseCase
	GetConfigUseCase      configUseCases.GetConfigUseCase
	ValidateConfigUseCase configUseCases.ValidateConfigUseCase
	ResetConfigUseCase    configUseCases.ResetConfigUseCase

	// Formatters
	AppFormatter    *formatters.AppFormatter
	AuthFormatter   *formatters.AuthFormatter
	ConfigFormatter *formatters.ConfigFormatter

	// TUI Factories
	AppFactory *factory.AppFactory

	// Handlers
	AppHandler    *appHandler.Handler
	AuthHandler   *authHandler.Handler
	ConfigHandler *configHandler.Handler
}

// buildDependencyContainer creates and wires all dependencies
func buildDependencyContainer() *Container {
	c := &Container{}

	// Initialize services
	c.AppService = services.NewAppService()
	c.AuthService = services.NewAuthService()

	// Initialize formatters
	c.AppFormatter = formatters.NewAppFormatter()
	c.AuthFormatter = formatters.NewAuthFormatter()
	c.ConfigFormatter = formatters.NewConfigFormatter()

	// Initialize use cases
	c.CreateAppUseCase = appUseCases.NewCreateAppUseCase(c.AppService)
	c.DeleteAppUseCase = appUseCases.NewDeleteAppUseCase(c.AppService)
	c.ListAppsUseCase = appUseCases.NewListAppsUseCase(c.AppService)
	c.GetAppUseCase = appUseCases.NewGetAppUseCase(c.AppService)

	c.LoginUseCase = authUseCases.NewLoginUseCase(c.AuthService)
	c.LogoutUseCase = authUseCases.NewLogoutUseCase(c.AuthService)
	c.WhoamiUseCase = authUseCases.NewWhoamiUseCase(c.AuthService)

	c.SetConfigUseCase = configUseCases.NewSetConfigUseCase()
	c.GetConfigUseCase = configUseCases.NewGetConfigUseCase()
	c.ValidateConfigUseCase = configUseCases.NewValidateConfigUseCase()
	c.ResetConfigUseCase = configUseCases.NewResetConfigUseCase()

	// Initialize TUI factories
	c.AppFactory = factory.NewAppFactory(
		c.CreateAppUseCase,
		c.DeleteAppUseCase,
		c.ListAppsUseCase,
		c.GetAppUseCase,
	)

	// Initialize handlers
	c.AppHandler = appHandler.NewHandler(
		c.CreateAppUseCase,
		c.DeleteAppUseCase,
		c.ListAppsUseCase,
		c.GetAppUseCase,
		c.AppFormatter,
		c.AppFactory,
	)

	c.AuthHandler = authHandler.NewHandler(
		c.LoginUseCase,
		c.LogoutUseCase,
		c.WhoamiUseCase,
		c.AuthFormatter,
	)

	c.ConfigHandler = configHandler.NewHandler(
		c.SetConfigUseCase,
		c.GetConfigUseCase,
		c.ValidateConfigUseCase,
		c.ResetConfigUseCase,
		c.ConfigFormatter,
	)

	return c
}
