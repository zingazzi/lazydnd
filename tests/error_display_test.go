// tests/error_display_test.go
package tests

import (
	"lazydnd/ui"
	"testing"
)

// TestErrorDisplay tests that errors are properly set and displayed
func TestErrorDisplay(t *testing.T) {
	m := ui.Model{}

	// Initially no error
	if m.ErrorVisible {
		t.Error("Error should not be visible initially")
	}
	if m.ErrorMessage != "" {
		t.Error("Error message should be empty initially")
	}

	// Set an error
	ui.SetError(&m, "Test error message")

	// Verify error is visible
	if !m.ErrorVisible {
		t.Error("Error should be visible after SetError")
	}
	if m.ErrorMessage != "Test error message" {
		t.Errorf("Expected 'Test error message', got '%s'", m.ErrorMessage)
	}

	// Clear error
	ui.ClearError(&m)

	// Verify error is cleared
	if m.ErrorVisible {
		t.Error("Error should not be visible after ClearError")
	}
	if m.ErrorMessage != "" {
		t.Errorf("Error message should be empty after ClearError, got '%s'", m.ErrorMessage)
	}
}

// TestErrorMessage tests the SetErrorMsg message type
func TestErrorMessage(t *testing.T) {
	msg := ui.SetErrorMsg{Message: "Test error"}

	if msg.Message != "Test error" {
		t.Errorf("Expected 'Test error', got '%s'", msg.Message)
	}
}

// TestClearErrorMessage tests the ClearErrorMsg message type
func TestClearErrorMessage(t *testing.T) {
	_ = ui.ClearErrorMsg{}
	// Just verify the type exists and can be created
}

// TestMultipleErrors tests setting multiple errors in sequence
func TestMultipleErrors(t *testing.T) {
	m := ui.Model{}

	// Set first error
	ui.SetError(&m, "First error")
	if m.ErrorMessage != "First error" {
		t.Errorf("Expected 'First error', got '%s'", m.ErrorMessage)
	}

	// Set second error (should replace first)
	ui.SetError(&m, "Second error")
	if m.ErrorMessage != "Second error" {
		t.Errorf("Expected 'Second error', got '%s'", m.ErrorMessage)
	}

	// Clear
	ui.ClearError(&m)
	if m.ErrorMessage != "" {
		t.Errorf("Error message should be empty after ClearError, got '%s'", m.ErrorMessage)
	}
}
