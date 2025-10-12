# CHANGELOG

## v2.9.0
- **Fix**
  - In history now you can reroll macro
  - Note improvment

- **UI Panel Enhancements**
  - Clearer input indicators and controls
  - Consistent styling across search modes
	- Clear interface

## 2.8.0
- **Dice Macros**
  - Save common rolls as macros: `fireball=8d6`, `goblin_attack=1d20+4`
  - Execute macros by name: type `fireball` to roll 8d6
  - **35+ preset macros for popular D&D spells built-in**:
    - Cantrips: `eldritch_blast`, `fire_bolt`, `toll_the_dead`, etc.
    - Damage spells: `fireball`, `lightning_bolt`, `cone_of_cold`, etc.
    - Healing: `cure_wounds`, `healing_word`
    - High-level: `disintegrate`, `chain_lightning`, `finger_of_death`
    - Attack combos: `sneak_attack_1d6` through `sneak_attack_4d6`
    - Divine smite: `divine_smite_1d8` through `divine_smite_3d8`
  - Macros persist with campaign saves
  - User macros override presets with same name
  - Simple syntax: `name=formula`
  - Supports all dice notation (advantage, disadvantage, modifiers)
  - Perfect for frequently used spells and attacks

- **Skill Check Shortcuts**
  - Type skill name to roll 1d20 for selected character
  - Supports all D&D 5e skills: stealth, perception, athletics, etc.
  - Add modifiers: `stealth+5` rolls 1d20+5
  - Automatically shows character name with result
  - Requires selected entry in Initiative Tracker

- **Group Initiative Rolls**
  - Command: `group`, `group init`, or `group initiative`
  - Rolls initiative for all monsters at once
  - Uses monster DEX modifiers automatically
  - Re-sorts initiative list after rolling
  - Saves time when starting combat with multiple monsters

- **Notes Panel (5th Panel)**
  - New dedicated panel for campaign notes (press `5` or `F5`)
  - Session notes, plot points, NPC names, important details
  - Full markdown-style formatting support:
    - `# Heading`, `## Subheading`, `### Smaller` - styled headings
    - `- item` or `* item` - bullet points
    - `**bold**`, `*italic*` - text emphasis (simplified for terminal)
  - Edit mode (`e` key) - multi-line text editing with Enter for new lines
  - Search functionality (`f` key) - find text within notes with highlighting
  - Auto-saves with campaign - notes persist through save/load
  - Empty state with helpful tips and markdown guide
  - Context-sensitive help and inline instructions

- **Dynamic Panel Layout**
  - Revolutionary adaptive layout system - panels resize based on focus
  - Top row: Dice Roller (1) and Initiative Tracker (2) - resize when active
  - Bottom row: Spells (3), Monsters (4), Notes (5) - resize when active
  - Active panel automatically expands for better visibility
  - Inactive panels shrink to save space
  - Smooth visual transitions when switching panels
  - Notes panel starts very small, expands to 40% when active
  - Tab/Shift+Tab now cycles through all 5 panels
  - Optimized space usage across all screen sizes

## 2.7.0
- **UI Error Messages**
  - Errors now display in a red banner above the status bar
  - Shows clear error messages for save/load failures
  - Automatically clears after 5 seconds or on any key press
  - No more silent failures - you'll always know when something goes wrong
  - Examples of errors shown:
    - "Failed to save campaign: [reason]"
    - "Failed to load campaign: [reason]"
    - "Failed to load campaign list: [reason]"

- **Max HP Editing for Monsters**
  - Press `H` (Shift+H) in edit mode to change a monster's maximum HP
  - Set absolute HP values (e.g., "150" sets max HP to 150)
  - Current HP automatically capped if it exceeds new maximum
  - Perfect for adjusting monster difficulty on the fly
  - Minimum max HP value of 1 enforced
  - All changes persist through save/load

- **Search Panel Improvements**
  - Added scrolling indicators for long suggestion lists
  - Shows "X more above/below" with arrow indicators
  - Maximum 5 visible suggestions at once
  - Selected item stays centered in scroll view
  - Smooth scrolling as you navigate suggestions
  - Clear visual separation with horizontal lines
  - Consistent styling and spacing
  - Better handling of long suggestion lists

- **Search Input Refinements**
  - Clearer input prompts based on context
  - Shows "Search:", "CR:", or "Level:" appropriately
  - Active cursor indicator with block character
  - Empty state handling shows minimal space
  - Control hints shown below input ("↑↓ Enter Esc")
  - Improved visual hierarchy and spacing
  - Better feedback for active/inactive states

