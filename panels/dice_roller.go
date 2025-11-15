// panels/dice_roller.go
package panels

import (
	"fmt"
	"lazydnd/config"
	"math/rand"
	"strconv"
	"strings"
	"time"

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
	currentCritEnabled    = true
	currentCritMode       = "double"
)

// ValidDiceTypes defines the standard D&D dice types plus Zocchi's dice
// Standard D&D: d4, d6, d8, d10, d12, d20, d100
// Zocchi's dice: d3, d5, d7, d14, d16, d24, d30
var ValidDiceTypes = []int{3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 20, 24, 30, 100}

// Note: Styling is now handled by TView widgets - these functions return plain text

// getCriticalHitBanner returns a simple critical hit indicator
func getCriticalHitBanner() string {
	return "★ CRIT"
}

// handleCriticalDamageRoll handles rolling critical damage based on the configured mode
// Supports two modes:
//   - "double": Roll all damage dice twice (standard D&D 5e)
//   - "max": Maximum damage + one roll (popular house rule)
func handleCriticalDamageRoll(command string) string {
	// Remove all spaces first
	command = strings.ReplaceAll(command, " ", "")

	// Parse the dice expression
	parts := strings.Split(command, "d")
	if len(parts) != 2 {
		return "Invalid crit format. Use: XdY crit (e.g., 2d8 crit)"
	}

	numDice := 1
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
	if err != nil || sides < MinDiceValue {
		return "Invalid dice type"
	}

	if err := validateDiceType(sides); err != nil {
		return err.Error()
	}

	// Roll critical damage based on mode
	var total int
	var rolls []int

	if currentCritMode == "max" {
		// Max damage mode: max value for all dice + one normal roll
		maxDamage := numDice * sides
		normalTotal, normalRolls := rollDiceSet(numDice, sides)
		total = maxDamage + normalTotal
		rolls = normalRolls
	} else {
		// Double dice mode: roll all dice twice
		total1, rolls1 := rollDiceSet(numDice, sides)
		total2, rolls2 := rollDiceSet(numDice, sides)
		total = total1 + total2
		rolls = append(rolls1, rolls2...)
	}

	finalTotal := total + modifier

	// Format: "TOTAL  rolls (formula) +mod CRIT"
	rollsStr := make([]string, len(rolls))
	for i, roll := range rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	result := fmt.Sprintf("%d", finalTotal) + "  "
	result += strings.Join(rollsStr, ", ") + " "

	formula := fmt.Sprintf("(%dd%d × 2)", numDice, sides)
	if modifier != 0 {
		formula += fmt.Sprintf(" %+d", modifier)
	}
	result += formula
	result += " " + getCriticalHitBanner()

	return result
}

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
	content += inputPrompt

	// Last result - already has styling applied from RollDice
	if diceResult != "" {
		content += "\n" + diceResult
	}

	// History (most recent first) - show only totals
	if len(diceHistory) > 0 {
		// Add separator line and extra spacing (styling handled by TView)
		content += "\n\n─────────────────────────────────"
		content += "\nRecent:"

		// Show max 3 most recent results
		maxHistory := 3
		count := 0
		for i := len(diceHistory) - 1; i >= 0 && count < maxHistory; i-- {
			content += "\n" + diceHistory[i]
			count++
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
		currentCritEnabled = cfg.DiceRoller.CriticalHitEnabled
		currentCritMode = cfg.DiceRoller.CriticalHitMode
	} else {
		currentMinValue = 1
		currentShowIndividual = true
		currentCritEnabled = true
		currentCritMode = "double"
	}

	// Check for comma-separated rolls (e.g., "2d8, 3d6")
	if strings.Contains(command, ",") {
		return handleMultipleRolls(command)
	}

	// Check for critical hit modifier
	var critRoll bool
	if strings.HasSuffix(command, " crit") || strings.HasSuffix(command, " critical") {
		critRoll = true
		command = strings.TrimSuffix(command, " crit")
		command = strings.TrimSuffix(command, " critical")
		command = strings.TrimSpace(command)
	}

	// Handle crit rolls for damage dice FIRST (before other parsing)
	if critRoll && currentCritEnabled {
		return handleCriticalDamageRoll(command)
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

	// Handle simple dice notation with standard format
	if command == "d4" || command == "1d4" {
		roll := rand.Intn(4) + 1
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d4)"
	}
	if command == "d6" || command == "1d6" {
		roll := rand.Intn(6) + 1
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d6)"
	}
	if command == "d8" || command == "1d8" {
		roll := rand.Intn(8) + 1
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d8)"
	}
	if command == "d10" || command == "1d10" {
		roll := rand.Intn(10) + 1
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d10)"
	}
	if command == "d12" || command == "1d12" {
		roll := rand.Intn(12) + 1
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d12)"
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

			// Check for crit
			isCrit := currentCritEnabled && result == 20
			resultStr := fmt.Sprintf("d20 %s: %d (%d, %d)", resultType, result, roll1, roll2)
			if isCrit {
				resultStr += " " + getCriticalHitBanner()
			}
			return resultStr
		} else {
			roll := rand.Intn(20) + 1
			isCrit := currentCritEnabled && roll == 20

			// Use standard format
			if isCrit {
				// Critical hit format: "20  (d20) ★ CRIT" - in red
				return fmt.Sprintf("%d", roll) + "  " + "(d20)" + " " + getCriticalHitBanner()
			} else {
				// Normal format: "15  15 (d20)"
				return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d20)"
			}
		}
	}
	if command == "d100" || command == "1d100" {
		roll := rand.Intn(100) + 1
		// Use standard format: "75  75 (d100)"
		return fmt.Sprintf("%d", roll) + "  " + fmt.Sprintf("%d", roll) + " " + "(d100)"
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

// rollSingleDiceExpressionWithDetails rolls a dice expression and returns total, rolls, and whether it's a crit
func rollSingleDiceExpressionWithDetails(expr string, advantage, disadvantage bool) (int, []int, bool) {
	// Parse the expression
	parts := strings.Split(expr, "d")
	if len(parts) != 2 {
		return -1, nil, false
	}

	numDice := 1
	if parts[0] != "" {
		var err error
		numDice, err = strconv.Atoi(parts[0])
		if err != nil || numDice < MinDiceValue || numDice > MaxDicePerRoll {
			return -1, nil, false
		}
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides < MinDiceValue {
		return -1, nil, false
	}

	if err := validateDiceType(sides); err != nil {
		return -1, nil, false
	}

	// Roll dice
	if advantage || disadvantage {
		total1, rolls1 := rollDiceSet(numDice, sides)
		total2, rolls2 := rollDiceSet(numDice, sides)

		if advantage {
			if total1 >= total2 {
				isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls1[0] == 20
				return total1, rolls1, isCrit
			}
			isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls2[0] == 20
			return total2, rolls2, isCrit
		} else {
			if total1 <= total2 {
				isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls1[0] == 20
				return total1, rolls1, isCrit
			}
			isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls2[0] == 20
			return total2, rolls2, isCrit
		}
	}

	total, rolls := rollDiceSet(numDice, sides)
	isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls[0] == 20
	return total, rolls, isCrit
}

// evaluateExpressions evaluates a list of dice/number expressions with operators.
// Returns the results array, formatted result strings, whether any roll was a crit, or an error.
func evaluateExpressions(expressions []string, advantage, disadvantage bool) ([]int, []string, bool, error) {
	var results []int
	var resultStrings []string
	hasCrit := false

	for _, expr := range expressions {
		// Check if it's a dice expression or just a number
		if strings.Contains(expr, "d") {
			result, _, isCrit := rollSingleDiceExpressionWithDetails(expr, advantage, disadvantage)
			if result == -1 {
				return nil, nil, false, fmt.Errorf("invalid dice expression: %s", expr)
			}
			if isCrit {
				hasCrit = true
			}
			results = append(results, result)
			resultStrings = append(resultStrings, fmt.Sprintf("%d (%s)", result, expr))
		} else {
			// It's just a number
			num, err := strconv.Atoi(expr)
			if err != nil {
				return nil, nil, false, fmt.Errorf("invalid number: %s", expr)
			}
			results = append(results, num)
			resultStrings = append(resultStrings, expr)
		}
	}

	return results, resultStrings, hasCrit, nil
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

		// Format: "TOTAL  diceRolls (formula) +mod TYPE"
		result := fmt.Sprintf("%d", finalResult) + "  "
		result += fmt.Sprintf("%s vs %s", strings.Join(rollsStr1, ", "), strings.Join(rollsStr2, ", ")) + " "

		formula := fmt.Sprintf("(%s)", diceExpr)
		if modifier != 0 {
			formula += fmt.Sprintf(" %+d", modifier)
		}
		result += formula + " " + resultType

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

		// Check for critical hit (natural 20 on d20)
		isCrit := currentCritEnabled && numDice == 1 && sides == 20 && rolls[0] == 20

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

		// Format: "TOTAL  rolls (formula) +mod" or "20  (1d20) CRIT" for natural 20
		var result string

		if isCrit {
			// Critical hit format: "20  (1d20) CRIT" - both 20 and CRIT in red
			result = fmt.Sprintf("%d", finalTotal) + "  "
			formula := fmt.Sprintf("(%s)", diceExpr)
			if modifier != 0 {
				formula += fmt.Sprintf(" %+d", modifier)
			}
			result += formula + " " + getCriticalHitBanner()
		} else {
			result = fmt.Sprintf("%d", finalTotal) + "  "
			result += strings.Join(rollsStr, ", ") + " "

			formula := fmt.Sprintf("(%s)", diceExpr)
			if modifier != 0 {
				formula += fmt.Sprintf(" %+d", modifier)
			}
			result += formula
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
	results, resultStrings, hasCrit, err := evaluateExpressions(expressions, advantage, disadvantage)
	if err != nil {
		return err.Error()
	}

	// Calculate total
	total := calculateTotal(results, operators)

	// Format each expression on its own line
	resultStr := ""

	// Build formula string
	formula := expressions[0]
	for i := 1; i < len(expressions); i++ {
		if i-1 < len(operators) {
			formula += " " + operators[i-1] + " " + expressions[i]
		}
	}

	// Single line with total
	resultStr += fmt.Sprintf("%d", total) + "  "

	// Show individual results
	breakdown := resultStrings[0]
	for i := 1; i < len(resultStrings); i++ {
		if i-1 < len(operators) {
			breakdown += " " + operators[i-1] + " " + resultStrings[i]
		}
	}
	resultStr += breakdown + " "
	resultStr += "(" + formula + ")"

	if hasCrit {
		resultStr += " " + getCriticalHitBanner()
	}
	if advantage {
		resultStr += " ADV"
	} else if disadvantage {
		resultStr += " DIS"
	}

	return resultStr
}

// rollSingleDiceExpression rolls a single dice expression and returns the total
// Returns -1 for invalid expressions
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

		// Call RollDice recursively for each roll to get proper formatting
		result := RollDice(roll, &config.Config{
			DiceRoller: config.DiceRollerConfig{
				MinimumValue:       currentMinValue,
				ShowIndividual:     currentShowIndividual,
				CriticalHitEnabled: currentCritEnabled,
				CriticalHitMode:    currentCritMode,
			},
		})
		results = append(results, result)
	}

	return strings.Join(results, "\n")
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
		results, resultStrings, _, err := evaluateExpressions(expressions, false, false)
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
