package ebitenui

import (
	"image/color"
	"strings"
)

type Theme struct {
	Name       string
	Palette    ThemePalette
	Spacing    ThemeSpacing
	Radius     ThemeRadius
	Stroke     ThemeStroke
	Typography ThemeTypography
	Components ThemeComponents
}

type ThemePalette struct {
	Surface ThemeSurfacePalette
	Text    ThemeTextPalette
	Accent  ThemeAccentPalette
	Status  ThemeStatusPalette
	Overlay color.Color
}

type ThemeSurfacePalette struct {
	Base     color.Color
	Elevated color.Color
	Sunken   color.Color
	Panel    color.Color
	Card     color.Color
	Field    color.Color
}

type ThemeTextPalette struct {
	Strong  color.Color
	Muted   color.Color
	Inverse color.Color
	Accent  color.Color
}

type ThemeAccentPalette struct {
	Primary   color.Color
	Secondary color.Color
	Warning   color.Color
}

type ThemeStatusPalette struct {
	Success color.Color
	Warning color.Color
	Error   color.Color
}

type ThemeSpacing struct {
	XS float64
	SM float64
	MD float64
	LG float64
	XL float64
}

type ThemeRadius struct {
	SM float64
	MD float64
	LG float64
}

type ThemeStroke struct {
	Thin   float64
	Medium float64
	Focus  float64
}

type ThemeTypography struct {
	LabelLineHeight float64
	BodyLineHeight  float64
	CodeLineHeight  float64
}

type ThemeComponents struct {
	InputField    FieldTheme
	Textarea      FieldTheme
	Dropdown      FieldTheme
	Checkbox      CheckboxTheme
	Toggle        ToggleTheme
	Slider        SliderTheme
	ProgressBar   ProgressBarTheme
	Panel         PanelTheme
	Card          CardTheme
	MenuButton    MenuButtonTheme
	Dialog        DialogTheme
	HUDBar        HUDBarTheme
	InventoryGrid InventoryGridTheme
	Modal         ModalTheme
	Tooltip       TooltipTheme
}

type FieldTheme struct {
	Background  color.Color
	Border      color.Color
	Text        color.Color
	Placeholder color.Color
	Label       color.Color
	Caret       color.Color
	Padding     float64
	Gap         float64
	Height      float64
	BorderWidth float64
	LineHeight  float64
}

type CheckboxTheme struct {
	Label       color.Color
	Border      color.Color
	Box         color.Color
	BoxSelected color.Color
	BoxDisabled color.Color
	Check       color.Color
	Size        float64
	Padding     float64
	BorderWidth float64
}

type ToggleTheme struct {
	Label         color.Color
	TrackOff      color.Color
	TrackOn       color.Color
	TrackDisabled color.Color
	Thumb         color.Color
	TrackWidth    float64
	TrackHeight   float64
	Padding       float64
	ThumbSize     float64
}

type SliderTheme struct {
	Label       color.Color
	Value       color.Color
	Track       color.Color
	Border      color.Color
	Fill        color.Color
	Thumb       color.Color
	TrackHeight float64
	ThumbWidth  float64
	BorderWidth float64
	Gap         float64
}

type ProgressBarTheme struct {
	Label       color.Color
	Value       color.Color
	Track       color.Color
	Border      color.Color
	Fill        color.Color
	TrackHeight float64
	BorderWidth float64
	Gap         float64
}

type PanelTheme struct {
	Background  color.Color
	Border      color.Color
	TitleText   color.Color
	Padding     float64
	Gap         float64
	BorderWidth float64
}

type CardTheme struct {
	Background  color.Color
	Border      color.Color
	TitleText   color.Color
	Padding     float64
	Gap         float64
	BorderWidth float64
}

type StateColors struct {
	Background color.Color
	Border     color.Color
	Text       color.Color
}

type MenuButtonTheme struct {
	Default     StateColors
	Focused     StateColors
	Selected    StateColors
	Disabled    StateColors
	Padding     float64
	Gap         float64
	BorderWidth float64
}

type DialogTheme struct {
	BodyText color.Color
}

type HUDBarTheme struct {
	Track       color.Color
	Border      color.Color
	Fill        color.Color
	Text        color.Color
	TrackHeight float64
	BorderWidth float64
	Gap         float64
}

type InventoryGridTheme struct {
	SlotBackground  color.Color
	SlotBorder      color.Color
	SlotText        color.Color
	SlotMuted       color.Color
	IconFill        color.Color
	SlotPadding     float64
	SlotGap         float64
	SlotBorderWidth float64
}

