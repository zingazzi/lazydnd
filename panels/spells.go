// panels/spells.go
package panels

import (
	"fmt"
	"reflect"
	"strings"

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
func GetSpellsContent(searchInput string, selectedSpell interface{}, suggestions []string, suggestionIndex int, searchMode, isActive bool) string {
	var contentLines []string

	// Header
	contentLines = append(contentLines, "✨ SPELL COMPENDIUM ✨")
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, "Search D&D 5e Spells")
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
		contentLines = append(contentLines, spellInputStyle.Render(prompt))
		contentLines = append(contentLines, "")

		// Show suggestions
		if len(suggestions) > 0 {
			contentLines = append(contentLines, "Suggestions:")
			for i, suggestion := range suggestions {
				if i == suggestionIndex {
					contentLines = append(contentLines, selectedSpellSuggestionStyle.Render("► "+suggestion))
				} else {
					contentLines = append(contentLines, spellSuggestionStyle.Render("  "+suggestion))
				}
			}
			contentLines = append(contentLines, "")
		}
	}

	// Show selected spell details
	if selectedSpell != nil {
		spellDetails := FormatSelectedSpell(selectedSpell)
		if spellDetails != "" {
			contentLines = append(contentLines, "Spell Details:")
			contentLines = append(contentLines, "")
			// Split the details into lines and add them
			detailLines := strings.Split(spellDetails, "\n")
			for _, line := range detailLines {
				contentLines = append(contentLines, line)
			}
		}
	} else if !searchMode {
		contentLines = append(contentLines, "No spell selected")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Press Enter to search for spells")
	}

	return strings.Join(contentLines, "\n")
}

// FormatSelectedSpell formats the selected spell for display
func FormatSelectedSpell(selectedSpell interface{}) string {
	if selectedSpell == nil {
		return ""
	}

	// Use reflection to handle the interface{} type
	v := reflect.ValueOf(selectedSpell)
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
	name := getSpellFieldString(v, "Name")
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

// getSpellFieldString gets a string field value using reflection
func getSpellFieldString(v reflect.Value, fieldName string) string {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}
	if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}
