// ui/layout_constants.go
package ui

// LayoutMetrics holds all calculated layout metrics for consistent rendering
type LayoutMetrics struct {
	// Border and padding
	BorderWidth  int // 2 (left + right borders)
	PanelPadding int // From config or default
	TitlePadding int // From config or default

	// Content area calculations
	ContentWidth  int // Available width for content
	ContentHeight int // Available height for content

	// Title constraints
	MaxTitleWidth  int // Maximum title width
	TitleBarHeight int // Always 1 line

	// Scroll indicators
	ScrollIndicatorHeight int // 2 lines (top + bottom)
}

// calculateLayoutMetrics calculates all layout metrics consistently
func (m Model) calculateLayoutMetrics(dimensions PanelDimensions) LayoutMetrics {
	// Get padding from config
	panelPadding := 2
	titlePadding := 1
	if m.Config != nil && m.Config.Display.CompactMode {
		panelPadding = 1
		titlePadding = 0
	}

	// Border takes 2 characters (left + right)
	borderWidth := 2

	// Available content width = panel width - borders - padding
	contentWidth := dimensions.Width - borderWidth - (panelPadding * 2)
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Available content height = panel height - title (1) - borders (2) - padding (2)
	// Borders take 2 lines total (top + bottom)
	contentHeight := dimensions.Height - 1 - 2 - (panelPadding * 2)
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Title width = content width - title padding
	maxTitleWidth := contentWidth - (titlePadding * 2)
	if maxTitleWidth < 5 {
		maxTitleWidth = 5 // Minimum readable width
	}

	return LayoutMetrics{
		BorderWidth:          borderWidth,
		PanelPadding:         panelPadding,
		TitlePadding:         titlePadding,
		ContentWidth:         contentWidth,
		ContentHeight:        contentHeight,
		MaxTitleWidth:        maxTitleWidth,
		TitleBarHeight:       1,
		ScrollIndicatorHeight: 2,
	}
}
