package factory

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/models/app_model"
)

type AppFactory struct{}

func NewAppFactory() *AppFactory {
	return &AppFactory{}
}

// CreateAppTUI runs the interactive app creation flow
func (f *AppFactory) CreateAppTUI(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Application Name").
				Description("Enter a unique name for your application").
				Placeholder("my-awesome-app").
				Value(&app.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("application name is required")
					}
					if len(str) > 100 {
						return fmt.Errorf("application name must be 100 characters or less")
					}
					return nil
				}),

			huh.NewText().
				Title("Description").
				Description("Provide a description for your application").
				Placeholder("A brief description of what this application does...").
				Value(&app.Description).
				Lines(3).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("application description is required")
					}
					if len(str) > 500 {
						return fmt.Errorf("description must be 500 characters or less")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Are you sure?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return nil, err
	}

	if !confirm {
		return nil, errors.New("application creation cancelled by user")
	}

	return app, nil
}

// DeleteAppTUI runs the interactive app deletion flow using Bubble Tea
func (f *AppFactory) DeleteAppsTUI(apps []domain.Application) ([]domain.Application, error) {
	m := app_model.NewDeleteAppModel(apps)

	p := tea.NewProgram(m, tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running delete app TUI: %w", err)
	}

	if deleteModel, ok := model.(*app_model.DeleteAppModel); ok {
		return deleteModel.GetSelectedApps(), nil
	}

	return nil, fmt.Errorf("unexpected model type: %T", model)
}

// ListAppsInteractive runs the interactive app listing flow
func (f *AppFactory) ListAppsInteractive(apps []domain.Application) error {
	// // Create the list model with loading capability
	model := app_model.NewListAppModelWithApps(apps)

	// Run the program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := program.Run()
	if err != nil {
		return fmt.Errorf("error running app list TUI: %w", err)
	}

	return nil
}
