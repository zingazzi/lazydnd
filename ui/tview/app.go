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
	application        *tview.Application
	grid               *tview.Grid
	pages              *tview.Pages // Pages widget for switching between main view and modals
	model              *ui.Model
	panels             map[ui.PanelType]tview.Primitive
	statusBar          *tview.TextView
	updateTimer        *time.Timer
	updateChan         chan bool // Channel to signal UI updates
	currentModal       tview.Primitive // Currently displayed modal
	previousActivePanel ui.PanelType   // Track previous active panel to avoid unnecessary layout recalculations
}

// NewApp creates a new TView application instance
func NewApp(model *ui.Model) *App {
	app := &App{
		application: tview.NewApplication(),
		model:       model,
		panels:      make(map[ui.PanelType]tview.Primitive),
		updateChan:  make(chan bool, 10), // Buffered channel for update signals
	}

	app.setupStatusBar()
	app.setupPanels()
	app.setupLayout()
	app.setupAutoSave()
	app.setupHandlers() // Set handlers AFTER layout so SetRoot is called first
	app.startUpdateLoop() // Start the update loop

	// Initialize previous active panel to an invalid value to force initial layout
	// This ensures the grid layout is calculated at least once
	app.previousActivePanel = ui.PanelType(-1)

	// Initialize panel borders and titles (but don't call Draw() yet)
	app.updatePanelBorders()
	app.updatePanelContent()
	app.updateStatusBar()
	app.updateModals() // Initialize modals

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
		SetColumns(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0). // 20 flexible columns for fine-grained control
		SetBorders(false)
	// Note: Input capture is handled at application level, not grid level

	// Initial layout will be set by updateGridLayout()
	// This ensures dynamic sizing based on active panel
	app.updateGridLayout()

	// Create pages widget to manage main view and modals
	app.pages = tview.NewPages()
	app.pages.AddPage("main", app.grid, true, true)

	app.application.SetRoot(app.pages, true).SetFocus(app.grid)
}

