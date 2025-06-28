package app_model

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
	"github.com/EnvSync-Cloud/envsync-cli/internal/features/usecases/app"
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

// Description returns the description for the list item
func (i ApplicationItem) Description() string {
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

// LoadingState represents the different loading states
type LoadingState int

const (
	LoadingStateIdle LoadingState = iota
	LoadingStateLoading
	LoadingStateLoaded
	LoadingStateError
)

// LoadAppsMsg is sent when apps are loaded
type LoadAppsMsg struct {
	Apps []domain.Application
	Err  error
}

// ListAppModel represents the state for application listing using bubbles/list with loading spinner
type ListAppModel struct {
	list         list.Model
	spinner      spinner.Model
	apps         []domain.Application
	loadingState LoadingState
	error        error
	quitting     bool
	ctx          context.Context
	listUseCase  app.ListAppsUseCase
}

// Define custom styles for the list
var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	loadingStyle = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	errorStyle = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Foreground(lipgloss.Color("196"))
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
	d.SetHeight(3) // Increase height to accommodate multi-line descriptions
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

// NewListAppModel creates a new list app model using bubbles/list with loading capability
func NewListAppModel(ctx context.Context, listUseCase app.ListAppsUseCase) *ListAppModel {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create empty list initially
	delegate := NewCustomDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)

	// Configure the list
	l.Title = "🚀 Applications List"
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	// Customize help with additional key bindings
	resetKey := key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset filter"),
	)
	refreshKey := key.NewBinding(
		key.WithKeys("F5"),
		key.WithHelp("F5", "refresh"),
	)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{resetKey, refreshKey}
	}

	model := &ListAppModel{
		list:         l,
		spinner:      s,
		loadingState: LoadingStateIdle,
		ctx:          ctx,
		listUseCase:  listUseCase,
	}

	return model
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
	l := list.New(items, delegate, 0, 0)

	// Configure the list
	l.Title = "🚀 Applications List"
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	// Customize help with additional key bindings
	resetKey := key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset filter"),
	)
	refreshKey := key.NewBinding(
		key.WithKeys("F5"),
		key.WithHelp("F5", "refresh"),
	)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{resetKey, refreshKey}
	}

	model := &ListAppModel{
		list:         l,
		apps:         apps,
		loadingState: LoadingStateLoaded,
	}

	return model
}

// Init initializes the model and starts loading if needed
func (m *ListAppModel) Init() tea.Cmd {
	if m.loadingState == LoadingStateIdle && m.listUseCase != nil {
		m.loadingState = LoadingStateLoading
		return tea.Batch(
			m.spinner.Tick,
			m.loadApps(),
		)
	}
	return nil
}

