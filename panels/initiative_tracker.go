// panels/initiative_tracker.go
package panels

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// InitiativeEntry represents a player or monster in the initiative tracker
// (This type is defined in ui/types.go to avoid import cycles)

var (
	initiativeInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Margin(1, 0)

	playerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	monsterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	selectedEntryStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7D56F4")).
				Foreground(lipgloss.Color("#FAFAFA")).
				Bold(true)
)

// GetInitiativeTrackerContent returns the content for the initiative tracker panel
func GetInitiativeTrackerContent(initiativeList interface{}, input string, inputMode bool, inputType string, selectedEntry int, isActive bool, listMode bool, editMode bool, editType string) string {
	var contentLines []string

	// Show different instructions based on mode
	if editMode {
		contentLines = append(contentLines, "EDIT MODE")
		switch editType {
		case "initiative":
			contentLines = append(contentLines, "Enter new initiative value:")
		case "hp":
			contentLines = append(contentLines, "Enter HP change (+heal/-damage):")
		case "delete":
			contentLines = append(contentLines, "Press Enter to confirm deletion")
		}
	} else if listMode {
		contentLines = append(contentLines, "LIST MODE - Use ↑↓ to select, i=initiative, h=HP, d=delete")
	} else {
		contentLines = append(contentLines, "Dungeon Master Panel")
		contentLines = append(contentLines, "Press 'p' to add player, 'm' to add monster, 'e' to edit")
	}
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, strings.Repeat("─", 40))
	contentLines = append(contentLines, "")

	// Input field (when adding entries or editing)
	if inputMode || editMode {
		var prompt string
		switch inputType {
		case "player_name":
			prompt = "Player Name: "
		case "player_initiative":
			prompt = "Player Initiative: "
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

		if isActive {
			prompt += input + "█"
		} else {
			prompt += input
		}
		contentLines = append(contentLines, initiativeInputStyle.Render(prompt))
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

				// Parse entries into a sortable structure
				type ParsedEntry struct {
					Name       string
					Type       string
					Initiative int
					HP         string
					MaxHP      string
					AC         string
					RawEntry   string
				}

				var parsedEntries []ParsedEntry

				// Parse each entry
				for _, entry := range entries {
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
						ac := extractField(entry, "AC:")

						parsedEntries = append(parsedEntries, ParsedEntry{
							Name:       name,
							Type:       entryType,
							Initiative: initiative,
							HP:         hp,
							MaxHP:      maxHP,
							AC:         ac,
							RawEntry:   entry,
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
					if entry.Type == "player" {
						line = fmt.Sprintf("%2d. %s (Initiative: %d)", i+1, entry.Name, entry.Initiative)
						if listMode && selectedEntry == i {
							line = selectedEntryStyle.Render("► " + line)
						} else {
							line = playerStyle.Render(line)
						}
					} else if entry.Type == "monster" {
						line = fmt.Sprintf("%2d. %s (Init: %d, HP: %s/%s, AC: %s)", i+1, entry.Name, entry.Initiative, entry.HP, entry.MaxHP, entry.AC)
						if listMode && selectedEntry == i {
							line = selectedEntryStyle.Render("► " + line)
						} else {
							line = monsterStyle.Render(line)
						}
					} else {
						line = fmt.Sprintf("%2d. %s (Initiative: %d)", i+1, entry.Name, entry.Initiative)
						if listMode && selectedEntry == i {
							line = selectedEntryStyle.Render("► " + line)
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
		nextFieldStart := -1
		fields := []string{" Type:", " Initiative:", " HP:", " MaxHP:", " AC:", "}"}

		for _, field := range fields {
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

// RollInitiative rolls a d20 for initiative
func RollInitiative() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(20) + 1
}

// ParseInput parses user input for different input types
func ParseInput(input string, inputType string) (interface{}, error) {
	switch inputType {
	case "player_name", "monster_name":
		if strings.TrimSpace(input) == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		return strings.TrimSpace(input), nil

	case "player_initiative", "monster_hp", "monster_ac":
		val, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 {
			return nil, fmt.Errorf("must be positive")
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
		if val < 0 {
			return nil, fmt.Errorf("must be positive")
		}
		return val, nil

	case "initiative":
		// Edit initiative value
		val, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, fmt.Errorf("must be a number")
		}
		if val < 0 {
			return nil, fmt.Errorf("must be positive")
		}
		return val, nil

	case "hp_change":
		// HP change (can be positive or negative)
		input = strings.TrimSpace(input)
		if input == "" {
			return nil, fmt.Errorf("enter a number (+ to heal, - to damage)")
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("must be a number (+ to heal, - to damage)")
		}
		return val, nil
	}

	return nil, fmt.Errorf("unknown input type")
}
