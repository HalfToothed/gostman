package main

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

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

func formatJSON(input string) string {
	var rawData interface{}

	err := json.Unmarshal([]byte(input), &rawData)
	if err != nil {
		return input
	}

	prettyJSON, err := json.MarshalIndent(rawData, "", "  ")
	if err != nil {
		return input
	}

	return string(prettyJSON)
}

func (m Model) appBoundaryView(textItems []string) string {
	builder := strings.Builder{}
	builder.WriteString(m.styles.HeaderDecoration.Render("+--"))
	for i, textItem := range textItems {
		if i > 0 {
			builder.WriteString(m.styles.HeaderDecoration.Render("/////"))
		}
		builder.WriteString(m.styles.HeaderText.Render(textItem))
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, builder.String(), lipgloss.WithWhitespaceChars("/"), lipgloss.WithWhitespaceForeground(indigo))
}
