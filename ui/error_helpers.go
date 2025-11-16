// ui/error_helpers.go
package ui

import (
	"time"
)

// SetError sets an error message to display in the UI
// Uses the Global state struct for error display
func SetError(m *Model, message string) {
	m.Global.ErrorMessage = message
	m.Global.ErrorVisible = true
	// Keep legacy fields in sync during migration
	m.ErrorMessage = message
	m.ErrorVisible = true
}

// ClearError clears the current error message
func ClearError(m *Model) {
	m.Global.ErrorMessage = ""
	m.Global.ErrorVisible = false
	// Keep legacy fields in sync during migration
	m.ErrorMessage = ""
	m.ErrorVisible = false
}

// SetSuccess displays a success message (uses the error display system with green text)
func SetSuccess(m *Model, message string) {
	// We'll use the error message system but with a success prefix
	successMsg := "✓ " + message
	m.Global.ErrorMessage = successMsg
	m.Global.ErrorVisible = true
	// Keep legacy fields in sync during migration
	m.ErrorMessage = successMsg
	m.ErrorVisible = true
}

// ClearErrorMsg is a message type for clearing errors after a delay
type ClearErrorMsg struct{}

// ShowErrorWithTimeout shows an error and automatically clears it after a delay
// For TView, this just sets the error (timeout handled by TView app)
func ShowErrorWithTimeout(m *Model, message string, duration time.Duration) {
	SetError(m, message)
	// TView app will handle timeout clearing
}

// SetErrorMsg is a message type for setting errors
type SetErrorMsg struct {
	Message string
}

// ScheduleClearError schedules an error to be cleared after a delay
// For TView, this is handled by the app's timer system
func ScheduleClearError(m *Model, duration time.Duration) {
	// TView app will handle timeout clearing via its timer
	// This is a no-op for now - timeout handled in TView app
}
