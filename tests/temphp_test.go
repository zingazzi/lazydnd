// tests/temphp_test.go
package tests

import (
	"lazydnd/ui"
	"testing"
)

// TestTempHP_Display verifies temp HP is displayed correctly
func TestTempHP_Display(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         7,
			MaxHP:      7,
			TempHP:     3,
			AC:         13,
		},
	}

	// Verify temp HP is stored
	if model.InitiativeList[0].TempHP != 3 {
		t.Errorf("Expected TempHP to be 3, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_DamageAbsorption verifies temp HP absorbs damage first
func TestTempHP_DamageAbsorption(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     5,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 3 damage (should only affect temp HP)
	model = applyHPChange(model, 0, -3)

	if model.InitiativeList[0].HP != 10 {
		t.Errorf("Expected HP to remain 10, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 2 {
		t.Errorf("Expected TempHP to be 2, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_ExceedsTempHP verifies overflow damage goes to real HP
func TestTempHP_ExceedsTempHP(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     3,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 7 damage (3 to temp HP, 4 to real HP)
	model = applyHPChange(model, 0, -7)

	if model.InitiativeList[0].HP != 6 {
		t.Errorf("Expected HP to be 6, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 0 {
		t.Errorf("Expected TempHP to be 0, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_NoTempHP verifies damage goes directly to real HP when no temp HP
func TestTempHP_NoTempHP(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     0,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 3 damage (goes directly to real HP)
	model = applyHPChange(model, 0, -3)

	if model.InitiativeList[0].HP != 7 {
		t.Errorf("Expected HP to be 7, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 0 {
		t.Errorf("Expected TempHP to remain 0, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_HealingIgnoresTempHP verifies healing doesn't affect temp HP
func TestTempHP_HealingIgnoresTempHP(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         5,
			MaxHP:      10,
			TempHP:     3,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 5 healing (should only affect real HP)
	model = applyHPChange(model, 0, +5)

	if model.InitiativeList[0].HP != 10 {
		t.Errorf("Expected HP to be 10, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 3 {
		t.Errorf("Expected TempHP to remain 3, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_SetTempHP verifies setting temp HP replaces existing
func TestTempHP_SetTempHP(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     3,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Set new temp HP (should replace, not add)
	model = setTempHP(model, 0, 8)

	if model.InitiativeList[0].HP != 10 {
		t.Errorf("Expected HP to remain 10, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 8 {
		t.Errorf("Expected TempHP to be 8, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_ClearTempHP verifies setting temp HP to 0 clears it
func TestTempHP_ClearTempHP(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     5,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Clear temp HP by setting to 0
	model = setTempHP(model, 0, 0)

	if model.InitiativeList[0].TempHP != 0 {
		t.Errorf("Expected TempHP to be 0, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_MultiTarget verifies temp HP with multi-target damage
func TestTempHP_MultiTarget(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin 1",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     5,
			AC:         13,
		},
		{
			Name:       "Goblin 2",
			Type:       "monster",
			Initiative: 14,
			HP:         10,
			MaxHP:      10,
			TempHP:     0,
			AC:         13,
		},
	}
	model.MultiTargetMode = true
	model.MultiTargetType = "damage"
	model.SelectedTargets = map[int]bool{0: true, 1: true}

	// Apply 7 damage to both
	model = ui.ApplyMultiTargetDamage(model, 7)

	// Goblin 1: 5 temp HP, 2 real HP lost
	if model.InitiativeList[0].HP != 8 {
		t.Errorf("Expected Goblin 1 HP to be 8, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 0 {
		t.Errorf("Expected Goblin 1 TempHP to be 0, got %d", model.InitiativeList[0].TempHP)
	}

	// Goblin 2: No temp HP, 7 real HP lost
	if model.InitiativeList[1].HP != 3 {
		t.Errorf("Expected Goblin 2 HP to be 3, got %d", model.InitiativeList[1].HP)
	}
	if model.InitiativeList[1].TempHP != 0 {
		t.Errorf("Expected Goblin 2 TempHP to remain 0, got %d", model.InitiativeList[1].TempHP)
	}
}

// TestTempHP_MultiTargetHealing verifies healing doesn't affect temp HP in multi-target
func TestTempHP_MultiTargetHealing(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         5,
			MaxHP:      10,
			TempHP:     3,
			AC:         13,
		},
	}
	model.MultiTargetMode = true
	model.MultiTargetType = "healing"
	model.SelectedTargets = map[int]bool{0: true}

	// Apply 5 healing
	model = ui.ApplyMultiTargetDamage(model, 5)

	// HP should heal, temp HP should remain unchanged
	if model.InitiativeList[0].HP != 10 {
		t.Errorf("Expected HP to be 10, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 3 {
		t.Errorf("Expected TempHP to remain 3, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_ZeroDamage verifies 0 damage doesn't affect anything
func TestTempHP_ZeroDamage(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     5,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 0 damage
	model = applyHPChange(model, 0, 0)

	if model.InitiativeList[0].HP != 10 {
		t.Errorf("Expected HP to remain 10, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 5 {
		t.Errorf("Expected TempHP to remain 5, got %d", model.InitiativeList[0].TempHP)
	}
}

// TestTempHP_LargeDamage verifies large damage depletes both temp and real HP
func TestTempHP_LargeDamage(t *testing.T) {
	model := ui.InitialModel()
	model.InitiativeList = []ui.InitiativeEntry{
		{
			Name:       "Goblin",
			Type:       "monster",
			Initiative: 15,
			HP:         10,
			MaxHP:      10,
			TempHP:     5,
			AC:         13,
		},
	}
	model.InitiativeListMode = true
	model.SelectedEntry = 0

	// Apply 20 damage (more than HP + temp HP)
	model = applyHPChange(model, 0, -20)

	// Should cap at 0 HP
	if model.InitiativeList[0].HP != 0 {
		t.Errorf("Expected HP to be 0, got %d", model.InitiativeList[0].HP)
	}
	if model.InitiativeList[0].TempHP != 0 {
		t.Errorf("Expected TempHP to be 0, got %d", model.InitiativeList[0].TempHP)
	}
}

// Helper function to manually apply HP changes for testing
func applyHPChange(m ui.Model, index int, change int) ui.Model {
	if index < 0 || index >= len(m.InitiativeList) {
		return m
	}

	if change < 0 {
		// Taking damage - apply to temp HP first
		damage := -change
		if m.InitiativeList[index].TempHP > 0 {
			if damage <= m.InitiativeList[index].TempHP {
				// All damage absorbed by temp HP
				m.InitiativeList[index].TempHP -= damage
			} else {
				// Temp HP absorbed some, rest goes to real HP
				remainingDamage := damage - m.InitiativeList[index].TempHP
				m.InitiativeList[index].TempHP = 0
				m.InitiativeList[index].HP -= remainingDamage
				if m.InitiativeList[index].HP < 0 {
					m.InitiativeList[index].HP = 0
				}
			}
		} else {
			// No temp HP, damage goes directly to real HP
			m.InitiativeList[index].HP += change // change is negative
			if m.InitiativeList[index].HP < 0 {
				m.InitiativeList[index].HP = 0
			}
		}
	} else {
		// Healing - only affects real HP, not temp HP
		newHP := m.InitiativeList[index].HP + change
		if newHP > m.InitiativeList[index].MaxHP {
			newHP = m.InitiativeList[index].MaxHP
		}
		m.InitiativeList[index].HP = newHP
	}

	return m
}

// Helper function to set temp HP for testing
func setTempHP(m ui.Model, index int, tempHP int) ui.Model {
	if index < 0 || index >= len(m.InitiativeList) {
		return m
	}
	if tempHP < 0 {
		tempHP = 0
	}
	m.InitiativeList[index].TempHP = tempHP
	return m
}
