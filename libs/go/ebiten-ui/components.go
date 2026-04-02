package ebitenui

import (
	"fmt"
	"image/color"
	"math"
	"unicode"
)

type Orientation = Direction

const (
	Vertical   Orientation = Column
	Horizontal Orientation = Row
)

type IconConfig struct {
	ID    string
	Image ImageSource
	Size  float64
}

type TextareaConfig struct {
	ID           string
	Label        string
	Value        string
	ValueBinding WritableValue[string]
	Placeholder  string
	Width        float64
	Height       float64
	State        InteractionState
	Theme        *Theme
	OnChange     func(string)
	OnSubmit     func(string)
}

type CheckboxConfig struct {
	ID             string
	Label          string
	Checked        bool
	CheckedBinding WritableValue[bool]
	State          InteractionState
	Theme          *Theme
	OnChange       func(bool)
}

type ToggleConfig struct {
	ID             string
	Label          string
	Checked        bool
	CheckedBinding WritableValue[bool]
	State          InteractionState
	Theme          *Theme
	OnChange       func(bool)
}

type SliderConfig struct {
	ID           string
	Label        string
	Min          float64
	Max          float64
	Value        float64
	ValueBinding WritableValue[float64]
	Width        float64
	State        InteractionState
	Tint         color.Color
	Theme        *Theme
	OnChange     func(float64)
}

type ScrollbarConfig struct {
	ID           string
	Orientation  Orientation
	Length       float64
	Thickness    float64
	ViewportSize float64
	ContentSize  float64
	ScrollOffset float64
}

type DropdownOption struct {
	ID    string
	Label string
	State InteractionState
}

type DropdownConfig struct {
	ID              string
	Label           string
	SelectedText    string
	SelectedBinding WritableValue[string]
	Open            bool
	OpenBinding     WritableValue[bool]
	Width           float64
	Options         []DropdownOption
	State           InteractionState
	Theme           *Theme
	OnOpenChange    func(bool)
	OnSelect        func(string)
}

type InputFieldConfig struct {
	ID           string
	Label        string
	Value        string
	ValueBinding WritableValue[string]
	Placeholder  string
	Width        float64
	State        InteractionState
	Theme        *Theme
	OnChange     func(string)
	OnSubmit     func(string)
}

type RadioOption struct {
	ID       string
	Label    string
	Selected bool
	State    InteractionState
}

type RadioGroupConfig struct {
	ID          string
	Label       string
	Orientation Orientation
	Options     []RadioOption
	OnChange    func(string)
}

type StepperConfig struct {
	ID       string
	Label    string
	Value    int
	Min      int
	Max      int
	Width    float64
	State    InteractionState
	OnChange func(int)
}

type ProgressBarConfig struct {
	ID      string
	Label   string
	Current float64
	Max     float64
	Width   float64
	Tint    color.Color
	Theme   *Theme
}

type DividerConfig struct {
	ID          string
	Orientation Orientation
	Length      float64
	Thickness   float64
	Color       color.Color
}

type GridConfig struct {
	ID       string
	Columns  int
	Gap      float64
	Children []*Node
}

type ListConfig struct {
	ID          string
	Orientation Orientation
	Gap         float64
	Items       []*Node
}

type VirtualListConfig struct {
	ID           string
	StartIndex   int
	VisibleCount int
	TotalCount   int
	Orientation  Orientation
	Gap          float64
	ItemBuilder  func(index int) *Node
}

type ModalConfig struct {
	ID      string
	Open    bool
	Title   string
	Width   float64
	Height  float64
	Theme   *Theme
	Content *Node
}

type TooltipConfig struct {
	ID          string
	Title       string
	Description string
	Width       float64
	Theme       *Theme
}

type ContextMenuItem struct {
	ID    string
	Label string
	State InteractionState
}

type ContextMenuConfig struct {
	ID    string
	Width float64
	Items []ContextMenuItem
}

type TabConfig struct {
	ID      string
	Label   string
	Content *Node
	State   InteractionState
}

type TabsConfig struct {
	ID            string
	SelectedIndex int
	Tabs          []TabConfig
	OnChange      func(int)
}

type AccordionSection struct {
	ID       string
	Label    string
	Expanded bool
	Content  *Node
	State    InteractionState
}

type AccordionConfig struct {
	ID       string
	Sections []AccordionSection
	OnToggle func(string, bool)
}

type BadgeConfig struct {
	ID    string
	Label string
	Color color.Color
}

type ChipConfig struct {
	ID          string
	Label       string
	Dismissible bool
	State       InteractionState
}

var defaultComponentTheme = DefaultTheme()

var (
	componentTextStrong = defaultComponentTheme.Palette.Text.Strong
	componentTextMuted  = defaultComponentTheme.Palette.Text.Muted
	componentPanel      = defaultComponentTheme.Palette.Surface.Elevated
	componentPanelAlt   = defaultComponentTheme.Palette.Surface.Sunken
	componentBorder     = defaultComponentTheme.Components.InputField.Border
	componentAccent     = defaultComponentTheme.Palette.Accent.Primary
	componentAccentAlt  = defaultComponentTheme.Palette.Accent.Secondary
	componentWarning    = defaultComponentTheme.Palette.Accent.Warning
	componentOverlay    = defaultComponentTheme.Palette.Overlay
)

func Icon(cfg IconConfig) *Node {
	size := cfg.Size
	if size == 0 {
		if width, _ := cfg.Image.intrinsicSize(); width > 0 {
			size = width
		} else {
			size = 16
		}
	}
	return Image(Props{
		ID:    cfg.ID,
		Image: cfg.Image,
		Style: Style{
			Width:  Px(size),
			Height: Px(size),
		},
	})
}

func Textarea(cfg TextareaConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	fieldTheme := theme.Components.Textarea
	currentValue := textBindingValue(cfg.ValueBinding, cfg.Value)
	height := cfg.Height
	if height == 0 {
		height = fieldTheme.Height
	}

	body := textLikeNodeWithTheme(theme, fieldTheme, currentValue, cfg.Placeholder, cfg.ID+"-body", cfg.ID+"-placeholder", true)
	children := []*Node{}
	if cfg.Label != "" {
		children = append(children, textLabelWithTheme(fieldTheme, cfg.ID+"-label", cfg.Label))
	}
	children = append(children, Div(Props{
		ID:        cfg.ID,
		Focusable: true,
		Handlers:  textInputHandlers(cfg.ID, func() string { return textBindingValue(cfg.ValueBinding, cfg.Value) }, bindTextValue(cfg.ValueBinding, cfg.OnChange), cfg.OnSubmit, true),
		State:     cfg.State,
		Style:     fieldContainerStyleWithTheme(fieldTheme, cfg.Width, height),
	},
		body,
		caretNodeWithTheme(fieldTheme, cfg.ID+"-caret", cfg.State, true),
	))
	return Div(Props{
		ID:    cfg.ID + "-wrapper",
		State: cfg.State,
		Style: Style{
			Width:     widthLength(cfg.Width),
			Direction: Column,
			Gap:       theme.Spacing.SM,
		},
	}, filterNil(children)...)
}

