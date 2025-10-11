// panels/spell_loader.go
package panels

import (
	"encoding/json"
	"fmt"
	"os"
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

// SearchSpells returns spells matching the search term
func SearchSpells(searchTerm string) []string {
	if err := LoadSpells(); err != nil {
		return []string{"Error loading spells"}
	}

	if searchTerm == "" {
		return []string{}
	}

	// Use fuzzy search to find matches
	results := fuzzy.Find(searchTerm, spellNames)

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
