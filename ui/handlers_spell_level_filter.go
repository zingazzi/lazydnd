// ui/handlers_spell_level_filter.go
package ui

import (
	"lazydnd/panels"
)

// handleSpellLevelFilter is called when 'f' is pressed in Spells panel
func handleSpellLevelFilter(m Model, msg KeyMsg) (Model, Cmd) {
	// Only in Spells panel and not in other modes
	if m.ActivePanel == Spells && !m.SpellSearchMode && !m.SpellLevelFilterMode && !m.ActiveSpellListMode {
		// Clear any selected spell and enter level filter mode
		m.SelectedSpell = nil
		m.SpellLevelFilterMode = true
		m.SpellLevelFilter = ""
		m.SpellSuggestions = []string{}
		m.SuggestionIndex = -1
		DebugLog("SPELL FILTER: Entering level filter mode")
	}
	return m, nil
}

// handleSpellLevelFilterInput handles text input for spell level filter
func handleSpellLevelFilterInput(m Model, key string) (Model, Cmd) {
	if !m.SpellLevelFilterMode {
		return m, nil
	}

	// Handle up/down arrow keys for navigation
	if key == "up" {
		if len(m.SpellSuggestions) > 0 {
			m.SuggestionIndex--
			if m.SuggestionIndex < 0 {
				m.SuggestionIndex = len(m.SpellSuggestions) - 1
			}
		}
		return m, nil
	}

	if key == "down" {
		if len(m.SpellSuggestions) > 0 {
			m.SuggestionIndex++
			if m.SuggestionIndex >= len(m.SpellSuggestions) {
				m.SuggestionIndex = 0
			}
		}
		return m, nil
	}

	// Handle backspace
	if key == "backspace" {
		if len(m.SpellLevelFilter) > 0 {
			m.SpellLevelFilter = m.SpellLevelFilter[:len(m.SpellLevelFilter)-1]
			DebugLog("SPELL FILTER: Backspace, filter='%s'", m.SpellLevelFilter)
			// Update spell list immediately
			if m.SpellLevelFilter != "" {
				m.SpellSuggestions = panels.SearchSpells("", m.SpellLevelFilter)
			} else {
				m.SpellSuggestions = []string{}
			}
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		}
		return m, nil
	}

	// Handle enter to select spell
	if key == "enter" {
		if len(m.SpellSuggestions) > 0 && m.SuggestionIndex >= 0 {
			// Select the spell and show details
			selectedSpellName := m.SpellSuggestions[m.SuggestionIndex]
			// Store the selected spell name (rendering code will find the full spell)
			m.SelectedSpell = &Spell{Name: selectedSpellName}
			m.SpellLevelFilterMode = false
			m.SpellLevelFilter = ""
			m.SpellSuggestions = []string{}
			m.SuggestionIndex = -1
			DebugLog("SPELL FILTER: Selected spell '%s'", selectedSpellName)
		}
		return m, nil
	}

	// Handle escape to cancel
	if key == "esc" {
		m.SpellLevelFilterMode = false
		m.SpellLevelFilter = ""
		m.SpellSuggestions = []string{}
		m.SuggestionIndex = -1
		DebugLog("SPELL FILTER: Cancelled level filter mode")
		return m, nil
	}

	// Handle valid input: numbers, dash, plus
	if len(key) == 1 {
		if (key >= "0" && key <= "9") || key == "-" || key == "+" {
			m.SpellLevelFilter += key
			DebugLog("SPELL FILTER: Added char, filter='%s'", m.SpellLevelFilter)
			// Update spell list immediately
			m.SpellSuggestions = panels.SearchSpells("", m.SpellLevelFilter)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		}
	}

	return m, nil
}
