// ui/tview_colors.go
package ui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// ColorConverter converts hex colors to tcell.Color
type ColorConverter struct{}

// HexToColor converts a hex color string (#RRGGBB) to tcell.Color
func HexToColor(hex string) tcell.Color {
	if len(hex) == 0 || hex[0] != '#' {
		return tcell.ColorDefault
	}

	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return tcell.ColorDefault
	}

	r, err1 := strconv.ParseUint(hex[0:2], 16, 8)
	g, err2 := strconv.ParseUint(hex[2:4], 16, 8)
	b, err3 := strconv.ParseUint(hex[4:6], 16, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return tcell.ColorDefault
	}

	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

// ColorNames provides named colors for common use cases
var ColorNames = struct {
	Default     tcell.Color
	White       tcell.Color
	Black       tcell.Color
	Yellow      tcell.Color
	Red         tcell.Color
	Green       tcell.Color
	Blue        tcell.Color
	Purple      tcell.Color
	DarkGray    tcell.Color
	LightGray   tcell.Color
	DarkPurple  tcell.Color
	DarkRed     tcell.Color
}{
	Default:     tcell.ColorDefault,
	White:       tcell.ColorWhite,
	Black:       tcell.ColorBlack,
	Yellow:      tcell.ColorYellow,
	Red:         tcell.ColorRed,
	Green:       tcell.ColorGreen,
	Blue:        tcell.ColorBlue,
	Purple:      HexToColor("#7D56F4"),
	DarkGray:    HexToColor("#333333"),
	LightGray:   HexToColor("#CCCCCC"),
	DarkPurple:  HexToColor("#5A3D9E"),
	DarkRed:     HexToColor("#8B0000"),
}
