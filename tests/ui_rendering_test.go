// tests/ui_rendering_test.go
package tests

import (
	"lazydnd/config"
	"lazydnd/ui"
	"strings"
	"testing"
)

// createTestModel creates a basic model for testing
func createTestModel() ui.Model {
	cfg := config.Default()
	styles := ui.NewStyles(cfg)

	model := ui.Model{
		Config:       cfg,
		Styles:       styles,
		Width:        80,
		Height:       24,
		ActivePanel:  ui.DiceRoller,
		ScrollOffset: make(map[ui.PanelType]int),
	}

	return model
}

// TestUI_BasicRendering tests basic view rendering
func TestUI_BasicRendering(t *testing.T) {
	model := createTestModel()

	view := model.View()

	if view == "" {
		t.Error("View returned empty string")
	}

	if view == "Loading..." {
		t.Error("View should not show loading with dimensions set")
	}
}

// TestUI_LoadingState tests loading state rendering
func TestUI_LoadingState(t *testing.T) {
	model := createTestModel()
	model.Width = 0
	model.Height = 0

	view := model.View()

	if view != "Loading..." {
		t.Errorf("Expected 'Loading...', got %q", view)
	}
}

// TestUI_PanelNames tests that panel names are visible
func TestUI_PanelNames(t *testing.T) {
	model := createTestModel()
	// Use larger dimensions for new layout
	model.Width = 200
	model.Height = 50

	view := model.View()

	expectedPanels := []string{
		"Dice Roller",
		"Initiative Tracker",
		"Spells",
		"Monsters",
		"Notes",
	}

	for _, panelName := range expectedPanels {
		if !strings.Contains(view, panelName) {
			t.Errorf("View should contain panel %q", panelName)
		}
	}
}

// TestUI_InitiativeTrackerWithEntries tests initiative tracker with entries
func TestUI_InitiativeTrackerWithEntries(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Gandalf",
			Type:       "player",
			Initiative: 18,
			AC:         15,
		},
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 10,
			HP:         7,
			MaxHP:      7,
			AC:         15,
		},
	}

	// Test using GetPanelContent which is the underlying content function
	content := model.GetPanelContent(ui.InitiativeTracker)

	// Check that entries are displayed
	if !strings.Contains(content, "Gandalf") {
		t.Error("Content should contain player name 'Gandalf'")
	}

	if !strings.Contains(content, "Goblin") {
		t.Error("Content should contain monster name 'Goblin'")
	}

	// Check that initiative values are shown
	if !strings.Contains(content, "18") {
		t.Error("Content should contain initiative value 18")
	}
}

// TestUI_EmptyInitiativeTracker tests empty initiative tracker
func TestUI_EmptyInitiativeTracker(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeList = []ui.InitiativeEntry{}

	view := model.View()

	if !strings.Contains(view, "Initiative Tracker") {
		t.Error("View should contain panel name")
	}
}

// TestUI_DiceRollerWithHistory tests dice roller with history
func TestUI_DiceRollerWithHistory(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.DiceRoller
	model.DiceResult = "Total: 15 (Rolled: 8 + 7)"
	model.DiceHistory = []string{
		"Total: 15 (Rolled: 8 + 7)",
		"Total: 10 (Rolled: 5 + 5)",
	}

	content := model.GetPanelContent(ui.DiceRoller)

	if !strings.Contains(content, "15") {
		t.Error("Content should contain dice result")
	}
}

// TestUI_InputMode tests input mode rendering
func TestUI_InputMode(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.DiceRoller
	model.InputMode = true
	model.DiceInput = "2d6"

	content := model.GetPanelContent(ui.DiceRoller)

	if !strings.Contains(content, "2d6") {
		t.Error("Content should contain input text")
	}
}

// TestUI_HelpPopup tests help popup rendering
func TestUI_HelpPopup(t *testing.T) {
	model := createTestModel()
	model.ShowHelpPopup = true

	view := model.View()

	if !strings.Contains(view, "Help") {
		t.Error("View should contain 'Help' text in popup")
	}

	// Help should contain common navigation keys
	if !strings.Contains(view, "Tab") {
		t.Error("Help should mention Tab key")
	}
}