- **Panel Content Organization**
  - Extracted search panel content logic to separate file
  - Created SearchContentConfig struct for configuration
  - Generic RenderSearchContent function for reuse
  - Consistent styling across search panels
  - Better separation of concerns
  - Improved code maintainability
  - Reduced code duplication
  - Cleaner panel rendering code

## 2.6.0
- **CR Filter for Monsters (Autocomplete Style)**
  - Browse monsters by Challenge Rating with live autocomplete
  - Press `f` in Monsters panel to open CR filter
  - Type CR value and monster list appears instantly
  - Supports multiple filter formats:
    - Exact: `5` → shows all CR 5 monsters
    - Range: `0-5` → shows CR 0 through 5
    - Minimum: `10+` → shows CR 10 and above
  - Handles fractional CRs: `1/4`, `1/2`, `0.5`
  - Navigate with `↑↓` arrows, select with `Enter`
  - Works just like the name search autocomplete
  - Parses CR from full strings like "5 (1,800 XP)"
  - Perfect for quickly finding level-appropriate encounters!

- **Fuzzy Search for Monsters and Spells**
  - Implemented intelligent fuzzy search using `github.com/sahilm/fuzzy`
  - Handles typos and partial matches (e.g., "frbl" finds "Fireball")
  - Results automatically sorted by match quality
  - Works in both Monsters and Spells panels
  - More forgiving search - no need for exact spelling
  - Examples:
    - "drag red" → finds "Adult Red Dragon"
    - "mgc msl" → finds "Magic Missile"
    - "phn wp" → finds "Phantom Weapon"

- **Rendering Fixes**
  - Fixed text wrapping in monster panel to prevent text overflow (60 → 35 chars)
  - Fixed separator lines to match panel width (40 → 35 chars)
  - Completely refactored search rendering to prevent overlap
  - Search mode and details mode are now mutually exclusive
  - Early return pattern ensures clean UI without any overlap
  - Prevents graphic glitches and text overlapping in all panels
  - All text now properly fits within panel boundaries

- **Debug Logging**
  - Added `--debug` flag to enable debug logging
  - Logs saved to `~/.config/lazydnd/debug.log`
  - Real-time logging of key events, state changes, and operations
  - Logs include: key presses, multi-target operations, condition management, spell tracking
  - Useful for troubleshooting and development
  - Example: `lazydnd --debug`

- **Zocchi's Dice Support**
  - Added support for Zocchi's dice: d3, d5, d7, d14, d16, d24, d30
  - All dice types work with standard notation (e.g., `2d7+3`, `1d30`, `3d5-1`)
  - Dice roller now supports 14 different dice types total

## 2.5.0
- **Multi-Target Damage/Healing System**
  - Select multiple targets in Initiative Tracker for simultaneous damage/healing
  - Press 't' to enter multi-target mode
  - Use Space to select/deselect targets (checkboxes appear)
  - Press Enter to apply damage or healing to all selected targets
  - Optional save mode: targets can succeed/fail saves for half damage
  - Toggle between damage and healing modes with 'h'
  - Visual feedback with checkboxes and counters
  - Works great for area spells like Fireball, mass healing, etc.

- **Spell Tracking System**
  - Cast spells and track their durations during combat
  - Active spell list showing all ongoing spell effects
  - Automatic duration countdown integrated with initiative tracker
  - Manual spell deletion capability
  - Duration parsing for rounds, minutes, hours, and days
  - Concentration spells marked with (C)
  - Active spells saved with campaign state
  - Human-readable time remaining display
  - Instantaneous spells not tracked
  - Press 'c' to cast spell, 'v' to view active spells, 'd' to delete

## 2.4.1
- **Version Display**
  - Added version number (v2.4.1) to status bar next to LazyDnD name
  - Added `--version` command-line flag to display version without launching app
  - Version centrally managed in `ui/text_content.go`

- **Configuration System Improvements**
  - Theme colors are now fully dynamic and applied throughout the UI
  - Active/inactive panel borders use configured colors
  - Status bar uses primary color from config
  - Help popup uses primary and highlight colors
  - All text highlights respect theme configuration
  - Auto-save now uses configured interval and enabled flag
  - Dice roller now respects `history_size` and `minimum_value` config
  - Removed deprecated `AutoSaveEnabled` model field
  - Removed `default_advantage` config option (not applicable)

