// panels/monsters.go
package panels

import (
	"fmt"
)

// Note: Styling is now handled by TView widgets - these functions return plain text

// GetMonstersContent returns the content for the monsters panel
func GetMonstersContent(searchInput string, selectedMonster interface{}, suggestions []string, suggestionIndex int, searchMode bool, isActive bool, crFilter string, crFilterMode bool) string {
	// If in CR filter mode, use CR filter as the "search" input
	if crFilterMode {
		return RenderSearchContent(SearchContentConfig{
			Title:           "Filter by Challenge Rating",
			ItemType:        "monster",
			SearchInput:     crFilter,
			SelectedItem:    nil, // Never show selected item in CR filter mode
			Suggestions:     suggestions,
			SuggestionIndex: suggestionIndex,
			SearchMode:      true, // Show as search mode
			IsActive:        isActive,
			FormatFunc:      FormatSelectedMonster,
			ShowAddPrompt:   false, // No "add to initiative" in CR mode
		})
	}

	// Regular search mode
	return RenderSearchContent(SearchContentConfig{
		Title:           "Search D&D 5e Monsters",
		ItemType:        "monster",
		SearchInput:     searchInput,
		SelectedItem:    selectedMonster,
		Suggestions:     suggestions,
		SuggestionIndex: suggestionIndex,
		SearchMode:      searchMode,
		IsActive:        isActive,
		FormatFunc:      FormatSelectedMonster,
		ShowAddPrompt:   true,
	})
}

// FormatSelectedMonster formats the selected monster for display
func FormatSelectedMonster(selectedMonster interface{}) string {
	name := ExtractNameFromInterface(selectedMonster)
	if name == "" {
		return ""
	}

	// Find the monster and format it
	monster := FindMonster(name)
	if monster == nil {
		return fmt.Sprintf("Monster '%s' not found in database", name)
	}

	return FormatMonster(monster)
}
