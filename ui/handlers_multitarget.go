// ui/handlers_multitarget.go
package ui

import (
	"strconv"
)

// handleT handles the 't' key to toggle multi-target mode
func handleT(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input
	if m.IsInputMode() {
		return handleSearchModeInput(m, "t"), nil
	}

	// Only work in Initiative Tracker
	if m.ActivePanel != InitiativeTracker {
		return m, nil
	}

	// Check if this is Shift+T (capital T) for temp HP
	isShiftT := msg.String() == KeyShiftT
	if isShiftT && m.InitiativeListMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 && m.InitiativeList[originalIndex].Type == "monster" {
			// Enter temp HP edit mode
			m.InitiativeEditMode = true
			m.InitiativeEditType = "temphp"
			m.InitiativeInput = ""
			return m, nil
		}
	}

	// Need to have entries in initiative list for multi-target
	if len(m.InitiativeList) == 0 {
		return m, nil
	}

	// Toggle multi-target mode (only if not in edit mode and lowercase t)
	if !m.InitiativeEditMode && !isShiftT {
		if !m.MultiTargetMode {
			// Enter multi-target mode
			m.MultiTargetMode = true
			m.InitiativeListMode = true // Ensure list mode is active

			// Ensure SelectedEntry is valid
			if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
				m.SelectedEntry = 0 // Start with first entry selected
			}

			m.SelectedTargets = make(map[int]bool)
			m.TargetSaveResults = make(map[int]string)
			// Exit any edit modes
			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		} else {
			// Exit multi-target mode but stay in list mode
			m.MultiTargetMode = false
			m.SelectedTargets = make(map[int]bool)
			m.ShowMultiTargetPopup = false
			m.MultiTargetInput = ""
			m.TargetSaveResults = make(map[int]string)
			// Stay in list mode with current selection
		}
	}

	return m, nil
}

// handleMultiTargetSpace handles space key in multi-target mode
func handleMultiTargetSpace(m Model) (Model, Cmd) {
	// Safety checks - return early if conditions not met
	if !m.MultiTargetMode {
		return m, nil
	}

	if m.SelectedEntry < 0 {
		return m, nil
	}

	if m.SelectedEntry >= len(m.InitiativeList) {
		return m, nil
	}

	// Initialize maps if nil (defensive programming)
	if m.SelectedTargets == nil {
		m.SelectedTargets = make(map[int]bool)
	}
	if m.TargetSaveResults == nil {
		m.TargetSaveResults = make(map[int]string)
	}

	// Toggle selection for current entry
	currentlySelected := m.SelectedTargets[m.SelectedEntry]
	if currentlySelected {
		// Deselect
		delete(m.SelectedTargets, m.SelectedEntry)
		delete(m.TargetSaveResults, m.SelectedEntry)
	} else {
		// Select
		m.SelectedTargets[m.SelectedEntry] = true
	}

	// IMPORTANT: Return the modified model, not the original
	return m, nil
}

// handleMultiTargetApply handles applying multi-target damage/healing
func handleMultiTargetApply(m Model) (Model, Cmd) {
	// Only if in multi-target mode with targets selected
	if !m.MultiTargetMode || len(m.SelectedTargets) == 0 {
		return m, nil
	}

	// Open the multi-target popup
	m.ShowMultiTargetPopup = true
	m.MultiTargetInput = ""
	m.MultiTargetType = "damage" // Default to damage

	return m, nil
}

// handleMultiTargetPopupInput handles input in the multi-target popup
func handleMultiTargetPopupInput(m Model, msg KeyMsg) (Model, Cmd) {
	if !m.ShowMultiTargetPopup {
		return m, nil
	}

	switch msg.String() {
	case "enter":
		// Apply damage/healing
		if m.MultiTargetInput != "" {
			// Parse input - supports +10 for healing, -10 for damage, or plain 10
			inputStr := m.MultiTargetInput
			var amount int
			var err error

			// Determine type from sign if present
			if len(inputStr) > 0 && (inputStr[0] == '+' || inputStr[0] == '-') {
				amount, err = strconv.Atoi(inputStr[1:]) // Parse number without sign
				if err == nil {
					if inputStr[0] == '+' {
						m.MultiTargetType = "healing"
					} else {
						m.MultiTargetType = "damage"
					}
				}
			} else {
				amount, err = strconv.Atoi(inputStr)
			}

			if err == nil && amount > 0 {
				// Check if save mode is on and all targets have save results
				canApply := true
				if m.MultiTargetSaveMode {
					for i := range m.SelectedTargets {
						if m.TargetSaveResults[i] == "" {
							canApply = false
							break
						}
					}
				}

				if canApply {
					m = ApplyMultiTargetDamage(m, amount)
					// Close popup and exit multi-target mode
					m.ShowMultiTargetPopup = false
					m.MultiTargetMode = false
					m.SelectedTargets = make(map[int]bool)
					m.MultiTargetInput = ""
					m.MultiTargetSaveMode = false
					m.TargetSaveResults = make(map[int]string)
				}
			}
		}

	case "esc":
		// Cancel
		m.ShowMultiTargetPopup = false

	case "x":
		// Toggle save mode
		m.MultiTargetSaveMode = !m.MultiTargetSaveMode
		// Clear save results when toggling
		m.TargetSaveResults = make(map[int]string)

	case "s":
		// Mark current target as save success (only if save mode is on)
		if m.MultiTargetSaveMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			if m.SelectedTargets[m.SelectedEntry] {
				m.TargetSaveResults[m.SelectedEntry] = "success"
			}
		}

	case "f":
		// Mark current target as save failure (only if save mode is on)
		if m.MultiTargetSaveMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			if m.SelectedTargets[m.SelectedEntry] {
				m.TargetSaveResults[m.SelectedEntry] = "failure"
			}
		}

	case "up":
		// Navigate to previous target
		if m.SelectedEntry > 0 {
			m.SelectedEntry--
		}

	case "down":
		// Navigate to next target
		if m.SelectedEntry < len(m.InitiativeList)-1 {
			m.SelectedEntry++
		}

	case "h":
		// Toggle between damage and healing
		if m.MultiTargetType == "damage" {
			m.MultiTargetType = "healing"
		} else {
			m.MultiTargetType = "damage"
		}

	case "backspace":
		if len(m.MultiTargetInput) > 0 {
			m.MultiTargetInput = m.MultiTargetInput[:len(m.MultiTargetInput)-1]
			// Update type based on remaining input
			if len(m.MultiTargetInput) > 0 && m.MultiTargetInput[0] == '+' {
				m.MultiTargetType = "healing"
			} else if len(m.MultiTargetInput) > 0 && m.MultiTargetInput[0] == '-' {
				m.MultiTargetType = "damage"
			}
		}

	default:
		// Add character to input if it's a digit, + or -
		if len(msg.String()) == 1 {
			char := msg.String()[0]
			// Allow + or - only at the start
			if (char == '+' || char == '-') && m.MultiTargetInput == "" {
				m.MultiTargetInput += msg.String()
				// Set mode based on sign
				if char == '+' {
					m.MultiTargetType = "healing"
				} else {
					m.MultiTargetType = "damage"
				}
			} else if char >= '0' && char <= '9' {
				m.MultiTargetInput += msg.String()
			}
		}
	}

	return m, nil
}
