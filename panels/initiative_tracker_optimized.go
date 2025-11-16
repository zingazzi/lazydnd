// panels/initiative_tracker_optimized.go
package panels

import (
	"fmt"
	"strconv"
	"strings"
)

// InitiativeRenderConfig holds configuration for rendering initiative tracker
type InitiativeRenderConfig struct {
	ViewportHeight  int  // Maximum number of entries to render
	ScrollOffset    int  // Starting entry index for viewport
	EnableWindowing bool // Enable viewport-based rendering
}

// DefaultRenderConfig returns default rendering configuration
func DefaultRenderConfig() InitiativeRenderConfig {
	return InitiativeRenderConfig{
		ViewportHeight:  20, // Default viewport height
		ScrollOffset:    0,
		EnableWindowing: true, // Enable by default for large lists
	}
}

// GetOptimizedInitiativeContent returns optimized initiative tracker content
// Uses viewport-based rendering for large lists (20+ entries)
func GetOptimizedInitiativeContent(
	initiativeList interface{},
	input string,
	inputMode bool,
	inputType string,
	selectedEntry int,
	isActive bool,
	listMode bool,
	editMode bool,
	editType string,
	currentTurn int,
	roundCounter int,
	multiTargetMode bool,
	selectedTargets map[int]bool,
	config InitiativeRenderConfig,
	healthyColor, mediumColor, lowColor, criticalColor, tempHPColor string,
) string {
	var content strings.Builder
	content.Grow(2048) // Pre-allocate reasonable buffer

	// Header section
	renderHeader(&content, multiTargetMode, selectedTargets, editMode, editType, listMode)

	// Round counter section
	if roundCounter > 0 {
		renderRoundCounterOptimized(&content, roundCounter)
	}

	// Separator
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", 35))
	content.WriteString("\n\n")

	// Input/edit field
	if inputMode || editMode {
		renderInputFieldOptimized(&content, inputType, editType, editMode, input, isActive)
		content.WriteString("\n")
	}

	// Parse and render initiative list
	if initiativeList != nil {
		listStr := fmt.Sprintf("%+v", initiativeList)
		if listStr != "[]" && listStr != "<nil>" {
		renderInitiativeListOptimized(&content, listStr, selectedEntry, listMode, currentTurn,
			multiTargetMode, selectedTargets, config, healthyColor, mediumColor, lowColor, criticalColor, tempHPColor)
		} else {
			renderEmptyList(&content)
		}
	} else {
		renderEmptyList(&content)
	}

	return content.String()
}

// renderHeader renders the header section with mode indicators
func renderHeader(content *strings.Builder, multiTargetMode bool, selectedTargets map[int]bool, editMode bool, editType string, listMode bool) {
	if multiTargetMode {
		selectedCount := len(selectedTargets)
		content.WriteString(fmt.Sprintf("🎯 MULTI-TARGET MODE - %d target(s) selected\n", selectedCount))
		content.WriteString("Space: select/deselect • Enter: apply damage/healing • t: exit")
	} else if editMode {
		content.WriteString("EDIT MODE\n")
		switch editType {
		case "initiative":
			content.WriteString("Enter new initiative value:")
		case "hp":
			content.WriteString("Enter HP change (+heal/-damage):")
		case "delete":
			content.WriteString("Press Enter to confirm deletion")
		}
	} else if listMode {
		content.WriteString("LIST MODE - Use ↑↓ to select, i=initiative, h=HP, d=delete, t=multi-target")
	} else {
		content.WriteString("Press 'p' to add player, 'm' to add monster, Enter to edit\n")
		content.WriteString("Press 'n' for next turn, 'x' to reset combat")
	}
}

// renderRoundCounterOptimized renders the round counter with time elapsed (optimized version)
func renderRoundCounterOptimized(content *strings.Builder, roundCounter int) {
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

	content.WriteString("\n\n⚔️  Round ")
	content.WriteString(strconv.Itoa(roundCounter))
	content.WriteString(" / ")
	content.WriteString(timeStr)
}

