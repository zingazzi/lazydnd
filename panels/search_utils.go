// panels/search_utils.go
package panels

import (
	"fmt"
	"reflect"
	"strings"
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
	FormatFunc      func(interface{}) string
	ShowAddPrompt   bool // Whether to show "Press 'a' to add to initiative"
}

// RenderSearchContent renders generic search panel content
func RenderSearchContent(cfg SearchContentConfig) string {
	var contentLines []string

	// SEARCH MODE: Show search interface
	if cfg.SearchMode {

		// Build the input prompt
		var prompt string
		if cfg.IsActive {
			prompt = cfg.SearchInput + "█"
		} else {
			prompt = cfg.SearchInput
		}
		if prompt == "" || prompt == "█" {
			prompt = " " // Show at least a space so the box isn't empty
		}

		// Determine label based on title
		label := "Search"
		if strings.Contains(cfg.Title, "Challenge Rating") {
			label = "CR"
		} else if strings.Contains(cfg.Title, "Spell Level") {
			label = "Level"
		}

		contentLines = append(contentLines, label+": "+prompt)
		contentLines = append(contentLines, "↑↓ Enter Esc")

		// Show suggestions with scrolling window
		if len(cfg.Suggestions) > 0 {
			contentLines = append(contentLines, "Suggestions:")

			// Show a scrolling window of suggestions (max 8 visible to prevent panel stretching)
			maxVisible := 5
			totalSuggestions := len(cfg.Suggestions)

			// Calculate which suggestions to show
			startIdx := 0
			endIdx := totalSuggestions

			// Only apply scrolling if there are more suggestions than visible space
			if totalSuggestions > maxVisible {
				// Keep selected item visible, preferring to show it in the middle
				startIdx = cfg.SuggestionIndex - (maxVisible / 2)

				// Ensure we don't go below 0
				if startIdx < 0 {
					startIdx = 0
				}

				// Calculate end index
				endIdx = startIdx + maxVisible

				// If we've gone past the end, adjust both indices
				if endIdx > totalSuggestions {
					endIdx = totalSuggestions
					startIdx = endIdx - maxVisible
					if startIdx < 0 {
						startIdx = 0
					}
				}
			}

			// Show scroll indicator at top if there are more items above
			if totalSuggestions > maxVisible && startIdx > 0 {
				scrollText := fmt.Sprintf("  ⬆ %s (%s more above)", strings.Repeat("─", 8), formatSuggestionCount(startIdx))
				contentLines = append(contentLines, scrollText)
			}

			// Show the visible window of suggestions
			for i := startIdx; i < endIdx; i++ {
				suggestion := cfg.Suggestions[i]
				if i == cfg.SuggestionIndex {
					contentLines = append(contentLines, "► "+suggestion)
				} else {
					contentLines = append(contentLines, "  "+suggestion)
				}
			}

			// Show scroll indicator at bottom if there are more items below
			if totalSuggestions > maxVisible && endIdx < totalSuggestions {
				remaining := totalSuggestions - endIdx
				scrollText := fmt.Sprintf("  ⬇ %s (%s more below)", strings.Repeat("─", 8), formatSuggestionCount(remaining))
				contentLines = append(contentLines, scrollText)
			}

			contentLines = append(contentLines, "")
		}

		return strings.Join(contentLines, "\n")
	}

	// NOT IN SEARCH MODE: Show item details or default message

	// Header for non-search mode
	contentLines = append(contentLines, cfg.Title)
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 35))
	contentLines = append(contentLines, "")

	// Show selected item details
	if cfg.SelectedItem != nil && cfg.FormatFunc != nil {
		itemDetails := cfg.FormatFunc(cfg.SelectedItem)
		if itemDetails != "" {
			// Split the details into lines and add them
			detailLines := strings.Split(itemDetails, "\n")
			for _, line := range detailLines {
				contentLines = append(contentLines, line)
			}
		}
	} else {
		contentLines = append(contentLines, "")
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
