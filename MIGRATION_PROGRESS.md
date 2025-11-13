# TView Migration Progress

## Phase 1: Setup & Dependencies ✅ COMPLETED

- [x] Added TView and TCell dependencies to `go.mod`
- [x] Created directory structure:
  - `ui/tview/` - Main TView package
  - `ui/tview/panels/` - Panel widgets
  - `ui/tview/popups/` - Popup widgets (to be implemented)
- [x] Created basic application skeleton:
  - `ui/tview/app.go` - Main TView application wrapper
  - `ui/tview/handlers.go` - Input handler routing
  - `ui/tview/panels/*.go` - All 6 panel widgets

## Phase 2: Core Application Structure 🔄 IN PROGRESS

- [x] Created `main_tview.go` - Alternative entry point for TView
- [x] Created TView application wrapper (`ui/tview/app.go`)
- [ ] Implement window resize handling
- [ ] Test basic application startup

## Phase 3: Panel Widgets ⏳ PENDING

- [x] Created panel widget stubs for all 6 panels
- [ ] Implement proper content rendering
- [ ] Handle panel borders and titles
- [ ] Implement active/inactive styling

## Phase 4: Input Handling ⏳ PENDING

- [x] Created basic input handler structure
- [ ] Convert handler chain to TView input capture
- [ ] Map all key handlers to TView events
- [ ] Implement input modes

## Phase 5-10: ⏳ PENDING

See `MIGRATION_PLAN.md` for detailed phases.

## Current Status

**Files Created:**
- `ui/tview/app.go` - Main application wrapper
- `ui/tview/handlers.go` - Input handler routing
- `ui/tview/panels/dice_roller.go` - Dice roller panel
- `ui/tview/panels/initiative.go` - Initiative tracker panel
- `ui/tview/panels/spells.go` - Spells panel
- `ui/tview/panels/monsters.go` - Monsters panel
- `ui/tview/panels/notes.go` - Notes panel
- `ui/tview/panels/encounter.go` - Encounter builder panel
- `main_tview.go` - Alternative entry point

**Next Steps:**
1. Run `go mod tidy` to fetch dependencies
2. Test compilation: `go build -o lazydnd-tview main_tview.go`
3. Implement proper handler routing (Phase 4)
4. Test basic functionality
5. Continue with remaining phases

## Testing

To test the TView implementation:

```bash
# Fetch dependencies
go mod tidy

# Build TView version
go build -o lazydnd-tview main_tview.go

# Run (if build succeeds)
./lazydnd-tview
```

## Notes

- The TView implementation runs in parallel with the existing Bubble Tea code
- Original `main.go` remains unchanged
- Both implementations can coexist during migration
- Once migration is complete, `main_tview.go` can replace `main.go`
