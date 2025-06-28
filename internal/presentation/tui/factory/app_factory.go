package factory

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/models/app_model"
)

type AppFactory struct {
	createUseCase app.CreateAppUseCase
	deleteUseCase app.DeleteAppUseCase
	listUseCase   app.ListAppsUseCase
	getUseCase    app.GetAppUseCase
}

func NewAppFactory(
	createUseCase app.CreateAppUseCase,
	deleteUseCase app.DeleteAppUseCase,
	listUseCase app.ListAppsUseCase,
	getUseCase app.GetAppUseCase,
) *AppFactory {
	return &AppFactory{
		createUseCase: createUseCase,
		deleteUseCase: deleteUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
	}
}

// CreateAppInteractive runs the interactive app creation flow
func (f *AppFactory) CreateAppInteractive(ctx context.Context) error {
	var name, description string
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Application Name").
				Description("Enter a unique name for your application").
				Placeholder("my-awesome-app").
				Value(&name).
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
				Value(&description).
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
		return err
	}

	if !confirm {
		return fmt.Errorf("application creation cancelled by user")
	}

	// Create the application
	req := app.CreateAppRequest{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Metadata:    make(map[string]any),
	}

	result, err := f.createUseCase.Execute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	// Show success confirmation with huh
	f.showSuccessConfirmation(result)

	return nil
}

// showSuccessConfirmation displays a success message and confirmation
func (f *AppFactory) showSuccessConfirmation(result *domain.Application) {
	successMsg := fmt.Sprintf("✅ Application '%s' created successfully!\n\n", result.Name)
	successMsg += fmt.Sprintf("📛 Name: %s\n", result.Name)
	successMsg += fmt.Sprintf("🆔 ID: %s\n", result.ID)
	if result.Description != "" {
		successMsg += fmt.Sprintf("📝 Description: %s\n", result.Description)
	}
}

// DeleteAppInteractive runs the interactive app deletion flow using Bubble Tea
func (f *AppFactory) DeleteAppInteractive(ctx context.Context) error {
	// Get list of apps first
	apps, err := f.listUseCase.Execute(ctx, app.ListAppsRequest{})
	if err != nil {
		return fmt.Errorf("failed to load applications: %w", err)
	}

	if len(apps) == 0 {
		return f.showNoAppsMessage("delete")
	}

	// Create the delete model
	model := app_model.NewDeleteAppModel(apps, f.deleteUseCase, ctx)

	// Run the program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	finalModel, err := program.Run()
	if err != nil {
		return err
	}

	// Check if deletion was successful
	if m, ok := finalModel.(*app_model.DeleteAppModel); ok {
		if m.GetError() != nil {
			return m.GetError()
		}
		if m.IsDeletionComplete() {
			fmt.Printf("✅ Successfully deleted %d application(s).\n", len(m.GetDeletedApps()))
		}
	}

	return nil
}

// // showNoAppsMessage displays a message when no apps are found
func (f *AppFactory) showNoAppsMessage(action string) error {
	var confirmed bool

	message := fmt.Sprintf("📭 No applications found.\n\nYou need to create an application first before you can %s one.", action)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("No Applications").
				Description(message).
				Affirmative("Create App").
				Negative("Exit").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return err
	}

	if confirmed {
		return f.CreateAppInteractive(context.Background())
	}

	return nil
}

// ListAppsInteractive runs the interactive app listing flow
func (f *AppFactory) ListAppsInteractive(ctx context.Context) error {
	// Create the list model with loading capability
	model := app_model.NewListAppModel(ctx, f.listUseCase)

	// Run the program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := program.Run()
	return err
}

// SelectAppInteractive runs the interactive app selection flow
func (f *AppFactory) SelectAppInteractive(ctx context.Context) error {
	// Create the list model with loading capability for selection
	model := app_model.NewListAppModel(ctx, f.listUseCase)

	// Run the program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	finalModel, err := program.Run()
	if err != nil {
		return err
	}

	// Check if a selection was made
	if listModel, ok := finalModel.(*app_model.ListAppModel); ok {
		selectedApp := listModel.GetSelectedApp()
		if selectedApp != nil {
			fmt.Printf("Selected application: %s (ID: %s)\n", selectedApp.Name, selectedApp.ID)
			return nil
		}
	}

	fmt.Println("No application selected")
	return nil
}

// Helper method for simple app listing (placeholder until full TUI is implemented)
func (f *AppFactory) runSimpleAppList(apps []domain.Application) error {
	// This method is now deprecated in favor of the new models
	// Keeping for backward compatibility
	return fmt.Errorf("runSimpleAppList is deprecated, use ListAppsInteractive instead")
}

// RunProgram is a helper to run any TUI program
func (f *AppFactory) RunProgram(model tea.Model) error {
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := program.Run()
	return err
}
