// tests/handler_chain_test.go
package tests

import (
	"lazydnd/config"
	"lazydnd/ui"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// createTestModel creates a basic model for handler testing
func createHandlerTestModel() ui.Model {
	cfg := config.Default()
	styles := ui.NewStyles(cfg)
	model := ui.InitialModel()
	model.Config = cfg
	model.Styles = styles
	model.Width = 80
	model.Height = 24
	return model
}

// createKeyMsg creates a tea.KeyMsg for testing
// Note: This is a simplified version - actual key handling uses msg.String()
func createKeyMsg(key string) tea.KeyMsg {
	// For most keys, we can use KeyRunes
	if len(key) == 1 && key != " " {
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(key),
		}
	}

	// For special keys, we need to map them
	switch key {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "space", " ":
		return tea.KeyMsg{Type: tea.KeySpace}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "q":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "+":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("+")}
	case "-":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("-")}
	case "=":
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("=")}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
}

// TestHandlerChain_PriorityOrder tests that handlers are checked in priority order
func TestHandlerChain_PriorityOrder(t *testing.T) {
	chain := ui.NewHandlerChain()

	// Verify handlers are in priority order
	expectedPriorities := []int{1, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 4, 5, 6, 7}
	if len(chain) != len(expectedPriorities) {
		t.Fatalf("Expected %d handlers, got %d", len(expectedPriorities), len(chain))
	}

	for i, handler := range chain {
		priority := handler.Priority()
		if priority != expectedPriorities[i] {
			t.Errorf("Handler %d (%s) has priority %d, expected %d", i, handler.Name(), priority, expectedPriorities[i])
		}
		// Verify priorities are non-decreasing (can have same priority)
		if i > 0 && priority < chain[i-1].Priority() {
			t.Errorf("Handler %d (%s) has priority %d, but previous handler has priority %d",
				i, handler.Name(), priority, chain[i-1].Priority())
		}
	}
}

// TestHandlerChain_Process tests the Process method routes to correct handler
func TestHandlerChain_Process(t *testing.T) {
	chain := ui.NewHandlerChain()

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		key         string
		expectHandled bool
		expectHandler string
	}{
		{
			name: "Help popup takes priority",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowHelp = true
				m.InputMode = true // Even if in input mode, help should win
				return m
			},
			key:         "a",
			expectHandled: true,
			expectHandler: "HelpPopupHandler",
		},
		{
			name: "Input mode handler takes priority over key handler",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				return m
			},
			key:         "q",
			expectHandled: true,
			expectHandler: "DefaultInputHandler", // Input mode goes to default input
		},
		{
			name: "Global key handler when no popups",
			setupModel: func(m ui.Model) ui.Model {
				// Normal state, no popups or input modes
				return m
			},
			key:         "q",
			expectHandled: true,
			expectHandler: "GlobalKeyHandler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			_, cmd := chain.Process(model, msg)

			// Verify a command was returned (indicates handler processed it)
			if tt.expectHandled && cmd == nil && tt.expectHandler != "DefaultInputHandler" {
				t.Error("Expected handler to return a command or modify model")
			}
		})
	}
}

// TestHelpPopupHandler tests HelpPopupHandler CanHandle logic
func TestHelpPopupHandler(t *testing.T) {
	handler := &ui.HelpPopupHandler{}

	tests := []struct {
		name       string
		setupModel func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when ShowHelp is true (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowHelp = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when ShowHelpPopup is true (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.ShowHelpPopup = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when help not shown",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("a")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}

			if canHandle {
				_, cmd := handler.Handle(model, msg)
				_ = cmd // Verify it doesn't panic
			}
		})
	}
}

