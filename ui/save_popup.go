// ui/save_popup.go
package ui

import (
	"fmt"
	"strings"
)

// RenderSavePopup renders the save campaign popup (plain text for TView)
func RenderSavePopup(m Model) string {
	var content strings.Builder

	content.WriteString("💾 Save Campaign\n\n")
	content.WriteString("Enter campaign name:\n\n")

	inputValue := m.SaveInput
	if inputValue == "" {
		inputValue = "my_campaign"
	}
	content.WriteString(fmt.Sprintf("[%s█]\n\n", inputValue))
	content.WriteString("Enter: Save  |  Esc: Cancel")

	return content.String()
}

// RenderLoadPopup renders the load campaign popup (plain text for TView)
func RenderLoadPopup(m Model) string {
	var content strings.Builder

	content.WriteString("📂 Load Campaign\n\n")

	// Check if there are campaigns
	if len(m.CampaignList) == 0 {
		content.WriteString("No saved campaigns found.\n\n")
		content.WriteString("Press Ctrl+S to save your first campaign!")
		return content.String()
	}

	content.WriteString("Select a campaign to load:\n\n")

	// Campaign list
	maxVisible := 8
	visibleStart := 0
	visibleEnd := len(m.CampaignList)

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
			content.WriteString(fmt.Sprintf("▶ %s\n", displayName))
		} else {
			content.WriteString(fmt.Sprintf("  %s\n", displayName))
		}
	}

	// Add scroll indicators if needed
	if len(m.CampaignList) > maxVisible {
		scrollInfo := fmt.Sprintf("\n(%d/%d)", m.CampaignListIndex+1, len(m.CampaignList))
		content.WriteString(scrollInfo)
	}

	content.WriteString("\n\n↑/↓: Navigate  |  Enter: Load  |  Esc: Cancel")

	return content.String()
}

// RenderSaveSuccessMessage renders a temporary success message (plain text)
func RenderSaveSuccessMessage(campaignName string) string {
	return fmt.Sprintf("✓ Campaign '%s' saved successfully!", campaignName)
}

// RenderLoadSuccessMessage renders a temporary success message (plain text)
func RenderLoadSuccessMessage(campaignName string) string {
	return fmt.Sprintf("✓ Campaign '%s' loaded successfully!", campaignName)
}

// RenderErrorMessage renders an error message (plain text)
func RenderErrorMessage(errorMsg string) string {
	return fmt.Sprintf("✗ Error: %s", errorMsg)
}

// RenderRenamePopup renders the rename campaign popup (plain text for TView)
func RenderRenamePopup(m Model) string {
	var content strings.Builder

	content.WriteString("✏️  Rename Campaign\n\n")

	// Current name display
	currentName := m.CurrentCampaignName
	if currentName == "" {
		currentName = "No campaign loaded"
	}
	content.WriteString(fmt.Sprintf("Current name: %s\n\n", currentName))
	content.WriteString("Enter new name:\n\n")

	inputValue := m.SaveInput
	if inputValue == "" {
		inputValue = currentName
	}
	content.WriteString(fmt.Sprintf("[%s█]\n\n", inputValue))
	content.WriteString("Enter: Rename  |  Esc: Cancel")

	return content.String()
}

// renderSavePopupOverlay is deprecated - TView handles overlays
func (m Model) renderSavePopupOverlay(mainView string) string {
	return mainView
}

// renderLoadPopupOverlay is deprecated - TView handles overlays
func (m Model) renderLoadPopupOverlay(mainView string) string {
	return mainView
}

// renderRenamePopupOverlay is deprecated - TView handles overlays
func (m Model) renderRenamePopupOverlay(mainView string) string {
	return mainView
}
