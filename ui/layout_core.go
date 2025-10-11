// ui/layout_core.go
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	return Model{
		ActivePanel:            DiceRoller,
		DiceInput:              "",
		DiceResult:             "",
		DiceHistory:            []string{},
		DiceCommands:           []string{},
		LastDiceCommand:        "",
		InputMode:              false,
		DiceHistoryMode:        false,
		HistoryIndex:           -1,
		ScrollOffset:           make(map[PanelType]int),
		SpellSearchInput:       "",
		SpellSearchMode:        false,
		SelectedSpell:          nil,
		SpellSuggestions:       []string{},
		SuggestionIndex:        -1,
		InitiativeList:         []InitiativeEntry{},
		InitiativeInput:        "",
		InitiativeInputMode:    false,
		InitiativeInputType:    "",
		SelectedEntry:          -1,
		TempEntry:              InitiativeEntry{},
		CurrentTurn:            -1,
		RoundCounter:           0,
		InitiativeEditMode:     false,
		InitiativeEditType:     "",
		InitiativeListMode:     false,
		MonsterSearchInput:     "",
		MonsterSearchMode:      false,
		SelectedMonster:        nil,
		MonsterSuggestions:     []string{},
		MonsterSuggestionIndex: -1,
		ActiveSpells:           []ActiveSpell{},
		ActiveSpellIndex:       -1,
		ActiveSpellListMode:    false,
		ShowCastSpellPrompt:    false,
		CastSpellInput:         "",
		CastSpellInputMode:     false,
		SpellToCast:            nil,
		ShowHelpPopup:          false,
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

	// Show save popup if active (highest priority)
	if m.ShowSavePopup {
		return m.renderSavePopupOverlay(mainView)
	}

	// Show load popup if active (highest priority)
	if m.ShowLoadPopup {
		return m.renderLoadPopupOverlay(mainView)
	}

	// Show rename popup if active (highest priority)
	if m.ShowRenamePopup {
		return m.renderRenamePopupOverlay(mainView)
	}

	// Show cast spell popup if active (takes priority over other popups)
	if m.ShowCastSpellPrompt {
		return renderCastSpellPopupOverlay(m, mainView)
	}

	// Show saving throw popup if active (takes priority over action popup)
	if m.ShowSavingThrowPopup {
		return m.renderSavingThrowPopupOverlay(mainView)
	}

	// Show action popup if active (takes priority over help popup)
	if m.ShowActionPopup {
		return m.renderActionPopupOverlay(mainView)
	}

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

// arrangeInGrid arranges the four panels in a 2x2 grid
func (m Model) arrangeInGrid(panelViews []string) string {
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[0], panelViews[1])
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[2], panelViews[3])

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

// renderSavePopupOverlay renders the save popup over the main view
func (m Model) renderSavePopupOverlay(mainView string) string {
	popup := RenderSavePopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}

// renderLoadPopupOverlay renders the load popup over the main view
func (m Model) renderLoadPopupOverlay(mainView string) string {
	popup := RenderLoadPopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}

// renderRenamePopupOverlay renders the rename popup over the main view
func (m Model) renderRenamePopupOverlay(mainView string) string {
	popup := RenderRenamePopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}
