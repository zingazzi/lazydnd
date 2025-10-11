// ui/layout_scroll.go
package ui

import (
	"strings"
)

// applyScrolling applies scrolling logic to panel content
func (m Model) applyScrolling(content string, panelType PanelType, dimensions PanelDimensions) string {
	// Skip panel-level scrolling for search/filter modes that handle their own scrolling
	if panelType == Monsters && (m.MonsterSearchMode || m.MonsterCRFilterMode) {
		return content // Search content already has its own scroll handling
	}
	if panelType == Spells && m.SpellSearchMode {
		return content // Search content already has its own scroll handling
	}

	contentLines := strings.Split(content, "\n")
	scrollOffset := m.ScrollOffset[panelType]

	heights := m.calculateContentHeights(dimensions, panelType)
	scrollOffset = m.normalizeScrollOffset(scrollOffset, contentLines, heights.Content, panelType)

	visibleLines := m.extractVisibleLines(contentLines, scrollOffset, heights.Content)
	finalContent := m.addScrollIndicators(visibleLines, scrollOffset, contentLines, heights, panelType)

	return m.truncateLongLines(finalContent, dimensions.Width)
}

// ContentHeights holds calculated content heights
type ContentHeights struct {
	Available int
	Content   int
}

// calculateContentHeights calculates available and content heights for a panel
func (m Model) calculateContentHeights(dimensions PanelDimensions, panelType PanelType) ContentHeights {
	// Account for: title (1) + borders (2) + padding (2) = 5 lines
	availableHeight := dimensions.Height - 5
	if availableHeight < 1 {
		availableHeight = 1 // Minimum 1 line of content
	}

	contentHeight := availableHeight

	// Reserve space for scroll indicators on active panels
	if panelType == m.ActivePanel {
		contentHeight = availableHeight - 2
		if contentHeight < 1 {
			contentHeight = 1
		}
	}

	// Extra buffer for spell panel to prevent resizing
	if panelType == Spells {
		contentHeight = contentHeight - 2
		if contentHeight < 1 {
			contentHeight = 1
		}
	}

	return ContentHeights{
		Available: availableHeight,
		Content:   contentHeight,
	}
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
func (m Model) addScrollIndicators(visibleLines []string, scrollOffset int, contentLines []string, heights ContentHeights, panelType PanelType) string {
	var finalLines []string

	// Add top indicator space
	if panelType == m.ActivePanel && scrollOffset > 0 && len(contentLines) > heights.Available {
		finalLines = append(finalLines, "▲ (more above)")
	} else {
		finalLines = append(finalLines, "") // Empty line to maintain spacing
	}

	// Add visible content
	finalLines = append(finalLines, visibleLines...)

	// Add bottom indicator space
	if panelType == m.ActivePanel && scrollOffset+len(visibleLines) < len(contentLines) && len(contentLines) > heights.Available {
		finalLines = append(finalLines, "▼ (more below)")
	} else {
		finalLines = append(finalLines, "") // Empty line to maintain spacing
	}

	// Ensure total matches available height exactly
	if len(finalLines) > heights.Available {
		finalLines = finalLines[:heights.Available]
	}
	for len(finalLines) < heights.Available {
		finalLines = append(finalLines, "")
	}

	return strings.Join(finalLines, "\n")
}

// truncateLongLines truncates lines that are too long to prevent wrapping
func (m Model) truncateLongLines(content string, panelWidth int) string {
	finalLines := strings.Split(content, "\n")
	maxLineWidth := panelWidth - 6 // Account for padding and borders

	for i, line := range finalLines {
		if len(line) > maxLineWidth {
			finalLines[i] = line[:maxLineWidth-3] + "..."
		}
	}

	return strings.Join(finalLines, "\n")
}
