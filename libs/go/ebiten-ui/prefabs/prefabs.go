package prefabs

import (
	"fmt"
	"image/color"
	"math"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

var (
	defaultPrefabTheme = ebitenui.DefaultTheme()
	panelBG            = defaultPrefabTheme.Components.Panel.Background
	panelBGAlt         = defaultPrefabTheme.Components.HUDBar.Track
	cardBG             = defaultPrefabTheme.Components.Card.Background
	borderMuted        = defaultPrefabTheme.Components.Panel.Border
	textStrong         = defaultPrefabTheme.Palette.Text.Strong
	textMuted          = defaultPrefabTheme.Palette.Text.Muted
	accentBlue         = defaultPrefabTheme.Palette.Accent.Primary
	accentGreen        = defaultPrefabTheme.Palette.Accent.Secondary
	accentYellow       = defaultPrefabTheme.Palette.Accent.Warning
)

type PanelConfig struct {
	ID       string
	Title    string
	Width    float64
	Theme    *ebitenui.Theme
	Children []*ebitenui.Node
}

type CardConfig struct {
	ID       string
	Title    string
	Width    float64
	Theme    *ebitenui.Theme
	Children []*ebitenui.Node
}

type StatusRowConfig struct {
	ID    string
	Label string
	Value string
	Icon  ebitenui.ImageSource
}

type MenuItem struct {
	ID    string
	Label string
	Hint  string
	State ebitenui.InteractionState
}

type MenuListConfig struct {
	ID    string
	Title string
	Width float64
	Theme *ebitenui.Theme
	Items []MenuItem
}

type DialogAction struct {
	ID    string
	Label string
	State ebitenui.InteractionState
}

type DialogConfig struct {
	ID      string
	Title   string
	Body    string
	Width   float64
	Theme   *ebitenui.Theme
	Actions []DialogAction
}

type HUDBarConfig struct {
	ID      string
	Label   string
	Current int
	Max     int
	Width   float64
	Tint    color.Color
	Theme   *ebitenui.Theme
}

type InventorySlot struct {
	ID       string
	Label    string
	Quantity int
	Icon     ebitenui.ImageSource
	State    ebitenui.InteractionState
}

type InventoryGridConfig struct {
	ID       string
	Title    string
	Columns  int
	CellSize float64
	Theme    *ebitenui.Theme
	Slots    []InventorySlot
}

type PauseMenuConfig struct {
	ID     string
	Title  string
	Width  float64
	Theme  *ebitenui.Theme
	Items  []MenuItem
	Footer string
}

type SettingOption struct {
	Label string
	Value string
	State ebitenui.InteractionState
}

type SettingSection struct {
	Title       string
	Description string
	Options     []SettingOption
}

type SettingsPanelConfig struct {
	ID       string
	Title    string
	Width    float64
	Height   float64
	Theme    *ebitenui.Theme
	Sections []SettingSection
}

type TooltipStat struct {
	Label string
	Value string
}

type TooltipConfig struct {
	ID          string
	Title       string
	Description string
	Width       float64
	Theme       *ebitenui.Theme
	Stats       []TooltipStat
}

func Panel(cfg PanelConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	panelTheme := theme.Components.Panel
	children := make([]*ebitenui.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID: cfg.ID + "-title",
			Style: ebitenui.Style{
				Color: panelTheme.TitleText,
			},
		}))
	}
	children = append(children, cfg.Children...)
	return ebitenui.Section(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:           widthLength(cfg.Width),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(panelTheme.Padding),
			Gap:             panelTheme.Gap,
			BackgroundColor: panelTheme.Background,
			BorderColor:     panelTheme.Border,
			BorderWidth:     panelTheme.BorderWidth,
		},
	}, children...)
}

func Card(cfg CardConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	cardTheme := theme.Components.Card
	children := make([]*ebitenui.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID:    cfg.ID + "-title",
			Style: ebitenui.Style{Color: cardTheme.TitleText},
		}))
	}
	children = append(children, cfg.Children...)
	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:           widthLength(cfg.Width),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(cardTheme.Padding),
			Gap:             cardTheme.Gap,
			BackgroundColor: cardTheme.Background,
			BorderColor:     cardTheme.Border,
			BorderWidth:     cardTheme.BorderWidth,
		},
	}, children...)
}