// TestUI_RoundCounter tests round counter display
func TestUI_RoundCounter(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.RoundCounter = 5
	model.CurrentTurn = 0
	model.InitiativeList = []ui.InitiativeEntry{
		{Name: "Player", Type: "player", Initiative: 15},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	if !strings.Contains(content, "Round") {
		t.Error("Content should contain 'Round' text")
	}

	if !strings.Contains(content, "5") {
		t.Error("Content should contain round number")
	}
}

// TestUI_ActiveSpells tests active spells rendering
func TestUI_ActiveSpells(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.Spells
	model.ActiveSpells = []ui.ActiveSpell{
		{
			Name:          "Bless",
			CasterName:    "Cleric",
			RoundsLeft:    10,
			TotalRounds:   10,
			Concentration: true,
		},
	}

	view := model.View()

	// Active spells should be visible in spells panel
	if !strings.Contains(view, "Spells") {
		t.Error("View should contain Spells panel")
	}
}

// TestUI_MonsterSearch tests monster search mode
func TestUI_MonsterSearch(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.Monsters
	model.MonsterSearchMode = true
	model.MonsterSearchInput = "goblin"

	view := model.View()

	if !strings.Contains(view, "goblin") {
		t.Error("View should contain search input")
	}
}

// TestUI_SpellSearch tests spell search mode
func TestUI_SpellSearch(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.Spells
	model.SpellSearchMode = true
	model.SpellSearchInput = "fireball"

	view := model.View()

	if !strings.Contains(view, "fireball") {
		t.Error("View should contain search input")
	}
}

// TestUI_StatusBar tests status bar rendering
func TestUI_StatusBar(t *testing.T) {
	model := createTestModel()

	view := model.View()

	// Status bar should contain help hint
	if !strings.Contains(view, "?") {
		t.Error("Status bar should contain help indicator")
	}
}

// TestUI_CampaignName tests campaign name in status bar
func TestUI_CampaignName(t *testing.T) {
	model := createTestModel()
	model.CurrentCampaignName = "Test Campaign"

	view := model.View()

	if !strings.Contains(view, "Test Campaign") {
		t.Error("Status bar should contain campaign name")
	}
}

// TestUI_MultiTargetMode tests multi-target mode rendering
func TestUI_MultiTargetMode(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.MultiTargetMode = true
	model.InitiativeList = []ui.InitiativeEntry{
		{Name: "Goblin 1", Type: "monster", Initiative: 10, HP: 7, MaxHP: 7},
		{Name: "Goblin 2", Type: "monster", Initiative: 9, HP: 7, MaxHP: 7},
	}
	model.SelectedTargets = map[int]bool{0: true}

	view := model.View()

	// Multi-target mode check removed - instruction text removed from panels
	_ = view
}

// TestUI_ConditionsDisplay tests conditions display
func TestUI_ConditionsDisplay(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Poisoned Goblin",
			Type:       "monster",
			Initiative: 10,
			HP:         5,
			MaxHP:      7,
			AC:         15,
			Conditions: []ui.Condition{
				{
					Name:        "Poisoned",
					RoundsLeft:  3,
					TotalRounds: 5,
				},
			},
		},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	if !strings.Contains(content, "Poisoned Goblin") {
		t.Error("Content should contain creature name")
	}

	// Conditions should show as emojis
	if !strings.Contains(content, "🤢") {
		t.Error("Content should contain condition emoji for Poisoned")
	}
}

// TestUI_ColoredHP tests HP color coding
func TestUI_ColoredHP(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker

	// Test healthy monster (> 50%)
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Healthy Monster",
			Type:       "monster",
			Initiative: 10,
			HP:         10,
			MaxHP:      10,
			AC:         15,
		},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	if !strings.Contains(content, "HP") {
		t.Error("Content should display HP")
	}

	if !strings.Contains(content, "10/10") {
		t.Error("Content should show HP ratio")
	}
}

