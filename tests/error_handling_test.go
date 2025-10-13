// tests/error_handling_test.go
package tests

import (
	"fmt"
	"lazydnd/ui"
	"strings"
	"testing"
)

func TestAppError_Creation(t *testing.T) {
	err := fmt.Errorf("file not found")
	appErr := ui.NewAppError("Test Operation", err, "User friendly message")

	if appErr.Operation != "Test Operation" {
		t.Errorf("Expected operation 'Test Operation', got '%s'", appErr.Operation)
	}

	if appErr.UserMessage() != "User friendly message" {
		t.Errorf("Expected user message 'User friendly message', got '%s'", appErr.UserMessage())
	}

	if appErr.Unwrap() != err {
		t.Error("Unwrap should return the underlying error")
	}
}

func TestAppError_CreationWithContext(t *testing.T) {
	err := fmt.Errorf("permission denied")
	appErr := ui.NewAppErrorWithContext("Save File", err, "config.json", "Could not save")

	if appErr.Context != "config.json" {
		t.Errorf("Expected context 'config.json', got '%s'", appErr.Context)
	}

	errorStr := appErr.Error()
	if !strings.Contains(errorStr, "Save File") {
		t.Error("Error string should contain operation name")
	}
	if !strings.Contains(errorStr, "config.json") {
		t.Error("Error string should contain context")
	}
}

func TestAppError_ErrorString(t *testing.T) {
	err := fmt.Errorf("underlying error")
	appErr := ui.NewAppError("Operation", err, "User message")

	errorStr := appErr.Error()
	if !strings.Contains(errorStr, "Operation") {
		t.Error("Error string should contain operation")
	}
	if !strings.Contains(errorStr, "underlying error") {
		t.Error("Error string should contain underlying error message")
	}
}

func TestAppError_UserMessage(t *testing.T) {
	tests := []struct {
		name     string
		userMsg  string
		expected string
	}{
		{
			name:     "With user message",
			userMsg:  "Could not save campaign",
			expected: "Could not save campaign",
		},
		{
			name:     "Without user message",
			userMsg:  "",
			expected: "Test Operation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf("some error")
			appErr := ui.NewAppError("Test Operation", err, tt.userMsg)

			if appErr.UserMessage() != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, appErr.UserMessage())
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	err := fmt.Errorf("original error")
	wrapped := ui.WrapError("Test Operation", err, "User message")

	if wrapped == nil {
		t.Error("WrapError should not return nil for non-nil error")
	}

	appErr, ok := wrapped.(*ui.AppError)
	if !ok {
		t.Error("WrapError should return an *AppError")
	}

	if appErr.Operation != "Test Operation" {
		t.Errorf("Expected operation 'Test Operation', got '%s'", appErr.Operation)
	}
}

func TestWrapError_NilError(t *testing.T) {
	wrapped := ui.WrapError("Test Operation", nil, "User message")

	if wrapped != nil {
		t.Error("WrapError should return nil for nil error")
	}
}

func TestWrapErrorf(t *testing.T) {
	err := fmt.Errorf("original error")
	wrapped := ui.WrapErrorf("Test Operation", err, "Failed for '%s'", "campaign1")

	appErr, ok := wrapped.(*ui.AppError)
	if !ok {
		t.Fatal("WrapErrorf should return an *AppError")
	}

	if appErr.UserMessage() != "Failed for 'campaign1'" {
		t.Errorf("Expected formatted message 'Failed for 'campaign1'', got '%s'", appErr.UserMessage())
	}
}

func TestWrapErrorf_NilError(t *testing.T) {
	wrapped := ui.WrapErrorf("Test Operation", nil, "Failed for '%s'", "test")

	if wrapped != nil {
		t.Error("WrapErrorf should return nil for nil error")
	}
}

func TestHandleError_WithAppError(t *testing.T) {
	err := fmt.Errorf("underlying error")
	appErr := ui.NewAppError("Operation", err, "User friendly message")

	cmd := ui.HandleError(appErr, "Fallback message")

	if cmd == nil {
		t.Error("HandleError should return a command for non-nil error")
	}

	// Execute the command to get the message
	msg := cmd()

	setErrorMsg, ok := msg.(ui.SetErrorMsg)
	if !ok {
		t.Error("Command should return SetErrorMsg")
	}

	if setErrorMsg.Message != "User friendly message" {
		t.Errorf("Expected 'User friendly message', got '%s'", setErrorMsg.Message)
	}
}

func TestHandleError_WithRegularError(t *testing.T) {
	err := fmt.Errorf("regular error")

	cmd := ui.HandleError(err, "Fallback message")

	if cmd == nil {
		t.Error("HandleError should return a command for non-nil error")
	}

	// Execute the command to get the message
	msg := cmd()

	setErrorMsg, ok := msg.(ui.SetErrorMsg)
	if !ok {
		t.Error("Command should return SetErrorMsg")
	}

	if setErrorMsg.Message != "Fallback message" {
		t.Errorf("Expected fallback 'Fallback message', got '%s'", setErrorMsg.Message)
	}
}

func TestHandleError_NilError(t *testing.T) {
	cmd := ui.HandleError(nil, "Fallback message")

	if cmd != nil {
		t.Error("HandleError should return nil for nil error")
	}
}

func TestErrorOperationConstants(t *testing.T) {
	// Verify important operation constants exist and are non-empty
	operations := []string{
		ui.OpSaveCampaign,
		ui.OpLoadCampaign,
		ui.OpDeleteCampaign,
		ui.OpRenameCampaign,
		ui.OpLoadMonsters,
		ui.OpLoadSpells,
		ui.OpRollDice,
		ui.OpParseInput,
		ui.OpValidateInput,
		ui.OpInitConfig,
		ui.OpSaveConfig,
		ui.OpLoadConfig,
	}

	for _, op := range operations {
		if op == "" {
			t.Error("Operation constant should not be empty")
		}
	}
}

func TestErrorOperationUsage(t *testing.T) {
	// Test using standard operation constants
	err := fmt.Errorf("test error")
	appErr := ui.NewAppError(ui.OpSaveCampaign, err, "Could not save")

	if appErr.Operation != ui.OpSaveCampaign {
		t.Errorf("Expected operation '%s', got '%s'", ui.OpSaveCampaign, appErr.Operation)
	}
}

func TestAppError_Unwrap(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	appErr := ui.NewAppError("Operation", originalErr, "User message")

	unwrapped := appErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Unwrap should return the original error")
	}
}

func TestErrorPropagation(t *testing.T) {
	// Simulate error propagation through layers
	baseErr := fmt.Errorf("disk full")
	wrappedErr := ui.WrapError(ui.OpSaveCampaign, baseErr, "Could not save campaign")

	// Verify we can still access the base error
	appErr, ok := wrappedErr.(*ui.AppError)
	if !ok {
		t.Fatal("Should be an AppError")
	}

	if appErr.Unwrap() != baseErr {
		t.Error("Should be able to unwrap to base error")
	}
}

func TestErrorMessages_UserFriendly(t *testing.T) {
	// Test that user messages don't contain technical details
	err := fmt.Errorf("syscall.ENOSPC: no space left on device")
	appErr := ui.NewAppError(ui.OpSaveCampaign, err, "Could not save campaign - disk full")

	userMsg := appErr.UserMessage()

	// User message should be friendly
	if strings.Contains(userMsg, "syscall") {
		t.Error("User message should not contain technical details like 'syscall'")
	}

	// But the full error should contain details for logging
	fullErr := appErr.Error()
	if !strings.Contains(fullErr, "no space left") {
		t.Error("Full error should contain original error details")
	}
}
