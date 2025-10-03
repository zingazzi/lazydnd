// main.go
package main

import (
	"fmt"
	"os"

	"lazydnd/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create initial model
	model := ui.InitialModel()

	// Create program with alt screen
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
