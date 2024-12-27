package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const commandsContent = `Tab = Move Around
shift + up = Send Request
ctrl + s = Save Request
ctrl + Arrow Keys = Change Tabs (Body/Param/Header)
shift + left = Open dashboard
ctrl + c = Quit`

func createCommandRows() []table.Row {
	lines := strings.Split(commandsContent, "\n")
	rows := make([]table.Row, len(lines))
	for i, line := range lines {
		cols := strings.Split(line, " = ")
		rows[i] = table.Row{cols[0], cols[1]}
	}
	return rows
}

type help struct {
	width         int
	height        int
	styles        *Styles
	returnModel   tea.Model
	commandsTable table.Model
}

func newHelp(width, height int, styles *Styles, returnModel tea.Model) help {

	columns := []table.Column{
		{Title: "Command", Width: (width / 2) - 13},
		{Title: "Description", Width: (width / 2) + 7},
	}
	commandsTable := table.New(
		table.WithColumns(columns),
		table.WithRows(createCommandRows()),
		table.WithFocused(false),
		table.WithHeight(height-7),
	)
	commandsTable.GotoBottom()
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("7")).
		Bold(false).
		Margin(0, 0, 0, 0).
		Padding(0, 0, 0, 0)
	commandsTable.SetStyles(s)

	return help{
		width:         width,
		height:        height,
		styles:        styles,
		returnModel:   returnModel,
		commandsTable: commandsTable,
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
		m.commandsTable.SetHeight(m.height - 7)
		m.commandsTable.SetWidth(m.width - 2)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.returnModel, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m help) View() string {

	header := m.appBoundaryView("gostman Help Page")
	body := borderStyle.Width(m.width - 2).Height(m.height - 4).
		Render("\n" + m.commandsTable.View() + "\n\n")

	footer := m.appBoundaryView("<ESC> to go back")

	return m.styles.Base.Render(header + "\n" + body + "\n" + footer)
}
