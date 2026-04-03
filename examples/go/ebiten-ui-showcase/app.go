package main

import (
	"fmt"
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
	fontPreset := applyShowcaseFontPreset(state.FontPreset)
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

	preset := showcaseThemePresetByID(initialShowcaseThemePreset(state.ThemePreset))

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
			BackgroundColor: preset.Chrome.RootBackground,
		},
	},
		buildShowcaseHeader(currentPage, preset, fontPreset, callbacks),
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
				BackgroundColor: preset.Chrome.RootBackground,
				Sidebar:         buildShowcaseSidebar(router, currentPageID, state.SidebarScroll, callbacks, preset.Chrome),
				Content:         buildShowcaseDetail(router, registry, currentPage, currentPageID, state.DetailScroll, callbacks, runtime, bindings, preset, fontPreset),
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

func buildShowcaseHeader(page ShowcasePageSpec, preset showcaseThemePreset, fontPreset showcaseFontPreset, callbacks *showcaseCallbacks) *ebitenui.Node {
	chrome := preset.Chrome

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
			BackgroundColor: chrome.PanelBackground,
			BorderColor:     chrome.PanelBorder,
			BorderWidth:     1,
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "header-copy",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Column,
				Gap:       8,
			},
		},
			ebitenui.Text("ebiten-ui showcase", ebitenui.Props{
				ID:    "hero-title",
				Style: ebitenui.Style{Color: chrome.TextStrong},
			}),
			ebitenui.Text("Page-based reference for components, prefabs, theme presets, and reactive state patterns.", ebitenui.Props{
				ID:    "hero-subtitle",
				Style: ebitenui.Style{Color: chrome.TextMuted},
			}),
			buildThemePresetSwitcher(preset, callbacks),
			buildFontPresetSwitcher(fontPreset, preset, callbacks),
		),
		ebitenui.Span(ebitenui.Props{
			ID: "header-badge",
			Style: ebitenui.Style{
				Padding:         ebitenui.All(10),
				BackgroundColor: chrome.Accent,
			},
		},
			ebitenui.Text(page.Title, ebitenui.Props{
				ID:    "badge-text",
				Style: ebitenui.Style{Color: chrome.BadgeText},
			}),
		),
	)
}

func buildFontPresetSwitcher(fontPreset showcaseFontPreset, preset showcaseThemePreset, callbacks *showcaseCallbacks) *ebitenui.Node {
	chrome := preset.Chrome
	buttons := make([]*ebitenui.Node, 0, len(showcaseFontPresets()))
	for _, option := range showcaseFontPresets() {
		option := option
		active := option.ID == fontPreset.ID
		background := chrome.PanelBackground
		textColor := chrome.TextMuted
		border := chrome.PanelBorder
		if active {
			background = chrome.AccentSoft
			textColor = chrome.BadgeText
			border = chrome.Accent
		}
		buttons = append(buttons, ebitenui.InteractiveButton(ebitenui.Props{
			ID: "font-preset-" + option.ID,
			Semantic: ebitenui.SemanticSpec{
				Screen:  "ebiten-ui-showcase",
				Element: "font-preset-" + option.ID,
				Role:    "action",
				Slot:    "font-preset",
			},
			State: ebitenui.InteractionState{
				Selected: active,
			},
			Handlers: ebitenui.EventHandlers{
				OnClick: func(ctx ebitenui.EventContext) {
					if callbacks != nil && callbacks.OnFontPresetChange != nil {
						callbacks.OnFontPresetChange(option.ID)
					}
				},
			},
			Style: ebitenui.Style{
				Padding:         ebitenui.Insets{Top: 8, Right: 12, Bottom: 8, Left: 12},
				BackgroundColor: background,
				BorderColor:     border,
				BorderWidth:     1,
			},
		},
			ebitenui.Text(option.Title, ebitenui.Props{
				ID:    "font-preset-label-" + option.ID,
				Style: ebitenui.Style{Color: textColor},
			}),
		))
	}

	return ebitenui.Div(ebitenui.Props{
		ID: "font-preset-switcher",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Column,
			Gap:       8,
		},
	},
		ebitenui.Text("Font Presets", ebitenui.Props{
			ID:    "font-preset-title",
			Style: ebitenui.Style{Color: chrome.TextMuted},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "font-preset-buttons",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       10,
			},
		}, buttons...),
	)
}

