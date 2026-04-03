package app

import (
	"fmt"
	"math"
	"strings"
)

type attackPacket struct {
	SourceUnitID string
	SourceDieID  string
	Value        int
	TargetID     string
}

type rolledDieResult struct {
	Selected SelectedDie
	Face     DieFace
}

type damageResolution struct {
	SourceUnitID string
	SourceDieID  string
	TargetID     string
	BaseDamage   int
	Absorbed     int
	DamageDealt  int
}

type enemyTurnAction struct {
	EnemyID   string
	EnemyName string
	Pattern   EncounterPattern
}

type queuedPlayerAttack struct {
	Packet    attackPacket
	LogPrefix string
}

func newCombatStateWithRandom(party []UnitState, encounter EncounterDefinition, random *RandomSource) *CombatState {
	combat := &CombatState{
		NodeID:               encounter.ID,
		EncounterID:          encounter.ID,
		EncounterName:        encounter.Name,
		EncounterKind:        encounter.Kind,
		PlayerUnits:          cloneUnits(party),
		EnemyUnits:           cloneUnits(encounter.Enemies),
		AvailableDice:        buildCombatDicePool(party),
		GraveyardDice:        nil,
		SelectedDice:         nil,
		EnemyStatuses:        nil,
		RevealedNextPatterns: nil,
		Outcome:              CombatOutcomeNone,
		Logs:                 nil,
		random:               random,
	}
	combat.prepareSelection()
	return combat
}

func newCombatState(seed int64, party []UnitState, encounter EncounterDefinition) *CombatState {
	return newCombatStateWithRandom(party, encounter, NewRandomSource(seed))
}

func (combat *CombatState) playerUnit(id string) *UnitState {
	for idx := range combat.PlayerUnits {
		if combat.PlayerUnits[idx].ID == id {
			return &combat.PlayerUnits[idx]
		}
	}
	return nil
}

func (combat *CombatState) enemyUnit(id string) *UnitState {
	for idx := range combat.EnemyUnits {
		if combat.EnemyUnits[idx].ID == id {
			return &combat.EnemyUnits[idx]
		}
	}
	return nil
}

func (combat *CombatState) alivePlayerUnits() []int {
	indices := make([]int, 0)
	for idx := range combat.PlayerUnits {
		if combat.PlayerUnits[idx].Downed || combat.PlayerUnits[idx].HP <= 0 {
			continue
		}
		indices = append(indices, idx)
	}
	return indices
}

func (combat *CombatState) aliveEnemyUnits() []int {
	indices := make([]int, 0)
	for idx := range combat.EnemyUnits {
		if combat.EnemyUnits[idx].Downed || combat.EnemyUnits[idx].HP <= 0 {
			continue
		}
		indices = append(indices, idx)
	}
	return indices
}

func (combat *CombatState) prepareSelection() {
	if combat == nil || combat.Outcome != CombatOutcomeNone {
		return
	}
	for len(combat.SelectedDice) < 3 {
		remainingSlots := 3 - len(combat.SelectedDice)
		if len(combat.AvailableDice) == 0 {
			if len(combat.GraveyardDice) == 0 {
				return
			}
			combat.AvailableDice = append(combat.AvailableDice, combat.GraveyardDice...)
			combat.GraveyardDice = nil
		}
		if len(combat.AvailableDice) <= remainingSlots {
			for _, die := range combat.AvailableDice {
				combat.SelectedDice = append(combat.SelectedDice, SelectedDie{
					Die:    die,
					Forced: true,
				})
			}
			combat.AvailableDice = nil
			continue
		}
		return
	}
}

func (combat *CombatState) selectDie(id string) error {
	if combat == nil {
		return fmt.Errorf("전투 상태를 찾을 수 없습니다")
	}
	if combat.Outcome != CombatOutcomeNone {
		return fmt.Errorf("이미 전투가 종료되었습니다")
	}
	combat.prepareSelection()
	if len(combat.SelectedDice) >= 3 {
		return fmt.Errorf("이번 턴에는 이미 주사위 3개가 선택되었습니다")
	}
	for idx, die := range combat.AvailableDice {
		if die.ID != id {
			continue
		}
		combat.SelectedDice = append(combat.SelectedDice, SelectedDie{Die: die})
		combat.AvailableDice = append(combat.AvailableDice[:idx], combat.AvailableDice[idx+1:]...)
		combat.prepareSelection()
		return nil
	}
	return fmt.Errorf("알 수 없는 주사위입니다: %s", id)
}

