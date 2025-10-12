// ui/layout_content.go
package ui

import (
	"fmt"
	"lazydnd/panels"
)

// PanelContentProvider defines a function type for getting panel content
type PanelContentProvider func(Model) string

// panelContentProviders maps panel types to their content provider functions
var panelContentProviders = map[PanelType]PanelContentProvider{
	DiceRoller:        getDiceRollerContent,
	InitiativeTracker: getInitiativeTrackerContent,
	Spells:            getSpellsContent,
	Monsters:          getMonstersContent,
	Notes:             getNotesContent,
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

// GetPanelContent is exported for testing
func (m Model) GetPanelContent(panelType PanelType) string {
	return m.getPanelContent(panelType)
}

// ========== PANEL CONTENT PROVIDERS ==========

// getDiceRollerContent gets content for the dice roller panel
func getDiceRollerContent(m Model) string {
	return panels.GetDiceRollerContent(
		m.DiceInput,
		m.DiceResult,
		m.DiceHistory,
		m.DiceCommands,
		m.LastDiceCommand,
		m.InputMode,
		m.ActivePanel == DiceRoller,
		m.DiceHistoryMode,
		m.HistoryIndex,
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
		m.CurrentTurn,
		m.RoundCounter,
		m.MultiTargetMode,
		m.SelectedTargets,
	)
}

// getSpellsContent gets content for the spells panel
func getSpellsContent(m Model) string {
	// If showing active spells list
	if m.ActiveSpellListMode {
		content := FormatActiveSpells(m.ActiveSpells, m.ActiveSpellIndex, m.ActivePanel == Spells)
		return content
	}

	// Normal spell search content
	content := panels.GetSpellsContent(
		m.SpellSearchInput,
		m.SelectedSpell,
		m.SpellSuggestions,
		m.SuggestionIndex,
		m.SpellSearchMode,
		m.ActivePanel == Spells,
		!m.CastSpellInputMode,
		m.SpellLevelFilter,
		m.SpellLevelFilterMode,
	)

	// If there are active spells, add a note
	if len(m.ActiveSpells) > 0 && !m.SpellSearchMode && !m.CastSpellInputMode {
		content += "\n\n📜 Press 'v' to view active spells (" + formatSpellCount(len(m.ActiveSpells)) + ")"
	}

	return content
}

// formatSpellCount formats the spell count
func formatSpellCount(count int) string {
	if count == 1 {
		return "1 active spell"
	}
	return fmt.Sprintf("%d active spells", count)
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
		m.MonsterCRFilter,
		m.MonsterCRFilterMode,
	)
}

// getNotesContent gets content for the notes panel
func getNotesContent(m Model) string {
	content := m.NotesContent
	if m.NotesEditMode {
		content = m.NotesInput
	}

	return panels.GetNotesContent(
		content,
		m.NotesEditMode,
		m.NotesSearchMode,
		m.NotesSearchInput,
		m.NotesSearchResult,
		m.Width,
	)
}

// ========== HELP TEXT PROVIDERS ==========

// HelpTextProvider defines a function type for getting help text
type HelpTextProvider func(Model) string

// helpTextProviders maps panel types to their help text provider functions
var helpTextProviders = map[PanelType]HelpTextProvider{
	DiceRoller:        getDiceRollerHelpText,
	InitiativeTracker: getInitiativeTrackerHelpText,
	Spells:            getSpellsHelpText,
	Monsters:          getMonstersHelpText,
	Notes:             getNotesHelpText,
}

// getHelpText returns context-sensitive help text
func (m Model) getHelpText(panelType PanelType) string {
	provider, exists := helpTextProviders[panelType]
	if !exists {
		return "\n" + m.Styles.HelpStyle.Render(DefaultInlineHelp())
	}

	return provider(m)
}

// getDiceRollerHelpText gets help text for the dice roller panel
func getDiceRollerHelpText(m Model) string {
	text := DiceRollerInlineHelp(m.InputMode, m.LastDiceCommand != "", m.DiceHistoryMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getInitiativeTrackerHelpText gets help text for the initiative tracker panel
func getInitiativeTrackerHelpText(m Model) string {
	text := InitiativeTrackerInlineHelp(m.InitiativeEditMode, m.InitiativeInputMode, m.InitiativeListMode, m.MultiTargetMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getSpellsHelpText gets help text for the spells panel
func getSpellsHelpText(m Model) string {
	text := SpellsInlineHelp(m)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getMonstersHelpText gets help text for the monsters panel
func getMonstersHelpText(m Model) string {
	text := MonstersInlineHelp(m)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getNotesHelpText gets help text for the notes panel
func getNotesHelpText(m Model) string {
	text := NotesInlineHelp(m.NotesEditMode, m.NotesSearchMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}
