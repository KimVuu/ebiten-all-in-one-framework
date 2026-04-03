package app

import (
	"strings"
	"testing"
)

func TestCombatSelectionForcesRemainingDiceThenRefills(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-warrior", "human-guard", "dwarf-smith"),
		testEncounter("normal-refill", EncounterKindNormal, idleEnemy("training-dummy", 40)),
		newRandomSourceWithScript(1),
	)

	allDice := append([]DieSpec(nil), combat.AvailableDice...)
	combat.AvailableDice = append([]DieSpec(nil), allDice[:2]...)
	combat.GraveyardDice = append([]DieSpec(nil), allDice[2:8]...)
	combat.SelectedDice = nil

	combat.prepareSelection()

	if got, want := len(combat.SelectedDice), 2; got != want {
		t.Fatalf("forced selected dice mismatch: got %d want %d", got, want)
	}
	for _, die := range combat.SelectedDice {
		if !die.Forced {
			t.Fatalf("expected forced selection for %s", die.Die.ID)
		}
	}
	if got, want := len(combat.AvailableDice), 6; got != want {
		t.Fatalf("available dice after refill mismatch: got %d want %d", got, want)
	}
	if got := len(combat.GraveyardDice); got != 0 {
		t.Fatalf("expected graveyard to be emptied after refill, got %d", got)
	}

	manualPick := combat.AvailableDice[0].ID
	if err := combat.selectDie(manualPick); err != nil {
		t.Fatalf("selectDie failed: %v", err)
	}
	if got, want := len(combat.SelectedDice), 3; got != want {
		t.Fatalf("selected dice mismatch after manual pick: got %d want %d", got, want)
	}
	if got, want := len(combat.AvailableDice), 5; got != want {
		t.Fatalf("available dice mismatch after manual pick: got %d want %d", got, want)
	}
}

func TestCombatDownedUnitDiceRemovedAndRevivedForNextBattle(t *testing.T) {
	party := mustPartyUnits("human-warrior", "human-guard", "human-priest")
	combat := newCombatStateWithRandom(
		party,
		testEncounter("normal-knockout", EncounterKindNormal, idleEnemy("training-dummy", 40)),
		newRandomSourceWithScript(1),
	)

	warriorDice := ownerDiceIDs(combat.AvailableDice, "human-warrior")
	if len(warriorDice) == 0 {
		t.Fatalf("expected warrior dice in available pool")
	}
	combat.SelectedDice = append(combat.SelectedDice, SelectedDie{Die: combat.AvailableDice[0], Forced: false})
	combat.GraveyardDice = append(combat.GraveyardDice, combat.AvailableDice[1])

	combat.applyDamageToUnit(sidePlayer, "human-warrior", 999)

	if unit := combat.playerUnit("human-warrior"); unit == nil || !unit.Downed || unit.HP != 0 {
		t.Fatalf("expected warrior to be downed, got %#v", unit)
	}
	if poolContainsOwner(combat.AvailableDice, "human-warrior") || poolContainsOwner(selectedToPool(combat.SelectedDice), "human-warrior") || poolContainsOwner(combat.GraveyardDice, "human-warrior") {
		t.Fatalf("expected warrior dice to be removed from all pools")
	}

	run := newRunState(42)
	run.SelectedPartyIDs = []string{"human-warrior", "human-guard", "human-priest"}
	run.startRun()
	run.PartyUnits[0].HP = 0
	run.PartyUnits[0].Downed = true

	nextCombat, err := run.startEncounterForNode("normal-a")
	if err != nil {
		t.Fatalf("startEncounterForNode failed: %v", err)
	}
	warrior := nextCombat.playerUnit("human-warrior")
	if warrior == nil {
		t.Fatalf("expected revived warrior in next combat")
	}
	if got, want := warrior.HP, 1; got != want || warrior.Downed {
		t.Fatalf("expected warrior revived at 1 HP, got hp=%d downed=%v", warrior.HP, warrior.Downed)
	}
	if !poolContainsOwner(nextCombat.AvailableDice, "human-warrior") {
		t.Fatalf("expected warrior dice to be restored for next battle")
	}
}

