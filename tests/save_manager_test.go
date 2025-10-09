// tests/save_manager_test.go
package tests

import (
	"encoding/json"
	"lazydnd/ui"
	"os"
	"path/filepath"
	"testing"
)

// TestSaveState tests the SaveState struct serialization
func TestSaveState(t *testing.T) {
	state := ui.SaveState{
		CampaignName: "Test Campaign",
		SavedAt:      "2024-01-01T12:00:00Z",
		InitiativeList: []ui.SavedInitiativeEntry{
			{
				Name:       "Fighter",
				Type:       "player",
				Initiative: 15,
				HP:         0,
				MaxHP:      0,
			},
			{
				Name:       "Goblin",
				Type:       "monster",
				Initiative: 10,
				HP:         7,
				MaxHP:      7,
			},
		},
		CurrentTurn:  1,
		RoundCounter: 2,
	}

	// Test JSON marshaling
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal SaveState: %v", err)
	}

	// Test JSON unmarshaling
	var loaded ui.SaveState
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal SaveState: %v", err)
	}

	// Verify all fields
	if len(loaded.InitiativeList) != len(state.InitiativeList) {
		t.Errorf("InitiativeList count = %d, want %d", len(loaded.InitiativeList), len(state.InitiativeList))
	}

	if loaded.CurrentTurn != state.CurrentTurn {
		t.Errorf("CurrentTurn = %d, want %d", loaded.CurrentTurn, state.CurrentTurn)
	}

	if loaded.RoundCounter != state.RoundCounter {
		t.Errorf("RoundCounter = %d, want %d", loaded.RoundCounter, state.RoundCounter)
	}

	if loaded.CampaignName != state.CampaignName {
		t.Errorf("CampaignName = %q, want %q", loaded.CampaignName, state.CampaignName)
	}

	// Verify specific entry details
	if loaded.InitiativeList[0].Name != "Fighter" {
		t.Errorf("First entry name = %q, want %q", loaded.InitiativeList[0].Name, "Fighter")
	}

	if loaded.InitiativeList[1].HP != 7 {
		t.Errorf("Second entry HP = %d, want %d", loaded.InitiativeList[1].HP, 7)
	}
}

// TestSaveLoadCycle tests full save and load cycle
func TestSaveLoadCycle(t *testing.T) {
	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "lazydnd_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test save state
	saveState := ui.SaveState{
		CampaignName: "Test Campaign",
		SavedAt:      "2024-01-01T12:00:00Z",
		InitiativeList: []ui.SavedInitiativeEntry{
			{
				Name:       "Wizard",
				Type:       "player",
				Initiative: 18,
				HP:         0,
				MaxHP:      0,
			},
			{
				Name:       "Orc",
				Type:       "monster",
				Initiative: 12,
				HP:         15,
				MaxHP:      15,
			},
		},
		CurrentTurn:  0,
		RoundCounter: 3,
	}

	// Save to temp file
	testFile := filepath.Join(tmpDir, "test_save.json")
	data, err := json.MarshalIndent(saveState, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	err = os.WriteFile(testFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Load from file
	loadedData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var loadedState ui.SaveState
	err = json.Unmarshal(loadedData, &loadedState)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify loaded data matches original
	if len(loadedState.InitiativeList) != len(saveState.InitiativeList) {
		t.Errorf("Loaded initiative count = %d, want %d", len(loadedState.InitiativeList), len(saveState.InitiativeList))
	}

	if loadedState.CurrentTurn != saveState.CurrentTurn {
		t.Errorf("Loaded CurrentTurn = %d, want %d", loadedState.CurrentTurn, saveState.CurrentTurn)
	}

	if loadedState.RoundCounter != saveState.RoundCounter {
		t.Errorf("Loaded RoundCounter = %d, want %d", loadedState.RoundCounter, saveState.RoundCounter)
	}

	if loadedState.InitiativeList[0].Name != "Wizard" {
		t.Errorf("Loaded first entry name = %q, want %q", loadedState.InitiativeList[0].Name, "Wizard")
	}

	if loadedState.InitiativeList[1].HP != 15 {
		t.Errorf("Loaded second entry HP = %d, want %d", loadedState.InitiativeList[1].HP, 15)
	}
}

// TestListSaves tests listing save files
func TestListSaves(t *testing.T) {
	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "lazydnd_test_saves_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test save files
	testFiles := []string{"save1.json", "save2.json", "not_json.txt"}
	for _, filename := range testFiles {
		path := filepath.Join(tmpDir, filename)
		err := os.WriteFile(path, []byte("{}"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// List files
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	// Count .json files
	jsonCount := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			jsonCount++
		}
	}

	if jsonCount != 2 {
		t.Errorf("Found %d .json files, want 2", jsonCount)
	}
}

// TestEmptySave tests saving empty state
func TestEmptySave(t *testing.T) {
	state := ui.SaveState{
		CampaignName:   "Empty Campaign",
		SavedAt:        "2024-01-01T12:00:00Z",
		InitiativeList: []ui.SavedInitiativeEntry{},
		CurrentTurn:    -1,
		RoundCounter:   0,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal empty state: %v", err)
	}

	var loaded ui.SaveState
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty state: %v", err)
	}

	if loaded.InitiativeList == nil {
		t.Error("Loaded InitiativeList is nil, want empty slice")
	}

	if loaded.CurrentTurn != -1 {
		t.Errorf("Loaded CurrentTurn = %d, want -1", loaded.CurrentTurn)
	}
}

// TestSaveWithMonsterData tests saving monsters with full data
func TestSaveWithMonsterData(t *testing.T) {
	state := ui.SaveState{
		CampaignName: "Monster Test",
		SavedAt:      "2024-01-01T12:00:00Z",
		InitiativeList: []ui.SavedInitiativeEntry{
			{
				Name:        "Goblin 1",
				Type:        "monster",
				Initiative:  10,
				HP:          7,
				MaxHP:       7,
				AC:          15,
				InstanceNum: 1,
				BaseName:    "Goblin",
				MonsterName: "Goblin",
			},
		},
		CurrentTurn:  0,
		RoundCounter: 1,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal with monster data: %v", err)
	}

	var loaded ui.SaveState
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal with monster data: %v", err)
	}

	if loaded.InitiativeList[0].InstanceNum != 1 {
		t.Errorf("InstanceNum = %d, want 1", loaded.InitiativeList[0].InstanceNum)
	}

	if loaded.InitiativeList[0].MonsterName != "Goblin" {
		t.Errorf("MonsterName = %q, want %q", loaded.InitiativeList[0].MonsterName, "Goblin")
	}
}
