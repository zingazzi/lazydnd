// ui/layout_panel.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderAllPanels renders all six panels with dynamic dimensions
func (m Model) renderAllPanels(dimensionsMap map[PanelType]PanelDimensions) []string {
	panelViews := make([]string, 6)

	for i := 0; i < 6; i++ {
		panelType := PanelType(i)
		dimensions := dimensionsMap[panelType]
		panelViews[i] = m.renderSinglePanel(panelType, dimensions, i+1)
	}

	return panelViews
}

// renderSinglePanel renders a single panel with all its content and styling
func (m Model) renderSinglePanel(panelType PanelType, dimensions PanelDimensions, panelNumber int) string {
	// Validate dimensions before rendering
	dimensions = m.validateDimensions(dimensions)

	content := m.getPanelContent(panelType)
	// Sanitize content to remove trailing newlines and handle empty content
	content = m.sanitizeContent(content)

	scrolledContent := m.applyScrolling(content, panelType, dimensions)
	styledContent := m.stylePanel(scrolledContent, panelType, dimensions, panelNumber)

	return styledContent
}

// validateDimensions ensures dimensions are within reasonable bounds
func (m Model) validateDimensions(dimensions PanelDimensions) PanelDimensions {
	// Minimum dimensions to ensure panels are visible
	const minWidth = 10
	const minHeight = 5

	if dimensions.Width < minWidth {
		dimensions.Width = minWidth
	}
	if dimensions.Height < minHeight {
		dimensions.Height = minHeight
	}

	// Maximum dimensions to prevent overflow
	if dimensions.Width > m.Width {
		dimensions.Width = m.Width
	}
	if dimensions.Height > m.Height {
		dimensions.Height = m.Height
	}

	return dimensions
}

// sanitizeContent removes trailing newlines and ensures content is properly formatted
func (m Model) sanitizeContent(content string) string {
	// Remove trailing newlines
	content = strings.TrimRight(content, "\n")

	// If content is completely empty, return a single space to maintain layout
	if content == "" {
		return " "
	}

	return content
}

// stylePanel applies styling to a panel
func (m Model) stylePanel(content string, panelType PanelType, dimensions PanelDimensions, panelNumber int) string {
	title := fmt.Sprintf(" %d. %s ", panelNumber, PanelNames[panelNumber-1])
	titleBar := m.Styles.PanelTitleStyle.Render(title)

	panelStyle := m.getPanelStyle(panelType, dimensions)

	panelContent := titleBar + "\n" + content
	return panelStyle.Render(panelContent)
}

// getPanelStyle returns the appropriate style for a panel based on whether it's active
func (m Model) getPanelStyle(panelType PanelType, dimensions PanelDimensions) lipgloss.Style {
	var baseStyle lipgloss.Style
	if panelType == m.ActivePanel {
		baseStyle = m.Styles.ActivePanelStyle
	} else {
		baseStyle = m.Styles.InactivePanelStyle
	}

	// Apply fixed dimensions to prevent resize
	return baseStyle.Width(dimensions.Width).Height(dimensions.Height)
}
