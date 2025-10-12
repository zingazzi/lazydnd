// panels/dice_roller.go
package panels

import (
	"fmt"
	"lazydnd/config"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Constants for dice rolling validation
const (
	MaxDicePerRoll = 100 // Maximum number of dice that can be rolled at once
	MinDiceValue   = 1   // Minimum value for dice count and sides
)

// Package-level config holder (set by RollDice function)
var (
	currentMinValue       = 1
	currentShowIndividual = true
)

// ValidDiceTypes defines the standard D&D dice types plus Zocchi's dice
// Standard D&D: d4, d6, d8, d10, d12, d20, d100
// Zocchi's dice: d3, d5, d7, d14, d16, d24, d30
var ValidDiceTypes = []int{3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 20, 24, 30, 100}

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
			Margin(1, 0).
			MaxWidth(35)
)

// GetDiceRollerContent returns the formatted content for the dice roller panel.
// It displays input mode, history mode, or the current dice result along with
// the dice roll history and available commands.
//
// Parameters:
//   - diceInput: Current user input (shown when inputMode is true)
//   - diceResult: Result of the last dice roll
//   - diceHistory: List of previous roll results
//   - diceCommands: List of commands that produced the history results
//   - lastCommand: The most recent dice command
//   - inputMode: Whether the user is currently typing
//   - isActive: Whether this panel is currently focused
//   - historyMode: Whether the user is browsing history to re-roll
//   - historyIndex: Selected history entry index (when historyMode is true)
func GetDiceRollerContent(diceInput, diceResult string, diceHistory []string, diceCommands []string, lastCommand string, inputMode, isActive bool, historyMode bool, historyIndex int) string {
	// Show different instructions based on mode
	if historyMode {
		content := ""

		// Show history with selection
		if len(diceHistory) > 0 {
			content += "History:"
			for i := len(diceHistory) - 1; i >= 0; i-- {
				var marker string
				if historyIndex == i {
					marker = "► "
				} else {
					marker = "  "
				}

				// Show command if available, otherwise show result
				displayText := diceHistory[i]
				if i < len(diceCommands) && diceCommands[i] != "" {
					displayText = diceCommands[i] + " → " + diceHistory[i]
				}

				content += "\n" + marker + displayText
			}
		} else {
			content += "\n\nNo history available"
		}

		return content
	}

	content := ""

	// Input field
	inputPrompt := "Dice: "
	if inputMode && isActive {
		inputPrompt += diceInput + "█"
	} else {
		inputPrompt += diceInput
	}
	content += inputStyle.Render(inputPrompt)

	// Last result
	if diceResult != "" {
		resultLines := strings.Split(diceResult, "\n")
		var allWrappedLines []string

		for _, line := range resultLines {
			wrappedLines := wrapText(line, 35)
			allWrappedLines = append(allWrappedLines, wrappedLines...)
		}

		wrappedResult := strings.Join(allWrappedLines, "\n")
		content += "\n" + diceResultStyle.Render(wrappedResult)
	}

	// History (most recent first)
	if len(diceHistory) > 0 {
		content += "\n\nHistory:"
		for i := len(diceHistory) - 1; i >= 0; i-- {
			content += "\n" + diceHistory[i]
		}
	}

	return content
}

