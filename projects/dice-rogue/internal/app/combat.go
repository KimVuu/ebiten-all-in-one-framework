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
	allyDefense := 0
	tankMultiplier := false
	guideReveal := false
	guideEscapeMain := false
	guideEscapeSupport := false
	smithBoost := 0
	attackPackets := make([]attackPacket, 0)
	archerPacketIndexes := make([]int, 0)
	archerCrits := 0
	defenseLogs := make([]string, 0)
	utilityLogs := make([]string, 0)
	attackLogs := make([]string, 0)

	for _, selected := range combat.SelectedDice {
		die := selected.Die
		face := die.Faces[combat.random.NextInt(len(die.Faces))]
		switch die.Kind {
		case DieKindAttack:
			attackLogs = append(attackLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
			if face.Kind == FaceKindValue && face.Value > 0 {
				attackPackets = append(attackPackets, attackPacket{
					SourceUnitID: die.OwnerID,
					SourceDieID:  die.ID,
					Value:        face.Value,
				})
			}
		case DieKindDefense:
			defenseLogs = append(defenseLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
			if face.Kind == FaceKindValue && face.Value > 0 {
				allyDefense += face.Value
			}
		case DieKindSkill:
			switch die.EffectID {
			case effectHeroGoddess:
				attackLogs = append(attackLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					unit := combat.playerUnit(die.OwnerID)
					unit.Counters[counterHeroGoddess]++
					attackLogs = append(attackLogs, fmt.Sprintf("%s의 여신 스택이 %d가 되었다.", ownerLabel(die.OwnerID), unit.Counters[counterHeroGoddess]))
					for unit.Counters[counterHeroGoddess] >= 3 {
						unit.Counters[counterHeroGoddess] -= 3
						attackLogs = append(attackLogs, fmt.Sprintf("%s가 여신 폭발을 발동했다.", ownerLabel(die.OwnerID)))
						for _, enemyIdx := range combat.aliveEnemyUnits() {
							attackPackets = append(attackPackets, attackPacket{
								SourceUnitID: die.OwnerID,
								SourceDieID:  die.ID,
								Value:        10,
								TargetID:     combat.EnemyUnits[enemyIdx].ID,
							})
						}
					}
				}
			case effectTankGuard:
				defenseLogs = append(defenseLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					tankMultiplier = true
					defenseLogs = append(defenseLogs, "방패병이 이번 턴 방어 2배를 준비했다.")
				}
			case effectPriestHeal:
				utilityLogs = append(utilityLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					for idx := range combat.PlayerUnits {
						unit := &combat.PlayerUnits[idx]
						if unit.Downed || unit.HP <= 0 {
							continue
						}
						unit.HP = minInt(unit.MaxHP, unit.HP+2)
					}
					utilityLogs = append(utilityLogs, "여신관이 생존한 아군을 회복시켰다.")
				}
			case effectGuideInfo:
				utilityLogs = append(utilityLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					guideReveal = true
					utilityLogs = append(utilityLogs, "길잡이가 다음 적 패턴을 간파했다.")
				}
				if face.Kind == FaceKindEscape {
					guideEscapeSupport = true
				}
			case effectGuideWeakness:
				utilityLogs = append(utilityLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					combat.addEnemyStatus(StatusEffect{
						ID:             statusWeakness,
						Name:           "약점",
						Magnitude:      30,
						RemainingTurns: 2,
					})
					utilityLogs = append(utilityLogs, "길잡이가 적의 약점을 드러냈다.")
				}
				if face.Kind == FaceKindEscape {
					guideEscapeSupport = true
				}
			case effectGuideEscape:
				utilityLogs = append(utilityLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindEscape {
					guideEscapeMain = true
				}
			case effectArcherShot:
				attackLogs = append(attackLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindValue && face.Value > 0 {
					attackPackets = append(attackPackets, attackPacket{
						SourceUnitID: die.OwnerID,
						SourceDieID:  die.ID,
						Value:        face.Value,
					})
					archerPacketIndexes = append(archerPacketIndexes, len(attackPackets)-1)
				}
				if face.Kind == FaceKindCritical {
					archerCrits++
				}
			case effectSmithForge:
				utilityLogs = append(utilityLogs, fmt.Sprintf("%s 주사위 결과: %s.", die.Name, face.Label))
				if face.Kind == FaceKindSuccess {
					smithBoost = 50
					utilityLogs = append(utilityLogs, "대장장이가 이번 턴 적이 받는 피해를 증폭시켰다.")
				}
			}
		}
	}

	for _, selected := range combat.SelectedDice {
		combat.GraveyardDice = append(combat.GraveyardDice, selected.Die)
	}
	combat.SelectedDice = nil

	if tankMultiplier {
		allyDefense *= 2
	}
	for _, line := range defenseLogs {
		combat.pushLog(line)
	}
	combat.AllyDefense = allyDefense
	summary.GeneratedAllyDefense = allyDefense
	if allyDefense > 0 {
		combat.pushLog(fmt.Sprintf("파티 방어막이 %d가 되었다.", allyDefense))
	}
	for _, line := range utilityLogs {
		combat.pushLog(line)
	}

	if archerCrits > 0 {
		for _, idx := range archerPacketIndexes {
			attackPackets[idx].Value *= 2
		}
	}
	for _, line := range attackLogs {
		combat.pushLog(line)
	}

	damageBoost := smithBoost + combat.enemyDamageBoostPercent()
	summary.DamageBoostPercent = damageBoost
	var playerResolutions []damageResolution
	summary.PlayerTargets, playerResolutions = combat.applyPacketsToEnemies(attackPackets, damageBoost)
	for _, resolution := range playerResolutions {
		combat.pushLog(playerAttackLogLine(resolution))
	}
	combat.tickEnemyStatusesAfterPlayerPhase()
	combat.EnemyDefense = 0

	if guideEscapeMain && guideEscapeSupport && combat.EncounterKind == EncounterKindNormal {
		combat.Outcome = CombatOutcomeEscape
		combat.pushLog("도주에 성공했다.")
	}
	if len(combat.aliveEnemyUnits()) == 0 {
		combat.Outcome = CombatOutcomeVictory
		combat.pushLog("전투를 돌파했다.")
	}
	if combat.Outcome != CombatOutcomeNone {
		summary.Outcome = combat.Outcome
		summary.Logs = append(summary.Logs, combat.Logs...)
		return summary
	}

	enemyDefense := 0
	enemyTargets := make([]string, 0)
	type enemyAction struct {
		enemyID   string
		enemyName string
		pattern   EncounterPattern
	}
	enemyActions := make([]enemyAction, 0)
	for _, enemyIdx := range combat.aliveEnemyUnits() {
		enemy := &combat.EnemyUnits[enemyIdx]
		if len(enemy.Patterns) == 0 {
			continue
		}
		pattern := enemy.Patterns[enemy.PatternIdx%len(enemy.Patterns)]
		enemy.PatternIdx++
		enemyDefense += pattern.Defense
		enemyActions = append(enemyActions, enemyAction{
			enemyID:   enemy.ID,
			enemyName: enemy.Name,
			pattern:   pattern,
		})
	}
	summary.GeneratedEnemyDefense = enemyDefense
	for _, action := range enemyActions {
		if action.pattern.Defense > 0 {
			combat.pushLog(fmt.Sprintf("%s 행동: 수비 %d.", action.enemyName, action.pattern.Defense))
		}
	}
	if enemyDefense > 0 {
		combat.pushLog(fmt.Sprintf("적 방어막이 %d가 되었다.", enemyDefense))
	}
	for _, action := range enemyActions {
		if len(action.pattern.Attacks) == 0 {
			continue
		}
		combat.pushLog(fmt.Sprintf("%s 행동: 공격 %s.", action.enemyName, joinAttackValues(action.pattern.Attacks)))
		for _, value := range action.pattern.Attacks {
			resolution := combat.applyPacketToPlayers(attackPacket{
				SourceUnitID: action.enemyID,
				Value:        value,
			})
			if resolution.TargetID != "" {
				enemyTargets = append(enemyTargets, resolution.TargetID)
			}
			combat.pushLog(enemyAttackLogLine(action.enemyName, resolution))
			if len(combat.alivePlayerUnits()) == 0 {
				break
			}
		}
		if len(combat.alivePlayerUnits()) == 0 {
			break
		}
	}
	combat.AllyDefense = 0
	combat.EnemyDefense = enemyDefense
	summary.EnemyTargets = enemyTargets

	if len(combat.alivePlayerUnits()) == 0 {
		combat.Outcome = CombatOutcomeDefeat
		combat.pushLog("파티가 전멸했다.")
	}

	if guideReveal {
		combat.RevealedNextPatterns = combat.nextPatternLabels()
	} else if hadReveal {
		combat.RevealedNextPatterns = nil
	}
	summary.RevealedNextPatterns = cloneStringMap(combat.RevealedNextPatterns)
	summary.Outcome = combat.Outcome
	summary.Logs = append(summary.Logs, combat.Logs...)

	if combat.Outcome == CombatOutcomeNone {
		combat.prepareSelection()
	}
	return summary
}

func (combat *CombatState) pushLog(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}
	combat.Logs = append(combat.Logs, line)
}

func joinAttackValues(values []int) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, fmt.Sprintf("%d", value))
	}
	return strings.Join(parts, ", ")
}

