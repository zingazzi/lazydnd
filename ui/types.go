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
	ActivePanel     PanelType
	Width           int
	Height          int
	DiceInput       string
	DiceResult      string
	DiceHistory     []string
	LastDiceCommand string
	InputMode       bool
	ScrollOffset    map[PanelType]int
	// Spell search state
	SpellSearchInput    string
	SpellSearchMode     bool
	SelectedSpell       *Spell
	SpellSuggestions    []string
	SuggestionIndex     int
}

// Spell represents a D&D spell
type Spell struct {
	Name            string   `json:"name"`
	Level           int      `json:"level"`
	School          string   `json:"school"`
	Classes         []string `json:"classes"`
	ActionType      string   `json:"actionType"`
	Concentration   bool     `json:"concentration"`
	Ritual          bool     `json:"ritual"`
	Range           string   `json:"range"`
	Components      []string `json:"components"`
	Material        string   `json:"material,omitempty"`
	Duration        string   `json:"duration"`
	Description     string   `json:"description"`
	CantripUpgrade  string   `json:"cantripUpgrade,omitempty"`
}

// Panel configuration
var PanelNames = []string{
	"Dice Roller",
	"Character Sheet",
	"Spells",
	"Campaign Notes",
}
