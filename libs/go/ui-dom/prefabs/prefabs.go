package prefabs

import (
	"fmt"
	"image/color"
	"math"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
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
	Children []*uidom.Node
}

type CardConfig struct {
	ID       string
	Title    string
	Width    float64
	Children []*uidom.Node
}

type StatusRowConfig struct {
	ID    string
	Label string
	Value string
	Icon  uidom.ImageSource
}

type MenuItem struct {
	ID    string
	Label string
	Hint  string
	State uidom.InteractionState
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
	State uidom.InteractionState
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
	Icon     uidom.ImageSource
	State    uidom.InteractionState
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
	State uidom.InteractionState
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

func Panel(cfg PanelConfig) *uidom.Node {
	children := make([]*uidom.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, uidom.Text(cfg.Title, uidom.Props{
			ID: cfg.ID + "-title",
			Style: uidom.Style{
				Color: textStrong,
			},
		}))
	}
	children = append(children, cfg.Children...)
	return uidom.Section(uidom.Props{
		ID: cfg.ID,
		Style: uidom.Style{
			Width:           widthLength(cfg.Width),
			Direction:       uidom.Column,
			Padding:         uidom.All(16),
			Gap:             12,
			BackgroundColor: panelBG,
			BorderColor:     borderMuted,
			BorderWidth:     1,
		},
	}, children...)
}

func Card(cfg CardConfig) *uidom.Node {
	children := make([]*uidom.Node, 0, len(cfg.Children)+1)
	if cfg.Title != "" {
		children = append(children, uidom.Text(cfg.Title, uidom.Props{
			ID:    cfg.ID + "-title",
			Style: uidom.Style{Color: textStrong},
		}))
	}
	children = append(children, cfg.Children...)
	return uidom.Div(uidom.Props{
		ID: cfg.ID,
		Style: uidom.Style{
			Width:           widthLength(cfg.Width),
			Direction:       uidom.Column,
			Padding:         uidom.All(12),
			Gap:             10,
			BackgroundColor: cardBG,
			BorderColor:     borderMuted,
			BorderWidth:     1,
		},
	}, children...)
}

func StatusRow(cfg StatusRowConfig) *uidom.Node {
	children := make([]*uidom.Node, 0, 3)
	if cfg.Icon.Width > 0 || cfg.Icon.Height > 0 || cfg.Icon.Image != nil || cfg.Icon.Fill != nil {
		children = append(children, uidom.Image(uidom.Props{
			ID:    cfg.ID + "-icon",
			Image: cfg.Icon,
			Style: uidom.Style{
				Width:  uidom.Px(cfg.Icon.Width),
				Height: uidom.Px(cfg.Icon.Height),
			},
		}))
	}
	children = append(children,
		uidom.Text(cfg.Label, uidom.Props{
			ID: cfg.ID + "-label",
			Style: uidom.Style{
				Width: uidom.Fill(),
				Color: textMuted,
			},
		}),
		uidom.Text(cfg.Value, uidom.Props{
			ID: cfg.ID + "-value",
			Style: uidom.Style{
				Color: textStrong,
			},
		}),
	)

	return uidom.Div(uidom.Props{
		ID: cfg.ID,
		Style: uidom.Style{
			Width:     uidom.Fill(),
			Direction: uidom.Row,
			Gap:       8,
		},
	}, children...)
}

func MenuList(cfg MenuListConfig) *uidom.Node {
	children := make([]*uidom.Node, 0, len(cfg.Items)+1)
	if cfg.Title != "" {
		children = append(children, uidom.Text(cfg.Title, uidom.Props{
			ID:    cfg.ID + "-title",
			Style: uidom.Style{Color: textStrong},
		}))
	}

	for _, item := range cfg.Items {
		itemChildren := []*uidom.Node{
			uidom.Text(item.Label, uidom.Props{
				ID: item.ID + "-label",
				Style: uidom.Style{
					Width: uidom.Fill(),
					Color: buttonTextColor(item.State),
				},
			}),
		}
		if item.Hint != "" {
			itemChildren = append(itemChildren, uidom.Text(item.Hint, uidom.Props{
				ID:    item.ID + "-hint",
				Style: uidom.Style{Color: buttonTextColor(item.State)},
			}))
		}

		children = append(children, uidom.InteractiveButton(uidom.Props{
			ID:    item.ID,
			State: item.State,
			Style: menuButtonStyle(item.State),
		}, itemChildren...))
	}

	return uidom.Div(uidom.Props{
		ID: cfg.ID,
		Style: uidom.Style{
			Width:     widthLength(cfg.Width),
			Direction: uidom.Column,
			Gap:       8,
		},
	}, children...)
}

