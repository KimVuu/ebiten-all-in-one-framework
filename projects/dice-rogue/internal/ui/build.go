package ui

import (
	"fmt"
	"image/color"
	"strings"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

var (
	rootBackground  = color.RGBA{R: 16, G: 19, B: 27, A: 255}
	panelBackground = color.RGBA{R: 27, G: 33, B: 48, A: 255}
	borderColor     = color.RGBA{R: 66, G: 76, B: 104, A: 255}
	accentColor     = color.RGBA{R: 224, G: 164, B: 74, A: 255}
	selectedColor   = color.RGBA{R: 62, G: 113, B: 184, A: 255}
	disabledColor   = color.RGBA{R: 49, G: 55, B: 70, A: 255}
	textStrong      = color.RGBA{R: 244, G: 246, B: 252, A: 255}
	textMuted       = color.RGBA{R: 181, G: 188, B: 205, A: 255}
	successColor    = color.RGBA{R: 84, G: 164, B: 118, A: 255}
	dangerColor     = color.RGBA{R: 176, G: 82, B: 82, A: 255}
)

const (
	rootPadding           = 14.0
	rootGap               = 12.0
	panelPadding          = 10.0
	panelGap              = 8.0
	gridGap               = 10.0
	combatColGap          = 12.0
	defaultScrollStep     = 48.0
	minTileWidth          = 210.0
	defaultViewportWidth  = 1600.0
	defaultViewportHeight = 960.0
)

type layoutMetrics struct {
	viewportWidth  float64
	viewportHeight float64
	contentWidth   float64
}

type buttonTrigger int

const (
	buttonTriggerClick buttonTrigger = iota
	buttonTriggerPointerDown
	buttonTriggerPointerHold
	buttonTriggerPointerUp
)

func BuildDOM(model Model, callbacks Callbacks, runtime *ebitenui.Runtime) *ebitenui.DOM {
	metrics := resolveLayoutMetrics(model)
	children := []*ebitenui.Node{
		buildHeader(model),
	}
	if model.CurrentScreen != "combat" {
		children = append(children, buildPartyRoster(model.PartyRoster, metrics))
	}
	children = append(children, ebitenui.Div(ebitenui.Props{
		ID: "current-screen-slot",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Height:    ebitenui.Fill(),
			Direction: ebitenui.Column,
		},
	}, buildCurrentScreen(model, callbacks, metrics, runtime)))

	root := ebitenui.Div(ebitenui.Props{
		ID: "dice-rogue-root",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Height:          ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(rootPadding),
			Gap:             rootGap,
			BackgroundColor: rootBackground,
		},
	}, children...)
	return ebitenui.New(root)
}

func buildHeader(model Model) *ebitenui.Node {
	return panel("screen-header", "주사위 로그",
		ebitenui.Text(model.HeaderTitle, ebitenui.Props{
			ID:    "screen-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.TextBlock(model.HeaderSubtitle, ebitenui.Props{
			ID:    "screen-subtitle",
			Style: ebitenui.Style{Color: textMuted},
		}),
	)
}

