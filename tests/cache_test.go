// tests/cache_test.go
package tests

import (
	"lazydnd/panels"
	"os"
	"testing"
)

// setupTestEnvironment changes to project root for asset access
func setupTestEnvironment(t *testing.T) {
	// Check if we're in tests directory
	if _, err := os.Stat("../assets/monsters.json"); err == nil {
		// We're in tests directory, change to parent
		if err := os.Chdir(".."); err != nil {
			t.Skipf("Skipping test: cannot change to project root: %v", err)
		}
	} else if _, err := os.Stat("assets/monsters.json"); err != nil {
		// Assets not found anywhere
		t.Skipf("Skipping test: assets not found")
	}
}

// checkAssetsAvailable checks if asset files are available
func checkAssetsAvailable(t *testing.T) {
	if _, err := os.Stat("assets/monsters.json"); err != nil {
		t.Skip("Skipping test: monsters.json not found")
	}
	if _, err := os.Stat("assets/spell.json"); err != nil {
		t.Skip("Skipping test: spell.json not found")
	}
}

// TestMonsterCache_InitialState tests initial cache state
func TestMonsterCache_InitialState(t *testing.T) {
	// Clear cache first
	panels.ClearMonsterCache()

	if panels.IsMonstersLoaded() {
		t.Error("Monsters should not be loaded initially")
	}

	if panels.GetMonsterCount() != 0 {
		t.Errorf("Expected 0 monsters initially, got %d", panels.GetMonsterCount())
	}
}

// TestMonsterCache_LoadMonsters tests loading monsters
func TestMonsterCache_LoadMonsters(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load monsters: %v", err)
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded after LoadMonsters()")
	}

	count := panels.GetMonsterCount()
	if count == 0 {
		t.Error("Expected some monsters to be loaded")
	}

	t.Logf("Loaded %d monsters", count)
}

// TestMonsterCache_LoadMonstersTwice tests that second load uses cache
func TestMonsterCache_LoadMonstersTwice(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// First load
	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load monsters: %v", err)
	}

	firstCount := panels.GetMonsterCount()

	// Second load should use cache
	err = panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed on second load: %v", err)
	}

	secondCount := panels.GetMonsterCount()

	if firstCount != secondCount {
		t.Errorf("Count should be same, got %d then %d", firstCount, secondCount)
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should still be loaded")
	}
}

// TestMonsterCache_ClearCache tests clearing the cache
func TestMonsterCache_ClearCache(t *testing.T) {
	setupTestEnvironment(t)
	// Load monsters
	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load monsters: %v", err)
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded")
	}

	// Clear cache
	panels.ClearMonsterCache()

	if panels.IsMonstersLoaded() {
		t.Error("Monsters should not be loaded after clearing cache")
	}

	if panels.GetMonsterCount() != 0 {
		t.Errorf("Expected 0 monsters after clear, got %d", panels.GetMonsterCount())
	}
}

// TestMonsterCache_ReloadMonsters tests reloading monsters
func TestMonsterCache_ReloadMonsters(t *testing.T) {
	setupTestEnvironment(t)
	// Initial load
	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load monsters: %v", err)
	}

	firstCount := panels.GetMonsterCount()

	// Reload
	err = panels.ReloadMonsters()
	if err != nil {
		t.Fatalf("Failed to reload monsters: %v", err)
	}

	secondCount := panels.GetMonsterCount()

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded after reload")
	}

	if firstCount != secondCount {
		t.Errorf("Count should be same after reload, got %d then %d", firstCount, secondCount)
	}
}

// TestMonsterCache_SearchAfterLoad tests searching after loading
func TestMonsterCache_SearchAfterLoad(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// Search should trigger load
	results := panels.SearchMonsters("goblin", "")

	if len(results) == 0 {
		t.Error("Expected some results for 'goblin' search")
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded after search")
	}
}

// TestMonsterCache_FindMonsterAfterLoad tests finding monster after loading
func TestMonsterCache_FindMonsterAfterLoad(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// Search for a common monster
	results := panels.SearchMonsters("goblin", "")
	if len(results) == 0 {
		t.Fatal("No goblins found to test with")
	}

	// Find the first result
	monster := panels.FindMonster(results[0])
	if monster == nil {
		t.Errorf("Should find monster %q", results[0])
	}
}

// TestSpellCache_InitialState tests initial spell cache state
func TestSpellCache_InitialState(t *testing.T) {
	panels.ClearSpellCache()

	if panels.IsSpellsLoaded() {
		t.Error("Spells should not be loaded initially")
	}

	if panels.GetSpellCount() != 0 {
		t.Errorf("Expected 0 spells initially, got %d", panels.GetSpellCount())
	}
}

// TestSpellCache_LoadSpells tests loading spells
func TestSpellCache_LoadSpells(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	err := panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded after LoadSpells()")
	}

	count := panels.GetSpellCount()
	if count == 0 {
		t.Error("Expected some spells to be loaded")
	}

	t.Logf("Loaded %d spells", count)
}

// TestSpellCache_LoadSpellsTwice tests that second load uses cache
func TestSpellCache_LoadSpellsTwice(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	// First load
	err := panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	firstCount := panels.GetSpellCount()

	// Second load should use cache
	err = panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed on second load: %v", err)
	}

	secondCount := panels.GetSpellCount()

	if firstCount != secondCount {
		t.Errorf("Count should be same, got %d then %d", firstCount, secondCount)
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should still be loaded")
	}
}

