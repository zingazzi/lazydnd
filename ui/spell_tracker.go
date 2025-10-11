// ui/spell_tracker.go
package ui

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseSpellDuration converts spell duration string to number of combat rounds
// Each round is 6 seconds
func ParseSpellDuration(duration string) (rounds int, isInstantaneous bool) {
	duration = strings.ToLower(strings.TrimSpace(duration))

	// Check for instantaneous spells
	if strings.Contains(duration, "instantaneous") ||
	   strings.Contains(duration, "special") ||
	   duration == "" {
		return 0, true
	}

	// Remove concentration prefix if present
	duration = strings.TrimPrefix(duration, "concentration, up to ")
	duration = strings.TrimPrefix(duration, "up to ")

	// Parse different duration formats
	// Format: "X round(s)", "X minute(s)", "X hour(s)", "X day(s)"

	// Rounds (6 seconds each)
	if matches := regexp.MustCompile(`(\d+)\s*round`).FindStringSubmatch(duration); len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val, false
		}
	}

	// Minutes (10 rounds per minute)
	if matches := regexp.MustCompile(`(\d+)\s*minute`).FindStringSubmatch(duration); len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val * 10, false
		}
	}

	// Hours (600 rounds per hour)
	if matches := regexp.MustCompile(`(\d+)\s*hour`).FindStringSubmatch(duration); len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val * 600, false
		}
	}

	// Days (14400 rounds per day = 24 hours * 600 rounds/hour)
	if matches := regexp.MustCompile(`(\d+)\s*day`).FindStringSubmatch(duration); len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val * 14400, false
		}
	}

	// If we can't parse it, treat as instantaneous
	return 0, true
}

// CastSpell adds a new active spell to the list
func CastSpell(m Model, spell *Spell, casterName string) Model {
	rounds, isInstantaneous := ParseSpellDuration(spell.Duration)

	// Don't track instantaneous spells
	if isInstantaneous {
		return m
	}

	activeSpell := ActiveSpell{
		Name:          spell.Name,
		CasterName:    casterName,
		RoundsLeft:    rounds,
		TotalRounds:   rounds,
		Concentration: spell.Concentration,
		StartRound:    m.RoundCounter,
	}

	m.ActiveSpells = append(m.ActiveSpells, activeSpell)
	return m
}

// UpdateSpellDurations decrements all active spell durations by one round
// Returns updated model and list of expired spell names
func UpdateSpellDurations(m Model) (Model, []string) {
	var expired []string
	var activeSpells []ActiveSpell

	for _, spell := range m.ActiveSpells {
		spell.RoundsLeft--
		if spell.RoundsLeft <= 0 {
			expired = append(expired, spell.Name)
		} else {
			activeSpells = append(activeSpells, spell)
		}
	}

	m.ActiveSpells = activeSpells

	// Reset selection if it's out of bounds
	if m.ActiveSpellIndex >= len(m.ActiveSpells) {
		m.ActiveSpellIndex = len(m.ActiveSpells) - 1
	}
	if m.ActiveSpellIndex < 0 {
		m.ActiveSpellIndex = 0
	}

	return m, expired
}

// RemoveActiveSpell removes a spell from the active spells list
func RemoveActiveSpell(m Model, index int) Model {
	if index < 0 || index >= len(m.ActiveSpells) {
		return m
	}

	m.ActiveSpells = append(m.ActiveSpells[:index], m.ActiveSpells[index+1:]...)

	// Adjust selection
	if m.ActiveSpellIndex >= len(m.ActiveSpells) {
		m.ActiveSpellIndex = len(m.ActiveSpells) - 1
	}
	if m.ActiveSpellIndex < 0 {
		m.ActiveSpellIndex = 0
	}

	return m
}

// FormatActiveSpells returns a formatted string of all active spells
func FormatActiveSpells(activeSpells []ActiveSpell, selectedIndex int, isActive bool) string {
	if len(activeSpells) == 0 {
		return "No active spells"
	}

	var lines []string
	lines = append(lines, "ACTIVE SPELLS:")
	lines = append(lines, "")

	for i, spell := range activeSpells {
		// Calculate time remaining
		timeStr := formatRoundsToTime(spell.RoundsLeft)

		// Build spell line
		prefix := "  "
		if isActive && i == selectedIndex {
			prefix = "→ "
		}

		concentrationMark := ""
		if spell.Concentration {
			concentrationMark = " (C)"
		}

		line := prefix + spell.Name + concentrationMark
		lines = append(lines, line)
		lines = append(lines, "    Caster: "+spell.CasterName)
		lines = append(lines, "    Time left: "+timeStr)
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// formatRoundsToTime converts rounds to human-readable time
func formatRoundsToTime(rounds int) string {
	if rounds < 10 {
		seconds := rounds * 6
		return strconv.Itoa(seconds) + " seconds (" + strconv.Itoa(rounds) + " rounds)"
	}

	minutes := rounds / 10
	remainingRounds := rounds % 10

	if minutes >= 60 {
		hours := minutes / 60
		remainingMinutes := minutes % 60
		if remainingMinutes == 0 && remainingRounds == 0 {
			return strconv.Itoa(hours) + " hours"
		}
		return strconv.Itoa(hours) + "h " + strconv.Itoa(remainingMinutes) + "m"
	}

	if remainingRounds == 0 {
		return strconv.Itoa(minutes) + " minutes"
	}

	return strconv.Itoa(minutes) + " min " + strconv.Itoa(remainingRounds*6) + " sec"
}
