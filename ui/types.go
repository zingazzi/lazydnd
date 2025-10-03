// ui/types.go
package ui

// PanelType represents the different panel types
type PanelType int

const (
	DiceRoller PanelType = iota
	InitiativeTracker
	Spells
	Monsters
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
	// Monster search state
	MonsterSearchInput    string
	MonsterSearchMode     bool
	SelectedMonster       *Monster
	MonsterSuggestions    []string
	MonsterSuggestionIndex int
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

// Monster represents a D&D monster
type Monster struct {
	Name             string `json:"name"`
	Meta             string `json:"meta"`
	ArmorClass       string `json:"Armor Class"`
	HitPoints        string `json:"Hit Points"`
	Speed            string `json:"Speed"`
	STR              string `json:"STR"`
	STRMod           string `json:"STR_mod"`
	DEX              string `json:"DEX"`
	DEXMod           string `json:"DEX_mod"`
	CON              string `json:"CON"`
	CONMod           string `json:"CON_mod"`
	INT              string `json:"INT"`
	INTMod           string `json:"INT_mod"`
	WIS              string `json:"WIS"`
	WISMod           string `json:"WIS_mod"`
	CHA              string `json:"CHA"`
	CHAMod           string `json:"CHA_mod"`
	SavingThrows     string `json:"Saving Throws,omitempty"`
	Skills           string `json:"Skills,omitempty"`
	Senses           string `json:"Senses,omitempty"`
	Languages        string `json:"Languages,omitempty"`
	Challenge        string `json:"Challenge"`
	Traits           string `json:"Traits,omitempty"`
	Actions          string `json:"Actions,omitempty"`
	LegendaryActions string `json:"Legendary Actions,omitempty"`
	ImgURL           string `json:"img_url,omitempty"`
}

// Panel configuration
var PanelNames = []string{
	"Dice Roller",
	"Initiative Tracker",
	"Spells",
	"Monsters",
}
