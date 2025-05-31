package cmd

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type startupModel struct {
	width    int
	height   int
	styles   *Styles
	showPrompt bool
	currentDir string
}

func NewStartupModel() startupModel {
	cwd, _ := os.Getwd()
	gostmanFile := filepath.Join(cwd, "gostman.json")
	
	m := startupModel{
		width:      80,
		height:     24,
		showPrompt: checkFileExists(gostmanFile),
		currentDir: cwd,
	}
	
	lg := lipgloss.DefaultRenderer()
	m.styles = NewStyles(lg)
	
	return m
}

type transitionMsg struct{}

func (m startupModel) Init() tea.Cmd {
	if !m.showPrompt {
		// File exists, transition immediately
		return func() tea.Msg {
			return transitionMsg{}
		}
	}
	return nil
}

func (m startupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case transitionMsg:
		// File exists, set current project and transition
		if err := SetCurrentProject(m.currentDir); err != nil {
			// Handle error - continue anyway
		}
		mainModel := NewModel()
		mainModel.width = m.width
		mainModel.height = m.height
		return mainModel, nil
	case tea.KeyMsg:
		if m.showPrompt {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "y", "Y":
				// Create gostman.json in current directory
				if err := CreateProjectInCurrentDir(); err != nil {
					// Handle error - for now just continue to main app
				}
				// Switch to main model
				mainModel := NewModel()
				mainModel.width = m.width
				mainModel.height = m.height
				return mainModel, nil
			case "n", "N":
				// Continue to main app without creating file - use latest project if available
				if config.CurrentProject != "" {
					// Latest project is already set in config, just ensure jsonfilePath is correct
					jsonfilePath = filepath.Join(config.CurrentProject, "gostman.json")
				} else if len(config.Projects) > 0 {
					// No current project but we have projects, use the first one
					if err := SetCurrentProject(config.Projects[0].Path); err != nil {
						// If error, fall back to default app folder
					}
				} else {
					// No projects at all, create a default one in app folder
					defaultProject := filepath.Join(appFolder, "default")
					if err := SetCurrentProject(defaultProject); err != nil {
						// If error, continue with app folder fallback
					}
				}
				mainModel := NewModel()
				mainModel.width = m.width
				mainModel.height = m.height
				return mainModel, nil
			}
		}
	}
	
	return m, nil
}

func (m startupModel) View() string {
	if !m.showPrompt {
		// File exists, return empty as we're transitioning
		return ""
	}

	// Show creation prompt
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Gostman")

	// Determine what will happen if user chooses "no"
	fallbackText := "Use default project"
	if config.CurrentProject != "" {
		fallbackText = "Use: " + filepath.Base(config.CurrentProject)
	} else if len(config.Projects) > 0 {
		fallbackText = "Use: " + config.Projects[0].Name
	}

	prompt := lipgloss.NewStyle().
		Width(60).
		Align(lipgloss.Center).
		Render("No gostman.json found in current directory.\n\nWould you like to create one?\n\n(y) Create new project here\n(n) " + fallbackText + "\n\nChoice:")

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press Ctrl+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"\n",
		prompt,
		"\n",
		footer,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}