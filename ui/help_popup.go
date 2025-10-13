// ui/help_popup.go
package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ========== HELP POPUP RENDERING ==========

// renderHelpPopupOverlay renders the help popup over the main view
func (m Model) renderHelpPopupOverlay(mainView string) string {
	helpContent := m.buildHelpContent()

	// Create the popup
	popup := m.Styles.HelpPopupStyle.Render(helpContent)

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

// buildHelpContent builds the help popup content in 2 columns
func (m Model) buildHelpContent() string {
	var content strings.Builder

	// Title
	title := m.Styles.HelpPopupTitleStyle.Render(HelpPopupTitle)
	content.WriteString(title)
	content.WriteString("\n\n")

	// Get help keys for both sections
	commonKeys := CommonNavigationKeys
	panelKeys := GetPanelHelpKeys(m.ActivePanel)

	// Build two-column layout
	leftColumn := m.buildColumnContent("Common Navigation:", commonKeys)
	rightColumn := m.buildColumnContent(fmt.Sprintf("%s Panel Keys:", PanelNames[m.ActivePanel]), panelKeys)

	// Calculate column width (half of available width minus spacing)
	columnWidth := 45 // Fixed width for each column

	// Style columns
	columnStyle := lipgloss.NewStyle().Width(columnWidth)
	leftStyled := columnStyle.Render(leftColumn)
	rightStyled := columnStyle.Render(rightColumn)

	// Join columns horizontally
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, "  ", rightStyled)
	content.WriteString(columns)

	// Footer
	content.WriteString("\n\n")
	footer := m.Styles.HelpPopupDescStyle.Render(HelpPopupFooter)
	content.WriteString(footer)

	return content.String()
}

// buildColumnContent builds content for a single column
func (m Model) buildColumnContent(sectionTitle string, keys []HelpKey) string {
	var content strings.Builder

	// Section title
	title := m.Styles.HelpPopupSectionStyle.Render(sectionTitle)
	content.WriteString(title)
	content.WriteString("\n")

	// Keys
	for _, helpKey := range keys {
		if helpKey.Key == "" && helpKey.Description == "" {
			content.WriteString("\n")
		} else if helpKey.Description == "" {
			// Section header (like "In Edit Mode:")
			content.WriteString(m.Styles.HelpPopupSectionStyle.Render(helpKey.Key))
			content.WriteString("\n")
		} else {
			// Regular key binding
			keyPart := m.Styles.HelpPopupKeyStyle.Render(helpKey.Key)
			descPart := m.Styles.HelpPopupDescStyle.Render(helpKey.Description)
			content.WriteString(keyPart + " " + descPart + "\n")
		}
	}

	return content.String()
}

// splitKeysForColumns splits help keys into two roughly equal columns
func splitKeysForColumns(keys []HelpKey) ([]HelpKey, []HelpKey) {
	if len(keys) <= 10 {
		return keys, []HelpKey{}
	}

	// Split roughly in half, but try to split at empty lines for better layout
	midpoint := int(math.Ceil(float64(len(keys)) / 2.0))

	// Find nearest empty line to split on
	for i := midpoint - 2; i <= midpoint + 2 && i < len(keys); i++ {
		if i > 0 && i < len(keys) && keys[i].Key == "" && keys[i].Description == "" {
			return keys[:i], keys[i+1:]
		}
	}

	// No good split point found, just split at midpoint
	return keys[:midpoint], keys[midpoint:]
}
