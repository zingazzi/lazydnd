// ui/validation.go
package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Input validation constants
const (
	MaxNameLength         = 50
	MaxCampaignNameLength = 50
	MaxNotesLength        = 10000
	MaxDiceCommandLength  = 100
	MaxConditionName      = 30

	MinHPValue         = 0
	MaxHPValue         = 9999
	MinACValue         = 0
	MaxACValue         = 99
	MinInitiativeValue = -10
	MaxInitiativeValue = 99
	MinSpellLevel      = 0
	MaxSpellLevel      = 9
)

// ValidateName validates character/monster names
// Returns sanitized name or error
func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	if len(name) > MaxNameLength {
		return "", fmt.Errorf("name too long (max %d characters)", MaxNameLength)
	}

	// Allow letters, numbers, spaces, hyphens, apostrophes, underscores
	validName := regexp.MustCompile(`^[a-zA-Z0-9 '\-_]+$`)
	if !validName.MatchString(name) {
		return "", fmt.Errorf("name contains invalid characters (use only letters, numbers, spaces, -, ', _)")
	}

	return name, nil
}

// ValidateCampaignName validates campaign save names
func ValidateCampaignName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", fmt.Errorf("campaign name cannot be empty")
	}

	if len(name) > MaxCampaignNameLength {
		return "", fmt.Errorf("campaign name too long (max %d characters)", MaxCampaignNameLength)
	}

	// Allow letters, numbers, spaces, hyphens, underscores
	validName := regexp.MustCompile(`^[a-zA-Z0-9 \-_]+$`)
	if !validName.MatchString(name) {
		return "", fmt.Errorf("campaign name contains invalid characters (use only letters, numbers, spaces, -, _)")
	}

	return name, nil
}

// ValidateHP validates hit point values
func ValidateHP(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("HP cannot be empty")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("HP must be a number")
	}

	if val < MinHPValue {
		return 0, fmt.Errorf("HP cannot be negative")
	}

	if val > MaxHPValue {
		return 0, fmt.Errorf("HP too high (max %d)", MaxHPValue)
	}

	return val, nil
}

// ValidateAC validates armor class values
func ValidateAC(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("AC cannot be empty")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("AC must be a number")
	}

	if val < MinACValue {
		return 0, fmt.Errorf("AC cannot be negative")
	}

	if val > MaxACValue {
		return 0, fmt.Errorf("AC too high (max %d)", MaxACValue)
	}

	return val, nil
}

// ValidateInitiative validates initiative values
func ValidateInitiative(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("initiative cannot be empty")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("initiative must be a number")
	}

	if val < MinInitiativeValue {
		return 0, fmt.Errorf("initiative too low (min %d)", MinInitiativeValue)
	}

	if val > MaxInitiativeValue {
		return 0, fmt.Errorf("initiative too high (max %d)", MaxInitiativeValue)
	}

	return val, nil
}

// ValidateHPChange validates HP modification values (+heal/-damage)
func ValidateHPChange(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("enter a number (+ to heal, - to damage)")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("must be a number (+ to heal, - to damage)")
	}

	// Allow reasonable range for damage/healing
	if val < -MaxHPValue || val > MaxHPValue {
		return 0, fmt.Errorf("value too extreme (max %d)", MaxHPValue)
	}

	return val, nil
}

// ValidateTempHP validates temporary hit points
func ValidateTempHP(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("enter a number (0 to clear)")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("temp HP must be a number")
	}

	if val < 0 {
		return 0, fmt.Errorf("temp HP cannot be negative")
	}

	if val > MaxHPValue {
		return 0, fmt.Errorf("temp HP too high (max %d)", MaxHPValue)
	}

	return val, nil
}

// ValidateDiceCommand validates dice notation commands
func ValidateDiceCommand(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", fmt.Errorf("dice command cannot be empty")
	}

	if len(input) > MaxDiceCommandLength {
		return "", fmt.Errorf("dice command too long (max %d characters)", MaxDiceCommandLength)
	}

	// Basic validation: check for dice notation or macro
	validDice := regexp.MustCompile(`^[0-9d+\-\*, advcritxgroup=_a-zA-Z]+$`)
	if !validDice.MatchString(input) {
		return "", fmt.Errorf("invalid dice command format")
	}

	return input, nil
}

// ValidateSpellLevel validates spell level values (0-9)
func ValidateSpellLevel(input string) (int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, fmt.Errorf("spell level cannot be empty")
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("spell level must be a number")
	}

	if val < MinSpellLevel || val > MaxSpellLevel {
		return 0, fmt.Errorf("spell level must be 0-9")
	}

	return val, nil
}

// ValidateCR validates challenge rating values
func ValidateCR(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", fmt.Errorf("CR cannot be empty")
	}

	// Allow numbers, fractions (1/8, 1/4, 1/2), ranges (1-5), and comparisons (5+, 10-)
	validCR := regexp.MustCompile(`^([0-9]+(/[0-9]+)?|\d+-\d+|\d+\+|\d+-)$`)
	if !validCR.MatchString(input) {
		return "", fmt.Errorf("invalid CR format (use: 1, 1/4, 1-5, 5+, 10-)")
	}

	return input, nil
}

// ValidateConditionName validates condition names
func ValidateConditionName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", fmt.Errorf("condition name cannot be empty")
	}

	if len(name) > MaxConditionName {
		return "", fmt.Errorf("condition name too long (max %d characters)", MaxConditionName)
	}

	// Allow letters, spaces, hyphens
	validName := regexp.MustCompile(`^[a-zA-Z \-]+$`)
	if !validName.MatchString(name) {
		return "", fmt.Errorf("condition name contains invalid characters")
	}

	return name, nil
}

// ValidateDuration validates duration values (positive integers or "permanent")
func ValidateDuration(input string) (string, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return "", fmt.Errorf("duration cannot be empty")
	}

	if input == "permanent" || input == "perm" {
		return "permanent", nil
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("duration must be a number or 'permanent'")
	}

	if val < 1 {
		return "", fmt.Errorf("duration must be at least 1")
	}

	if val > 1000 {
		return "", fmt.Errorf("duration too long (max 1000 rounds)")
	}

	return input, nil
}

// ValidateNotes validates notes content
func ValidateNotes(input string) (string, error) {
	if len(input) > MaxNotesLength {
		return "", fmt.Errorf("notes too long (max %d characters)", MaxNotesLength)
	}

	return input, nil
}

// SanitizeFilename removes/replaces characters that could cause filesystem issues
func SanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove any characters that aren't alphanumeric, underscore, or hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	name = reg.ReplaceAllString(name, "")

	// Truncate if too long (leave room for .json extension)
	if len(name) > MaxCampaignNameLength-5 {
		name = name[:MaxCampaignNameLength-5]
	}

	return name
}
