// ui/handlers_quick_hp.go
package ui

import (
	"lazydnd/panels"

	tea "github.com/charmbracelet/bubbletea"
)

// handleQuickAddHP opens the quick add HP popup
func handleQuickAddHP(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Only in initiative tracker when not in edit mode
	if m.ActivePanel == InitiativeTracker && !m.InitiativeEditMode && !m.InitiativeInputMode {
		// Check if we have targets (either selected or in multi-target mode)
		hasTargets := false
		if m.MultiTargetMode {
			for _, selected := range m.SelectedTargets {
				if selected {
					hasTargets = true
					break
				}
			}
		} else if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			hasTargets = true
		}

		if hasTargets {
			m.ShowQuickHPPopup = true
			m.QuickHPInput = ""
			m.QuickHPMode = "add"
		}
	}
	return m, nil
}

// handleQuickRemoveHP opens the quick remove HP popup
func handleQuickRemoveHP(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Only in initiative tracker when not in edit mode
	if m.ActivePanel == InitiativeTracker && !m.InitiativeEditMode && !m.InitiativeInputMode {
		// Check if we have targets (either selected or in multi-target mode)
		hasTargets := false
		if m.MultiTargetMode {
			for _, selected := range m.SelectedTargets {
				if selected {
					hasTargets = true
					break
				}
			}
		} else if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			hasTargets = true
		}

		if hasTargets {
			m.ShowQuickHPPopup = true
			m.QuickHPInput = ""
			m.QuickHPMode = "remove"
		}
	}
	return m, nil
}

// handleQuickHPInput handles input in the quick HP popup
func handleQuickHPInput(m Model, key string) (Model, tea.Cmd) {
	if key == "esc" {
		// Cancel and close popup
		m.ShowQuickHPPopup = false
		m.QuickHPInput = ""
		m.QuickHPMode = ""
		return m, nil
	}

	if key == "enter" {
		// Apply HP change
		if m.QuickHPInput == "" {
			return m, nil
		}

		// Parse the input
		var hpChange int
		val, err := panels.ParseInput(m.QuickHPInput, "hp_change")
		if err != nil {
			// Show error and stay in popup
			SetError(&m, err.Error())
			return m, nil
		}
		hpChange = val.(int)

		// Adjust sign based on mode
		if m.QuickHPMode == "remove" && hpChange > 0 {
			hpChange = -hpChange
		} else if m.QuickHPMode == "add" && hpChange < 0 {
			hpChange = -hpChange
		}

		// Apply HP change to all selected targets
		if m.MultiTargetMode {
			// Multi-target mode - apply to all selected
			for i, selected := range m.SelectedTargets {
				if selected && i >= 0 && i < len(m.InitiativeList) {
					// Save current HP
					oldHP := m.InitiativeList[i].HP
					oldTempHP := m.InitiativeList[i].TempHP

					// Update HP using the HP calculator
					newHP, newTempHP, _ := HPCalc.CalculateHPChange(
						m.InitiativeList[i].HP,
						m.InitiativeList[i].MaxHP,
						m.InitiativeList[i].TempHP,
						hpChange,
					)

					m.InitiativeList[i].HP = newHP
					m.InitiativeList[i].TempHP = newTempHP

					// Save to undo history
					if m.InitiativeList[i].HP != oldHP || m.InitiativeList[i].TempHP != oldTempHP {
						pushHPHistory(&m, i, oldHP, newHP)
					}
				}
			}

			// Exit multi-target mode after applying
			m.MultiTargetMode = false
			m.SelectedTargets = make(map[int]bool)
		} else {
			// Single target mode
			if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
				// Save current HP to undo stack
				oldHP := m.InitiativeList[m.SelectedEntry].HP
				oldTempHP := m.InitiativeList[m.SelectedEntry].TempHP

				// Update HP using the HP calculator
				newHP, newTempHP, _ := HPCalc.CalculateHPChange(
					m.InitiativeList[m.SelectedEntry].HP,
					m.InitiativeList[m.SelectedEntry].MaxHP,
					m.InitiativeList[m.SelectedEntry].TempHP,
					hpChange,
				)

				m.InitiativeList[m.SelectedEntry].HP = newHP
				m.InitiativeList[m.SelectedEntry].TempHP = newTempHP

				// Save to undo history
				if m.InitiativeList[m.SelectedEntry].HP != oldHP || m.InitiativeList[m.SelectedEntry].TempHP != oldTempHP {
					pushHPHistory(&m, m.SelectedEntry, oldHP, newHP)
				}
			}
		}

		// Close popup
		m.ShowQuickHPPopup = false
		m.QuickHPInput = ""
		m.QuickHPMode = ""
		return m, nil
	}

	// Handle backspace
	if key == "backspace" || key == "ctrl+h" {
		if len(m.QuickHPInput) > 0 {
			m.QuickHPInput = m.QuickHPInput[:len(m.QuickHPInput)-1]
		}
		return m, nil
	}

	// Handle number input and minus sign
	if len(key) == 1 {
		ch := key[0]
		if (ch >= '0' && ch <= '9') || ch == '-' {
			// Limit length
			if len(m.QuickHPInput) < 5 {
				m.QuickHPInput += key
			}
		}
	}

	return m, nil
}
