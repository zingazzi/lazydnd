// ui/handlers_notes.go
package ui

import (
	"lazydnd/panels"
)

// handleE handles the 'e' key - for edit mode in Notes panel
func handleNotesE(m Model, msg KeyMsg) (Model, Cmd) {
	if m.ActivePanel != Notes {
		return m, nil
	}

	// Enter edit mode
	if !m.NotesEditMode && !m.NotesSearchMode {
		m.NotesEditMode = true
		m.NotesInput = m.NotesContent
	}

	return m, nil
}

// handleNotesF handles the 'f' key - for search mode in Notes panel
func handleNotesF(m Model, msg KeyMsg) (Model, Cmd) {
	if m.ActivePanel != Notes {
		return m, nil
	}

	// Enter search mode
	if !m.NotesEditMode && !m.NotesSearchMode {
		m.NotesSearchMode = true
		m.NotesSearchInput = ""
		m.NotesSearchResult = []int{}
	}

	return m, nil
}

// handleNotesInput handles text input in Notes panel
func handleNotesInput(m Model, key string) Model {
	if m.ActivePanel != Notes {
		return m
	}

	if m.NotesEditMode {
		// Handle text input in edit mode
		switch key {
		case "enter":
			// Add newline
			m.NotesInput += "\n"
		case "backspace", "ctrl+h":
			if len(m.NotesInput) > 0 {
				m.NotesInput = m.NotesInput[:len(m.NotesInput)-1]
			}
		default:
			// Add character if it's a single character
			if len(key) == 1 {
				m.NotesInput += key
			}
		}
	} else if m.NotesSearchMode {
		// Handle text input in search mode
		switch key {
		case "backspace", "ctrl+h":
			if len(m.NotesSearchInput) > 0 {
				m.NotesSearchInput = m.NotesSearchInput[:len(m.NotesSearchInput)-1]
				// Update search results
				m.NotesSearchResult = panels.SearchNotes(m.NotesContent, m.NotesSearchInput)
			}
		default:
			// Add character if it's a single character
			if len(key) == 1 {
				m.NotesSearchInput += key
				// Update search results
				m.NotesSearchResult = panels.SearchNotes(m.NotesContent, m.NotesSearchInput)
			}
		}
	}

	return m
}

// handleNotesEnter handles Enter key in Notes panel
func handleNotesEnter(m Model) Model {
	if m.ActivePanel != Notes {
		return m
	}

	if m.NotesEditMode {
		// Save and exit edit mode
		m.NotesContent = m.NotesInput
		m.NotesEditMode = false
		m.NotesInput = ""
	}

	return m
}

// handleNotesEscape handles Escape key in Notes panel
func handleNotesEscape(m Model) Model {
	if m.ActivePanel != Notes {
		return m
	}

	if m.NotesEditMode {
		// Cancel edit mode
		m.NotesEditMode = false
		m.NotesInput = ""
	} else if m.NotesSearchMode {
		// Cancel search mode
		m.NotesSearchMode = false
		m.NotesSearchInput = ""
		m.NotesSearchResult = []int{}
	}

	return m
}
