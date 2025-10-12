// ui/dice_presets.go
package ui

// getDefaultDiceMacros returns a map of preset macros for popular D&D spells
func getDefaultDiceMacros() map[string]string {
	return map[string]string{
		// Cantrips
		"eldritch_blast":  "1d10",
		"fire_bolt":       "1d10",
		"ray_of_frost":    "1d8",
		"sacred_flame":    "1d8",
		"toll_the_dead":   "1d12",

		// 1st Level Spells
		"magic_missile":   "3d4+3",
		"guiding_bolt":    "4d6",
		"inflict_wounds":  "3d10",
		"cure_wounds":     "1d8",
		"healing_word":    "1d4",
		"chromatic_orb":   "3d8",
		"burning_hands":   "3d6",
		"thunderwave":     "2d8",
		"ice_knife":       "1d10",

		// 2nd Level Spells
		"scorching_ray":   "2d6",
		"shatter":         "3d8",
		"spiritual_weapon": "1d8+3",

		// 3rd Level Spells
		"fireball":        "8d6",
		"lightning_bolt":  "8d6",
		"vampiric_touch":  "3d6",
		"spirit_guardians": "3d8",

		// 4th Level Spells
		"blight":          "8d8",
		"ice_storm":       "4d6",

		// 5th Level Spells
		"cone_of_cold":    "8d8",
		"flame_strike":    "8d6",

		// 6th Level+ Spells
		"chain_lightning": "10d8",
		"disintegrate":    "10d6+40",
		"finger_of_death": "7d8+30",

		// Common Attack Macros
		"sneak_attack_1d6": "1d20+1d6",
		"sneak_attack_2d6": "1d20+2d6",
		"sneak_attack_3d6": "1d20+3d6",
		"sneak_attack_4d6": "1d20+4d6",
		"divine_smite_1d8": "1d8",
		"divine_smite_2d8": "2d8",
		"divine_smite_3d8": "3d8",
	}
}

// mergeDefaultMacros merges default macros with user macros
// User macros take precedence over defaults
func mergeDefaultMacros(userMacros map[string]string) map[string]string {
	if userMacros == nil {
		return getDefaultDiceMacros()
	}

	// Start with defaults
	merged := getDefaultDiceMacros()

	// Overwrite with user macros
	for name, formula := range userMacros {
		merged[name] = formula
	}

	return merged
}
