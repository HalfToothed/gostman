package main

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/textarea"
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
