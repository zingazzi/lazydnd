// ui/layout.go
package ui

import (
	"fmt"
	"strings"

	"lazydnd/panels"

	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates the initial application model
func InitialModel() Model {
	return Model{
		ActivePanel:  DiceRoller,
		DiceInput:    "",
		DiceResult:   "",
		DiceHistory:  []string{},
		InputMode:    false,
		ScrollOffset: make(map[PanelType]int),
	}
}

// View renders the main application view with 2x2 panel layout
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// Calculate panel dimensions
	panelWidth := (m.Width - 6) / 2
	panelHeight := (m.Height - 4) / 2

	// Create panels
	panelViews := make([]string, 4)

	for i := 0; i < 4; i++ {
		panelType := PanelType(i)
		content := m.getPanelContent(panelType)

		// Apply scrolling
		contentLines := strings.Split(content, "\n")
		scrollOffset := m.ScrollOffset[panelType]

		// Calculate available content height (panel height minus title and padding)
		availableHeight := panelHeight - 4 // Account for title, borders, and padding

		// Apply scroll offset
		if scrollOffset > 0 && scrollOffset < len(contentLines) {
			if scrollOffset+availableHeight < len(contentLines) {
				contentLines = contentLines[scrollOffset : scrollOffset+availableHeight]
			} else {
				contentLines = contentLines[scrollOffset:]
			}
		} else if scrollOffset >= len(contentLines) {
			// Reset scroll if we've gone too far
			m.ScrollOffset[panelType] = 0
		} else {
			// Show from beginning
			if len(contentLines) > availableHeight {
				contentLines = contentLines[:availableHeight]
			}
		}

		scrolledContent := strings.Join(contentLines, "\n")

		// Style the panel
		title := fmt.Sprintf(" %d. %s ", i+1, PanelNames[i])
		titleBar := PanelTitleStyle.Render(title)

		var panelStyle lipgloss.Style
		if panelType == m.ActivePanel {
			panelStyle = ActivePanelStyle
		} else {
			panelStyle = InactivePanelStyle
		}

		panelContent := titleBar + "\n" + scrolledContent
		panelViews[i] = panelStyle.Width(panelWidth).Height(panelHeight).Render(panelContent)
	}

	// Arrange panels in 2x2 grid
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[0], panelViews[1])
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, panelViews[2], panelViews[3])

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

// getPanelContent returns the content for a specific panel
func (m Model) getPanelContent(panelType PanelType) string {
	var content string

	switch panelType {
	case DiceRoller:
		content = panels.GetDiceRollerContent(m.DiceInput, m.DiceResult, m.DiceHistory, m.InputMode, m.ActivePanel == DiceRoller)
	case CharacterSheet:
		content = panels.GetCharacterSheetContent()
	case Spells:
		content = panels.GetSpellsContent()
	case CampaignNotes:
		content = panels.GetCampaignNotesContent()
	}

	// Add help text for active panel
	if panelType == m.ActivePanel {
		content += m.getHelpText(panelType)
	}

	return content
}

// getHelpText returns context-sensitive help text
func (m Model) getHelpText(panelType PanelType) string {
	if panelType == DiceRoller {
		if m.InputMode {
			return "\n" + HelpStyle.Render("Enter: roll • Esc: cancel • F1-F4: switch panels")
		} else {
			return "\n" + HelpStyle.Render("Enter: input dice • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
		}
	} else {
		return "\n" + HelpStyle.Render("↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
	}
}
