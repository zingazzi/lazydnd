// ui/tview/colors.go
package tview

import (
	"fmt"
	"lazydnd/config"

	"github.com/gdamore/tcell/v2"
)

// ColorConverter converts hex colors to tcell.Color
type ColorConverter struct {
	primaryColor   tcell.Color
	borderColor    tcell.Color
	highlightColor tcell.Color
	errorColor     tcell.Color
	successColor   tcell.Color
}

// NewColorConverter creates a new color converter from config
func NewColorConverter(cfg *config.Config) *ColorConverter {
	cc := &ColorConverter{}

	// Default colors
	primaryColor := "#7D56F4"
	borderColor := "#444444"
	highlightColor := "#00FF00"
	errorColor := "#FF0000"
	successColor := "#00FF00"

	if cfg != nil {
		primaryColor = cfg.Theme.PrimaryColor
		borderColor = cfg.Theme.BorderColor
		highlightColor = cfg.Theme.HighlightColor
		errorColor = cfg.Theme.ErrorColor
		successColor = cfg.Theme.SuccessColor
	}

	cc.primaryColor = hexToColor(primaryColor)
	cc.borderColor = hexToColor(borderColor)
	cc.highlightColor = hexToColor(highlightColor)
	cc.errorColor = hexToColor(errorColor)
	cc.successColor = hexToColor(successColor)

	return cc
}

// hexToColor converts a hex color string to tcell.Color
func hexToColor(hex string) tcell.Color {
	if len(hex) == 0 || hex[0] != '#' {
		return tcell.ColorDefault
	}

	if len(hex) == 7 {
		// Parse #RRGGBB format
		var r, g, b uint8
		_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return tcell.ColorDefault
		}
		return tcell.NewRGBColor(int32(r), int32(g), int32(b))
	}

	return tcell.ColorDefault
}

// PrimaryColor returns the primary color
func (cc *ColorConverter) PrimaryColor() tcell.Color {
	return cc.primaryColor
}

// BorderColor returns the border color
func (cc *ColorConverter) BorderColor() tcell.Color {
	return cc.borderColor
}

// HighlightColor returns the highlight color
func (cc *ColorConverter) HighlightColor() tcell.Color {
	return cc.highlightColor
}

// ErrorColor returns the error color
func (cc *ColorConverter) ErrorColor() tcell.Color {
	return cc.errorColor
}

// SuccessColor returns the success color
func (cc *ColorConverter) SuccessColor() tcell.Color {
	return cc.successColor
}

// DarkenColor returns a darker version of a color
func (cc *ColorConverter) DarkenColor(color tcell.Color) tcell.Color {
	// Simple darkening - reduce RGB values by 30%
	r, g, b := color.RGB()
	return tcell.NewRGBColor(
		int32(float64(r)*0.7),
		int32(float64(g)*0.7),
		int32(float64(b)*0.7),
	)
}

// HexToColor is a standalone function to convert hex to tcell.Color
func HexToColor(hex string) tcell.Color {
	return hexToColor(hex)
}
