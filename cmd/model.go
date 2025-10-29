package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
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
	activeTab        int
	response         string
	status           string
	id               string
	focused          int
	fields           []string
	spinner          spinner.Model
	message          string
	loading          bool
	apiResponse      string

	// Histories for auto-completion (derived from saved requests)
	nameHistory []string
	urlHistory  []string
	// Autocomplete state
	autoCompleteIndex int
	autoCompleteField string // "name" or "url"

	// HTTP method selection options/index
	methodOptions []string
	methodIndex   int
}

func NewModel() Model {

	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.id = ""
	m.tabs = []string{"Body", "Params", "Headers"}

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
	m.methodField.CharLimit = 7
	m.methodField.Cursor.Blink = false

	// Initialize HTTP method list and default selection
	m.methodOptions = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	m.methodIndex = 0
	m.methodField.SetValue(m.methodOptions[m.methodIndex])

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

	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E03535"))
	m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
	m.loading = false

	// Build histories from saved requests (for auto-complete)
	if saved := getSavedData(); len(saved.Requests) > 0 {
		dedup := func(list []string) []string {
			m := map[string]struct{}{}
			out := make([]string, 0, len(list))
			for _, s := range list {
				if s == "" {
					continue
				}
				if _, ok := m[s]; !ok {
					m[s] = struct{}{}
					out = append(out, s)
				}
			}
			return out
		}
		names := make([]string, 0, len(saved.Requests))
		urls := make([]string, 0, len(saved.Requests))
		for _, r := range saved.Requests {
			names = append(names, r.Name)
			urls = append(urls, r.URL)
		}
		m.nameHistory = dedup(names)
		m.urlHistory = dedup(urls)
	}

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
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "shift+left":
			m.activeTab = max(m.activeTab-1, 0)
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
					response, status := send(m) // Simulate the send function
					formattedResponse := formatJSON(response)
					return responseMsg{
						response: formattedResponse,
						status:   status,
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
			m.focused = (m.focused + 1) % len(m.fields)
			m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
			// reset autocomplete cycling when focus changes
			m.autoCompleteIndex = 0
			m.autoCompleteField = ""

		case "shift+tab":
			m.focused = (m.focused - 1) % len(m.fields)
			m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
			// reset autocomplete cycling when focus changes
			m.autoCompleteIndex = 0
			m.autoCompleteField = ""

		case "ctrl+f":
			if m.focused == 0 {
				m = m.applyAutoComplete("name")
				return m, nil
			}
			if m.focused == 2 {
				m = m.applyAutoComplete("url")
				return m, nil
			}

		case "up":
			if m.focused == 1 { // method field: cycle up
				m.methodIndex = (m.methodIndex - 1 + len(m.methodOptions)) % len(m.methodOptions)
				m.methodField.SetValue(m.methodOptions[m.methodIndex])
				return m, nil
			}
		case "down":
			if m.focused == 1 { // method field: cycle down
				m.methodIndex = (m.methodIndex + 1) % len(m.methodOptions)
				m.methodField.SetValue(m.methodOptions[m.methodIndex])
				return m, nil
			}

		case "ctrl+y":
			txt := stripANSI(m.apiResponse)
			clipboard.WriteAll(txt)
			m.message = m.appBoundaryMessage("Response Copied ...")
			return m, nil
		}

	}

	m.sizeInputs()

	// Handle custom messages for async tasks
	switch msg := msg.(type) {
	case responseMsg:
		m.response = msg.response
		m.status = msg.status
		m.loading = false
		m.message = m.appBoundaryMessage("Request Sent!")

		// Update in-session histories from current inputs
		if v := strings.TrimSpace(m.nameField.Value()); v != "" {
			exists := false
			for _, s := range m.nameHistory {
				if s == v {
					exists = true
					break
				}
			}
			if !exists {
				m.nameHistory = append(m.nameHistory, v)
			}
		}
		if v := strings.TrimSpace(m.urlField.Value()); v != "" {
			exists := false
			for _, s := range m.urlHistory {
				if s == v {
					exists = true
					break
				}
			}
			if !exists {
				m.urlHistory = append(m.urlHistory, v)
			}
		}

		wrappedContent := wordwrap.String(m.response, m.responseViewport.Width)
		m.apiResponse = wrappedContent
		m.responseViewport.SetContent(wrappedContent)
		m.responseViewport.GotoTop()
	case saveMsg:
		m.loading = false
		m.message = m.appBoundaryMessage(msg.message)
		// Also update histories on save
		if v := strings.TrimSpace(m.nameField.Value()); v != "" {
			exists := false
			for _, s := range m.nameHistory {
				if s == v {
					exists = true
					break
				}
			}
			if !exists {
				m.nameHistory = append(m.nameHistory, v)
			}
		}
		if v := strings.TrimSpace(m.urlField.Value()); v != "" {
			exists := false
			for _, s := range m.urlHistory {
				if s == v {
					exists = true
					break
				}
			}
			if !exists {
				m.urlHistory = append(m.urlHistory, v)
			}
		}
	}

	// Update based on focus
	switch m.focused {
	case 0:
		m.nameField, cmd = m.nameField.Update(msg)
		// Update inline suggestion hint for name
		m.recomputeInlineSuggestion("name")
		cmds = append(cmds, cmd)
	case 1:
		m.methodField, cmd = m.methodField.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		m.urlField, cmd = m.urlField.Update(msg)
		// Update inline suggestion hint for URL
		m.recomputeInlineSuggestion("url")
		cmds = append(cmds, cmd)
	case 3:
		// Update the active tab in the tabContent array
		m.tabContent[m.activeTab], cmd = m.tabContent[m.activeTab].Update(msg)
		cmds = append(cmds, cmd)
	}

	updateViewport := true
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "up" || keyMsg.String() == "down" {
			if m.focused == 3 {
				updateViewport = false
			}
		}
	}
	if _, ok := msg.(tea.MouseMsg); ok {
		if m.focused == 3 {
			updateViewport = false
		}
	}

	if updateViewport {
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

	m.responseViewport.Height = m.height - 7
	m.responseViewport.Width = m.width - tabContentWidth - 2

	responsePanel := borderStyle.Width(m.width - tabContentWidth - 2).Height(m.height - 6).Render(titleStyle.Render(" Response: ") + headingStyle.Render(m.status) + "\n" + m.responseViewport.View())
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

// applyAutoComplete tries to complete the current field from history. Repeated triggers cycle through matches.
func (m Model) applyAutoComplete(field string) Model {
	pick := m.getBestSuggestion(field)
	if pick == "" {
		return m
	}
	if field == "name" {
		m.nameField.SetValue(pick)
	} else {
		m.urlField.SetValue(pick)
	}
	m.autoCompleteField = field
	// don't cycle indices with inline accept; keep index at 0
	return m
}

// getBestSuggestion returns the first history entry that matches the current prefix for the given field.
func (m Model) getBestSuggestion(field string) string {
	var current string
	var history []string
	switch field {
	case "name":
		current = strings.TrimSpace(m.nameField.Value())
		history = m.nameHistory
	case "url":
		current = strings.TrimSpace(m.urlField.Value())
		history = m.urlHistory
	default:
		return ""
	}
	if len(history) == 0 {
		return ""
	}
	lower := strings.ToLower(current)
	for _, h := range history {
		if lower != "" && strings.HasPrefix(strings.ToLower(h), lower) && h != current {
			return h
		}
	}
	return ""
}

// recomputeInlineSuggestion updates the footer hint to show the best suggestion for the field.
func (m *Model) recomputeInlineSuggestion(field string) {
	suggest := m.getBestSuggestion(field)
	if suggest != "" && !m.loading {
		m.message = m.appBoundaryView(fmt.Sprintf("\u21B3 Suggest: %s  (Ctrl+F to accept)", suggest))
	} else if !m.loading {
		m.message = m.appBoundaryView("Ctrl+c to quit, Ctrl+h for help")
	}
}