func Checkbox(cfg CheckboxConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	checkboxTheme := theme.Components.Checkbox
	state := cfg.State
	state.Selected = boolBindingValue(cfg.CheckedBinding, cfg.Checked) || state.Selected
	return InteractiveButton(Props{
		ID:    cfg.ID,
		State: state,
		Handlers: EventHandlers{
			OnClick: func(ctx EventContext) {
				value := !boolBindingOrRuntime(ctx.Runtime, cfg.ID, cfg.CheckedBinding, cfg.Checked)
				ctx.Runtime.SetBoolValue(cfg.ID, value)
				if cfg.CheckedBinding != nil {
					cfg.CheckedBinding.Set(value)
				}
				if cfg.OnChange != nil {
					cfg.OnChange(value)
				}
			},
		},
		Style: Style{
			Width:     Fill(),
			Direction: Row,
			Gap:       theme.Spacing.LG,
		},
	},
		Div(Props{
			ID:    cfg.ID + "-box",
			State: state,
			Style: checkboxBoxStyleWithTheme(checkboxTheme, state),
		},
			checkmarkNodeWithTheme(checkboxTheme, cfg.ID+"-check", state),
		),
		Text(cfg.Label, Props{
			ID:    cfg.ID + "-label",
			Style: Style{Color: checkboxTheme.Label},
		}),
	)
}

func Toggle(cfg ToggleConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	toggleTheme := theme.Components.Toggle
	state := cfg.State
	state.Selected = boolBindingValue(cfg.CheckedBinding, cfg.Checked) || state.Selected

	trackChildren := []*Node{}
	if state.Selected {
		trackChildren = append(trackChildren, Spacer(Props{
			ID:    cfg.ID + "-leading-gap",
			Style: Style{Width: Fill(), Height: Fill()},
		}))
	}
	trackChildren = append(trackChildren, Div(Props{
		ID:    cfg.ID + "-thumb",
		State: state,
		Style: Style{
			Width:           Px(toggleTheme.ThumbSize),
			Height:          Px(toggleTheme.ThumbSize),
			BackgroundColor: toggleTheme.Thumb,
		},
	}))
	if !state.Selected {
		trackChildren = append(trackChildren, Spacer(Props{
			ID:    cfg.ID + "-trailing-gap",
			Style: Style{Width: Fill(), Height: Fill()},
		}))
	}

	return InteractiveButton(Props{
		ID:    cfg.ID,
		State: state,
		Handlers: EventHandlers{
			OnClick: func(ctx EventContext) {
				value := !boolBindingOrRuntime(ctx.Runtime, cfg.ID, cfg.CheckedBinding, cfg.Checked)
				ctx.Runtime.SetBoolValue(cfg.ID, value)
				if cfg.CheckedBinding != nil {
					cfg.CheckedBinding.Set(value)
				}
				if cfg.OnChange != nil {
					cfg.OnChange(value)
				}
			},
		},
		Style: Style{
			Width:     Fill(),
			Direction: Row,
			Gap:       theme.Spacing.LG,
		},
	},
		Text(cfg.Label, Props{
			ID: cfg.ID + "-label",
			Style: Style{
				Width: Fill(),
				Color: toggleTheme.Label,
			},
		}),
		Div(Props{
			ID:    cfg.ID + "-track",
			State: state,
			Style: Style{
				Width:           Px(toggleTheme.TrackWidth),
				Height:          Px(toggleTheme.TrackHeight),
				Direction:       Row,
				Padding:         All(toggleTheme.Padding),
				BackgroundColor: toggleTrackColorWithTheme(toggleTheme, state),
			},
		}, trackChildren...),
	)
}

func Slider(cfg SliderConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	sliderTheme := theme.Components.Slider
	currentValue := numberBindingValue(cfg.ValueBinding, cfg.Value)
	width := cfg.Width
	if width == 0 {
		width = 200
	}
	ratio := clampRatio(currentValue-cfg.Min, cfg.Max-cfg.Min)
	fillWidth := width * ratio
	thumbWidth := sliderTheme.ThumbWidth
	restWidth := maxFloat(0, width-fillWidth-thumbWidth)
	tint := cfg.Tint
	if tint == nil {
		tint = sliderTheme.Fill
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Px(width),
			Direction: Column,
			Gap:       sliderTheme.Gap,
		},
	},
		statusTextWithTheme(theme, sliderTheme.Label, sliderTheme.Value, cfg.ID+"-label", cfg.Label, fmt.Sprintf("%.0f", currentValue)),
		Div(Props{
			ID:       cfg.ID + "-track",
			Handlers: sliderHandlers(cfg),
			Style: Style{
				Width:           Px(width),
				Height:          Px(sliderTheme.TrackHeight),
				Direction:       Row,
				BackgroundColor: sliderTheme.Track,
				BorderColor:     sliderTheme.Border,
				BorderWidth:     sliderTheme.BorderWidth,
			},
		},
			Div(Props{
				ID: cfg.ID + "-fill",
				Style: Style{
					Width:           Px(fillWidth),
					Height:          Fill(),
					BackgroundColor: tint,
				},
			}),
			Div(Props{
				ID:    cfg.ID + "-thumb",
				State: cfg.State,
				Style: Style{
					Width:           Px(thumbWidth),
					Height:          Fill(),
					BackgroundColor: sliderTheme.Thumb,
				},
			}),
			Spacer(Props{
				ID: cfg.ID + "-rest",
				Style: Style{
					Width:  Px(restWidth),
					Height: Fill(),
				},
			}),
		),
	)
}

func Scrollbar(cfg ScrollbarConfig) *Node {
	length := cfg.Length
	if length == 0 {
		length = 100
	}
	thickness := cfg.Thickness
	if thickness == 0 {
		thickness = 12
	}

	trackLength := length
	thumbLength := maxFloat(12, trackLength*clampRatio(cfg.ViewportSize, cfg.ContentSize))
	offsetRange := maxFloat(0, cfg.ContentSize-cfg.ViewportSize)
	offsetRatio := 0.0
	if offsetRange > 0 {
		offsetRatio = clampRatio(cfg.ScrollOffset, offsetRange)
	}
	thumbOffset := maxFloat(0, (trackLength-thumbLength)*offsetRatio)
	tailLength := maxFloat(0, trackLength-thumbLength-thumbOffset)

	if cfg.Orientation == Horizontal {
		return Div(Props{
			ID: cfg.ID,
			Style: Style{
				Width:           Px(trackLength),
				Height:          Px(thickness),
				Direction:       Row,
				BackgroundColor: componentPanelAlt,
				BorderColor:     componentBorder,
				BorderWidth:     1,
			},
		},
			Spacer(Props{ID: cfg.ID + "-lead", Style: Style{Width: Px(thumbOffset), Height: Fill()}}),
			Div(Props{
				ID: cfg.ID + "-thumb",
				Style: Style{
					Width:           Px(thumbLength),
					Height:          Fill(),
					BackgroundColor: componentAccentAlt,
				},
			}),
			Spacer(Props{ID: cfg.ID + "-trail", Style: Style{Width: Px(tailLength), Height: Fill()}}),
		)
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:           Px(thickness),
			Height:          Px(trackLength),
			Direction:       Column,
			BackgroundColor: componentPanelAlt,
			BorderColor:     componentBorder,
			BorderWidth:     1,
		},
	},
		Spacer(Props{ID: cfg.ID + "-lead", Style: Style{Width: Fill(), Height: Px(thumbOffset)}}),
		Div(Props{
			ID: cfg.ID + "-thumb",
			Style: Style{
				Width:           Fill(),
				Height:          Px(thumbLength),
				BackgroundColor: componentAccentAlt,
			},
		}),
		Spacer(Props{ID: cfg.ID + "-trail", Style: Style{Width: Fill(), Height: Px(tailLength)}}),
	)
}

