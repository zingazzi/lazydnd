// panels/monsters.go
package panels

import (
	"fmt"
	"reflect"
	"strings"

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
	var contentLines []string

	// Header
	contentLines = append(contentLines, "🐲 MONSTER COMPENDIUM 🐲")
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, "Search D&D 5e Monsters")
	contentLines = append(contentLines, "Press Enter to start searching")
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 40))
	contentLines = append(contentLines, "")

	// Search input
	if searchMode {
		var prompt string
		if isActive {
			prompt = "Search: " + searchInput + "█"
		} else {
			prompt = "Search: " + searchInput
		}
		contentLines = append(contentLines, monsterInputStyle.Render(prompt))
		contentLines = append(contentLines, "")

		// Show suggestions
		if len(suggestions) > 0 {
			contentLines = append(contentLines, "Suggestions:")
			for i, suggestion := range suggestions {
				if i == suggestionIndex {
					contentLines = append(contentLines, selectedMonsterSuggestionStyle.Render("► "+suggestion))
				} else {
					contentLines = append(contentLines, monsterSuggestionStyle.Render("  "+suggestion))
				}
			}
			contentLines = append(contentLines, "")
		}
	}

	// Show selected monster details
	if selectedMonster != nil {
		monsterDetails := FormatSelectedMonster(selectedMonster)
		if monsterDetails != "" {
			contentLines = append(contentLines, "Monster Details:")
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, "Press 'a' to Add to Initiative")
			contentLines = append(contentLines, "")
			// Split the details into lines and add them
			detailLines := strings.Split(monsterDetails, "\n")
			for _, line := range detailLines {
				contentLines = append(contentLines, line)
			}
		}
	} else if !searchMode {
		contentLines = append(contentLines, "No monster selected")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Press Enter to search for monsters")
	}

	return strings.Join(contentLines, "\n")
}

// FormatSelectedMonster formats the selected monster for display
func FormatSelectedMonster(selectedMonster interface{}) string {
	if selectedMonster == nil {
		return ""
	}

	// Use reflection to handle the interface{} type
	v := reflect.ValueOf(selectedMonster)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	// Extract fields using reflection
	name := getFieldString(v, "Name")
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

// getFieldString gets a string field value using reflection
func getFieldString(v reflect.Value, fieldName string) string {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}
	if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}
