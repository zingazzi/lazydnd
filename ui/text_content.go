// ui/text_content.go
package ui

// Application text content centralized for easy maintenance and localization

// ========== APPLICATION INFO ==========

const (
	AppName      = "LazyDnD"
	AppNameEmoji = "🎲 LazyDnD"
)

// ========== STATUS BAR TEXT ==========

type StatusBarText struct {
	ProjectName     string
	TabKey          string
	TabDesc         string
	ArrowKeys       string
	ArrowDesc       string
	NumbersKey      string
	NumbersDesc     string
	HelpKey         string
	HelpDesc        string
	QuitKey         string
	QuitDesc        string
}

var DefaultStatusBarText = StatusBarText{
	ProjectName: AppNameEmoji,
	TabKey:      "Tab",
	TabDesc:     "Switch Panel",
	ArrowKeys:   "↑↓←→",
	ArrowDesc:   "Navigate",
	NumbersKey:  "1-4",
	NumbersDesc: "Quick Switch",
	HelpKey:     "?",
	HelpDesc:    "Help",
	QuitKey:     "q",
	QuitDesc:    "Quit",
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
	{"Esc", "Cancel/Exit current mode"},
	{"?", "Toggle this help"},
	{"q", "Quit application"},
}

// ========== PANEL-SPECIFIC HELP TEXT ==========

// DiceRollerHelp contains help text for the Dice Roller panel
var DiceRollerHelp = []HelpKey{
	{"Enter", "Start/confirm dice input"},
	{"r", "Reroll last command"},
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
	{"e", "Enter edit mode"},
	{"", ""},
	{"In Edit Mode:", ""},
	{"  ↑/↓", "Select entry"},
	{"  i", "Edit initiative value"},
	{"  h", "Edit HP (monsters only)"},
	{"  l", "View linked monster details"},
	{"  a", "Show monster actions (if linked)"},
	{"  c", "Copy/duplicate entry"},
	{"  d", "Delete entry"},
}

// SpellsHelp contains help text for the Spells panel
var SpellsHelp = []HelpKey{
	{"Enter", "Start spell search"},
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
func DiceRollerInlineHelp(inputMode bool, hasLastCommand bool) string {
	if inputMode {
		return "Enter: roll • Esc: cancel • F1-F4: switch panels"
	}
	if hasLastCommand {
		return "Enter: input dice • r: reroll • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
	}
	return "Enter: input dice • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
}

// InitiativeTrackerInlineHelp returns inline help for initiative tracker panel
func InitiativeTrackerInlineHelp(editMode, inputMode, listMode bool) string {
	if editMode {
		return "Enter: confirm edit • Esc: cancel • F1-F4: switch panels"
	}
	if inputMode {
		return "Enter: confirm • Esc: cancel • F1-F4: switch panels"
	}
	if listMode {
		return "↑↓: select • i: init • h: HP • l: view • a: actions • c: copy • d: delete • Esc: exit"
	}
	return "p: add player • m: add monster • e: edit list • ↑↓: scroll • 1-4/F1-F4: switch • q: quit"
}

// SpellsInlineHelp returns inline help for spells panel
func SpellsInlineHelp(searchMode bool) string {
	if searchMode {
		return "Enter: select spell • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch"
	}
	return "Enter: search spells • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}

// MonstersInlineHelp returns inline help for monsters panel
func MonstersInlineHelp(searchMode bool) string {
	if searchMode {
		return "Enter: select monster • ↑↓: navigate suggestions • Esc: cancel • F1-F4: switch"
	}
	return "Enter: search monsters • ↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}

// DefaultInlineHelp returns default inline help
func DefaultInlineHelp() string {
	return "↑↓: scroll • 1-4/F1-F4: switch panels • q: quit"
}