func Dropdown(cfg DropdownConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	fieldTheme := theme.Components.Dropdown
	selectedID := textBindingValue(cfg.SelectedBinding, "")
	selectedText := cfg.SelectedText
	if selectedID != "" {
		selectedText = dropdownSelectedText(cfg.Options, selectedID, selectedText)
	}
	open := boolBindingValue(cfg.OpenBinding, cfg.Open)
	children := []*Node{}
	if cfg.Label != "" {
		children = append(children, textLabelWithTheme(fieldTheme, cfg.ID+"-label", cfg.Label))
	}
	children = append(children, InteractiveButton(Props{
		ID:    cfg.ID + "-trigger",
		State: cfg.State,
		Handlers: EventHandlers{
			OnClick: func(ctx EventContext) {
				nextOpen := !boolBindingOrRuntime(ctx.Runtime, cfg.ID+"-open", cfg.OpenBinding, cfg.Open)
				ctx.Runtime.SetBoolValue(cfg.ID+"-open", nextOpen)
				if cfg.OpenBinding != nil {
					cfg.OpenBinding.Set(nextOpen)
				}
				if cfg.OnOpenChange != nil {
					cfg.OnOpenChange(nextOpen)
				}
			},
		},
		Style: fieldContainerStyleWithTheme(fieldTheme, cfg.Width, fieldTheme.Height),
	},
		Text(selectedText, Props{
			ID: cfg.ID + "-value",
			Style: Style{
				Width: Fill(),
				Color: fieldTheme.Text,
			},
		}),
		Text("v", Props{
			ID:    cfg.ID + "-chevron",
			Style: Style{Color: fieldTheme.Placeholder},
		}),
	))

	if open {
		optionNodes := make([]*Node, 0, len(cfg.Options))
		for _, option := range cfg.Options {
			option := option
			optionNodes = append(optionNodes, InteractiveButton(Props{
				ID:    option.ID,
				State: option.State,
				Handlers: EventHandlers{
					OnClick: func(ctx EventContext) {
						ctx.Runtime.SetTextValue(cfg.ID+"-selected", option.ID)
						if cfg.SelectedBinding != nil {
							cfg.SelectedBinding.Set(option.ID)
						}
						if cfg.OnSelect != nil {
							cfg.OnSelect(option.ID)
						}
					},
				},
				Style: menuLikeButtonStyle(option.State),
			},
				Text(option.Label, Props{
					ID:    option.ID + "-label",
					Style: Style{Color: interactiveTextColor(option.State)},
				}),
			))
		}
		children = append(children, Div(Props{
			ID: cfg.ID + "-options",
			Style: Style{
				Width:           widthLength(cfg.Width),
				Direction:       Column,
				Gap:             6,
				Padding:         All(8),
				BackgroundColor: componentPanel,
				BorderColor:     componentBorder,
				BorderWidth:     1,
			},
		}, optionNodes...))
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     widthLength(cfg.Width),
			Direction: Column,
			Gap:       6,
		},
	}, filterNil(children)...)
}

func InputField(cfg InputFieldConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	fieldTheme := theme.Components.InputField
	currentValue := textBindingValue(cfg.ValueBinding, cfg.Value)
	children := []*Node{}
	if cfg.Label != "" {
		children = append(children, textLabelWithTheme(fieldTheme, cfg.ID+"-label", cfg.Label))
	}

	valueNode := textLikeNodeWithTheme(theme, fieldTheme, currentValue, cfg.Placeholder, cfg.ID+"-value", cfg.ID+"-placeholder", false)
	children = append(children, Div(Props{
		ID:        cfg.ID,
		Focusable: true,
		Handlers:  textInputHandlers(cfg.ID, func() string { return textBindingValue(cfg.ValueBinding, cfg.Value) }, bindTextValue(cfg.ValueBinding, cfg.OnChange), cfg.OnSubmit, false),
		State:     cfg.State,
		Style:     fieldContainerStyleWithTheme(fieldTheme, cfg.Width, fieldTheme.Height),
	},
		valueNode,
		caretNodeWithTheme(fieldTheme, cfg.ID+"-caret", cfg.State, false),
	))

	return Div(Props{
		ID: cfg.ID + "-wrapper",
		Style: Style{
			Width:     widthLength(cfg.Width),
			Direction: Column,
			Gap:       theme.Spacing.SM,
		},
	}, filterNil(children)...)
}

func RadioGroup(cfg RadioGroupConfig) *Node {
	direction := directionForOrientation(cfg.Orientation)
	optionNodes := make([]*Node, 0, len(cfg.Options))
	for _, option := range cfg.Options {
		option := option
		state := option.State
		state.Selected = option.Selected || state.Selected
		optionNodes = append(optionNodes, InteractiveButton(Props{
			ID:    option.ID,
			State: state,
			Handlers: EventHandlers{
				OnClick: func(ctx EventContext) {
					ctx.Runtime.SetTextValue(cfg.ID, option.ID)
					if cfg.OnChange != nil {
						cfg.OnChange(option.ID)
					}
				},
			},
			Style: Style{
				Width:           Auto(),
				Direction:       Row,
				Gap:             8,
				Padding:         All(8),
				BackgroundColor: choiceBackground(state),
				BorderColor:     componentBorder,
				BorderWidth:     1,
			},
		},
			Div(Props{
				ID:    option.ID + "-dot",
				State: state,
				Style: radioDotStyle(state),
			}),
			Text(option.Label, Props{
				ID:    option.ID + "-label",
				Style: Style{Color: interactiveTextColor(state)},
			}),
		))
	}

	children := []*Node{}
	if cfg.Label != "" {
		children = append(children, textLabel(cfg.ID+"-label", cfg.Label))
	}
	children = append(children, Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Fill(),
			Direction: direction,
			Gap:       8,
		},
	}, optionNodes...))

	return Div(Props{
		ID: cfg.ID + "-group",
		Style: Style{
			Width:     Fill(),
			Direction: Column,
			Gap:       6,
		},
	}, children...)
}

