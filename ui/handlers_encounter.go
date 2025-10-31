// ui/handlers_encounter.go
package ui

import (
	"fmt"
	"lazydnd/encounters"
	"lazydnd/panels"
	"math/rand"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// handleEncounterBuilderInput handles all input for the encounter builder panel
func handleEncounterBuilderInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Don't handle if we're in a search mode from another panel
	if m.MonsterSearchMode || m.SpellSearchMode {
		return m, nil
	}

	key := msg.String()
	DebugLog("ENCOUNTER BUILDER: Mode=%s, Key=%s", m.EncounterBuilderMode, key)

	switch m.EncounterBuilderMode {
	case "party_setup":
		return handlePartySetupInput(m, msg)
	case "building":
		return handleBuildingInput(m, msg)
	case "templates":
		return handleTemplatesInput(m, msg)
	case "template_detail":
		return handleTemplateDetailInput(m, msg)
	default:
		return m, nil
	}
}

// handlePartySetupInput handles input in party setup mode
// Returns (model, cmd, handled) - if not handled, let other handlers process it
func handlePartySetupInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Set party size
		size, _ := strconv.Atoi(msg.String())
		m.PartySize = size
		return m, nil

	case "!": // Shift+1
		m.PartyLevel = 1
		return m, nil
	case "@": // Shift+2
		m.PartyLevel = 2
		return m, nil
	case "#": // Shift+3
		m.PartyLevel = 3
		return m, nil
	case "$": // Shift+4
		m.PartyLevel = 4
		return m, nil
	case "%": // Shift+5
		m.PartyLevel = 5
		return m, nil
	case "^": // Shift+6
		m.PartyLevel = 6
		return m, nil
	case "&": // Shift+7
		m.PartyLevel = 7
		return m, nil
	case "*": // Shift+8
		m.PartyLevel = 8
		return m, nil
	case "(": // Shift+9
		m.PartyLevel = 9
		return m, nil
	case "0":
		m.PartyLevel = 10
		return m, nil

	case "+", "=":
		if m.PartyLevel < 20 {
			m.PartyLevel++
		}
		return m, nil

	case "-", "_":
		if m.PartyLevel > 1 {
			m.PartyLevel--
		}
		return m, nil

	case "n", "N":
		// Next: Switch to building mode
		m.EncounterBuilderMode = "building"
		return m, nil

	case "t", "T":
		// View saved templates
		m.EncounterBuilderMode = "templates"
		m.EncounterListMode = true // Enable list mode for visual selection
		m.EncounterSelectedSaved = 0
		// Load saved encounters
		savedEncs, err := encounters.LoadEncounters("")
		if err == nil {
			m.SavedEncounters = convertToUIEncounters(savedEncs)
		}
		return m, nil

	case "g", "G":
		// Open generator popup
		m.EncounterGenerating = true
		m.EncounterDifficultyIndex = 1 // Medium
		m.EncounterEnvironmentIndex = 0 // Any
		m.EncounterGeneratorFocus = "difficulty" // Start with difficulty focused
		return m, nil

	default:
		// Let other handlers process unhandled keys (like tab for navigation)
		return m, nil
	}
}

