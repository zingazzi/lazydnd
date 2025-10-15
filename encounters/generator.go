// encounters/generator.go
package encounters

import (
	"math/rand"
	"strings"
)

// Environment represents different terrain types
type Environment string

const (
	EnvAny       Environment = "Any"
	EnvForest    Environment = "Forest"
	EnvMountain  Environment = "Mountain"
	EnvDesert    Environment = "Desert"
	EnvSwamp     Environment = "Swamp"
	EnvUnderdark Environment = "Underdark"
	EnvUrban     Environment = "Urban"
	EnvCoast     Environment = "Coast"
	EnvArctic    Environment = "Arctic"
	EnvJungle    Environment = "Jungle"
	EnvPlains    Environment = "Plains"
)

// GetAllEnvironments returns all available environments
func GetAllEnvironments() []string {
	return []string{
		string(EnvAny),
		string(EnvForest),
		string(EnvMountain),
		string(EnvDesert),
		string(EnvSwamp),
		string(EnvUnderdark),
		string(EnvUrban),
		string(EnvCoast),
		string(EnvArctic),
		string(EnvJungle),
		string(EnvPlains),
	}
}

// MonsterInfo represents monster data for encounter generation
type MonsterInfo struct {
	Name        string
	CR          string
	HP          int
	MaxHP       int
	AC          int
	XP          int
	Meta        string
	Environment []string
}

// GenerateEncounterRequest contains parameters for encounter generation
type GenerateEncounterRequest struct {
	PartySize   int
	PartyLevel  int
	Difficulty  string      // "easy", "medium", "hard", "deadly"
	Environment Environment // Terrain filter
	Monsters    []MonsterInfo
}

// GenerateEncounterResult contains the generated encounter
type GenerateEncounterResult struct {
	Monsters       []EncounterMonster
	TotalXP        int
	AdjustedXP     int
	TargetXP       int
	Difficulty     string // Requested difficulty
	ActualDiff     string // Actual difficulty achieved
	Environment    string // Selected environment
	EnvironmentMsg string // Message about environment filtering
}

// GenerateEncounter creates a balanced encounter based on criteria
func GenerateEncounter(req GenerateEncounterRequest) GenerateEncounterResult {
	result := GenerateEncounterResult{
		Monsters:       []EncounterMonster{},
		Environment:    string(req.Environment),
		EnvironmentMsg: "",
	}

	// Calculate target XP based on difficulty
	partyLevels := make([]int, req.PartySize)
	for i := 0; i < req.PartySize; i++ {
		partyLevels[i] = req.PartyLevel
	}

	analysis := CalculateDifficulty(partyLevels, []string{})

	var targetXP int
	switch strings.ToLower(req.Difficulty) {
	case "easy":
		targetXP = analysis.EasyThreshold
	case "medium":
		targetXP = analysis.MediumThreshold
	case "hard":
		targetXP = analysis.HardThreshold
	case "deadly":
		targetXP = analysis.DeadlyThreshold
	default:
		targetXP = analysis.MediumThreshold
	}
	result.TargetXP = targetXP
	result.Difficulty = req.Difficulty

	// Filter monsters by environment and appropriate CR
	filteredMonsters := filterMonstersByEnvironmentAndCR(req.Monsters, req.Environment, req.PartyLevel)

	if len(filteredMonsters) == 0 {
		result.EnvironmentMsg = "No suitable monsters found for " + string(req.Environment)
		return result
	}

	// Generate encounter using different strategies
	result.Monsters = generateBalancedEncounter(filteredMonsters, targetXP, req.PartyLevel)

	// Calculate actual difficulty
	if len(result.Monsters) > 0 {
		crs := []string{}
		totalMonsterXP := 0
		for _, m := range result.Monsters {
			for i := 0; i < m.Quantity; i++ {
				crs = append(crs, m.CR)
				totalMonsterXP += m.XP
			}
		}

		finalAnalysis := CalculateDifficulty(partyLevels, crs)
		result.TotalXP = totalMonsterXP
		result.AdjustedXP = finalAnalysis.AdjustedXP
		result.ActualDiff = finalAnalysis.Difficulty
	}

	return result
}

// filterMonstersByEnvironmentAndCR filters monsters by environment and appropriate CR
func filterMonstersByEnvironmentAndCR(monsters []MonsterInfo, env Environment, partyLevel int) []MonsterInfo {
	filtered := []MonsterInfo{}

	// Determine CR range based on party level
	minCR := maxInt(0, partyLevel-4)
	maxCR := partyLevel + 4

	for _, monster := range monsters {
		// Check CR range
		cr := parseCRToInt(monster.CR)
		if cr < minCR || cr > maxCR {
			continue
		}

		// Check environment (if not "Any")
		if env != EnvAny {
			if !monsterMatchesEnvironment(monster, env) {
				continue
			}
		}

		filtered = append(filtered, monster)
	}

	return filtered
}

// generateBalancedEncounter creates a balanced encounter from filtered monsters
func generateBalancedEncounter(monsters []MonsterInfo, targetXP int, partyLevel int) []EncounterMonster {
	if len(monsters) == 0 {
		return []EncounterMonster{}
	}

	// Strategy 1: Single strong monster (boss fight)
	if rand.Float32() < 0.3 {
		boss := findMonsterNearXP(monsters, int(float64(targetXP)*0.7), partyLevel+2)
		if boss != nil {
			return []EncounterMonster{{
				Name:     boss.Name,
				CR:       boss.CR,
				HP:       boss.HP,
				MaxHP:    boss.MaxHP,
				AC:       boss.AC,
				Quantity: 1,
				XP:       boss.XP,
			}}
		}
	}

	// Strategy 2: Small group (2-4 monsters)
	if rand.Float32() < 0.5 {
		group := generateSmallGroup(monsters, targetXP, 2, 4)
		if len(group) > 0 {
			return group
		}
	}

	// Strategy 3: Horde (5+ monsters)
	group := generateHorde(monsters, targetXP, 5, 10)
	if len(group) > 0 {
		return group
	}

	// Fallback: Just pick a random appropriate monster
	if len(monsters) > 0 {
		monster := monsters[rand.Intn(len(monsters))]
		return []EncounterMonster{{
			Name:     monster.Name,
			CR:       monster.CR,
			HP:       monster.HP,
			MaxHP:    monster.MaxHP,
			AC:       monster.AC,
			Quantity: 1,
			XP:       monster.XP,
		}}
	}

	return []EncounterMonster{}
}

