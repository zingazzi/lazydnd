// ui/navigation.go
package ui

import (
	"lazydnd/panels"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleNavigation processes navigation-related key presses
func (m Model) HandleNavigation(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if !m.InputMode && !m.InitiativeInputMode && !m.SpellSearchMode && !m.MonsterSearchMode {
			return m, tea.Quit
		} else if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'q' to spell search input
			m.SpellSearchInput += "q"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'q' to monster search input
			m.MonsterSearchInput += "q"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		}

	case "1":
		if m.InputMode {
			m.DiceInput += "1"
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			m.InitiativeInput += "1"
		} else {
			m.ActivePanel = DiceRoller
		}
	case "2":
		if m.InputMode {
			m.DiceInput += "2"
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			m.InitiativeInput += "2"
		} else {
			m.ActivePanel = InitiativeTracker
		}
	case "3":
		if m.InputMode {
			m.DiceInput += "3"
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			m.InitiativeInput += "3"
		} else {
			m.ActivePanel = Spells
		}
	case "4":
		if m.InputMode {
			m.DiceInput += "4"
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			m.InitiativeInput += "4"
		} else {
			m.ActivePanel = Monsters
		}

	case "tab":
		if !m.InputMode {
			m.ActivePanel = (m.ActivePanel + 1) % 4
		}

	case "f1":
		m.ActivePanel = DiceRoller
		m.InputMode = false
	case "f2":
		m.ActivePanel = InitiativeTracker
		m.InputMode = false
	case "f3":
		m.ActivePanel = Spells
		m.InputMode = false
	case "f4":
		m.ActivePanel = Monsters
		m.InputMode = false

	case "up":
		// Handle spell suggestion navigation first (highest priority)
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

	case "down":
		// Handle spell suggestion navigation first (highest priority)
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

	case "enter":
		if m.ActivePanel == DiceRoller {
			if m.InputMode && m.DiceInput != "" {
				// Roll the dice
				result := panels.RollDice(m.DiceInput)
				m.DiceResult = result
				m.DiceHistory = append(m.DiceHistory, result)
				m.LastDiceCommand = m.DiceInput
				if len(m.DiceHistory) > 15 {
					m.DiceHistory = m.DiceHistory[1:]
				}
				m.DiceInput = ""
				m.InputMode = false
			} else {
				m.InputMode = true
			}
		} else if m.ActivePanel == InitiativeTracker {
			if m.InitiativeEditMode {
				// Process edit action
				m = m.processInitiativeEdit()
			} else if m.InitiativeInputMode {
				// Process initiative tracker input
				m = m.processInitiativeInput()
			}
		} else if m.ActivePanel == Spells {
			if m.SpellSearchMode {
				// Select spell from suggestions
				if len(m.SpellSuggestions) > 0 && m.SuggestionIndex >= 0 && m.SuggestionIndex < len(m.SpellSuggestions) {
					selectedSpellName := m.SpellSuggestions[m.SuggestionIndex]
					foundSpell := panels.FindSpell(selectedSpellName)
					if foundSpell != nil {
						// Convert panels.Spell to ui.Spell
						m.SelectedSpell = &Spell{
							Name:            foundSpell.Name,
							Level:           foundSpell.Level,
							School:          foundSpell.School,
							Classes:         foundSpell.Classes,
							ActionType:      foundSpell.ActionType,
							Concentration:   foundSpell.Concentration,
							Ritual:          foundSpell.Ritual,
							Range:           foundSpell.Range,
							Components:      foundSpell.Components,
							Material:        foundSpell.Material,
							Duration:        foundSpell.Duration,
							Description:     foundSpell.Description,
							CantripUpgrade:  foundSpell.CantripUpgrade,
						}
					}
					m.SpellSearchInput = selectedSpellName
					m.SpellSearchMode = false
					m.SpellSuggestions = []string{}
					m.SuggestionIndex = -1
				}
			} else {
				// Start spell search mode
				m.SpellSearchMode = true
				m.SpellSuggestions = []string{}
				m.SuggestionIndex = -1
			}
		} else if m.ActivePanel == Monsters {
			if m.MonsterSearchMode {
				// Select monster from suggestions
				if len(m.MonsterSuggestions) > 0 && m.MonsterSuggestionIndex >= 0 && m.MonsterSuggestionIndex < len(m.MonsterSuggestions) {
					selectedMonsterName := m.MonsterSuggestions[m.MonsterSuggestionIndex]
					foundMonster := panels.FindMonster(selectedMonsterName)
					if foundMonster != nil {
						// Convert panels.Monster to ui.Monster
						m.SelectedMonster = &Monster{
							Name:             foundMonster.Name,
							Meta:             foundMonster.Meta,
							ArmorClass:       foundMonster.ArmorClass,
							HitPoints:        foundMonster.HitPoints,
							Speed:            foundMonster.Speed,
							STR:              foundMonster.STR,
							STRMod:           foundMonster.STRMod,
							DEX:              foundMonster.DEX,
							DEXMod:           foundMonster.DEXMod,
							CON:              foundMonster.CON,
							CONMod:           foundMonster.CONMod,
							INT:              foundMonster.INT,
							INTMod:           foundMonster.INTMod,
							WIS:              foundMonster.WIS,
							WISMod:           foundMonster.WISMod,
							CHA:              foundMonster.CHA,
							CHAMod:           foundMonster.CHAMod,
							SavingThrows:     foundMonster.SavingThrows,
							Skills:           foundMonster.Skills,
							Senses:           foundMonster.Senses,
							Languages:        foundMonster.Languages,
							Challenge:        foundMonster.Challenge,
							Traits:           foundMonster.Traits,
							Actions:          foundMonster.Actions,
							LegendaryActions: foundMonster.LegendaryActions,
							ImgURL:           foundMonster.ImgURL,
						}
					}
					m.MonsterSearchInput = selectedMonsterName
					m.MonsterSearchMode = false
					m.MonsterSuggestions = []string{}
					m.MonsterSuggestionIndex = -1
				}
			} else {
				// Start monster search mode
				m.MonsterSearchMode = true
				m.MonsterSuggestions = []string{}
				m.MonsterSuggestionIndex = -1
			}
		}

	case "r":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'r' to spell search input
			m.SpellSearchInput += "r"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'r' to monster search input
			m.MonsterSearchInput += "r"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == DiceRoller && !m.InputMode && !m.InitiativeInputMode && m.LastDiceCommand != "" {
			// Reroll the last dice command
			result := panels.RollDice(m.LastDiceCommand)
			m.DiceResult = result
			m.DiceHistory = append(m.DiceHistory, result)
			if len(m.DiceHistory) > 15 {
				m.DiceHistory = m.DiceHistory[1:]
			}
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Add 'r' to input when in input mode
			m.InitiativeInput += "r"
		}

	case "esc":
		if m.ActivePanel == DiceRoller {
			m.DiceInput = ""
			m.InputMode = false
		} else if m.ActivePanel == InitiativeTracker {
			if m.InitiativeEditMode {
				// Exit edit mode
				m.InitiativeEditMode = false
				m.InitiativeEditType = ""
				m.InitiativeInput = ""
			} else if m.InitiativeListMode {
				// Exit list mode
				m.InitiativeListMode = false
				m.SelectedEntry = -1
			} else {
				// Exit input mode
				m.InitiativeInput = ""
				m.InitiativeInputMode = false
				m.InitiativeInputType = ""
			}
		} else if m.ActivePanel == Spells {
			m.SpellSearchInput = ""
			m.SpellSearchMode = false
			m.SpellSuggestions = []string{}
			m.SuggestionIndex = -1
		} else if m.ActivePanel == Monsters {
			m.MonsterSearchInput = ""
			m.MonsterSearchMode = false
			m.MonsterSuggestions = []string{}
			m.MonsterSuggestionIndex = -1
		}

	case "backspace", "ctrl+h":
		if m.InputMode && len(m.DiceInput) > 0 {
			m.DiceInput = m.DiceInput[:len(m.DiceInput)-1]
		} else if m.InitiativeInputMode && len(m.InitiativeInput) > 0 {
			m.InitiativeInput = m.InitiativeInput[:len(m.InitiativeInput)-1]
		} else if m.SpellSearchMode && len(m.SpellSearchInput) > 0 {
			m.SpellSearchInput = m.SpellSearchInput[:len(m.SpellSearchInput)-1]
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && len(m.MonsterSearchInput) > 0 {
			m.MonsterSearchInput = m.MonsterSearchInput[:len(m.MonsterSearchInput)-1]
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		}

	case "p":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'p' to spell search input
			m.SpellSearchInput += "p"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'p' to monster search input
			m.MonsterSearchInput += "p"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && !m.InitiativeInputMode {
			// Start adding a player
			m.InitiativeInputMode = true
			m.InitiativeInputType = "player_name"
			m.InitiativeInput = ""
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Add 'p' to input when in input mode
			m.InitiativeInput += "p"
		}

	case "m":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'm' to spell search input
			m.SpellSearchInput += "m"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'm' to monster search input
			m.MonsterSearchInput += "m"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && !m.InitiativeInputMode {
			// Start adding a monster
			m.InitiativeInputMode = true
			m.InitiativeInputType = "monster_name"
			m.InitiativeInput = ""
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Add 'm' to input when in input mode
			m.InitiativeInput += "m"
		}

	case "e":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'e' to spell search input
			m.SpellSearchInput += "e"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'e' to monster search input
			m.MonsterSearchInput += "e"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && !m.InitiativeInputMode && !m.InitiativeEditMode && len(m.InitiativeList) > 0 {
			// Enter list edit mode
			m.InitiativeListMode = true
			if m.SelectedEntry == -1 {
				m.SelectedEntry = 0
			}
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			// Add 'e' to input when in input/edit mode
			m.InitiativeInput += "e"
		}

	case "i":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'i' to spell search input
			m.SpellSearchInput += "i"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'i' to monster search input
			m.MonsterSearchInput += "i"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && !m.InitiativeInputMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			// Edit initiative
			m.InitiativeEditMode = true
			m.InitiativeEditType = "initiative"
			m.InitiativeInput = ""
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			// Add 'i' to input when in input/edit mode
			m.InitiativeInput += "i"
		}

	case "h":
		if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'h' to spell search input
			m.SpellSearchInput += "h"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'h' to monster search input
			m.MonsterSearchInput += "h"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && !m.InitiativeInputMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			// Edit HP (only for monsters)
			originalIndex := m.findOriginalIndex(m.SelectedEntry)
			if originalIndex >= 0 && m.InitiativeList[originalIndex].Type == "monster" {
				m.InitiativeEditMode = true
				m.InitiativeEditType = "hp"
				m.InitiativeInput = ""
			}
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			// Add 'h' to input when in input/edit mode
			m.InitiativeInput += "h"
		}

	case "d":
		if m.InputMode && m.ActivePanel == DiceRoller {
			// Add 'd' to dice input (for dice notation like "2d6")
			m.DiceInput += "d"
		} else if m.SpellSearchMode && m.ActivePanel == Spells {
			// Add 'd' to spell search input
			m.SpellSearchInput += "d"
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Add 'd' to monster search input
			m.MonsterSearchInput += "d"
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		} else if m.ActivePanel == InitiativeTracker && m.InitiativeListMode && !m.InitiativeInputMode && !m.InitiativeEditMode && m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			// Delete entry
			m.InitiativeEditMode = true
			m.InitiativeEditType = "delete"
		} else if m.InitiativeInputMode || m.InitiativeEditMode {
			// Add 'd' to input when in input/edit mode
			m.InitiativeInput += "d"
		}

	case "space":
		if m.InputMode && m.ActivePanel == DiceRoller {
			m.DiceInput += " "
		} else if (m.InitiativeInputMode || m.InitiativeEditMode) && m.ActivePanel == InitiativeTracker {
			m.InitiativeInput += " "
		} else if m.SpellSearchMode && m.ActivePanel == Spells {
			m.SpellSearchInput += " "
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		}

	default:
		// Handle text input for dice commands
		if m.InputMode && m.ActivePanel == DiceRoller {
			key := msg.String()
			// Allow alphanumeric characters and common symbols for dice notation
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				(key >= "0" && key <= "9") ||
				key == "+" || key == "-" || key == "d") {
				m.DiceInput += key
			}
		} else if (m.InitiativeInputMode || m.InitiativeEditMode) && m.ActivePanel == InitiativeTracker {
			// Handle text input for initiative tracker (both input and edit modes)
			key := msg.String()
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				(key >= "0" && key <= "9") ||
				key == " " || key == "'" || key == "-" || key == "." || key == "_" || key == "+") {
				m.InitiativeInput += key
			}
		} else if m.SpellSearchMode && m.ActivePanel == Spells {
			// Handle text input for spell search
			key := msg.String()
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				key == "'" || key == "-" || key == " ") {
				m.SpellSearchInput += key
				// Update suggestions
				m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
				if len(m.SpellSuggestions) > 0 {
					m.SuggestionIndex = 0
				} else {
					m.SuggestionIndex = -1
				}
			}
		} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
			// Handle text input for monster search
			key := msg.String()
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				key == "'" || key == "-" || key == " ") {
				m.MonsterSearchInput += key
				// Update suggestions
				m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput)
				if len(m.MonsterSuggestions) > 0 {
					m.MonsterSuggestionIndex = 0
				} else {
					m.MonsterSuggestionIndex = -1
				}
			}
		}
	}

	return m, nil
}

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
		entry InitiativeEntry
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
