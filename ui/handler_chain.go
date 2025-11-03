// ui/handler_chain.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Handler defines the interface for all input handlers
type Handler interface {
	CanHandle(m Model, msg tea.KeyMsg) bool
	Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd)
	Priority() int // Lower number = higher priority
	Name() string  // For debugging/logging
}

// HandlerChain is a prioritized list of handlers
type HandlerChain []Handler

// NewHandlerChain creates a new handler chain with all handlers in priority order
func NewHandlerChain() HandlerChain {
	return HandlerChain{
		&HelpPopupHandler{},
		&CastSpellHandler{},
		&MultiTargetHandler{},
		&QuickHPHandler{},
		&SpellFilterHandler{},
		&MonsterFilterHandler{},
		&ConditionPopupHandler{},
		&EncounterPromptHandler{},
		&GeneratorPopupHandler{},
		&ActionPopupHandler{},
		&SavePopupHandler{},
		&LoadPopupHandler{},
		&RenamePopupHandler{},
		&ActiveSpellListHandler{},
		&EncounterBuilderHandler{},
		&GlobalKeyHandler{},
		&DefaultInputHandler{},
	}
}

// Process processes a key message through the handler chain
func (chain HandlerChain) Process(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	for _, handler := range chain {
		if handler.CanHandle(m, msg) {
			return handler.Handle(m, msg)
		}
	}
	// Fallback (should never reach here if DefaultInputHandler is in chain)
	return handleDefaultInput(m, msg)
}

// ========== Priority 1: Help Popup Handler ==========

// HelpPopupHandler handles help popup scrolling (Priority 1)
type HelpPopupHandler struct{}

func (h *HelpPopupHandler) Priority() int { return 1 }
func (h *HelpPopupHandler) Name() string  { return "HelpPopupHandler" }

func (h *HelpPopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowHelp || m.ShowHelpPopup
}

func (h *HelpPopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleHelpPopupInput(m, msg)
}

// ========== Priority 2: Input Mode Handlers ==========

// CastSpellHandler handles cast spell popup input (Priority 2)
type CastSpellHandler struct{}

func (h *CastSpellHandler) Priority() int { return 2 }
func (h *CastSpellHandler) Name() string  { return "CastSpellHandler" }

func (h *CastSpellHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.GetInputMode() == ModeCastSpell
}

func (h *CastSpellHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleCastSpellInput(m, msg)
}

// MultiTargetHandler handles multi-target popup input (Priority 2)
type MultiTargetHandler struct{}

func (h *MultiTargetHandler) Priority() int { return 2 }
func (h *MultiTargetHandler) Name() string  { return "MultiTargetHandler" }

func (h *MultiTargetHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.GetInputMode() == ModeMultiTarget
}

func (h *MultiTargetHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleMultiTargetPopupInput(m, msg)
}

// QuickHPHandler handles quick HP popup input (Priority 2)
type QuickHPHandler struct{}

func (h *QuickHPHandler) Priority() int { return 2 }
func (h *QuickHPHandler) Name() string  { return "QuickHPHandler" }

func (h *QuickHPHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.GetInputMode() == ModeQuickHP
}

func (h *QuickHPHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleQuickHPInput(m, key)
}

// SpellFilterHandler handles spell level filter input (Priority 2)
type SpellFilterHandler struct{}

func (h *SpellFilterHandler) Priority() int { return 2 }
func (h *SpellFilterHandler) Name() string  { return "SpellFilterHandler" }

func (h *SpellFilterHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.GetInputMode() == ModeSpellFilter
}

func (h *SpellFilterHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleSpellLevelFilterInput(m, key)
}

// MonsterFilterHandler handles monster CR filter input (Priority 2)
type MonsterFilterHandler struct{}

func (h *MonsterFilterHandler) Priority() int { return 2 }
func (h *MonsterFilterHandler) Name() string  { return "MonsterFilterHandler" }

func (h *MonsterFilterHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.GetInputMode() == ModeMonsterFilter
}

func (h *MonsterFilterHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleCRFilterInput(m, key)
}

// ========== Priority 3: Popup Handlers ==========

// ConditionPopupHandler handles condition popup input (Priority 3)
type ConditionPopupHandler struct{}

func (h *ConditionPopupHandler) Priority() int { return 3 }
func (h *ConditionPopupHandler) Name() string  { return "ConditionPopupHandler" }

func (h *ConditionPopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowCondition || m.ShowConditionPopup
}

func (h *ConditionPopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleConditionPopupInput(m, msg)
}

// EncounterPromptHandler handles encounter save prompt input (Priority 3)
type EncounterPromptHandler struct{}

func (h *EncounterPromptHandler) Priority() int { return 3 }
func (h *EncounterPromptHandler) Name() string  { return "EncounterPromptHandler" }

func (h *EncounterPromptHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowEncounterPrompt || m.ShowEncounterPrompt
}

func (h *EncounterPromptHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleEncounterPromptInput(m, key)
}

// GeneratorPopupHandler handles encounter generator popup input (Priority 3)
type GeneratorPopupHandler struct{}

