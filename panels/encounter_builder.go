// panels/encounter_builder.go
package panels

import (
	"fmt"
	"lazydnd/encounters"
	"strings"
)

// GetEncounterBuilderContent generates the encounter builder panel content
func GetEncounterBuilderContent(
	mode string,
	partySize, partyLevel int,
	encounterMonsters []EncounterMonster,
	selectedIndex int,
	savedEncounters []Encounter,
	encounterListMode bool,
	selectedSaved int,
	crFilter string,
	width, height int,
	activeStyle, inactiveStyle, titleStyle interface{}, // Deprecated - kept for compatibility
) string {
	content := ""

	switch mode {
	case "party_setup":
		content = renderPartySetup(partySize, partyLevel, width, height)
	case "building":
		content = renderEncounterBuilder(partySize, partyLevel, encounterMonsters, selectedIndex, crFilter, width, height)
	case "templates":
		content = renderTemplates(savedEncounters, encounterListMode, selectedSaved, width, height)
	case "template_detail":
		content = renderTemplateDetail(savedEncounters, selectedSaved, partySize, partyLevel, width, height)
	default:
		content = renderPartySetup(partySize, partyLevel, width, height)
	}

	return content
}

// EncounterMonster matches the type in ui/types.go
type EncounterMonster struct {
	Name     string
	CR       string
	HP       int
	MaxHP    int
	AC       int
	Quantity int
	XP       int
}

// Encounter matches the type in ui/types.go
type Encounter struct {
	Name     string
	Monsters []EncounterMonster
}

func renderPartySetup(partySize, partyLevel, width, height int) string {
	var b strings.Builder

	b.WriteString("⚔️  Party Setup\n\n")
	b.WriteString(fmt.Sprintf("Party Size: %d players\n", partySize))
	b.WriteString(fmt.Sprintf("Party Level: %d\n", partyLevel))

	return b.String()
}

func renderEncounterBuilder(
	partySize, partyLevel int,
	encounterMonsters []EncounterMonster,
	selectedIndex int,
	crFilter string,
	width, height int,
) string {
	var b strings.Builder

	// Header with party info
	b.WriteString("⚔️  Encounter Builder\n\n")
	b.WriteString(fmt.Sprintf("Party: %d × Level %d\n\n", partySize, partyLevel))

	// Calculate difficulty
	if partySize > 0 && partyLevel > 0 {
		partyLevels := make([]int, partySize)
		for i := 0; i < partySize; i++ {
			partyLevels[i] = partyLevel
		}

		monsterCRs := []string{}
		for _, m := range encounterMonsters {
			for i := 0; i < m.Quantity; i++ {
				monsterCRs = append(monsterCRs, m.CR)
			}
		}

		analysis := encounters.CalculateDifficulty(partyLevels, monsterCRs)
		diffColor := encounters.GetDifficultyColor(analysis.Difficulty)

		// Use ANSI color codes
		colorCode := ""
		switch diffColor {
		case "#AAAAAA":
			colorCode = "\033[90m" // Bright black (gray)
		case "#00FF00":
			colorCode = "\033[92m" // Bright green
		case "#FFD700":
			colorCode = "\033[93m" // Bright yellow
		case "#FFA500":
			colorCode = "\033[33m" // Yellow/orange
		case "#FF0000":
			colorCode = "\033[91m" // Bright red
		default:
			colorCode = "\033[97m" // Bright white
		}
		resetCode := "\033[0m"

		b.WriteString(fmt.Sprintf("Difficulty: %s%s%s\n", colorCode, analysis.Difficulty, resetCode))
		b.WriteString(fmt.Sprintf("XP: %d (×%.1f = %d adj.)\n\n",
			analysis.TotalMonsterXP, analysis.Multiplier, analysis.AdjustedXP))
	}

	// Monsters list
	b.WriteString("Encounter Monsters:\n")
	if len(encounterMonsters) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for i, m := range encounterMonsters {
			prefix := "  "
			if i == selectedIndex {
				prefix = "→ "
			}
			b.WriteString(fmt.Sprintf("%s%dx %s (CR %s, AC %d, HP %d)\n",
				prefix, m.Quantity, m.Name, m.CR, m.AC, m.HP))
		}
	}

	return b.String()
}

func renderTemplates(
	savedEncounters []Encounter,
	encounterListMode bool,
	selectedSaved int,
	width, height int,
) string {
	var b strings.Builder

	b.WriteString("⚔️  Saved Encounter Templates\n\n")

	if len(savedEncounters) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for i, enc := range savedEncounters {
			prefix := "  "
			if encounterListMode && i == selectedSaved {
				prefix = "→ "
			}

			// Count total monsters
			totalMonsters := 0
			for _, m := range enc.Monsters {
				totalMonsters += m.Quantity
			}

			b.WriteString(fmt.Sprintf("%s%s (%d monsters)\n", prefix, enc.Name, totalMonsters))
		}
	}

	return b.String()
}

func renderTemplateDetail(
	savedEncounters []Encounter,
	selectedSaved int,
	partySize, partyLevel int,
	width, height int,
) string {
	var b strings.Builder

	if selectedSaved < 0 || selectedSaved >= len(savedEncounters) {
		b.WriteString("⚔️  No Template Selected\n\n")
		b.WriteString("Press [esc] to go back\n")
		return b.String()
	}

	enc := savedEncounters[selectedSaved]

	b.WriteString(fmt.Sprintf("⚔️  %s\n\n", enc.Name))

	// Count total monsters and calculate XP
	totalMonsters := 0
	totalXP := 0
	for _, m := range enc.Monsters {
		totalMonsters += m.Quantity
		totalXP += m.XP * m.Quantity
	}

	b.WriteString(fmt.Sprintf("Total Monsters: %d\n", totalMonsters))
	b.WriteString(fmt.Sprintf("Total XP: %d\n\n", totalXP))

	// Show monsters
	b.WriteString("Monsters:\n")
	for _, m := range enc.Monsters {
		b.WriteString(fmt.Sprintf("  %dx %s (CR %s, %d XP, AC %d, HP %d)\n",
			m.Quantity, m.Name, m.CR, m.XP, m.AC, m.HP))
	}

	b.WriteString("\n")

	// If we have party info, calculate difficulty
	if partySize > 0 && partyLevel > 0 {
		b.WriteString("For your current party:\n")
		b.WriteString(fmt.Sprintf("Party: %d × Level %d\n", partySize, partyLevel))

		// Calculate difficulty
		partyLevels := make([]int, partySize)
		for i := 0; i < partySize; i++ {
			partyLevels[i] = partyLevel
		}

		monsterCRs := []string{}
		for _, m := range enc.Monsters {
			for i := 0; i < m.Quantity; i++ {
				monsterCRs = append(monsterCRs, m.CR)
			}
		}

		// We need to import the encounters package
		// For now, show basic info
		multiplier := getMultiplier(len(monsterCRs))
		adjustedXP := int(float64(totalXP) * multiplier)

		b.WriteString(fmt.Sprintf("Adjusted XP: %d (×%.1f multiplier)\n", adjustedXP, multiplier))
		b.WriteString(fmt.Sprintf("Difficulty: %s\n", estimateDifficulty(partySize, partyLevel, adjustedXP)))
	}

	return b.String()
}