// RollDice processes a dice roll command and returns the formatted result.
//
// Supported formats:
//   - Simple rolls: "d20", "2d6", "3d8"
//   - With modifiers: "1d20+5", "2d8-2"
//   - Advantage: "1d20 adv" or "1d20 advantage"
//   - Disadvantage: "1d20 dis" or "1d20 disadvantage"
//   - Complex expressions: "2d8+3d6", "1d20+3+2d4"
//   - Multiple rolls: "2d8, 3d6, 1d20" (comma-separated)
//
// Returns a formatted string with the roll result, or an error message if the command is invalid.
// All results include the dice notation and breakdown of individual rolls where applicable.
func RollDice(command string, cfg *config.Config) string {
	rand.Seed(time.Now().UnixNano())
	command = strings.TrimSpace(strings.ToLower(command))

	// Update package-level config values
	if cfg != nil {
		currentMinValue = cfg.DiceRoller.MinimumValue
		if currentMinValue < 0 {
			currentMinValue = 1
		}
		currentShowIndividual = cfg.DiceRoller.ShowIndividual
	} else {
		currentMinValue = 1
		currentShowIndividual = true
	}

	// Check for comma-separated rolls (e.g., "2d8, 3d6")
	if strings.Contains(command, ",") {
		return handleMultipleRolls(command)
	}

	// Check for advantage/disadvantage
	var advantage, disadvantage bool
	if strings.HasSuffix(command, " adv") || strings.HasSuffix(command, " advantage") || strings.HasSuffix(command, "adv") {
		advantage = true
		command = strings.TrimSuffix(command, " adv")
		command = strings.TrimSuffix(command, " advantage")
		command = strings.TrimSuffix(command, "adv")
	} else if strings.HasSuffix(command, " dis") || strings.HasSuffix(command, " disadvantage") || strings.HasSuffix(command, "dis") {
		disadvantage = true
		command = strings.TrimSuffix(command, " dis")
		command = strings.TrimSuffix(command, " disadvantage")
		command = strings.TrimSuffix(command, "dis")
	}

	// Clean up the command after removing adv/dis
	command = strings.TrimSpace(command)

	// Check for multiple dice expressions with + (e.g., "2d8+3d6")
	if hasMultipleDiceExpressions(command) {
		return handleMultipleDiceExpressions(command, advantage, disadvantage)
	}

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
		if advantage || disadvantage {
			roll1 := rand.Intn(20) + 1
			roll2 := rand.Intn(20) + 1
			var result int
			var resultType string

			if advantage {
				if roll1 >= roll2 {
					result = roll1
				} else {
					result = roll2
				}
				resultType = "ADV"
			} else {
				if roll1 <= roll2 {
					result = roll1
				} else {
					result = roll2
				}
				resultType = "DIS"
			}

			return fmt.Sprintf("d20 %s: %d (%d, %d)", resultType, result, roll1, roll2)
		} else {
			result := rand.Intn(20) + 1
			return fmt.Sprintf("d20: %d", result)
		}
	}
	if command == "d100" || command == "1d100" {
		result := rand.Intn(100) + 1
		return fmt.Sprintf("d100: %d", result)
	}

	// Handle complex dice notation like "2d6", "3d8+2"
	if strings.Contains(command, "d") {
		return parseComplexDice(command, advantage, disadvantage)
	}

	return "Invalid dice command"
}

// validateDiceType checks if the given dice type is a valid D&D die.
// Returns an error if the dice type is not in the ValidDiceTypes list.
func validateDiceType(sides int) error {
	for _, valid := range ValidDiceTypes {
		if sides == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid dice type: d%d (allowed: d3, d4, d5, d6, d7, d8, d10, d12, d14, d16, d20, d24, d30, d100)", sides)
}

// evaluateExpressions evaluates a list of dice/number expressions with operators.
// Returns the results array and formatted result strings, or an error.
func evaluateExpressions(expressions []string, advantage, disadvantage bool) ([]int, []string, error) {
	var results []int
	var resultStrings []string

	for _, expr := range expressions {
		// Check if it's a dice expression or just a number
		if strings.Contains(expr, "d") {
			result := rollSingleDiceExpression(expr, advantage, disadvantage)
			if result == -1 {
				return nil, nil, fmt.Errorf("invalid dice expression: %s", expr)
			}
			results = append(results, result)
			resultStrings = append(resultStrings, fmt.Sprintf("%d (%s)", result, expr))
		} else {
			// It's just a number
			num, err := strconv.Atoi(expr)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid number: %s", expr)
			}
			results = append(results, num)
			resultStrings = append(resultStrings, expr)
		}
	}

	return results, resultStrings, nil
}

// calculateTotal calculates the total from results array using the provided operators.
// Ensures minimum value of 1 (D&D rule: no negative results).
func calculateTotal(results []int, operators []string) int {
	if len(results) == 0 {
		return 1
	}

	total := results[0]
	for i := 1; i < len(results); i++ {
		if i-1 < len(operators) {
			if operators[i-1] == "+" {
				total += results[i]
			} else if operators[i-1] == "-" {
				total -= results[i]
			}
		}
	}

	// Ensure minimum value (configurable, D&D standard is 1)
	if total < currentMinValue {
		total = currentMinValue
	}

	return total
}

