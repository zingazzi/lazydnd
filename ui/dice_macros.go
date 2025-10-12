// ui/dice_macros.go
package ui

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"lazydnd/panels"
)

// List of D&D 5e skills for skill check shortcuts
var dnd5eSkills = []string{
	"acrobatics", "animal", "arcana", "athletics", "deception",
	"history", "insight", "intimidation", "investigation",
	"medicine", "nature", "perception", "performance", "persuasion",
	"religion", "sleight", "stealth", "survival",
}

// processDiceInput processes dice input, handling macros, skill checks, and group rolls
func (m *Model) processDiceInput(input string) (string, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	// Check if creating a macro (contains '=')
	if strings.Contains(input, "=") {
		return m.createMacro(input)
	}

	// Check for group initiative command
	if input == "group" || input == "group init" || input == "group initiative" {
		return m.rollGroupInitiative()
	}

	// Check if input is a macro name
	if formula, exists := m.DiceMacros[input]; exists {
		result := panels.RollDice(formula, m.Config)
		// Store the expanded formula for reroll, not the macro name
		m.addToHistory(fmt.Sprintf("%s: %s", input, result), formula)
		m.LastDiceCommand = formula
		return fmt.Sprintf("Macro '%s' (%s): %s", input, formula, result), true
	}

	// Check if input is a skill check
	if m.isSkillCheck(input) {
		return m.rollSkillCheck(input)
	}

	// Otherwise, treat as normal dice expression
	result := panels.RollDice(input, m.Config)
	m.addToHistory(result, input)
	return result, true
}

// createMacro creates a new dice macro from input like "fireball=8d6"
func (m *Model) createMacro(input string) (string, bool) {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return "Invalid macro format. Use: name=formula (e.g., fireball=8d6)", false
	}

	name := strings.TrimSpace(parts[0])
	formula := strings.TrimSpace(parts[1])

	if name == "" || formula == "" {
		return "Macro name and formula cannot be empty", false
	}

	// Validate formula by trying to parse it (using a default config)
	// We can't validate fully without config, so we'll just check basic syntax
	if !strings.Contains(formula, "d") && !strings.Contains(formula, "D") {
		return fmt.Sprintf("Invalid dice formula '%s': must contain 'd' (e.g., 2d6)", formula), false
	}

	m.DiceMacros[name] = formula
	return fmt.Sprintf("✓ Macro '%s' saved: %s", name, formula), true
}

// isSkillCheck checks if input is a skill name
func (m *Model) isSkillCheck(input string) bool {
	for _, skill := range dnd5eSkills {
		if strings.HasPrefix(input, skill) {
			return true
		}
	}
	return false
}

// rollSkillCheck rolls a skill check for the selected initiative entry
func (m *Model) rollSkillCheck(skillInput string) (string, bool) {
	// Extract skill name
	skill := ""
	for _, s := range dnd5eSkills {
		if strings.HasPrefix(skillInput, s) {
			skill = s
			break
		}
	}

	// Check if there's a selected entry in initiative tracker
	if len(m.InitiativeList) == 0 || m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
		return fmt.Sprintf("No character selected. Select a character in Initiative Tracker to roll %s", skill), false
	}

	entry := m.InitiativeList[m.SelectedEntry]

	// Roll 1d20
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(20) + 1

	// Try to get modifier from monster data or parse from input
	modifier := 0
	modifierStr := strings.TrimPrefix(skillInput, skill)
	modifierStr = strings.TrimSpace(modifierStr)

	if modifierStr != "" && (modifierStr[0] == '+' || modifierStr[0] == '-') {
		if mod, err := strconv.Atoi(modifierStr); err == nil {
			modifier = mod
		}
	}

	total := roll + modifier

	result := fmt.Sprintf("%s - %s: 1d20%+d = %d%+d = %d",
		entry.Name,
		strings.Title(skill),
		modifier,
		roll,
		modifier,
		total)

	// Store the dice expression for reroll
	diceExpr := "1d20"
	if modifier != 0 {
		diceExpr = fmt.Sprintf("1d20%+d", modifier)
	}
	m.addToHistory(result, diceExpr)
	m.LastDiceCommand = diceExpr
	return result, true
}

// rollGroupInitiative rolls initiative for all monsters in the initiative tracker
func (m *Model) rollGroupInitiative() (string, bool) {
	if len(m.InitiativeList) == 0 {
		return "No entries in Initiative Tracker to roll for", false
	}

	rand.Seed(time.Now().UnixNano())
	rolled := 0
	results := []string{}

	for i := range m.InitiativeList {
		// Only roll for monsters
		if m.InitiativeList[i].Type != "monster" {
			continue
		}

		// Get DEX modifier if available
		dexMod := 0
		if m.InitiativeList[i].MonsterData != nil {
			dexModStr := m.InitiativeList[i].MonsterData.DEXMod
			if mod, err := strconv.Atoi(strings.TrimPrefix(dexModStr, "+")); err == nil {
				dexMod = mod
			}
		}

		// Roll 1d20 + DEX
		roll := rand.Intn(20) + 1
		total := roll + dexMod

		// Update initiative
		m.InitiativeList[i].Initiative = total

		results = append(results, fmt.Sprintf("%s: %d", m.InitiativeList[i].Name, total))
		rolled++
	}

	// Re-sort initiative list
	sort.SliceStable(m.InitiativeList, func(i, j int) bool {
		return m.InitiativeList[i].Initiative > m.InitiativeList[j].Initiative
	})

	// Reset selected entry
	if m.SelectedEntry >= len(m.InitiativeList) {
		m.SelectedEntry = 0
	}

	if rolled == 0 {
		return "No monsters found to roll initiative for", false
	}

	resultStr := fmt.Sprintf("Rolled initiative for %d monster(s):\n%s",
		rolled,
		strings.Join(results, ", "))

	// Don't set LastDiceCommand for group init (not rerollable)
	m.addToHistory(resultStr, "")
	return resultStr, true
}

// getMacroList returns a sorted list of macro names
func (m *Model) getMacroList() []string {
	macros := make([]string, 0, len(m.DiceMacros))
	for name := range m.DiceMacros {
		macros = append(macros, name)
	}
	sort.Strings(macros)
	return macros
}

// deleteMacro removes a macro by name
func (m *Model) deleteMacro(name string) bool {
	if _, exists := m.DiceMacros[name]; exists {
		delete(m.DiceMacros, name)
		return true
	}
	return false
}