func Stepper(cfg StepperConfig) *Node {
	width := cfg.Width
	if width == 0 {
		width = 180
	}
	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Px(width),
			Direction: Column,
			Gap:       6,
		},
	},
		StatusText(cfg.ID+"-label", cfg.Label, fmt.Sprintf("%d", cfg.Value)),
		Div(Props{
			ID: cfg.ID + "-controls",
			Style: Style{
				Width:     Fill(),
				Direction: Row,
				Gap:       8,
			},
		},
			InteractiveButton(Props{
				ID:    cfg.ID + "-decrement",
				State: InteractionState{Disabled: cfg.Value <= cfg.Min},
				Handlers: EventHandlers{
					OnClick: func(ctx EventContext) {
						current := int(ctx.Runtime.NumberValueOrDefault(cfg.ID, float64(cfg.Value)))
						if current <= cfg.Min {
							return
						}
						current--
						ctx.Runtime.SetNumberValue(cfg.ID, float64(current))
						if cfg.OnChange != nil {
							cfg.OnChange(current)
						}
					},
				},
				Style: menuLikeButtonStyle(InteractionState{Disabled: cfg.Value <= cfg.Min}),
			}, Text("-", Props{ID: cfg.ID + "-decrement-label", Style: Style{Color: interactiveTextColor(InteractionState{Disabled: cfg.Value <= cfg.Min})}})),
			Div(Props{
				ID: cfg.ID + "-value",
				Style: Style{
					Width:           Fill(),
					Padding:         All(12),
					BackgroundColor: componentPanelAlt,
					BorderColor:     componentBorder,
					BorderWidth:     1,
				},
			}, Text(fmt.Sprintf("%d", cfg.Value), Props{ID: cfg.ID + "-value-text", Style: Style{Color: componentTextStrong}})),
			InteractiveButton(Props{
				ID:    cfg.ID + "-increment",
				State: InteractionState{Disabled: cfg.Value >= cfg.Max},
				Handlers: EventHandlers{
					OnClick: func(ctx EventContext) {
						current := int(ctx.Runtime.NumberValueOrDefault(cfg.ID, float64(cfg.Value)))
						if current >= cfg.Max {
							return
						}
						current++
						ctx.Runtime.SetNumberValue(cfg.ID, float64(current))
						if cfg.OnChange != nil {
							cfg.OnChange(current)
						}
					},
				},
				Style: menuLikeButtonStyle(InteractionState{Disabled: cfg.Value >= cfg.Max}),
			}, Text("+", Props{ID: cfg.ID + "-increment-label", Style: Style{Color: interactiveTextColor(InteractionState{Disabled: cfg.Value >= cfg.Max})}})),
		),
	)
}

func ProgressBar(cfg ProgressBarConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	progressTheme := theme.Components.ProgressBar
	width := cfg.Width
	if width == 0 {
		width = 180
	}
	tint := cfg.Tint
	if tint == nil {
		tint = progressTheme.Fill
	}
	fillWidth := width * clampRatio(cfg.Current, cfg.Max)
	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Px(width),
			Direction: Column,
			Gap:       progressTheme.Gap,
		},
	},
		statusTextWithTheme(theme, progressTheme.Label, progressTheme.Value, cfg.ID+"-label", cfg.Label, fmt.Sprintf("%.0f/%.0f", cfg.Current, cfg.Max)),
		Div(Props{
			ID: cfg.ID + "-track",
			Style: Style{
				Width:           Px(width),
				Height:          Px(progressTheme.TrackHeight),
				Direction:       Row,
				BackgroundColor: progressTheme.Track,
				BorderColor:     progressTheme.Border,
				BorderWidth:     progressTheme.BorderWidth,
			},
		},
			Div(Props{
				ID: cfg.ID + "-fill",
				Style: Style{
					Width:           Px(fillWidth),
					Height:          Fill(),
					BackgroundColor: tint,
				},
			}),
			Spacer(Props{
				ID: cfg.ID + "-rest",
				Style: Style{
					Width:  Px(maxFloat(0, width-fillWidth)),
					Height: Fill(),
				},
			}),
		),
	)
}

func Divider(cfg DividerConfig) *Node {
	length := cfg.Length
	if length == 0 {
		length = 80
	}
	thickness := cfg.Thickness
	if thickness == 0 {
		thickness = 1
	}
	fill := cfg.Color
	if fill == nil {
		fill = componentBorder
	}

	style := Style{
		BackgroundColor: fill,
	}
	if cfg.Orientation == Horizontal {
		style.Width = Px(length)
		style.Height = Px(thickness)
	} else {
		style.Width = Px(thickness)
		style.Height = Px(length)
	}
	return Div(Props{
		ID:    cfg.ID,
		Style: style,
	})
}

func Grid(cfg GridConfig) *Node {
	columns := cfg.Columns
	if columns <= 0 {
		columns = 2
	}
	rows := make([]*Node, 0, int(math.Ceil(float64(len(cfg.Children))/float64(columns))))
	for rowIndex := 0; rowIndex*columns < len(cfg.Children); rowIndex++ {
		rowChildren := make([]*Node, 0, columns)
		for column := 0; column < columns; column++ {
			index := rowIndex*columns + column
			if index >= len(cfg.Children) {
				break
			}
			rowChildren = append(rowChildren, cfg.Children[index])
		}
		rows = append(rows, Div(Props{
			ID: cfg.ID + fmt.Sprintf("-row-%d", rowIndex),
			Style: Style{
				Width:     Fill(),
				Direction: Row,
				Gap:       cfg.Gap,
			},
		}, rowChildren...))
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Fill(),
			Direction: Column,
			Gap:       cfg.Gap,
		},
	}, rows...)
}

func List(cfg ListConfig) *Node {
	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Fill(),
			Direction: directionForOrientation(cfg.Orientation),
			Gap:       cfg.Gap,
		},
	}, cfg.Items...)
}

func VirtualList(cfg VirtualListConfig) *Node {
	end := minInt(cfg.TotalCount, cfg.StartIndex+cfg.VisibleCount)
	items := make([]*Node, 0, maxInt(0, end-cfg.StartIndex))
	for i := cfg.StartIndex; i < end; i++ {
		if cfg.ItemBuilder == nil {
			continue
		}
		items = append(items, cfg.ItemBuilder(i))
	}
	return List(ListConfig{
		ID:          cfg.ID,
		Orientation: cfg.Orientation,
		Gap:         cfg.Gap,
		Items:       items,
	})
}

