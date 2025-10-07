# Docker Setup Summary

## ✅ Complete Docker Configuration

Your LazyDnD project now has full Docker support! Here's what has been created:

### 📁 Files Created

1. **`Dockerfile`** - Multi-stage build configuration
   - Builder stage: Compiles Go application
   - Runtime stage: Minimal Alpine Linux image (~15MB)
   - Non-root user for security
   - Optimized with `-ldflags "-s -w"` for smaller binary

2. **`.dockerignore`** - Excludes unnecessary files from Docker context
   - Reduces build time
   - Smaller image size
   - Excludes .git, docs, build artifacts

3. **`docker-compose.yml`** - Easy orchestration
   - Simple `docker-compose up` to run
   - TTY and STDIN configured for TUI
   - Ready for volume mounting

4. **`docker-build.sh`** - Build script
   - Automatic version tagging from git
   - Creates both versioned and latest tags
   - Color output and status messages

5. **`docker-run.sh`** - Run script
   - Auto-builds if image doesn't exist
   - Interactive mode with proper TTY
   - Clean exit with --rm flag

6. **`.github/workflows/docker.yml`** - Automated CI/CD
   - Builds on every tag push
   - Multi-platform (amd64, arm64)
   - Pushes to GitHub Container Registry
   - Build caching for faster builds

7. **`DOCKER.md`** - Comprehensive documentation
   - All Docker commands explained
   - Troubleshooting section
   - Advanced usage examples
   - Publishing to registries

8. **`DOCKER_QUICKSTART.md`** - Quick reference
   - Most common commands
   - Copy-paste ready
   - Troubleshooting tips

### 🚀 Usage

#### For Users

**Easiest - Pull from GHCR:**
```bash
docker run -it --rm ghcr.io/zingazzi/lazydnd:latest
```

**Build locally:**
```bash
git clone https://github.com/zingazzi/lazydnd
cd lazydnd
./docker-build.sh
./docker-run.sh
```

**Using Docker Compose:**
```bash
docker-compose up
```

#### For Developers

**Build:**
```bash
docker build -t lazydnd:dev .
```

**Run:**
```bash
docker run -it --rm lazydnd:dev
```

**Development with live reload:**
```bash
docker run -it --rm \
    -v $(pwd):/app \
    -w /app \
    golang:1.21-alpine \
    go run .
```

### 🔄 Automated Workflow

When you push a tag (e.g., `v1.0.1`):

1. ✅ GitHub Actions triggers
2. ✅ Builds multi-platform Docker image
3. ✅ Pushes to GitHub Container Registry
4. ✅ Tags: `latest`, `v1.0.1`, `v1.0`, `v1`
5. ✅ Available at: `ghcr.io/zingazzi/lazydnd:latest`

### 📊 Image Details

- **Base**: Alpine Linux (minimal)
- **Size**: ~15MB (vs ~500MB with full Go image)
- **Platforms**: linux/amd64, linux/arm64
- **User**: Non-root (lazydnd:lazydnd)
- **Security**: Minimal attack surface

### 🎯 Key Features

1. **Multi-stage build**: Small final image
2. **Non-root user**: Enhanced security
3. **Multi-platform**: Works on ARM and x86
4. **Automated builds**: CI/CD integrated
5. **Helper scripts**: Easy to use
6. **Comprehensive docs**: DOCKER.md and DOCKER_QUICKSTART.md

### 📝 Documentation Updates

- ✅ README.md updated with Docker section
- ✅ CHANGELOG.md updated with Docker features
- ✅ .gitignore updated to exclude Docker temp files

### 🔧 Before Publishing

**Important:** Make sure to update the repository name in:
- `.github/workflows/docker.yml` (line 39: `images:`)
- `DOCKER.md`

### 🧪 Testing

**Local test:**
```bash
# Build
./docker-build.sh

# Run
./docker-run.sh

# Test it works
# (You should see the LazyDnD TUI interface)
```

**Test from GHCR (after pushing):**
```bash
docker pull ghcr.io/zingazzi/lazydnd:latest
docker run -it --rm ghcr.io/zingazzi/lazydnd:latest
```

### 🚀 Publishing to GHCR

**Manual push:**
```bash
# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u zingazzi --password-stdin

# Tag for GHCR
docker tag lazydnd:latest ghcr.io/zingazzi/lazydnd:latest
docker tag lazydnd:latest ghcr.io/zingazzi/lazydnd:v1.0.1

# Push
docker push ghcr.io/zingazzi/lazydnd:latest
docker push ghcr.io/zingazzi/lazydnd:v1.0.1
```

**Or automatic (recommended):**
Just push a tag and GitHub Actions handles everything:
```bash
git tag v1.0.1
git push origin v1.0.1
```

### 🎉 What Users Get

Users can now run LazyDnD with a single command, no installation needed:

```bash
docker run -it --rm ghcr.io/zingazzi/lazydnd:latest
```

This works on:
- ✅ Linux (any distro)
- ✅ macOS (Intel and Apple Silicon)
- ✅ Windows with Docker Desktop
- ✅ Raspberry Pi (ARM)
- ✅ Cloud servers
- ✅ CI/CD pipelines

### 📚 Further Reading

- [DOCKER.md](DOCKER.md) - Complete Docker documentation
- [Official Docker Docs](https://docs.docker.com/)
- [GitHub Container Registry Docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)

🎲 **Your D&D tool is now fully containerized!** 🐳
