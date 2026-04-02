package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

type showcaseChrome struct {
	RootBackground  color.Color
	PanelBackground color.Color
	PanelBorder     color.Color
	TextStrong      color.Color
	TextMuted       color.Color
	Accent          color.Color
	AccentSoft      color.Color
	BadgeText       color.Color
	CodeBackground  color.Color
	CodeBorder      color.Color
	CodeText        color.Color
}

type showcaseThemePreset struct {
	ID     string
	Title  string
	Theme  ebitenui.Theme
	Chrome showcaseChrome
}

func showcaseThemePresets() []showcaseThemePreset {
	return []showcaseThemePreset{
		buildDefaultShowcasePreset(),
		buildForestShowcasePreset(),
		buildEmberShowcasePreset(),
	}
}

func showcaseThemePresetByID(id string) showcaseThemePreset {
	for _, preset := range showcaseThemePresets() {
		if preset.ID == id {
			return preset
		}
	}
	return buildDefaultShowcasePreset()
}

func initialShowcaseThemePreset(id string) string {
	if id == "" {
		return "default"
	}
	return showcaseThemePresetByID(id).ID
}

func buildDefaultShowcasePreset() showcaseThemePreset {
	theme := ebitenui.DefaultTheme()
	return showcaseThemePreset{
		ID:    "default",
		Title: "Default",
		Theme: theme,
		Chrome: showcaseChrome{
			RootBackground:  color.RGBA{R: 13, G: 18, B: 27, A: 255},
			PanelBackground: color.RGBA{R: 24, G: 31, B: 43, A: 255},
			PanelBorder:     color.RGBA{R: 88, G: 110, B: 140, A: 255},
			TextStrong:      color.RGBA{R: 242, G: 246, B: 252, A: 255},
			TextMuted:       color.RGBA{R: 178, G: 190, B: 207, A: 255},
			Accent:          color.RGBA{R: 80, G: 160, B: 255, A: 255},
			AccentSoft:      color.RGBA{R: 72, G: 211, B: 161, A: 255},
			BadgeText:       color.RGBA{R: 8, G: 12, B: 20, A: 255},
			CodeBackground:  color.RGBA{R: 17, G: 22, B: 31, A: 255},
			CodeBorder:      color.RGBA{R: 63, G: 78, B: 101, A: 255},
			CodeText:        color.RGBA{R: 213, G: 223, B: 238, A: 255},
		},
	}
}

func buildForestShowcasePreset() showcaseThemePreset {
	theme := ebitenui.NewTheme("forest")
	theme.Components.Panel.Background = color.RGBA{R: 20, G: 39, B: 34, A: 255}
	theme.Components.Panel.Border = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.Panel.TitleText = color.RGBA{R: 232, G: 246, B: 238, A: 255}
	theme.Components.Card.Background = color.RGBA{R: 16, G: 31, B: 28, A: 255}
	theme.Components.Card.Border = color.RGBA{R: 66, G: 142, B: 118, A: 255}
	theme.Components.Card.TitleText = color.RGBA{R: 226, G: 243, B: 235, A: 255}
	theme.Components.InputField.Background = color.RGBA{R: 12, G: 27, B: 24, A: 255}
	theme.Components.InputField.Border = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.InputField.Text = color.RGBA{R: 232, G: 246, B: 238, A: 255}
	theme.Components.InputField.Placeholder = color.RGBA{R: 144, G: 185, B: 170, A: 255}
	theme.Components.InputField.Label = color.RGBA{R: 144, G: 185, B: 170, A: 255}
	theme.Components.Textarea = theme.Components.InputField
	theme.Components.Dropdown = theme.Components.InputField
	theme.Components.Toggle.TrackOn = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.Toggle.TrackOff = color.RGBA{R: 28, G: 55, B: 48, A: 255}
	theme.Components.ProgressBar.Fill = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.ProgressBar.Track = color.RGBA{R: 16, G: 31, B: 28, A: 255}
	theme.Components.HUDBar.Fill = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.HUDBar.Track = color.RGBA{R: 16, G: 31, B: 28, A: 255}
	theme.Components.MenuButton.Default.Background = color.RGBA{R: 20, G: 39, B: 34, A: 255}
	theme.Components.MenuButton.Default.Border = color.RGBA{R: 66, G: 142, B: 118, A: 255}
	theme.Components.MenuButton.Default.Text = color.RGBA{R: 232, G: 246, B: 238, A: 255}
	theme.Components.MenuButton.Selected.Background = color.RGBA{R: 88, G: 187, B: 152, A: 255}
	theme.Components.MenuButton.Selected.Border = color.RGBA{R: 226, G: 206, B: 120, A: 255}
	theme.Components.MenuButton.Selected.Text = color.RGBA{R: 14, G: 26, B: 22, A: 255}
	theme.Components.MenuButton.Focused.Background = color.RGBA{R: 135, G: 212, B: 180, A: 255}
	theme.Components.MenuButton.Focused.Border = color.RGBA{R: 226, G: 206, B: 120, A: 255}
	theme.Components.MenuButton.Focused.Text = color.RGBA{R: 14, G: 26, B: 22, A: 255}
	return showcaseThemePreset{
		ID:    "forest",
		Title: "Forest",
		Theme: theme,
		Chrome: showcaseChrome{
			RootBackground:  color.RGBA{R: 8, G: 20, B: 18, A: 255},
			PanelBackground: color.RGBA{R: 20, G: 39, B: 34, A: 255},
			PanelBorder:     color.RGBA{R: 88, G: 187, B: 152, A: 255},
			TextStrong:      color.RGBA{R: 232, G: 246, B: 238, A: 255},
			TextMuted:       color.RGBA{R: 144, G: 185, B: 170, A: 255},
			Accent:          color.RGBA{R: 88, G: 187, B: 152, A: 255},
			AccentSoft:      color.RGBA{R: 173, G: 220, B: 143, A: 255},
			BadgeText:       color.RGBA{R: 14, G: 26, B: 22, A: 255},
			CodeBackground:  color.RGBA{R: 10, G: 24, B: 22, A: 255},
			CodeBorder:      color.RGBA{R: 66, G: 142, B: 118, A: 255},
			CodeText:        color.RGBA{R: 214, G: 236, B: 227, A: 255},
		},
	}
}

