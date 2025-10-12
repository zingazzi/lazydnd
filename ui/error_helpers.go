// ui/error_helpers.go
package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// SetError sets an error message to display in the UI
func SetError(m *Model, message string) {
	m.ErrorMessage = message
	m.ErrorVisible = true
}

// ClearError clears the current error message
func ClearError(m *Model) {
	m.ErrorMessage = ""
	m.ErrorVisible = false
}

// ClearErrorMsg is a message type for clearing errors after a delay
type ClearErrorMsg struct{}

// ShowErrorWithTimeout shows an error and automatically clears it after a delay
func ShowErrorWithTimeout(message string, duration time.Duration) tea.Cmd {
	return func() tea.Msg {
		// First, return a message to set the error
		return SetErrorMsg{Message: message}
	}
}

// SetErrorMsg is a message type for setting errors
type SetErrorMsg struct {
	Message string
}

// ScheduleClearError schedules an error to be cleared after a delay
func ScheduleClearError(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}
