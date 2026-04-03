package ui

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	fontassets "github.com/kimyechan/ebiten-aio-framework/projects/dice-rogue/asset/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

const defaultFontSize = 13

type Renderer struct {
	face font.Face

	mu    sync.RWMutex
	cache map[textCacheKey]cachedText
}

type textCacheKey struct {
	Value string
	Color color.RGBA
}

type cachedText struct {
	Image  *ebiten.Image
	Width  int
	Height int
}

func NewRenderer() (*Renderer, error) {
	face, err := newFaceFromBytes(fontassets.NeoDunggeunmoProRegular, defaultFontSize)
	if err != nil {
		return nil, err
	}
	return newRendererWithFace(face), nil
}

func NewFallbackRenderer() *Renderer {
	return newRendererWithFace(basicfont.Face7x13)
}

func newRendererWithFace(face font.Face) *Renderer {
	if face == nil {
		face = basicfont.Face7x13
	}
	return &Renderer{
		face:  face,
		cache: map[textCacheKey]cachedText{},
	}
}

func newFaceFromBytes(data []byte, size float64) (font.Face, error) {
	fontData, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse dice-rogue font: %w", err)
	}
	face, err := opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create dice-rogue font face: %w", err)
	}
	return face, nil
}

func (r *Renderer) Draw(screen *ebiten.Image, dom *ebitenui.DOM, viewport ebitenui.Viewport) *ebitenui.LayoutNode {
	if r == nil || screen == nil || dom == nil {
		return nil
	}

	layout := dom.Layout(viewport)
	if layout == nil {
		return nil
	}

	r.drawNode(screen, layout)
	return layout
}

func (r *Renderer) drawNode(screen *ebiten.Image, layout *ebitenui.LayoutNode) {
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
		r.drawTextLines(screen, layout, []string{layout.Node.Text})
	case ebitenui.TagTextBlock:
		r.drawTextLines(screen, layout, layout.TextLines)
	case ebitenui.TagImage:
		drawImage(screen, layout)
	}

	for _, child := range layout.Children {
		if layout.ClipChildren && !intersects(layout.Frame, child.Frame) {
			continue
		}
		r.drawNode(screen, child)
	}
}

func (r *Renderer) drawTextLines(screen *ebiten.Image, layout *ebitenui.LayoutNode, lines []string) {
	if len(lines) == 0 {
		return
	}

	slotHeight := layout.Frame.Height / float64(len(lines))
	if slotHeight <= 0 {
		slotHeight = float64(r.face.Metrics().Height.Ceil())
	}

	for index, line := range lines {
		cached := r.cachedLineImage(line, layout.Node.Props.Style.Color)
		if cached.Image == nil || cached.Width <= 0 || cached.Height <= 0 {
			continue
		}

		scaleX := 1.0
		if layout.Frame.Width > 0 {
			scaleX = minRendererFloat(1, layout.Frame.Width/float64(cached.Width))
		}
		scaleY := minRendererFloat(1, slotHeight/float64(cached.Height))
		scale := minRendererFloat(scaleX, scaleY)
		if scale <= 0 {
			continue
		}

		drawWidth := float64(cached.Width) * scale
		drawHeight := float64(cached.Height) * scale
		x := layout.Frame.X
		switch layout.Node.Props.Style.TextAlign {
		case ebitenui.TextAlignCenter:
			x = layout.Frame.X + maxFloat(0, (layout.Frame.Width-drawWidth)*0.5)
		case ebitenui.TextAlignEnd:
			x = layout.Frame.X + maxFloat(0, layout.Frame.Width-drawWidth)
		}

		y := layout.Frame.Y + float64(index)*slotHeight + maxFloat(0, (slotHeight-drawHeight)*0.5)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(scale, scale)
		options.GeoM.Translate(x, y)
		screen.DrawImage(cached.Image, options)
	}
}

func (r *Renderer) cachedLineImage(value string, fill color.Color) cachedText {
	key := textCacheKey{
		Value: value,
		Color: rgbaColor(fill),
	}

	r.mu.RLock()
	cached, ok := r.cache[key]
	r.mu.RUnlock()
	if ok {
		return cached
	}

	image := r.renderLineImage(value, key.Color)
	cached = cachedText{
		Image:  image,
		Width:  image.Bounds().Dx(),
		Height: image.Bounds().Dy(),
	}

	r.mu.Lock()
	r.cache[key] = cached
	r.mu.Unlock()
	return cached
}

func (r *Renderer) renderLineImage(value string, fill color.RGBA) *ebiten.Image {
	bounds := text.BoundString(r.face, value)
	width := bounds.Dx()
	if width <= 0 {
		width = 1
	}
	height := r.face.Metrics().Height.Ceil()
	if height <= 0 {
		height = 1
	}

	image := ebiten.NewImage(width, height)
	text.Draw(image, value, r.face, -bounds.Min.X, r.face.Metrics().Ascent.Ceil(), fill)
	return image
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

func rgbaColor(fill color.Color) color.RGBA {
	if fill == nil {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	converted, ok := color.RGBAModel.Convert(fill).(color.RGBA)
	if ok {
		return converted
	}
	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

func minRendererFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
