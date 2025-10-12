// tests/campaign_integration_test.go
package tests

import (
	"encoding/json"
	"lazydnd/ui"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestIntegration_FullSaveLoadCycle tests complete save and load workflow
func TestIntegration_FullSaveLoadCycle(t *testing.T) {
	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "lazydnd_integration_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override save directory for testing
	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Create test model with comprehensive data
	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "Gandalf",
				Type:       "player",
				Initiative: 18,
				HP:         0,
				MaxHP:      0,
				AC:         15,
				Conditions: []ui.Condition{
					{
						Name:        "Blessed",
						RoundsLeft:  10,
						TotalRounds: 10,
						Description: "Test blessing",
					},
				},
			},
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
				Conditions: []ui.Condition{
					{
						Name:        "Poisoned",
						RoundsLeft:  3,
						TotalRounds: 5,
						Description: "Poisoned by arrow",
					},
				},
			},
			{
				Name:        "Orc",
				Type:        "monster",
				Initiative:  12,
				HP:          15,
				MaxHP:       15,
				AC:          13,
				MonsterName: "Orc",
			},
		},
		CurrentTurn:  1,
		RoundCounter: 3,
		ActiveSpells: []ui.ActiveSpell{
			{
				Name:          "Bless",
				CasterName:    "Gandalf",
				RoundsLeft:    8,
				TotalRounds:   10,
				Concentration: true,
				StartRound:    1,
			},
		},
	}

	// Save campaign
	campaignName := "Test Campaign"
	err = ui.SaveCampaign(model, campaignName)
	if err != nil {
		t.Fatalf("Failed to save campaign: %v", err)
	}

	// Verify file was created (filename will be same as campaign name + .json)
	expectedFile := campaignName + ".json"
	savePath := filepath.Join(tmpDir, expectedFile)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Fatalf("Save file was not created at %s", savePath)
	}

	// Load campaign
	saveState, initiativeList, err := ui.LoadCampaign(expectedFile)
	if err != nil {
		t.Fatalf("Failed to load campaign: %v", err)
	}

	// Verify loaded data matches original
	if saveState.CampaignName != campaignName {
		t.Errorf("Loaded campaign name = %q, want %q", saveState.CampaignName, campaignName)
	}

	if saveState.CurrentTurn != model.CurrentTurn {
		t.Errorf("Loaded CurrentTurn = %d, want %d", saveState.CurrentTurn, model.CurrentTurn)
	}

	if saveState.RoundCounter != model.RoundCounter {
		t.Errorf("Loaded RoundCounter = %d, want %d", saveState.RoundCounter, model.RoundCounter)
	}

	if len(initiativeList) != len(model.InitiativeList) {
		t.Errorf("Loaded initiative count = %d, want %d", len(initiativeList), len(model.InitiativeList))
	}

	// Verify player data
	if initiativeList[0].Name != "Gandalf" {
		t.Errorf("Player name = %q, want %q", initiativeList[0].Name, "Gandalf")
	}
	if initiativeList[0].AC != 15 {
		t.Errorf("Player AC = %d, want %d", initiativeList[0].AC, 15)
	}

	// Verify monster data
	if initiativeList[1].Name != "Goblin 1" {
		t.Errorf("Monster name = %q, want %q", initiativeList[1].Name, "Goblin 1")
	}
	if initiativeList[1].HP != 7 {
		t.Errorf("Monster HP = %d, want %d", initiativeList[1].HP, 7)
	}
	if initiativeList[1].InstanceNum != 1 {
		t.Errorf("Monster InstanceNum = %d, want %d", initiativeList[1].InstanceNum, 1)
	}

	// Verify active spells
	if len(saveState.ActiveSpells) != 1 {
		t.Errorf("Active spells count = %d, want %d", len(saveState.ActiveSpells), 1)
	} else {
		spell := saveState.ActiveSpells[0]
		if spell.Name != "Bless" {
			t.Errorf("Spell name = %q, want %q", spell.Name, "Bless")
		}
		if spell.CasterName != "Gandalf" {
			t.Errorf("Caster name = %q, want %q", spell.CasterName, "Gandalf")
		}
		if spell.Concentration != true {
			t.Errorf("Concentration = %v, want %v", spell.Concentration, true)
		}
	}
}

