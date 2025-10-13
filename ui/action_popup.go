// ui/action_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Action popup styles
var (
	ActionPopupStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFFFFF")).
				Padding(1, 2).
				Width(70).
				MaxHeight(30)

	ActionPopupTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Align(lipgloss.Center)

	ActionItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	ActionItemSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("#5555FF")).
				Foreground(lipgloss.Color("#FFFFFF"))
)

const (
	ActionPopupTitle = "🗡️  Monster Actions"
)

// renderActionPopupOverlay renders the action popup over the main view
func (m Model) renderActionPopupOverlay(mainView string) string {
	actionContent := m.buildActionContent()

	// Create the popup
	popup := ActionPopupStyle.Render(actionContent)

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

// buildActionContent builds the action popup content
func (m Model) buildActionContent() string {
	var content strings.Builder

	// Title
	title := ActionPopupTitleStyle.Render(fmt.Sprintf("%s - %s", ActionPopupTitle, m.ActionPopupMonster))
	content.WriteString(title)
	content.WriteString("\n\n")

	// Show ActionNumber if available (attacks per round)
	if len(m.ActionPopupActions) > 0 {
		// Get the monster from initiative list to access ActionNumber
		for _, entry := range m.InitiativeList {
			if entry.MonsterData != nil && entry.MonsterData.Name == m.ActionPopupMonster {
				if entry.MonsterData.ActionNumber > 0 {
					attacksPerRound := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FFD700")).
						Bold(true).
						Render(fmt.Sprintf("⚔️  Attacks per Round: %d", entry.MonsterData.ActionNumber))
					content.WriteString(attacksPerRound)
					content.WriteString("\n\n")
				}
				break
			}
		}
	}

	// Actions list
	for i, action := range m.ActionPopupActions {
		isSelected := i == m.ActionPopupIndex
		actionLine := m.formatActionLine(action, isSelected)
		content.WriteString(actionLine)
		content.WriteString("\n")
	}

	// Show current advantage/disadvantage mode
	if m.ActionPopupAdvantage || m.ActionPopupDisadvantage {
		content.WriteString("\n")
		modeText := ""
		if m.ActionPopupAdvantage {
			modeText = "⚡ ADVANTAGE ⚡"
		} else if m.ActionPopupDisadvantage {
			modeText = "⚠️  DISADVANTAGE ⚠️"
		}
		modeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
		content.WriteString(modeStyle.Render(modeText))
	}

	// Instructions
	content.WriteString("\n")
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("↑/↓: Navigate | a: Advantage | d: Disadvantage | Enter: Roll | Esc: Close")
	content.WriteString(instructions)

	return content.String()
}

// formatActionLine formats a single action line for display
func (m Model) formatActionLine(action MonsterAction, isSelected bool) string {
	var line strings.Builder

	// Action name and type
	actionHeader := fmt.Sprintf("%s (%s)", action.Name, action.Type)

	// Build details
	var details []string
	if action.Roll != "" {
		details = append(details, fmt.Sprintf("Roll: 1d20%s", action.Roll))
	}
	if action.Reach != "" {
		details = append(details, fmt.Sprintf("Reach: %s", action.Reach))
	}
	if action.Range != "" {
		details = append(details, fmt.Sprintf("Range: %s", action.Range))
	}
	if action.Damage != "" {
		damageInfo := action.Damage
		if action.DamageType != "" {
			damageInfo += " " + action.DamageType
		}
		details = append(details, fmt.Sprintf("Damage: %s", damageInfo))
	}
	if action.SaveDC != "" {
		saveInfo := action.SaveDC
		if action.SaveType != "" {
			saveInfo += " " + action.SaveType
		}
		details = append(details, fmt.Sprintf("Save: %s", saveInfo))
	}

	// Format the line
	if len(details) > 0 {
		line.WriteString(fmt.Sprintf("%s\n  %s", actionHeader, strings.Join(details, " | ")))
	} else {
		line.WriteString(actionHeader)
	}

	// Apply style based on selection
	if isSelected {
		return ActionItemSelectedStyle.Render(line.String())
	}
	return ActionItemStyle.Render(line.String())
}