func (h *GeneratorPopupHandler) Priority() int { return 3 }
func (h *GeneratorPopupHandler) Name() string  { return "GeneratorPopupHandler" }

func (h *GeneratorPopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Encounter.Generating || m.EncounterGenerating
}

func (h *GeneratorPopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleGeneratorPopupInput(m, msg)
}

// ActionPopupHandler handles action popup input except Enter key (Priority 3)
type ActionPopupHandler struct{}

func (h *ActionPopupHandler) Priority() int { return 3 }
func (h *ActionPopupHandler) Name() string  { return "ActionPopupHandler" }

func (h *ActionPopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	key := msg.String()
	// Action popup handles all keys except Enter (which goes to handleEnter)
	return (m.Popup.ShowAction || m.ShowActionPopup) && key != "enter"
}

func (h *ActionPopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleActionPopupInput(m, key)
}

// SavePopupHandler handles save popup input (Priority 3)
type SavePopupHandler struct{}

func (h *SavePopupHandler) Priority() int { return 3 }
func (h *SavePopupHandler) Name() string  { return "SavePopupHandler" }

func (h *SavePopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowSave || m.ShowSavePopup
}

func (h *SavePopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleSavePopupInput(m, msg)
}

// LoadPopupHandler handles load popup input (Priority 3)
type LoadPopupHandler struct{}

func (h *LoadPopupHandler) Priority() int { return 3 }
func (h *LoadPopupHandler) Name() string  { return "LoadPopupHandler" }

func (h *LoadPopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowLoad || m.ShowLoadPopup
}

func (h *LoadPopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleLoadPopupInput(m, msg)
}

// RenamePopupHandler handles rename popup input (Priority 3)
type RenamePopupHandler struct{}

func (h *RenamePopupHandler) Priority() int { return 3 }
func (h *RenamePopupHandler) Name() string  { return "RenamePopupHandler" }

func (h *RenamePopupHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	return m.Popup.ShowRename || m.ShowRenamePopup
}

func (h *RenamePopupHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleRenamePopupInput(m, msg)
}

// ========== Priority 4: Special Navigation Handlers ==========

// ActiveSpellListHandler handles active spell list navigation (Priority 4)
type ActiveSpellListHandler struct{}

func (h *ActiveSpellListHandler) Priority() int { return 4 }
func (h *ActiveSpellListHandler) Name() string  { return "ActiveSpellListHandler" }

func (h *ActiveSpellListHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	key := msg.String()
	return (m.Spells.ActiveSpellListMode || m.ActiveSpellListMode) && (key == "up" || key == "down")
}

func (h *ActiveSpellListHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	return handleActiveSpellNavigation(m, key)
}

// ========== Priority 5: Panel-Specific Handlers ==========

// EncounterBuilderHandler handles encounter builder panel input (Priority 5)
type EncounterBuilderHandler struct{}

func (h *EncounterBuilderHandler) Priority() int { return 5 }
func (h *EncounterBuilderHandler) Name() string  { return "EncounterBuilderHandler" }

func (h *EncounterBuilderHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	key := msg.String()
	mode := m.GetInputMode()

	// Only handle if in encounter builder panel and not in search modes
	if m.ActivePanel != EncounterBuilder || mode == ModeMonsterSearch || mode == ModeSpellSearch {
		return false
	}

	// Skip encounter builder for global keys - let them go to global handlers
	if key == "tab" || key == "shift+tab" || key == "?" || key == "q" || key == "ctrl+c" {
		return false
	}

	return true
}

func (h *EncounterBuilderHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleEncounterBuilderInput(m, msg)
}

// ========== Priority 6: Key-Based Handlers ==========

// GlobalKeyHandler handles keys from keyHandlers map (Priority 6)
type GlobalKeyHandler struct{}

func (h *GlobalKeyHandler) Priority() int { return 6 }
func (h *GlobalKeyHandler) Name() string  { return "GlobalKeyHandler" }

func (h *GlobalKeyHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	key := msg.String()

	// Check if we have a handler for this key
	_, exists := keyHandlers[key]
	if !exists {
		return false
	}

	// For certain keys like + and -, prioritize input mode over quick HP shortcuts
	// This allows typing formulas like "2d6+5" in the dice roller
	if (key == "+" || key == "=" || key == "-" || key == "_") && m.IsInputMode() {
		return false // Let DefaultInputHandler handle it
	}

	return true
}

func (h *GlobalKeyHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	handler := keyHandlers[key]
	return handler(m, msg)
}

// ========== Priority 7: Default Handler ==========

// DefaultInputHandler handles all remaining text input (Priority 7)
type DefaultInputHandler struct{}

func (h *DefaultInputHandler) Priority() int { return 7 }
func (h *DefaultInputHandler) Name() string  { return "DefaultInputHandler" }

func (h *DefaultInputHandler) CanHandle(m Model, msg tea.KeyMsg) bool {
	// Default handler always handles if reached
	return true
}

func (h *DefaultInputHandler) Handle(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	return handleDefaultInput(m, msg)
}
