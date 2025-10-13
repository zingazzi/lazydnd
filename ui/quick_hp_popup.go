// ui/quick_hp_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderQuickHPPopup renders the quick HP adjustment popup
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
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5555FF")).
		Padding(0, 1).
		Width(42).
		Align(lipgloss.Center)

	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	// Target info (compact)
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")).
		Italic(true)

	if selectedCount == 1 {
		entry := m.InitiativeList[m.SelectedEntry]
		targetInfo := fmt.Sprintf("%s • HP: %d/%d", entry.Name, entry.HP, entry.MaxHP)
		if entry.TempHP > 0 {
			targetInfo += fmt.Sprintf(" +%d", entry.TempHP)
		}
		content.WriteString(infoStyle.Render(targetInfo))
	} else {
		content.WriteString(infoStyle.Render(strings.Join(targetNames, ", ")))
	}
	content.WriteString("\n\n")

	// Prompt
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)
	content.WriteString(promptStyle.Render(prompt))
	content.WriteString("\n\n")

	// Input field - PROMINENT and VISIBLE
	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#00FF00")).
		Background(lipgloss.Color("#1a1a1a")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(1, 2).
		Width(36).
		Align(lipgloss.Center)

	// Show the input with a visible cursor
	displayText := m.QuickHPInput
	if displayText == "" {
		displayText = "0"
	}
	displayText += " █" // Cursor

	content.WriteString(inputBoxStyle.Render(displayText))
	content.WriteString("\n\n")

	// Instructions
	instructStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)

	content.WriteString(instructStyle.Render("↵ Enter: Apply • Esc: Cancel"))

	// Create popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5555FF")).
		Padding(1, 2).
		Width(46)

	return popupStyle.Render(content.String())
}
