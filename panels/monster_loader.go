// panels/monster_loader.go
package panels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
	"lazydnd/config"
)

// MonsterAction represents a single action a monster can take
type MonsterAction struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Roll        string `json:"roll,omitempty"`
	Reach       string `json:"reach,omitempty"`
	Range       string `json:"range,omitempty"`
	Damage      string `json:"damage,omitempty"`
	DamageType  string `json:"damage_type,omitempty"`
	SaveDC      string `json:"save_dc,omitempty"`
	SaveType    string `json:"save_type,omitempty"`
}

// Monster represents a D&D monster (duplicate to avoid import cycles)
type Monster struct {
	Name             string          `json:"name"`
	Meta             string          `json:"meta"`
	ArmorClass       string          `json:"Armor Class"`
	HitPoints        string          `json:"Hit Points"`
	Speed            string          `json:"Speed"`
	STR              string          `json:"STR"`
	STRMod           string          `json:"STR_mod"`
	DEX              string          `json:"DEX"`
	DEXMod           string          `json:"DEX_mod"`
	CON              string          `json:"CON"`
	CONMod           string          `json:"CON_mod"`
	INT              string          `json:"INT"`
	INTMod           string          `json:"INT_mod"`
	WIS              string          `json:"WIS"`
	WISMod           string          `json:"WIS_mod"`
	CHA              string          `json:"CHA"`
	CHAMod           string          `json:"CHA_mod"`
	SavingThrows     string          `json:"Saving Throws,omitempty"`
	Skills           string          `json:"Skills,omitempty"`
	Senses           string          `json:"Senses,omitempty"`
	Languages        string          `json:"Languages,omitempty"`
	Challenge        string          `json:"Challenge"`
	Traits           string          `json:"Traits,omitempty"`
	Actions          string          `json:"Actions,omitempty"`
	LegendaryActions string          `json:"Legendary Actions,omitempty"`
	ImgURL           string          `json:"img_url,omitempty"`
	ActionNumber     int             `json:"ActionNumber"`
	ActionList       []MonsterAction `json:"ActionList"`
}

var monsters []Monster
var monstersLoaded bool
var monsterMutex sync.Mutex

// ClearMonsterCache clears the cached monster data
func ClearMonsterCache() {
	monsterMutex.Lock()
	defer monsterMutex.Unlock()
	monsters = nil
	monstersLoaded = false
}

// ReloadMonsters forces a reload of monster data from disk
func ReloadMonsters() error {
	ClearMonsterCache()
	return LoadMonsters()
}

// IsMonstersLoaded returns true if monsters are currently cached
func IsMonstersLoaded() bool {
	return monstersLoaded
}

// GetMonsterCount returns the number of cached monsters
func GetMonsterCount() int {
	return len(monsters)
}

// GetMonsterByName searches for a monster by name
func GetMonsterByName(name string) *Monster {
	if err := LoadMonsters(); err != nil {
		return nil
	}

	monsterMutex.Lock()
	defer monsterMutex.Unlock()

	nameLower := strings.ToLower(name)
	for i := range monsters {
		if strings.ToLower(monsters[i].Name) == nameLower {
			return &monsters[i]
		}
	}

	return nil
}

// getCustomMonstersDir returns the path to the custom monsters directory
// Custom monsters are stored in the configured directory (default: ~/.lazydnd/monsters/)
func getCustomMonstersDir() (string, error) {
	// Load config to get custom monster directory
	cfg, err := config.Load()
	if err != nil {
		// If config fails to load, use default path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, ".lazydnd", "monsters"), nil
	}
	
	return cfg.GetMonsterDirectory()
}

