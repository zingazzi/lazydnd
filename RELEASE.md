# Release Process

This document explains how to create releases for LazyDnD.

## Automated Releases (Recommended)

The project uses GitHub Actions to automatically build and release binaries when you push a tag.

### Creating a Release

1. **Update the CHANGELOG.md** with your changes

2. **Commit your changes:**
   ```bash
   git add .
   git commit -m "Release v0.1.3"
   ```

3. **Create and push a tag:**
   ```bash
   git tag v0.1.3
   git push origin v0.1.3
   ```

4. **GitHub Actions will automatically:**
   - Build binaries for all platforms (Linux, macOS, Windows)
   - Create checksums for verification
   - Create a GitHub Release with all binaries attached
   - Generate installation instructions in the release notes

### What Gets Built

The automated release builds:
- `lazydnd-linux-amd64` - Linux (Intel/AMD 64-bit)
- `lazydnd-linux-arm64` - Linux (ARM 64-bit, e.g., Raspberry Pi)
- `lazydnd-macos-amd64` - macOS (Intel)
- `lazydnd-macos-arm64` - macOS (Apple Silicon M1/M2/M3)
- `lazydnd-windows-amd64.exe` - Windows (64-bit)
- `checksums.txt` - SHA256 checksums for all binaries

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version (v1.0.0): Incompatible API changes
- **MINOR** version (v0.1.0): New functionality, backwards compatible
- **PATCH** version (v0.0.1): Bug fixes, backwards compatible

Examples:
- `v0.1.0` - Initial release
- `v0.1.1` - Bug fix
- `v0.2.0` - New feature
- `v1.0.0` - Stable release

## Manual Release (Alternative)

If you prefer to create releases manually:

### 1. Build Binaries

```bash
./build.sh
```

This creates binaries in the `build/` directory.

### 2. Create Checksums

```bash
cd build
sha256sum * > checksums.txt
cd ..
```

### 3. Create GitHub Release

1. Go to your repository on GitHub
2. Click "Releases" → "Draft a new release"
3. Create a new tag (e.g., `v0.1.3`)
4. Set the release title (e.g., "LazyDnD v0.1.3")
5. Add release notes from CHANGELOG.md
6. Upload all files from `build/` directory
7. Publish release

## Testing a Release

Before publishing, test the binaries:

```bash
# Test local build
./build/lazydnd-linux-amd64

# Test downloaded binary
curl -L -o lazydnd-test https://github.com/zingazzi/lazydnd/releases/download/v0.1.3/lazydnd-linux-amd64
chmod +x lazydnd-test
./lazydnd-test
```

## Updating Installation Script

If you change the repository name or structure, update these files:
- `.github/workflows/release.yml` - Update repository references
- `install.sh` - Update `REPO` variable
- `README.md` - Update installation URLs
- `INSTALL.md` - Update installation URLs

Replace `zingazzi` with your actual GitHub username in all files.

## Release Checklist

Before creating a release:

- [ ] Update CHANGELOG.md with all changes
- [ ] Update version in documentation if needed
- [ ] Test locally with `./build.sh`
- [ ] Test all major features work
- [ ] Commit all changes
- [ ] Create and push tag
- [ ] Verify GitHub Actions workflow succeeds
- [ ] Test installation from release
- [ ] Announce release (optional)

## Troubleshooting

### GitHub Actions fails
- Check the Actions tab for error logs
- Ensure `go.mod` is up to date
- Verify all dependencies are available

### Binaries don't work
- Test locally first with `./build.sh`
- Check for platform-specific issues
- Verify Go version compatibility

### Release not appearing
- Make sure you pushed the tag: `git push origin v0.1.3`
- Check GitHub Actions status
- Verify workflow file is in `.github/workflows/`

## Example Release Flow

```bash
# 1. Make your changes
git add .
git commit -m "Add new dice rolling features"

# 2. Update CHANGELOG
vim CHANGELOG.md
git add CHANGELOG.md
git commit -m "Update CHANGELOG for v0.1.3"

# 3. Create and push tag
git tag v0.1.3
git push origin main
git push origin v0.1.3

# 4. Wait for GitHub Actions to complete
# 5. Check the Releases page on GitHub
# 6. Test the installation
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash
```

## Post-Release

After a successful release:
- Update any documentation that references version numbers
- Announce on social media, Discord, etc. (optional)
- Monitor issues for bug reports
- Plan next release based on feedback
