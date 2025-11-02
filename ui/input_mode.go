// ui/input_mode.go
package ui

// InputMode represents the current input/mode state of the application
type InputMode int

const (
	ModeNormal InputMode = iota
	ModeDiceInput
	ModeDiceHistory
	ModeSpellSearch
	ModeSpellFilter
	ModeMonsterSearch
	ModeMonsterFilter
	ModeInitiativeInput
	ModeInitiativeEdit
	ModeInitiativeList
	ModeNotesEdit
	ModeNotesSearch
	ModeEncounterBuilding
	ModeEncounterFilter
	ModeCastSpell
	ModeMultiTarget
	ModeQuickHP
	ModeMacroCreate
)

// GetInputMode returns the current input mode based on the model's state
// Checks both new state structs and legacy fields for compatibility during migration
func (m Model) GetInputMode() InputMode {
	// Check popups first (they take priority)
	// Use legacy fields during migration (they're kept in sync)
	if (m.Spells.ShowCastSpellPrompt || m.ShowCastSpellPrompt) && (m.Spells.CastSpellInputMode || m.CastSpellInputMode) {
		return ModeCastSpell
	}
	if m.Initiative.ShowMultiTargetPopup || m.ShowMultiTargetPopup {
		return ModeMultiTarget
	}
	if m.Initiative.ShowQuickHPPopup || m.ShowQuickHPPopup {
		return ModeQuickHP
	}
	if m.DiceRoller.ShowMacroPrompt || m.ShowMacroPrompt {
		return ModeMacroCreate
	}
	if m.Popup.ShowCondition || m.ShowConditionPopup {
		return ModeNormal // Condition popup handles its own input
	}
	if m.Popup.ShowSave || m.Popup.ShowLoad || m.Popup.ShowRename || m.ShowSavePopup || m.ShowLoadPopup || m.ShowRenamePopup {
		return ModeNormal // Popup handles its own input
	}

	// Check dice roller modes (check both new and legacy)
	if m.DiceRoller.HistoryMode || m.DiceHistoryMode {
		return ModeDiceHistory
	}
	if m.DiceRoller.InputMode || m.InputMode {
		return ModeDiceInput
	}

	// Check spell panel modes
	if m.Spells.LevelFilterMode || m.SpellLevelFilterMode {
		return ModeSpellFilter
	}
	if m.Spells.SearchMode || m.SpellSearchMode {
		return ModeSpellSearch
	}
	if m.Spells.ActiveSpellListMode || m.ActiveSpellListMode {
		return ModeNormal // Navigation mode
	}

	// Check monster panel modes
	if m.Monsters.CRFilterMode || m.MonsterCRFilterMode {
		return ModeMonsterFilter
	}
	if m.Monsters.SearchMode || m.MonsterSearchMode {
		return ModeMonsterSearch
	}

	// Check initiative tracker modes
	if m.Initiative.MultiTargetMode || m.MultiTargetMode {
		return ModeMultiTarget
	}
	if m.Initiative.EditMode || m.InitiativeEditMode {
		return ModeInitiativeEdit
	}
	if m.Initiative.ListMode || m.InitiativeListMode {
		return ModeInitiativeList
	}
	if m.Initiative.InputMode || m.InitiativeInputMode {
		return ModeInitiativeInput
	}

	// Check notes panel modes
	if m.Notes.SearchMode || m.NotesSearchMode {
		return ModeNotesSearch
	}
	if m.Notes.EditMode || m.NotesEditMode {
		return ModeNotesEdit
	}

	// Check encounter builder modes
	if m.Encounter.Generating || m.EncounterGenerating {
		return ModeNormal // Generator popup handles its own input
	}
	if m.Encounter.FilterActive || m.EncounterFilterActive {
		return ModeEncounterFilter
	}
	if m.ActivePanel == EncounterBuilder && (m.Encounter.BuilderMode != "" || m.EncounterBuilderMode != "") {
		return ModeEncounterBuilding
	}

	return ModeNormal
}

