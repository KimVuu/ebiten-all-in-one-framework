package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func showcaseGroupStyle() ebitenui.Style {
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Column,
		Padding:         ebitenui.All(16),
		Gap:             12,
		BackgroundColor: color.RGBA{R: 24, G: 31, B: 43, A: 255},
		BorderColor:     color.RGBA{R: 88, G: 110, B: 140, A: 255},
		BorderWidth:     1,
	}
}

func showcaseGroupTitleStyle() ebitenui.Style {
	return ebitenui.Style{
		Color: color.RGBA{R: 242, G: 246, B: 252, A: 255},
	}
}

func showcaseGroupCopyStyle() ebitenui.Style {
	return ebitenui.Style{
		Width:      ebitenui.Fill(),
		Color:      color.RGBA{R: 178, G: 190, B: 207, A: 255},
		LineHeight: 16,
	}
}

func detailSectionStyle() ebitenui.Style {
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Column,
		Padding:         ebitenui.All(16),
		Gap:             12,
		BackgroundColor: color.RGBA{R: 24, G: 31, B: 43, A: 255},
		BorderColor:     color.RGBA{R: 86, G: 104, B: 128, A: 255},
		BorderWidth:     1,
	}
}

func detailTitleStyle() ebitenui.Style {
	return ebitenui.Style{
		Color: color.RGBA{R: 240, G: 244, B: 252, A: 255},
	}
}
