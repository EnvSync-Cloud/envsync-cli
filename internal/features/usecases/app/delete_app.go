package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type deleteAppUseCase struct {
	appService services.ApplicationService
}

func NewDeleteAppUseCase(appService services.ApplicationService) DeleteAppUseCase {
	return &deleteAppUseCase{
		appService: appService,
	}
}

func (uc *deleteAppUseCase) Execute(ctx context.Context, req DeleteAppRequest) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return NewValidationError("invalid delete app request", err)
	}

	// Find the application to delete
	app, err := uc.findApplication(req)
	if err != nil {
		return err
	}

	// Check if application can be deleted (business rules)
	if err := uc.validateDeletion(app); err != nil {
		return err
	}

	// Delete application via service
	if err := uc.appService.DeleteApp(*app); err != nil {
		return NewServiceError("failed to delete application", err)
	}

	return nil
}

func (uc *deleteAppUseCase) findApplication(req DeleteAppRequest) (*domain.Application, error) {
	// Get all applications
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return nil, NewServiceError("failed to retrieve applications", err)
	}

	// Find by ID or name
	for _, app := range apps {
		if (req.ID != "" && app.ID == req.ID) ||
			(req.Name != "" && strings.EqualFold(app.Name, req.Name)) {
			return &app, nil
		}
	}

	// Application not found
	identifier := req.ID
	if identifier == "" {
		identifier = req.Name
	}

	return nil, NewNotFoundError(
		fmt.Sprintf("application '%s' not found", identifier),
		ErrAppNotFound,
	)
}

func (uc *deleteAppUseCase) validateDeletion(app *domain.Application) error {
	// Check if application has environments
	if err := uc.checkEnvironmentDependencies(app); err != nil {
		return err
	}

	// Check if application is currently in use
	if err := uc.checkApplicationInUse(app); err != nil {
		return err
	}

	// Additional business rules can be added here

	return nil
}

func (uc *deleteAppUseCase) checkEnvironmentDependencies(app *domain.Application) error {
	// Check if application has any environment types
	if len(app.EnvTypes) > 0 {
		return NewInUseError(
			fmt.Sprintf("cannot delete application '%s': it has %d environment type(s)",
				app.Name, len(app.EnvTypes)),
			ErrAppHasEnvironments,
		)
	}

	// Check environment count if available
	if app.EnvCount != "" && app.EnvCount != "0" {
		return NewInUseError(
			fmt.Sprintf("cannot delete application '%s': it has %s environment(s)",
				app.Name, app.EnvCount),
			ErrAppHasEnvironments,
		)
	}

	return nil
}

func (uc *deleteAppUseCase) checkApplicationInUse(app *domain.Application) error {
	// In a real implementation, you might check:
	// - Active deployments
	// - Running processes
	// - Active user sessions
	// - Scheduled jobs
	// - External integrations

	// For now, we'll implement a basic check
	// This could be extended to call other services

	// Example: Check if app has recent activity
	// if uc.hasRecentActivity(app) {
	//     return NewInUseError(
	//         fmt.Sprintf("cannot delete application '%s': it has recent activity", app.Name),
	//         ErrAppInUse,
	//     )
	// }

	return nil
}

// Helper method to check recent activity (placeholder)
func (uc *deleteAppUseCase) hasRecentActivity(app *domain.Application) bool {
	// This would typically call an analytics or activity service
	// to check if the application has been used recently

	// For now, return false (no recent activity check)
	return false
}

// Additional helper methods for extended validation

func (uc *deleteAppUseCase) validateUserPermissions(ctx context.Context, app *domain.Application) error {
	// In a real implementation, you would:
	// 1. Extract user context from ctx
	// 2. Check if user has delete permissions for this app
	// 3. Check organization-level permissions

	// For now, we'll assume permissions are handled at a higher level
	return nil
}

func (uc *deleteAppUseCase) createDeletionAuditLog(ctx context.Context, app *domain.Application) error {
	// In a real implementation, you would:
	// 1. Create an audit log entry
	// 2. Record who deleted the app
	// 3. Record when it was deleted
	// 4. Store app metadata for potential recovery

	// For now, this is a placeholder
	return nil
}

func (uc *deleteAppUseCase) notifyStakeholders(ctx context.Context, app *domain.Application) error {
	// In a real implementation, you might:
	// 1. Send notifications to app owners
	// 2. Update external systems
	// 3. Trigger cleanup workflows

	// For now, this is a placeholder
	return nil
}
