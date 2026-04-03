package ui

type PartyMember struct {
	ID          string
	Name        string
	Role        string
	DiceSummary string
	Status      string
	HP          int
	MaxHP       int
	Selected    bool
	Downed      bool
}

type PartySelectionModel struct {
	Candidates    []PartyMember
	SelectedCount int
	CanStart      bool
}

type MapNode struct {
	ID     string
	Name   string
	Kind   string
	Detail string
}

type MapModel struct {
	CurrentNodeID string
	Nodes         []MapNode
}

type DieView struct {
	ID     string
	Label  string
	Detail string
	Forced bool
}

type CombatModel struct {
	EncounterName    string
	Turn             int
	Party            []PartyMember
	Enemies          []PartyMember
	AvailableDice    []DieView
	SelectedDice     []DieView
	RevealedPatterns []string
	Logs             []string
	CanResolve       bool
	AllyDefense      int
	EnemyDefense     int
	DamageBoost      int
}

type OutcomeModel struct {
	Title       string
	Body        string
	CanContinue bool
	RunEnded    bool
}

type Model struct {
	CurrentScreen  string
	HeaderTitle    string
	HeaderSubtitle string
	ViewportWidth  float64
	ViewportHeight float64
	PartyRoster    []PartyMember
	PartySelection PartySelectionModel
	Map            MapModel
	Combat         CombatModel
	Outcome        OutcomeModel
}

type Callbacks struct {
	OnToggleParty   func(string)
	OnStartRun      func()
	OnSelectMapNode func(string)
	OnSelectDie     func(string)
	OnResolveTurn   func()
	OnContinue      func()
	OnRestart       func()
}
