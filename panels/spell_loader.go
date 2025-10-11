// panels/spell_loader.go
package panels

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sahilm/fuzzy"
)

// Spell represents a D&D spell (duplicate from ui package for panels use)
type Spell struct {
	Name           string   `json:"name"`
	Level          int      `json:"level"`
	School         string   `json:"school"`
	Classes        []string `json:"classes"`
	ActionType     string   `json:"actionType"`
	Concentration  bool     `json:"concentration"`
	Ritual         bool     `json:"ritual"`
	Range          string   `json:"range"`
	Components     []string `json:"components"`
	Material       string   `json:"material,omitempty"`
	Duration       string   `json:"duration"`
	Description    string   `json:"description"`
	CantripUpgrade string   `json:"cantripUpgrade,omitempty"`
}

var spellDatabase []Spell
var spellNames []string

// LoadSpells loads spells from the JSON file
func LoadSpells() error {
	if len(spellDatabase) > 0 {
		return nil // Already loaded
	}

	data, err := os.ReadFile("assets/spell.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &spellDatabase)
	if err != nil {
		return err
	}

	// Create spell names list for autocomplete
	spellNames = make([]string, len(spellDatabase))
	for i, spell := range spellDatabase {
		spellNames[i] = spell.Name
	}

	return nil
}

// matchesLevelFilter checks if a spell's level matches the filter
// Filter formats: "0" (cantrips), "3" (exactly 3), "1-3" (range), "5+" (5 and above)
func matchesLevelFilter(spell *Spell, filter string) bool {
	if filter == "" {
		return true // No filter
	}
	
	filter = strings.TrimSpace(filter)
	spellLevel := spell.Level
	
	// Handle "X+" format (e.g., "5+")
	if strings.HasSuffix(filter, "+") {
		minLevel, err := strconv.Atoi(strings.TrimSuffix(filter, "+"))
		if err != nil {
			return false
		}
		return spellLevel >= minLevel
	}
	
	// Handle "X-Y" format (e.g., "1-3")
	if strings.Contains(filter, "-") {
		parts := strings.Split(filter, "-")
		if len(parts) == 2 {
			minLevel, err1 := strconv.Atoi(parts[0])
			maxLevel, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				return spellLevel >= minLevel && spellLevel <= maxLevel
			}
		}
	}
	
	// Handle single value (e.g., "0" for cantrips, "3" for level 3)
	targetLevel, err := strconv.Atoi(filter)
	if err != nil {
		return false
	}
	return spellLevel == targetLevel
}

// SearchSpells returns spells matching the search term and level filter
func SearchSpells(searchTerm string, levelFilter string) []string {
	if err := LoadSpells(); err != nil {
		return []string{"Error loading spells"}
	}

	if searchTerm == "" && levelFilter == "" {
		return []string{}
	}

	// Build list of all spells for filtering
	allSpells := make([]*Spell, 0, len(spellDatabase))
	for i := range spellDatabase {
		allSpells = append(allSpells, &spellDatabase[i])
	}
	
	// Filter by level first if specified
	filteredSpells := allSpells
	if levelFilter != "" {
		filteredSpells = make([]*Spell, 0)
		for _, spell := range allSpells {
			if matchesLevelFilter(spell, levelFilter) {
				filteredSpells = append(filteredSpells, spell)
			}
		}
	}

	// If no search term, return filtered spells (limited to 50)
	if searchTerm == "" {
		matches := make([]string, 0, 50)
		for i, spell := range filteredSpells {
			if i >= 50 {
				break
			}
			matches = append(matches, spell.Name)
		}
		return matches
	}

	// Build list of filtered spell names for fuzzy search
	filteredNames := make([]string, 0, len(filteredSpells))
	for _, spell := range filteredSpells {
		filteredNames = append(filteredNames, spell.Name)
	}

	// Use fuzzy search to find matches
	results := fuzzy.Find(searchTerm, filteredNames)

	// Extract matched names (already sorted by score)
	matches := make([]string, 0, 50)
	for i, result := range results {
		if i >= 50 { // Limit to 50 suggestions for better browsing
			break
		}
		matches = append(matches, result.Str)
	}

	return matches
}

// FindSpell returns a spell by exact name match
func FindSpell(name string) *Spell {
	if err := LoadSpells(); err != nil {
		return nil
	}

	for _, spell := range spellDatabase {
		if strings.EqualFold(spell.Name, name) {
			return &spell
		}
	}
	return nil
}

// FormatSpell formats a spell for display
func FormatSpell(spell *Spell) string {
	if spell == nil {
		return "No spell selected"
	}

	var result strings.Builder

	// Header
	result.WriteString("📜 " + spell.Name + "\n")
	result.WriteString(strings.Repeat("─", len(spell.Name)+3) + "\n\n")

	// Basic info
	levelText := "Cantrip"
	if spell.Level > 0 {
		levelText = ordinal(spell.Level) + " level"
	}
	result.WriteString("Level: " + levelText + " " + spell.School + "\n")

	// Classes
	result.WriteString("Classes: " + strings.Join(spell.Classes, ", ") + "\n")

	// Casting details
	result.WriteString("Casting Time: " + spell.ActionType + "\n")
	result.WriteString("Range: " + spell.Range + "\n")
	result.WriteString("Duration: " + spell.Duration)
	if spell.Concentration {
		result.WriteString(" (Concentration)")
	}
	result.WriteString("\n")

	// Components
	components := make([]string, len(spell.Components))
	for i, comp := range spell.Components {
		switch comp {
		case "v":
			components[i] = "Verbal"
		case "s":
			components[i] = "Somatic"
		case "m":
			components[i] = "Material"
		default:
			components[i] = comp
		}
	}
	result.WriteString("Components: " + strings.Join(components, ", "))
	if spell.Material != "" {
		result.WriteString(" (" + spell.Material + ")")
	}
	result.WriteString("\n")

	if spell.Ritual {
		result.WriteString("Ritual: Yes\n")
	}

	result.WriteString("\n")

	// Description
	result.WriteString("Description:\n")
	wrappedDesc := wrapSpellText(spell.Description, 60)
	for _, line := range wrappedDesc {
		result.WriteString(line + "\n")
	}

	// Cantrip upgrade
	if spell.CantripUpgrade != "" {
		result.WriteString("\nAt Higher Levels:\n")
		wrappedUpgrade := wrapSpellText(spell.CantripUpgrade, 60)
		for _, line := range wrappedUpgrade {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

// ordinal converts a number to its ordinal form (1st, 2nd, 3rd, etc.)
func ordinal(n int) string {
	suffix := "th"
	switch n % 10 {
	case 1:
		if n%100 != 11 {
			suffix = "st"
		}
	case 2:
		if n%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if n%100 != 13 {
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suffix)
}

// wrapSpellText wraps text to specified line length
func wrapSpellText(text string, maxWidth int) []string {
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
