// tests/multi_target_test.go
package tests

import (
	"fmt"
	"lazydnd/config"
	"lazydnd/ui"
	"testing"
)

// createMultiTargetTestModel creates a model with initiative entries for testing
func createMultiTargetTestModel() ui.Model {
	cfg := &config.Config{}
	styles := ui.NewStyles(cfg)

	model := ui.Model{
		Config:              cfg,
		Styles:              styles,
		Width:               80,
		Height:              24,
		ActivePanel:         ui.InitiativeTracker,
		ScrollOffset:        make(map[ui.PanelType]int),
		SelectedTargets:     make(map[int]bool),
		TargetSaveResults:   make(map[int]string),
		HPUndoStack:         []ui.HPHistoryEntry{},
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "Goblin 1",
				Type:       "monster",
				Initiative: 15,
				HP:         7,
				MaxHP:      7,
				AC:         15,
			},
			{
				Name:       "Goblin 2",
				Type:       "monster",
				Initiative: 14,
				HP:         7,
				MaxHP:      7,
				AC:         15,
			},
			{
				Name:       "Orc",
				Type:       "monster",
				Initiative: 12,
				HP:         15,
				MaxHP:      15,
				AC:         13,
			},
			{
				Name:       "Fighter",
				Type:       "player",
				Initiative: 18,
				AC:         18,
			},
		},
	}

	return model
}

