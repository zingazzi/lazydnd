// panels/initiative_tracker.go
package panels

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// InitiativeEntry represents a player or monster in the initiative tracker
// (This type is defined in ui/types.go to avoid import cycles)

// Note: Styling is now handled by TView widgets - these functions return plain text

// getColoredHP returns HP text with color coding based on percentage
// Returns a styled string like "HP: 7/10" with appropriate color
// Thresholds: >50% = grey, ≤50% and >20% = orange, ≤20% = red
func getColoredHP(hp, maxHP int, healthyColor, mediumColor, criticalColor string) string {
	if maxHP <= 0 {
		return ""
	}

	// Calculate percentage
	percentage := float64(hp) / float64(maxHP) * 100
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	// Color code based on percentage (TView format: [color]text[white])
	hpText := fmt.Sprintf("HP: %d/%d", hp, maxHP)
	if percentage <= 20 {
		return criticalColor + hpText + "[white]" // Critical (≤20%) - red
	} else if percentage <= 50 {
		return mediumColor + hpText + "[white]" // Medium (≤50% and >20%) - orange
	}
	return healthyColor + hpText + "[white]" // Healthy (>50%) - grey
}

// getColoredHPWithTemp returns HP text with temp HP in cyan
// Returns a styled string like "HP: 7/10 +5" with appropriate colors
func getColoredHPWithTemp(hp, maxHP, tempHP int, healthyColor, mediumColor, criticalColor, tempHPColor string) string {
	if maxHP <= 0 {
		return ""
	}

	// Get the colored HP part
	hpPart := getColoredHP(hp, maxHP, healthyColor, mediumColor, criticalColor)

	// Add temp HP in configured color if present
	if tempHP > 0 {
		hpPart += " " + tempHPColor + "+" + fmt.Sprintf("%d", tempHP) + "[white]" // Temp HP color
	}

	return hpPart
}

// GetInitiativeTrackerContent returns the formatted content for the initiative tracker panel.
// It renders the round counter, input fields, and the sorted initiative list with all entries.
// The function delegates to smaller helper functions for each rendering component.
func GetInitiativeTrackerContent(initiativeList interface{}, input string, inputMode bool, inputType string, selectedEntry int, isActive bool, listMode bool, editMode bool, editType string, currentTurn int, roundCounter int, multiTargetMode bool, selectedTargets map[int]bool, showRoundCounter bool, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor, playerNameColor, textColor string) string {
	var contentLines []string

	// Render round counter if enabled
	if showRoundCounter && roundCounter > 0 {
		contentLines = append(contentLines, renderRoundCounter(roundCounter, textColor)...)
	}

	// Render input field if in input or edit mode
	if inputMode || editMode {
		contentLines = append(contentLines, renderInputField(input, inputMode, inputType, editMode, editType, isActive)...)
	}

	// Display initiative list
	if initiativeList != nil {
		listLines := renderInitiativeList(initiativeList, selectedEntry, isActive, listMode, currentTurn, multiTargetMode, selectedTargets, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor, playerNameColor, textColor)
		contentLines = append(contentLines, listLines...)
	} else {
		contentLines = append(contentLines, "No entries in initiative order")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Add players and monsters to begin combat!")
	}

	return strings.Join(contentLines, "\n")
}

// renderRoundCounter renders the round counter and elapsed time display
func renderRoundCounter(roundCounter int, textColor string) []string {
	totalSeconds := roundCounter * 6
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	var timeStr string
	if minutes > 0 {
		if seconds > 0 {
			timeStr = fmt.Sprintf("%d minute%s %d second%s", minutes, pluralize(minutes), seconds, pluralize(seconds))
		} else {
			timeStr = fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
		}
	} else {
		timeStr = fmt.Sprintf("%d second%s", seconds, pluralize(seconds))
	}

	roundInfo := fmt.Sprintf("⚔️  Round %d / %s", roundCounter, timeStr)
	return []string{"", textColor + roundInfo + "[white]"}
}

// renderInputField renders the input field prompt for adding or editing entries
func renderInputField(input string, inputMode bool, inputType string, editMode bool, editType string, isActive bool) []string {
	var prompt string

	if editMode {
		prompt = getEditPrompt(editType)
	} else {
		prompt = getInputPrompt(inputType)
	}

	if isActive {
		prompt += input + "█"
	} else {
		prompt += input
	}

	return []string{prompt, ""}
}

