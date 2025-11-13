// ui/layout_core.go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	return Model{
		ActivePanel:  DiceRoller,
		Width:        0,
		Height:       0,
		ScrollOffset: make(map[PanelType]int),

		// Initialize new state structs
		DiceRoller: DiceRollerState{
			Input:            "",
			Result:           "",
			History:          []string{},
			Commands:         []string{},
			LastCommand:      "",
			InputMode:        false,
			HistoryMode:      false,
			HistoryIndex:     -1,
			Macros:           getDefaultDiceMacros(),
			MacroListMode:    false,
			SelectedMacro:    -1,
			ShowMacroPrompt:  false,
			MacroNameInput:   "",
			MacroFormulaInput: "",
			MacroInputStep:   0,
		},
		Initiative: InitiativeState{
			List:                []InitiativeEntry{},
			Input:               "",
			InputMode:           false,
			InputType:           "",
			SelectedEntry:       -1,
			TempEntry:           InitiativeEntry{},
			EditMode:            false,
			EditType:            "",
			ListMode:            false,
			CurrentTurn:         -1,
			RoundCounter:        0,
			ShowQuickHPPopup:    false,
			QuickHPInput:        "",
			QuickHPMode:         "",
			MultiTargetMode:     false,
			SelectedTargets:     make(map[int]bool),
			ShowMultiTargetPopup: false,
			MultiTargetInput:    "",
			MultiTargetType:     "damage",
			MultiTargetSaveMode: false,
			TargetSaveResults:   make(map[int]string),
			HPUndoStack:         []HPHistoryEntry{},
			HPRedoStack:         []HPHistoryEntry{},
		},
		Spells: SpellState{
			SearchInput:         "",
			SearchMode:          false,
			SelectedSpell:       nil,
			Suggestions:         []string{},
			SuggestionIndex:     -1,
			LevelFilter:         "",
			LevelFilterMode:     false,
			ActiveSpells:        []ActiveSpell{},
			ActiveSpellIndex:    -1,
			ActiveSpellListMode: false,
			ShowCastSpellPrompt: false,
			CastSpellInput:      "",
			CastSpellInputMode:  false,
			SpellToCast:         nil,
		},
		Monsters: MonsterState{
			SearchInput:     "",
			SearchMode:      false,
			SelectedMonster: nil,
			Suggestions:     []string{},
			SuggestionIndex: -1,
			CRFilter:        "",
			CRFilterMode:    false,
		},
		Notes: NotesState{
			Content:     "",
			Input:       "",
			EditMode:    false,
			SearchMode:  false,
			SearchInput: "",
			SearchResult: []int{},
		},
		Encounter: EncounterState{
			PartySize:            4,
			PartyLevel:           3,
			Monsters:             []EncounterMonster{},
			SelectedIndex:        -1,
			SavedEncounters:      []Encounter{},
			ListMode:             false,
			NameInput:            "",
			ShowPrompt:           false,
			LoadedTemplateName:   "",
			BuilderMode:          "party_setup",
			CRFilter:             "",
			FilterActive:         false,
			SelectedSaved:        -1,
			AddingMonster:        false,
			Environment:          "Any",
			Difficulty:           "medium",
			Generating:           false,
			EnvironmentIndex:     0,
			DifficultyIndex:      1, // medium
			GeneratorFocus:       "",
			AvailableEnvironments: []string{"Any", "Forest", "Mountain", "Desert", "Swamp", "Underdark", "Urban", "Coast", "Arctic", "Jungle", "Plains"},
		},
		Popup: PopupState{
			ShowHelp:                false,
			HelpScrollOffset:        0,
			ShowAction:              false,
			ActionActions:           []MonsterAction{},
			ActionIndex:             0,
			ActionMonster:           "",
			ActionAdvantage:         false,
			ActionDisadvantage:      false,
			ShowSavingThrow:         false,
			ShowCondition:           false,
			ConditionMode:           "",
			ConditionInput:          "",
			ConditionDurationInput:  "",
			ConditionInputStep:      0,
			SelectedConditionIdx:    0,
			SelectedConditionNameIdx: 0,
			ShowSave:                false,
			ShowLoad:                false,
			ShowRename:              false,
			SaveInput:               "",
			CurrentCampaignFile:     "",
			CurrentCampaignName:     "",
			ShowEncounterPrompt:     false,
		},
		Global: GlobalState{
			DebugMode:         false,
			CampaignList:      []string{},
			CampaignListIndex: 0,
			LastAutoSave:      "",
			ErrorMessage:      "",
			ErrorVisible:      false,
		},

		// Initialize legacy fields for backward compatibility
		DiceInput:                "",
		DiceResult:               "",
		DiceHistory:              []string{},
		DiceCommands:             []string{},
		LastDiceCommand:          "",
		InputMode:                false,
		DiceHistoryMode:          false,
		HistoryIndex:             -1,
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
		PartySize:                4,
		PartyLevel:               3,
		EncounterMonsters:        []EncounterMonster{},
		SelectedEncounterIndex:   -1,
		SavedEncounters:          []Encounter{},
		EncounterListMode:        false,
		EncounterNameInput:       "",
		ShowEncounterPrompt:      false,
		EncounterBuilderMode:     "party_setup",
		EncounterCRFilter:        "",
		EncounterFilterActive:    false,
		EncounterSelectedSaved:   -1,
		AddingMonsterToEncounter: false,
		EncounterEnvironment:     "Any",
		EncounterDifficulty:      "medium",
		EncounterGenerating:      false,
		EncounterEnvironmentIndex: 0,
		EncounterDifficultyIndex:  1,
		AvailableEnvironments:    []string{"Any", "Forest", "Mountain", "Desert", "Swamp", "Underdark", "Urban", "Coast", "Arctic", "Jungle", "Plains"},
		HPUndoStack:              []HPHistoryEntry{},
		HPRedoStack:              []HPHistoryEntry{},
		ErrorMessage:             "",
		ErrorVisible:             false,
		NotesContent:             "",
		NotesInput:               "",
		NotesEditMode:            false,
		NotesSearchMode:          false,
		NotesSearchInput:         "",
		NotesSearchResult:        []int{},
		DebugMode:                false,
		CampaignList:             []string{},
		CampaignListIndex:        0,
		LastAutoSave:             "",
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

	// Show save popup if active (highest priority) - check both new and legacy
	if m.Popup.ShowSave || m.ShowSavePopup {
		return m.renderSavePopupOverlay(mainView)
	}

	// Show load popup if active (highest priority) - check both new and legacy
	if m.Popup.ShowLoad || m.ShowLoadPopup {
		return m.renderLoadPopupOverlay(mainView)
	}

	// Show quick HP popup - check both new and legacy
	if m.Initiative.ShowQuickHPPopup || m.ShowQuickHPPopup {
		return m.renderQuickHPPopupOverlay(mainView)
	}

	// Show encounter prompt (save) - check both new and legacy
	if m.Popup.ShowEncounterPrompt || m.ShowEncounterPrompt {
		return m.renderEncounterPromptOverlay(mainView)
	}

	// Show encounter generator popup - check both new and legacy
	if m.Encounter.Generating || m.EncounterGenerating {
		return m.renderGeneratorPopupOverlay(mainView)
	}

	// Show rename popup if active (highest priority) - check both new and legacy
	if m.Popup.ShowRename || m.ShowRenamePopup {
		return m.renderRenamePopupOverlay(mainView)
	}

	// Show cast spell popup if active (takes priority over other popups) - check both new and legacy
	if m.Spells.ShowCastSpellPrompt || m.ShowCastSpellPrompt {
		return renderCastSpellPopupOverlay(m, mainView)
	}

	// Show multi-target popup if active (takes priority over other popups) - check both new and legacy
	if m.Initiative.ShowMultiTargetPopup || m.ShowMultiTargetPopup {
		return renderMultiTargetPopupOverlay(m, mainView)
	}

	// Show condition popup if active (takes priority over other popups) - check both new and legacy
	if m.Popup.ShowCondition || m.ShowConditionPopup {
		return renderConditionPopupOverlay(m, mainView)
	}

	// Show saving throw popup if active (takes priority over action popup) - check both new and legacy
	if m.Popup.ShowSavingThrow || m.ShowSavingThrowPopup {
		return m.renderSavingThrowPopupOverlay(mainView)
	}

	// Show action popup if active (takes priority over help popup) - check both new and legacy
	if m.Popup.ShowAction || m.ShowActionPopup {
		return m.renderActionPopupOverlay(mainView)
	}

	// Show help popup if active - check both new and legacy
	if m.Popup.ShowHelp || m.ShowHelpPopup {
		return m.renderHelpPopupOverlay(mainView)
	}

	return mainView
}

