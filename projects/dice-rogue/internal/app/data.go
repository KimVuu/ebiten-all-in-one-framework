package app

import "fmt"

var characterCatalog = []UnitState{
	{
		ID:    "human-warrior",
		Name:  "Human Warrior",
		Role:  "Balanced Dealer",
		MaxHP: 24,
		HP:    24,
		Dice: []DieSpec{
			valueDie("human-warrior-attack-1", "human-warrior", "Attack 1", DieKindAttack, []int{1, 2, 2, 3, 4, 5}),
			valueDie("human-warrior-attack-2", "human-warrior", "Attack 2", DieKindAttack, []int{1, 2, 2, 3, 4, 5}),
			valueDie("human-warrior-defense-1", "human-warrior", "Defense 1", DieKindDefense, []int{1, 1, 2, 2, 3, 4}),
			skillDie("human-warrior-goddess-1", "human-warrior", "Goddess 1", effectHeroGoddess,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
			skillDie("human-warrior-goddess-2", "human-warrior", "Goddess 2", effectHeroGoddess,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-guard",
		Name:  "Human Guard",
		Role:  "Tank",
		MaxHP: 32,
		HP:    32,
		Dice: []DieSpec{
			valueDie("human-guard-defense-1", "human-guard", "Defense 1", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-2", "human-guard", "Defense 2", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-3", "human-guard", "Defense 3", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-4", "human-guard", "Defense 4", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			skillDie("human-guard-tank-1", "human-guard", "Tank 1", effectTankGuard,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-priest",
		Name:  "Human Priest",
		Role:  "Healer",
		MaxHP: 20,
		HP:    20,
		Dice: []DieSpec{
			skillDie("human-priest-priest-1", "human-priest", "Priest 1", effectPriestHeal,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
			skillDie("human-priest-priest-2", "human-priest", "Priest 2", effectPriestHeal,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
			skillDie("human-priest-priest-3", "human-priest", "Priest 3", effectPriestHeal,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-guide",
		Name:  "Human Guide",
		Role:  "Support",
		MaxHP: 18,
		HP:    18,
		Dice: []DieSpec{
			skillDie("human-guide-info-1", "human-guide", "Info 1", effectGuideInfo,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindEscape, "Escape", 0),
			),
			skillDie("human-guide-weakness-1", "human-guide", "Weakness 1", effectGuideWeakness,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindEscape, "Escape", 0),
			),
			skillDie("human-guide-escape-1", "human-guide", "Escape 1", effectGuideEscape,
				face(FaceKindEscape, "Escape", 0),
				face(FaceKindEscape, "Escape", 0),
				face(FaceKindEscape, "Escape", 0),
				face(FaceKindEscape, "Escape", 0),
				face(FaceKindEscape, "Escape", 0),
				face(FaceKindFailure, "Failure", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "elf-archer",
		Name:  "Elf Archer",
		Role:  "Dealer",
		MaxHP: 20,
		HP:    20,
		Dice: []DieSpec{
			skillDie("elf-archer-shot-1", "elf-archer", "Shot 1", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "Critical", 0),
				face(FaceKindCritical, "Critical", 0),
			),
			skillDie("elf-archer-shot-2", "elf-archer", "Shot 2", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "Critical", 0),
				face(FaceKindCritical, "Critical", 0),
			),
			skillDie("elf-archer-shot-3", "elf-archer", "Shot 3", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "Critical", 0),
				face(FaceKindCritical, "Critical", 0),
			),
			skillDie("elf-archer-shot-4", "elf-archer", "Shot 4", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "Critical", 0),
				face(FaceKindCritical, "Critical", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "dwarf-smith",
		Name:  "Dwarf Smith",
		Role:  "Support Tank",
		MaxHP: 26,
		HP:    26,
		Dice: []DieSpec{
			valueDie("dwarf-smith-attack-1", "dwarf-smith", "Attack 1", DieKindAttack, []int{1, 2, 2, 3, 3, 4}),
			valueDie("dwarf-smith-defense-1", "dwarf-smith", "Defense 1", DieKindDefense, []int{1, 1, 2, 2, 3, 3}),
			valueDie("dwarf-smith-defense-2", "dwarf-smith", "Defense 2", DieKindDefense, []int{1, 1, 2, 2, 3, 3}),
			skillDie("dwarf-smith-forge-1", "dwarf-smith", "Forge 1", effectSmithForge,
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindSuccess, "Success", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
				face(FaceKindFailure, "Failure", 0),
			),
		},
		Counters: map[string]int{},
	},
}

var encounterCatalog = map[string]EncounterDefinition{
	"normal-a": {
		ID:   "normal-a",
		Name: "Normal A",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("spike-wolf", "Spike Wolf", 14,
				pattern("wolf-bite", "Bite", []int{3}, 0),
				pattern("wolf-rush", "Rush", []int{4}, 0),
				pattern("wolf-hide", "Hide", nil, 2),
			),
			enemyTemplate("shell-bug", "Shell Bug", 20,
				pattern("bug-shell", "Shell", nil, 4),
				pattern("bug-sting", "Sting", []int{3}, 0),
				pattern("bug-shell-2", "Shell", nil, 3),
			),
		},
	},
	"normal-b": {
		ID:   "normal-b",
		Name: "Normal B",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("blade-imp-a", "Blade Imp A", 16,
				pattern("imp-slice", "Slice", []int{4}, 0),
				pattern("imp-flurry", "Flurry", []int{2, 2}, 0),
				pattern("imp-guard", "Guard", nil, 1),
			),
			enemyTemplate("blade-imp-b", "Blade Imp B", 16,
				pattern("imp-slice", "Slice", []int{4}, 0),
				pattern("imp-flurry", "Flurry", []int{2, 2}, 0),
				pattern("imp-guard", "Guard", nil, 1),
			),
		},
	},
	"normal-c": {
		ID:   "normal-c",
		Name: "Normal C",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("shell-bug", "Shell Bug", 20,
				pattern("bug-shell", "Shell", nil, 4),
				pattern("bug-sting", "Sting", []int{3}, 0),
				pattern("bug-shell-2", "Shell", nil, 3),
			),
			enemyTemplate("bat-seer", "Bat Seer", 12,
				pattern("bat-peck", "Peck", []int{2}, 0),
				pattern("bat-veil", "Veil", nil, 2),
				pattern("bat-curse", "Curse", []int{3}, 0),
			),
		},
	},
	"elite-1": {
		ID:   "elite-1",
		Name: "Elite",
		Kind: EncounterKindElite,
		Enemies: []UnitState{
			enemyTemplate("guardian-idol", "Guardian Idol", 24,
				pattern("idol-guard", "Guard", nil, 5),
				pattern("idol-slam", "Slam", []int{3}, 0),
				pattern("idol-wall", "Wall", nil, 5),
			),
			enemyTemplate("blade-imp", "Blade Imp", 16,
				pattern("imp-slice", "Slice", []int{4}, 0),
				pattern("imp-flurry", "Flurry", []int{2, 2}, 0),
				pattern("imp-guard", "Guard", nil, 1),
			),
		},
	},
	"boss-1": {
		ID:   "boss-1",
		Name: "Boss",
		Kind: EncounterKindBoss,
		Enemies: []UnitState{
			enemyTemplate("ruin-ogre", "Ruin Ogre", 58,
				pattern("ogre-crush", "Crush", []int{4, 4}, 0),
				pattern("ogre-hide", "Hide", nil, 6),
				pattern("ogre-roar", "Roar", []int{7}, 0),
			),
		},
	},
}

var actMapNodes = map[string]EncounterNode{
	"start": {
		ID:      "start",
		Name:    "Start",
		Kind:    NodeKindStart,
		NextIDs: []string{"normal-a", "normal-b"},
	},
	"normal-a": {
		ID:          "normal-a",
		Name:        "Fork A",
		Kind:        NodeKindNormal,
		EncounterID: "normal-a",
		NextIDs:     []string{"rest-1", "normal-c"},
	},
	"normal-b": {
		ID:          "normal-b",
		Name:        "Fork B",
		Kind:        NodeKindNormal,
		EncounterID: "normal-b",
		NextIDs:     []string{"rest-1", "normal-c"},
	},
	"normal-c": {
		ID:          "normal-c",
		Name:        "Ambush",
		Kind:        NodeKindNormal,
		EncounterID: "normal-c",
		NextIDs:     []string{"elite-1"},
	},
	"rest-1": {
		ID:      "rest-1",
		Name:    "Camp",
		Kind:    NodeKindRest,
		NextIDs: []string{"elite-1"},
	},
	"elite-1": {
		ID:          "elite-1",
		Name:        "Elite Gate",
		Kind:        NodeKindElite,
		EncounterID: "elite-1",
		NextIDs:     []string{"rest-2"},
	},
	"rest-2": {
		ID:      "rest-2",
		Name:    "Shrine",
		Kind:    NodeKindRest,
		NextIDs: []string{"boss-1"},
	},
	"boss-1": {
		ID:          "boss-1",
		Name:        "Boss",
		Kind:        NodeKindBoss,
		EncounterID: "boss-1",
		NextIDs:     nil,
	},
}

func newCharacterState(id string) (UnitState, bool) {
	for _, unit := range characterCatalog {
		if unit.ID == id {
			return cloneUnits([]UnitState{unit})[0], true
		}
	}
	return UnitState{}, false
}

func characterChoices() []UnitState {
	return cloneUnits(characterCatalog)
}

func encounterByID(id string) (EncounterDefinition, bool) {
	encounter, ok := encounterCatalog[id]
	if !ok {
		return EncounterDefinition{}, false
	}
	encounter.Enemies = cloneUnits(encounter.Enemies)
	return encounter, true
}

func mapNodeByID(id string) (EncounterNode, bool) {
	node, ok := actMapNodes[id]
	return node, ok
}

func valueDie(id string, ownerID string, name string, kind DieKind, values []int) DieSpec {
	faces := make([]DieFace, 0, len(values))
	for _, value := range values {
		faces = append(faces, face(FaceKindValue, fmt.Sprintf("%d", value), value))
	}
	return DieSpec{
		ID:              id,
		OwnerID:         ownerID,
		Name:            name,
		Kind:            kind,
		Faces:           faces,
		EnabledInBattle: true,
	}
}

func skillDie(id string, ownerID string, name string, effectID string, faces ...DieFace) DieSpec {
	return DieSpec{
		ID:              id,
		OwnerID:         ownerID,
		Name:            name,
		Kind:            DieKindSkill,
		Faces:           append([]DieFace(nil), faces...),
		EffectID:        effectID,
		EnabledInBattle: true,
	}
}

func face(kind FaceKind, label string, value int) DieFace {
	return DieFace{Kind: kind, Label: label, Value: value}
}

func pattern(id string, label string, attacks []int, defense int) EncounterPattern {
	return EncounterPattern{
		ID:      id,
		Label:   label,
		Attacks: append([]int(nil), attacks...),
		Defense: defense,
	}
}

func enemyTemplate(id string, name string, hp int, patterns ...EncounterPattern) UnitState {
	return UnitState{
		ID:       id,
		Name:     name,
		Role:     "Enemy",
		MaxHP:    hp,
		HP:       hp,
		Patterns: append([]EncounterPattern(nil), patterns...),
		Counters: map[string]int{},
	}
}