// generateSmallGroup creates a small group encounter
func generateSmallGroup(monsters []MonsterInfo, targetXP int, minCount, maxCount int) []EncounterMonster {
	count := minCount + rand.Intn(maxCount-minCount+1)

	// Account for multiplier (2-4 monsters = 2x multiplier)
	baseXP := int(float64(targetXP) / 2.0)
	perMonsterXP := baseXP / count

	monster := findMonsterNearXP(monsters, perMonsterXP, 999)
	if monster != nil {
		return []EncounterMonster{{
			Name:     monster.Name,
			CR:       monster.CR,
			HP:       monster.HP,
			MaxHP:    monster.MaxHP,
			AC:       monster.AC,
			Quantity: count,
			XP:       monster.XP,
		}}
	}

	return []EncounterMonster{}
}

// generateHorde creates a horde encounter
func generateHorde(monsters []MonsterInfo, targetXP int, minCount, maxCount int) []EncounterMonster {
	count := minCount + rand.Intn(maxCount-minCount+1)

	// Account for multiplier (7-10 monsters = 2.5x multiplier)
	baseXP := int(float64(targetXP) / 2.5)
	perMonsterXP := baseXP / count

	monster := findMonsterNearXP(monsters, perMonsterXP, 999)
	if monster != nil {
		return []EncounterMonster{{
			Name:     monster.Name,
			CR:       monster.CR,
			HP:       monster.HP,
			MaxHP:    monster.MaxHP,
			AC:       monster.AC,
			Quantity: count,
			XP:       monster.XP,
		}}
	}

	return []EncounterMonster{}
}

// findMonsterNearXP finds a monster with XP close to target
func findMonsterNearXP(monsters []MonsterInfo, targetXP int, maxCR int) *MonsterInfo {
	var best *MonsterInfo
	bestDiff := 999999

	for i := range monsters {
		m := &monsters[i]
		cr := parseCRToInt(m.CR)
		if cr > maxCR {
			continue
		}

		diff := abs(m.XP - targetXP)
		if diff < bestDiff {
			bestDiff = diff
			best = m
		}
	}

	return best
}

// monsterMatchesEnvironment checks if monster is suitable for environment
func monsterMatchesEnvironment(monster MonsterInfo, env Environment) bool {
	meta := strings.ToLower(monster.Meta)
	envStr := strings.ToLower(string(env))

	// Check meta field for environment keywords
	if strings.Contains(meta, envStr) {
		return true
	}

	// Additional environment matching logic
	switch env {
	case EnvForest:
		return strings.Contains(meta, "forest") ||
		       strings.Contains(meta, "wood") ||
		       strings.Contains(meta, "beast") ||
		       strings.Contains(meta, "fey")
	case EnvMountain:
		return strings.Contains(meta, "mountain") ||
		       strings.Contains(meta, "giant") ||
		       strings.Contains(meta, "dwarf")
	case EnvDesert:
		return strings.Contains(meta, "desert") ||
		       strings.Contains(meta, "sand") ||
		       strings.Contains(meta, "scorpion")
	case EnvSwamp:
		return strings.Contains(meta, "swamp") ||
		       strings.Contains(meta, "bog") ||
		       strings.Contains(meta, "lizard") ||
		       strings.Contains(meta, "frog")
	case EnvUnderdark:
		return strings.Contains(meta, "underdark") ||
		       strings.Contains(meta, "drow") ||
		       strings.Contains(meta, "duergar") ||
		       strings.Contains(meta, "aberration")
	case EnvUrban:
		return strings.Contains(meta, "humanoid") ||
		       strings.Contains(meta, "human") ||
		       strings.Contains(meta, "thief")
	case EnvCoast:
		return strings.Contains(meta, "aquatic") ||
		       strings.Contains(meta, "pirate") ||
		       strings.Contains(meta, "sahuagin")
	case EnvArctic:
		return strings.Contains(meta, "arctic") ||
		       strings.Contains(meta, "frost") ||
		       strings.Contains(meta, "ice") ||
		       strings.Contains(meta, "white")
	case EnvJungle:
		return strings.Contains(meta, "jungle") ||
		       strings.Contains(meta, "tropical") ||
		       strings.Contains(meta, "ape") ||
		       strings.Contains(meta, "snake")
	case EnvPlains:
		return strings.Contains(meta, "plains") ||
		       strings.Contains(meta, "grassland") ||
		       strings.Contains(meta, "horse")
	}

	return false
}

// parseCRToInt converts CR string to approximate integer for comparison
func parseCRToInt(cr string) int {
	switch cr {
	case "0":
		return 0
	case "1/8":
		return 0
	case "1/4":
		return 0
	case "1/2":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "11":
		return 11
	case "12":
		return 12
	case "13":
		return 13
	case "14":
		return 14
	case "15":
		return 15
	case "16":
		return 16
	case "17":
		return 17
	case "18":
		return 18
	case "19":
		return 19
	case "20":
		return 20
	default:
		// For higher CRs, try to parse
		// This is a simplification
		return 20
	}
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
