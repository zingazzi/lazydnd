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
			Margin(1, 0).
			MaxWidth(35)
)

// GetDiceRollerContent returns the content for the dice roller panel
func GetDiceRollerContent(diceInput, diceResult string, diceHistory []string, diceCommands []string, lastCommand string, inputMode, isActive bool, historyMode bool, historyIndex int) string {
	// Show different instructions based on mode
	if historyMode {
		content := "HISTORY MODE - Select a roll to repeat\nUse ↑↓ to navigate, Enter to re-roll, Esc to exit"
		content += "\n\n" + strings.Repeat("─", 30)

		// Show history with selection
		if len(diceHistory) > 0 {
			content += "\n\nSelect a roll to repeat:"
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

	content := "Press Enter to roll, 'h' for history, Esc to clear\nExamples: 1d20, 2d8+3d6, 1d8-2, 1d6-1d4"

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
		// Split result by newlines first (for multi-line results)
		resultLines := strings.Split(diceResult, "\n")
		var allWrappedLines []string

		for _, line := range resultLines {
			// Wrap at 35 characters to fit within panel width
			wrappedLines := wrapText(line, 35)
			allWrappedLines = append(allWrappedLines, wrappedLines...)
		}

		wrappedResult := strings.Join(allWrappedLines, "\n")
		content += "\n" + diceResultStyle.Render("Last Roll:\n"+wrappedResult)
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

// parseComplexDice handles complex dice notation with single modifier (e.g., "2d6+3")
func parseComplexDice(command string, advantage, disadvantage bool) string {
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

		// Ensure minimum value of 1 (D&D rule: no negative results)
		if finalTotal < 1 {
			finalTotal = 1
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

	// Roll each expression
	var results []int
	var resultStrings []string

	for _, expr := range expressions {
		// Check if it's a dice expression or just a number
		if strings.Contains(expr, "d") {
			result := rollSingleDiceExpression(expr, advantage, disadvantage)
			if result == -1 {
				return "Invalid dice expression: " + expr
			}
			results = append(results, result)
			resultStrings = append(resultStrings, fmt.Sprintf("%d (%s)", result, expr))
		} else {
			// It's just a number
			num, err := strconv.Atoi(expr)
			if err != nil {
				return "Invalid number: " + expr
			}
			results = append(results, num)
			resultStrings = append(resultStrings, expr)
		}
	}

	// Calculate total
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

	// Ensure minimum value of 1 (D&D rule: no negative results)
	if total < 1 {
		total = 1
	}

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

	numDice := 1
	if parts[0] != "" {
		var err error
		numDice, err = strconv.Atoi(parts[0])
		if err != nil || numDice <= 0 || numDice > 100 {
			return -1
		}
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides <= 0 {
		return -1
	}

	// Validate dice type
	validDice := []int{4, 6, 8, 10, 12, 20, 100}
	isValid := false
	for _, valid := range validDice {
		if sides == valid {
			isValid = true
			break
		}
	}
	if !isValid {
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

	// Roll each expression
	var results []int
	var resultStrings []string

	for _, expr := range expressions {
		// Check if it's a dice expression or just a number
		if strings.Contains(expr, "d") {
			result := rollSingleDiceExpression(expr, false, false)
			if result == -1 {
				return "Invalid: " + expr
			}
			results = append(results, result)
			resultStrings = append(resultStrings, fmt.Sprintf("%d(%s)", result, expr))
		} else {
			// It's just a number
			num, err := strconv.Atoi(expr)
			if err != nil {
				return "Invalid: " + expr
			}
			results = append(results, num)
			resultStrings = append(resultStrings, expr)
		}
	}

	// Calculate total
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

	// Ensure minimum value of 1 (D&D rule: no negative results)
	if total < 1 {
		total = 1
	}

	// Format result compactly for comma-separated display
	breakdown := resultStrings[0]
	for i := 1; i < len(resultStrings); i++ {
		if i-1 < len(operators) {
			breakdown += operators[i-1] + resultStrings[i]
		}
	}

	return fmt.Sprintf("%d (%s)", total, breakdown)
}
