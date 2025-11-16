// ui/tview/handlers.go
package tview

import (
	"lazydnd/ui"

	"github.com/gdamore/tcell/v2"
)

// HandleInput routes input events to the appropriate handler
// Returns true if handled, false if not handled, and a special "quit" signal
func HandleInput(model *ui.Model, key tcell.Key, rune rune) (handled bool, shouldQuit bool) {
	// Check for quit keys first
	if key == tcell.KeyCtrlC {
		// Check if we're in input mode
		if !model.InputMode && !model.InitiativeInputMode && !model.SpellSearchMode && !model.MonsterSearchMode {
			return true, true // Handled, should quit
		}
	}
	if key == tcell.KeyRune && rune == 'q' {
		// Check if we're in input mode
		if !model.InputMode && !model.InitiativeInputMode && !model.SpellSearchMode && !model.MonsterSearchMode {
			return true, true // Handled, should quit
		}
	}

	// Convert TCell key to handler chain format
	keyStr := convertKeyToString(key, rune)
	if keyStr == "" {
		return false, false
	}

	// Create a KeyMsg adapter
	keyMsg := ui.NewTViewKeyMsg(keyStr)

	// Use the existing handler chain
	updatedModel, _ := ui.HandleNavigation(*model, keyMsg)
	*model = updatedModel

	// Always return handled=true since we processed the key
	// The handler chain will determine if any action was taken
	return true, false
}

// convertKeyToString converts TCell key events to string format used by handlers
func convertKeyToString(key tcell.Key, rune rune) string {
	// Handle special keys
	switch key {
	case tcell.KeyCtrlC:
		return "ctrl+c"
	case tcell.KeyTab:
		return "tab"
	case tcell.KeyBacktab:
		return "shift+tab"
	case tcell.KeyUp:
		return "up"
	case tcell.KeyDown:
		return "down"
	case tcell.KeyEnter:
		return "enter"
	case tcell.KeyEsc:
		return "esc"
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return "backspace"
	case tcell.KeyCtrlS:
		return "ctrl+s"
	case tcell.KeyCtrlL:
		return "ctrl+l"
	case tcell.KeyCtrlN:
		return "ctrl+n"
	case tcell.KeyCtrlZ:
		return "ctrl+z"
	case tcell.KeyCtrlY:
		return "ctrl+y"
	case tcell.KeyF1:
		return "f1"
	case tcell.KeyF2:
		return "f2"
	case tcell.KeyF3:
		return "f3"
	case tcell.KeyF4:
		return "f4"
	case tcell.KeyRune:
		// Handle regular runes
		switch rune {
		case ' ':
			return " "
		case '?':
			return "?"
		case '+', '=':
			return "+"
		case '-', '_':
			return "-"
		default:
			return string(rune)
		}
	}

	return ""
}