func StatusRow(cfg StatusRowConfig) *ebitenui.Node {
	children := make([]*ebitenui.Node, 0, 3)
	if cfg.Icon.Width > 0 || cfg.Icon.Height > 0 || cfg.Icon.Image != nil || cfg.Icon.Fill != nil {
		children = append(children, ebitenui.Image(ebitenui.Props{
			ID:    cfg.ID + "-icon",
			Image: cfg.Icon,
			Style: ebitenui.Style{
				Width:  ebitenui.Px(cfg.Icon.Width),
				Height: ebitenui.Px(cfg.Icon.Height),
			},
		}))
	}
	children = append(children,
		ebitenui.Text(cfg.Label, ebitenui.Props{
			ID: cfg.ID + "-label",
			Style: ebitenui.Style{
				Width: ebitenui.Fill(),
				Color: textMuted,
			},
		}),
		ebitenui.Text(cfg.Value, ebitenui.Props{
			ID: cfg.ID + "-value",
			Style: ebitenui.Style{
				Color: textStrong,
			},
		}),
	)

	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Row,
			Gap:       8,
		},
	}, children...)
}

func MenuList(cfg MenuListConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	children := make([]*ebitenui.Node, 0, len(cfg.Items)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID:    cfg.ID + "-title",
			Style: ebitenui.Style{Color: theme.Components.Panel.TitleText},
		}))
	}

	for _, item := range cfg.Items {
		itemChildren := []*ebitenui.Node{
			ebitenui.Text(item.Label, ebitenui.Props{
				ID: item.ID + "-label",
				Style: ebitenui.Style{
					Width: ebitenui.Fill(),
					Color: buttonTextColor(theme, item.State),
				},
			}),
		}
		if item.Hint != "" {
			itemChildren = append(itemChildren, ebitenui.Text(item.Hint, ebitenui.Props{
				ID:    item.ID + "-hint",
				Style: ebitenui.Style{Color: buttonTextColor(theme, item.State)},
			}))
		}

		children = append(children, ebitenui.InteractiveButton(ebitenui.Props{
			ID:    item.ID,
			State: item.State,
			Style: menuButtonStyle(theme, item.State),
		}, itemChildren...))
	}

	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:     widthLength(cfg.Width),
			Direction: ebitenui.Column,
			Gap:       8,
		},
	}, children...)
}

func Dialog(cfg DialogConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	children := []*ebitenui.Node{
		ebitenui.TextBlock(cfg.Body, ebitenui.Props{
			ID: cfg.ID + "-body",
			Style: ebitenui.Style{
				Width:      ebitenui.Fill(),
				Color:      theme.Components.Dialog.BodyText,
				LineHeight: 16,
			},
		}),
	}

	actionNodes := make([]*ebitenui.Node, 0, len(cfg.Actions))
	for _, action := range cfg.Actions {
		actionNodes = append(actionNodes, ebitenui.InteractiveButton(ebitenui.Props{
			ID:    action.ID,
			State: action.State,
			Style: menuButtonStyle(theme, action.State),
		},
			ebitenui.Text(action.Label, ebitenui.Props{
				ID: action.ID + "-label",
				Style: ebitenui.Style{
					Color: buttonTextColor(theme, action.State),
				},
			}),
		))
	}

	children = append(children, ebitenui.Div(ebitenui.Props{
		ID: cfg.ID + "-actions",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Row,
			Gap:       10,
		},
	}, actionNodes...))

	return Panel(PanelConfig{
		ID:       cfg.ID,
		Title:    cfg.Title,
		Width:    cfg.Width,
		Theme:    cfg.Theme,
		Children: children,
	})
}

func HUDBar(cfg HUDBarConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	hudTheme := theme.Components.HUDBar
	width := cfg.Width
	if width == 0 {
		width = 200
	}
	tint := cfg.Tint
	if tint == nil {
		tint = hudTheme.Fill
	}

	fillWidth := width * clampRatio(cfg.Current, cfg.Max)
	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:     ebitenui.Px(width),
			Direction: ebitenui.Column,
			Gap:       hudTheme.Gap,
		},
	},
		StatusRow(StatusRowConfig{
			ID:    cfg.ID + "-status",
			Label: cfg.Label,
			Value: fmt.Sprintf("%d/%d", cfg.Current, cfg.Max),
		}),
		ebitenui.Stack(ebitenui.Props{
			ID: cfg.ID + "-track",
			Style: ebitenui.Style{
				Width:           ebitenui.Px(width),
				Height:          ebitenui.Px(hudTheme.TrackHeight),
				BackgroundColor: hudTheme.Track,
				BorderColor:     hudTheme.Border,
				BorderWidth:     hudTheme.BorderWidth,
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: cfg.ID + "-fill",
				Style: ebitenui.Style{
					Width:           ebitenui.Px(fillWidth),
					Height:          ebitenui.Fill(),
					BackgroundColor: tint,
				},
			}),
			ebitenui.Text(fmt.Sprintf("%s %d/%d", cfg.Label, cfg.Current, cfg.Max), ebitenui.Props{
				ID: cfg.ID + "-text",
				Style: ebitenui.Style{
					Color: hudTheme.Text,
				},
			}),
		),
	)
}

