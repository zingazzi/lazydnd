# CHANGELOG

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
