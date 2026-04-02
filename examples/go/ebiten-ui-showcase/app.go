package main

import (
	"image/color"
	"strings"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

const showcaseScrollStep = 48

func buildShowcaseDOM() *ebitenui.DOM {
	return buildShowcaseDOMWithState(showcaseLayoutState{}, nil, nil, newShowcaseBindings())
}

func buildShowcaseDOMWithState(state showcaseLayoutState, callbacks *showcaseCallbacks, runtime *ebitenui.Runtime, bindings *showcaseBindings) *ebitenui.DOM {
	if bindings == nil {
		bindings = newShowcaseBindings()
	}
	registry := buildShowcasePageRegistry()
	router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{
		Routes:        registry.Routes,
		InitialPageID: initialShowcasePageID(state.CurrentPageID),
	})
	currentPageID := router.CurrentPageID()
	currentPage, ok := registry.Pages[currentPageID]
	if !ok {
		currentPageID = "overview"
		_ = router.Navigate(currentPageID)
		currentPage = registry.Pages[currentPageID]
	}

	root := ebitenui.Div(ebitenui.Props{
		ID: "showcase-root",
		Semantic: ebitenui.SemanticSpec{
			Screen:  "ebiten-ui-showcase",
			Element: "showcase-root",
			Role:    "screen",
			Slot:    "root",
		},
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeFlowVertical,
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
			Padding: ebitenui.All(24),
			Gap:     16,
		},
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Height:          ebitenui.Fill(),
			Direction:       ebitenui.Column,
			Padding:         ebitenui.All(24),
			Gap:             16,
			BackgroundColor: color.RGBA{R: 13, G: 18, B: 27, A: 255},
		},
	},
		buildShowcaseHeader(currentPage),
		ebitenui.Main(ebitenui.Props{
			ID: "showcase-main",
			Semantic: ebitenui.SemanticSpec{
				Screen:  "ebiten-ui-showcase",
				Element: "showcase-main",
				Role:    "main",
				Slot:    "content",
			},
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeFlowHorizontal,
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Fill(),
					Height: ebitenui.Fill(),
				},
				Gap: 16,
			},
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Height:    ebitenui.Fill(),
				Direction: ebitenui.Column,
			},
		},
			ebitenui.PageScreen(ebitenui.PageScreenConfig{
				ID:              "showcase-page-screen",
				SidebarWidth:    320,
				Gap:             16,
				BackgroundColor: color.RGBA{R: 13, G: 18, B: 27, A: 255},
				Sidebar:         buildShowcaseSidebar(router, registry, currentPageID, state.SidebarScroll, callbacks),
				Content:         buildShowcaseDetail(router, registry, currentPage, currentPageID, state.DetailScroll, callbacks, runtime, bindings),
			}),
		),
	)

	return ebitenui.New(root)
}

func initialShowcasePageID(pageID string) string {
	if strings.TrimSpace(pageID) == "" {
		return "overview"
	}
	return pageID
}