func InventoryGrid(cfg InventoryGridConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	gridTheme := theme.Components.InventoryGrid
	columns := cfg.Columns
	if columns <= 0 {
		columns = 4
	}
	cellSize := cfg.CellSize
	if cellSize == 0 {
		cellSize = 48
	}

	rows := make([]*ebitenui.Node, 0, int(math.Ceil(float64(len(cfg.Slots))/float64(columns))))
	for rowIndex := 0; rowIndex*columns < len(cfg.Slots); rowIndex++ {
		rowChildren := make([]*ebitenui.Node, 0, columns)
		for column := 0; column < columns; column++ {
			index := rowIndex*columns + column
			if index >= len(cfg.Slots) {
				break
			}
			slot := cfg.Slots[index]
			icon := slot.Icon
			if icon.Width == 0 && icon.Height == 0 && icon.Fill == nil && icon.Image == nil {
				icon = ebitenui.SolidImage(18, 18, gridTheme.IconFill)
			}
			rowChildren = append(rowChildren, ebitenui.InteractiveButton(ebitenui.Props{
				ID:    slot.ID,
				State: slot.State,
				Style: ebitenui.Style{
					Width:           ebitenui.Px(cellSize),
					Height:          ebitenui.Px(cellSize),
					Direction:       ebitenui.Column,
					Padding:         ebitenui.All(gridTheme.SlotPadding),
					Gap:             gridTheme.SlotGap,
					BackgroundColor: gridTheme.SlotBackground,
					BorderColor:     gridTheme.SlotBorder,
					BorderWidth:     gridTheme.SlotBorderWidth,
				},
			},
				ebitenui.Image(ebitenui.Props{
					ID:    slot.ID + "-icon",
					Image: icon,
				}),
				ebitenui.Text(slot.Label, ebitenui.Props{
					ID: slot.ID + "-label",
					Style: ebitenui.Style{
						Color: gridTheme.SlotText,
					},
				}),
				ebitenui.Text(fmt.Sprintf("x%d", slot.Quantity), ebitenui.Props{
					ID: slot.ID + "-qty",
					Style: ebitenui.Style{
						Color: gridTheme.SlotMuted,
					},
				}),
			))
		}
		rows = append(rows, ebitenui.Div(ebitenui.Props{
			ID: cfg.ID + fmt.Sprintf("-row-%d", rowIndex),
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       8,
			},
		}, rowChildren...))
	}

	return Panel(PanelConfig{
		ID:    cfg.ID,
		Title: cfg.Title,
		Width: float64(columns)*cellSize + float64(columns-1)*8 + 32,
		Theme: cfg.Theme,
		Children: []*ebitenui.Node{
			ebitenui.Div(ebitenui.Props{
				ID: cfg.ID + "-grid",
				Style: ebitenui.Style{
					Width:     ebitenui.Fill(),
					Direction: ebitenui.Column,
					Gap:       8,
				},
			}, rows...),
		},
	})
}

func PauseMenu(cfg PauseMenuConfig) *ebitenui.Node {
	children := []*ebitenui.Node{
		MenuList(MenuListConfig{
			ID:    cfg.ID + "-menu",
			Width: cfg.Width,
			Theme: cfg.Theme,
			Items: cfg.Items,
		}),
	}
	if cfg.Footer != "" {
		children = append(children, ebitenui.Text(cfg.Footer, ebitenui.Props{
			ID:    cfg.ID + "-footer",
			Style: ebitenui.Style{Color: textMuted},
		}))
	}

	return Panel(PanelConfig{
		ID:       cfg.ID,
		Title:    cfg.Title,
		Width:    cfg.Width,
		Theme:    cfg.Theme,
		Children: children,
	})
}

