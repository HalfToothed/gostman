package main

import (
	"github.com/charmbracelet/bubbles/table"
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
	columns := []table.Column{
		{Title: "Command", Width: (width / 2) - 13},
		{Title: "Description", Width: (width / 2) + 7},
	}
	commandsTable := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(height-12),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		BorderBottom(false).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0")).
		Bold(false).
		Margin(0, 0, 0, 0).
		Padding(0, 0, 0, 0)
	commandsTable.SetStyles(s)

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

func (m help) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, m.styles.HeaderText.Render("+-- "+text), lipgloss.WithWhitespaceChars("/"), lipgloss.WithWhitespaceForeground(indigo))
}
