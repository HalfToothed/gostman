package cmd

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
	returnModel *Model
	showMsg     bool
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
			m.returnModel.height = m.height
			m.returnModel.width = m.width
			return m.returnModel, nil
		}
		if msg.String() == "enter" {
			data := m.list.SelectedItem().(listItem)
			newModel := NewModel()
			load(data.request, &newModel)
			newModel.width = m.width
			newModel.height = m.height
			newModel.styles = m.styles
			return newModel, nil
		}
		if msg.String() == "n" {
			if !m.showMsg {
				newModel := NewModel()
				newModel.width = m.width
				newModel.height = m.height
				newModel.styles = m.styles
				return newModel, nil
			}
		}
		if msg.String() == "d" && !m.showMsg {
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m board) View() string {
	footer := m.appBoundaryMessage("Ctrl+c to quit, <ESC> to go back")
	if m.showMsg {
		footer = m.appBoundaryMessage("Delete selected item? : (Y/N)")
	}

	body := borderStyle.Width(m.width - 2).Render(m.list.View())
	return m.styles.Base.Render(body + "\n" + footer)
}