type ModalTheme struct {
	Backdrop    color.Color
	Surface     color.Color
	Border      color.Color
	TitleText   color.Color
	BodyText    color.Color
	Padding     float64
	Gap         float64
	BorderWidth float64
}

type TooltipTheme struct {
	Surface     color.Color
	Border      color.Color
	TitleText   color.Color
	BodyText    color.Color
	Padding     float64
	Gap         float64
	BorderWidth float64
	LineHeight  float64
}

func DefaultTheme() Theme {
	return Theme{
		Name: "default",
		Palette: ThemePalette{
			Surface: ThemeSurfacePalette{
				Base:     color.RGBA{R: 13, G: 18, B: 27, A: 255},
				Elevated: color.RGBA{R: 26, G: 32, B: 44, A: 255},
				Sunken:   color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Panel:    color.RGBA{R: 25, G: 32, B: 42, A: 255},
				Card:     color.RGBA{R: 31, G: 40, B: 54, A: 255},
				Field:    color.RGBA{R: 19, G: 25, B: 35, A: 255},
			},
			Text: ThemeTextPalette{
				Strong:  color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Muted:   color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Inverse: color.RGBA{R: 12, G: 18, B: 26, A: 255},
				Accent:  color.RGBA{R: 250, G: 245, B: 236, A: 255},
			},
			Accent: ThemeAccentPalette{
				Primary:   color.RGBA{R: 92, G: 162, B: 255, A: 255},
				Secondary: color.RGBA{R: 82, G: 205, B: 150, A: 255},
				Warning:   color.RGBA{R: 255, G: 193, B: 82, A: 255},
			},
			Status: ThemeStatusPalette{
				Success: color.RGBA{R: 82, G: 205, B: 150, A: 255},
				Warning: color.RGBA{R: 255, G: 193, B: 82, A: 255},
				Error:   color.RGBA{R: 230, G: 70, B: 90, A: 255},
			},
			Overlay: color.RGBA{R: 6, G: 10, B: 18, A: 190},
		},
		Spacing: ThemeSpacing{
			XS: 4,
			SM: 6,
			MD: 8,
			LG: 12,
			XL: 16,
		},
		Radius: ThemeRadius{
			SM: 6,
			MD: 10,
			LG: 16,
		},
		Stroke: ThemeStroke{
			Thin:   1,
			Medium: 2,
			Focus:  2,
		},
		Typography: ThemeTypography{
			LabelLineHeight: 16,
			BodyLineHeight:  16,
			CodeLineHeight:  16,
		},
		Components: ThemeComponents{
			InputField: FieldTheme{
				Background:  color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Text:        color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Placeholder: color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Label:       color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Caret:       color.RGBA{R: 255, G: 193, B: 82, A: 255},
				Padding:     12,
				Gap:         8,
				Height:      40,
				BorderWidth: 1,
				LineHeight:  16,
			},
			Textarea: FieldTheme{
				Background:  color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Text:        color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Placeholder: color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Label:       color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Caret:       color.RGBA{R: 255, G: 193, B: 82, A: 255},
				Padding:     12,
				Gap:         8,
				Height:      96,
				BorderWidth: 1,
				LineHeight:  16,
			},
			Dropdown: FieldTheme{
				Background:  color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Text:        color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Placeholder: color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Label:       color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Caret:       color.RGBA{R: 255, G: 193, B: 82, A: 255},
				Padding:     12,
				Gap:         8,
				Height:      40,
				BorderWidth: 1,
				LineHeight:  16,
			},
			Checkbox: CheckboxTheme{
				Label:       color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Box:         color.RGBA{R: 230, G: 234, B: 242, A: 255},
				BoxSelected: color.RGBA{R: 92, G: 162, B: 255, A: 255},
				BoxDisabled: color.RGBA{R: 84, G: 90, B: 100, A: 255},
				Check:       color.RGBA{R: 12, G: 18, B: 26, A: 255},
				Size:        20,
				Padding:     2,
				BorderWidth: 1,
			},
			Toggle: ToggleTheme{
				Label:         color.RGBA{R: 239, G: 244, B: 250, A: 255},
				TrackOff:      color.RGBA{R: 19, G: 25, B: 35, A: 255},
				TrackOn:       color.RGBA{R: 82, G: 205, B: 150, A: 255},
				TrackDisabled: color.RGBA{R: 84, G: 90, B: 100, A: 255},
				Thumb:         color.RGBA{R: 244, G: 246, B: 250, A: 255},
				TrackWidth:    40,
				TrackHeight:   20,
				Padding:       2,
				ThumbSize:     16,
			},
			Slider: SliderTheme{
				Label:       color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Value:       color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Track:       color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Fill:        color.RGBA{R: 92, G: 162, B: 255, A: 255},
				Thumb:       color.RGBA{R: 255, G: 193, B: 82, A: 255},
				TrackHeight: 18,
				ThumbWidth:  12,
				BorderWidth: 1,
				Gap:         6,
			},
			ProgressBar: ProgressBarTheme{
				Label:       color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Value:       color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Track:       color.RGBA{R: 19, G: 25, B: 35, A: 255},
				Border:      color.RGBA{R: 88, G: 106, B: 132, A: 255},
				Fill:        color.RGBA{R: 92, G: 162, B: 255, A: 255},
				TrackHeight: 14,
				BorderWidth: 1,
				Gap:         6,
			},
			Panel: PanelTheme{
				Background:  color.RGBA{R: 25, G: 32, B: 42, A: 255},
				Border:      color.RGBA{R: 85, G: 103, B: 128, A: 255},
				TitleText:   color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Padding:     16,
				Gap:         12,
				BorderWidth: 1,
			},
			Card: CardTheme{
				Background:  color.RGBA{R: 31, G: 40, B: 54, A: 255},
				Border:      color.RGBA{R: 85, G: 103, B: 128, A: 255},
				TitleText:   color.RGBA{R: 239, G: 244, B: 250, A: 255},
				Padding:     12,
				Gap:         10,
				BorderWidth: 1,
			},
			MenuButton: MenuButtonTheme{
				Default: StateColors{
					Background: color.RGBA{R: 31, G: 40, B: 54, A: 255},
					Border:     color.RGBA{R: 85, G: 103, B: 128, A: 255},
					Text:       color.RGBA{R: 239, G: 244, B: 250, A: 255},
				},
				Focused: StateColors{
					Background: color.RGBA{R: 82, G: 205, B: 150, A: 255},
					Border:     color.RGBA{R: 255, G: 194, B: 82, A: 255},
					Text:       color.RGBA{R: 18, G: 24, B: 33, A: 255},
				},
				Selected: StateColors{
					Background: color.RGBA{R: 91, G: 162, B: 255, A: 255},
					Border:     color.RGBA{R: 255, G: 194, B: 82, A: 255},
					Text:       color.RGBA{R: 18, G: 24, B: 33, A: 255},
				},
				Disabled: StateColors{
					Background: color.RGBA{R: 50, G: 56, B: 68, A: 255},
					Border:     color.RGBA{R: 74, G: 80, B: 90, A: 255},
					Text:       color.RGBA{R: 150, G: 156, B: 168, A: 255},
				},
				Padding:     12,
				Gap:         8,
				BorderWidth: 1,
			},
			Dialog: DialogTheme{
				BodyText: color.RGBA{R: 178, G: 188, B: 204, A: 255},
			},
			HUDBar: HUDBarTheme{
				Track:       color.RGBA{R: 18, G: 24, B: 33, A: 255},
				Border:      color.RGBA{R: 85, G: 103, B: 128, A: 255},
				Fill:        color.RGBA{R: 91, G: 162, B: 255, A: 255},
				Text:        color.RGBA{R: 239, G: 244, B: 250, A: 255},
				TrackHeight: 18,
				BorderWidth: 1,
				Gap:         6,
			},
			InventoryGrid: InventoryGridTheme{
				SlotBackground:  color.RGBA{R: 31, G: 40, B: 54, A: 255},
				SlotBorder:      color.RGBA{R: 85, G: 103, B: 128, A: 255},
				SlotText:        color.RGBA{R: 239, G: 244, B: 250, A: 255},
				SlotMuted:       color.RGBA{R: 178, G: 188, B: 204, A: 255},
				IconFill:        color.RGBA{R: 91, G: 162, B: 255, A: 255},
				SlotPadding:     8,
				SlotGap:         4,
				SlotBorderWidth: 1,
			},
			Modal: ModalTheme{
				Backdrop:    color.RGBA{R: 6, G: 10, B: 18, A: 190},
				Surface:     color.RGBA{R: 26, G: 32, B: 44, A: 255},
				Border:      color.RGBA{R: 255, G: 193, B: 82, A: 255},
				TitleText:   color.RGBA{R: 239, G: 244, B: 250, A: 255},
				BodyText:    color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Padding:     16,
				Gap:         12,
				BorderWidth: 2,
			},
			Tooltip: TooltipTheme{
				Surface:     color.RGBA{R: 26, G: 32, B: 44, A: 255},
				Border:      color.RGBA{R: 255, G: 193, B: 82, A: 255},
				TitleText:   color.RGBA{R: 255, G: 193, B: 82, A: 255},
				BodyText:    color.RGBA{R: 176, G: 188, B: 204, A: 255},
				Padding:     12,
				Gap:         8,
				BorderWidth: 1,
				LineHeight:  16,
			},
		},
	}
}

