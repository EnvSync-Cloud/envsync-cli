package app_model

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
)

// ApplicationItem implements list.Item interface for domain.Application
type ApplicationItem struct {
	domain.Application
}

// FilterValue returns the string that should be matched when filtering
func (i ApplicationItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.Name, i.Description, i.ID)
}

// Title returns the title for the list item
func (i ApplicationItem) Title() string {
	return fmt.Sprintf("📛 %s", i.Name)
}

// Desc returns the description for the list item
func (i ApplicationItem) Desc() string {
	var desc strings.Builder

	desc.WriteString(fmt.Sprintf("🆔 %s\n", i.ID))

	if i.Application.Description != "" {
		maxLen := 60
		appDesc := i.Application.Description
		if len(appDesc) > maxLen {
			appDesc = appDesc[:maxLen-3] + "..."
		}
		desc.WriteString(fmt.Sprintf("📝 %s\n", appDesc))
	}

	if i.EnvCount != "" && i.EnvCount != "0" {
		desc.WriteString(fmt.Sprintf("🌍 %s envs", i.EnvCount))
	}

	return desc.String()
}

// ListAppModel represents the state for application listing using bubbles/list
type ListAppModel struct {
	list list.Model
	apps []domain.Application
}

// Define custom styles for the list
var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

// CustomDelegate extends the default list delegate with custom styling
type CustomDelegate struct {
	list.DefaultDelegate
}

// NewCustomDelegate creates a new custom delegate
func NewCustomDelegate() *CustomDelegate {
	d := &CustomDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
	}

	// Customize the delegate styles
	d.ShowDescription = true
	d.SetHeight(4) // Increase height to accommodate multi-line descriptions
	d.SetSpacing(1)

	// Define custom styles
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Bold(true)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#AD58B4", Dark: "#AD58B4"})

	return d
}

// Render implements the list.ItemDelegate interface
func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	d.DefaultDelegate.Render(w, m, index, listItem)
}

// NewListAppModelWithApps creates a new list app model with pre-loaded apps
func NewListAppModelWithApps(apps []domain.Application) *ListAppModel {
	// Convert domain.Application slice to list.Item slice
	items := make([]list.Item, len(apps))
	for i, app := range apps {
		items[i] = ApplicationItem{Application: app}
	}

	// Create the list with custom delegate
	delegate := NewCustomDelegate()
	l := list.New(items, delegate, 80, 24)

	// Configure the list
	l.Title = "🚀 Applications List"
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	model := &ListAppModel{
		list: l,
		apps: apps,
	}

	return model
}

// Init initializes the model
func (m *ListAppModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *ListAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// Update the list and get any commands
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the current view
func (m *ListAppModel) View() string {
	if len(m.apps) == 0 {
		return docStyle.Render(`
🚀 Applications List

📭 No applications found.

💡 Create your first application to get started!

Press Q to quit.
		`)
	}

	return docStyle.Render(m.list.View())
}

// GetApps returns all applications
func (m *ListAppModel) GetApps() []domain.Application {
	return m.apps
}

// GetSelectedApp returns the currently selected application
func (m *ListAppModel) GetSelectedApp() *domain.Application {
	if item, ok := m.list.SelectedItem().(ApplicationItem); ok {
		return &item.Application
	}
	return nil
}
