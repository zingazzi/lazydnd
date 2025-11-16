// ui/tview/panels/notes.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewNotesPanel creates a new notes panel widget
func NewNotesPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(true) // Enable scrolling for long content
	panel.SetBorder(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.Notes)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	return panel
}
