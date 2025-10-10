// ui/dice_helpers.go
package ui

// addToHistory adds a dice result to history, respecting the configured history size
func (m *Model) addToHistory(result, command string) {
	m.DiceHistory = append(m.DiceHistory, result)
	m.DiceCommands = append(m.DiceCommands, command)

	historySize := m.Config.DiceRoller.HistorySize
	if historySize < 1 {
		historySize = 15 // Fallback to default
	}

	// Trim history if it exceeds the configured size
	if len(m.DiceHistory) > historySize {
		m.DiceHistory = m.DiceHistory[1:]
		m.DiceCommands = m.DiceCommands[1:]
	}
}
