# LazyD&D - Terminal-based D&D Panel System for Dungeon Master

A lazygit-inspired terminal UI for managing your D&D game sessions, built with Go and Bubble Tea.

![LazyD&D Screenshot](screenshot.png)

## Features

🎲 **Dice Roller Panel** - Roll any dice with simple commands (2d6, 1d20+5, etc.)
⚔️ **Initiative Tracker Panel** - Manage combat initiative for players and monsters
✨ **Spells Panel** - Search and browse D&D 5e spells with autocomplete
🐲 **Monsters Panel** - Search and view detailed monster stat blocks

## Installation

### Docker (Easiest)

**Run with Docker (no installation needed):**
```bash
docker run -it --rm ghcr.io/zingazzi/lazydnd:latest
```

**Or build locally:**
```bash
git clone https://github.com/zingazzi/lazydnd
cd lazydnd
./docker-build.sh
./docker-run.sh
```

See [DOCKER.md](DOCKER.md) for detailed Docker documentation.

### Quick Install (Recommended)

**One-line installer (Linux/macOS):**
```bash
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash
```

Or download and run the installer:
```bash
curl -O https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh
chmod +x install.sh
./install.sh
```

### Manual Installation

Download the latest release for your platform from the [Releases page](https://github.com/zingazzi/lazydnd/releases/latest):

**Linux (Intel/AMD):**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-amd64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**Linux (ARM):**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-arm64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-macos-amd64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-macos-arm64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**Windows:**
Download `lazydnd-windows-amd64.exe` from the [Releases page](https://github.com/zingazzi/lazydnd/releases/latest) and add it to your PATH.

### Build from Source

```bash
git clone https://github.com/zingazzi/lazydnd
cd lazydnd
go build -o lazydnd
./lazydnd
```

### Cross-Platform Build

Build executables for all platforms:
```bash
./build.sh
```

This creates executables in the `build/` directory for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Usage

### Navigation
- **Numbers 1-4** or **F1-F4**: Jump directly to panels
- **Tab**: Cycle forward through panels
- **Shift+Tab**: Cycle backward through panels
- **↑↓**: Navigate lists and scroll content
- **Enter**: Activate input modes or select items
- **Esc**: Cancel input or exit modes
- **?**: Show help popup
- **q**: Quit application

### Panel Features

#### 1. Dice Roller
- **Simple Rolls**: `1d20`, `2d6`, `3d8`
- **With Modifiers**: `2d6+3`, `1d20-1`, `1d8-2`
- **Multiple Dice**: `2d8+3d6`, `1d6-1d4`, `2d8+3d6-1d4`
- **Complex Expressions**: `1d6+3+2d8-5`, `2d6-1d4+3`
- **Comma-Separated**: `1d8+3, 3d6-1` (rolls multiple expressions)
- **Advantage/Disadvantage**: `1d20 adv`, `2d6 dis`
- **Minimum Value**: Results never go below 1 (D&D rule)
- **Available Dice**: d4, d6, d8, d10, d12, d20, d100
- **History**: View recent rolls
- **Reroll**: Press 'r' to reroll last command

#### 2. Initiative Tracker (Dungeon Master Panel)
- **Add Players**: Press 'p' → enter name and initiative
- **Add Monsters**: Press 'm' → enter name, HP, AC, and initiative
- **Edit Mode**: Press 'e' to edit existing entries
  - 'i': Edit initiative
  - 'h': Edit monster HP (+heal/-damage)
  - 'd': Delete entry
- **Auto-Sort**: Entries sorted by initiative (highest first)

#### 3. Spells Panel
- **Search**: Press Enter → type spell name
- **Autocomplete**: Real-time suggestions as you type
- **Details**: Full spell descriptions, components, duration, etc.
- **Database**: Complete D&D 5e spell compendium

#### 4. Monsters Panel
- **Search**: Press Enter → type monster name
- **Autocomplete**: Find monsters quickly
- **Stat Blocks**: Complete monster information including:
  - Ability scores and modifiers
  - AC, HP, Speed, Challenge Rating
  - Traits, Actions, Legendary Actions
- **Database**: 8750+ D&D 5e monsters

## Requirements

- Go 1.19 or higher
- Terminal with color support

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling

## Contribute

We welcome contributions! Whether it's bug fixes, new features, documentation, or suggestions, your help is appreciated.

**How to contribute:**
1. [Fork the repository](https://github.com/zingazzi/lazydnd/fork)
2. Create a new branch: `git checkout -b feature/your-feature`
3. Make your changes and commit: `git commit -am 'Add new feature'`
4. Push to your fork: `git push origin feature/your-feature`
5. [Open a Pull Request](https://github.com/zingazzi/lazydnd/pulls)

Please review the [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for code organization and best practices.

### Reporting Issues

If you find a bug, have a feature request, or need help:
- [Open an Issue](https://github.com/zingazzi/lazydnd/issues)
- Include clear steps to reproduce, expected behavior, and screenshots/logs if possible.

Thank you for helping improve LazyDnD!

## Enjoy LazyDnd? Buy me a coffee

If you find LazyDnD helpful, consider supporting development!  
[☕ Buy me a coffee](https://www.buymeacoffee.com/zingazzi)  
Your support helps keep the project alive and growing. Thank you!


## License

MIT License
