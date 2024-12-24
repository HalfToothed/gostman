package main

import tea "github.com/charmbracelet/bubbletea"

type board struct {
	width       int
	height      int
	styles      *Styles
	returnModel tea.Model
}

func dashboard(width, height int, styles *Styles, returnModel tea.Model) board {
	return board{
		width:       width,
		height:      height,
		styles:      styles,
		returnModel: returnModel,
	}
}

func (m board) Init() tea.Cmd {
	return nil
}

func (m board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return nil, nil
}

func (m board) View() string {

	return ""
}