// Layout constants
const (
	MinPanelWidth   = 20 // Minimum panel width
	MinPanelHeight  = 8  // Minimum panel height
	StatusBarHeight = 2  // Status bar height
	GridSpacing     = 2  // Space between panels
	TopMargin       = 1  // Margin at the top of the grid
	RightMargin     = 1  // Margin on the right side of the grid
)

// PanelDimensions holds calculated panel dimensions
type PanelDimensions struct {
	Width  int
	Height int
}

// calculatePanelDimensions calculates panel dimensions with dynamic sizing
func (m Model) calculatePanelDimensions() map[PanelType]PanelDimensions {
	// Calculate available space (account for margins)
	// Target row width: terminal width minus right margin
	// Note: Each panel width includes its borders (2 chars: left + right)
	// When panels are joined, adjacent borders are separate, so we need to account for them
	// Top row has 2 panels = 1 border between them (already included in panel widths)
	// Bottom row has 4 panels = 3 borders between them (already included in panel widths)
	// Since borders are already in panel widths, targetRowWidth is just m.Width - RightMargin
	targetRowWidth := m.Width - RightMargin
	availableHeight := m.Height - StatusBarHeight - TopMargin

	// Ensure minimum sizes
	if targetRowWidth < MinPanelWidth*2 {
		targetRowWidth = MinPanelWidth * 2
	}
	if availableHeight < MinPanelHeight*2 {
		availableHeight = MinPanelHeight * 2
	}

	dimensions := make(map[PanelType]PanelDimensions)

	// Top row height (Dice and Initiative)
	topHeight := (availableHeight - GridSpacing) / 2

	// Bottom row height (Spells, Monsters, Notes, Encounter Builder)
	bottomHeight := availableHeight - topHeight - GridSpacing

	// Ensure minimum heights
	if topHeight < MinPanelHeight {
		topHeight = MinPanelHeight
		bottomHeight = availableHeight - topHeight - GridSpacing
	}
	if bottomHeight < MinPanelHeight {
		bottomHeight = MinPanelHeight
		topHeight = availableHeight - bottomHeight - GridSpacing
	}

	// Width allocation based on active panel
	// Use targetRowWidth and calculate last panel as remainder to ensure exact fit
	switch m.ActivePanel {
	case DiceRoller:
		// Dice enlarged
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 6 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}

	case InitiativeTracker:
		// Initiative enlarged
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 4 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}

	case Spells:
		// Spells enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}

	case Monsters:
		// Monsters enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}

	case Notes:
		// Notes enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}

	case EncounterBuilder:
		// Encounter Builder enlarged - top row splits evenly
		dimensions[DiceRoller] = PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		// Initiative Tracker gets remainder to ensure exact fit
		dimensions[InitiativeTracker] = PanelDimensions{Width: targetRowWidth - dimensions[DiceRoller].Width, Height: topHeight}
		dimensions[Spells] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Monsters] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[Notes] = PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		// Encounter Builder gets remainder to ensure exact fit
		dimensions[EncounterBuilder] = PanelDimensions{Width: targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width, Height: bottomHeight}
	}

	// Enforce minimum dimensions for all panels
	for panelType := range dimensions {
		dim := dimensions[panelType]
		if dim.Width < MinPanelWidth {
			dim.Width = MinPanelWidth
		}
		if dim.Height < MinPanelHeight {
			dim.Height = MinPanelHeight
		}
		dimensions[panelType] = dim
	}

	// Both rows should now sum to exactly targetRowWidth since we calculate last panel as remainder
	// Verify and ensure minimum widths are maintained, then enforce exact row width
	initTracker := dimensions[InitiativeTracker]
	if initTracker.Width < MinPanelWidth {
		initTracker.Width = MinPanelWidth
		dimensions[InitiativeTracker] = initTracker
		// Recalculate DiceRoller to maintain total
		diceRoller := dimensions[DiceRoller]
		diceRoller.Width = targetRowWidth - initTracker.Width
		if diceRoller.Width < MinPanelWidth {
			diceRoller.Width = MinPanelWidth
		}
		dimensions[DiceRoller] = diceRoller
	}
	encounterBuilder := dimensions[EncounterBuilder]
	if encounterBuilder.Width < MinPanelWidth {
		encounterBuilder.Width = MinPanelWidth
		dimensions[EncounterBuilder] = encounterBuilder
		// Recalculate other bottom panels proportionally to maintain total
		remainingWidth := targetRowWidth - encounterBuilder.Width
		// Redistribute remaining width proportionally
		totalRatio := 3 + 3 + 2 // Spells + Monsters + Notes ratios
		spells := dimensions[Spells]
		spells.Width = remainingWidth * 3 / totalRatio
		if spells.Width < MinPanelWidth {
			spells.Width = MinPanelWidth
		}
		dimensions[Spells] = spells
		monsters := dimensions[Monsters]
		monsters.Width = remainingWidth * 3 / totalRatio
		if monsters.Width < MinPanelWidth {
			monsters.Width = MinPanelWidth
		}
		dimensions[Monsters] = monsters
		notes := dimensions[Notes]
		notes.Width = remainingWidth - spells.Width - monsters.Width
		if notes.Width < MinPanelWidth {
			notes.Width = MinPanelWidth
		}
		dimensions[Notes] = notes
	}

	// Final verification: ensure both rows sum to exactly targetRowWidth
	// This accounts for borders between panels (each panel width includes its borders)
	topRowSum := dimensions[DiceRoller].Width + dimensions[InitiativeTracker].Width
	if topRowSum != targetRowWidth {
		// Adjust Initiative Tracker to make exact fit
		initTracker = dimensions[InitiativeTracker]
		initTracker.Width = targetRowWidth - dimensions[DiceRoller].Width
		if initTracker.Width < MinPanelWidth {
			initTracker.Width = MinPanelWidth
			// If still too small, adjust DiceRoller
			diceRoller := dimensions[DiceRoller]
			diceRoller.Width = targetRowWidth - initTracker.Width
			dimensions[DiceRoller] = diceRoller
		}
		dimensions[InitiativeTracker] = initTracker
	}

	bottomRowSum := dimensions[Spells].Width + dimensions[Monsters].Width + dimensions[Notes].Width + dimensions[EncounterBuilder].Width
	if bottomRowSum != targetRowWidth {
		// Adjust Encounter Builder to make exact fit
		encounterBuilder = dimensions[EncounterBuilder]
		encounterBuilder.Width = targetRowWidth - dimensions[Spells].Width - dimensions[Monsters].Width - dimensions[Notes].Width
		if encounterBuilder.Width < MinPanelWidth {
			encounterBuilder.Width = MinPanelWidth
			// If still too small, we need to reduce other panels proportionally
			// This is a fallback - should rarely happen
			remainingWidth := targetRowWidth - encounterBuilder.Width
			totalRatio := 3 + 3 + 2 // Spells + Monsters + Notes ratios
			spells := dimensions[Spells]
			spells.Width = remainingWidth * 3 / totalRatio
			if spells.Width < MinPanelWidth {
				spells.Width = MinPanelWidth
			}
			dimensions[Spells] = spells
			monsters := dimensions[Monsters]
			monsters.Width = remainingWidth * 3 / totalRatio
			if monsters.Width < MinPanelWidth {
				monsters.Width = MinPanelWidth
			}
			dimensions[Monsters] = monsters
			notes := dimensions[Notes]
			notes.Width = remainingWidth - spells.Width - monsters.Width
			if notes.Width < MinPanelWidth {
				notes.Width = MinPanelWidth
			}
			dimensions[Notes] = notes
			// Recalculate Encounter Builder with new values
			encounterBuilder.Width = targetRowWidth - spells.Width - monsters.Width - notes.Width
		}
		dimensions[EncounterBuilder] = encounterBuilder
	}

	return dimensions
}

