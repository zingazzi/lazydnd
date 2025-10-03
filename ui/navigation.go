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
		} else {
			m.ActivePanel = DiceRoller
		}
	case "2":
		if m.InputMode {
			m.DiceInput += "2"
		} else {
			m.ActivePanel = CharacterSheet
		}
	case "3":
		if m.InputMode {
			m.DiceInput += "3"
		} else {
			m.ActivePanel = Spells
		}
	case "4":
		if m.InputMode {
			m.DiceInput += "4"
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
		m.ActivePanel = CharacterSheet
		m.InputMode = false
	case "f3":
		m.ActivePanel = Spells
		m.InputMode = false
	case "f4":
		m.ActivePanel = CampaignNotes
		m.InputMode = false

	case "up":
		if !m.InputMode {
			if m.ScrollOffset[m.ActivePanel] > 0 {
				m.ScrollOffset[m.ActivePanel]--
			}
		}

	case "down":
		if !m.InputMode {
			m.ScrollOffset[m.ActivePanel]++
		}

	case "enter":
		if m.ActivePanel == DiceRoller {
			if m.InputMode && m.DiceInput != "" {
				// Roll the dice
				result := panels.RollDice(m.DiceInput)
				m.DiceResult = result
				m.DiceHistory = append(m.DiceHistory, result)
				if len(m.DiceHistory) > 15 {
					m.DiceHistory = m.DiceHistory[1:]
				}
				m.DiceInput = ""
				m.InputMode = false
			} else {
				m.InputMode = true
			}
		}

	case "esc":
		if m.ActivePanel == DiceRoller {
			m.DiceInput = ""
			m.InputMode = false
		}

	case "backspace", "ctrl+h":
		if m.InputMode && len(m.DiceInput) > 0 {
			m.DiceInput = m.DiceInput[:len(m.DiceInput)-1]
		}

	case "space":
		if m.InputMode && m.ActivePanel == DiceRoller {
			m.DiceInput += " "
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
		}
	}

	return m, nil
}
