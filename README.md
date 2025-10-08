# LazyD&D - Terminal-based D&D Panel System for Dungeon Master

A lazygit-inspired terminal UI for managing your D&D game sessions, built with Go and Bubble Tea.

![LazyD&D Screenshot](screenshot/screenshot.png)

## Features

🎲 **Dice Roller Panel** - Roll any dice with simple commands (2d6, 1d20+5, etc.)
⚔️ **Initiative Tracker Panel** - Manage combat initiative for players and monsters
✨ **Spells Panel** - Search and browse D&D 5e spells with autocomplete
🐲 **Monsters Panel** - Search and view detailed monster stat blocks with actions
💾 **Campaign Save/Load** - Save your game state and resume later
🔄 **Auto-Save** - Automatic saving every 5 minutes
🔗 **Monster Integration** - Link monsters to initiative with full action support

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

### Global Navigation & Controls

#### Panel Navigation
| Key | Action |
|-----|--------|
| `1-4` or `F1-F4` | Jump directly to specific panel |
| `Tab` | Cycle forward through panels |
| `Shift+Tab` | Cycle backward through panels |
| `↑` `↓` | Navigate lists and scroll content |
| `Esc` | Cancel input or exit current mode |
| `?` | Show help popup with all keybindings |
| `q` or `Ctrl+C` | Quit application |

#### Campaign Management
| Key | Action |
|-----|--------|
| `Ctrl+S` | Save campaign (quick save if already saved, or open save dialog) |
| `Ctrl+L` | Load campaign (shows list of saved campaigns) |
| `Ctrl+N` | Rename current campaign |

**Campaign Features:**
- 💾 **Auto-Save**: Campaigns automatically save every 5 minutes
- 📁 **Save Location**: All campaigns stored in `~/.lazydnd/`
- 📊 **Status Bar**: Shows current campaign name and last save time
- 🔄 **Full State**: Saves entire initiative tracker with all monster links

---

### Panel-Specific Features

#### 🎲 Panel 1: Dice Roller

**Keybindings:**
| Key | Action |
|-----|--------|
| `Enter` | Start dice input mode |
| `r` | Reroll last command |
| Type dice notation and press `Enter` | Roll dice |

**Dice Notation Examples:**
- **Simple Rolls**: `1d20`, `2d6`, `3d8`, `4d10`
- **With Modifiers**: `2d6+3`, `1d20-1`, `1d8+5`
- **Multiple Dice**: `2d8+3d6`, `1d6-1d4`, `2d8+3d6-1d4`
- **Complex Expressions**: `1d6+3+2d8-5`, `2d6-1d4+3`
- **Comma-Separated**: `1d8+3, 3d6-1` (rolls multiple expressions)
- **Advantage/Disadvantage**: `1d20 adv`, `2d6 dis`

**Features:**
- ✅ Minimum value of 1 (D&D rule)
- ✅ Available dice: d4, d6, d8, d10, d12, d20, d100
- ✅ Roll history displayed
- ✅ Quick reroll with 'r' key

---

#### ⚔️ Panel 2: Initiative Tracker

![LazyD&D Attack](screenshot/attack.png)

**Keybindings:**
| Key | Action |
|-----|--------|
| `p` | Add player to initiative |
| `m` | Add monster manually (enter name, HP, AC, initiative) |
| `e` | Enter edit mode (navigate entries with ↑↓) |
| `i` | Edit initiative value (in edit mode) |
| `h` | Edit HP - add/remove HP with +/- (in edit mode) |
| `d` | Delete selected entry (in edit mode) |
| `l` | View linked monster details (if added from Monster panel) |
| `a` | Show monster actions popup (if monster has actions) |
| `c` | Copy/duplicate selected entry |

**Features:**
- ✅ **Auto-Sort**: Entries sorted by initiative (highest first)
- ✅ **Monster Linking**: Monsters added from Monster panel retain full data
- ✅ **Action Integration**: Press 'a' on linked monsters to see available actions
- ✅ **Quick Actions**: Select action to auto-roll damage in Dice Roller
- ✅ **HP Tracking**: Real-time HP management for monsters
- ✅ **Duplicate Monsters**: Copy entries with automatic numbering (Goblin 1, Goblin 2, etc.)
- ✅ **Campaign Save**: All initiative data saved with campaign

**Adding Monsters:**
1. **Manual**: Press 'm' and enter details (no action support)
2. **From Monster Panel**: Search monster, press 'a' to add with full stats and actions

---

#### ✨ Panel 3: Spells

**Keybindings:**
| Key | Action |
|-----|--------|
| `Enter` | Start spell search |
| Type spell name | Real-time autocomplete suggestions |
| `↑` `↓` | Navigate spell suggestions |
| `Enter` | Select spell to view details |
| `Esc` | Exit search mode |

**Features:**
- ✅ **Complete D&D 5e Spell Database**
- ✅ **Real-time Autocomplete**: Suggestions appear as you type
- ✅ **Full Details**: Level, school, casting time, range, components, duration, description
- ✅ **Class Information**: Shows which classes can cast the spell
- ✅ **Ritual & Concentration**: Clearly marked

---

#### 🐉 Panel 4: Monsters

**Keybindings:**
| Key | Action |
|-----|--------|
| `Enter` | Start monster search |
| Type monster name | Real-time autocomplete suggestions |
| `↑` `↓` | Navigate monster suggestions |
| `Enter` | Select monster to view stat block |
| `a` | Add selected monster to Initiative Tracker (with full actions) |
| `Esc` | Exit search mode |

**Features:**
- ✅ **8750+ D&D 5e Monsters**
- ✅ **Complete Stat Blocks**: AC, HP, Speed, Ability Scores, Saves, Skills
- ✅ **Traits & Actions**: All special abilities and attacks
- ✅ **Legendary Actions**: Full legendary action details
- ✅ **Challenge Rating**: CR and XP values
- ✅ **Action Parser**: Structured action data with damage, reach, save DCs
- ✅ **Initiative Integration**: Press 'a' to add monster to initiative with all data linked

**Monster Actions:**
When a monster is added to initiative from the Monster panel:
- Press 'a' in Initiative Tracker to see action list
- Select action to automatically roll damage in Dice Roller
- Actions include: attack bonus, reach/range, damage dice, damage type, save DCs

---

### Campaign Management Details

![LazyD&D Save](screenshot/save.png)

**Saving Your Campaign:**
1. Press `Ctrl+S` to save
2. Enter campaign name (e.g., "dragon_heist")
3. Campaign saved to `~/.lazydnd/dragon_heist.json`
4. Status bar shows: `📁 dragon_heist (💾 Just now)`

**Loading a Campaign:**
1. Press `Ctrl+L` to open load dialog
2. Use `↑` `↓` to navigate saved campaigns
3. Press `Enter` to load
4. All initiative data and monster links restored

**Renaming a Campaign:**
1. Press `Ctrl+N` (only available when campaign is loaded)
2. Enter new name
3. Old file deleted, new file created

**Auto-Save:**
- Automatically saves every 5 minutes when campaign is loaded
- Status bar shows time since last save: `💾 3m ago`
- Manual save with `Ctrl+S` updates to `💾 Just now`

**What Gets Saved:**
- ✅ All initiative tracker entries (players and monsters)
- ✅ HP, AC, initiative values
- ✅ Monster links (actions remain available after load)
- ✅ Instance numbers for duplicated monsters
- ❌ Dice history (session-specific)
- ❌ Current panel or scroll position

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
