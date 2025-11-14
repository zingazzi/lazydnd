// ui/multi_target_popup.go
package ui

import (
	"fmt"
	"strings"
)

// RenderMultiTargetPopup renders the multi-target damage/healing popup (plain text for TView)
func RenderMultiTargetPopup(m Model) string {
	var content strings.Builder

	// Title
	title := "🎯 Multi-Target "
	if m.MultiTargetType == "healing" {
		title += "Healing"
	} else {
		title += "Damage"
	}
	content.WriteString(title)
	content.WriteString("\n\n")

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
					line += " [SAVED]"
				} else if saveResult == "failure" {
					line += " [FAILED]"
				} else {
					line += " [Press 's' for success, 'f' for failure]"
				}
			}

			content.WriteString(line + "\n")
		}
	}

	if targetCount == 0 {
		content.WriteString("  (No targets selected)\n")
	}
	content.WriteString("\n")

	// Amount input
	actionVerb := "Damage"
	if m.MultiTargetType == "healing" {
		actionVerb = "Healing"
	}
	content.WriteString("Enter Amount:\n")
	content.WriteString("(Use -10 for damage, +10 for healing, or plain 10)\n")
	content.WriteString(fmt.Sprintf("[%s█]\n\n", m.MultiTargetInput))

	// Show current mode
	modeText := fmt.Sprintf("Current Mode: %s", actionVerb)
	content.WriteString(modeText + "\n")

	// Save mode toggle
	saveModeText := "Save Mode: OFF (full damage)"
	if m.MultiTargetSaveMode {
		saveModeText = "Save Mode: ON (half damage on success)"
	}
	content.WriteString(saveModeText + "\n")
	content.WriteString("Press 'x' to toggle save mode\n\n")

	// Help text
	if m.MultiTargetSaveMode {
		content.WriteString("s: mark save success • f: mark save failure\n")
	}
	content.WriteString("h: toggle mode • x: toggle save mode\n")
	content.WriteString("Enter: apply • Esc: cancel")

	return content.String()
}

// renderMultiTargetPopupOverlay is deprecated - TView handles overlays
func renderMultiTargetPopupOverlay(m Model, baseView string) string {
	return baseView
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

		// Save old HP for undo history
		oldHP := entry.HP

		// Apply damage or healing
		if m.MultiTargetType == "healing" {
			// Healing - only affects real HP
			newHP := entry.HP + actualAmount
			// Cap at max HP
			if newHP > entry.MaxHP {
				newHP = entry.MaxHP
			}
			entry.HP = newHP
		} else {
			// Damage - apply to temp HP first
			if entry.TempHP > 0 {
				if actualAmount <= entry.TempHP {
					// All damage absorbed by temp HP
					entry.TempHP -= actualAmount
				} else {
					// Temp HP absorbed some, rest goes to real HP
					remainingDamage := actualAmount - entry.TempHP
					entry.TempHP = 0
					entry.HP -= remainingDamage
					// Cap at 0 HP
					if entry.HP < 0 {
						entry.HP = 0
					}
				}
			} else {
				// No temp HP, damage goes directly to real HP
				entry.HP -= actualAmount
				// Cap at 0 HP
				if entry.HP < 0 {
					entry.HP = 0
				}
			}
		}

		// Save to undo history
		pushHPHistory(&m, i, oldHP, entry.HP)
	}

	return m
}
