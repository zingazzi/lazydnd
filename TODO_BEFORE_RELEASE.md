# TODO Before First Release

## ✅ Setup Complete

The automated release system is configured! Here's what you need to do before creating your first release:

## 🔴 CRITICAL - Must Do First

### 1. Replace zingazzi with your GitHub username

Run these commands (replace `your-actual-username` with your GitHub username):

```bash
cd /home/zingo/PLAYGROUND/lazydnd

# Update all files at once
find . -type f \( -name "*.yml" -o -name "*.sh" -o -name "*.md" \) -exec sed -i 's/zingazzi/your-actual-username/g' {} +

# Verify the changes
grep -r "zingazzi" . --include="*.yml" --include="*.sh" --include="*.md"
# This should return no results if successful
```

Or manually edit these files:
- [ ] `.github/workflows/release.yml`
- [ ] `install.sh` (line 10)
- [ ] `README.md` (all installation URLs)
- [ ] `INSTALL.md` (all installation URLs)  
- [ ] `RELEASE.md` (example commands)
- [ ] `SETUP_SUMMARY.md` (example commands)

### 2. Test Local Build

```bash
./build.sh
```

Verify all binaries are created in `build/` directory.

### 3. Test the Application

```bash
./lazydnd
```

Make sure everything works as expected.

## 📝 Recommended Before Release

### 4. Update CHANGELOG.md

Make sure `CHANGELOG.md` reflects all changes in v0.1.3:
- [ ] Dice roller improvements documented
- [ ] Shift+Tab navigation mentioned
- [ ] All new features listed

### 5. Verify Documentation

- [ ] README.md has correct installation instructions
- [ ] Examples in README match actual functionality
- [ ] Screenshots are up to date (if any)

### 6. Clean Up

```bash
# Remove any test files
rm -f lazydnd-test

# Ensure build directory is in .gitignore
echo "build/" >> .gitignore
```

## 🚀 Ready to Release

Once everything above is done:

### 1. Commit Everything

```bash
git add .
git commit -m "Setup automated releases and update documentation"
git push origin main
```

### 2. Create First Release

```bash
# Create and push tag
git tag v0.1.3
git push origin v0.1.3
```

### 3. Monitor Release

1. Go to: `https://github.com/zingazzi/lazydnd/actions`
2. Watch the "Release" workflow run
3. Should complete in 2-3 minutes
4. Check: `https://github.com/zingazzi/lazydnd/releases`

### 4. Test Installation

```bash
# Test the installer
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash

# Or test manual installation
curl -L -o lazydnd-test https://github.com/zingazzi/lazydnd/releases/download/v0.1.3/lazydnd-linux-amd64
chmod +x lazydnd-test
./lazydnd-test
```

## 📋 Post-Release

After successful release:

- [ ] Test installation on different platforms (if possible)
- [ ] Update any external documentation
- [ ] Announce release (Discord, Reddit, etc.)
- [ ] Monitor GitHub Issues for bug reports
- [ ] Plan next release features

## 🐛 If Something Goes Wrong

### Release workflow fails
1. Check Actions tab for error logs
2. Fix the issue
3. Delete the tag: `git tag -d v0.1.3 && git push origin :refs/tags/v0.1.3`
4. Fix and try again

### Binaries don't work
1. Test locally first with `./build.sh`
2. Check Go version (needs 1.21+)
3. Review build logs in Actions

### Installation script fails
1. Verify zingazzi was replaced
2. Check if release exists
3. Test URL manually in browser

## 📚 Reference

- **Semantic Versioning**: https://semver.org/
- **GitHub Actions**: https://docs.github.com/en/actions
- **Go Cross Compilation**: https://go.dev/doc/install/source#environment

## 🎉 You're Ready!

Once you've completed the checklist above, you're ready to create your first automated release!

Good luck! 🎲