func buildThemePresetSwitcher(preset showcaseThemePreset, callbacks *showcaseCallbacks) *ebitenui.Node {
	chrome := preset.Chrome
	buttons := make([]*ebitenui.Node, 0, len(showcaseThemePresets()))
	for _, option := range showcaseThemePresets() {
		option := option
		active := option.ID == preset.ID
		background := chrome.PanelBackground
		textColor := chrome.TextMuted
		border := chrome.PanelBorder
		if active {
			background = chrome.Accent
			textColor = chrome.BadgeText
			border = chrome.AccentSoft
		}
		buttons = append(buttons, ebitenui.InteractiveButton(ebitenui.Props{
			ID: "theme-preset-" + option.ID,
			Semantic: ebitenui.SemanticSpec{
				Screen:  "ebiten-ui-showcase",
				Element: "theme-preset-" + option.ID,
				Role:    "action",
				Slot:    "theme-preset",
			},
			State: ebitenui.InteractionState{
				Selected: active,
			},
			Handlers: ebitenui.EventHandlers{
				OnClick: func(ctx ebitenui.EventContext) {
					if callbacks != nil && callbacks.OnThemePresetChange != nil {
						callbacks.OnThemePresetChange(option.ID)
					}
				},
			},
			Style: ebitenui.Style{
				Padding:         ebitenui.Insets{Top: 8, Right: 12, Bottom: 8, Left: 12},
				BackgroundColor: background,
				BorderColor:     border,
				BorderWidth:     1,
			},
		},
			ebitenui.Text(option.Title, ebitenui.Props{
				ID:    "theme-preset-label-" + option.ID,
				Style: ebitenui.Style{Color: textColor},
			}),
		))
	}

	return ebitenui.Div(ebitenui.Props{
		ID: "theme-preset-switcher",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Direction: ebitenui.Column,
			Gap:       8,
		},
	},
		ebitenui.Text("Theme Presets", ebitenui.Props{
			ID:    "theme-preset-title",
			Style: ebitenui.Style{Color: chrome.TextMuted},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "theme-preset-buttons",
			Style: ebitenui.Style{
				Width:     ebitenui.Fill(),
				Direction: ebitenui.Row,
				Gap:       10,
			},
		}, buttons...),
	)
}

func buildShowcaseSidebar(router *ebitenui.PageRouter, currentPageID string, scrollOffset float64, callbacks *showcaseCallbacks, chrome showcaseChrome) *ebitenui.Node {
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
			Style: showcaseGroupTitleStyleForChrome(chrome),
		}),
		ebitenui.TextBlock("Browse groups on the left, then inspect usage, code, live demo state, and theme-aware variants on the right.", ebitenui.Props{
			ID:    "showcase-sidebar-copy",
			Style: showcaseGroupCopyStyleForChrome(chrome),
		}),
	}

	for _, route := range router.VisibleNavTree() {
		children = append(children, buildNavRoute(route, currentPageID, callbacks, 0, chrome))
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
		Style: showcaseGroupStyleForChrome(chrome),
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
			Scroll:   ebitenui.ScrollState{OffsetY: scrollOffset},
			Handlers: scrollHandlers,
		}, children...),
	)
}

func buildNavRoute(route ebitenui.PageRoute, currentPageID string, callbacks *showcaseCallbacks, depth int, chrome showcaseChrome) *ebitenui.Node {
	active := route.ID == currentPageID
	expanded := routeContainsPage(route, currentPageID)
	paddingLeft := 12 + (depth * 18)

	buttonBackground := chrome.PanelBackground
	buttonText := chrome.TextMuted
	if expanded {
		buttonBackground = chrome.CodeBackground
	}
	if active {
		buttonBackground = chrome.Accent
		buttonText = chrome.BadgeText
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
			BorderColor:     chrome.PanelBorder,
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

	if !expanded && !active {
		return item
	}

	children := make([]*ebitenui.Node, 0, len(route.Children)+1)
	children = append(children, item)
	for _, child := range route.Children {
		children = append(children, buildNavRoute(child, currentPageID, callbacks, depth+1, chrome))
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

func buildShowcaseDetail(router *ebitenui.PageRouter, registry ShowcasePageRegistry, page ShowcasePageSpec, currentPageID string, scrollOffset float64, callbacks *showcaseCallbacks, runtime *ebitenui.Runtime, bindings *showcaseBindings, preset showcaseThemePreset, fontPreset showcaseFontPreset) *ebitenui.Node {
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
		ThemePresetID: preset.ID,
		FontPresetID:  fontPreset.ID,
		Theme:         preset.Theme,
		Chrome:        preset.Chrome,
	}
	demo := page.DemoBuilder
	var demoNode *ebitenui.Node
	if demo != nil {
		demoNode = demo(ctx)
	}
	if demoNode == nil {
		demoNode = ebitenui.Text("No demo available for this page yet.", ebitenui.Props{
			ID:    "page-demo-empty",
			Style: showcaseGroupCopyStyleForChrome(preset.Chrome),
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
		Style: showcaseGroupStyleForChrome(preset.Chrome),
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
					Style: detailSectionStyleForChrome(preset.Chrome),
				},
					ebitenui.Text(breadcrumbLabel, ebitenui.Props{
						ID: "page-breadcrumb",
						Style: ebitenui.Style{
							Color: preset.Chrome.Accent,
						},
					}),
					ebitenui.Text(page.Title, ebitenui.Props{
						ID:    "page-title",
						Style: detailTitleStyleForChrome(preset.Chrome),
					}),
					ebitenui.TextBlock(page.Description, ebitenui.Props{
						ID:    "page-description",
						Style: showcaseGroupCopyStyleForChrome(preset.Chrome),
					}),
				),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-demo",
					Style: detailSectionStyleForChrome(preset.Chrome),
				},
					ebitenui.Text("Live Demo", ebitenui.Props{
						ID:    "page-demo-title",
						Style: detailTitleStyleForChrome(preset.Chrome),
					}),
					demoNode,
				),
				buildShowcaseLiveState(currentPageID, bindings, preset, fontPreset),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-usage",
					Style: detailSectionStyleForChrome(preset.Chrome),
				},
					ebitenui.Text("How To Use", ebitenui.Props{
						ID:    "page-usage-title",
						Style: detailTitleStyleForChrome(preset.Chrome),
					}),
					ebitenui.TextBlock(page.UsageNotes, ebitenui.Props{
						ID:    "page-usage-copy",
						Style: showcaseGroupCopyStyleForChrome(preset.Chrome),
					}),
				),
				ebitenui.Div(ebitenui.Props{
					ID:    "page-code",
					Style: detailSectionStyleForChrome(preset.Chrome),
				},
					ebitenui.Text("Code Example", ebitenui.Props{
						ID:    "page-code-title",
						Style: detailTitleStyleForChrome(preset.Chrome),
					}),
					ebitenui.TextBlock(page.CodeExample, ebitenui.Props{
						ID: "page-code-block",
						Style: ebitenui.Style{
							Width:           ebitenui.Fill(),
							Padding:         ebitenui.All(14),
							Color:           preset.Chrome.CodeText,
							LineHeight:      16,
							BackgroundColor: preset.Chrome.CodeBackground,
							BorderColor:     preset.Chrome.CodeBorder,
							BorderWidth:     1,
						},
					}),
				),
			),
		),
	)
}

