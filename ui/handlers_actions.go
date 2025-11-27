// ui/handlers/actions.go
package ui

import (
	"fmt"
	"lazydnd/panels"
	"reflect"
	"strings"
)

// ========== LETTER HANDLERS ==========

// handleR handles the 'r' key
func handleR(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input first
	if m.IsInputMode() {
		return handleSearchModeInput(m, "r"), nil
	}

	// Check if this is Shift+R (capital R) for reaction toggle
	isShiftR := msg.String() == KeyShiftR

	if isShiftR {
		// Toggle reaction status in initiative tracker
		if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			m.InitiativeList[m.SelectedEntry].ReactionUsed = !m.InitiativeList[m.SelectedEntry].ReactionUsed
		}
		return m, nil
	}

	// Reroll dice command (lowercase 'r')
	if m.ActivePanel == DiceRoller && !m.InputMode && !m.DiceHistoryMode && m.LastDiceCommand != "" {
		m = rerollDiceCommand(m)
	}

	return m, nil
}

// rerollDiceCommand rerolls the last dice command with critical hit detection for monster attacks
func rerollDiceCommand(m Model) Model {
	command := m.LastDiceCommand

	// Check if this is a monster attack format: "1d20+X, damage"
	if strings.Contains(command, ",") && strings.HasPrefix(command, "1d20") {
		// Split into attack and damage
		parts := strings.SplitN(command, ",", 2)
		if len(parts) == 2 {
			attackCommand := strings.TrimSpace(parts[0])
			damageCommand := strings.TrimSpace(parts[1])

			// Roll attack first
			attackResult := panels.RollDice(attackCommand, m.Config)

			// Check for critical hit
			isCrit := strings.Contains(attackResult, "CRIT") || (strings.Contains(attackResult, "d20: 20") && m.Config.DiceRoller.CriticalHitEnabled)

			// Roll damage (with crit if needed)
			var damageResult string
			if isCrit {
				damageResult = panels.RollDice(damageCommand+" crit", m.Config)
			} else {
				damageResult = panels.RollDice(damageCommand, m.Config)
			}

			// Combine results
			result := attackResult + "\n\n" + damageResult
			m.DiceResult = result
			m.addToHistory(result, command)
			return m
		}
	}

	// Normal reroll for non-attack commands
	result := panels.RollDice(command, m.Config)
	m.DiceResult = result
	m.addToHistory(result, command)
	return m
}

// handleP handles the 'p' key
func handleP(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input first
	if m.IsInputMode() {
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
func handleM(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input first
	if m.IsInputMode() {
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
func handleE(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input first
	if m.IsInputMode() {
		return handleSearchModeInput(m, "e"), nil
	}

	// Handle Notes panel edit mode
	if m.ActivePanel == Notes && !m.NotesEditMode && !m.NotesSearchMode {
		return handleNotesE(m, msg)
	}

	// 'e' key no longer enters edit mode (now using Enter key)
	// Keep the handler for typing 'e' in search modes

	return m, nil
}

// handleS handles the 's' key
func handleS(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input first
	if m.IsInputMode() {
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
func handleI(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "i"
		return m, nil
	}

	// Handle search mode input
	if m.IsInputMode() {
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

// handleH handles the 'h' and 'H' keys
func handleH(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "h"
		return m, nil
	}

	// Handle Initiative Tracker name input mode (typing player/monster names)
	if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
		m.InitiativeInput += msg.String() // Add "h" or "H" to the name
		return m, nil
	}

	// Handle search mode input
	if m.IsInputMode() {
		return handleSearchModeInput(m, "h"), nil
	}

	// Enter dice history mode in dice roller
	if m.ActivePanel == DiceRoller && !m.InputMode && !m.DiceHistoryMode && len(m.DiceHistory) > 0 {
		m.DiceHistoryMode = true
		// Start at the most recent (last index)
		m.HistoryIndex = len(m.DiceHistory) - 1
		return m, nil
	}

	// Check if this is Shift+H (capital H) for max HP editing
	isShiftH := msg.String() == KeyShiftH

	// Edit HP or Max HP in list mode (for both monsters and players)
	// This works when in list mode, whether or not already in edit mode
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 {
			m.InitiativeEditMode = true
			if isShiftH {
				m.InitiativeEditType = "maxhp"
			} else {
				m.InitiativeEditType = "hp"
			}
			m.InitiativeInput = ""
		}
		return m, nil
	}

	// When NOT in list mode, 'h' opens quick HP popup (damage/remove by default)
	// This allows quick HP adjustments without entering edit mode
	if m.ActivePanel == InitiativeTracker && !m.InitiativeListMode && !m.InitiativeEditMode && !m.InitiativeInputMode {
		return handleQuickRemoveHP(m, msg)
	}

	return m, nil
}

// handleK handles the 'k' key (AC editing)
func handleK(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "k"
		return m, nil
	}

	// Handle Initiative Tracker name input mode (typing player/monster names)
	if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
		m.InitiativeInput += msg.String() // Add "k" or "K" to the name
		return m, nil
	}

	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "k"), nil
	}

	// Edit AC in list mode (for both monsters and players)
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		originalIndex := findOriginalIndex(m, m.SelectedEntry)
		if originalIndex >= 0 {
			m.InitiativeEditMode = true
			m.InitiativeEditType = "ac"
			m.InitiativeInput = ""
		}
	}

	return m, nil
}

