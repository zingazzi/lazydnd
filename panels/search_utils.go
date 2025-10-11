// panels/search_utils.go
package panels

import (
	"fmt"
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

	// SEARCH MODE: Show search interface
	if cfg.SearchMode {
		contentLines = append(contentLines, cfg.Title)
		contentLines = append(contentLines, "Type to search, Esc to cancel")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, strings.Repeat("─", 35))
		contentLines = append(contentLines, "")
		var prompt string
		if cfg.IsActive {
			prompt = "Search: " + cfg.SearchInput + "█"
		} else {
			prompt = "Search: " + cfg.SearchInput
		}
		contentLines = append(contentLines, cfg.InputStyle.Render(prompt))
		contentLines = append(contentLines, "")

		// Show suggestions with scrolling window
		if len(cfg.Suggestions) > 0 {
			contentLines = append(contentLines, "Suggestions:")

			// Show a scrolling window of suggestions (max 8 visible at a time)
			maxVisible := 8
			totalSuggestions := len(cfg.Suggestions)

			// Calculate which suggestions to show
			startIdx := 0
			endIdx := totalSuggestions

			if totalSuggestions > maxVisible {
				// Center the selected item in the window when possible
				startIdx = cfg.SuggestionIndex - (maxVisible / 2)
				if startIdx < 0 {
					startIdx = 0
				}
				endIdx = startIdx + maxVisible
				if endIdx > totalSuggestions {
					endIdx = totalSuggestions
					startIdx = endIdx - maxVisible
					if startIdx < 0 {
						startIdx = 0
					}
				}
			}

			// Show scroll indicator at top if there are more items above
			if startIdx > 0 {
				contentLines = append(contentLines, cfg.SuggestionStyle.Render("  ⬆ "+strings.Repeat("─", 10)+" ("+formatSuggestionCount(startIdx)+" more above)"))
			}

			// Show the visible window of suggestions
			for i := startIdx; i < endIdx; i++ {
				suggestion := cfg.Suggestions[i]
				if i == cfg.SuggestionIndex {
					contentLines = append(contentLines, cfg.SelectedStyle.Render("► "+suggestion))
				} else {
					contentLines = append(contentLines, cfg.SuggestionStyle.Render("  "+suggestion))
				}
			}

			// Show scroll indicator at bottom if there are more items below
			if endIdx < totalSuggestions {
				remaining := totalSuggestions - endIdx
				contentLines = append(contentLines, cfg.SuggestionStyle.Render("  ⬇ "+strings.Repeat("─", 10)+" ("+formatSuggestionCount(remaining)+" more below)"))
			}

			contentLines = append(contentLines, "")
		}
		
		return strings.Join(contentLines, "\n")
	}

	// NOT IN SEARCH MODE: Show item details or default message
	
	// Header for non-search mode
	contentLines = append(contentLines, cfg.Title)
	contentLines = append(contentLines, "Press Enter to start searching")
	if cfg.ShowAddPrompt {
		contentLines = append(contentLines, "Press 'a' to add to initiative")
	}
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 35))
	contentLines = append(contentLines, "")

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
	} else {
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

// formatSuggestionCount formats the count for suggestion scroll indicators
func formatSuggestionCount(count int) string {
	if count == 1 {
		return "1 item"
	}
	return fmt.Sprintf("%d items", count)
}
