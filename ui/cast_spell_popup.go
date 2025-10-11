// ui/cast_spell_popup.go
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderCastSpellPopup renders the popup for casting a spell
func RenderCastSpellPopup(spellName, casterInput string) string {
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(50).
		Background(lipgloss.Color("#1a1a1a"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center).
		Width(46)

	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(44).
		Foreground(lipgloss.Color("#FAFAFA"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Align(lipgloss.Center).
		Width(46)

	content := titleStyle.Render("💫 Cast Spell") + "\n\n"
	content += "Spell: " + spellName + "\n\n"
	content += "Enter caster name:\n"
	content += inputBoxStyle.Render(casterInput+"│") + "\n\n"
	content += helpStyle.Render("Press Enter to cast • Esc to cancel")

	return popupStyle.Render(content)
}

// renderCastSpellPopupOverlay renders the cast spell popup as an overlay
func renderCastSpellPopupOverlay(m Model, baseView string) string {
	if !m.ShowCastSpellPrompt || m.SpellToCast == nil {
		return baseView
	}

	popup := RenderCastSpellPopup(m.SpellToCast.Name, m.CastSpellInput)
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
	)
}