// renderInputFieldOptimized renders the input/edit field (optimized version)
func renderInputFieldOptimized(content *strings.Builder, inputType, editType string, editMode bool, input string, isActive bool) {
	var prompt string
	switch inputType {
	case "player_name":
		prompt = "Player Name: "
	case "player_initiative":
		prompt = "Player Initiative: "
	case "player_ac":
		prompt = "Player AC: "
	case "monster_name":
		prompt = "Monster Name: "
	case "monster_hp":
		prompt = "Monster HP: "
	case "monster_ac":
		prompt = "Monster AC: "
	case "monster_initiative":
		prompt = "Monster Initiative (or 'r' to roll): "
	default:
		prompt = "Input: "
	}

	// Override prompt for edit modes
	if editMode {
		switch editType {
		case "initiative":
			prompt = "New Initiative: "
		case "hp":
			prompt = "HP Change (+heal/-damage): "
		case "delete":
			prompt = "Press Enter to confirm deletion"
		}
	}

	content.WriteString(prompt)
	content.WriteString(input)
	if isActive {
		content.WriteString("█")
	}
}

// renderEmptyList renders the empty list message
func renderEmptyList(content *strings.Builder) {
	content.WriteString("No entries in initiative order\n\n")
	content.WriteString("Add players and monsters to begin combat!")
}

// ParsedInitiativeEntry represents a parsed initiative entry for optimized rendering
type ParsedInitiativeEntry struct {
	Name         string
	Type         string
	Initiative   int
	HP           int
	MaxHP        int
	AC           string
	ConditionStr string // Pre-computed condition icons
	RawEntry     string
}

// parseInitiativeEntries parses and sorts initiative entries efficiently
func parseInitiativeEntriesOptimized(listStr string) []ParsedInitiativeEntry {
	// Remove outer brackets
	listStr = strings.TrimPrefix(listStr, "[")
	listStr = strings.TrimSuffix(listStr, "]")

	if listStr == "" {
		return nil
	}

	// Extract entries using brace matching
	entries := extractEntries(listStr)
	if len(entries) == 0 {
		return nil
	}

	// Parse each entry
	parsed := make([]ParsedInitiativeEntry, 0, len(entries))
	for _, entry := range entries {
		if pe := parseSingleEntry(entry); pe != nil {
			parsed = append(parsed, *pe)
		}
	}

	// Sort by initiative (highest first) - inline bubble sort is fine for initiative lists
	for i := 0; i < len(parsed); i++ {
		for j := i + 1; j < len(parsed); j++ {
			if parsed[j].Initiative > parsed[i].Initiative {
				parsed[i], parsed[j] = parsed[j], parsed[i]
			}
		}
	}

	return parsed
}

// extractEntries extracts individual entry strings from the list
func extractEntries(listStr string) []string {
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

	return entries
}

// parseSingleEntry parses a single initiative entry
func parseSingleEntry(entry string) *ParsedInitiativeEntry {
	if !strings.Contains(entry, "Name:") {
		return nil
	}

	pe := &ParsedInitiativeEntry{
		Name:       extractField(entry, "Name:"),
		Type:       extractField(entry, "Type:"),
		AC:         extractField(entry, "AC:"),
		RawEntry:   entry,
	}

	// Parse initiative
	if initStr := extractField(entry, "Initiative:"); initStr != "" {
		pe.Initiative, _ = strconv.Atoi(initStr)
	}

	// Parse HP values
	if hpStr := extractField(entry, "HP:"); hpStr != "" {
		pe.HP, _ = strconv.Atoi(hpStr)
	}
	if maxHPStr := extractField(entry, "MaxHP:"); maxHPStr != "" {
		pe.MaxHP, _ = strconv.Atoi(maxHPStr)
	}

	// Pre-compute condition icons
	pe.ConditionStr = extractConditionIcons(entry)

	return pe
}

// extractConditionIcons extracts and formats condition icons from entry
func extractConditionIcons(entry string) string {
	if !strings.Contains(entry, "Conditions:[") {
		return ""
	}

	// Find condition names and convert to emojis
	var emojis []string

	// Simple approach: look for common condition names in the Conditions array
	conditionMap := map[string]string{
		"Poisoned":        "🤢",
		"Stunned":         "😵",
		"Paralyzed":       "😵",
		"Incapacitated":   "😵",
		"Unconscious":     "😵",
		"Frightened":      "😱",
		"Charmed":         "😍",
		"Invisible":       "👻",
		"Prone":           "🤕",
		"Grappled":        "🔗",
		"Restrained":      "🔗",
		"Blinded":         "🙈",
		"Deafened":        "🙉",
		"Petrified":       "🗿",
		"Exhausted":       "😫",
	}

	for condName, emoji := range conditionMap {
		if strings.Contains(entry, "Name:"+condName) {
			emojis = append(emojis, emoji)
		}
	}

	if len(emojis) > 0 {
		return " " + strings.Join(emojis, "")
	}

	return ""
}