// setupHandlers configures input handlers
func (app *App) setupHandlers() {
	// Set input capture on the application level so it works for both grid and modals
	app.application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check if a modal is active
		if app.currentModal != nil {
			// Special handling for TextView-based popups (action, load, condition) - need handler chain access
			if app.model.ShowActionPopup || app.model.ShowLoadPopup || app.model.ShowConditionPopup {
				// For TextView-based popups, only intercept Esc to close
				// Let handler chain process arrow keys, Enter, etc.
				key := event.Key()
				if key == tcell.KeyEscape {
					if app.model.ShowActionPopup {
						app.model.ShowActionPopup = false
					}
					if app.model.ShowLoadPopup {
						app.model.ShowLoadPopup = false
					}
					if app.model.ShowConditionPopup {
						app.model.ShowConditionPopup = false
					}
					app.updateModals()
					return nil
				}
				// Pass through to handler chain for navigation/selection
			} else {
				// For other modals (with buttons), only intercept Esc to close
				// Let the modal handle Tab (button navigation), Enter (button selection), etc.
				key := event.Key()
				if key == tcell.KeyEscape {
					// Close current modal
					app.closeCurrentModal()
					return nil
				}
				// Pass through all other input to the modal so it can handle button navigation
				// Don't process through handler chain when modal is active
				return event
			}
		}

		// No modal active, or action popup (which needs handler chain) - process input normally
		key := event.Key()
		rune := event.Rune()

		// Allow scrolling keys to pass through to TextViews for scrolling
		// For action/load popups: arrow keys are used for navigation (via handler chain), but PageUp/Down can scroll
		// For panels: all scroll keys should pass through to the active panel's TextView
		scrollKeys := []tcell.Key{
			tcell.KeyPgUp, tcell.KeyPgDn,
			tcell.KeyHome, tcell.KeyEnd,
			tcell.KeyCtrlF, tcell.KeyCtrlB,
		}

		// Check if we should allow Up/Down for scrolling
		// Don't allow Up/Down for scrolling if:
		// - Action/Load/Condition popups are active (arrow keys navigate lists)
		// - Monster/Spell search mode is active (arrow keys navigate suggestions)
		// - Initiative input/edit mode (arrow keys have other functions)
		// - Initiative list mode is active (arrow keys navigate list)
		allowUpDownForScrolling := !app.model.ShowActionPopup &&
			!app.model.ShowLoadPopup &&
			!app.model.ShowConditionPopup &&
			!app.model.MonsterSearchMode &&
			!app.model.SpellSearchMode &&
			!app.model.MonsterCRFilterMode &&
			!app.model.InitiativeInputMode &&
			!app.model.InitiativeEditMode &&
			!app.model.InitiativeListMode

		if allowUpDownForScrolling {
			scrollKeys = append(scrollKeys, tcell.KeyUp, tcell.KeyDown)
		}

		for _, scrollKey := range scrollKeys {
			if key == scrollKey {
				// For panels, set focus on the active panel's TextView and pass event through
				if !app.model.ShowActionPopup && !app.model.ShowLoadPopup && !app.model.ShowConditionPopup {
					if panel, ok := app.panels[app.model.ActivePanel]; ok {
						if textView, ok := panel.(*tview.TextView); ok {
							// Set focus on the TextView so it can handle scrolling
							app.application.SetFocus(textView)
							// Pass the event through to the TextView
							return event
						}
					}
				} else {
					// For action/load/condition popups, pass through PageUp/Down for scrolling
					// Arrow keys will be handled by handler chain for navigation
					return event
				}
			}
		}

		// Route to appropriate handler
		handled, shouldQuit := app.handleInput(key, rune)

		if shouldQuit {
			// Stop in a goroutine to avoid blocking
			go func() {
				app.Stop()
			}()
			return nil
		}
		if handled {
			// Update UI immediately in the input handler
			// Then signal for a redraw via channel (non-blocking)
			app.updateStatusBar()
			app.updatePanelBorders()
			app.updatePanelContent()
			app.updateModals() // Update modals based on popup state

			// Signal for redraw (non-blocking, don't wait)
			select {
			case app.updateChan <- true:
			default:
				// Channel full, skip - UI already updated
			}

			return nil // Event handled
		}

		return event // Pass through if not handled
	})
}

// closeCurrentModal closes the currently active modal
func (app *App) closeCurrentModal() {
	if app.model.ShowLoadPopup {
		app.model.ShowLoadPopup = false
	} else if app.model.ShowActionPopup {
		app.model.ShowActionPopup = false
	} else if app.model.ShowSavePopup {
		app.model.ShowSavePopup = false
	} else if app.model.ShowHelpPopup {
		app.model.ShowHelpPopup = false
	} else if app.model.ShowQuickHPPopup {
		app.model.ShowQuickHPPopup = false
	} else if app.model.ShowMultiTargetPopup {
		app.model.ShowMultiTargetPopup = false
	} else if app.model.ShowConditionPopup {
		app.model.ShowConditionPopup = false
	} else if app.model.ShowCastSpellPrompt {
		app.model.ShowCastSpellPrompt = false
	} else if app.model.ShowSavingThrowPopup {
		app.model.ShowSavingThrowPopup = false
	} else if app.model.ShowEncounterPrompt {
		app.model.ShowEncounterPrompt = false
	} else if app.model.EncounterGenerating {
		app.model.EncounterGenerating = false
	} else if app.model.ShowRenamePopup {
		app.model.ShowRenamePopup = false
	}
	app.updateModals()
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
	// Update grid layout only when active panel changes
	// This prevents unnecessary recalculations that cause visual glitches
	if app.model.ActivePanel != app.previousActivePanel {
		app.updateGridLayout()
		app.previousActivePanel = app.model.ActivePanel
	}

	for panelType, panel := range app.panels {
		if textView, ok := panel.(*tview.TextView); ok {
			isActive := app.model.ActivePanel == panelType
			app.stylePanel(textView, panelType, isActive)
		}
	}
}

