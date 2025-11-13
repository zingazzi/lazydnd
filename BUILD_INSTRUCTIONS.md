# Building the TView Migration

## Quick Build

To build the TView version, run:

```bash
go build -tags tview -o lazydnd
```

Or use the build script:

```bash
bash build_tview.sh
```

## Build Tags

The project uses build tags to separate the two implementations:

- **Bubble Tea (default)**: `go build -o lazydnd` or `go build -tags '!tview' -o lazydnd`
- **TView**: `go build -tags tview -o lazydnd`

## Setup Steps

1. **Get dependencies:**
   ```bash
   go get github.com/rivo/tview@latest
   go get github.com/gdamore/tcell/v2@latest
   go mod tidy
   ```

2. **Build:**
   ```bash
   go build -tags tview -o lazydnd
   ```

3. **Run:**
   ```bash
   ./lazydnd
   ```

## Troubleshooting

If you get "missing go.sum entry" errors:
```bash
go mod tidy
```

If you get "main redeclared" errors:
- Make sure you're using build tags: `go build -tags tview -o lazydnd`
- The build tags ensure only one main function is compiled

If dependencies aren't found:
```bash
go clean -modcache
go mod download
go mod tidy
```

## Current Status

✅ Code structure complete
✅ Build tags configured
✅ Dependencies added to go.mod
⏳ Handler implementation (Phase 4) - placeholder
⏳ Full functionality (Phases 5-10) - pending
