package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const shortcutContent = `
	Tab = move around
	Shift + Up = Send Request
	Shift + Left/Right = Change Tabs- Body/Params/Auth/Headers
	Ctrl + c = Quit gostman
	`

type help struct {
	width       int
	height      int
	styles      *Styles
	returnModel tea.Model
}

func newHelp(width, height int, styles *Styles, returnModel tea.Model) help {
	return help{
		width:       width,
		height:      height,
		styles:      styles,
		returnModel: returnModel,
	}
}

// Init is run once when the program starts
func (m help) Init() tea.Cmd {
	return nil
}

// Update handles game state changes
func (m help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.returnModel, nil
		}
	}
	return m, nil
}

func (m help) View() string {
	title := "Help"
	var titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 2).
		MarginLeft((m.width / 2) - len(title) - 1)

	header := m.appBoundaryView("gostman Help Page")
	body := borderStyle.Width(m.width - 2).Height(m.height - 4).
		Render(
			titleStyle.Render(title) + "\n" +
				headingStyle.Render("Shortcuts") + "\n" +
				shortcutContent + "\n\n")

	footer := m.appBoundaryView("<ESC> to go back")

	return m.styles.Base.Render(header + "\n" + body + "\n" + footer)
}
