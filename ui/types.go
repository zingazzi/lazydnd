// ui/types.go
package ui

// PanelType represents the different panel types
type PanelType int

const (
	DiceRoller PanelType = iota
	InitiativeTracker
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
	// Initiative tracker state
	InitiativeList      []InitiativeEntry
	InitiativeInput     string
	InitiativeInputMode bool
	InitiativeInputType string // "player", "monster", "initiative"
	SelectedEntry       int
	TempEntry           InitiativeEntry // Temporary storage while building entry
	// Initiative edit state
	InitiativeEditMode  bool
	InitiativeEditType  string // "initiative", "hp", "delete"
	InitiativeListMode  bool   // When true, navigating the list instead of adding entries
}

// InitiativeEntry represents a player or monster in the initiative tracker
type InitiativeEntry struct {
	Name       string
	Type       string // "player" or "monster"
	Initiative int
	HP         int    // Only for monsters
	MaxHP      int    // Only for monsters
	AC         int    // Only for monsters
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
	"Initiative Tracker",
	"Spells",
	"Campaign Notes",
}
