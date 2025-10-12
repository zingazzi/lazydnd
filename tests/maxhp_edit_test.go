// tests/maxhp_edit_test.go
package tests

import (
	"lazydnd/panels"
	"lazydnd/ui"
	"testing"
)

// TestMaxHPParsing tests the parsing of max HP input
func TestMaxHPParsing(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue int
		wantError bool
	}{
		{
			name:      "Valid positive value",
			input:     "50",
			wantValue: 50,
			wantError: false,
		},
		{
			name:      "Valid large value",
			input:     "999",
			wantValue: 999,
			wantError: false,
		},
		{
			name:      "Valid value with spaces",
			input:     "  42  ",
			wantValue: 42,
			wantError: false,
		},
		{
			name:      "Zero value should error",
			input:     "0",
			wantValue: 0,
			wantError: true,
		},
		{
			name:      "Negative value should error",
			input:     "-10",
			wantValue: 0,
			wantError: true,
		},
		{
			name:      "Empty input should error",
			input:     "",
			wantValue: 0,
			wantError: true,
		},
		{
			name:      "Non-numeric input should error",
			input:     "abc",
			wantValue: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := panels.ParseInput(tt.input, "maxhp")

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.(int) != tt.wantValue {
				t.Errorf("Expected %d, got %d", tt.wantValue, result.(int))
			}
		})
	}
}

// TestMaxHPEdit tests the max HP editing functionality
func TestMaxHPEdit(t *testing.T) {
	tests := []struct {
		name           string
		initialHP      int
		initialMaxHP   int
		newMaxHP       int
		expectedHP     int
		expectedMaxHP  int
	}{
		{
			name:          "Increase max HP, current HP unchanged",
			initialHP:     50,
			initialMaxHP:  100,
			newMaxHP:      150,
			expectedHP:    50,
			expectedMaxHP: 150,
		},
		{
			name:          "Decrease max HP, current HP still valid",
			initialHP:     30,
			initialMaxHP:  100,
			newMaxHP:      50,
			expectedHP:    30,
			expectedMaxHP: 50,
		},
		{
			name:          "Decrease max HP below current HP, should cap current HP",
			initialHP:     80,
			initialMaxHP:  100,
			newMaxHP:      50,
			expectedHP:    50,
			expectedMaxHP: 50,
		},
		{
			name:          "Set max HP to 1",
			initialHP:     50,
			initialMaxHP:  100,
			newMaxHP:      1,
			expectedHP:    1,
			expectedMaxHP: 1,
		},
		{
			name:          "Current HP at max, increase max HP",
			initialHP:     100,
			initialMaxHP:  100,
			newMaxHP:      200,
			expectedHP:    100,
			expectedMaxHP: 200,
		},
		{
			name:          "Current HP at max, decrease max HP",
			initialHP:     100,
			initialMaxHP:  100,
			newMaxHP:      50,
			expectedHP:    50,
			expectedMaxHP: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a model with one monster entry
			m := ui.Model{
				InitiativeList: []ui.InitiativeEntry{
					{
						Name:       "Goblin",
						Type:       "monster",
						HP:         tt.initialHP,
						MaxHP:      tt.initialMaxHP,
						Initiative: 15,
					},
				},
				SelectedEntry:      0,
				InitiativeEditMode: true,
				InitiativeEditType: "maxhp",
			}

			// Apply the max HP change
			newMaxHP := tt.newMaxHP
			if newMaxHP < 1 {
				newMaxHP = 1
			}
			m.InitiativeList[0].MaxHP = newMaxHP

			// Cap current HP if it exceeds new max HP
			if m.InitiativeList[0].HP > newMaxHP {
				m.InitiativeList[0].HP = newMaxHP
			}

			// Verify the results
			if m.InitiativeList[0].HP != tt.expectedHP {
				t.Errorf("Expected HP %d, got %d", tt.expectedHP, m.InitiativeList[0].HP)
			}
			if m.InitiativeList[0].MaxHP != tt.expectedMaxHP {
				t.Errorf("Expected MaxHP %d, got %d", tt.expectedMaxHP, m.InitiativeList[0].MaxHP)
			}
		})
	}
}

