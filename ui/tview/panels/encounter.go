// ui/tview/panels/encounter.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewEncounterBuilderPanel creates a new encounter builder panel widget
func NewEncounterBuilderPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(true) // Enable scrolling for long content
	panel.SetBorder(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.EncounterBuilder)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	return panel
}
