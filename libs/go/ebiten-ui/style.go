package ebitenui

import "image/color"

type Direction int

const (
	Column Direction = iota
	Row
)

type TextAlign int

const (
	TextAlignStart TextAlign = iota
	TextAlignCenter
	TextAlignEnd
)

type LengthKind int

const (
	LengthAuto LengthKind = iota
	LengthPx
	LengthFill
)

type Length struct {
	Kind  LengthKind
	Value float64
}

func Auto() Length {
	return Length{Kind: LengthAuto}
}

func Px(value float64) Length {
	return Length{Kind: LengthPx, Value: value}
}

func Fill() Length {
	return Length{Kind: LengthFill}
}

type Insets struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func All(value float64) Insets {
	return Insets{
		Top:    value,
		Right:  value,
		Bottom: value,
		Left:   value,
	}
}

func (i Insets) Horizontal() float64 {
	return i.Left + i.Right
}

func (i Insets) Vertical() float64 {
	return i.Top + i.Bottom
}

type Style struct {
	Width           Length
	Height          Length
	Direction       Direction
	Padding         Insets
	Gap             float64
	BackgroundColor color.Color
	BorderColor     color.Color
	BorderWidth     float64
	Color           color.Color
	TextAlign       TextAlign
	LineHeight      float64
}

func (s Style) directionOrDefault() Direction {
	if s.Direction == Row {
		return Row
	}
	return Column
}
