// ui/tview/app.go
package tview

import (
	"lazydnd/ui"
	tviewpanels "lazydnd/ui/tview/panels"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


// App wraps the TView application and manages the UI
type App struct {
	application *tview.Application
	grid        *tview.Grid
	model       *ui.Model
	panels      map[ui.PanelType]tview.Primitive
	statusBar   *tview.TextView
	updateTimer *time.Timer
}

// NewApp creates a new TView application instance
func NewApp(model *ui.Model) *App {
	app := &App{
		application: tview.NewApplication(),
		model:       model,
		panels:      make(map[ui.PanelType]tview.Primitive),
	}

	app.setupStatusBar()
	app.setupPanels()
	app.setupLayout()
	app.setupHandlers()
	app.setupAutoSave()

	return app
}

// setupStatusBar creates the status bar widget
func (app *App) setupStatusBar() {
	app.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	app.updateStatusBar()
}

// setupPanels creates all panel widgets
func (app *App) setupPanels() {
	app.panels[ui.DiceRoller] = tviewpanels.NewDiceRollerPanel(app.model)
	app.panels[ui.InitiativeTracker] = tviewpanels.NewInitiativeTrackerPanel(app.model)
	app.panels[ui.Spells] = tviewpanels.NewSpellsPanel(app.model)
	app.panels[ui.Monsters] = tviewpanels.NewMonstersPanel(app.model)
	app.panels[ui.Notes] = tviewpanels.NewNotesPanel(app.model)
	app.panels[ui.EncounterBuilder] = tviewpanels.NewEncounterBuilderPanel(app.model)
}

// setupLayout configures the grid layout
func (app *App) setupLayout() {
	app.grid = tview.NewGrid().
		SetRows(0, 0, 1). // Two panel rows + status bar
		SetColumns(0, 0, 0, 0, 0). // Flexible columns
		SetBorders(false)

	// Top row: Dice Roller (60%) + Initiative Tracker (40%)
	app.grid.AddItem(app.panels[ui.DiceRoller], 0, 0, 1, 3, 0, 0, false).
		AddItem(app.panels[ui.InitiativeTracker], 0, 3, 1, 2, 0, 0, false).
		// Bottom row: Spells + Monsters + Notes + Encounter Builder
		AddItem(app.panels[ui.Spells], 1, 0, 1, 1, 0, 0, false).
		AddItem(app.panels[ui.Monsters], 1, 1, 1, 1, 0, 0, false).
		AddItem(app.panels[ui.Notes], 1, 2, 1, 1, 0, 0, false).
		AddItem(app.panels[ui.EncounterBuilder], 1, 3, 1, 2, 0, 0, false).
		// Status bar
		AddItem(app.statusBar, 2, 0, 1, 5, 0, 0, false)

	app.application.SetRoot(app.grid, true)
}

// setupHandlers configures input handlers
func (app *App) setupHandlers() {
	app.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Convert TCell event to handler
		key := event.Key()
		rune := event.Rune()

		// Route to appropriate handler
		handled, shouldQuit := app.handleInput(key, rune)
		if shouldQuit {
			app.Stop()
			return nil
		}
		if handled {
			app.refreshUI()
			return nil // Event handled
		}

		return event // Pass through if not handled
	})
}

// handleInput routes input to the appropriate handler
func (app *App) handleInput(key tcell.Key, rune rune) (bool, bool) {
	// Convert to handler chain format
	// This will be implemented in handlers.go
	return HandleInput(app.model, key, rune)
}

// setupAutoSave sets up the autosave timer
func (app *App) setupAutoSave() {
	app.updateTimer = time.NewTimer(time.Minute)
	go func() {
		for range app.updateTimer.C {
			app.handleAutoSave()
			app.updateTimer.Reset(time.Minute)
		}
	}()
}

// handleAutoSave performs autosave
func (app *App) handleAutoSave() {
	// Call the existing autosave handler
	*app.model = ui.HandleAutoSave(*app.model)
	app.refreshUI()
}

// refreshUI updates all UI elements
func (app *App) refreshUI() {
	app.updateStatusBar()
	app.updatePanelBorders()
	app.updatePanelContent()
	app.application.Draw()
}

// updateStatusBar updates the status bar content
func (app *App) updateStatusBar() {
	// Get status bar text (contains ANSI codes from Lipgloss)
	statusText := app.model.RenderStatusBar()
	// TView TextView supports ANSI colors via SetDynamicColors(true)
	app.statusBar.SetText(statusText)
}

// updatePanelBorders updates panel border styling based on active panel
func (app *App) updatePanelBorders() {
	for panelType, panel := range app.panels {
		if textView, ok := panel.(*tview.TextView); ok {
			isActive := app.model.ActivePanel == panelType
			app.stylePanel(textView, panelType, isActive)
		}
	}
}

// stylePanel applies styling to a panel
func (app *App) stylePanel(textView *tview.TextView, panelType ui.PanelType, isActive bool) {
	title := ui.PanelNames[panelType]

	textView.SetTitle(" " + title + " ")

	if isActive {
		textView.SetBorderColor(tcell.ColorYellow)
		textView.SetTitleColor(tcell.ColorYellow)
	} else {
		textView.SetBorderColor(tcell.ColorWhite)
		textView.SetTitleColor(tcell.ColorWhite)
	}
}

// updatePanelContent updates the content of all panels
func (app *App) updatePanelContent() {
	for panelType, panel := range app.panels {
		if textView, ok := panel.(*tview.TextView); ok {
			content := app.model.GetPanelContent(panelType)
			textView.SetText(content)
		}
	}
}

// Run starts the TView application
func (app *App) Run() error {
	return app.application.Run()
}

// Stop stops the TView application
func (app *App) Stop() {
	app.application.Stop()
	if app.updateTimer != nil {
		app.updateTimer.Stop()
	}
}
