package main

import (
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func overviewSection() *ebitenui.Node {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}
	accent3 := color.RGBA{R: 255, G: 180, B: 72, A: 255}

	return ebitenui.Section(ebitenui.Props{
		ID: "overview-section",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(16),
			Gap:             12,
			BackgroundColor: panel,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		ebitenui.Text("Overview", ebitenui.Props{
			ID:    "overview-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "overview-cards",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       12,
			},
		},
			infoCard("layout-card", accent, textStrong, textMuted, "Layout", "Row, column, fill, padding, gap"),
			infoCard("tags-card", accent2, textStrong, textMuted, "Tags", "div, header, main, section, footer, button, span, text"),
			infoCard("renderer-card", accent3, textStrong, textMuted, "Renderer", "Ebiten image rendering with DOM tree layout"),
		),
	)
}

func buttonSection() *ebitenui.Node {
	panelAlt := color.RGBA{R: 30, G: 40, B: 56, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}
	accent3 := color.RGBA{R: 255, G: 180, B: 72, A: 255}

	return ebitenui.Section(ebitenui.Props{
		ID: "button-section",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(16),
			Gap:             12,
			BackgroundColor: panelAlt,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		ebitenui.Text("Interactive-looking Buttons", ebitenui.Props{
			ID:    "buttons-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "button-row",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       12,
			},
		},
			actionButton("action-primary", accent, color.RGBA{R: 8, G: 12, B: 20, A: 255}, "Primary", "Launch scene", ebitenui.InteractionState{Hovered: true}),
			actionButton("action-secondary", accent2, color.RGBA{R: 8, G: 16, B: 12, A: 255}, "Secondary", "Open panel", ebitenui.InteractionState{Focused: true}),
			actionButton("action-warning", accent3, color.RGBA{R: 24, G: 14, B: 6, A: 255}, "Warning", "Reset state", ebitenui.InteractionState{Pressed: true}),
		),
	)
}

func foundationSection() *ebitenui.Node {
	panelAlt := color.RGBA{R: 30, G: 40, B: 56, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}

	return ebitenui.Section(ebitenui.Props{
		ID: "foundation-section",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(16),
			Gap:             12,
			BackgroundColor: panelAlt,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		ebitenui.Text("New foundation nodes", ebitenui.Props{
			ID:    "foundation-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "foundation-row",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       12,
			},
		},
			previewImageCard("preview-image", accent),
			stackPreview("stack-preview", accent2),
		),
		ebitenui.TextBlock("TextBlock can wrap a longer UI description inside a constrained width while still staying in the DOM tree like any other node.", ebitenui.Props{
			ID: "foundation-copy",
			Style: ebitenui.Style{
				Width:      ebitenui.Px(360),
				Color:      textMuted,
				LineHeight: 16,
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "spacer-row",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       8,
			},
		},
			ebitenui.Text("Left", ebitenui.Props{
				ID:    "spacer-left",
				Style: ebitenui.Style{Color: textStrong},
			}),
			ebitenui.Spacer(ebitenui.Props{
				ID: "spacer-demo",
				Style: ebitenui.Style{
					Width:  ebitenui.Fill(),
					Height: ebitenui.Px(1),
				},
			}),
			ebitenui.Text("Right", ebitenui.Props{
				ID:    "spacer-right",
				Style: ebitenui.Style{Color: textStrong},
			}),
		),
		scrollPreview("scroll-preview"),
	)
}

func domListSection() *ebitenui.Node {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}

	return ebitenui.Section(ebitenui.Props{
		ID: "dom-list-section",
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(16),
			Gap:             10,
			BackgroundColor: panel,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		ebitenui.Text("Supported DOM tags in this library", ebitenui.Props{
			ID:    "dom-list-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "dom-list-row",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       8,
			},
		},
			tagChip("chip-div", "div"),
			tagChip("chip-header", "header"),
			tagChip("chip-main", "main"),
			tagChip("chip-section", "section"),
			tagChip("chip-footer", "footer"),
			tagChip("chip-button", "button"),
			tagChip("chip-span", "span"),
			tagChip("chip-text", "text"),
			tagChip("chip-image", "img"),
			tagChip("chip-text-block", "text-block"),
			tagChip("chip-spacer", "spacer"),
			tagChip("chip-stack", "stack"),
			tagChip("chip-scroll", "scroll-view"),
		),
	)
}