// handleBuildingInput handles input in encounter building mode
func handleBuildingInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "m", "M":
		// Add monster from monster list - switch to monster panel
		m.ActivePanel = Monsters
		m.MonsterSearchMode = true
		m.AddingMonsterToEncounter = true // Flag that we're adding to encounter
		return m, nil

	case "up", "k":
		if len(m.EncounterMonsters) > 0 {
			if m.SelectedEncounterIndex <= 0 {
				m.SelectedEncounterIndex = len(m.EncounterMonsters) - 1
			} else {
				m.SelectedEncounterIndex--
			}
		}
		return m, nil

	case "down", "j":
		if len(m.EncounterMonsters) > 0 {
			if m.SelectedEncounterIndex >= len(m.EncounterMonsters)-1 {
				m.SelectedEncounterIndex = 0
			} else {
				m.SelectedEncounterIndex++
			}
		}
		return m, nil

	case "+", "=":
		// Increase quantity of selected monster
		if m.SelectedEncounterIndex >= 0 && m.SelectedEncounterIndex < len(m.EncounterMonsters) {
			m.EncounterMonsters[m.SelectedEncounterIndex].Quantity++
		}
		return m, nil

	case "-", "_":
		// Decrease quantity of selected monster
		if m.SelectedEncounterIndex >= 0 && m.SelectedEncounterIndex < len(m.EncounterMonsters) {
			if m.EncounterMonsters[m.SelectedEncounterIndex].Quantity > 1 {
				m.EncounterMonsters[m.SelectedEncounterIndex].Quantity--
			}
		}
		return m, nil

	case "delete", "backspace", "x", "X":
		// Remove selected monster
		if m.SelectedEncounterIndex >= 0 && m.SelectedEncounterIndex < len(m.EncounterMonsters) {
			m.EncounterMonsters = append(
				m.EncounterMonsters[:m.SelectedEncounterIndex],
				m.EncounterMonsters[m.SelectedEncounterIndex+1:]...,
			)
			if m.SelectedEncounterIndex >= len(m.EncounterMonsters) {
				m.SelectedEncounterIndex = len(m.EncounterMonsters) - 1
			}
		}
		return m, nil

	case "c", "C":
		// Clear encounter
		m.EncounterMonsters = []EncounterMonster{}
		m.SelectedEncounterIndex = -1
		m.LoadedTemplateName = "" // Clear loaded template name when clearing encounter
		return m, nil

	case "p", "P":
		// Back to party setup
		m.EncounterBuilderMode = "party_setup"
		return m, nil

	case "d", "D", "l", "L":
		// Deploy/Load to initiative tracker
		if len(m.EncounterMonsters) == 0 {
			SetError(&m, "No monsters in encounter to deploy")
			return m, nil
		}
		return deployEncounterToInitiative(m)

	case "s", "S":
		// Show save prompt
		m.ShowEncounterPrompt = true
		// Pre-fill with loaded template name if editing an existing template
		if m.LoadedTemplateName != "" {
			m.EncounterNameInput = m.LoadedTemplateName
		} else {
			m.EncounterNameInput = ""
		}
		return m, nil

	case "t", "T":
		// View saved templates
		// Load templates first
		savedEncs, err := encounters.LoadEncounters("")
		if err != nil {
			SetError(&m, fmt.Sprintf("Failed to load templates: %v", err))
			return m, nil
		}
		m.SavedEncounters = convertToUIEncounters(savedEncs)
		m.EncounterSelectedSaved = 0
		m.EncounterBuilderMode = "templates"
		m.EncounterListMode = true // Enable list mode for visual selection
		return m, nil

	case "g", "G":
		// Open generator popup (from building mode)
		m.EncounterGenerating = true
		m.EncounterDifficultyIndex = 1 // Medium
		m.EncounterEnvironmentIndex = 0 // Any
		m.EncounterGeneratorFocus = "difficulty" // Start with difficulty focused
		return m, nil

	default:
		// Let other handlers process unhandled keys (like tab for navigation)
		return m, nil
	}
}

// handleTemplatesInput handles input in templates viewing mode
func handleTemplatesInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	DebugLog("TEMPLATES LIST: Key=%s, Mode=%s", key, m.EncounterBuilderMode)

	switch key {
	case "up", "k":
		if len(m.SavedEncounters) > 0 {
			if m.EncounterSelectedSaved <= 0 {
				m.EncounterSelectedSaved = len(m.SavedEncounters) - 1
			} else {
				m.EncounterSelectedSaved--
			}
		}
		return m, nil

	case "down", "j":
		if len(m.SavedEncounters) > 0 {
			if m.EncounterSelectedSaved >= len(m.SavedEncounters)-1 {
				m.EncounterSelectedSaved = 0
			} else {
				m.EncounterSelectedSaved++
			}
		}
		return m, nil

	case "enter":
		// View selected template details
		if m.EncounterSelectedSaved >= 0 && m.EncounterSelectedSaved < len(m.SavedEncounters) {
			DebugLog("TEMPLATES: Switching to template_detail mode")
			m.EncounterBuilderMode = "template_detail"
			DebugLog("TEMPLATES: Mode is now: %s", m.EncounterBuilderMode)
		}
		return m, nil

	case "delete", "backspace", "x", "X":
		// Delete selected template
		if m.EncounterSelectedSaved >= 0 && m.EncounterSelectedSaved < len(m.SavedEncounters) {
			selected := m.SavedEncounters[m.EncounterSelectedSaved]
			err := encounters.DeleteEncounter(selected.Name, "")
			if err != nil {
				SetError(&m, fmt.Sprintf("Failed to delete: %v", err))
			} else {
				// Reload encounters
				savedEncs, err := encounters.LoadEncounters("")
				if err == nil {
					m.SavedEncounters = convertToUIEncounters(savedEncs)
					if m.EncounterSelectedSaved >= len(m.SavedEncounters) {
						m.EncounterSelectedSaved = len(m.SavedEncounters) - 1
					}
				}
				SetSuccess(&m, fmt.Sprintf("Deleted template: %s", selected.Name))
			}
		}
		return m, nil

	case "p", "P":
		// Back to party setup
		m.EncounterBuilderMode = "party_setup"
		m.EncounterListMode = false
		return m, nil

	case "n", "N":
		// New encounter
		m.EncounterBuilderMode = "building"
		m.EncounterListMode = false
		m.EncounterMonsters = []EncounterMonster{}
		m.SelectedEncounterIndex = -1
		m.LoadedTemplateName = "" // Clear loaded template name for new encounter
		return m, nil

	case "esc", "b", "B":
		// Back to building mode
		m.EncounterBuilderMode = "building"
		m.EncounterListMode = false
		return m, nil

	default:
		// Let other handlers process unhandled keys (like tab for navigation)
		return m, nil
	}
}

