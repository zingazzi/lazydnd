# LazyD&D - Terminal-based D&D Panel System

A lazygit-inspired terminal UI for managing your D&D game sessions, built with Go and Bubble Tea.

## Features

🎲 **Dice Roller Panel** - Roll any dice with simple commands
⚔️ **Character Sheet Panel** - View character stats and information
✨ **Spells Panel** - Browse available spells and spell slots
📖 **Campaign Notes Panel** - Keep track of session notes and NPCs

## Installation

```bash
# Clone and build
git clone <your-repo>
cd lazydnd
go build -o lazydnd
./lazydnd
```

## Usage

### Navigation
- **Numbers 1-4**: Jump directly to panels
- **Tab**: Cycle through panels
- **q**: Quit application

### Dice Roller Commands
- `d4`, `d6`, `d8`, `d10`, `d12`, `d20` - Roll single dice
- `2d6`, `3d8` - Roll multiple dice
- `2d20+5` - Roll with modifiers
- `3d6-1` - Roll with negative modifiers

### Panels

1. **Dice Roller** - Interactive dice rolling with history
2. **Character Sheet** - Character stats and information
3. **Spells** - Spell lists and slot tracking
4. **Campaign Notes** - Session notes and NPC tracking

## Requirements

- Go 1.19 or higher
- Terminal with color support

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling

## License

MIT License
