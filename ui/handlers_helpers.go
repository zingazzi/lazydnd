// ui/handlers_helpers.go
package ui

import (
	"fmt"
	"lazydnd/panels"
	"reflect"
	"strings"
)

// ========== INITIATIVE PROCESSING ==========

// processInitiativeInput handles the multi-step process of adding players/monsters
func (m Model) processInitiativeInput() Model {
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
			// Complete player entry
			m.TempEntry.Initiative = val.(int)
			m.InitiativeList = append(m.InitiativeList, m.TempEntry)
			m.InitiativeInputMode = false
			m.InitiativeInputType = ""
			m.InitiativeInput = ""
			m.TempEntry = InitiativeEntry{} // Reset temp entry
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
			m.InitiativeList = append(m.InitiativeList, m.TempEntry)
			m.InitiativeInputMode = false
			m.InitiativeInputType = ""
			m.InitiativeInput = ""
			m.TempEntry = InitiativeEntry{} // Reset temp entry
		}
	}

	return m
}

// processInitiativeEdit handles editing existing initiative entries
func (m Model) processInitiativeEdit() Model {
	if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
		return m
	}

	// Find the original array index for the selected display position
	originalIndex := m.findOriginalIndex(m.SelectedEntry)
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
			// Update HP (can be positive or negative)
			change := val.(int)
			newHP := m.InitiativeList[originalIndex].HP + change
			if newHP < 0 {
				newHP = 0
			}
			if newHP > m.InitiativeList[originalIndex].MaxHP {
				newHP = m.InitiativeList[originalIndex].MaxHP
			}
			m.InitiativeList[originalIndex].HP = newHP
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

		// Renumber instances after deletion
		m = m.renumberMonsterInstances()
	}

	return m
}

// findOriginalIndex finds the original array index for a sorted display position
func (m Model) findOriginalIndex(sortedIndex int) int {
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

// renumberMonsterInstances renumbers all monster instances based on duplicates
func (m Model) renumberMonsterInstances() Model {
	// First pass: normalize all base names (strip existing numbers)
	for i := range m.InitiativeList {
		if m.InitiativeList[i].BaseName == "" {
			// Extract base name by removing " N" suffix if present
			name := m.InitiativeList[i].Name
			// Simple approach: if BaseName is empty, use current Name as BaseName
			m.InitiativeList[i].BaseName = name
		}
	}

	// Count instances of each base name
	nameCounts := make(map[string]int)
	for _, entry := range m.InitiativeList {
		nameCounts[entry.BaseName]++
	}

	// Track current instance number for each base name
	instanceNumbers := make(map[string]int)

	// Update all entries
	for i := range m.InitiativeList {
		baseName := m.InitiativeList[i].BaseName

		// If there's more than one instance, add numbers
		if nameCounts[baseName] > 1 {
			instanceNumbers[baseName]++
			m.InitiativeList[i].InstanceNum = instanceNumbers[baseName]
			m.InitiativeList[i].Name = fmt.Sprintf("%s %d", baseName, instanceNumbers[baseName])
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
func (m Model) addToSpellSearch(char string) Model {
	m.SpellSearchInput += char
	m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
	if len(m.SpellSuggestions) > 0 {
		m.SuggestionIndex = 0
	} else {
		m.SuggestionIndex = -1
	}
	return m
}

// addToMonsterSearch adds a character to monster search input and updates suggestions
func (m Model) addToMonsterSearch(char string) Model {
	m.MonsterSearchInput += char
	m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
	if len(m.MonsterSuggestions) > 0 {
		m.MonsterSuggestionIndex = 0
	} else {
		m.MonsterSuggestionIndex = -1
	}
	return m
}

// addToInitiativeInput adds a character to initiative input
func (m Model) addToInitiativeInput(char string) Model {
	m.InitiativeInput += char
	return m
}

// handleSearchModeInput handles input when in search mode for spells or monsters
func (m Model) handleSearchModeInput(char string) Model {
	if m.SpellSearchMode && m.ActivePanel == Spells {
		return m.addToSpellSearch(char)
	} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
		return m.addToMonsterSearch(char)
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		return m.addToInitiativeInput(char)
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

// convertMonsterActions converts panels.MonsterAction to ui.MonsterAction
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