// parseComplexDice handles complex dice notation with single modifier (e.g., "2d6+3").
// Supports advantage/disadvantage modifiers for d20 rolls.
// Returns a formatted string with the roll result or an error message.
func parseComplexDice(command string, advantage, disadvantage bool) string {
	// Remove spaces
	command = strings.ReplaceAll(command, " ", "")

	// Split by 'd'
	parts := strings.Split(command, "d")
	if len(parts) != 2 {
		return "Invalid format. Use: XdY or XdY+Z"
	}

	// Parse number of dice
	numDice := MinDiceValue
	if parts[0] != "" {
		var err error
		numDice, err = strconv.Atoi(parts[0])
		if err != nil || numDice < MinDiceValue || numDice > MaxDicePerRoll {
			return fmt.Sprintf("Invalid number of dice (%d-%d)", MinDiceValue, MaxDicePerRoll)
		}
	}

	// Parse dice type and modifier
	modifier := 0
	diceType := parts[1]
	diceExpr := fmt.Sprintf("%dd%s", numDice, parts[1])
	if parts[0] == "" {
		diceExpr = "d" + parts[1]
	}

	if strings.Contains(diceType, "+") {
		modParts := strings.Split(diceType, "+")
		if len(modParts) == 2 {
			diceType = modParts[0]
			var err error
			modifier, err = strconv.Atoi(modParts[1])
			if err != nil {
				return "Invalid modifier"
			}
			diceExpr = fmt.Sprintf("%dd%s", numDice, diceType)
			if parts[0] == "" {
				diceExpr = "d" + diceType
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
			diceExpr = fmt.Sprintf("%dd%s", numDice, diceType)
			if parts[0] == "" {
				diceExpr = "d" + diceType
			}
		}
	}

	sides, err := strconv.Atoi(diceType)
	if err != nil || sides < MinDiceValue {
		return "Invalid dice type"
	}

	// Validate dice type
	if err := validateDiceType(sides); err != nil {
		return err.Error()
	}

	// Handle advantage/disadvantage
	if advantage || disadvantage {
		// Roll twice and take higher/lower
		total1, rolls1 := rollDiceSet(numDice, sides)
		total2, rolls2 := rollDiceSet(numDice, sides)

		var finalTotal int
		var resultType string

		if advantage {
			if total1 >= total2 {
				finalTotal = total1
			} else {
				finalTotal = total2
			}
			resultType = "ADV"
		} else {
			if total1 <= total2 {
				finalTotal = total1
			} else {
				finalTotal = total2
			}
			resultType = "DIS"
		}

		// Add modifier
		finalResult := finalTotal + modifier

		// Format result
		rollsStr1 := make([]string, len(rolls1))
		rollsStr2 := make([]string, len(rolls2))
		for i := range rolls1 {
			rollsStr1[i] = strconv.Itoa(rolls1[i])
			rollsStr2[i] = strconv.Itoa(rolls2[i])
		}

		// Format result - Total first, then breakdown
		result := fmt.Sprintf("TOTAL: %d (%s)\n", finalResult, resultType)
		result += fmt.Sprintf("%d (%s): [%s] vs [%s]", finalTotal, diceExpr,
			strings.Join(rollsStr1, ", "), strings.Join(rollsStr2, ", "))
		if modifier != 0 {
			result += fmt.Sprintf(" %+d", modifier)
		}

		return result
	} else {
		// Normal roll
		total := 0
		rolls := make([]int, numDice)
		for i := 0; i < numDice; i++ {
			roll := rand.Intn(sides) + 1
			rolls[i] = roll
			total += roll
		}

		// Format result
		rollsStr := make([]string, len(rolls))
		for i, roll := range rolls {
			rollsStr[i] = strconv.Itoa(roll)
		}

		finalTotal := total + modifier

		// Ensure minimum value (configurable, D&D standard is 1)
		if finalTotal < currentMinValue {
			finalTotal = currentMinValue
		}

		// Format result - Total first, then breakdown
		result := fmt.Sprintf("TOTAL: %d\n", finalTotal)
		result += fmt.Sprintf("%d (%s): [%s]", total, diceExpr, strings.Join(rollsStr, ", "))
		if modifier != 0 {
			result += fmt.Sprintf(" %+d", modifier)
		}

		return result
	}
}

// rollDiceSet rolls a set of dice and returns total and individual rolls
func rollDiceSet(numDice, sides int) (int, []int) {
	total := 0
	rolls := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		roll := rand.Intn(sides) + 1
		rolls[i] = roll
		total += roll
	}
	return total, rolls
}

// hasMultipleDiceExpressions checks if command has multiple dice/number expressions (e.g., "2d8+3d6" or "2d8+3")
func hasMultipleDiceExpressions(command string) bool {
	// Remove spaces
	command = strings.ReplaceAll(command, " ", "")

	// Check if there's a '+' or '-' in the command
	return strings.ContainsAny(command, "+-") && strings.Contains(command, "d")
}

