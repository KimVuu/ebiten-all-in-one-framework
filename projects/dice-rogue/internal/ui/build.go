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
	minTileWidth          = 210.0
	defaultViewportWidth  = 1600.0
	defaultViewportHeight = 960.0
)

type layoutMetrics struct {
	viewportWidth  float64
	viewportHeight float64
	contentWidth   float64
}

func BuildDOM(model Model, callbacks Callbacks) *ebitenui.DOM {
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
	}, buildCurrentScreen(model, callbacks, metrics)))

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

func buildCurrentScreen(model Model, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
	switch model.CurrentScreen {
	case "map":
		return buildMapScreen(model.Map, callbacks, metrics)
	case "combat":
		return buildCombatScreen(model.Combat, callbacks, metrics)
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

func buildCombatScreen(model CombatModel, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
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
	visibleDice, extraDice := limitDieViews(model.AvailableDice, 4)
	for _, die := range visibleDice {
		detail := die.Detail
		if die.Forced {
			detail += " / 강제 선택"
		}
		availableChildren = append(availableChildren, compactButton(
			fmt.Sprintf("available-die-%s", die.ID),
			die.Label,
			detail,
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

	selectedChildren := make([]*ebitenui.Node, 0, len(model.SelectedDice))
	for _, die := range model.SelectedDice {
		text := die.Label
		if die.Forced {
			text += " / 강제 선택"
		}
		selectedChildren = append(selectedChildren, infoCard(
			fmt.Sprintf("selected-die-%s", die.ID),
			text,
			die.Detail,
			textMuted,
		))
	}
	if len(selectedChildren) == 0 {
		selectedChildren = append(selectedChildren, infoCard(
			"selected-dice-empty",
			"주사위 선택",
			"턴을 진행하려면 주사위 3개를 선택하세요.",
			textMuted,
		))
	}

	logChildren := make([]*ebitenui.Node, 0, len(model.Logs))
	visibleLogs, extraLogs := limitStrings(model.Logs, 4)
	if len(visibleLogs) == 0 {
		logChildren = append(logChildren, infoCard(
			"combat-log-empty",
			"전투 기록 없음",
			"턴이 끝나면 전투 기록이 여기에 표시됩니다.",
			textMuted,
		))
	} else {
		for index, line := range visibleLogs {
			logChildren = append(logChildren, infoCard(
				fmt.Sprintf("combat-log-%d", index),
				fmt.Sprintf("기록 %d", index+1),
				line,
				textMuted,
			))
		}
	}
	combatColumns := combatColumnCount(metrics)
	columnWidth := combatColumnWidth(metrics, combatColumns)

	leftColumn := fixedWidthColumn("combat-left-column", columnWidth,
		panel("combat-party-panel", "파티", buildUnitCards(model.Party, "party-card")...),
		panel("selected-dice-panel", "선택된 주사위", cardGridWithWidth("selected-dice-grid", columnWidth, 1, selectedChildren...)...),
		button(
			"resolve-turn-button",
			"턴 진행",
			"선택한 주사위 3개를 처리한 뒤 적이 행동합니다.",
			false,
			!model.CanResolve,
			func() {
				if callbacks.OnResolveTurn != nil {
					callbacks.OnResolveTurn()
				}
			},
		),
	)

	centerChildren := []*ebitenui.Node{
		panel("combat-enemy-panel", "적", buildUnitCards(model.Enemies, "enemy-card")...),
		panel("revealed-patterns-panel", "공개된 패턴", cardGridWithWidth("revealed-pattern-grid", columnWidth, 1, revealChildren...)...),
	}
	rightChildren := []*ebitenui.Node{
		panelWithExtra("available-dice-panel", "사용 가능 주사위",
			cardGridWithWidth("available-dice-grid", columnWidth, 2, availableChildren...),
			extraDiceSummary(extraDice),
		),
		panelWithExtra("combat-log-panel", "전투 기록",
			cardGridWithWidth("combat-log-grid", columnWidth, 1, logChildren...),
			extraLogSummary(extraLogs),
		),
	}

	combatRowChildren := []*ebitenui.Node{leftColumn}
	if combatColumns >= 3 {
		combatRowChildren = append(combatRowChildren, fixedWidthColumn("combat-center-column", columnWidth, centerChildren...))
		combatRowChildren = append(combatRowChildren, fixedWidthColumn("combat-right-column", columnWidth, rightChildren...))
	} else if combatColumns == 2 {
		secondColumnChildren := append([]*ebitenui.Node{}, centerChildren...)
		secondColumnChildren = append(secondColumnChildren, rightChildren...)
		combatRowChildren = append(combatRowChildren, fixedWidthColumn("combat-right-column", columnWidth, secondColumnChildren...))
	} else {
		leftColumn.Children = append(leftColumn.Children, centerChildren...)
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
	children := append([]*ebitenui.Node{}, main...)
	if extra != nil {
		children = append(children, extra)
	}
	return panel(id, title, children...)
}

func button(id string, label string, detail string, selected bool, disabled bool, onClick func()) *ebitenui.Node {
	background := panelBackground
	if selected {
		background = selectedColor
	}
	if disabled {
		background = disabledColor
	}
	handlers := ebitenui.EventHandlers{}
	if !disabled && onClick != nil {
		handlers.OnClick = func(ctx ebitenui.EventContext) {
			onClick()
		}
	}
	return ebitenui.InteractiveButton(ebitenui.Props{
		ID: id,
		State: ebitenui.InteractionState{
			Selected: selected,
			Disabled: disabled,
		},
		Handlers: handlers,
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

func compactButton(id string, label string, detail string, onClick func()) *ebitenui.Node {
	handlers := ebitenui.EventHandlers{}
	if onClick != nil {
		handlers.OnClick = func(ctx ebitenui.EventContext) {
			onClick()
		}
	}
	return ebitenui.InteractiveButton(ebitenui.Props{
		ID:       id,
		Handlers: handlers,
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

func extraDiceSummary(extra int) *ebitenui.Node {
	if extra <= 0 {
		return nil
	}
	return ebitenui.TextBlock(
		fmt.Sprintf("아직 %d개의 주사위가 풀에 남아 있습니다.", extra),
		ebitenui.Props{
			ID:    "available-dice-more",
			Style: ebitenui.Style{Color: textMuted},
		},
	)
}

func extraLogSummary(extra int) *ebitenui.Node {
	if extra <= 0 {
		return nil
	}
	return ebitenui.TextBlock(
		fmt.Sprintf("이전 전투 기록 %d줄이 숨겨져 있습니다.", extra),
		ebitenui.Props{
			ID:    "combat-log-more",
			Style: ebitenui.Style{Color: textMuted},
		},
	)
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

func combatColumnWidth(metrics layoutMetrics, columns int) float64 {
	innerWidth := metrics.contentWidth - (panelPadding * 2) - (combatColGap * float64(columns-1))
	if innerWidth < minTileWidth {
		innerWidth = minTileWidth
	}
	return innerWidth / float64(columns)
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

func limitDieViews(dice []DieView, limit int) ([]DieView, int) {
	if len(dice) <= limit {
		return dice, 0
	}
	return dice[:limit], len(dice) - limit
}

func limitStrings(values []string, limit int) ([]string, int) {
	if len(values) <= limit {
		return values, 0
	}
	return values[len(values)-limit:], len(values) - limit
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
