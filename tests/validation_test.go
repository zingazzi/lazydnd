// tests/validation_test.go
package tests

import (
	"lazydnd/panels"
	"testing"
)

func TestValidation_PlayerName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid name", "Gandalf", false},
		{"Name with space", "Frodo Baggins", false},
		{"Name with apostrophe", "O'Brien", false},
		{"Name with hyphen", "Jean-Luc", false},
		{"Name with underscore", "Dark_Knight", false},
		{"Empty name", "", true},
		{"Only spaces", "   ", true},
		{"Too long", "ThisNameIsWayTooLongAndExceedsTheMaximumCharacterLimit", true},
		{"Invalid characters", "Player@123", true},
		{"Special characters", "Name$%^", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "player_name")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(player_name, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_Initiative(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid positive", "15", false},
		{"Valid zero", "0", false},
		{"Valid negative", "-5", false},
		{"Min value", "-10", false},
		{"Max value", "99", false},
		{"Too low", "-11", true},
		{"Too high", "100", true},
		{"Way too high", "999", true},
		{"Not a number", "abc", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "player_initiative")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(player_initiative, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_HP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid HP", "50", false},
		{"Zero HP", "0", false},
		{"Max HP", "9999", false},
		{"Too high", "10000", true},
		{"Way too high", "999999", true},
		{"Negative", "-10", true},
		{"Not a number", "fifty", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "monster_hp")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(monster_hp, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_AC(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid AC", "15", false},
		{"Zero AC", "0", false},
		{"High AC", "30", false},
		{"Max AC", "99", false},
		{"Too high", "100", true},
		{"Negative", "-5", true},
		{"Not a number", "heavy", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "player_ac")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(player_ac, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_HPChange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Healing", "+10", false},
		{"Damage", "-15", false},
		{"Zero", "0", false},
		{"Large healing", "+100", false},
		{"Large damage", "-100", false},
		{"Max value", "9999", false},
		{"Min value", "-9999", false},
		{"Too high", "10000", true},
		{"Too low", "-10000", true},
		{"Not a number", "heal", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "hp_change")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(hp_change, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_TempHP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid temp HP", "20", false},
		{"Zero (clear)", "0", false},
		{"Max temp HP", "9999", false},
		{"Negative", "-10", true},
		{"Too high", "10000", true},
		{"Not a number", "temp", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "temphp")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(temphp, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_MaxHP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid max HP", "100", false},
		{"Min value", "1", false},
		{"Max value", "9999", false},
		{"Zero", "0", true},
		{"Negative", "-10", true},
		{"Too high", "10000", true},
		{"Not a number", "max", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "maxhp")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(maxhp, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_MonsterInitiative(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid number", "10", false},
		{"Roll command lowercase", "r", false},
		{"Roll command uppercase", "R", false},
		{"Roll command word", "roll", false},
		{"Valid negative", "-5", false},
		{"Too high", "100", true},
		{"Too low", "-11", true},
		{"Invalid input", "abc", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := panels.ParseInput(tt.input, "monster_initiative")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInput(monster_initiative, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidation_EdgeCases(t *testing.T) {
	t.Run("Name with mixed characters", func(t *testing.T) {
		_, err := panels.ParseInput("Thr4x_D4rkb4n3", "player_name")
		if err != nil {
			t.Errorf("Should allow alphanumeric with underscore: %v", err)
		}
	})

	t.Run("Whitespace trimming", func(t *testing.T) {
		result, err := panels.ParseInput("  Gandalf  ", "player_name")
		if err != nil {
			t.Errorf("Should trim whitespace: %v", err)
		}
		if result != "Gandalf" {
			t.Errorf("Expected 'Gandalf', got '%v'", result)
		}
	})

	t.Run("HP exactly at limit", func(t *testing.T) {
		_, err := panels.ParseInput("9999", "monster_hp")
		if err != nil {
			t.Errorf("Should allow HP of exactly 9999: %v", err)
		}
	})

	t.Run("Initiative at boundaries", func(t *testing.T) {
		_, err1 := panels.ParseInput("-10", "player_initiative")
		_, err2 := panels.ParseInput("99", "player_initiative")
		if err1 != nil || err2 != nil {
			t.Errorf("Should allow initiative at boundaries -10 and 99")
		}
	})
}

func TestValidation_ErrorMessages(t *testing.T) {
	t.Run("Name too long error message", func(t *testing.T) {
		_, err := panels.ParseInput("ThisNameIsWayTooLongAndExceedsTheMaximumCharacterLimitForNames", "player_name")
		if err == nil {
			t.Error("Expected error for name too long")
		} else if err.Error() != "name too long (max 50 characters)" {
			t.Errorf("Expected specific error message, got: %v", err.Error())
		}
	})

	t.Run("Invalid characters error message", func(t *testing.T) {
		_, err := panels.ParseInput("Name@123", "player_name")
		if err == nil {
			t.Error("Expected error for invalid characters")
		} else if err.Error() != "name contains invalid characters" {
			t.Errorf("Expected specific error message, got: %v", err.Error())
		}
	})

	t.Run("Initiative out of range error message", func(t *testing.T) {
		_, err := panels.ParseInput("100", "player_initiative")
		if err == nil {
			t.Error("Expected error for initiative out of range")
		} else if err.Error() != "initiative must be -10 to 99" {
			t.Errorf("Expected specific error message, got: %v", err.Error())
		}
	})

	t.Run("HP too high error message", func(t *testing.T) {
		_, err := panels.ParseInput("10000", "monster_hp")
		if err == nil {
			t.Error("Expected error for HP too high")
		} else if err.Error() != "HP too high (max 9999)" {
			t.Errorf("Expected specific error message, got: %v", err.Error())
		}
	})
}
