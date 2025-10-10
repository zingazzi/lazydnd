// ui/handlers/actions.go
package ui

import (
	"lazydnd/panels"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

// ========== LETTER HANDLERS ==========

// handleR handles the 'r' key
func handleR(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "r"), nil
	}

	// Reroll dice command
	if m.ActivePanel == DiceRoller && !m.InputMode && !m.DiceHistoryMode && m.LastDiceCommand != "" {
		result := panels.RollDice(m.LastDiceCommand, m.Config)
		m.DiceResult = result
		m.addToHistory(result, m.LastDiceCommand)
	}

	return m, nil
}

// handleP handles the 'p' key
func handleP(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "p"), nil
	}

	// Start adding a player to initiative
	if m.ActivePanel == InitiativeTracker {
		m.InitiativeInputMode = true
		m.InitiativeInputType = "player_name"
		m.InitiativeInput = ""
	}

	return m, nil
}

// handleM handles the 'm' key
func handleM(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "m"), nil
	}

	// Start adding a monster to initiative
	if m.ActivePanel == InitiativeTracker {
		m.InitiativeInputMode = true
		m.InitiativeInputType = "monster_name"
		m.InitiativeInput = ""
	}

	return m, nil
}

// handleE handles the 'e' key
func handleE(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "e"), nil
	}

	// 'e' key no longer enters edit mode (now using Enter key)
	// Keep the handler for typing 'e' in search modes

	return m, nil
}

// handleS handles the 's' key
func handleS(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input first
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "s"), nil
	}

	// Open saving throw popup in initiative tracker when a monster is selected
	if m.ActivePanel == InitiativeTracker && !m.InputMode {
		if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			entry := m.InitiativeList[m.SelectedEntry]
			// Only works for monsters with full data
			if entry.Type == "monster" && entry.MonsterData != nil {
				m.ShowSavingThrowPopup = true
			}
		}
	}

	return m, nil
}

// handleI handles the 'i' key
func handleI(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "i"
		return m, nil
	}

	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "i"), nil
	}

	// Edit initiative value in list mode
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		m.InitiativeEditMode = true
		m.InitiativeEditType = "initiative"
		m.InitiativeInput = ""
	}

	return m, nil
}

// handleH handles the 'h' key
func handleH(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "h"
		return m, nil
	}

	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "h"), nil
	}

	// Enter dice history mode in dice roller
	if m.ActivePanel == DiceRoller && !m.InputMode && !m.DiceHistoryMode && len(m.DiceHistory) > 0 {
		m.DiceHistoryMode = true
		// Start at the most recent (last index)
		m.HistoryIndex = len(m.DiceHistory) - 1
		return m, nil
	}

	// Edit HP in list mode (only for monsters)
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 && m.InitiativeList[originalIndex].Type == "monster" {
			m.InitiativeEditMode = true
			m.InitiativeEditType = "hp"
			m.InitiativeInput = ""
		}
	}

	return m, nil
}

// handleA handles the 'a' key
func handleA(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "a"
		return m, nil
	}

	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "a"), nil
	}

	// Add monster to initiative from Monster panel
	if m.ActivePanel == Monsters && m.SelectedMonster != nil {
		// Add selected monster to initiative tracker
		monsterName := getMonsterFieldString(reflect.ValueOf(m.SelectedMonster).Elem(), "Name")
		if monsterName != "" {
			// Extract monster stats
			hp, ac, err := panels.ExtractMonsterStats(monsterName)
			if err == nil && hp > 0 && ac > 0 {
				// Roll initiative for the monster
				initiative := panels.RollInitiative()

				// Create initiative entry with full monster data link
				newEntry := InitiativeEntry{
					Name:        monsterName,
					Type:        "monster",
					Initiative:  initiative,
					HP:          hp,
					MaxHP:       hp,
					AC:          ac,
					MonsterData: m.SelectedMonster, // Link to full monster data
					BaseName:    monsterName,
					MonsterName: monsterName, // Store for save/load persistence
				}

				// Add to initiative list
				m.InitiativeList = append(m.InitiativeList, newEntry)

				// Renumber instances if there are duplicates
				m = renumberMonsterInstances(m)
			}
		}
	} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && !m.InitiativeInputMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		// Show action popup for selected monster (if it has actions)
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 && m.InitiativeList[originalIndex].MonsterData != nil {
			monster := m.InitiativeList[originalIndex].MonsterData
			if len(monster.ActionList) > 0 {
				m.ShowActionPopup = true
				m.ActionPopupActions = monster.ActionList
				m.ActionPopupIndex = 0
				m.ActionPopupMonster = monster.Name
			}
		}
	}

	return m, nil
}

// handleD handles the 'd' key
func handleD(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "d"
		return m, nil
	}

	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "d"), nil
	}

	// Delete entry in list mode
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		m.InitiativeEditMode = true
		m.InitiativeEditType = "delete"
	}

	return m, nil
}

// handleL handles the 'l' key
func handleL(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "l"), nil
	}

	// Show linked monster details in Monster panel
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 && m.InitiativeList[originalIndex].MonsterData != nil {
			m.SelectedMonster = m.InitiativeList[originalIndex].MonsterData
			m.ActivePanel = Monsters
		}
	}

	return m, nil
}

// handleC handles the 'c' key
func handleC(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle search mode input
	if m.SpellSearchMode || m.MonsterSearchMode || m.InitiativeInputMode || m.InitiativeEditMode {
		return handleSearchModeInput(m, "c"), nil
	}

	// Duplicate entry in list mode
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		// Duplicate the selected entry
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 {
			original := m.InitiativeList[originalIndex]

			// Use the BaseName from the original (not the numbered Name)
			baseName := original.BaseName
			if baseName == "" {
				baseName = original.Name
			}

			// Create a duplicate with same stats
			duplicate := InitiativeEntry{
				Name:        baseName, // Use base name, will be renumbered
				Type:        original.Type,
				Initiative:  original.Initiative,
				HP:          original.HP,
				MaxHP:       original.MaxHP,
				AC:          original.AC,
				MonsterData: original.MonsterData,
				BaseName:    baseName,
			}

			// Add to initiative list
			m.InitiativeList = append(m.InitiativeList, duplicate)

			// Renumber all instances of this monster
			m = renumberMonsterInstances(m)
		}
	}

	return m, nil
}
