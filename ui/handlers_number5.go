// ui/handlers_number5.go
package ui

// handleNumber5 handles the '5' key
func handleNumber5(m Model, msg KeyMsg) (Model, Cmd) {
	if m.InputMode {
		m.DiceInput += "5"
	} else if m.InitiativeInputMode || m.InitiativeEditMode {
		m.InitiativeInput += "5"
	} else {
		m.ActivePanel = Notes
	}
	return m, nil
}
