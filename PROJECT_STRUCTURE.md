# LazyD&D Project Structure

## 📁 Directory Organization

```
lazydnd/
├── main.go                    # Entry point - minimal, just starts the app
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies
├── README.md                  # Project documentation
├── PROJECT_STRUCTURE.md       # This file
├── ui/                        # UI package - handles interface and layout
│   ├── types.go              # Data structures and constants
│   ├── styles.go             # Lipgloss styling definitions
│   ├── layout.go             # Panel layout and rendering logic
│   ├── navigation.go         # Navigation and key handling
│   └── model.go              # Bubble Tea model implementation
└── panels/                    # Panels package - content for each panel
    ├── dice_roller.go        # Dice rolling logic and content
    ├── character_sheet.go    # Character sheet content
    ├── spells.go             # Spells panel content
    └── campaign_notes.go     # Campaign notes content
```

## 📦 Package Responsibilities

### `main` package
- **Purpose**: Application entry point
- **Files**: `main.go`
- **Responsibilities**:
  - Initialize the UI model
  - Start the Bubble Tea program
  - Handle application-level errors

### `ui` package
- **Purpose**: User interface management
- **Files**: `types.go`, `styles.go`, `layout.go`, `navigation.go`, `model.go`
- **Responsibilities**:
  - Define application state and panel types
  - Handle visual styling and theming
  - Manage 2x2 panel layout and scrolling
  - Process navigation and keyboard input
  - Implement Bubble Tea model interface

### `panels` package
- **Purpose**: Panel content and logic
- **Files**: `dice_roller.go`, `character_sheet.go`, `spells.go`, `campaign_notes.go`
- **Responsibilities**:
  - Generate content for each panel
  - Handle panel-specific logic (e.g., dice rolling)
  - Maintain panel data and state

## 🔄 Data Flow

1. **User Input** → `ui/navigation.go` → processes keys and updates model
2. **Model Update** → `ui/model.go` → handles Bubble Tea messages
3. **Rendering** → `ui/layout.go` → creates 2x2 panel layout
4. **Content Generation** → `panels/*.go` → provides panel-specific content
5. **Styling** → `ui/styles.go` → applies visual styling

## 🎯 Benefits of This Structure

- **Separation of Concerns**: UI logic separate from content logic
- **Modularity**: Each panel is self-contained
- **Maintainability**: Easy to add new panels or modify existing ones
- **Testability**: Each package can be tested independently
- **Scalability**: Simple to extend with new features

## 🚀 Adding New Panels

1. Create new panel file in `panels/` directory
2. Add panel type to `ui/types.go`
3. Update panel names array in `ui/types.go`
4. Add content function call in `ui/layout.go`
5. Update navigation in `ui/navigation.go`