// SetInputMode sets the input mode and updates the corresponding boolean flags
// Updates both new state structs and legacy fields for compatibility during migration
func (m *Model) SetInputMode(mode InputMode) {
	// Reset all mode flags first (both new and legacy)
	m.resetInputModes()

	switch mode {
	case ModeDiceInput:
		m.DiceRoller.InputMode = true
		m.InputMode = true // Legacy
	case ModeDiceHistory:
		m.DiceRoller.HistoryMode = true
		m.DiceHistoryMode = true // Legacy
	case ModeSpellSearch:
		m.Spells.SearchMode = true
		m.SpellSearchMode = true // Legacy
	case ModeSpellFilter:
		m.Spells.LevelFilterMode = true
		m.SpellLevelFilterMode = true // Legacy
	case ModeMonsterSearch:
		m.Monsters.SearchMode = true
		m.MonsterSearchMode = true // Legacy
	case ModeMonsterFilter:
		m.Monsters.CRFilterMode = true
		m.MonsterCRFilterMode = true // Legacy
	case ModeInitiativeInput:
		m.Initiative.InputMode = true
		m.InitiativeInputMode = true // Legacy
	case ModeInitiativeEdit:
		m.Initiative.EditMode = true
		m.InitiativeEditMode = true // Legacy
	case ModeInitiativeList:
		m.Initiative.ListMode = true
		m.InitiativeListMode = true // Legacy
	case ModeNotesEdit:
		m.Notes.EditMode = true
		m.NotesEditMode = true // Legacy
	case ModeNotesSearch:
		m.Notes.SearchMode = true
		m.NotesSearchMode = true // Legacy
	case ModeEncounterFilter:
		m.Encounter.FilterActive = true
		m.EncounterFilterActive = true // Legacy
	case ModeCastSpell:
		m.Spells.ShowCastSpellPrompt = true
		m.Spells.CastSpellInputMode = true
		m.ShowCastSpellPrompt = true // Legacy
		m.CastSpellInputMode = true  // Legacy
	case ModeMultiTarget:
		m.Initiative.MultiTargetMode = true
		m.MultiTargetMode = true // Legacy
	case ModeQuickHP:
		m.Initiative.ShowQuickHPPopup = true
		m.ShowQuickHPPopup = true // Legacy
	case ModeMacroCreate:
		m.DiceRoller.ShowMacroPrompt = true
		m.ShowMacroPrompt = true // Legacy
	case ModeNormal:
		// Already reset above
	}
}

// resetInputModes resets all input mode flags (both new and legacy)
func (m *Model) resetInputModes() {
	// Reset new state structs
	m.DiceRoller.InputMode = false
	m.DiceRoller.HistoryMode = false
	m.Spells.SearchMode = false
	m.Spells.LevelFilterMode = false
	m.Monsters.SearchMode = false
	m.Monsters.CRFilterMode = false
	m.Initiative.InputMode = false
	m.Initiative.EditMode = false
	m.Initiative.ListMode = false
	m.Notes.EditMode = false
	m.Notes.SearchMode = false
	m.Encounter.FilterActive = false
	m.Spells.ShowCastSpellPrompt = false
	m.Spells.CastSpellInputMode = false
	m.Initiative.MultiTargetMode = false
	m.Initiative.ShowMultiTargetPopup = false
	m.Initiative.ShowQuickHPPopup = false
	m.DiceRoller.ShowMacroPrompt = false

	// Reset legacy fields
	m.InputMode = false
	m.DiceHistoryMode = false
	m.SpellSearchMode = false
	m.SpellLevelFilterMode = false
	m.MonsterSearchMode = false
	m.MonsterCRFilterMode = false
	m.InitiativeInputMode = false
	m.InitiativeEditMode = false
	m.InitiativeListMode = false
	m.NotesEditMode = false
	m.NotesSearchMode = false
	m.EncounterFilterActive = false
	m.ShowCastSpellPrompt = false
	m.CastSpellInputMode = false
	m.MultiTargetMode = false
	m.ShowMultiTargetPopup = false
	m.ShowQuickHPPopup = false
	m.ShowMacroPrompt = false
}

// IsInputMode returns true if the model is in any input mode
func (m Model) IsInputMode() bool {
	return m.GetInputMode() != ModeNormal
}

// IsSearchMode returns true if the model is in a search mode
func (m Model) IsSearchMode() bool {
	mode := m.GetInputMode()
	return mode == ModeSpellSearch || mode == ModeMonsterSearch || mode == ModeNotesSearch
}

// IsEditMode returns true if the model is in an edit mode
func (m Model) IsEditMode() bool {
	mode := m.GetInputMode()
	return mode == ModeInitiativeEdit || mode == ModeNotesEdit || mode == ModeDiceInput || mode == ModeMacroCreate
}

// IsPopupMode returns true if a popup is currently showing
// Checks both new state structs and legacy fields
func (m Model) IsPopupMode() bool {
	return m.Popup.ShowHelp || m.ShowHelpPopup ||
		m.Popup.ShowAction || m.ShowActionPopup ||
		m.Popup.ShowSavingThrow || m.ShowSavingThrowPopup ||
		m.Popup.ShowCondition || m.ShowConditionPopup ||
		m.Popup.ShowSave || m.ShowSavePopup ||
		m.Popup.ShowLoad || m.ShowLoadPopup ||
		m.Popup.ShowRename || m.ShowRenamePopup ||
		m.Initiative.ShowQuickHPPopup || m.ShowQuickHPPopup ||
		m.Initiative.ShowMultiTargetPopup || m.ShowMultiTargetPopup ||
		m.Spells.ShowCastSpellPrompt || m.ShowCastSpellPrompt ||
		m.Popup.ShowEncounterPrompt || m.ShowEncounterPrompt ||
		m.Encounter.Generating || m.EncounterGenerating
}