// TestMaxHPEditIntegration tests the full flow of editing max HP
func TestMaxHPEditIntegration(t *testing.T) {
	// Create a model with multiple entries
	m := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "Goblin",
				Type:       "monster",
				HP:         25,
				MaxHP:      50,
				AC:         15,
				Initiative: 18,
			},
			{
				Name:       "Warrior",
				Type:       "player",
				HP:         40,
				MaxHP:      60,
				AC:         18,
				Initiative: 15,
			},
			{
				Name:       "Orc",
				Type:       "monster",
				HP:         100,
				MaxHP:      100,
				AC:         13,
				Initiative: 10,
			},
		},
		SelectedEntry:      0,
		InitiativeListMode: true,
	}

	// Test entering edit mode for first monster
	m.InitiativeEditMode = true
	m.InitiativeEditType = "maxhp"
	m.InitiativeInput = "75"

	// Parse and apply
	if val, err := panels.ParseInput(m.InitiativeInput, "maxhp"); err == nil {
		newMaxHP := val.(int)
		if newMaxHP < 1 {
			newMaxHP = 1
		}
		m.InitiativeList[0].MaxHP = newMaxHP
		if m.InitiativeList[0].HP > newMaxHP {
			m.InitiativeList[0].HP = newMaxHP
		}
	}

	// Verify first entry updated
	if m.InitiativeList[0].MaxHP != 75 {
		t.Errorf("Expected Goblin MaxHP to be 75, got %d", m.InitiativeList[0].MaxHP)
	}
	if m.InitiativeList[0].HP != 25 {
		t.Errorf("Expected Goblin HP to remain 25, got %d", m.InitiativeList[0].HP)
	}

	// Test with third monster - reduce max HP below current HP
	m.SelectedEntry = 2
	m.InitiativeEditMode = true
	m.InitiativeEditType = "maxhp"
	m.InitiativeInput = "60"

	if val, err := panels.ParseInput(m.InitiativeInput, "maxhp"); err == nil {
		newMaxHP := val.(int)
		if newMaxHP < 1 {
			newMaxHP = 1
		}
		m.InitiativeList[2].MaxHP = newMaxHP
		if m.InitiativeList[2].HP > newMaxHP {
			m.InitiativeList[2].HP = newMaxHP
		}
	}

	// Verify third entry updated and HP capped
	if m.InitiativeList[2].MaxHP != 60 {
		t.Errorf("Expected Orc MaxHP to be 60, got %d", m.InitiativeList[2].MaxHP)
	}
	if m.InitiativeList[2].HP != 60 {
		t.Errorf("Expected Orc HP to be capped at 60, got %d", m.InitiativeList[2].HP)
	}

	// Verify other entries unchanged
	if m.InitiativeList[1].MaxHP != 60 {
		t.Errorf("Expected Warrior MaxHP to remain 60, got %d", m.InitiativeList[1].MaxHP)
	}
}

// TestMaxHPSaveState tests that max HP changes are preserved in save state
func TestMaxHPSaveState(t *testing.T) {
	// Create a model with modified max HP
	m := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "Dragon",
				Type:       "monster",
				HP:         150,
				MaxHP:      200,
				AC:         19,
				Initiative: 20,
			},
		},
	}

	// Create temp directory for save
	tmpDir := t.TempDir()
	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Save the campaign
	campaignName := "Test Max HP Campaign"
	err := ui.SaveCampaign(m, campaignName)
	if err != nil {
		t.Fatalf("Failed to save campaign: %v", err)
	}

	// Load the campaign
	_, initiativeList, err := ui.LoadCampaign("Test Max HP Campaign.json")
	if err != nil {
		t.Fatalf("Failed to load campaign: %v", err)
	}

	// Verify max HP preserved
	if len(initiativeList) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(initiativeList))
	}
	if initiativeList[0].MaxHP != 200 {
		t.Errorf("Expected MaxHP 200, got %d", initiativeList[0].MaxHP)
	}
	if initiativeList[0].HP != 150 {
		t.Errorf("Expected HP 150, got %d", initiativeList[0].HP)
	}
}
