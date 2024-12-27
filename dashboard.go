package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
	model       *Model
}

func dashboard(width, height int, styles *Styles, returnModel tea.Model, model *Model) board {
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
		model:       model,
	}

	board.list.Title = "List of Requests "

	board.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Back,
		}
	}

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
		if msg.String() == "enter" {
			data := m.list.SelectedItem().(listItem)
			load(data.request, m.model)
			return m.model, nil
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m board) View() string {

	body := borderStyle.Width(m.width - 4).Render(m.list.View())
	return m.styles.Base.Render(body)
}
