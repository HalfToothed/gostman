package cmd

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type env struct {
	width       int
	height      int
	styles      *Styles
	returnModel tea.Model
	content     textarea.Model
}

func environment(board Model) env {

	env := env{
		width:       board.width,
		height:      board.height,
		styles:      board.styles,
		returnModel: board,
		content:     newTextarea(),
	}

	env.content.SetValue(loadVariables())
	env.content.SetWidth(env.width - 2)
	env.content.SetHeight(env.height - 5)
	env.content.Focus()
	env.content.Placeholder = `
	{
		"Key":"Value",
	}`

	return env

}

func (en *env) sizeInputs() {
	en.content.SetWidth(en.width - 2)
	en.content.SetHeight(en.height - 5)
}

func (en env) Init() tea.Cmd {
	return nil
}

func (en env) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return en, tea.Quit
		}
		if msg.String() == "esc" {
			return en.returnModel, nil
		}

	case tea.WindowSizeMsg:
		en.width = msg.Width
		en.height = msg.Height
	}

	en.sizeInputs()

	en.content, cmd = en.content.Update(msg)
	cmds = append(cmds, cmd)

	return en, tea.Batch(cmds...)
}

func (en env) View() string {
	footer := en.appBoundaryMessage("Ctrl+c to quit, F2 for help")

	body := borderStyle.Width(en.width - 2).Height(en.height - 4).Render(en.content.View())
	return en.styles.Base.Render(body + "\n" + footer)
}
