package app

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	ebitenuidebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui-debug"
	renderer "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/renderer"
	gameui "github.com/kimyechan/ebiten-aio-framework/projects/dice-rogue/internal/ui"
)

type Game struct {
	mu sync.RWMutex

	width        int
	height       int
	frame        int
	tick         int
	debugEnabled bool

	renderer       *renderer.Renderer
	runtime        *ebitenui.Runtime
	dom            *ebitenui.DOM
	lastInput      ebitenui.InputSnapshot
	overlayEnabled bool

	run         *RunState
	uiDebug     *ebitenuidebug.Adapter
	debugBridge *ebitendebug.Bridge
}

const (
	DefaultWindowWidth  = 1600
	DefaultWindowHeight = 960
)

func NewGame(config GameConfig) *Game {
	return newGame(config)
}

func newGame(config GameConfig) *Game {
	seed := config.Seed
	if seed == 0 {
		seed = 1
	}
	_ = gameui.ApplyTextFace()
	game := &Game{
		width:          DefaultWindowWidth,
		height:         DefaultWindowHeight,
		debugEnabled:   config.DebugEnabled,
		renderer:       renderer.New(),
		runtime:        ebitenui.NewRuntime(),
		run:            newRunState(seed),
		overlayEnabled: false,
	}
	game.uiDebug = ebitenuidebug.NewAdapter(ebitenuidebug.Config{
		GameID:         "dice-rogue",
		ScreenshotsDir: diceRogueScreenshotsDir(),
	}, ebitenuidebug.Callbacks{
		CurrentLayout:   game.currentLayout,
		CurrentViewport: game.currentViewport,
		CurrentRuntime: func() *ebitenui.Runtime {
			return game.runtime
		},
		CurrentInput: func() ebitenui.InputSnapshot {
			game.mu.RLock()
			defer game.mu.RUnlock()
			return game.lastInput
		},
		CurrentFrame: game.currentFrame,
		OverlayEnabled: func() bool {
			game.mu.RLock()
			defer game.mu.RUnlock()
			return game.overlayEnabled
		},
		SetOverlay: func(enabled bool) {
			game.overlayEnabled = enabled
		},
	})
	return game
}

func diceRogueScreenshotsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join(".", "screenshots", "dice-rogue")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", "screenshots", "dice-rogue"))
}

func (game *Game) Update() error {
	return game.step(game.collectInput())
}

func (game *Game) step(input ebitenui.InputSnapshot) error {
	game.mu.Lock()
	defer game.mu.Unlock()

	game.tick++
	game.frame++
	viewport := game.currentViewportLocked()

	dom := gameui.BuildDOM(game.currentModelLocked(), game.callbacksLocked(), game.runtime)
	layout := dom.Layout(viewport)
	input = game.uiDebug.ApplyQueuedInput(game.frame, dom, game.runtime, layout, input)
	game.runtime.Update(dom, viewport, input)

	dom = gameui.BuildDOM(game.currentModelLocked(), game.callbacksLocked(), game.runtime)
	game.normalizeNonButtonFocus(dom)
	game.applyRuntimeVisualStates(dom, input)
	game.dom = dom
	game.lastInput = input
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	size := screen.Bounds().Size()

	game.mu.Lock()
	game.width = size.X
	game.height = size.Y
	dom := game.dom
	overlayEnabled := game.overlayEnabled
	game.mu.Unlock()

	if dom == nil {
		game.mu.RLock()
		dom = gameui.BuildDOM(game.currentModelLocked(), game.callbacksLocked(), game.runtime)
		game.mu.RUnlock()
		game.normalizeNonButtonFocus(dom)
		game.applyRuntimeVisualStates(dom, game.lastInput)
	}

	viewport := ebitenui.Viewport{
		Width:  float64(size.X),
		Height: float64(size.Y),
	}
	layout := game.renderer.Draw(screen, dom, viewport)
	game.uiDebug.DrawOverlay(screen, layout, overlayEnabled)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return DefaultWindowWidth, DefaultWindowHeight
}