func NewTheme(name string) Theme {
	theme := DefaultTheme()
	if trimmed := strings.TrimSpace(name); trimmed != "" {
		theme.Name = trimmed
	}
	return theme
}

func ResolveTheme(theme *Theme) Theme {
	base := DefaultTheme()
	if theme == nil {
		return base
	}
	return mergeTheme(base, *theme)
}

func mergeTheme(base Theme, override Theme) Theme {
	if override.Name != "" {
		base.Name = override.Name
	}
	base.Palette = mergePalette(base.Palette, override.Palette)
	base.Spacing = mergeSpacing(base.Spacing, override.Spacing)
	base.Radius = mergeRadius(base.Radius, override.Radius)
	base.Stroke = mergeStroke(base.Stroke, override.Stroke)
	base.Typography = mergeTypography(base.Typography, override.Typography)
	base.Components = mergeComponents(base.Components, override.Components)
	return base
}

func mergePalette(base ThemePalette, override ThemePalette) ThemePalette {
	base.Surface = ThemeSurfacePalette{
		Base:     colorOrDefault(override.Surface.Base, base.Surface.Base),
		Elevated: colorOrDefault(override.Surface.Elevated, base.Surface.Elevated),
		Sunken:   colorOrDefault(override.Surface.Sunken, base.Surface.Sunken),
		Panel:    colorOrDefault(override.Surface.Panel, base.Surface.Panel),
		Card:     colorOrDefault(override.Surface.Card, base.Surface.Card),
		Field:    colorOrDefault(override.Surface.Field, base.Surface.Field),
	}
	base.Text = ThemeTextPalette{
		Strong:  colorOrDefault(override.Text.Strong, base.Text.Strong),
		Muted:   colorOrDefault(override.Text.Muted, base.Text.Muted),
		Inverse: colorOrDefault(override.Text.Inverse, base.Text.Inverse),
		Accent:  colorOrDefault(override.Text.Accent, base.Text.Accent),
	}
	base.Accent = ThemeAccentPalette{
		Primary:   colorOrDefault(override.Accent.Primary, base.Accent.Primary),
		Secondary: colorOrDefault(override.Accent.Secondary, base.Accent.Secondary),
		Warning:   colorOrDefault(override.Accent.Warning, base.Accent.Warning),
	}
	base.Status = ThemeStatusPalette{
		Success: colorOrDefault(override.Status.Success, base.Status.Success),
		Warning: colorOrDefault(override.Status.Warning, base.Status.Warning),
		Error:   colorOrDefault(override.Status.Error, base.Status.Error),
	}
	base.Overlay = colorOrDefault(override.Overlay, base.Overlay)
	return base
}

