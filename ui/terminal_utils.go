// ui/terminal_utils.go
package ui

import (
	"os"
	"strconv"
)

// detectTerminalSize attempts to detect terminal size with fallbacks
func detectTerminalSize() (width, height int) {
	// Try environment variables first (most terminals set these)
	if w := os.Getenv("COLUMNS"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			// Height
			if h := os.Getenv("LINES"); h != "" {
				if height, err := strconv.Atoi(h); err == nil && height > 0 {
					return width, height
				}
			}
		}
	}

	// Fallback: use reasonable defaults
	return 120, 40
}

// isUnicodeTerminal checks if terminal supports Unicode
func isUnicodeTerminal() bool {
	// Check locale
	lang := os.Getenv("LANG")
	if lang != "" {
		return true // Assume Unicode support if locale is set
	}

	// Check terminal type
	term := os.Getenv("TERM")
	unicodeTerms := []string{"xterm-256color", "screen-256color", "tmux-256color"}
	for _, ut := range unicodeTerms {
		if term == ut {
			return true
		}
	}

	return false
}