func buildShowcaseHeader(page ShowcasePageSpec) *ebitenui.Node {
	panel := color.RGBA{R: 24, G: 31, B: 43, A: 255}
	stroke := color.RGBA{R: 88, G: 110, B: 140, A: 255}
	textStrong := color.RGBA{R: 242, G: 246, B: 252, A: 255}
	textMuted := color.RGBA{R: 178, G: 190, B: 207, A: 255}
	accent := color.RGBA{R: 80, G: 160, B: 255, A: 255}

	return ebitenui.Header(ebitenui.Props{
		ID: "showcase-header",
		Semantic: ebitenui.SemanticSpec{
			Screen:  "ebiten-ui-showcase",
			Element: "showcase-header",
			Role:    "header",
			Slot:    "hero",
		},
		Layout: ebitenui.LayoutSpec{
			Mode:     ebitenui.LayoutModeFlowHorizontal,
			ParentID: "showcase-root",
		},
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
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
			ebitenui.Text("ebiten-ui showcase", ebitenui.Props{
				ID:    "hero-title",
				Style: ebitenui.Style{Color: textStrong},
			}),
			ebitenui.Text("Page-based reference for components, prefabs, and UI patterns.", ebitenui.Props{
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
			ebitenui.Text(page.Title, ebitenui.Props{
				ID:    "badge-text",
				Style: ebitenui.Style{Color: color.RGBA{R: 8, G: 12, B: 20, A: 255}},
			}),
		),
	)
}

func buildShowcaseSidebar(router *ebitenui.PageRouter, registry ShowcasePageRegistry, currentPageID string, scrollOffset float64, callbacks *showcaseCallbacks) *ebitenui.Node {
	scrollHandlers := ebitenui.EventHandlers{}
	if callbacks != nil && callbacks.OnSidebarScrollChange != nil {
		scrollHandlers.OnScroll = func(ctx ebitenui.EventContext) {
			maxOffset := maxFloat64(0, ctx.Layout.ContentHeight-ctx.Layout.Frame.Height)
			nextOffset := clampFloat64(scrollOffset-(ctx.ScrollY*showcaseScrollStep), 0, maxOffset)
			if nextOffset != scrollOffset {
				callbacks.OnSidebarScrollChange(nextOffset)
			}
		}
	}

	children := []*ebitenui.Node{
		ebitenui.Text("Pages", ebitenui.Props{
			ID:    "showcase-sidebar-title",
			Style: showcaseGroupTitleStyle(),
		}),
		ebitenui.TextBlock("Browse groups on the left, then inspect usage, code, and live demos on the right.", ebitenui.Props{
			ID:    "showcase-sidebar-copy",
			Style: showcaseGroupCopyStyle(),
		}),
	}

	for _, route := range router.VisibleNavTree() {
		children = append(children, buildNavRoute(route, currentPageID, callbacks, 0))
	}

	return ebitenui.Div(ebitenui.Props{
		ID: "showcase-sidebar",
		Semantic: ebitenui.SemanticSpec{
			Screen:  "ebiten-ui-showcase",
			Element: "showcase-sidebar",
			Role:    "navigation",
			Slot:    "sidebar",
		},
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeFlowVertical,
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		},
		Style: showcaseGroupStyle(),
	},
		ebitenui.ScrollView(ebitenui.Props{
			ID: "showcase-sidebar-scroll",
			Semantic: ebitenui.SemanticSpec{
				Screen:  "ebiten-ui-showcase",
				Element: "showcase-sidebar-scroll",
				Role:    "scroll",
				Slot:    "sidebar-scroll",
			},
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeFlowVertical,
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Fill(),
					Height: ebitenui.Fill(),
				},
				Constraints: ebitenui.LayoutConstraints{
					ClipChildren:     true,
					KeepInsideParent: true,
				},
			},
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Height:    ebitenui.Fill(),
				Direction: ebitenui.Column,
				Gap:       10,
			},
			Scroll: ebitenui.ScrollState{
				OffsetY: scrollOffset,
			},
			Handlers: scrollHandlers,
		}, children...),
	)
}

