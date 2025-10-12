// panels/spells.go
package panels

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	spellInputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	spellSuggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

	selectedSpellSuggestionStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#7D56F4")).
					Foreground(lipgloss.Color("#FAFAFA"))

	spellDetailStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#444444")).
				Padding(1, 2).
				Margin(1, 0)
)

// GetSpellsContent returns the content for the spells panel
func GetSpellsContent(searchInput string, selectedSpell interface{}, suggestions []string, suggestionIndex int, searchMode, isActive bool, showCastPrompt bool, levelFilter string, levelFilterMode bool) string {
	// If in level filter mode, use level filter as the "search" input
	if levelFilterMode {
		return RenderSearchContent(SearchContentConfig{
			Title:           "Filter by Spell Level",
			ItemType:        "spell",
			SearchInput:     levelFilter,
			SelectedItem:    nil, // Never show selected item in level filter mode
			Suggestions:     suggestions,
			SuggestionIndex: suggestionIndex,
			SearchMode:      true, // Show as search mode
			IsActive:        isActive,
			InputStyle:      spellInputStyle,
			SuggestionStyle: spellSuggestionStyle,
			SelectedStyle:   selectedSpellSuggestionStyle,
			FormatFunc:      FormatSelectedSpell,
			ShowAddPrompt:   false,
		})
	}

	// Regular search mode
	content := RenderSearchContent(SearchContentConfig{
		Title:           "Search D&D 5e Spells",
		ItemType:        "spell",
		SearchInput:     searchInput,
		SelectedItem:    selectedSpell,
		Suggestions:     suggestions,
		SuggestionIndex: suggestionIndex,
		SearchMode:      searchMode,
		IsActive:        isActive,
		InputStyle:      spellInputStyle,
		SuggestionStyle: spellSuggestionStyle,
		SelectedStyle:   selectedSpellSuggestionStyle,
		FormatFunc:      FormatSelectedSpell,
		ShowAddPrompt:   false,
	})

	// Cast spell prompt removed for cleaner interface

	return content
}

// FormatSelectedSpell formats the selected spell for display
func FormatSelectedSpell(selectedSpell interface{}) string {
	name := ExtractNameFromInterface(selectedSpell)
	if name == "" {
		return ""
	}

	// Find the spell and format it
	spell := FindSpell(name)
	if spell == nil {
		return fmt.Sprintf("Spell '%s' not found in database", name)
	}

	return FormatSpell(spell)
}
