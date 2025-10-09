// ui/handlers/core.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// KeyHandler defines a function type for handling key presses
type KeyHandler func(Model, tea.KeyMsg) (Model, tea.Cmd)

// keyHandlers maps key strings to their handler functions
var keyHandlers = map[string]KeyHandler{
	// Quit handlers
	"ctrl+c": handleQuit,
	"q":      handleQuit,

	// Navigation handlers
	"tab":       handleTab,
	"shift+tab": handleShiftTab,
	"up":        handleUp,
	"down":      handleDown,

	// Function key handlers
	"f1": handleF1,
	"f2": handleF2,
	"f3": handleF3,
	"f4": handleF4,

	// Number key handlers
	"1": handleNumber1,
	"2": handleNumber2,
	"3": handleNumber3,
	"4": handleNumber4,

	// Action handlers
	"enter":     handleEnter,
	"esc":       handleEscape,
	"backspace": handleBackspace,
	"ctrl+h":    handleBackspace,
	"space":     handleSpace,

	// Letter handlers
	"r": handleR,
	"p": handleP,
	"m": handleM,
	"e": handleE,
	"i": handleI,
	"h": handleH,
	"a": handleA,
	"d": handleD,
	"l": handleL,
	"c": handleC,
	"n": handleNextTurn,
	"x": handleResetCombat,

	// Special handlers
	"?": handleHelp,

	// Save/Load handlers
	"ctrl+s": handleCtrlS,
	"ctrl+l": handleCtrlL,
	"ctrl+n": handleCtrlN,
}

// HandleNavigation processes navigation-related key presses
func HandleNavigation(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	// Handle save popup input
	if m.ShowSavePopup {
		return handleSavePopupInput(m, msg)
	}

	// Handle load popup input
	if m.ShowLoadPopup {
		return handleLoadPopupInput(m, msg)
	}

	// Handle rename popup input
	if m.ShowRenamePopup {
		return handleRenamePopupInput(m, msg)
	}

	// Check if we have a specific handler for this key
	if handler, exists := keyHandlers[key]; exists {
		return handler(m, msg)
	}

	// Handle default text input
	return handleDefaultInput(m, msg)
}
