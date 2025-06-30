package app_model

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory/shared"
)

// DeleteAppModel represents the state for multi-select app deletion
type DeleteAppModel struct {
	apps         []domain.Application
	selectedApps []domain.Application
	cursor       int
	state        deleteState
}

type deleteState int

const (
	stateSelecting deleteState = iota
	stateSelected
)

// NewDeleteAppModel creates a new delete app model
func NewDeleteAppModel(
	apps []domain.Application,
) *DeleteAppModel {
	return &DeleteAppModel{
		apps:         apps,
		selectedApps: make([]domain.Application, 0),
		cursor:       0,
		state:        stateSelecting,
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
		case stateSelected:
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
		m.toggleSelectedApp()
	case "enter":
		// Check if any apps are selected
		m.state = stateSelected
	case "a":
		// Select all
		m.selectedApps = m.apps
	case "n":
		// Select none
		m.selectedApps = []domain.Application{}
	}
	return m, nil
}

// View renders the current view based on the state
func (m *DeleteAppModel) View() string {
	switch m.state {
	case stateSelecting:
		return m.viewSelecting()
	}

	return ""
}

// viewSelecting renders the selection view
func (m *DeleteAppModel) viewSelecting() string {
	var s strings.Builder

	s.WriteString("🗑️  Delete Applications (Multi-Select)\n\n")
	s.WriteString(fmt.Sprintf("📋 %s\n", shared.NavigationHelp()))
	s.WriteString(fmt.Sprintf("📋 %s, ENTER to confirm\n\n", shared.MultiSelectHelp()))

	selectedCount := len(m.selectedApps)
	if selectedCount > 0 {
		s.WriteString(fmt.Sprintf("✅ Selected: %d application(s)\n\n", selectedCount))
	}

	for i, application := range m.apps {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		checkbox := "[ ] "
		// Compare if application is present in selectedApps
		if slices.ContainsFunc(m.selectedApps, func(a domain.Application) bool {
			return a.ID == application.ID
		}) {
			checkbox = "[x] "
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

func (m *DeleteAppModel) GetSelectedApps() []domain.Application {
	return m.selectedApps
}

func (m *DeleteAppModel) toggleSelectedApp() {
	if m.cursor < 0 || m.cursor >= len(m.apps) {
		return
	}

	app := m.apps[m.cursor]
	// Check if app is already selected
	for i, selectedApp := range m.selectedApps {
		if selectedApp.ID == app.ID {
			// Remove from selection
			m.selectedApps = append(m.selectedApps[:i], m.selectedApps[i+1:]...)
			return
		}
	}

	// Add to selection
	m.selectedApps = append(m.selectedApps, app)

}

// viewConfirming renders the confirmation view
// func (m *DeleteAppModel) viewConfirming() string {
// 	var s strings.Builder

// 	s.WriteString("⚠️  Confirm Deletion\n\n")
// 	s.WriteString("🚨 You are about to delete the following applications:\n\n")

// 	selectedApps := shared.GetSelectedApps(m.apps, m.selected)
// 	for _, app := range selectedApps {
// 		s.WriteString(fmt.Sprintf("• 📛 %s (ID: %s)\n", app.Name, app.ID))
// 		if app.EnvCount != "" && app.EnvCount != "0" {
// 			s.WriteString(fmt.Sprintf("  🌍 %s environments will also be deleted\n", app.EnvCount))
// 		}
// 	}

// 	s.WriteString("\n🚨 This action cannot be undone!\n\n")
// 	s.WriteString("Are you sure you want to proceed?\n\n")

// 	// Confirmation buttons
// 	yesStyle := "  Yes  "
// 	noStyle := "  No   "

// 	if m.confirmCursor == 0 {
// 		yesStyle = "► Yes ◄"
// 	} else {
// 		noStyle = "► No  ◄"
// 	}

// 	s.WriteString(fmt.Sprintf("%s    %s\n\n", yesStyle, noStyle))
// 	s.WriteString(shared.ConfirmationHelp())

// 	return s.String()
// }

// // viewDeleting renders the deletion progress view
// func (m *DeleteAppModel) viewDeleting() string {
// 	var s strings.Builder

// 	s.WriteString("🔄 Deleting Applications...\n\n")

// 	if len(m.deletedApps) > 0 {
// 		s.WriteString("✅ Successfully deleted:\n")
// 		for _, name := range m.deletedApps {
// 			s.WriteString(fmt.Sprintf("  • %s\n", name))
// 		}
// 	}

// 	if m.err != nil {
// 		s.WriteString(fmt.Sprintf("\n❌ Error: %s\n", m.err.Error()))
// 	}

// 	s.WriteString("\nPress any key to exit...")
// 	return s.String()
// }

// // GetError returns any error that occurred during deletion
// func (m *DeleteAppModel) GetError() error {
// 	return m.err
// }

// // IsDeletionComplete returns whether the deletion was completed successfully
// func (m *DeleteAppModel) IsDeletionComplete() bool {
// 	return m.deletionComplete
// }

// // GetDeletedApps returns the list of successfully deleted app names
// func (m *DeleteAppModel) GetDeletedApps() []string {
// 	return m.deletedApps
// }
