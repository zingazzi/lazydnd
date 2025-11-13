#!/bin/bash
# Build script for TView migration
set -e

echo "=== Building TView Migration ==="
echo ""

echo "1. Getting latest dependencies..."
go get github.com/rivo/tview@latest
go get github.com/gdamore/tcell/v2@latest

echo ""
echo "2. Running go mod tidy..."
go mod tidy

echo ""
echo "3. Building with standard command: go build -tags tview -o lazydnd"
go build -tags tview -o lazydnd

echo ""
echo "✅ Build successful!"
echo "Binary created: ./lazydnd"
