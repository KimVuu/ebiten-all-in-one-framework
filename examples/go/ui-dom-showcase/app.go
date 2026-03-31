package main

import (
	"fmt"
	"image/color"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

type showcaseLayoutState struct {
	PageScroll float64
}

func buildShowcaseDOM() *uidom.DOM {
	return buildShowcaseDOMWithState(showcaseLayoutState{}, nil, nil)
}

func buildShowcaseDOMWithState(state showcaseLayoutState, onPageScrollChange func(float64), runtime *uidom.Runtime) *uidom.DOM {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}

	dom := uidom.New(
		uidom.PageLayout(uidom.PageLayoutConfig{
			ID:              "showcase-root",
			ScrollID:        "showcase-scroll",
			Padding:         24,
			Gap:             16,
			BackgroundColor: color.RGBA{R: 13, G: 18, B: 27, A: 255},
			ScrollOffsetY:   state.PageScroll,
			OnScrollChange:  onPageScrollChange,
			Header: uidom.Header(uidom.Props{
				ID: "showcase-header",
				Style: uidom.Style{
					Direction:       uidom.Row,
					Padding:         uidom.All(16),
					Gap:             12,
					BackgroundColor: panel,
					BorderColor:     stroke,
					BorderWidth:     1,
				},
			},
				uidom.Div(uidom.Props{
					ID: "header-copy",
					Style: uidom.Style{
						Width:     uidom.Fill(),
						Direction: uidom.Column,
						Gap:       6,
					},
				},
					uidom.Text("ebiten ui-dom showcase", uidom.Props{
						ID:    "hero-title",
						Style: uidom.Style{Color: textStrong},
					}),
					uidom.Text("HTML-like tags, stateful components, and prefab UI in one Ebiten showcase.", uidom.Props{
						ID:    "hero-subtitle",
						Style: uidom.Style{Color: textMuted},
					}),
				),
				uidom.Span(uidom.Props{
					ID: "header-badge",
					Style: uidom.Style{
						Padding:         uidom.All(10),
						BackgroundColor: accent,
					},
				},
					uidom.Text("single example", uidom.Props{
						ID:    "badge-text",
						Style: uidom.Style{Color: color.RGBA{R: 8, G: 12, B: 20, A: 255}},
					}),
				),
			),
			Content: uidom.Main(uidom.Props{
				ID: "showcase-main",
				Style: uidom.Style{
					Width:     uidom.Fill(),
					Direction: uidom.Column,
					Gap:       16,
				},
			},
				overviewSection(),
				buttonSection(),
				foundationSection(),
				domListSection(),
				componentsSection(),
				prefabsSection(),
			),
			Footer: uidom.Footer(uidom.Props{
				ID: "showcase-footer",
				Style: uidom.Style{
					Direction:       uidom.Row,
					Padding:         uidom.All(14),
					Gap:             12,
					BackgroundColor: panel,
					BorderColor:     stroke,
					BorderWidth:     1,
				},
			},
				uidom.Text("Footer", uidom.Props{
					ID:    "footer-title",
					Style: uidom.Style{Color: textStrong},
				}),
				uidom.Text(fmt.Sprintf("Unified showcase for %d core tags, components, and prefabs.", 13), uidom.Props{
					ID:    "footer-copy",
					Style: uidom.Style{Color: textMuted},
				}),
			),
		}),
	)

	applyShowcaseMetadata(dom)
	applyShowcaseRuntimeValues(dom, runtime)
	return dom
}