// handleA handles the 'a' key
func handleA(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "a"
		return m, nil
	}

	// Add monster to initiative from Monster panel
	// Check this BEFORE search mode so 'a' works even if search mode is still active
	if m.ActivePanel == Monsters && m.SelectedMonster != nil {
		// Add selected monster to initiative tracker
		monsterName := getMonsterFieldString(reflect.ValueOf(m.SelectedMonster).Elem(), "Name")
		if monsterName != "" {
			// Extract monster stats
			hp, ac, err := panels.ExtractMonsterStats(monsterName)
			if err == nil && hp > 0 && ac > 0 {
				// Roll initiative for the monster
				initiative := panels.RollInitiative()

			// Parse legendary actions count
			legendaryMax := 0
			if m.SelectedMonster != nil && m.SelectedMonster.LegendaryActions != "" {
				legendaryMax = panels.ParseLegendaryActionCount(m.SelectedMonster.LegendaryActions)
			}

			// Create initiative entry with full monster data link
			newEntry := InitiativeEntry{
				Name:                monsterName,
				Type:                "monster",
				Initiative:          initiative,
				HP:                  hp,
				MaxHP:               hp,
				AC:                  ac,
				ReactionUsed:        false, // Initialize reaction as available
				MonsterData:         m.SelectedMonster, // Link to full monster data
				BaseName:            monsterName,
				MonsterName:         monsterName, // Store for save/load persistence
				LegendaryActionsMax: legendaryMax,
				LegendaryActionsUsed: 0, // Start with all legendary actions available
			}

				// Add to initiative list
				m.InitiativeList = append(m.InitiativeList, newEntry)

				// Renumber instances if there are duplicates
				m = renumberMonsterInstances(m)

				// Sort initiative list by initiative value (descending)
				sortInitiativeList(&m.InitiativeList)

				// Set up initiative tracker state so 'h' key works
				if len(m.InitiativeList) > 0 {
					m.InitiativeListMode = true
					// Set selected entry to 0 (first entry after sorting)
					// This ensures the tracker is in the correct state for editing
					m.SelectedEntry = 0
					// Switch to initiative tracker panel so 'h' key works immediately
					m.ActivePanel = InitiativeTracker
				}
			}
		}
		return m, nil
	}

	// Handle search mode input (only if not adding monster to initiative)
	if m.isInInputMode() {
		return handleSearchModeInput(m, "a"), nil
	}

	// Show action popup in initiative tracker
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && !m.InitiativeInputMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
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
func handleD(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle dice input mode
	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += "d"
		return m, nil
	}

	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "d"), nil
	}

	// Delete active spell in active spell list mode
	if m.ActivePanel == Spells && m.ActiveSpellListMode {
		return handleDeleteActiveSpell(m)
	}

	// Delete entry in list mode
	if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		m.InitiativeEditMode = true
		m.InitiativeEditType = "delete"
	}

	return m, nil
}

