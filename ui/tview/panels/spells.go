// ui/tview/panels/spells.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewSpellsPanel creates a new spells panel widget
func NewSpellsPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)
	panel.SetBorder(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.Spells)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	return panel
}