// calculatePanelDimensionsForGrid calculates panel dimensions for grid layout
// This is a simplified version that matches the logic in layout_core.go
func (app *App) calculatePanelDimensionsForGrid() map[ui.PanelType]ui.PanelDimensions {
	// Use the same calculation logic as Model.calculatePanelDimensions
	// but we need to access it through a helper or duplicate the logic
	// For now, we'll use a simplified version based on active panel

	targetRowWidth := app.model.Width - 2 // Approximate margin
	if targetRowWidth < 40 {
		targetRowWidth = 40
	}

	dimensions := make(map[ui.PanelType]ui.PanelDimensions)
	topHeight := 10 // Approximate
	bottomHeight := 10 // Approximate

	// Width allocation based on active panel (matching layout_core.go logic)
	switch app.model.ActivePanel {
	case ui.DiceRoller:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 6 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	case ui.InitiativeTracker:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 4 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 3 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	case ui.Spells:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	case ui.Monsters:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	case ui.Notes:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 4 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	case ui.EncounterBuilder:
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth * 2 / 10, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	default:
		// Default equal distribution
		dimensions[ui.DiceRoller] = ui.PanelDimensions{Width: targetRowWidth * 5 / 10, Height: topHeight}
		dimensions[ui.InitiativeTracker] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.DiceRoller].Width, Height: topHeight}
		dimensions[ui.Spells] = ui.PanelDimensions{Width: targetRowWidth / 4, Height: bottomHeight}
		dimensions[ui.Monsters] = ui.PanelDimensions{Width: targetRowWidth / 4, Height: bottomHeight}
		dimensions[ui.Notes] = ui.PanelDimensions{Width: targetRowWidth / 4, Height: bottomHeight}
		dimensions[ui.EncounterBuilder] = ui.PanelDimensions{Width: targetRowWidth - dimensions[ui.Spells].Width - dimensions[ui.Monsters].Width - dimensions[ui.Notes].Width, Height: bottomHeight}
	}

	return dimensions
}

// updateGridLayout updates the grid layout based on active panel dimensions
func (app *App) updateGridLayout() {
	// Get dynamic panel dimensions
	dimensions := app.calculatePanelDimensionsForGrid()

	// Grid uses 20 columns for fine-grained control
	// Convert panel widths to column spans proportionally
	totalCols := 20

	// Calculate column spans for top row
	diceWidth := dimensions[ui.DiceRoller].Width
	initWidth := dimensions[ui.InitiativeTracker].Width
	topRowTotal := diceWidth + initWidth

	var diceCols, initCols int
	if topRowTotal > 0 {
		diceCols = (diceWidth * totalCols) / topRowTotal
		initCols = totalCols - diceCols
	} else {
		diceCols = totalCols / 2
		initCols = totalCols / 2
	}

	// Calculate column spans for bottom row
	spellsWidth := dimensions[ui.Spells].Width
	monstersWidth := dimensions[ui.Monsters].Width
	notesWidth := dimensions[ui.Notes].Width
	encounterWidth := dimensions[ui.EncounterBuilder].Width
	bottomRowTotal := spellsWidth + monstersWidth + notesWidth + encounterWidth

	var spellsCols, monstersCols, notesCols, encounterCols int
	if bottomRowTotal > 0 {
		spellsCols = (spellsWidth * totalCols) / bottomRowTotal
		monstersCols = (monstersWidth * totalCols) / bottomRowTotal
		notesCols = (notesWidth * totalCols) / bottomRowTotal
		encounterCols = totalCols - spellsCols - monstersCols - notesCols
	} else {
		// Default equal distribution
		spellsCols = totalCols / 4
		monstersCols = totalCols / 4
		notesCols = totalCols / 4
		encounterCols = totalCols - spellsCols - monstersCols - notesCols
	}

	// Remove all items from grid
	app.grid.Clear()

	// Re-add items with new column spans
	// Top row: Dice Roller + Initiative Tracker
	app.grid.AddItem(app.panels[ui.DiceRoller], 0, 0, 1, diceCols, 0, 0, false).
		AddItem(app.panels[ui.InitiativeTracker], 0, diceCols, 1, initCols, 0, 0, false).
		// Bottom row: Spells + Monsters + Notes + Encounter Builder
		AddItem(app.panels[ui.Spells], 1, 0, 1, spellsCols, 0, 0, false).
		AddItem(app.panels[ui.Monsters], 1, spellsCols, 1, monstersCols, 0, 0, false).
		AddItem(app.panels[ui.Notes], 1, spellsCols+monstersCols, 1, notesCols, 0, 0, false).
		AddItem(app.panels[ui.EncounterBuilder], 1, spellsCols+monstersCols+notesCols, 1, encounterCols, 0, 0, false).
		// Status bar spans all columns
		AddItem(app.statusBar, 2, 0, 1, totalCols, 0, 0, false)
}

