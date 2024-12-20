package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type Model struct {
	lg         *lipgloss.Renderer
	styles     *Styles
	width      int
	height     int
	inputField textinput.Model
	tabs       []string
	tabContent []string
	activeTab  int
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.inputField = textinput.New()
	m.inputField.Focus()
	m.tabs = []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	m.tabContent = []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}

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
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {

	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	doc.WriteString(row)
	doc.WriteString("\n")
	//doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabContent[m.activeTab]))

	tabContent := windowStyle.Width(m.width - 44).Height(m.height - 9).Render(m.tabContent[m.activeTab])
	doc.WriteString(tabContent)

	cmdDesc := doc.String()
	resPanel := borderStyle.Width(40).Height(m.height - 7).Render(titleStyle.Render("Response"))
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left, cmdDesc, resPanel)
	cmdInput := borderStyle.Height(1).Width(m.width - 2).Render(m.inputField.View())
	body := lipgloss.JoinVertical(lipgloss.Top, cmdInput, mainPanel)

	return m.styles.Base.Render("gostman"+"\n"+body+"\n"+"Ctrl+h to view help,", "Ctrl+c to quit")

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
