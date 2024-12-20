package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
		Padding(0, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
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
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type Model struct {
	lg         *lipgloss.Renderer
	styles     *Styles
	width      int
	height     int
	inputField textinput.Model
	tabs       []string
	tabContent []textarea.Model
	activeTab  int
	focused    string
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.tabs = []string{"Body", "Params", "Authorization", "Headers"}

	m.inputField = textinput.New()
	m.inputField.Placeholder = "URL"
	m.inputField.Focus()
	m.inputField.Cursor.Blink = false
	m.focused = "inputField"

	// Initialize tab contents
	for _, tab := range m.tabs {
		ta := newTextarea()
		ta.Placeholder = fmt.Sprintf("Write something in %s...", tab)
		ta.Focus()
		ta.Cursor.Blink = false
		m.tabContent = append(m.tabContent, ta)
	}

	return m
}

// Init is run once when the program starts
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "right":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil

		case "tab": // Example key to switch focus
			if m.focused == "inputField" {
				m.focused = "tabContent"
			} else {
				m.focused = "inputField"
			}
		}
	}

	m.sizeInputs()

	// Update based on focus
	if m.focused == "inputField" {
		m.inputField, cmd = m.inputField.Update(msg)
	} else if m.focused == "tabContent" {
		m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg)
	}

	return m, cmd
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
	tabContent := windowStyle.Width(m.width - 60).Height(m.height - 9).Render(m.tabContent[m.activeTab].View())
	doc.WriteString(tabContent)

	cmdDesc := doc.String()
	resPanel := borderStyle.Width(m.width - 102).Height(m.height - 7).Render(titleStyle.Render("Response"))
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left, cmdDesc, resPanel)

	methodInput := borderStyle.Width(15).Height(1).Render(m.inputField.View())
	cmdInput := borderStyle.Height(1).Width(m.width - 19).Render(m.inputField.View())

	topPanel := lipgloss.JoinHorizontal(lipgloss.Left, methodInput, cmdInput)

	body := lipgloss.JoinVertical(lipgloss.Top, topPanel, mainPanel)

	return m.styles.Base.Render("gostman"+"\n"+body+"\n"+"Ctrl+h to view help,", "Ctrl+c to quit")

}

func (m *Model) sizeInputs() {
	for i := range m.tabContent {
		m.tabContent[i].SetWidth(m.width - 65)
		m.tabContent[i].SetHeight(m.height - 16)
	}
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

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()
	return t
}
