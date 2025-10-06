# CHANGELOG

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
