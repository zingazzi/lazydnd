// ui/help_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ========== HELP POPUP RENDERING ==========

// renderHelpPopupOverlay renders the help popup over the main view
func (m Model) renderHelpPopupOverlay(mainView string) string {
	helpContent := m.buildHelpContent()

	// Create the popup
	popup := HelpPopupStyle.Render(helpContent)

	// Place popup centered over the main view
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333333")),
	)
}

// buildHelpContent builds the help popup content
func (m Model) buildHelpContent() string {
	var content strings.Builder

	// Title
	title := HelpPopupTitleStyle.Render(HelpPopupTitle)
	content.WriteString(title)
	content.WriteString("\n\n")

	// Common navigation section
	content.WriteString(m.renderCommonNavigationSection())

	// Panel-specific section
	content.WriteString(m.renderPanelSpecificSection())

	// Footer
	content.WriteString("\n")
	footer := HelpPopupDescStyle.Render(HelpPopupFooter)
	content.WriteString(footer)

	return content.String()
}

// renderCommonNavigationSection renders the common navigation keys section
func (m Model) renderCommonNavigationSection() string {
	var content strings.Builder

	sectionTitle := HelpPopupSectionStyle.Render("Common Navigation:")
	content.WriteString(sectionTitle)
	content.WriteString("\n")

	for _, helpKey := range CommonNavigationKeys {
		content.WriteString(m.formatHelpLine(helpKey.Key, helpKey.Description))
	}

	return content.String()
}

// renderPanelSpecificSection renders the panel-specific keys section
func (m Model) renderPanelSpecificSection() string {
	var content strings.Builder

	sectionTitle := HelpPopupSectionStyle.Render(
		fmt.Sprintf("\n%s Panel Keys:", PanelNames[m.ActivePanel]),
	)
	content.WriteString(sectionTitle)
	content.WriteString("\n")

	panelKeys := GetPanelHelpKeys(m.ActivePanel)
	for _, helpKey := range panelKeys {
		content.WriteString(m.formatHelpLine(helpKey.Key, helpKey.Description))
	}

	return content.String()
}

// formatHelpLine formats a single help line with key and description
func (m Model) formatHelpLine(key, description string) string {
	keyPart := HelpPopupKeyStyle.Render(key)
	descPart := HelpPopupDescStyle.Render(description)
	return keyPart + " " + descPart + "\n"
}