func TestCombatHeroGoddessStackExplodesAtThree(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-warrior", "human-guard", "human-priest"),
		testEncounter("normal-aoe", EncounterKindNormal,
			idleEnemy("slime-a", 10),
			idleEnemy("slime-b", 10),
		),
		newRandomSourceWithScript(1, 0, 0, 0),
	)
	combat.playerUnit("human-warrior").Counters[counterHeroGoddess] = 2

	mustSelectDice(t, combat,
		"human-warrior-goddess-1",
		"human-guard-defense-1",
		"human-guard-defense-2",
	)

	combat.resolveTurn()

	if got := combat.playerUnit("human-warrior").Counters[counterHeroGoddess]; got != 0 {
		t.Fatalf("expected goddess stack to reset, got %d", got)
	}
	if enemy := combat.enemyUnit("slime-a"); enemy == nil || !enemy.Downed {
		t.Fatalf("expected slime-a to be defeated, got %#v", enemy)
	}
	if enemy := combat.enemyUnit("slime-b"); enemy == nil || !enemy.Downed {
		t.Fatalf("expected slime-b to be defeated, got %#v", enemy)
	}
}

func TestCombatShieldDefenseMultiplierAppliesAfterSummedDefense(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-guard", "human-warrior", "human-priest"),
		testEncounter("normal-defense", EncounterKindNormal, idleEnemy("training-dummy", 40)),
		newRandomSourceWithScript(1, 0, 1, 0),
	)

	mustSelectDice(t, combat,
		"human-guard-defense-1",
		"human-guard-defense-2",
		"human-guard-tank-1",
	)

	summary := combat.resolveTurn()
	if got, want := summary.GeneratedAllyDefense, 6; got != want {
		t.Fatalf("ally defense mismatch: got %d want %d", got, want)
	}
}

func TestCombatPriestHealCapsAtMaxAndDoesNotRevive(t *testing.T) {
	party := mustPartyUnits("human-priest", "human-warrior", "human-guard")
	findPartyUnit(party, "human-warrior").HP = 10
	findPartyUnit(party, "human-guard").HP = 0
	findPartyUnit(party, "human-guard").Downed = true

	combat := newCombatStateWithRandom(
		party,
		testEncounter("normal-heal", EncounterKindNormal, idleEnemy("training-dummy", 40)),
		newRandomSourceWithScript(1, 0, 0, 0),
	)

	mustSelectDice(t, combat,
		"human-priest-priest-1",
		"human-warrior-attack-1",
		"human-warrior-defense-1",
	)

	combat.resolveTurn()

	if got, want := combat.playerUnit("human-warrior").HP, 12; got != want {
		t.Fatalf("warrior HP mismatch: got %d want %d", got, want)
	}
	guard := combat.playerUnit("human-guard")
	if guard == nil || !guard.Downed || guard.HP != 0 {
		t.Fatalf("expected guard to stay downed without revive, got %#v", guard)
	}
}

