#!/bin/bash
# Test build script for TView migration

set -e

echo "=== Testing TView Migration Build ==="
echo ""

echo "1. Getting latest tview version..."
go get github.com/rivo/tview@latest

echo ""
echo "2. Running go mod tidy..."
go mod tidy

echo ""
echo "3. Checking dependencies..."
go list -m github.com/gdamore/tcell/v2
go list -m github.com/rivo/tview

echo ""
echo "4. Building TView package..."
go build ./ui/tview/...

echo ""
echo "5. Building main_tview.go..."
go build -o /tmp/lazydnd-tview-test main_tview.go

echo ""
echo "✅ Build successful!"
echo "Binary created at: /tmp/lazydnd-tview-test"
