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
	DiceCommands    []string // Store original commands for history
	LastDiceCommand string
	InputMode       bool
	DiceHistoryMode bool // When true, navigating dice history to select
	HistoryIndex    int  // Selected index in history (-1 = none)
	ScrollOffset    map[PanelType]int
	// Spell search state
	SpellSearchInput string
	SpellSearchMode  bool
	SelectedSpell    *Spell
	SpellSuggestions []string
	SuggestionIndex  int
	// Initiative tracker state
	InitiativeList      []InitiativeEntry
	InitiativeInput     string
	InitiativeInputMode bool
	InitiativeInputType string // "player", "monster", "initiative"
	SelectedEntry       int
	TempEntry           InitiativeEntry // Temporary storage while building entry
	CurrentTurn         int             // Index of current turn in initiative order (-1 = no combat started)
	RoundCounter        int             // Current combat round (0 = combat not started, 1+ = active rounds)
	// Initiative edit state
	InitiativeEditMode bool
	InitiativeEditType string // "initiative", "hp", "delete"
	InitiativeListMode bool   // When true, navigating the list instead of adding entries
	// Monster search state
	MonsterSearchInput     string
	MonsterSearchMode      bool
	SelectedMonster        *Monster
	MonsterSuggestions     []string
	MonsterSuggestionIndex int
	// Help popup state
	ShowHelpPopup bool
	// Action popup state
	ShowActionPopup    bool
	ActionPopupActions []MonsterAction
	ActionPopupIndex   int
	ActionPopupMonster string // Name of the monster whose actions are shown
	// Saving throw popup state
	ShowSavingThrowPopup bool
	// Save/Load state
	ShowSavePopup       bool
	ShowLoadPopup       bool
	ShowRenamePopup     bool
	SaveInput           string
	CurrentCampaignFile string
	CurrentCampaignName string
	CampaignList        []string
	CampaignListIndex   int
	LastAutoSave        string
	AutoSaveEnabled     bool
}

// InitiativeEntry represents a player or monster in the initiative tracker
type InitiativeEntry struct {
	Name        string
	Type        string // "player" or "monster"
	Initiative  int
	HP          int      // Only for monsters
	MaxHP       int      // Only for monsters
	AC          int      // Only for monsters
	MonsterData *Monster // Link to full monster data for actions
	InstanceNum int      // Instance number for duplicates (0 = no number shown)
	BaseName    string   // Original name without number
	MonsterName string   // Original monster name for save/load persistence
}

// Spell represents a D&D spell
type Spell struct {
	Name           string   `json:"name"`
	Level          int      `json:"level"`
	School         string   `json:"school"`
	Classes        []string `json:"classes"`
	ActionType     string   `json:"actionType"`
	Concentration  bool     `json:"concentration"`
	Ritual         bool     `json:"ritual"`
	Range          string   `json:"range"`
	Components     []string `json:"components"`
	Material       string   `json:"material,omitempty"`
	Duration       string   `json:"duration"`
	Description    string   `json:"description"`
	CantripUpgrade string   `json:"cantripUpgrade,omitempty"`
}

// MonsterAction represents a single action a monster can take
type MonsterAction struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Roll        string `json:"roll,omitempty"`
	Reach       string `json:"reach,omitempty"`
	Range       string `json:"range,omitempty"`
	Damage      string `json:"damage,omitempty"`
	DamageType  string `json:"damage_type,omitempty"`
	SaveDC      string `json:"save_dc,omitempty"`
	SaveType    string `json:"save_type,omitempty"`
}

// Monster represents a D&D monster
type Monster struct {
	Name             string          `json:"name"`
	Meta             string          `json:"meta"`
	ArmorClass       string          `json:"Armor Class"`
	HitPoints        string          `json:"Hit Points"`
	Speed            string          `json:"Speed"`
	STR              string          `json:"STR"`
	STRMod           string          `json:"STR_mod"`
	DEX              string          `json:"DEX"`
	DEXMod           string          `json:"DEX_mod"`
	CON              string          `json:"CON"`
	CONMod           string          `json:"CON_mod"`
	INT              string          `json:"INT"`
	INTMod           string          `json:"INT_mod"`
	WIS              string          `json:"WIS"`
	WISMod           string          `json:"WIS_mod"`
	CHA              string          `json:"CHA"`
	CHAMod           string          `json:"CHA_mod"`
	SavingThrows     string          `json:"Saving Throws,omitempty"`
	Skills           string          `json:"Skills,omitempty"`
	Senses           string          `json:"Senses,omitempty"`
	Languages        string          `json:"Languages,omitempty"`
	Challenge        string          `json:"Challenge"`
	Traits           string          `json:"Traits,omitempty"`
	Actions          string          `json:"Actions,omitempty"`
	LegendaryActions string          `json:"Legendary Actions,omitempty"`
	ImgURL           string          `json:"img_url,omitempty"`
	ActionNumber     int             `json:"ActionNumber"`
	ActionList       []MonsterAction `json:"ActionList"`
}

// SaveState represents a saved campaign state
type SaveState struct {
	CampaignName   string                 `json:"campaign_name"`
	SavedAt        string                 `json:"saved_at"`
	InitiativeList []SavedInitiativeEntry `json:"initiative_list"`
	CurrentTurn    int                    `json:"current_turn"`
	RoundCounter   int                    `json:"round_counter"`
}

// SavedInitiativeEntry represents an initiative entry for persistence
type SavedInitiativeEntry struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Initiative  int    `json:"initiative"`
	HP          int    `json:"hp"`
	MaxHP       int    `json:"max_hp"`
	AC          int    `json:"ac"`
	MonsterName string `json:"monster_name"`
	InstanceNum int    `json:"instance_num"`
	BaseName    string `json:"base_name"`
}

// Panel configuration
var PanelNames = []string{
	"🎲 Dice Roller",
	"⚔️  Initiative Tracker",
	"✨ Spells",
	"🐉 Monsters",
}
