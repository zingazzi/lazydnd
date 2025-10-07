#!/bin/bash
# install.sh - Installation script for LazyDnD

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# GitHub repository
REPO="zingazzi/lazydnd"
INSTALL_DIR="/usr/local/bin"

echo -e "${GREEN}LazyDnD Installer${NC}"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="macos"
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

BINARY_NAME="lazydnd-${OS}-${ARCH}"
echo -e "Detected system: ${YELLOW}${OS} ${ARCH}${NC}"
echo -e "Binary to download: ${YELLOW}${BINARY_NAME}${NC}"
echo ""

# Get latest release version
echo "Fetching latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "Latest version: ${GREEN}${LATEST_VERSION}${NC}"
echo ""

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"

# Download binary
echo "Downloading from: $DOWNLOAD_URL"
TMP_FILE=$(mktemp)
if curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "${GREEN}✓ Download successful${NC}"
else
    echo -e "${RED}✗ Download failed${NC}"
    rm -f "$TMP_FILE"
    exit 1
fi

# Make executable
chmod +x "$TMP_FILE"

# Install binary
echo ""
echo "Installing to ${INSTALL_DIR}/lazydnd"
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "${INSTALL_DIR}/lazydnd"
else
    echo "Requesting sudo permissions to install to ${INSTALL_DIR}..."
    sudo mv "$TMP_FILE" "${INSTALL_DIR}/lazydnd"
fi

echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "Run 'lazydnd' to start the application"
