// ui/panel_state.go
package ui

// DiceRollerState holds all state related to the dice roller panel
type DiceRollerState struct {
	Input            string
	Result           string
	History          []string
	Commands         []string
	LastCommand      string
	InputMode        bool
	HistoryMode      bool
	HistoryIndex     int
	Macros           map[string]string
	MacroListMode    bool
	SelectedMacro    int
	ShowMacroPrompt  bool
	MacroNameInput   string
	MacroFormulaInput string
	MacroInputStep   int
}

// InitiativeState holds all state related to the initiative tracker panel
type InitiativeState struct {
	List                []InitiativeEntry
	Input               string
	InputMode           bool
	InputType           string
	SelectedEntry       int
	TempEntry           InitiativeEntry
	EditMode            bool
	EditType            string
	ListMode            bool
	CurrentTurn         int
	RoundCounter        int
	ShowQuickHPPopup    bool
	QuickHPInput        string
	QuickHPMode         string
	MultiTargetMode     bool
	SelectedTargets     map[int]bool
	ShowMultiTargetPopup bool
	MultiTargetInput    string
	MultiTargetType     string
	MultiTargetSaveMode bool
	TargetSaveResults   map[int]string
	HPUndoStack         []HPHistoryEntry
	HPRedoStack         []HPHistoryEntry
}

// SpellState holds all state related to the spells panel
type SpellState struct {
	SearchInput        string
	SearchMode         bool
	SelectedSpell      *Spell
	Suggestions        []string
	SuggestionIndex    int
	LevelFilter        string
	LevelFilterMode    bool
	ActiveSpells       []ActiveSpell
	ActiveSpellIndex   int
	ActiveSpellListMode bool
	ShowCastSpellPrompt bool
	CastSpellInput     string
	CastSpellInputMode bool
	SpellToCast        *Spell
}

// MonsterState holds all state related to the monsters panel
type MonsterState struct {
	SearchInput      string
	SearchMode       bool
	SelectedMonster  *Monster
	Suggestions      []string
	SuggestionIndex  int
	CRFilter         string
	CRFilterMode     bool
}

// NotesState holds all state related to the notes panel
type NotesState struct {
	Content      string
	Input        string
	EditMode     bool
	SearchMode   bool
	SearchInput  string
	SearchResult []int
}

// EncounterState holds all state related to the encounter builder panel
type EncounterState struct {
	PartySize              int
	PartyLevel             int
	Monsters               []EncounterMonster
	SelectedIndex          int
	SavedEncounters        []Encounter
	ListMode               bool
	NameInput              string
	ShowPrompt             bool
	LoadedTemplateName     string
	BuilderMode            string
	CRFilter               string
	FilterActive           bool
	SelectedSaved          int
	AddingMonster          bool
	Environment            string
	Difficulty             string
	Generating             bool
	EnvironmentIndex       int
	DifficultyIndex        int
	GeneratorFocus         string
	AvailableEnvironments  []string
}

// PopupState holds all state related to popups and overlays
type PopupState struct {
	ShowHelp              bool
	HelpScrollOffset      int
	ShowAction            bool
	ActionActions         []MonsterAction
	ActionIndex           int
	ActionMonster         string
	ActionAdvantage       bool
	ActionDisadvantage    bool
	ShowSavingThrow       bool
	ShowCondition         bool
	ConditionMode         string
	ConditionInput        string
	ConditionDurationInput string
	ConditionInputStep    int
	SelectedConditionIdx  int
	SelectedConditionNameIdx int
	ShowSave              bool
	ShowLoad              bool
	ShowRename            bool
	SaveInput             string
	CurrentCampaignFile   string
	CurrentCampaignName   string
	ShowEncounterPrompt   bool
}

// GlobalState holds application-wide state
type GlobalState struct {
	DebugMode         bool
	CampaignList      []string
	CampaignListIndex int
	LastAutoSave      string
	ErrorMessage      string
	ErrorVisible      bool
}
