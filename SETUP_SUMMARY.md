# GitHub Release Setup - Quick Summary

## What's Been Set Up

Your LazyDnD project now has automated releases configured! Here's what you have:

### 📁 Files Created/Updated

1. **`.github/workflows/release.yml`** - GitHub Actions workflow for automated releases
2. **`install.sh`** - One-line installer script for users
3. **`INSTALL.md`** - Comprehensive installation guide
4. **`RELEASE.md`** - Release process documentation
5. **`README.md`** - Updated with installation instructions
6. **`build.sh`** - Updated to include Windows builds

### 🚀 How It Works

When you push a tag (e.g., `v0.1.3`), GitHub Actions automatically:
1. Builds binaries for Linux, macOS, and Windows
2. Creates SHA256 checksums
3. Creates a GitHub Release
4. Attaches all binaries to the release
5. Generates installation instructions

### 📝 Before You Start

**Important:** Replace `zingazzi` in these files with your actual GitHub username:
- `.github/workflows/release.yml`
- `install.sh` (line 10: `REPO="zingazzi/lazydnd"`)
- `README.md` (all installation URLs)
- `INSTALL.md` (all installation URLs)
- `RELEASE.md` (example commands)

### 🎯 Quick Start

1. **Update your GitHub username:**
   ```bash
   # Replace zingazzi with your actual username
   sed -i 's/zingazzi/your-actual-username/g' .github/workflows/release.yml
   sed -i 's/zingazzi/your-actual-username/g' install.sh
   sed -i 's/zingazzi/your-actual-username/g' README.md
   sed -i 's/zingazzi/your-actual-username/g' INSTALL.md
   sed -i 's/zingazzi/your-actual-username/g' RELEASE.md
   ```

2. **Commit everything:**
   ```bash
   git add .
   git commit -m "Add automated release workflow"
   git push origin main
   ```

3. **Create your first release:**
   ```bash
   git tag v0.1.3
   git push origin v0.1.3
   ```

4. **Watch it build:**
   - Go to your GitHub repository
   - Click "Actions" tab
   - Watch the release workflow run
   - Check "Releases" when done

### 👥 User Installation

Once released, users can install with:

```bash
# One-line install
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash

# Or manual download
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-amd64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

### 🔧 What Gets Built

Each release includes:
- `lazydnd-linux-amd64` (Linux Intel/AMD)
- `lazydnd-linux-arm64` (Linux ARM/Raspberry Pi)
- `lazydnd-macos-amd64` (macOS Intel)
- `lazydnd-macos-arm64` (macOS Apple Silicon)
- `lazydnd-windows-amd64.exe` (Windows)
- `checksums.txt` (SHA256 verification)

### 📋 Release Checklist

1. Update `CHANGELOG.md`
2. Commit changes
3. Create tag: `git tag v0.1.3`
4. Push tag: `git push origin v0.1.3`
5. Wait for GitHub Actions
6. Test installation
7. Announce release

### 🐛 Troubleshooting

**Actions not running?**
- Make sure `.github/workflows/release.yml` is committed
- Check Actions tab for errors
- Verify tag was pushed: `git push origin v0.1.3`

**Binaries not working?**
- Test locally first: `./build.sh`
- Check Go version (needs 1.21+)
- Review Actions logs for build errors

**Installation fails?**
- Update zingazzi in all files
- Verify release exists on GitHub
- Check file permissions

### 📚 Documentation

- **README.md** - User-facing documentation with installation
- **INSTALL.md** - Detailed installation guide for all platforms
- **RELEASE.md** - Developer guide for creating releases
- **CHANGELOG.md** - Version history and changes

### 🎉 Next Steps

1. Replace `zingazzi` in all files
2. Commit and push to GitHub
3. Create your first tag and release
4. Share with your D&D group!

### 💡 Tips

- Use semantic versioning (v0.1.0, v0.2.0, v1.0.0)
- Always update CHANGELOG.md before releasing
- Test locally before creating a tag
- Keep installation instructions up to date
- Monitor GitHub Issues after releases

---

**Need help?** Check `RELEASE.md` for detailed release instructions or open an issue on GitHub.