func (combat *CombatState) resolveTurn() TurnResolution {
	summary := TurnResolution{
		Turn:    combat.Turn + 1,
		Outcome: CombatOutcomeNone,
	}
	if combat == nil {
		return summary
	}
	combat.prepareSelection()
	if len(combat.SelectedDice) != 3 {
		summary.Logs = append(summary.Logs, "주사위 3개가 선택되어야 턴을 진행할 수 있습니다.")
		return summary
	}

	hadReveal := len(combat.RevealedNextPatterns) > 0
	combat.Turn++
	guideReveal := false
	combat.pushLog(fmt.Sprintf("---턴%d---", combat.Turn))

	finishTurn := func() TurnResolution {
		if guideReveal {
			combat.RevealedNextPatterns = combat.nextPatternLabels()
		} else if hadReveal {
			combat.RevealedNextPatterns = nil
		}
		summary.RevealedNextPatterns = cloneStringMap(combat.RevealedNextPatterns)
		summary.Outcome = combat.Outcome
		summary.Logs = append(summary.Logs, combat.Logs...)
		combat.AllyDefense = 0
		combat.EnemyDefense = 0
		if combat.Outcome == CombatOutcomeNone {
			combat.prepareSelection()
		}
		return summary
	}

	rolled := make([]rolledDieResult, 0, len(combat.SelectedDice))
	for _, selected := range combat.SelectedDice {
		die := selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) {
			continue
		}
		face := die.Faces[combat.random.NextInt(len(die.Faces))]
		rolled = append(rolled, rolledDieResult{
			Selected: selected,
			Face:     face,
		})
	}
	for _, selected := range combat.SelectedDice {
		combat.GraveyardDice = append(combat.GraveyardDice, selected.Die)
	}
	combat.SelectedDice = nil

	enemyActions := combat.planEnemyTurnActions()

	allyDefense := 0
	guideEscapeMain := false
	guideEscapeSupport := false
	smithBoost := 0
	queuedSkillAttacks := make([]queuedPlayerAttack, 0)
	enemyDefense := 0

	for _, result := range rolled {
		die := result.Selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) || die.Kind != DieKindDefense {
			continue
		}
		if result.Face.Kind == FaceKindValue && result.Face.Value > 0 {
			allyDefense += result.Face.Value
		}
		combat.pushLog(playerDefenseLogLine(die, result.Face, allyDefense))
	}
	for _, result := range rolled {
		die := result.Selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) || die.EffectID != effectTankGuard {
			continue
		}
		if result.Face.Kind == FaceKindSuccess {
			allyDefense *= 2
		}
		combat.pushLog(playerTankGuardLogLine(die, result.Face, allyDefense))
	}
	for _, action := range enemyActions {
		if action.Pattern.Defense <= 0 {
			continue
		}
		enemyDefense += action.Pattern.Defense
		combat.pushLog(enemyDefenseLogLine(action.EnemyName, action.Pattern.Defense, enemyDefense))
	}
	combat.AllyDefense = allyDefense
	summary.GeneratedAllyDefense = allyDefense
	combat.EnemyDefense = enemyDefense
	summary.GeneratedEnemyDefense = enemyDefense

	for _, result := range rolled {
		die := result.Selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) {
			continue
		}
		switch die.EffectID {
		case effectHeroGoddess:
			if result.Face.Kind == FaceKindSuccess {
				unit := combat.playerUnit(die.OwnerID)
				unit.Counters[counterHeroGoddess]++
				line := fmt.Sprintf("%s 결과 %s, %s의 여신 스택이 %d가 되었다.", dieLogLabel(die), result.Face.Label, ownerLabel(die.OwnerID), unit.Counters[counterHeroGoddess])
				triggered := false
				for unit.Counters[counterHeroGoddess] >= 3 {
					triggered = true
					unit.Counters[counterHeroGoddess] -= 3
					for _, enemyIdx := range combat.aliveEnemyUnits() {
						queuedSkillAttacks = append(queuedSkillAttacks, queuedPlayerAttack{
							Packet: attackPacket{
								SourceUnitID: die.OwnerID,
								SourceDieID:  die.ID,
								Value:        10,
								TargetID:     combat.EnemyUnits[enemyIdx].ID,
							},
							LogPrefix: fmt.Sprintf("%s의 여신 폭발", ownerLabel(die.OwnerID)),
						})
					}
				}
				if triggered {
					line += " 여신 폭발이 발동했다."
				}
				combat.pushLog(line)
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "아무 일도 일어나지 않았다."))
		case effectPriestHeal:
			if result.Face.Kind == FaceKindSuccess {
				for idx := range combat.PlayerUnits {
					unit := &combat.PlayerUnits[idx]
					if unit.Downed || unit.HP <= 0 {
						continue
					}
					unit.HP = minInt(unit.MaxHP, unit.HP+2)
				}
				combat.pushLog(skillResultLogLine(die, result.Face, "생존한 아군 전체가 2 회복되었다."))
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "회복은 발생하지 않았다."))
		case effectGuideInfo:
			if result.Face.Kind == FaceKindSuccess {
				guideReveal = true
				combat.pushLog(skillResultLogLine(die, result.Face, "다음 적 패턴을 간파했다."))
				continue
			}
			if result.Face.Kind == FaceKindEscape {
				guideEscapeSupport = true
				combat.pushLog(skillResultLogLine(die, result.Face, "도주 지원이 준비되었다."))
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "정보를 얻지 못했다."))
		case effectGuideWeakness:
			if result.Face.Kind == FaceKindSuccess {
				combat.addEnemyStatus(StatusEffect{
					ID:             statusWeakness,
					Name:           "약점",
					Magnitude:      30,
					RemainingTurns: 2,
				})
				combat.pushLog(skillResultLogLine(die, result.Face, "적이 받는 피해가 2턴 동안 30% 증가한다."))
				continue
			}
			if result.Face.Kind == FaceKindEscape {
				guideEscapeSupport = true
				combat.pushLog(skillResultLogLine(die, result.Face, "도주 지원이 준비되었다."))
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "약점을 찾지 못했다."))
		case effectGuideEscape:
			if result.Face.Kind == FaceKindEscape {
				guideEscapeMain = true
				combat.pushLog(skillResultLogLine(die, result.Face, "도주 시도를 시작했다."))
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "도주 시도에 실패했다."))
		case effectSmithForge:
			if result.Face.Kind == FaceKindSuccess {
				smithBoost += 50
				combat.pushLog(skillResultLogLine(die, result.Face, "이번 턴 적이 받는 피해가 50% 증가한다."))
				continue
			}
			combat.pushLog(skillResultLogLine(die, result.Face, "피해 증폭은 발생하지 않았다."))
		}
	}

	damageBoost := smithBoost + combat.enemyDamageBoostPercent()
	summary.DamageBoostPercent = damageBoost

	if guideEscapeMain && guideEscapeSupport && combat.EncounterKind == EncounterKindNormal {
		combat.Outcome = CombatOutcomeEscape
		combat.pushLog("도주에 성공했다.")
		return finishTurn()
	}

	archerCrits := 0
	for _, result := range rolled {
		die := result.Selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) || die.EffectID != effectArcherShot {
			continue
		}
		if result.Face.Kind == FaceKindCritical {
			archerCrits++
		}
	}

	playerTargets := make([]string, 0)
	for _, result := range rolled {
		die := result.Selected.Die
		if !combat.isUnitActive(sidePlayer, die.OwnerID) {
			continue
		}
		switch {
		case die.Kind == DieKindAttack:
			resolution := combat.applyPacketToEnemy(attackPacket{
				SourceUnitID: die.OwnerID,
				SourceDieID:  die.ID,
				Value:        result.Face.Value,
			}, damageBoost)
			if resolution.TargetID != "" {
				playerTargets = append(playerTargets, resolution.TargetID)
			}
			combat.pushLog(playerDieAttackLogLine(die, result.Face, resolution, false))
		case die.EffectID == effectArcherShot:
			if result.Face.Kind == FaceKindCritical {
				combat.pushLog(archerCriticalLogLine(die))
				continue
			}
			value := result.Face.Value
			if archerCrits > 0 {
				value *= 2
			}
			resolution := combat.applyPacketToEnemy(attackPacket{
				SourceUnitID: die.OwnerID,
				SourceDieID:  die.ID,
				Value:        value,
			}, damageBoost)
			if resolution.TargetID != "" {
				playerTargets = append(playerTargets, resolution.TargetID)
			}
			combat.pushLog(playerDieAttackLogLine(die, result.Face, resolution, archerCrits > 0))
		}
		if len(combat.aliveEnemyUnits()) == 0 {
			break
		}
	}
	for _, queued := range queuedSkillAttacks {
		if !combat.isUnitActive(sidePlayer, queued.Packet.SourceUnitID) {
			continue
		}
		resolution := combat.applyPacketToEnemy(queued.Packet, damageBoost)
		if resolution.TargetID != "" {
			playerTargets = append(playerTargets, resolution.TargetID)
		}
		combat.pushLog(skillAttackLogLine(queued.LogPrefix, resolution))
		if len(combat.aliveEnemyUnits()) == 0 {
			break
		}
	}
	summary.PlayerTargets = playerTargets
	combat.tickEnemyStatusesAfterPlayerPhase()

	if len(combat.aliveEnemyUnits()) == 0 {
		combat.Outcome = CombatOutcomeVictory
		combat.pushLog("전투를 돌파했다.")
		return finishTurn()
	}

	enemyTargets := make([]string, 0)
	for _, action := range enemyActions {
		if !combat.isUnitActive(sideEnemy, action.EnemyID) {
			continue
		}
		for _, value := range action.Pattern.Attacks {
			resolution := combat.applyPacketToPlayers(attackPacket{
				SourceUnitID: action.EnemyID,
				Value:        value,
			})
			if resolution.TargetID != "" {
				enemyTargets = append(enemyTargets, resolution.TargetID)
			}
			combat.pushLog(enemyAttackLogLine(action.EnemyName, value, resolution))
			if len(combat.alivePlayerUnits()) == 0 {
				break
			}
		}
		if len(combat.alivePlayerUnits()) == 0 {
			break
		}
	}
	summary.EnemyTargets = enemyTargets

	if len(combat.alivePlayerUnits()) == 0 {
		combat.Outcome = CombatOutcomeDefeat
		combat.pushLog("파티가 전멸했다.")
	}
	return finishTurn()
}

