package app

import (
	"context"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
)

// CreateAppUseCase defines the interface for creating applications
type CreateAppUseCase interface {
	Execute(ctx context.Context, req CreateAppRequest) (*domain.Application, error)
}

// DeleteAppUseCase defines the interface for deleting applications
type DeleteAppUseCase interface {
	Execute(ctx context.Context, req DeleteAppRequest) error
}

// ListAppsUseCase defines the interface for listing applications
type ListAppsUseCase interface {
	Execute(ctx context.Context, req ListAppsRequest) ([]domain.Application, error)
}

// GetAppUseCase defines the interface for getting a single application
type GetAppUseCase interface {
	Execute(ctx context.Context, req GetAppRequest) (*domain.Application, error)
}

// UpdateAppUseCase defines the interface for updating applications
type UpdateAppUseCase interface {
	Execute(ctx context.Context, req UpdateAppRequest) (*domain.Application, error)
}

// Request/Response types

type CreateAppRequest struct {
	Name        string
	Description string
	Metadata    map[string]any
}

type DeleteAppRequest struct {
	ID   string
	Name string // Alternative identifier
}

type ListAppsRequest struct {
	// Add filtering options if needed
	OrgID  string
	Limit  int
	Offset int
}

type GetAppRequest struct {
	ID   string
	Name string // Alternative identifier
}

type UpdateAppRequest struct {
	ID          string
	Name        *string
	Description *string
	Metadata    map[string]any
}

// Validation interface for requests
type Validator interface {
	Validate() error
}

// Implement validation for each request type
func (r CreateAppRequest) Validate() error {
	if r.Name == "" {
		return ErrAppNameRequired
	}
	if r.Description == "" {
		return ErrAppDescriptionRequired
	}
	return nil
}

func (r DeleteAppRequest) Validate() error {
	if r.ID == "" && r.Name == "" {
		return ErrAppIdentifierRequired
	}
	return nil
}

func (r ListAppsRequest) Validate() error {
	if r.Limit < 0 {
		return ErrInvalidLimit
	}
	if r.Offset < 0 {
		return ErrInvalidOffset
	}
	return nil
}

func (r GetAppRequest) Validate() error {
	if r.ID == "" && r.Name == "" {
		return ErrAppIdentifierRequired
	}
	return nil
}

func (r UpdateAppRequest) Validate() error {
	if r.ID == "" {
		return ErrAppIDRequired
	}
	return nil
}
