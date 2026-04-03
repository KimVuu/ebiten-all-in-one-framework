package app

import "fmt"

var characterCatalog = []UnitState{
	{
		ID:    "human-warrior",
		Name:  "인간 용사",
		Role:  "밸런스 딜러",
		MaxHP: 24,
		HP:    24,
		Dice: []DieSpec{
			valueDie("human-warrior-attack-1", "human-warrior", "공격 주사위 1", DieKindAttack, []int{1, 2, 2, 3, 4, 5}),
			valueDie("human-warrior-attack-2", "human-warrior", "공격 주사위 2", DieKindAttack, []int{1, 2, 2, 3, 4, 5}),
			valueDie("human-warrior-defense-1", "human-warrior", "방어 주사위 1", DieKindDefense, []int{1, 1, 2, 2, 3, 4}),
			skillDie("human-warrior-goddess-1", "human-warrior", "여신의 주사위 1", effectHeroGoddess,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
			skillDie("human-warrior-goddess-2", "human-warrior", "여신의 주사위 2", effectHeroGoddess,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-guard",
		Name:  "인간 방패병",
		Role:  "탱커",
		MaxHP: 32,
		HP:    32,
		Dice: []DieSpec{
			valueDie("human-guard-defense-1", "human-guard", "방어 주사위 1", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-2", "human-guard", "방어 주사위 2", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-3", "human-guard", "방어 주사위 3", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			valueDie("human-guard-defense-4", "human-guard", "방어 주사위 4", DieKindDefense, []int{1, 2, 2, 3, 3, 4}),
			skillDie("human-guard-tank-1", "human-guard", "탱커의 주사위 1", effectTankGuard,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-priest",
		Name:  "인간 여신관",
		Role:  "힐러",
		MaxHP: 20,
		HP:    20,
		Dice: []DieSpec{
			skillDie("human-priest-priest-1", "human-priest", "신관의 주사위 1", effectPriestHeal,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
			skillDie("human-priest-priest-2", "human-priest", "신관의 주사위 2", effectPriestHeal,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
			skillDie("human-priest-priest-3", "human-priest", "신관의 주사위 3", effectPriestHeal,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "human-guide",
		Name:  "인간 길잡이",
		Role:  "서포터",
		MaxHP: 18,
		HP:    18,
		Dice: []DieSpec{
			skillDie("human-guide-info-1", "human-guide", "정보의 주사위 1", effectGuideInfo,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindEscape, "도주", 0),
			),
			skillDie("human-guide-weakness-1", "human-guide", "약점의 주사위 1", effectGuideWeakness,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindEscape, "도주", 0),
			),
			skillDie("human-guide-escape-1", "human-guide", "도주의 주사위 1", effectGuideEscape,
				face(FaceKindEscape, "도주", 0),
				face(FaceKindEscape, "도주", 0),
				face(FaceKindEscape, "도주", 0),
				face(FaceKindEscape, "도주", 0),
				face(FaceKindEscape, "도주", 0),
				face(FaceKindFailure, "실패", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "elf-archer",
		Name:  "엘프 궁수",
		Role:  "딜러",
		MaxHP: 20,
		HP:    20,
		Dice: []DieSpec{
			skillDie("elf-archer-shot-1", "elf-archer", "궁수의 주사위 1", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "치명타", 0),
				face(FaceKindCritical, "치명타", 0),
			),
			skillDie("elf-archer-shot-2", "elf-archer", "궁수의 주사위 2", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "치명타", 0),
				face(FaceKindCritical, "치명타", 0),
			),
			skillDie("elf-archer-shot-3", "elf-archer", "궁수의 주사위 3", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "치명타", 0),
				face(FaceKindCritical, "치명타", 0),
			),
			skillDie("elf-archer-shot-4", "elf-archer", "궁수의 주사위 4", effectArcherShot,
				face(FaceKindValue, "1", 1),
				face(FaceKindValue, "2", 2),
				face(FaceKindValue, "3", 3),
				face(FaceKindValue, "4", 4),
				face(FaceKindCritical, "치명타", 0),
				face(FaceKindCritical, "치명타", 0),
			),
		},
		Counters: map[string]int{},
	},
	{
		ID:    "dwarf-smith",
		Name:  "드워프 대장장이",
		Role:  "서포팅 탱커",
		MaxHP: 26,
		HP:    26,
		Dice: []DieSpec{
			valueDie("dwarf-smith-attack-1", "dwarf-smith", "공격 주사위 1", DieKindAttack, []int{1, 2, 2, 3, 3, 4}),
			valueDie("dwarf-smith-defense-1", "dwarf-smith", "방어 주사위 1", DieKindDefense, []int{1, 1, 2, 2, 3, 3}),
			valueDie("dwarf-smith-defense-2", "dwarf-smith", "방어 주사위 2", DieKindDefense, []int{1, 1, 2, 2, 3, 3}),
			skillDie("dwarf-smith-forge-1", "dwarf-smith", "대장장이의 주사위 1", effectSmithForge,
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindSuccess, "성공", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
				face(FaceKindFailure, "실패", 0),
			),
		},
		Counters: map[string]int{},
	},
}

var encounterCatalog = map[string]EncounterDefinition{
	"normal-a": {
		ID:   "normal-a",
		Name: "일반전 A",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("spike-wolf", "송곳늑대", 14,
				pattern("wolf-bite", "물어뜯기", []int{3}, 0),
				pattern("wolf-rush", "돌진", []int{4}, 0),
				pattern("wolf-hide", "은신", nil, 2),
			),
			enemyTemplate("shell-bug", "껍질벌레", 20,
				pattern("bug-shell", "껍질", nil, 4),
				pattern("bug-sting", "찌르기", []int{3}, 0),
				pattern("bug-shell-2", "껍질", nil, 3),
			),
		},
	},
	"normal-b": {
		ID:   "normal-b",
		Name: "일반전 B",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("blade-imp-a", "칼날임프 A", 16,
				pattern("imp-slice", "베기", []int{4}, 0),
				pattern("imp-flurry", "난타", []int{2, 2}, 0),
				pattern("imp-guard", "수비", nil, 1),
			),
			enemyTemplate("blade-imp-b", "칼날임프 B", 16,
				pattern("imp-slice", "베기", []int{4}, 0),
				pattern("imp-flurry", "난타", []int{2, 2}, 0),
				pattern("imp-guard", "수비", nil, 1),
			),
		},
	},
	"normal-c": {
		ID:   "normal-c",
		Name: "일반전 C",
		Kind: EncounterKindNormal,
		Enemies: []UnitState{
			enemyTemplate("shell-bug", "껍질벌레", 20,
				pattern("bug-shell", "껍질", nil, 4),
				pattern("bug-sting", "찌르기", []int{3}, 0),
				pattern("bug-shell-2", "껍질", nil, 3),
			),
			enemyTemplate("bat-seer", "박쥐예언자", 12,
				pattern("bat-peck", "쪼기", []int{2}, 0),
				pattern("bat-veil", "장막", nil, 2),
				pattern("bat-curse", "저주", []int{3}, 0),
			),
		},
	},
	"elite-1": {
		ID:   "elite-1",
		Name: "엘리트",
		Kind: EncounterKindElite,
		Enemies: []UnitState{
			enemyTemplate("guardian-idol", "수호우상", 24,
				pattern("idol-guard", "수비", nil, 5),
				pattern("idol-slam", "강타", []int{3}, 0),
				pattern("idol-wall", "장벽", nil, 5),
			),
			enemyTemplate("blade-imp", "칼날임프", 16,
				pattern("imp-slice", "베기", []int{4}, 0),
				pattern("imp-flurry", "난타", []int{2, 2}, 0),
				pattern("imp-guard", "수비", nil, 1),
			),
		},
	},
	"boss-1": {
		ID:   "boss-1",
		Name: "보스",
		Kind: EncounterKindBoss,
		Enemies: []UnitState{
			enemyTemplate("ruin-ogre", "폐허 오우거", 58,
				pattern("ogre-crush", "분쇄", []int{4, 4}, 0),
				pattern("ogre-hide", "숨기", nil, 6),
				pattern("ogre-roar", "포효", []int{7}, 0),
			),
		},
	},
}

var actMapNodes = map[string]EncounterNode{
	"start": {
		ID:      "start",
		Name:    "시작",
		Kind:    NodeKindStart,
		NextIDs: []string{"normal-a", "normal-b"},
	},
	"normal-a": {
		ID:          "normal-a",
		Name:        "갈림길 A",
		Kind:        NodeKindNormal,
		EncounterID: "normal-a",
		NextIDs:     []string{"rest-1", "normal-c"},
	},
	"normal-b": {
		ID:          "normal-b",
		Name:        "갈림길 B",
		Kind:        NodeKindNormal,
		EncounterID: "normal-b",
		NextIDs:     []string{"rest-1", "normal-c"},
	},
	"normal-c": {
		ID:          "normal-c",
		Name:        "매복",
		Kind:        NodeKindNormal,
		EncounterID: "normal-c",
		NextIDs:     []string{"elite-1"},
	},
	"rest-1": {
		ID:      "rest-1",
		Name:    "야영지",
		Kind:    NodeKindRest,
		NextIDs: []string{"elite-1"},
	},
	"elite-1": {
		ID:          "elite-1",
		Name:        "엘리트 관문",
		Kind:        NodeKindElite,
		EncounterID: "elite-1",
		NextIDs:     []string{"rest-2"},
	},
	"rest-2": {
		ID:      "rest-2",
		Name:    "성소",
		Kind:    NodeKindRest,
		NextIDs: []string{"boss-1"},
	},
	"boss-1": {
		ID:          "boss-1",
		Name:        "보스",
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
		Role:     "적",
		MaxHP:    hp,
		HP:       hp,
		Patterns: append([]EncounterPattern(nil), patterns...),
		Counters: map[string]int{},
	}
}
