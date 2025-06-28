package app

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type listAppsUseCase struct {
	appService services.ApplicationService
}

func NewListAppsUseCase(appService services.ApplicationService) ListAppsUseCase {
	return &listAppsUseCase{
		appService: appService,
	}
}

func (uc *listAppsUseCase) Execute(ctx context.Context, req ListAppsRequest) ([]domain.Application, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, NewValidationError("invalid list apps request", err)
	}

	// Get applications from service
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return nil, NewServiceError("failed to retrieve applications", err)
	}

	// Apply business logic transformations
	filteredApps := uc.applyFilters(apps, req)
	sortedApps := uc.applySorting(filteredApps)
	paginatedApps := uc.applyPagination(sortedApps, req)

	return paginatedApps, nil
}

func (uc *listAppsUseCase) applyFilters(apps []domain.Application, req ListAppsRequest) []domain.Application {
	if req.OrgID == "" {
		return apps
	}

	var filtered []domain.Application
	for _, app := range apps {
		if app.OrgID == req.OrgID {
			filtered = append(filtered, app)
		}
	}

	return filtered
}

func (uc *listAppsUseCase) applySorting(apps []domain.Application) []domain.Application {
	// Sort by name alphabetically
	// Note: In a real implementation, you might want to use a proper sorting library
	// or allow configurable sorting options

	// Simple bubble sort for demonstration
	for i := 0; i < len(apps)-1; i++ {
		for j := 0; j < len(apps)-i-1; j++ {
			if apps[j].Name > apps[j+1].Name {
				apps[j], apps[j+1] = apps[j+1], apps[j]
			}
		}
	}

	return apps
}

func (uc *listAppsUseCase) applyPagination(apps []domain.Application, req ListAppsRequest) []domain.Application {
	if req.Limit == 0 {
		return apps // No pagination
	}

	start := req.Offset
	if start >= len(apps) {
		return []domain.Application{} // Empty result if offset is beyond range
	}

	end := start + req.Limit
	if end > len(apps) {
		end = len(apps)
	}

	return apps[start:end]
}
