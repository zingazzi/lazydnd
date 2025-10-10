# Scripts Documentation

This directory contains utility scripts for version management and release automation.

## Available Scripts

### `extract_version.sh`
**Purpose**: Extract the current version from Git tags.

**Usage**:
```bash
./scripts/extract_version.sh
```

**Output**:
- Current Git tag version (e.g., `v2.4.0`) if on a tagged commit
- `dev` if not on a tagged commit

**Example**:
```bash
$ git tag -a v2.4.0 -m "Release v2.4.0"
$ ./scripts/extract_version.sh
v2.4.0

$ git checkout main
$ ./scripts/extract_version.sh
dev
```

---

### `extract_changelog.sh`
**Purpose**: Extract changelog entries for a specific version from `CHANGELOG.md`.

**Usage**:
```bash
./scripts/extract_changelog.sh <version>
```

**Parameters**:
- `<version>`: Version number (with or without `v` prefix)

**Output**: All changelog entries between the version header and the next version header.

**Example**:
```bash
$ ./scripts/extract_changelog.sh v2.4.0
- **Version Display**
  - Added version number (v2.4.0) to status bar next to LazyDnD name
  - Added `--version` command-line flag to display version without launching app
...

$ ./scripts/extract_changelog.sh 2.4.0  # Works without 'v' prefix too
- **Version Display**
...
```

---

### `update_version.sh`
**Purpose**: Update the version in `ui/text_content.go` based on the current Git tag.

**Usage**:
```bash
./scripts/update_version.sh
```

**What it does**:
1. Extracts current version from Git tags
2. Updates `AppVersion` constant in `ui/text_content.go`
3. Creates a backup file (`.bak`) and removes it after successful update

**Example**:
```bash
$ git tag -a v2.5.0 -m "Release v2.5.0"
$ ./scripts/update_version.sh
Updating version to: v2.5.0
✅ Version updated to v2.5.0 in ui/text_content.go
```

**Note**: This script is automatically run by GitHub Actions during the release process. Manual execution is only needed for local testing.

---

## Workflow Integration

These scripts are used by the GitHub Actions workflow (`.github/workflows/release.yml`) during the automated release process:

1. **Extract Version**: Gets version from the pushed Git tag
2. **Update Version**: Updates `ui/text_content.go` with the tag version
3. **Extract Changelog**: Pulls relevant changelog entries for the release notes
4. **Build**: Compiles binaries with the correct version embedded
5. **Release**: Creates GitHub release with changelog as description

## Development Workflow

### Local Testing
```bash
# Test version extraction
./scripts/extract_version.sh

# Test changelog extraction
./scripts/extract_changelog.sh v2.4.0

# Update version in code (for testing)
./scripts/update_version.sh
go build -o lazydnd
./lazydnd --version
```

### Creating a Release
See [VERSIONING.md](../VERSIONING.md) for the complete release process.

**Quick steps**:
```bash
# 1. Update CHANGELOG.md
# 2. Commit changes
git add CHANGELOG.md
git commit -m "Prepare release v2.5.0"
git push

# 3. Create and push tag
git tag -a v2.5.0 -m "Release v2.5.0"
git push origin v2.5.0

# GitHub Actions will automatically:
# - Run update_version.sh to update the code
# - Run extract_changelog.sh to get release notes
# - Build all platform binaries
# - Create GitHub release
```

## Requirements

All scripts require:
- **Bash**: `/bin/bash`
- **Git**: For version extraction and tagging
- **sed**: For text replacement (standard on Linux/macOS)
- **awk**: For changelog extraction (standard on Linux/macOS)

## Troubleshooting

### "Permission denied"
Make scripts executable:
```bash
chmod +x scripts/*.sh
```

### Wrong version extracted
Ensure you're on a tagged commit:
```bash
git describe --tags --exact-match
```

### Changelog extraction returns nothing
- Check that version exists in `CHANGELOG.md`
- Ensure version format matches: `## X.Y.Z` (no `v` prefix in changelog)
- Test manually: `./scripts/extract_changelog.sh vX.Y.Z`

### sed error on macOS
The scripts use GNU sed syntax. On macOS, ensure compatibility or install GNU sed:
```bash
brew install gnu-sed
# Then alias sed to gsed in your shell
```

## Contributing

When modifying these scripts:
1. Test locally before pushing
2. Ensure compatibility with GitHub Actions Ubuntu runners
3. Update this README with any new parameters or behavior
4. Add examples for new functionality

