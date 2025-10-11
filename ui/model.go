// ui/model.go
package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AutoSaveTickMsg is sent every minute to check for autosave
type AutoSaveTickMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// tickCmd returns a command that sends a tick every minute
func tickCmd() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return AutoSaveTickMsg(t)
	})
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		return HandleNavigation(m, msg)

	case AutoSaveTickMsg:
		// Handle autosave
		m = handleAutoSave(m)
		// Schedule next tick
		return m, tickCmd()

	case AnimationTickMsg:
		// Handle animation updates
		return HandleAnimationTick(m)
	}

	return m, nil
}
