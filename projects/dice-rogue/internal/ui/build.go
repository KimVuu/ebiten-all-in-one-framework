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
	rootPadding  = 12.0
	rootGap      = 10.0
	panelPadding = 8.0
	panelGap     = 6.0
	gridGap      = 8.0
	combatColGap = 10.0
	minTileWidth = 180.0
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
	children = append(children, buildCurrentScreen(model, callbacks, metrics))

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
	return panel("screen-header", "Dice Rogue",
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
		ebitenui.Text("Party", ebitenui.Props{
			ID:    "party-roster-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
	}
	if len(roster) == 0 {
		children = append(children, ebitenui.TextBlock("Select three members to start the act.", ebitenui.Props{
			ID:    "party-roster-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	} else {
		cards := make([]*ebitenui.Node, 0, len(roster))
		for _, member := range roster {
			label := fmt.Sprintf("%s (%s) %d/%d", member.Name, member.Role, member.HP, member.MaxHP)
			if member.Downed {
				label += " / downed"
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
	return panel("party-roster-panel", "Roster", children...)
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
		ebitenui.TextBlock(fmt.Sprintf("Selected %d / 3", model.SelectedCount), ebitenui.Props{
			ID:    "party-selection-count",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}
	gridChildren := make([]*ebitenui.Node, 0, len(model.Candidates))
	for _, candidate := range model.Candidates {
		label := fmt.Sprintf("%s / %s / HP %d", candidate.Name, candidate.Role, candidate.MaxHP)
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
		"Start Run",
		"Open the first act map with the selected party.",
		false,
		!model.CanStart,
		func() {
			if callbacks.OnStartRun != nil {
				callbacks.OnStartRun()
			}
		},
	))
	return panel("party-selection-screen", "Party Selection", children...)
}

func buildMapScreen(model MapModel, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.TextBlock("Choose the next route in the fixed first act.", ebitenui.Props{
			ID:    "map-copy",
			Style: ebitenui.Style{Color: textMuted},
		}),
		ebitenui.Text(fmt.Sprintf("Current node: %s", fallback(model.CurrentNodeID, "start")), ebitenui.Props{
			ID:    "map-current-node",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}
	gridChildren := make([]*ebitenui.Node, 0, len(model.Nodes))
	for _, node := range model.Nodes {
		gridChildren = append(gridChildren, button(
			fmt.Sprintf("map-node-%s", node.ID),
			fmt.Sprintf("%s / %s", strings.Title(node.Kind), node.Name),
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
	return panel("map-screen", "Act Map", children...)
}

func buildCombatScreen(model CombatModel, callbacks Callbacks, metrics layoutMetrics) *ebitenui.Node {
	children := []*ebitenui.Node{
		ebitenui.Text(fmt.Sprintf("Turn %d", maxOne(model.Turn)), ebitenui.Props{
			ID:    "combat-turn",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Text(fmt.Sprintf("Ally defense %d / Enemy defense %d / Damage +%d%%", model.AllyDefense, model.EnemyDefense, model.DamageBoost), ebitenui.Props{
			ID:    "combat-defense-summary",
			Style: ebitenui.Style{Color: textMuted},
		}),
	}

	revealChildren := []*ebitenui.Node{}
	if len(model.RevealedPatterns) == 0 {
		revealChildren = append(revealChildren, ebitenui.TextBlock("No revealed enemy intent.", ebitenui.Props{
			ID:    "revealed-patterns-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	} else {
		for index, text := range model.RevealedPatterns {
			revealChildren = append(revealChildren, infoCard(
				fmt.Sprintf("revealed-pattern-%d", index),
				fmt.Sprintf("Enemy %d", index+1),
				text,
				accentColor,
			))
		}
	}

	availableChildren := make([]*ebitenui.Node, 0, len(model.AvailableDice))
	visibleDice, extraDice := limitDieViews(model.AvailableDice, 6)
	for _, die := range visibleDice {
		detail := die.Detail
		if die.Forced {
			detail += " / forced"
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
			"No choices",
			"No manual choices left in the current pool.",
			textMuted,
		))
	}

	selectedChildren := make([]*ebitenui.Node, 0, len(model.SelectedDice))
	for _, die := range model.SelectedDice {
		text := die.Label
		if die.Forced {
			text += " / forced"
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
			"Select Dice",
			"Select three dice to resolve the turn.",
			textMuted,
		))
	}

	logChildren := make([]*ebitenui.Node, 0, len(model.Logs))
	visibleLogs, extraLogs := limitStrings(model.Logs, 4)
	if len(visibleLogs) == 0 {
		logChildren = append(logChildren, infoCard(
			"combat-log-empty",
			"No logs yet",
			"Turn logs appear here after combat resolves.",
			textMuted,
		))
	} else {
		for index, line := range visibleLogs {
			logChildren = append(logChildren, infoCard(
				fmt.Sprintf("combat-log-%d", index),
				fmt.Sprintf("Log %d", index+1),
				line,
				textMuted,
			))
		}
	}
	combatColumns := combatColumnCount(metrics)
	columnWidth := combatColumnWidth(metrics, combatColumns)

	leftColumn := fixedWidthColumn("combat-left-column", columnWidth,
		panel("combat-party-panel", "Party", buildUnitCards(model.Party, "party-card")...),
		panel("selected-dice-panel", "Selected Dice", cardGridWithWidth("selected-dice-grid", columnWidth, 1, selectedChildren...)...),
		button(
			"resolve-turn-button",
			"Resolve Turn",
			"Resolve the selected three dice, then let enemies act.",
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
		panel("combat-enemy-panel", "Enemies", buildUnitCards(model.Enemies, "enemy-card")...),
		panel("revealed-patterns-panel", "Revealed Patterns", cardGridWithWidth("revealed-pattern-grid", columnWidth, 1, revealChildren...)...),
	}
	rightChildren := []*ebitenui.Node{
		panelWithExtra("available-dice-panel", "Available Dice",
			cardGridWithWidth("available-dice-grid", columnWidth, 2, availableChildren...),
			extraDiceSummary(extraDice),
		),
		panelWithExtra("combat-log-panel", "Log",
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
			Direction: ebitenui.Row,
			Gap:       combatColGap,
		},
	}, combatRowChildren...))

	return panel("combat-screen", fallback(model.EncounterName, "Combat"), children...)
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
			"Continue",
			"Return to the act map.",
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
			"Restart",
			"Start a fresh run from party selection.",
			false,
			false,
			func() {
				if callbacks.OnRestart != nil {
					callbacks.OnRestart()
				}
			},
		))
	}
	return panel("outcome-screen", fallback(model.Title, "Outcome"), children...)
}

func buildUnitList(units []PartyMember, prefix string) []*ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(units))
	for _, unit := range units {
		label := fmt.Sprintf("%s / %s / %d/%d", unit.Name, unit.Role, unit.HP, unit.MaxHP)
		if unit.Downed {
			label += " / downed"
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
		children = append(children, ebitenui.TextBlock("No units available.", ebitenui.Props{
			ID:    prefix + "-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	}
	return children
}

func buildUnitCards(units []PartyMember, prefix string) []*ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(units))
	for _, unit := range units {
		label := fmt.Sprintf("%s / %d/%d", unit.Role, unit.HP, unit.MaxHP)
		if unit.Downed {
			label += " / downed"
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
		children = append(children, infoCard(prefix+"-empty", "Empty", "No units available.", textMuted))
	}
	return children
}

func panel(id string, title string, children ...*ebitenui.Node) *ebitenui.Node {
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
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Color: textMuted},
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
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Color: textMuted},
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
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.TextBlock(detail, ebitenui.Props{
			ID:    id + "-detail",
			Style: ebitenui.Style{Color: textColor},
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
		fmt.Sprintf("%d more dice are still in the pool.", extra),
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
		fmt.Sprintf("%d older log lines are hidden.", extra),
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
		width = 1280
	}
	if height <= 0 {
		height = 720
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
	for index, child := range children {
		if child == nil {
			continue
		}
		wrapped = append(wrapped, ebitenui.Div(ebitenui.Props{
			ID: fmt.Sprintf("%s-wrap-%d", id, index),
			Style: ebitenui.Style{
				Width:     ebitenui.Px(maxFloat(width, minTileWidth)),
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