// getInputPrompt returns the prompt string for input mode
func getInputPrompt(inputType string) string {
	switch inputType {
	case "player_name":
		return "Player Name: "
	case "player_initiative":
		return "Player Initiative: "
	case "player_ac":
		return "Player AC: "
	case "monster_name":
		return "Monster Name: "
	case "monster_hp":
		return "Monster HP: "
	case "monster_ac":
		return "Monster AC: "
	case "monster_initiative":
		return "Monster Initiative (or 'r' to roll): "
	default:
		return "Input: "
	}
}

// getEditPrompt returns the prompt string for edit mode
func getEditPrompt(editType string) string {
	switch editType {
	case "initiative":
		return "New Initiative: "
	case "hp":
		return "HP Change (+heal/-damage): "
	case "maxhp":
		return "Max HP: "
	case "ac":
		return "AC: "
	case "temphp":
		return "Temporary HP: "
	case "delete":
		return "Press Enter to confirm deletion"
	default:
		return "Input: "
	}
}

// renderInitiativeList renders the initiative list with all entries
func renderInitiativeList(initiativeList interface{}, selectedEntry int, isActive bool, listMode bool, currentTurn int, multiTargetMode bool, selectedTargets map[int]bool, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor, playerNameColor, textColor string) []string {
	var contentLines []string

	// Try to convert interface{} to string and parse it
	listStr := fmt.Sprintf("%+v", initiativeList)

	if listStr == "[]" || listStr == "<nil>" {
		contentLines = append(contentLines, "No entries in initiative order")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Add players and monsters to begin combat!")
		return contentLines
	}

	contentLines = append(contentLines, textColor+"Initiative Order:[white]")
	contentLines = append(contentLines, "")

	// Parse entries from string representation
	parsedEntries := parseInitiativeEntries(listStr)
	if len(parsedEntries) == 0 {
		contentLines = append(contentLines, "No entries in initiative order")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Add players and monsters to begin combat!")
		return contentLines
	}

	// Sort entries by initiative
	sortEntriesByInitiative(parsedEntries)

	// Render each entry
	for i, entry := range parsedEntries {
		line := renderInitiativeEntry(entry, i, selectedEntry, listMode, currentTurn, multiTargetMode, selectedTargets, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor, playerNameColor)
		contentLines = append(contentLines, line)
	}

	return contentLines
}

// ParsedEntry represents a parsed initiative entry
type ParsedEntry struct {
	Name                string
	Type                string
	Initiative          int
	HP                  string
	MaxHP               string
	TempHP              string
	AC                  string
	ReactionUsed        string
	LegendaryActionsMax string
	LegendaryActionsUsed string
	RawEntry            string
	Conditions          []Condition
}

// Condition represents a parsed condition
type Condition struct {
	Name        string
	RoundsLeft  int
	TotalRounds int
	Description string
}

// parseInitiativeEntries parses the string representation of initiative list into ParsedEntry structs
func parseInitiativeEntries(listStr string) []ParsedEntry {
	// Remove the outer brackets
	listStr = strings.TrimPrefix(listStr, "[")
	listStr = strings.TrimSuffix(listStr, "]")

	if listStr == "" {
		return nil
	}

	// Split entries - each entry is wrapped in {}
	var entries []string
	depth := 0
	start := 0

	for i, char := range listStr {
		if char == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if char == '}' {
			depth--
			if depth == 0 {
				entries = append(entries, listStr[start:i+1])
			}
		}
	}

	// Parse each entry
	var parsedEntries []ParsedEntry
	for _, entry := range entries {
		if strings.Contains(entry, "Name:") {
			parsedEntry := parseInitiativeEntry(entry)
			parsedEntries = append(parsedEntries, parsedEntry)
		}
	}

	return parsedEntries
}

// parseInitiativeEntry parses a single initiative entry string into a ParsedEntry
func parseInitiativeEntry(entry string) ParsedEntry {
	name := extractField(entry, "Name:")
	entryType := extractField(entry, "Type:")
	initiativeStr := extractField(entry, "Initiative:")

	// Convert initiative to int for sorting
	initiative := 0
	if initiativeStr != "" {
		if val, err := strconv.Atoi(initiativeStr); err == nil {
			initiative = val
		}
	}

	hp := extractField(entry, "HP:")
	maxHP := extractField(entry, "MaxHP:")
	tempHP := extractField(entry, "TempHP:")
	ac := extractField(entry, "AC:")
	reactionUsed := extractField(entry, "ReactionUsed:")
	legendaryMax := extractField(entry, "LegendaryActionsMax:")
	legendaryUsed := extractField(entry, "LegendaryActionsUsed:")

	// Extract conditions if present
	conditions := parseConditions(entry)

	return ParsedEntry{
		Name:                name,
		Type:                entryType,
		Initiative:          initiative,
		HP:                  hp,
		MaxHP:               maxHP,
		TempHP:              tempHP,
		AC:                  ac,
		ReactionUsed:        reactionUsed,
		LegendaryActionsMax: legendaryMax,
		LegendaryActionsUsed: legendaryUsed,
		RawEntry:            entry,
		Conditions:          conditions,
	}
}

