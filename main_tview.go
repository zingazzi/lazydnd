// +build tview

// main_tview.go - TView implementation (alternative to main.go)
// Build with: go build -tags tview -o lazydnd main_tview.go
package main

import (
	"flag"
	"fmt"
	"os"

	"lazydnd/config"
	"lazydnd/ui"
	tviewapp "lazydnd/ui/tview"
)

func main() {
	// Initialize debug logger
	if err := ui.InitDebugLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not initialize debug logger: %v\n", err)
	}
	defer ui.CloseDebugLogger()

	// Parse command-line flags
	versionFlag := flag.Bool("version", false, "Print version and exit")
	debugFlag := flag.Bool("debug", false, "Enable debug logging to ~/.config/lazydnd/debug.log")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println("LazyDnD " + ui.AppVersion + " (TView)")
		os.Exit(0)
	}

	// Enable debug mode if flag is set
	if *debugFlag {
		ui.EnableDebugMode()
		fmt.Println("🐛 Debug mode enabled - Logging to ~/.config/lazydnd/debug.log")
	}

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

	// Set initial dimensions (TView will handle resize)
	model.Width = 120
	model.Height = 40

	// Create TView application
	app := tviewapp.NewApp(&model)

	// Run the application
	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
