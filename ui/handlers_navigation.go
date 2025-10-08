// ui/handlers/navigation.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ========== NAVIGATION HANDLERS ==========

// handleTab handles tab key navigation (forward)
func handleTab(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.InputMode {
		m.ActivePanel = (m.ActivePanel + 1) % 4
	}
	return m, nil
}

// handleShiftTab handles shift+tab key navigation (backward)
func handleShiftTab(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.InputMode {
		m.ActivePanel = (m.ActivePanel - 1 + 4) % 4
	}
	return m, nil
}

// handleUp handles up arrow key
func handleUp(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle action popup navigation first (highest priority)
	if m.ShowActionPopup && len(m.ActionPopupActions) > 0 {
		if m.ActionPopupIndex > 0 {
			m.ActionPopupIndex--
		}
		return m, nil
	}

	// Handle dice history navigation
	if m.ActivePanel == DiceRoller && m.DiceHistoryMode && len(m.DiceHistory) > 0 {
		if m.HistoryIndex < len(m.DiceHistory)-1 {
			m.HistoryIndex++
		}
		return m, nil
	}

	// Handle spell suggestion navigation
	if m.ActivePanel == Spells && m.SpellSearchMode && len(m.SpellSuggestions) > 0 {
		if m.SuggestionIndex > 0 {
			m.SuggestionIndex--
		}
	} else if m.ActivePanel == Monsters && m.MonsterSearchMode && len(m.MonsterSuggestions) > 0 {
		// Navigate monster suggestions
		if m.MonsterSuggestionIndex > 0 {
			m.MonsterSuggestionIndex--
		}
	} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && len(m.InitiativeList) > 0 {
		// Navigate initiative list
		if m.SelectedEntry > 0 {
			m.SelectedEntry--
		}
	} else if !m.InputMode && !m.InitiativeInputMode {
		// Normal panel scrolling when not in input mode (allow even in spell search mode if no suggestions)
		if m.ScrollOffset[m.ActivePanel] > 0 {
			m.ScrollOffset[m.ActivePanel]--
		}
	}
	return m, nil
}

// handleDown handles down arrow key
func handleDown(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle action popup navigation first (highest priority)
	if m.ShowActionPopup && len(m.ActionPopupActions) > 0 {
		if m.ActionPopupIndex < len(m.ActionPopupActions)-1 {
			m.ActionPopupIndex++
		}
		return m, nil
	}

	// Handle dice history navigation
	if m.ActivePanel == DiceRoller && m.DiceHistoryMode && len(m.DiceHistory) > 0 {
		if m.HistoryIndex > 0 {
			m.HistoryIndex--
		}
		return m, nil
	}

	// Handle spell suggestion navigation
	if m.ActivePanel == Spells && m.SpellSearchMode && len(m.SpellSuggestions) > 0 {
		if m.SuggestionIndex < len(m.SpellSuggestions)-1 {
			m.SuggestionIndex++
		}
	} else if m.ActivePanel == Monsters && m.MonsterSearchMode && len(m.MonsterSuggestions) > 0 {
		// Navigate monster suggestions
		if m.MonsterSuggestionIndex < len(m.MonsterSuggestions)-1 {
			m.MonsterSuggestionIndex++
		}
	} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && len(m.InitiativeList) > 0 {
		// Navigate initiative list
		if m.SelectedEntry < len(m.InitiativeList)-1 {
			m.SelectedEntry++
		}
	} else if !m.InputMode && !m.InitiativeInputMode {
		// Normal panel scrolling when not in input mode (allow even in spell search mode if no suggestions)
		m.ScrollOffset[m.ActivePanel]++
	}
	return m, nil
}

// ========== FUNCTION KEY HANDLERS ==========

// handleF1 switches to dice roller panel
func handleF1(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	m.ActivePanel = DiceRoller
	m.InputMode = false
	return m, nil
}

// handleF2 switches to initiative tracker panel
func handleF2(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	m.ActivePanel = InitiativeTracker
	m.InputMode = false
	return m, nil
}

// handleF3 switches to spells panel
func handleF3(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	m.ActivePanel = Spells
	m.InputMode = false
	return m, nil
}

// handleF4 switches to monsters panel
func handleF4(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	m.ActivePanel = Monsters
	m.InputMode = false
	return m, nil
}

// ========== NUMBER KEY HANDLERS ==========

// handleNumber1 handles the '1' key
func handleNumber1(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InputMode {
		m.DiceInput += "1"
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		m.InitiativeInput += "1"
	} else {
		m.ActivePanel = DiceRoller
	}
	return m, nil
}

// handleNumber2 handles the '2' key
func handleNumber2(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InputMode {
		m.DiceInput += "2"
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		m.InitiativeInput += "2"
	} else {
		m.ActivePanel = InitiativeTracker
	}
	return m, nil
}

// handleNumber3 handles the '3' key
func handleNumber3(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InputMode {
		m.DiceInput += "3"
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		m.InitiativeInput += "3"
	} else {
		m.ActivePanel = Spells
	}
	return m, nil
}

// handleNumber4 handles the '4' key
func handleNumber4(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InputMode {
		m.DiceInput += "4"
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		m.InitiativeInput += "4"
	} else {
		m.ActivePanel = Monsters
	}
	return m, nil
}

// ========== TURN TRACKING HANDLERS ==========

// handleResetCombat resets the combat turn and round counter
func handleResetCombat(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// If in any input mode, pass through to default input handler
	if m.InputMode || m.InitiativeInputMode || m.InitiativeEditMode || m.SpellSearchMode || m.MonsterSearchMode {
		return handleDefaultInput(m, msg)
	}

	// Only works in initiative tracker panel when not in list mode
	if m.ActivePanel != InitiativeTracker || m.InitiativeListMode {
		return m, nil
	}

	// Reset combat state
	m.CurrentTurn = -1
	m.RoundCounter = 0

	return m, nil
}

// handleNextTurn advances to the next turn in initiative order
func handleNextTurn(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// If in any input mode, pass through to default input handler
	if m.InputMode || m.InitiativeInputMode || m.InitiativeEditMode || m.SpellSearchMode || m.MonsterSearchMode {
		return handleDefaultInput(m, msg)
	}

	// Only works in initiative tracker panel when not in list mode
	if m.ActivePanel != InitiativeTracker || m.InitiativeListMode {
		return m, nil
	}

	// If there are no entries, do nothing
	if len(m.InitiativeList) == 0 {
		return m, nil
	}

	// If combat hasn't started (CurrentTurn == -1), start at 0 and begin round 1
	if m.CurrentTurn == -1 {
		m.CurrentTurn = 0
		m.RoundCounter = 1
	} else {
		// Advance to next turn
		nextTurn := (m.CurrentTurn + 1) % len(m.InitiativeList)

		// If we wrapped around to 0, increment round counter
		if nextTurn == 0 {
			m.RoundCounter++
		}

		m.CurrentTurn = nextTurn
	}

	return m, nil
}
