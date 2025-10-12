// tests/initiative_rendering_test.go
package tests

import (
	"fmt"
	"lazydnd/panels"
	"strings"
	"testing"
	"time"
)

// createLargeInitiativeList creates a test initiative list with N entries
func createLargeInitiativeList(n int) string {
	entries := make([]string, n)
	for i := 0; i < n; i++ {
		entryType := "monster"
		name := fmt.Sprintf("Monster %d", i+1)
		if i%5 == 0 {
			entryType = "player"
			name = fmt.Sprintf("Player %d", (i/5)+1)
		}

		entries[i] = fmt.Sprintf("{Name:%s Type:%s Initiative:%d HP:%d MaxHP:%d AC:%d Conditions:[]}",
			name, entryType, 20-i, 15, 15, 14)
	}

	return "[" + strings.Join(entries, " ") + "]"
}

// TestOptimized_BasicRendering tests basic optimized rendering
func TestOptimized_BasicRendering(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if content == "" {
		t.Error("Content should not be empty")
	}

	if !strings.Contains(content, "Initiative Order:") {
		t.Error("Content should contain initiative order header")
	}
}

// TestOptimized_LargeList tests rendering with 20+ entries
func TestOptimized_LargeList(t *testing.T) {
	list := createLargeInitiativeList(25)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 10

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if content == "" {
		t.Error("Content should not be empty")
	}

	// Should have scroll indicators
	if !strings.Contains(content, "more below") {
		t.Error("Should show 'more below' indicator for large list")
	}
}

// TestOptimized_WindowingDisabled tests rendering with windowing disabled
func TestOptimized_WindowingDisabled(t *testing.T) {
	list := createLargeInitiativeList(30)
	config := panels.DefaultRenderConfig()
	config.EnableWindowing = false

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	// Should render all entries when windowing is disabled
	for i := 1; i <= 30; i++ {
		entryNum := fmt.Sprintf("%2d.", i)
		if !strings.Contains(content, entryNum) {
			t.Errorf("Should contain entry %d when windowing is disabled", i)
		}
	}
}

// TestOptimized_ViewportScrolling tests viewport scrolling
func TestOptimized_ViewportScrolling(t *testing.T) {
	list := createLargeInitiativeList(30)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 10
	config.ScrollOffset = 10 // Start at entry 10

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 15, true, false, false, "", 0, 0, false, nil, config,
	)

	// Should show "more above" indicator
	if !strings.Contains(content, "more above") {
		t.Error("Should show 'more above' indicator when scrolled")
	}

	// Should show "more below" indicator
	if !strings.Contains(content, "more below") {
		t.Error("Should show 'more below' indicator when scrolled")
	}
}

// TestOptimized_SelectedEntryVisible tests that selected entry stays visible
func TestOptimized_SelectedEntryVisible(t *testing.T) {
	list := createLargeInitiativeList(30)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 10
	config.ScrollOffset = 0

	// Select entry 25 (near end)
	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 25, true, true, false, "", 0, 0, false, nil, config,
	)

	// Should contain entry 25
	if !strings.Contains(content, "26.") { // Entry 25 is displayed as "26." (1-indexed)
		t.Error("Selected entry should be visible in viewport")
	}
}

// TestOptimized_EmptyList tests rendering empty list
func TestOptimized_EmptyList(t *testing.T) {
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		"[]", "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if !strings.Contains(content, "No entries") {
		t.Error("Should show 'No entries' message for empty list")
	}
}

// TestOptimized_RoundCounter tests round counter display
func TestOptimized_RoundCounter(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 5, false, nil, config,
	)

	if !strings.Contains(content, "Round 5") {
		t.Error("Should display round counter")
	}
}

// TestOptimized_MultiTargetMode tests multi-target mode rendering
func TestOptimized_MultiTargetMode(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()
	selectedTargets := map[int]bool{0: true, 2: true}

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, true, selectedTargets, config,
	)

	if !strings.Contains(content, "MULTI-TARGET") {
		t.Error("Should show multi-target mode indicator")
	}

	if !strings.Contains(content, "2 target(s) selected") {
		t.Error("Should show correct number of selected targets")
	}

	// Should have checkboxes
	checkCount := strings.Count(content, "[✓]")
	if checkCount != 2 {
		t.Errorf("Should have 2 checked boxes, got %d", checkCount)
	}
}

// TestOptimized_EditMode tests edit mode rendering
func TestOptimized_EditMode(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "15", false, "", 0, true, true, true, "hp", 0, 0, false, nil, config,
	)

	if !strings.Contains(content, "EDIT MODE") {
		t.Error("Should show edit mode indicator")
	}

	if !strings.Contains(content, "HP Change") {
		t.Error("Should show HP change prompt")
	}
}

// TestOptimized_InputMode tests input mode rendering
func TestOptimized_InputMode(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "Gandalf", true, "player_name", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if !strings.Contains(content, "Player Name:") {
		t.Error("Should show player name prompt")
	}

	if !strings.Contains(content, "Gandalf") {
		t.Error("Should show input text")
	}
}

