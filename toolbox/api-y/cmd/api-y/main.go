package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"api-y/internal/model"
)

func main() {
	p := tea.NewProgram(
		model.NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
