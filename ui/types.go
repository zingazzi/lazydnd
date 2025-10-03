// ui/types.go
package ui

// PanelType represents the different panel types
type PanelType int

const (
	DiceRoller PanelType = iota
	CharacterSheet
	Spells
	CampaignNotes
)

// Model represents the main application state
type Model struct {
	ActivePanel  PanelType
	Width        int
	Height       int
	DiceInput    string
	DiceResult   string
	DiceHistory  []string
	InputMode    bool
	ScrollOffset map[PanelType]int
}

// Panel configuration
var PanelNames = []string{
	"Dice Roller",
	"Character Sheet",
	"Spells",
	"Campaign Notes",
}
