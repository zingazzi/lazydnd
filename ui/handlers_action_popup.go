// ui/handlers_action_popup.go
package ui

import tea "github.com/charmbracelet/bubbletea"

// handleActionPopupInput handles input specifically for the action popup
func handleActionPopupInput(m Model, key string) (Model, tea.Cmd) {
	// Close popup on Escape
	if key == "esc" {
		m.ShowActionPopup = false
		m.ActionPopupActions = []MonsterAction{}
		m.ActionPopupIndex = 0
		m.ActionPopupMonster = ""
		m.ActionPopupAdvantage = false
		m.ActionPopupDisadvantage = false
		return m, nil
	}

	// Toggle advantage
	if key == "a" {
		if m.ActionPopupAdvantage {
			m.ActionPopupAdvantage = false // Turn off
		} else {
			m.ActionPopupAdvantage = true
			m.ActionPopupDisadvantage = false // Can't have both
		}
		return m, nil
	}

	// Toggle disadvantage
	if key == "d" {
		if m.ActionPopupDisadvantage {
			m.ActionPopupDisadvantage = false // Turn off
		} else {
			m.ActionPopupDisadvantage = true
			m.ActionPopupAdvantage = false // Can't have both
		}
		return m, nil
	}

	// Navigate actions
	if key == "up" {
		if m.ActionPopupIndex > 0 {
			m.ActionPopupIndex--
		}
		return m, nil
	}

	if key == "down" {
		if m.ActionPopupIndex < len(m.ActionPopupActions)-1 {
			m.ActionPopupIndex++
		}
		return m, nil
	}

	// For Enter key, let it pass through to handleEnter
	// For any other key, don't consume it
	return m, nil
}