// stylePanel applies styling to a panel
func (app *App) stylePanel(textView *tview.TextView, panelType ui.PanelType, isActive bool) {
	title := ui.PanelNames[panelType] // Title already includes icon

	textView.SetTitle(" " + title + " ")

	// Ensure border is enabled
	textView.SetBorder(true)

	// Get colors from config
	colorConverter := NewColorConverter(app.model.Config)

	// Set border colors based on active state
	var borderColor, titleColor tcell.Color
	if isActive {
		// Active panel: use primary color (violet by default)
		borderColor = colorConverter.PrimaryColor()
		titleColor = colorConverter.PrimaryColor()
	} else {
		// Inactive panel: use border color (light grey by default)
		borderColor = colorConverter.BorderColor()
		titleColor = colorConverter.BorderColor()
	}

	// Set panel background to black
	// Note: Text color is set via color tags in content, not via SetTextColor
	textView.SetBackgroundColor(tcell.ColorBlack)

	textView.SetBorderColor(borderColor)
	textView.SetTitleColor(titleColor)
	// Note: TView doesn't support SetTitleBackgroundColor directly
	// Title background is handled by the border color system
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

// startUpdateLoop starts a goroutine that processes UI redraw signals
func (app *App) startUpdateLoop() {
	go func() {
		for {
			select {
			case <-app.updateChan:
				// Call Draw() directly from goroutine - this is safe and faster than QueueUpdateDraw
				// UI is already updated in the input handler, we just need to redraw
				app.application.Draw()
			}
		}
	}()
}

// Stop stops the TView application
func (app *App) Stop() {
	app.application.Stop()
	if app.updateTimer != nil {
		app.updateTimer.Stop()
	}
	if app.updateChan != nil {
		close(app.updateChan)
	}
}

// updateModals checks popup state and shows/hides modals accordingly
func (app *App) updateModals() {
	// Check which popup should be shown
	var modal tview.Primitive = nil
	var modalName string

	if app.model.ShowLoadPopup {
		modal = app.createLoadTextView()
		modalName = "load"
	} else if app.model.ShowActionPopup {
		modal = app.createActionTextView()
		modalName = "action"
	} else if app.model.ShowSavePopup {
		modal = app.createSaveModal()
		modalName = "save"
	} else if app.model.ShowHelpPopup {
		modal = app.createHelpTextView()
		modalName = "help"
	} else if app.model.ShowQuickHPPopup {
		modal = app.createQuickHPModal()
		modalName = "quickhp"
	} else if app.model.ShowMultiTargetPopup {
		modal = app.createMultiTargetModal()
		modalName = "multitarget"
	} else if app.model.ShowConditionPopup {
		modal = app.createConditionTextView()
		modalName = "condition"
	} else if app.model.ShowCastSpellPrompt {
		modal = app.createCastSpellModal()
		modalName = "castspell"
	} else if app.model.ShowSavingThrowPopup {
		modal = app.createSavingThrowModal()
		modalName = "savingthrow"
	} else if app.model.ShowEncounterPrompt {
		modal = app.createEncounterPromptModal()
		modalName = "encounterprompt"
	} else if app.model.EncounterGenerating {
		modal = app.createEncounterGeneratorModal()
		modalName = "encountergen"
	} else if app.model.ShowRenamePopup {
		modal = app.createRenameModal()
		modalName = "rename"
	}

	// Show or hide modal using Pages widget
	if modal != nil {
		if app.currentModal != modal {
			// New modal to show
			app.currentModal = modal
			// For action, load, and condition popups, use fullScreen=false to allow custom sizing via Flex
			// For other modals, use fullScreen=true (default modal behavior)
			fullScreen := modalName != "action" && modalName != "load" && modalName != "condition"
			app.pages.AddPage(modalName, modal, true, fullScreen)
			app.pages.SwitchToPage(modalName)
			app.application.SetFocus(modal)
		}
	} else {
		if app.currentModal != nil {
			// Hide modal, return to main view
			app.currentModal = nil
			app.pages.SwitchToPage("main")
			app.application.SetFocus(app.grid)
		}
	}
}

// createLoadTextView creates a TextView widget for the load popup (allows keyboard navigation)
func (app *App) createLoadTextView() tview.Primitive {
	content := ui.RenderLoadPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	colorConverter := NewColorConverter(app.model.Config)

	loadView := tview.NewTextView()
	loadView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true). // Enable scrolling for long lists
		SetTextAlign(tview.AlignLeft)
	loadView.SetBorder(true).
		SetTitle(" Load Campaign ").
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderColor(colorConverter.PrimaryColor()).
		SetTitleColor(colorConverter.PrimaryColor())
	loadView.SetText(coloredContent)

	// Don't set SetInputCapture on the TextView - let application-level handler process all input
	// This allows the handler chain to process arrow keys, Enter, Esc, etc.

	// Wrap in a Flex container to center it and size it appropriately
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer (flexible)
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false). // Left spacer (flexible)
				AddItem(loadView, 0, 1, false). // Load view (flexible, will size to content)
				AddItem(nil, 0, 1, false), // Right spacer (flexible)
			0, 1, false). // Center row (flexible height)
		AddItem(nil, 0, 1, false) // Bottom spacer (flexible)

	return flex
}