func (game *Game) currentModelLocked() gameui.Model {
	partyRoster := make([]gameui.PartyMember, 0)
	if game.run.CurrentCombat != nil {
		for _, unit := range game.run.CurrentCombat.PlayerUnits {
			partyRoster = append(partyRoster, buildPartyMemberView(unit, false))
		}
	} else if len(game.run.PartyUnits) > 0 {
		for _, unit := range game.run.PartyUnits {
			partyRoster = append(partyRoster, buildPartyMemberView(unit, false))
		}
	} else {
		selected := map[string]bool{}
		for _, id := range game.run.SelectedPartyIDs {
			selected[id] = true
		}
		for _, unit := range characterChoices() {
			if selected[unit.ID] {
				partyRoster = append(partyRoster, buildPartyMemberView(unit, true))
			}
		}
	}

	model := gameui.Model{
		CurrentScreen:  string(game.run.Screen),
		HeaderTitle:    screenTitle(game.run.Screen),
		HeaderSubtitle: fmt.Sprintf("시드 %d / 현재 노드 %s", game.run.Seed, nodeDisplayName(game.run.CurrentNodeID)),
		ViewportWidth:  float64(game.width),
		ViewportHeight: float64(game.height),
		PartyRoster:    partyRoster,
	}

	switch game.run.Screen {
	case ScreenMap:
		nodes := make([]gameui.MapNode, 0, len(game.run.NextNodeIDs))
		for _, id := range game.run.NextNodeIDs {
			node, ok := mapNodeByID(id)
			if !ok {
				continue
			}
			nodes = append(nodes, gameui.MapNode{
				ID:     node.ID,
				Name:   node.Name,
				Kind:   string(node.Kind),
				Detail: mapNodeDetail(node),
			})
		}
		model.Map = gameui.MapModel{
			CurrentNodeID: nodeDisplayName(game.run.CurrentNodeID),
			Nodes:         nodes,
		}
	case ScreenCombat:
		model.Combat = game.currentCombatModelLocked()
	case ScreenOutcome:
		model.Outcome = gameui.OutcomeModel{
			Title:       game.run.Outcome.Title,
			Body:        game.run.Outcome.Body,
			CanContinue: game.run.Outcome.CanContinue,
			RunEnded:    game.run.Outcome.RunEnded,
		}
	default:
		candidates := make([]gameui.PartyMember, 0, len(characterCatalog))
		selected := map[string]bool{}
		for _, id := range game.run.SelectedPartyIDs {
			selected[id] = true
		}
		for _, unit := range characterChoices() {
			candidates = append(candidates, buildPartyMemberView(unit, selected[unit.ID]))
		}
		model.PartySelection = gameui.PartySelectionModel{
			Candidates:    candidates,
			SelectedCount: len(game.run.SelectedPartyIDs),
			CanStart:      len(game.run.SelectedPartyIDs) == 3,
		}
	}
	return model
}

func (game *Game) currentCombatModelLocked() gameui.CombatModel {
	if game.run.CurrentCombat == nil {
		return gameui.CombatModel{}
	}
	combat := game.run.CurrentCombat
	party := make([]gameui.PartyMember, 0, len(combat.PlayerUnits))
	for _, unit := range combat.PlayerUnits {
		status := ""
		if unit.ID == "human-warrior" && unit.Counters[counterHeroGoddess] > 0 {
			status = fmt.Sprintf("여신 스택 %d", unit.Counters[counterHeroGoddess])
		}
		party = append(party, buildPartyMemberViewWithStatus(unit, false, status))
	}
	enemies := make([]gameui.PartyMember, 0, len(combat.EnemyUnits))
	for _, unit := range combat.EnemyUnits {
		enemies = append(enemies, buildPartyMemberView(unit, false))
	}
	available := make([]gameui.DieView, 0, len(combat.AvailableDice))
	for _, die := range combat.AvailableDice {
		available = append(available, gameui.DieView{
			ID:     die.ID,
			Label:  fmt.Sprintf("%s / %s", diceOwnerLabel(die.OwnerID), die.Name),
			Detail: dieDetail(die),
		})
	}
	selected := make([]gameui.DieView, 0, len(combat.SelectedDice))
	for _, die := range combat.SelectedDice {
		selected = append(selected, gameui.DieView{
			ID:     die.Die.ID,
			Label:  fmt.Sprintf("%s / %s", diceOwnerLabel(die.Die.OwnerID), die.Die.Name),
			Detail: dieDetail(die.Die),
			Forced: die.Forced,
		})
	}
	patterns := make([]string, 0, len(combat.RevealedNextPatterns))
	for enemyID, label := range combat.RevealedNextPatterns {
		patterns = append(patterns, fmt.Sprintf("%s 다음 패턴: %s", ownerLabel(enemyID), label))
	}
	logs := append([]string(nil), combat.Logs...)
	return gameui.CombatModel{
		EncounterName:    combat.EncounterName,
		Turn:             maxInt(combat.Turn, 1),
		Party:            party,
		Enemies:          enemies,
		AvailableDice:    available,
		SelectedDice:     selected,
		RevealedPatterns: patterns,
		Logs:             logs,
		CanResolve:       len(combat.SelectedDice) == 3,
		AllyDefense:      combat.AllyDefense,
		EnemyDefense:     combat.EnemyDefense,
		DamageBoost:      combat.enemyDamageBoostPercent(),
	}
}

