package ebitenui_test

import (
	"image/color"
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestDefaultThemeProvidesNonZeroTokens(t *testing.T) {
	theme := ebitenui.DefaultTheme()

	if theme.Name == "" {
		t.Fatalf("expected default theme name")
	}
	if theme.Palette.Surface.Panel == nil {
		t.Fatalf("expected panel surface token")
	}
	if theme.Palette.Text.Strong == nil {
		t.Fatalf("expected strong text token")
	}
	if theme.Spacing.MD <= 0 {
		t.Fatalf("expected medium spacing token")
	}
	if theme.Radius.MD <= 0 {
		t.Fatalf("expected medium radius token")
	}
	if theme.Stroke.Thin <= 0 {
		t.Fatalf("expected thin stroke token")
	}
}

func TestInputFieldUsesThemeColors(t *testing.T) {
	fieldBackground := color.RGBA{R: 45, G: 51, B: 68, A: 255}
	fieldBorder := color.RGBA{R: 120, G: 132, B: 170, A: 255}
	fieldText := color.RGBA{R: 248, G: 244, B: 238, A: 255}
	fieldPlaceholder := color.RGBA{R: 160, G: 164, B: 180, A: 255}
	caret := color.RGBA{R: 255, G: 196, B: 82, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.InputField.Background = fieldBackground
	theme.Components.InputField.Border = fieldBorder
	theme.Components.InputField.Text = fieldText
	theme.Components.InputField.Placeholder = fieldPlaceholder
	theme.Components.InputField.Caret = caret

	withValue := ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:    "theme-input",
		Value: "Rune",
		State: ebitenui.InteractionState{Focused: true},
		Theme: &theme,
	})
	input, ok := ebitenui.New(withValue).FindByID("theme-input")
	if !ok {
		t.Fatalf("expected input node")
	}
	if !sameColor(input.Props.Style.BackgroundColor, fieldBackground) {
		t.Fatalf("expected themed input background")
	}
	if !sameColor(input.Props.Style.BorderColor, fieldBorder) {
		t.Fatalf("expected themed input border")
	}
	valueNode, ok := ebitenui.New(withValue).FindByID("theme-input-value")
	if !ok {
		t.Fatalf("expected value node")
	}
	if !sameColor(valueNode.Props.Style.Color, fieldText) {
		t.Fatalf("expected themed input text color")
	}
	caretNode, ok := ebitenui.New(withValue).FindByID("theme-input-caret")
	if !ok {
		t.Fatalf("expected caret node")
	}
	if !sameColor(caretNode.Props.Style.BackgroundColor, caret) {
		t.Fatalf("expected themed caret color")
	}

	withPlaceholder := ebitenui.InputField(ebitenui.InputFieldConfig{
		ID:          "theme-placeholder",
		Placeholder: "name@example.com",
		Theme:       &theme,
	})
	placeholderNode, ok := ebitenui.New(withPlaceholder).FindByID("theme-placeholder-placeholder")
	if !ok {
		t.Fatalf("expected placeholder node")
	}
	if !sameColor(placeholderNode.Props.Style.Color, fieldPlaceholder) {
		t.Fatalf("expected themed placeholder color")
	}
}

func TestToggleUsesThemeTrackColors(t *testing.T) {
	trackOn := color.RGBA{R: 96, G: 202, B: 156, A: 255}
	trackOff := color.RGBA{R: 38, G: 44, B: 58, A: 255}
	thumb := color.RGBA{R: 244, G: 246, B: 250, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.Toggle.TrackOn = trackOn
	theme.Components.Toggle.TrackOff = trackOff
	theme.Components.Toggle.Thumb = thumb

	onNode := ebitenui.Toggle(ebitenui.ToggleConfig{
		ID:      "music",
		Checked: true,
		Theme:   &theme,
	})
	track, ok := ebitenui.New(onNode).FindByID("music-track")
	if !ok {
		t.Fatalf("expected toggle track")
	}
	if !sameColor(track.Props.Style.BackgroundColor, trackOn) {
		t.Fatalf("expected selected track color")
	}
	thumbNode, ok := ebitenui.New(onNode).FindByID("music-thumb")
	if !ok {
		t.Fatalf("expected toggle thumb")
	}
	if !sameColor(thumbNode.Props.Style.BackgroundColor, thumb) {
		t.Fatalf("expected themed thumb color")
	}

	offNode := ebitenui.Toggle(ebitenui.ToggleConfig{
		ID:    "sfx",
		Theme: &theme,
	})
	offTrack, ok := ebitenui.New(offNode).FindByID("sfx-track")
	if !ok {
		t.Fatalf("expected off toggle track")
	}
	if !sameColor(offTrack.Props.Style.BackgroundColor, trackOff) {
		t.Fatalf("expected off track color")
	}
}

func TestProgressBarUsesThemeTintWhenTintUnset(t *testing.T) {
	fill := color.RGBA{R: 197, G: 116, B: 255, A: 255}
	track := color.RGBA{R: 31, G: 38, B: 52, A: 255}
	border := color.RGBA{R: 110, G: 118, B: 150, A: 255}

	theme := ebitenui.DefaultTheme()
	theme.Components.ProgressBar.Fill = fill
	theme.Components.ProgressBar.Track = track
	theme.Components.ProgressBar.Border = border

	node := ebitenui.ProgressBar(ebitenui.ProgressBarConfig{
		ID:      "mana",
		Label:   "Mana",
		Current: 60,
		Max:     100,
		Width:   200,
		Theme:   &theme,
	})

	dom := ebitenui.New(node)
	trackNode, ok := dom.FindByID("mana-track")
	if !ok {
		t.Fatalf("expected progress track")
	}
	if !sameColor(trackNode.Props.Style.BackgroundColor, track) {
		t.Fatalf("expected themed track color")
	}
	if !sameColor(trackNode.Props.Style.BorderColor, border) {
		t.Fatalf("expected themed border color")
	}
	fillNode, ok := dom.FindByID("mana-fill")
	if !ok {
		t.Fatalf("expected progress fill")
	}
	if !sameColor(fillNode.Props.Style.BackgroundColor, fill) {
		t.Fatalf("expected themed fill color")
	}
}

func sameColor(got color.Color, want color.Color) bool {
	if got == nil || want == nil {
		return got == want
	}
	gotRGBA := color.NRGBAModel.Convert(got).(color.NRGBA)
	wantRGBA := color.NRGBAModel.Convert(want).(color.NRGBA)
	return gotRGBA == wantRGBA
}
