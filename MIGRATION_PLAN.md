# TView Migration Plan

## Overview

This document outlines the plan to migrate LazyDnD from **Bubble Tea + Lipgloss** to **TView** for the terminal UI framework.

**Current Stack:**
- Bubble Tea (MVU architecture)
- Lipgloss (styling)
- Manual layout calculations
- Custom border/panel management

**Target Stack:**
- TView (widget-based architecture)
- TCell (low-level terminal library, used by TView)
- Automatic layout management
- Built-in panel/widget system

## Goals

1. ✅ Eliminate manual border/width calculations
2. ✅ Simplify layout management
3. ✅ Maintain all existing functionality
4. ✅ Preserve user experience
5. ✅ Improve code maintainability

## Architecture Comparison

### Current Architecture (Bubble Tea)

```
main.go
  └─> tea.NewProgram(model)
      └─> Model.Update(msg) → (Model, Cmd)
      └─> Model.View() → string
          └─> Manual layout calculations
          └─> Lipgloss styling
          └─> String concatenation
```

**Key Characteristics:**
- Functional/declarative (MVU pattern)
- State is immutable
- Commands for side effects
- Manual rendering pipeline

### Target Architecture (TView)

```
main.go
  └─> tview.NewApplication()
      └─> Grid/Primitive widgets
          └─> SetInputCapture() for handlers
          └─> Automatic layout
          └─> Built-in styling
```

**Key Characteristics:**
- Imperative/widget-based
- State is mutable
- Direct event handlers
- Automatic rendering

## Migration Phases

### Phase 1: Setup & Dependencies (1-2 hours)

**Tasks:**
- [ ] Add TView and TCell dependencies
- [ ] Create new `ui/tview/` package structure
- [ ] Set up basic TView application skeleton
- [ ] Create parallel implementation alongside existing code

**Files to Create:**
- `ui/tview/app.go` - Main TView application
- `ui/tview/panels.go` - Panel widget definitions
- `ui/tview/handlers.go` - Input handlers
- `ui/tview/layout.go` - Grid layout setup

**Dependencies:**
```go
go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2
```

### Phase 2: Core Application Structure (4-6 hours)

**Tasks:**
- [ ] Convert `main.go` to use TView
- [ ] Create TView application wrapper
- [ ] Implement basic window resize handling
- [ ] Set up grid layout structure

**Changes:**
- `main.go`: Replace `tea.NewProgram` with `tview.NewApplication`
- Create `ui/tview/app.go` with application lifecycle
- Implement terminal resize handling

### Phase 3: Panel Widgets (8-12 hours)

**Tasks:**
- [ ] Convert each panel to TView widget
- [ ] Implement panel content rendering
- [ ] Set up panel borders and titles
- [ ] Handle active/inactive panel styling

**Panel Migration Order:**
1. Dice Roller (simplest)
2. Notes (simple text)
3. Initiative Tracker (table-like)
4. Spells (list with search)
5. Monsters (list with search)
6. Encounter Builder (complex)

**Files to Modify:**
- `ui/tview/panels/dice_roller.go`
- `ui/tview/panels/initiative.go`
- `ui/tview/panels/spells.go`
- `ui/tview/panels/monsters.go`
- `ui/tview/panels/notes.go`
- `ui/tview/panels/encounter.go`

### Phase 4: Input Handling (6-8 hours)

**Tasks:**
- [ ] Convert handler chain to TView input capture
- [ ] Map all key handlers to TView events
- [ ] Implement input modes (edit, search, etc.)
- [ ] Handle special input states

**Key Changes:**
- Replace `HandleNavigation()` with `SetInputCapture()`
- Convert `tea.KeyMsg` to `*tcell.EventKey`
- Map handler chain to TView event routing

**Files to Modify:**
- `ui/tview/handlers.go` - Main input handler
- `ui/tview/handlers_*.go` - Specific handler groups
- Convert all `ui/handlers_*.go` files

### Phase 5: State Management (4-6 hours)

**Tasks:**
- [ ] Adapt Model struct for TView (make mutable)
- [ ] Convert state updates to direct mutations
- [ ] Remove Bubble Tea command system
- [ ] Implement state synchronization

**Key Changes:**
- Model becomes mutable (TView pattern)
- Remove `tea.Cmd` return types
- Direct state updates instead of returning new Model
- Keep existing state structs (DiceRollerState, etc.)

### Phase 6: Popups & Modals (6-8 hours)

**Tasks:**
- [ ] Convert all popups to TView modals/pages
- [ ] Implement popup overlay system
- [ ] Handle popup input separately
- [ ] Maintain popup state management