// loadMonstersFromFile loads monsters from a specific JSON file
func loadMonstersFromFile(filePath string) ([]Monster, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var loadedMonsters []Monster
	err = json.Unmarshal(data, &loadedMonsters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	return loadedMonsters, nil
}

// loadCustomMonsters loads all custom monster files from the custom directory
func loadCustomMonsters() ([]Monster, error) {
	customDir, err := getCustomMonstersDir()
	if err != nil {
		return nil, err
	}

	// Check if directory exists
	if _, err := os.Stat(customDir); os.IsNotExist(err) {
		return []Monster{}, nil // No custom monsters directory, return empty list
	}

	// Read all JSON files from the directory
	files, err := filepath.Glob(filepath.Join(customDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list custom monster files: %w", err)
	}

	var customMonsters []Monster
	for _, file := range files {
		fileMonsters, err := loadMonstersFromFile(file)
		if err != nil {
			// Log error but continue loading other files
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			continue
		}
		customMonsters = append(customMonsters, fileMonsters...)
	}

	return customMonsters, nil
}

// mergeMonsters merges custom monsters with default monsters
// Custom monsters with the same name override default monsters
func mergeMonsters(defaultMonsters, customMonsters []Monster) []Monster {
	// Create a map for quick lookup
	monsterMap := make(map[string]Monster)

	// Add default monsters
	for _, monster := range defaultMonsters {
		monsterMap[strings.ToLower(monster.Name)] = monster
	}

	// Override with custom monsters
	for _, monster := range customMonsters {
		monsterMap[strings.ToLower(monster.Name)] = monster
	}

	// Convert map back to slice
	merged := make([]Monster, 0, len(monsterMap))
	for _, monster := range monsterMap {
		merged = append(merged, monster)
	}

	return merged
}

// LoadMonsters loads monsters from the default JSON file and custom monster files
func LoadMonsters() error {
	monsterMutex.Lock()
	defer monsterMutex.Unlock()

	if monstersLoaded && len(monsters) > 0 {
		return nil // Already loaded
	}

	// Load default monsters
	defaultMonsters, err := loadMonstersFromFile("assets/monsters.json")
	if err != nil {
		return err
	}

	// Load custom monsters
	customMonsters, err := loadCustomMonsters()
	if err != nil {
		// Log warning but don't fail - custom monsters are optional
		fmt.Fprintf(os.Stderr, "Warning: failed to load custom monsters: %v\n", err)
		monsters = defaultMonsters
		monstersLoaded = true
		return nil
	}

	// Merge default and custom monsters
	monsters = mergeMonsters(defaultMonsters, customMonsters)
	monstersLoaded = true

	// Print info about loaded monsters
	if len(customMonsters) > 0 {
		fmt.Fprintf(os.Stderr, "Loaded %d default monsters and %d custom monsters\n",
			len(defaultMonsters), len(customMonsters))
	}

	return nil
}

// parseCR converts CR string to a float64 for comparison
// Handles fractions like "1/8", "1/4", "1/2" and integers
func parseCR(crStr string) float64 {
	crStr = strings.TrimSpace(crStr)

	// Extract just the CR value (before any parentheses or XP info)
	if idx := strings.Index(crStr, "("); idx > 0 {
		crStr = strings.TrimSpace(crStr[:idx])
	}

	// Handle fractions
	if strings.Contains(crStr, "/") {
		parts := strings.Split(crStr, "/")
		if len(parts) == 2 {
			num, err1 := strconv.ParseFloat(parts[0], 64)
			den, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil && den != 0 {
				return num / den
			}
		}
	}

	// Handle regular numbers
	cr, err := strconv.ParseFloat(crStr, 64)
	if err != nil {
		return -1 // Invalid CR
	}
	return cr
}

// matchesCRFilter checks if a monster's CR matches the filter
// Filter formats: "5" (exactly 5), "0-5" (range), "10+" (10 and above)
func matchesCRFilter(monsterCR string, filter string) bool {
	if filter == "" {
		return true // No filter
	}

	cr := parseCR(monsterCR)
	if cr < 0 {
		return false // Invalid CR
	}

	filter = strings.TrimSpace(filter)

	// Handle "X+" format (e.g., "10+")
	if strings.HasSuffix(filter, "+") {
		minCR, err := strconv.ParseFloat(strings.TrimSuffix(filter, "+"), 64)
		if err != nil {
			return false
		}
		return cr >= minCR
	}

	// Handle "X-Y" format (e.g., "0-5")
	if strings.Contains(filter, "-") {
		parts := strings.Split(filter, "-")
		if len(parts) == 2 {
			minCR, err1 := strconv.ParseFloat(parts[0], 64)
			maxCR, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil {
				return cr >= minCR && cr <= maxCR
			}
		}
	}

	// Handle single value (e.g., "5")
	targetCR, err := strconv.ParseFloat(filter, 64)
	if err != nil {
		return false
	}
	return cr == targetCR
}

// SearchMonsters returns monster names that match the search term and CR filter
// GetAllMonstersForEncounter returns all monster names without limit (for encounter generation)
func GetAllMonstersForEncounter() []string {
	if err := LoadMonsters(); err != nil {
		return []string{}
	}

	names := make([]string, 0, len(monsters))
	for _, monster := range monsters {
		names = append(names, monster.Name)
	}
	return names
}

func SearchMonsters(searchTerm string, crFilter string) []string {
	if err := LoadMonsters(); err != nil {
		return []string{}
	}

	if searchTerm == "" && crFilter == "" {
		return []string{}
	}

	// Filter by CR first if specified
	filteredMonsters := monsters
	if crFilter != "" {
		filteredMonsters = make([]Monster, 0)
		for _, monster := range monsters {
			if matchesCRFilter(monster.Challenge, crFilter) {
				filteredMonsters = append(filteredMonsters, monster)
			}
		}
	}

	// If no search term, return filtered monsters (limited to 50 for better browsing)
	if searchTerm == "" {
		matches := make([]string, 0, 50)
		for i, monster := range filteredMonsters {
			if i >= 50 {
				break
			}
			matches = append(matches, monster.Name)
		}
		return matches
	}

	// Build list of filtered monster names for fuzzy search
	monsterNames := make([]string, 0, len(filteredMonsters))
	for _, monster := range filteredMonsters {
		monsterNames = append(monsterNames, monster.Name)
	}

	// Use fuzzy search to find matches
	results := fuzzy.Find(searchTerm, monsterNames)

	// Extract matched names (already sorted by score)
	matches := make([]string, 0, 50)
	for i, result := range results {
		if i >= 50 { // Limit to 50 suggestions for better browsing
			break
		}
		matches = append(matches, result.Str)
	}

	return matches
}

// FindMonster returns a monster by name
func FindMonster(name string) *Monster {
	if err := LoadMonsters(); err != nil {
		return nil
	}

	for _, monster := range monsters {
		if strings.EqualFold(monster.Name, name) {
			return &monster
		}
	}

	return nil
}

// FormatMonster formats monster details for display
func FormatMonster(monster *Monster) string {
	if monster == nil {
		return "Monster not found"
	}

	var lines []string

	// Header
	lines = append(lines, fmt.Sprintf("🐲 %s", monster.Name))
	lines = append(lines, monster.Meta)
	lines = append(lines, "")

	// Basic Stats
	lines = append(lines, fmt.Sprintf("AC: %s", monster.ArmorClass))
	lines = append(lines, fmt.Sprintf("HP: %s", monster.HitPoints))
	lines = append(lines, fmt.Sprintf("Speed: %s", monster.Speed))
	lines = append(lines, "")

	// Ability Scores
	lines = append(lines, "ABILITY SCORES:")
	lines = append(lines, fmt.Sprintf("STR: %s %s  DEX: %s %s  CON: %s %s",
		monster.STR, monster.STRMod, monster.DEX, monster.DEXMod, monster.CON, monster.CONMod))
	lines = append(lines, fmt.Sprintf("INT: %s %s  WIS: %s %s  CHA: %s %s",
		monster.INT, monster.INTMod, monster.WIS, monster.WISMod, monster.CHA, monster.CHAMod))
	lines = append(lines, "")

	// Optional fields
	if monster.SavingThrows != "" {
		lines = append(lines, fmt.Sprintf("Saving Throws: %s", monster.SavingThrows))
	}
	if monster.Skills != "" {
		lines = append(lines, fmt.Sprintf("Skills: %s", monster.Skills))
	}
	if monster.Senses != "" {
		lines = append(lines, fmt.Sprintf("Senses: %s", monster.Senses))
	}
	if monster.Languages != "" {
		lines = append(lines, fmt.Sprintf("Languages: %s", monster.Languages))
	}
	lines = append(lines, fmt.Sprintf("Challenge: %s", monster.Challenge))
	lines = append(lines, "")

	// Traits (clean HTML)
	if monster.Traits != "" {
		lines = append(lines, "TRAITS:")
		cleanTraits := cleanHTML(monster.Traits)
		wrappedTraits := wrapText(cleanTraits, 35) // Wrap at 35 characters to fit panel
		lines = append(lines, wrappedTraits...)
		lines = append(lines, "")
	}

	// Actions (clean HTML)
	if monster.Actions != "" {
		lines = append(lines, "ACTIONS:")
		cleanActions := cleanHTML(monster.Actions)
		wrappedActions := wrapText(cleanActions, 35) // Wrap at 35 characters to fit panel
		lines = append(lines, wrappedActions...)
		lines = append(lines, "")
	}

	// Legendary Actions (clean HTML)
	if monster.LegendaryActions != "" {
		lines = append(lines, "LEGENDARY ACTIONS:")
		cleanLegendary := cleanHTML(monster.LegendaryActions)
		wrappedLegendary := wrapText(cleanLegendary, 35) // Wrap at 35 characters to fit panel
		lines = append(lines, wrappedLegendary...)
	}

	return strings.Join(lines, "\n")
}

// cleanHTML removes HTML tags and formats text for terminal display
func cleanHTML(html string) string {
	// Simple HTML cleaning - remove tags and format
	text := html

	// Replace common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")

	// Remove HTML tags (simple approach)
	for strings.Contains(text, "<") && strings.Contains(text, ">") {
		start := strings.Index(text, "<")
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Clean up extra whitespace
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.TrimSpace(text)

	return text
}

// wrapText wraps text to specified line length
func wrapText(text string, maxWidth int) []string {
	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed the line length, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > maxWidth {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// ExtractMonsterStats extracts HP and AC from a monster for initiative tracker
func ExtractMonsterStats(monsterName string) (hp int, ac int, err error) {
	monster := FindMonster(monsterName)
	if monster == nil {
		return 0, 0, fmt.Errorf("monster not found: %s", monsterName)
	}

	// Parse HP from string like "58 (9d10 + 9)" - extract the first number
	hpStr := strings.TrimSpace(monster.HitPoints)
	if hpStr != "" {
		// Find the first number in the HP string
		var hpPart strings.Builder
		for _, char := range hpStr {
			if char >= '0' && char <= '9' {
				hpPart.WriteRune(char)
			} else if hpPart.Len() > 0 {
				break // Stop at first non-digit after we've started collecting digits
			}
		}
		if hpPart.Len() > 0 {
			if parsedHP, parseErr := strconv.Atoi(hpPart.String()); parseErr == nil {
				hp = parsedHP
			}
		}
	}

	// Parse AC from string like "16 (Natural Armor)" - extract the first number
	acStr := strings.TrimSpace(monster.ArmorClass)
	if acStr != "" {
		// Find the first number in the AC string
		var acPart strings.Builder
		for _, char := range acStr {
			if char >= '0' && char <= '9' {
				acPart.WriteRune(char)
			} else if acPart.Len() > 0 {
				break // Stop at first non-digit after we've started collecting digits
			}
		}
		if acPart.Len() > 0 {
			if parsedAC, parseErr := strconv.Atoi(acPart.String()); parseErr == nil {
				ac = parsedAC
			}
		}
	}

	if hp == 0 && ac == 0 {
		return 0, 0, fmt.Errorf("could not parse HP or AC for monster: %s", monsterName)
	}

	return hp, ac, nil
}
