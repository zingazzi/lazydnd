# CHANGELOG

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
