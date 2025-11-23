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
	if panelType == Spells && (m.SpellSearchMode || m.SpellLevelFilterMode) {
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
	// Account for: title (1) + borders (2) = 3 lines
	// Note: padding is horizontal only (Padding(0, panelPadding)), not vertical
	availableHeight := dimensions.Height - 3
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
// Returns the normalized offset value without mutating state
func (m Model) normalizeScrollOffset(scrollOffset int, contentLines []string, contentHeight int, panelType PanelType) int {
	maxScroll := len(contentLines) - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	if scrollOffset > maxScroll {
		scrollOffset = maxScroll
	}
	if scrollOffset < 0 {
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
	// Defensive check: ensure heights are valid
	if heights.Available < 1 {
		heights.Available = 1
	}
	if heights.Content < 1 {
		heights.Content = 1
	}

	// Defensive check: ensure visibleLines doesn't exceed content height
	if len(visibleLines) > heights.Content {
		visibleLines = visibleLines[:heights.Content]
	}

	var finalLines []string

	// Calculate how many lines we can use for content (accounting for indicators)
	needsTopIndicator := panelType == m.ActivePanel && scrollOffset > 0 && len(contentLines) > heights.Available
	needsBottomIndicator := panelType == m.ActivePanel && scrollOffset+len(visibleLines) < len(contentLines) && len(contentLines) > heights.Available

	// Reserve space for indicators
	indicatorCount := 0
	if needsTopIndicator {
		indicatorCount++
	}
	if needsBottomIndicator {
		indicatorCount++
	}

	// Calculate available space for content
	contentSpace := heights.Available - indicatorCount
	if contentSpace < 1 {
		contentSpace = 1
	}

	// Add top indicator if needed
	if needsTopIndicator {
		finalLines = append(finalLines, "▲ (more above)")
	}

	// Add visible content, ensuring we don't exceed available space
	contentToAdd := visibleLines
	if len(contentToAdd) > contentSpace {
		contentToAdd = contentToAdd[:contentSpace]
	}
	finalLines = append(finalLines, contentToAdd...)

	// Add bottom indicator if needed
	if needsBottomIndicator {
		finalLines = append(finalLines, "▼ (more below)")
	}

	// Ensure total matches available height exactly - pad or truncate as needed
	if len(finalLines) > heights.Available {
		// Truncate to exact height
		finalLines = finalLines[:heights.Available]
	}
	for len(finalLines) < heights.Available {
		finalLines = append(finalLines, "")
	}

	// Final defensive check: ensure we have exactly the right number of lines
	if len(finalLines) != heights.Available {
		// Force to exact height
		if len(finalLines) > heights.Available {
			finalLines = finalLines[:heights.Available]
		} else {
			for len(finalLines) < heights.Available {
				finalLines = append(finalLines, "")
			}
		}
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
