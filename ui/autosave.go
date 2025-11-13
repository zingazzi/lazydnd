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

// HandleAutoSave performs autosave if conditions are met (exported for TView)
func HandleAutoSave(m Model) Model {
	return handleAutoSave(m)
}

// handleAutoSave performs autosave if conditions are met
func handleAutoSave(m Model) Model {
	// Only autosave if:
	// 1. AutoSave is enabled in config
	// 2. A campaign is loaded
	// 3. Interval time has passed since last autosave
	if !m.Config.AutoSave.Enabled || m.CurrentCampaignName == "" {
		return m
	}

	now := time.Now()
	interval := time.Duration(m.Config.AutoSave.IntervalMinutes) * time.Minute

	if now.Sub(lastAutoSaveTime) < interval {
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
