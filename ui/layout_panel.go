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
	// Step 1: Calculate layout metrics
	metrics := m.calculateLayoutMetrics(dimensions)

	// Step 2: Get raw content
	rawContent := m.getPanelContent(panelType)

	// Step 3: Process content (wrap/truncate based on panel type)
	processedContent := m.processPanelContent(rawContent, panelType, metrics)

	// Step 4: Apply scrolling
	scrolledContent := m.applyScrolling(processedContent, panelType, dimensions, metrics)

	// Step 5: Create title bar (constrained)
	titleBar := m.createTitleBar(panelNumber, metrics)

	// Step 6: Combine and style
	return m.assemblePanel(titleBar, scrolledContent, panelType, dimensions)
}

// processPanelContent processes content based on panel type
func (m Model) processPanelContent(content string, panelType PanelType, metrics LayoutMetrics) string {
	// Determine if panel should wrap or truncate
	shouldWrap := m.shouldWrapPanel(panelType)

	if shouldWrap {
		return m.wrapContent(content, metrics.ContentWidth)
	}

	return m.truncateContent(content, metrics.ContentWidth)
}

// shouldWrapPanel determines if a panel should wrap text
func (m Model) shouldWrapPanel(panelType PanelType) bool {
	// Only wrap for panels that benefit from it
	wrapPanels := map[PanelType]bool{
		DiceRoller:        true,
		InitiativeTracker: true,
		Notes:             true,
	}
	return wrapPanels[panelType]
}

// wrapContent wraps content to fit width
func (m Model) wrapContent(content string, maxWidth int) string {
	lines := strings.Split(content, "\n")
	var wrappedLines []string

	for _, line := range lines {
		wrapped := wrapText(line, maxWidth)
		wrappedLines = append(wrappedLines, wrapped...)
	}

	return strings.Join(wrappedLines, "\n")
}

// truncateContent truncates content to fit width
func (m Model) truncateContent(content string, maxWidth int) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = truncateText(line, maxWidth)
	}
	return strings.Join(lines, "\n")
}

// createTitleBar creates a constrained title bar
func (m Model) createTitleBar(panelNumber int, metrics LayoutMetrics) string {
	title := fmt.Sprintf(" %d. %s ", panelNumber, PanelNames[panelNumber-1])

	// Constrain title width
	constrainedTitle := truncateText(title, metrics.MaxTitleWidth)

	return m.Styles.PanelTitleStyle.Render(constrainedTitle)
}

// assemblePanel combines title and content with proper styling
func (m Model) assemblePanel(titleBar, content string, panelType PanelType, dimensions PanelDimensions) string {
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
