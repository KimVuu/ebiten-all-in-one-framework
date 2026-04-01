package main

import (
	"fmt"
	"image/color"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

type showcaseLayoutState struct {
	PageScroll float64
}

func buildShowcaseDOM() *ebitenui.DOM {
	return buildShowcaseDOMWithState(showcaseLayoutState{}, nil, nil)
}

func buildShowcaseDOMWithState(state showcaseLayoutState, onPageScrollChange func(float64), runtime *ebitenui.Runtime) *ebitenui.DOM {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}

	dom := ebitenui.New(
		ebitenui.PageLayout(ebitenui.PageLayoutConfig{
			ID:              "showcase-root",
			ScrollID:        "showcase-scroll",
			Padding:         24,
			Gap:             16,
			BackgroundColor: color.RGBA{R: 13, G: 18, B: 27, A: 255},
			ScrollOffsetY:   state.PageScroll,
			OnScrollChange:  onPageScrollChange,
			Header: ebitenui.Header(ebitenui.Props{
				ID: "showcase-header",
				Style: ebitenui.Style{
					Direction:       ebitenui.Row,
					Padding:         ebitenui.All(16),
					Gap:             12,
					BackgroundColor: panel,
					BorderColor:     stroke,
					BorderWidth:     1,
				},
			},
				ebitenui.Div(ebitenui.Props{
					ID: "header-copy",
					Style: ebitenui.Style{
						Width:     ebitenui.Fill(),
						Direction: ebitenui.Column,
						Gap:       6,
					},
				},
					ebitenui.Text("ebiten ebiten-ui showcase", ebitenui.Props{
						ID:    "hero-title",
						Style: ebitenui.Style{Color: textStrong},
					}),
					ebitenui.Text("HTML-like tags, stateful components, and prefab UI in one Ebiten showcase.", ebitenui.Props{
						ID:    "hero-subtitle",
						Style: ebitenui.Style{Color: textMuted},
					}),
				),
				ebitenui.Span(ebitenui.Props{
					ID: "header-badge",
					Style: ebitenui.Style{
						Padding:         ebitenui.All(10),
						BackgroundColor: accent,
					},
				},
					ebitenui.Text("single example", ebitenui.Props{
						ID:    "badge-text",
						Style: ebitenui.Style{Color: color.RGBA{R: 8, G: 12, B: 20, A: 255}},
					}),
				),
			),
			Content: ebitenui.Main(ebitenui.Props{
				ID: "showcase-main",
				Style: ebitenui.Style{
					Width:     ebitenui.Fill(),
					Direction: ebitenui.Column,
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
			Footer: ebitenui.Footer(ebitenui.Props{
				ID: "showcase-footer",
				Style: ebitenui.Style{
					Direction:       ebitenui.Row,
					Padding:         ebitenui.All(14),
					Gap:             12,
					BackgroundColor: panel,
					BorderColor:     stroke,
					BorderWidth:     1,
				},
			},
				ebitenui.Text("Footer", ebitenui.Props{
					ID:    "footer-title",
					Style: ebitenui.Style{Color: textStrong},
				}),
				ebitenui.Text(fmt.Sprintf("Unified showcase for %d core tags, components, and prefabs.", 13), ebitenui.Props{
					ID:    "footer-copy",
					Style: ebitenui.Style{Color: textMuted},
				}),
			),
		}),
	)

	applyShowcaseMetadata(dom)
	applyShowcaseRuntimeValues(dom, runtime)
	return dom
}

func applyShowcaseMetadata(dom *ebitenui.DOM) {
	if dom == nil || dom.Root == nil {
		return
	}

	annotate := func(id string, semantic ebitenui.SemanticSpec, layout ebitenui.LayoutSpec) {
		node, ok := dom.FindByID(id)
		if !ok || node == nil {
			return
		}
		node.Props.Semantic = semantic
		if layout != (ebitenui.LayoutSpec{}) {
			node.Props.Layout = mergeLayoutSpec(node.Props.Layout, layout)
		}
	}

	annotate("showcase-root", ebitenui.SemanticSpec{
		Screen:  "ebiten-ui-showcase",
		Element: "showcase-root",
		Role:    "screen",
		Slot:    "root",
	}, ebitenui.LayoutSpec{
		Mode:    ebitenui.LayoutModeFlowVertical,
		Size:    ebitenui.LayoutSize{Width: ebitenui.Fill(), Height: ebitenui.Fill()},
		Padding: ebitenui.All(24),
		Gap:     16,
	})
	annotate("showcase-header", ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: "showcase-header", Role: "header", Slot: "hero"}, ebitenui.LayoutSpec{Mode: ebitenui.LayoutModeFlowHorizontal, ParentID: "showcase-root"})
	annotate("showcase-scroll", ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: "showcase-scroll", Role: "scroll", Slot: "page"}, ebitenui.LayoutSpec{
		Mode:     ebitenui.LayoutModeFlowVertical,
		ParentID: "showcase-root",
		Size:     ebitenui.LayoutSize{Width: ebitenui.Fill(), Height: ebitenui.Fill()},
		Constraints: ebitenui.LayoutConstraints{
			ClipChildren:     true,
			KeepInsideParent: true,
		},
	})
	annotate("showcase-main", ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: "showcase-main", Role: "main", Slot: "content"}, ebitenui.LayoutSpec{Mode: ebitenui.LayoutModeFlowVertical, ParentID: "showcase-scroll", Size: ebitenui.LayoutSize{Width: ebitenui.Fill()}})
	annotate("showcase-footer", ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: "showcase-footer", Role: "footer", Slot: "footer"}, ebitenui.LayoutSpec{Mode: ebitenui.LayoutModeFlowHorizontal, ParentID: "showcase-scroll"})

	sectionIDs := []string{
		"overview-section", "button-section", "foundation-section", "dom-list-section", "components-section", "prefabs-section",
		"form-section", "layout-section", "overlay-section", "data-section", "status-section",
	}
	for _, id := range sectionIDs {
		annotate(id, ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: id, Role: "section", Slot: "content"}, ebitenui.LayoutSpec{})
	}

	controlIDs := []string{
		"name-input", "resolution-dropdown", "bio-textarea", "mode-radio", "party-stepper", "difficulty-toggle", "music-slider",
		"inventory-scrollbar", "dialog-demo", "hud-demo", "inventory-demo", "pause-demo", "settings-demo", "tooltip-demo",
	}
	for _, id := range controlIDs {
		annotate(id, ebitenui.SemanticSpec{Screen: "ebiten-ui-showcase", Element: id, Role: "control", Slot: "interactive"}, ebitenui.LayoutSpec{})
	}
}