// Update handles messages and updates the model state
func (m *ListAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case LoadAppsMsg:
		m.loadingState = LoadingStateLoaded
		if msg.Err != nil {
			m.loadingState = LoadingStateError
			m.error = msg.Err
			return m, nil
		}

		m.apps = msg.Apps
		m.error = nil

		// Update list items
		items := make([]list.Item, len(msg.Apps))
		for i, app := range msg.Apps {
			items[i] = ApplicationItem{Application: app}
		}
		m.list.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		// Handle loading state key bindings
		if m.loadingState == LoadingStateLoading {
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.quitting = true
				return m, tea.Quit
			}
			// Update spinner during loading
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

		// Handle error state key bindings
		if m.loadingState == LoadingStateError {
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				m.quitting = true
				return m, tea.Quit
			case "F5", "r":
				// Retry loading
				m.loadingState = LoadingStateLoading
				m.error = nil
				return m, tea.Batch(
					m.spinner.Tick,
					m.loadApps(),
				)
			}
			return m, nil
		}

		// Handle normal state key bindings
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// Reset filter
			m.list.ResetFilter()
			return m, nil
		case "F5":
			// Refresh apps
			if m.listUseCase != nil {
				m.loadingState = LoadingStateLoading
				return m, tea.Batch(
					m.spinner.Tick,
					m.loadApps(),
				)
			}
			return m, nil
		case "enter":
			// Handle selection
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				if appItem, ok := selectedItem.(ApplicationItem); ok {
					m.list.NewStatusMessage(fmt.Sprintf("Selected: %s", appItem.Name))
				}
			}
			return m, nil
		}
	}

	// Update list if in loaded state
	if m.loadingState == LoadingStateLoaded {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// Update spinner if in loading state
	if m.loadingState == LoadingStateLoading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the current view
func (m *ListAppModel) View() string {
	if m.quitting {
		return ""
	}

	switch m.loadingState {
	case LoadingStateLoading:
		content := fmt.Sprintf("%s Loading applications...\n\n💡 Press Ctrl+C, Q, or ESC to cancel",
			m.spinner.View(),
		)
		return loadingStyle.Render(content)

	case LoadingStateError:
		content := fmt.Sprintf("❌ Error loading applications:\n\n%s\n\n💡 Press F5 or R to retry, Q to quit",
			m.error.Error(),
		)
		return errorStyle.Render(content)

	case LoadingStateLoaded:
		if len(m.apps) == 0 {
			return docStyle.Render(`
					🚀 Applications List

					📭 No applications found.

					💡 Create your first application to get started!

					Press F5 to refresh or Q to quit.
				`)
		}

		return docStyle.Render(m.list.View())

	default:
		return docStyle.Render("Initializing...")
	}
}

// loadApps performs the actual loading of applications
func (m *ListAppModel) loadApps() tea.Cmd {
	if m.listUseCase == nil {
		return nil
	}

	return func() tea.Msg {
		apps, err := m.listUseCase.Execute(m.ctx, app.ListAppsRequest{})
		return LoadAppsMsg{
			Apps: apps,
			Err:  err,
		}
	}
}

// GetSelectedApp returns the currently selected application
func (m *ListAppModel) GetSelectedApp() *domain.Application {
	if m.loadingState != LoadingStateLoaded {
		return nil
	}

	selectedItem := m.list.SelectedItem()
	if selectedItem != nil {
		if appItem, ok := selectedItem.(ApplicationItem); ok {
			return &appItem.Application
		}
	}
	return nil
}

// GetApps returns all applications
func (m *ListAppModel) GetApps() []domain.Application {
	return m.apps
}

// GetFilteredApps returns filtered applications
func (m *ListAppModel) GetFilteredApps() []domain.Application {
	if m.loadingState != LoadingStateLoaded {
		return nil
	}

	filteredItems := m.list.VisibleItems()
	filteredApps := make([]domain.Application, len(filteredItems))

	for i, item := range filteredItems {
		if appItem, ok := item.(ApplicationItem); ok {
			filteredApps[i] = appItem.Application
		}
	}

	return filteredApps
}

// SetSize sets the size of the list
func (m *ListAppModel) SetSize(width, height int) {
	h, v := docStyle.GetFrameSize()
	m.list.SetSize(width-h, height-v)
}

// SetItems updates the list items (useful for refreshing data)
func (m *ListAppModel) SetItems(apps []domain.Application) {
	m.apps = apps
	items := make([]list.Item, len(apps))
	for i, app := range apps {
		items[i] = ApplicationItem{Application: app}
	}
	m.list.SetItems(items)
	m.loadingState = LoadingStateLoaded
}

// IsFiltering returns true if the list is currently being filtered
func (m *ListAppModel) IsFiltering() bool {
	if m.loadingState != LoadingStateLoaded {
		return false
	}
	return m.list.FilterState() == list.Filtering
}

// GetFilterValue returns the current filter value
func (m *ListAppModel) GetFilterValue() string {
	if m.loadingState != LoadingStateLoaded {
		return ""
	}
	return m.list.FilterValue()
}

// IsLoading returns true if the model is currently loading
func (m *ListAppModel) IsLoading() bool {
	return m.loadingState == LoadingStateLoading
}

// HasError returns true if there was an error during loading
func (m *ListAppModel) HasError() bool {
	return m.loadingState == LoadingStateError
}

// GetError returns the loading error
func (m *ListAppModel) GetError() error {
	return m.error
}

// IsLoaded returns true if the model has finished loading successfully
func (m *ListAppModel) IsLoaded() bool {
	return m.loadingState == LoadingStateLoaded
}
