// ui/styles.go
package ui

import (
	"lazydnd/config"

	"github.com/gdamore/tcell/v2"
)

// Styles returns all UI colors based on configuration (TView compatible)
type Styles struct {
	PrimaryColor    tcell.Color
	BorderColor     tcell.Color
	HighlightColor  tcell.Color
	ErrorColor      tcell.Color
	SuccessColor    tcell.Color
	BackgroundColor tcell.Color
	TextColor       tcell.Color
	StatusBarKeyBg  tcell.Color
	CompactMode     bool
	PanelPadding    int
	TitlePadding    int
	InputPadding    int
	ResultPadding   int
	HelpPadding     int
	StatusPadding   int
}

// NewStyles creates styles from configuration
func NewStyles(cfg *config.Config) *Styles {
	// Use default colors if config is nil
	primaryColor := "#7D56F4"
	borderColor := "#444444"
	highlightColor := "#00FF00"
	errorColor := "#FF0000"
	successColor := "#00FF00"
	compactMode := false

	if cfg != nil {
		primaryColor = cfg.Theme.PrimaryColor
		borderColor = cfg.Theme.BorderColor
		highlightColor = cfg.Theme.HighlightColor
		errorColor = cfg.Theme.ErrorColor
		successColor = cfg.Theme.SuccessColor
		compactMode = cfg.Display.CompactMode
	}

	// Adjust padding based on compact mode
	panelPadding := 2
	titlePadding := 1
	inputPadding := 1
	resultPadding := 1
	helpPadding := 3
	statusPadding := 1

	if compactMode {
		panelPadding = 1
		titlePadding = 0
		inputPadding = 0
		resultPadding = 0
		helpPadding = 1
		statusPadding = 0
	}

	// Convert hex colors to tcell.Color
	primaryColorTView := HexToColor(primaryColor)
	borderColorTView := HexToColor(borderColor)
	highlightColorTView := HexToColor(highlightColor)
	errorColorTView := HexToColor(errorColor)
	successColorTView := HexToColor(successColor)

	// Derive a darker shade for status bar key background
	statusBarKeyBg := HexToColor(DarkenColor(primaryColor))

	return &Styles{
		PrimaryColor:    primaryColorTView,
		BorderColor:     borderColorTView,
		HighlightColor:  highlightColorTView,
		ErrorColor:      errorColorTView,
		SuccessColor:    successColorTView,
		BackgroundColor: HexToColor("#1A1A1A"),
		TextColor:       ColorNames.White,
		StatusBarKeyBg:  statusBarKeyBg,
		CompactMode:     compactMode,
		PanelPadding:    panelPadding,
		TitlePadding:    titlePadding,
		InputPadding:    inputPadding,
		ResultPadding:   resultPadding,
		HelpPadding:     helpPadding,
		StatusPadding:   statusPadding,
	}
}

// DarkenColor returns a darker version of the hex color (simple approximation)
func DarkenColor(hex string) string {
	// Simple darkening by reducing hex values
	// This is a basic implementation - for production you'd want proper color manipulation
	if len(hex) == 7 && hex[0] == '#' {
		// For now, return a hardcoded darker purple for the default
		if hex == "#7D56F4" {
			return "#5A3D9E"
		}
		// For other colors, just return a generic dark color
		return "#444444"
	}
	return hex
}
