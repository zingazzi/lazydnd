// ui/cast_spell_popup.go
package ui

import (
	"fmt"
)

// RenderCastSpellPopup renders the popup for casting a spell (plain text for TView)
func RenderCastSpellPopup(spellName, casterInput string) string {
	var content string

	content += "💫 Cast Spell\n\n"
	content += fmt.Sprintf("Spell: %s\n\n", spellName)
	content += "Enter caster name:\n"
	content += fmt.Sprintf("[%s█]\n\n", casterInput)
	content += "Press Enter to cast • Esc to cancel"

	return content
}

// renderCastSpellPopupOverlay is deprecated - TView handles overlays
func renderCastSpellPopupOverlay(m Model, baseView string) string {
	return baseView
}
