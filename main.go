// main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Panel types
type PanelType int

const (
	DiceRoller PanelType = iota
	CharacterSheet
	Spells
	CampaignNotes
)

// Model represents the main application state
type Model struct {
	activePanel    PanelType
	width          int
	height         int
	diceInput      string
	diceResult     string
	diceHistory    []string
	inputMode      bool
	scrollOffset   map[PanelType]int
}

// Panel configuration
var panelNames = []string{
	"Dice Roller",
	"Character Sheet",
	"Spells",
	"Campaign Notes",
}

var panelContent = []string{
	"🎲 DICE ROLLER 🎲\n\nType dice commands like:\n• d20 (single die)\n• 2d8 (multiple dice)\n• 3d6+4 (with modifier)\n• 1d20-2 (with penalty)\n\nPress Enter to roll, Esc to clear input",
	"⚔️ CHARACTER SHEET ⚔️\n\nName: Thorin Oakenshield\nClass: Fighter\nLevel: 5\nRace: Dwarf\n\nStats:\nSTR: 16 (+3)\nDEX: 12 (+1)\nCON: 15 (+2)\nINT: 10 (+0)\nWIS: 13 (+1)\nCHA: 8 (-1)\n\nHP: 45/45\nAC: 18",
	"✨ SPELLS ✨\n\nCantrips:\n• Firebolt (1d10 fire)\n• Mage Hand\n• Prestidigitation\n\nLevel 1 (3/3):\n• Magic Missile\n• Shield\n• Detect Magic\n\nLevel 2 (2/2):\n• Misty Step\n• Scorching Ray",
	"📖 CAMPAIGN NOTES 📖\n\nSession 12: The Lost Temple\n\n• Found ancient dwarven ruins\n• Encountered goblin patrol (defeated)\n• Discovered magical artifact\n• Next: Explore deeper chambers\n\nNPCs:\n• Elara the Wise (Elf Wizard)\n• Gareth Ironforge (Dwarf Blacksmith)\n\nQuests:\n• Retrieve the Crystal of Power\n• Rescue the missing villagers",
}

