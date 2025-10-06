#!/bin/bash
# build.sh - Cross-platform build script for LazyD&D

echo "Building LazyD&D for multiple platforms..."

# Create build directory if it doesn't exist
mkdir -p build

# Build for Linux (amd64)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o build/lazydnd-linux-amd64 .
if [ $? -eq 0 ]; then
    echo "✓ Linux (amd64) build successful: build/lazydnd-linux-amd64"
else
    echo "✗ Linux (amd64) build failed"
fi

# Build for Linux (arm64)
echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o build/lazydnd-linux-arm64 .
if [ $? -eq 0 ]; then
    echo "✓ Linux (arm64) build successful: build/lazydnd-linux-arm64"
else
    echo "✗ Linux (arm64) build failed"
fi

# Build for macOS (amd64 - Intel)
echo "Building for macOS (amd64 - Intel)..."
GOOS=darwin GOARCH=amd64 go build -o build/lazydnd-macos-amd64 .
if [ $? -eq 0 ]; then
    echo "✓ macOS (amd64) build successful: build/lazydnd-macos-amd64"
else
    echo "✗ macOS (amd64) build failed"
fi

# Build for macOS (arm64 - Apple Silicon)
echo "Building for macOS (arm64 - Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o build/lazydnd-macos-arm64 .
if [ $? -eq 0 ]; then
    echo "✓ macOS (arm64) build successful: build/lazydnd-macos-arm64"
else
    echo "✗ macOS (arm64) build failed"
fi

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o build/lazydnd-windows-amd64.exe .
if [ $? -eq 0 ]; then
    echo "✓ Windows (amd64) build successful: build/lazydnd-windows-amd64.exe"
else
    echo "✗ Windows (amd64) build failed"
fi

echo ""
echo "Build complete! Executables are in the 'build/' directory:"
ls -lh build/

echo ""
echo "Usage:"
echo "  Linux (Intel/AMD):     ./build/lazydnd-linux-amd64"
echo "  Linux (ARM):           ./build/lazydnd-linux-arm64"
echo "  macOS (Intel):         ./build/lazydnd-macos-amd64"
echo "  macOS (Apple Silicon): ./build/lazydnd-macos-arm64"
echo "  Windows:               ./build/lazydnd-windows-amd64.exe"

