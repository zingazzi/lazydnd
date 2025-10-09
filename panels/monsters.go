// panels/monsters.go
package panels

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	monsterInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	monsterSuggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

	selectedMonsterSuggestionStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#7D56F4")).
					Foreground(lipgloss.Color("#FAFAFA"))

	monsterDetailStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#444444")).
				Padding(1, 2).
				Margin(1, 0)
)

// GetMonstersContent returns the content for the monsters panel
func GetMonstersContent(searchInput string, selectedMonster interface{}, suggestions []string, suggestionIndex int, searchMode bool, isActive bool) string {
	return RenderSearchContent(SearchContentConfig{
		Title:           "Search D&D 5e Monsters",
		ItemType:        "monster",
		SearchInput:     searchInput,
		SelectedItem:    selectedMonster,
		Suggestions:     suggestions,
		SuggestionIndex: suggestionIndex,
		SearchMode:      searchMode,
		IsActive:        isActive,
		InputStyle:      monsterInputStyle,
		SuggestionStyle: monsterSuggestionStyle,
		SelectedStyle:   selectedMonsterSuggestionStyle,
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
