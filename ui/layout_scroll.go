// ui/layout_scroll.go
package ui

import (
	"strings"
)

// applyScrolling applies scrolling logic to panel content
func (m Model) applyScrolling(content string, panelType PanelType, dimensions PanelDimensions, metrics LayoutMetrics) string {
	// Skip panel-level scrolling for search/filter modes that handle their own scrolling
	if panelType == Monsters && (m.MonsterSearchMode || m.MonsterCRFilterMode) {
		return content // Search content already has its own scroll handling
	}
	if panelType == Spells && (m.SpellSearchMode || m.SpellLevelFilterMode) {
		return content // Search content already has its own scroll handling
	}

	contentLines := strings.Split(content, "\n")
	scrollOffset := m.ScrollOffset[panelType]

	// Calculate available height (accounting for scroll indicators)
	availableHeight := metrics.ContentHeight
	if panelType == m.ActivePanel {
		availableHeight -= metrics.ScrollIndicatorHeight
	}
	if availableHeight < 1 {
		availableHeight = 1
	}

	scrollOffset = m.normalizeScrollOffset(scrollOffset, contentLines, availableHeight, panelType)

	visibleLines := m.extractVisibleLines(contentLines, scrollOffset, availableHeight)
	finalContent := m.addScrollIndicators(visibleLines, scrollOffset, contentLines, availableHeight, panelType, metrics)

	// Content should already be processed (wrapped/truncated), but ensure it fits
	return m.truncateContent(finalContent, metrics.ContentWidth)
}

// ContentHeights holds calculated content heights (deprecated - use LayoutMetrics instead)
// Kept for backward compatibility during migration
type ContentHeights struct {
	Available int
	Content   int
}

// normalizeScrollOffset ensures scroll offset is within valid bounds
func (m Model) normalizeScrollOffset(scrollOffset int, contentLines []string, contentHeight int, panelType PanelType) int {
	maxScroll := len(contentLines) - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if scrollOffset > maxScroll {
		m.ScrollOffset[panelType] = maxScroll
		scrollOffset = maxScroll
	}
	if scrollOffset < 0 {
		m.ScrollOffset[panelType] = 0
		scrollOffset = 0
	}

	return scrollOffset
}

// extractVisibleLines extracts the lines that should be visible based on scroll offset
func (m Model) extractVisibleLines(contentLines []string, scrollOffset, contentHeight int) []string {
	if len(contentLines) <= contentHeight {
		return contentLines
	}

	endIndex := scrollOffset + contentHeight
	if endIndex > len(contentLines) {
		endIndex = len(contentLines)
	}

	visibleLines := contentLines[scrollOffset:endIndex]

	// Ensure we don't exceed the available height
	if len(visibleLines) > contentHeight {
		visibleLines = visibleLines[:contentHeight]
	}

	return visibleLines
}

// addScrollIndicators adds scroll indicators and maintains consistent spacing
func (m Model) addScrollIndicators(visibleLines []string, scrollOffset int, contentLines []string, availableHeight int, panelType PanelType, metrics LayoutMetrics) string {
	var finalLines []string

	// Calculate total height needed
	totalHeight := metrics.ContentHeight
	if panelType == m.ActivePanel {
		totalHeight = availableHeight + metrics.ScrollIndicatorHeight
	}

	// Add top indicator
	if panelType == m.ActivePanel && scrollOffset > 0 && len(contentLines) > availableHeight {
		indicator := truncateText("▲ (more above)", metrics.ContentWidth)
		finalLines = append(finalLines, indicator)
	} else {
		finalLines = append(finalLines, "") // Empty line to maintain spacing
	}

	// Add visible content
	finalLines = append(finalLines, visibleLines...)

	// Add bottom indicator
	if panelType == m.ActivePanel && scrollOffset+len(visibleLines) < len(contentLines) && len(contentLines) > availableHeight {
		indicator := truncateText("▼ (more below)", metrics.ContentWidth)
		finalLines = append(finalLines, indicator)
	} else {
		finalLines = append(finalLines, "") // Empty line to maintain spacing
	}

	// Ensure exact height match
	if len(finalLines) > totalHeight {
		finalLines = finalLines[:totalHeight]
	}
	for len(finalLines) < totalHeight {
		finalLines = append(finalLines, "")
	}

	return strings.Join(finalLines, "\n")
}

// truncateLongLines truncates lines that are too long to prevent wrapping
// Deprecated: Use truncateContent() instead. Kept for backward compatibility.
func (m Model) truncateLongLines(content string, panelWidth int) string {
	// Estimate content width (this is approximate - proper usage should pass metrics)
	estimatedContentWidth := panelWidth - 6
	if estimatedContentWidth < 1 {
		estimatedContentWidth = 1
	}
	return m.truncateContent(content, estimatedContentWidth)
}
