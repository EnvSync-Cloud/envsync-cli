package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type getAppUseCase struct {
	appService services.ApplicationService
}

func NewGetAppUseCase(appService services.ApplicationService) GetAppUseCase {
	return &getAppUseCase{
		appService: appService,
	}
}

func (uc *getAppUseCase) Execute(ctx context.Context, req GetAppRequest) (*domain.Application, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, NewValidationError("invalid get app request", err)
	}

	// Find the application
	app, err := uc.findApplication(req)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (uc *getAppUseCase) findApplication(req GetAppRequest) (*domain.Application, error) {
	// If ID is provided, try to get by ID first
	if req.ID != "" {
		app, err := uc.appService.GetAppByID(req.ID)
		if err != nil {
			// If not found by ID, fall back to searching by name if name is also provided
			if req.Name != "" {
				return uc.findApplicationByName(req.Name)
			}
			return nil, NewServiceError("failed to get application by ID", err)
		}
		return &app, nil
	}

	// If only name is provided, search by name
	if req.Name != "" {
		return uc.findApplicationByName(req.Name)
	}

	// This should not happen due to validation, but just in case
	return nil, NewValidationError("no identifier provided", ErrAppIdentifierRequired)
}

func (uc *getAppUseCase) findApplicationByName(name string) (*domain.Application, error) {
	// Get all applications and find by name
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return nil, NewServiceError("failed to retrieve applications", err)
	}

	// Find by name (case-insensitive)
	for _, app := range apps {
		if strings.EqualFold(app.Name, name) {
			return &app, nil
		}
	}

	// Application not found
	return nil, NewNotFoundError(
		fmt.Sprintf("application '%s' not found", name),
		ErrAppNotFound,
	)
}
