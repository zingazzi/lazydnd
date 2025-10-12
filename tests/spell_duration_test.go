// tests/spell_duration_test.go
package tests

import (
	"lazydnd/ui"
	"testing"
)

// TestSpellDuration_Rounds tests parsing of round-based durations
func TestSpellDuration_Rounds(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
		isInstant      bool
	}{
		{"1 round", 1, false},
		{"2 rounds", 2, false},
		{"10 rounds", 10, false},
		{"5 round", 5, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("Expected instantaneous=%v, got %v", tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_Minutes tests parsing of minute-based durations
func TestSpellDuration_Minutes(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
		isInstant      bool
	}{
		{"1 minute", 10, false},
		{"2 minutes", 20, false},
		{"10 minutes", 100, false},
		{"5 minute", 50, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("Expected instantaneous=%v, got %v", tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_Hours tests parsing of hour-based durations
func TestSpellDuration_Hours(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
		isInstant      bool
	}{
		{"1 hour", 600, false},
		{"2 hours", 1200, false},
		{"8 hours", 4800, false},
		{"24 hours", 14400, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("Expected instantaneous=%v, got %v", tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_Days tests parsing of day-based durations
func TestSpellDuration_Days(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
		isInstant      bool
	}{
		{"1 day", 14400, false},
		{"2 days", 28800, false},
		{"7 days", 100800, false},
		{"10 days", 144000, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("Expected instantaneous=%v, got %v", tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_Concentration tests concentration prefix handling
func TestSpellDuration_Concentration(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
		isInstant      bool
	}{
		{"Concentration, up to 1 minute", 10, false},
		{"Concentration, up to 10 minutes", 100, false},
		{"Concentration, up to 1 hour", 600, false},
		{"up to 5 minutes", 50, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("Expected instantaneous=%v, got %v", tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_Instantaneous tests instantaneous spell handling
func TestSpellDuration_Instantaneous(t *testing.T) {
	testCases := []string{
		"Instantaneous",
		"instantaneous",
		"INSTANTANEOUS",
		"Special",
		"special",
		"",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(input)
			if rounds != 0 {
				t.Errorf("Expected 0 rounds for instantaneous spell, got %d", rounds)
			}
			if !isInstant {
				t.Error("Expected instantaneous=true")
			}
		})
	}
}

// TestSpellDuration_UnparsableFormats tests unparsable duration formats
func TestSpellDuration_UnparsableFormats(t *testing.T) {
	testCases := []string{
		"Until dispelled",
		"Permanent",
		"1 week",
		"Until dawn",
		"Random text",
		"123",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(input)
			if !isInstant {
				t.Errorf("Unparsable duration '%s' should be treated as instantaneous", input)
			}
			if rounds != 0 {
				t.Errorf("Unparsable duration should return 0 rounds, got %d", rounds)
			}
		})
	}
}

// TestSpellDuration_CaseInsensitivity tests case insensitive parsing
func TestSpellDuration_CaseInsensitivity(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
	}{
		{"1 MINUTE", 10},
		{"2 MiNuTes", 20},
		{"1 ROUND", 1},
		{"3 HoUrS", 1800},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant {
				t.Error("Expected non-instantaneous spell")
			}
		})
	}
}

// TestSpellDuration_Whitespace tests whitespace handling
func TestSpellDuration_Whitespace(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
	}{
		{"  1 minute  ", 10},
		{"\t2 rounds\n", 2},
		{"1minute", 10},
		{"5  rounds", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, _ := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
		})
	}
}

// TestSpellDuration_ZeroValues tests edge case of zero duration
func TestSpellDuration_ZeroValues(t *testing.T) {
	testCases := []string{
		"0 rounds",
		"0 minutes",
		"0 hours",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(input)
			if rounds != 0 {
				t.Errorf("Expected 0 rounds, got %d", rounds)
			}
			if isInstant {
				t.Error("0 duration should not be instantaneous")
			}
		})
	}
}

// TestSpellDuration_LargeValues tests very large duration values
func TestSpellDuration_LargeValues(t *testing.T) {
	testCases := []struct {
		input         string
		expectedRounds int
	}{
		{"99 rounds", 99},
		{"999 minutes", 9990},
		{"100 hours", 60000},
		{"30 days", 432000},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.input)
			if rounds != tc.expectedRounds {
				t.Errorf("Expected %d rounds, got %d", tc.expectedRounds, rounds)
			}
			if isInstant {
				t.Error("Expected non-instantaneous spell")
			}
		})
	}
}

// TestCastSpell tests the CastSpell function
func TestCastSpell(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{},
		RoundCounter: 5,
	}

	spell := &ui.Spell{
		Name:          "Bless",
		Duration:      "1 minute",
		Concentration: true,
	}

	model = ui.CastSpell(model, spell, "Cleric")

	if len(model.ActiveSpells) != 1 {
		t.Fatalf("Expected 1 active spell, got %d", len(model.ActiveSpells))
	}

	activeSpell := model.ActiveSpells[0]
	if activeSpell.Name != "Bless" {
		t.Errorf("Expected spell name 'Bless', got '%s'", activeSpell.Name)
	}
	if activeSpell.CasterName != "Cleric" {
		t.Errorf("Expected caster 'Cleric', got '%s'", activeSpell.CasterName)
	}
	if activeSpell.RoundsLeft != 10 {
		t.Errorf("Expected 10 rounds, got %d", activeSpell.RoundsLeft)
	}
	if activeSpell.TotalRounds != 10 {
		t.Errorf("Expected 10 total rounds, got %d", activeSpell.TotalRounds)
	}
	if !activeSpell.Concentration {
		t.Error("Expected concentration to be true")
	}
	if activeSpell.StartRound != 5 {
		t.Errorf("Expected start round 5, got %d", activeSpell.StartRound)
	}
}

