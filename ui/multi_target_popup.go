// ui/multi_target_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderMultiTargetPopup renders the multi-target damage/healing popup
func RenderMultiTargetPopup(m Model) string {
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(60).
		Background(lipgloss.Color("#1a1a1a"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center).
		Width(56)

	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(54).
		Foreground(lipgloss.Color("#FAFAFA"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Width(56)

	targetStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	saveSuccessStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	saveFailureStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	// Build content
	var content strings.Builder

	// Title
	title := "🎯 Multi-Target "
	if m.MultiTargetType == "healing" {
		title += "Healing"
	} else {
		title += "Damage"
	}
	content.WriteString(titleStyle.Render(title) + "\n\n")

	// Selected targets
	content.WriteString("Selected Targets:\n")
	targetCount := 0
	for i, entry := range m.InitiativeList {
		if m.SelectedTargets[i] {
			targetCount++
			line := fmt.Sprintf("  • %s", entry.Name)
			if entry.Type == "monster" {
				line += fmt.Sprintf(" (HP: %d/%d)", entry.HP, entry.MaxHP)
			}

			// Show save result if in save mode
			if m.MultiTargetSaveMode {
				saveResult := m.TargetSaveResults[i]
				if saveResult == "success" {
					line += " " + saveSuccessStyle.Render("[SAVED]")
				} else if saveResult == "failure" {
					line += " " + saveFailureStyle.Render("[FAILED]")
				} else {
					line += " [Press 's' for success, 'f' for failure]"
				}
			}

			content.WriteString(targetStyle.Render(line) + "\n")
		}
	}

	if targetCount == 0 {
		content.WriteString(targetStyle.Render("  (No targets selected)") + "\n")
	}
	content.WriteString("\n")

	// Amount input
	actionVerb := "Damage"
	if m.MultiTargetType == "healing" {
		actionVerb = "Healing"
	}
	content.WriteString("Enter Amount:\n")
	content.WriteString(helpStyle.Render("(Use -10 for damage, +10 for healing, or plain 10)") + "\n")
	content.WriteString(inputBoxStyle.Render(m.MultiTargetInput+"│") + "\n\n")

	// Show current mode
	modeText := fmt.Sprintf("Current Mode: %s", actionVerb)
	content.WriteString(helpStyle.Render(modeText) + "\n")

	// Save mode toggle
	saveModeText := "Save Mode: OFF (full damage)"
	if m.MultiTargetSaveMode {
		saveModeText = "Save Mode: ON (half damage on success)"
	}
	content.WriteString(helpStyle.Render(saveModeText) + "\n")
	content.WriteString(helpStyle.Render("Press 'x' to toggle save mode") + "\n\n")

	// Help text
	if m.MultiTargetSaveMode {
		content.WriteString(helpStyle.Render("s: mark save success • f: mark save failure") + "\n")
	}
	content.WriteString(helpStyle.Render("h: toggle mode • x: toggle save mode") + "\n")
	content.WriteString(helpStyle.Render("Enter: apply • Esc: cancel"))

	return popupStyle.Render(content.String())
}

// renderMultiTargetPopupOverlay renders the multi-target popup as an overlay
func renderMultiTargetPopupOverlay(m Model, baseView string) string {
	if !m.ShowMultiTargetPopup {
		return baseView
	}

	popup := RenderMultiTargetPopup(m)
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}

// ApplyMultiTargetDamage applies damage or healing to all selected targets
func ApplyMultiTargetDamage(m Model, amount int) Model {
	for i := range m.InitiativeList {
		if !m.SelectedTargets[i] {
			continue
		}

		entry := &m.InitiativeList[i]

		// Only apply to monsters (players manage their own HP)
		if entry.Type != "monster" {
			continue
		}

		actualAmount := amount

		// Apply save logic if in save mode
		if m.MultiTargetSaveMode {
			saveResult := m.TargetSaveResults[i]
			if saveResult == "success" && m.MultiTargetType == "damage" {
				// Half damage on successful save
				actualAmount = amount / 2
			} else if saveResult == "failure" {
				// Full damage/healing on failed save
				actualAmount = amount
			} else {
				// No save result recorded, skip this target
				continue
			}
		}

		// Apply damage or healing
		if m.MultiTargetType == "healing" {
			entry.HP += actualAmount
			// Cap at max HP
			if entry.HP > entry.MaxHP {
				entry.HP = entry.MaxHP
			}
		} else {
			// Damage
			entry.HP -= actualAmount
			// Cap at 0 HP
			if entry.HP < 0 {
				entry.HP = 0
			}
		}
	}

	return m
}
