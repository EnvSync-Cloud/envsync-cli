package app

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type listAppsUseCase struct {
	appService services.ApplicationService
	tui        *factory.AppFactory
}

func NewListAppsUseCase() ListAppsUseCase {
	tui := factory.NewAppFactory()
	service := services.NewAppService()
	return &listAppsUseCase{
		appService: service,
		tui:        tui,
	}
}

func (uc *listAppsUseCase) Execute(ctx context.Context) error {
	// Get applications from service
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return NewServiceError("failed to retrieve applications", err)
	}

	_ = uc.tui.ListAppsInteractive(apps)

	return nil
}