func mergeSpacing(base ThemeSpacing, override ThemeSpacing) ThemeSpacing {
	base.XS = floatOrDefault(override.XS, base.XS)
	base.SM = floatOrDefault(override.SM, base.SM)
	base.MD = floatOrDefault(override.MD, base.MD)
	base.LG = floatOrDefault(override.LG, base.LG)
	base.XL = floatOrDefault(override.XL, base.XL)
	return base
}

func mergeRadius(base ThemeRadius, override ThemeRadius) ThemeRadius {
	base.SM = floatOrDefault(override.SM, base.SM)
	base.MD = floatOrDefault(override.MD, base.MD)
	base.LG = floatOrDefault(override.LG, base.LG)
	return base
}

func mergeStroke(base ThemeStroke, override ThemeStroke) ThemeStroke {
	base.Thin = floatOrDefault(override.Thin, base.Thin)
	base.Medium = floatOrDefault(override.Medium, base.Medium)
	base.Focus = floatOrDefault(override.Focus, base.Focus)
	return base
}

func mergeTypography(base ThemeTypography, override ThemeTypography) ThemeTypography {
	base.LabelLineHeight = floatOrDefault(override.LabelLineHeight, base.LabelLineHeight)
	base.BodyLineHeight = floatOrDefault(override.BodyLineHeight, base.BodyLineHeight)
	base.CodeLineHeight = floatOrDefault(override.CodeLineHeight, base.CodeLineHeight)
	return base
}