func Modal(cfg ModalConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	modalTheme := theme.Components.Modal
	if !cfg.Open {
		return Div(Props{
			ID: cfg.ID,
			Style: Style{
				Width:  Px(1),
				Height: Px(1),
			},
		})
	}

	width := cfg.Width
	if width == 0 {
		width = 320
	}
	height := cfg.Height
	if height == 0 {
		height = 180
	}

	content := cfg.Content
	if content == nil {
		content = Text("", Props{ID: cfg.ID + "-empty", Style: Style{Color: modalTheme.BodyText}})
	}

	return Stack(Props{
		ID: cfg.ID,
		Style: Style{
			Width:  Fill(),
			Height: Fill(),
		},
	},
		Div(Props{
			ID: cfg.ID + "-overlay",
			Style: Style{
				Width:           Fill(),
				Height:          Fill(),
				BackgroundColor: modalTheme.Backdrop,
			},
		}),
		Div(Props{
			ID: cfg.ID + "-stage",
			Style: Style{
				Width:     Fill(),
				Height:    Fill(),
				Direction: Column,
			},
		},
			Spacer(Props{ID: cfg.ID + "-top-spacer", Style: Style{Width: Fill(), Height: Fill()}}),
			Div(Props{
				ID: cfg.ID + "-center-row",
				Style: Style{
					Width:     Fill(),
					Direction: Row,
				},
			},
				Spacer(Props{ID: cfg.ID + "-left-spacer", Style: Style{Width: Fill(), Height: Fill()}}),
				Div(Props{
					ID: cfg.ID + "-content",
					Style: Style{
						Width:           Px(width),
						Height:          Px(height),
						Direction:       Column,
						Padding:         All(modalTheme.Padding),
						Gap:             modalTheme.Gap,
						BackgroundColor: modalTheme.Surface,
						BorderColor:     modalTheme.Border,
						BorderWidth:     modalTheme.BorderWidth,
					},
				},
					Text(cfg.Title, Props{ID: cfg.ID + "-title", Style: Style{Color: modalTheme.TitleText}}),
					content,
				),
				Spacer(Props{ID: cfg.ID + "-right-spacer", Style: Style{Width: Fill(), Height: Fill()}}),
			),
			Spacer(Props{ID: cfg.ID + "-bottom-spacer", Style: Style{Width: Fill(), Height: Fill()}}),
		),
	)
}

func Tooltip(cfg TooltipConfig) *Node {
	theme := ResolveTheme(cfg.Theme)
	tooltipTheme := theme.Components.Tooltip
	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:           widthLength(cfg.Width),
			Direction:       Column,
			Padding:         All(tooltipTheme.Padding),
			Gap:             tooltipTheme.Gap,
			BackgroundColor: tooltipTheme.Surface,
			BorderColor:     tooltipTheme.Border,
			BorderWidth:     tooltipTheme.BorderWidth,
		},
	},
		Text(cfg.Title, Props{
			ID:    cfg.ID + "-title",
			Style: Style{Color: tooltipTheme.TitleText},
		}),
		TextBlock(cfg.Description, Props{
			ID: cfg.ID + "-description",
			Style: Style{
				Width:      Fill(),
				Color:      tooltipTheme.BodyText,
				LineHeight: tooltipTheme.LineHeight,
			},
		}),
	)
}

func ContextMenu(cfg ContextMenuConfig) *Node {
	itemNodes := make([]*Node, 0, len(cfg.Items))
	for _, item := range cfg.Items {
		itemNodes = append(itemNodes, InteractiveButton(Props{
			ID:    item.ID,
			State: item.State,
			Style: menuLikeButtonStyle(item.State),
		},
			Text(item.Label, Props{
				ID:    item.ID + "-label",
				Style: Style{Color: interactiveTextColor(item.State)},
			}),
		))
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:           widthLength(cfg.Width),
			Direction:       Column,
			Padding:         All(8),
			Gap:             6,
			BackgroundColor: componentPanel,
			BorderColor:     componentBorder,
			BorderWidth:     1,
		},
	}, itemNodes...)
}

func Tabs(cfg TabsConfig) *Node {
	tabButtons := make([]*Node, 0, len(cfg.Tabs))
	var selectedContent *Node
	for i, tab := range cfg.Tabs {
		i := i
		tab := tab
		state := tab.State
		state.Selected = i == cfg.SelectedIndex || state.Selected
		tabButtons = append(tabButtons, InteractiveButton(Props{
			ID:    tab.ID,
			State: state,
			Handlers: EventHandlers{
				OnClick: func(ctx EventContext) {
					ctx.Runtime.SetNumberValue(cfg.ID, float64(i))
					if cfg.OnChange != nil {
						cfg.OnChange(i)
					}
				},
			},
			Style: menuLikeButtonStyle(state),
		},
			Text(tab.Label, Props{
				ID:    tab.ID + "-label",
				Style: Style{Color: interactiveTextColor(state)},
			}),
		))
		if i == cfg.SelectedIndex {
			selectedContent = tab.Content
		}
	}

	children := []*Node{
		Div(Props{
			ID: cfg.ID + "-headers",
			Style: Style{
				Width:     Fill(),
				Direction: Row,
				Gap:       8,
			},
		}, tabButtons...),
	}
	if selectedContent != nil {
		children = append(children, selectedContent)
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Fill(),
			Direction: Column,
			Gap:       8,
		},
	}, children...)
}

func Accordion(cfg AccordionConfig) *Node {
	children := make([]*Node, 0, len(cfg.Sections)*2)
	for _, section := range cfg.Sections {
		section := section
		state := section.State
		state.Selected = section.Expanded || state.Selected
		children = append(children, InteractiveButton(Props{
			ID:    section.ID,
			State: state,
			Handlers: EventHandlers{
				OnClick: func(ctx EventContext) {
					expanded := !ctx.Runtime.BoolValueOrDefault(section.ID, section.Expanded)
					ctx.Runtime.SetBoolValue(section.ID, expanded)
					if cfg.OnToggle != nil {
						cfg.OnToggle(section.ID, expanded)
					}
				},
			},
			Style: menuLikeButtonStyle(state),
		},
			Text(section.Label, Props{
				ID: section.ID + "-label",
				Style: Style{
					Width: Fill(),
					Color: interactiveTextColor(state),
				},
			}),
			Text(collapseMark(section.Expanded), Props{
				ID:    section.ID + "-chevron",
				Style: Style{Color: interactiveTextColor(state)},
			}),
		))
		if section.Expanded && section.Content != nil {
			children = append(children, Div(Props{
				ID: section.ID + "-panel",
				Style: Style{
					Width:           Fill(),
					Padding:         All(10),
					BackgroundColor: componentPanelAlt,
					BorderColor:     componentBorder,
					BorderWidth:     1,
				},
			}, section.Content))
		}
	}

	return Div(Props{
		ID: cfg.ID,
		Style: Style{
			Width:     Fill(),
			Direction: Column,
			Gap:       8,
		},
	}, children...)
}

func Badge(cfg BadgeConfig) *Node {
	fill := cfg.Color
	if fill == nil {
		fill = componentAccent
	}
	return Span(Props{
		ID: cfg.ID,
		Style: Style{
			Padding:         All(8),
			BackgroundColor: fill,
		},
	},
		Text(cfg.Label, Props{
			ID: cfg.ID + "-label",
			Style: Style{
				Color: color.RGBA{R: 10, G: 16, B: 24, A: 255},
			},
		}),
	)
}

func Chip(cfg ChipConfig) *Node {
	children := []*Node{
		Text(cfg.Label, Props{
			ID: cfg.ID + "-label",
			Style: Style{
				Color: interactiveTextColor(cfg.State),
			},
		}),
	}
	if cfg.Dismissible {
		children = append(children, InteractiveButton(Props{
			ID:    cfg.ID + "-dismiss",
			State: cfg.State,
			Style: Style{
				Padding:         All(4),
				BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 24},
			},
		},
			Text("x", Props{
				ID:    cfg.ID + "-dismiss-text",
				Style: Style{Color: interactiveTextColor(cfg.State)},
			}),
		))
	}
	return Span(Props{
		ID:    cfg.ID,
		State: cfg.State,
		Style: Style{
			Direction:       Row,
			Padding:         All(8),
			Gap:             6,
			BackgroundColor: choiceBackground(cfg.State),
			BorderColor:     componentBorder,
			BorderWidth:     1,
		},
	}, children...)
}

