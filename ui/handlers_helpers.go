// ui/handlers/helpers.go
package ui

import (
	"fmt"
	"lazydnd/panels"
	"reflect"
	"strings"
)

// ========== INPUT MODE HELPERS ==========

// isInInputMode returns true if any input mode is active
func (m Model) isInInputMode() bool {
	return m.InputMode || m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode || m.NotesEditMode || m.NotesSearchMode
}

// ========== INITIATIVE PROCESSING ==========

// completeInitiativeEntry finishes an initiative entry and resets the input state
func completeInitiativeEntry(m *Model) {
	m.InitiativeList = append(m.InitiativeList, m.TempEntry)
	m.InitiativeInputMode = false
	m.InitiativeInputType = ""
	m.InitiativeInput = ""
	m.TempEntry = InitiativeEntry{} // Reset temp entry
}

// processInitiativeInput handles the multi-step process of adding players/monsters
func processInitiativeInput(m Model) Model {
	switch m.InitiativeInputType {
	case "player_name":
		if val, err := panels.ParseInput(m.InitiativeInput, "player_name"); err == nil {
			// Store name and move to initiative input
			m.TempEntry.Name = val.(string)
			m.TempEntry.Type = "player"
			m.InitiativeInputType = "player_initiative"
			m.InitiativeInput = ""
		}

	case "player_initiative":
		if val, err := panels.ParseInput(m.InitiativeInput, "player_initiative"); err == nil {
			// Store initiative and move to AC input
			m.TempEntry.Initiative = val.(int)
			m.InitiativeInputType = "player_ac"
			m.InitiativeInput = ""
		}

	case "player_ac":
		if val, err := panels.ParseInput(m.InitiativeInput, "player_ac"); err == nil {
			// Complete player entry
			m.TempEntry.AC = val.(int)
			m.TempEntry.ReactionUsed = false // Initialize reaction as available
			completeInitiativeEntry(&m)
		}

	case "monster_name":
		if val, err := panels.ParseInput(m.InitiativeInput, "monster_name"); err == nil {
			// Store name and move to HP input
			m.TempEntry.Name = val.(string)
			m.TempEntry.Type = "monster"
			m.InitiativeInputType = "monster_hp"
			m.InitiativeInput = ""
		}

	case "monster_hp":
		if val, err := panels.ParseInput(m.InitiativeInput, "monster_hp"); err == nil {
			// Store HP and move to AC input
			m.TempEntry.HP = val.(int)
			m.TempEntry.MaxHP = val.(int)
			m.InitiativeInputType = "monster_ac"
			m.InitiativeInput = ""
		}

	case "monster_ac":
		if val, err := panels.ParseInput(m.InitiativeInput, "monster_ac"); err == nil {
			// Store AC and move to initiative input
			m.TempEntry.AC = val.(int)
			m.InitiativeInputType = "monster_initiative"
			m.InitiativeInput = ""
		}

	case "monster_initiative":
		if val, err := panels.ParseInput(m.InitiativeInput, "monster_initiative"); err == nil {
			// Complete monster entry
			m.TempEntry.Initiative = val.(int)
			m.TempEntry.ReactionUsed = false // Initialize reaction as available
			completeInitiativeEntry(&m)
		}
	}

	return m
}

