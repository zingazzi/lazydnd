// ui/layout_panel.go
package ui

import (
	"fmt"
	"strings"
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

// createTitleBar creates a constrained title bar (deprecated - TView handles this)
func (m Model) createTitleBar(panelNumber int, metrics LayoutMetrics) string {
	title := fmt.Sprintf(" %d. %s ", panelNumber, PanelNames[panelNumber-1])
	// Constrain title width
	return truncateText(title, metrics.MaxTitleWidth)
}

// assemblePanel combines title and content (deprecated - TView handles this)
func (m Model) assemblePanel(titleBar, content string, panelType PanelType, dimensions PanelDimensions) string {
	return titleBar + "\n" + content
}

// getPanelStyle is deprecated - TView handles panel styling
func (m Model) getPanelStyle(panelType PanelType, dimensions PanelDimensions) interface{} {
	// Return nil - TView handles styling
	return nil
}