// handleMultipleDiceExpressions handles expressions like "2d8+3d6" or "2d8+3+3d6"
func handleMultipleDiceExpressions(command string, advantage, disadvantage bool) string {
	// Remove spaces
	command = strings.ReplaceAll(command, " ", "")

	// Split by + and - while keeping track of operators
	var expressions []string
	var operators []string
	currentExpr := ""

	for i := 0; i < len(command); i++ {
		if (command[i] == '+' || command[i] == '-') && i > 0 && currentExpr != "" {
			expressions = append(expressions, currentExpr)
			operators = append(operators, string(command[i]))
			currentExpr = ""
		} else {
			currentExpr += string(command[i])
		}
	}
	if currentExpr != "" {
		expressions = append(expressions, currentExpr)
	}

	// Evaluate all expressions
	results, resultStrings, err := evaluateExpressions(expressions, advantage, disadvantage)
	if err != nil {
		return err.Error()
	}

	// Calculate total
	total := calculateTotal(results, operators)

	// Format result - Total first, then breakdown
	resultStr := fmt.Sprintf("TOTAL: %d", total)

	if advantage {
		resultStr += " (ADV)"
	} else if disadvantage {
		resultStr += " (DIS)"
	}

	// Add breakdown on new line
	resultStr += "\n"
	resultStr += resultStrings[0]
	for i := 1; i < len(resultStrings); i++ {
		if i-1 < len(operators) {
			resultStr += " " + operators[i-1] + " " + resultStrings[i]
		}
	}

	return resultStr
}

// rollSingleDiceExpression rolls a single dice expression and returns the total
func rollSingleDiceExpression(expr string, advantage, disadvantage bool) int {
	// Parse the expression
	parts := strings.Split(expr, "d")
	if len(parts) != 2 {
		return -1
	}

	numDice := MinDiceValue
	if parts[0] != "" {
		var err error
		numDice, err = strconv.Atoi(parts[0])
		if err != nil || numDice < MinDiceValue || numDice > MaxDicePerRoll {
			return -1
		}
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides < MinDiceValue {
		return -1
	}

	// Validate dice type
	if err := validateDiceType(sides); err != nil {
		return -1
	}

	// Roll dice
	if advantage || disadvantage {
		total1, _ := rollDiceSet(numDice, sides)
		total2, _ := rollDiceSet(numDice, sides)

		if advantage {
			if total1 >= total2 {
				return total1
			}
			return total2
		} else {
			if total1 <= total2 {
				return total1
			}
			return total2
		}
	}

	total, _ := rollDiceSet(numDice, sides)
	return total
}

// handleMultipleRolls handles comma-separated rolls (e.g., "2d8, 3d6" or "1d8+3, 3d6+1")
func handleMultipleRolls(command string) string {
	// Split by comma
	rolls := strings.Split(command, ",")

	var results []string
	for _, roll := range rolls {
		roll = strings.TrimSpace(roll)
		if roll == "" {
			continue
		}

		// Check if this roll is a complex expression (has + or -)
		if hasMultipleDiceExpressions(roll) || (strings.ContainsAny(roll, "+-") && strings.Contains(roll, "d")) {
			// Handle as complex expression
			result := rollComplexExpression(roll)
			results = append(results, result)
		} else {
			// Simple single dice expression
			total := rollSingleDiceExpression(roll, false, false)
			if total == -1 {
				return "Invalid dice expression: " + roll
			}
			results = append(results, fmt.Sprintf("%d (%s)", total, roll))
		}
	}

	return strings.Join(results, " | ")
}

// rollComplexExpression rolls a complex expression and returns formatted result
func rollComplexExpression(command string) string {
	// Remove spaces
	command = strings.ReplaceAll(command, " ", "")

	// Split by + and - while keeping track of operators
	var expressions []string
	var operators []string
	currentExpr := ""

	for i := 0; i < len(command); i++ {
		if (command[i] == '+' || command[i] == '-') && i > 0 && currentExpr != "" {
			expressions = append(expressions, currentExpr)
			operators = append(operators, string(command[i]))
			currentExpr = ""
		} else {
			currentExpr += string(command[i])
		}
	}
	if currentExpr != "" {
		expressions = append(expressions, currentExpr)
	}

	// Evaluate all expressions (no advantage/disadvantage for comma-separated rolls)
	results, resultStrings, err := evaluateExpressions(expressions, false, false)
	if err != nil {
		return "Invalid: " + err.Error()
	}

	// Calculate total
	total := calculateTotal(results, operators)

	// Format result compactly for comma-separated display
	breakdown := resultStrings[0]
	for i := 1; i < len(resultStrings); i++ {
		if i-1 < len(operators) {
			breakdown += operators[i-1] + resultStrings[i]
		}
	}

	return fmt.Sprintf("%d (%s)", total, breakdown)
}
