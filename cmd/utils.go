package cmd

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

type responseMsg struct {
	response string
	status   string
}

type saveMsg struct {
	success bool
	message string
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
	t.FocusedStyle.Base = lipgloss.NewStyle()
	t.BlurredStyle.Base = lipgloss.NewStyle()
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()
	return t
}

func formatJSON(input string) string {
	var rawData interface{}

	err := json.Unmarshal([]byte(input), &rawData)
	if err != nil {
		return input
	}

	f := colorjson.NewFormatter()
	f.Indent = 2

	prettyJSON, err := f.Marshal(rawData)
	if err != nil {
		return input
	}

	return string(prettyJSON)
}

func initializeVariables() string {
	return ""
}

func createHeaders() string {

	var rawData interface{}

	headers := `
	{
		"Content-Type":"application/json",
		"Accept":"*/*",
		"Accept-Encoding":"gzip, deflate, br",
		"Connection":"keep-alive"
	}`

	err := json.Unmarshal([]byte(headers), &rawData)
	if err != nil {
		return headers
	}

	prettyJSON, err := json.MarshalIndent(rawData, "", "  ")
	if err != nil {
		return headers
	}

	return string(prettyJSON)
}

type keymap struct {
	Create key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Create: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "create"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
