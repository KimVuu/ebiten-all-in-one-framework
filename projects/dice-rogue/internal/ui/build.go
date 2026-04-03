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

func BuildDOM(model Model, callbacks Callbacks) *ebitenui.DOM {
	root := ebitenui.Div(ebitenui.Props{
		ID: "dice-rogue-root",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Height:          ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(20),
			Gap:             14,
			BackgroundColor: rootBackground,
		},
	},
		buildHeader(model),
		buildPartyRoster(model.PartyRoster),
		buildCurrentScreen(model, callbacks),
	)
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

func buildPartyRoster(roster []PartyMember) *ebitenui.Node {
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
		children = append(children, ebitenui.Grid(ebitenui.GridConfig{
			ID:       "party-roster-grid",
			Columns:  3,
			Gap:      8,
			Children: cards,
		}))
	}
	return panel("party-roster-panel", "Roster", children...)
}

func buildCurrentScreen(model Model, callbacks Callbacks) *ebitenui.Node {
	switch model.CurrentScreen {
	case "map":
		return buildMapScreen(model.Map, callbacks)
	case "combat":
		return buildCombatScreen(model.Combat, callbacks)
	case "outcome":
		return buildOutcomeScreen(model.Outcome, callbacks)
	default:
		return buildPartySelectionScreen(model.PartySelection, callbacks)
	}
}

func buildPartySelectionScreen(model PartySelectionModel, callbacks Callbacks) *ebitenui.Node {
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
	children = append(children, ebitenui.Grid(ebitenui.GridConfig{
		ID:       "party-selection-grid",
		Columns:  2,
		Gap:      10,
		Children: gridChildren,
	}))
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

func buildMapScreen(model MapModel, callbacks Callbacks) *ebitenui.Node {
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
	children = append(children, ebitenui.Grid(ebitenui.GridConfig{
		ID:       "map-node-grid",
		Columns:  2,
		Gap:      10,
		Children: gridChildren,
	}))
	return panel("map-screen", "Act Map", children...)
}

func buildCombatScreen(model CombatModel, callbacks Callbacks) *ebitenui.Node {
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

	topRow := ebitenui.Div(ebitenui.Props{
		ID: "combat-top-row",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Row,
			Gap:       10,
		},
	},
		fixedWidthPanel("combat-party-panel", "Party", 560, buildUnitList(model.Party, "party-card")...),
		fixedWidthPanel("combat-enemy-panel", "Enemies", 560, buildUnitList(model.Enemies, "enemy-card")...),
	)
	children = append(children, topRow)

	revealChildren := []*ebitenui.Node{}
	if len(model.RevealedPatterns) == 0 {
		revealChildren = append(revealChildren, ebitenui.TextBlock("No revealed enemy intent.", ebitenui.Props{
			ID:    "revealed-patterns-empty",
			Style: ebitenui.Style{Color: textMuted},
		}))
	} else {
		for index, text := range model.RevealedPatterns {
			revealChildren = append(revealChildren, ebitenui.Text(text, ebitenui.Props{
				ID:    fmt.Sprintf("revealed-pattern-%d", index),
				Style: ebitenui.Style{Color: accentColor},
			}))
		}
	}
	children = append(children, panel("revealed-patterns-panel", "Revealed Patterns", revealChildren...))

	availableChildren := make([]*ebitenui.Node, 0, len(model.AvailableDice))
	for _, die := range model.AvailableDice {
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
	children = append(children, panel("available-dice-panel", "Available Dice",
		ebitenui.Grid(ebitenui.GridConfig{
			ID:       "available-dice-grid",
			Columns:  3,
			Gap:      8,
			Children: availableChildren,
		}),
	))

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
	children = append(children, panel("selected-dice-panel", "Selected Dice",
		ebitenui.Grid(ebitenui.GridConfig{
			ID:       "selected-dice-grid",
			Columns:  3,
			Gap:      8,
			Children: selectedChildren,
		}),
	))

	logChildren := make([]*ebitenui.Node, 0, len(model.Logs))
	if len(model.Logs) == 0 {
		logChildren = append(logChildren, infoCard(
			"combat-log-empty",
			"No logs yet",
			"Turn logs appear here after combat resolves.",
			textMuted,
		))
	} else {
		for index, line := range model.Logs {
			logChildren = append(logChildren, infoCard(
				fmt.Sprintf("combat-log-%d", index),
				fmt.Sprintf("Log %d", index+1),
				line,
				textMuted,
			))
		}
	}
	children = append(children, panel("combat-log-panel", "Log",
		ebitenui.Grid(ebitenui.GridConfig{
			ID:       "combat-log-grid",
			Columns:  2,
			Gap:      8,
			Children: logChildren,
		}),
	))
	children = append(children, button(
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
	))

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
			Padding:         ebitenui.All(12),
			Gap:             8,
			BackgroundColor: panelBackground,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	}, content...)
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
			Padding:         ebitenui.All(10),
			Gap:             4,
			BackgroundColor: background,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.Text(label, ebitenui.Props{
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
			Padding:         ebitenui.All(8),
			Gap:             4,
			BackgroundColor: color.RGBA{R: 34, G: 41, B: 59, A: 255},
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.Text(label, ebitenui.Props{
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
			Padding:         ebitenui.All(8),
			Gap:             4,
			BackgroundColor: color.RGBA{R: 34, G: 41, B: 59, A: 255},
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	},
		ebitenui.Text(title, ebitenui.Props{
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
			Padding:         ebitenui.All(12),
			Gap:             8,
			BackgroundColor: panelBackground,
			BorderColor:     borderColor,
			BorderWidth:     1,
		},
	}, content...)
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