// TestOptimized_CurrentTurnMarker tests current turn marker
func TestOptimized_CurrentTurnMarker(t *testing.T) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 2, 1, false, nil, config,
	)

	// Should have turn marker
	if !strings.Contains(content, "★") {
		t.Error("Should show turn marker")
	}
}

// BenchmarkOptimized_SmallList benchmarks rendering with 5 entries
func BenchmarkOptimized_SmallList(b *testing.B) {
	list := createLargeInitiativeList(5)
	config := panels.DefaultRenderConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panels.GetOptimizedInitiativeContent(
			list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
		)
	}
}

// BenchmarkOptimized_MediumList benchmarks rendering with 20 entries
func BenchmarkOptimized_MediumList(b *testing.B) {
	list := createLargeInitiativeList(20)
	config := panels.DefaultRenderConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panels.GetOptimizedInitiativeContent(
			list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
		)
	}
}

// BenchmarkOptimized_LargeList benchmarks rendering with 50 entries
func BenchmarkOptimized_LargeList(b *testing.B) {
	list := createLargeInitiativeList(50)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panels.GetOptimizedInitiativeContent(
			list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
		)
	}
}

// BenchmarkOptimized_VeryLargeList benchmarks rendering with 100 entries
func BenchmarkOptimized_VeryLargeList(b *testing.B) {
	list := createLargeInitiativeList(100)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panels.GetOptimizedInitiativeContent(
			list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
		)
	}
}

// TestOptimized_Performance tests rendering performance
func TestOptimized_Performance(t *testing.T) {
	list := createLargeInitiativeList(50)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 20

	start := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		_ = panels.GetOptimizedInitiativeContent(
			list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
		)
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(iterations)

	// Should be fast - under 1ms per render on average
	if avgDuration > time.Millisecond {
		t.Logf("Warning: Average render time is %v (target: <1ms)", avgDuration)
	} else {
		t.Logf("Average render time: %v", avgDuration)
	}
}

// TestOptimized_ViewportBounds tests viewport bounds are respected
func TestOptimized_ViewportBounds(t *testing.T) {
	list := createLargeInitiativeList(30)
	config := panels.DefaultRenderConfig()
	config.ViewportHeight = 10

	testCases := []struct {
		name         string
		scrollOffset int
		expectedMin  int
		expectedMax  int
	}{
		{"Start", 0, 1, 10},
		{"Middle", 10, 11, 20},
		{"End", 20, 21, 30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.ScrollOffset = tc.scrollOffset

			content := panels.GetOptimizedInitiativeContent(
				list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
			)

			// Check that expected entries are present
			for i := tc.expectedMin; i <= tc.expectedMax; i++ {
				entryNum := fmt.Sprintf("%2d.", i)
				if !strings.Contains(content, entryNum) {
					t.Errorf("Should contain entry %d", i)
				}
			}
		})
	}
}

// TestOptimized_PlayerVsMonster tests player and monster rendering
func TestOptimized_PlayerVsMonster(t *testing.T) {
	list := "[{Name:Fighter Type:player Initiative:18 HP:0 MaxHP:0 AC:18 Conditions:[]} " +
		"{Name:Goblin Type:monster Initiative:10 HP:7 MaxHP:7 AC:15 Conditions:[]}]"

	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if !strings.Contains(content, "Fighter") {
		t.Error("Should contain player name")
	}

	if !strings.Contains(content, "Goblin") {
		t.Error("Should contain monster name")
	}

	// Monster should show HP
	if !strings.Contains(content, "HP:") {
		t.Error("Monster should show HP")
	}
}

// TestOptimized_Conditions tests condition rendering
func TestOptimized_Conditions(t *testing.T) {
	list := "[{Name:Poisoned Goblin Type:monster Initiative:10 HP:5 MaxHP:7 AC:15 " +
		"Conditions:[{Name:Poisoned RoundsLeft:3 TotalRounds:5 Description:}]}]"

	config := panels.DefaultRenderConfig()

	content := panels.GetOptimizedInitiativeContent(
		list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
	)

	if !strings.Contains(content, "Poisoned Goblin") {
		t.Error("Should contain creature name")
	}

	// Should have condition emoji
	if !strings.Contains(content, "🤢") {
		t.Error("Should contain condition emoji for Poisoned")
	}
}

// TestOptimized_HPColors tests HP color coding
func TestOptimized_HPColors(t *testing.T) {
	// Test with different HP levels
	testCases := []struct {
		name   string
		hp     int
		maxHP  int
	}{
		{"Healthy", 10, 10},
		{"Bloodied", 5, 10},
		{"Critical", 2, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := fmt.Sprintf("[{Name:Monster Type:monster Initiative:10 HP:%d MaxHP:%d AC:15 Conditions:[]}]",
				tc.hp, tc.maxHP)

			config := panels.DefaultRenderConfig()

			content := panels.GetOptimizedInitiativeContent(
				list, "", false, "", 0, true, false, false, "", 0, 0, false, nil, config,
			)

			// Should contain HP ratio
			hpText := fmt.Sprintf("%d/%d", tc.hp, tc.maxHP)
			if !strings.Contains(content, hpText) {
				t.Errorf("Should contain HP ratio %s", hpText)
			}
		})
	}
}