func buildPartyMemberView(unit UnitState, selected bool) gameui.PartyMember {
	return buildPartyMemberViewWithStatus(unit, selected, "")
}

func buildPartyMemberViewWithStatus(unit UnitState, selected bool, status string) gameui.PartyMember {
	return gameui.PartyMember{
		ID:          unit.ID,
		Name:        unit.Name,
		Role:        unit.Role,
		DiceSummary: unitDiceSummary(unit),
		Status:      status,
		HP:          unit.HP,
		MaxHP:       unit.MaxHP,
		Selected:    selected,
		Downed:      unit.Downed,
	}
}

func (game *Game) callbacksLocked() gameui.Callbacks {
	return gameui.Callbacks{
		OnToggleParty: func(id string) {
			game.run.togglePartySelection(id)
		},
		OnStartRun: func() {
			game.run.startRun()
			game.resetCombatScrollPositions()
		},
		OnSelectMapNode: func(id string) {
			_ = game.run.selectMapNode(id)
			game.resetCombatScrollPositions()
		},
		OnSelectDie: func(id string) {
			if game.run.CurrentCombat != nil {
				_ = game.run.CurrentCombat.selectDie(id)
			}
		},
		OnResolveTurn: func() {
			game.run.resolveCombatTurn()
			game.resetCombatScrollPositions()
		},
		OnContinue: func() {
			game.run.continueAfterOutcome()
		},
		OnRestart: func() {
			game.run.restart()
			game.resetCombatScrollPositions()
		},
	}
}

func (game *Game) frameSnapshot() ebitendebug.FrameSnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return ebitendebug.FrameSnapshot{
		Frame:        game.frame,
		Tick:         game.tick,
		FPS:          ebiten.ActualFPS(),
		TPS:          ebiten.ActualTPS(),
		Paused:       false,
		DebugEnabled: game.debugEnabled,
	}
}

func (game *Game) sceneSnapshot() ebitendebug.SceneSnapshot {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return ebitendebug.SceneSnapshot{
		Current: ebitendebug.SceneRef{ID: "dice-rogue", Name: "주사위 로그"},
		Known:   []ebitendebug.SceneRef{{ID: "dice-rogue", Name: "주사위 로그"}},
		Summary: game.debugSummaryLocked(),
	}
}

func (game *Game) worldSnapshot() ebitendebug.WorldSnapshot {
	layout := game.currentLayout()
	if layout == nil {
		return ebitendebug.WorldSnapshot{}
	}

	game.mu.RLock()
	model := game.currentModelLocked()
	game.mu.RUnlock()

	ids := []struct {
		ID   string
		Type string
	}{
		{ID: "dice-rogue-root", Type: "screen"},
		{ID: model.CurrentScreen + "-screen", Type: "screen-body"},
	}
	for _, member := range model.PartyRoster {
		ids = append(ids, struct {
			ID   string
			Type string
		}{ID: "party-summary-" + member.ID, Type: "party-summary"})
	}
	for _, member := range model.Combat.Party {
		ids = append(ids, struct {
			ID   string
			Type string
		}{ID: "party-card-" + member.ID, Type: "party-unit"})
	}
	for _, member := range model.Combat.Enemies {
		ids = append(ids, struct {
			ID   string
			Type string
		}{ID: "enemy-card-" + member.ID, Type: "enemy-unit"})
	}
	for _, die := range model.Combat.AvailableDice {
		ids = append(ids, struct {
			ID   string
			Type string
		}{ID: "available-die-" + die.ID, Type: "die"})
	}

	entities := make([]ebitendebug.EntitySnapshot, 0, len(ids))
	for _, entry := range ids {
		node, ok := layout.FindByID(entry.ID)
		if !ok {
			continue
		}
		entities = append(entities, ebitendebug.EntitySnapshot{
			ID:      entry.ID,
			Type:    entry.Type,
			Visible: true,
			Enabled: true,
			Tags:    []string{"dice-rogue", string(game.currentScreen())},
			Position: ebitendebug.Vector2{
				X: node.Frame.X,
				Y: node.Frame.Y,
			},
			Size: ebitendebug.Vector2{
				X: node.Frame.Width,
				Y: node.Frame.Height,
			},
		})
	}
	return ebitendebug.WorldSnapshot{Entities: entities}
}

