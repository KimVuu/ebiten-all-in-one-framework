package ui

import (
	"fmt"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	fontassets "github.com/kimyechan/ebiten-aio-framework/projects/dice-rogue/asset/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const defaultTextFontSize = 16

func ApplyTextFace() error {
	face, err := LoadTextFace()
	if err != nil {
		return err
	}
	ebitenui.SetTextFace(face)
	return nil
}

func LoadTextFace() (font.Face, error) {
	fontData, err := opentype.Parse(fontassets.NeoDunggeunmoProRegular)
	if err != nil {
		return nil, fmt.Errorf("parse dice-rogue font: %w", err)
	}
	face, err := opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    defaultTextFontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create dice-rogue font face: %w", err)
	}
	return face, nil
}
