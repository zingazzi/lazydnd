// ui/debug_logger.go
package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var debugFile *os.File
var debugEnabled bool

// InitDebugLogger initializes the debug log file
func InitDebugLogger() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "lazydnd")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	debugPath := filepath.Join(configDir, "debug.log")

	// Open file in append mode, create if doesn't exist
	debugFile, err = os.OpenFile(debugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	debugEnabled = false // Start disabled, toggle with F12
	return nil
}

// CloseDebugLogger closes the debug log file
func CloseDebugLogger() {
	if debugFile != nil {
		debugFile.Close()
	}
}

// DebugLog writes a debug message to the log file if debug mode is enabled
func DebugLog(format string, args ...interface{}) {
	if !debugEnabled || debugFile == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	debugFile.WriteString(logLine)
	debugFile.Sync() // Flush immediately for real-time viewing
}

// EnableDebugMode enables debug logging
func EnableDebugMode() {
	debugEnabled = true
	DebugLog("========================================")
	DebugLog("DEBUG MODE ENABLED")
	DebugLog("========================================")
}

// DisableDebugMode disables debug logging
func DisableDebugMode() {
	DebugLog("DEBUG MODE DISABLED")
	DebugLog("========================================")
	debugEnabled = false
}

// IsDebugEnabled returns the current debug state
func IsDebugEnabled() bool {
	return debugEnabled
}

// LogModelState logs the current model state
func LogModelState(m Model, context string) {
	if !debugEnabled {
		return
	}

	DebugLog("--- MODEL STATE: %s ---", context)
	DebugLog("  ActivePanel: %d (0=Dice, 1=Initiative, 2=Spells, 3=Monsters, 4=Debug)", m.ActivePanel)
	DebugLog("  InputMode: %v", m.InputMode)
	DebugLog("  SpellSearchMode: %v", m.SpellSearchMode)
	DebugLog("  MonsterSearchMode: %v", m.MonsterSearchMode)
	DebugLog("  InitiativeInputMode: %v", m.InitiativeInputMode)
	DebugLog("  InitiativeListMode: %v", m.InitiativeListMode)
	DebugLog("  InitiativeEditMode: %v", m.InitiativeEditMode)
	DebugLog("  MultiTargetMode: %v", m.MultiTargetMode)
	DebugLog("  SelectedEntry: %d", m.SelectedEntry)
	DebugLog("  CurrentTurn: %d", m.CurrentTurn)
	DebugLog("  RoundCounter: %d", m.RoundCounter)
	DebugLog("  InitiativeList length: %d", len(m.InitiativeList))
	DebugLog("  SelectedTargets: %d items", len(m.SelectedTargets))
	DebugLog("  ActiveSpells: %d items", len(m.ActiveSpells))
	DebugLog("  ShowMultiTargetPopup: %v", m.ShowMultiTargetPopup)
	DebugLog("  ShowConditionPopup: %v", m.ShowConditionPopup)
	DebugLog("  ShowActionPopup: %v", m.ShowActionPopup)
}

// LogKeyPress logs a key press event
func LogKeyPress(key string, context string) {
	if !debugEnabled {
		return
	}
	DebugLog("KEY: %s | Context: %s", key, context)
}

// LogCondition logs condition-related operations
func LogCondition(operation string, details string) {
	if !debugEnabled {
		return
	}
	DebugLog("CONDITION: %s | %s", operation, details)
}

// LogMultiTarget logs multi-target operations
func LogMultiTarget(operation string, details string) {
	if !debugEnabled {
		return
	}
	DebugLog("MULTI-TARGET: %s | %s", operation, details)
}

// LogSpell logs spell-related operations
func LogSpell(operation string, details string) {
	if !debugEnabled {
		return
	}
	DebugLog("SPELL: %s | %s", operation, details)
}