// parseConditions extracts conditions from an entry string
func parseConditions(entry string) []Condition {
	var conditions []Condition
	if !strings.Contains(entry, "Conditions:[") {
		return conditions
	}

	// Find the Conditions array in the string
	condStart := strings.Index(entry, "Conditions:[")
	if condStart == -1 {
		return conditions
	}

	// Find the matching closing bracket
	condStr := entry[condStart+len("Conditions:"):]
	depth := 0
	endIdx := 0
	for i, char := range condStr {
		if char == '[' {
			depth++
		} else if char == ']' {
			depth--
			if depth == 0 {
				endIdx = i
				break
			}
		}
	}

	if endIdx == 0 {
		return conditions
	}

	// Extract the content between [ and ]
	condArrayStr := condStr[1:endIdx]
	if condArrayStr == "" {
		return conditions
	}

	// Split by condition entries {Name:... RoundsLeft:... TotalRounds:... Description:...}
	condDepth := 0
	condStartIdx := 0
	for i, char := range condArrayStr {
		if char == '{' {
			if condDepth == 0 {
				condStartIdx = i
			}
			condDepth++
		} else if char == '}' {
			condDepth--
			if condDepth == 0 {
				condEntryStr := condArrayStr[condStartIdx : i+1]
				// Parse individual condition
				condName := extractField(condEntryStr, "Name:")
				roundsLeftStr := extractField(condEntryStr, "RoundsLeft:")
				totalRoundsStr := extractField(condEntryStr, "TotalRounds:")

				roundsLeft := 0
				totalRounds := 0
				if roundsLeftStr != "" {
					roundsLeft, _ = strconv.Atoi(roundsLeftStr)
				}
				if totalRoundsStr != "" {
					totalRounds, _ = strconv.Atoi(totalRoundsStr)
				}

				conditions = append(conditions, Condition{
					Name:        condName,
					RoundsLeft:  roundsLeft,
					TotalRounds: totalRounds,
				})
			}
		}
	}

	return conditions
}

