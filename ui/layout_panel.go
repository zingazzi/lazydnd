// ui/layout_panel.go
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// renderAllPanels renders all four panels
func (m Model) renderAllPanels(dimensions PanelDimensions) []string {
	panelViews := make([]string, 4)

	for i := 0; i < 4; i++ {
		panelType := PanelType(i)
		panelViews[i] = m.renderSinglePanel(panelType, dimensions, i+1)
	}

	return panelViews
}

// renderSinglePanel renders a single panel with all its content and styling
func (m Model) renderSinglePanel(panelType PanelType, dimensions PanelDimensions, panelNumber int) string {
	content := m.getPanelContent(panelType)
	scrolledContent := m.applyScrolling(content, panelType, dimensions)
	styledContent := m.stylePanel(scrolledContent, panelType, dimensions, panelNumber)

	return styledContent
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
