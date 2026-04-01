package renderer

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"golang.org/x/image/font/basicfont"
)

type Renderer struct{}

func New() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Draw(screen *ebiten.Image, dom *ebitenui.DOM, viewport ebitenui.Viewport) *ebitenui.LayoutNode {
	if r == nil || screen == nil || dom == nil {
		return nil
	}

	layout := dom.Layout(viewport)
	if layout == nil {
		return nil
	}

	drawNode(screen, layout)
	return layout
}

func drawNode(screen *ebiten.Image, layout *ebitenui.LayoutNode) {
	if layout == nil || layout.Node == nil {
		return
	}

	style := layout.Node.Props.Style
	if style.BackgroundColor != nil {
		vector.DrawFilledRect(
			screen,
			float32(layout.Frame.X),
			float32(layout.Frame.Y),
			float32(layout.Frame.Width),
			float32(layout.Frame.Height),
			style.BackgroundColor,
			false,
		)
	}

	if style.BorderColor != nil && style.BorderWidth > 0 {
		vector.StrokeRect(
			screen,
			float32(layout.Frame.X),
			float32(layout.Frame.Y),
			float32(layout.Frame.Width),
			float32(layout.Frame.Height),
			float32(style.BorderWidth),
			style.BorderColor,
			false,
		)
	}

	drawInteractionState(screen, layout)

	switch layout.Node.Tag {
	case ebitenui.TagText:
		drawTextLines(screen, layout, []string{layout.Node.Text})
	case ebitenui.TagTextBlock:
		drawTextLines(screen, layout, layout.TextLines)
	case ebitenui.TagImage:
		drawImage(screen, layout)
	}

	for _, child := range layout.Children {
		if layout.ClipChildren && !intersects(layout.Frame, child.Frame) {
			continue
		}
		drawNode(screen, child)
	}
}

func drawTextLines(screen *ebiten.Image, layout *ebitenui.LayoutNode, lines []string) {
	textColor := layout.Node.Props.Style.Color
	if textColor == nil {
		textColor = color.White
	}

	ascent := basicfont.Face7x13.Metrics().Ascent.Ceil()
	lineHeight := int(uidomLineHeight(layout.Node.Props.Style))
	for i, line := range lines {
		x := int(layout.Frame.X)
		lineWidth := fontWidth(line)
		switch layout.Node.Props.Style.TextAlign {
		case ebitenui.TextAlignCenter:
			x = int(layout.Frame.X + maxFloat(0, (layout.Frame.Width-float64(lineWidth))*0.5))
		case ebitenui.TextAlignEnd:
			x = int(layout.Frame.X + maxFloat(0, layout.Frame.Width-float64(lineWidth)))
		}

		y := int(layout.Frame.Y) + ascent + i*lineHeight
		text.Draw(screen, line, basicfont.Face7x13, x, y, textColor)
	}
}

func drawImage(screen *ebiten.Image, layout *ebitenui.LayoutNode) {
	source := layout.Node.Props.Image
	if source.Fill != nil {
		vector.DrawFilledRect(
			screen,
			float32(layout.Frame.X),
			float32(layout.Frame.Y),
			float32(layout.Frame.Width),
			float32(layout.Frame.Height),
			source.Fill,
			false,
		)
	}

	if source.Image == nil {
		return
	}

	bounds := source.Image.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		return
	}

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(layout.Frame.Width/float64(bounds.Dx()), layout.Frame.Height/float64(bounds.Dy()))
	options.GeoM.Translate(layout.Frame.X, layout.Frame.Y)
	screen.DrawImage(source.Image, options)
}

func drawInteractionState(screen *ebiten.Image, layout *ebitenui.LayoutNode) {
	state := layout.Node.Props.State
	frame := image.Rect(
		int(layout.Frame.X),
		int(layout.Frame.Y),
		int(layout.Frame.X+layout.Frame.Width),
		int(layout.Frame.Y+layout.Frame.Height),
	)

	if state.Hovered {
		vector.DrawFilledRect(screen, float32(frame.Min.X), float32(frame.Min.Y), float32(frame.Dx()), float32(frame.Dy()), color.RGBA{R: 255, G: 255, B: 255, A: 24}, false)
	}
	if state.Pressed {
		vector.DrawFilledRect(screen, float32(frame.Min.X), float32(frame.Min.Y), float32(frame.Dx()), float32(frame.Dy()), color.RGBA{R: 0, G: 0, B: 0, A: 48}, false)
	}
	if state.Focused || state.Selected {
		vector.StrokeRect(screen, float32(frame.Min.X), float32(frame.Min.Y), float32(frame.Dx()), float32(frame.Dy()), 2, color.RGBA{R: 255, G: 214, B: 92, A: 255}, false)
	}
	if state.Disabled {
		vector.DrawFilledRect(screen, float32(frame.Min.X), float32(frame.Min.Y), float32(frame.Dx()), float32(frame.Dy()), color.RGBA{R: 20, G: 20, B: 24, A: 90}, false)
	}
}

func intersects(a, b ebitenui.Rect) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X && a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}

func fontWidth(value string) int {
	return text.BoundString(basicfont.Face7x13, value).Dx()
}

func uidomLineHeight(style ebitenui.Style) float64 {
	if style.LineHeight > 0 {
		return style.LineHeight
	}
	return float64(basicfont.Face7x13.Metrics().Height.Ceil())
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
