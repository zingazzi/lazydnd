# Building the TView Migration

## Quick Build

To build the application, run:

```bash
go build -o lazydnd
```

The project now uses TView exclusively. No build tags needed.

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
