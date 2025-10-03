// panels/dice_roller.go
package panels

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
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
)

// GetDiceRollerContent returns the content for the dice roller panel
func GetDiceRollerContent(diceInput, diceResult string, diceHistory []string, lastCommand string, inputMode, isActive bool) string {
	content := "🎲 DICE ROLLER 🎲\nPress Enter to roll, Esc to clear input"

	content += "\n\n" + strings.Repeat("─", 30)

	// Input field
	inputPrompt := "Dice Command: "
	if inputMode && isActive {
		inputPrompt += diceInput + "█"
	} else {
		inputPrompt += diceInput
	}
	content += "\n\n" + inputStyle.Render(inputPrompt)

	// Last result
	if diceResult != "" {
		content += "\n" + diceResultStyle.Render("Last Roll: "+diceResult)
	}

	// Show last command for reroll
	if lastCommand != "" && !inputMode {
		content += "\n\nLast Command: " + lastCommand + " (press 'r' to reroll)"
	}

	// History (most recent first)
	if len(diceHistory) > 0 {
		content += "\n\nHistory:"
		for i := len(diceHistory) - 1; i >= 0; i-- {
			content += "\n• " + diceHistory[i]
		}
	}

	return content
}

// RollDice handles dice rolling logic
func RollDice(command string) string {
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
	if command == "d100" || command == "1d100" {
		result := rand.Intn(100) + 1
		return fmt.Sprintf("d100: %d", result)
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