func (combat *CombatState) pushLog(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}
	combat.Logs = append(combat.Logs, line)
}

func dieLogLabel(die DieSpec) string {
	return fmt.Sprintf("%s의 %s", ownerLabel(die.OwnerID), die.Name)
}

func playerDefenseLogLine(die DieSpec, face DieFace, total int) string {
	return fmt.Sprintf("%s 결과 %s, 파티 방어막이 %d가 되었다.", dieLogLabel(die), face.Label, total)
}

func playerTankGuardLogLine(die DieSpec, face DieFace, total int) string {
	switch face.Kind {
	case FaceKindSuccess:
		return fmt.Sprintf("%s 결과 %s, 파티 방어막이 %d가 되었다.", dieLogLabel(die), face.Label, total)
	default:
		return fmt.Sprintf("%s 결과 %s, 파티 방어막은 %d다.", dieLogLabel(die), face.Label, total)
	}
}

func enemyDefenseLogLine(enemyName string, value int, total int) string {
	return fmt.Sprintf("%s 행동 수비 %d, 적 방어막이 %d가 되었다.", enemyName, value, total)
}

func skillResultLogLine(die DieSpec, face DieFace, detail string) string {
	if strings.TrimSpace(detail) == "" {
		return fmt.Sprintf("%s 결과 %s.", dieLogLabel(die), face.Label)
	}
	return fmt.Sprintf("%s 결과 %s, %s", dieLogLabel(die), face.Label, detail)
}

