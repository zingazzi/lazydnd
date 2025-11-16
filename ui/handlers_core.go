// ui/handlers/core.go
package ui

// KeyHandler defines a function type for handling key presses
type KeyHandler func(Model, KeyMsg) (Model, Cmd)

// keyHandlers maps key strings to their handler functions
// Uses key constants for type safety
var keyHandlers = map[string]KeyHandler{
	// Quit handlers
	KeyCtrlC: handleQuit,
	KeyQuit:  handleQuit,

	// Navigation handlers
	KeyTab:      handleTab,
	KeyShiftTab: handleShiftTab,
	KeyUp:       handleUp,
	KeyDown:     handleDown,

	// Function key handlers
	KeyF1: handleF1,
	KeyF2: handleF2,
	KeyF3: handleF3,
	KeyF4: handleF4,

	// Number key handlers
	Key1: handleNumber1,
	Key2: handleNumber2,
	Key3: handleNumber3,
	Key4: handleNumber4,
	Key5: handleNumber5,

	// Action handlers
	KeyEnter:     handleEnter,
	KeyEscape:    handleEscape,
	KeyBackspace: handleBackspace,
	KeyCtrlH:     handleBackspace,
	KeySpace:     handleSpace, // Space key is represented as " " not "space"

	// Letter handlers
	KeyR:     handleR,
	KeyShiftR: handleR, // Shift+R for reaction toggle
	KeyP:     handleP,
	KeyM:     handleM,
	KeyE:     handleE,
	KeyI:     handleI,
	KeyH:     handleH,
	KeyShiftH: handleH, // Shift+H for max HP editing
	KeyK:     handleK, // k for AC editing
	KeyA:     handleA,
	KeyD:     handleD,
	KeyL:     handleL,
	KeyShiftL: handleL, // Shift+L for restoring legendary actions
	KeyC:     handleC,
	KeyS:     handleS,
	KeyN:     handleNextTurn,
	KeyX:     handleResetCombat,
	KeyV:     handleV,
	KeyT:     handleT,
	KeyShiftT: handleT, // Shift+T for temp HP
	KeyO:     handleO,
	KeyF:     handleF,

	// Special handlers
	KeyHelp:       handleHelp,
	KeyPlus:       handleQuickAddHP,
	KeyEquals:     handleQuickAddHP,    // = is + without shift
	KeyMinus:      handleQuickRemoveHP,
	KeyUnderscore: handleQuickRemoveHP, // _ is - with shift

	// Save/Load handlers
	KeyCtrlS: handleCtrlS,
	KeyCtrlL: handleCtrlL,
	KeyCtrlN: handleCtrlN,

	// Undo/Redo handlers
	KeyCtrlZ: handleCtrlZ,
	KeyCtrlY: handleCtrlY,
}

// HandleNavigation processes navigation-related key presses
// Uses a prioritized handler chain to route input to appropriate handlers
func HandleNavigation(m Model, msg KeyMsg) (Model, Cmd) {
	chain := NewHandlerChain()
	return chain.Process(m, msg)
}
