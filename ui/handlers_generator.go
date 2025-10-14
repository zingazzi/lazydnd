// ui/handlers_generator.go
package ui

import (
	"lazydnd/encounters"
	"lazydnd/panels"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// handleGeneratorPopupInput handles input in the generator popup
func handleGeneratorPopupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.EncounterGenerating = false
		return m, nil

	case "up", "k":
		if m.EncounterDifficultyIndex > 0 {
			m.EncounterDifficultyIndex--
		}
		return m, nil

	case "down", "j":
		difficulties := 4 // easy, medium, hard, deadly
		if m.EncounterDifficultyIndex < difficulties-1 {
			m.EncounterDifficultyIndex++
		}
		return m, nil

	case "left", "h":
		if m.EncounterEnvironmentIndex > 0 {
			m.EncounterEnvironmentIndex--
		}
		return m, nil

	case "right", "l":
		if m.EncounterEnvironmentIndex < len(m.AvailableEnvironments)-1 {
			m.EncounterEnvironmentIndex++
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

	// Load all monsters
	err := panels.LoadMonsters()
	if err != nil {
		SetError(&m, "Failed to load monsters")
		return m
	}

	// Convert monsters to MonsterInfo
	allMonsters := panels.SearchMonsters("", "") // Get all monsters
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

	// Generate encounter
	result := encounters.GenerateEncounter(encounters.GenerateEncounterRequest{
		PartySize:   m.PartySize,
		PartyLevel:  m.PartyLevel,
		Difficulty:  selectedDifficulty,
		Environment: encounters.Environment(selectedEnvironment),
		Monsters:    monsterInfos,
	})

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
		SetError(&m, "No suitable monsters found for "+selectedEnvironment)
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

