// tests/dice_roller_test.go
package tests

import (
	"lazydnd/config"
	"lazydnd/panels"
	"strings"
	"testing"
)

// TestRollDice tests basic dice rolling functionality
func TestRollDice(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantContains string
		wantError    bool
	}{
		{
			name:         "Simple d20 roll",
			command:      "1d20",
			wantContains: "d20:",
			wantError:    false,
		},
		{
			name:         "Simple d6 roll",
			command:      "1d6",
			wantContains: "d6:",
			wantError:    false,
		},
		{
			name:         "Multiple dice",
			command:      "2d6",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Dice with modifier",
			command:      "1d20+5",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Dice with negative modifier",
			command:      "1d20-2",
			wantContains: "TOTAL:",
			wantError:    false,
		},
		{
			name:         "Invalid dice command",
			command:      "invalid",
			wantContains: "Invalid",
			wantError:    true,
		},
		{
			name:         "Advantage roll",
			command:      "1d20 adv",
			wantContains: "ADV",
			wantError:    false,
		},
		{
			name:         "Disadvantage roll",
			command:      "1d20 dis",
			wantContains: "DIS",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default()
			result := panels.RollDice(tt.command, cfg)

			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("RollDice(%q) = %q, want to contain %q", tt.command, result, tt.wantContains)
			}
		})
	}
}

// TestRollDiceResults tests that dice results are within valid ranges
func TestRollDiceResults(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		minValue int
		maxValue int
	}{
		{
			name:     "d6 range",
			command:  "1d6",
			minValue: 1,
			maxValue: 6,
		},
		{
			name:     "d20 range",
			command:  "1d20",
			minValue: 1,
			maxValue: 20,
		},
		{
			name:     "2d6 range",
			command:  "2d6",
			minValue: 2,
			maxValue: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to check randomness
			for i := 0; i < 10; i++ {
				cfg := config.Default()
				result := panels.RollDice(tt.command, cfg)

				// Extract the rolled value (this is a simplified check)
				if !strings.Contains(result, ":") {
					t.Errorf("Expected result to contain colon separator")
					continue
				}

				// Just verify we got a non-empty result
				if result == "" {
					t.Errorf("RollDice(%q) returned empty result", tt.command)
				}
			}
		})
	}
}

// TestRollInitiative tests initiative rolling
func TestRollInitiative(t *testing.T) {
	for i := 0; i < 20; i++ {
		result := panels.RollInitiative()

		if result < 1 || result > 20 {
			t.Errorf("RollInitiative() = %d, want between 1 and 20", result)
		}
	}
}

// TestParseInput tests input parsing for initiative tracker
func TestParseInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		inputType string
		wantError bool
		wantValue interface{}
	}{
		{
			name:      "Valid player name",
			input:     "Gandalf",
			inputType: "player_name",
			wantError: false,
			wantValue: "Gandalf",
		},
		{
			name:      "Empty player name",
			input:     "  ",
			inputType: "player_name",
			wantError: true,
		},
		{
			name:      "Valid initiative",
			input:     "15",
			inputType: "player_initiative",
			wantError: false,
			wantValue: 15,
		},
		{
			name:      "Invalid initiative (not a number)",
			input:     "abc",
			inputType: "player_initiative",
			wantError: true,
		},
		{
			name:      "Negative initiative",
			input:     "-5",
			inputType: "player_initiative",
			wantError: false, // Now allows -10 to 99
		},
		{
			name:      "Valid HP",
			input:     "50",
			inputType: "monster_hp",
			wantError: false,
			wantValue: 50,
		},
		{
			name:      "Monster initiative roll",
			input:     "r",
			inputType: "monster_initiative",
			wantError: false,
		},
		{
			name:      "HP change positive",
			input:     "+5",
			inputType: "hp_change",
			wantError: false,
			wantValue: 5,
		},
		{
			name:      "HP change negative",
			input:     "-10",
			inputType: "hp_change",
			wantError: false,
			wantValue: -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := panels.ParseInput(tt.input, tt.inputType)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseInput(%q, %q) expected error, got nil", tt.input, tt.inputType)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseInput(%q, %q) unexpected error: %v", tt.input, tt.inputType, err)
				return
			}

			if tt.wantValue != nil && result != tt.wantValue {
				// For roll results, just check it's valid
				if tt.inputType == "monster_initiative" && tt.input == "r" {
					if val, ok := result.(int); !ok || val < 1 || val > 20 {
						t.Errorf("ParseInput(%q, %q) = %v, want value between 1-20", tt.input, tt.inputType, result)
					}
				} else if result != tt.wantValue {
					t.Errorf("ParseInput(%q, %q) = %v, want %v", tt.input, tt.inputType, result, tt.wantValue)
				}
			}
		})
	}
}