func ComponentID(prefix string, index int) string {
	return fmt.Sprintf("%s-%d", prefix, index)
}

func StatusText(id string, label string, value string) *Node {
	theme := DefaultTheme()
	return statusTextWithTheme(theme, theme.Palette.Text.Muted, theme.Palette.Text.Strong, id, label, value)
}

func statusTextWithTheme(_ Theme, labelColor color.Color, valueColor color.Color, id string, label string, value string) *Node {
	return Div(Props{
		ID: id,
		Style: Style{
			Width:     Fill(),
			Direction: Row,
			Gap:       8,
		},
	},
		Text(label, Props{
			ID: id + "-left",
			Style: Style{
				Width: Fill(),
				Color: labelColor,
			},
		}),
		Text(value, Props{
			ID:    id + "-right",
			Style: Style{Color: valueColor},
		}),
	)
}

func textLabel(id string, value string) *Node {
	return textLabelWithTheme(DefaultTheme().Components.InputField, id, value)
}

func textLabelWithTheme(fieldTheme FieldTheme, id string, value string) *Node {
	return Text(value, Props{
		ID: id,
		Style: Style{
			Color: fieldTheme.Label,
		},
	})
}

func textLikeNode(value, placeholder, valueID, placeholderID string, multiline bool) *Node {
	return textLikeNodeWithTheme(DefaultTheme(), DefaultTheme().Components.InputField, value, placeholder, valueID, placeholderID, multiline)
}

func textLikeNodeWithTheme(theme Theme, fieldTheme FieldTheme, value, placeholder, valueID, placeholderID string, multiline bool) *Node {
	if value != "" {
		if multiline {
			return TextBlock(value, Props{
				ID: valueID,
				Style: Style{
					Width:      Fill(),
					Color:      fieldTheme.Text,
					LineHeight: floatOrDefault(fieldTheme.LineHeight, theme.Typography.BodyLineHeight),
				},
			})
		}
		return Text(value, Props{
			ID: valueID,
			Style: Style{
				Width: Fill(),
				Color: fieldTheme.Text,
			},
		})
	}

	if multiline {
		return TextBlock(placeholder, Props{
			ID: placeholderID,
			Style: Style{
				Width:      Fill(),
				Color:      fieldTheme.Placeholder,
				LineHeight: floatOrDefault(fieldTheme.LineHeight, theme.Typography.BodyLineHeight),
			},
		})
	}
	return Text(placeholder, Props{
		ID: placeholderID,
		Style: Style{
			Width: Fill(),
			Color: fieldTheme.Placeholder,
		},
	})
}

func caretNode(id string, state InteractionState, multiline bool) *Node {
	return caretNodeWithTheme(DefaultTheme().Components.InputField, id, state, multiline)
}

func caretNodeWithTheme(fieldTheme FieldTheme, id string, state InteractionState, multiline bool) *Node {
	if !state.Focused {
		return nil
	}
	height := 18.0
	if multiline {
		height = 28
	}
	return Divider(DividerConfig{
		ID:          id,
		Orientation: Vertical,
		Length:      height,
		Thickness:   2,
		Color:       fieldTheme.Caret,
	})
}

func fieldContainerStyle(width, height float64) Style {
	return fieldContainerStyleWithTheme(DefaultTheme().Components.InputField, width, height)
}

func fieldContainerStyleWithTheme(fieldTheme FieldTheme, width, height float64) Style {
	if height == 0 {
		height = fieldTheme.Height
	}
	return Style{
		Width:           widthLength(width),
		Height:          Px(height),
		Direction:       Row,
		Padding:         All(fieldTheme.Padding),
		Gap:             fieldTheme.Gap,
		BackgroundColor: fieldTheme.Background,
		BorderColor:     fieldTheme.Border,
		BorderWidth:     fieldTheme.BorderWidth,
	}
}

func checkmarkNode(id string, state InteractionState) *Node {
	return checkmarkNodeWithTheme(DefaultTheme().Components.Checkbox, id, state)
}

func checkmarkNodeWithTheme(checkboxTheme CheckboxTheme, id string, state InteractionState) *Node {
	if !state.Selected {
		return nil
	}
	return Text("v", Props{
		ID:    id,
		Style: Style{Color: checkboxTheme.Check},
	})
}

func checkboxBoxStyle(state InteractionState) Style {
	return checkboxBoxStyleWithTheme(DefaultTheme().Components.Checkbox, state)
}

func checkboxBoxStyleWithTheme(checkboxTheme CheckboxTheme, state InteractionState) Style {
	fill := checkboxTheme.Box
	if state.Selected {
		fill = checkboxTheme.BoxSelected
	}
	if state.Disabled {
		fill = checkboxTheme.BoxDisabled
	}
	return Style{
		Width:           Px(checkboxTheme.Size),
		Height:          Px(checkboxTheme.Size),
		Padding:         All(checkboxTheme.Padding),
		BackgroundColor: fill,
		BorderColor:     checkboxTheme.Border,
		BorderWidth:     checkboxTheme.BorderWidth,
	}
}

func radioDotStyle(state InteractionState) Style {
	var fill color.Color = color.RGBA{R: 255, G: 255, B: 255, A: 24}
	if state.Selected {
		fill = componentWarning
	}
	return Style{
		Width:           Px(14),
		Height:          Px(14),
		BackgroundColor: fill,
		BorderColor:     componentBorder,
		BorderWidth:     1,
	}
}

func toggleTrackColor(state InteractionState) color.Color {
	return toggleTrackColorWithTheme(DefaultTheme().Components.Toggle, state)
}

func toggleTrackColorWithTheme(toggleTheme ToggleTheme, state InteractionState) color.Color {
	if state.Disabled {
		return toggleTheme.TrackDisabled
	}
	if state.Selected {
		return toggleTheme.TrackOn
	}
	return toggleTheme.TrackOff
}

func choiceBackground(state InteractionState) color.Color {
	if state.Selected {
		return componentAccent
	}
	if state.Focused {
		return componentAccentAlt
	}
	if state.Disabled {
		return color.RGBA{R: 66, G: 72, B: 84, A: 255}
	}
	return componentPanelAlt
}

func menuLikeButtonStyle(state InteractionState) Style {
	return Style{
		Width:           Fill(),
		Direction:       Row,
		Padding:         All(10),
		Gap:             8,
		BackgroundColor: choiceBackground(state),
		BorderColor:     componentBorder,
		BorderWidth:     1,
	}
}

func interactiveTextColor(state InteractionState) color.Color {
	if state.Selected || state.Focused {
		return color.RGBA{R: 12, G: 18, B: 26, A: 255}
	}
	if state.Disabled {
		return color.RGBA{R: 150, G: 156, B: 166, A: 255}
	}
	return componentTextStrong
}

