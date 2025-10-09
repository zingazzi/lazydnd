// tests/dice_validation_test.go
package tests

import (
	"lazydnd/panels"
	"strings"
	"testing"
)

// TestValidateDiceType tests the dice type validation function
func TestRollDiceWithModifiers(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantContains string
		wantError    bool
	}{
		{
			name:         "Positive modifier",
			command:      "1d20+5",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Negative modifier",
			command:      "1d20-3",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Large modifier",
			command:      "2d6+10",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Zero modifier (invalid syntax)",
			command:      "1d20+0",
			wantContains: "TOTAL:",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := panels.RollDice(tt.command)

			if tt.wantError && !strings.Contains(result, "Invalid") {
				t.Errorf("Expected error for command %q, got %q", tt.command, result)
			}

			if !tt.wantError && !strings.Contains(result, tt.wantContains) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", tt.command, result, tt.wantContains)
			}
		})
	}
}

// TestRollDiceInvalidInputs tests error handling
func TestRollDiceInvalidInputs(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantContains string
	}{
		{
			name:         "Invalid dice type d7",
			command:      "1d7",
			wantContains: "invalid",
		},
		{
			name:         "Invalid dice type d15",
			command:      "2d15",
			wantContains: "invalid",
		},
		{
			name:         "Too many dice",
			command:      "101d6",
			wantContains: "Invalid",
		},
		{
			name:         "Zero dice",
			command:      "0d20",
			wantContains: "Invalid",
		},
		{
			name:         "Negative dice",
			command:      "-1d20",
			wantContains: "Invalid",
		},
		{
			name:         "Empty command",
			command:      "",
			wantContains: "Invalid",
		},
		{
			name:         "Just text",
			command:      "hello",
			wantContains: "Invalid",
		},
		{
			name:         "Missing dice type",
			command:      "2d",
			wantContains: "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := panels.RollDice(tt.command)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", tt.command, result, tt.wantContains)
			}
		})
	}
}

// TestRollDiceComplexExpressions tests complex dice expressions
func TestRollDiceComplexExpressions(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantContains string
	}{
		{
			name:         "Two dice expressions",
			command:      "2d6+1d4",
			wantContains: "TOTAL:",
		},
		{
			name:         "Multiple expressions with modifiers",
			command:      "1d20+2d6+5",
			wantContains: "TOTAL:",
		},
		{
			name:         "Subtraction",
			command:      "3d6-1d4",
			wantContains: "TOTAL:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := panels.RollDice(tt.command)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", tt.command, result, tt.wantContains)
			}

			// Verify result is not an error
			if strings.Contains(result, "Invalid") {
				t.Errorf("RollDice(%q) returned error: %q", tt.command, result)
			}
		})
	}
}

// TestRollDiceMultipleRolls tests comma-separated rolls
func TestRollDiceMultipleRolls(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantLen int // Number of results expected
	}{
		{
			name:    "Two rolls",
			command: "2d6, 1d20",
			wantLen: 2,
		},
		{
			name:    "Three rolls",
			command: "1d4, 1d6, 1d8",
			wantLen: 3,
		},
		{
			name:    "Rolls with modifiers",
			command: "1d20+5, 2d6-1",
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := panels.RollDice(tt.command)

			// Count number of "|" separators
			separators := strings.Count(result, "|")
			if separators != tt.wantLen-1 {
				t.Errorf("RollDice(%q) returned %d results, want %d", tt.command, separators+1, tt.wantLen)
			}
		})
	}
}

// TestRollDiceAdvantageDisadvantage tests advantage/disadvantage mechanics
func TestRollDiceAdvantageDisadvantage(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantContains string
	}{
		{
			name:         "Advantage full word",
			command:      "1d20 advantage",
			wantContains: "ADV",
		},
		{
			name:         "Advantage abbreviation",
			command:      "1d20 adv",
			wantContains: "ADV",
		},
		{
			name:         "Disadvantage full word",
			command:      "1d20 disadvantage",
			wantContains: "DIS",
		},
		{
			name:         "Disadvantage abbreviation",
			command:      "1d20 dis",
			wantContains: "DIS",
		},
		{
			name:         "Advantage with modifier",
			command:      "1d20+5 adv",
			wantContains: "ADV",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := panels.RollDice(tt.command)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", tt.command, result, tt.wantContains)
			}
		})
	}
}

// TestRollDiceStandardDiceTypes tests all standard D&D dice
func TestRollDiceStandardDiceTypes(t *testing.T) {
	diceTypes := []string{"d4", "d6", "d8", "d10", "d12", "d20", "d100"}

	for _, diceType := range diceTypes {
		t.Run(diceType, func(t *testing.T) {
			result := panels.RollDice("1" + diceType)

			// Should not be an error
			if strings.Contains(result, "Invalid") {
				t.Errorf("RollDice(%q) returned error: %q", diceType, result)
			}

			// Should contain the dice notation
			if !strings.Contains(result, diceType) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", diceType, result, diceType)
			}
		})
	}
}

// TestRollInitiativeRange tests that initiative is always 1-20
func TestRollInitiativeRange(t *testing.T) {
	// Run 100 times to test randomness
	for i := 0; i < 100; i++ {
		result := panels.RollInitiative()

		if result < 1 || result > 20 {
			t.Errorf("RollInitiative() = %d, want between 1 and 20", result)
		}
	}
}

// TestRollInitiativeDistribution tests that all values are possible
func TestRollInitiativeDistribution(t *testing.T) {
	results := make(map[int]bool)

	// Run many times to hit all values
	for i := 0; i < 500; i++ {
		result := panels.RollInitiative()
		results[result] = true
	}

	// Check we got a reasonable spread (at least 15 different values)
	if len(results) < 15 {
		t.Errorf("RollInitiative() only produced %d different values in 500 rolls, expected at least 15", len(results))
	}
}

// TestParseInputPlayerName tests player name validation
func TestParseInputPlayerName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		wantValue string
	}{
		{
			name:      "Valid name",
			input:     "Gandalf",
			wantError: false,
			wantValue: "Gandalf",
		},
		{
			name:      "Name with spaces",
			input:     "Bilbo Baggins",
			wantError: false,
			wantValue: "Bilbo Baggins",
		},
		{
			name:      "Empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "Only spaces",
			input:     "   ",
			wantError: true,
		},
		{
			name:      "Name with special characters",
			input:     "Drizzt Do'Urden",
			wantError: false,
			wantValue: "Drizzt Do'Urden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := panels.ParseInput(tt.input, "player_name")

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseInput(%q, player_name) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseInput(%q, player_name) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.wantValue {
				t.Errorf("ParseInput(%q, player_name) = %v, want %v", tt.input, result, tt.wantValue)
			}
		})
	}
}

// TestParseInputInitiativeRoll tests the "r" auto-roll feature
func TestParseInputInitiativeRoll(t *testing.T) {
	// Test the "r" command multiple times
	for i := 0; i < 20; i++ {
		result, err := panels.ParseInput("r", "monster_initiative")

		if err != nil {
			t.Errorf("ParseInput('r', monster_initiative) unexpected error: %v", err)
			continue
		}

		// Should be an integer between 1-20
		if val, ok := result.(int); !ok || val < 1 || val > 20 {
			t.Errorf("ParseInput('r', monster_initiative) = %v, want int between 1-20", result)
		}
	}
}