func infoCard(id string, accent color.Color, titleColor color.Color, bodyColor color.Color, title string, body string) *ebitenui.Node {
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(12),
			Gap:             8,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		ebitenui.Span(ebitenui.Props{
			ID: id + "-eyebrow",
			Style: ebitenui.Style{
				Padding:         ebitenui.All(8),
				BackgroundColor: accent,
			},
		},
			ebitenui.Text(title, ebitenui.Props{
				ID:    id + "-title",
				Style: ebitenui.Style{Color: color.RGBA{R: 10, G: 14, B: 22, A: 255}},
			}),
		),
		ebitenui.Text(body, ebitenui.Props{
			ID:    id + "-body",
			Style: ebitenui.Style{Color: bodyColor},
		}),
	)
}

func actionButton(id string, bg color.Color, fg color.Color, label string, caption string, state ebitenui.InteractionState) *ebitenui.Node {
	return ebitenui.InteractiveButton(ebitenui.Props{
		ID:    id,
		State: state,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(12),
			Gap:             6,
			BackgroundColor: bg,
			BorderColor:     color.RGBA{R: 248, G: 248, B: 248, A: 100},
			BorderWidth:     1,
		},
	},
		ebitenui.Text(label, ebitenui.Props{
			ID:    id + "-label",
			Style: ebitenui.Style{Color: fg},
		}),
		ebitenui.Span(ebitenui.Props{
			ID: id + "-caption-wrap",
			Style: ebitenui.Style{
				Padding:         ebitenui.All(8),
				BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 48},
			},
		},
			ebitenui.Text(caption, ebitenui.Props{
				ID:    id + "-caption",
				Style: ebitenui.Style{Color: fg},
			}),
		),
	)
}

func previewImageCard(id string, accent color.Color) *ebitenui.Node {
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(12),
			Gap:             8,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		ebitenui.Image(ebitenui.Props{
			ID:    id + "-image",
			Image: ebitenui.SolidImage(72, 54, accent),
		}),
		ebitenui.Text("Image", ebitenui.Props{
			ID:    id + "-title",
			Style: ebitenui.Style{Color: textStrong},
		}),
	)
}

func stackPreview(id string, accent color.Color) *ebitenui.Node {
	return ebitenui.Stack(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Height:          ebitenui.Px(96),
			Padding:         ebitenui.All(10),
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		ebitenui.Image(ebitenui.Props{
			ID:    id + "-background",
			Image: ebitenui.SolidImage(120, 70, color.RGBA{R: 40, G: 60, B: 92, A: 255}),
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Span(ebitenui.Props{
			ID: id + "-badge",
			Style: ebitenui.Style{
				Padding:         ebitenui.All(8),
				BackgroundColor: accent,
			},
		},
			ebitenui.Text("Stack", ebitenui.Props{
				ID:    id + "-badge-text",
				Style: ebitenui.Style{Color: color.RGBA{R: 12, G: 18, B: 28, A: 255}},
			}),
		),
	)
}

func scrollPreview(id string) *ebitenui.Node {
	borderMuted := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	return ebitenui.ScrollView(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Px(360),
			Height:          ebitenui.Px(84),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(10),
			Gap:             6,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     borderMuted,
			BorderWidth:     1,
		},
		Scroll: ebitenui.ScrollState{
			OffsetY: 16,
		},
	},
		ebitenui.Text("Scroll row 01", ebitenui.Props{ID: id + "-1", Style: ebitenui.Style{Color: textStrong}}),
		ebitenui.Text("Scroll row 02", ebitenui.Props{ID: id + "-2", Style: ebitenui.Style{Color: textMuted}}),
		ebitenui.Text("Scroll row 03", ebitenui.Props{ID: id + "-3", Style: ebitenui.Style{Color: textMuted}}),
		ebitenui.Text("Scroll row 04", ebitenui.Props{ID: id + "-4", Style: ebitenui.Style{Color: textMuted}}),
	)
}

func tagChip(id string, label string) *ebitenui.Node {
	return ebitenui.Span(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Padding:         ebitenui.All(10),
			BackgroundColor: color.RGBA{R: 46, G: 58, B: 78, A: 255},
			BorderColor:     color.RGBA{R: 96, G: 120, B: 156, A: 255},
			BorderWidth:     1,
		},
	},
		ebitenui.Text(label, ebitenui.Props{
			ID: id + "-text",
			Style: ebitenui.Style{
				Color: color.RGBA{R: 234, G: 239, B: 248, A: 255},
			},
		}),
	)
}
