// ui/color_helpers.go
package ui

import (
	"lazydnd/config"
	"strings"
)

// HexToTViewColor converts a hex color to TView color tag format
// TView supports named colors like [red], [green], [yellow], etc.
// For custom hex colors, we map to the closest named color
func HexToTViewColor(hex string) string {
	hex = strings.ToUpper(strings.TrimSpace(hex))

	// Map common hex colors to TView named colors
	colorMap := map[string]string{
		// Red variants
		"#FF0000": "red",
		"#FF4500": "red",
		"#DC143C": "red",
		"#8B0000": "red",
		"#F92672": "red",

		// Green variants
		"#00FF00": "green",
		"#00FF7F": "green",
		"#228B22": "green",
		"#32CD32": "green",
		"#A6E22E": "green",

		// Yellow variants
		"#FFFF00": "yellow",
		"#FFD700": "yellow",
		"#FF8C00": "yellow",
		// Orange variants
		"#FFA500": "orange",
		"#FF7F00": "orange",

		// Blue variants
		"#0000FF": "blue",
		"#00BFFF": "blue",
		"#4169E1": "blue",
		"#1E90FF": "blue",
		"#7AA2F7": "blue",

		// Cyan variants
		"#00FFFF": "cyan",
		"#66D9EF": "cyan",

		// Magenta/Purple variants
		"#FF00FF": "magenta",
		"#7D56F4": "magenta",
		"#5A3D9E": "magenta",

		// White/Gray variants
		"#FFFFFF": "white",
		"#CCCCCC": "grey",
		"#AAAAAA": "grey",
		"#999999": "grey",
		"#666666": "grey",
		"#444444": "grey",
		"#333333": "grey",
	}

	// Check exact match first
	if color, ok := colorMap[hex]; ok {
		return color
	}

	// Default fallback based on hex value
	if strings.HasPrefix(hex, "#FF") || strings.HasPrefix(hex, "#DC") || strings.HasPrefix(hex, "#8B") {
		return "red"
	}
	if strings.HasPrefix(hex, "#00") && strings.Contains(hex, "FF") {
		return "green"
	}
	if strings.HasPrefix(hex, "#FF") && (strings.Contains(hex, "FF00") || strings.Contains(hex, "D700")) {
		return "yellow"
	}
	if strings.HasPrefix(hex, "#00") && strings.Contains(hex, "FFFF") {
		return "cyan"
	}

	// Default to white if no match
	return "white"
}

// GetDiceColors returns color tags for dice roller based on config
func GetDiceColors(cfg *config.Config) (critColor, goodRollColor, mediumRollColor string) {
	if cfg == nil {
		return "[red]", "[green]", "[yellow]"
	}

	critColor = "[" + HexToTViewColor(cfg.Theme.CritColor) + "]"
	goodRollColor = "[" + HexToTViewColor(cfg.Theme.GoodRollColor) + "]"
	mediumRollColor = "[" + HexToTViewColor(cfg.Theme.MediumRollColor) + "]"

	return critColor, goodRollColor, mediumRollColor
}

// GetHPColors returns color tags for HP display based on config
// Thresholds: >50% = healthy (grey), ≤50% and >20% = medium (orange), ≤20% = critical (red)
func GetHPColors(cfg *config.Config) (healthyColor, mediumColor, criticalColor, tempHPColor string) {
	if cfg == nil {
		return "[grey]", "[orange]", "[red]", "[cyan]"
	}

	healthyColor = "[" + HexToTViewColor(cfg.Theme.HPHealthyColor) + "]"
	mediumColor = "[" + HexToTViewColor(cfg.Theme.HPMediumColor) + "]"
	criticalColor = "[" + HexToTViewColor(cfg.Theme.HPCriticalColor) + "]"
	tempHPColor = "[" + HexToTViewColor(cfg.Theme.TempHPColor) + "]"

	return healthyColor, mediumColor, criticalColor, tempHPColor
}

// GetInitiativeNameColors returns color tags for monster and player names
func GetInitiativeNameColors(cfg *config.Config) (monsterNameColor, playerNameColor string) {
	if cfg == nil {
		return "[red]", "[green]"
	}

	// Use defaults if config colors are empty
	if cfg.Theme.MonsterNameColor == "" {
		monsterNameColor = "[red]"
	} else {
		monsterNameColor = "[" + HexToTViewColor(cfg.Theme.MonsterNameColor) + "]"
	}

	if cfg.Theme.PlayerNameColor == "" {
		playerNameColor = "[green]"
	} else {
		playerNameColor = "[" + HexToTViewColor(cfg.Theme.PlayerNameColor) + "]"
	}

	return monsterNameColor, playerNameColor
}

// GetTextColor returns color tag for regular text
func GetTextColor(cfg *config.Config) string {
	if cfg == nil {
		return "[grey]"
	}

	return "[" + HexToTViewColor(cfg.Theme.TextColor) + "]"
}
