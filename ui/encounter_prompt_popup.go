// ui/encounter_prompt_popup.go
package ui

import (
	"fmt"
)

// RenderEncounterPromptPopup renders the save encounter prompt popup (plain text for TView)
func RenderEncounterPromptPopup(m Model) string {
	var content string
	
	content += "💾 Save Encounter\n\n"
	content += "Name:\n"
	content += fmt.Sprintf("[%s█]\n\n", m.EncounterNameInput)
	content += "[enter] save  [esc] cancel"
	
	return content
}

// renderEncounterPromptOverlay is deprecated - TView handles overlays
func (m Model) renderEncounterPromptOverlay(mainView string) string {
	return mainView
}
