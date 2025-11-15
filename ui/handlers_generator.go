// ui/handlers_generator.go
package ui

import (
	"lazydnd/encounters"
	"lazydnd/panels"
	"strconv"
	"strings"
)

// handleGeneratorPopupInput handles input in the generator popup
func handleGeneratorPopupInput(m Model, msg KeyMsg) (Model, Cmd) {
	key := msg.String()

	// Allow 'q' to quit the application (handled by global handler)
	if key == "q" {
		// Quit handled by TView app
		return m, nil
	}

	switch key {
	case "esc":
		m.EncounterGenerating = false
		return m, nil

	case "tab":
		// Switch focus between difficulty and environment
		if m.EncounterGeneratorFocus == "difficulty" {
			m.EncounterGeneratorFocus = "environment"
		} else {
			m.EncounterGeneratorFocus = "difficulty"
		}
		return m, nil

	case "up", "k":
		if m.EncounterGeneratorFocus == "difficulty" {
			if m.EncounterDifficultyIndex > 0 {
				m.EncounterDifficultyIndex--
			}
		} else if m.EncounterGeneratorFocus == "environment" {
			if m.EncounterEnvironmentIndex > 0 {
				m.EncounterEnvironmentIndex--
			}
		}
		return m, nil

	case "down", "j":
		if m.EncounterGeneratorFocus == "difficulty" {
			difficulties := 4 // easy, medium, hard, deadly
			if m.EncounterDifficultyIndex < difficulties-1 {
				m.EncounterDifficultyIndex++
			}
		} else if m.EncounterGeneratorFocus == "environment" {
			if m.EncounterEnvironmentIndex < len(m.AvailableEnvironments)-1 {
				m.EncounterEnvironmentIndex++
			}
		}
		return m, nil

	case "enter":
		// Generate encounter
		m = generateEncounter(m)
		m.EncounterGenerating = false
		return m, nil
	}

	return m, nil
}

// generateEncounter generates an encounter based on current settings
func generateEncounter(m Model) Model {
	difficulties := []string{"easy", "medium", "hard", "deadly"}
	selectedDifficulty := difficulties[m.EncounterDifficultyIndex]
	selectedEnvironment := m.AvailableEnvironments[m.EncounterEnvironmentIndex]

	// Load all monsters for encounter generation (no limit)
	err := panels.LoadMonsters()
	if err != nil {
		DebugLog("ENCOUNTER GEN: Failed to load monsters: %v", err)
		SetError(&m, "Failed to load monsters")
		return m
	}

	// Get all monsters (using dedicated function that returns ALL monsters)
	allMonsters := panels.GetAllMonstersForEncounter()
	DebugLog("ENCOUNTER GEN: Total monsters loaded: %d", len(allMonsters))

	monsterInfos := []encounters.MonsterInfo{}

	for _, name := range allMonsters {
		panelsMonster := panels.GetMonsterByName(name)
		if panelsMonster == nil {
			continue
		}

		hp := parseMonsterHP(panelsMonster.HitPoints)
		ac := parseMonsterAC(panelsMonster.ArmorClass)
		xp := encounters.GetCRXP(panelsMonster.Challenge)

		monsterInfos = append(monsterInfos, encounters.MonsterInfo{
			Name:        panelsMonster.Name,
			CR:          panelsMonster.Challenge,
			HP:          hp,
			MaxHP:       hp,
			AC:          ac,
			XP:          xp,
			Meta:        panelsMonster.Meta,
			Environment: []string{}, // We'll use meta field for matching
		})
	}

	DebugLog("ENCOUNTER GEN: MonsterInfos created: %d", len(monsterInfos))
	DebugLog("ENCOUNTER GEN: Party level: %d, Party size: %d", m.PartyLevel, m.PartySize)
	DebugLog("ENCOUNTER GEN: Difficulty: %s, Environment: %s", selectedDifficulty, selectedEnvironment)

	// Generate encounter
	result := encounters.GenerateEncounter(encounters.GenerateEncounterRequest{
		PartySize:   m.PartySize,
		PartyLevel:  m.PartyLevel,
		Difficulty:  selectedDifficulty,
		Environment: encounters.Environment(selectedEnvironment),
		Monsters:    monsterInfos,
	})

	DebugLog("ENCOUNTER GEN: Generation result - Monsters count: %d", len(result.Monsters))
	DebugLog("ENCOUNTER GEN: EnvironmentMsg: %s", result.EnvironmentMsg)

	// Convert result to UI EncounterMonsters
	m.EncounterMonsters = []EncounterMonster{}
	for _, monster := range result.Monsters {
		m.EncounterMonsters = append(m.EncounterMonsters, EncounterMonster{
			Name:     monster.Name,
			CR:       monster.CR,
			HP:       monster.HP,
			MaxHP:    monster.MaxHP,
			AC:       monster.AC,
			Quantity: monster.Quantity,
			XP:       monster.XP,
		})
	}

	// Set mode to building
	m.EncounterBuilderMode = "building"
	if len(m.EncounterMonsters) > 0 {
		m.SelectedEncounterIndex = 0
	}

	// Show success message
	if len(result.Monsters) > 0 {
		SetSuccess(&m, "Generated "+result.ActualDiff+" encounter!")
	} else {
		if result.EnvironmentMsg != "" {
			SetError(&m, result.EnvironmentMsg)
		} else {
			SetError(&m, "Could not generate encounter - try different settings")
		}
	}

	return m
}

// parseMonsterHP extracts HP from monster data
func parseMonsterHP(hpStr string) int {
	fields := strings.Fields(hpStr)
	if len(fields) == 0 {
		return 1
	}
	hp, err := strconv.Atoi(fields[0])
	if err != nil {
		return 1
	}
	return hp
}

// parseMonsterAC extracts AC from monster data
func parseMonsterAC(acStr string) int {
	fields := strings.Fields(acStr)
	if len(fields) == 0 {
		return 10
	}
	ac, err := strconv.Atoi(fields[0])
	if err != nil {
		return 10
	}
	return ac
}