func applyShowcaseMetadata(dom *uidom.DOM) {
	if dom == nil || dom.Root == nil {
		return
	}

	annotate := func(id string, semantic uidom.SemanticSpec, layout uidom.LayoutSpec) {
		node, ok := dom.FindByID(id)
		if !ok || node == nil {
			return
		}
		node.Props.Semantic = semantic
		if layout != (uidom.LayoutSpec{}) {
			node.Props.Layout = mergeLayoutSpec(node.Props.Layout, layout)
		}
	}

	annotate("showcase-root", uidom.SemanticSpec{
		Screen:  "ui-dom-showcase",
		Element: "showcase-root",
		Role:    "screen",
		Slot:    "root",
	}, uidom.LayoutSpec{
		Mode:    uidom.LayoutModeFlowVertical,
		Size:    uidom.LayoutSize{Width: uidom.Fill(), Height: uidom.Fill()},
		Padding: uidom.All(24),
		Gap:     16,
	})
	annotate("showcase-header", uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: "showcase-header", Role: "header", Slot: "hero"}, uidom.LayoutSpec{Mode: uidom.LayoutModeFlowHorizontal, ParentID: "showcase-root"})
	annotate("showcase-scroll", uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: "showcase-scroll", Role: "scroll", Slot: "page"}, uidom.LayoutSpec{
		Mode:     uidom.LayoutModeFlowVertical,
		ParentID: "showcase-root",
		Size:     uidom.LayoutSize{Width: uidom.Fill(), Height: uidom.Fill()},
		Constraints: uidom.LayoutConstraints{
			ClipChildren:     true,
			KeepInsideParent: true,
		},
	})
	annotate("showcase-main", uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: "showcase-main", Role: "main", Slot: "content"}, uidom.LayoutSpec{Mode: uidom.LayoutModeFlowVertical, ParentID: "showcase-scroll", Size: uidom.LayoutSize{Width: uidom.Fill()}})
	annotate("showcase-footer", uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: "showcase-footer", Role: "footer", Slot: "footer"}, uidom.LayoutSpec{Mode: uidom.LayoutModeFlowHorizontal, ParentID: "showcase-scroll"})

	sectionIDs := []string{
		"overview-section", "button-section", "foundation-section", "dom-list-section", "components-section", "prefabs-section",
		"form-section", "layout-section", "overlay-section", "data-section", "status-section",
	}
	for _, id := range sectionIDs {
		annotate(id, uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: id, Role: "section", Slot: "content"}, uidom.LayoutSpec{})
	}

	controlIDs := []string{
		"name-input", "resolution-dropdown", "bio-textarea", "mode-radio", "party-stepper", "difficulty-toggle", "music-slider",
		"inventory-scrollbar", "dialog-demo", "hud-demo", "inventory-demo", "pause-demo", "settings-demo", "tooltip-demo",
	}
	for _, id := range controlIDs {
		annotate(id, uidom.SemanticSpec{Screen: "ui-dom-showcase", Element: id, Role: "control", Slot: "interactive"}, uidom.LayoutSpec{})
	}
}

func applyShowcaseRuntimeValues(dom *uidom.DOM, runtime *uidom.Runtime) {
	if dom == nil || runtime == nil {
		return
	}

	setNodeText(dom, "name-input-value", runtime.TextValueOrDefault("name-input", "Kim"))
	setNodeText(dom, "bio-textarea-value", runtime.TextValueOrDefault("bio-textarea", "Explorer of the ember valley.\nSpecializes in bows and traps."))
	setNodeText(dom, "resolution-dropdown-value", resolutionLabel(runtime.TextValueOrDefault("resolution-dropdown-selected", "resolution-720")))
}

func setNodeText(dom *uidom.DOM, id string, value string) {
	node, ok := dom.FindByID(id)
	if !ok || node == nil {
		return
	}
	node.Text = value
}

func resolutionLabel(id string) string {
	switch id {
	case "resolution-1080":
		return "1920x1080"
	case "resolution-720":
		return "1280x720"
	default:
		if id == "" {
			return "1280x720"
		}
		return id
	}
}

func mergeLayoutSpec(base, overlay uidom.LayoutSpec) uidom.LayoutSpec {
	if overlay.Mode != "" {
		base.Mode = overlay.Mode
	}
	if overlay.ParentID != "" {
		base.ParentID = overlay.ParentID
	}
	if overlay.Anchor != "" {
		base.Anchor = overlay.Anchor
	}
	if overlay.Pivot != "" {
		base.Pivot = overlay.Pivot
	}
	if overlay.Offset != (uidom.Point{}) {
		base.Offset = overlay.Offset
	}
	if overlay.Size != (uidom.LayoutSize{}) {
		base.Size = overlay.Size
	}
	if overlay.MinSize != (uidom.LayoutSize{}) {
		base.MinSize = overlay.MinSize
	}
	if overlay.MaxSize != (uidom.LayoutSize{}) {
		base.MaxSize = overlay.MaxSize
	}
	if overlay.Margin != (uidom.Insets{}) {
		base.Margin = overlay.Margin
	}
	if overlay.Padding != (uidom.Insets{}) {
		base.Padding = overlay.Padding
	}
	if overlay.Gap != 0 {
		base.Gap = overlay.Gap
	}
	if overlay.ZIndex != 0 {
		base.ZIndex = overlay.ZIndex
	}
	if overlay.Constraints != (uidom.LayoutConstraints{}) {
		base.Constraints = overlay.Constraints
	}
	if overlay.Grid != (uidom.LayoutGrid{}) {
		base.Grid = overlay.Grid
	}
	return base
}