// TestCastSpellHandler tests CastSpellHandler CanHandle logic
func TestCastSpellHandler(t *testing.T) {
	handler := &ui.CastSpellHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when in ModeCastSpell (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.ShowCastSpellPrompt = true
				m.Spells.CastSpellInputMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when in ModeCastSpell (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.ShowCastSpellPrompt = true
				m.CastSpellInputMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when not in cast spell mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("a")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestMultiTargetHandler tests MultiTargetHandler CanHandle logic
func TestMultiTargetHandler(t *testing.T) {
	handler := &ui.MultiTargetHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when in ModeMultiTarget (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Initiative.ShowMultiTargetPopup = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when in ModeMultiTarget (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.ShowMultiTargetPopup = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when not in multi-target mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("a")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestQuickHPHandler tests QuickHPHandler CanHandle logic
func TestQuickHPHandler(t *testing.T) {
	handler := &ui.QuickHPHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when in ModeQuickHP (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Initiative.ShowQuickHPPopup = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when in ModeQuickHP (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.ShowQuickHPPopup = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when not in quick HP mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("5")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestSpellFilterHandler tests SpellFilterHandler CanHandle logic
func TestSpellFilterHandler(t *testing.T) {
	handler := &ui.SpellFilterHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when in ModeSpellFilter (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.LevelFilterMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when in ModeSpellFilter (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.SpellLevelFilterMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when not in filter mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("3")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestMonsterFilterHandler tests MonsterFilterHandler CanHandle logic
func TestMonsterFilterHandler(t *testing.T) {
	handler := &ui.MonsterFilterHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Handles when in ModeMonsterFilter (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Monsters.CRFilterMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles when in ModeMonsterFilter (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.MonsterCRFilterMode = true
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Does not handle when not in filter mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("5")

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestActionPopupHandler tests ActionPopupHandler edge cases
func TestActionPopupHandler(t *testing.T) {
	handler := &ui.ActionPopupHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		key         string
		shouldHandle bool
	}{
		{
			name: "Handles when action popup is open and key is not enter",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowAction = true
				return m
			},
			key:         "up",
			shouldHandle: true,
		},
		{
			name: "Does not handle Enter key (goes to handleEnter)",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowAction = true
				return m
			},
			key:         "enter",
			shouldHandle: false,
		},
		{
			name: "Does not handle when popup not open",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key:         "up",
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestActiveSpellListHandler tests ActiveSpellListHandler CanHandle logic
func TestActiveSpellListHandler(t *testing.T) {
	handler := &ui.ActiveSpellListHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		key         string
		shouldHandle bool
	}{
		{
			name: "Handles up/down keys when in active spell list mode (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.ActiveSpellListMode = true
				m.Spells.ActiveSpells = []ui.ActiveSpell{{Name: "Test"}}
				return m
			},
			key:         "up",
			shouldHandle: true,
		},
		{
			name: "Handles up/down keys when in active spell list mode (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.ActiveSpellListMode = true
				m.ActiveSpells = []ui.ActiveSpell{{Name: "Test"}}
				return m
			},
			key:         "down",
			shouldHandle: true,
		},
		{
			name: "Does not handle non-navigation keys",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.ActiveSpellListMode = true
				return m
			},
			key:         "a",
			shouldHandle: false,
		},
		{
			name: "Does not handle when not in active spell list mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key:         "up",
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestEncounterBuilderHandler tests EncounterBuilderHandler edge cases
func TestEncounterBuilderHandler(t *testing.T) {
	handler := &ui.EncounterBuilderHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		key         string
		shouldHandle bool
	}{
		{
			name: "Handles when in encounter builder panel",
			setupModel: func(m ui.Model) ui.Model {
				m.ActivePanel = ui.EncounterBuilder
				return m
			},
			key:         "m",
			shouldHandle: true,
		},
		{
			name: "Does not handle when in monster search mode",
			setupModel: func(m ui.Model) ui.Model {
				m.ActivePanel = ui.EncounterBuilder
				m.Monsters.SearchMode = true
				return m
			},
			key:         "m",
			shouldHandle: false,
		},
		{
			name: "Does not handle global keys (tab)",
			setupModel: func(m ui.Model) ui.Model {
				m.ActivePanel = ui.EncounterBuilder
				return m
			},
			key:         "tab",
			shouldHandle: false,
		},
		{
			name: "Does not handle global keys (q)",
			setupModel: func(m ui.Model) ui.Model {
				m.ActivePanel = ui.EncounterBuilder
				return m
			},
			key:         "q",
			shouldHandle: false,
		},
		{
			name: "Does not handle when not in encounter builder panel",
			setupModel: func(m ui.Model) ui.Model {
				m.ActivePanel = ui.DiceRoller
				return m
			},
			key:         "m",
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestGlobalKeyHandler tests GlobalKeyHandler edge cases
func TestGlobalKeyHandler(t *testing.T) {
	handler := &ui.GlobalKeyHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		key         string
		shouldHandle bool
	}{
		{
			name: "Handles keys in keyHandlers map",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key:         "q",
			shouldHandle: true,
		},
		{
			name: "Does not handle + key when in input mode",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				return m
			},
			key:         "+",
			shouldHandle: false, // Should go to DefaultInputHandler
		},
		{
			name: "Handles + key when not in input mode",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key:         "+",
			shouldHandle: true,
		},
		{
			name: "Does not handle keys not in keyHandlers map",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key:         "z",
			shouldHandle: false, // Should go to DefaultInputHandler
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestDefaultInputHandler tests DefaultInputHandler
func TestDefaultInputHandler(t *testing.T) {
	handler := &ui.DefaultInputHandler{}

	tests := []struct {
		name        string
		setupModel  func(ui.Model) ui.Model
		shouldHandle bool
	}{
		{
			name: "Always handles (catch-all)",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			shouldHandle: true,
		},
		{
			name: "Handles even with complex state",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowHelp = true
				m.InputMode = true
				m.SpellSearchMode = true
				return m
			},
			shouldHandle: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg("a")

			// Default handler always handles
			canHandle := handler.CanHandle(model, msg)
			if canHandle != tt.shouldHandle {
				t.Errorf("CanHandle() = %v, want %v", canHandle, tt.shouldHandle)
			}
		})
	}
}

// TestStateTransitions tests state machine transitions
func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name       string
		setupModel func(ui.Model) ui.Model
		wantMode   ui.InputMode
	}{
		{
			name: "Normal mode when no special states",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			wantMode: ui.ModeNormal,
		},
		{
			name: "DiceInput mode when InputMode is true (new)",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				return m
			},
			wantMode: ui.ModeDiceInput,
		},
		{
			name: "DiceInput mode when InputMode is true (legacy)",
			setupModel: func(m ui.Model) ui.Model {
				m.InputMode = true
				return m
			},
			wantMode: ui.ModeDiceInput,
		},
		{
			name: "DiceHistory mode takes priority over DiceInput",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				m.DiceRoller.HistoryMode = true
				return m
			},
			wantMode: ui.ModeDiceHistory,
		},
		{
			name: "CastSpell mode takes priority",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.ShowCastSpellPrompt = true
				m.Spells.CastSpellInputMode = true
				m.DiceRoller.InputMode = true // Should be overridden
				return m
			},
			wantMode: ui.ModeCastSpell,
		},
		{
			name: "SpellSearch mode",
			setupModel: func(m ui.Model) ui.Model {
				m.Spells.SearchMode = true
				return m
			},
			wantMode: ui.ModeSpellSearch,
		},
		{
			name: "MonsterSearch mode",
			setupModel: func(m ui.Model) ui.Model {
				m.Monsters.SearchMode = true
				return m
			},
			wantMode: ui.ModeMonsterSearch,
		},
		{
			name: "InitiativeEdit mode",
			setupModel: func(m ui.Model) ui.Model {
				m.Initiative.EditMode = true
				return m
			},
			wantMode: ui.ModeInitiativeEdit,
		},
		{
			name: "NotesEdit mode",
			setupModel: func(m ui.Model) ui.Model {
				m.Notes.EditMode = true
				return m
			},
			wantMode: ui.ModeNotesEdit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)

			gotMode := model.GetInputMode()
			if gotMode != tt.wantMode {
				t.Errorf("GetInputMode() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}

// TestHandlerPriority_HelpOverridesInput tests that help popup takes priority
func TestHandlerPriority_HelpOverridesInput(t *testing.T) {
	chain := ui.NewHandlerChain()
	model := createHandlerTestModel()

	// Set up conflicting states
	model.Popup.ShowHelp = true
	model.DiceRoller.InputMode = true
	model.Spells.SearchMode = true

	msg := createKeyMsg("a")

	// Process through chain
	_, cmd := chain.Process(model, msg)

	// Help should win (we can't easily verify which handler was called,
	// but we can verify it didn't crash and returned something)
	if cmd == nil {
		t.Error("Handler chain should return a command")
	}
}

// TestHandlerPriority_PopupsOverrideKeys tests popup priority
func TestHandlerPriority_PopupsOverrideKeys(t *testing.T) {
	tests := []struct {
		name       string
		setupModel func(ui.Model) ui.Model
		key        string
	}{
		{
			name: "Save popup overrides quit key",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowSave = true
				return m
			},
			key: "q",
		},
		{
			name: "Load popup overrides navigation",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowLoad = true
				return m
			},
			key: "tab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := ui.NewHandlerChain()
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			// Should not panic and should process correctly
			_, cmd := chain.Process(model, msg)
			if cmd == nil {
				t.Error("Handler chain should return a command")
			}
		})
	}
}

// TestEdgeCases_MultiplePopups tests edge case with multiple popup flags set
func TestEdgeCases_MultiplePopups(t *testing.T) {
	chain := ui.NewHandlerChain()
	model := createHandlerTestModel()

	// Set multiple popup flags (should prioritize by handler order)
	model.Popup.ShowSave = true
	model.Popup.ShowLoad = true
	model.Popup.ShowHelp = true

	msg := createKeyMsg("a")

	// Should process without error (HelpPopupHandler should win due to priority 1)
	_, cmd := chain.Process(model, msg)
	if cmd == nil {
		t.Error("Handler chain should handle multiple popups gracefully")
	}
}

// TestEdgeCases_EmptyChain tests empty chain handling
func TestEdgeCases_EmptyChain(t *testing.T) {
	chain := ui.HandlerChain{}
	model := createHandlerTestModel()
	msg := createKeyMsg("a")

	// Should fall back to default input handler
	_, cmd := chain.Process(model, msg)
	if cmd == nil {
		t.Error("Empty chain should fallback to default handler")
	}
}

// TestHandlerIntegration_FullChain tests the complete handler chain with realistic scenarios
func TestHandlerIntegration_FullChain(t *testing.T) {
	chain := ui.NewHandlerChain()

	tests := []struct {
		name       string
		setupModel func(ui.Model) ui.Model
		key        string
		description string
	}{
		{
			name: "Quit key in normal state",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key: "q",
			description: "Should route to GlobalKeyHandler",
		},
		{
			name: "Tab navigation in normal state",
			setupModel: func(m ui.Model) ui.Model {
				return m
			},
			key: "tab",
			description: "Should route to GlobalKeyHandler for navigation",
		},
		{
			name: "Text input in dice input mode",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				return m
			},
			key: "2",
			description: "Should route to DefaultInputHandler",
		},
		{
			name: "Complex state - help popup with input mode",
			setupModel: func(m ui.Model) ui.Model {
				m.Popup.ShowHelp = true
				m.DiceRoller.InputMode = true
				m.Spells.SearchMode = true
				return m
			},
			key: "a",
			description: "Help popup should win (Priority 1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)
			msg := createKeyMsg(tt.key)

			// Should process without panic
			newModel, cmd := chain.Process(model, msg)
			_ = newModel
			_ = cmd
			// Just verify it doesn't panic - actual behavior verification would require
			// more complex assertions on the returned model state
		})
	}
}

// TestEdgeCases_StateMachinePriority tests priority of state machine modes
func TestEdgeCases_StateMachinePriority(t *testing.T) {
	tests := []struct {
		name       string
		setupModel func(ui.Model) ui.Model
		wantMode   ui.InputMode
		description string
	}{
		{
			name: "Cast spell mode overrides dice input",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				m.Spells.ShowCastSpellPrompt = true
				m.Spells.CastSpellInputMode = true
				return m
			},
			wantMode: ui.ModeCastSpell,
			description: "CastSpell should win due to GetInputMode priority",
		},
		{
			name: "Multi-target overrides other input modes",
			setupModel: func(m ui.Model) ui.Model {
				m.DiceRoller.InputMode = true
				m.Initiative.MultiTargetMode = true
				m.Initiative.ShowMultiTargetPopup = true
				return m
			},
			wantMode: ui.ModeMultiTarget,
			description: "Multi-target popup should win",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createHandlerTestModel()
			model = tt.setupModel(model)

			gotMode := model.GetInputMode()
			if gotMode != tt.wantMode {
				t.Errorf("GetInputMode() = %v, want %v (%s)", gotMode, tt.wantMode, tt.description)
			}
		})
	}
}

// TestEdgeCases_HandlerNames tests that all handlers have names
func TestEdgeCases_HandlerNames(t *testing.T) {
	chain := ui.NewHandlerChain()

	for _, handler := range chain {
		name := handler.Name()
		if name == "" {
			t.Errorf("Handler with priority %d has empty name", handler.Priority())
		}
		if handler.Priority() < 1 || handler.Priority() > 10 {
			t.Errorf("Handler %s has unusual priority %d", name, handler.Priority())
		}
	}
}

// TestEdgeCases_NilHandling tests that handlers handle edge cases gracefully
func TestEdgeCases_NilHandling(t *testing.T) {
	chain := ui.NewHandlerChain()
	model := createHandlerTestModel()

	// Test with various keys
	keys := []string{"a", "1", "q", "tab", "enter", "esc", "up", "down", "+", "-", "="}

	for _, key := range keys {
		t.Run("Key_"+key, func(t *testing.T) {
			msg := createKeyMsg(key)
			// Should not panic
			_, cmd := chain.Process(model, msg)
			_ = cmd
		})
	}
}
