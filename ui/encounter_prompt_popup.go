// ui/encounter_prompt_popup.go
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderEncounterPromptPopup renders the save encounter prompt popup
func RenderEncounterPromptPopup(m Model) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA"))

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true).
		MarginTop(1)

	// Build content
	title := titleStyle.Render("💾 Save Encounter")
	label := labelStyle.Render("Name:")
	input := inputStyle.Render(m.EncounterNameInput + "█")
	hint := hintStyle.Render("[enter] save  [esc] cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		label,
		input,
		hint,
	)

	// Popup box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(1, 2).
		Width(50)

	return boxStyle.Render(content)
}


