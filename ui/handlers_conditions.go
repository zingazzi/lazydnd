// ui/handlers_conditions.go
package ui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// handleO handles the 'o' key to open condition management (cOnditions)
func handleO(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "o"), nil
	}

	// Only work in Initiative Tracker with list mode active
	if m.ActivePanel != InitiativeTracker || !m.InitiativeListMode {
		return m, nil
	}

	// If in multi-target mode, check if any targets are selected
	if m.MultiTargetMode {
		if len(m.SelectedTargets) == 0 {
			return m, nil
		}
	} else {
		// Single target mode - need a selected entry
		if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
			return m, nil
		}
	}

	// Open condition popup in list mode
	m.ShowConditionPopup = true
	m.ConditionPopupMode = "list"
	m.SelectedConditionIdx = 0
	m.ConditionInput = ""
	m.ConditionDurationInput = ""
	m.ConditionInputStep = 0

	return m, nil
}

// handleConditionPopupInput handles input in the condition popup
func handleConditionPopupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.ShowConditionPopup {
		return m, nil
	}

	key := msg.String()

	// List mode
	if m.ConditionPopupMode == "list" {
		switch key {
		case "esc":
			// Close popup
			m.ShowConditionPopup = false
			m.ConditionPopupMode = ""
			m.SelectedConditionIdx = 0
			return m, nil

		case "a":
			// Switch to add mode
			m.ConditionPopupMode = "add"
			m.ConditionInputStep = 0
			m.ConditionInput = ""
			m.ConditionDurationInput = ""
			m.SelectedConditionNameIdx = 0 // Start at first condition in list
			return m, nil

		case "up":
			// Navigate up in condition list
			if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
				if m.SelectedConditionIdx > 0 {
					m.SelectedConditionIdx--
				}
			}
			return m, nil

		case "down":
			// Navigate down in condition list
			if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
				entry := m.InitiativeList[m.SelectedEntry]
				if m.SelectedConditionIdx < len(entry.Conditions)-1 {
					m.SelectedConditionIdx++
				}
			}
			return m, nil

		case "d":
			// Delete selected condition
			if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
				entry := &m.InitiativeList[m.SelectedEntry]
				if m.SelectedConditionIdx >= 0 && m.SelectedConditionIdx < len(entry.Conditions) {
					// Remove condition
					entry.Conditions = append(entry.Conditions[:m.SelectedConditionIdx],
						entry.Conditions[m.SelectedConditionIdx+1:]...)
					// Adjust index
					if m.SelectedConditionIdx > 0 && m.SelectedConditionIdx >= len(entry.Conditions) {
						m.SelectedConditionIdx = len(entry.Conditions) - 1
					}
					if len(entry.Conditions) == 0 {
						m.SelectedConditionIdx = 0
					}
				}
			}
			return m, nil
		}
	}

	// Add mode
	if m.ConditionPopupMode == "add" {
		switch key {
		case "esc":
			// Cancel and return to list mode
			m.ConditionPopupMode = "list"
			m.ConditionInput = ""
			m.ConditionDurationInput = ""
			m.ConditionInputStep = 0
			m.SelectedConditionNameIdx = 0
			return m, nil

		case "up":
			// Navigate up in condition selection
			if m.ConditionInputStep == 0 {
				if m.SelectedConditionNameIdx > 0 {
					m.SelectedConditionNameIdx--
				}
			}
			return m, nil

		case "down":
			// Navigate down in condition selection
			if m.ConditionInputStep == 0 {
				maxIdx := len(CommonConditions) // +1 for "Custom" option
				if m.SelectedConditionNameIdx < maxIdx {
					m.SelectedConditionNameIdx++
				}
			}
			return m, nil

		case "enter":
			if m.ConditionInputStep == 0 {
				// Condition selected from list
				if m.SelectedConditionNameIdx < len(CommonConditions) {
					// Selected a common condition, move to duration
					m.ConditionInput = CommonConditions[m.SelectedConditionNameIdx]
					m.ConditionInputStep = 1
				} else {
					// Selected custom condition, ask for name
					m.ConditionInputStep = 2
					m.ConditionInput = ""
				}
		} else if m.ConditionInputStep == 1 {
			// Duration entered, apply condition
			duration := 0
			if m.ConditionDurationInput != "" {
				duration, _ = strconv.Atoi(m.ConditionDurationInput)
				if duration < 0 {
					duration = 0
				}
			}

			newCondition := Condition{
				Name:        m.ConditionInput,
				RoundsLeft:  duration,
				TotalRounds: duration,
				Description: "",
			}

			// Apply to all selected targets in multi-target mode, or single target
			if m.MultiTargetMode {
				// Apply to all selected targets
				for idx := range m.SelectedTargets {
					if idx >= 0 && idx < len(m.InitiativeList) {
						m.InitiativeList[idx].Conditions = append(
							m.InitiativeList[idx].Conditions,
							newCondition,
						)
					}
				}
			} else {
				// Apply to single selected entry
				if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
					m.InitiativeList[m.SelectedEntry].Conditions = append(
						m.InitiativeList[m.SelectedEntry].Conditions,
						newCondition,
					)
				}
			}

			// Reset and return to list mode
			m.ConditionPopupMode = "list"
			m.ConditionInput = ""
			m.ConditionDurationInput = ""
			m.ConditionInputStep = 0
			m.SelectedConditionNameIdx = 0

			// Set selected condition index based on mode
			if !m.MultiTargetMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
				m.SelectedConditionIdx = len(m.InitiativeList[m.SelectedEntry].Conditions) - 1
			}
			} else if m.ConditionInputStep == 2 {
				// Custom name entered, move to duration
				if m.ConditionInput != "" {
					m.ConditionInputStep = 1
				}
			}
			return m, nil

		case "backspace":
			if m.ConditionInputStep == 1 {
				if len(m.ConditionDurationInput) > 0 {
					m.ConditionDurationInput = m.ConditionDurationInput[:len(m.ConditionDurationInput)-1]
				}
			} else if m.ConditionInputStep == 2 {
				if len(m.ConditionInput) > 0 {
					m.ConditionInput = m.ConditionInput[:len(m.ConditionInput)-1]
				}
			}
			return m, nil

		default:
			// Handle text input
			if len(key) == 1 {
				if m.ConditionInputStep == 1 {
					// Only accept digits for duration
					if key >= "0" && key <= "9" {
						m.ConditionDurationInput += key
					}
				} else if m.ConditionInputStep == 2 {
					// Accept any character for custom condition name
					m.ConditionInput += key
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// UpdateConditionDurations decrements all condition durations by one round
func UpdateConditionDurations(m Model) Model {
	for i := range m.InitiativeList {
		entry := &m.InitiativeList[i]
		var activeConditions []Condition

		for _, condition := range entry.Conditions {
			if condition.RoundsLeft > 0 {
				condition.RoundsLeft--
				if condition.RoundsLeft > 0 {
					activeConditions = append(activeConditions, condition)
				}
				// If RoundsLeft reaches 0, the condition expires
			} else {
				// Indefinite duration (RoundsLeft == 0), keep it
				activeConditions = append(activeConditions, condition)
			}
		}

		entry.Conditions = activeConditions
	}

	return m
}
