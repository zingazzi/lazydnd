// ui/styles.go
package ui

import (
	"lazydnd/config"

	"github.com/charmbracelet/lipgloss"
)

// Styles returns all UI styles based on configuration
type Styles struct {
	ActivePanelStyle      lipgloss.Style
	InactivePanelStyle    lipgloss.Style
	PanelTitleStyle       lipgloss.Style
	InputStyle            lipgloss.Style
	DiceResultStyle       lipgloss.Style
	HelpStyle             lipgloss.Style
	StatusBarStyle        lipgloss.Style
	StatusBarKeyStyle     lipgloss.Style
	StatusBarTextStyle    lipgloss.Style
	HelpPopupStyle        lipgloss.Style
	HelpPopupTitleStyle   lipgloss.Style
	HelpPopupSectionStyle lipgloss.Style
	HelpPopupKeyStyle     lipgloss.Style
	HelpPopupDescStyle    lipgloss.Style
	HelpPopupOverlayStyle lipgloss.Style
	ErrorStyle            lipgloss.Style
	SuccessStyle          lipgloss.Style
	HighlightStyle        lipgloss.Style
}

// NewStyles creates styles from configuration
func NewStyles(cfg *config.Config) *Styles {
	// Use default colors if config is nil
	primaryColor := "#7D56F4"
	borderColor := "#444444"
	highlightColor := "#00FF00"
	errorColor := "#FF0000"
	successColor := "#00FF00"

	if cfg != nil {
		primaryColor = cfg.Theme.PrimaryColor
		borderColor = cfg.Theme.BorderColor
		highlightColor = cfg.Theme.HighlightColor
		errorColor = cfg.Theme.ErrorColor
		successColor = cfg.Theme.SuccessColor
	}

	// Derive a darker shade for status bar key background
	statusBarKeyBg := darkenColor(primaryColor)

	return &Styles{
		ActivePanelStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryColor)).
			Padding(0, 2),

		InactivePanelStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(borderColor)).
			Padding(0, 2),

		PanelTitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color(primaryColor)).
			Padding(0, 1),

		InputStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(primaryColor)).
			Padding(0, 1),

		DiceResultStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(highlightColor)).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1),

		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),

		StatusBarStyle: lipgloss.NewStyle().
			Background(lipgloss.Color(primaryColor)).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1).
			Bold(true),

		StatusBarKeyStyle: lipgloss.NewStyle().
			Background(lipgloss.Color(statusBarKeyBg)).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true).
			Margin(0, 1),

		StatusBarTextStyle: lipgloss.NewStyle().
			Background(lipgloss.Color(primaryColor)).
			Foreground(lipgloss.Color("#FAFAFA")),

		HelpPopupStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryColor)).
			Background(lipgloss.Color("#1A1A1A")).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(1, 2).
			Margin(1, 2),

		HelpPopupTitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(primaryColor)).
			Underline(true).
			Margin(0, 0, 1, 0),

		HelpPopupSectionStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(highlightColor)).
			Margin(1, 0, 0, 0),

		HelpPopupKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(primaryColor)).
			Width(15),

		HelpPopupDescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")),

		HelpPopupOverlayStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("#000000")).
			Foreground(lipgloss.Color("#000000")),

		ErrorStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(errorColor)),

		SuccessStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(successColor)),

		HighlightStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(highlightColor)),
	}
}

// darkenColor returns a darker version of the hex color (simple approximation)
func darkenColor(hex string) string {
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