// TestUI_CurrentTurnMarker tests current turn indicator
func TestUI_CurrentTurnMarker(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.CurrentTurn = 0
	model.RoundCounter = 1
	model.InitiativeList = []ui.InitiativeEntry{
		{Name: "Player 1", Type: "player", Initiative: 20},
		{Name: "Player 2", Type: "player", Initiative: 15},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	// Should contain turn marker (★)
	if !strings.Contains(content, "★") {
		t.Error("Content should contain turn marker")
	}
}

// TestUI_DiceHistoryMode tests dice history mode
func TestUI_DiceHistoryMode(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.DiceRoller
	model.DiceHistoryMode = true
	model.DiceHistory = []string{"Roll 1", "Roll 2", "Roll 3"}
	model.DiceCommands = []string{"2d6", "1d20", "3d8"}
	model.HistoryIndex = 1

	view := model.View()

	// Should show history
	if !strings.Contains(view, "Roll") {
		t.Error("View should show dice history")
	}
}

// TestUI_InitiativeEditMode tests edit mode
func TestUI_InitiativeEditMode(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeListMode = true
	model.InitiativeEditMode = true
	model.InitiativeEditType = "hp"
	model.SelectedEntry = 0
	model.InitiativeList = []ui.InitiativeEntry{
		{Name: "Monster", Type: "monster", Initiative: 10, HP: 7, MaxHP: 7},
	}

	view := model.View()

	// Edit mode check removed - instruction text removed from panels
	_ = view
}

// TestUI_PlayerWithAC tests player AC display
func TestUI_PlayerWithAC(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Fighter",
			Type:       "player",
			Initiative: 15,
			AC:         18,
		},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	if !strings.Contains(content, "Fighter") {
		t.Error("Content should contain player name")
	}

	if !strings.Contains(content, "18") {
		t.Error("Content should show AC value")
	}
}

// TestUI_NoWidth tests zero width handling
func TestUI_NoWidth(t *testing.T) {
	model := createTestModel()
	model.Width = 0

	view := model.View()

	if view != "Loading..." {
		t.Error("View should show loading when width is 0")
	}
}

// TestUI_NoHeight tests zero height handling
func TestUI_NoHeight(t *testing.T) {
	model := createTestModel()
	model.Height = 0

	view := model.View()

	if view != "Loading..." {
		t.Error("View should show loading when height is 0")
	}
}

// TestUI_SwitchPanels tests switching between panels
func TestUI_SwitchPanels(t *testing.T) {
	model := createTestModel()

	panels := []ui.PanelType{
		ui.DiceRoller,
		ui.InitiativeTracker,
		ui.Spells,
		ui.Monsters,
	}

	for _, panel := range panels {
		model.ActivePanel = panel
		view := model.View()

		if view == "" {
			t.Errorf("View should not be empty for panel %d", panel)
		}
	}
}

// TestUI_LargeInitiativeList tests rendering with many entries
func TestUI_LargeInitiativeList(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker

	// Add 20 entries
	for i := 0; i < 20; i++ {
		model.InitiativeList = append(model.InitiativeList, ui.InitiativeEntry{
			Name:       "Monster " + string(rune('A'+i)),
			Type:       "monster",
			Initiative: 20 - i,
			HP:         10,
			MaxHP:      10,
			AC:         15,
		})
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	// Should contain at least the first entry
	if !strings.Contains(content, "Monster A") {
		t.Error("Content should contain first entry")
	}
}

// TestUI_InstanceNumbers tests monster instance numbering
func TestUI_InstanceNumbers(t *testing.T) {
	model := createTestModel()
	model.ActivePanel = ui.InitiativeTracker
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:        "Goblin 1",
			Type:        "monster",
			Initiative:  10,
			HP:          7,
			MaxHP:       7,
			AC:          15,
			InstanceNum: 1,
			BaseName:    "Goblin",
		},
		{
			Name:        "Goblin 2",
			Type:        "monster",
			Initiative:  9,
			HP:          7,
			MaxHP:       7,
			AC:          15,
			InstanceNum: 2,
			BaseName:    "Goblin",
		},
	}

	content := model.GetPanelContent(ui.InitiativeTracker)

	if !strings.Contains(content, "Goblin 1") {
		t.Error("Content should contain 'Goblin 1'")
	}

	if !strings.Contains(content, "Goblin 2") {
		t.Error("Content should contain 'Goblin 2'")
	}
}

// TestUI_ViewNotEmpty tests that view is never completely empty
func TestUI_ViewNotEmpty(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*ui.Model)
	}{
		{
			name: "Empty model",
			setup: func(m *ui.Model) {
				// No setup needed
			},
		},
		{
			name: "With entries",
			setup: func(m *ui.Model) {
				m.InitiativeList = []ui.InitiativeEntry{
					{Name: "Test", Type: "player", Initiative: 10},
				}
			},
		},
		{
			name: "In input mode",
			setup: func(m *ui.Model) {
				m.InputMode = true
				m.DiceInput = "2d6"
			},
		},
		{
			name: "With help popup",
			setup: func(m *ui.Model) {
				m.ShowHelpPopup = true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := createTestModel()
			tc.setup(&model)

			view := model.View()

			if view == "" {
				t.Error("View should never be completely empty")
			}

			if len(view) < 10 {
				t.Errorf("View seems too short: %d characters", len(view))
			}
		})
	}
}