func (game *Game) uiSnapshot() ebitendebug.UISnapshot {
	snapshot := game.uiDebug.UISnapshot()
	if snapshot.Root.Props == nil {
		snapshot.Root.Props = map[string]any{}
	}
	game.mu.RLock()
	for key, value := range game.debugSummaryLocked() {
		snapshot.Root.Props[key] = value
	}
	game.mu.RUnlock()
	return snapshot
}

func (game *Game) debugSummaryLocked() map[string]any {
	summary := map[string]any{
		"currentScreen":        string(game.run.Screen),
		"currentNodeID":        game.run.CurrentNodeID,
		"partyIDs":             []string{},
		"partyHP":              map[string]int{},
		"downedUnitIDs":        []string{},
		"availableDiceCount":   0,
		"graveyardDiceCount":   0,
		"selectedDiceIDs":      []string{},
		"enemyIDs":             []string{},
		"enemyHP":              map[string]int{},
		"revealedNextPatterns": map[string]string{},
		"turn":                 0,
	}

	partyIDs := make([]string, 0, len(game.run.PartyUnits))
	partyHP := map[string]int{}
	downed := make([]string, 0)
	partySource := game.run.PartyUnits
	if game.run.CurrentCombat != nil {
		partySource = game.run.CurrentCombat.PlayerUnits
	}
	for _, unit := range partySource {
		partyIDs = append(partyIDs, unit.ID)
		partyHP[unit.ID] = unit.HP
		if unit.Downed || unit.HP <= 0 {
			downed = append(downed, unit.ID)
		}
	}
	summary["partyIDs"] = partyIDs
	summary["partyHP"] = partyHP
	summary["downedUnitIDs"] = downed

	if combat := game.run.CurrentCombat; combat != nil {
		summary["availableDiceCount"] = len(combat.AvailableDice)
		summary["graveyardDiceCount"] = len(combat.GraveyardDice)
		selectedIDs := make([]string, 0, len(combat.SelectedDice))
		for _, die := range combat.SelectedDice {
			selectedIDs = append(selectedIDs, die.Die.ID)
		}
		summary["selectedDiceIDs"] = selectedIDs
		enemyIDs := make([]string, 0, len(combat.EnemyUnits))
		enemyHP := map[string]int{}
		for _, unit := range combat.EnemyUnits {
			enemyIDs = append(enemyIDs, unit.ID)
			enemyHP[unit.ID] = unit.HP
		}
		summary["enemyIDs"] = enemyIDs
		summary["enemyHP"] = enemyHP
		summary["revealedNextPatterns"] = cloneStringMap(combat.RevealedNextPatterns)
		summary["turn"] = combat.Turn
	}
	return summary
}

func (game *Game) currentViewport() ebitenui.Viewport {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return game.currentViewportLocked()
}

func (game *Game) currentViewportLocked() ebitenui.Viewport {
	width := game.width
	height := game.height
	if width <= 0 {
		width = DefaultWindowWidth
	}
	if height <= 0 {
		height = DefaultWindowHeight
	}
	return ebitenui.Viewport{
		Width:  float64(width),
		Height: float64(height),
	}
}

