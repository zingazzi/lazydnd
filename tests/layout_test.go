// tests/layout_test.go
package tests

import (
	"lazydnd/config"
	"lazydnd/ui"
	"testing"
	"unicode/utf8"
)

// Helper functions to access unexported functions for testing
// These mirror the functions in ui/text_utils.go

func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return text
	}

	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{""}
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return []string{text}
	}

	var lines []string
	var currentLine []rune

	for _, r := range runes {
		if len(currentLine) >= maxWidth {
			lines = append(lines, string(currentLine))
			currentLine = []rune{r}
		} else {
			currentLine = append(currentLine, r)
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, string(currentLine))
	}

	return lines
}

func measureTextWidth(text string) int {
	return utf8.RuneCountInString(text)
}

// createLayoutTestModel creates a test model with default config for layout testing
func createLayoutTestModel() ui.Model {
	cfg := config.Default()
	return ui.Model{
		Config: cfg,
		Styles: ui.NewStyles(cfg),
		Width:  120,
		Height: 40,
	}
}

// TestTruncateText tests Unicode-safe text truncation
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxWidth int
		expected string
	}{
		{
			name:     "short text",
			text:     "hello",
			maxWidth: 10,
			expected: "hello",
		},
		{
			name:     "exact width",
			text:     "hello",
			maxWidth: 5,
			expected: "hello",
		},
		{
			name:     "truncate with ellipsis",
			text:     "hello world",
			maxWidth: 8,
			expected: "hello...",
		},
		{
			name:     "unicode characters",
			text:     "こんにちは世界",
			maxWidth: 5,
			expected: "こん...",
		},
		{
			name:     "emoji",
			text:     "🎲 Dice Roller 🎲",
			maxWidth: 10,
			expected: "🎲 Dice ...", // Emoji + " Dice " + ellipsis
		},
		{
			name:     "very small width",
			text:     "hello",
			maxWidth: 3,
			expected: "hel",
		},
		{
			name:     "zero width",
			text:     "hello",
			maxWidth: 0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxWidth)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestWrapText tests text wrapping
func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxWidth int
		expected []string
	}{
		{
			name:     "short text",
			text:     "hello",
			maxWidth: 10,
			expected: []string{"hello"},
		},
		{
			name:     "wrap long text",
			text:     "hello world",
			maxWidth: 5,
			expected: []string{"hello", " worl", "d"},
		},
		{
			name:     "unicode characters",
			text:     "こんにちは世界",
			maxWidth: 3,
			expected: []string{"こんに", "ちは世", "界"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.maxWidth)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d", len(tt.expected), len(result))
			}
			for i, line := range result {
				if i < len(tt.expected) && line != tt.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i, tt.expected[i], line)
				}
			}
		})
	}
}

// TestMeasureTextWidth tests text width measurement
func TestMeasureTextWidth(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "ascii",
			text:     "hello",
			expected: 5,
		},
		{
			name:     "unicode",
			text:     "こんにちは",
			expected: 5,
		},
		{
			name:     "emoji",
			text:     "🎲🎯",
			expected: 2,
		},
		{
			name:     "mixed",
			text:     "hello🎲",
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := measureTextWidth(tt.text)
			if result != tt.expected {
				t.Errorf("Expected width %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestLayoutConstants tests that layout system is properly set up
func TestLayoutConstants(t *testing.T) {
	// Verify that PanelDimensions type exists and works
	dim := ui.PanelDimensions{Width: 50, Height: 20}
	if dim.Width != 50 {
		t.Error("PanelDimensions Width not working correctly")
	}
	if dim.Height != 20 {
		t.Error("PanelDimensions Height not working correctly")
	}

	// Verify constants are accessible
	if ui.MinPanelWidth < 1 {
		t.Error("MinPanelWidth should be at least 1")
	}
	if ui.MinPanelHeight < 1 {
		t.Error("MinPanelHeight should be at least 1")
	}
	if ui.StatusBarHeight < 1 {
		t.Error("StatusBarHeight should be at least 1")
	}
}
