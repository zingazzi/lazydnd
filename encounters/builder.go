// encounters/builder.go
package encounters

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EncounterMonster represents a monster in an encounter
type EncounterMonster struct {
	Name     string `json:"name"`
	CR       string `json:"cr"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"max_hp"`
	AC       int    `json:"ac"`
	Quantity int    `json:"quantity"`
	XP       int    `json:"xp"`
}

// Encounter represents a pre-built encounter
type Encounter struct {
	Name       string             `json:"name"`
	PartySize  int                `json:"party_size"`
	PartyLevel int                `json:"party_level"` // Average level if all same
	Monsters   []EncounterMonster `json:"monsters"`
	CreatedAt  time.Time          `json:"created_at"`
	ModifiedAt time.Time          `json:"modified_at"`
}

// Party represents the player party
type Party struct {
	Size  int   `json:"size"`
	Level int   `json:"level"` // For simplicity, assuming all same level
	// Future: Support individual levels per player
}

// GetTotalMonsterCount returns total number of monster instances
func (e *Encounter) GetTotalMonsterCount() int {
	total := 0
	for _, monster := range e.Monsters {
		total += monster.Quantity
	}
	return total
}

// GetMonsterCRs returns a slice of CRs for all monster instances
func (e *Encounter) GetMonsterCRs() []string {
	crs := []string{}
	for _, monster := range e.Monsters {
		for i := 0; i < monster.Quantity; i++ {
			crs = append(crs, monster.CR)
		}
	}
	return crs
}

// SaveEncounter saves an encounter template to disk
func SaveEncounter(enc *Encounter, saveDir string) error {
	if saveDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		saveDir = filepath.Join(homeDir, ".lazydnd", "encounters")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create encounters directory: %w", err)
	}

	// Sanitize filename
	filename := sanitizeFilename(enc.Name)
	if filename == "" {
		filename = "encounter"
	}
	filepath := filepath.Join(saveDir, filename+".json")

	// Update modified time
	enc.ModifiedAt = time.Now()

	// Marshal to JSON
	data, err := json.MarshalIndent(enc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal encounter: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write encounter file: %w", err)
	}

	return nil
}

// LoadEncounters loads all saved encounters from disk
func LoadEncounters(saveDir string) ([]Encounter, error) {
	if saveDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		saveDir = filepath.Join(homeDir, ".lazydnd", "encounters")
	}

	// Check if directory exists
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		return []Encounter{}, nil // No encounters yet
	}

	// Read directory
	files, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read encounters directory: %w", err)
	}

	encounters := []Encounter{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Read file
		filepath := filepath.Join(saveDir, file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			continue // Skip files that can't be read
		}

		// Unmarshal
		var enc Encounter
		if err := json.Unmarshal(data, &enc); err != nil {
			continue // Skip invalid JSON
		}

		encounters = append(encounters, enc)
	}

	return encounters, nil
}

// DeleteEncounter deletes an encounter template
func DeleteEncounter(name string, saveDir string) error {
	if saveDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		saveDir = filepath.Join(homeDir, ".lazydnd", "encounters")
	}

	filename := sanitizeFilename(name)
	filepath := filepath.Join(saveDir, filename+".json")

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete encounter: %w", err)
	}

	return nil
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	// Replace invalid characters with underscore
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return strings.ToLower(replacer.Replace(name))
}