func mergeComponents(base ThemeComponents, override ThemeComponents) ThemeComponents {
	base.InputField = mergeFieldTheme(base.InputField, override.InputField)
	base.Textarea = mergeFieldTheme(base.Textarea, override.Textarea)
	base.Dropdown = mergeFieldTheme(base.Dropdown, override.Dropdown)
	base.Checkbox = mergeCheckboxTheme(base.Checkbox, override.Checkbox)
	base.Toggle = mergeToggleTheme(base.Toggle, override.Toggle)
	base.Slider = mergeSliderTheme(base.Slider, override.Slider)
	base.ProgressBar = mergeProgressBarTheme(base.ProgressBar, override.ProgressBar)
	base.Panel = mergePanelTheme(base.Panel, override.Panel)
	base.Card = mergeCardTheme(base.Card, override.Card)
	base.MenuButton = mergeMenuButtonTheme(base.MenuButton, override.MenuButton)
	base.Dialog = mergeDialogTheme(base.Dialog, override.Dialog)
	base.HUDBar = mergeHUDBarTheme(base.HUDBar, override.HUDBar)
	base.InventoryGrid = mergeInventoryGridTheme(base.InventoryGrid, override.InventoryGrid)
	base.Modal = mergeModalTheme(base.Modal, override.Modal)
	base.Tooltip = mergeTooltipTheme(base.Tooltip, override.Tooltip)
	return base
}

func mergeFieldTheme(base FieldTheme, override FieldTheme) FieldTheme {
	base.Background = colorOrDefault(override.Background, base.Background)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Text = colorOrDefault(override.Text, base.Text)
	base.Placeholder = colorOrDefault(override.Placeholder, base.Placeholder)
	base.Label = colorOrDefault(override.Label, base.Label)
	base.Caret = colorOrDefault(override.Caret, base.Caret)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.Height = floatOrDefault(override.Height, base.Height)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	base.LineHeight = floatOrDefault(override.LineHeight, base.LineHeight)
	return base
}

func mergeCheckboxTheme(base CheckboxTheme, override CheckboxTheme) CheckboxTheme {
	base.Label = colorOrDefault(override.Label, base.Label)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Box = colorOrDefault(override.Box, base.Box)
	base.BoxSelected = colorOrDefault(override.BoxSelected, base.BoxSelected)
	base.BoxDisabled = colorOrDefault(override.BoxDisabled, base.BoxDisabled)
	base.Check = colorOrDefault(override.Check, base.Check)
	base.Size = floatOrDefault(override.Size, base.Size)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	return base
}

func mergeToggleTheme(base ToggleTheme, override ToggleTheme) ToggleTheme {
	base.Label = colorOrDefault(override.Label, base.Label)
	base.TrackOff = colorOrDefault(override.TrackOff, base.TrackOff)
	base.TrackOn = colorOrDefault(override.TrackOn, base.TrackOn)
	base.TrackDisabled = colorOrDefault(override.TrackDisabled, base.TrackDisabled)
	base.Thumb = colorOrDefault(override.Thumb, base.Thumb)
	base.TrackWidth = floatOrDefault(override.TrackWidth, base.TrackWidth)
	base.TrackHeight = floatOrDefault(override.TrackHeight, base.TrackHeight)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.ThumbSize = floatOrDefault(override.ThumbSize, base.ThumbSize)
	return base
}

func mergeSliderTheme(base SliderTheme, override SliderTheme) SliderTheme {
	base.Label = colorOrDefault(override.Label, base.Label)
	base.Value = colorOrDefault(override.Value, base.Value)
	base.Track = colorOrDefault(override.Track, base.Track)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Fill = colorOrDefault(override.Fill, base.Fill)
	base.Thumb = colorOrDefault(override.Thumb, base.Thumb)
	base.TrackHeight = floatOrDefault(override.TrackHeight, base.TrackHeight)
	base.ThumbWidth = floatOrDefault(override.ThumbWidth, base.ThumbWidth)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	return base
}

