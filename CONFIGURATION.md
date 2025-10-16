# Configuration Guide

LazyDnD uses a JSON configuration file to customize its behavior and appearance.

## Configuration Location

The configuration file is located at:
```
~/.config/lazydnd/config.json
```

On first run, LazyDnD will automatically create this file with default values.

## Configuration Options

### Theme Settings

Control the application's color scheme. Colors use hex format (`#RRGGBB`).

```json
{
  "theme": {
    "primary_color": "#7D56F4",     // Main accent color (purple) - active panels, title bars, status bar
    "border_color": "#444444",       // Inactive panel border color (dark gray)
    "highlight_color": "#00FF00",    // Highlighted text, dice results (green)
    "error_color": "#FF0000",        // Error messages (red)
    "success_color": "#00FF00"       // Success messages (green)
  }
}
```

**What Changes When:**
- `primary_color`: Active panel borders, panel title backgrounds, status bar, input borders, help popup accents
- `border_color`: Inactive panel borders
- `highlight_color`: Dice results, selected items in lists, section headers in help
- `error_color`: Error messages
- `success_color`: Success notifications

**Popular Color Schemes:**

**Purple Theme (Default):**
- Primary: `#7D56F4` (purple)
- Border: `#444444` (dark gray)
- Highlight: `#00FF00` (green)

**Blue Theme:**
- Primary: `#00BFFF` (deep sky blue)
- Border: `#4169E1` (royal blue)
- Highlight: `#FFD700` (gold)

**Green Theme:**
- Primary: `#00FF7F` (spring green)
- Border: `#228B22` (forest green)
- Highlight: `#FFFF00` (yellow)

**Red Theme:**
- Primary: `#FF4500` (orange red)
- Border: `#DC143C` (crimson)
- Highlight: `#FFD700` (gold)

**Dark Minimal:**
- Primary: `#FFFFFF` (white)
- Border: `#666666` (gray)
- Highlight: `#AAAAAA` (light gray)

**Monokai:**
- Primary: `#66D9EF` (cyan)
- Border: `#75715E` (brown-gray)
- Highlight: `#A6E22E` (lime green)
- Error: `#F92672` (pink-red)
- Success: `#A6E22E` (lime green)

**Tokyo Night:**
- Primary: `#7aa2f7` (bright blue)
- Border: `#414868` (dark blue-gray)
- Highlight: `#9ece6a` (green)
- Error: `#f7768e` (red)
- Success: `#9ece6a` (green)

### Auto-Save Settings

Configure campaign auto-save behavior:

```json
{
  "auto_save": {
    "enabled": true,           // Enable/disable auto-save
    "interval_minutes": 5      // Save every N minutes
  }
}
```

**Options:**
- `enabled`: `true` or `false`
- `interval_minutes`: 1-60 (recommended: 5)

### Dice Roller Settings

Customize dice rolling behavior:

```json
{
  "dice_roller": {
    "history_size": 15,         // Number of rolls to keep in history
    "show_individual": true,    // Show individual die results
    "minimum_value": 1          // Minimum roll value (D&D standard)
  }
}
```

**Options:**
- `history_size`: 1-100 (recommended: 10-20) - Controls how many previous rolls are kept in history
- `show_individual`: `true` or `false` - Show individual die results like `[3, 5, 2]` (future feature)
- `minimum_value`: 0 or higher (D&D standard is 1) - The minimum value any dice roll can return, preventing negative results

### Initiative Tracker Settings

Configure combat tracker behavior:

```json
{
  "initiative_tracker": {
    "auto_sort": true,          // Automatically sort by initiative
    "show_hp": true,            // Display HP in list
    "show_ac": true,            // Display AC in list
    "highlight_active": true,   // Highlight current turn
    "round_counter": true       // Show round counter
  }
}
```

**All options:** `true` or `false`

### Display Settings

Control UI display preferences:

```json
{
  "display": {
    "show_help_hints": true,       // Show help text at bottom
    "compact_mode": false,         // Use compact panel layout
    "animate_transitions": false,  // Animate panel switching (experimental)
    "line_wrap": true,             // Wrap long text lines
    "max_line_length": 50          // Maximum characters before wrap
  }
}
```

**Options:**
- `show_help_hints`: Show keyboard shortcuts
- `compact_mode`: Reduce spacing and padding
- `animate_transitions`: Smooth panel transitions (may affect performance)
- `line_wrap`: Wrap long text to next line
- `max_line_length`: 20-200 characters (recommended: 50)