**Popups to Convert:**
- Save/Load popup
- Quick HP popup
- Encounter prompt popup
- Encounter generator popup
- Help popup
- Action popup
- Condition popup
- Multi-target popup
- Cast spell popup
- Saving throw popup

**Files:**
- `ui/tview/popups/` - New popup implementations
- Convert all `ui/*_popup.go` files

### Phase 7: Layout & Styling (4-6 hours)

**Tasks:**
- [ ] Set up grid layout with proper proportions
- [ ] Configure panel borders and colors
- [ ] Implement active panel highlighting
- [ ] Add status bar widget
- [ ] Handle responsive sizing

**Key Features:**
- Use TView Grid for layout
- Automatic border management
- Color themes from config
- Responsive column/row sizing

### Phase 8: Content Rendering (6-8 hours)

**Tasks:**
- [ ] Convert content providers to TView widgets
- [ ] Implement scrolling for panels
- [ ] Handle text wrapping/truncation
- [ ] Add search highlighting
- [ ] Implement table/list views

**Content Types:**
- Text content (Dice, Notes)
- Lists (Spells, Monsters)
- Tables (Initiative)
- Forms (Encounter Builder)

### Phase 9: Advanced Features (8-10 hours)

**Tasks:**
- [ ] Autosave integration
- [ ] Error message display
- [ ] Debug mode support
- [ ] Configuration integration
- [ ] Save/load functionality
- [ ] Undo/redo system

### Phase 10: Testing & Polish (6-8 hours)

**Tasks:**
- [ ] Test all functionality
- [ ] Fix layout issues
- [ ] Verify all keybindings
- [ ] Test on different terminals
- [ ] Performance optimization
- [ ] Code cleanup

## File-by-File Migration Strategy

### Core Files

| Current File | TView Equivalent | Priority | Complexity |
|-------------|------------------|----------|------------|
| `main.go` | `main.go` + `ui/tview/app.go` | High | Medium |
| `ui/model.go` | `ui/tview/app.go` | High | Medium |
| `ui/types.go` | `ui/tview/types.go` | High | Low |
| `ui/layout_core.go` | `ui/tview/layout.go` | High | Low |
| `ui/layout_panel.go` | `ui/tview/panels/*.go` | High | High |
| `ui/layout_content.go` | `ui/tview/panels/*.go` | High | Medium |
| `ui/layout_scroll.go` | Built into TView | Medium | Low |
| `ui/layout_statusbar.go` | `ui/tview/statusbar.go` | Medium | Low |
| `ui/styles.go` | `ui/tview/styles.go` | Medium | Low |

### Handler Files

| Current File | TView Equivalent | Priority | Complexity |
|-------------|------------------|----------|------------|
| `ui/handlers_core.go` | `ui/tview/handlers.go` | High | Medium |
| `ui/handler_chain.go` | `ui/tview/handlers.go` | High | Medium |
| `ui/handlers_*.go` (15 files) | `ui/tview/handlers_*.go` | High | High |

### Popup Files

| Current File | TView Equivalent | Priority | Complexity |
|-------------|------------------|----------|------------|
| `ui/*_popup.go` (10 files) | `ui/tview/popups/*.go` | Medium | Medium |

### Utility Files

| Current File | TView Equivalent | Priority | Complexity |
|-------------|------------------|----------|------------|
| `ui/text_utils.go` | Keep or adapt | Low | Low |
| `ui/terminal_utils.go` | Keep or adapt | Low | Low |
| `ui/layout_constants.go` | `ui/tview/constants.go` | Low | Low |

## Key Implementation Details

### 1. Application Structure

```go
// ui/tview/app.go
type App struct {
    application *tview.Application
    grid        *tview.Grid
    panels      []*tview.TextView
    model       *Model  // Mutable reference
    activePanel PanelType
}

func NewApp(model *Model) *App {
    app := &App{
        application: tview.NewApplication(),
        model:       model,
    }
    app.setupLayout()
    app.setupHandlers()
    return app
}
```

### 2. Panel Creation

```go
// ui/tview/panels/dice_roller.go
func NewDiceRollerPanel(model *Model) *tview.TextView {
    panel := tview.NewTextView()
    panel.SetTitle(" 1. Dice Roller ")
    panel.SetBorder(true)
    panel.SetDynamicColors(true)
    panel.SetWrap(true)

    // Update content based on model state
    updateContent := func() {
        content := getDiceRollerContent(model)
        panel.SetText(content)
    }

    // Register update function
    model.RegisterUpdateCallback(updateContent)

    return panel
}
```

