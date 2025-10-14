// encounters/difficulty.go
package encounters

// XP thresholds per character level (D&D 5e DMG page 82)
var xpThresholds = map[int]map[string]int{
	1:  {"easy": 25, "medium": 50, "hard": 75, "deadly": 100},
	2:  {"easy": 50, "medium": 100, "hard": 150, "deadly": 200},
	3:  {"easy": 75, "medium": 150, "hard": 225, "deadly": 400},
	4:  {"easy": 125, "medium": 250, "hard": 375, "deadly": 500},
	5:  {"easy": 250, "medium": 500, "hard": 750, "deadly": 1100},
	6:  {"easy": 300, "medium": 600, "hard": 900, "deadly": 1400},
	7:  {"easy": 350, "medium": 750, "hard": 1100, "deadly": 1700},
	8:  {"easy": 450, "medium": 900, "hard": 1400, "deadly": 2100},
	9:  {"easy": 550, "medium": 1100, "hard": 1600, "deadly": 2400},
	10: {"easy": 600, "medium": 1200, "hard": 1900, "deadly": 2800},
	11: {"easy": 800, "medium": 1600, "hard": 2400, "deadly": 3600},
	12: {"easy": 1000, "medium": 2000, "hard": 3000, "deadly": 4500},
	13: {"easy": 1100, "medium": 2200, "hard": 3400, "deadly": 5100},
	14: {"easy": 1250, "medium": 2500, "hard": 3800, "deadly": 5700},
	15: {"easy": 1400, "medium": 2800, "hard": 4300, "deadly": 6400},
	16: {"easy": 1600, "medium": 3200, "hard": 4800, "deadly": 7200},
	17: {"easy": 2000, "medium": 3900, "hard": 5900, "deadly": 8800},
	18: {"easy": 2100, "medium": 4200, "hard": 6300, "deadly": 9500},
	19: {"easy": 2400, "medium": 4900, "hard": 7300, "deadly": 10900},
	20: {"easy": 2800, "medium": 5700, "hard": 8500, "deadly": 12700},
}

// CR to XP conversion (DMG page 275)
var crToXP = map[string]int{
	"0":     10,
	"1/8":   25,
	"1/4":   50,
	"1/2":   100,
	"1":     200,
	"2":     450,
	"3":     700,
	"4":     1100,
	"5":     1800,
	"6":     2300,
	"7":     2900,
	"8":     3900,
	"9":     5000,
	"10":    5900,
	"11":    7200,
	"12":    8400,
	"13":    10000,
	"14":    11500,
	"15":    13000,
	"16":    15000,
	"17":    18000,
	"18":    20000,
	"19":    22000,
	"20":    25000,
	"21":    33000,
	"22":    41000,
	"23":    50000,
	"24":    62000,
	"25":    75000,
	"26":    90000,
	"27":    105000,
	"28":    120000,
	"29":    135000,
	"30":    155000,
}

// GetEncounterMultiplier returns the multiplier based on number of monsters (DMG page 82)
func GetEncounterMultiplier(numMonsters int) float64 {
	switch {
	case numMonsters == 1:
		return 1.0
	case numMonsters == 2:
		return 1.5
	case numMonsters >= 3 && numMonsters <= 6:
		return 2.0
	case numMonsters >= 7 && numMonsters <= 10:
		return 2.5
	case numMonsters >= 11 && numMonsters <= 14:
		return 3.0
	case numMonsters >= 15:
		return 4.0
	default:
		return 1.0
	}
}

// EncounterAnalysis contains the result of difficulty calculation
type EncounterAnalysis struct {
	TotalMonsterXP  int
	AdjustedXP      int
	Multiplier      float64
	Difficulty      string // "Trivial", "Easy", "Medium", "Hard", "Deadly"
	EasyThreshold   int
	MediumThreshold int
	HardThreshold   int
	DeadlyThreshold int
	NextThreshold   string // What difficulty level is next
	XPToNext        int    // XP needed to reach next difficulty
}

// CalculateDifficulty determines encounter difficulty based on party and monsters
func CalculateDifficulty(partyLevels []int, monsterCRs []string) EncounterAnalysis {
	analysis := EncounterAnalysis{}

	// Handle empty party
	if len(partyLevels) == 0 {
		analysis.Difficulty = "Unknown"
		return analysis
	}

	// Calculate party XP thresholds
	totalEasy, totalMedium, totalHard, totalDeadly := 0, 0, 0, 0
	for _, level := range partyLevels {
		if level < 1 {
			level = 1
		}
		if level > 20 {
			level = 20
		}
		thresholds := xpThresholds[level]
		totalEasy += thresholds["easy"]
		totalMedium += thresholds["medium"]
		totalHard += thresholds["hard"]
		totalDeadly += thresholds["deadly"]
	}

	analysis.EasyThreshold = totalEasy
	analysis.MediumThreshold = totalMedium
	analysis.HardThreshold = totalHard
	analysis.DeadlyThreshold = totalDeadly

	// Handle no monsters
	if len(monsterCRs) == 0 {
		analysis.Difficulty = "None"
		analysis.NextThreshold = "Easy"
		analysis.XPToNext = totalEasy
		return analysis
	}

	// Calculate monster XP
	totalMonsterXP := 0
	for _, cr := range monsterCRs {
		if xp, exists := crToXP[cr]; exists {
			totalMonsterXP += xp
		}
	}
	analysis.TotalMonsterXP = totalMonsterXP

	// Apply encounter multiplier
	multiplier := GetEncounterMultiplier(len(monsterCRs))
	analysis.Multiplier = multiplier
	adjustedXP := int(float64(totalMonsterXP) * multiplier)
	analysis.AdjustedXP = adjustedXP

	// Determine difficulty
	if adjustedXP < totalEasy {
		analysis.Difficulty = "Trivial"
		analysis.NextThreshold = "Easy"
		analysis.XPToNext = totalEasy - adjustedXP
	} else if adjustedXP < totalMedium {
		analysis.Difficulty = "Easy"
		analysis.NextThreshold = "Medium"
		analysis.XPToNext = totalMedium - adjustedXP
	} else if adjustedXP < totalHard {
		analysis.Difficulty = "Medium"
		analysis.NextThreshold = "Hard"
		analysis.XPToNext = totalHard - adjustedXP
	} else if adjustedXP < totalDeadly {
		analysis.Difficulty = "Hard"
		analysis.NextThreshold = "Deadly"
		analysis.XPToNext = totalDeadly - adjustedXP
	} else {
		analysis.Difficulty = "Deadly"
		analysis.NextThreshold = "Beyond Deadly"
		analysis.XPToNext = 0
	}

	return analysis
}

// GetCRXP returns the XP value for a given CR
func GetCRXP(cr string) int {
	if xp, exists := crToXP[cr]; exists {
		return xp
	}
	return 0
}

// GetDifficultyColor returns a color code for the difficulty level
func GetDifficultyColor(difficulty string) string {
	switch difficulty {
	case "Trivial":
		return "#AAAAAA" // Gray
	case "Easy":
		return "#00FF00" // Green
	case "Medium":
		return "#FFD700" // Gold
	case "Hard":
		return "#FFA500" // Orange
	case "Deadly":
		return "#FF0000" // Red
	default:
		return "#FFFFFF" // White
	}
}