func Dialog(cfg DialogConfig) *uidom.Node {
	children := []*uidom.Node{
		uidom.TextBlock(cfg.Body, uidom.Props{
			ID: cfg.ID + "-body",
			Style: uidom.Style{
				Width:      uidom.Fill(),
				Color:      textMuted,
				LineHeight: 16,
			},
		}),
	}

	actionNodes := make([]*uidom.Node, 0, len(cfg.Actions))
	for _, action := range cfg.Actions {
		actionNodes = append(actionNodes, uidom.InteractiveButton(uidom.Props{
			ID:    action.ID,
			State: action.State,
			Style: menuButtonStyle(action.State),
		},
			uidom.Text(action.Label, uidom.Props{
				ID: action.ID + "-label",
				Style: uidom.Style{
					Color: buttonTextColor(action.State),
				},
			}),
		))
	}

	children = append(children, uidom.Div(uidom.Props{
		ID: cfg.ID + "-actions",
		Style: uidom.Style{
			Width:     uidom.Fill(),
			Direction: uidom.Row,
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

func HUDBar(cfg HUDBarConfig) *uidom.Node {
	width := cfg.Width
	if width == 0 {
		width = 200
	}
	tint := cfg.Tint
	if tint == nil {
		tint = accentBlue
	}

	fillWidth := width * clampRatio(cfg.Current, cfg.Max)
	return uidom.Div(uidom.Props{
		ID: cfg.ID,
		Style: uidom.Style{
			Width:     uidom.Px(width),
			Direction: uidom.Column,
			Gap:       6,
		},
	},
		StatusRow(StatusRowConfig{
			ID:    cfg.ID + "-status",
			Label: cfg.Label,
			Value: fmt.Sprintf("%d/%d", cfg.Current, cfg.Max),
		}),
		uidom.Stack(uidom.Props{
			ID: cfg.ID + "-track",
			Style: uidom.Style{
				Width:           uidom.Px(width),
				Height:          uidom.Px(18),
				BackgroundColor: panelBGAlt,
				BorderColor:     borderMuted,
				BorderWidth:     1,
			},
		},
			uidom.Div(uidom.Props{
				ID: cfg.ID + "-fill",
				Style: uidom.Style{
					Width:           uidom.Px(fillWidth),
					Height:          uidom.Fill(),
					BackgroundColor: tint,
				},
			}),
			uidom.Text(fmt.Sprintf("%s %d/%d", cfg.Label, cfg.Current, cfg.Max), uidom.Props{
				ID: cfg.ID + "-text",
				Style: uidom.Style{
					Color: textStrong,
				},
			}),
		),
	)
}

func InventoryGrid(cfg InventoryGridConfig) *uidom.Node {
	columns := cfg.Columns
	if columns <= 0 {
		columns = 4
	}
	cellSize := cfg.CellSize
	if cellSize == 0 {
		cellSize = 48
	}

	rows := make([]*uidom.Node, 0, int(math.Ceil(float64(len(cfg.Slots))/float64(columns))))
	for rowIndex := 0; rowIndex*columns < len(cfg.Slots); rowIndex++ {
		rowChildren := make([]*uidom.Node, 0, columns)
		for column := 0; column < columns; column++ {
			index := rowIndex*columns + column
			if index >= len(cfg.Slots) {
				break
			}
			slot := cfg.Slots[index]
			icon := slot.Icon
			if icon.Width == 0 && icon.Height == 0 && icon.Fill == nil && icon.Image == nil {
				icon = uidom.SolidImage(18, 18, accentBlue)
			}
			rowChildren = append(rowChildren, uidom.InteractiveButton(uidom.Props{
				ID:    slot.ID,
				State: slot.State,
				Style: uidom.Style{
					Width:           uidom.Px(cellSize),
					Height:          uidom.Px(cellSize),
					Direction:       uidom.Column,
					Padding:         uidom.All(8),
					Gap:             4,
					BackgroundColor: cardBG,
					BorderColor:     borderMuted,
					BorderWidth:     1,
				},
			},
				uidom.Image(uidom.Props{
					ID:    slot.ID + "-icon",
					Image: icon,
				}),
				uidom.Text(slot.Label, uidom.Props{
					ID: slot.ID + "-label",
					Style: uidom.Style{
						Color: textStrong,
					},
				}),
				uidom.Text(fmt.Sprintf("x%d", slot.Quantity), uidom.Props{
					ID: slot.ID + "-qty",
					Style: uidom.Style{
						Color: textMuted,
					},
				}),
			))
		}
		rows = append(rows, uidom.Div(uidom.Props{
			ID: cfg.ID + fmt.Sprintf("-row-%d", rowIndex),
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       8,
			},
		}, rowChildren...))
	}

	return Panel(PanelConfig{
		ID:    cfg.ID,
		Title: cfg.Title,
		Width: float64(columns)*cellSize + float64(columns-1)*8 + 32,
		Children: []*uidom.Node{
			uidom.Div(uidom.Props{
				ID: cfg.ID + "-grid",
				Style: uidom.Style{
					Width:     uidom.Fill(),
					Direction: uidom.Column,
					Gap:       8,
				},
			}, rows...),
		},
	})
}

func PauseMenu(cfg PauseMenuConfig) *uidom.Node {
	children := []*uidom.Node{
		MenuList(MenuListConfig{
			ID:    cfg.ID + "-menu",
			Width: cfg.Width,
			Items: cfg.Items,
		}),
	}
	if cfg.Footer != "" {
		children = append(children, uidom.Text(cfg.Footer, uidom.Props{
			ID:    cfg.ID + "-footer",
			Style: uidom.Style{Color: textMuted},
		}))
	}

	return Panel(PanelConfig{
		ID:       cfg.ID,
		Title:    cfg.Title,
		Width:    cfg.Width,
		Children: children,
	})
}

func SettingsPanel(cfg SettingsPanelConfig) *uidom.Node {
	sectionNodes := make([]*uidom.Node, 0, len(cfg.Sections))
	for i, section := range cfg.Sections {
		optionNodes := make([]*uidom.Node, 0, len(section.Options)+2)
		optionNodes = append(optionNodes,
			uidom.Text(section.Title, uidom.Props{
				ID:    fmt.Sprintf("%s-section-%d-title", cfg.ID, i),
				Style: uidom.Style{Color: textStrong},
			}),
			uidom.TextBlock(section.Description, uidom.Props{
				ID: fmt.Sprintf("%s-section-%d-copy", cfg.ID, i),
				Style: uidom.Style{
					Width:      uidom.Fill(),
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

	scroll := uidom.ScrollView(uidom.Props{
		ID: cfg.ID + "-scroll",
		Style: uidom.Style{
			Width:     uidom.Fill(),
			Height:    uidom.Px(cfg.Height - 80),
			Direction: uidom.Column,
			Gap:       10,
		},
	}, sectionNodes...)

	return Panel(PanelConfig{
		ID:    cfg.ID,
		Title: cfg.Title,
		Width: cfg.Width,
		Children: []*uidom.Node{
			scroll,
		},
	})
}

func Tooltip(cfg TooltipConfig) *uidom.Node {
	statRows := make([]*uidom.Node, 0, len(cfg.Stats)+2)
	statRows = append(statRows,
		uidom.TextBlock(cfg.Description, uidom.Props{
			ID: cfg.ID + "-description",
			Style: uidom.Style{
				Width:      uidom.Fill(),
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

func buttonTextColor(state uidom.InteractionState) color.Color {
	if state.Disabled {
		return color.RGBA{R: 150, G: 156, B: 168, A: 255}
	}
	if state.Selected || state.Focused {
		return color.RGBA{R: 18, G: 24, B: 33, A: 255}
	}
	return textStrong
}

func menuButtonStyle(state uidom.InteractionState) uidom.Style {
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
	return uidom.Style{
		Width:           uidom.Fill(),
		Direction:       uidom.Row,
		Padding:         uidom.All(12),
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

func widthLength(width float64) uidom.Length {
	if width > 0 {
		return uidom.Px(width)
	}
	return uidom.Fill()
}
