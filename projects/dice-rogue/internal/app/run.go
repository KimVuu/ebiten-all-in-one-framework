package app

import "fmt"

func newRunState(seed int64) *RunState {
	return &RunState{
		Seed:             seed,
		Screen:           ScreenPartySelection,
		SelectedPartyIDs: nil,
		PartyUnits:       nil,
		CurrentNodeID:    "start",
		NextNodeIDs:      nil,
	}
}

func (run *RunState) startRun() {
	if run == nil || len(run.SelectedPartyIDs) != 3 {
		return
	}
	party := make([]UnitState, 0, len(run.SelectedPartyIDs))
	for _, id := range run.SelectedPartyIDs {
		unit, ok := newCharacterState(id)
		if !ok {
			return
		}
		party = append(party, unit)
	}
	run.PartyUnits = party
	run.CurrentNodeID = "start"
	run.NextNodeIDs = append([]string(nil), actMapNodes["start"].NextIDs...)
	run.Screen = ScreenMap
	run.Outcome = OutcomeState{}
	run.CurrentCombat = nil
}

func (run *RunState) restart() {
	if run == nil {
		return
	}
	run.Screen = ScreenPartySelection
	run.SelectedPartyIDs = nil
	run.PartyUnits = nil
	run.CurrentNodeID = "start"
	run.NextNodeIDs = nil
	run.CurrentCombat = nil
	run.Outcome = OutcomeState{}
}

func (run *RunState) togglePartySelection(id string) {
	if run == nil || run.Screen != ScreenPartySelection {
		return
	}
	for idx, selected := range run.SelectedPartyIDs {
		if selected != id {
			continue
		}
		run.SelectedPartyIDs = append(run.SelectedPartyIDs[:idx], run.SelectedPartyIDs[idx+1:]...)
		return
	}
	if len(run.SelectedPartyIDs) >= 3 {
		return
	}
	run.SelectedPartyIDs = append(run.SelectedPartyIDs, id)
}

func (run *RunState) selectMapNode(id string) error {
	if run == nil || run.Screen != ScreenMap {
		return fmt.Errorf("맵 화면이 아닙니다")
	}
	allowed := false
	for _, nextID := range run.NextNodeIDs {
		if nextID == id {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("선택할 수 없는 노드입니다: %s", id)
	}
	node, ok := mapNodeByID(id)
	if !ok {
		return fmt.Errorf("알 수 없는 노드입니다: %s", id)
	}
	switch node.Kind {
	case NodeKindRest:
		run.applyRestNode(node)
		return nil
	case NodeKindNormal, NodeKindElite, NodeKindBoss:
		_, err := run.startEncounterForNode(id)
		return err
	default:
		return fmt.Errorf("지원하지 않는 노드 종류입니다: %s", node.Kind)
	}
}

func (run *RunState) startEncounterForNode(id string) (*CombatState, error) {
	if run == nil {
		return nil, fmt.Errorf("런 상태를 찾을 수 없습니다")
	}
	node, ok := mapNodeByID(id)
	if !ok {
		return nil, fmt.Errorf("알 수 없는 노드입니다: %s", id)
	}
	encounter, ok := encounterByID(node.EncounterID)
	if !ok {
		return nil, fmt.Errorf("알 수 없는 전투입니다: %s", node.EncounterID)
	}
	party := cloneUnits(run.PartyUnits)
	for idx := range party {
		if party[idx].Downed || party[idx].HP <= 0 {
			party[idx].HP = 1
			party[idx].Downed = false
		}
		party[idx].Statuses = nil
	}
	combat := newCombatState(run.Seed+int64(len(run.CurrentNodeID))+int64(len(id))+int64(len(run.PartyUnits)), party, encounter)
	combat.NodeID = id
	run.CurrentCombat = combat
	run.Screen = ScreenCombat
	return combat, nil
}

func (run *RunState) resolveCombatTurn() TurnResolution {
	if run == nil || run.CurrentCombat == nil {
		return TurnResolution{}
	}
	summary := run.CurrentCombat.resolveTurn()
	if summary.Outcome == CombatOutcomeNone {
		return summary
	}

	run.PartyUnits = cloneUnits(run.CurrentCombat.PlayerUnits)
	run.CurrentNodeID = run.CurrentCombat.NodeID
	node, _ := mapNodeByID(run.CurrentCombat.NodeID)
	run.NextNodeIDs = append([]string(nil), node.NextIDs...)
	run.CurrentCombat = nil
	run.Screen = ScreenOutcome

	switch summary.Outcome {
	case CombatOutcomeEscape:
		run.Outcome = OutcomeState{
			Title:       "도주 성공",
			Body:        "길잡이가 전장을 빠져나갈 길을 찾아냈습니다.",
			CanContinue: true,
		}
	case CombatOutcomeVictory:
		if node.Kind == NodeKindBoss {
			run.Outcome = OutcomeState{
				Title:    "1막 클리어",
				Body:     "보스를 쓰러뜨렸습니다. 1막을 돌파했습니다.",
				RunEnded: true,
			}
		} else {
			run.Outcome = OutcomeState{
				Title:       "승리",
				Body:        "파티가 전투를 돌파했습니다.",
				CanContinue: true,
			}
		}
	case CombatOutcomeDefeat:
		run.Outcome = OutcomeState{
			Title:    "패배",
			Body:     "파티원 3명이 모두 쓰러졌습니다.",
			RunEnded: true,
		}
	}
	return summary
}

func (run *RunState) continueAfterOutcome() {
	if run == nil || run.Screen != ScreenOutcome {
		return
	}
	if run.Outcome.CanContinue {
		run.Screen = ScreenMap
		run.Outcome = OutcomeState{}
		return
	}
	if run.Outcome.RunEnded {
		run.restart()
	}
}

func (run *RunState) applyRestNode(node EncounterNode) {
	run.CurrentNodeID = node.ID
	run.NextNodeIDs = append([]string(nil), node.NextIDs...)
	for idx := range run.PartyUnits {
		unit := &run.PartyUnits[idx]
		if unit.Downed || unit.HP <= 0 {
			unit.HP = 1
			unit.Downed = false
		}
		heal := int(float64(unit.MaxHP) * 0.30)
		if heal < 1 {
			heal = 1
		}
		unit.HP = minInt(unit.MaxHP, unit.HP+heal)
	}
	run.Screen = ScreenOutcome
	run.Outcome = OutcomeState{
		Title:       "휴식",
		Body:        "파티가 휴식 노드에서 회복했습니다.",
		CanContinue: true,
	}
}