func applyShowcaseRuntimeValues(dom *ebitenui.DOM, runtime *ebitenui.Runtime) {
	if dom == nil || runtime == nil {
		return
	}

	setNodeText(dom, "name-input-value", runtime.TextValueOrDefault("name-input", "Kim"))
	setNodeText(dom, "bio-textarea-value", runtime.TextValueOrDefault("bio-textarea", "Explorer of the ember valley.\nSpecializes in bows and traps."))
	setNodeText(dom, "resolution-dropdown-value", resolutionLabel(runtime.TextValueOrDefault("resolution-dropdown-selected", "resolution-720")))
}

func setNodeText(dom *ebitenui.DOM, id string, value string) {
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

func mergeLayoutSpec(base, overlay ebitenui.LayoutSpec) ebitenui.LayoutSpec {
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
	if overlay.Offset != (ebitenui.Point{}) {
		base.Offset = overlay.Offset
	}
	if overlay.Size != (ebitenui.LayoutSize{}) {
		base.Size = overlay.Size
	}
	if overlay.MinSize != (ebitenui.LayoutSize{}) {
		base.MinSize = overlay.MinSize
	}
	if overlay.MaxSize != (ebitenui.LayoutSize{}) {
		base.MaxSize = overlay.MaxSize
	}
	if overlay.Margin != (ebitenui.Insets{}) {
		base.Margin = overlay.Margin
	}
	if overlay.Padding != (ebitenui.Insets{}) {
		base.Padding = overlay.Padding
	}
	if overlay.Gap != 0 {
		base.Gap = overlay.Gap
	}
	if overlay.ZIndex != 0 {
		base.ZIndex = overlay.ZIndex
	}
	if overlay.Constraints != (ebitenui.LayoutConstraints{}) {
		base.Constraints = overlay.Constraints
	}
	if overlay.Grid != (ebitenui.LayoutGrid{}) {
		base.Grid = overlay.Grid
	}
	return base
}
