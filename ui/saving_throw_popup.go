// ui/saving_throw_popup.go
package ui

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SavingThrowResult holds the result of a saving throw roll
type SavingThrowResult struct {
	Ability      string
	Roll         int
	Modifier     int
	Total        int
	IsProficient bool
}

// SkillCheckResult holds the result of a skill check roll
type SkillCheckResult struct {
	Skill    string
	Roll     int
	Modifier int
	Total    int
}

// renderSavingThrowPopupOverlay renders the saving throw popup over the main view
func (m Model) renderSavingThrowPopupOverlay(mainView string) string {
	popup := RenderSavingThrowPopup(m)

	// Place popup centered over the main view
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333333")),
	)
}

// RenderSavingThrowPopup renders a popup showing saving throw rolls for a monster
func RenderSavingThrowPopup(m Model) string {
	if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
		return ""
	}

	entry := m.InitiativeList[m.SelectedEntry]

	// Only works for monsters with data
	if entry.Type != "monster" || entry.MonsterData == nil {
		return renderErrorPopup("Saving throws only available for monsters from the monster panel")
	}

	// Get monster data
	monster := entry.MonsterData

	// Parse saving throw totals (includes base modifier + proficiency)
	savingThrowTotals := parseSavingThrowBonuses(monster.SavingThrows)

	// Roll saving throws for all abilities
	results := []SavingThrowResult{
		rollSavingThrow("STR", monster.STRMod, savingThrowTotals),
		rollSavingThrow("DEX", monster.DEXMod, savingThrowTotals),
		rollSavingThrow("CON", monster.CONMod, savingThrowTotals),
		rollSavingThrow("INT", monster.INTMod, savingThrowTotals),
		rollSavingThrow("WIS", monster.WISMod, savingThrowTotals),
		rollSavingThrow("CHA", monster.CHAMod, savingThrowTotals),
	}

	// Parse and roll skill checks (Stealth and Perception)
	skillBonuses := parseSkillBonuses(monster.Skills)
	skillChecks := []SkillCheckResult{}

	// Add Stealth if monster has it
	if bonus, hasSkill := skillBonuses["STEALTH"]; hasSkill {
		skillChecks = append(skillChecks, rollSkillCheck("Stealth", bonus))
	}

	// Add Perception if monster has it
	if bonus, hasSkill := skillBonuses["PERCEPTION"]; hasSkill {
		skillChecks = append(skillChecks, rollSkillCheck("Perception", bonus))
	}

	// Build popup content
	content := buildSavingThrowContent(entry.Name, results, skillChecks)

	// Style the popup
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	return popupStyle.Render(content)
}

// rollSavingThrow rolls a d20 and adds the modifier
func rollSavingThrow(ability, modifierStr string, savingThrowTotals map[string]int) SavingThrowResult {
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(20) + 1

	// Check if this ability has a saving throw proficiency
	var modifier int
	var isProficient bool

	if totalBonus, hasProficiency := savingThrowTotals[ability]; hasProficiency {
		// Monster has proficiency in this save - use the total bonus from the stat block
		modifier = totalBonus
		isProficient = true
	} else {
		// No proficiency - just use the base ability modifier
		modifier = 0
		if modifierStr != "" {
			// Remove parentheses and parse: "(+5)" -> 5, "(-1)" -> -1
			cleanMod := strings.TrimSpace(modifierStr)
			cleanMod = strings.Trim(cleanMod, "()")
			if val, err := strconv.Atoi(cleanMod); err == nil {
				modifier = val
			}
		}
		isProficient = false
	}

	total := roll + modifier

	return SavingThrowResult{
		Ability:      ability,
		Roll:         roll,
		Modifier:     modifier,
		Total:        total,
		IsProficient: isProficient,
	}
}

// parseSavingThrowBonuses parses the "Saving Throws" field from monster data
// Format: "Str +6, Dex +4, Con +8" or "Dex +5, Wis +3"
// Returns a map of ability -> total bonus (includes base modifier + proficiency)
func parseSavingThrowBonuses(savingThrowsStr string) map[string]int {
	bonuses := make(map[string]int)

	if savingThrowsStr == "" {
		return bonuses
	}

	// Split by comma
	parts := strings.Split(savingThrowsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Split by space to get ability and bonus
		fields := strings.Fields(part)
		if len(fields) >= 2 {
			ability := strings.ToUpper(fields[0])
			bonusStr := fields[1]

			// Parse the total bonus (includes base modifier + proficiency)
			bonusStr = strings.TrimPrefix(bonusStr, "+")
			if val, err := strconv.Atoi(bonusStr); err == nil {
				bonuses[ability] = val
			}
		}
	}

	return bonuses
}