func buildNavRoute(route ebitenui.PageRoute, currentPageID string, callbacks *showcaseCallbacks, depth int) *ebitenui.Node {
	active := route.ID == currentPageID
	expanded := routeContainsPage(route, currentPageID)
	paddingLeft := 12 + (depth * 18)

	buttonBackground := color.RGBA{R: 18, G: 24, B: 34, A: 255}
	buttonText := color.RGBA{R: 194, G: 204, B: 219, A: 255}
	if expanded {
		buttonBackground = color.RGBA{R: 26, G: 35, B: 49, A: 255}
	}
	if active {
		buttonBackground = color.RGBA{R: 80, G: 160, B: 255, A: 255}
		buttonText = color.RGBA{R: 12, G: 18, B: 28, A: 255}
	}

	item := ebitenui.InteractiveButton(ebitenui.Props{
		ID: "nav-item-" + sanitizeID(route.ID),
		Semantic: ebitenui.SemanticSpec{
			Screen:  "ebiten-ui-showcase",
			Element: route.ID,
			Role:    "nav-item",
			Slot:    "sidebar-item",
		},
		State: ebitenui.InteractionState{
			Selected: active,
		},
		Handlers: ebitenui.EventHandlers{
			OnClick: func(ctx ebitenui.EventContext) {
				if callbacks != nil && callbacks.OnNavigate != nil {
					callbacks.OnNavigate(route.ID)
				}
			},
		},
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Row,
			Padding:         ebitenui.Insets{Top: 10, Right: 12, Bottom: 10, Left: float64(paddingLeft)},
			Gap:             8,
			BackgroundColor: buttonBackground,
			BorderColor:     color.RGBA{R: 74, G: 92, B: 120, A: 255},
			BorderWidth:     1,
		},
	},
		ebitenui.Text(route.Title, ebitenui.Props{
			ID: "nav-item-label-" + sanitizeID(route.ID),
			Style: ebitenui.Style{
				Color: buttonText,
			},
		}),
	)

	if len(route.Children) == 0 {
		return item
	}

	children := make([]*ebitenui.Node, 0, len(route.Children)+1)
	children = append(children, item)
	for _, child := range route.Children {
		children = append(children, buildNavRoute(child, currentPageID, callbacks, depth+1))
	}

	return ebitenui.Div(ebitenui.Props{
		ID: "nav-group-" + sanitizeID(route.ID),
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Column,
			Gap:       8,
		},
	}, children...)
}