func collapseMark(expanded bool) string {
	if expanded {
		return "-"
	}
	return "+"
}

func directionForOrientation(orientation Orientation) Direction {
	if orientation == Horizontal {
		return Row
	}
	return Column
}

func textInputHandlers(id string, initial func() string, onChange func(string), onSubmit func(string), multiline bool) EventHandlers {
	return EventHandlers{
		OnFocus: func(ctx EventContext) {
			current := ""
			if initial != nil {
				current = initial()
			}
			value := ctx.Runtime.TextValueOrDefault(id, current)
			ctx.Runtime.SetTextValue(id, value)
			ctx.Runtime.SetTextCursor(id, len([]rune(value)))
			ctx.Runtime.SetTextSelection(id, TextSelection{})
		},
		OnBlur: func(ctx EventContext) {
			ctx.Runtime.SetTextSelection(id, TextSelection{})
		},
		OnShortcut: func(ctx EventContext) {
			switch ctx.Shortcut {
			case "ctrl+a", "cmd+a", "meta+a":
				value := ctx.Runtime.TextValueOrDefault(id, textInitialValue(initial))
				length := len([]rune(value))
				ctx.Runtime.SetTextCursor(id, length)
				ctx.Runtime.SetTextSelection(id, TextSelection{Start: 0, End: length})
			}
		},
		OnSelectAll: func(ctx EventContext) {
			value := ctx.Runtime.TextValueOrDefault(id, textInitialValue(initial))
			length := len([]rune(value))
			ctx.Runtime.SetTextCursor(id, length)
			ctx.Runtime.SetTextSelection(id, TextSelection{Start: 0, End: length})
		},
		OnCursorMove: func(ctx EventContext) {
			value, cursor, selection := ctx.Runtime.textValueAndCursor(id, textInitialValue(initial))
			cursor, selection = moveTextCursor(value, cursor, selection, ctx.Input, multiline)
			ctx.Runtime.SetTextCursor(id, cursor)
			ctx.Runtime.SetTextSelection(id, selection)
		},
		OnTextInput: func(ctx EventContext) {
			value, cursor, selection := ctx.Runtime.textValueAndCursor(id, textInitialValue(initial))
			if ctx.Input.SelectAll || isSelectAllShortcut(ctx.Shortcut) {
				selection = TextSelection{Start: 0, End: len([]rune(value))}
				cursor = len([]rune(value))
			}
			value, cursor = insertTextAtCursor(value, cursor, selection, ctx.Text)
			ctx.Runtime.SetTextValue(id, value)
			ctx.Runtime.SetTextCursor(id, cursor)
			ctx.Runtime.SetTextSelection(id, TextSelection{})
			if onChange != nil {
				onChange(value)
			}
		},
		OnBackspace: func(ctx EventContext) {
			value, cursor, selection := ctx.Runtime.textValueAndCursor(id, textInitialValue(initial))
			value, cursor = deleteTextBackward(value, cursor, selection, ctx.Input)
			ctx.Runtime.SetTextValue(id, value)
			ctx.Runtime.SetTextCursor(id, cursor)
			ctx.Runtime.SetTextSelection(id, TextSelection{})
			if onChange != nil {
				onChange(value)
			}
		},
		OnDelete: func(ctx EventContext) {
			value, cursor, selection := ctx.Runtime.textValueAndCursor(id, textInitialValue(initial))
			value, cursor = deleteTextForward(value, cursor, selection, ctx.Input)
			ctx.Runtime.SetTextValue(id, value)
			ctx.Runtime.SetTextCursor(id, cursor)
			ctx.Runtime.SetTextSelection(id, TextSelection{})
			if onChange != nil {
				onChange(value)
			}
		},
		OnSubmit: func(ctx EventContext) {
			value, cursor, selection := ctx.Runtime.textValueAndCursor(id, textInitialValue(initial))
			if multiline && !(ctx.Input.Control || ctx.Input.Meta) {
				value, cursor = insertTextAtCursor(value, cursor, selection, "\n")
				ctx.Runtime.SetTextValue(id, value)
				ctx.Runtime.SetTextCursor(id, cursor)
				ctx.Runtime.SetTextSelection(id, TextSelection{})
				if onChange != nil {
					onChange(value)
				}
				return
			}
			if onSubmit != nil {
				onSubmit(value)
			}
		},
	}
}

func textInitialValue(initial func() string) string {
	if initial == nil {
		return ""
	}
	return initial()
}

func insertTextAtCursor(value string, cursor int, selection TextSelection, text string) (string, int) {
	if text == "" {
		return value, cursor
	}
	runes := []rune(value)
	start, end := normalizeSelectionRange(selection.Start, selection.End, cursor, len(runes))
	insert := []rune(text)
	out := make([]rune, 0, len(runes)-maxInt(0, end-start)+len(insert))
	out = append(out, runes[:start]...)
	out = append(out, insert...)
	out = append(out, runes[end:]...)
	cursor = start + len(insert)
	return string(out), cursor
}

func isSelectAllShortcut(shortcut string) bool {
	switch shortcut {
	case "ctrl+a", "cmd+a", "meta+a", "select-all":
		return true
	default:
		return false
	}
}

func deleteTextBackward(value string, cursor int, selection TextSelection, input InputSnapshot) (string, int) {
	runes := []rune(value)
	start, end := normalizeSelectionRange(selection.Start, selection.End, cursor, len(runes))
	if start != end {
		out := append([]rune{}, runes[:start]...)
		out = append(out, runes[end:]...)
		return string(out), start
	}
	if cursor <= 0 || len(runes) == 0 {
		return value, cursor
	}
	cursor = clampInt(cursor, 0, len(runes))
	if input.Control || input.Meta {
		start := previousDeletionBoundary(runes, cursor)
		out := append([]rune{}, runes[:start]...)
		out = append(out, runes[cursor:]...)
		return string(out), start
	}
	out := append([]rune{}, runes[:cursor-1]...)
	out = append(out, runes[cursor:]...)
	return string(out), cursor - 1
}

func deleteTextForward(value string, cursor int, selection TextSelection, input InputSnapshot) (string, int) {
	runes := []rune(value)
	start, end := normalizeSelectionRange(selection.Start, selection.End, cursor, len(runes))
	if start != end {
		out := append([]rune{}, runes[:start]...)
		out = append(out, runes[end:]...)
		return string(out), start
	}
	cursor = clampInt(cursor, 0, len(runes))
	if cursor >= len(runes) {
		return value, cursor
	}
	if input.Control || input.Meta {
		end := nextDeletionBoundary(runes, cursor)
		out := append([]rune{}, runes[:cursor]...)
		out = append(out, runes[end:]...)
		return string(out), cursor
	}
	out := append([]rune{}, runes[:cursor]...)
	out = append(out, runes[cursor+1:]...)
	return string(out), cursor
}

