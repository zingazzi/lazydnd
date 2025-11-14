// ui/action_popup.go
package ui

import (
	"fmt"
	"strings"
)

const (
	ActionPopupTitle = "🗡️  Monster Actions"
)

// renderActionPopupOverlay is deprecated - TView handles overlays
func (m Model) renderActionPopupOverlay(mainView string) string {
	return mainView
}

// buildActionContent builds the action popup content (plain text for TView)
func (m Model) buildActionContent() string {
	var content strings.Builder

	content.WriteString(ActionPopupTitle)
	content.WriteString("\n\n")

	if len(m.ActionPopupActions) == 0 {
		content.WriteString("No actions available")
		return content.String()
	}

	for i, action := range m.ActionPopupActions {
		prefix := "  "
		if i == m.ActionPopupIndex {
			prefix = "▶ "
		}

		content.WriteString(fmt.Sprintf("%s%s\n", prefix, action.Name))
		if action.Description != "" {
			content.WriteString(fmt.Sprintf("    %s\n", action.Description))
		}
		content.WriteString("\n")
	}

	content.WriteString("\nEnter: Select  |  Esc: Cancel")
	return content.String()
}

// RenderActionPopup renders the action popup (plain text for TView)
func RenderActionPopup(m Model) string {
	return m.buildActionContent()
}
