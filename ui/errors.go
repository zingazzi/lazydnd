// ui/errors.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// AppError represents an application-level error with context
type AppError struct {
	Operation string // What operation was being performed
	Err       error  // The underlying error
	Context   string // Additional context
	UserMsg   string // User-friendly message
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Operation, e.Err.Error(), e.Context)
	}
	return fmt.Sprintf("%s: %s", e.Operation, e.Err.Error())
}

// UserMessage returns a user-friendly error message
func (e *AppError) UserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	// Fallback to operation + simplified error
	return fmt.Sprintf("%s failed", e.Operation)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(operation string, err error, userMsg string) *AppError {
	return &AppError{
		Operation: operation,
		Err:       err,
		UserMsg:   userMsg,
	}
}

// NewAppErrorWithContext creates a new AppError with additional context
func NewAppErrorWithContext(operation string, err error, context string, userMsg string) *AppError {
	return &AppError{
		Operation: operation,
		Err:       err,
		Context:   context,
		UserMsg:   userMsg,
	}
}

// Common error operations
const (
	OpSaveCampaign   = "Save campaign"
	OpLoadCampaign   = "Load campaign"
	OpDeleteCampaign = "Delete campaign"
	OpRenameCampaign = "Rename campaign"
	OpLoadMonsters   = "Load monsters"
	OpLoadSpells     = "Load spells"
	OpRollDice       = "Roll dice"
	OpParseInput     = "Parse input"
	OpValidateInput  = "Validate input"
	OpInitConfig     = "Initialize config"
	OpSaveConfig     = "Save config"
	OpLoadConfig     = "Load config"
)

// HandleError is a helper to handle errors consistently
// Returns a tea.Cmd that displays the error message
func HandleError(err error, fallbackMsg string) tea.Cmd {
	if err == nil {
		return nil
	}

	var userMsg string
	if appErr, ok := err.(*AppError); ok {
		userMsg = appErr.UserMessage()
	} else {
		userMsg = fallbackMsg
	}

	return func() tea.Msg {
		return SetErrorMsg{Message: userMsg}
	}
}

// WrapError wraps an error with an operation name and user message
func WrapError(operation string, err error, userMsg string) error {
	if err == nil {
		return nil
	}
	return NewAppError(operation, err, userMsg)
}

// WrapErrorf wraps an error with formatted user message
func WrapErrorf(operation string, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewAppError(operation, err, fmt.Sprintf(format, args...))
}

