package domain

import "time"

type Application struct {
	ID          string
	Name        string
	Description string
	Metadata    map[string]any
	OrgID       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewApplication(name, description string, metadata map[string]any) *Application {
	return &Application{
		Name:        name,
		Description: description,
		Metadata:    metadata,
	}
}

// EnvironmentType represents an environment configuration for an application
type EnvironmentType struct {
	ID        string
	OrgID     string
	AppID     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