func TestCombatGuideRevealAndEscapeRules(t *testing.T) {
	revealCombat := newCombatStateWithRandom(
		mustPartyUnits("human-guide", "human-warrior", "human-guard"),
		testEncounter("normal-reveal", EncounterKindNormal, enemyWithPatterns("seer", 30,
			EncounterPattern{ID: "swing", Label: "Swing", Attacks: []int{3}},
			EncounterPattern{ID: "brace", Label: "Brace", Defense: 4},
		)),
		newRandomSourceWithScript(1, 0, 0, 0, 0),
	)

	mustSelectDice(t, revealCombat,
		"human-guide-info-1",
		"human-warrior-defense-1",
		"human-guard-defense-1",
	)
	summary := revealCombat.resolveTurn()
	if got, want := summary.RevealedNextPatterns["seer"], "Brace"; got != want {
		t.Fatalf("revealed next pattern mismatch: got %q want %q", got, want)
	}

	escapeCombat := newCombatStateWithRandom(
		mustPartyUnits("human-guide", "human-warrior", "human-guard"),
		testEncounter("normal-escape", EncounterKindNormal, idleEnemy("runner", 30)),
		newRandomSourceWithScript(1, 0, 5, 0),
	)
	mustSelectDice(t, escapeCombat,
		"human-guide-escape-1",
		"human-guide-info-1",
		"human-warrior-defense-1",
	)
	if got := escapeCombat.resolveTurn().Outcome; got != CombatOutcomeEscape {
		t.Fatalf("expected normal encounter escape, got %q", got)
	}

	bossCombat := newCombatStateWithRandom(
		mustPartyUnits("human-guide", "human-warrior", "human-guard"),
		testEncounter("boss-no-escape", EncounterKindBoss, idleEnemy("boss", 60)),
		newRandomSourceWithScript(1, 0, 5, 0, 0),
	)
	mustSelectDice(t, bossCombat,
		"human-guide-escape-1",
		"human-guide-info-1",
		"human-warrior-defense-1",
	)
	if got := bossCombat.resolveTurn().Outcome; got == CombatOutcomeEscape {
		t.Fatalf("expected boss encounter to reject escape")
	}
}

func TestCombatDamageBoostsStackAndFloorAcrossTwoPlayerTurns(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-guide", "dwarf-smith", "human-warrior"),
		testEncounter("normal-boost", EncounterKindNormal, idleEnemy("dummy", 30)),
		newRandomSourceWithScript(1,
			0, 0, 5, 0,
			3, 5, 0, 0,
		),
	)

	mustSelectDice(t, combat,
		"human-guide-weakness-1",
		"dwarf-smith-forge-1",
		"human-warrior-attack-1",
	)
	combat.resolveTurn()
	if got, want := combat.enemyUnit("dummy").HP, 21; got != want {
		t.Fatalf("enemy HP mismatch after stacked boost turn: got %d want %d", got, want)
	}

	mustSelectDice(t, combat,
		"human-guide-info-1",
		"human-warrior-attack-2",
		"human-warrior-defense-1",
	)
	combat.resolveTurn()
	if got, want := combat.enemyUnit("dummy").HP, 15; got != want {
		t.Fatalf("enemy HP mismatch after second player turn: got %d want %d", got, want)
	}
}

func TestCombatArcherCriticalDoublesOnlyNumericResults(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("elf-archer", "human-priest", "human-guard"),
		testEncounter("normal-archer", EncounterKindNormal, idleEnemy("dummy", 20)),
		newRandomSourceWithScript(1,
			3, 4, 5,
			4, 5, 5,
		),
	)

	mustSelectDice(t, combat,
		"elf-archer-shot-1",
		"elf-archer-shot-2",
		"human-priest-priest-1",
	)
	combat.resolveTurn()
	if got, want := combat.enemyUnit("dummy").HP, 12; got != want {
		t.Fatalf("enemy HP mismatch after crit turn: got %d want %d", got, want)
	}

	mustSelectDice(t, combat,
		"elf-archer-shot-3",
		"elf-archer-shot-4",
		"human-priest-priest-2",
	)
	combat.resolveTurn()
	if got, want := combat.enemyUnit("dummy").HP, 12; got != want {
		t.Fatalf("enemy HP should stay unchanged on crit-only turn: got %d want %d", got, want)
	}
}

func TestCombatRandomTargetingSkipsDownedUnitsAndUsesDeterministicSeed(t *testing.T) {
	party := mustPartyUnits("human-warrior", "human-guard", "human-priest")
	findPartyUnit(party, "human-priest").HP = 0
	findPartyUnit(party, "human-priest").Downed = true

	encounter := testEncounter("normal-target", EncounterKindNormal, enemyWithPatterns("raider", 40,
		EncounterPattern{ID: "flurry", Label: "Flurry", Attacks: []int{3, 3, 3}},
	))

	first := newCombatStateWithRandom(cloneUnits(party), encounter, NewRandomSource(99))
	second := newCombatStateWithRandom(cloneUnits(party), encounter, NewRandomSource(99))

	mustSelectDice(t, first, "human-warrior-attack-1", "human-warrior-attack-2", "human-guard-tank-1")
	mustSelectDice(t, second, "human-warrior-attack-1", "human-warrior-attack-2", "human-guard-tank-1")

	firstSummary := first.resolveTurn()
	secondSummary := second.resolveTurn()

	if joinStrings(firstSummary.EnemyTargets, ",") != joinStrings(secondSummary.EnemyTargets, ",") {
		t.Fatalf("expected deterministic enemy targets, got %v and %v", firstSummary.EnemyTargets, secondSummary.EnemyTargets)
	}
	for _, target := range firstSummary.EnemyTargets {
		if target == "human-priest" {
			t.Fatalf("expected downed priest to be excluded from enemy targeting")
		}
	}
}

