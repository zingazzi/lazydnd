// ui/quick_hp_popup.go
package ui

import (
	"fmt"
	"strings"
)

// RenderQuickHPPopup renders the quick HP adjustment popup (plain text for TView)
func RenderQuickHPPopup(m Model) string {
	if !m.ShowQuickHPPopup {
		return ""
	}

	// Count selected targets (for multi-target mode)
	selectedCount := 0
	var targetNames []string

	if m.MultiTargetMode {
		// Multi-target mode - show all selected
		for i, selected := range m.SelectedTargets {
			if selected && i >= 0 && i < len(m.InitiativeList) {
				selectedCount++
				targetNames = append(targetNames, m.InitiativeList[i].Name)
				if selectedCount >= 3 {
					targetNames = append(targetNames, "...")
					break
				}
			}
		}
	} else {
		// Single target mode
		if m.SelectedEntry < 0 || m.SelectedEntry >= len(m.InitiativeList) {
			return ""
		}
		selectedCount = 1
		targetNames = []string{m.InitiativeList[m.SelectedEntry].Name}
	}

	// Determine title and prompt based on mode
	var title string
	var prompt string
	if m.QuickHPMode == "add" {
		title = "⚕️  ADD HP"
		if selectedCount > 1 {
			prompt = fmt.Sprintf("Amount to add to %d targets:", selectedCount)
		} else {
			prompt = "Amount to add:"
		}
	} else {
		title = "💥 REMOVE HP"
		if selectedCount > 1 {
			prompt = fmt.Sprintf("Amount to remove from %d targets:", selectedCount)
		} else {
			prompt = "Amount to remove:"
		}
	}

	// Build popup content
	var content strings.Builder

	// Title
	content.WriteString(title)
	content.WriteString("\n\n")

	// Target info
	if selectedCount == 1 {
		entry := m.InitiativeList[m.SelectedEntry]
		targetInfo := fmt.Sprintf("%s • HP: %d/%d", entry.Name, entry.HP, entry.MaxHP)
		if entry.TempHP > 0 {
			targetInfo += fmt.Sprintf(" +%d", entry.TempHP)
		}
		content.WriteString(targetInfo)
	} else {
		content.WriteString(strings.Join(targetNames, ", "))
	}
	content.WriteString("\n\n")

	// Prompt
	content.WriteString(prompt)
	content.WriteString("\n\n")

	// Input field
	displayText := m.QuickHPInput
	if displayText == "" {
		displayText = "0"
	}
	displayText += " █" // Cursor
	content.WriteString(fmt.Sprintf("[%s]\n\n", displayText))

	// Instructions
	content.WriteString("↵ Enter: Apply • Esc: Cancel")

	return content.String()
}

// renderQuickHPPopupOverlay is deprecated - TView handles overlays
func (m Model) renderQuickHPPopupOverlay(mainView string) string {
	return mainView
}
