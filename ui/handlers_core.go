// ui/handlers_core.go
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

	// Special handlers
	"?": handleHelp,
}

// HandleNavigation processes navigation-related key presses
func (m Model) HandleNavigation(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	// Check if we have a specific handler for this key
	if handler, exists := keyHandlers[key]; exists {
		return handler(m, msg)
	}

	// Handle default text input
	return handleDefaultInput(m, msg)
}
