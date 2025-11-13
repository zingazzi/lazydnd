# Quick Start: TView Migration

This guide will help you start the migration from Bubble Tea to TView.

## Current Branch

You're on: `feature/tview-migration`

## Step 1: Install Dependencies

```bash
go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2
go mod tidy
```

## Step 2: Create Directory Structure

```bash
mkdir -p ui/tview/panels
mkdir -p ui/tview/popups
```

## Step 3: Start with Phase 1

Create the basic application skeleton:

1. **Create `ui/tview/app.go`** - Main TView application wrapper
2. **Create `ui/tview/types.go`** - TView-specific types
3. **Create `ui/tview/layout.go`** - Grid layout setup

## Step 4: Test Basic Setup

Create a minimal test to verify TView works:

```go
// ui/tview/app_test.go (or just run a simple main)
package tview

import (
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
)

func TestBasicApp() {
    app := tview.NewApplication()
    text := tview.NewTextView().SetText("Hello TView!")
    if err := app.SetRoot(text, true).Run(); err != nil {
        panic(err)
    }
}
```

## Step 5: Begin Panel Migration

Start with the simplest panel (Dice Roller):

1. Create `ui/tview/panels/dice_roller.go`
2. Convert content rendering
3. Test in isolation
4. Integrate into main app

## Development Workflow

1. **Work incrementally**: One panel at a time
2. **Test frequently**: Run the app after each change
3. **Keep original code**: Don't delete Bubble Tea code until migration is complete
4. **Use feature flags**: Consider build tags to switch between implementations

## Build Tags (Optional)

To support both implementations during migration:

```go
// +build !tview

// Original Bubble Tea code
```

```go
// +build tview

// New TView code
```

Build with: `go build -tags tview`

## Reference Files

- `MIGRATION_PLAN.md` - Full migration plan
- `examples/tview_example.go` - Working TView example
- `examples/README.md` - Comparison guide

## Getting Help

- TView Docs: https://pkg.go.dev/github.com/rivo/tview
- TView Examples: https://github.com/rivo/tview/tree/master/_examples
- TCell Docs: https://pkg.go.dev/github.com/gdamore/tcell/v2

## Next Steps

1. Review `MIGRATION_PLAN.md` for detailed phases
2. Start with Phase 1 (Setup & Dependencies)
3. Create basic skeleton
4. Test with one simple panel
5. Iterate and expand

