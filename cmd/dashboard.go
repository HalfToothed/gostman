package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listItem struct {
	request Request
}

func (i listItem) Title() string       { return i.request.Name }
func (i listItem) Description() string { return i.request.Method }
func (i listItem) FilterValue() string { return i.request.Name }

type projectItem struct {
	project Project
}

func (i projectItem) Title() string       { 
	title := i.project.Name
	if i.project.Path == GetCurrentProject() {
		return "• " + title + " (current)"
	}
	return "  " + title
}
func (i projectItem) Description() string { 
	if i.project.Path == GetCurrentProject() {
		return "→ " + i.project.Path
	}
	return "  " + i.project.Path 
}
func (i projectItem) FilterValue() string { return i.project.Name }

type board struct {
	width       int
	height      int
	styles      *Styles
	list        list.Model
	returnModel *Model
	showMsg        bool
	showProjects   bool
	projectList    list.Model
	showInput      bool
	projectInput   textinput.Model
}

func dashboard(width, height int, styles *Styles, returnModel *Model) board {
	var saved_data SavedData

	if !checkFileExists(jsonfilePath) {
		file, err := os.ReadFile(jsonfilePath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &saved_data)
	}

	savedRequests := saved_data.Requests

	// Convert requests to list items
	var items []list.Item
	for _, req := range savedRequests {
		items = append(items, listItem{request: req})
	}

	board := board{
		width:       width,
		height:      height,
		styles:      styles,
		list:        list.New(items, list.NewDefaultDelegate(), width, height-3),
		returnModel: returnModel,
		showMsg:     false,
	}

	board.list.Title = "List of Requests "

	board.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Create,
			Keymap.Delete,
			Keymap.Paths,
			Keymap.Back,
		}
	}

	// Initialize project list
	var projectItems []list.Item
	for _, project := range GetProjects() {
		projectItems = append(projectItems, projectItem{project: project})
	}
	
	board.projectList = list.New(projectItems, list.NewDefaultDelegate(), width, height-3)
	board.projectList.Title = "Projects"
	board.projectList.SetShowStatusBar(false)
	board.projectList.SetFilteringEnabled(false)
	board.projectList.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	// Initialize project input
	board.projectInput = textinput.New()
	board.projectInput.Placeholder = "Enter project name..."
	board.projectInput.Width = width - 10
	board.projectInput.CharLimit = 200

	return board
}

func (m board) Init() tea.Cmd {
	return nil
}

func (m board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "esc" {
			if m.showInput {
				m.showInput = false
				m.projectInput.Blur()
				return m, nil
			}
			m.returnModel.height = m.height
			m.returnModel.width = m.width
			return m.returnModel, nil
		}
		if msg.String() == "enter" {
			if m.showInput {
				// Create new project
				projectName := m.projectInput.Value()
				if projectName != "" {
					if err := CreateProjectInCurrentDir(); err != nil {
						// Show error message and stay in input mode
						m.projectInput.SetValue("")
						m.projectInput.Placeholder = "Error: " + err.Error() + " - Try again..."
						return m, nil
					}
					m.showInput = false
					m.showProjects = false
					// Refresh the dashboard with new data
					return dashboard(m.width, m.height, m.styles, m.returnModel), nil
				}
			} else if m.showProjects {
				// Switch project
				data := m.projectList.SelectedItem().(projectItem)
				if err := SetCurrentProject(data.project.Path); err != nil {
					// If there's an error, just stay on current selection
					return m, nil
				}
				m.showProjects = false
				// Refresh the dashboard with new data
				return dashboard(m.width, m.height, m.styles, m.returnModel), nil
			} else {
				// Load request
				data := m.list.SelectedItem().(listItem)
				newModel := NewModel()
				load(data.request, &newModel)
				newModel.width = m.width
				newModel.height = m.height
				newModel.styles = m.styles
				return newModel, nil
			}
		}
		if msg.String() == "n" {
			if !m.showMsg && !m.showProjects && !m.showInput {
				newModel := NewModel()
				newModel.width = m.width
				newModel.height = m.height
				newModel.styles = m.styles
				return newModel, nil
			}
		}
		if msg.String() == "p" && !m.showMsg {
			m.showProjects = !m.showProjects
			m.showInput = false
			return m, nil
		}
		if msg.String() == "a" && !m.showMsg && m.showProjects {
			m.showInput = true
			m.projectInput.SetValue("")
			m.projectInput.Focus()
			return m, nil
		}
		if msg.String() == "r" && !m.showMsg && m.showProjects && !m.showInput {
			// Remove selected project
			data := m.projectList.SelectedItem().(projectItem)
			if err := RemoveProject(data.project.Path); err != nil {
				// Show error message briefly (could implement a temporary message system)
				return m, nil
			}
			// Refresh the dashboard
			return dashboard(m.width, m.height, m.styles, m.returnModel), nil
		}
		if msg.String() == "d" && !m.showMsg && !m.showProjects && !m.showInput {
			// Show message before deleting
			m.showMsg = true
			return m, nil
		}
		// Handle Yes/No input when the message is active
		if m.showMsg {
			if msg.String() == "y" {
				data := m.list.SelectedItem().(listItem)
				err := delete(data.request.Id)
				if err != nil {
					m.showMsg = false
					return m.returnModel, nil
				}

				// Remove the item from the list
				newItems := []list.Item{}
				for i, item := range m.list.Items() {
					if i != m.list.Index() {
						newItems = append(newItems, item)
					}
				}

				// Update the list with new items
				m.list.SetItems(newItems)
				m.showMsg = false
				return m, nil
			}
			if msg.String() == "n" {
				m.showMsg = false // Cancel delete
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-3)
		m.height = msg.Height
		m.width = msg.Width
	}

	var cmd tea.Cmd
	if m.showInput {
		m.projectInput, cmd = m.projectInput.Update(msg)
	} else if m.showProjects {
		m.projectList, cmd = m.projectList.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m board) View() string {
	var footer string
	var body string
	
	if m.showInput {
		footer = m.appBoundaryMessage("Create project in current dir • <ESC> to cancel")
		inputView := "Create New Project:\n\n" + m.projectInput.View()
		body = borderStyle.Width(m.width - 2).Render(inputView)
	} else if m.showProjects {
		footer = m.appBoundaryMessage("↑↓ navigate • enter to select • a to create • r to remove • esc to cancel")
		
		projectHeader := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Width(m.width - 4).
			Render("Switch Project")
		
		projectContent := lipgloss.JoinVertical(lipgloss.Left, projectHeader, "", m.projectList.View())
		body = borderStyle.Width(m.width - 2).Render(projectContent)
	} else {
		// Main view with project selector
		currentProject := GetCurrentProject()
		var projectName string
		if currentProject == "" {
			projectName = "Default"
		} else {
			projectName = filepath.Base(currentProject)
			// Find project name from config
			for _, p := range GetProjects() {
				if p.Path == currentProject {
					projectName = p.Name
					break
				}
			}
		}
		
		// Create project selector header similar to lazygit
		projectSelector := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Width(m.width - 4).
			Render("Project: " + projectName + " (press 'p' to switch)")
		
		footer = m.appBoundaryMessage("Ctrl+c to quit • <ESC> to go back • p for projects")
		if m.showMsg {
			footer = m.appBoundaryMessage("Delete selected item? : (Y/N)")
		}
		
		listView := m.list.View()
		content := lipgloss.JoinVertical(lipgloss.Left, projectSelector, "", listView)
		body = borderStyle.Width(m.width - 2).Render(content)
	}
	
	return m.styles.Base.Render(body + "\n" + footer)
}
