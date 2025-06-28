package app

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/services"
)

type createAppUseCase struct {
	appService services.ApplicationService
}

func NewCreateAppUseCase(appService services.ApplicationService) CreateAppUseCase {
	return &createAppUseCase{
		appService: appService,
	}
}

func (uc *createAppUseCase) Execute(ctx context.Context, req CreateAppRequest) (*domain.Application, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, NewValidationError("invalid create app request", err)
	}

	// Additional business validation
	if err := uc.validateBusinessRules(req); err != nil {
		return nil, err
	}

	// Create domain object
	app := domain.NewApplication(req.Name, req.Description, req.Metadata)

	// Check if application with same name already exists
	if exists, err := uc.checkApplicationExists(req.Name); err != nil {
		return nil, NewServiceError("failed to check application existence", err)
	} else if exists {
		return nil, NewAlreadyExistsError(
			fmt.Sprintf("application with name '%s' already exists", req.Name),
			ErrAppAlreadyExists,
		)
	}

	// Create application via service
	createdApp, err := uc.appService.CreateApp(app)
	if err != nil {
		return nil, NewServiceError("failed to create application", err)
	}

	return &createdApp, nil
}

func (uc *createAppUseCase) validateBusinessRules(req CreateAppRequest) error {
	// Validate name length
	if len(req.Name) > 100 {
		return NewValidationError("application name too long", ErrAppNameTooLong)
	}

	// Validate description length
	if len(req.Description) > 500 {
		return NewValidationError("application description too long", ErrAppDescriptionTooLong)
	}

	// Validate name format (alphanumeric, hyphens, underscores only)
	if !uc.isValidAppName(req.Name) {
		return NewValidationError("invalid application name format", ErrInvalidAppName)
	}

	// Validate metadata size
	if err := uc.validateMetadata(req.Metadata); err != nil {
		return NewValidationError("invalid metadata", err)
	}

	return nil
}

func (uc *createAppUseCase) isValidAppName(name string) bool {
	// Allow alphanumeric characters, hyphens, and underscores
	// Must start with a letter or number
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9_-]*$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}

func (uc *createAppUseCase) validateMetadata(metadata map[string]any) error {
	if len(metadata) > 20 {
		return fmt.Errorf("too many metadata entries (max 20)")
	}

	for key, value := range metadata {
		// Validate key format
		if len(key) > 50 {
			return fmt.Errorf("metadata key '%s' is too long (max 50 characters)", key)
		}

		if !uc.isValidMetadataKey(key) {
			return fmt.Errorf("metadata key '%s' contains invalid characters", key)
		}

		// Validate value
		if err := uc.validateMetadataValue(key, value); err != nil {
			return err
		}
	}

	return nil
}

func (uc *createAppUseCase) isValidMetadataKey(key string) bool {
	// Allow alphanumeric characters, hyphens, underscores, and dots
	pattern := `^[a-zA-Z0-9._-]+$`
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func (uc *createAppUseCase) validateMetadataValue(key string, value any) error {
	switch v := value.(type) {
	case string:
		if len(v) > 200 {
			return fmt.Errorf("metadata value for key '%s' is too long (max 200 characters)", key)
		}
	case int, int32, int64, float32, float64, bool:
		// These types are acceptable
	default:
		return fmt.Errorf("metadata value for key '%s' has unsupported type", key)
	}

	return nil
}

func (uc *createAppUseCase) checkApplicationExists(name string) (bool, error) {
	// Get all applications and check if name exists
	apps, err := uc.appService.GetAllApps()
	if err != nil {
		return false, err
	}

	for _, app := range apps {
		if strings.EqualFold(app.Name, name) {
			return true, nil
		}
	}

	return false, nil
}