// sortEntriesByInitiative sorts parsed entries by initiative (highest first)
func sortEntriesByInitiative(entries []ParsedEntry) {
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Initiative > entries[i].Initiative {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

// renderInitiativeEntry renders a single initiative entry as a formatted line
func renderInitiativeEntry(entry ParsedEntry, index int, selectedEntry int, listMode bool, currentTurn int, multiTargetMode bool, selectedTargets map[int]bool, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor, playerNameColor string) string {
	checkbox := getCheckbox(multiTargetMode, selectedTargets, index)
	turnMarker := getTurnMarker(currentTurn, index)
	conditionIcons := getConditionIcons(entry.Conditions)

	if entry.Type == "player" {
		return renderPlayerEntry(entry, index, selectedEntry, listMode, checkbox, turnMarker, conditionIcons, playerNameColor)
	} else if entry.Type == "monster" {
		return renderMonsterEntry(entry, index, selectedEntry, listMode, checkbox, turnMarker, conditionIcons, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor)
	}

	// Fallback for unknown types
	return renderPlayerEntry(entry, index, selectedEntry, listMode, checkbox, turnMarker, conditionIcons, playerNameColor)
}

// getCheckbox returns the checkbox string for multi-target mode
func getCheckbox(multiTargetMode bool, selectedTargets map[int]bool, index int) string {
	if multiTargetMode {
		if selectedTargets[index] {
			return "[✓] "
		}
		return "[ ] "
	}
	return ""
}

// getTurnMarker returns the turn marker string
func getTurnMarker(currentTurn int, index int) string {
	if currentTurn == index {
		return "★ "
	}
	return "  "
}

// getConditionIcons returns the condition icons string
func getConditionIcons(conditions []Condition) string {
	if len(conditions) == 0 {
		return ""
	}

	var emojis []string
	for _, cond := range conditions {
		emoji := getConditionEmoji(cond.Name)
		emojis = append(emojis, emoji)
	}

	if len(emojis) > 0 {
		return " " + strings.Join(emojis, "")
	}
	return ""
}

// getConditionEmoji returns the emoji for a condition name
func getConditionEmoji(condName string) string {
	switch condName {
	case "Poisoned":
		return "🤢"
	case "Stunned", "Paralyzed", "Incapacitated", "Unconscious":
		return "😵"
	case "Frightened":
		return "😱"
	case "Charmed":
		return "😍"
	case "Invisible":
		return "👻"
	case "Prone":
		return "🤕"
	case "Grappled", "Restrained":
		return "🔗"
	case "Blinded":
		return "🙈"
	case "Deafened":
		return "🙉"
	case "Petrified":
		return "🗿"
	case "Exhausted":
		return "😫"
	default:
		return "🔮"
	}
}

// renderPlayerEntry renders a player entry line
func renderPlayerEntry(entry ParsedEntry, index int, selectedEntry int, listMode bool, checkbox, turnMarker, conditionIcons, playerNameColor string) string {
	coloredName := "[green]" + entry.Name + "[white]"
	var line string

	if entry.AC != "" && entry.AC != "0" {
		line = fmt.Sprintf("%s%s%2d. %s (Init: %d, AC: %s)%s", checkbox, turnMarker, index+1, coloredName, entry.Initiative, entry.AC, conditionIcons)
	} else {
		line = fmt.Sprintf("%s%s%2d. %s (Init: %d)%s", checkbox, turnMarker, index+1, coloredName, entry.Initiative, conditionIcons)
	}

	// Apply selection marker
	if listMode && selectedEntry == index {
		line = "► " + line
	}

	return line
}

// renderMonsterEntry renders a monster entry line
func renderMonsterEntry(entry ParsedEntry, index int, selectedEntry int, listMode bool, checkbox, turnMarker, conditionIcons, healthyColor, mediumColor, criticalColor, tempHPColor, monsterNameColor string) string {
	// Convert HP strings to integers for color coding
	hpInt := 0
	maxHPInt := 0
	if entry.HP != "" {
		hpInt, _ = strconv.Atoi(entry.HP)
	}
	if entry.MaxHP != "" {
		maxHPInt, _ = strconv.Atoi(entry.MaxHP)
	}

	// Get colored HP text with temp HP
	tempHPInt := 0
	if entry.TempHP != "" {
		tempHPInt, _ = strconv.Atoi(entry.TempHP)
	}
	coloredHP := getColoredHPWithTemp(hpInt, maxHPInt, tempHPInt, healthyColor, mediumColor, criticalColor, tempHPColor)

	// Get reaction indicator
	reactionIcon := getReactionIcon(entry.ReactionUsed)

	// Get legendary action counter
	legendaryCounter := getLegendaryCounter(entry.LegendaryActionsMax, entry.LegendaryActionsUsed)

	// Format line with colored HP and legendary counter
	if monsterNameColor == "" {
		monsterNameColor = "[red]"
	}
	coloredName := monsterNameColor + entry.Name + "[white]" + legendaryCounter
	line := fmt.Sprintf("%s%s%2d. %s (Init: %d, %s, AC: %s)%s%s",
		checkbox, turnMarker, index+1, coloredName, entry.Initiative, coloredHP, entry.AC, reactionIcon, conditionIcons)

	// Apply selection marker
	if listMode && selectedEntry == index {
		line = "► " + line
	}

	return line
}

// getReactionIcon returns the reaction indicator string
func getReactionIcon(reactionUsed string) string {
	if reactionUsed == "true" {
		return " [✗]" // Reaction used
	} else if reactionUsed == "false" {
		return " [✓]" // Reaction available
	}
	return ""
}

// getLegendaryCounter returns the legendary action counter string
func getLegendaryCounter(legendaryMaxStr, legendaryUsedStr string) string {
	if legendaryMaxStr == "" {
		return ""
	}

	max, err := strconv.Atoi(legendaryMaxStr)
	if err != nil || max <= 0 {
		return ""
	}

	used := 0
	if legendaryUsedStr != "" {
		if u, err := strconv.Atoi(legendaryUsedStr); err == nil {
			used = u
		}
	}

	// Show counter at end of name: "[3/3]" format
	return fmt.Sprintf(" [%d/%d]", max-used, max)
}

// extractField extracts a field value from a struct string representation
func extractField(entry, fieldName string) string {
	start := strings.Index(entry, fieldName)
	if start == -1 {
		return ""
	}

	start += len(fieldName)
	end := start

	// For Name field, we need to handle spaces differently
	if fieldName == "Name:" {
		// Find the next field or end of struct
		// Check if this is an InitiativeEntry or a Condition by looking for different fields
		nextFieldStart := -1

		// Try initiative entry fields first
		initiativeFields := []string{" Type:", " Initiative:", " HP:", " MaxHP:", " AC:", " MonsterData:", " InstanceNum:", " BaseName:", " MonsterName:", " Conditions:"}
		// Try condition fields
		conditionFields := []string{" RoundsLeft:", " TotalRounds:", " Description:"}

		// Combine all possible next fields
		allFields := append(initiativeFields, conditionFields...)
		allFields = append(allFields, "}")

		for _, field := range allFields {
			if idx := strings.Index(entry[start:], field); idx != -1 {
				if nextFieldStart == -1 || idx < nextFieldStart {
					nextFieldStart = idx
				}
			}
		}

		if nextFieldStart != -1 {
			end = start + nextFieldStart
		} else {
			end = len(entry)
		}
	} else {
		// For other fields, stop at space or closing brace
		for end < len(entry) && entry[end] != ' ' && entry[end] != '}' {
			end++
		}
	}

	return strings.TrimSpace(entry[start:end])
}

// RollInitiative rolls a standard d20 die for initiative order.
// Returns a random integer between 1 and 20 (inclusive).
// Used when adding monsters or players to the initiative tracker.
func RollInitiative() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(20) + 1
}

// ParseInput parses and validates user input based on the specified input type.
//
// Supported input types:
//   - "player_name", "monster_name": Validates non-empty trimmed strings with length/character limits
//   - "player_initiative", "monster_initiative", "initiative": Validates initiative values (-10 to 99)
//   - "player_ac", "monster_ac": Validates armor class (0 to 99)
//   - "monster_hp": Validates hit points (0 to 9999)
//   - "hp_change": Parses signed integers ("+5" or "-10") for HP modifications
//   - "maxhp": Validates maximum HP (1 to 9999)
//   - "temphp": Validates temporary HP (0 to 9999)
//
// Returns the parsed value (string or int) or an error if validation fails.
// For "monster_initiative" with input "r", automatically rolls a d20.
func ParseInput(input string, inputType string) (interface{}, error) {
	switch inputType {
	case "player_name", "monster_name":
		name := strings.TrimSpace(input)
		if name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		if len(name) > 50 {
			return nil, fmt.Errorf("name too long (max 50 characters)")
		}
		// Allow letters, numbers, spaces, hyphens, apostrophes, underscores
		for _, char := range name {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == ' ' || char == '-' ||
				char == '\'' || char == '_') {
				return nil, fmt.Errorf("name contains invalid characters")
			}
		}
		return name, nil

	case "player_initiative", "initiative":
		val, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < -10 || val > 99 {
			return nil, fmt.Errorf("initiative must be -10 to 99")
		}
		return val, nil

	case "monster_initiative":
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "r" || input == "roll" {
			return RollInitiative(), nil
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("enter a number or 'r' to roll")
		}
		if val < -10 || val > 99 {
			return nil, fmt.Errorf("initiative must be -10 to 99")
		}
		return val, nil

	case "player_ac", "monster_ac":
		val, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 || val > 99 {
			return nil, fmt.Errorf("AC must be 0 to 99")
		}
		return val, nil

	case "monster_hp":
		val, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 {
			return nil, fmt.Errorf("HP cannot be negative")
		}
		if val > 9999 {
			return nil, fmt.Errorf("HP too high (max 9999)")
		}
		return val, nil

	case "hp_change":
		input = strings.TrimSpace(input)
		if input == "" {
			return nil, fmt.Errorf("enter a number (+ to heal, - to damage)")
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("must be a number (+ to heal, - to damage)")
		}
		if val < -9999 || val > 9999 {
			return nil, fmt.Errorf("value too extreme (max 9999)")
		}
		return val, nil

	case "maxhp":
		input = strings.TrimSpace(input)
		if input == "" {
			return nil, fmt.Errorf("enter a number")
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 1 {
			return nil, fmt.Errorf("must be at least 1")
		}
		if val > 9999 {
			return nil, fmt.Errorf("must be at most 9999")
		}
		return val, nil

	case "ac":
		input = strings.TrimSpace(input)
		if input == "" {
			return nil, fmt.Errorf("enter a number")
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 {
			return nil, fmt.Errorf("must be at least 0")
		}
		if val > 99 {
			return nil, fmt.Errorf("AC too high (max 99)")
		}
		return val, nil

	case "temphp":
		input = strings.TrimSpace(input)
		if input == "" {
			return nil, fmt.Errorf("enter a number (0 to clear)")
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 {
			return nil, fmt.Errorf("temp HP cannot be negative")
		}
		if val > 9999 {
			return nil, fmt.Errorf("temp HP too high (max 9999)")
		}
		return val, nil
	}

	return nil, fmt.Errorf("unknown input type")
}

// pluralize returns "s" if count is not 1, otherwise empty string
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
