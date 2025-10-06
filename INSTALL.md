# Installation Guide

## Quick Install (Recommended)

The easiest way to install LazyDnD is using our installation script:

```bash
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash
```

This will:
- Detect your operating system and architecture
- Download the latest release
- Install to `/usr/local/bin/lazydnd`
- Make it executable and ready to use

## Manual Installation

### Step 1: Download

Visit the [Releases page](https://github.com/zingazzi/lazydnd/releases/latest) and download the appropriate binary for your system:

- **Linux (Intel/AMD)**: `lazydnd-linux-amd64`
- **Linux (ARM)**: `lazydnd-linux-arm64`
- **macOS (Intel)**: `lazydnd-macos-amd64`
- **macOS (Apple Silicon)**: `lazydnd-macos-arm64`
- **Windows**: `lazydnd-windows-amd64.exe`

### Step 2: Install

**Linux/macOS:**

```bash
# Download (replace URL with your actual release URL)
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-amd64

# Make executable
chmod +x lazydnd

# Move to system path (requires sudo)
sudo mv lazydnd /usr/local/bin/

# Or install to user directory (no sudo required)
mkdir -p ~/.local/bin
mv lazydnd ~/.local/bin/
# Add ~/.local/bin to your PATH if not already there
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Windows:**

1. Download `lazydnd-windows-amd64.exe`
2. Rename to `lazydnd.exe`
3. Move to a directory in your PATH, or:
   - Create a folder like `C:\Program Files\LazyDnD`
   - Move `lazydnd.exe` there
   - Add that folder to your system PATH

### Step 3: Verify

```bash
lazydnd
```

## Platform-Specific Instructions

### Linux

**Debian/Ubuntu:**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-amd64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**Arch Linux:**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-amd64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**Raspberry Pi (ARM):**
```bash
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-linux-arm64
chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

### macOS

**Using Homebrew (coming soon):**
```bash
brew install zingazzi/tap/lazydnd
```

**Manual Installation:**
```bash
# For Apple Silicon (M1/M2/M3)
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-macos-arm64

# For Intel Macs
curl -L -o lazydnd https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-macos-amd64

chmod +x lazydnd
sudo mv lazydnd /usr/local/bin/
```

**Note for macOS users:** On first run, macOS may block the app because it's not signed. To allow it:
1. Go to System Preferences → Security & Privacy
2. Click "Allow Anyway" next to the blocked app message
3. Run `lazydnd` again

Or use this command to bypass Gatekeeper:
```bash
xattr -d com.apple.quarantine /usr/local/bin/lazydnd
```

### Windows

**PowerShell:**
```powershell
# Download
Invoke-WebRequest -Uri "https://github.com/zingazzi/lazydnd/releases/latest/download/lazydnd-windows-amd64.exe" -OutFile "lazydnd.exe"

# Move to a directory in PATH (example)
Move-Item lazydnd.exe C:\Windows\System32\
```

**Windows Terminal / WSL:**
If you're using WSL (Windows Subsystem for Linux), follow the Linux instructions instead.

## Build from Source

If you prefer to build from source:

### Prerequisites
- Go 1.21 or higher
- Git

### Steps

```bash
# Clone the repository
git clone https://github.com/zingazzi/lazydnd.git
cd lazydnd

# Build
go build -o lazydnd

# Run
./lazydnd

# Optional: Install to system
sudo mv lazydnd /usr/local/bin/
```

### Cross-Platform Build

To build for all platforms:

```bash
./build.sh
```

This creates binaries in the `build/` directory for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)

## Updating

To update to the latest version, simply run the installation command again:

```bash
curl -sSL https://raw.githubusercontent.com/zingazzi/lazydnd/main/install.sh | bash
```

Or manually download the latest release and replace your existing binary.

## Uninstalling

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/lazydnd

# If installed to ~/.local/bin
rm ~/.local/bin/lazydnd
```

## Troubleshooting

### "Command not found"
- Make sure the binary is in your PATH
- Try running with full path: `/usr/local/bin/lazydnd`
- Check if `~/.local/bin` is in your PATH: `echo $PATH`

### "Permission denied"
- Make sure the binary is executable: `chmod +x lazydnd`
- You may need sudo to install to `/usr/local/bin`

### macOS "Cannot be opened because the developer cannot be verified"
- Run: `xattr -d com.apple.quarantine /usr/local/bin/lazydnd`
- Or go to System Preferences → Security & Privacy and click "Allow Anyway"

### Terminal colors not working
- Make sure your terminal supports 256 colors
- Try a modern terminal like iTerm2 (macOS), Windows Terminal (Windows), or GNOME Terminal (Linux)

## Support

If you encounter any issues, please open an issue on [GitHub](https://github.com/zingazzi/lazydnd/issues).
