package main

import (
	"sync"

	showcasefonts "github.com/kimyechan/ebiten-aio-framework/examples/go/ebiten-ui-showcase/asset/fonts"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

type showcaseFontPreset struct {
	ID    string
	Title string
	Face  font.Face
}

var (
	showcaseFontPresetsOnce sync.Once
	showcaseFontPresetList  []showcaseFontPreset
)

func showcaseFontPresets() []showcaseFontPreset {
	showcaseFontPresetsOnce.Do(func() {
		showcaseFontPresetList = []showcaseFontPreset{
			{
				ID:    "default",
				Title: "Default",
				Face:  basicfont.Face7x13,
			},
			{
				ID:    "neo-dunggeunmo",
				Title: "Neo Dunggeunmo",
				Face:  mustLoadShowcaseFont(showcasefonts.NeoDunggeunmoProRegular, basicfont.Face7x13),
			},
		}
	})

	cloned := make([]showcaseFontPreset, len(showcaseFontPresetList))
	copy(cloned, showcaseFontPresetList)
	return cloned
}

func showcaseFontPresetByID(id string) showcaseFontPreset {
	for _, preset := range showcaseFontPresets() {
		if preset.ID == id {
			return preset
		}
	}
	return showcaseFontPresets()[0]
}

func initialShowcaseFontPreset(id string) string {
	if id == "" {
		return "default"
	}
	return showcaseFontPresetByID(id).ID
}

func applyShowcaseFontPreset(id string) showcaseFontPreset {
	preset := showcaseFontPresetByID(initialShowcaseFontPreset(id))
	ebitenui.SetTextFace(preset.Face)
	return preset
}

func mustLoadShowcaseFont(data []byte, fallback font.Face) font.Face {
	if len(data) == 0 {
		return fallback
	}
	parsed, err := opentype.Parse(data)
	if err != nil {
		return fallback
	}
	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fallback
	}
	return face
}