// processInitiativeEdit handles editing existing initiative entries
func processInitiativeEdit(m Model) Model {
	if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
		return m
	}

	// Find the original array index for the selected display position
	originalIndex := findOriginalIndex(m, m.SelectedEntry)
	if originalIndex == -1 {
		// Invalid selection, exit edit mode
		m.InitiativeEditMode = false
		m.InitiativeEditType = ""
		m.InitiativeInput = ""
		return m
	}

	switch m.InitiativeEditType {
	case "initiative":
		if val, err := panels.ParseInput(m.InitiativeInput, "initiative"); err == nil {
			// Update initiative
			m.InitiativeList[originalIndex].Initiative = val.(int)
			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		}

	case "hp":
		if val, err := panels.ParseInput(m.InitiativeInput, "hp_change"); err == nil {
			// Save current HP to undo stack before making changes
			oldHP := m.InitiativeList[originalIndex].HP
			oldTempHP := m.InitiativeList[originalIndex].TempHP

			// Update HP (can be positive or negative)
			change := val.(int)

			if change < 0 {
				// Taking damage - apply to temp HP first
				damage := -change
				if m.InitiativeList[originalIndex].TempHP > 0 {
					if damage <= m.InitiativeList[originalIndex].TempHP {
						// All damage absorbed by temp HP
						m.InitiativeList[originalIndex].TempHP -= damage
					} else {
						// Temp HP absorbed some, rest goes to real HP
						remainingDamage := damage - m.InitiativeList[originalIndex].TempHP
						m.InitiativeList[originalIndex].TempHP = 0
						m.InitiativeList[originalIndex].HP -= remainingDamage
						if m.InitiativeList[originalIndex].HP < 0 {
							m.InitiativeList[originalIndex].HP = 0
						}
					}
				} else {
					// No temp HP, damage goes directly to real HP
					m.InitiativeList[originalIndex].HP += change // change is negative
					if m.InitiativeList[originalIndex].HP < 0 {
						m.InitiativeList[originalIndex].HP = 0
					}
				}
			} else {
				// Healing - only affects real HP, not temp HP
				newHP := m.InitiativeList[originalIndex].HP + change
				if newHP > m.InitiativeList[originalIndex].MaxHP {
					newHP = m.InitiativeList[originalIndex].MaxHP
				}
				m.InitiativeList[originalIndex].HP = newHP
			}

			// Save to undo history (only if HP changed)
			if m.InitiativeList[originalIndex].HP != oldHP || m.InitiativeList[originalIndex].TempHP != oldTempHP {
				pushHPHistory(&m, originalIndex, oldHP, m.InitiativeList[originalIndex].HP)
			}

			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		}

	case "maxhp":
		if val, err := panels.ParseInput(m.InitiativeInput, "maxhp"); err == nil {
			// Update Max HP (absolute value)
			newMaxHP := val.(int)
			if newMaxHP < 1 {
				newMaxHP = 1 // Ensure max HP is at least 1
			}
			m.InitiativeList[originalIndex].MaxHP = newMaxHP

			// Cap current HP if it exceeds new max HP
			if m.InitiativeList[originalIndex].HP > newMaxHP {
				m.InitiativeList[originalIndex].HP = newMaxHP
			}

			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		}

	case "temphp":
		if val, err := panels.ParseInput(m.InitiativeInput, "temphp"); err == nil {
			// Set Temp HP (absolute value, replaces existing temp HP)
			newTempHP := val.(int)
			if newTempHP < 0 {
				newTempHP = 0 // Temp HP cannot be negative
			}
			// Temp HP doesn't stack - just set the new value
			m.InitiativeList[originalIndex].TempHP = newTempHP

			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		} else {
			// Handle error
			SetError(&m, err.Error())
			m.InitiativeEditMode = false
			m.InitiativeEditType = ""
			m.InitiativeInput = ""
		}

	case "delete":
		// Delete the selected entry
		m.InitiativeList = append(m.InitiativeList[:originalIndex], m.InitiativeList[originalIndex+1:]...)
		// Adjust selected entry if needed
		if m.SelectedEntry >= len(m.InitiativeList) && len(m.InitiativeList) > 0 {
			m.SelectedEntry = len(m.InitiativeList) - 1
		} else if len(m.InitiativeList) == 0 {
			m.SelectedEntry = -1
			m.InitiativeListMode = false
		}
		m.InitiativeEditMode = false
		m.InitiativeEditType = ""

		// Update display names after deletion (keeps instance numbers stable)
		m = renumberMonsterInstances(m)
	}

	return m
}

// findOriginalIndex finds the original array index for a sorted display position
func findOriginalIndex(m Model, sortedIndex int) int {
	if sortedIndex < 0 || sortedIndex >= len(m.InitiativeList) {
		return -1
	}

	// Create a copy of the list with original indices
	type indexedEntry struct {
		entry         InitiativeEntry
		originalIndex int
	}

	var indexed []indexedEntry
	for i, entry := range m.InitiativeList {
		indexed = append(indexed, indexedEntry{entry: entry, originalIndex: i})
	}

	// Sort by initiative (highest first) - same logic as display
	for i := 0; i < len(indexed); i++ {
		for j := i + 1; j < len(indexed); j++ {
			if indexed[j].entry.Initiative > indexed[i].entry.Initiative {
				indexed[i], indexed[j] = indexed[j], indexed[i]
			}
		}
	}

	// Return the original index for the sorted position
	if sortedIndex < len(indexed) {
		return indexed[sortedIndex].originalIndex
	}

	return -1
}

