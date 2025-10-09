// ui/styles.go
package ui

import "github.com/charmbracelet/lipgloss"

// Styles for the UI components
var (
	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(1, 2)

	InactivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#444444")).
				Padding(1, 2)

	PanelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Margin(1, 0)

	DiceResultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1).
			Margin(1, 0)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Margin(1, 0, 0, 0)

	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1).
			Bold(true)

	StatusBarKeyStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#5A3D9E")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1).
				Bold(true).
				Margin(0, 1)

	StatusBarTextStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FAFAFA"))

	// Help popup styles
	HelpPopupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1A1A1A")).
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(1, 2).
			Margin(1, 2)

	HelpPopupTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				Underline(true).
				Margin(0, 0, 1, 0)

	HelpPopupSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FF00")).
				Margin(1, 0, 0, 0)

	HelpPopupKeyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				Width(15)

	HelpPopupDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC"))

	HelpPopupOverlayStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#000000")).
				Foreground(lipgloss.Color("#000000"))
)
