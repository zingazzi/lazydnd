// ui/autosave.go
package ui

import (
	"fmt"
	"time"
)

var (
	lastAutoSaveTime time.Time
	autoSaveCounter  int
)

// handleAutoSave performs autosave if conditions are met
func handleAutoSave(m Model) Model {
	// Only autosave if:
	// 1. AutoSave is enabled
	// 2. A campaign is loaded
	// 3. At least 5 minutes have passed since last autosave
	if !m.AutoSaveEnabled || m.CurrentCampaignName == "" {
		return m
	}

	now := time.Now()
	if now.Sub(lastAutoSaveTime) < 5*time.Minute {
		// Update time display
		elapsed := now.Sub(lastAutoSaveTime)
		minutes := int(elapsed.Minutes())
		if minutes == 0 {
			m.LastAutoSave = "Just now"
		} else {
			m.LastAutoSave = fmt.Sprintf("%dm ago", minutes)
		}
		return m
	}

	// Perform autosave
	err := SaveCampaign(m, m.CurrentCampaignName)
	if err == nil {
		lastAutoSaveTime = now
		autoSaveCounter++
		m.LastAutoSave = "Just now"
	}

	return m
}