func TestCombatEnemyActionsAppearInLogs(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-guard", "human-warrior", "human-priest"),
		testEncounter("normal-enemy-log", EncounterKindNormal, enemyWithPatterns("raider", 40,
			EncounterPattern{ID: "combo", Label: "Combo", Attacks: []int{3}, Defense: 4},
		)),
		newRandomSourceWithScript(0, 0, 1, 0),
	)

	mustSelectDice(t, combat,
		"human-guard-defense-1",
		"human-warrior-defense-1",
		"human-priest-priest-1",
	)

	summary := combat.resolveTurn()
	logs := joinStrings(summary.Logs, " | ")
	if !strings.Contains(logs, "raider 행동 수비 4, 적 방어막이 4가 되었다.") {
		t.Fatalf("expected enemy defense log, got %q", logs)
	}
	if !strings.Contains(logs, "raider 행동 공격 3, 파티 방어막이 2 막았고, 인간 방패병에게 1 피해를 입혔다.") {
		t.Fatalf("expected enemy action log, got %q", logs)
	}
	if strings.Index(logs, "raider 행동 수비 4, 적 방어막이 4가 되었다.") > strings.Index(logs, "raider 행동 공격 3, 파티 방어막이 2 막았고, 인간 방패병에게 1 피해를 입혔다.") {
		t.Fatalf("expected defense log before attack log, got %q", logs)
	}
}

func TestCombatPlayerAttackLogsTargetAndShieldAbsorb(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-warrior", "human-guard", "human-priest"),
		testEncounter("normal-player-log", EncounterKindNormal, enemyWithPatterns("dummy", 20,
			EncounterPattern{ID: "brace", Label: "Brace", Defense: 2},
		)),
		newRandomSourceWithScript(1, 3, 0, 1, 0),
	)

	mustSelectDice(t, combat,
		"human-warrior-attack-1",
		"human-guard-defense-1",
		"human-priest-priest-1",
	)

	summary := combat.resolveTurn()
	logs := joinStrings(summary.Logs, " | ")
	if !strings.Contains(logs, "인간 용사의 공격 주사위 1 결과 3, 적 방어막이 2 막았고, dummy에게 1 피해를 입혔다.") {
		t.Fatalf("expected player target and shield log, got %q", logs)
	}
}