func buildPartyRoster(roster []PartyMember, metrics layoutMetrics) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.Text("파티", ebitenui.Props{
			ID:    "party-roster-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
	}
	if len(roster) == 0 {
		children = append(children, ebitenui.TextBlock("3명의 파티원을 선택해 1막을 시작하세요.", ebitenui.Props{
			ID:    "party-roster-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	} else {
		cards := make([]*ebitenui.Node, 0, len(roster))
		for _, member := range roster {
			label := fmt.Sprintf("역할 %s / 체력 %d/%d", member.Role, member.HP, member.MaxHP)
			if member.Downed {
				label += " / 쓰러짐"
			}
			cards = append(cards, infoCard(
				fmt.Sprintf("party-summary-%s", member.ID),
				member.Name,
				label,
				colorForUnit(member),
			))
		}
		children = append(children, cardGrid("party-roster-grid", metrics, gridColumnCount(metrics, 3), cards...))
	}
	return panel("party-roster-panel", "파티 현황", children...)
}

func buildCurrentScreen(model Model, callbacks Callbacks, metrics layoutMetrics, runtime *ebitenui.Runtime) *ebitenui.Node {
	switch model.CurrentScreen {
	case "map":
		return buildMapScreen(model.Map, callbacks, metrics)
	case "combat":
		return buildCombatScreen(model.Combat, callbacks, metrics, runtime)
	case "outcome":
		return buildOutcomeScreen(model.Outcome, callbacks)
	default:
		return buildPartySelectionScreen(model.PartySelection, callbacks, metrics)
	}
}

func buildPartySelectionScreen(model PartySelectionModel, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.TextBlock(fmt.Sprintf("선택 %d / 3", model.SelectedCount), ebitenui.Props{
			ID:    "party-selection-count",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}
	gridChildren := make([]*ebitenui.Node, 0, len(model.Candidates))
	for _, candidate := range model.Candidates {
		label := fmt.Sprintf("%s / %s / 체력 %d", candidate.Name, candidate.Role, candidate.MaxHP)
		gridChildren = append(gridChildren, button(
			fmt.Sprintf("party-option-%s", candidate.ID),
			label,
			candidate.DiceSummary,
			candidate.Selected,
			false,
			buttonTriggerClick,
			func(id string) func() {
				return func() {
					if callbacks.OnToggleParty != nil {
						callbacks.OnToggleParty(id)
					}
				}
			}(candidate.ID),
		))
	}
	children = append(children, cardGrid("party-selection-grid", metrics, gridColumnCount(metrics, 2), gridChildren...))
	children = append(children, button(
		"start-run-button",
		"출발하기",
		"선택한 파티로 1막 지도에 진입합니다.",
		false,
		!model.CanStart,
		buttonTriggerClick,
		func() {
			if callbacks.OnStartRun != nil {
				callbacks.OnStartRun()
			}
		},
	))
	return screenPanel("party-selection-screen", "파티 선택", children...)
}

func buildMapScreen(model MapModel, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.TextBlock("고정된 1막 경로 중 다음 노드를 선택하세요.", ebitenui.Props{
			ID:    "map-copy",
			Style: ebitenui.Style{Color: textMuted},
		}),
		ebitenui.Text(fmt.Sprintf("현재 노드: %s", fallback(model.CurrentNodeID, "시작")), ebitenui.Props{
			ID:    "map-current-node",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}
	gridChildren := make([]*ebitenui.Node, 0, len(model.Nodes))
	for _, node := range model.Nodes {
		gridChildren = append(gridChildren, button(
			fmt.Sprintf("map-node-%s", node.ID),
			fmt.Sprintf("%s / %s", nodeKindLabel(node.Kind), node.Name),
			node.Detail,
			false,
			false,
			buttonTriggerClick,
			func(id string) func() {
				return func() {
					if callbacks.OnSelectMapNode != nil {
						callbacks.OnSelectMapNode(id)
					}
				}
			}(node.ID),
		))
	}
	children = append(children, cardGrid("map-node-grid", metrics, gridColumnCount(metrics, 2), gridChildren...))
	return screenPanel("map-screen", "1막 지도", children...)
}

func buildCombatScreen(model CombatModel, callbacks Callbacks, metrics layoutMetrics, runtime *ebitenui.Runtime) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.Text(fmt.Sprintf("%d턴", maxOne(model.Turn)), ebitenui.Props{
			ID:    "combat-turn",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Text(fmt.Sprintf("아군 방어 %d / 적 방어 %d / 피해 증가 +%d%%", model.AllyDefense, model.EnemyDefense, model.DamageBoost), ebitenui.Props{
			ID:    "combat-defense-summary",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}

	revealChildren := []*ebitenui.Node{}
	if len(model.RevealedPatterns) == 0 {
		revealChildren = append(revealChildren, ebitenui.TextBlock("공개된 다음 적 패턴이 없습니다.", ebitenui.Props{
			ID:    "revealed-patterns-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	} else {
		for index, text := range model.RevealedPatterns {
			revealChildren = append(revealChildren, infoCard(
				fmt.Sprintf("revealed-pattern-%d", index),
				fmt.Sprintf("적 %d", index+1),
				text,
				accentColor,
			))
		}
	}

	availableChildren := make([]*ebitenui.Node, 0, len(model.AvailableDice))
	for _, die := range model.AvailableDice {
		detail := die.Detail
		if die.Forced {
			detail += " / 강제 선택"
		}
		availableChildren = append(availableChildren, compactButton(
			fmt.Sprintf("available-die-%s", die.ID),
			die.Label,
			detail,
			buttonTriggerClick,
			func(id string) func() {
				return func() {
					if callbacks.OnSelectDie != nil {
						callbacks.OnSelectDie(id)
					}
				}
			}(die.ID),
		))
	}
	if len(availableChildren) == 0 {
		availableChildren = append(availableChildren, infoCard(
			"available-dice-empty",
			"선택 없음",
			"현재 풀에서 직접 고를 주사위가 남아 있지 않습니다.",
			textMuted,
		))
	}

	usedChildren := make([]*ebitenui.Node, 0, len(model.SelectedDice))
	for _, die := range model.SelectedDice {
		detail := die.Detail + " / 선택됨"
		if die.Forced {
			detail = die.Detail + " / 강제 선택"
		}
		usedChildren = append(usedChildren, infoCard(
			fmt.Sprintf("used-die-%s", die.ID),
			die.Label,
			detail,
			textMuted,
		))
	}
	if len(usedChildren) == 0 {
		usedChildren = append(usedChildren, infoCard(
			"used-dice-empty",
			"사용한 주사위 없음",
			"주사위를 고르면 이 칸에 표시됩니다.",
			textMuted,
		))
	}

	logChildren := make([]*ebitenui.Node, 0, len(model.Logs))
	if len(model.Logs) == 0 {
		logChildren = append(logChildren, infoCard(
			"combat-log-empty",
			"전투 기록 없음",
			"턴이 끝나면 전투 기록이 여기에 표시됩니다.",
			textMuted,
		))
	} else {
		for index, line := range model.Logs {
			logChildren = append(logChildren, infoCard(
				fmt.Sprintf("combat-log-%d", index),
				fmt.Sprintf("기록 %d", index+1),
				line,
				textMuted,
			))
		}
	}
	combatColumns := combatColumnCount(metrics)
	columnWidths := combatColumnWidths(metrics, combatColumns)
	leftWidth := columnWidths[0]
	centerWidth := columnWidths[minInt(len(columnWidths)-1, 1)]
	rightWidth := columnWidths[minInt(len(columnWidths)-1, 2)]

	buildAvailableDicePanel := func(panelWidth float64) *ebitenui.Node {
		return panelWithHeight("available-dice-panel", "사용 가능 주사위", ebitenui.Fill(),
			persistentScrollView("available-dice-scroll", runtime,
				cardColumnWithWidth("available-dice-grid", panelWidth, availableChildren...),
			),
		)
	}
	buildUsedDicePanel := func(panelWidth float64) *ebitenui.Node {
		return panelWithHeight("used-dice-panel", "사용한 주사위", ebitenui.Fill(),
			persistentScrollView("used-dice-scroll", runtime,
				cardColumnWithWidth("used-dice-grid", panelWidth, usedChildren...),
			),
		)
	}
	buildCombatLogPanel := func(panelWidth float64) *ebitenui.Node {
		return panelWithHeight("combat-log-panel", "전투 기록", ebitenui.Fill(),
			persistentScrollView("combat-log-scroll", runtime,
				cardColumnWithWidth("combat-log-grid", panelWidth, logChildren...),
			),
		)
	}

	leftColumn := fixedWidthColumn("combat-left-column", leftWidth,
		panel("combat-party-panel", "파티", buildUnitCards(model.Party, "party-card")...),
		ebitenui.Div(ebitenui.Props{
			ID: "combat-dice-stack",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Height:    ebitenui.Fill(),
				Direction: ebitenui.Column,
				Gap:       combatColGap,
			},
		},
			buildAvailableDicePanel(leftWidth),
			buildUsedDicePanel(leftWidth),
		),
	)

	centerColumn := fixedWidthColumn("combat-center-column", centerWidth,
		buildCombatLogPanel(centerWidth),
	)
	rightChildren := []*ebitenui.Node{
		panel("combat-enemy-panel", "적", buildUnitCards(model.Enemies, "enemy-card")...),
		panelWithHeight("revealed-patterns-panel", "공개된 패턴", ebitenui.Fill(), cardGridWithWidth("revealed-pattern-grid", rightWidth, 1, revealChildren...)...),
		button(
			"resolve-turn-button",
			"턴 진행",
			"선택한 주사위 3개를 처리한 뒤 적이 행동합니다.",
			false,
			!model.CanResolve,
			buttonTriggerClick,
			func() {
				if callbacks.OnResolveTurn != nil {
					callbacks.OnResolveTurn()
				}
			},
		),
	}
	rightColumn := fixedWidthColumn("combat-right-column", rightWidth, rightChildren...)

	combatRowChildren := []*ebitenui.Node{leftColumn}
	if combatColumns >= 3 {
		combatRowChildren = append(combatRowChildren, centerColumn, rightColumn)
	} else if combatColumns == 2 {
		secondColumnChildren := []*ebitenui.Node{
			buildCombatLogPanel(rightWidth),
		}
		secondColumnChildren = append(secondColumnChildren, rightChildren...)
		combatRowChildren = append(combatRowChildren, fixedWidthColumn("combat-right-column", rightWidth, secondColumnChildren...))
	} else {
		leftColumn.Children = append(leftColumn.Children, buildCombatLogPanel(leftWidth))
		leftColumn.Children = append(leftColumn.Children, rightChildren...)
	}

	children = append(children, ebitenui.Div(ebitenui.Props{
		ID: "combat-dashboard",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Height:    ebitenui.Fill(),
			Direction: ebitenui.Row,
			Gap:       combatColGap,
		},
	}, combatRowChildren...))

	return screenPanel("combat-screen", fallback(model.EncounterName, "전투"), children...)
}

func buildOutcomeScreen(model OutcomeModel, callbacks Callbacks) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.TextBlock(model.Body, ebitenui.Props{
			ID:    "outcome-body",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}
	if model.CanContinue {
		children = append(children, button(
			"continue-button",
			"계속",
			"1막 지도로 돌아갑니다.",
			false,
			false,
			buttonTriggerClick,
			func() {
				if callbacks.OnContinue != nil {
					callbacks.OnContinue()
				}
			},
		))
	}
	if model.RunEnded {
		children = append(children, button(
			"restart-button",
			"다시 시작",
			"파티 선택부터 새 런을 시작합니다.",
			false,
			false,
			buttonTriggerClick,
			func() {
				if callbacks.OnRestart != nil {
					callbacks.OnRestart()
				}
			},
		))
	}
	return screenPanel("outcome-screen", fallback(model.Title, "결과"), children...)
}

func buildUnitList(units []PartyMember, prefix string) []*ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(units))
	for _, unit := range units {
		label := fmt.Sprintf("%s / %s / 체력 %d/%d", unit.Name, unit.Role, unit.HP, unit.MaxHP)
		if unit.Downed {
			label += " / 쓰러짐"
		}
		if unit.Status != "" {
			label += " / " + unit.Status
		}
		children = append(children, ebitenui.Text(label, ebitenui.Props{
			ID:    fmt.Sprintf("%s-%s", prefix, unit.ID),
			Style: ebitenui.Style{Color: colorForUnit(unit)},
		}))
	}
	if len(children) == 0 {
		children = append(children, ebitenui.TextBlock("표시할 유닛이 없습니다.", ebitenui.Props{
			ID:    prefix + "-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	}
	return children
}

func buildUnitCards(units []PartyMember, prefix string) []*ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(units))
	for _, unit := range units {
		label := fmt.Sprintf("%s / 체력 %d/%d", unit.Role, unit.HP, unit.MaxHP)
		if unit.Downed {
			label += " / 쓰러짐"
		}
		if unit.Status != "" {
			label += " / " + unit.Status
		}
		children = append(children, infoCard(
			fmt.Sprintf("%s-%s", prefix, unit.ID),
			unit.Name,
			label,
			colorForUnit(unit),
		))
	}
	if len(children) == 0 {
		children = append(children, infoCard(prefix+"-empty", "비어 있음", "표시할 유닛이 없습니다.", textMuted))
	}
	return children
}

func panel(id string, title string, children ...*ebitenui.Node) *ebitenui.Node {
	return panelWithHeight(id, title, ebitenui.Auto(), children...)
}

func screenPanel(id string, title string, children ...*ebitenui.Node) *ebitenui.Node {
	return panelWithHeight(id, title, ebitenui.Fill(), children...)
}

func panelWithHeight(id string, title string, height ebitenui.Length, children ...*ebitenui.Node) *ebitenui.Node {
	content := []*ebitenui.Node{
		ebitenui.Text(title, ebitenui.Props{
			ID:    id + "-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
	}
	content = append(content, children...)
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Height:          height,
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(panelPadding),
			Gap:             panelGap,
			BackgroundColor: panelBackground,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	}, content...)
}

func panelWithExtra(id string, title string, main []*ebitenui.Node, extra *ebitenui.Node) *ebitenui.Node {
	return panelWithExtraHeight(id, title, ebitenui.Auto(), main, extra)
}

func panelWithExtraHeight(id string, title string, height ebitenui.Length, main []*ebitenui.Node, extra *ebitenui.Node) *ebitenui.Node {
	children := append([]*ebitenui.Node{}, main...)
	if extra != nil {
		children = append(children, extra)
	}
	return panelWithHeight(id, title, height, children...)
}

func button(id string, label string, detail string, selected bool, disabled bool, trigger buttonTrigger, onAction func()) *ebitenui.Node {
	background := panelBackground
	if selected {
		background = selectedColor
	}
	if disabled {
		background = disabledColor
	}
	return ebitenui.InteractiveButton(ebitenui.Props{
		ID: id,
		State: ebitenui.InteractionState{
			Selected: selected,
			Disabled: disabled,
		},
		Handlers: bindButtonHandlers(disabled, trigger, onAction),
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(6),
			Gap:             4,
			BackgroundColor: background,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.TextBlock(label, ebitenui.Props{
			ID:    id + "-label",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textMuted},
		}),
	)
}

func compactButton(id string, label string, detail string, trigger buttonTrigger, onAction func()) *ebitenui.Node {
	return ebitenui.InteractiveButton(ebitenui.Props{
		ID:       id,
		Handlers: bindButtonHandlers(false, trigger, onAction),
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(6),
			Gap:             4,
			BackgroundColor: color.RGBA{R: 34, G: 41, B: 59, A: 255},
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.TextBlock(label, ebitenui.Props{
			ID:    id + "-label",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textMuted},
		}),
	)
}

func bindButtonHandlers(disabled bool, trigger buttonTrigger, onAction func()) ebitenui.EventHandlers {
	if disabled || onAction == nil {
		return ebitenui.EventHandlers{}
	}

	handlers := ebitenui.EventHandlers{}
	switch trigger {
	case buttonTriggerPointerDown:
		handlers.OnPointerDown = func(ctx ebitenui.EventContext) {
			onAction()
		}
	case buttonTriggerPointerHold:
		handlers.OnPointerHold = func(ctx ebitenui.EventContext) {
			onAction()
		}
	case buttonTriggerPointerUp:
		handlers.OnPointerUp = func(ctx ebitenui.EventContext) {
			onAction()
		}
	default:
		handlers.OnClick = func(ctx ebitenui.EventContext) {
			onAction()
		}
	}
	return handlers
}

func infoCard(id string, title string, detail string, textColor color.Color) *ebitenui.Node {
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(6),
			Gap:             4,
			BackgroundColor: color.RGBA{R: 34, G: 41, B: 59, A: 255},
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.TextBlock(title, ebitenui.Props{
			ID:    id + "-title",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Color: textColor},
		}),
	)
}

func fixedWidthPanel(id string, title string, width float64, children ...*ebitenui.Node) *ebitenui.Node {
	content := []*ebitenui.Node{
		ebitenui.Text(title, ebitenui.Props{
			ID:    id + "-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
	}
	content = append(content, children...)
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Px(width),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(panelPadding),
			Gap:             panelGap,
			BackgroundColor: panelBackground,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	}, content...)
}

func fixedWidthColumn(id string, width float64, children ...*ebitenui.Node) *ebitenui.Node {
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:     ebitenui.Px(width),
			Height:    ebitenui.Fill(),
			Direction: ebitenui.Column,
			Gap:       combatColGap,
		},
	}, children...)
}

func resolveLayoutMetrics(model Model) layoutMetrics {
	width := model.ViewportWidth
	height := model.ViewportHeight
	if width <= 0 {
		width = defaultViewportWidth
	}
	if height <= 0 {
		height = defaultViewportHeight
	}
	return layoutMetrics{
		viewportWidth:  width,
		viewportHeight: height,
		contentWidth:   width - (rootPadding * 2),
	}
}

func gridColumnCount(metrics layoutMetrics, preferred int) int {
	if preferred <= 1 {
		return 1
	}
	if metrics.contentWidth < 900 {
		return 1
	}
	return preferred
}

func combatColumnCount(metrics layoutMetrics) int {
	if metrics.contentWidth >= 1140 {
		return 3
	}
	if metrics.contentWidth >= 760 {
		return 2
	}
	return 1
}

func combatColumnWidths(metrics layoutMetrics, columns int) []float64 {
	if columns <= 1 {
		return []float64{maxFloat(metrics.contentWidth-(panelPadding*2), minTileWidth)}
	}
	innerWidth := metrics.contentWidth - (panelPadding * 2) - (combatColGap * float64(columns-1))
	if innerWidth < minTileWidth {
		innerWidth = minTileWidth
	}
	if columns == 2 {
		left := innerWidth * 0.36
		right := innerWidth - left
		return []float64{left, right}
	}
	left := innerWidth * 0.28
	center := innerWidth * 0.44
	right := innerWidth - left - center
	return []float64{left, center, right}
}

func cardGrid(id string, metrics layoutMetrics, columns int, children ...*ebitenui.Node) *ebitenui.Node {
	innerWidth := maxFloat(metrics.contentWidth-(panelPadding*2), minTileWidth)
	cardWidth := innerWidth
	if columns > 1 {
		cardWidth = (innerWidth - (gridGap * float64(columns-1))) / float64(columns)
	}
	return ebitenui.Grid(ebitenui.GridConfig{
		ID:       id,
		Columns:  columns,
		Gap:      gridGap,
		Children: wrapCards(id, cardWidth, children...),
	})
}

func cardGridWithWidth(id string, width float64, columns int, children ...*ebitenui.Node) []*ebitenui.Node {
	innerWidth := maxFloat(width-(panelPadding*2), minTileWidth)
	cardWidth := innerWidth
	if columns > 1 {
		cardWidth = (innerWidth - (gridGap * float64(columns-1))) / float64(columns)
	}
	return []*ebitenui.Node{
		ebitenui.Grid(ebitenui.GridConfig{
			ID:       id,
			Columns:  columns,
			Gap:      gridGap,
			Children: wrapCards(id, cardWidth, children...),
		}),
	}
}

func cardColumnWithWidth(id string, width float64, children ...*ebitenui.Node) *ebitenui.Node {
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Column,
			Gap:       gridGap,
		},
	}, wrapCards(id, maxFloat(width-(panelPadding*2), minTileWidth), children...)...)
}

