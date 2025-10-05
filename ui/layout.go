// ui/layout.go
package ui

import (
	"fmt"
	"strings"

	"lazydnd/panels"

	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	return Model{
		ActivePanel:         DiceRoller,
		DiceInput:           "",
		DiceResult:          "",
		DiceHistory:         []string{},
		LastDiceCommand:     "",
		InputMode:           false,
		ScrollOffset:        make(map[PanelType]int),
		SpellSearchInput:    "",
		SpellSearchMode:     false,
		SelectedSpell:       nil,
		SpellSuggestions:    []string{},
		SuggestionIndex:     -1,
		InitiativeList:      []InitiativeEntry{},
		InitiativeInput:     "",
		InitiativeInputMode: false,
		InitiativeInputType: "",
		SelectedEntry:       -1,
		TempEntry:           InitiativeEntry{},
		InitiativeEditMode:  false,
		InitiativeEditType:  "",
		InitiativeListMode:  false,
		MonsterSearchInput:    "",
		MonsterSearchMode:     false,
		SelectedMonster:       nil,
		MonsterSuggestions:    []string{},
		MonsterSuggestionIndex: -1,
		ShowHelpPopup:         false,
	}
}

// View renders the main application view with 2x2 panel layout and status bar
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	dimensions := m.calculatePanelDimensions()
	panelViews := m.renderAllPanels(dimensions)
	grid := m.arrangeInGrid(panelViews)
	statusBar := m.renderStatusBar()

	mainView := lipgloss.JoinVertical(lipgloss.Left, grid, statusBar)

	// Show help popup if active
	if m.ShowHelpPopup {
		return m.renderHelpPopupOverlay(mainView)
	}

	return mainView
}

// PanelDimensions holds calculated panel dimensions
type PanelDimensions struct {
	Width  int
	Height int
}

// calculatePanelDimensions calculates panel dimensions to fill the screen
func (m Model) calculatePanelDimensions() PanelDimensions {
	// Reserve 2 lines for status bar at the bottom
	availableHeight := m.Height - 2

	return PanelDimensions{
		Width:  (m.Width - 6) / 2,
		Height: (availableHeight - 4) / 2,
	}
}

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

// applyScrolling applies scrolling logic to panel content
func (m Model) applyScrolling(content string, panelType PanelType, dimensions PanelDimensions) string {
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

// arrangeInGrid arranges the four panels in a 2x2 grid
func (m Model) arrangeInGrid(panelViews []string) string {
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[0], panelViews[1])
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[2], panelViews[3])

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

// PanelContentProvider defines a function type for getting panel content
type PanelContentProvider func(Model) string

// panelContentProviders maps panel types to their content provider functions
var panelContentProviders = map[PanelType]PanelContentProvider{
	DiceRoller:        getDiceRollerContent,
	InitiativeTracker: getInitiativeTrackerContent,
	Spells:           getSpellsContent,
	Monsters:         getMonstersContent,
}

// getPanelContent returns the content for a specific panel
func (m Model) getPanelContent(panelType PanelType) string {
	// Get base content from provider
	provider, exists := panelContentProviders[panelType]
	if !exists {
		return "Unknown panel type"
	}

	content := provider(m)

	// Add help text for active panel
	if panelType == m.ActivePanel {
		content += m.getHelpText(panelType)
	}

	return content
}

// ========== PANEL CONTENT PROVIDERS ==========

// getDiceRollerContent gets content for the dice roller panel
func getDiceRollerContent(m Model) string {
	return panels.GetDiceRollerContent(
		m.DiceInput,
		m.DiceResult,
		m.DiceHistory,
		m.LastDiceCommand,
		m.InputMode,
		m.ActivePanel == DiceRoller,
	)
}

// getInitiativeTrackerContent gets content for the initiative tracker panel
func getInitiativeTrackerContent(m Model) string {
	return panels.GetInitiativeTrackerContent(
		m.InitiativeList,
		m.InitiativeInput,
		m.InitiativeInputMode,
		m.InitiativeInputType,
		m.SelectedEntry,
		m.ActivePanel == InitiativeTracker,
		m.InitiativeListMode,
		m.InitiativeEditMode,
		m.InitiativeEditType,
	)
}

// getSpellsContent gets content for the spells panel
func getSpellsContent(m Model) string {
	return panels.GetSpellsContent(
		m.SpellSearchInput,
		m.SelectedSpell,
		m.SpellSuggestions,
		m.SuggestionIndex,
		m.SpellSearchMode,
		m.ActivePanel == Spells,
	)
}

// getMonstersContent gets content for the monsters panel
func getMonstersContent(m Model) string {
	return panels.GetMonstersContent(
		m.MonsterSearchInput,
		m.SelectedMonster,
		m.MonsterSuggestions,
		m.MonsterSuggestionIndex,
		m.MonsterSearchMode,
		m.ActivePanel == Monsters,
	)
}

