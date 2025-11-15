// ui/tview/panels/dice_roller.go
package panels

import (
	"lazydnd/ui"

	"github.com/rivo/tview"
)

// NewDiceRollerPanel creates a new dice roller panel widget
func NewDiceRollerPanel(model *ui.Model) *tview.TextView {
	panel := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true)
	panel.SetBorder(true)

	// Set initial content
	updateContent := func() {
		content := model.GetPanelContent(ui.DiceRoller)
		panel.SetText(content)
	}

	// Initial update
	updateContent()

	// TODO: Register update callback when state management is implemented
	// model.RegisterUpdateCallback(updateContent)

	return panel
}
