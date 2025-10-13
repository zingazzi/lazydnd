// ui/layout_core.go
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	return Model{
		ActivePanel:              DiceRoller,
		DiceInput:                "",
		DiceResult:               "",
		DiceHistory:              []string{},
		DiceCommands:             []string{},
		LastDiceCommand:          "",
		InputMode:                false,
		DiceHistoryMode:          false,
		HistoryIndex:             -1,
		ScrollOffset:             make(map[PanelType]int),
		DiceMacros:               getDefaultDiceMacros(),
		MacroListMode:            false,
		SelectedMacro:            -1,
		ShowMacroPrompt:          false,
		MacroNameInput:           "",
		MacroFormulaInput:        "",
		MacroInputStep:           0,
		SpellSearchInput:         "",
		SpellSearchMode:          false,
		SelectedSpell:            nil,
		SpellSuggestions:         []string{},
		SuggestionIndex:          -1,
		InitiativeList:           []InitiativeEntry{},
		InitiativeInput:          "",
		InitiativeInputMode:      false,
		InitiativeInputType:      "",
		SelectedEntry:            -1,
		TempEntry:                InitiativeEntry{},
		CurrentTurn:              -1,
		RoundCounter:             0,
		InitiativeEditMode:       false,
		InitiativeEditType:       "",
		InitiativeListMode:       false,
		MonsterSearchInput:       "",
		MonsterSearchMode:        false,
		SelectedMonster:          nil,
		MonsterSuggestions:       []string{},
		MonsterSuggestionIndex:   -1,
		MonsterCRFilter:          "",
		MonsterCRFilterMode:      false,
		SpellLevelFilter:         "",
		SpellLevelFilterMode:     false,
		ActiveSpells:             []ActiveSpell{},
		ActiveSpellIndex:         -1,
		ActiveSpellListMode:      false,
		ShowCastSpellPrompt:      false,
		CastSpellInput:           "",
		CastSpellInputMode:       false,
		SpellToCast:              nil,
		MultiTargetMode:          false,
		SelectedTargets:          make(map[int]bool),
		ShowMultiTargetPopup:     false,
		MultiTargetInput:         "",
		MultiTargetType:          "damage",
		MultiTargetSaveMode:      false,
		TargetSaveResults:        make(map[int]string),
		ShowConditionPopup:       false,
		ConditionPopupMode:       "",
		ConditionInput:           "",
		ConditionDurationInput:   "",
		ConditionInputStep:       0,
		SelectedConditionIdx:     0,
		SelectedConditionNameIdx: 0,
		ShowHelpPopup:            false,
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

	// Show quick HP popup
	if m.ShowQuickHPPopup {
		return m.renderQuickHPPopupOverlay(mainView)
	}

	// Show rename popup if active (highest priority)
	if m.ShowRenamePopup {
		return m.renderRenamePopupOverlay(mainView)
	}

	// Show cast spell popup if active (takes priority over other popups)
	if m.ShowCastSpellPrompt {
		return renderCastSpellPopupOverlay(m, mainView)
	}

	// Show multi-target popup if active (takes priority over other popups)
	if m.ShowMultiTargetPopup {
		return renderMultiTargetPopupOverlay(m, mainView)
	}

	// Show condition popup if active (takes priority over other popups)
	if m.ShowConditionPopup {
		return renderConditionPopupOverlay(m, mainView)
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

// calculatePanelDimensions calculates panel dimensions with dynamic sizing
func (m Model) calculatePanelDimensions() map[PanelType]PanelDimensions {
	// Reserve 2 lines for status bar at the bottom
	availableHeight := m.Height - 2
	availableWidth := m.Width - 10 // Account for borders and spacing

	dimensions := make(map[PanelType]PanelDimensions)

	// Top row height (Dice and Initiative)
	topHeight := (availableHeight - 4) / 2

	// Bottom row height (Spells, Monsters, Notes)
	bottomHeight := availableHeight - topHeight - 4

	// Width allocation based on active panel
	switch m.ActivePanel {
	case DiceRoller:
		// Dice enlarged
		dimensions[DiceRoller] = PanelDimensions{Width: availableWidth*6/10, Height: topHeight}
		dimensions[InitiativeTracker] = PanelDimensions{Width: availableWidth*4/10, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: availableWidth*4/10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: availableWidth*4/10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: availableWidth*2/10, Height: bottomHeight}

	case InitiativeTracker:
		// Initiative enlarged
		dimensions[DiceRoller] = PanelDimensions{Width: availableWidth*4/10, Height: topHeight}
		dimensions[InitiativeTracker] = PanelDimensions{Width: availableWidth*6/10, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: availableWidth*4/10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: availableWidth*4/10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: availableWidth*2/10, Height: bottomHeight}

	case Spells:
		// Spells enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[InitiativeTracker] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: availableWidth*5/10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: availableWidth*3/10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: availableWidth*2/10, Height: bottomHeight}

	case Monsters:
		// Monsters enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[InitiativeTracker] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: availableWidth*3/10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: availableWidth*5/10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: availableWidth*2/10, Height: bottomHeight}

	case Notes:
		// Notes enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[InitiativeTracker] = PanelDimensions{Width: availableWidth*5/10, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: availableWidth*3/10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: availableWidth*3/10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: availableWidth*4/10, Height: bottomHeight}
	}

	return dimensions
}

// arrangeInGrid arranges the five panels in custom layout
func (m Model) arrangeInGrid(panelViews []string) string {
	// Top row: Dice (0) and Initiative (1)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[0], panelViews[1])

	// Bottom row: Spells (2), Monsters (3), Notes (4)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[2], panelViews[3], panelViews[4])

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

// renderQuickHPPopupOverlay renders the quick HP popup over the main view
func (m Model) renderQuickHPPopupOverlay(mainView string) string {
	popup := RenderQuickHPPopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}

// renderRenamePopupOverlay renders the rename popup over the main view
func (m Model) renderRenamePopupOverlay(mainView string) string {
	popup := RenderRenamePopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}