// HelpTextProvider defines a function type for getting help text
type HelpTextProvider func(Model) string

// helpTextProviders maps panel types to their help text provider functions
var helpTextProviders = map[PanelType]HelpTextProvider{
	DiceRoller:        getDiceRollerHelpText,
	InitiativeTracker: getInitiativeTrackerHelpText,
	Spells:           getSpellsHelpText,
	Monsters:         getMonstersHelpText,
}

// getHelpText returns context-sensitive help text
func (m Model) getHelpText(panelType PanelType) string {
	provider, exists := helpTextProviders[panelType]
	if !exists {
		return "\n" + HelpStyle.Render("↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
	}

	return provider(m)
}

// ========== HELP TEXT PROVIDERS ==========

// getDiceRollerHelpText gets help text for the dice roller panel
func getDiceRollerHelpText(m Model) string {
		if m.InputMode {
			return "\n" + HelpStyle.Render("Enter: roll • Esc: cancel • F1-F4: switch panels")
	}

			if m.LastDiceCommand != "" {
				return "\n" + HelpStyle.Render("Enter: input dice • r: reroll • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
	}

				return "\n" + HelpStyle.Render("Enter: input dice • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
			}

// getInitiativeTrackerHelpText gets help text for the initiative tracker panel
func getInitiativeTrackerHelpText(m Model) string {
		if m.InitiativeEditMode {
			return "\n" + HelpStyle.Render("Enter: confirm edit • Esc: cancel • F1-F4: switch panels")
	}

	if m.InitiativeInputMode {
			return "\n" + HelpStyle.Render("Enter: confirm • Esc: cancel • F1-F4: switch panels")
	}

	if m.InitiativeListMode {
			return "\n" + HelpStyle.Render("↑↓: select • i: edit initiative • h: edit HP • d: delete • Esc: exit edit • F1-F4: switch")
	}

			return "\n" + HelpStyle.Render("p: add player • m: add monster • e: edit list • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
		}

// getSpellsHelpText gets help text for the spells panel
func getSpellsHelpText(m Model) string {
		if m.SpellSearchMode {
			return "\n" + HelpStyle.Render("Enter: select spell • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch")
	}

			return "\n" + HelpStyle.Render("Enter: search spells • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
		}

// getMonstersHelpText gets help text for the monsters panel
func getMonstersHelpText(m Model) string {
		if m.MonsterSearchMode {
			return "\n" + HelpStyle.Render("Enter: select monster • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch")
	}

			return "\n" + HelpStyle.Render("Enter: search monsters • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
		}

// ========== STATUS BAR ==========

// renderStatusBar renders the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	// Project name
	projectName := StatusBarTextStyle.Render("🎲 LazyDnD")

	// Navigation hints
	tabKey := StatusBarKeyStyle.Render("Tab")
	tabText := StatusBarTextStyle.Render("Switch Panel")

	arrowKeys := StatusBarKeyStyle.Render("↑↓←→")
	arrowText := StatusBarTextStyle.Render("Navigate")

	numbersKey := StatusBarKeyStyle.Render("1-4")
	numbersText := StatusBarTextStyle.Render("Quick Switch")

	helpKey := StatusBarKeyStyle.Render("?")
	helpText := StatusBarTextStyle.Render("Help")

	quitKey := StatusBarKeyStyle.Render("q")
	quitText := StatusBarTextStyle.Render("Quit")

	// Build the status bar content
	leftSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectName,
		"  ",
	)

	middleSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		tabKey,
		tabText,
		arrowKeys,
		arrowText,
		numbersKey,
		numbersText,
	)

	rightSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		helpKey,
		helpText,
		quitKey,
		quitText,
	)

	// Calculate spacing to distribute sections across the width
	leftWidth := lipgloss.Width(leftSection)
	middleWidth := lipgloss.Width(middleSection)
	rightWidth := lipgloss.Width(rightSection)

	totalContentWidth := leftWidth + middleWidth + rightWidth
	availableSpace := m.Width - totalContentWidth

	// Distribute space evenly
	spacing1 := availableSpace / 3
	spacing2 := availableSpace / 3
	if spacing1 < 2 {
		spacing1 = 2
	}
	if spacing2 < 2 {
		spacing2 = 2
	}

	statusBarContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSection,
		strings.Repeat(" ", spacing1),
		middleSection,
		strings.Repeat(" ", spacing2),
		rightSection,
	)

	// Apply full-width background style
	statusBar := StatusBarStyle.
		Width(m.Width).
		Render(statusBarContent)

	return statusBar
}

// ========== HELP POPUP ==========

// renderHelpPopupOverlay renders the help popup over the main view
func (m Model) renderHelpPopupOverlay(mainView string) string {
	helpContent := m.buildHelpContent()

	// Create the popup
	popup := HelpPopupStyle.Render(helpContent)

	// Calculate position to center the popup
	popupWidth := lipgloss.Width(popup)
	popupHeight := lipgloss.Height(popup)

	// Center horizontally and vertically
	leftPadding := (m.Width - popupWidth) / 2
	topPadding := (m.Height - popupHeight) / 2

	if leftPadding < 0 {
		leftPadding = 0
	}
	if topPadding < 0 {
		topPadding = 0
	}

	// Place popup over the main view
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333333")),
	)
}

