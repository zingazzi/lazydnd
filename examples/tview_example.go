package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// PanelType represents the different panel types
type PanelType int

const (
	DiceRoller PanelType = iota
	InitiativeTracker
	Spells
	Monsters
	Notes
	EncounterBuilder
)

var panelNames = []string{
	"Dice Roller",
	"Initiative Tracker",
	"Spells",
	"Monsters",
	"Notes",
	"Encounter Builder",
}

func main() {
	app := tview.NewApplication()

	// Create a grid layout
	grid := tview.NewGrid().
		SetRows(0, 0). // Two rows, equal height
		SetColumns(0, 0, 0, 0, 0). // Flexible columns for bottom row
		SetBorders(false)

	// Create panels
	panels := make([]*tview.TextView, 6)

	// Panel 1: Dice Roller
	panels[DiceRoller] = tview.NewTextView()
	panels[DiceRoller].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 1. Dice Roller ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[DiceRoller].SetText("🎲 Roll dice here\n\nExample: 1d20+5\n2d6+3")

	// Panel 2: Initiative Tracker
	panels[InitiativeTracker] = tview.NewTextView()
	panels[InitiativeTracker].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 2. Initiative Tracker ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[InitiativeTracker].SetText("Initiative order:\n\n1. Player 1 (20)\n2. Monster (15)\n3. Player 2 (12)")

	// Panel 3: Spells
	panels[Spells] = tview.NewTextView()
	panels[Spells].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 3. Spells ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[Spells].SetText("Spell list:\n\n• Fireball\n• Magic Missile\n• Shield")

	// Panel 4: Monsters
	panels[Monsters] = tview.NewTextView()
	panels[Monsters].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 4. Monsters ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[Monsters].SetText("Monster list:\n\n• Goblin\n• Orc\n• Dragon")

	// Panel 5: Notes
	panels[Notes] = tview.NewTextView()
	panels[Notes].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 5. Notes ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[Notes].SetText("Quick notes:\n\n• Important clue\n• NPC name\n• Location")

	// Panel 6: Encounter Builder
	panels[EncounterBuilder] = tview.NewTextView()
	panels[EncounterBuilder].SetDynamicColors(true).
		SetWrap(true).
		SetTitle(" 6. Encounter Builder ").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	panels[EncounterBuilder].SetText("Encounter setup:\n\n• 3 Goblins\n• 1 Orc\n• Difficulty: Medium")

	// Arrange panels in grid
	// Top row: Dice Roller (60%) + Initiative Tracker (40%)
	grid.AddItem(panels[DiceRoller], 0, 0, 1, 3, 0, 0, false).
		AddItem(panels[InitiativeTracker], 0, 3, 1, 2, 0, 0, false).
		// Bottom row: Spells (30%) + Monsters (30%) + Notes (20%) + Encounter (20%)
		AddItem(panels[Spells], 1, 0, 1, 1, 0, 0, false).
		AddItem(panels[Monsters], 1, 1, 1, 1, 0, 0, false).
		AddItem(panels[Notes], 1, 2, 1, 1, 0, 0, false).
		AddItem(panels[EncounterBuilder], 1, 3, 1, 2, 0, 0, false)

	// Track active panel
	activePanel := DiceRoller

	// Function to update active panel border color
	updateActivePanel := func(panel PanelType) {
		for i, p := range panels {
			if i == int(panel) {
				p.SetBorderColor(tview.Styles.PrimaryTextColor) // Active: bright color
			} else {
				p.SetBorderColor(tview.Styles.PrimitiveBackgroundColor) // Inactive: dim color
			}
		}
		app.Draw()
	}

	// Set initial active panel
	updateActivePanel(activePanel)

	// Handle keyboard input for navigation
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Cycle through panels
			activePanel = (activePanel + 1) % 6
			updateActivePanel(activePanel)
			return nil
		case tcell.KeyBacktab:
			// Cycle backwards
			activePanel = (activePanel + 5) % 6
			updateActivePanel(activePanel)
			return nil
		case tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	// Create status bar
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[::b]Tab[:] Next Panel | [::b]Shift+Tab[:] Previous | [::b]Esc[:] Quit")

	// Create main flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(grid, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	// Run the application
	if err := app.SetRoot(flex, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