// TestCastSpell_Instantaneous tests that instantaneous spells are not tracked
func TestCastSpell_Instantaneous(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{},
		RoundCounter: 1,
	}

	spell := &ui.Spell{
		Name:     "Fireball",
		Duration: "Instantaneous",
	}

	model = ui.CastSpell(model, spell, "Wizard")

	if len(model.ActiveSpells) != 0 {
		t.Errorf("Instantaneous spell should not be tracked, got %d active spells", len(model.ActiveSpells))
	}
}

// TestUpdateSpellDurations tests spell duration updates
func TestUpdateSpellDurations(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{
			{Name: "Bless", RoundsLeft: 3, TotalRounds: 10},
			{Name: "Haste", RoundsLeft: 1, TotalRounds: 10},
			{Name: "Shield of Faith", RoundsLeft: 5, TotalRounds: 100},
		},
		ActiveSpellIndex: 0,
	}

	newModel, expired := ui.UpdateSpellDurations(model)

	// Should have 2 remaining spells (Bless and Shield of Faith)
	if len(newModel.ActiveSpells) != 2 {
		t.Errorf("Expected 2 active spells, got %d", len(newModel.ActiveSpells))
	}

	// Haste should have expired
	if len(expired) != 1 || expired[0] != "Haste" {
		t.Errorf("Expected 'Haste' to expire, got %v", expired)
	}

	// Check that rounds were decremented
	for _, spell := range newModel.ActiveSpells {
		if spell.Name == "Bless" && spell.RoundsLeft != 2 {
			t.Errorf("Expected Bless to have 2 rounds left, got %d", spell.RoundsLeft)
		}
		if spell.Name == "Shield of Faith" && spell.RoundsLeft != 4 {
			t.Errorf("Expected Shield of Faith to have 4 rounds left, got %d", spell.RoundsLeft)
		}
	}
}

// TestUpdateSpellDurations_AllExpire tests all spells expiring at once
func TestUpdateSpellDurations_AllExpire(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{
			{Name: "Spell1", RoundsLeft: 1, TotalRounds: 10},
			{Name: "Spell2", RoundsLeft: 1, TotalRounds: 10},
		},
		ActiveSpellIndex: 0,
	}

	newModel, expired := ui.UpdateSpellDurations(model)

	if len(newModel.ActiveSpells) != 0 {
		t.Errorf("Expected 0 active spells, got %d", len(newModel.ActiveSpells))
	}

	if len(expired) != 2 {
		t.Errorf("Expected 2 expired spells, got %d", len(expired))
	}
}

// TestUpdateSpellDurations_EmptyList tests updating with no active spells
func TestUpdateSpellDurations_EmptyList(t *testing.T) {
	model := ui.Model{
		ActiveSpells:     []ui.ActiveSpell{},
		ActiveSpellIndex: -1,
	}

	newModel, expired := ui.UpdateSpellDurations(model)

	if len(newModel.ActiveSpells) != 0 {
		t.Errorf("Expected 0 active spells, got %d", len(newModel.ActiveSpells))
	}

	if len(expired) != 0 {
		t.Errorf("Expected 0 expired spells, got %d", len(expired))
	}
}

// TestRemoveActiveSpell tests removing a spell from the list
func TestRemoveActiveSpell(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{
			{Name: "Spell1", RoundsLeft: 10},
			{Name: "Spell2", RoundsLeft: 5},
			{Name: "Spell3", RoundsLeft: 8},
		},
		ActiveSpellIndex: 1,
	}

	// Remove spell at index 1 (Spell2)
	newModel := ui.RemoveActiveSpell(model, 1)

	if len(newModel.ActiveSpells) != 2 {
		t.Errorf("Expected 2 active spells, got %d", len(newModel.ActiveSpells))
	}

	// Verify remaining spells
	if newModel.ActiveSpells[0].Name != "Spell1" {
		t.Errorf("Expected first spell to be 'Spell1', got '%s'", newModel.ActiveSpells[0].Name)
	}
	if newModel.ActiveSpells[1].Name != "Spell3" {
		t.Errorf("Expected second spell to be 'Spell3', got '%s'", newModel.ActiveSpells[1].Name)
	}
}

