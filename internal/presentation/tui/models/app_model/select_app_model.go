package app_model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/presentation/tui/factory/shared"
)

// viewport represents a scrollable viewport
type viewport struct {
	top    int
	height int
}

// SelectAppModel represents the state for application selection
type SelectAppModel struct {
	apps          []domain.Application
	cursor        int
	selectedApp   *domain.Application
	filterText    string
	filteredApps  []domain.Application
	showFilter    bool
	viewport      viewport
	selectionMade bool
}

// NewSelectAppModel creates a new select app model
func NewSelectAppModel(apps []domain.Application) *SelectAppModel {
	model := &SelectAppModel{
		apps:         apps,
		cursor:       0,
		filteredApps: apps,
		viewport: viewport{
			top:    0,
			height: 10, // Default height, can be adjusted
		},
	}
	return model
}

// Init initializes the model
func (m *SelectAppModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *SelectAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.height = msg.Height - 10 // Leave space for header and footer
		return m, nil
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

// handleKeyPress processes keyboard input
func (m *SelectAppModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle filter mode
	if m.showFilter {
		switch msg.String() {
		case "esc":
			m.showFilter = false
			m.filterText = ""
			m.filteredApps = m.apps
			m.cursor = 0
			return m, nil
		case "enter":
			m.showFilter = false
			return m, nil
		case "backspace":
			if len(m.filterText) > 0 {
				m.filterText = m.filterText[:len(m.filterText)-1]
				m.applyFilter()
			}
			return m, nil
		default:
			// Add character to filter
			if len(msg.String()) == 1 && msg.String() != " " {
				m.filterText += msg.String()
				m.applyFilter()
			}
		}
		return m, nil
	}

	// Handle navigation mode
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.adjustViewport()
		}
	case "down", "j":
		if m.cursor < len(m.filteredApps)-1 {
			m.cursor++
			m.adjustViewport()
		}
	case "home", "g":
		m.cursor = 0
		m.viewport.top = 0
	case "end", "G":
		m.cursor = len(m.filteredApps) - 1
		m.adjustViewport()
	case "pageup":
		m.cursor = shared.Max(0, m.cursor-m.viewport.height)
		m.adjustViewport()
	case "pagedown":
		m.cursor = shared.Min(len(m.filteredApps)-1, m.cursor+m.viewport.height)
		m.adjustViewport()
	case "enter", " ":
		// Select the current application
		if m.cursor >= 0 && m.cursor < len(m.filteredApps) {
			m.selectedApp = &m.filteredApps[m.cursor]
			m.selectionMade = true
			return m, tea.Quit
		}
	case "/":
		m.showFilter = true
		return m, nil
	case "r":
		// Reset filter
		m.filterText = ""
		m.filteredApps = m.apps
		m.cursor = 0
		m.viewport.top = 0
	}
	return m, nil
}

// applyFilter filters the applications based on the filter text
func (m *SelectAppModel) applyFilter() {
	if m.filterText == "" {
		m.filteredApps = m.apps
	} else {
		m.filteredApps = nil
		filterLower := strings.ToLower(m.filterText)
		for _, app := range m.apps {
			if strings.Contains(strings.ToLower(app.Name), filterLower) ||
				strings.Contains(strings.ToLower(app.Description), filterLower) ||
				strings.Contains(strings.ToLower(app.ID), filterLower) {
				m.filteredApps = append(m.filteredApps, app)
			}
		}
	}

	// Adjust cursor if it's out of bounds
	if m.cursor >= len(m.filteredApps) {
		m.cursor = shared.Max(0, len(m.filteredApps)-1)
	}
	m.adjustViewport()
}

// adjustViewport adjusts the viewport to keep the cursor visible
func (m *SelectAppModel) adjustViewport() {
	if m.cursor < m.viewport.top {
		m.viewport.top = m.cursor
	} else if m.cursor >= m.viewport.top+m.viewport.height {
		m.viewport.top = m.cursor - m.viewport.height + 1
	}
}

// View renders the current view
func (m *SelectAppModel) View() string {
	var s strings.Builder

	// Header
	s.WriteString("🎯 Select Application\n\n")

	if len(m.apps) == 0 {
		s.WriteString("📭 No applications found.\n\n")
		s.WriteString("💡 Create your first application to get started!")
		return s.String()
	}

	// Filter status
	if m.showFilter {
		s.WriteString(fmt.Sprintf("🔍 Filter: %s_\n", m.filterText))
		s.WriteString("Type to filter, ESC to cancel, ENTER to apply\n\n")
	} else if m.filterText != "" {
		s.WriteString(fmt.Sprintf("🔍 Filtered by: '%s' (%d/%d apps)\n",
			m.filterText, len(m.filteredApps), len(m.apps)))
		s.WriteString("Press 'r' to reset filter\n\n")
	}

	// Applications list
	if len(m.filteredApps) == 0 {
		s.WriteString("📭 No applications match your filter.\n\n")
		s.WriteString("💡 Try a different filter or press 'r' to reset")
		return s.String()
	}

	// Calculate visible range
	start := m.viewport.top
	end := shared.Min(start+m.viewport.height, len(m.filteredApps))

	for i := start; i < end; i++ {
		app := m.filteredApps[i]
		cursor := "   "
		if i == m.cursor {
			cursor = "►  "
		}

		s.WriteString(shared.FormatAppListItem(app, cursor, ""))
		s.WriteString("\n")
	}

	// Footer with navigation help
	s.WriteString("\n")
	s.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	s.WriteString(fmt.Sprintf("📊 Showing %d-%d of %d applications",
		start+1, end, len(m.filteredApps)))

	if len(m.filteredApps) != len(m.apps) {
		s.WriteString(fmt.Sprintf(" (filtered from %d total)", len(m.apps)))
	}
	s.WriteString("\n")

	if !m.showFilter {
		s.WriteString("💡 Navigate: ↑/↓ • Select: ENTER/SPACE • Filter: / • Reset: r • Quit: q")
	}

	return s.String()
}

// GetSelectedApp returns the selected application
func (m *SelectAppModel) GetSelectedApp() *domain.Application {
	return m.selectedApp
}

// IsSelectionMade returns whether a selection has been made
func (m *SelectAppModel) IsSelectionMade() bool {
	return m.selectionMade
}

// GetApps returns all applications
func (m *SelectAppModel) GetApps() []domain.Application {
	return m.apps
}

// GetFilteredApps returns filtered applications
func (m *SelectAppModel) GetFilteredApps() []domain.Application {
	return m.filteredApps
}

// GetCurrentApp returns the currently highlighted application
func (m *SelectAppModel) GetCurrentApp() *domain.Application {
	if m.cursor >= 0 && m.cursor < len(m.filteredApps) {
		return &m.filteredApps[m.cursor]
	}
	return nil
}

// Helper functions moved to shared/utils.go
