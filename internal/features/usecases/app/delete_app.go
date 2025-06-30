package app

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type deleteAppUseCase struct {
	appService services.ApplicationService
	tui        *factory.AppFactory
}

func NewDeleteAppUseCase() DeleteAppUseCase {
	service := services.NewAppService()
	tui := factory.NewAppFactory()
	return &deleteAppUseCase{
		appService: service,
		tui:        tui,
	}
}

func (uc *deleteAppUseCase) Execute(ctx context.Context) error {
	// Find the application to delete
	apps, err := uc.findAllApplications()
	if err != nil {
		return err
	}

	deleteApps, _ := uc.tui.DeleteAppsTUI(apps)

	for _, app := range deleteApps {
		_ = uc.appService.DeleteApp(app)
	}

	return nil
}

func (uc *deleteAppUseCase) findAllApplications() ([]domain.Application, error) {
	// Retrieve all applications from the service
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return nil, NewServiceError("failed to retrieve applications", err)
	}
	return apps, nil
}
