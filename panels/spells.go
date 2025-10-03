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
			Padding(0, 1).
			Margin(1, 0)

	suggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Margin(0, 1)

	selectedSuggestionStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FAFAFA")).
					Background(lipgloss.Color("#7D56F4")).
					Margin(0, 1)

	spellDetailStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(1).
				Margin(1, 0)
)

// GetSpellsContent returns the content for the spells panel
func GetSpellsContent(searchInput string, selectedSpell interface{}, suggestions []string, suggestionIndex int, searchMode, isActive bool) string {
	var contentLines []string

	// Fixed header section
	contentLines = append(contentLines, "✨ SPELL SEARCH ✨")
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, "Search for D&D spells by name")
	contentLines = append(contentLines, "Press Enter to start searching")
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 30))
	contentLines = append(contentLines, "")

	// Search input field
	searchPrompt := "Search: "
	if searchMode && isActive {
		searchPrompt += searchInput + "█"
	} else {
		searchPrompt += searchInput
	}
	contentLines = append(contentLines, searchPrompt)
	contentLines = append(contentLines, "")

	// Suggestions section (always reserve space)
	if len(suggestions) > 0 && searchMode {
		contentLines = append(contentLines, "Suggestions:")
		for i, suggestion := range suggestions {
			if i == suggestionIndex {
				contentLines = append(contentLines, "▶ "+suggestion)
			} else {
				contentLines = append(contentLines, "  "+suggestion)
			}
		}
	} else {
		// Add empty space to maintain consistent structure
		contentLines = append(contentLines, "")
	}

	// Separator before spell details
	contentLines = append(contentLines, "")
	if selectedSpell != nil {
		contentLines = append(contentLines, strings.Repeat("═", 40))

		// Format spell details and split into lines
		spellDetails := FormatSelectedSpell(selectedSpell)
		detailLines := strings.Split(spellDetails, "\n")
		contentLines = append(contentLines, detailLines...)
	}

	return strings.Join(contentLines, "\n")
}

// FormatSelectedSpell formats a spell from interface{} for display
func FormatSelectedSpell(selectedSpell interface{}) string {
	if selectedSpell == nil {
		return "No spell selected"
	}

	// Use reflection to extract spell data
	v := reflect.ValueOf(selectedSpell)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("Invalid spell data: %+v", selectedSpell)
	}

	// Extract fields using reflection
	var spell Spell

	if nameField := v.FieldByName("Name"); nameField.IsValid() && nameField.Kind() == reflect.String {
		spell.Name = nameField.String()
	}
	if levelField := v.FieldByName("Level"); levelField.IsValid() && levelField.Kind() == reflect.Int {
		spell.Level = int(levelField.Int())
	}
	if schoolField := v.FieldByName("School"); schoolField.IsValid() && schoolField.Kind() == reflect.String {
		spell.School = schoolField.String()
	}
	if actionField := v.FieldByName("ActionType"); actionField.IsValid() && actionField.Kind() == reflect.String {
		spell.ActionType = actionField.String()
	}
	if rangeField := v.FieldByName("Range"); rangeField.IsValid() && rangeField.Kind() == reflect.String {
		spell.Range = rangeField.String()
	}
	if durationField := v.FieldByName("Duration"); durationField.IsValid() && durationField.Kind() == reflect.String {
		spell.Duration = durationField.String()
	}
	if descField := v.FieldByName("Description"); descField.IsValid() && descField.Kind() == reflect.String {
		spell.Description = descField.String()
	}
	if materialField := v.FieldByName("Material"); materialField.IsValid() && materialField.Kind() == reflect.String {
		spell.Material = materialField.String()
	}
	if cantripField := v.FieldByName("CantripUpgrade"); cantripField.IsValid() && cantripField.Kind() == reflect.String {
		spell.CantripUpgrade = cantripField.String()
	}
	if concField := v.FieldByName("Concentration"); concField.IsValid() && concField.Kind() == reflect.Bool {
		spell.Concentration = concField.Bool()
	}
	if ritualField := v.FieldByName("Ritual"); ritualField.IsValid() && ritualField.Kind() == reflect.Bool {
		spell.Ritual = ritualField.Bool()
	}

	// Handle slice fields
	if classesField := v.FieldByName("Classes"); classesField.IsValid() && classesField.Kind() == reflect.Slice {
		spell.Classes = make([]string, classesField.Len())
		for i := 0; i < classesField.Len(); i++ {
			if elem := classesField.Index(i); elem.Kind() == reflect.String {
				spell.Classes[i] = elem.String()
			}
		}
	}
	if compField := v.FieldByName("Components"); compField.IsValid() && compField.Kind() == reflect.Slice {
		spell.Components = make([]string, compField.Len())
		for i := 0; i < compField.Len(); i++ {
			if elem := compField.Index(i); elem.Kind() == reflect.String {
				spell.Components[i] = elem.String()
			}
		}
	}

	return FormatSpell(&spell)
}