func moveTextCursor(value string, cursor int, selection TextSelection, input InputSnapshot, multiline bool) (int, TextSelection) {
	runes := []rune(value)
	cursor = clampInt(cursor, 0, len(runes))
	next := cursor
	switch {
	case input.Home:
		if multiline && !(input.Control || input.Meta) {
			next = lineStartBoundary(runes, cursor)
		} else {
			next = 0
		}
	case input.End:
		if multiline && !(input.Control || input.Meta) {
			next = lineEndBoundary(runes, cursor)
		} else {
			next = len(runes)
		}
	case input.Control && input.ArrowLeft:
		next = previousWordBoundary(runes, cursor)
	case input.Control && input.ArrowRight:
		next = nextWordBoundary(runes, cursor)
	case multiline && input.ArrowUp:
		next = lineVerticalBoundary(runes, cursor, -1)
	case multiline && input.ArrowDown:
		next = lineVerticalBoundary(runes, cursor, 1)
	case input.ArrowLeft:
		next = cursor - 1
	case input.ArrowRight:
		next = cursor + 1
	default:
		return cursor, selection
	}
	next = clampInt(next, 0, len(runes))
	if input.Shift {
		if selection.Start == selection.End {
			return next, TextSelection{Start: cursor, End: next}
		}
		if cursor == selection.Start {
			return next, TextSelection{Start: next, End: selection.End}
		}
		return next, TextSelection{Start: selection.Start, End: next}
	}
	return next, TextSelection{}
}

func previousWordBoundary(runes []rune, cursor int) int {
	if cursor <= 0 {
		return 0
	}
	i := clampInt(cursor, 0, len(runes)) - 1
	for i > 0 && runes[i] == ' ' {
		i--
	}
	for i > 0 {
		if runes[i-1] == ' ' {
			break
		}
		if unicode.IsLower(runes[i-1]) != unicode.IsLower(runes[i]) || unicode.IsUpper(runes[i-1]) != unicode.IsUpper(runes[i]) {
			break
		}
		i--
	}
	return i
}

func nextWordBoundary(runes []rune, cursor int) int {
	if cursor >= len(runes) {
		return len(runes)
	}
	i := clampInt(cursor, 0, len(runes))
	for i < len(runes) && runes[i] != ' ' {
		if i > 0 && (unicode.IsLower(runes[i-1]) != unicode.IsLower(runes[i]) || unicode.IsUpper(runes[i-1]) != unicode.IsUpper(runes[i])) {
			break
		}
		i++
	}
	for i < len(runes) && runes[i] == ' ' {
		i++
	}
	return i
}

func previousDeletionBoundary(runes []rune, cursor int) int {
	if cursor <= 0 {
		return 0
	}
	i := clampInt(cursor, 0, len(runes))
	for i > 0 && runes[i-1] == ' ' {
		i--
	}
	for i > 0 && runes[i-1] != ' ' {
		i--
	}
	for i > 0 && runes[i-1] == ' ' {
		i--
	}
	return i
}

func nextDeletionBoundary(runes []rune, cursor int) int {
	if cursor >= len(runes) {
		return len(runes)
	}
	i := clampInt(cursor, 0, len(runes))
	for i < len(runes) && runes[i] == ' ' {
		i++
	}
	for i < len(runes) && runes[i] != ' ' {
		i++
	}
	for i < len(runes) && runes[i] == ' ' {
		i++
	}
	return i
}

func lineStartBoundary(runes []rune, cursor int) int {
	cursor = clampInt(cursor, 0, len(runes))
	for i := cursor - 1; i >= 0; i-- {
		if runes[i] == '\n' {
			return i + 1
		}
	}
	return 0
}

func lineEndBoundary(runes []rune, cursor int) int {
	cursor = clampInt(cursor, 0, len(runes))
	for i := cursor; i < len(runes); i++ {
		if runes[i] == '\n' {
			return i
		}
	}
	return len(runes)
}

func lineVerticalBoundary(runes []rune, cursor int, direction int) int {
	cursor = clampInt(cursor, 0, len(runes))
	start := lineStartBoundary(runes, cursor)
	column := cursor - start
	if direction < 0 {
		if start == 0 {
			return 0
		}
		prevEnd := start - 1
		prevStart := lineStartBoundary(runes, prevEnd)
		return clampInt(prevStart+column, prevStart, prevEnd)
	}
	end := lineEndBoundary(runes, cursor)
	if end >= len(runes) {
		return len(runes)
	}
	nextStart := end + 1
	nextEnd := lineEndBoundary(runes, nextStart)
	return clampInt(nextStart+column, nextStart, nextEnd)
}

func normalizeSelectionRange(start, end, cursor, length int) (int, int) {
	if start == end {
		cursor = clampInt(cursor, 0, length)
		return cursor, cursor
	}
	start = clampInt(start, 0, length)
	end = clampInt(end, 0, length)
	if start > end {
		start, end = end, start
	}
	return start, end
}

func sliderHandlers(cfg SliderConfig) EventHandlers {
	apply := func(ctx EventContext) {
		width := ctx.Layout.Frame.Width
		if width <= 0 {
			return
		}
		ratio := clampRatio(ctx.LocalX, width)
		value := cfg.Min + (cfg.Max-cfg.Min)*ratio
		ctx.Runtime.SetNumberValue(cfg.ID, value)
		if cfg.ValueBinding != nil {
			cfg.ValueBinding.Set(value)
		}
		if cfg.OnChange != nil {
			cfg.OnChange(value)
		}
	}
	return EventHandlers{
		OnPointerDown: apply,
		OnPointerMove: func(ctx EventContext) {
			if ctx.Input.PointerDown {
				apply(ctx)
			}
		},
		OnClick: apply,
	}
}

func textBindingValue(binding Value[string], fallback string) string {
	if binding != nil {
		return binding.Get()
	}
	return fallback
}

func bindTextValue(binding WritableValue[string], observer func(string)) func(string) {
	return func(value string) {
		if binding != nil {
			binding.Set(value)
		}
		if observer != nil {
			observer(value)
		}
	}
}

func boolBindingValue(binding Value[bool], fallback bool) bool {
	if binding != nil {
		return binding.Get()
	}
	return fallback
}

func boolBindingOrRuntime(runtime *Runtime, id string, binding Value[bool], fallback bool) bool {
	if binding != nil {
		return binding.Get()
	}
	return runtime.BoolValueOrDefault(id, fallback)
}

func numberBindingValue(binding Value[float64], fallback float64) float64 {
	if binding != nil {
		return binding.Get()
	}
	return fallback
}

func dropdownSelectedText(options []DropdownOption, selectedID string, fallback string) string {
	for _, option := range options {
		if option.ID == selectedID {
			return option.Label
		}
	}
	if selectedID != "" {
		return selectedID
	}
	return fallback
}

func trimLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}
	return string(runes[:len(runes)-1])
}

func widthLength(width float64) Length {
	if width > 0 {
		return Px(width)
	}
	return Fill()
}

func clampRatio(current, max float64) float64 {
	if max <= 0 {
		return 0
	}
	ratio := current / max
	return math.Max(0, math.Min(1, ratio))
}

func filterNil(nodes []*Node) []*Node {
	result := make([]*Node, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		result = append(result, node)
	}
	return result
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
