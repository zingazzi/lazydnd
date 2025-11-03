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
	"5": handleNumber5,

	// Action handlers
	"enter":     handleEnter,
	"esc":       handleEscape,
	"backspace": handleBackspace,
	"ctrl+h":    handleBackspace,
	" ":         handleSpace, // Space key is represented as " " not "space"

	// Letter handlers
	"r": handleR,
	"R": handleR, // Shift+R for reaction toggle
	"p": handleP,
	"m": handleM,
	"e": handleE,
	"i": handleI,
	"h": handleH,
	"H": handleH, // Shift+H for max HP editing
	"k": handleK, // k for AC editing
	"a": handleA,
	"d": handleD,
	"l": handleL,
	"L": handleL, // Shift+L for restoring legendary actions
	"c": handleC,
	"s": handleS,
	"n": handleNextTurn,
	"x": handleResetCombat,
	"v": handleV,
	"t": handleT,
	"T": handleT, // Shift+T for temp HP
	"o": handleO,
	"f": handleF,

	// Special handlers
	"?": handleHelp,
	"+": handleQuickAddHP,
	"=": handleQuickAddHP,    // = is + without shift
	"-": handleQuickRemoveHP,
	"_": handleQuickRemoveHP, // _ is - with shift

	// Save/Load handlers
	"ctrl+s": handleCtrlS,
	"ctrl+l": handleCtrlL,
	"ctrl+n": handleCtrlN,

	// Undo/Redo handlers
	"ctrl+z": handleCtrlZ,
	"ctrl+y": handleCtrlY,
}

// HandleNavigation processes navigation-related key presses
// Uses a prioritized handler chain to route input to appropriate handlers
func HandleNavigation(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	chain := NewHandlerChain()
	return chain.Process(m, msg)
}
