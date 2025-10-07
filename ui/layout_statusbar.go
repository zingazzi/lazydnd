// ui/layout_statusbar.go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderStatusBar renders the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	text := DefaultStatusBarText

	// Project name with campaign info
	projectNameText := text.ProjectName
	if m.CurrentCampaignName != "" {
		projectNameText += " | 📁 " + m.CurrentCampaignName
		if m.LastAutoSave != "" {
			projectNameText += " (💾 " + m.LastAutoSave + ")"
		}
	}
	projectName := StatusBarTextStyle.Render(projectNameText)

	// Navigation hints
	tabKey := StatusBarKeyStyle.Render(text.TabKey)
	tabText := StatusBarTextStyle.Render(text.TabDesc)

	arrowKeys := StatusBarKeyStyle.Render(text.ArrowKeys)
	arrowText := StatusBarTextStyle.Render(text.ArrowDesc)

	numbersKey := StatusBarKeyStyle.Render(text.NumbersKey)
	numbersText := StatusBarTextStyle.Render(text.NumbersDesc)

	helpKey := StatusBarKeyStyle.Render(text.HelpKey)
	helpText := StatusBarTextStyle.Render(text.HelpDesc)

	quitKey := StatusBarKeyStyle.Render(text.QuitKey)
	quitText := StatusBarTextStyle.Render(text.QuitDesc)

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
	statusBar := StatusBarStyle.
		Width(m.Width).
		Render(statusBarContent)

	return statusBar
}