func mergeProgressBarTheme(base ProgressBarTheme, override ProgressBarTheme) ProgressBarTheme {
	base.Label = colorOrDefault(override.Label, base.Label)
	base.Value = colorOrDefault(override.Value, base.Value)
	base.Track = colorOrDefault(override.Track, base.Track)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Fill = colorOrDefault(override.Fill, base.Fill)
	base.TrackHeight = floatOrDefault(override.TrackHeight, base.TrackHeight)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	return base
}

func mergePanelTheme(base PanelTheme, override PanelTheme) PanelTheme {
	base.Background = colorOrDefault(override.Background, base.Background)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.TitleText = colorOrDefault(override.TitleText, base.TitleText)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	return base
}

func mergeCardTheme(base CardTheme, override CardTheme) CardTheme {
	base.Background = colorOrDefault(override.Background, base.Background)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.TitleText = colorOrDefault(override.TitleText, base.TitleText)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	return base
}

func mergeMenuButtonTheme(base MenuButtonTheme, override MenuButtonTheme) MenuButtonTheme {
	base.Default = mergeStateColors(base.Default, override.Default)
	base.Focused = mergeStateColors(base.Focused, override.Focused)
	base.Selected = mergeStateColors(base.Selected, override.Selected)
	base.Disabled = mergeStateColors(base.Disabled, override.Disabled)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	return base
}

func mergeDialogTheme(base DialogTheme, override DialogTheme) DialogTheme {
	base.BodyText = colorOrDefault(override.BodyText, base.BodyText)
	return base
}

func mergeHUDBarTheme(base HUDBarTheme, override HUDBarTheme) HUDBarTheme {
	base.Track = colorOrDefault(override.Track, base.Track)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Fill = colorOrDefault(override.Fill, base.Fill)
	base.Text = colorOrDefault(override.Text, base.Text)
	base.TrackHeight = floatOrDefault(override.TrackHeight, base.TrackHeight)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	return base
}

func mergeInventoryGridTheme(base InventoryGridTheme, override InventoryGridTheme) InventoryGridTheme {
	base.SlotBackground = colorOrDefault(override.SlotBackground, base.SlotBackground)
	base.SlotBorder = colorOrDefault(override.SlotBorder, base.SlotBorder)
	base.SlotText = colorOrDefault(override.SlotText, base.SlotText)
	base.SlotMuted = colorOrDefault(override.SlotMuted, base.SlotMuted)
	base.IconFill = colorOrDefault(override.IconFill, base.IconFill)
	base.SlotPadding = floatOrDefault(override.SlotPadding, base.SlotPadding)
	base.SlotGap = floatOrDefault(override.SlotGap, base.SlotGap)
	base.SlotBorderWidth = floatOrDefault(override.SlotBorderWidth, base.SlotBorderWidth)
	return base
}

func mergeModalTheme(base ModalTheme, override ModalTheme) ModalTheme {
	base.Backdrop = colorOrDefault(override.Backdrop, base.Backdrop)
	base.Surface = colorOrDefault(override.Surface, base.Surface)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.TitleText = colorOrDefault(override.TitleText, base.TitleText)
	base.BodyText = colorOrDefault(override.BodyText, base.BodyText)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	return base
}

func mergeTooltipTheme(base TooltipTheme, override TooltipTheme) TooltipTheme {
	base.Surface = colorOrDefault(override.Surface, base.Surface)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.TitleText = colorOrDefault(override.TitleText, base.TitleText)
	base.BodyText = colorOrDefault(override.BodyText, base.BodyText)
	base.Padding = floatOrDefault(override.Padding, base.Padding)
	base.Gap = floatOrDefault(override.Gap, base.Gap)
	base.BorderWidth = floatOrDefault(override.BorderWidth, base.BorderWidth)
	base.LineHeight = floatOrDefault(override.LineHeight, base.LineHeight)
	return base
}

func mergeStateColors(base StateColors, override StateColors) StateColors {
	base.Background = colorOrDefault(override.Background, base.Background)
	base.Border = colorOrDefault(override.Border, base.Border)
	base.Text = colorOrDefault(override.Text, base.Text)
	return base
}

func colorOrDefault(override color.Color, fallback color.Color) color.Color {
	if override != nil {
		return override
	}
	return fallback
}

func floatOrDefault(override float64, fallback float64) float64 {
	if override != 0 {
		return override
	}
	return fallback
}
