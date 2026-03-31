package main

import (
	"image/color"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func overviewSection() *uidom.Node {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}
	accent3 := color.RGBA{R: 255, G: 180, B: 72, A: 255}

	return uidom.Section(uidom.Props{
		ID: "overview-section",
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(16),
			Gap:             12,
			BackgroundColor: panel,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		uidom.Text("Overview", uidom.Props{
			ID:    "overview-title",
			Style: uidom.Style{Color: textStrong},
		}),
		uidom.Div(uidom.Props{
			ID: "overview-cards",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       12,
			},
		},
			infoCard("layout-card", accent, textStrong, textMuted, "Layout", "Row, column, fill, padding, gap"),
			infoCard("tags-card", accent2, textStrong, textMuted, "Tags", "div, header, main, section, footer, button, span, text"),
			infoCard("renderer-card", accent3, textStrong, textMuted, "Renderer", "Ebiten image rendering with DOM tree layout"),
		),
	)
}

func buttonSection() *uidom.Node {
	panelAlt := color.RGBA{R: 30, G: 40, B: 56, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}
	accent3 := color.RGBA{R: 255, G: 180, B: 72, A: 255}

	return uidom.Section(uidom.Props{
		ID: "button-section",
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(16),
			Gap:             12,
			BackgroundColor: panelAlt,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		uidom.Text("Interactive-looking Buttons", uidom.Props{
			ID:    "buttons-title",
			Style: uidom.Style{Color: textStrong},
		}),
		uidom.Div(uidom.Props{
			ID: "button-row",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       12,
			},
		},
			actionButton("action-primary", accent, color.RGBA{R: 8, G: 12, B: 20, A: 255}, "Primary", "Launch scene", uidom.InteractionState{Hovered: true}),
			actionButton("action-secondary", accent2, color.RGBA{R: 8, G: 16, B: 12, A: 255}, "Secondary", "Open panel", uidom.InteractionState{Focused: true}),
			actionButton("action-warning", accent3, color.RGBA{R: 24, G: 14, B: 6, A: 255}, "Warning", "Reset state", uidom.InteractionState{Pressed: true}),
		),
	)
}

func foundationSection() *uidom.Node {
	panelAlt := color.RGBA{R: 30, G: 40, B: 56, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}
	accent2 := color.RGBA{R: 72, G: 211, B: 161, A: 255}

	return uidom.Section(uidom.Props{
		ID: "foundation-section",
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(16),
			Gap:             12,
			BackgroundColor: panelAlt,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		uidom.Text("New foundation nodes", uidom.Props{
			ID:    "foundation-title",
			Style: uidom.Style{Color: textStrong},
		}),
		uidom.Div(uidom.Props{
			ID: "foundation-row",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       12,
			},
		},
			previewImageCard("preview-image", accent),
			stackPreview("stack-preview", accent2),
		),
		uidom.TextBlock("TextBlock can wrap a longer UI description inside a constrained width while still staying in the DOM tree like any other node.", uidom.Props{
			ID: "foundation-copy",
			Style: uidom.Style{
				Width:      uidom.Px(360),
				Color:      textMuted,
				LineHeight: 16,
			},
		}),
		uidom.Div(uidom.Props{
			ID: "spacer-row",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
				Gap:       8,
			},
		},
			uidom.Text("Left", uidom.Props{
				ID:    "spacer-left",
				Style: uidom.Style{Color: textStrong},
			}),
			uidom.Spacer(uidom.Props{
				ID: "spacer-demo",
				Style: uidom.Style{
					Width:  uidom.Fill(),
					Height: uidom.Px(1),
				},
			}),
			uidom.Text("Right", uidom.Props{
				ID:    "spacer-right",
				Style: uidom.Style{Color: textStrong},
			}),
		),
		scrollPreview("scroll-preview"),
	)
}