// TestIntegration_MultipleConsecutiveSaves tests saving multiple times
func TestIntegration_MultipleConsecutiveSaves(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_multi_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "Fighter",
				Type:       "player",
				Initiative: 15,
			},
		},
		CurrentTurn:  0,
		RoundCounter: 1,
	}

	// Save three times with modifications
	for i := 1; i <= 3; i++ {
		model.RoundCounter = i
		err := ui.SaveCampaign(model, "Progressive Campaign")
		if err != nil {
			t.Fatalf("Save %d failed: %v", i, err)
		}

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Load and verify final state
	saveState, _, err := ui.LoadCampaign("Progressive Campaign.json")
	if err != nil {
		t.Fatalf("Failed to load campaign: %v", err)
	}

	if saveState.RoundCounter != 3 {
		t.Errorf("Final round counter = %d, want %d", saveState.RoundCounter, 3)
	}
}

// TestIntegration_ListCampaigns tests listing saved campaigns
func TestIntegration_ListCampaigns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_list_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Create multiple campaigns
	campaigns := []string{"Campaign One", "Campaign Two", "Test Campaign"}
	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{{Name: "Test", Type: "player", Initiative: 10}},
	}

	for _, name := range campaigns {
		err := ui.SaveCampaign(model, name)
		if err != nil {
			t.Fatalf("Failed to save %q: %v", name, err)
		}
	}

	// List campaigns
	list, err := ui.ListCampaigns()
	if err != nil {
		t.Fatalf("Failed to list campaigns: %v", err)
	}

	if len(list) != len(campaigns) {
		t.Errorf("Listed %d campaigns, want %d", len(list), len(campaigns))
	}

	// Verify all campaigns are in the list
	for _, expectedName := range campaigns {
		found := false
		expectedFilename := expectedName + ".json"
		for _, filename := range list {
			if filename == expectedFilename {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Campaign %q not found in list", expectedName)
		}
	}
}

// TestIntegration_SpecialCharactersInNames tests handling special characters
func TestIntegration_SpecialCharactersInNames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_special_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Campaign with special characters that should be sanitized
	specialNames := []string{
		"Campaign: The Beginning",
		"Player's Quest",
		"Test/Campaign",
		"Campaign?!",
		"<Test> Campaign",
	}

	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{Name: "Test Player", Type: "player", Initiative: 10},
		},
	}

	for _, name := range specialNames {
		err := ui.SaveCampaign(model, name)
		if err != nil {
			t.Errorf("Failed to save campaign with special name %q: %v", name, err)
			continue
		}

		// List campaigns to verify file was created
		list, err := ui.ListCampaigns()
		if err != nil {
			t.Errorf("Failed to list campaigns: %v", err)
			continue
		}

		if len(list) < 1 {
			t.Errorf("No campaigns found after saving %q", name)
		}
	}
}

// TestIntegration_EmptyCampaignSaveLoad tests saving and loading empty campaign
func TestIntegration_EmptyCampaignSaveLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_empty_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Empty model (no initiative entries)
	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{},
		CurrentTurn:    -1,
		RoundCounter:   0,
		ActiveSpells:   []ui.ActiveSpell{},
	}

	// Save empty campaign
	err = ui.SaveCampaign(model, "Empty Campaign")
	if err != nil {
		t.Fatalf("Failed to save empty campaign: %v", err)
	}

	// Load empty campaign
	saveState, initiativeList, err := ui.LoadCampaign("Empty Campaign.json")
	if err != nil {
		t.Fatalf("Failed to load empty campaign: %v", err)
	}

	if len(initiativeList) != 0 {
		t.Errorf("Loaded initiative count = %d, want 0", len(initiativeList))
	}

	if saveState.CurrentTurn != -1 {
		t.Errorf("CurrentTurn = %d, want -1", saveState.CurrentTurn)
	}

	if saveState.RoundCounter != 0 {
		t.Errorf("RoundCounter = %d, want 0", saveState.RoundCounter)
	}
}

