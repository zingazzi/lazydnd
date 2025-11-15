// ui/tview/panels/monsters.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewMonstersPanel creates a new monsters panel widget
func NewMonstersPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)
	panel.SetBorder(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.Monsters)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	return panel
}
