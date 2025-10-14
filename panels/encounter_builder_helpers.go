// panels/encounter_builder_helpers.go
package panels

// getMultiplier returns encounter multiplier based on monster count (D&D 5e DMG rules)
func getMultiplier(count int) float64 {
	if count == 1 {
		return 1.0
	} else if count == 2 {
		return 1.5
	} else if count <= 6 {
		return 2.0
	} else if count <= 10 {
		return 2.5
	} else if count <= 14 {
		return 3.0
	}
	return 4.0
}

// estimateDifficulty provides a rough difficulty estimate based on adjusted XP
func estimateDifficulty(partySize, partyLevel, adjustedXP int) string {
	// Rough XP thresholds per character (simplified from DMG)
	easyPerChar := 25 * partyLevel
	mediumPerChar := 50 * partyLevel
	hardPerChar := 75 * partyLevel
	deadlyPerChar := 100 * partyLevel

	totalEasy := easyPerChar * partySize
	totalMedium := mediumPerChar * partySize
	totalHard := hardPerChar * partySize
	totalDeadly := deadlyPerChar * partySize

	if adjustedXP < totalEasy {
		return "Trivial"
	} else if adjustedXP < totalMedium {
		return "Easy"
	} else if adjustedXP < totalHard {
		return "Medium"
	} else if adjustedXP < totalDeadly {
		return "Hard"
	}
	return "Deadly"
}

