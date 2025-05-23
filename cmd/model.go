package cmd

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
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
	consoleViewport  viewport.Model
	responseTabs     []string
	activeTab        int
	activeResponseTab int
	response         string
	consoleContent   string
	status           string
	id               string
	focused          int
	fields           []string
	spinner          spinner.Model
	message          string
	loading          bool
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.id = ""
	m.tabs = []string{"Body", "Params", "Headers"}
	m.responseTabs = []string{"Response", "Console"}

	m.nameField = textinput.New()
	m.nameField.Cursor.Blink = false
	m.nameField.SetValue("New Request")
	m.nameField.Placeholder = "Name"
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
	
	cvp := viewport.New(m.width, m.height)
	m.consoleViewport = cvp

	m.focused = 0
	m.fields = []string{"nameField", "methodField", "urlField", "tabContent", "response"}

	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E03535"))
	m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
	m.loading = false

	return m
}

// Init is run once when the program starts
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+h":
			cmd = tea.EnterAltScreen
			help := newHelp(m.width, m.height, m.styles, &m)
			return help, nil
		case "shift+right":
			if m.focused == 3 { // Request tabs
				m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			} else if m.focused == 4 { // Response tabs
				m.activeResponseTab = min(m.activeResponseTab+1, len(m.responseTabs)-1)
			}
			return m, nil
		case "shift+left":
			if m.focused == 3 { // Request tabs
				m.activeTab = max(m.activeTab-1, 0)
			} else if m.focused == 4 { // Response tabs
				m.activeResponseTab = max(m.activeResponseTab-1, 0)
			}
			return m, nil
		case "ctrl+e":
			environment := environment(m)
			return environment, nil
		case "enter":

			if m.focused != 3 {

				m.loading = true
				m.message = m.appBoundaryMessage("Sending Request....")

				m.spinner, cmd = m.spinner.Update(msg)
				cmds = append(cmds, cmd)

				// Perform the async operation in a goroutine
				return m, func() tea.Msg {
					response, status, console := send(m) // Simulate the send function
					formattedResponse := formatJSON(response)
					return responseMsg{
						response: formattedResponse,
						status:   status,
						console:  console,
					}
				}
			}

		case "ctrl+d":
			dashboard := dashboard(m.width, m.height, m.styles, &m)
			return dashboard, nil
		case "ctrl+s":

			m.loading = true
			m.message = m.appBoundaryMessage("Saving Request....")
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)

			// Perform the async save operation in a goroutine
			return m, func() tea.Msg {
				save(m)
				return saveMsg{
					success: true,
					message: "Request Saved Successfully!",
				}
			}

		case "tab":
			m.focused = (m.focused + 1) % 5 // 0-4: name, method, url, tabContent, response
			m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
		}
	}

	m.sizeInputs()

	// Handle custom messages for async tasks
	switch msg := msg.(type) {
	case responseMsg:
		m.response = msg.response
		m.status = msg.status
		m.consoleContent = msg.console
		m.loading = false
		m.message = m.appBoundaryMessage("Request Sent!")

		wrappedContent := wordwrap.String(m.response, m.responseViewport.Width)
		m.responseViewport.SetContent(wrappedContent)
		m.responseViewport.GotoTop()
		
		wrappedConsole := wordwrap.String(m.consoleContent, m.consoleViewport.Width)
		m.consoleViewport.SetContent(wrappedConsole)
		m.consoleViewport.GotoTop()
	case saveMsg:
		m.loading = false
		m.message = m.appBoundaryMessage(msg.message)
	}

	// Update based on focus
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
	case 4:
		// Handle response area - viewport scrolling
		if m.activeResponseTab == 0 {
			m.responseViewport, cmd = m.responseViewport.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.consoleViewport, cmd = m.consoleViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	updateViewport := true
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "up" || keyMsg.String() == "down" {
			if m.focused == 3 || m.focused == 4 {
				updateViewport = false
			}
		}
	}
	if _, ok := msg.(tea.MouseMsg); ok {
		if m.focused == 3 || m.focused == 4 {
			updateViewport = false
		}
	}

	if updateViewport && m.focused != 4 {
		m.responseViewport, cmd = m.responseViewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	// Combine all commands into a single tea.Cmd
	return m, tea.Batch(cmds...)
}
func (m Model) View() string {

	var footer string

	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isActive := i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		renderedTabs = append(renderedTabs, style.Render(t))
	}

	tabContentWidth := int(float64(m.width) * 0.5)

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	tabStyle := borderStyle
	if m.focused == 3 {
		m.tabContent[m.activeTab].Focus()
		tabStyle = focusedBorder
	} else {
		m.tabContent[m.activeTab].Blur()
	}

	tabContent := lipgloss.NewStyle().
		Width(tabContentWidth - 2).
		Height(m.height - 8).
		Render(m.tabContent[m.activeTab].View())

	combined := lipgloss.JoinVertical(lipgloss.Left, tabRow, tabContent)

	// Now, wrap the entire combined layout in a border.
	finalPanel := tabStyle.Render(combined)

	doc.WriteString(finalPanel)
	requestPanel := doc.String()

	m.responseViewport.Height = m.height - 9
	m.responseViewport.Width = m.width - tabContentWidth - 4
	m.consoleViewport.Height = m.height - 9
	m.consoleViewport.Width = m.width - tabContentWidth - 4

	// Render response tabs
	var renderedResponseTabs []string
	for i, t := range m.responseTabs {
		var style lipgloss.Style
		isActive := i == m.activeResponseTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		renderedResponseTabs = append(renderedResponseTabs, style.Render(t))
	}
	
	responseTabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedResponseTabs...)
	
	// Display content based on active response tab
	var responseContent string
	if m.activeResponseTab == 0 { // Response tab
		responseContent = m.responseViewport.View()
	} else { // Console tab
		responseContent = m.consoleViewport.View()
	}
	
	// Combine tabs and status on same line (like left panel)
	statusLine := titleStyle.Render(" Status: ") + headingStyle.Render(m.status)
	
	// Calculate available width for status alignment
	availableWidth := m.width - tabContentWidth - 4 - lipgloss.Width(responseTabRow)
	if availableWidth < 0 {
		availableWidth = 0
	}
	
	// Create tabs row with status aligned to the right
	tabsAndStatus := lipgloss.JoinHorizontal(lipgloss.Top, 
		responseTabRow,
		lipgloss.NewStyle().Width(availableWidth).Align(lipgloss.Right).Render(statusLine))
	
	// Join tabs+status line with content
	responseCombined := lipgloss.JoinVertical(lipgloss.Left, tabsAndStatus, responseContent)
	
	responseStyle := borderStyle
	if m.focused == 4 {
		responseStyle = focusedBorder
	}
	
	responsePanel := responseStyle.Width(m.width - tabContentWidth - 2).Height(m.height - 6).Render(responseCombined)
	mainPanel := lipgloss.JoinHorizontal(lipgloss.Left, requestPanel, responsePanel)

	nameStyle := borderStyle
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
	urlInput := urlStyle.Height(1).Width(m.width - 46).Render(m.urlField.View())

	topPanel := lipgloss.JoinHorizontal(lipgloss.Left, nameInput, methodInput, urlInput)

	body := lipgloss.JoinVertical(lipgloss.Top, topPanel, mainPanel)

	if m.loading {
		spinnerView := m.spinner.View()
		footer = spinnerView + m.appBoundaryMessage(m.message)
	} else {
		footer = m.appBoundaryMessage(m.message)
	}

	return m.styles.Base.Render(body + "\n" + footer)
}

func (m *Model) sizeInputs() {
	for i := range m.tabContent {
		m.tabContent[i].SetWidth(int(float64(m.width)*0.5) - 2)
		m.tabContent[i].SetHeight(m.height - 8)
	}
}
