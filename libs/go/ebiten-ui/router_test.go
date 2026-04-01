package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestPageRouterNavigatesNestedRoutes(t *testing.T) {
	router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{
		InitialPageID: "inputs",
		Routes: []ebitenui.PageRoute{
			{
				ID:             "inputs",
				Title:          "Inputs",
				DefaultChildID: "inputs/input-field",
				Children: []ebitenui.PageRoute{
					{ID: "inputs/input-field", Title: "Input Field"},
					{ID: "inputs/dropdown", Title: "Dropdown"},
				},
			},
			{
				ID:    "prefabs",
				Title: "Prefabs",
				Children: []ebitenui.PageRoute{
					{ID: "prefabs/dialog", Title: "Dialog"},
				},
			},
		},
	})

	if got, want := router.CurrentPageID(), "inputs/input-field"; got != want {
		t.Fatalf("current page mismatch: got %q want %q", got, want)
	}
	if !router.Navigate("prefabs/dialog") {
		t.Fatalf("expected navigate to leaf page to succeed")
	}
	if got, want := router.CurrentPageID(), "prefabs/dialog"; got != want {
		t.Fatalf("current page mismatch after navigate: got %q want %q", got, want)
	}
	if router.Navigate("missing") {
		t.Fatalf("expected navigate to unknown page to fail")
	}
}

func TestPageRouterUsesDefaultChildAndBuildsBreadcrumb(t *testing.T) {
	router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{
		InitialPageID: "overlay",
		Routes: []ebitenui.PageRoute{
			{
				ID:             "overlay",
				Title:          "Overlay",
				DefaultChildID: "overlay/modal",
				Children: []ebitenui.PageRoute{
					{ID: "overlay/modal", Title: "Modal"},
					{ID: "overlay/tooltip", Title: "Tooltip"},
				},
			},
		},
	})

	if got, want := router.CurrentPageID(), "overlay/modal"; got != want {
		t.Fatalf("default child mismatch: got %q want %q", got, want)
	}

	breadcrumb := router.Breadcrumb("overlay/tooltip")
	if got, want := len(breadcrumb), 2; got != want {
		t.Fatalf("breadcrumb length mismatch: got %d want %d", got, want)
	}
	if got, want := breadcrumb[0].ID, "overlay"; got != want {
		t.Fatalf("breadcrumb root mismatch: got %q want %q", got, want)
	}
	if got, want := breadcrumb[1].ID, "overlay/tooltip"; got != want {
		t.Fatalf("breadcrumb leaf mismatch: got %q want %q", got, want)
	}

	children := router.ChildrenOf("overlay")
	if got, want := len(children), 2; got != want {
		t.Fatalf("children length mismatch: got %d want %d", got, want)
	}
	if got, want := router.ParentPageID("overlay/modal"), "overlay"; got != want {
		t.Fatalf("parent page mismatch: got %q want %q", got, want)
	}
}

func TestVisibleNavTreeFiltersHiddenRoutes(t *testing.T) {
	router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{
		InitialPageID: "overview",
		Routes: []ebitenui.PageRoute{
			{ID: "overview", Title: "Overview"},
			{
				ID:    "inputs",
				Title: "Inputs",
				Children: []ebitenui.PageRoute{
					{ID: "inputs/input-field", Title: "Input Field"},
					{ID: "inputs/internal", Title: "Internal", Hidden: true},
				},
			},
		},
	})

	tree := router.VisibleNavTree()
	if got, want := len(tree), 2; got != want {
		t.Fatalf("nav tree length mismatch: got %d want %d", got, want)
	}
	if got, want := len(tree[1].Children), 1; got != want {
		t.Fatalf("visible child count mismatch: got %d want %d", got, want)
	}
	if got, want := tree[1].Children[0].ID, "inputs/input-field"; got != want {
		t.Fatalf("visible child mismatch: got %q want %q", got, want)
	}
}

func TestPageScreenBuildsSidebarAndContentColumns(t *testing.T) {
	dom := ebitenui.New(ebitenui.PageScreen(ebitenui.PageScreenConfig{
		ID: "page-screen",
		Sidebar: ebitenui.Div(ebitenui.Props{
			ID: "page-nav",
			Style: ebitenui.Style{Width: ebitenui.Px(240), Height: ebitenui.Fill()},
		}),
		Content: ebitenui.Div(ebitenui.Props{
			ID: "page-content",
			Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Fill()},
		}),
	}))

	layout := dom.Layout(ebitenui.Viewport{Width: 1280, Height: 720})
	nav, ok := layout.FindByID("page-nav")
	if !ok {
		t.Fatalf("expected nav node")
	}
	content, ok := layout.FindByID("page-content")
	if !ok {
		t.Fatalf("expected content node")
	}
	if content.Frame.X <= nav.Frame.X {
		t.Fatalf("expected content to be placed to the right of nav")
	}
}