func TestCombatLogsFollowDefenseSkillAttackOrder(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-warrior", "human-guard", "human-priest"),
		testEncounter("normal-order", EncounterKindNormal, enemyWithPatterns("raider", 20,
			EncounterPattern{ID: "brace-hit", Label: "Brace Hit", Attacks: []int{3}, Defense: 2},
		)),
		newRandomSourceWithScript(1, 3, 1, 0, 0, 0),
	)

	mustSelectDice(t, combat,
		"human-warrior-attack-1",
		"human-guard-defense-1",
		"human-priest-priest-1",
	)

	summary := combat.resolveTurn()
	logs := joinStrings(summary.Logs, " | ")
	playerDefense := "인간 방패병의 방어 주사위 1 결과 2, 파티 방어막이 2가 되었다."
	enemyDefense := "raider 행동 수비 2, 적 방어막이 2가 되었다."
	skill := "인간 여신관의 신관의 주사위 1 결과 성공, 생존한 아군 전체가 2 회복되었다."
	playerAttack := "인간 용사의 공격 주사위 1 결과 3, 적 방어막이 2 막았고, raider에게 1 피해를 입혔다."
	enemyAttack := "raider 행동 공격 3, 파티 방어막이 2 막았고, 인간 용사에게 1 피해를 입혔다."

	if strings.Index(logs, playerDefense) == -1 || strings.Index(logs, enemyDefense) == -1 || strings.Index(logs, skill) == -1 || strings.Index(logs, playerAttack) == -1 || strings.Index(logs, enemyAttack) == -1 {
		t.Fatalf("expected ordered logs, got %q", logs)
	}
	if strings.Index(logs, playerDefense) > strings.Index(logs, skill) {
		t.Fatalf("expected player defense before skill, got %q", logs)
	}
	if strings.Index(logs, enemyDefense) > strings.Index(logs, skill) {
		t.Fatalf("expected enemy defense before skill, got %q", logs)
	}
	if strings.Index(logs, skill) > strings.Index(logs, playerAttack) {
		t.Fatalf("expected skill before player attack, got %q", logs)
	}
	if strings.Index(logs, playerAttack) > strings.Index(logs, enemyAttack) {
		t.Fatalf("expected player attack before enemy attack, got %q", logs)
	}
}

func TestCombatDeadEnemyDoesNotActAfterBeingDefeated(t *testing.T) {
	combat := newCombatStateWithRandom(
		mustPartyUnits("human-warrior", "human-guard", "human-priest"),
		testEncounter("normal-dead-enemy", EncounterKindNormal, enemyWithPatterns("raider", 1,
			EncounterPattern{ID: "hit", Label: "Hit", Attacks: []int{3}},
		)),
		newRandomSourceWithScript(1, 3, 0, 1, 0),
	)

	mustSelectDice(t, combat,
		"human-warrior-attack-1",
		"human-guard-defense-1",
		"human-priest-priest-1",
	)

	summary := combat.resolveTurn()
	logs := joinStrings(summary.Logs, " | ")
	if strings.Contains(logs, "raider 행동 공격 3") {
		t.Fatalf("dead enemy should not act, got %q", logs)
	}
	if got := combat.playerUnit("human-warrior").HP; got != 24 {
		t.Fatalf("expected warrior HP unchanged, got %d", got)
	}
	if summary.Outcome != CombatOutcomeVictory {
		t.Fatalf("expected victory after defeating enemy before attack, got %q", summary.Outcome)
	}
}

func mustPartyUnits(ids ...string) []UnitState {
	units := make([]UnitState, 0, len(ids))
	for _, id := range ids {
		unit, ok := newCharacterState(id)
		if !ok {
			panic("unknown character: " + id)
		}
		units = append(units, unit)
	}
	return units
}

func mustSelectDice(t *testing.T, combat *CombatState, ids ...string) {
	t.Helper()
	for _, id := range ids {
		if err := combat.selectDie(id); err != nil {
			t.Fatalf("selectDie(%q) failed: %v", id, err)
		}
	}
}

func testEncounter(id string, kind EncounterKind, enemies ...UnitState) EncounterDefinition {
	return EncounterDefinition{
		ID:      id,
		Name:    id,
		Kind:    kind,
		Enemies: cloneUnits(enemies),
	}
}

func idleEnemy(id string, hp int) UnitState {
	return enemyWithPatterns(id, hp, EncounterPattern{ID: "idle", Label: "Idle"})
}

func enemyWithPatterns(id string, hp int, patterns ...EncounterPattern) UnitState {
	return UnitState{
		ID:         id,
		Name:       id,
		Role:       "enemy",
		MaxHP:      hp,
		HP:         hp,
		Patterns:   append([]EncounterPattern(nil), patterns...),
		Counters:   map[string]int{},
		Statuses:   nil,
		Dice:       nil,
		Downed:     false,
		PatternIdx: 0,
	}
}

func findPartyUnit(units []UnitState, id string) *UnitState {
	for idx := range units {
		if units[idx].ID == id {
			return &units[idx]
		}
	}
	return nil
}
