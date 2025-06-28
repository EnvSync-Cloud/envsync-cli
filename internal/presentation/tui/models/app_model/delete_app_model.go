package app_model

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory/shared"
)

// DeleteAppModel represents the state for multi-select app deletion
type DeleteAppModel struct {
	apps             []domain.Application
	selected         map[int]bool
	cursor           int
	deleteUseCase    app.DeleteAppUseCase
	ctx              context.Context
	state            deleteState
	confirmCursor    int
	deletedApps      []string
	deletionComplete bool
	err              error
}

type deleteState int

const (
	stateSelecting deleteState = iota
	stateConfirming
	stateDeleting
)

// NewDeleteAppModel creates a new delete app model
func NewDeleteAppModel(
	apps []domain.Application,
	deleteUseCase app.DeleteAppUseCase,
	ctx context.Context,
) *DeleteAppModel {
	return &DeleteAppModel{
		apps:          apps,
		selected:      make(map[int]bool),
		cursor:        0,
		deleteUseCase: deleteUseCase,
		ctx:           ctx,
		state:         stateSelecting,
		confirmCursor: 1, // Default to "No"
	}
}

// Init initializes the model
func (m *DeleteAppModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *DeleteAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateSelecting:
			return m.updateSelecting(msg)
		case stateConfirming:
			return m.updateConfirming(msg)
		case stateDeleting:
			return m, tea.Quit
		}
	}
	return m, nil
}

// updateSelecting handles key events in the selection state
func (m *DeleteAppModel) updateSelecting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.apps)-1 {
			m.cursor++
		}
	case " ":
		// Toggle selection
		shared.ToggleSelection(m.selected, m.cursor)
	case "enter":
		// Check if any apps are selected
		if shared.CountSelected(m.selected) > 0 {
			m.state = stateConfirming
		}
	case "a":
		// Select all
		m.selected = shared.SelectAll(m.apps)
	case "n":
		// Select none
		m.selected = shared.SelectNone()
	}
	return m, nil
}

// updateConfirming handles key events in the confirmation state
func (m *DeleteAppModel) updateConfirming(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.state = stateSelecting
		return m, nil
	case "left", "h":
		m.confirmCursor = 0 // Yes
	case "right", "l":
		m.confirmCursor = 1 // No
	case "enter":
		if m.confirmCursor == 0 { // Yes
			return m, m.deleteSelectedApps()
		} else { // No
			m.state = stateSelecting
		}
	}
	return m, nil
}

// deleteSelectedApps performs the deletion of selected applications
func (m *DeleteAppModel) deleteSelectedApps() tea.Cmd {
	return func() tea.Msg {
		m.state = stateDeleting
		var deletedApps []string
		var errors []string

		selectedApps := shared.GetSelectedApps(m.apps, m.selected)
		for _, application := range selectedApps {
			err := m.deleteUseCase.Execute(m.ctx, app.DeleteAppRequest{ID: application.ID})
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to delete %s: %v", application.Name, err))
			} else {
				deletedApps = append(deletedApps, application.Name)
			}
		}

		m.deletedApps = deletedApps
		if len(errors) > 0 {
			m.err = fmt.Errorf("deletion errors: %s", strings.Join(errors, "; "))
		} else {
			m.deletionComplete = true
		}

		return tea.Quit()
	}
}

// View renders the current view based on the state
func (m *DeleteAppModel) View() string {
	switch m.state {
	case stateSelecting:
		return m.viewSelecting()
	case stateConfirming:
		return m.viewConfirming()
	case stateDeleting:
		return m.viewDeleting()
	}
	return ""
}

// viewSelecting renders the selection view
func (m *DeleteAppModel) viewSelecting() string {
	var s strings.Builder

	s.WriteString("🗑️  Delete Applications (Multi-Select)\n\n")
	s.WriteString(fmt.Sprintf("📋 %s\n", shared.NavigationHelp()))
	s.WriteString(fmt.Sprintf("📋 %s, ENTER to confirm\n\n", shared.MultiSelectHelp()))

	selectedCount := shared.CountSelected(m.selected)
	if selectedCount > 0 {
		s.WriteString(fmt.Sprintf("✅ Selected: %d application(s)\n\n", selectedCount))
	}

	for i, application := range m.apps {
		cursor := "  "
		if m.cursor == i {
			cursor = "► "
		}

		checkbox := "☐ "
		if m.selected[i] {
			checkbox = "☑️ "
		}

		s.WriteString(shared.FormatAppListItem(application, cursor, checkbox))
		s.WriteString("\n")
	}

	if selectedCount == 0 {
		s.WriteString("💡 Select applications to delete, then press ENTER to continue")
	} else {
		s.WriteString("💡 Press ENTER to proceed with deletion")
	}

	return s.String()
}

// viewConfirming renders the confirmation view
func (m *DeleteAppModel) viewConfirming() string {
	var s strings.Builder

	s.WriteString("⚠️  Confirm Deletion\n\n")
	s.WriteString("🚨 You are about to delete the following applications:\n\n")

	selectedApps := shared.GetSelectedApps(m.apps, m.selected)
	for _, app := range selectedApps {
		s.WriteString(fmt.Sprintf("• 📛 %s (ID: %s)\n", app.Name, app.ID))
		if app.EnvCount != "" && app.EnvCount != "0" {
			s.WriteString(fmt.Sprintf("  🌍 %s environments will also be deleted\n", app.EnvCount))
		}
	}

	s.WriteString("\n🚨 This action cannot be undone!\n\n")
	s.WriteString("Are you sure you want to proceed?\n\n")

	// Confirmation buttons
	yesStyle := "  Yes  "
	noStyle := "  No   "

	if m.confirmCursor == 0 {
		yesStyle = "► Yes ◄"
	} else {
		noStyle = "► No  ◄"
	}

	s.WriteString(fmt.Sprintf("%s    %s\n\n", yesStyle, noStyle))
	s.WriteString(shared.ConfirmationHelp())

	return s.String()
}

// viewDeleting renders the deletion progress view
func (m *DeleteAppModel) viewDeleting() string {
	var s strings.Builder

	s.WriteString("🔄 Deleting Applications...\n\n")

	if len(m.deletedApps) > 0 {
		s.WriteString("✅ Successfully deleted:\n")
		for _, name := range m.deletedApps {
			s.WriteString(fmt.Sprintf("  • %s\n", name))
		}
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf("\n❌ Error: %s\n", m.err.Error()))
	}

	s.WriteString("\nPress any key to exit...")
	return s.String()
}

// GetError returns any error that occurred during deletion
func (m *DeleteAppModel) GetError() error {
	return m.err
}

// IsDeletionComplete returns whether the deletion was completed successfully
func (m *DeleteAppModel) IsDeletionComplete() bool {
	return m.deletionComplete
}

// GetDeletedApps returns the list of successfully deleted app names
func (m *DeleteAppModel) GetDeletedApps() []string {
	return m.deletedApps
}
