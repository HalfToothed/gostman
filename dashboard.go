package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 0)

type listItem struct {
	request Request
}

func (i listItem) Title() string       { return i.request.Name }
func (i listItem) Description() string { return i.request.Method }
func (i listItem) FilterValue() string { return i.request.Name }

type board struct {
	width       int
	height      int
	styles      *Styles
	list        list.Model
	returnModel tea.Model
}

func dashboard(width, height int, styles *Styles, returnModel tea.Model) board {
	var savedRequests []Request

	if !checkFileExists(jsonfilePath) {

		file, err := os.ReadFile(jsonfilePath)

		if err != nil {
			panic(err)
		}

		json.Unmarshal(file, &savedRequests)
	}

	// Convert requests to list items
	var items []list.Item
	for _, req := range savedRequests {
		items = append(items, listItem{request: req})
	}

	board := board{
		width:       width,
		height:      height,
		styles:      styles,
		list:        list.New(items, list.NewDefaultDelegate(), width, height-2),
		returnModel: returnModel,
	}

	board.list.Title = "Requests "

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
			return m.returnModel, nil
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m board) View() string {

	header := m.appBoundaryView("Dashboard")
	body := docStyle.Render(m.list.View())

	footer := m.appBoundaryView("<ESC> to go back")

	return m.styles.Base.Render(header + "\n" + body + "\n" + footer)
}