// parseSkillBonuses parses the "Skills" field from monster data
// Format: "Stealth +6, Perception +3" or "History +12, Perception +10"
// Returns a map of skill -> total bonus
func parseSkillBonuses(skillsStr string) map[string]int {
	bonuses := make(map[string]int)

	if skillsStr == "" {
		return bonuses
	}

	// Split by comma
	parts := strings.Split(skillsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Split by space to get skill and bonus
		fields := strings.Fields(part)
		if len(fields) >= 2 {
			skill := strings.ToUpper(fields[0])
			bonusStr := fields[1]

			// Parse the bonus
			bonusStr = strings.TrimPrefix(bonusStr, "+")
			if val, err := strconv.Atoi(bonusStr); err == nil {
				bonuses[skill] = val
			}
		}
	}

	return bonuses
}

// rollSkillCheck rolls a d20 skill check with the given bonus
func rollSkillCheck(skillName string, bonus int) SkillCheckResult {
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(20) + 1
	total := roll + bonus

	return SkillCheckResult{
		Skill:    skillName,
		Roll:     roll,
		Modifier: bonus,
		Total:    total,
	}
}

// getRollStyle returns the appropriate style based on the roll value
func getRollStyle(roll int) lipgloss.Style {
	if roll == 20 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	} else if roll == 1 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
}

// buildSavingThrowContent formats the saving throw results and skill checks
func buildSavingThrowContent(monsterName string, results []SavingThrowResult, skillChecks []SkillCheckResult) string {
	var lines []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Align(lipgloss.Center)

	lines = append(lines, titleStyle.Render("🎲 SAVING THROWS & SKILLS 🎲"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(monsterName))
	lines = append(lines, strings.Repeat("─", 46))
	lines = append(lines, "")

	// Column headers
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#888888"))

	header := fmt.Sprintf("%-8s  %-10s  %-12s  %-8s", "Ability", "Roll", "Modifier", "Total")
	lines = append(lines, headerStyle.Render(header))
	lines = append(lines, strings.Repeat("─", 46))

	// Results
	for _, result := range results {
		// Style based on roll value
		rollStyle := getRollStyle(result.Roll)

		// Format modifier with + or -
		modStr := fmt.Sprintf("%+d", result.Modifier)

		// Add proficiency indicator
		profIndicator := ""
		if result.IsProficient {
			profIndicator = " ⭐"
		}

		line := fmt.Sprintf("%-8s  %-10s  %-12s  ",
			result.Ability,
			rollStyle.Render(fmt.Sprintf("d20: %d", result.Roll)),
			modStr,
		)

		// Total with special styling
		totalStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
		line += totalStyle.Render(fmt.Sprintf("%d", result.Total)) + profIndicator

		lines = append(lines, line)
	}

	// Add skill checks if any
	if len(skillChecks) > 0 {
		lines = append(lines, "")
		lines = append(lines, strings.Repeat("─", 46))
		lines = append(lines, "")

		// Skill checks header
		skillHeaderStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CED1"))
		lines = append(lines, skillHeaderStyle.Render("SKILL CHECKS"))
		lines = append(lines, "")

		// Skill checks results
		for _, skill := range skillChecks {
			// Style based on roll value
			rollStyle := getRollStyle(skill.Roll)

			// Format modifier
			modStr := fmt.Sprintf("%+d", skill.Modifier)

			line := fmt.Sprintf("%-8s  %-10s  %-12s  ",
				skill.Skill,
				rollStyle.Render(fmt.Sprintf("d20: %d", skill.Roll)),
				modStr,
			)

			// Total with special styling
			totalStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00CED1"))
			line += totalStyle.Render(fmt.Sprintf("%d", skill.Total))

			lines = append(lines, line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, strings.Repeat("─", 46))

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	lines = append(lines, helpStyle.Render("⭐ = Proficient    Enter to reroll    Esc to close"))

	return strings.Join(lines, "\n")
}

// renderErrorPopup renders an error message popup
func renderErrorPopup(message string) string {
	content := fmt.Sprintf("❌ %s", message)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF0000")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	return style.Render(content + "\n\nPress Esc to close")
}
