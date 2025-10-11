// ui/handlers_cr_filter.go
package ui

import (
	"lazydnd/panels"

	tea "github.com/charmbracelet/bubbletea"
)

// handleF toggles CR filter mode in monster panel OR level filter in spells panel
func handleF(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.isInInputMode() {
		return handleSearchModeInput(m, "f"), nil
	}
	
	// Monster panel: CR filter
	if m.ActivePanel == Monsters && !m.MonsterSearchMode && !m.MonsterCRFilterMode {
		// Clear any selected monster and enter CR filter mode
		m.SelectedMonster = nil
		m.MonsterCRFilterMode = true
		m.MonsterCRFilter = ""
		m.MonsterSuggestions = []string{}
		m.MonsterSuggestionIndex = -1
		DebugLog("FILTER: Entering CR filter mode")
		return m, nil
	}
	
	// Spells panel: Level filter
	if m.ActivePanel == Spells && !m.SpellSearchMode && !m.SpellLevelFilterMode && !m.ActiveSpellListMode && !m.CastSpellInputMode {
		// Clear any selected spell and enter level filter mode
		m.SelectedSpell = nil
		m.SpellLevelFilterMode = true
		m.SpellLevelFilter = ""
		m.SpellSuggestions = []string{}
		m.SuggestionIndex = -1
		DebugLog("FILTER: Entering spell level filter mode")
		return m, nil
	}
	
	return m, nil
}

// handleCRFilterInput handles text input for CR filter
func handleCRFilterInput(m Model, key string) (Model, tea.Cmd) {
	if !m.MonsterCRFilterMode {
		return m, nil
	}

	// Handle up/down arrow keys for navigation
	if key == "up" {
		if len(m.MonsterSuggestions) > 0 {
			m.MonsterSuggestionIndex--
			if m.MonsterSuggestionIndex < 0 {
				m.MonsterSuggestionIndex = len(m.MonsterSuggestions) - 1
			}
		}
		return m, nil
	}

	if key == "down" {
		if len(m.MonsterSuggestions) > 0 {
			m.MonsterSuggestionIndex++
			if m.MonsterSuggestionIndex >= len(m.MonsterSuggestions) {
				m.MonsterSuggestionIndex = 0
			}
		}
		return m, nil
	}

	// Handle backspace
	if key == "backspace" {
		if len(m.MonsterCRFilter) > 0 {
			m.MonsterCRFilter = m.MonsterCRFilter[:len(m.MonsterCRFilter)-1]
			DebugLog("FILTER: Backspace, filter='%s'", m.MonsterCRFilter)
			// Update monster list immediately
			if m.MonsterCRFilter != "" {
				m.MonsterSuggestions = panels.SearchMonsters("", m.MonsterCRFilter)
			} else {
				m.MonsterSuggestions = []string{}
			}
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		}
		return m, nil
	}

	// Handle enter to select monster
	if key == "enter" {
		if len(m.MonsterSuggestions) > 0 && m.MonsterSuggestionIndex >= 0 {
			// Select the monster and show details
			selectedMonsterName := m.MonsterSuggestions[m.MonsterSuggestionIndex]
			// Store the selected monster name (rendering code will find the full monster)
			m.SelectedMonster = &Monster{Name: selectedMonsterName}
			m.MonsterCRFilterMode = false
			m.MonsterCRFilter = ""
			m.MonsterSuggestions = []string{}
			m.MonsterSuggestionIndex = -1
			DebugLog("FILTER: Selected monster '%s'", selectedMonsterName)
		}
		return m, nil
	}

	// Handle escape to cancel
	if key == "esc" {
		m.MonsterCRFilterMode = false
		m.MonsterCRFilter = ""
		m.MonsterSuggestions = []string{}
		m.MonsterSuggestionIndex = -1
		DebugLog("FILTER: Cancelled CR filter mode")
		return m, nil
	}

	// Handle valid input: numbers, dash, plus, slash, dot
	if len(key) == 1 {
		if (key >= "0" && key <= "9") || key == "-" || key == "+" || key == "/" || key == "." {
			m.MonsterCRFilter += key
			DebugLog("FILTER: Added char, filter='%s'", m.MonsterCRFilter)
			// Update monster list immediately
			m.MonsterSuggestions = panels.SearchMonsters("", m.MonsterCRFilter)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		}
	}

	return m, nil
}