func (game *Game) currentLayout() *ebitenui.LayoutNode {
	game.mu.RLock()
	dom := game.dom
	viewport := game.currentViewportLocked()
	model := game.currentModelLocked()
	callbacks := game.callbacksLocked()
	game.mu.RUnlock()

	if dom != nil {
		return dom.Layout(viewport)
	}
	if dom == nil {
		dom = gameui.BuildDOM(model, callbacks, game.runtime)
	}
	return dom.Layout(viewport)
}

func (game *Game) currentFrame() int {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return game.frame
}

func (game *Game) currentScreen() ScreenID {
	game.mu.RLock()
	defer game.mu.RUnlock()
	return game.run.Screen
}

func (game *Game) resetCombatScrollPositions() {
	if game.runtime == nil {
		return
	}
	game.runtime.SetNumberValue("available-dice-scroll-offset", 0)
	game.runtime.SetNumberValue("used-dice-scroll-offset", 0)
	game.runtime.SetNumberValue("combat-log-scroll-offset", 0)
}

func (game *Game) normalizeNonButtonFocus(dom *ebitenui.DOM) {
	if game.runtime == nil || dom == nil {
		return
	}
	focusedID := game.runtime.FocusedID()
	if focusedID == "" {
		return
	}
	node, ok := dom.FindByID(focusedID)
	if ok && node != nil && node.Tag == ebitenui.TagButton {
		return
	}
	game.runtime.ClearFocus(dom, ebitenui.InputSnapshot{})
}

func (game *Game) applyRuntimeVisualStates(dom *ebitenui.DOM, input ebitenui.InputSnapshot) {
	if game.runtime == nil || dom == nil {
		return
	}
	hoveredID := game.runtime.HoveredID()
	focusedID := game.runtime.FocusedID()

	if hoveredID != "" {
		if node, ok := dom.FindByID(hoveredID); ok && node.Tag == ebitenui.TagButton {
			state := node.Props.State
			state.Hovered = true
			node.Props.State = state
		}
	}
	if focusedID != "" {
		if node, ok := dom.FindByID(focusedID); ok && node.Tag == ebitenui.TagButton {
			state := node.Props.State
			state.Focused = true
			if input.PointerDown {
				state.Pressed = true
			}
			node.Props.State = state
		}
	}
}

func (game *Game) collectInput() ebitenui.InputSnapshot {
	pointerX, pointerY := ebiten.CursorPosition()
	scrollX, scrollY := ebiten.Wheel()
	textInput := ebiten.AppendInputChars(nil)

	input := ebitenui.InputSnapshot{
		PointerX:    float64(pointerX),
		PointerY:    float64(pointerY),
		PointerDown: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		ScrollX:     scrollX,
		ScrollY:     scrollY,
	}
	if len(textInput) > 0 {
		input.Text = string(textInput)
	}
	input.Backspace = inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
	input.Delete = inpututil.IsKeyJustPressed(ebiten.KeyDelete)
	input.Home = inpututil.IsKeyJustPressed(ebiten.KeyHome)
	input.End = inpututil.IsKeyJustPressed(ebiten.KeyEnd)
	input.Submit = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	input.Space = inpututil.IsKeyJustPressed(ebiten.KeySpace)
	input.Tab = inpututil.IsKeyJustPressed(ebiten.KeyTab)
	input.Escape = inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	input.ArrowUp = inpututil.IsKeyJustPressed(ebiten.KeyArrowUp)
	input.ArrowDown = inpututil.IsKeyJustPressed(ebiten.KeyArrowDown)
	input.ArrowLeft = inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft)
	input.ArrowRight = inpututil.IsKeyJustPressed(ebiten.KeyArrowRight)
	input.Shift = ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	input.Control = ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)
	input.Alt = ebiten.IsKeyPressed(ebiten.KeyAltLeft) || ebiten.IsKeyPressed(ebiten.KeyAltRight)
	input.Meta = ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight)
	input.SelectAll = (input.Control || input.Meta) && inpututil.IsKeyJustPressed(ebiten.KeyA)
	return input
}

func (game *Game) startDebugBridge(addr string) error {
	bridge := ebitendebug.New(ebitendebug.Config{
		Enabled: true,
		Addr:    addr,
		GameID:  "dice-rogue",
		Version: "v1",
	})
	bridge.SetFrameProvider(game.frameSnapshot)
	bridge.SetSceneProvider(game.sceneSnapshot)
	bridge.SetWorldProvider(game.worldSnapshot)
	game.uiDebug.Attach(bridge)
	if err := bridge.Start(); err != nil {
		return err
	}

	game.mu.Lock()
	game.debugBridge = bridge
	game.mu.Unlock()
	return nil
}

