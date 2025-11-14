// ui/help_popup.go
package ui

import (
	"fmt"
	"math"
	"strings"
)

// ========== HELP POPUP RENDERING ==========

// renderHelpPopupOverlay is deprecated - TView handles overlays
func (m Model) renderHelpPopupOverlay(mainView string) string {
	return mainView
}

// buildHelpContent builds the help popup content in 2 columns
func (m Model) buildHelpContent() string {
	var content strings.Builder

	// Title
	content.WriteString(HelpPopupTitle)
	content.WriteString("\n\n")

	// Get help keys for both sections
	commonKeys := CommonNavigationKeys
	panelKeys := GetPanelHelpKeys(m.ActivePanel)

	// Build two-column layout
	leftColumn := m.buildColumnContent("Common Navigation:", commonKeys)
	rightColumn := m.buildColumnContent(fmt.Sprintf("%s Panel Keys:", PanelNames[m.ActivePanel]), panelKeys)

	// Join columns horizontally (simple text layout)
	leftLines := strings.Split(leftColumn, "\n")
	rightLines := strings.Split(rightColumn, "\n")
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		leftLine := ""
		rightLine := ""
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}
		// Pad left column to fixed width
		leftPadded := fmt.Sprintf("%-45s", leftLine)
		content.WriteString(leftPadded + "  " + rightLine + "\n")
	}

	// Footer
	content.WriteString("\n\n")
	content.WriteString(HelpPopupFooter)

	return content.String()
}

// buildColumnContent builds content for a single column
func (m Model) buildColumnContent(sectionTitle string, keys []HelpKey) string {
	var content strings.Builder

	// Section title
	content.WriteString(sectionTitle)
	content.WriteString("\n")

	// Keys
	for _, helpKey := range keys {
		if helpKey.Key == "" && helpKey.Description == "" {
			content.WriteString("\n")
		} else if helpKey.Description == "" {
			// Section header (like "In Edit Mode:")
			content.WriteString(helpKey.Key)
			content.WriteString("\n")
		} else {
			// Regular key binding
			keyPart := fmt.Sprintf("%-12s", helpKey.Key)
			content.WriteString(keyPart + " " + helpKey.Description + "\n")
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
