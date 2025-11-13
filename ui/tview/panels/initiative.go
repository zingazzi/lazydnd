// ui/tview/panels/initiative.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewInitiativeTrackerPanel creates a new initiative tracker panel widget
func NewInitiativeTrackerPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.InitiativeTracker)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	return panel
}