### 3. Input Handling

```go
// ui/tview/handlers.go
func (app *App) setupHandlers() {
    app.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        // Convert to handler chain
        key := event.Key()
        rune := event.Rune()

        // Route to appropriate handler
        if handler := getHandler(key, rune); handler != nil {
            handler(app.model)  // Direct mutation
            app.application.Draw()  // Trigger redraw
            return nil
        }

        return event
    })
}
```

### 4. State Updates

```go
// Instead of: return newModel, cmd
// Use: model.UpdateField(...)

// Old (Bubble Tea):
func handleTab(m Model) (Model, tea.Cmd) {
    m.ActivePanel = (m.ActivePanel + 1) % 6
    return m, nil
}

// New (TView):
func handleTab(model *Model) {
    model.ActivePanel = (model.ActivePanel + 1) % 6
    updateActivePanelBorder(model.ActivePanel)
}
```

### 5. Layout Management

```go
// ui/tview/layout.go
func (app *App) setupLayout() {
    app.grid = tview.NewGrid().
        SetRows(0, 0, 1).  // Two panel rows + status bar
        SetColumns(0, 0, 0, 0, 0).  // Flexible columns
        SetBorders(false)

    // Top row: Dice (60%) + Initiative (40%)
    app.grid.AddItem(app.panels[DiceRoller], 0, 0, 1, 3, 0, 0, false).
        AddItem(app.panels[InitiativeTracker], 0, 3, 1, 2, 0, 0, false).
        // Bottom row: Spells + Monsters + Notes + Encounter
        AddItem(app.panels[Spells], 1, 0, 1, 1, 0, 0, false).
        AddItem(app.panels[Monsters], 1, 1, 1, 1, 0, 0, false).
        AddItem(app.panels[Notes], 1, 2, 1, 1, 0, 0, false).
        AddItem(app.panels[EncounterBuilder], 1, 3, 1, 2, 0, 0, false).
        // Status bar
        AddItem(app.statusBar, 2, 0, 1, 5, 0, 0, false)
}
```

## Testing Strategy

### Unit Tests
- Test individual panel widgets
- Test input handlers
- Test state mutations
- Test layout calculations

### Integration Tests
- Test full application flow
- Test all keybindings
- Test popup interactions
- Test save/load functionality

### Manual Testing Checklist
- [ ] All panels render correctly
- [ ] Navigation works (Tab, Shift+Tab, number keys)
- [ ] All input modes work
- [ ] All popups display correctly
- [ ] Save/load works
- [ ] Autosave works
- [ ] Error messages display
- [ ] Status bar updates
- [ ] Terminal resize works
- [ ] Works on different terminals (iTerm, Terminal.app, Alacritty, etc.)

## Rollback Plan

1. **Keep original branch**: `feature/improvment` remains untouched
2. **Feature flag**: Add build tag to switch between implementations
3. **Gradual migration**: Migrate one panel at a time
4. **Parallel implementation**: Keep both implementations during migration

## Estimated Timeline

- **Total Time**: 50-70 hours
- **Phases 1-3** (Setup + Core + Panels): 15-20 hours
- **Phases 4-5** (Handlers + State): 10-14 hours
- **Phases 6-7** (Popups + Layout): 10-14 hours
- **Phases 8-9** (Content + Features): 14-18 hours
- **Phase 10** (Testing): 6-8 hours

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Loss of functionality | High | Comprehensive testing, feature parity checklist |
| Performance issues | Medium | Benchmark both implementations |
| Different UX | Medium | Maintain same keybindings and behavior |
| Migration complexity | High | Phased approach, keep original branch |
| Terminal compatibility | Low | Test on multiple terminals early |

## Success Criteria

1. ✅ All existing functionality works
2. ✅ No manual border calculations needed
3. ✅ Layout is automatic and responsive
4. ✅ Code is more maintainable
5. ✅ Performance is equal or better
6. ✅ User experience is unchanged or improved

## Next Steps

1. Review and approve this plan
2. Start Phase 1: Setup & Dependencies
3. Create initial TView skeleton
4. Begin panel-by-panel migration
5. Test incrementally

## Notes

- Keep existing state structs (DiceRollerState, etc.) - they're well-designed
- TView widgets are mutable, but we can still use functional patterns for state updates
- Consider using TView's `Pages` primitive for popups
- Use `Flex` for status bar and other simple layouts
- Use `Table` widget for Initiative Tracker
- Use `List` widget for Spells/Monsters with search
