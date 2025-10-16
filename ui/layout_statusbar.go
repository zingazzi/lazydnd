// ui/layout_statusbar.go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// getPanelCommands returns key shortcuts for the active panel
func (m Model) getPanelCommands() string {
	switch m.ActivePanel {
	case DiceRoller:
		if m.InputMode {
			return "Enter:roll Esc:cancel"
		}
		return "Enter:input r:reroll h:history"
	case InitiativeTracker:
		if m.InitiativeInputMode || m.InitiativeEditMode {
			return "Enter:confirm Esc:cancel"
		}
		if m.MultiTargetMode {
			return "Space:select t:exit Enter:apply"
		}
		if m.InitiativeListMode {
			return "p:player m:monster n:next i:init h:HP d:del"
		}
		return "p:player m:monster n:next"
	case Spells:
		if m.SpellSearchMode {
			return "↑↓:select Enter:view Esc:cancel"
		}
		if m.ActiveSpellListMode {
			return "v:view c:cast d:delete Esc:exit"
		}
		return "Enter:search f:filter v:active c:cast"
	case Monsters:
		if m.MonsterSearchMode {
			return "↑↓:select Enter:view Esc:cancel"
		}
		if m.MonsterCRFilterMode {
			return "Enter:filter Esc:cancel"
		}
		return "Enter:search f:CR a:add"
	case Notes:
		if m.NotesEditMode {
			return "Enter:newline Esc:save"
		}
		if m.NotesSearchMode {
			return "Enter:search Esc:exit"
		}
		return "e:edit f:search"
	case EncounterBuilder:
		// Show different keys based on current mode
		switch m.EncounterBuilderMode {
		case "party_setup":
			return "1-9:size Shift+1-9:level b:building"
		case "building":
			if m.ShowEncounterPrompt {
				return "Enter:confirm Esc:cancel"
			}
			return "m:monster +/-:qty g:generate s:save t:templates l:deploy c:clear"
		case "templates":
			return "↑↓:select Enter:view Delete:delete b:back"
		case "template_detail":
			return "l:load d:delete Esc/b:back"
		default:
			return "b:building t:templates"
		}
	default:
		return ""
	}
}

// renderStatusBar renders the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	var result string

	// Render error banner if there's an error
	if m.ErrorVisible && m.ErrorMessage != "" {
		errorBanner := m.Styles.ErrorStyle.
			Width(m.Width).
			Background(lipgloss.Color("#8B0000")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1).
			Render("❌ " + m.ErrorMessage)
		result = errorBanner + "\n"
	}

	// LEFT SECTION: LazyD&D info
	var leftParts []string
	leftParts = append(leftParts, "🎲LazyDnD")
	if m.CurrentCampaignName != "" {
		leftParts = append(leftParts, "│", m.CurrentCampaignName)
	}
	leftSection := m.Styles.StatusBarTextStyle.Render(strings.Join(leftParts, " "))

	// Check if help hints are enabled in config
	showHelp := m.Config.Display.ShowHelpHints

	// MIDDLE SECTION: Panel-specific commands (if help enabled)
	var middleSection string
	if showHelp {
		panelName := PanelNames[m.ActivePanel]
		panelCommands := m.getPanelCommands()
		middleParts := []string{panelName}
		if panelCommands != "" {
			middleParts = append(middleParts, "│", panelCommands)
		}
		middleSection = m.Styles.StatusBarKeyStyle.Render(strings.Join(middleParts, " "))
	} else {
		// Just show panel name without commands
		middleSection = m.Styles.StatusBarTextStyle.Render(PanelNames[m.ActivePanel])
	}

	// RIGHT SECTION: Shared commands (if help enabled)
	var rightSection string
	if showHelp {
		rightSection = m.Styles.StatusBarTextStyle.Render("Tab:switch │ ?:help │ Ctrl+S:save │ q:quit")
	}

	// Calculate spacing for sections
	leftWidth := lipgloss.Width(leftSection)
	middleWidth := lipgloss.Width(middleSection)
	rightWidth := lipgloss.Width(rightSection)

	var statusBarContent string
	if showHelp && rightWidth > 0 {
		// Three sections with help
		totalContentWidth := leftWidth + middleWidth + rightWidth
		availableSpace := m.Width - totalContentWidth - 4 // Account for separators

		// Minimal spacing between sections
		spacing1 := 2
		spacing2 := availableSpace - spacing1
		if spacing2 < 2 {
			spacing2 = 2
		}

		statusBarContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftSection,
			strings.Repeat(" ", spacing1),
			middleSection,
			strings.Repeat(" ", spacing2),
			rightSection,
		)
	} else {
		// Two sections without help
		totalContentWidth := leftWidth + middleWidth
		availableSpace := m.Width - totalContentWidth - 2
		spacing := availableSpace / 2
		if spacing < 2 {
			spacing = 2
		}

		statusBarContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftSection,
			strings.Repeat(" ", spacing),
			middleSection,
		)
	}

	// Apply full-width background style
	statusBar := m.Styles.StatusBarStyle.
		Width(m.Width).
		Render(statusBarContent)

	result += statusBar
	return result
}
