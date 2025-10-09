// panels/search_utils.go
package panels

import (
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SearchContentConfig holds configuration for rendering search panel content
type SearchContentConfig struct {
	Title           string // e.g., "Search D&D 5e Monsters"
	ItemType        string // e.g., "monster", "spell"
	SearchInput     string
	SelectedItem    interface{}
	Suggestions     []string
	SuggestionIndex int
	SearchMode      bool
	IsActive        bool
	InputStyle      lipgloss.Style
	SuggestionStyle lipgloss.Style
	SelectedStyle   lipgloss.Style
	FormatFunc      func(interface{}) string
	ShowAddPrompt   bool // Whether to show "Press 'a' to add to initiative"
}

// RenderSearchContent renders generic search panel content
func RenderSearchContent(cfg SearchContentConfig) string {
	var contentLines []string

	// Header
	contentLines = append(contentLines, cfg.Title)
	contentLines = append(contentLines, "Press Enter to start searching")
	if cfg.ShowAddPrompt {
		contentLines = append(contentLines, "Press 'a' to add to initiative")
	}
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 40))
	contentLines = append(contentLines, "")

	// Search input
	if cfg.SearchMode {
		var prompt string
		if cfg.IsActive {
			prompt = "Search: " + cfg.SearchInput + "█"
		} else {
			prompt = "Search: " + cfg.SearchInput
		}
		contentLines = append(contentLines, cfg.InputStyle.Render(prompt))
		contentLines = append(contentLines, "")

		// Show suggestions
		if len(cfg.Suggestions) > 0 {
			contentLines = append(contentLines, "Suggestions:")
			for i, suggestion := range cfg.Suggestions {
				if i == cfg.SuggestionIndex {
					contentLines = append(contentLines, cfg.SelectedStyle.Render("► "+suggestion))
				} else {
					contentLines = append(contentLines, cfg.SuggestionStyle.Render("  "+suggestion))
				}
			}
			contentLines = append(contentLines, "")
		}
	}

	// Show selected item details
	if cfg.SelectedItem != nil && cfg.FormatFunc != nil {
		itemDetails := cfg.FormatFunc(cfg.SelectedItem)
		if itemDetails != "" {
			// Capitalize first letter of item type
			itemTypeTitle := strings.Title(cfg.ItemType)
			contentLines = append(contentLines, itemTypeTitle+" Details:")
			contentLines = append(contentLines, "")
			if cfg.ShowAddPrompt {
				contentLines = append(contentLines, "Press 'a' to Add to Initiative")
				contentLines = append(contentLines, "")
			}
			// Split the details into lines and add them
			detailLines := strings.Split(itemDetails, "\n")
			for _, line := range detailLines {
				contentLines = append(contentLines, line)
			}
		}
	} else if !cfg.SearchMode {
		contentLines = append(contentLines, "No "+cfg.ItemType+" selected")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Press Enter to search for "+cfg.ItemType+"s")
	}

	return strings.Join(contentLines, "\n")
}

// GetFieldString gets a string field value using reflection.
// This is a shared utility for extracting struct fields.
func GetFieldString(v reflect.Value, fieldName string) string {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}
	if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

// ExtractNameFromInterface extracts the Name field from an interface{} using reflection
func ExtractNameFromInterface(item interface{}) string {
	if item == nil {
		return ""
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	return GetFieldString(v, "Name")
}
