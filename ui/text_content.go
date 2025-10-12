// ui/text_content.go
package ui

// Application text content centralized for easy maintenance and localization

// ========== APPLICATION INFO ==========

const (
	AppName      = "LazyDnD"
	AppNameEmoji = "🎲 LazyDnD"
	AppVersion   = "v2.4.0"
)

// ========== STATUS BAR TEXT ==========

type StatusBarText struct {
	ProjectName string
	TabKey      string
	TabDesc     string
	ArrowKeys   string
	ArrowDesc   string
	NumbersKey  string
	NumbersDesc string
	HelpKey     string
	HelpDesc    string
	QuitKey     string
	QuitDesc    string
}

var DefaultStatusBarText = StatusBarText{
	ProjectName: AppNameEmoji,
	TabKey:      "Tab",
	TabDesc:     "Switch",
	ArrowKeys:   "↑↓",
	ArrowDesc:   "Navigate",
	NumbersKey:  "1-4",
	NumbersDesc: "Panels",
	HelpKey:     "?",
	HelpDesc:    "Help",
	QuitKey:     "Ctrl+S/L",
	QuitDesc:    "Save/Load",
}

// ========== HELP POPUP TEXT ==========

const (
	HelpPopupTitle  = "🎲 LazyDnD - Help"
	HelpPopupFooter = "Press ? or Esc to close this help"
)

// CommonNavigationKeys defines the common navigation keys shown in help
type HelpKey struct {
	Key         string
	Description string
}

var CommonNavigationKeys = []HelpKey{
	{"Tab", "Switch to next panel"},
	{"1-4", "Quick switch to panel"},
	{"F1-F4", "Switch to specific panel"},
	{"↑/↓", "Scroll panel content"},
	{"Ctrl+S", "Save campaign"},
	{"Ctrl+L", "Load campaign"},
	{"Ctrl+N", "Rename campaign"},
	{"Esc", "Cancel/Exit current mode"},
	{"?", "Toggle this help"},
	{"q", "Quit application"},
}

// ========== PANEL-SPECIFIC HELP TEXT ==========

// DiceRollerHelp contains help text for the Dice Roller panel
var DiceRollerHelp = []HelpKey{
	{"Enter", "Start/confirm dice input"},
	{"r", "Reroll last command"},
	{"h", "Open history to select any roll"},
	{"Examples:", ""},
	{"  2d6", "Roll 2 six-sided dice"},
	{"  1d20+5", "Roll d20 with +5 modifier"},
	{"  adv", "Roll with advantage"},
	{"  dis", "Roll with disadvantage"},
}

// InitiativeTrackerHelp contains help text for the Initiative Tracker panel
var InitiativeTrackerHelp = []HelpKey{
	{"p", "Add player to initiative"},
	{"m", "Add monster to initiative"},
	{"Enter", "Enter edit mode"},
	{"n", "Next turn (works in all modes)"},
	{"x", "Reset combat (works in all modes)"},
	{"Ctrl+Z", "Undo HP change (up to 3 actions)"},
	{"Ctrl+Y", "Redo HP change"},
	{"", ""},
	{"In Edit Mode:", ""},
	{"  ↑/↓", "Select entry"},
	{"  i", "Edit initiative value"},
	{"  h", "Edit HP (monsters only)"},
	{"  s", "Roll saving throws (monsters only)"},
	{"  l", "View linked monster details"},
	{"  a", "Show monster actions (if linked)"},
	{"  c", "Copy/duplicate entry"},
	{"  d", "Delete entry"},
}

// SpellsHelp contains help text for the Spells panel
var SpellsHelp = []HelpKey{
	{"Enter", "Start spell search"},
	{"f", "Filter by spell level (0-9, ranges)"},
	{"c", "Cast selected spell"},
	{"v", "View active spells"},
	{"", ""},
	{"In Search Mode:", ""},
	{"  Type", "Search for spells"},
	{"  ↑/↓", "Navigate suggestions"},
	{"  Enter", "Select spell"},
	{"  Backspace", "Delete character"},
}

