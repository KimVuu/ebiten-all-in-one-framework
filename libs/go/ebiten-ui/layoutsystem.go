package ebitenui

import "image/color"

const defaultPageScrollStep = 48

type PageLayoutConfig struct {
	ID              string
	ScrollID        string
	Header          *Node
	Content         *Node
	Footer          *Node
	Padding         float64
	Gap             float64
	BackgroundColor color.Color
	ScrollOffsetY   float64
	OnScrollChange  func(float64)
}

func PageLayout(cfg PageLayoutConfig) *Node {
	rootID := cfg.ID
	if rootID == "" {
		rootID = "page-root"
	}

	scrollID := cfg.ScrollID
	if scrollID == "" {
		scrollID = "page-scroll"
	}

	children := make([]*Node, 0, 2)
	if cfg.Header != nil {
		children = append(children, cfg.Header)
	}

	scrollChildren := make([]*Node, 0, 2)
	if cfg.Content != nil {
		scrollChildren = append(scrollChildren, cfg.Content)
	}
	if cfg.Footer != nil {
		scrollChildren = append(scrollChildren, cfg.Footer)
	}

	scrollProps := Props{
		ID: scrollID,
		Style: Style{
			Width:     Fill(),
			Height:    Fill(),
			Direction: Column,
			Gap:       cfg.Gap,
		},
		Scroll: ScrollState{
			OffsetY: cfg.ScrollOffsetY,
		},
	}
	if cfg.OnScrollChange != nil {
		scrollProps.Handlers.OnScroll = func(ctx EventContext) {
			maxOffset := maxFloat(0, ctx.Layout.ContentHeight-ctx.Layout.Frame.Height)
			nextOffset := clampFloat(cfg.ScrollOffsetY-(ctx.ScrollY*defaultPageScrollStep), 0, maxOffset)
			if nextOffset == cfg.ScrollOffsetY {
				return
			}
			cfg.OnScrollChange(nextOffset)
		}
	}

	children = append(children, ScrollView(scrollProps, scrollChildren...))

	return Div(Props{
		ID: rootID,
		Style: Style{
			Width:           Fill(),
			Height:          Fill(),
			Direction:       Column,
			Padding:         All(cfg.Padding),
			Gap:             cfg.Gap,
			BackgroundColor: cfg.BackgroundColor,
		},
	}, children...)
}

func clampFloat(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
