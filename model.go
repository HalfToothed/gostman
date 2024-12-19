package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	borderStyle = lipgloss.NewStyle().
			Padding(0, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	allowedDir    = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	disallowedDir = lipgloss.NewStyle().Foreground(lipgloss.Color("173"))
	commandStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
)

type Model struct {
	lg          *lipgloss.Renderer
	styles      *Styles
	width       int
	height      int
	inputField  textinput.Model
	descContent string
	contentPort *viewport.Model
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.inputField = textinput.New()
	m.inputField.Focus()

	return m
}

// Init is run once when the program starts
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {

	vpHeight := 18
	if m.height > 7 {
		vpHeight = m.height - 7
	}
	if m.contentPort == nil {
		vp := viewport.New(m.width-26, vpHeight)
		m.contentPort = &vp
		m.contentPort.SetContent(wordwrap.String(m.descContent, m.width-26))
		m.contentPort.GotoBottom()
	}

	cmdDesc := borderStyle.Width(m.width - 44).Height(m.height - 7).Render(m.contentPort.View())
	resPanel := borderStyle.Width(40).Height(m.height - 7).Render(titleStyle.Render("Response"))
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left, cmdDesc, resPanel)
	cmdInput := borderStyle.Height(1).Width(m.width - 2).Render(m.inputField.View())
	body := lipgloss.JoinVertical(lipgloss.Top, cmdInput, mainPanel)

	return m.styles.Base.Render("gostman"+"\n"+body+"\n"+"Ctrl+h to view help,", "Ctrl+c to quit")

}
