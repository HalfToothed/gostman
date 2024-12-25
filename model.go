package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	lg               *lipgloss.Renderer
	styles           *Styles
	width            int
	height           int
	nameField        textinput.Model
	urlField         textinput.Model
	methodField      textinput.Model
	tabs             []string
	tabContent       []textarea.Model
	responseViewport viewport.Model
	activeTab        int
	response         string
	status           string

	focused int
	fields  []string
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.tabs = []string{"Body", "Params", "Headers"}

	m.nameField = textinput.New()
	m.nameField.Cursor.Blink = false
	m.nameField.SetValue("New Request")
	m.nameField.Focus()
	m.nameField.CharLimit = 22

	m.urlField = textinput.New()
	m.urlField.Placeholder = "URL"
	m.urlField.Focus()
	m.urlField.Cursor.Blink = false

	m.methodField = textinput.New()
	m.methodField.Placeholder = "METHOD"
	m.methodField.Focus()
	m.methodField.CharLimit = 6
	m.methodField.Cursor.Blink = false

	// Initialize tab contents
	for range m.tabs {
		ta := newTextarea()
		ta.Cursor.Blink = false
		m.tabContent = append(m.tabContent, ta)
	}

	m.tabContent[1].Placeholder = `
	write Query Params in key-value format

{
	"key":"value"
}`

	m.tabContent[2].SetValue(createHeaders())

	vp := viewport.New(m.width, m.height)
	m.responseViewport = vp

	m.focused = 0
	m.fields = []string{"nameField", "methodField", "urlField", "tabContnet"}

	return m
}

// Init is run once when the program starts
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "alt+`":
			help := newHelp(m.width, m.height, m.styles, m)
			return help, nil
		case "ctrl+right":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "ctrl+left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "shift+up":
			m.response, m.status = send(m)
			formattedResponse := formatJSON(m.response)
			wrappedContent := wordwrap.String(formattedResponse, m.responseViewport.Width)
			m.responseViewport.SetContent(wrappedContent)
			m.responseViewport.GotoTop()
			return m, nil
		case "shift+left":
			dashboard := dashboard(m.width, m.height, m.styles, m)
			return dashboard, nil

		case "ctrl+s":
			save(m)

		case "tab":
			m.focused = (m.focused + 1) % len(m.fields)
		case "shift+tab":
			m.focused = (m.focused - 1 + len(m.fields)) % len(m.fields)

		case "up":
			m.responseViewport.LineUp(1) // Scroll up
		case "down":
			m.responseViewport.LineDown(1) // Scroll down
		case "pgup":
			m.responseViewport.ViewUp() // Scroll up a full page
		case "pgdown":
			m.responseViewport.ViewDown() // Scroll down a full page
		}

	}

	m.sizeInputs()

	// Update based on focus
	var cmd tea.Cmd
	switch m.focused {
	case 0:
		m.nameField, cmd = m.nameField.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		m.methodField, cmd = m.methodField.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		m.urlField, cmd = m.urlField.Update(msg)
		cmds = append(cmds, cmd)
	case 3:
		// Update the active tab in the tabContent array
		m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg)
		cmds = append(cmds, cmd)
	}

	m.responseViewport, cmd = m.responseViewport.Update(msg)
	cmds = append(cmds, cmd)

	// Combine all commands into a single tea.Cmd
	return m, tea.Batch(cmds...)
}
func (m Model) View() string {

	footer := m.appBoundaryView("Ctrl+c to quit, alt+` to help")

	focusedBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("205"))

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

	tabContentWidth := int(float64(m.width) * 0.5)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	tabContent := windowStyle.Width(tabContentWidth).Height(m.height - 9).Render(m.tabContent[m.activeTab].View())
	doc.WriteString(tabContent)
	requestPanel := doc.String()

	if m.focused == 3 {
		m.tabContent[m.activeTab].Focus()
	} else {
		m.tabContent[m.activeTab].Blur()
	}

	m.responseViewport.Height = m.height - 7
	m.responseViewport.Width = m.width - tabContentWidth - 2

	responsePanel := borderStyle.Width(m.width - tabContentWidth - 2).Height(m.height - 6).Render(titleStyle.Render(" Response: ") + headingStyle.Render(m.status) + "\n" + m.responseViewport.View())
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left, requestPanel, responsePanel)

	nameStyle := darkStyle
	if m.focused == 0 {
		m.nameField.Focus()
		nameStyle = focusedBorder

	} else {
		m.nameField.Blur()
	}
	nameInput := nameStyle.Width(25).Height(1).Render(m.nameField.View())

	// Render the Method input field
	methodStyle := borderStyle
	if m.focused == 1 {
		m.methodField.Focus()
		methodStyle = focusedBorder
	} else {
		m.methodField.Blur()
	}
	methodInput := methodStyle.Width(15).Height(1).Render(m.methodField.View())

	// Render the URL input field
	urlStyle := borderStyle
	if m.focused == 2 {
		m.urlField.Focus()
		urlStyle = focusedBorder
	} else {
		m.urlField.Blur()
	}
	urlInput := urlStyle.Height(1).Width(m.width - 50).Render(m.urlField.View())

	topPanel := lipgloss.JoinHorizontal(lipgloss.Left, nameInput, methodInput, urlInput)

	body := lipgloss.JoinVertical(lipgloss.Top, topPanel, mainPanel)

	return m.styles.Base.Render(body + "\n" + footer)
}

func (m *Model) sizeInputs() {
	for i := range m.tabContent {
		m.tabContent[i].SetWidth(int(float64(m.width) * 0.5))
		m.tabContent[i].SetHeight(m.height - 9)
	}
}