// handleTemplateDetailInput handles input when viewing a template's details
func handleTemplateDetailInput(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()
	DebugLog("TEMPLATE DETAIL: Key pressed: %s, Mode: %s", key, m.EncounterBuilderMode)

	switch key {
	case "l", "L", "enter":
		// Load selected template
		if m.EncounterSelectedSaved >= 0 && m.EncounterSelectedSaved < len(m.SavedEncounters) {
			selected := m.SavedEncounters[m.EncounterSelectedSaved]
			m.EncounterMonsters = selected.Monsters
			m.SelectedEncounterIndex = 0
			m.EncounterBuilderMode = "building"
			m.EncounterListMode = false
			m.LoadedTemplateName = selected.Name // Remember the loaded template name
			SetSuccess(&m, fmt.Sprintf("Loaded template: %s", selected.Name))
		}
		return m, nil

	case "d", "D", "delete", "backspace", "x", "X":
		// Delete this template
		if m.EncounterSelectedSaved >= 0 && m.EncounterSelectedSaved < len(m.SavedEncounters) {
			selected := m.SavedEncounters[m.EncounterSelectedSaved]
			err := encounters.DeleteEncounter(selected.Name, "")
			if err != nil {
				SetError(&m, fmt.Sprintf("Failed to delete: %v", err))
			} else {
				// Reload encounters
				savedEncs, err := encounters.LoadEncounters("")
				if err == nil {
					m.SavedEncounters = convertToUIEncounters(savedEncs)
					if m.EncounterSelectedSaved >= len(m.SavedEncounters) {
						m.EncounterSelectedSaved = len(m.SavedEncounters) - 1
					}
				}
				SetSuccess(&m, fmt.Sprintf("Deleted template: %s", selected.Name))
				// Go back to templates list
				m.EncounterBuilderMode = "templates"
			}
		}
		return m, nil

	case "esc", "b", "B":
		// Back to templates list
		DebugLog("TEMPLATE DETAIL: Returning to templates list")
		m.EncounterBuilderMode = "templates"
		m.EncounterListMode = true
		return m, nil

	default:
		return m, nil
	}
}

