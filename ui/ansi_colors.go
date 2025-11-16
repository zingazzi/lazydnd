// ui/ansi_colors.go
package ui

import "fmt"

// ANSI color codes for TView TextView (with SetDynamicColors(true))
// TView uses square bracket format: [color:XX] or [::XX] for colors
// Format: [color:foreground,background:attributes]text[::]
// Or use standard ANSI with \x1b[ prefix
const (
	ANSIReset      = "\x1b[0m"
	ANSIBold       = "\x1b[1m"
	ANSIDim        = "\x1b[2m"
	ANSIItalic     = "\x1b[3m"
	ANSIUnderline  = "\x1b[4m"
	ANSIBlink      = "\x1b[5m"
	ANSIReverse    = "\x1b[7m"
	ANSICrossedOut = "\x1b[9m"

	// Standard colors
	ANSIFgBlack   = "\x1b[30m"
	ANSIFgRed     = "\x1b[31m"
	ANSIFgGreen   = "\x1b[32m"
	ANSIFgYellow  = "\x1b[33m"
	ANSIFgBlue    = "\x1b[34m"
	ANSIFgMagenta = "\x1b[35m"
	ANSIFgCyan    = "\x1b[36m"
	ANSIFgWhite   = "\x1b[37m"

	// Bright colors
	ANSIFgBrightBlack   = "\x1b[90m"
	ANSIFgBrightRed     = "\x1b[91m"
	ANSIFgBrightGreen   = "\x1b[92m"
	ANSIFgBrightYellow  = "\x1b[93m"
	ANSIFgBrightBlue   = "\x1b[94m"
	ANSIFgBrightMagenta = "\x1b[95m"
	ANSIFgBrightCyan    = "\x1b[96m"
	ANSIFgBrightWhite   = "\x1b[97m"

	// Background colors
	ANSIBgBlack   = "\x1b[40m"
	ANSIBgRed     = "\x1b[41m"
	ANSIBgGreen   = "\x1b[42m"
	ANSIBgYellow  = "\x1b[43m"
	ANSIBgBlue    = "\x1b[44m"
	ANSIBgMagenta = "\x1b[45m"
	ANSIBgCyan    = "\x1b[46m"
	ANSIBgWhite   = "\x1b[47m"
)

// Colorize wraps text with ANSI color codes
func Colorize(text string, colorCode string) string {
	return "\x1b[" + colorCode + text + ANSIReset
}

// ColorizeBold wraps text with ANSI color codes and bold
func ColorizeBold(text string, colorCode string) string {
	return "\x1b[" + colorCode + "\x1b[1m" + text + ANSIReset
}

// ColorizeRGB wraps text with RGB color (for 24-bit color support)
// r, g, b should be 0-255
func ColorizeRGB(text string, r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s%s", r, g, b, text, ANSIReset)
}

// Common color helpers
func Red(text string) string {
	return Colorize(text, "31")
}

func Green(text string) string {
	return Colorize(text, "32")
}

func Yellow(text string) string {
	return Colorize(text, "33")
}

func Blue(text string) string {
	return Colorize(text, "34")
}

func Magenta(text string) string {
	return Colorize(text, "35")
}

func Cyan(text string) string {
	return Colorize(text, "36")
}

func White(text string) string {
	return Colorize(text, "37")
}

func BrightRed(text string) string {
	return Colorize(text, "91")
}

func BrightGreen(text string) string {
	return Colorize(text, "92")
}

func BrightYellow(text string) string {
	return Colorize(text, "93")
}

func BrightBlue(text string) string {
	return Colorize(text, "94")
}

func BrightCyan(text string) string {
	return Colorize(text, "96")
}

// ColorizeWithStyle wraps text with multiple ANSI codes
func ColorizeWithStyle(text string, codes ...string) string {
	result := ""
	for _, code := range codes {
		result += "\x1b[" + code
	}
	result += text + ANSIReset
	return result
}