// renumberMonsterInstances assigns instance numbers to monsters, keeping existing numbers stable
func renumberMonsterInstances(m Model) Model {
	// First pass: normalize all base names
	for i := range m.InitiativeList {
		if m.InitiativeList[i].BaseName == "" {
			// If BaseName is empty, use current Name as BaseName
			m.InitiativeList[i].BaseName = m.InitiativeList[i].Name
		}
	}

	// Count instances of each base name and find max instance number
	nameCounts := make(map[string]int)
	maxInstanceNum := make(map[string]int)

	for _, entry := range m.InitiativeList {
		nameCounts[entry.BaseName]++
		if entry.InstanceNum > maxInstanceNum[entry.BaseName] {
			maxInstanceNum[entry.BaseName] = entry.InstanceNum
		}
	}

	// Assign instance numbers to entries that don't have one yet
	for i := range m.InitiativeList {
		baseName := m.InitiativeList[i].BaseName

		// If there's more than one instance of this monster
		if nameCounts[baseName] > 1 {
			// If this entry doesn't have an instance number yet, assign the next available one
			if m.InitiativeList[i].InstanceNum == 0 {
				maxInstanceNum[baseName]++
				m.InitiativeList[i].InstanceNum = maxInstanceNum[baseName]
			}
			// Update display name with instance number
			m.InitiativeList[i].Name = fmt.Sprintf("%s %d", baseName, m.InitiativeList[i].InstanceNum)
		} else {
			// Only one instance, no number needed
			m.InitiativeList[i].InstanceNum = 0
			m.InitiativeList[i].Name = baseName
		}
	}

	return m
}

// ========== SEARCH INPUT HELPERS ==========

// addToSpellSearch adds a character to spell search input and updates suggestions
func addToSpellSearch(m Model, char string) Model {
	m.SpellSearchInput += char
	m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput, "")
	if len(m.SpellSuggestions) > 0 {
		m.SuggestionIndex = 0
	} else {
		m.SuggestionIndex = -1
	}
	return m
}

// addToMonsterSearch adds a character to monster search input and updates suggestions
func addToMonsterSearch(m Model, char string) Model {
	m.MonsterSearchInput += char
	m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput, m.MonsterCRFilter)
	if len(m.MonsterSuggestions) > 0 {
		m.MonsterSuggestionIndex = 0
	} else {
		m.MonsterSuggestionIndex = -1
	}
	return m
}

// addToInitiativeInput adds a character to initiative input
func addToInitiativeInput(m Model, char string) Model {
	m.InitiativeInput += char
	return m
}

// handleSearchModeInput handles input when in search mode for spells or monsters
func handleSearchModeInput(m Model, char string) Model {
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += char
		return m
	} else if m.SpellSearchMode && m.ActivePanel == Spells {
		return addToSpellSearch(m, char)
	} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
		return addToMonsterSearch(m, char)
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		return addToInitiativeInput(m, char)
	} else if m.NotesEditMode && m.ActivePanel == Notes {
		m.NotesInput += char
		return m
	} else if m.NotesSearchMode && m.ActivePanel == Notes {
		m.NotesSearchInput += char
		return m
	}
	return m
}

// ========== UTILITY HELPERS ==========

// getMonsterFieldString gets a string field value using reflection (helper for monster operations)
func getMonsterFieldString(v reflect.Value, fieldName string) string {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}
	if field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

// convertMonsterActions converts panels.MonsterAction to MonsterAction
func convertMonsterActions(panelActions []panels.MonsterAction) []MonsterAction {
	uiActions := make([]MonsterAction, len(panelActions))
	for i, action := range panelActions {
		uiActions[i] = MonsterAction{
			Name:        action.Name,
			Type:        action.Type,
			Description: action.Description,
			Roll:        action.Roll,
			Reach:       action.Reach,
			Range:       action.Range,
			Damage:      action.Damage,
			DamageType:  action.DamageType,
			SaveDC:      action.SaveDC,
			SaveType:    action.SaveType,
		}
	}
	return uiActions
}