// TestRemoveActiveSpell_InvalidIndex tests removing with invalid indices
func TestRemoveActiveSpell_InvalidIndex(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{
			{Name: "Spell1", RoundsLeft: 10},
		},
		ActiveSpellIndex: 0,
	}

	// Try to remove with negative index
	newModel := ui.RemoveActiveSpell(model, -1)
	if len(newModel.ActiveSpells) != 1 {
		t.Error("Should not remove spell with negative index")
	}

	// Try to remove with out-of-bounds index
	newModel = ui.RemoveActiveSpell(model, 10)
	if len(newModel.ActiveSpells) != 1 {
		t.Error("Should not remove spell with out-of-bounds index")
	}
}

// TestRemoveActiveSpell_LastSpell tests removing the last spell
func TestRemoveActiveSpell_LastSpell(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{
			{Name: "Spell1", RoundsLeft: 10},
		},
		ActiveSpellIndex: 0,
	}

	newModel := ui.RemoveActiveSpell(model, 0)

	if len(newModel.ActiveSpells) != 0 {
		t.Errorf("Expected 0 active spells, got %d", len(newModel.ActiveSpells))
	}

	if newModel.ActiveSpellIndex != 0 {
		t.Errorf("Expected index to be 0, got %d", newModel.ActiveSpellIndex)
	}
}

// TestSpellDuration_RealWorldExamples tests actual D&D spell durations
func TestSpellDuration_RealWorldExamples(t *testing.T) {
	testCases := []struct {
		spellName      string
		duration       string
		expectedRounds int
		isInstant      bool
	}{
		{"Fireball", "Instantaneous", 0, true},
		{"Shield", "1 round", 1, false},
		{"Bless", "Concentration, up to 1 minute", 10, false},
		{"Haste", "Concentration, up to 1 minute", 10, false},
		{"Greater Invisibility", "Concentration, up to 1 minute", 10, false},
		{"Mage Armor", "8 hours", 4800, false},
		{"Aid", "8 hours", 4800, false},
		{"Heroes' Feast", "24 hours", 14400, false},
		{"Mind Blank", "24 hours", 14400, false},
		{"Delayed Blast Fireball", "Concentration, up to 1 minute", 10, false},
		{"Detect Magic", "Concentration, up to 10 minutes", 100, false},
	}

	for _, tc := range testCases {
		t.Run(tc.spellName, func(t *testing.T) {
			rounds, isInstant := ui.ParseSpellDuration(tc.duration)
			if rounds != tc.expectedRounds {
				t.Errorf("%s: Expected %d rounds, got %d", tc.spellName, tc.expectedRounds, rounds)
			}
			if isInstant != tc.isInstant {
				t.Errorf("%s: Expected instantaneous=%v, got %v", tc.spellName, tc.isInstant, isInstant)
			}
		})
	}
}

// TestSpellDuration_MultipleSpells tests tracking multiple spells simultaneously
func TestSpellDuration_MultipleSpells(t *testing.T) {
	model := ui.Model{
		ActiveSpells: []ui.ActiveSpell{},
		RoundCounter: 1,
	}

	spells := []*ui.Spell{
		{Name: "Bless", Duration: "1 minute", Concentration: true},
		{Name: "Shield of Faith", Duration: "10 minutes", Concentration: false},
		{Name: "Mage Armor", Duration: "8 hours", Concentration: false},
	}

	casters := []string{"Cleric", "Cleric", "Wizard"}

	for i, spell := range spells {
		model = ui.CastSpell(model, spell, casters[i])
	}

	if len(model.ActiveSpells) != 3 {
		t.Fatalf("Expected 3 active spells, got %d", len(model.ActiveSpells))
	}

	// Verify durations
	if model.ActiveSpells[0].RoundsLeft != 10 {
		t.Errorf("Bless should have 10 rounds, got %d", model.ActiveSpells[0].RoundsLeft)
	}
	if model.ActiveSpells[1].RoundsLeft != 100 {
		t.Errorf("Shield of Faith should have 100 rounds, got %d", model.ActiveSpells[1].RoundsLeft)
	}
	if model.ActiveSpells[2].RoundsLeft != 4800 {
		t.Errorf("Mage Armor should have 4800 rounds, got %d", model.ActiveSpells[2].RoundsLeft)
	}
}