// TestIntegration_LargeCampaignSaveLoad tests with many entries
func TestIntegration_LargeCampaignSaveLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_large_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Create large campaign with many entries
	model := ui.Model{
		InitiativeList: make([]ui.InitiativeEntry, 0),
		CurrentTurn:    5,
		RoundCounter:   10,
	}

	// Add 50 entries (mix of players and monsters)
	for i := 0; i < 50; i++ {
		var entry ui.InitiativeEntry
		if i%2 == 0 {
			entry = ui.InitiativeEntry{
				Name:       "Player " + string(rune('A'+i/2)),
				Type:       "player",
				Initiative: 20 - i,
				AC:         10 + i%10,
			}
		} else {
			entry = ui.InitiativeEntry{
				Name:        "Monster " + string(rune('1'+i/2)),
				Type:        "monster",
				Initiative:  20 - i,
				HP:          10 + i,
				MaxHP:       20 + i,
				AC:          10 + i%10,
				InstanceNum: (i / 2) + 1,
			}
		}
		model.InitiativeList = append(model.InitiativeList, entry)
	}

	// Save large campaign
	err = ui.SaveCampaign(model, "Large Campaign")
	if err != nil {
		t.Fatalf("Failed to save large campaign: %v", err)
	}

	// Load large campaign
	saveState, initiativeList, err := ui.LoadCampaign("Large Campaign.json")
	if err != nil {
		t.Fatalf("Failed to load large campaign: %v", err)
	}

	if len(initiativeList) != 50 {
		t.Errorf("Loaded entry count = %d, want 50", len(initiativeList))
	}

	if saveState.CurrentTurn != model.CurrentTurn {
		t.Errorf("CurrentTurn = %d, want %d", saveState.CurrentTurn, model.CurrentTurn)
	}

	// Verify a sample of entries
	if initiativeList[0].Name != "Player A" {
		t.Errorf("First entry name = %q, want %q", initiativeList[0].Name, "Player A")
	}

	if initiativeList[1].Name != "Monster 1" {
		t.Errorf("Second entry name = %q, want %q", initiativeList[1].Name, "Monster 1")
	}
}

// TestIntegration_CorruptedFileHandling tests error handling for corrupted files
func TestIntegration_CorruptedFileHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_corrupt_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	// Write corrupted JSON file
	corruptedFile := filepath.Join(tmpDir, "corrupted.json")
	err = os.WriteFile(corruptedFile, []byte("{ invalid json }"), 0644)
	if err != nil {
		t.Fatalf("Failed to write corrupted file: %v", err)
	}

	// Try to load corrupted file
	_, _, err = ui.LoadCampaign("corrupted.json")
	if err == nil {
		t.Error("Expected error when loading corrupted file, got nil")
	}
}

// TestIntegration_SaveStateStructure tests the JSON structure
func TestIntegration_SaveStateStructure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_structure_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{
				Name:       "TestPlayer",
				Type:       "player",
				Initiative: 15,
				AC:         18,
			},
		},
		CurrentTurn:  0,
		RoundCounter: 1,
		ActiveSpells: []ui.ActiveSpell{
			{
				Name:       "Shield",
				CasterName: "TestPlayer",
				RoundsLeft: 1,
			},
		},
	}

	// Save campaign
	err = ui.SaveCampaign(model, "Structure Test")
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Read raw JSON
	data, err := os.ReadFile(filepath.Join(tmpDir, "Structure Test.json"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Parse JSON to verify structure
	var rawData map[string]interface{}
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify top-level fields exist
	requiredFields := []string{"campaign_name", "saved_at", "initiative_list", "current_turn", "round_counter", "active_spells"}
	for _, field := range requiredFields {
		if _, exists := rawData[field]; !exists {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Verify initiative_list is an array
	initiativeList, ok := rawData["initiative_list"].([]interface{})
	if !ok {
		t.Fatal("initiative_list is not an array")
	}

	if len(initiativeList) != 1 {
		t.Errorf("initiative_list length = %d, want 1", len(initiativeList))
	}

	// Verify first entry structure
	if len(initiativeList) > 0 {
		entry, ok := initiativeList[0].(map[string]interface{})
		if !ok {
			t.Fatal("First entry is not an object")
		}

		entryFields := []string{"name", "type", "initiative", "ac"}
		for _, field := range entryFields {
			if _, exists := entry[field]; !exists {
				t.Errorf("Entry missing field: %s", field)
			}
		}
	}
}

// TestIntegration_ConcurrentSaves tests saving from multiple goroutines
func TestIntegration_ConcurrentSaves(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lazydnd_concurrent_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ui.SetTestSaveDirectory(tmpDir)
	defer ui.ClearTestSaveDirectory()

	model := ui.Model{
		InitiativeList: []ui.InitiativeEntry{
			{Name: "Test", Type: "player", Initiative: 10},
		},
	}

	// Save concurrently from multiple goroutines
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			err := ui.SaveCampaign(model, "Concurrent Test")
			if err != nil {
				t.Errorf("Goroutine %d save failed: %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify file was created and is valid
	_, _, err = ui.LoadCampaign("Concurrent Test.json")
	if err != nil {
		t.Errorf("Failed to load after concurrent saves: %v", err)
	}
}