// buildHelpContent builds the help popup content
func (m Model) buildHelpContent() string {
	var content strings.Builder

	// Title
	title := HelpPopupTitleStyle.Render("🎲 LazyDnD - Help")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Common navigation section
	commonSection := HelpPopupSectionStyle.Render("Common Navigation:")
	content.WriteString(commonSection)
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine("Tab", "Switch to next panel"))
	content.WriteString(m.formatHelpLine("1-4", "Quick switch to panel"))
	content.WriteString(m.formatHelpLine("F1-F4", "Switch to specific panel"))
	content.WriteString(m.formatHelpLine("↑/↓", "Scroll panel content"))
	content.WriteString(m.formatHelpLine("Esc", "Cancel/Exit current mode"))
	content.WriteString(m.formatHelpLine("?", "Toggle this help"))
	content.WriteString(m.formatHelpLine("q", "Quit application"))

	// Panel-specific section
	panelSection := HelpPopupSectionStyle.Render(fmt.Sprintf("\n%s Panel Keys:", PanelNames[m.ActivePanel]))
	content.WriteString(panelSection)
	content.WriteString("\n")
	content.WriteString(m.getPanelSpecificHelp(m.ActivePanel))

	// Footer
	content.WriteString("\n")
	footer := HelpPopupDescStyle.Render("Press ? or Esc to close this help")
	content.WriteString(footer)

	return content.String()
}

// formatHelpLine formats a single help line with key and description
func (m Model) formatHelpLine(key, description string) string {
	keyPart := HelpPopupKeyStyle.Render(key)
	descPart := HelpPopupDescStyle.Render(description)
	return keyPart + " " + descPart + "\n"
}

// getPanelSpecificHelp returns panel-specific help content
func (m Model) getPanelSpecificHelp(panelType PanelType) string {
	var content strings.Builder

	switch panelType {
	case DiceRoller:
		content.WriteString(m.formatHelpLine("Enter", "Start/confirm dice input"))
		content.WriteString(m.formatHelpLine("r", "Reroll last command"))
		content.WriteString(m.formatHelpLine("Examples:", ""))
		content.WriteString(m.formatHelpLine("  2d6", "Roll 2 six-sided dice"))
		content.WriteString(m.formatHelpLine("  1d20+5", "Roll d20 with +5 modifier"))
		content.WriteString(m.formatHelpLine("  adv", "Roll with advantage"))
		content.WriteString(m.formatHelpLine("  dis", "Roll with disadvantage"))

	case InitiativeTracker:
		content.WriteString(m.formatHelpLine("p", "Add player to initiative"))
		content.WriteString(m.formatHelpLine("m", "Add monster to initiative"))
		content.WriteString(m.formatHelpLine("e", "Enter edit mode"))
		content.WriteString(m.formatHelpLine("", ""))
		content.WriteString(m.formatHelpLine("In Edit Mode:", ""))
		content.WriteString(m.formatHelpLine("  ↑/↓", "Select entry"))
		content.WriteString(m.formatHelpLine("  i", "Edit initiative value"))
		content.WriteString(m.formatHelpLine("  h", "Edit HP (monsters only)"))
		content.WriteString(m.formatHelpLine("  d", "Delete entry"))

	case Spells:
		content.WriteString(m.formatHelpLine("Enter", "Start spell search"))
		content.WriteString(m.formatHelpLine("", ""))
		content.WriteString(m.formatHelpLine("In Search Mode:", ""))
		content.WriteString(m.formatHelpLine("  Type", "Search for spells"))
		content.WriteString(m.formatHelpLine("  ↑/↓", "Navigate suggestions"))
		content.WriteString(m.formatHelpLine("  Enter", "Select spell"))
		content.WriteString(m.formatHelpLine("  Backspace", "Delete character"))

	case Monsters:
		content.WriteString(m.formatHelpLine("Enter", "Start monster search"))
		content.WriteString(m.formatHelpLine("a", "Add to initiative tracker"))
		content.WriteString(m.formatHelpLine("", ""))
		content.WriteString(m.formatHelpLine("In Search Mode:", ""))
		content.WriteString(m.formatHelpLine("  Type", "Search for monsters"))
		content.WriteString(m.formatHelpLine("  ↑/↓", "Navigate suggestions"))
		content.WriteString(m.formatHelpLine("  Enter", "Select monster"))
		content.WriteString(m.formatHelpLine("  Backspace", "Delete character"))
	}

	return content.String()
}
