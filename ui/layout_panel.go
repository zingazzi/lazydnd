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
	titleBar := PanelTitleStyle.Render(title)

	borderColor := m.getBorderColor(panelType)
	panelStyle := m.createPanelStyle(borderColor, dimensions)

	panelContent := titleBar + "\n" + content
	return panelStyle.Render(panelContent)
}

// getBorderColor returns the appropriate border color for a panel
func (m Model) getBorderColor(panelType PanelType) string {
	if panelType == m.ActivePanel {
		return "#7D56F4"
	}
	return "#444444"
}

// createPanelStyle creates the lipgloss style for a panel
func (m Model) createPanelStyle(borderColor string, dimensions PanelDimensions) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(1, 2).
		Width(dimensions.Width).
		Height(dimensions.Height)
}