func buildEmberShowcasePreset() showcaseThemePreset {
	theme := ebitenui.NewTheme("ember")
	theme.Components.Panel.Background = color.RGBA{R: 52, G: 26, B: 22, A: 255}
	theme.Components.Panel.Border = color.RGBA{R: 235, G: 131, B: 72, A: 255}
	theme.Components.Panel.TitleText = color.RGBA{R: 252, G: 239, B: 228, A: 255}
	theme.Components.Card.Background = color.RGBA{R: 42, G: 21, B: 18, A: 255}
	theme.Components.Card.Border = color.RGBA{R: 186, G: 97, B: 54, A: 255}
	theme.Components.Card.TitleText = color.RGBA{R: 252, G: 239, B: 228, A: 255}
	theme.Components.InputField.Background = color.RGBA{R: 34, G: 18, B: 16, A: 255}
	theme.Components.InputField.Border = color.RGBA{R: 235, G: 131, B: 72, A: 255}
	theme.Components.InputField.Text = color.RGBA{R: 252, G: 239, B: 228, A: 255}
	theme.Components.InputField.Placeholder = color.RGBA{R: 214, G: 173, B: 146, A: 255}
	theme.Components.InputField.Label = color.RGBA{R: 214, G: 173, B: 146, A: 255}
	theme.Components.Textarea = theme.Components.InputField
	theme.Components.Dropdown = theme.Components.InputField
	theme.Components.Toggle.TrackOn = color.RGBA{R: 235, G: 131, B: 72, A: 255}
	theme.Components.Toggle.TrackOff = color.RGBA{R: 73, G: 37, B: 30, A: 255}
	theme.Components.ProgressBar.Fill = color.RGBA{R: 244, G: 159, B: 85, A: 255}
	theme.Components.ProgressBar.Track = color.RGBA{R: 42, G: 21, B: 18, A: 255}
	theme.Components.HUDBar.Fill = color.RGBA{R: 244, G: 159, B: 85, A: 255}
	theme.Components.HUDBar.Track = color.RGBA{R: 42, G: 21, B: 18, A: 255}
	theme.Components.MenuButton.Default.Background = color.RGBA{R: 52, G: 26, B: 22, A: 255}
	theme.Components.MenuButton.Default.Border = color.RGBA{R: 186, G: 97, B: 54, A: 255}
	theme.Components.MenuButton.Default.Text = color.RGBA{R: 252, G: 239, B: 228, A: 255}
	theme.Components.MenuButton.Selected.Background = color.RGBA{R: 244, G: 159, B: 85, A: 255}
	theme.Components.MenuButton.Selected.Border = color.RGBA{R: 255, G: 205, B: 132, A: 255}
	theme.Components.MenuButton.Selected.Text = color.RGBA{R: 34, G: 18, B: 16, A: 255}
	theme.Components.MenuButton.Focused.Background = color.RGBA{R: 252, G: 197, B: 128, A: 255}
	theme.Components.MenuButton.Focused.Border = color.RGBA{R: 255, G: 205, B: 132, A: 255}
	theme.Components.MenuButton.Focused.Text = color.RGBA{R: 34, G: 18, B: 16, A: 255}
	return showcaseThemePreset{
		ID:    "ember",
		Title: "Ember",
		Theme: theme,
		Chrome: showcaseChrome{
			RootBackground:  color.RGBA{R: 22, G: 12, B: 12, A: 255},
			PanelBackground: color.RGBA{R: 52, G: 26, B: 22, A: 255},
			PanelBorder:     color.RGBA{R: 235, G: 131, B: 72, A: 255},
			TextStrong:      color.RGBA{R: 252, G: 239, B: 228, A: 255},
			TextMuted:       color.RGBA{R: 214, G: 173, B: 146, A: 255},
			Accent:          color.RGBA{R: 244, G: 159, B: 85, A: 255},
			AccentSoft:      color.RGBA{R: 255, G: 205, B: 132, A: 255},
			BadgeText:       color.RGBA{R: 34, G: 18, B: 16, A: 255},
			CodeBackground:  color.RGBA{R: 30, G: 16, B: 15, A: 255},
			CodeBorder:      color.RGBA{R: 186, G: 97, B: 54, A: 255},
			CodeText:        color.RGBA{R: 247, G: 226, B: 208, A: 255},
		},
	}
}
