// ui/tview_key_adapter.go
// Adapter to convert TCell events to handler-compatible format
package ui

// TViewKeyMsg is a minimal adapter for tea.KeyMsg used by handlers
// This allows existing handlers to work with TView without importing bubbletea
type TViewKeyMsg struct {
	keyStr string
}

// String returns the key string representation
func (k TViewKeyMsg) String() string {
	return k.keyStr
}

// Type returns "key" to match tea.KeyMsg behavior
func (k TViewKeyMsg) Type() string {
	return "key"
}

// NewTViewKeyMsg creates a TViewKeyMsg from a key string
func NewTViewKeyMsg(keyStr string) TViewKeyMsg {
	return TViewKeyMsg{keyStr: keyStr}
}