func playerAttackLogLine(resolution damageResolution) string {
	attacker := ownerLabel(resolution.SourceUnitID)
	switch {
	case resolution.Absorbed > 0 && resolution.DamageDealt > 0:
		return fmt.Sprintf("%s의 공격은 적 방어막이 %d 막았고, %s에게 %d 피해를 입혔다.", attacker, resolution.Absorbed, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	case resolution.Absorbed > 0:
		return fmt.Sprintf("%s의 공격은 적 방어막이 %d 막았다.", attacker, resolution.Absorbed)
	case resolution.DamageDealt > 0:
		return fmt.Sprintf("%s의 공격이 %s에게 %d 피해를 입혔다.", attacker, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	default:
		return fmt.Sprintf("%s의 공격은 피해를 주지 못했다.", attacker)
	}
}

func enemyAttackLogLine(enemyName string, resolution damageResolution) string {
	switch {
	case resolution.Absorbed > 0 && resolution.DamageDealt > 0:
		return fmt.Sprintf("%s의 공격은 파티 방어막이 %d 막았고, %s에게 %d 피해를 입혔다.", enemyName, resolution.Absorbed, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	case resolution.Absorbed > 0:
		return fmt.Sprintf("%s의 공격은 파티 방어막이 %d 막았다.", enemyName, resolution.Absorbed)
	case resolution.DamageDealt > 0:
		return fmt.Sprintf("%s의 공격이 %s에게 %d 피해를 입혔다.", enemyName, ownerLabel(resolution.TargetID), resolution.DamageDealt)
	default:
		return fmt.Sprintf("%s의 공격은 피해를 주지 못했다.", enemyName)
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

func (combat *CombatState) applyPacketsToEnemies(packets []attackPacket, boostPercent int) ([]string, []damageResolution) {
	targets := make([]string, 0)
	resolutions := make([]damageResolution, 0, len(packets))
	for _, packet := range packets {
		damage := packet.Value
		if boostPercent > 0 {
			damage = int(math.Floor(float64(damage) * float64(100+boostPercent) / 100.0))
		}
		if damage <= 0 {
			continue
		}
		resolution := damageResolution{
			SourceUnitID: packet.SourceUnitID,
			SourceDieID:  packet.SourceDieID,
			BaseDamage:   damage,
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
			resolutions = append(resolutions, resolution)
			continue
		}
		target := combat.resolveEnemyTarget(packet.TargetID)
		if target == nil {
			break
		}
		combat.damageUnit(target, damage)
		resolution.TargetID = target.ID
		resolution.DamageDealt = damage
		resolutions = append(resolutions, resolution)
		targets = append(targets, target.ID)
		if len(combat.aliveEnemyUnits()) == 0 {
			break
		}
	}
	return targets, resolutions
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
