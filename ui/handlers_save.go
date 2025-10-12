// ui/handlers_save.go
package ui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// handleCtrlS handles Ctrl+S key press (save campaign)
func handleCtrlS(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Don't open save popup if already in another popup or input mode
	if m.ShowHelpPopup || m.ShowActionPopup || m.ShowLoadPopup || m.ShowRenamePopup {
		return m, nil
	}

	// If campaign already exists, just save it
	if m.CurrentCampaignFile != "" && m.CurrentCampaignName != "" {
		err := SaveCampaign(m, m.CurrentCampaignName)
		if err != nil {
			return m, func() tea.Msg {
				return SetErrorMsg{Message: "Failed to save campaign: " + err.Error()}
			}
		}
		m.LastAutoSave = "Just now"
		return m, nil
	}

	// Otherwise open save popup
	m.ShowSavePopup = true
	m.SaveInput = ""

	return m, nil
}

// handleCtrlL handles Ctrl+L key press (load campaign)
func handleCtrlL(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Don't open load popup if already in another popup or input mode
	if m.ShowHelpPopup || m.ShowActionPopup || m.ShowSavePopup || m.ShowRenamePopup {
		return m, nil
	}

	// Load campaign list
	campaigns, err := ListCampaigns()
	if err != nil {
		return m, func() tea.Msg {
			return SetErrorMsg{Message: "Failed to load campaign list: " + err.Error()}
		}
	}

	// Open load popup
	m.ShowLoadPopup = true
	m.CampaignList = campaigns
	m.CampaignListIndex = 0

	return m, nil
}

// handleCtrlN handles Ctrl+N key press (rename campaign)
func handleCtrlN(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Don't open rename popup if already in another popup or no campaign loaded
	if m.ShowHelpPopup || m.ShowActionPopup || m.ShowSavePopup || m.ShowLoadPopup {
		return m, nil
	}

	// Only allow rename if a campaign is loaded
	if m.CurrentCampaignName == "" {
		return m, nil
	}

	// Open rename popup
	m.ShowRenamePopup = true
	m.SaveInput = m.CurrentCampaignName

	return m, nil
}

// handleSavePopupInput handles input when save popup is active
func handleSavePopupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		// Save campaign
		campaignName := m.SaveInput
		if campaignName == "" {
			campaignName = "my_campaign"
		}

		err := SaveCampaign(m, campaignName)
		if err != nil {
			m.ShowSavePopup = false
			m.SaveInput = ""
			return m, func() tea.Msg {
				return SetErrorMsg{Message: "Failed to save campaign: " + err.Error()}
			}
		}

		// Update current campaign file and name
		m.CurrentCampaignFile = sanitizeFilename(campaignName) + ".json"
		m.CurrentCampaignName = campaignName
		m.LastAutoSave = "Just now"

		// Close popup
		m.ShowSavePopup = false
		m.SaveInput = ""

		return m, nil

	case "esc":
		// Cancel save
		m.ShowSavePopup = false
		m.SaveInput = ""
		return m, nil

	case "backspace", "ctrl+h":
		// Remove last character
		if len(m.SaveInput) > 0 {
			m.SaveInput = m.SaveInput[:len(m.SaveInput)-1]
		}
		return m, nil

	default:
		// Add character to input
		if len(key) == 1 {
			m.SaveInput += key
		}
		return m, nil
	}
}

// handleLoadPopupInput handles input when load popup is active
func handleLoadPopupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		// Load selected campaign
		if len(m.CampaignList) == 0 {
			m.ShowLoadPopup = false
			return m, nil
		}

		selectedFile := m.CampaignList[m.CampaignListIndex]
		saveState, initiativeList, err := LoadCampaign(selectedFile)
		if err != nil {
			m.ShowLoadPopup = false
			return m, func() tea.Msg {
				return SetErrorMsg{Message: "Failed to load campaign: " + err.Error()}
			}
		}

		// Update model with loaded data
		m.InitiativeList = initiativeList
		m.CurrentCampaignFile = selectedFile
		m.CurrentCampaignName = saveState.CampaignName
		m.LastAutoSave = "Loaded"
		m.CurrentTurn = saveState.CurrentTurn
		m.RoundCounter = saveState.RoundCounter
		m.ActiveSpells = saveState.ActiveSpells
		m.NotesContent = saveState.Notes
		// Merge saved macros with defaults (user macros take precedence)
		m.DiceMacros = mergeDefaultMacros(saveState.DiceMacros)

		// Reset initiative tracker state
		m.InitiativeListMode = len(initiativeList) > 0
		m.SelectedEntry = 0
		m.InitiativeInputMode = false
		m.InitiativeEditMode = false

		// Close popup
		m.ShowLoadPopup = false
		m.CampaignList = nil
		m.CampaignListIndex = 0

		// Switch to initiative tracker panel
		m.ActivePanel = InitiativeTracker

		return m, nil

	case "esc":
		// Cancel load
		m.ShowLoadPopup = false
		m.CampaignList = nil
		m.CampaignListIndex = 0
		return m, nil

	case "up":
		// Navigate up in campaign list
		if m.CampaignListIndex > 0 {
			m.CampaignListIndex--
		}
		return m, nil

	case "down":
		// Navigate down in campaign list
		if m.CampaignListIndex < len(m.CampaignList)-1 {
			m.CampaignListIndex++
		}
		return m, nil

	default:
		return m, nil
	}
}

// handleRenamePopupInput handles input when rename popup is active
func handleRenamePopupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		// Rename campaign
		newName := m.SaveInput
		if newName == "" || newName == m.CurrentCampaignName {
			// No change, just close
			m.ShowRenamePopup = false
			m.SaveInput = ""
			return m, nil
		}

		// Delete old file
		if m.CurrentCampaignFile != "" {
			saveDir, err := getSaveDirectory()
			if err == nil {
				oldPath := filepath.Join(saveDir, m.CurrentCampaignFile)
				_ = os.Remove(oldPath) // Ignore error if file doesn't exist
			}
		}

		// Save with new name
		err := SaveCampaign(m, newName)
		if err != nil {
			m.ShowRenamePopup = false
			m.SaveInput = ""
			return m, func() tea.Msg {
				return SetErrorMsg{Message: "Failed to rename campaign: " + err.Error()}
			}
		}

		// Update current campaign file and name
		m.CurrentCampaignFile = sanitizeFilename(newName) + ".json"
		m.CurrentCampaignName = newName
		m.LastAutoSave = "Renamed"

		// Close popup
		m.ShowRenamePopup = false
		m.SaveInput = ""

		return m, nil

	case "esc":
		// Cancel rename
		m.ShowRenamePopup = false
		m.SaveInput = ""
		return m, nil

	case "backspace", "ctrl+h":
		// Remove last character
		if len(m.SaveInput) > 0 {
			m.SaveInput = m.SaveInput[:len(m.SaveInput)-1]
		}
		return m, nil

	default:
		// Add character to input
		if len(key) == 1 {
			m.SaveInput += key
		}
		return m, nil
	}
}
