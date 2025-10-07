// ui/save_popup.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderSavePopup renders the save campaign popup
func RenderSavePopup(m Model) string {
	// Popup dimensions
	popupWidth := 60
	popupHeight := 10

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5555FF")).
		Padding(0, 1).
		Width(popupWidth - 2).
		Align(lipgloss.Center).
		Render("💾 Save Campaign")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(1, 2).
		Render("Enter campaign name:")

	// Input field
	inputValue := m.SaveInput
	if inputValue == "" {
		inputValue = "my_campaign"
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5555FF")).
		Padding(0, 1).
		Width(popupWidth - 8)

	inputField := inputStyle.Render(inputValue + "█")

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(1, 2).
		Render("Enter: Save  |  Esc: Cancel")

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		instructions,
		inputField,
		helpText,
	)

	// Create popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#5555FF")).
		Width(popupWidth).
		Height(popupHeight).
		Padding(0)

	popup := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)

	return popupStyle.Render(popup)
}

// RenderLoadPopup renders the load campaign popup
func RenderLoadPopup(m Model) string {
	// Popup dimensions
	popupWidth := 60
	maxPopupHeight := 20

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#55AA55")).
		Padding(0, 1).
		Width(popupWidth - 2).
		Align(lipgloss.Center).
		Render("📂 Load Campaign")

	// Check if there are campaigns
	if len(m.CampaignList) == 0 {
		instructions := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(2, 2).
			Align(lipgloss.Center).
			Render("No saved campaigns found.\n\nPress Ctrl+S to save your first campaign!")

		helpText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(1, 2).
			Render("Esc: Cancel")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			instructions,
			helpText,
		)

		popupStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#55AA55")).
			Width(popupWidth).
			Padding(0)

		popup := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			content,
		)

		return popupStyle.Render(popup)
	}

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(1, 2).
		Render("Select a campaign to load:")

	// Campaign list
	var campaignLines []string
	visibleStart := 0
	visibleEnd := len(m.CampaignList)
	maxVisible := 8

	// Calculate visible range if list is too long
	if len(m.CampaignList) > maxVisible {
		if m.CampaignListIndex >= maxVisible {
			visibleStart = m.CampaignListIndex - maxVisible + 1
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(m.CampaignList) {
			visibleEnd = len(m.CampaignList)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		campaign := m.CampaignList[i]
		displayName := GetCampaignDisplayName(campaign)

		if i == m.CampaignListIndex {
			// Selected campaign
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#55AA55")).
				Padding(0, 1).
				Width(popupWidth - 6).
				Render("▶ " + displayName)
			campaignLines = append(campaignLines, line)
		} else {
			// Unselected campaign
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Padding(0, 1).
				Width(popupWidth - 6).
				Render("  " + displayName)
			campaignLines = append(campaignLines, line)
		}
	}

	// Add scroll indicators if needed
	if len(m.CampaignList) > maxVisible {
		scrollInfo := fmt.Sprintf("(%d/%d)", m.CampaignListIndex+1, len(m.CampaignList))
		scrollLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Padding(0, 2).
			Render(scrollInfo)
		campaignLines = append(campaignLines, scrollLine)
	}

	campaignListContent := lipgloss.NewStyle().
		Padding(0, 2).
		Render(strings.Join(campaignLines, "\n"))

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(1, 2).
		Render("↑/↓: Navigate  |  Enter: Load  |  Esc: Cancel")

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		instructions,
		campaignListContent,
		helpText,
	)

	// Create popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#55AA55")).
		Width(popupWidth).
		MaxHeight(maxPopupHeight).
		Padding(0)

	popup := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)

	return popupStyle.Render(popup)
}

// RenderSaveSuccessMessage renders a temporary success message
func RenderSaveSuccessMessage(campaignName string) string {
	message := fmt.Sprintf("✓ Campaign '%s' saved successfully!", campaignName)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#55AA55")).
		Padding(0, 2).
		Bold(true)

	return style.Render(message)
}

// RenderLoadSuccessMessage renders a temporary success message
func RenderLoadSuccessMessage(campaignName string) string {
	message := fmt.Sprintf("✓ Campaign '%s' loaded successfully!", campaignName)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#55AA55")).
		Padding(0, 2).
		Bold(true)

	return style.Render(message)
}

// RenderErrorMessage renders an error message
func RenderErrorMessage(errorMsg string) string {
	message := fmt.Sprintf("✗ Error: %s", errorMsg)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#AA5555")).
		Padding(0, 2).
		Bold(true)

	return style.Render(message)
}

// RenderRenamePopup renders the rename campaign popup
func RenderRenamePopup(m Model) string {
	// Popup dimensions
	popupWidth := 60
	popupHeight := 12

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF9500")).
		Padding(0, 1).
		Width(popupWidth - 2).
		Align(lipgloss.Center).
		Render("✏️  Rename Campaign")

	// Current name display
	currentName := m.CurrentCampaignName
	if currentName == "" {
		currentName = "No campaign loaded"
	}

	currentNameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(1, 2).
		Render(fmt.Sprintf("Current: %s", currentName))

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 2).
		Render("Enter new campaign name:")

	// Input field
	inputValue := m.SaveInput
	if inputValue == "" {
		inputValue = currentName
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF9500")).
		Padding(0, 1).
		Width(popupWidth - 8)

	inputField := inputStyle.Render(inputValue + "█")

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(1, 2).
		Render("Enter: Rename  |  Esc: Cancel")

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		currentNameText,
		instructions,
		inputField,
		helpText,
	)

	// Create popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF9500")).
		Width(popupWidth).
		Height(popupHeight).
		Padding(0)

	popup := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)

	return popupStyle.Render(popup)
}
