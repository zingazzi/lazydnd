// ui/encounter_generator_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderGeneratorPopup renders the encounter generator configuration popup
func RenderGeneratorPopup(m Model) string {
	difficulties := []string{"easy", "medium", "hard", "deadly"}

	var content strings.Builder
	content.WriteString("⚔️  Generate Encounter\n\n")
	content.WriteString(fmt.Sprintf("Party: %d × Level %d\n\n", m.PartySize, m.PartyLevel))

	// Difficulty selection
	content.WriteString("Difficulty:\n")
	for i, diff := range difficulties {
		prefix := "  "
		if i == m.EncounterDifficultyIndex {
			prefix = "→ "
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, strings.Title(diff)))
	}
	content.WriteString("\n")

	// Environment selection
	content.WriteString("Environment:\n")
	for i, env := range m.AvailableEnvironments {
		if i >= 8 {
			break // Limit display
		}
		prefix := "  "
		if i == m.EncounterEnvironmentIndex {
			prefix = "→ "
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, env))
	}

	if len(m.AvailableEnvironments) > 8 {
		content.WriteString("  ...\n")
	}

	content.WriteString("\n")
	content.WriteString("Commands:\n")
	content.WriteString("  [↑/↓]   Navigate\n")
	content.WriteString("  [tab]   Switch difficulty/environment\n")
	content.WriteString("  [enter] Generate\n")
	content.WriteString("  [esc]   Cancel\n")

	// Style the popup
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.Config.Theme.PrimaryColor)).
		Padding(1, 2).
		Width(50)

	return popupStyle.Render(content.String())
}


