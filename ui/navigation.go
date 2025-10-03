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
		if !m.InputMode {
			return m, tea.Quit
		}

	case "1":
		if m.InputMode {
			m.DiceInput += "1"
		} else if m.InitiativeInputMode {
			m.InitiativeInput += "1"
		} else {
			m.ActivePanel = DiceRoller
		}
	case "2":
		if m.InputMode {
			m.DiceInput += "2"
		} else if m.InitiativeInputMode {
			m.InitiativeInput += "2"
		} else {
			m.ActivePanel = InitiativeTracker
		}
	case "3":
		if m.InputMode {
			m.DiceInput += "3"
		} else if m.InitiativeInputMode {
			m.InitiativeInput += "3"
		} else {
			m.ActivePanel = Spells
		}
	case "4":
		if m.InputMode {
			m.DiceInput += "4"
		} else if m.InitiativeInputMode {
			m.InitiativeInput += "4"
		} else {
			m.ActivePanel = CampaignNotes
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
		m.ActivePanel = CampaignNotes
		m.InputMode = false

	case "up":
		// Handle spell suggestion navigation first (highest priority)
		if m.ActivePanel == Spells && m.SpellSearchMode && len(m.SpellSuggestions) > 0 {
			if m.SuggestionIndex > 0 {
				m.SuggestionIndex--
			}
		} else if !m.InputMode {
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
		} else if !m.InputMode {
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
			if m.InitiativeInputMode {
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
		}

	case "r":
		if m.ActivePanel == DiceRoller && !m.InputMode && !m.InitiativeInputMode && m.LastDiceCommand != "" {
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
			m.InitiativeInput = ""
			m.InitiativeInputMode = false
			m.InitiativeInputType = ""
		} else if m.ActivePanel == Spells {
			m.SpellSearchInput = ""
			m.SpellSearchMode = false
			m.SpellSuggestions = []string{}
			m.SuggestionIndex = -1
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
		}

	case "p":
		if m.ActivePanel == InitiativeTracker && !m.InitiativeInputMode {
			// Start adding a player
			m.InitiativeInputMode = true
			m.InitiativeInputType = "player_name"
			m.InitiativeInput = ""
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Add 'p' to input when in input mode
			m.InitiativeInput += "p"
		}

	case "m":
		if m.ActivePanel == InitiativeTracker && !m.InitiativeInputMode {
			// Start adding a monster
			m.InitiativeInputMode = true
			m.InitiativeInputType = "monster_name"
			m.InitiativeInput = ""
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Add 'm' to input when in input mode
			m.InitiativeInput += "m"
		}

	case "space":
		if m.InputMode && m.ActivePanel == DiceRoller {
			m.DiceInput += " "
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
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
		} else if m.InitiativeInputMode && m.ActivePanel == InitiativeTracker {
			// Handle text input for initiative tracker
			key := msg.String()
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				(key >= "0" && key <= "9") ||
				key == " " || key == "'" || key == "-" || key == "." || key == "_") {
				m.InitiativeInput += key
			}
		} else if m.SpellSearchMode && m.ActivePanel == Spells {
			// Handle text input for spell search
			key := msg.String()
			if len(key) == 1 && (
				(key >= "a" && key <= "z") ||
				(key >= "A" && key <= "Z") ||
				key == "'" || key == "-") {
				m.SpellSearchInput += key
				// Update suggestions
				m.SpellSuggestions = panels.SearchSpells(m.SpellSearchInput)
				if len(m.SpellSuggestions) > 0 {
					m.SuggestionIndex = 0
				} else {
					m.SuggestionIndex = -1
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
