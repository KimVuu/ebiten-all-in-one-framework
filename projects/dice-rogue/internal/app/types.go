package app

import (
	"sort"
	"strings"
)

type ScreenID string

const (
	ScreenPartySelection ScreenID = "party-selection"
	ScreenMap            ScreenID = "map"
	ScreenCombat         ScreenID = "combat"
	ScreenOutcome        ScreenID = "outcome"
)

type DieKind string

const (
	DieKindAttack  DieKind = "attack"
	DieKindDefense DieKind = "defense"
	DieKindSkill   DieKind = "skill"
)

type FaceKind string

const (
	FaceKindValue    FaceKind = "value"
	FaceKindSuccess  FaceKind = "success"
	FaceKindFailure  FaceKind = "failure"
	FaceKindCritical FaceKind = "critical"
	FaceKindEscape   FaceKind = "escape"
)

type EncounterKind string

const (
	EncounterKindNormal EncounterKind = "normal"
	EncounterKindElite  EncounterKind = "elite"
	EncounterKindBoss   EncounterKind = "boss"
)

type NodeKind string

const (
	NodeKindStart  NodeKind = "start"
	NodeKindNormal NodeKind = "normal"
	NodeKindElite  NodeKind = "elite"
	NodeKindBoss   NodeKind = "boss"
	NodeKindRest   NodeKind = "rest"
)

type CombatOutcomeType string

const (
	CombatOutcomeNone    CombatOutcomeType = "none"
	CombatOutcomeVictory CombatOutcomeType = "victory"
	CombatOutcomeEscape  CombatOutcomeType = "escape"
	CombatOutcomeDefeat  CombatOutcomeType = "defeat"
)

type DieFace struct {
	Label string
	Kind  FaceKind
	Value int
}

type DieSpec struct {
	ID              string
	OwnerID         string
	Name            string
	Kind            DieKind
	Faces           []DieFace
	EffectID        string
	EnabledInBattle bool
}

type StatusEffect struct {
	ID             string
	Name           string
	Magnitude      int
	RemainingTurns int
}

type EncounterPattern struct {
	ID      string
	Label   string
	Attacks []int
	Defense int
}

type UnitState struct {
	ID         string
	Name       string
	Role       string
	MaxHP      int
	HP         int
	Downed     bool
	Dice       []DieSpec
	Counters   map[string]int
	Statuses   []StatusEffect
	Patterns   []EncounterPattern
	PatternIdx int
}

type EncounterDefinition struct {
	ID      string
	Name    string
	Kind    EncounterKind
	Enemies []UnitState
}

type EncounterNode struct {
	ID          string
	Name        string
	Kind        NodeKind
	EncounterID string
	NextIDs     []string
}

type SelectedDie struct {
	Die    DieSpec
	Forced bool
}

type TurnResolution struct {
	Turn                 int
	GeneratedAllyDefense int
	GeneratedEnemyDefense int
	DamageBoostPercent   int
	EnemyTargets         []string
	PlayerTargets        []string
	RevealedNextPatterns map[string]string
	Outcome              CombatOutcomeType
	Logs                 []string
}

type CombatState struct {
	NodeID       string
	EncounterID  string
	EncounterName string
	EncounterKind EncounterKind
	Turn         int
	PlayerUnits  []UnitState
	EnemyUnits   []UnitState
	AvailableDice []DieSpec
	GraveyardDice []DieSpec
	SelectedDice []SelectedDie
	AllyDefense  int
	EnemyDefense int
	EnemyStatuses []StatusEffect
	RevealedNextPatterns map[string]string
	Outcome      CombatOutcomeType
	Logs         []string
	random       *RandomSource
}

type OutcomeState struct {
	Title       string
	Body        string
	CanContinue bool
	RunEnded    bool
}

type RunState struct {
	Seed             int64
	Screen           ScreenID
	SelectedPartyIDs []string
	PartyUnits        []UnitState
	CurrentNodeID     string
	NextNodeIDs       []string
	CurrentCombat     *CombatState
	Outcome           OutcomeState
}

type GameConfig struct {
	DebugEnabled bool
	Seed         int64
}

type combatSide string

const (
	sidePlayer combatSide = "player"
	sideEnemy  combatSide = "enemy"
)

const (
	counterHeroGoddess = "hero_goddess_stack"
	effectHeroGoddess  = "hero.goddess"
	effectTankGuard    = "guard.tank"
	effectPriestHeal   = "priest.heal"
	effectGuideInfo    = "guide.info"
	effectGuideWeakness = "guide.weakness"
	effectGuideEscape  = "guide.escape"
	effectArcherShot   = "archer.shot"
	effectSmithForge   = "smith.forge"
	statusWeakness     = "weakness"
)

func cloneUnits(units []UnitState) []UnitState {
	cloned := make([]UnitState, 0, len(units))
	for _, unit := range units {
		next := unit
		next.Dice = append([]DieSpec(nil), unit.Dice...)
		next.Patterns = append([]EncounterPattern(nil), unit.Patterns...)
		if unit.Counters != nil {
			next.Counters = map[string]int{}
			for key, value := range unit.Counters {
				next.Counters[key] = value
			}
		} else {
			next.Counters = map[string]int{}
		}
		next.Statuses = append([]StatusEffect(nil), unit.Statuses...)
		cloned = append(cloned, next)
	}
	return cloned
}

func cloneDicePool(dice []DieSpec) []DieSpec {
	return append([]DieSpec(nil), dice...)
}

func buildCombatDicePool(units []UnitState) []DieSpec {
	pool := make([]DieSpec, 0)
	for _, unit := range units {
		if unit.Downed || unit.HP <= 0 {
			continue
		}
		for _, die := range unit.Dice {
			if !die.EnabledInBattle {
				continue
			}
			pool = append(pool, die)
		}
	}
	sort.SliceStable(pool, func(i, j int) bool {
		return pool[i].ID < pool[j].ID
	})
	return pool
}

func selectedToPool(selected []SelectedDie) []DieSpec {
	pool := make([]DieSpec, 0, len(selected))
	for _, die := range selected {
		pool = append(pool, die.Die)
	}
	return pool
}

func poolContainsOwner(pool []DieSpec, ownerID string) bool {
	for _, die := range pool {
		if die.OwnerID == ownerID {
			return true
		}
	}
	return false
}

func ownerDiceIDs(pool []DieSpec, ownerID string) []string {
	ids := make([]string, 0)
	for _, die := range pool {
		if die.OwnerID == ownerID {
			ids = append(ids, die.ID)
		}
	}
	return ids
}

func joinStrings(parts []string, sep string) string {
	return strings.Join(parts, sep)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
