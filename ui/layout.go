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
		return "\n" + HelpStyle.Render(DefaultInlineHelp())
	}

	return provider(m)
}

// ========== HELP TEXT PROVIDERS ==========

// getDiceRollerHelpText gets help text for the dice roller panel
func getDiceRollerHelpText(m Model) string {
	text := DiceRollerInlineHelp(m.InputMode, m.LastDiceCommand != "")
	return "\n" + HelpStyle.Render(text)
}

// getInitiativeTrackerHelpText gets help text for the initiative tracker panel
func getInitiativeTrackerHelpText(m Model) string {
	text := InitiativeTrackerInlineHelp(m.InitiativeEditMode, m.InitiativeInputMode, m.InitiativeListMode)
	return "\n" + HelpStyle.Render(text)
}

// getSpellsHelpText gets help text for the spells panel
func getSpellsHelpText(m Model) string {
	text := SpellsInlineHelp(m.SpellSearchMode)
	return "\n" + HelpStyle.Render(text)
}

// getMonstersHelpText gets help text for the monsters panel
func getMonstersHelpText(m Model) string {
	text := MonstersInlineHelp(m.MonsterSearchMode)
	return "\n" + HelpStyle.Render(text)
}

// ========== STATUS BAR ==========

// renderStatusBar renders the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	text := DefaultStatusBarText

	// Project name
	projectName := StatusBarTextStyle.Render(text.ProjectName)

	// Navigation hints
	tabKey := StatusBarKeyStyle.Render(text.TabKey)
	tabText := StatusBarTextStyle.Render(text.TabDesc)

	arrowKeys := StatusBarKeyStyle.Render(text.ArrowKeys)
	arrowText := StatusBarTextStyle.Render(text.ArrowDesc)

	numbersKey := StatusBarKeyStyle.Render(text.NumbersKey)
	numbersText := StatusBarTextStyle.Render(text.NumbersDesc)

	helpKey := StatusBarKeyStyle.Render(text.HelpKey)
	helpText := StatusBarTextStyle.Render(text.HelpDesc)

	quitKey := StatusBarKeyStyle.Render(text.QuitKey)
	quitText := StatusBarTextStyle.Render(text.QuitDesc)

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