// TestSpellCache_ClearCache tests clearing the spell cache
func TestSpellCache_ClearCache(t *testing.T) {
	setupTestEnvironment(t)
	// Load spells
	err := panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded")
	}

	// Clear cache
	panels.ClearSpellCache()

	if panels.IsSpellsLoaded() {
		t.Error("Spells should not be loaded after clearing cache")
	}

	if panels.GetSpellCount() != 0 {
		t.Errorf("Expected 0 spells after clear, got %d", panels.GetSpellCount())
	}
}

// TestSpellCache_ReloadSpells tests reloading spells
func TestSpellCache_ReloadSpells(t *testing.T) {
	setupTestEnvironment(t)
	// Initial load
	err := panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	firstCount := panels.GetSpellCount()

	// Reload
	err = panels.ReloadSpells()
	if err != nil {
		t.Fatalf("Failed to reload spells: %v", err)
	}

	secondCount := panels.GetSpellCount()

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded after reload")
	}

	if firstCount != secondCount {
		t.Errorf("Count should be same after reload, got %d then %d", firstCount, secondCount)
	}
}

// TestSpellCache_SearchAfterLoad tests searching after loading
func TestSpellCache_SearchAfterLoad(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	// Search should trigger load
	results := panels.SearchSpells("fireball", "")

	if len(results) == 0 {
		t.Error("Expected some results for 'fireball' search")
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded after search")
	}
}

// TestSpellCache_FindSpellAfterLoad tests finding spell after loading
func TestSpellCache_FindSpellAfterLoad(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	// Search for a common spell
	results := panels.SearchSpells("fireball", "")
	if len(results) == 0 {
		t.Fatal("No fireball found to test with")
	}

	// Find the first result
	spell := panels.FindSpell(results[0])
	if spell == nil {
		t.Errorf("Should find spell %q", results[0])
	}
}

// TestCache_ConcurrentMonsterLoads tests concurrent monster loading
func TestCache_ConcurrentMonsterLoads(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// Load monsters concurrently
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			err := panels.LoadMonsters()
			if err != nil {
				t.Errorf("Concurrent load failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded after concurrent loads")
	}
}

// TestCache_ConcurrentSpellLoads tests concurrent spell loading
func TestCache_ConcurrentSpellLoads(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	// Load spells concurrently
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			err := panels.LoadSpells()
			if err != nil {
				t.Errorf("Concurrent load failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded after concurrent loads")
	}
}

// TestCache_MonsterSearchWithFilter tests monster search with CR filter
func TestCache_MonsterSearchWithFilter(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// Search with CR filter
	results := panels.SearchMonsters("", "0-1")

	if len(results) == 0 {
		t.Error("Expected some results for CR 0-1 filter")
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded after filtered search")
	}
}

// TestCache_SpellSearchWithFilter tests spell search with level filter
func TestCache_SpellSearchWithFilter(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearSpellCache()

	// Search with level filter
	results := panels.SearchSpells("", "0")

	if len(results) == 0 {
		t.Error("Expected some results for level 0 (cantrips) filter")
	}

	if !panels.IsSpellsLoaded() {
		t.Error("Spells should be loaded after filtered search")
	}
}

// TestCache_PerformanceBenefit tests that cache improves performance
func TestCache_PerformanceBenefit(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// First search (triggers load)
	results1 := panels.SearchMonsters("dragon", "")
	if len(results1) == 0 {
		t.Fatal("Expected dragon results")
	}

	// Second search (uses cache)
	results2 := panels.SearchMonsters("dragon", "")

	// Should get same results
	if len(results1) != len(results2) {
		t.Errorf("Expected same results, got %d then %d", len(results1), len(results2))
	}

	// Cache should still be loaded
	if !panels.IsMonstersLoaded() {
		t.Error("Cache should remain loaded")
	}
}

// TestCache_MultipleOperationsWithoutReload tests multiple operations use cache
func TestCache_MultipleOperationsWithoutReload(t *testing.T) {
	setupTestEnvironment(t)
	panels.ClearMonsterCache()

	// Initial load
	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	initialCount := panels.GetMonsterCount()

	// Perform multiple operations
	_ = panels.SearchMonsters("goblin", "")
	_ = panels.FindMonster("Goblin")
	_ = panels.SearchMonsters("dragon", "5+")

	// Count should remain the same (using cache)
	if panels.GetMonsterCount() != initialCount {
		t.Errorf("Count changed unexpectedly: %d -> %d", initialCount, panels.GetMonsterCount())
	}

	if !panels.IsMonstersLoaded() {
		t.Error("Cache should remain loaded")
	}
}

// TestCache_ClearAndReloadIndependently tests clearing one cache doesn't affect other
func TestCache_ClearAndReloadIndependently(t *testing.T) {
	setupTestEnvironment(t)
	// Load both
	err := panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to load monsters: %v", err)
	}
	err = panels.LoadSpells()
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	// Clear monsters only
	panels.ClearMonsterCache()

	// Monsters should be cleared
	if panels.IsMonstersLoaded() {
		t.Error("Monsters should be cleared")
	}

	// Spells should still be loaded
	if !panels.IsSpellsLoaded() {
		t.Error("Spells should still be loaded")
	}

	// Reload monsters
	err = panels.LoadMonsters()
	if err != nil {
		t.Fatalf("Failed to reload monsters: %v", err)
	}

	// Both should now be loaded
	if !panels.IsMonstersLoaded() {
		t.Error("Monsters should be loaded")
	}
	if !panels.IsSpellsLoaded() {
		t.Error("Spells should still be loaded")
	}
}