// handleL handles the 'l' key
func handleL(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "l"), nil
	}

	// Only work in Initiative Tracker
	if m.ActivePanel != InitiativeTracker || !m.InitiativeListMode || m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
		return m, nil
	}

	// Check if this is Shift+L (capital L) for restoring legendary action
	isShiftL := msg.String() == KeyShiftL
	originalIndex := findOriginalIndex(m, m.SelectedEntry)
	if originalIndex < 0 {
		return m, nil
	}

	entry := &m.InitiativeList[originalIndex]

	// Handle legendary actions for monsters
	if entry.Type == "monster" && entry.LegendaryActionsMax > 0 {
		if isShiftL {
			// Shift+L: Restore one legendary action
			if entry.LegendaryActionsUsed > 0 {
				entry.LegendaryActionsUsed--
				SetSuccess(&m, fmt.Sprintf("Restored legendary action for %s (%d/%d)", entry.Name, entry.LegendaryActionsMax-entry.LegendaryActionsUsed, entry.LegendaryActionsMax))
			}
		} else {
			// l: Use one legendary action
			if entry.LegendaryActionsUsed < entry.LegendaryActionsMax {
				entry.LegendaryActionsUsed++
				SetSuccess(&m, fmt.Sprintf("Used legendary action for %s (%d/%d)", entry.Name, entry.LegendaryActionsMax-entry.LegendaryActionsUsed, entry.LegendaryActionsMax))
			} else {
				SetError(&m, fmt.Sprintf("%s has no legendary actions remaining", entry.Name))
			}
		}
		return m, nil
	}

	// For monsters without legendary actions, or if Shift+L pressed on non-monster, show linked monster details
	if !isShiftL && entry.MonsterData != nil {
		m.SelectedMonster = entry.MonsterData
		m.ActivePanel = Monsters
	}

	return m, nil
}

// handleC handles the 'c' key
func handleC(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "c"), nil
	}

	// Cast spell in Spells panel
	if m.ActivePanel == Spells && !m.ActiveSpellListMode {
		return handleCastSpell(m)
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
				Name:                baseName, // Use base name, will be renumbered
				Type:                original.Type,
				Initiative:          original.Initiative,
				HP:                  original.HP,
				MaxHP:               original.MaxHP,
				TempHP:              original.TempHP,
				AC:                  original.AC,
				ReactionUsed:        false, // Reset reaction for new duplicate
				MonsterData:         original.MonsterData,
				BaseName:            baseName,
				MonsterName:         original.MonsterName,
				LegendaryActionsMax: original.LegendaryActionsMax,
				LegendaryActionsUsed: 0, // Reset legendary actions for new duplicate
			}

			// Add to initiative list
			m.InitiativeList = append(m.InitiativeList, duplicate)

			// Renumber all instances of this monster
			m = renumberMonsterInstances(m)

			// Keep selection valid after renumbering
			if m.SelectedEntry >= len(m.InitiativeList) {
				m.SelectedEntry = len(m.InitiativeList) - 1
			}
			if m.SelectedEntry < 0 && len(m.InitiativeList) > 0 {
				m.SelectedEntry = 0
			}
		}
	}

	return m, nil
}

// handleV handles the 'v' key
func handleV(m Model, msg KeyMsg) (Model, Cmd) {
	// Handle search mode input
	if m.isInInputMode() {
		return handleSearchModeInput(m, "v"), nil
	}

	// View active spells in Spells panel
	if m.ActivePanel == Spells {
		return handleViewActiveSpells(m)
	}

	return m, nil
}
