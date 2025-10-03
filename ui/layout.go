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
		ActivePanel:         DiceRoller,
		DiceInput:           "",
		DiceResult:          "",
		DiceHistory:         []string{},
		LastDiceCommand:     "",
		InputMode:           false,
		ScrollOffset:        make(map[PanelType]int),
		SpellSearchInput:    "",
		SpellSearchMode:     false,
		SelectedSpell:       nil,
		SpellSuggestions:    []string{},
		SuggestionIndex:     -1,
	}
}

// View renders the main application view with 2x2 panel layout
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// Calculate panel dimensions to fill the screen
	panelWidth := (m.Width - 6) / 2
	panelHeight := (m.Height - 4) / 2

	// Create panels
	panelViews := make([]string, 4)

	for i := 0; i < 4; i++ {
		panelType := PanelType(i)
		content := m.getPanelContent(panelType)

		// Apply scrolling with strict height enforcement
		contentLines := strings.Split(content, "\n")
		scrollOffset := m.ScrollOffset[panelType]

		// Calculate available content height - match original approach
		// Account for: title (1) + borders (2) + padding (2) = 5 lines
		availableHeight := panelHeight - 5
		if availableHeight < 1 {
			availableHeight = 1 // Minimum 1 line of content
		}

		// Always reserve space for scroll indicators to maintain consistent sizing
		contentHeight := availableHeight
		if panelType == m.ActivePanel {
			// Always reserve 2 lines for potential scroll indicators on active panels
			contentHeight = availableHeight - 2
			if contentHeight < 1 {
				contentHeight = 1
			}
		}

		// For spell panel, use even more conservative height to prevent resizing
		if panelType == Spells {
			contentHeight = contentHeight - 2 // Extra buffer for dynamic content
			if contentHeight < 1 {
				contentHeight = 1
			}
		}

		// Ensure scroll offset is within bounds
		maxScroll := len(contentLines) - contentHeight
		if maxScroll < 0 {
			maxScroll = 0
		}
		if scrollOffset > maxScroll {
			m.ScrollOffset[panelType] = maxScroll
			scrollOffset = maxScroll
		}
		if scrollOffset < 0 {
			m.ScrollOffset[panelType] = 0
			scrollOffset = 0
		}

		// Apply scroll offset and strictly limit content
		var scrolledContent string
		var visibleLines []string

		if len(contentLines) <= contentHeight {
			// Content fits, no scrolling needed
			visibleLines = contentLines
		} else {
			// Content needs scrolling - strictly limit to contentHeight
			endIndex := scrollOffset + contentHeight
			if endIndex > len(contentLines) {
				endIndex = len(contentLines)
			}
			visibleLines = contentLines[scrollOffset:endIndex]
		}

		// Ensure we don't exceed the available height
		if len(visibleLines) > contentHeight {
			visibleLines = visibleLines[:contentHeight]
		}

		// Add scroll indicators for active panels to maintain consistent sizing
		if panelType == m.ActivePanel {
			var finalLines []string

			// Always add top indicator space (empty or with indicator)
			if scrollOffset > 0 && len(contentLines) > availableHeight {
				finalLines = append(finalLines, "▲ (more above)")
			} else {
				finalLines = append(finalLines, "") // Empty line to maintain spacing
			}

			// Add visible content
			finalLines = append(finalLines, visibleLines...)

			// Always add bottom indicator space (empty or with indicator)
			if scrollOffset+len(visibleLines) < len(contentLines) && len(contentLines) > availableHeight {
				finalLines = append(finalLines, "▼ (more below)")
			} else {
				finalLines = append(finalLines, "") // Empty line to maintain spacing
			}

			// Ensure total matches available height exactly
			if len(finalLines) > availableHeight {
				finalLines = finalLines[:availableHeight]
			}
			for len(finalLines) < availableHeight {
				finalLines = append(finalLines, "")
			}

			scrolledContent = strings.Join(finalLines, "\n")
		} else {
			// For inactive panels, just use visible lines with padding
			finalLines := visibleLines
			for len(finalLines) < availableHeight {
				finalLines = append(finalLines, "")
			}
			if len(finalLines) > availableHeight {
				finalLines = finalLines[:availableHeight]
			}
			scrolledContent = strings.Join(finalLines, "\n")
		}

		// Ensure content doesn't exceed available height
		finalLines := strings.Split(scrolledContent, "\n")
		if len(finalLines) > availableHeight {
			finalLines = finalLines[:availableHeight]
			scrolledContent = strings.Join(finalLines, "\n")
		}

		// Style the panel - use consistent styling approach
		title := fmt.Sprintf(" %d. %s ", i+1, PanelNames[i])
		titleBar := PanelTitleStyle.Render(title)

		// Always use the same base style, just change border color
		var borderColor string
		if panelType == m.ActivePanel {
			borderColor = "#7D56F4"
		} else {
			borderColor = "#444444"
		}

		panelStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(borderColor)).
			Padding(1, 2).
			Width(panelWidth).
			Height(panelHeight)

		panelContent := titleBar + "\n" + scrolledContent
		panelViews[i] = panelStyle.Render(panelContent)
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
		content = panels.GetDiceRollerContent(m.DiceInput, m.DiceResult, m.DiceHistory, m.LastDiceCommand, m.InputMode, m.ActivePanel == DiceRoller)
	case CharacterSheet:
		content = panels.GetCharacterSheetContent()
	case Spells:
		content = panels.GetSpellsContent(m.SpellSearchInput, m.SelectedSpell, m.SpellSuggestions, m.SuggestionIndex, m.SpellSearchMode, m.ActivePanel == Spells)
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
			if m.LastDiceCommand != "" {
				return "\n" + HelpStyle.Render("Enter: input dice • r: reroll • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
			} else {
				return "\n" + HelpStyle.Render("Enter: input dice • ↑↓: scroll • 1-4/F1-F4: switch • q: quit")
			}
		}
	} else if panelType == Spells {
		if m.SpellSearchMode {
			return "\n" + HelpStyle.Render("Enter: select spell • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch")
		} else {
			return "\n" + HelpStyle.Render("Enter: search spells • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
		}
	} else {
		return "\n" + HelpStyle.Render("↑↓: scroll • 1-4/F1-F4: switch panels • q: quit")
	}
}
