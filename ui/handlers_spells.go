// ui/handlers_spells.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleCastSpell handles the 'c' key to cast a spell
func handleCastSpell(m Model) (Model, tea.Cmd) {
	// Only in Spells panel when spell is selected
	if m.ActivePanel != Spells || m.SelectedSpell == nil || m.SpellSearchMode {
		return m, nil
	}

	// Check if spell has a duration (not instantaneous)
	rounds, isInstantaneous := ParseSpellDuration(m.SelectedSpell.Duration)
	if isInstantaneous || rounds == 0 {
		// Don't track instantaneous spells
		return m, nil
	}

	// Show cast spell popup
	m.ShowCastSpellPrompt = true
	m.CastSpellInputMode = true
	m.CastSpellInput = ""
	m.SpellToCast = m.SelectedSpell

	return m, nil
}

// handleViewActiveSpells handles the 'v' key to view active spells
func handleViewActiveSpells(m Model) (Model, tea.Cmd) {
	// Only in Spells panel when not in search mode
	if m.ActivePanel != Spells || m.SpellSearchMode || m.CastSpellInputMode {
		return m, nil
	}

	// Toggle active spell list mode
	m.ActiveSpellListMode = !m.ActiveSpellListMode
	if m.ActiveSpellListMode && len(m.ActiveSpells) > 0 {
		m.ActiveSpellIndex = 0
	}

	return m, nil
}

// handleDeleteActiveSpell handles the 'd' key to delete active spell
func handleDeleteActiveSpell(m Model) (Model, tea.Cmd) {
	// Only in Spells panel in active spell list mode
	if m.ActivePanel != Spells || !m.ActiveSpellListMode {
		return m, nil
	}

	if len(m.ActiveSpells) == 0 || m.ActiveSpellIndex < 0 {
		return m, nil
	}

	// Remove the selected spell
	m = RemoveActiveSpell(m, m.ActiveSpellIndex)

	// Exit list mode if no more spells
	if len(m.ActiveSpells) == 0 {
		m.ActiveSpellListMode = false
	}

	return m, nil
}

// handleCastSpellInput handles input in the cast spell popup
func handleCastSpellInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.ShowCastSpellPrompt || !m.CastSpellInputMode {
		return m, nil
	}

	switch msg.String() {
	case "enter":
		// Cast the spell if caster name is provided
		if m.CastSpellInput != "" && m.SpellToCast != nil {
			m = CastSpell(m, m.SpellToCast, m.CastSpellInput)
		}
		// Close popup
		m.ShowCastSpellPrompt = false
		m.CastSpellInputMode = false
		m.CastSpellInput = ""
		m.SpellToCast = nil

	case "esc":
		// Cancel casting
		m.ShowCastSpellPrompt = false
		m.CastSpellInputMode = false
		m.CastSpellInput = ""
		m.SpellToCast = nil

	case "backspace":
		if len(m.CastSpellInput) > 0 {
			m.CastSpellInput = m.CastSpellInput[:len(m.CastSpellInput)-1]
		}

	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.CastSpellInput += msg.String()
		}
	}

	return m, nil
}

// handleActiveSpellNavigation handles navigation in active spell list
func handleActiveSpellNavigation(m Model, key string) (Model, tea.Cmd) {
	if !m.ActiveSpellListMode || len(m.ActiveSpells) == 0 {
		return m, nil
	}

	switch key {
	case "up":
		if m.ActiveSpellIndex > 0 {
			m.ActiveSpellIndex--
		}
	case "down":
		if m.ActiveSpellIndex < len(m.ActiveSpells)-1 {
			m.ActiveSpellIndex++
		}
	}

	return m, nil
}
