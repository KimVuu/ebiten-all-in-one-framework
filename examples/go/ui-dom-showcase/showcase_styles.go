package main

import (
	"image/color"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func showcaseGroupStyle() uidom.Style {
	return uidom.Style{
		Width:           uidom.Fill(),
		Direction:       uidom.Column,
		Padding:         uidom.All(16),
		Gap:             12,
		BackgroundColor: color.RGBA{R: 24, G: 31, B: 43, A: 255},
		BorderColor:     color.RGBA{R: 88, G: 110, B: 140, A: 255},
		BorderWidth:     1,
	}
}

func showcaseGroupTitleStyle() uidom.Style {
	return uidom.Style{
		Color: color.RGBA{R: 242, G: 246, B: 252, A: 255},
	}
}

func showcaseGroupCopyStyle() uidom.Style {
	return uidom.Style{
		Width:      uidom.Fill(),
		Color:      color.RGBA{R: 178, G: 190, B: 207, A: 255},
		LineHeight: 16,
	}
}

func detailSectionStyle() uidom.Style {
	return uidom.Style{
		Width:           uidom.Fill(),
		Direction:       uidom.Column,
		Padding:         uidom.All(16),
		Gap:             12,
		BackgroundColor: color.RGBA{R: 24, G: 31, B: 43, A: 255},
		BorderColor:     color.RGBA{R: 86, G: 104, B: 128, A: 255},
		BorderWidth:     1,
	}
}

func detailTitleStyle() uidom.Style {
	return uidom.Style{
		Color: color.RGBA{R: 240, G: 244, B: 252, A: 255},
	}
}
