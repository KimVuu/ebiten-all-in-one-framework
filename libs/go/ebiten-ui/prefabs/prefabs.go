package prefabs

import (
	"fmt"
	"image/color"
	"math"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

var (
	panelBG      = color.RGBA{R: 25, G: 32, B: 42, A: 255}
	panelBGAlt   = color.RGBA{R: 18, G: 24, B: 33, A: 255}
	cardBG       = color.RGBA{R: 31, G: 40, B: 54, A: 255}
	borderMuted  = color.RGBA{R: 85, G: 103, B: 128, A: 255}
	textStrong   = color.RGBA{R: 239, G: 244, B: 250, A: 255}
	textMuted    = color.RGBA{R: 178, G: 188, B: 204, A: 255}
	accentBlue   = color.RGBA{R: 91, G: 162, B: 255, A: 255}
	accentGreen  = color.RGBA{R: 82, G: 205, B: 150, A: 255}
	accentYellow = color.RGBA{R: 255, G: 194, B: 82, A: 255}
)

type PanelConfig struct {
	ID       string
	Title    string
	Width    float64
	Children []*ebitenui.Node
}

type CardConfig struct {
	ID       string
	Title    string
	Width    float64
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
	Actions []DialogAction
}

type HUDBarConfig struct {
	ID      string
	Label   string
	Current int
	Max     int
	Width   float64
	Tint    color.Color
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
	Slots    []InventorySlot
}

type PauseMenuConfig struct {
	ID     string
	Title  string
	Width  float64
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
	Stats       []TooltipStat
}

func Panel(cfg PanelConfig) *ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID: cfg.ID + "-title",
			Style: ebitenui.Style{
				Color: textStrong,
			},
		}))
	}
	children = append(children, cfg.Children...)
	return ebitenui.Section(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:           widthLength(cfg.Width),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(16),
			Gap:             12,
			BackgroundColor: panelBG,
			BorderColor:     borderMuted,
			BorderWidth:     1,
		},
	}, children...)
}

