// main.go
package main

import (
	"fmt"
	"os"

	"lazydnd/config"
	"lazydnd/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load config, using defaults: %v\n", err)
		cfg = config.Default()
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Warning: Invalid config, using defaults: %v\n", err)
		cfg = config.Default()
		// Try to save valid config
		_ = cfg.Save()
	}

	// Ensure all configured directories exist
	if err := cfg.EnsureDirectoriesExist(); err != nil {
		fmt.Printf("Warning: Failed to create directories: %v\n", err)
	}

	// Create initial model with configuration
	model := ui.InitialModel()
	model.Config = cfg
	model.Styles = ui.NewStyles(cfg)

	// Create program with alt screen
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