// deployEncounterToInitiative deploys all monsters from encounter to initiative tracker
func deployEncounterToInitiative(m Model) (Model, tea.Cmd) {
	addedCount := 0

	for _, monster := range m.EncounterMonsters {
		// Find monster data
		monsterData := findMonsterByName(monster.Name)
		if monsterData == nil {
			continue
		}

		// Add each instance to initiative
		for i := 0; i < monster.Quantity; i++ {
			// Roll initiative
			initiative := rollDice(1, 20) + getDexModifier(monsterData)

			// Determine instance number
			instanceNum := 0
			if monster.Quantity > 1 {
				instanceNum = i + 1
			}

			displayName := monster.Name
			if instanceNum > 0 {
				displayName = fmt.Sprintf("%s %d", monster.Name, instanceNum)
			}

			// Parse legendary actions count
			legendaryMax := 0
			if monsterData != nil && monsterData.LegendaryActions != "" {
				legendaryMax = panels.ParseLegendaryActionCount(monsterData.LegendaryActions)
			}

			entry := InitiativeEntry{
				Name:                displayName,
				Type:                "monster",
				Initiative:          initiative,
				HP:                  monster.HP,
				MaxHP:               monster.MaxHP,
				TempHP:              0,
				AC:                  monster.AC,
				ReactionUsed:        false,
				MonsterData:         monsterData,
				InstanceNum:         instanceNum,
				BaseName:            monster.Name,
				MonsterName:         monster.Name,
				Conditions:          []Condition{},
				LegendaryActionsMax: legendaryMax,
				LegendaryActionsUsed: 0, // Start with all legendary actions available
			}

			m.InitiativeList = append(m.InitiativeList, entry)
			addedCount++
		}
	}

	// Sort initiative list by initiative value (descending)
	sortInitiativeList(&m.InitiativeList)

	// Switch to initiative tracker
	m.ActivePanel = InitiativeTracker
	m.InitiativeListMode = true
	m.SelectedEntry = 0

	SetSuccess(&m, fmt.Sprintf("Deployed %d monsters to initiative", addedCount))
	return m, nil
}

// convertToUIEncounters converts encounters.Encounter to ui.Encounter
func convertToUIEncounters(encs []encounters.Encounter) []Encounter {
	result := make([]Encounter, len(encs))
	for i, enc := range encs {
		uiMonsters := make([]EncounterMonster, len(enc.Monsters))
		for j, mon := range enc.Monsters {
			uiMonsters[j] = EncounterMonster{
				Name:     mon.Name,
				CR:       mon.CR,
				HP:       mon.HP,
				MaxHP:    mon.MaxHP,
				AC:       mon.AC,
				Quantity: mon.Quantity,
				XP:       mon.XP,
			}
		}
		result[i] = Encounter{
			Name:     enc.Name,
			Monsters: uiMonsters,
		}
	}
	return result
}

// sortInitiativeList sorts the initiative list by initiative value (descending)
func sortInitiativeList(list *[]InitiativeEntry) {
	// Simple bubble sort
	for i := 0; i < len(*list); i++ {
		for j := i + 1; j < len(*list); j++ {
			if (*list)[j].Initiative > (*list)[i].Initiative {
				(*list)[i], (*list)[j] = (*list)[j], (*list)[i]
			}
		}
	}
}

// getDexModifier extracts and calculates the dexterity modifier from a monster
func getDexModifier(monster *Monster) int {
	// Parse DEXMod field
	modStr := strings.TrimSpace(monster.DEXMod)
	if modStr == "" {
		return 0
	}

	// Remove "+" prefix if present
	modStr = strings.TrimPrefix(modStr, "+")

	mod, err := strconv.Atoi(modStr)
	if err != nil {
		return 0
	}

	return mod
}

// AddMonsterToEncounter adds a selected monster to the current encounter
func AddMonsterToEncounter(m Model, monsterName string) Model {
	// Find monster data
	monster := findMonsterByName(monsterName)
	if monster == nil {
		SetError(&m, "Monster not found")
		return m
	}

	// Parse HP and AC
	hp := parseHP(monster.HitPoints)
	ac := parseAC(monster.ArmorClass)

	// Get XP from CR
	xp := encounters.GetCRXP(monster.Challenge)

	// Add to encounter
	encounterMonster := EncounterMonster{
		Name:     monster.Name,
		CR:       monster.Challenge,
		HP:       hp,
		MaxHP:    hp,
		AC:       ac,
		Quantity: 1,
		XP:       xp,
	}

	m.EncounterMonsters = append(m.EncounterMonsters, encounterMonster)
	m.SelectedEncounterIndex = len(m.EncounterMonsters) - 1

	// Switch back to encounter builder
	m.ActivePanel = EncounterBuilder
	m.EncounterBuilderMode = "building"
	m.MonsterSearchMode = false
	m.MonsterSearchInput = ""
	m.MonsterSuggestions = []string{}
	m.MonsterSuggestionIndex = -1

	SetSuccess(&m, fmt.Sprintf("Added %s to encounter", monsterName))
	return m
}

// parseHP extracts the average HP from a HP string like "45 (6d10 + 12)"
func parseHP(hpStr string) int {
	// Extract number before space or parenthesis
	parts := strings.Fields(hpStr)
	if len(parts) == 0 {
		return 1
	}

	hp, err := strconv.Atoi(parts[0])
	if err != nil {
		return 1
	}

	return hp
}

