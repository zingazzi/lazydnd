// ui/animation.go
package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AnimationTickMsg is sent on each animation frame
type AnimationTickMsg time.Time

// AnimationState tracks the current animation state
type AnimationState struct {
	Active         bool
	Progress       float64 // 0.0 to 1.0
	SourcePanel    PanelType
	TargetPanel    PanelType
	StartTime      time.Time
	Duration       time.Duration
	AnimationType  string // "slide_left", "slide_right", "fade"
}

// Animation constants
const (
	DefaultAnimationDuration = 200 * time.Millisecond
	AnimationFPS             = 60
	AnimationTickInterval    = time.Second / AnimationFPS
)

// StartPanelTransition initializes a panel transition animation
func (m Model) StartPanelTransition(targetPanel PanelType) (Model, tea.Cmd) {
	// Skip animation if already on target panel
	if m.ActivePanel == targetPanel {
		return m, nil
	}

	// Determine animation direction based on panel order
	animType := "slide_right"
	if targetPanel < m.ActivePanel {
		animType = "slide_left"
	}

	m.Animation = AnimationState{
		Active:        true,
		Progress:      0.0,
		SourcePanel:   m.ActivePanel,
		TargetPanel:   targetPanel,
		StartTime:     time.Now(),
		Duration:      DefaultAnimationDuration,
		AnimationType: animType,
	}

	// Start animation ticker
	return m, animationTickCmd()
}

// UpdateAnimation updates the animation state and returns true if animation is complete
func (m *Model) UpdateAnimation() bool {
	if !m.Animation.Active {
		return false
	}

	elapsed := time.Since(m.Animation.StartTime)
	m.Animation.Progress = float64(elapsed) / float64(m.Animation.Duration)

	if m.Animation.Progress >= 1.0 {
		// Animation complete
		m.Animation.Progress = 1.0
		m.ActivePanel = m.Animation.TargetPanel
		m.Animation.Active = false
		return true
	}

	return false
}

// EaseOutCubic provides smooth deceleration
func EaseOutCubic(t float64) float64 {
	t = t - 1
	return t*t*t + 1
}

// EaseInOutCubic provides smooth acceleration and deceleration
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	t = t - 1
	return 1 + 4*t*t*t
}

// GetAnimationOffset calculates the horizontal offset for panel animation
func (m Model) GetAnimationOffset() int {
	if !m.Animation.Active {
		return 0
	}

	// Apply easing
	easedProgress := EaseOutCubic(m.Animation.Progress)

	// Calculate offset based on panel width
	maxOffset := m.Width
	offset := int(float64(maxOffset) * (1.0 - easedProgress))

	// Reverse direction for slide_left
	if m.Animation.AnimationType == "slide_left" {
		offset = -offset
	}

	return offset
}

// GetPanelOpacity calculates opacity for fade animations (0.0 to 1.0)
func (m Model) GetPanelOpacity() float64 {
	if !m.Animation.Active {
		return 1.0
	}

	if m.Animation.AnimationType == "fade" {
		return EaseInOutCubic(m.Animation.Progress)
	}

	return 1.0
}

// animationTickCmd returns a command that sends animation tick messages
func animationTickCmd() tea.Cmd {
	return tea.Tick(AnimationTickInterval, func(t time.Time) tea.Msg {
		return AnimationTickMsg(t)
	})
}

// HandleAnimationTick processes animation tick messages
func HandleAnimationTick(m Model) (Model, tea.Cmd) {
	if !m.Animation.Active {
		return m, nil
	}

	// Update animation state
	complete := m.UpdateAnimation()

	// Continue animation or stop
	if !complete {
		return m, animationTickCmd()
	}

	return m, nil
}
