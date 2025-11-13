# TView Example

This example demonstrates how to recreate a similar panel layout using `tview` instead of `Bubble Tea + Lipgloss`.

## Running the Example

First, install the required dependencies:

```bash
go mod init tview-example  # if not already initialized
go get github.com/rivo/tview
go get github.com/gdamore/tcell/v2
```

Then run:

```bash
go run tview_example.go
```

## Key Differences: TView vs Bubble Tea + Lipgloss

### TView Approach

**Pros:**
- ✅ Built-in panel/widget system - no manual border calculations
- ✅ Automatic layout management with Grid
- ✅ Less code for basic layouts
- ✅ Built-in focus management
- ✅ More widget-focused (tables, forms, etc.)

**Cons:**
- ❌ More imperative style (less functional)
- ❌ Less control over exact rendering
- ❌ Different state management pattern
- ❌ Would require significant refactoring of existing code

### Current Bubble Tea + Lipgloss Approach

**Pros:**
- ✅ Functional/declarative architecture (MVU pattern)
- ✅ Full control over rendering
- ✅ Consistent with current codebase
- ✅ Better for complex state management
- ✅ Lipgloss provides fine-grained styling control

**Cons:**
- ❌ Manual layout calculations (borders, spacing, etc.)
- ❌ More code for basic layouts
- ❌ Need to handle focus/input routing manually

## Example Features

The example demonstrates:
1. **Grid Layout**: Two rows with flexible column widths
2. **Panel Borders**: Automatic border rendering with tview
3. **Active Panel Highlighting**: Border color changes based on focus
4. **Navigation**: Tab/Shift+Tab to cycle through panels
5. **Status Bar**: Bottom status bar with keybindings

## Layout Structure

```
┌─────────────────────┬──────────────┐
│  1. Dice Roller     │  2. Initiative│
│                     │   Tracker    │
├──────┬──────┬───────┼──────────────┤
│ 3.   │ 4.   │ 5.    │  6. Encounter│
│Spells│Monst.│ Notes │   Builder    │
└──────┴──────┴───────┴──────────────┘
```

## Migration Considerations

If you were to migrate to tview, you would need to:

1. **Refactor State Management**: Convert from MVU pattern to tview's widget-based approach
2. **Rewrite Handlers**: Convert input handlers to tview's SetInputCapture pattern
3. **Rework Layout Logic**: Replace manual calculations with Grid.SetColumns/SetRows
4. **Update Content Rendering**: Use tview widgets (TextView, Table, etc.) instead of manual string building
5. **Change Styling**: Use tview's styling system instead of Lipgloss

**Estimated Effort**: Significant (weeks of work for a full migration)

## Recommendation

Unless you're experiencing specific issues with Bubble Tea + Lipgloss, staying with the current approach is recommended because:
- The codebase is already well-structured
- The border alignment issues can be fixed (as we're doing)
- The functional architecture fits well with the application's complexity
- Migration would be a major undertaking with limited benefits