// renderInitiativeList renders the initiative list with optional windowing
func renderInitiativeListOptimized(
	content *strings.Builder,
	listStr string,
	selectedEntry int,
	listMode bool,
	currentTurn int,
	multiTargetMode bool,
	selectedTargets map[int]bool,
	config InitiativeRenderConfig,
	healthyColor, mediumColor, lowColor, criticalColor, tempHPColor string,
) {
	content.WriteString("Initiative Order:\n\n")

	// Parse entries
	entries := parseInitiativeEntriesOptimized(listStr)
	if len(entries) == 0 {
		content.WriteString("No valid entries found")
		return
	}

	// Determine viewport window
	startIdx := 0
	endIdx := len(entries)

	if config.EnableWindowing && len(entries) > config.ViewportHeight {
		// Calculate viewport window based on scroll offset
		startIdx = config.ScrollOffset
		if startIdx < 0 {
			startIdx = 0
		}
		if startIdx >= len(entries) {
			startIdx = len(entries) - config.ViewportHeight
			if startIdx < 0 {
				startIdx = 0
			}
		}

		endIdx = startIdx + config.ViewportHeight
		if endIdx > len(entries) {
			endIdx = len(entries)
		}

		// Ensure selected entry is visible if needed
		if selectedEntry >= 0 && selectedEntry < len(entries) && listMode {
			if selectedEntry < startIdx {
				startIdx = selectedEntry
				endIdx = startIdx + config.ViewportHeight
				if endIdx > len(entries) {
					endIdx = len(entries)
				}
			} else if selectedEntry >= endIdx {
				endIdx = selectedEntry + 1
				startIdx = endIdx - config.ViewportHeight
				if startIdx < 0 {
					startIdx = 0
				}
			}
		}

		// Add scroll indicators
		if startIdx > 0 {
			content.WriteString(fmt.Sprintf("▲ (%d more above)\n\n", startIdx))
		}
	}

	// Render visible entries
	for i := startIdx; i < endIdx; i++ {
		renderEntry(content, entries[i], i, selectedEntry, listMode, currentTurn, multiTargetMode, selectedTargets, healthyColor, mediumColor, lowColor, criticalColor, tempHPColor)
	}

	// Bottom scroll indicator
	if config.EnableWindowing && endIdx < len(entries) {
		content.WriteString(fmt.Sprintf("\n▼ (%d more below)", len(entries)-endIdx))
	}
}

// renderEntry renders a single initiative entry
func renderEntry(
	content *strings.Builder,
	entry ParsedInitiativeEntry,
	index int,
	selectedEntry int,
	listMode bool,
	currentTurn int,
	multiTargetMode bool,
	selectedTargets map[int]bool,
	healthyColor, mediumColor, lowColor, criticalColor, tempHPColor string,
) {
	// Checkbox for multi-target mode
	if multiTargetMode {
		if selectedTargets[index] {
			content.WriteString("[✓] ")
		} else {
			content.WriteString("[ ] ")
		}
	}

	// Turn marker
	if currentTurn == index {
		content.WriteString("★ ")
	} else {
		content.WriteString("  ")
	}

	// Entry number and name
	content.WriteString(fmt.Sprintf("%2d. %s (Init: %d", index+1, entry.Name, entry.Initiative))

	// Type-specific info
	if entry.Type == "monster" && entry.MaxHP > 0 {
		// Add colored HP
		coloredHP := getColoredHP(entry.HP, entry.MaxHP, healthyColor, mediumColor, criticalColor)
		content.WriteString(", ")
		content.WriteString(coloredHP)
	}

	// AC
	if entry.AC != "" && entry.AC != "0" {
		content.WriteString(", AC: ")
		content.WriteString(entry.AC)
	}

	content.WriteString(")")

	// Conditions
	if entry.ConditionStr != "" {
		content.WriteString(entry.ConditionStr)
	}

	// Selection indicator
	if listMode && selectedEntry == index {
		// For selected entry, prepend with arrow (done after line is complete)
		line := content.String()
		lastNewline := strings.LastIndex(line, "\n")
		if lastNewline >= 0 {
			lineContent := line[lastNewline+1:]
			content.Reset()
			content.WriteString(line[:lastNewline+1])
			content.WriteString("► ")
			content.WriteString(lineContent)
		}
	}

	content.WriteString("\n")
}