- **Code Quality Improvements**
  - Reduced code duplication in UI handlers
  - Added `isInInputMode()` helper to eliminate 11 duplicate checks
  - Added `addToHistory()` helper for consistent dice history management
  - Updated all tests to use new config-based RollDice signature
  - All tests passing ✅

  **Add Custom Monster**
  - Now you can add your custom monster in a json file

  - Update release script

## 2.3.0
- Add quality code check (duplication and lint)
- Add Tests
- Saving throw and ability check for monsters

## 2.2.0
- Add turn counter
- Minor navigation fix

## 2.1.0
- Reroll a specific dice
- Minor fix
- Show attacks x round
- Add initiative turn tracker

## 2.0.0
- **Major UI Overhaul**
  - Redesigned panel layout system for better space utilization
  - Improved panel resizing and scrolling behavior
  - Enhanced visual hierarchy and readability
  - Added status bar with campaign info and auto-save status
  - Cleaner borders and spacing between panels

- **Campaign Management**
  - New campaign save/load system with auto-save
  - Campaign state persistence (initiative, monster links, etc.)
  - Campaign renaming support
  - Auto-save every 5 minutes with status indicator
  - Save files stored in ~/.lazydnd/

- **Initiative Tracker Enhancements**
  - Monster linking from Monster panel with full data
  - Action integration - view and roll monster actions
  - HP tracking with +/- modification
  - Entry duplication with auto-numbering
  - Edit mode improvements
  - Initiative sorting

- **Monster Panel Improvements**
  - 8750+ D&D 5e monsters database
  - Structured action parsing
  - Complete stat blocks
  - Quick add to initiative with 'a' key
  - Real-time search with suggestions

- **Spell Panel Updates**
  - Complete D&D 5e spell database
  - Real-time search with autocomplete
  - Detailed spell information display
  - Class availability indicators
  - Ritual and concentration markers

- **Dice Roller Updates**
  - Support for advantage/disadvantage rolls
  - Complex dice expressions
  - Roll history display
  - Quick reroll functionality
  - D&D minimum value rule enforcement

- **Performance & Stability**
  - Optimized panel rendering
  - Reduced memory usage
  - Improved error handling
  - Better state management
  - Faster search response times



## v1.1.0
- **Code Refactoring & Organization**
  - Refactored layout.go from 287 lines to focused functions
  - Extracted panel dimension calculation
  - Separated scrolling logic
  - Created provider pattern for panel content and help text
  - Added struct types for better data organization
  - Organized functions into clear sections
  - Refactored navigation.go from 824-line switch to KeyHandler map pattern
  - Organized handlers by category (quit, navigation, function keys, etc.)
  - Each handler has clear naming and focused purpose
  - Improved code maintainability and readability
  - Better adherence to single responsibility principle

- **Documentation & Standards**
  - Added file path/name comments
  - Improved commenting standards
  - Comments focus on purpose over effect
  - Added descriptive section headers
  - Clearer function documentation


## v1.0.1
- **Docker Support**
  - Complete Docker containerization
  - Multi-stage build for optimized image size (~15MB)
  - Docker Compose configuration
  - Automated Docker builds via GitHub Actions
  - Multi-platform support (amd64, arm64)
  - Published to GitHub Container Registry (GHCR)
  - Helper scripts: `docker-build.sh`, `docker-run.sh`
  - Comprehensive Docker documentation

- **Installation & Release System**
  - Automated GitHub Actions workflow for releases
  - Cross-platform builds (Linux, macOS, Windows)
  - One-line installer script
  - Comprehensive installation documentation

- **Navigation Improvements**
  - Added Shift+Tab for backward panel navigation
  - Tab cycles forward, Shift+Tab cycles backward

- **Dice Roller Enhancements**
  - Support for multiple dice expressions: `2d8+3d6`, `1d6-1d4`
  - Support for subtraction: `1d8-2`, `2d6-1d4+3`
  - Complex expressions: `1d6+3+2d8-5`
  - Comma-separated rolls: `1d8+3, 3d6-1`
  - Minimum value enforcement (results never go below 1, following D&D rules)
  - Improved result display: Total shown first, then detailed breakdown
  - Automatic text wrapping for long results
  - Results show dice expressions with rolled values
  - Support for advantage/disadvantage with complex expressions

## v0.1.2
- First version of LazyDnd
- Update Snaphsot
