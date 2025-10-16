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

	// Difficulty selection (with focus indicator)
	diffFocused := m.EncounterGeneratorFocus == "difficulty"
	if diffFocused {
		content.WriteString("► Difficulty: (↑/↓ to change)\n")
	} else {
		content.WriteString("  Difficulty:\n")
	}
	for i, diff := range difficulties {
		prefix := "    "
		if i == m.EncounterDifficultyIndex {
			prefix = "  → " // Always show selection arrow
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, strings.Title(diff)))
	}
	content.WriteString("\n")

	// Environment selection (with focus indicator and scrolling)
	envFocused := m.EncounterGeneratorFocus == "environment"
	if envFocused {
		content.WriteString("► Environment: (↑/↓ to change)\n")
	} else {
		content.WriteString("  Environment:\n")
	}

	// Show scrollable window of environments (8 visible at a time)
	visibleCount := 8
	totalEnvs := len(m.AvailableEnvironments)
	selectedIdx := m.EncounterEnvironmentIndex

	// Calculate scroll window
	startIdx := 0
	endIdx := totalEnvs

	if totalEnvs > visibleCount {
		// Center the selected item in the visible window
		startIdx = selectedIdx - visibleCount/2
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + visibleCount
		if endIdx > totalEnvs {
			endIdx = totalEnvs
			startIdx = endIdx - visibleCount
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	// Show indicator if there are items above
	if startIdx > 0 {
		content.WriteString("    ▲ more above\n")
	}

	// Show visible environments
	for i := startIdx; i < endIdx; i++ {
		prefix := "    "
		if i == selectedIdx {
			prefix = "  → " // Always show selection arrow
		}
		content.WriteString(fmt.Sprintf("%s%s\n", prefix, m.AvailableEnvironments[i]))
	}

	// Show indicator if there are items below
	if endIdx < totalEnvs {
		content.WriteString("    ▼ more below\n")
	}

	content.WriteString("\n")
	content.WriteString("Commands:\n")
	content.WriteString("  [tab]   Switch between fields\n")
	content.WriteString("  [↑/↓]   Navigate options\n")
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
