// ui/condition_popup.go
package ui

import (
	"fmt"
	"strings"
)

// Common D&D 5e conditions for quick selection
var CommonConditions = []string{
	"Blinded",
	"Charmed",
	"Deafened",
	"Exhausted",
	"Frightened",
	"Grappled",
	"Incapacitated",
	"Invisible",
	"Paralyzed",
	"Petrified",
	"Poisoned",
	"Prone",
	"Restrained",
	"Stunned",
	"Unconscious",
}

// RenderConditionPopup renders the condition management popup (plain text for TView)
func RenderConditionPopup(m Model) string {
	var content strings.Builder

	// Get selected entry name(s)
	entryName := ""
	if m.MultiTargetMode {
		// Multi-target mode: show count
		selectedCount := 0
		for _, selected := range m.SelectedTargets {
			if selected {
				selectedCount++
			}
		}
		entryName = fmt.Sprintf("%d targets selected", selectedCount)
	} else if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
		entryName = m.InitiativeList[m.SelectedEntry].Name
	}

	if m.ConditionPopupMode == "list" {
		// List mode: show current conditions
		content.WriteString(fmt.Sprintf("🔮 Conditions - %s\n\n", entryName))

		// In multi-target mode, don't show existing conditions (just allow adding)
		if m.MultiTargetMode {
			content.WriteString("Press 'a' to add a condition to all targets\n\n")
		} else if m.SelectedEntry >= 0 && m.SelectedEntry < len(m.InitiativeList) {
			if len(m.InitiativeList[m.SelectedEntry].Conditions) == 0 {
				content.WriteString("No active conditions\n")
				content.WriteString("Press 'a' to add a condition\n\n")
			} else {
				content.WriteString("Active Conditions:\n\n")
				for i, cond := range m.InitiativeList[m.SelectedEntry].Conditions {
					var line string
					var emoji string

					// Choose emoji based on condition type
					switch cond.Name {
					case "Poisoned":
						emoji = "🤢"
					case "Stunned", "Paralyzed", "Incapacitated", "Unconscious":
						emoji = "😵"
					case "Frightened":
						emoji = "😱"
					case "Charmed":
						emoji = "😍"
					case "Invisible":
						emoji = "👻"
					case "Prone":
						emoji = "🤕"
					case "Grappled", "Restrained":
						emoji = "🔗"
					case "Blinded":
						emoji = "🙈"
					case "Deafened":
						emoji = "🙉"
					default:
						emoji = "🔮"
					}

					if cond.RoundsLeft == 0 {
						line = fmt.Sprintf("%s %s (Indefinite)", emoji, cond.Name)
					} else if cond.RoundsLeft == 1 {
						line = fmt.Sprintf("%s %s (1 round left)", emoji, cond.Name)
					} else {
						line = fmt.Sprintf("%s %s (%d rounds left)", emoji, cond.Name, cond.RoundsLeft)
					}

					if i == m.SelectedConditionIdx {
						content.WriteString(fmt.Sprintf("► %s\n", line))
					} else {
						content.WriteString(fmt.Sprintf("  %s\n", line))
					}
				}
				content.WriteString("\n")
			}
		}

		// Show appropriate help text based on mode
		if m.MultiTargetMode {
			content.WriteString("a: add new • Esc: close")
		} else {
			content.WriteString("↑↓: navigate • d: remove • a: add new • Esc: close")
		}

	} else if m.ConditionPopupMode == "add" {
		// Add mode: select from list
		content.WriteString(fmt.Sprintf("➕ Add Condition - %s\n\n", entryName))

		if m.ConditionInputStep == 0 {
			// Step 1: Select condition from list
			content.WriteString("Select Condition:\n\n")

			for i, cond := range CommonConditions {
				var line string
				// Add emoji to each condition
				switch cond {
				case "Poisoned":
					line = fmt.Sprintf("🤢 %s", cond)
				case "Stunned", "Paralyzed", "Incapacitated", "Unconscious":
					line = fmt.Sprintf("😵 %s", cond)
				case "Frightened":
					line = fmt.Sprintf("😱 %s", cond)
				case "Charmed":
					line = fmt.Sprintf("😍 %s", cond)
				case "Invisible":
					line = fmt.Sprintf("👻 %s", cond)
				case "Prone":
					line = fmt.Sprintf("🤕 %s", cond)
				case "Grappled", "Restrained":
					line = fmt.Sprintf("🔗 %s", cond)
				case "Blinded":
					line = fmt.Sprintf("🙈 %s", cond)
				case "Deafened":
					line = fmt.Sprintf("🙉 %s", cond)
				default:
					line = fmt.Sprintf("🔮 %s", cond)
				}

				if i == m.SelectedConditionNameIdx {
					content.WriteString(fmt.Sprintf("► %s\n", line))
				} else {
					content.WriteString(fmt.Sprintf("  %s\n", line))
				}
			}
			content.WriteString("\n")

			// Add option for custom condition
			customLine := "✏️  Custom condition..."
			if m.SelectedConditionNameIdx == len(CommonConditions) {
				content.WriteString(fmt.Sprintf("► %s\n\n", customLine))
			} else {
				content.WriteString(fmt.Sprintf("  %s\n\n", customLine))
			}

			content.WriteString("↑↓: navigate • Enter: select • Esc: cancel")

		} else if m.ConditionInputStep == 1 {
			// Step 2: Enter duration
			selectedCondition := ""
			if m.SelectedConditionNameIdx < len(CommonConditions) {
				selectedCondition = CommonConditions[m.SelectedConditionNameIdx]
			} else {
				selectedCondition = m.ConditionInput
			}

			content.WriteString(fmt.Sprintf("Condition: %s\n\n", selectedCondition))
			content.WriteString("Duration (in rounds):\n")
			content.WriteString("(Enter 0 for indefinite duration)\n")
			content.WriteString(fmt.Sprintf("[%s█]\n\n", m.ConditionDurationInput))

			content.WriteString("Enter: apply • Esc: cancel")

		} else if m.ConditionInputStep == 2 {
			// Step 3: Custom condition name
			content.WriteString("Enter custom condition name:\n")
			content.WriteString(fmt.Sprintf("[%s█]\n\n", m.ConditionInput))
			content.WriteString("Enter: next • Esc: cancel")
		}
	}

	return content.String()
}

// renderConditionPopupOverlay is deprecated - TView handles overlays
func renderConditionPopupOverlay(m Model, baseView string) string {
	return baseView
}
