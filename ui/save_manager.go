// ui/save_manager.go
package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lazydnd/panels"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// testSaveDir can be set during tests to override the default save directory
var testSaveDir string

// getSaveDirectory returns the path to the save directory in user's home
func getSaveDirectory() (string, error) {
	// Use test directory if set (for testing purposes)
	if testSaveDir != "" {
		return testSaveDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".lazydnd"), nil
}

// SetTestSaveDirectory sets a custom save directory for testing
func SetTestSaveDirectory(dir string) {
	testSaveDir = dir
}

// ClearTestSaveDirectory clears the test save directory setting
func ClearTestSaveDirectory() {
	testSaveDir = ""
}

// SaveCampaign saves the current campaign state to a JSON file
func SaveCampaign(m Model, campaignName string) error {
	// Get save directory path
	saveDir, err := getSaveDirectory()
	if err != nil {
		return err
	}

	// Create saves directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create saves directory: %w", err)
	}

	// Convert initiative list to saveable format
	savedEntries := make([]SavedInitiativeEntry, len(m.InitiativeList))
	for i, entry := range m.InitiativeList {
		monsterName := ""
		if entry.MonsterData != nil {
			monsterName = entry.MonsterData.Name
		} else if entry.MonsterName != "" {
			monsterName = entry.MonsterName
		}

		savedEntries[i] = SavedInitiativeEntry{
			Name:        entry.Name,
			Type:        entry.Type,
			Initiative:  entry.Initiative,
			HP:          entry.HP,
			MaxHP:       entry.MaxHP,
			TempHP:      entry.TempHP,
			AC:          entry.AC,
			MonsterName: monsterName,
			InstanceNum: entry.InstanceNum,
			BaseName:    entry.BaseName,
		}
	}

	// Create save state
	saveState := SaveState{
		CampaignName:   campaignName,
		SavedAt:        time.Now().Format(time.RFC3339),
		InitiativeList: savedEntries,
		CurrentTurn:    m.CurrentTurn,
		RoundCounter:   m.RoundCounter,
		ActiveSpells:   m.ActiveSpells,
		Notes:          m.NotesContent,
		DiceMacros:     m.DiceMacros,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(saveState, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save state: %w", err)
	}

	// Create filename from campaign name
	filename := sanitizeFilename(campaignName) + ".json"
	filePath := filepath.Join(saveDir, filename)

	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadCampaign loads a campaign state from a JSON file
func LoadCampaign(filename string) (SaveState, []InitiativeEntry, error) {
	var saveState SaveState

	// Get save directory path
	saveDir, err := getSaveDirectory()
	if err != nil {
		return saveState, nil, err
	}

	// Read file
	filePath := filepath.Join(saveDir, filename)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return saveState, nil, fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &saveState); err != nil {
		return saveState, nil, fmt.Errorf("failed to unmarshal save state: %w", err)
	}

	// Build monster map by searching for each unique monster name
	uniqueMonsterNames := make(map[string]bool)
	for _, saved := range saveState.InitiativeList {
		if saved.MonsterName != "" {
			uniqueMonsterNames[saved.MonsterName] = true
		}
	}

	monsterMap := make(map[string]*Monster)

	// Only load monsters if there are monsters to link
	if len(uniqueMonsterNames) > 0 {
		// Load monsters for re-linking
		if err := panels.LoadMonsters(); err != nil {
			// If monster loading fails, continue without monster data
			// This allows basic save/load to work without the assets directory
			uniqueMonsterNames = make(map[string]bool) // Clear monster names to skip linking
		}
	}

	// Find and convert each unique monster
	for monsterName := range uniqueMonsterNames {
		panelsMonster := panels.FindMonster(monsterName)
		if panelsMonster == nil {
			continue
		}

		// Convert panels.Monster to ui.Monster
		uiMonster := &Monster{
			Name:             panelsMonster.Name,
			Meta:             panelsMonster.Meta,
			ArmorClass:       panelsMonster.ArmorClass,
			HitPoints:        panelsMonster.HitPoints,
			Speed:            panelsMonster.Speed,
			STR:              panelsMonster.STR,
			STRMod:           panelsMonster.STRMod,
			DEX:              panelsMonster.DEX,
			DEXMod:           panelsMonster.DEXMod,
			CON:              panelsMonster.CON,
			CONMod:           panelsMonster.CONMod,
			INT:              panelsMonster.INT,
			INTMod:           panelsMonster.INTMod,
			WIS:              panelsMonster.WIS,
			WISMod:           panelsMonster.WISMod,
			CHA:              panelsMonster.CHA,
			CHAMod:           panelsMonster.CHAMod,
			SavingThrows:     panelsMonster.SavingThrows,
			Skills:           panelsMonster.Skills,
			Senses:           panelsMonster.Senses,
			Languages:        panelsMonster.Languages,
			Challenge:        panelsMonster.Challenge,
			Traits:           panelsMonster.Traits,
			Actions:          panelsMonster.Actions,
			LegendaryActions: panelsMonster.LegendaryActions,
			ImgURL:           panelsMonster.ImgURL,
			ActionNumber:     panelsMonster.ActionNumber,
		}

		// Convert actions
		uiMonster.ActionList = make([]MonsterAction, len(panelsMonster.ActionList))
		for j, action := range panelsMonster.ActionList {
			uiMonster.ActionList[j] = MonsterAction{
				Name:        action.Name,
				Type:        action.Type,
				Description: action.Description,
				Roll:        action.Roll,
				Reach:       action.Reach,
				Range:       action.Range,
				Damage:      action.Damage,
				DamageType:  action.DamageType,
				SaveDC:      action.SaveDC,
				SaveType:    action.SaveType,
			}
		}

		monsterMap[panelsMonster.Name] = uiMonster
	}

	// Convert saved entries back to initiative entries
	initiativeList := make([]InitiativeEntry, len(saveState.InitiativeList))
	for i, saved := range saveState.InitiativeList {
		entry := InitiativeEntry{
			Name:        saved.Name,
			Type:        saved.Type,
			Initiative:  saved.Initiative,
			HP:          saved.HP,
			MaxHP:       saved.MaxHP,
			AC:          saved.AC,
			InstanceNum: saved.InstanceNum,
			BaseName:    saved.BaseName,
			MonsterName: saved.MonsterName,
		}

		// Re-link monster data if available
		if saved.MonsterName != "" {
			if monster, exists := monsterMap[saved.MonsterName]; exists {
				entry.MonsterData = monster
			}
		}

		initiativeList[i] = entry
	}

	return saveState, initiativeList, nil
}

// ListCampaigns returns a list of all saved campaign files
func ListCampaigns() ([]string, error) {
	// Get save directory path
	saveDir, err := getSaveDirectory()
	if err != nil {
		return nil, err
	}

	// Create saves directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create saves directory: %w", err)
	}

	// Read directory
	files, err := ioutil.ReadDir(saveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read saves directory: %w", err)
	}

	// Filter JSON files
	var campaigns []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			campaigns = append(campaigns, file.Name())
		}
	}

	return campaigns, nil
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	// Replace invalid characters with underscore
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	// Trim spaces and limit length
	result = strings.TrimSpace(result)
	if len(result) > 100 {
		result = result[:100]
	}
	if result == "" {
		result = "campaign"
	}
	return result
}

// GetCampaignDisplayName extracts the campaign name from a filename
func GetCampaignDisplayName(filename string) string {
	// Remove .json extension
	name := strings.TrimSuffix(filename, ".json")
	// Replace underscores with spaces
	name = strings.ReplaceAll(name, "_", " ")
	return name
}
