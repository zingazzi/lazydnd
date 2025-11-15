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
	app.setupAutoSave()
	app.setupHandlers() // Set handlers AFTER layout so SetRoot is called first

	// Initialize panel borders and titles (but don't call Draw() yet)
	app.updatePanelBorders()
	app.updatePanelContent()
	app.updateStatusBar()

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
	// Note: Input capture is handled at application level, not grid level

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

	app.application.SetRoot(app.grid, true).SetFocus(app.grid)
}

// setupHandlers configures input handlers
func (app *App) setupHandlers() {
	// Set input capture at application level to intercept all input before TView's focus system
	app.application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Debug: Log that we received an event
		ui.DebugLog("INPUT CAPTURE: Received key event - Key: %d, Rune: %c", event.Key(), event.Rune())

		// Convert TCell event to handler
		key := event.Key()
		rune := event.Rune()

		// Route to appropriate handler
		handled, shouldQuit := app.handleInput(key, rune)
		ui.DebugLog("INPUT CAPTURE: handled=%v, shouldQuit=%v", handled, shouldQuit)

		if shouldQuit {
			ui.DebugLog("INPUT CAPTURE: Quit requested, stopping application")
			// Stop in a goroutine to avoid blocking
			go func() {
				app.Stop()
			}()
			return nil
		}
		if handled {
			// Move all updates INSIDE QueueUpdateDraw callback
			// This ensures they happen on the main event loop
			ui.DebugLog("INPUT CAPTURE: Queuing UI update - ActivePanel=%d", app.model.ActivePanel)
			app.application.QueueUpdateDraw(func() {
				ui.DebugLog("QUEUE UPDATE DRAW: Callback executing - ActivePanel=%d", app.model.ActivePanel)
				app.updateStatusBar()
				app.updatePanelBorders()
				app.updatePanelContent()
				ui.DebugLog("QUEUE UPDATE DRAW: Updates complete")
			})
			ui.DebugLog("INPUT CAPTURE: QueueUpdateDraw scheduled, returning nil")

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
// This can be called from any goroutine - it uses QueueUpdateDraw for thread safety
func (app *App) refreshUI() {
	if app.application != nil {
		app.application.QueueUpdateDraw(func() {
			app.updateStatusBar()
			app.updatePanelBorders()
			app.updatePanelContent()
		})
	}
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
	ui.DebugLog("updatePanelBorders: ActivePanel=%d", app.model.ActivePanel)
	for panelType, panel := range app.panels {
		if textView, ok := panel.(*tview.TextView); ok {
			isActive := app.model.ActivePanel == panelType
			ui.DebugLog("updatePanelBorders: Panel %d, isActive=%v", panelType, isActive)
			app.stylePanel(textView, panelType, isActive)
		}
	}
}

// stylePanel applies styling to a panel
func (app *App) stylePanel(textView *tview.TextView, panelType ui.PanelType, isActive bool) {
	title := ui.PanelNames[panelType]

	textView.SetTitle(" " + title + " ")

	// Ensure border is enabled
	textView.SetBorder(true)

	var borderColor, titleColor tcell.Color
	if isActive {
		borderColor = tcell.ColorYellow
		titleColor = tcell.ColorYellow
		ui.DebugLog("stylePanel: Panel %d set to ACTIVE (Yellow) - borderColor=%d, titleColor=%d", panelType, borderColor, titleColor)
	} else {
		borderColor = tcell.ColorWhite
		titleColor = tcell.ColorWhite
		ui.DebugLog("stylePanel: Panel %d set to INACTIVE (White) - borderColor=%d, titleColor=%d", panelType, borderColor, titleColor)
	}

	textView.SetBorderColor(borderColor)
	textView.SetTitleColor(titleColor)

	ui.DebugLog("stylePanel: Colors set for panel %d - borderColor=%d, titleColor=%d", panelType, borderColor, titleColor)
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
	// Ensure UI is up to date before starting
	app.updatePanelBorders()
	app.updatePanelContent()
	app.updateStatusBar()

	// Ensure grid has focus for input capture
	app.application.SetFocus(app.grid)

	// TView will automatically draw when Run() starts
	return app.application.Run()
}

// Stop stops the TView application
func (app *App) Stop() {
	app.application.Stop()
	if app.updateTimer != nil {
		app.updateTimer.Stop()
	}
}