// createActionTextView creates a TextView widget for the action popup (allows keyboard navigation)
// Wrapped in a Grid to center it and make it a reasonable size (not 100%)
func (app *App) createActionTextView() tview.Primitive {
	content := ui.RenderActionPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	colorConverter := NewColorConverter(app.model.Config)

	actionView := tview.NewTextView()
	actionView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true). // Enable scrolling for long action lists
		SetTextAlign(tview.AlignLeft)
	actionView.SetBorder(true).
		SetTitle(" Monster Actions ").
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderColor(colorConverter.PrimaryColor()).
		SetTitleColor(colorConverter.PrimaryColor())

	// Set the action content with grey color
	actionView.SetText(coloredContent)

	// Don't set SetInputCapture on the TextView - let application-level handler process all input
	// This allows the handler chain to process arrow keys, Enter, etc.

	// Wrap in a Flex container to center it and size it appropriately
	// This creates a centered modal that's about 60% width and 70% height
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer (flexible)
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false). // Left spacer (flexible)
				AddItem(actionView, 0, 1, false). // Action view (flexible, will size to content)
				AddItem(nil, 0, 1, false), // Right spacer (flexible)
			0, 1, false). // Center row (flexible height)
		AddItem(nil, 0, 1, false) // Bottom spacer (flexible)

	return flex
}

// createSaveModal creates a modal for the save popup
func (app *App) createSaveModal() *tview.Modal {
	content := ui.RenderSavePopup(*app.model)

	modal := tview.NewModal().
		SetText(content).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				// Handle save action - this will be handled by input handler
			}
			app.model.ShowSavePopup = false
			app.updateModals()
		})

	return modal
}

// createHelpTextView creates a TextView widget for the help popup (two columns)
func (app *App) createHelpTextView() *tview.TextView {
	// Get help content (two columns: left=generic, right=panel-specific)
	content := ui.RenderHelpPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	colorConverter := NewColorConverter(app.model.Config)

	helpView := tview.NewTextView()
	helpView.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	helpView.SetBorder(true).
		SetTitle(" Help ").
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderColor(colorConverter.PrimaryColor()).
		SetTitleColor(colorConverter.PrimaryColor())

	// Set the help content with grey color
	helpView.SetText(coloredContent)

	// Handle input for closing
	helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyEscape || (key == tcell.KeyRune && event.Rune() == '?') {
			app.model.ShowHelpPopup = false
			app.updateModals()
			return nil
		}
		return event
	})

	return helpView
}

