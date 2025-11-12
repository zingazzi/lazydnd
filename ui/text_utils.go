// ui/text_utils.go
package ui

import (
	"strings"
	"unicode/utf8"
)

// truncateText truncates text to maxWidth using runes (Unicode-safe)
func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return text
	}

	// Leave room for ellipsis
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

// wrapText wraps text to maxWidth, respecting word boundaries when possible
func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{""}
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return []string{text}
	}

	var lines []string
	var currentLine []rune

	for _, r := range runes {
		// Check if adding this rune would exceed width
		if len(currentLine) >= maxWidth {
			lines = append(lines, string(currentLine))
			currentLine = []rune{r}
		} else {
			currentLine = append(currentLine, r)
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, string(currentLine))
	}

	return lines
}

// measureTextWidth returns the display width of text (rune count)
func measureTextWidth(text string) int {
	return utf8.RuneCountInString(text)
}

// padTextRight pads text to the right to reach target width
func padTextRight(text string, targetWidth int) string {
	width := measureTextWidth(text)
	if width >= targetWidth {
		return text
	}

	return text + strings.Repeat(" ", targetWidth-width)
}
