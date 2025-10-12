// panels/notes.go
package panels

import (
	"fmt"
	"strings"
)

// GetNotesContent returns the formatted notes panel content
func GetNotesContent(content string, editMode bool, searchMode bool, searchInput string, searchResults []int, width int) string {
	var lines []string

	// Header based on mode
	if searchMode {
		lines = append(lines, "🔍 Search Notes")
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Search: %s█", searchInput))
		lines = append(lines, "")
		if len(searchResults) > 0 {
			lines = append(lines, fmt.Sprintf("Found %d match(es)", len(searchResults)))
		} else if searchInput != "" {
			lines = append(lines, "No matches found")
		}
		lines = append(lines, "")
		lines = append(lines, strings.Repeat("─", 35))
		lines = append(lines, "")
	} else 	if editMode {
		lines = append(lines, "✏️  Edit Mode")
		lines = append(lines, "")
		lines = append(lines, "")
	} else {
		lines = append(lines, "📝 Campaign Notes")
		lines = append(lines, "")
	}

	// Display notes content
	if content == "" {
		lines = append(lines, "")
		lines = append(lines, "  (Empty)")
	} else {
		// Render notes with basic markdown-style formatting
		notesLines := strings.Split(content, "\n")
		for i, line := range notesLines {
			lineNum := i + 1
			highlighted := false

			// Check if this line is in search results
			if searchMode && len(searchResults) > 0 {
				for _, resultLine := range searchResults {
					if resultLine == lineNum {
						highlighted = true
						break
					}
				}
			}

			// Apply basic markdown-style rendering
			formatted := formatMarkdown(line)

			// Add highlight marker for search results
			if highlighted {
				formatted = "→ " + formatted
			} else if content != "" {
				formatted = "  " + formatted
			}

			lines = append(lines, formatted)
		}
	}

	return strings.Join(lines, "\n")
}

// formatMarkdown applies basic markdown-style formatting
func formatMarkdown(line string) string {
	line = strings.TrimSpace(line)

	// Headings
	if strings.HasPrefix(line, "# ") {
		return "═══ " + strings.TrimPrefix(line, "# ") + " ═══"
	}
	if strings.HasPrefix(line, "## ") {
		return "─── " + strings.TrimPrefix(line, "## ") + " ───"
	}
	if strings.HasPrefix(line, "### ") {
		return "··· " + strings.TrimPrefix(line, "### ") + " ···"
	}

	// Bullet points
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
		return "  • " + line[2:]
	}

	// Bold (simplified - just remove markers for terminal)
	line = strings.ReplaceAll(line, "**", "")

	// Italic (simplified - just remove markers for terminal)
	line = strings.ReplaceAll(line, "*", "")

	return line
}

// SearchNotes finds lines matching the search query
func SearchNotes(content string, query string) []int {
	if query == "" {
		return []int{}
	}

	var results []int
	lines := strings.Split(content, "\n")
	queryLower := strings.ToLower(query)

	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), queryLower) {
			results = append(results, i+1) // Line numbers are 1-based
		}
	}

	return results
}