func SettingsPanel(cfg SettingsPanelConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	sectionNodes := make([]*ebitenui.Node, 0, len(cfg.Sections))
	for i, section := range cfg.Sections {
		optionNodes := make([]*ebitenui.Node, 0, len(section.Options)+2)
		optionNodes = append(optionNodes,
			ebitenui.Text(section.Title, ebitenui.Props{
				ID:    fmt.Sprintf("%s-section-%d-title", cfg.ID, i),
				Style: ebitenui.Style{Color: theme.Components.Card.TitleText},
			}),
			ebitenui.TextBlock(section.Description, ebitenui.Props{
				ID: fmt.Sprintf("%s-section-%d-copy", cfg.ID, i),
				Style: ebitenui.Style{
					Width:      ebitenui.Fill(),
					Color:      textMuted,
					LineHeight: 16,
				},
			}),
		)

		for j, option := range section.Options {
			optionNodes = append(optionNodes, StatusRow(StatusRowConfig{
				ID:    fmt.Sprintf("%s-section-%d-option-%d", cfg.ID, i, j),
				Label: option.Label,
				Value: option.Value,
			}))
		}

		sectionNodes = append(sectionNodes, Card(CardConfig{
			ID:       fmt.Sprintf("%s-section-%d", cfg.ID, i),
			Theme:    cfg.Theme,
			Children: optionNodes,
		}))
	}

	scroll := ebitenui.ScrollView(ebitenui.Props{
		ID: cfg.ID + "-scroll",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Height:    ebitenui.Px(cfg.Height - 80),
			Direction: ebitenui.Column,
			Gap:       10,
		},
	}, sectionNodes...)

	return Panel(PanelConfig{
		ID:    cfg.ID,
		Title: cfg.Title,
		Width: cfg.Width,
		Theme: cfg.Theme,
		Children: []*ebitenui.Node{
			scroll,
		},
	})
}

func Tooltip(cfg TooltipConfig) *ebitenui.Node {
	theme := ebitenui.ResolveTheme(cfg.Theme)
	statRows := make([]*ebitenui.Node, 0, len(cfg.Stats)+2)
	statRows = append(statRows,
		ebitenui.TextBlock(cfg.Description, ebitenui.Props{
			ID: cfg.ID + "-description",
			Style: ebitenui.Style{
				Width:      ebitenui.Fill(),
				Color:      theme.Components.Tooltip.BodyText,
				LineHeight: 16,
			},
		}),
	)

	for i, stat := range cfg.Stats {
		statRows = append(statRows, StatusRow(StatusRowConfig{
			ID:    fmt.Sprintf("%s-stat-%d", cfg.ID, i),
			Label: stat.Label,
			Value: stat.Value,
		}))
	}

	return Card(CardConfig{
		ID:       cfg.ID,
		Title:    cfg.Title,
		Width:    cfg.Width,
		Theme:    cfg.Theme,
		Children: statRows,
	})
}

func buttonTextColor(theme ebitenui.Theme, state ebitenui.InteractionState) color.Color {
	buttonTheme := theme.Components.MenuButton
	switch {
	case state.Disabled:
		return buttonTheme.Disabled.Text
	case state.Selected:
		return buttonTheme.Selected.Text
	case state.Focused:
		return buttonTheme.Focused.Text
	default:
		return buttonTheme.Default.Text
	}
}

func menuButtonStyle(theme ebitenui.Theme, state ebitenui.InteractionState) ebitenui.Style {
	buttonTheme := theme.Components.MenuButton
	colors := buttonTheme.Default
	switch {
	case state.Disabled:
		colors = buttonTheme.Disabled
	case state.Selected:
		colors = buttonTheme.Selected
	case state.Focused:
		colors = buttonTheme.Focused
	}
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Row,
		Padding:         ebitenui.All(buttonTheme.Padding),
		Gap:             buttonTheme.Gap,
		BackgroundColor: colors.Background,
		BorderColor:     colors.Border,
		BorderWidth:     buttonTheme.BorderWidth,
	}
}

func clampRatio(current, max int) float64 {
	if max <= 0 {
		return 0
	}
	ratio := float64(current) / float64(max)
	return math.Max(0, math.Min(1, ratio))
}

func widthLength(width float64) ebitenui.Length {
	if width > 0 {
		return ebitenui.Px(width)
	}
	return ebitenui.Fill()
}