// parseAC extracts AC value from AC string like "13 (natural armor)"
func parseAC(acStr string) int {
	// Extract number before space or parenthesis
	parts := strings.Fields(acStr)
	if len(parts) == 0 {
		return 10
	}

	ac, err := strconv.Atoi(parts[0])
	if err != nil {
		return 10
	}

	return ac
}

// handleEncounterPromptInput handles input in the save encounter prompt
func handleEncounterPromptInput(m Model, key string) (Model, tea.Cmd) {
	if key == "esc" {
		// Cancel save
		m.ShowEncounterPrompt = false
		m.EncounterNameInput = ""
		return m, nil
	}

	if key == "enter" {
		// Save encounter
		if strings.TrimSpace(m.EncounterNameInput) == "" {
			SetError(&m, "Encounter name cannot be empty")
			return m, nil
		}

		// Convert to encounters.Encounter for saving
		encounterMonsters := make([]encounters.EncounterMonster, len(m.EncounterMonsters))
		for i, mon := range m.EncounterMonsters {
			encounterMonsters[i] = encounters.EncounterMonster{
				Name:     mon.Name,
				CR:       mon.CR,
				HP:       mon.HP,
				MaxHP:    mon.MaxHP,
				AC:       mon.AC,
				Quantity: mon.Quantity,
				XP:       mon.XP,
			}
		}

		enc := &encounters.Encounter{
			Name:       m.EncounterNameInput,
			PartySize:  m.PartySize,
			PartyLevel: m.PartyLevel,
			Monsters:   encounterMonsters,
		}

		err := encounters.SaveEncounter(enc, "")
		if err != nil {
			SetError(&m, fmt.Sprintf("Failed to save: %v", err))
		} else {
			SetSuccess(&m, fmt.Sprintf("Saved encounter: %s", m.EncounterNameInput))
			m.LoadedTemplateName = m.EncounterNameInput // Update loaded template name after saving
			m.ShowEncounterPrompt = false
			m.EncounterNameInput = ""
		}

		return m, nil
	}

	// Handle text input
	if key == "backspace" {
		if len(m.EncounterNameInput) > 0 {
			m.EncounterNameInput = m.EncounterNameInput[:len(m.EncounterNameInput)-1]
		}
	} else if len(key) == 1 {
		m.EncounterNameInput += key
	}

	return m, nil
}

// rollDice rolls a dice with given number of sides (e.g., rollDice(1, 20) for 1d20)
func rollDice(count int, sides int) int {
	total := 0
	for i := 0; i < count; i++ {
		total += rand.Intn(sides) + 1
	}
	return total
}

// findMonsterByName searches for a monster by name in the loaded monsters list
func findMonsterByName(name string) *Monster {
	panelsMonster := panels.GetMonsterByName(name)
	if panelsMonster == nil {
		return nil
	}

	// Convert panels.Monster to ui.Monster
	return &Monster{
		Name:             panelsMonster.Name,
		Meta:             panelsMonster.Meta,
		ArmorClass:       panelsMonster.ArmorClass,
		HitPoints:        panelsMonster.HitPoints,
		Speed:            panelsMonster.Speed,
		STR:              panelsMonster.STR,
		STRMod:           panelsMonster.STRMod,
		DEX:              panelsMonster.DEX,
		DEXMod:           panelsMonster.DEXMod,
		CON:              panelsMonster.CON,
		CONMod:           panelsMonster.CONMod,
		INT:              panelsMonster.INT,
		INTMod:           panelsMonster.INTMod,
		WIS:              panelsMonster.WIS,
		WISMod:           panelsMonster.WISMod,
		CHA:              panelsMonster.CHA,
		CHAMod:           panelsMonster.CHAMod,
		SavingThrows:     panelsMonster.SavingThrows,
		Skills:           panelsMonster.Skills,
		Senses:           panelsMonster.Senses,
		Languages:        panelsMonster.Languages,
		Challenge:        panelsMonster.Challenge,
		Traits:           panelsMonster.Traits,
		Actions:          panelsMonster.Actions,
		LegendaryActions: panelsMonster.LegendaryActions,
		ImgURL:           panelsMonster.ImgURL,
		ActionNumber:     panelsMonster.ActionNumber,
		ActionList:       convertMonsterActions(panelsMonster.ActionList),
	}
}