// TestMultiTarget_BasicDamage tests basic damage application to multiple targets
func TestMultiTarget_BasicDamage(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select first two goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	// Apply 5 damage
	model = ui.ApplyMultiTargetDamage(model, 5)

	// Check HP was reduced
	if model.InitiativeList[0].HP != 2 {
		t.Errorf("Expected Goblin 1 HP to be 2, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[1].HP != 2 {
		t.Errorf("Expected Goblin 2 HP to be 2, got %d", model.InitiativeList[1].HP)
	}

	// Check Orc was not affected
	if model.InitiativeList[2].HP != 15 {
		t.Errorf("Expected Orc HP to remain 15, got %d", model.InitiativeList[2].HP)
	}
}

// TestMultiTarget_BasicHealing tests basic healing application to multiple targets
func TestMultiTarget_BasicHealing(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"

	// Damage goblins first
	model.InitiativeList[0].HP = 3
	model.InitiativeList[1].HP = 2

	// Select both goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	// Apply 4 healing
	model = ui.ApplyMultiTargetDamage(model, 4)

	// Check HP was increased
	if model.InitiativeList[0].HP != 7 {
		t.Errorf("Expected Goblin 1 HP to be 7, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[1].HP != 6 {
		t.Errorf("Expected Goblin 2 HP to be 6, got %d", model.InitiativeList[1].HP)
	}
}

// TestMultiTarget_HealingCappedAtMaxHP tests that healing doesn't exceed max HP
func TestMultiTarget_HealingCappedAtMaxHP(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"

	// Damage goblin
	model.InitiativeList[0].HP = 5

	// Select goblin
	model.SelectedTargets[0] = true

	// Apply 10 healing (more than max HP)
	model = ui.ApplyMultiTargetDamage(model, 10)

	// Check HP was capped at max HP
	if model.InitiativeList[0].HP != 7 {
		t.Errorf("Expected Goblin 1 HP to be capped at 7, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_DamageCappedAtZero tests that damage doesn't go below 0 HP
func TestMultiTarget_DamageCappedAtZero(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select goblin
	model.SelectedTargets[0] = true

	// Apply 100 damage (way more than HP)
	model = ui.ApplyMultiTargetDamage(model, 100)

	// Check HP was capped at 0
	if model.InitiativeList[0].HP != 0 {
		t.Errorf("Expected Goblin 1 HP to be capped at 0, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_SaveSuccess tests half damage on successful save
func TestMultiTarget_SaveSuccess(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"
	model.MultiTargetSaveMode = true

	// Select both goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	// Goblin 1 succeeds on save, Goblin 2 fails
	model.TargetSaveResults[0] = "success"
	model.TargetSaveResults[1] = "failure"

	// Apply 6 damage
	model = ui.ApplyMultiTargetDamage(model, 6)

	// Goblin 1 should take half damage (3)
	if model.InitiativeList[0].HP != 4 {
		t.Errorf("Expected Goblin 1 HP to be 4 (half damage), got %d", model.InitiativeList[0].HP)
	}

	// Goblin 2 should take full damage (6)
	if model.InitiativeList[1].HP != 1 {
		t.Errorf("Expected Goblin 2 HP to be 1 (full damage), got %d", model.InitiativeList[1].HP)
	}
}

// TestMultiTarget_SaveSuccessHealing tests save mode with healing
func TestMultiTarget_SaveSuccessHealing(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"
	model.MultiTargetSaveMode = true

	// Damage goblins
	model.InitiativeList[0].HP = 3
	model.InitiativeList[1].HP = 3

	// Select both goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	// One succeeds on save, one fails
	model.TargetSaveResults[0] = "success"
	model.TargetSaveResults[1] = "failure"

	// Apply 4 healing
	model = ui.ApplyMultiTargetDamage(model, 4)

	// Goblin 1 with save success and healing is skipped (not a valid combination)
	if model.InitiativeList[0].HP != 3 {
		t.Errorf("Expected Goblin 1 HP to remain 3 (success+healing skipped), got %d", model.InitiativeList[0].HP)
	}

	// Goblin 2 with save failure receives full healing
	if model.InitiativeList[1].HP != 7 {
		t.Errorf("Expected Goblin 2 HP to be 7 (failure allows healing), got %d", model.InitiativeList[1].HP)
	}
}

// TestMultiTarget_NoSaveResultSkipsTarget tests targets without save results are skipped
func TestMultiTarget_NoSaveResultSkipsTarget(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"
	model.MultiTargetSaveMode = true

	// Select both goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	// Only Goblin 1 has save result
	model.TargetSaveResults[0] = "failure"
	// Goblin 2 has no save result (empty string)

	// Apply 5 damage
	model = ui.ApplyMultiTargetDamage(model, 5)

	// Goblin 1 should take damage
	if model.InitiativeList[0].HP != 2 {
		t.Errorf("Expected Goblin 1 HP to be 2, got %d", model.InitiativeList[0].HP)
	}

	// Goblin 2 should be unchanged (no save result)
	if model.InitiativeList[1].HP != 7 {
		t.Errorf("Expected Goblin 2 HP to remain 7, got %d", model.InitiativeList[1].HP)
	}
}

// TestMultiTarget_PlayerTargetsIgnored tests that player targets are ignored
func TestMultiTarget_PlayerTargetsIgnored(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select player
	model.SelectedTargets[3] = true

	// Apply 10 damage
	model = ui.ApplyMultiTargetDamage(model, 10)

	// Player HP should not be tracked/changed (players manage their own HP)
	// This test just ensures no panic or error occurs
}

// TestMultiTarget_MixedMonstersAndPlayers tests mixed target selection
func TestMultiTarget_MixedMonstersAndPlayers(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select goblin and player
	model.SelectedTargets[0] = true
	model.SelectedTargets[3] = true

	// Apply 5 damage
	model = ui.ApplyMultiTargetDamage(model, 5)

	// Only monster should be affected
	if model.InitiativeList[0].HP != 2 {
		t.Errorf("Expected Goblin 1 HP to be 2, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_AllTargets tests selecting all targets
func TestMultiTarget_AllTargets(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select all monsters
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true
	model.SelectedTargets[2] = true

	// Apply 3 damage
	model = ui.ApplyMultiTargetDamage(model, 3)

	// Check all monsters were damaged
	if model.InitiativeList[0].HP != 4 {
		t.Errorf("Expected Goblin 1 HP to be 4, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[1].HP != 4 {
		t.Errorf("Expected Goblin 2 HP to be 4, got %d", model.InitiativeList[1].HP)
	}
	if model.InitiativeList[2].HP != 12 {
		t.Errorf("Expected Orc HP to be 12, got %d", model.InitiativeList[2].HP)
	}
}

// TestMultiTarget_NoTargetsSelected tests behavior with no targets
func TestMultiTarget_NoTargetsSelected(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Don't select any targets
	originalHP := []int{
		model.InitiativeList[0].HP,
		model.InitiativeList[1].HP,
		model.InitiativeList[2].HP,
	}

	// Apply 10 damage
	model = ui.ApplyMultiTargetDamage(model, 10)

	// All HP should remain unchanged
	if model.InitiativeList[0].HP != originalHP[0] {
		t.Errorf("Expected Goblin 1 HP to remain %d, got %d", originalHP[0], model.InitiativeList[0].HP)
	}
	if model.InitiativeList[1].HP != originalHP[1] {
		t.Errorf("Expected Goblin 2 HP to remain %d, got %d", originalHP[1], model.InitiativeList[1].HP)
	}
	if model.InitiativeList[2].HP != originalHP[2] {
		t.Errorf("Expected Orc HP to remain %d, got %d", originalHP[2], model.InitiativeList[2].HP)
	}
}

// TestMultiTarget_ZeroDamage tests applying zero damage
func TestMultiTarget_ZeroDamage(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	model.SelectedTargets[0] = true

	originalHP := model.InitiativeList[0].HP

	// Apply 0 damage
	model = ui.ApplyMultiTargetDamage(model, 0)

	// HP should remain unchanged
	if model.InitiativeList[0].HP != originalHP {
		t.Errorf("Expected Goblin 1 HP to remain %d, got %d", originalHP, model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_ZeroHealing tests applying zero healing
func TestMultiTarget_ZeroHealing(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"

	model.InitiativeList[0].HP = 3
	model.SelectedTargets[0] = true

	// Apply 0 healing
	model = ui.ApplyMultiTargetDamage(model, 0)

	// HP should remain unchanged
	if model.InitiativeList[0].HP != 3 {
		t.Errorf("Expected Goblin 1 HP to remain 3, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_LargeDamageValue tests very large damage values
func TestMultiTarget_LargeDamageValue(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	model.SelectedTargets[2] = true // Select Orc with 15 HP

	// Apply 1000 damage
	model = ui.ApplyMultiTargetDamage(model, 1000)

	// Should be capped at 0
	if model.InitiativeList[2].HP != 0 {
		t.Errorf("Expected Orc HP to be capped at 0, got %d", model.InitiativeList[2].HP)
	}
}

// TestMultiTarget_OddDamageHalving tests odd damage values with saves
func TestMultiTarget_OddDamageHalving(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"
	model.MultiTargetSaveMode = true

	model.SelectedTargets[0] = true
	model.TargetSaveResults[0] = "success"

	// Apply 7 damage (odd number)
	model = ui.ApplyMultiTargetDamage(model, 7)

	// Should take 3 damage (7/2 = 3 with integer division)
	if model.InitiativeList[0].HP != 4 {
		t.Errorf("Expected Goblin 1 HP to be 4 (7/2=3 damage), got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_SaveModeOffIgnoresSaveResults tests that save results are ignored when save mode is off
func TestMultiTarget_SaveModeOffIgnoresSaveResults(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"
	model.MultiTargetSaveMode = false // Save mode OFF

	model.SelectedTargets[0] = true
	// Set save result even though save mode is off
	model.TargetSaveResults[0] = "success"

	// Apply 6 damage
	model = ui.ApplyMultiTargetDamage(model, 6)

	// Should take full damage (save mode is off)
	if model.InitiativeList[0].HP != 1 {
		t.Errorf("Expected Goblin 1 HP to be 1 (full damage), got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_HPHistory tests that HP changes are recorded in undo history
func TestMultiTarget_HPHistory(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select two goblins
	model.SelectedTargets[0] = true
	model.SelectedTargets[1] = true

	initialUndoStackSize := len(model.HPUndoStack)

	// Apply 5 damage
	model = ui.ApplyMultiTargetDamage(model, 5)

	// Should have 2 new entries in undo stack (one per target)
	if len(model.HPUndoStack) != initialUndoStackSize+2 {
		t.Errorf("Expected undo stack to grow by 2, got %d entries", len(model.HPUndoStack)-initialUndoStackSize)
	}
}

// TestMultiTarget_DifferentMaxHP tests targets with different max HP
func TestMultiTarget_DifferentMaxHP(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Select goblin (7 HP) and orc (15 HP)
	model.SelectedTargets[0] = true
	model.SelectedTargets[2] = true

	// Apply 10 damage
	model = ui.ApplyMultiTargetDamage(model, 10)

	// Goblin should be at 0 HP (7 - 10 = -3, capped at 0)
	if model.InitiativeList[0].HP != 0 {
		t.Errorf("Expected Goblin HP to be 0, got %d", model.InitiativeList[0].HP)
	}

	// Orc should be at 5 HP (15 - 10 = 5)
	if model.InitiativeList[2].HP != 5 {
		t.Errorf("Expected Orc HP to be 5, got %d", model.InitiativeList[2].HP)
	}
}

// TestMultiTarget_SequentialDamage tests applying damage multiple times
func TestMultiTarget_SequentialDamage(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	model.SelectedTargets[0] = true

	// Apply 2 damage three times
	model = ui.ApplyMultiTargetDamage(model, 2)
	model = ui.ApplyMultiTargetDamage(model, 2)
	model = ui.ApplyMultiTargetDamage(model, 2)

	// Should have 1 HP left (7 - 6 = 1)
	if model.InitiativeList[0].HP != 1 {
		t.Errorf("Expected Goblin HP to be 1, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_HealingAfterDamage tests healing after taking damage
func TestMultiTarget_HealingAfterDamage(t *testing.T) {
	model := createMultiTargetTestModel()

	model.SelectedTargets[0] = true

	// Apply 5 damage
	model.MultiTargetType = "damage"
	model = ui.ApplyMultiTargetDamage(model, 5)

	if model.InitiativeList[0].HP != 2 {
		t.Errorf("Expected Goblin HP to be 2 after damage, got %d", model.InitiativeList[0].HP)
	}

	// Apply 3 healing
	model.MultiTargetType = "healing"
	model = ui.ApplyMultiTargetDamage(model, 3)

	if model.InitiativeList[0].HP != 5 {
		t.Errorf("Expected Goblin HP to be 5 after healing, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_MaxHPBoundary tests healing at exact max HP boundary
func TestMultiTarget_MaxHPBoundary(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"

	// Set HP to 1 below max
	model.InitiativeList[0].HP = 6 // Max is 7

	model.SelectedTargets[0] = true

	// Heal exactly to max
	model = ui.ApplyMultiTargetDamage(model, 1)

	if model.InitiativeList[0].HP != 7 {
		t.Errorf("Expected Goblin HP to be exactly 7, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_MinHPBoundary tests damage at exact 0 HP boundary
func TestMultiTarget_MinHPBoundary(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	model.SelectedTargets[0] = true

	// Apply exactly enough damage to reach 0
	model = ui.ApplyMultiTargetDamage(model, 7)

	if model.InitiativeList[0].HP != 0 {
		t.Errorf("Expected Goblin HP to be exactly 0, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_AlreadyAtZeroHP tests further damage to 0 HP target
func TestMultiTarget_AlreadyAtZeroHP(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"

	// Set HP to 0
	model.InitiativeList[0].HP = 0

	model.SelectedTargets[0] = true

	// Apply more damage
	model = ui.ApplyMultiTargetDamage(model, 10)

	// Should remain at 0
	if model.InitiativeList[0].HP != 0 {
		t.Errorf("Expected Goblin HP to remain 0, got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_AlreadyAtMaxHP tests healing at max HP
func TestMultiTarget_AlreadyAtMaxHP(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "healing"

	// Already at max HP (7/7)
	model.SelectedTargets[0] = true

	// Try to heal
	model = ui.ApplyMultiTargetDamage(model, 5)

	// Should remain at max
	if model.InitiativeList[0].HP != 7 {
		t.Errorf("Expected Goblin HP to remain at max (7), got %d", model.InitiativeList[0].HP)
	}
}

// TestMultiTarget_SaveSuccessRoundingDown tests save success with rounding
func TestMultiTarget_SaveSuccessRoundingDown(t *testing.T) {
	model := createMultiTargetTestModel()
	model.MultiTargetType = "damage"
	model.MultiTargetSaveMode = true

	model.SelectedTargets[0] = true
	model.TargetSaveResults[0] = "success"

	testCases := []struct {
		damage   int
		expected int // Expected HP after damage
	}{
		{1, 7},  // 1/2 = 0 damage, HP: 7-0 = 7
		{2, 6},  // 2/2 = 1 damage, HP: 7-1 = 6
		{3, 6},  // 3/2 = 1 damage, HP: 7-1 = 6
		{4, 5},  // 4/2 = 2 damage, HP: 7-2 = 5
		{5, 5},  // 5/2 = 2 damage, HP: 7-2 = 5
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("damage_%d", tc.damage), func(t *testing.T) {
			// Reset HP
			testModel := createMultiTargetTestModel()
			testModel.MultiTargetType = "damage"
			testModel.MultiTargetSaveMode = true
			testModel.SelectedTargets[0] = true
			testModel.TargetSaveResults[0] = "success"

			testModel = ui.ApplyMultiTargetDamage(testModel, tc.damage)

			if testModel.InitiativeList[0].HP != tc.expected {
				t.Errorf("With %d damage and save success, expected HP %d, got %d",
					tc.damage, tc.expected, testModel.InitiativeList[0].HP)
			}
		})
	}
}