func buildShowcaseDetail(router *ebitenui.PageRouter, registry ShowcasePageRegistry, page ShowcasePageSpec, currentPageID string, scrollOffset float64, callbacks *showcaseCallbacks, runtime *ebitenui.Runtime, bindings *showcaseBindings) *ebitenui.Node {
	scrollHandlers := ebitenui.EventHandlers{}
	if callbacks != nil && callbacks.OnDetailScrollChange != nil {
		scrollHandlers.OnScroll = func(ctx ebitenui.EventContext) {
			maxOffset := maxFloat64(0, ctx.Layout.ContentHeight-ctx.Layout.Frame.Height)
			nextOffset := clampFloat64(scrollOffset-(ctx.ScrollY*showcaseScrollStep), 0, maxOffset)
			if nextOffset != scrollOffset {
				callbacks.OnDetailScrollChange(nextOffset)
			}
		}
	}

	ctx := ShowcaseDemoContext{
		Runtime:       runtime,
		Registry:      registry,
		CurrentPageID: currentPageID,
		Bindings:      bindings,
	}
	demo := page.DemoBuilder
	var demoNode *ebitenui.Node
	if demo != nil {
		demoNode = demo(ctx)
	}
	if demoNode == nil {
		demoNode = ebitenui.Text("No demo available for this page yet.", ebitenui.Props{
			ID:    "page-demo-empty",
			Style: showcaseGroupCopyStyle(),
		})
	}

	breadcrumbLabel := breadcrumbText(router.Breadcrumb(currentPageID))

	return ebitenui.Div(ebitenui.Props{
		ID: "showcase-detail",
		Semantic: ebitenui.SemanticSpec{
			Screen:  "ebiten-ui-showcase",
			Element: "showcase-detail",
			Role:    "content",
			Slot:    "detail",
		},
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeFlowVertical,
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		},
		Style: showcaseGroupStyle(),
	},
		ebitenui.ScrollView(ebitenui.Props{
			ID: "showcase-detail-scroll",
			Semantic: ebitenui.SemanticSpec{
				Screen:  "ebiten-ui-showcase",
				Element: "showcase-detail-scroll",
				Role:    "scroll",
				Slot:    "detail-scroll",
			},
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeFlowVertical,
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Fill(),
					Height: ebitenui.Fill(),
				},
				Constraints: ebitenui.LayoutConstraints{
					ClipChildren:     true,
					KeepInsideParent: true,
				},
			},
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Height:    ebitenui.Fill(),
				Direction: ebitenui.Column,
				Gap:       16,
			},
			Scroll:   ebitenui.ScrollState{OffsetY: scrollOffset},
			Handlers: scrollHandlers,
		},
			ebitenui.Div(ebitenui.Props{
				ID: "showcase-detail-body",
				Style: ebitenui.Style{
					Width:     ebitenui.Fill(),
					Direction: ebitenui.Column,
					Gap:       16,
				},
			},
				ebitenui.Div(ebitenui.Props{
					ID:    "page-summary",
					Style: detailSectionStyle(),
				},
					ebitenui.Text(breadcrumbLabel, ebitenui.Props{
						ID: "page-breadcrumb",
						Style: ebitenui.Style{
							Color: color.RGBA{R: 80, G: 160, B: 255, A: 255},
						},
					}),
					ebitenui.Text(page.Title, ebitenui.Props{
						ID:    "page-title",
						Style: detailTitleStyle(),
					}),
					ebitenui.TextBlock(page.Description, ebitenui.Props{
						ID:    "page-description",
						Style: showcaseGroupCopyStyle(),
					}),
				),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-demo",
					Style: detailSectionStyle(),
				},
					ebitenui.Text("Live Demo", ebitenui.Props{
						ID:    "page-demo-title",
						Style: detailTitleStyle(),
					}),
					demoNode,
				),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-usage",
					Style: detailSectionStyle(),
				},
					ebitenui.Text("How To Use", ebitenui.Props{
						ID:    "page-usage-title",
						Style: detailTitleStyle(),
					}),
					ebitenui.TextBlock(page.UsageNotes, ebitenui.Props{
						ID:    "page-usage-copy",
						Style: showcaseGroupCopyStyle(),
					}),
				),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-code",
					Style: detailSectionStyle(),
				},
					ebitenui.Text("Code Example", ebitenui.Props{
						ID:    "page-code-title",
						Style: detailTitleStyle(),
					}),
					ebitenui.TextBlock(page.CodeExample, ebitenui.Props{
						ID: "page-code-block",
						Style: ebitenui.Style{
							Width:           ebitenui.Fill(),
							Padding:         ebitenui.All(14),
							Color:           color.RGBA{R: 213, G: 223, B: 238, A: 255},
							LineHeight:      16,
							BackgroundColor: color.RGBA{R: 17, G: 22, B: 31, A: 255},
							BorderColor:     color.RGBA{R: 63, G: 78, B: 101, A: 255},
							BorderWidth:     1,
						},
					}),
				),
			),
		),
	)
}

func newShowcaseBindings() *showcaseBindings {
	return &showcaseBindings{
		NameInput:      ebitenui.NewRef("Kim"),
		Resolution:     ebitenui.NewRef("resolution-720"),
		ResolutionOpen: ebitenui.NewRef(true),
		Bio:            ebitenui.NewRef("Explorer of the ember valley."),
		Hardcore:       ebitenui.NewRef(true),
		MusicVolume:    ebitenui.NewRef(65.0),
	}
}

func breadcrumbText(routes []ebitenui.PageRoute) string {
	if len(routes) == 0 {
		return "Showcase"
	}
	parts := make([]string, 0, len(routes))
	for _, route := range routes {
		parts = append(parts, route.Title)
	}
	return strings.Join(parts, " / ")
}

func routeContainsPage(route ebitenui.PageRoute, currentPageID string) bool {
	if route.ID == currentPageID {
		return true
	}
	for _, child := range route.Children {
		if routeContainsPage(child, currentPageID) {
			return true
		}
	}
	return false
}

func clampFloat64(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func maxFloat64(left, right float64) float64 {
	if left > right {
		return left
	}
	return right
}