// arrangeInGrid arranges the six panels in custom layout
func (m Model) arrangeInGrid(panelViews []string) string {
	// Top row: Dice (0) and Initiative (1)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[0], panelViews[1])

	// Bottom row: Spells (2), Monsters (3), Notes (4), Encounter Builder (5)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[2], panelViews[3], panelViews[4], panelViews[5])

	// Combine rows
	grid := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// Add top margin
	if TopMargin > 0 {
		grid = strings.Repeat("\n", TopMargin) + grid
	}

	// Add right margin by appending spaces to each line (preserves borders)
	if RightMargin > 0 {
		lines := strings.Split(grid, "\n")
		marginSpaces := strings.Repeat(" ", RightMargin)
		for i, line := range lines {
			lines[i] = line + marginSpaces
		}
		grid = strings.Join(lines, "\n")
	}

	return grid
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

// renderEncounterPromptOverlay renders the encounter save prompt popup over the main view
func (m Model) renderEncounterPromptOverlay(mainView string) string {
	popup := RenderEncounterPromptPopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}

// renderGeneratorPopupOverlay renders the encounter generator popup over the main view
func (m Model) renderGeneratorPopupOverlay(mainView string) string {
	popup := RenderGeneratorPopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}

// renderRenamePopupOverlay renders the rename popup over the main view
func (m Model) renderRenamePopupOverlay(mainView string) string {
	popup := RenderRenamePopup(m)

	// Place popup over main view with centered positioning
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, popup, lipgloss.WithWhitespaceChars("░"))
}