func wrapCards(id string, width float64, children ...*ebitenui.Node) []*ebitenui.Node {
	wrapped := make([]*ebitenui.Node, 0, len(children))
	safeWidth := width
	if safeWidth > minTileWidth {
		safeWidth -= 1
	}
	for index, child := range children {
		if child == nil {
			continue
		}
		wrapped = append(wrapped, ebitenui.Div(ebitenui.Props{
			ID: fmt.Sprintf("%s-wrap-%d", id, index),
			Style: ebitenui.Style{
				Width:     ebitenui.Px(maxFloat(safeWidth, minTileWidth)),
				Direction: ebitenui.Column,
			},
		}, child))
	}
	return wrapped
}

func persistentScrollView(id string, runtime *ebitenui.Runtime, children ...*ebitenui.Node) *ebitenui.Node {
	offsetKey := id + "-offset"
	offsetY := 0.0
	if runtime != nil {
		offsetY = runtime.NumberValueOrDefault(offsetKey, 0)
	}

	handlers := ebitenui.EventHandlers{}
	if runtime != nil {
		handlers.OnScroll = func(ctx ebitenui.EventContext) {
			if ctx.Runtime == nil {
				return
			}
			maxOffset := maxFloat(0, ctx.Layout.ContentHeight-ctx.Layout.Frame.Height)
			nextOffset := clampFloat(offsetY-(ctx.ScrollY*defaultScrollStep), 0, maxOffset)
			if nextOffset == offsetY {
				return
			}
			ctx.Runtime.SetNumberValue(offsetKey, nextOffset)
		}
	}

	return ebitenui.ScrollView(ebitenui.Props{
		ID:        id,
		Focusable: true,
		Handlers:  handlers,
		Scroll: ebitenui.ScrollState{
			OffsetY: offsetY,
		},
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Height:    ebitenui.Fill(),
			Direction: ebitenui.Column,
		},
	}, children...)
}

func fallback(value string, fallbackValue string) string {
	if strings.TrimSpace(value) == "" {
		return fallbackValue
	}
	return value
}

func nodeKindLabel(kind string) string {
	switch kind {
	case "start":
		return "시작"
	case "normal":
		return "일반"
	case "elite":
		return "엘리트"
	case "boss":
		return "보스"
	case "rest":
		return "휴식"
	default:
		return kind
	}
}

func maxOne(value int) int {
	if value < 1 {
		return 1
	}
	return value
}

func colorForUnit(unit PartyMember) color.Color {
	if unit.Downed {
		return dangerColor
	}
	if unit.HP > 0 && unit.HP*2 <= unit.MaxHP {
		return accentColor
	}
	return successColor
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clampFloat(value, minimum, maximum float64) float64 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
