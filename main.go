package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/halftoothed/gostman/cmd"
)

func main() {
	p := tea.NewProgram(cmd.NewStartupModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