### Paths Settings

Customize file locations:

```json
{
  "paths": {
    "save_directory": "",       // Campaign save location (empty = ~/.lazydnd)
    "monster_directory": "",    // Custom monsters location (empty = ~/.lazydnd/monsters)
    "backup_enabled": true,     // Enable automatic backups
    "backup_directory": "",     // Backup location (empty = saves/.backups)
    "max_backups": 10          // Maximum backups per campaign
  }
}
```

**Options:**
- `save_directory`: Custom save location
  - Empty string: Use default `~/.lazydnd`
  - Custom path: Use absolute path, e.g., `/path/to/saves`
  - Supports `~` expansion: `~/Documents/DnD/LazyDnD`
  - Supports environment variables: `$HOME/dnd-saves`

- `monster_directory`: Custom monsters location
  - Empty string: Use default `~/.lazydnd/monsters`
  - Custom path: Absolute path to your monster collection
  - All `.json` files in this directory are loaded automatically
  - Supports `~` expansion: `~/Documents/DnD/Monsters`
  - Supports environment variables: `$HOME/custom-monsters`

- `backup_enabled`: Create automatic backups
  - `true`: Create backup before each save
  - `false`: No automatic backups

- `backup_directory`: Where to store backups
  - Empty string: Use `[save_directory]/.backups`
  - Custom path: Absolute path to backup folder

- `max_backups`: Number of backups to keep
  - Range: 1-100 (recommended: 5-20)
  - Oldest backups deleted when limit reached

**Examples:**

Store saves in Documents:
```json
{
  "paths": {
    "save_directory": "~/Documents/DnD/LazyDnD",
    "backup_enabled": true,
    "max_backups": 15
  }
}
```

Use cloud storage (Dropbox/Google Drive):
```json
{
  "paths": {
    "save_directory": "~/Dropbox/DnD/LazyDnD",
    "backup_directory": "~/Dropbox/DnD/LazyDnD/backups",
    "backup_enabled": true,
    "max_backups": 20
  }
}
```

Custom monster directory:
```json
{
  "paths": {
    "monster_directory": "~/Documents/DnD/CustomMonsters"
  }
}
```

Shared monster library (team of DMs):
```json
{
  "paths": {
    "monster_directory": "~/Dropbox/DnD-Team/SharedMonsters"
  }
}
```

Disable backups:
```json
{
  "paths": {
    "backup_enabled": false
  }
}
```

## Example Configurations

### Monokai Theme

```json
{
  "theme": {
    "primary_color": "#66D9EF",
    "border_color": "#75715E",
    "highlight_color": "#A6E22E",
    "error_color": "#F92672",
    "success_color": "#A6E22E"
  },
  "auto_save": {
    "enabled": true,
    "interval_minutes": 5
  },
  "dice_roller": {
    "history_size": 15,
    "show_individual": true,
    "minimum_value": 1,
    "critical_hit_enabled": true,
    "critical_hit_mode": "double"
  },
  "initiative_tracker": {
    "auto_sort": true,
    "show_hp": true,
    "show_ac": true,
    "highlight_active": true,
    "round_counter": true
  },
  "display": {
    "show_help_hints": true,
    "compact_mode": false,
    "line_wrap": true,
    "max_line_length": 50
  },
  "paths": {
    "save_directory": "",
    "monster_directory": "",
    "backup_enabled": true,
    "backup_directory": "",
    "max_backups": 10
  }
}
```

### Tokyo Night Theme

```json
{
  "theme": {
    "primary_color": "#7aa2f7",
    "border_color": "#414868",
    "highlight_color": "#9ece6a",
    "error_color": "#f7768e",
    "success_color": "#9ece6a"
  },
  "auto_save": {
    "enabled": true,
    "interval_minutes": 5
  },
  "dice_roller": {
    "history_size": 15,
    "show_individual": true,
    "minimum_value": 1,
    "critical_hit_enabled": true,
    "critical_hit_mode": "double"
  },
  "initiative_tracker": {
    "auto_sort": true,
    "show_hp": true,
    "show_ac": true,
    "highlight_active": true,
    "round_counter": true
  },
  "display": {
    "show_help_hints": true,
    "compact_mode": false,
    "line_wrap": true,
    "max_line_length": 50
  },
  "paths": {
    "save_directory": "",
    "monster_directory": "",
    "backup_enabled": true,
    "backup_directory": "",
    "max_backups": 10
  }
}
```

