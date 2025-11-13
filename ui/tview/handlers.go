// ui/tview/handlers.go
package tview

import (
	"lazydnd/ui"

	"github.com/gdamore/tcell/v2"
)

// HandleInput routes input events to the appropriate handler
func HandleInput(model *ui.Model, key tcell.Key, rune rune) bool {
	// Convert TCell key to handler chain format
	_ = convertKeyToString(key, rune)

	// Use the existing handler chain
	// For now, we'll need to adapt the handler chain to work with TView
	// This is a placeholder - will be fully implemented in Phase 4
	return false
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
	case tcell.KeyCtrlH:
		return "ctrl+h"
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