func buildShowcaseLiveState(currentPageID string, bindings *showcaseBindings, preset showcaseThemePreset, fontPreset showcaseFontPreset) *ebitenui.Node {
	name := ebitenui.NewComputed(func() string {
		return bindings.NameInput.Get()
	})
	resolution := ebitenui.NewComputed(func() string {
		switch bindings.Resolution.Get() {
		case "resolution-1080":
			return "1920x1080"
		default:
			return "1280x720"
		}
	})
	hardcore := ebitenui.NewComputed(func() string {
		if bindings.Hardcore.Get() {
			return "enabled"
		}
		return "disabled"
	})
	bioLines := ebitenui.NewComputed(func() string {
		count := strings.Count(bindings.Bio.Get(), "\n") + 1
		return fmt.Sprintf("%d lines", count)
	})
	summary := ebitenui.NewComputed(func() string {
		return fmt.Sprintf("%s · %s · music %.0f%%", name.Get(), resolution.Get(), bindings.MusicVolume.Get())
	})

	return ebitenui.Div(ebitenui.Props{
		ID:    "page-live-state",
		Style: detailSectionStyleForChrome(preset.Chrome),
	},
		ebitenui.Text("Live State", ebitenui.Props{
			ID:    "page-live-state-title",
			Style: detailTitleStyleForChrome(preset.Chrome),
		}),
		ebitenui.TextBlock("The showcase now exposes the current controlled values so you can see theme and reactive updates without inspecting code first.", ebitenui.Props{
			ID:    "page-live-state-copy",
			Style: showcaseGroupCopyStyleForChrome(preset.Chrome),
		}),
		liveStateRow("live-state-current-page", "Current page", currentPageID, preset),
		liveStateRow("live-state-theme-preset", "Theme preset", preset.Title, preset),
		liveStateRow("live-state-font-preset", "Font preset", fontPreset.Title, preset),
		liveStateRow("live-state-name", "Player name", name.Get(), preset),
		liveStateRow("live-state-resolution", "Resolution", resolution.Get(), preset),
		liveStateRow("live-state-hardcore", "Hardcore", hardcore.Get(), preset),
		liveStateRow("live-state-bio-lines", "Bio", bioLines.Get(), preset),
		liveStateRow("live-state-volume", "Music volume", fmt.Sprintf("%.0f%%", bindings.MusicVolume.Get()), preset),
		liveStateRow("live-state-derived-summary", "Derived summary", summary.Get(), preset),
	)
}

func liveStateRow(id, label, value string, preset showcaseThemePreset) *ebitenui.Node {
	return ebitenui.Div(ebitenui.Props{
		ID: id,
		Style: ebitenui.Style{
			Width:           ebitenui.Fill(),
			Direction:       ebitenui.Row,
			Padding:         ebitenui.Insets{Top: 10, Right: 12, Bottom: 10, Left: 12},
			Gap:             10,
			BackgroundColor: preset.Chrome.CodeBackground,
			BorderColor:     preset.Chrome.CodeBorder,
			BorderWidth:     1,
		},
	},
		ebitenui.Text(label, ebitenui.Props{
			ID:    id + "-label",
			Style: ebitenui.Style{Color: preset.Chrome.TextMuted},
		}),
		ebitenui.Spacer(ebitenui.Props{
			ID:    id + "-spacer",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(1)},
		}),
		ebitenui.Text(value, ebitenui.Props{
			ID:    id + "-value",
			Style: ebitenui.Style{Color: preset.Chrome.TextStrong},
		}),
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
