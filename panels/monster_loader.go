// panels/monster_loader.go
package panels

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Monster represents a D&D monster (duplicate to avoid import cycles)
type Monster struct {
	Name             string `json:"name"`
	Meta             string `json:"meta"`
	ArmorClass       string `json:"Armor Class"`
	HitPoints        string `json:"Hit Points"`
	Speed            string `json:"Speed"`
	STR              string `json:"STR"`
	STRMod           string `json:"STR_mod"`
	DEX              string `json:"DEX"`
	DEXMod           string `json:"DEX_mod"`
	CON              string `json:"CON"`
	CONMod           string `json:"CON_mod"`
	INT              string `json:"INT"`
	INTMod           string `json:"INT_mod"`
	WIS              string `json:"WIS"`
	WISMod           string `json:"WIS_mod"`
	CHA              string `json:"CHA"`
	CHAMod           string `json:"CHA_mod"`
	SavingThrows     string `json:"Saving Throws,omitempty"`
	Skills           string `json:"Skills,omitempty"`
	Senses           string `json:"Senses,omitempty"`
	Languages        string `json:"Languages,omitempty"`
	Challenge        string `json:"Challenge"`
	Traits           string `json:"Traits,omitempty"`
	Actions          string `json:"Actions,omitempty"`
	LegendaryActions string `json:"Legendary Actions,omitempty"`
	ImgURL           string `json:"img_url,omitempty"`
}

var monsters []Monster

// LoadMonsters loads monsters from the JSON file
func LoadMonsters() error {
	if len(monsters) > 0 {
		return nil // Already loaded
	}

	data, err := os.ReadFile("assets/monsters.json")
	if err != nil {
		return fmt.Errorf("failed to read monsters.json: %w", err)
	}

	err = json.Unmarshal(data, &monsters)
	if err != nil {
		return fmt.Errorf("failed to parse monsters.json: %w", err)
	}

	return nil
}

// SearchMonsters returns monster names that match the search term
func SearchMonsters(searchTerm string) []string {
	if err := LoadMonsters(); err != nil {
		return []string{}
	}

	if searchTerm == "" {
		return []string{}
	}

	var matches []string
	searchLower := strings.ToLower(searchTerm)

	for _, monster := range monsters {
		if strings.Contains(strings.ToLower(monster.Name), searchLower) {
			matches = append(matches, monster.Name)
			if len(matches) >= 10 { // Limit suggestions
				break
			}
		}
	}

	return matches
}

// FindMonster returns a monster by name
func FindMonster(name string) *Monster {
	if err := LoadMonsters(); err != nil {
		return nil
	}

	for _, monster := range monsters {
		if strings.EqualFold(monster.Name, name) {
			return &monster
		}
	}

	return nil
}