func (game *Game) stopDebugBridge() error {
	game.mu.Lock()
	bridge := game.debugBridge
	game.debugBridge = nil
	game.mu.Unlock()
	if bridge == nil {
		return nil
	}
	return bridge.Close(context.Background())
}

func (game *Game) StartDebugBridge(addr string) error {
	return game.startDebugBridge(addr)
}

func (game *Game) StopDebugBridge() error {
	return game.stopDebugBridge()
}

func (game *Game) debugBridgeLikeCommand(name string, args map[string]any) ebitendebug.CommandResult {
	game.mu.RLock()
	bridge := game.debugBridge
	game.mu.RUnlock()
	if bridge != nil {
		return bridge.InvokeCommand(name, args)
	}
	bridge = ebitendebug.New(ebitendebug.Config{Enabled: true, GameID: "dice-rogue", Version: "v1"})
	game.uiDebug.Attach(bridge)
	return bridge.InvokeCommand(name, args)
}

func screenTitle(screen ScreenID) string {
	switch screen {
	case ScreenMap:
		return "1막 지도"
	case ScreenCombat:
		return "전투"
	case ScreenOutcome:
		return "휴식 / 결과"
	default:
		return "파티 선택"
	}
}

func unitDiceSummary(unit UnitState) string {
	counts := map[DieKind]int{}
	for _, die := range unit.Dice {
		counts[die.Kind]++
	}
	return fmt.Sprintf("공격 %d / 방어 %d / 스킬 %d", counts[DieKindAttack], counts[DieKindDefense], counts[DieKindSkill])
}

func ownerLabel(id string) string {
	if unit, ok := newCharacterState(id); ok {
		return unit.Name
	}
	for _, encounter := range encounterCatalog {
		for _, enemy := range encounter.Enemies {
			if enemy.ID == id {
				return enemy.Name
			}
		}
	}
	return strings.ReplaceAll(id, "-", " ")
}

func diceOwnerLabel(id string) string {
	switch id {
	case "human-warrior":
		return "용사"
	case "human-guard":
		return "방패병"
	case "human-priest":
		return "여신관"
	case "human-guide":
		return "길잡이"
	case "elf-archer":
		return "궁수"
	case "dwarf-smith":
		return "대장장이"
	default:
		return ownerLabel(id)
	}
}

func dieDetail(die DieSpec) string {
	switch die.Kind {
	case DieKindAttack:
		return "무작위 적 1명 피해"
	case DieKindDefense:
		return "파티 공유 방어"
	default:
		switch die.EffectID {
		case effectHeroGoddess:
			return "성공2 실패4 / 3스택 전체10"
		case effectTankGuard:
			return "성공2 실패4 / 방어 2배"
		case effectPriestHeal:
			return "성공1 실패5 / 전체 2회복"
		case effectGuideInfo:
			return "성공3 실패2 도주1 / 다음 패턴 공개"
		case effectGuideWeakness:
			return "성공3 실패2 도주1 / 2턴 피해+30%"
		case effectGuideEscape:
			return "도주5 실패1 / 일반전 조건부 도주"
		case effectArcherShot:
			return "1 2 3 4 치명 치명 / 치명 시 2배"
		case effectSmithForge:
			return "성공3 실패3 / 피해+50%"
		default:
			return die.EffectID
		}
	}
}

func mapNodeDetail(node EncounterNode) string {
	switch node.Kind {
	case NodeKindRest:
		return "최대 체력의 30%를 회복하고 쓰러진 아군은 먼저 1 체력으로 복귀합니다."
	case NodeKindElite:
		return "엘리트 전투입니다."
	case NodeKindBoss:
		return "보스 전투입니다. 도주할 수 없습니다."
	default:
		return "일반 전투입니다."
	}
}

func nodeDisplayName(id string) string {
	if node, ok := mapNodeByID(id); ok {
		return node.Name
	}
	return fallbackString(id, "시작")
}

func fallbackString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
