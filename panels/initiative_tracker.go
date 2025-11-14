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
func getColoredHP(hp, maxHP int) string {
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

	// Format HP text (styling handled by TView)
	return fmt.Sprintf("HP: %d/%d", hp, maxHP)
}

// getColoredHPWithTemp returns HP text with temp HP in cyan
// Returns a styled string like "HP: 7/10 +5" with appropriate colors
func getColoredHPWithTemp(hp, maxHP, tempHP int) string {
	if maxHP <= 0 {
		return ""
	}

	// Get the colored HP part
	hpPart := getColoredHP(hp, maxHP)

	// Add temp HP if present
	if tempHP > 0 {
		tempPart := fmt.Sprintf(" +%d", tempHP)
		return hpPart + tempPart
	}

	return hpPart
}

// GetInitiativeTrackerContent returns the content for the initiative tracker panel
func GetInitiativeTrackerContent(initiativeList interface{}, input string, inputMode bool, inputType string, selectedEntry int, isActive bool, listMode bool, editMode bool, editType string, currentTurn int, roundCounter int, multiTargetMode bool, selectedTargets map[int]bool, showRoundCounter bool) string {
	var contentLines []string

	// Removed instruction text for cleaner interface

	// Show round counter and elapsed time if combat has started (if enabled in config)
	if showRoundCounter && roundCounter > 0 {
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
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, roundInfo)
	}


	// Input field (when adding entries or editing)
	if inputMode || editMode {
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
			case "maxhp":
				prompt = "Max HP: "
			case "ac":
				prompt = "AC: "
			case "temphp":
				prompt = "Temporary HP: "
			case "delete":
				prompt = "Press Enter to confirm deletion"
			}
		}

		if isActive {
			prompt += input + "█"
		} else {
			prompt += input
		}
		contentLines = append(contentLines, prompt)
		contentLines = append(contentLines, "")
	}

	// Display initiative list
	if initiativeList != nil {
		// Try to convert interface{} to string and parse it
		listStr := fmt.Sprintf("%+v", initiativeList)

		if listStr != "[]" && listStr != "<nil>" {
			contentLines = append(contentLines, "Initiative Order:")
			contentLines = append(contentLines, "")

			// Better parsing approach - handle the slice format properly
			// The format is like: [{Name:Player1 Type:player Initiative:15 HP:0 MaxHP:0 AC:0} {Name:Monster1 Type:monster Initiative:12 HP:25 MaxHP:25 AC:14}]

			// Remove the outer brackets and split by individual entries
			listStr = strings.TrimPrefix(listStr, "[")
			listStr = strings.TrimSuffix(listStr, "]")

			if listStr != "" {
				// Split entries - each entry is wrapped in {}
				// Count braces carefully to handle nested structures like Conditions:[{...}]
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

				// Import ui package type for Conditions
				type Condition struct {
					Name        string
					RoundsLeft  int
					TotalRounds int
					Description string
				}

				// Parse entries into a sortable structure
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

				var parsedEntries []ParsedEntry

				// Parse each entry
				for entryIdx, entry := range entries {
					if strings.Contains(entry, "Name:") {
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

					// DEBUG: Print what we're about to parse
					_ = entryIdx // Use the variable to avoid unused error

					// Extract conditions if present
						var conditions []Condition
						if strings.Contains(entry, "Conditions:[") {
							// Find the Conditions array in the string
							condStart := strings.Index(entry, "Conditions:[")
							if condStart != -1 {
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

								if endIdx > 0 {
									// Extract the content between [ and ]
									condArrayStr := condStr[1:endIdx]
									if condArrayStr != "" {
										// Split by condition entries {Name:... RoundsLeft:... TotalRounds:... Description:...}
										condDepth := 0
										condStart := 0
										for i, char := range condArrayStr {
											if char == '{' {
												if condDepth == 0 {
													condStart = i
												}
												condDepth++
											} else if char == '}' {
												condDepth--
												if condDepth == 0 {
													condEntryStr := condArrayStr[condStart : i+1]
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
									}
								}
							}
						}

					parsedEntries = append(parsedEntries, ParsedEntry{
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
					})
					}
				}

				// Sort by initiative (highest first)
				for i := 0; i < len(parsedEntries); i++ {
					for j := i + 1; j < len(parsedEntries); j++ {
						if parsedEntries[j].Initiative > parsedEntries[i].Initiative {
							parsedEntries[i], parsedEntries[j] = parsedEntries[j], parsedEntries[i]
						}
					}
				}

				// Display sorted entries
				for i, entry := range parsedEntries {
					var line string
					var turnMarker string
					var checkbox string
					var conditionIcons string

					// Build condition icons string (compact display)
					if len(entry.Conditions) > 0 {
						var emojis []string
						for _, cond := range entry.Conditions {
							var emoji string
							switch cond.Name {
							case "Poisoned":
								emoji = "🤢"
							case "Stunned", "Paralyzed", "Incapacitated", "Unconscious":
								emoji = "😵"
							case "Frightened":
								emoji = "😱"
							case "Charmed":
								emoji = "😍"
							case "Invisible":
								emoji = "👻"
							case "Prone":
								emoji = "🤕"
							case "Grappled", "Restrained":
								emoji = "🔗"
							case "Blinded":
								emoji = "🙈"
							case "Deafened":
								emoji = "🙉"
							case "Petrified":
								emoji = "🗿"
							case "Exhausted":
								emoji = "😫"
							default:
								emoji = "🔮"
							}
							emojis = append(emojis, emoji)
						}
						if len(emojis) > 0 {
							conditionIcons = " " + strings.Join(emojis, "")
						}
					}

					// Add checkbox for multi-target mode
					if multiTargetMode {
						if selectedTargets[i] {
							checkbox = "[✓] "
						} else {
							checkbox = "[ ] "
						}
					} else {
						checkbox = ""
					}

					// Add turn marker if this is the current turn
					if currentTurn == i {
						turnMarker = "★ "
					} else {
						turnMarker = "  "
					}

					if entry.Type == "player" {
						// Format player line with AC if available
						if entry.AC != "" && entry.AC != "0" {
							line = fmt.Sprintf("%s%s%2d. %s (Init: %d, AC: %s)%s", checkbox, turnMarker, i+1, entry.Name, entry.Initiative, entry.AC, conditionIcons)
						} else {
							line = fmt.Sprintf("%s%s%2d. %s (Init: %d)%s", checkbox, turnMarker, i+1, entry.Name, entry.Initiative, conditionIcons)
						}
						// Apply selection marker (styling handled by TView)
						if listMode && selectedEntry == i {
							line = "► " + line
						}
					} else if entry.Type == "monster" {
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
					coloredHP := getColoredHPWithTemp(hpInt, maxHPInt, tempHPInt)

					// Get reaction indicator
					reactionIcon := ""
					if entry.ReactionUsed == "true" {
						reactionIcon = " [✗]" // Reaction used
					} else if entry.ReactionUsed == "false" {
						reactionIcon = " [✓]" // Reaction available
					}

					// Get legendary action counter
					legendaryCounter := ""
					legendaryMax := 0
					legendaryUsed := 0
					if entry.LegendaryActionsMax != "" {
						if max, err := strconv.Atoi(entry.LegendaryActionsMax); err == nil && max > 0 {
							legendaryMax = max
							if entry.LegendaryActionsUsed != "" {
								if used, err := strconv.Atoi(entry.LegendaryActionsUsed); err == nil {
									legendaryUsed = used
								}
							}
							// Show counter at end of name: "[3/3]" format
							legendaryCounter = fmt.Sprintf(" [%d/%d]", legendaryMax-legendaryUsed, legendaryMax)
						}
					}

					// Format line with colored HP and legendary counter
					displayName := entry.Name + legendaryCounter
					line = fmt.Sprintf("%s%s%2d. %s (Init: %d, %s, AC: %s)%s%s",
							checkbox, turnMarker, i+1, displayName, entry.Initiative, coloredHP, entry.AC, reactionIcon, conditionIcons)

					// Apply selection marker (styling handled by TView)
					if listMode && selectedEntry == i {
						line = "► " + line
					}
					} else {
						// Get reaction indicator for players too
						reactionIcon := ""
						if entry.ReactionUsed == "true" {
							reactionIcon = " [✗]" // Reaction used
						} else if entry.ReactionUsed == "false" {
							reactionIcon = " [✓]" // Reaction available
						}

					line = fmt.Sprintf("%s%s%2d. %s (Initiative: %d)%s%s", checkbox, turnMarker, i+1, entry.Name, entry.Initiative, reactionIcon, conditionIcons)
					if listMode && selectedEntry == i {
						line = "► " + line
					}
					}

					contentLines = append(contentLines, line)
				}
			}
		} else {
			contentLines = append(contentLines, "No entries in initiative order")
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, "Add players and monsters to begin combat!")
		}
	} else {
		contentLines = append(contentLines, "No entries in initiative order")
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Add players and monsters to begin combat!")
	}

	return strings.Join(contentLines, "\n")
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
