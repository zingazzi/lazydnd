// ui/handlers/input.go
package ui

import (
	"lazydnd/panels"

	tea "github.com/charmbracelet/bubbletea"
)

// ========== QUIT HANDLERS ==========

// handleQuit handles quit commands (ctrl+c, q)
func handleQuit(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	if !m.InputMode && !m.InitiativeInputMode && !m.SpellSearchMode && !m.MonsterSearchMode {
		return m, tea.Quit
	} else if m.SpellSearchMode && m.ActivePanel == Spells && key == "q" {
		// Add 'q' to spell search input
		m.SpellSearchInput += "q"
		// Update suggestions
		m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput, "")
		if len(m.SpellSuggestions) > 0 {
			m.SuggestionIndex = 0
		} else {
			m.SuggestionIndex = -1
		}
	} else if m.MonsterSearchMode && m.ActivePanel == Monsters && key == "q" {
		// Add 'q' to monster search input
		m.MonsterSearchInput += "q"
		// Update suggestions
		m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput, m.MonsterCRFilter)
		if len(m.MonsterSuggestions) > 0 {
			m.MonsterSuggestionIndex = 0
		} else {
			m.MonsterSuggestionIndex = -1
		}
	}

	return m, nil
}

// handleEscape handles escape key presses
func handleEscape(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Close saving throw popup if open
	if m.ShowSavingThrowPopup {
		m.ShowSavingThrowPopup = false
		return m, nil
	}

	// Close action popup if open
	if m.ShowActionPopup {
		m.ShowActionPopup = false
		m.ActionPopupActions = []MonsterAction{}
		m.ActionPopupIndex = 0
		m.ActionPopupMonster = ""
		return m, nil
	}

	// Close help popup if open
	if m.ShowHelpPopup {
		m.ShowHelpPopup = false
		return m, nil
	}

	if m.ActivePanel == DiceRoller {
		if m.DiceHistoryMode {
			// Exit history mode
			m.DiceHistoryMode = false
			m.HistoryIndex = -1
		} else {
			// Clear input
			m.DiceInput = ""
			m.InputMode = false
		}
	} else if m.ActivePanel == InitiativeTracker {
		if m.MultiTargetMode {
			// Exit multi-target mode
			m.MultiTargetMode = false
			m.SelectedTargets = make(map[int]bool)
			m.ShowMultiTargetPopup = false
			m.MultiTargetInput = ""
			m.TargetSaveResults = make(map[int]string)
		} else if m.InitiativeEditMode {
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
		if m.ActiveSpellListMode {
			// Exit active spell list mode
			m.ActiveSpellListMode = false
			m.ActiveSpellIndex = -1
		} else {
			m.SpellSearchInput = ""
			m.SpellSearchMode = false
			m.SpellSuggestions = []string{}
			m.SuggestionIndex = -1
		}
	} else if m.ActivePanel == Monsters {
		m.MonsterSearchInput = ""
		m.MonsterSearchMode = false
		m.MonsterSuggestions = []string{}
		m.MonsterSuggestionIndex = -1
	} else if m.ActivePanel == Notes {
		// Handle Notes panel Escape key
		m = handleNotesEscape(m)
	}

	return m, nil
}

// ========== ACTION HANDLERS ==========

// handleEnter handles enter key presses
func handleEnter(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Multi-target mode: open popup to input damage/healing
	if m.MultiTargetMode && m.ActivePanel == InitiativeTracker && !m.ShowMultiTargetPopup {
		return handleMultiTargetApply(m)
	}

	// Handle saving throw popup - reroll on Enter (highest priority)
	if m.ShowSavingThrowPopup {
		// Just keep the popup open - the RenderSavingThrowPopup function
		// will generate new rolls each time it's called (random seed)
		// We don't need to do anything, the next render will show new rolls
		return m, nil
	}

	// Handle action popup selection first
	if m.ShowActionPopup && len(m.ActionPopupActions) > 0 && m.ActionPopupIndex >= 0 && m.ActionPopupIndex < len(m.ActionPopupActions) {
		selectedAction := m.ActionPopupActions[m.ActionPopupIndex]

		// Build dice command from action
		var diceCommand string
		if selectedAction.Roll != "" && selectedAction.Damage != "" {
			// Clean damage string: remove spaces and keep comma-separated format
			cleanDamage := cleanDiceNotation(selectedAction.Damage)
			// Format: "1d20+7, 2d6+4" (comma-separated for dice roller)
			diceCommand = "1d20" + selectedAction.Roll + ", " + cleanDamage
		} else if selectedAction.Damage != "" {
			// Just damage roll
			diceCommand = cleanDiceNotation(selectedAction.Damage)
		} else if selectedAction.Roll != "" {
			// Just attack roll
			diceCommand = "1d20" + selectedAction.Roll
		}

		// Execute the dice roll if we have a command
		if diceCommand != "" {
			result := panels.RollDice(diceCommand, m.Config)
			m.DiceResult = result
			m.LastDiceCommand = diceCommand
			m.addToHistory(result, diceCommand)
		}

		// Close the popup
		m.ShowActionPopup = false
		m.ActionPopupActions = []MonsterAction{}
		m.ActionPopupIndex = 0
		m.ActionPopupMonster = ""

		return m, nil
	}

	if m.ActivePanel == DiceRoller {
		if m.DiceHistoryMode && m.HistoryIndex >= 0 && m.HistoryIndex < len(m.DiceCommands) {
			// Re-roll selected history command
			command := m.DiceCommands[m.HistoryIndex]
			result := panels.RollDice(command, m.Config)
			m.DiceResult = result
			m.LastDiceCommand = command
			m.addToHistory(result, command)
			// Exit history mode
			m.DiceHistoryMode = false
			m.HistoryIndex = -1
		} else if m.InputMode && m.DiceInput != "" {
			// Roll the dice
			result := panels.RollDice(m.DiceInput, m.Config)
			m.DiceResult = result
			m.LastDiceCommand = m.DiceInput
			m.addToHistory(result, m.DiceInput)
			m.DiceInput = ""
			m.InputMode = false
		} else if !m.DiceHistoryMode {
			m.InputMode = true
		}
	} else if m.ActivePanel == InitiativeTracker {
		if m.InitiativeEditMode {
			// Process edit action
			m = processInitiativeEdit(m)
		} else if m.InitiativeInputMode {
			// Process initiative tracker input
			m = processInitiativeInput(m)
		} else if len(m.InitiativeList) > 0 {
			// Enter list edit mode when there are entries
			m.InitiativeListMode = true
			if m.SelectedEntry == -1 {
				m.SelectedEntry = 0
			}
		}
	} else if m.ActivePanel == Spells {
		if m.SpellSearchMode {
			// Select spell from suggestions
			if len(m.SpellSuggestions) > 0 && m.SuggestionIndex >= 0 && m.SuggestionIndex < len(m.SpellSuggestions) {
				selectedSpellName := m.SpellSuggestions[m.SuggestionIndex]
				foundSpell := panels.FindSpell(selectedSpellName)
				if foundSpell != nil {
					// Convert panels.Spell to Spell
					m.SelectedSpell = &Spell{
						Name:           foundSpell.Name,
						Level:          foundSpell.Level,
						School:         foundSpell.School,
						Classes:        foundSpell.Classes,
						ActionType:     foundSpell.ActionType,
						Concentration:  foundSpell.Concentration,
						Ritual:         foundSpell.Ritual,
						Range:          foundSpell.Range,
						Components:     foundSpell.Components,
						Material:       foundSpell.Material,
						Duration:       foundSpell.Duration,
						Description:    foundSpell.Description,
						CantripUpgrade: foundSpell.CantripUpgrade,
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
					// Convert panels.Monster to Monster
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
						ActionNumber:     foundMonster.ActionNumber,
						ActionList:       convertMonsterActions(foundMonster.ActionList),
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
	} else if m.ActivePanel == Notes {
		// Handle Notes panel Enter key
		if m.NotesEditMode {
			m = handleNotesEnter(m)
		}
	}

	return m, nil
}

// handleBackspace handles backspace and ctrl+h key presses
func handleBackspace(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InputMode && len(m.DiceInput) > 0 {
		m.DiceInput = m.DiceInput[:len(m.DiceInput)-1]
	} else if (m.InitiativeInputMode || m.InitiativeEditMode) && len(m.InitiativeInput) > 0 {
		m.InitiativeInput = m.InitiativeInput[:len(m.InitiativeInput)-1]
	} else if m.SpellSearchMode && len(m.SpellSearchInput) > 0 {
		m.SpellSearchInput = m.SpellSearchInput[:len(m.SpellSearchInput)-1]
		// Update suggestions
		m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput, "")
		if len(m.SpellSuggestions) > 0 {
			m.SuggestionIndex = 0
		} else {
			m.SuggestionIndex = -1
		}
	} else if m.MonsterSearchMode && len(m.MonsterSearchInput) > 0 {
		m.MonsterSearchInput = m.MonsterSearchInput[:len(m.MonsterSearchInput)-1]
		// Update suggestions
		m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput, m.MonsterCRFilter)
		if len(m.MonsterSuggestions) > 0 {
			m.MonsterSuggestionIndex = 0
		} else {
			m.MonsterSuggestionIndex = -1
		}
	}

	return m, nil
}

// handleSpace handles space key presses
func handleSpace(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Debug: log when space is pressed
	_ = msg // Keep msg parameter even if unused

	// IMPORTANT: Multi-target mode must return the modified model
	// Multi-target mode: toggle selection of current entry (HIGHEST PRIORITY)
	if m.MultiTargetMode {
		newModel, cmd := handleMultiTargetSpace(m)
		return newModel, cmd
	}

	if m.InputMode && m.ActivePanel == DiceRoller {
		m.DiceInput += " "
	} else if (m.InitiativeInputMode || m.InitiativeEditMode) && m.ActivePanel == InitiativeTracker {
		m.InitiativeInput += " "
	} else if m.SpellSearchMode && m.ActivePanel == Spells {
		m.SpellSearchInput += " "
		// Update suggestions
		m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput, "")
		if len(m.SpellSuggestions) > 0 {
			m.SuggestionIndex = 0
		} else {
			m.SuggestionIndex = -1
		}
	}

	return m, nil
}

// ========== DEFAULT INPUT HANDLER ==========

// handleDefaultInput handles default text input for various modes
func handleDefaultInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	// Handle text input for dice commands
	if m.InputMode && m.ActivePanel == DiceRoller {
		// Allow alphanumeric characters and common symbols for dice notation
		if len(key) == 1 && ((key >= "a" && key <= "z") ||
			(key >= "A" && key <= "Z") ||
			(key >= "0" && key <= "9") ||
			key == "+" || key == "-" || key == "d" || key == " " || key == ",") {
			m.DiceInput += key
		}
	} else if (m.InitiativeInputMode || m.InitiativeEditMode) && m.ActivePanel == InitiativeTracker {
		// Handle text input for initiative tracker (both input and edit modes)
		if len(key) == 1 && ((key >= "a" && key <= "z") ||
			(key >= "A" && key <= "Z") ||
			(key >= "0" && key <= "9") ||
			key == " " || key == "'" || key == "-" || key == "." || key == "_" || key == "+") {
			m.InitiativeInput += key
		}
	} else if m.SpellSearchMode && m.ActivePanel == Spells {
		// Handle text input for spell search
		if len(key) == 1 && ((key >= "a" && key <= "z") ||
			(key >= "A" && key <= "Z") ||
			key == "'" || key == "-" || key == " ") {
			m.SpellSearchInput += key
			// Update suggestions
			m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput, "")
			if len(m.SpellSuggestions) > 0 {
				m.SuggestionIndex = 0
			} else {
				m.SuggestionIndex = -1
			}
		}
	} else if m.MonsterSearchMode && m.ActivePanel == Monsters {
		// Handle text input for monster search
		if len(key) == 1 && ((key >= "a" && key <= "z") ||
			(key >= "A" && key <= "Z") ||
			key == "'" || key == "-" || key == " ") {
			m.MonsterSearchInput += key
			// Update suggestions
			m.MonsterSuggestions = panels.SearchMonsters(m.MonsterSearchInput, m.MonsterCRFilter)
			if len(m.MonsterSuggestions) > 0 {
				m.MonsterSuggestionIndex = 0
			} else {
				m.MonsterSuggestionIndex = -1
			}
		}
	} else if m.ActivePanel == Notes && (m.NotesEditMode || m.NotesSearchMode) {
		// Handle text input for notes panel
		m = handleNotesInput(m, key)
	}

	return m, nil
}

// ========== SPECIAL HANDLERS ==========

// handleHelp toggles the help popup
func handleHelp(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	m.ShowHelpPopup = !m.ShowHelpPopup
	return m, nil
}
