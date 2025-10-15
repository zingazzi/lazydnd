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
func HandleNavigation(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	// Handle help popup scrolling
	if m.ShowHelpPopup {
		return handleHelpPopupInput(m, msg)
	}

	// Handle cast spell popup input (highest priority)
	if m.ShowCastSpellPrompt && m.CastSpellInputMode {
		return handleCastSpellInput(m, msg)
	}

	// Handle multi-target popup input (highest priority)
	if m.ShowMultiTargetPopup {
		return handleMultiTargetPopupInput(m, msg)
	}

	// Handle condition popup input
	if m.ShowConditionPopup {
		return handleConditionPopupInput(m, msg)
	}

	// Handle quick HP popup input
	if m.ShowQuickHPPopup {
		return handleQuickHPInput(m, key)
	}

	// Handle encounter name prompt input
	if m.ShowEncounterPrompt {
		return handleEncounterPromptInput(m, key)
	}

	// Handle encounter generator popup
	if m.EncounterGenerating {
		return handleGeneratorPopupInput(m, msg)
	}

	// Handle action popup input (but not Enter key - that goes to handleEnter)
	if m.ShowActionPopup && key != "enter" {
		return handleActionPopupInput(m, key)
	}

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

	// Handle active spell list navigation
	if m.ActiveSpellListMode && (key == "up" || key == "down") {
		return handleActiveSpellNavigation(m, key)
	}

	// Handle CR filter input (when in CR filter mode)
	if m.MonsterCRFilterMode {
		return handleCRFilterInput(m, key)
	}

	// Handle spell level filter input (when in spell level filter mode)
	if m.SpellLevelFilterMode {
		return handleSpellLevelFilterInput(m, key)
	}

	// Handle encounter builder input first (when in encounter builder panel AND not in other modes)
	// EXCEPT for tab navigation, help keys, and quit keys which should always work globally
	if m.ActivePanel == EncounterBuilder && !m.MonsterSearchMode && !m.SpellSearchMode {
		// Skip encounter builder for global keys - let them go to global handlers
		if key != "tab" && key != "shift+tab" && key != "?" && key != "q" && key != "ctrl+c" {
			return handleEncounterBuilderInput(m, msg)
		}
	}

	// Check if we have a specific handler for this key (like tab for navigation, quit, etc.)
	if handler, exists := keyHandlers[key]; exists {
		return handler(m, msg)
	}

	// Handle default text input
	return handleDefaultInput(m, msg)
}