func Card(cfg CardConfig) *ebitenui.Node {
	children := make([]*ebitenui.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID:    cfg.ID + "-title",
			Style: ebitenui.Style{Color: textStrong},
		}))
	}
	children = append(children, cfg.Children...)
	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:           widthLength(cfg.Width),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(12),
			Gap:             10,
			BackgroundColor: cardBG,
			BorderColor:     borderMuted,
			BorderWidth:     1,
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
	children := make([]*ebitenui.Node, 0, len(cfg.Items)+1)
	if cfg.Title != "" {
		children = append(children, ebitenui.Text(cfg.Title, ebitenui.Props{
			ID:    cfg.ID + "-title",
			Style: ebitenui.Style{Color: textStrong},
		}))
	}

	for _, item := range cfg.Items {
		itemChildren := []*ebitenui.Node{
			ebitenui.Text(item.Label, ebitenui.Props{
				ID: item.ID + "-label",
				Style: ebitenui.Style{
					Width: ebitenui.Fill(),
					Color: buttonTextColor(item.State),
				},
			}),
		}
		if item.Hint != "" {
			itemChildren = append(itemChildren, ebitenui.Text(item.Hint, ebitenui.Props{
				ID:    item.ID + "-hint",
				Style: ebitenui.Style{Color: buttonTextColor(item.State)},
			}))
		}

		children = append(children, ebitenui.InteractiveButton(ebitenui.Props{
			ID:    item.ID,
			State: item.State,
			Style: menuButtonStyle(item.State),
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
	children := []*ebitenui.Node{
		ebitenui.TextBlock(cfg.Body, ebitenui.Props{
			ID: cfg.ID + "-body",
			Style: ebitenui.Style{
				Width:      ebitenui.Fill(),
				Color:      textMuted,
				LineHeight: 16,
			},
		}),
	}

	actionNodes := make([]*ebitenui.Node, 0, len(cfg.Actions))
	for _, action := range cfg.Actions {
		actionNodes = append(actionNodes, ebitenui.InteractiveButton(ebitenui.Props{
			ID:    action.ID,
			State: action.State,
			Style: menuButtonStyle(action.State),
		},
			ebitenui.Text(action.Label, ebitenui.Props{
				ID: action.ID + "-label",
				Style: ebitenui.Style{
					Color: buttonTextColor(action.State),
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
		Children: children,
	})
}

func HUDBar(cfg HUDBarConfig) *ebitenui.Node {
	width := cfg.Width
	if width == 0 {
		width = 200
	}
	tint := cfg.Tint
	if tint == nil {
		tint = accentBlue
	}

	fillWidth := width * clampRatio(cfg.Current, cfg.Max)
	return ebitenui.Div(ebitenui.Props{
		ID: cfg.ID,
		Style: ebitenui.Style{
			Width:     ebitenui.Px(width),
			Direction: ebitenui.Column,
			Gap:       6,
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
				Height:          ebitenui.Px(18),
				BackgroundColor: panelBGAlt,
				BorderColor:     borderMuted,
				BorderWidth:     1,
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
					Color: textStrong,
				},
			}),
		),
	)
}

func InventoryGrid(cfg InventoryGridConfig) *ebitenui.Node {
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
				icon = ebitenui.SolidImage(18, 18, accentBlue)
			}
			rowChildren = append(rowChildren, ebitenui.InteractiveButton(ebitenui.Props{
				ID:    slot.ID,
				State: slot.State,
				Style: ebitenui.Style{
					Width:           ebitenui.Px(cellSize),
					Height:          ebitenui.Px(cellSize),
					Direction:       ebitenui.Column,
					Padding:         ebitenui.All(8),
					Gap:             4,
					BackgroundColor: cardBG,
					BorderColor:     borderMuted,
					BorderWidth:     1,
				},
			},
				ebitenui.Image(ebitenui.Props{
					ID:    slot.ID + "-icon",
					Image: icon,
				}),
				ebitenui.Text(slot.Label, ebitenui.Props{
					ID: slot.ID + "-label",
					Style: ebitenui.Style{
						Color: textStrong,
					},
				}),
				ebitenui.Text(fmt.Sprintf("x%d", slot.Quantity), ebitenui.Props{
					ID: slot.ID + "-qty",
					Style: ebitenui.Style{
						Color: textMuted,
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
		Children: children,
	})
}

func SettingsPanel(cfg SettingsPanelConfig) *ebitenui.Node {
	sectionNodes := make([]*ebitenui.Node, 0, len(cfg.Sections))
	for i, section := range cfg.Sections {
		optionNodes := make([]*ebitenui.Node, 0, len(section.Options)+2)
		optionNodes = append(optionNodes,
			ebitenui.Text(section.Title, ebitenui.Props{
				ID:    fmt.Sprintf("%s-section-%d-title", cfg.ID, i),
				Style: ebitenui.Style{Color: textStrong},
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
		Children: []*ebitenui.Node{
			scroll,
		},
	})
}

func Tooltip(cfg TooltipConfig) *ebitenui.Node {
	statRows := make([]*ebitenui.Node, 0, len(cfg.Stats)+2)
	statRows = append(statRows,
		ebitenui.TextBlock(cfg.Description, ebitenui.Props{
			ID: cfg.ID + "-description",
			Style: ebitenui.Style{
				Width:      ebitenui.Fill(),
				Color:      textMuted,
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
		Children: statRows,
	})
}

func buttonTextColor(state ebitenui.InteractionState) color.Color {
	if state.Disabled {
		return color.RGBA{R: 150, G: 156, B: 168, A: 255}
	}
	if state.Selected || state.Focused {
		return color.RGBA{R: 18, G: 24, B: 33, A: 255}
	}
	return textStrong
}

func menuButtonStyle(state ebitenui.InteractionState) ebitenui.Style {
	background := cardBG
	border := borderMuted
	if state.Selected {
		background = accentBlue
		border = accentYellow
	} else if state.Focused {
		background = accentGreen
		border = accentYellow
	} else if state.Disabled {
		background = color.RGBA{R: 50, G: 56, B: 68, A: 255}
		border = color.RGBA{R: 74, G: 80, B: 90, A: 255}
	}
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Row,
		Padding:         ebitenui.All(12),
		Gap:             8,
		BackgroundColor: background,
		BorderColor:     border,
		BorderWidth:     1,
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