// Styles
var (
	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(1, 2)

	inactivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#444444")).
				Padding(1, 2)

	panelTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Margin(0, 0, 1, 0)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Margin(1, 0)

	diceResultStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FF00")).
				Background(lipgloss.Color("#1A1A1A")).
				Padding(0, 1).
				Margin(1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Margin(1, 0, 0, 0)
)

func initialModel() Model {
	return Model{
		activePanel:  DiceRoller,
		diceInput:    "",
		diceResult:   "",
		diceHistory:  []string{},
		inputMode:    false,
		scrollOffset: make(map[PanelType]int),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.inputMode {
				return m, tea.Quit
			}

		case "1":
			if m.inputMode {
				m.diceInput += "1"
			} else {
				m.activePanel = DiceRoller
			}
		case "2":
			if m.inputMode {
				m.diceInput += "2"
			} else {
				m.activePanel = CharacterSheet
			}
		case "3":
			if m.inputMode {
				m.diceInput += "3"
			} else {
				m.activePanel = Spells
			}
		case "4":
			if m.inputMode {
				m.diceInput += "4"
			} else {
				m.activePanel = CampaignNotes
			}

		case "tab":
			if !m.inputMode {
				m.activePanel = (m.activePanel + 1) % 4
			}

		case "f1":
			m.activePanel = DiceRoller
			m.inputMode = false
		case "f2":
			m.activePanel = CharacterSheet
			m.inputMode = false
		case "f3":
			m.activePanel = Spells
			m.inputMode = false
		case "f4":
			m.activePanel = CampaignNotes
			m.inputMode = false

		case "up":
			if !m.inputMode {
				if m.scrollOffset[m.activePanel] > 0 {
					m.scrollOffset[m.activePanel]--
				}
			}

		case "down":
			if !m.inputMode {
				m.scrollOffset[m.activePanel]++
			}

		case "enter":
			if m.activePanel == DiceRoller {
				if m.inputMode && m.diceInput != "" {
					// Roll the dice
					result := rollDice(m.diceInput)
					m.diceResult = result
					m.diceHistory = append(m.diceHistory, result)
					if len(m.diceHistory) > 15 {
						m.diceHistory = m.diceHistory[1:]
					}
					m.diceInput = ""
					m.inputMode = false
				} else {
					m.inputMode = true
				}
			}

		case "esc":
			if m.activePanel == DiceRoller {
				m.diceInput = ""
				m.inputMode = false
			}

		case "backspace", "ctrl+h":
			if m.inputMode && len(m.diceInput) > 0 {
				m.diceInput = m.diceInput[:len(m.diceInput)-1]
			}

		case "space":
			if m.inputMode && m.activePanel == DiceRoller {
				m.diceInput += " "
			}

		default:
			// Handle text input for dice commands
			if m.inputMode && m.activePanel == DiceRoller {
				key := msg.String()
				// Allow alphanumeric characters and common symbols for dice notation
				if len(key) == 1 && (
					(key >= "a" && key <= "z") ||
					(key >= "A" && key <= "Z") ||
					(key >= "0" && key <= "9") ||
					key == "+" || key == "-" || key == "d") {
					m.diceInput += key
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate panel dimensions
	panelWidth := (m.width - 6) / 2
	panelHeight := (m.height - 4) / 2

	// Create panels
	panels := make([]string, 4)

	for i := 0; i < 4; i++ {
		panelType := PanelType(i)
		content := panelContent[i]

		// Add special content for dice roller
		if panelType == DiceRoller {
			content += "\n\n" + strings.Repeat("─", 30)

			// Input field
			inputPrompt := "Dice Command: "
			if m.inputMode && m.activePanel == DiceRoller {
				inputPrompt += m.diceInput + "█"
			} else {
				inputPrompt += m.diceInput
			}
			content += "\n\n" + inputStyle.Render(inputPrompt)

			// Last result
			if m.diceResult != "" {
				content += "\n" + diceResultStyle.Render("Last Roll: "+m.diceResult)
			}

			// History
			if len(m.diceHistory) > 0 {
				content += "\n\nHistory:"
				for _, roll := range m.diceHistory {
					content += "\n• " + roll
				}
			}
		}

		// Add help text
		if panelType == m.activePanel {
			if panelType == DiceRoller {
				if m.inputMode {
					content += "\n" + helpStyle.Render("Enter: roll • Esc: cancel • F1-F4: switch panels")
				} else {
					content += "\n" + helpStyle.Render("Enter: input dice • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
				}
			} else {
				content += "\n" + helpStyle.Render("↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
			}
		}

		// Apply scrolling
		contentLines := strings.Split(content, "\n")
		scrollOffset := m.scrollOffset[panelType]

		// Calculate available content height (panel height minus title and padding)
		availableHeight := panelHeight - 4 // Account for title, borders, and padding

		// Apply scroll offset
		if scrollOffset > 0 && scrollOffset < len(contentLines) {
			if scrollOffset+availableHeight < len(contentLines) {
				contentLines = contentLines[scrollOffset : scrollOffset+availableHeight]
			} else {
				contentLines = contentLines[scrollOffset:]
			}
		} else if scrollOffset >= len(contentLines) {
			// Reset scroll if we've gone too far
			m.scrollOffset[panelType] = 0
		} else {
			// Show from beginning
			if len(contentLines) > availableHeight {
				contentLines = contentLines[:availableHeight]
			}
		}

		scrolledContent := strings.Join(contentLines, "\n")

		// Style the panel
		title := fmt.Sprintf(" %d. %s ", i+1, panelNames[i])
		titleBar := panelTitleStyle.Render(title)

		var panelStyle lipgloss.Style
		if panelType == m.activePanel {
			panelStyle = activePanelStyle
		} else {
			panelStyle = inactivePanelStyle
		}

		panelContent := titleBar + "\n" + scrolledContent
		panels[i] = panelStyle.Width(panelWidth).Height(panelHeight).Render(panelContent)
	}

	// Arrange panels in 2x2 grid
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panels[0], panels[1])
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panels[2], panels[3])

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

// rollDice handles dice rolling logic
func rollDice(command string) string {
	rand.Seed(time.Now().UnixNano())
	command = strings.TrimSpace(strings.ToLower(command))

	// Handle simple dice notation
	if command == "d4" || command == "1d4" {
		result := rand.Intn(4) + 1
		return fmt.Sprintf("d4: %d", result)
	}
	if command == "d6" || command == "1d6" {
		result := rand.Intn(6) + 1
		return fmt.Sprintf("d6: %d", result)
	}
	if command == "d8" || command == "1d8" {
		result := rand.Intn(8) + 1
		return fmt.Sprintf("d8: %d", result)
	}
	if command == "d10" || command == "1d10" {
		result := rand.Intn(10) + 1
		return fmt.Sprintf("d10: %d", result)
	}
	if command == "d12" || command == "1d12" {
		result := rand.Intn(12) + 1
		return fmt.Sprintf("d12: %d", result)
	}
	if command == "d20" || command == "1d20" {
		result := rand.Intn(20) + 1
		return fmt.Sprintf("d20: %d", result)
	}

	// Handle complex dice notation like "2d6", "3d8+2"
	if strings.Contains(command, "d") {
		return parseComplexDice(command)
	}

	return "Invalid dice command"
}

// parseComplexDice handles complex dice notation
func parseComplexDice(command string) string {
	// Remove spaces
	command = strings.ReplaceAll(command, " ", "")

	// Split by 'd'
	parts := strings.Split(command, "d")
	if len(parts) != 2 {
		return "Invalid format. Use: XdY or XdY+Z"
	}

	// Parse number of dice
	numDice := 1
	if parts[0] != "" {
		var err error
		numDice, err = strconv.Atoi(parts[0])
		if err != nil || numDice <= 0 || numDice > 100 {
			return "Invalid number of dice (1-100)"
		}
	}

	// Parse dice type and modifier
	modifier := 0
	diceType := parts[1]

	if strings.Contains(diceType, "+") {
		modParts := strings.Split(diceType, "+")
		if len(modParts) == 2 {
			diceType = modParts[0]
			var err error
			modifier, err = strconv.Atoi(modParts[1])
			if err != nil {
				return "Invalid modifier"
			}
		}
	} else if strings.Contains(diceType, "-") {
		modParts := strings.Split(diceType, "-")
		if len(modParts) == 2 {
			diceType = modParts[0]
			mod, err := strconv.Atoi(modParts[1])
			if err != nil {
				return "Invalid modifier"
			}
			modifier = -mod
		}
	}

	sides, err := strconv.Atoi(diceType)
	if err != nil || sides <= 0 {
		return "Invalid dice type"
	}

	// Validate allowed dice types
	validDice := []int{4, 6, 8, 10, 12, 20, 100}
	isValid := false
	for _, valid := range validDice {
		if sides == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Sprintf("Invalid dice type: d%d (allowed: d4, d6, d8, d10, d12, d20, d100)", sides)
	}

	// Roll the dice
	total := 0
	rolls := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		roll := rand.Intn(sides) + 1
		rolls[i] = roll
		total += roll
	}

	total += modifier

	// Format result
	rollsStr := make([]string, len(rolls))
	for i, roll := range rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	result := fmt.Sprintf("%s: [%s]", command, strings.Join(rollsStr, ", "))
	if modifier != 0 {
		result += fmt.Sprintf(" %+d", modifier)
	}
	result += fmt.Sprintf(" = %d", total)

	return result
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