func playerDieAttackLogLine(die DieSpec, face DieFace, resolution damageResolution, critBoosted bool) string {
	label := fmt.Sprintf("%s 결과 %s", dieLogLabel(die), face.Label)
	if critBoosted && die.EffectID == effectArcherShot && face.Kind == FaceKindValue {
		label += "(치명타 적용)"
	}
	return fmt.Sprintf("%s, %s", label, playerAttackOutcomeText(resolution))
}

func archerCriticalLogLine(die DieSpec) string {
	return fmt.Sprintf("%s 결과 치명타, 이번 턴 엘프 궁수의 공격이 2배가 된다.", dieLogLabel(die))
}

func skillAttackLogLine(prefix string, resolution damageResolution) string {
	return fmt.Sprintf("%s, %s", prefix, playerAttackOutcomeText(resolution))
}

func playerAttackOutcomeText(resolution damageResolution) string {
	switch {
	case resolution.Absorbed > 0 && resolution.DamageDealt > 0:
		return fmt.Sprintf("적 방어막이 %d 막았고, %s에게 %d 피해를 입혔다.", resolution.Absorbed, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	case resolution.Absorbed > 0:
		return fmt.Sprintf("적 방어막이 %d 막았다.", resolution.Absorbed)
	case resolution.DamageDealt > 0:
		return fmt.Sprintf("%s에게 %d 피해를 입혔다.", ownerLabel(resolution.TargetID), resolution.DamageDealt)
	default:
		return "피해를 주지 못했다."
	}
}

func enemyAttackLogLine(enemyName string, value int, resolution damageResolution) string {
	switch {
	case resolution.Absorbed > 0 && resolution.DamageDealt > 0:
		return fmt.Sprintf("%s 행동 공격 %d, 파티 방어막이 %d 막았고, %s에게 %d 피해를 입혔다.", enemyName, value, resolution.Absorbed, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	case resolution.Absorbed > 0:
		return fmt.Sprintf("%s 행동 공격 %d, 파티 방어막이 %d 막았다.", enemyName, value, resolution.Absorbed)
	case resolution.DamageDealt > 0:
		return fmt.Sprintf("%s 행동 공격 %d, %s에게 %d 피해를 입혔다.", enemyName, value, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	default:
		return fmt.Sprintf("%s 행동 공격 %d, 피해를 주지 못했다.", enemyName, value)
	}
}

func (combat *CombatState) addEnemyStatus(status StatusEffect) {
	combat.EnemyStatuses = append(combat.EnemyStatuses, status)
}

func (combat *CombatState) enemyDamageBoostPercent() int {
	total := 0
	for _, status := range combat.EnemyStatuses {
		if status.ID == statusWeakness && status.RemainingTurns > 0 {
			total += status.Magnitude
		}
	}
	return total
}

func (combat *CombatState) tickEnemyStatusesAfterPlayerPhase() {
	filtered := make([]StatusEffect, 0, len(combat.EnemyStatuses))
	for _, status := range combat.EnemyStatuses {
		status.RemainingTurns--
		if status.RemainingTurns > 0 {
			filtered = append(filtered, status)
		}
	}
	combat.EnemyStatuses = filtered
}

func (combat *CombatState) planEnemyTurnActions() []enemyTurnAction {
	actions := make([]enemyTurnAction, 0)
	for _, enemyIdx := range combat.aliveEnemyUnits() {
		enemy := &combat.EnemyUnits[enemyIdx]
		if len(enemy.Patterns) == 0 {
			continue
		}
		pattern := enemy.Patterns[enemy.PatternIdx%len(enemy.Patterns)]
		enemy.PatternIdx++
		actions = append(actions, enemyTurnAction{
			EnemyID:   enemy.ID,
			EnemyName: enemy.Name,
			Pattern:   pattern,
		})
	}
	return actions
}

func (combat *CombatState) isUnitActive(side combatSide, unitID string) bool {
	var unit *UnitState
	switch side {
	case sidePlayer:
		unit = combat.playerUnit(unitID)
	case sideEnemy:
		unit = combat.enemyUnit(unitID)
	}
	return unit != nil && !unit.Downed && unit.HP > 0
}

func (combat *CombatState) applyPacketToEnemy(packet attackPacket, boostPercent int) damageResolution {
	damage := packet.Value
	if boostPercent > 0 {
		damage = int(math.Floor(float64(damage) * float64(100+boostPercent) / 100.0))
	}
	resolution := damageResolution{
		SourceUnitID: packet.SourceUnitID,
		SourceDieID:  packet.SourceDieID,
		BaseDamage:   damage,
	}
	if damage <= 0 {
		return resolution
	}
	if packet.TargetID != "" {
		if target := combat.resolveEnemyTarget(packet.TargetID); target != nil {
			resolution.TargetID = target.ID
		}
	}
	if combat.EnemyDefense > 0 {
		absorbed := minInt(combat.EnemyDefense, damage)
		combat.EnemyDefense -= absorbed
		damage -= absorbed
		resolution.Absorbed = absorbed
	}
	if damage <= 0 {
		return resolution
	}
	target := combat.resolveEnemyTarget(packet.TargetID)
	if target == nil {
		return resolution
	}
	combat.damageUnit(target, damage)
	resolution.TargetID = target.ID
	resolution.DamageDealt = damage
	return resolution
}

func (combat *CombatState) applyPacketToPlayers(packet attackPacket) damageResolution {
	damage := packet.Value
	resolution := damageResolution{
		SourceUnitID: packet.SourceUnitID,
		SourceDieID:  packet.SourceDieID,
		BaseDamage:   damage,
	}
	if damage <= 0 {
		return resolution
	}
	if combat.AllyDefense > 0 {
		absorbed := minInt(combat.AllyDefense, damage)
		combat.AllyDefense -= absorbed
		damage -= absorbed
		resolution.Absorbed = absorbed
	}
	if damage <= 0 {
		return resolution
	}
	target := combat.resolvePlayerTarget(packet.TargetID)
	if target == nil {
		return resolution
	}
	combat.damageUnit(target, damage)
	resolution.TargetID = target.ID
	resolution.DamageDealt = damage
	return resolution
}

func (combat *CombatState) resolveEnemyTarget(targetID string) *UnitState {
	if targetID != "" {
		target := combat.enemyUnit(targetID)
		if target != nil && !target.Downed && target.HP > 0 {
			return target
		}
		return nil
	}
	alive := combat.aliveEnemyUnits()
	if len(alive) == 0 {
		return nil
	}
	index := alive[combat.random.NextInt(len(alive))]
	return &combat.EnemyUnits[index]
}

func (combat *CombatState) resolvePlayerTarget(targetID string) *UnitState {
	if targetID != "" {
		target := combat.playerUnit(targetID)
		if target != nil && !target.Downed && target.HP > 0 {
			return target
		}
		return nil
	}
	alive := combat.alivePlayerUnits()
	if len(alive) == 0 {
		return nil
	}
	index := alive[combat.random.NextInt(len(alive))]
	return &combat.PlayerUnits[index]
}

func (combat *CombatState) damageUnit(unit *UnitState, damage int) {
	if unit == nil || unit.Downed || damage <= 0 {
		return
	}
	unit.HP -= damage
	if unit.HP <= 0 {
		unit.HP = 0
		unit.Downed = true
		combat.removeOwnerDice(unit.ID)
	}
}

func (combat *CombatState) removeOwnerDice(ownerID string) {
	filterPool := func(pool []DieSpec) []DieSpec {
		filtered := pool[:0]
		for _, die := range pool {
			if die.OwnerID == ownerID {
				continue
			}
			filtered = append(filtered, die)
		}
		return filtered
	}
	filterSelected := func(pool []SelectedDie) []SelectedDie {
		filtered := pool[:0]
		for _, die := range pool {
			if die.Die.OwnerID == ownerID {
				continue
			}
			filtered = append(filtered, die)
		}
		return filtered
	}
	combat.AvailableDice = filterPool(combat.AvailableDice)
	combat.GraveyardDice = filterPool(combat.GraveyardDice)
	combat.SelectedDice = filterSelected(combat.SelectedDice)
}

func (combat *CombatState) applyDamageToUnit(side combatSide, unitID string, damage int) {
	switch side {
	case sidePlayer:
		combat.damageUnit(combat.playerUnit(unitID), damage)
	case sideEnemy:
		combat.damageUnit(combat.enemyUnit(unitID), damage)
	}
}

func (combat *CombatState) nextPatternLabels() map[string]string {
	if len(combat.EnemyUnits) == 0 {
		return nil
	}
	labels := map[string]string{}
	for idx := range combat.EnemyUnits {
		unit := &combat.EnemyUnits[idx]
		if unit.Downed || unit.HP <= 0 || len(unit.Patterns) == 0 {
			continue
		}
		next := unit.Patterns[unit.PatternIdx%len(unit.Patterns)]
		labels[unit.ID] = next.Label
	}
	return labels
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	cloned := map[string]string{}
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