// cleanDiceNotation cleans up dice notation string for parsing
// Removes spaces around operators but keeps comma-separated format
// Example: "2d6 + 4, 7d6" -> "2d6+4, 7d6"
func cleanDiceNotation(damage string) string {
	result := damage

	// Remove spaces around operators (+, -) but NOT around commas
	result = strings.ReplaceAll(result, " + ", "+")
	result = strings.ReplaceAll(result, " - ", "-")
	result = strings.ReplaceAll(result, " +", "+")
	result = strings.ReplaceAll(result, "+ ", "+")
	result = strings.ReplaceAll(result, " -", "-")
	result = strings.ReplaceAll(result, "- ", "-")

	return strings.TrimSpace(result)
}

// ========== UNDO/REDO HP HISTORY ==========

// pushHPHistory saves an HP change to the undo stack (limited to 3 actions)
func pushHPHistory(m *Model, entryIndex int, oldHP int, newHP int) {
	if entryIndex < 0 || entryIndex >= len(m.InitiativeList) {
		return
	}

	entry := HPHistoryEntry{
		EntryIndex: entryIndex,
		OldHP:      oldHP,
		NewHP:      newHP,
		MaxHP:      m.InitiativeList[entryIndex].MaxHP,
		EntryName:  m.InitiativeList[entryIndex].Name,
	}

	// Add to undo stack
	m.HPUndoStack = append(m.HPUndoStack, entry)

	// Limit to 3 actions
	if len(m.HPUndoStack) > 3 {
		m.HPUndoStack = m.HPUndoStack[1:]
	}

	// Clear redo stack when new action is performed
	m.HPRedoStack = []HPHistoryEntry{}
}

// undoHPChange reverts the last HP change
func undoHPChange(m Model) Model {
	if len(m.HPUndoStack) == 0 {
		return m
	}

	// Pop from undo stack
	lastAction := m.HPUndoStack[len(m.HPUndoStack)-1]
	m.HPUndoStack = m.HPUndoStack[:len(m.HPUndoStack)-1]

	// Apply old HP and temp HP values
	if lastAction.EntryIndex >= 0 && lastAction.EntryIndex < len(m.InitiativeList) {
		// Validate values before applying
		m.InitiativeList[lastAction.EntryIndex].HP = HPCalc.ValidateHP(lastAction.OldHP, m.InitiativeList[lastAction.EntryIndex].MaxHP)
		m.InitiativeList[lastAction.EntryIndex].TempHP = HPCalc.SetTempHP(lastAction.OldTempHP)

		// Add to redo stack
		m.HPRedoStack = append(m.HPRedoStack, lastAction)
		if len(m.HPRedoStack) > 3 {
			m.HPRedoStack = m.HPRedoStack[1:]
		}
	}

	return m
}

// redoHPChange reapplies a previously undone HP change
func redoHPChange(m Model) Model {
	if len(m.HPRedoStack) == 0 {
		return m
	}

	// Pop from redo stack
	lastUndo := m.HPRedoStack[len(m.HPRedoStack)-1]
	m.HPRedoStack = m.HPRedoStack[:len(m.HPRedoStack)-1]

	// Apply new HP and temp HP values
	if lastUndo.EntryIndex >= 0 && lastUndo.EntryIndex < len(m.InitiativeList) {
		// Validate values before applying
		m.InitiativeList[lastUndo.EntryIndex].HP = HPCalc.ValidateHP(lastUndo.NewHP, m.InitiativeList[lastUndo.EntryIndex].MaxHP)
		m.InitiativeList[lastUndo.EntryIndex].TempHP = HPCalc.SetTempHP(lastUndo.NewTempHP)

		// Add back to undo stack
		m.HPUndoStack = append(m.HPUndoStack, lastUndo)
		if len(m.HPUndoStack) > 3 {
			m.HPUndoStack = m.HPUndoStack[1:]
		}
	}

	return m
}