// createQuickHPModal creates a modal for the quick HP popup
func (app *App) createQuickHPModal() *tview.Modal {
	content := ui.RenderQuickHPPopup(*app.model)

	modal := tview.NewModal().
		SetText(content).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Apply", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Apply" {
				// Handle HP change - this will be handled by input handler
			}
			app.model.ShowQuickHPPopup = false
			app.updateModals()
		})

	return modal
}

// createMultiTargetModal creates a modal for the multi-target popup
func (app *App) createMultiTargetModal() *tview.Modal {
	content := ui.RenderMultiTargetPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Apply", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Apply" {
				// Handle multi-target action - this will be handled by input handler
			}
			app.model.ShowMultiTargetPopup = false
			app.updateModals()
		})

	return modal
}

// createConditionTextView creates a TextView widget for the condition popup (allows keyboard navigation)
func (app *App) createConditionTextView() tview.Primitive {
	content := ui.RenderConditionPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	colorConverter := NewColorConverter(app.model.Config)

	conditionView := tview.NewTextView()
	conditionView.SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true). // Enable scrolling for long lists
		SetTextAlign(tview.AlignLeft)
	conditionView.SetBorder(true).
		SetTitle(" Conditions ").
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderColor(colorConverter.PrimaryColor()).
		SetTitleColor(colorConverter.PrimaryColor())
	conditionView.SetText(coloredContent)

	// Don't set SetInputCapture on the TextView - let application-level handler process all input
	// This allows the handler chain to process arrow keys, Enter, Esc, etc.

	// Wrap in a Flex container to center it and make it larger
	// Make it take up more of the screen (about 70% height, 60% width)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer (flexible, smaller)
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 2, false). // Left spacer (flexible)
				AddItem(conditionView, 0, 3, false). // Condition view (larger, 3x weight for width)
				AddItem(nil, 0, 2, false), // Right spacer (flexible)
			0, 4, false). // Center row (larger, 4x weight for more height)
		AddItem(nil, 0, 1, false) // Bottom spacer (flexible, smaller)

	return flex
}

// createCastSpellModal creates a modal for the cast spell popup
func (app *App) createCastSpellModal() *tview.Modal {
	content := ui.RenderCastSpellPopup(app.model.SpellSearchInput, app.model.SpellSearchInput)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Cast", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cast" {
				// Handle spell casting - this will be handled by input handler
			}
			app.model.ShowCastSpellPrompt = false
			app.updateModals()
		})

	return modal
}

// createSavingThrowModal creates a modal for the saving throw popup
func (app *App) createSavingThrowModal() *tview.Modal {
	content := ui.RenderSavingThrowPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Roll", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Roll" {
				// Handle saving throw - this will be handled by input handler
			}
			app.model.ShowSavingThrowPopup = false
			app.updateModals()
		})

	return modal
}

// createEncounterPromptModal creates a modal for the encounter prompt popup
func (app *App) createEncounterPromptModal() *tview.Modal {
	content := ui.RenderEncounterPromptPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				// Handle encounter save - this will be handled by input handler
			}
			app.model.ShowEncounterPrompt = false
			app.updateModals()
		})

	return modal
}

// createEncounterGeneratorModal creates a modal for the encounter generator popup
func (app *App) createEncounterGeneratorModal() *tview.Modal {
	content := ui.RenderGeneratorPopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Generate", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Generate" {
				// Handle encounter generation - this will be handled by input handler
			}
			app.model.EncounterGenerating = false
			app.updateModals()
		})

	return modal
}

// createRenameModal creates a modal for the rename popup
func (app *App) createRenameModal() *tview.Modal {
	content := ui.RenderRenamePopup(*app.model)

	// Wrap content with grey color tag for TView
	coloredContent := "[grey]" + content + "[white]"

	modal := tview.NewModal().
		SetText(coloredContent).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Rename", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Rename" {
				// Handle rename - this will be handled by input handler
			}
			app.model.ShowRenamePopup = false
			app.updateModals()
		})

	return modal
}