// FormatMonster formats monster details for display
func FormatMonster(monster *Monster) string {
	if monster == nil {
		return "Monster not found"
	}

	var lines []string

	// Header
	lines = append(lines, fmt.Sprintf("🐲 %s", monster.Name))
	lines = append(lines, monster.Meta)
	lines = append(lines, "")

	// Basic Stats
	lines = append(lines, fmt.Sprintf("AC: %s", monster.ArmorClass))
	lines = append(lines, fmt.Sprintf("HP: %s", monster.HitPoints))
	lines = append(lines, fmt.Sprintf("Speed: %s", monster.Speed))
	lines = append(lines, "")

	// Ability Scores
	lines = append(lines, "ABILITY SCORES:")
	lines = append(lines, fmt.Sprintf("STR: %s %s  DEX: %s %s  CON: %s %s",
		monster.STR, monster.STRMod, monster.DEX, monster.DEXMod, monster.CON, monster.CONMod))
	lines = append(lines, fmt.Sprintf("INT: %s %s  WIS: %s %s  CHA: %s %s",
		monster.INT, monster.INTMod, monster.WIS, monster.WISMod, monster.CHA, monster.CHAMod))
	lines = append(lines, "")

	// Optional fields
	if monster.SavingThrows != "" {
		lines = append(lines, fmt.Sprintf("Saving Throws: %s", monster.SavingThrows))
	}
	if monster.Skills != "" {
		lines = append(lines, fmt.Sprintf("Skills: %s", monster.Skills))
	}
	if monster.Senses != "" {
		lines = append(lines, fmt.Sprintf("Senses: %s", monster.Senses))
	}
	if monster.Languages != "" {
		lines = append(lines, fmt.Sprintf("Languages: %s", monster.Languages))
	}
	lines = append(lines, fmt.Sprintf("Challenge: %s", monster.Challenge))
	lines = append(lines, "")

	// Traits (clean HTML)
	if monster.Traits != "" {
		lines = append(lines, "TRAITS:")
		cleanTraits := cleanHTML(monster.Traits)
		wrappedTraits := wrapText(cleanTraits, 60) // Wrap at 60 characters
		lines = append(lines, wrappedTraits...)
		lines = append(lines, "")
	}

	// Actions (clean HTML)
	if monster.Actions != "" {
		lines = append(lines, "ACTIONS:")
		cleanActions := cleanHTML(monster.Actions)
		wrappedActions := wrapText(cleanActions, 60) // Wrap at 60 characters
		lines = append(lines, wrappedActions...)
		lines = append(lines, "")
	}

	// Legendary Actions (clean HTML)
	if monster.LegendaryActions != "" {
		lines = append(lines, "LEGENDARY ACTIONS:")
		cleanLegendary := cleanHTML(monster.LegendaryActions)
		wrappedLegendary := wrapText(cleanLegendary, 60) // Wrap at 60 characters
		lines = append(lines, wrappedLegendary...)
	}

	return strings.Join(lines, "\n")
}

// cleanHTML removes HTML tags and formats text for terminal display
func cleanHTML(html string) string {
	// Simple HTML cleaning - remove tags and format
	text := html

	// Replace common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")

	// Remove HTML tags (simple approach)
	for strings.Contains(text, "<") && strings.Contains(text, ">") {
		start := strings.Index(text, "<")
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Clean up extra whitespace
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.TrimSpace(text)

	return text
}

// wrapText wraps text to specified line length
func wrapText(text string, maxWidth int) []string {
	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed the line length, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > maxWidth {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// ExtractMonsterStats extracts HP and AC from a monster for initiative tracker
func ExtractMonsterStats(monsterName string) (hp int, ac int, err error) {
	monster := FindMonster(monsterName)
	if monster == nil {
		return 0, 0, fmt.Errorf("monster not found: %s", monsterName)
	}

	// Parse HP from string like "58 (9d10 + 9)" - extract the first number
	hpStr := strings.TrimSpace(monster.HitPoints)
	if hpStr != "" {
		// Find the first number in the HP string
		var hpPart strings.Builder
		for _, char := range hpStr {
			if char >= '0' && char <= '9' {
				hpPart.WriteRune(char)
			} else if hpPart.Len() > 0 {
				break // Stop at first non-digit after we've started collecting digits
			}
		}
		if hpPart.Len() > 0 {
			if parsedHP, parseErr := strconv.Atoi(hpPart.String()); parseErr == nil {
				hp = parsedHP
			}
		}
	}

	// Parse AC from string like "16 (Natural Armor)" - extract the first number
	acStr := strings.TrimSpace(monster.ArmorClass)
	if acStr != "" {
		// Find the first number in the AC string
		var acPart strings.Builder
		for _, char := range acStr {
			if char >= '0' && char <= '9' {
				acPart.WriteRune(char)
			} else if acPart.Len() > 0 {
				break // Stop at first non-digit after we've started collecting digits
			}
		}
		if acPart.Len() > 0 {
			if parsedAC, parseErr := strconv.Atoi(acPart.String()); parseErr == nil {
				ac = parsedAC
			}
		}
	}

	if hp == 0 && ac == 0 {
		return 0, 0, fmt.Errorf("could not parse HP or AC for monster: %s", monsterName)
	}

	return hp, ac, nil
}