// MonstersHelp contains help text for the Monsters panel
var MonstersHelp = []HelpKey{
	{"Enter", "Start monster search"},
	{"f", "Filter by CR (exact, range, min)"},
	{"a", "Add to initiative tracker"},
	{"", ""},
	{"In Search Mode:", ""},
	{"  Type", "Search for monsters"},
	{"  ↑/↓", "Navigate suggestions"},
	{"  Enter", "Select monster"},
	{"  Backspace", "Delete character"},
}

// GetPanelHelpKeys returns the help keys for a specific panel
func GetPanelHelpKeys(panelType PanelType) []HelpKey {
	switch panelType {
	case DiceRoller:
		return DiceRollerHelp
	case InitiativeTracker:
		return InitiativeTrackerHelp
	case Spells:
		return SpellsHelp
	case Monsters:
		return MonstersHelp
	default:
		return []HelpKey{}
	}
}

// ========== HELP TEXT (INLINE) ==========

// InlineHelpText contains the inline help text shown at the bottom of panels

// DiceRollerInlineHelp returns inline help for dice roller panel
func DiceRollerInlineHelp(inputMode bool, hasLastCommand bool, historyMode bool) string {
	if historyMode {
		return "↑↓: select • Enter: re-roll • Esc: exit history • F1-F4: switch panels"
	}
	if inputMode {
		return "Enter: roll • Esc: cancel • F1-F4: switch panels"
	}
	if hasLastCommand {
		return "Enter: input dice • r: reroll • h: history • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
	}
	return "Enter: input dice • h: history • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
}

// InitiativeTrackerInlineHelp returns inline help for initiative tracker panel
func InitiativeTrackerInlineHelp(editMode, inputMode, listMode, multiTargetMode bool) string {
	if multiTargetMode {
		return "Space: select/deselect • Enter: apply damage/healing • t: exit multi-target • ↑↓: navigate"
	}
	if editMode {
		return "Enter: confirm edit • Esc: cancel • F1-F4: switch panels"
	}
	if inputMode {
		return "Enter: confirm • Esc: cancel • F1-F4: switch panels"
	}
	if listMode {
		return "↑↓: select • i: init • h: HP • o: conditions • s: saves • l: view • a: actions • c: copy • d: delete • t: multi-target • Esc: exit"
	}
	return "p: add player • m: add monster • e: edit list • n: next turn • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
}

// SpellsInlineHelp returns inline help for spells panel
func SpellsInlineHelp(m Model) string {
	if m.CastSpellInputMode {
		return "Type caster name • Enter: confirm • Esc: cancel"
	}
	if m.SpellLevelFilterMode {
		return "Type level (0-9, 1-3, 5+) • ↑↓: navigate • Enter: select • Esc: cancel"
	}
	if m.SpellSearchMode {
		return "Enter: select spell • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch"
	}
	if m.ActiveSpellListMode {
		return "↑↓: select spell • d: delete spell • Esc: back to spells • F1-F4: switch panels"
	}
	if m.SelectedSpell != nil {
		return "c: Cast Spell • Enter: Search • f: Filter by Level • Esc: Back • v: View Active • 1-4/F1-F4: Switch"
	}
	if len(m.ActiveSpells) > 0 {
		return "Enter: search spell • f: filter by level • v: view active spells • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
	}
	return "Enter: search spell • f: filter by level • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}

// MonstersInlineHelp returns inline help for monsters panel
func MonstersInlineHelp(m Model) string {
	if m.MonsterCRFilterMode {
		return "Type CR (0-5, 10+, or 5) • ↑↓: navigate • Enter: select • Esc: cancel"
	}
	if m.MonsterSearchMode {
		return "Enter: select monster • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch"
	}
	if m.SelectedMonster != nil {
		return "a: Add to Initiative • Enter: Search • f: Filter by CR • Esc: Back"
	}
	return "Enter: search by name • f: filter by CR • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}

// DefaultInlineHelp returns default inline help
func DefaultInlineHelp() string {
	return "↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}
