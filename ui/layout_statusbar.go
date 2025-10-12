// ui/layout_statusbar.go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderStatusBar renders the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	var result string

	// Render error banner if there's an error
	if m.ErrorVisible && m.ErrorMessage != "" {
		errorBanner := m.Styles.ErrorStyle.
			Width(m.Width).
			Background(lipgloss.Color("#8B0000")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1).
			Render("❌ " + m.ErrorMessage)
		result = errorBanner + "\n"
	}

	text := DefaultStatusBarText

	// Project name with version and campaign info
	projectNameText := text.ProjectName + " " + AppVersion
	if m.CurrentCampaignName != "" {
		projectNameText += " | 📁 " + m.CurrentCampaignName
		if m.LastAutoSave != "" {
			projectNameText += " (💾 " + m.LastAutoSave + ")"
		}
	}
	projectName := m.Styles.StatusBarTextStyle.Render(projectNameText)

	// Navigation hints
	tabKey := m.Styles.StatusBarKeyStyle.Render(text.TabKey)
	tabText := m.Styles.StatusBarTextStyle.Render(text.TabDesc)

	arrowKeys := m.Styles.StatusBarKeyStyle.Render(text.ArrowKeys)
	arrowText := m.Styles.StatusBarTextStyle.Render(text.ArrowDesc)

	numbersKey := m.Styles.StatusBarKeyStyle.Render(text.NumbersKey)
	numbersText := m.Styles.StatusBarTextStyle.Render(text.NumbersDesc)

	helpKey := m.Styles.StatusBarKeyStyle.Render(text.HelpKey)
	helpText := m.Styles.StatusBarTextStyle.Render(text.HelpDesc)

	quitKey := m.Styles.StatusBarKeyStyle.Render(text.QuitKey)
	quitText := m.Styles.StatusBarTextStyle.Render(text.QuitDesc)

	// Build the status bar content
	leftSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectName,
		"  ",
	)

	middleSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		tabKey,
		tabText,
		arrowKeys,
		arrowText,
		numbersKey,
		numbersText,
	)

	rightSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		helpKey,
		helpText,
		quitKey,
		quitText,
	)

	// Calculate spacing to distribute sections across the width
	leftWidth := lipgloss.Width(leftSection)
	middleWidth := lipgloss.Width(middleSection)
	rightWidth := lipgloss.Width(rightSection)

	totalContentWidth := leftWidth + middleWidth + rightWidth
	availableSpace := m.Width - totalContentWidth

	// Distribute space evenly
	spacing1 := availableSpace / 3
	spacing2 := availableSpace / 3
	if spacing1 < 2 {
		spacing1 = 2
	}
	if spacing2 < 2 {
		spacing2 = 2
	}

	statusBarContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSection,
		strings.Repeat(" ", spacing1),
		middleSection,
		strings.Repeat(" ", spacing2),
		rightSection,
	)

	// Apply full-width background style
	statusBar := m.Styles.StatusBarStyle.
		Width(m.Width).
		Render(statusBarContent)

	result += statusBar
	return result
}
