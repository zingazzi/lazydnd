// ui/layout_content.go
package ui

import (
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

// ========== HELP TEXT PROVIDERS ==========

// HelpTextProvider defines a function type for getting help text
type HelpTextProvider func(Model) string

// helpTextProviders maps panel types to their help text provider functions
var helpTextProviders = map[PanelType]HelpTextProvider{
	DiceRoller:        getDiceRollerHelpText,
	InitiativeTracker: getInitiativeTrackerHelpText,
	Spells:            getSpellsHelpText,
	Monsters:          getMonstersHelpText,
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
	text := InitiativeTrackerInlineHelp(m.InitiativeEditMode, m.InitiativeInputMode, m.InitiativeListMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getSpellsHelpText gets help text for the spells panel
func getSpellsHelpText(m Model) string {
	text := SpellsInlineHelp(m.SpellSearchMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}

// getMonstersHelpText gets help text for the monsters panel
func getMonstersHelpText(m Model) string {
	text := MonstersInlineHelp(m.MonsterSearchMode)
	return "\n" + m.Styles.HelpStyle.Render(text)
}