### Minimalist Setup

```json
{
  "theme": {
    "primary_color": "#FFFFFF",
    "border_color": "#666666",
    "highlight_color": "#AAAAAA",
    "error_color": "#FF6666",
    "success_color": "#66FF66"
  },
  "auto_save": {
    "enabled": false,
    "interval_minutes": 10
  },
  "dice_roller": {
    "history_size": 5,
    "show_individual": false,
    "minimum_value": 1,
    "default_advantage": false
  },
  "initiative_tracker": {
    "auto_sort": true,
    "show_hp": false,
    "show_ac": false,
    "highlight_active": true,
    "round_counter": false
  },
  "display": {
    "show_help_hints": false,
    "compact_mode": true,
    "animate_transitions": false,
    "line_wrap": true,
    "max_line_length": 40
  }
}
```

### Power User Setup

```json
{
  "theme": {
    "primary_color": "#00BFFF",
    "border_color": "#4169E1",
    "highlight_color": "#FFD700",
    "error_color": "#FF1493",
    "success_color": "#00FF7F"
  },
  "auto_save": {
    "enabled": true,
    "interval_minutes": 2
  },
  "dice_roller": {
    "history_size": 30,
    "show_individual": true,
    "minimum_value": 1,
    "default_advantage": false
  },
  "initiative_tracker": {
    "auto_sort": true,
    "show_hp": true,
    "show_ac": true,
    "highlight_active": true,
    "round_counter": true
  },
  "display": {
    "show_help_hints": true,
    "compact_mode": false,
    "animate_transitions": false,
    "line_wrap": true,
    "max_line_length": 60
  }
}
```

## Managing Configuration

### View Current Configuration

```bash
cat ~/.config/lazydnd/config.json
```

### Edit Configuration

```bash
nano ~/.config/lazydnd/config.json
# or
vim ~/.config/lazydnd/config.json
```

### Reset to Defaults

Delete the config file and restart LazyDnD:

```bash
rm ~/.config/lazydnd/config.json
lazydnd
```

### Backup Configuration

```bash
cp ~/.config/lazydnd/config.json ~/.config/lazydnd/config.backup.json
```

### Copy Configuration to Another Machine

```bash
# On source machine
cat ~/.config/lazydnd/config.json

# Copy output, then on destination machine
mkdir -p ~/.config/lazydnd
nano ~/.config/lazydnd/config.json
# Paste and save
```

## Validation

LazyDnD validates your configuration on startup. Invalid values will:
1. Show an error message
2. Fall back to defaults
3. Create a new valid config file

### Common Validation Rules

- `auto_save.interval_minutes`: Must be ≥ 1
- `dice_roller.history_size`: Must be 1-100
- `dice_roller.minimum_value`: Must be ≥ 0
- `display.max_line_length`: Must be ≥ 20
- Color values: Must be valid hex colors (e.g., `#FF0000`)

## Troubleshooting

### Config File Not Working

1. Check JSON syntax (use a JSON validator)
2. Verify file permissions: `chmod 644 ~/.config/lazydnd/config.json`
3. Check LazyDnD has write access: `ls -la ~/.config/lazydnd/`

### Changes Not Applied

1. Restart LazyDnD after editing config
2. Check for validation errors in terminal
3. Verify config file location is correct

### Colors Look Wrong

1. Ensure your terminal supports 256 colors
2. Try different hex color values
3. Use a terminal color picker to test values

### Auto-Save Not Working

1. Check `auto_save.enabled` is `true`
2. Verify interval is reasonable (1-60 minutes)
3. Check disk space and write permissions

## Tips

1. **Start with defaults**: Only change what you need
2. **Test one change at a time**: Easier to troubleshoot
3. **Keep a backup**: Before major changes
4. **Use example configs**: Copy from examples above
5. **Validate JSON**: Use online JSON validators before saving

## Future Configuration Options

Planned for future releases:
- Custom keybindings
- Font size adjustments
- Custom dice sets
- Panel layout customization
- Network/API integrations
- Campaign templates

## Support

If you encounter issues with configuration:
- Check this documentation
- Validate your JSON syntax
- Review the example config: `config.example.json`
- Open an issue: https://github.com/zingazzi/lazydnd/issues

---

**Note:** Configuration is loaded once at startup. Changes require restarting LazyDnD to take effect.
