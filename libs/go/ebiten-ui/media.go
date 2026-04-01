package ebitenui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type InteractionState struct {
	Hovered  bool
	Focused  bool
	Pressed  bool
	Disabled bool
	Selected bool
}

type ScrollState struct {
	OffsetX float64
	OffsetY float64
}

type ImageSource struct {
	Image  *ebiten.Image
	Width  float64
	Height float64
	Fill   color.Color
}

func SolidImage(width, height float64, fill color.Color) ImageSource {
	return ImageSource{
		Width:  width,
		Height: height,
		Fill:   fill,
	}
}

func RasterImage(image *ebiten.Image) ImageSource {
	if image == nil {
		return ImageSource{}
	}
	bounds := image.Bounds()
	return ImageSource{
		Image:  image,
		Width:  float64(bounds.Dx()),
		Height: float64(bounds.Dy()),
	}
}

func (s ImageSource) intrinsicSize() (float64, float64) {
	width := s.Width
	height := s.Height

	if s.Image != nil {
		bounds := s.Image.Bounds()
		if width == 0 {
			width = float64(bounds.Dx())
		}
		if height == 0 {
			height = float64(bounds.Dy())
		}
	}

	return width, height
}