func domListSection() *uidom.Node {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}

	return uidom.Section(uidom.Props{
		ID: "dom-list-section",
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(16),
			Gap:             10,
			BackgroundColor: panel,
			BorderColor:     stroke,
			BorderWidth:     1,
		},
	},
		uidom.Text("Supported DOM tags in this library", uidom.Props{
			ID:    "dom-list-title",
			Style: uidom.Style{Color: textStrong},
		}),
		uidom.Div(uidom.Props{
			ID: "dom-list-row",
			Style: uidom.Style{
				Width:     uidom.Fill(),
				Direction: uidom.Row,
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

func infoCard(id string, accent color.Color, titleColor color.Color, bodyColor color.Color, title string, body string) *uidom.Node {
	return uidom.Div(uidom.Props{
		ID: id,
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(12),
			Gap:             8,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		uidom.Span(uidom.Props{
			ID: id + "-eyebrow",
			Style: uidom.Style{
				Padding:         uidom.All(8),
				BackgroundColor: accent,
			},
		},
			uidom.Text(title, uidom.Props{
				ID:    id + "-title",
				Style: uidom.Style{Color: color.RGBA{R: 10, G: 14, B: 22, A: 255}},
			}),
		),
		uidom.Text(body, uidom.Props{
			ID:    id + "-body",
			Style: uidom.Style{Color: bodyColor},
		}),
	)
}

func actionButton(id string, bg color.Color, fg color.Color, label string, caption string, state uidom.InteractionState) *uidom.Node {
	return uidom.InteractiveButton(uidom.Props{
		ID:    id,
		State: state,
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(12),
			Gap:             6,
			BackgroundColor: bg,
			BorderColor:     color.RGBA{R: 248, G: 248, B: 248, A: 100},
			BorderWidth:     1,
		},
	},
		uidom.Text(label, uidom.Props{
			ID:    id + "-label",
			Style: uidom.Style{Color: fg},
		}),
		uidom.Span(uidom.Props{
			ID: id + "-caption-wrap",
			Style: uidom.Style{
				Padding:         uidom.All(8),
				BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 48},
			},
		},
			uidom.Text(caption, uidom.Props{
				ID:    id + "-caption",
				Style: uidom.Style{Color: fg},
			}),
		),
	)
}

func previewImageCard(id string, accent color.Color) *uidom.Node {
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	return uidom.Div(uidom.Props{
		ID: id,
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Direction:       uidom.Column,
			Padding:         uidom.All(12),
			Gap:             8,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		uidom.Image(uidom.Props{
			ID:    id + "-image",
			Image: uidom.SolidImage(72, 54, accent),
		}),
		uidom.Text("Image", uidom.Props{
			ID:    id + "-title",
			Style: uidom.Style{Color: textStrong},
		}),
	)
}

func stackPreview(id string, accent color.Color) *uidom.Node {
	return uidom.Stack(uidom.Props{
		ID: id,
		Style: uidom.Style{
			Width:           uidom.Fill(),
			Height:          uidom.Px(96),
			Padding:         uidom.All(10),
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     accent,
			BorderWidth:     1,
		},
	},
		uidom.Image(uidom.Props{
			ID:    id + "-background",
			Image: uidom.SolidImage(120, 70, color.RGBA{R: 40, G: 60, B: 92, A: 255}),
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Fill(),
			},
		}),
		uidom.Span(uidom.Props{
			ID: id + "-badge",
			Style: uidom.Style{
				Padding:         uidom.All(8),
				BackgroundColor: accent,
			},
		},
			uidom.Text("Stack", uidom.Props{
				ID:    id + "-badge-text",
				Style: uidom.Style{Color: color.RGBA{R: 12, G: 18, B: 28, A: 255}},
			}),
		),
	)
}

func scrollPreview(id string) *uidom.Node {
	borderMuted := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	return uidom.ScrollView(uidom.Props{
		ID: id,
		Style: uidom.Style{
			Width:           uidom.Px(360),
			Height:          uidom.Px(84),
			Direction:       uidom.Column,
			Padding:         uidom.All(10),
			Gap:             6,
			BackgroundColor: color.RGBA{R: 19, G: 25, B: 35, A: 255},
			BorderColor:     borderMuted,
			BorderWidth:     1,
		},
		Scroll: uidom.ScrollState{
			OffsetY: 16,
		},
	},
		uidom.Text("Scroll row 01", uidom.Props{ID: id + "-1", Style: uidom.Style{Color: textStrong}}),
		uidom.Text("Scroll row 02", uidom.Props{ID: id + "-2", Style: uidom.Style{Color: textMuted}}),
		uidom.Text("Scroll row 03", uidom.Props{ID: id + "-3", Style: uidom.Style{Color: textMuted}}),
		uidom.Text("Scroll row 04", uidom.Props{ID: id + "-4", Style: uidom.Style{Color: textMuted}}),
	)
}

func tagChip(id string, label string) *uidom.Node {
	return uidom.Span(uidom.Props{
		ID: id,
		Style: uidom.Style{
			Padding:         uidom.All(10),
			BackgroundColor: color.RGBA{R: 46, G: 58, B: 78, A: 255},
			BorderColor:     color.RGBA{R: 96, G: 120, B: 156, A: 255},
			BorderWidth:     1,
		},
	},
		uidom.Text(label, uidom.Props{
			ID: id + "-text",
			Style: uidom.Style{
				Color: color.RGBA{R: 234, G: 239, B: 248, A: 255},
			},
		}),
	)
}
