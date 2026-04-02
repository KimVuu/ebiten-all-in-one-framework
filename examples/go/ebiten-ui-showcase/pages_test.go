package main

import (
	"testing"

	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestShowcaseRegistryContainsGroupsAndLeafPages(t *testing.T) {
	registry := buildShowcasePageRegistry()
	if len(registry.Routes) == 0 {
		t.Fatalf("expected showcase routes")
	}
	if _, ok := registry.Pages["reactive/ref-and-computed"]; !ok {
		t.Fatalf("expected reactive ref-and-computed page")
	}
	if _, ok := registry.Pages["reactive/controlled-inputs"]; !ok {
		t.Fatalf("expected reactive controlled-inputs page")
	}
	if _, ok := registry.Pages["inputs/input-field"]; !ok {
		t.Fatalf("expected input-field page")
	}
	if _, ok := registry.Pages["foundations/theme"]; !ok {
		t.Fatalf("expected theme page")
	}
	if _, ok := registry.Pages["prefabs/dialog"]; !ok {
		t.Fatalf("expected dialog page")
	}
}

func TestBuildShowcaseDOMBuildsSidebarAndDetailForCurrentPage(t *testing.T) {
	dom := buildShowcaseDOMWithState(showcaseLayoutState{
		CurrentPageID: "inputs/input-field",
	}, nil, nil, nil)

	for _, id := range []string{
		"showcase-root",
		"showcase-sidebar",
		"showcase-detail",
		"page-title",
		"page-code-block",
		"page-live-state",
		"theme-preset-default",
		"theme-preset-forest",
		"theme-preset-ember",
	} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected node %q", id)
		}
	}
	if _, ok := dom.FindByID("nav-item-inputs-input-field"); !ok {
		t.Fatalf("expected nav item for current page")
	}
	if inputNode, ok := dom.FindByID("name-input"); !ok || inputNode == nil {
		t.Fatalf("expected current page demo node")
	}
}

func TestBuildShowcaseDOMBuildsReactivePage(t *testing.T) {
	dom := buildShowcaseDOMWithState(showcaseLayoutState{
		CurrentPageID: "reactive/ref-and-computed",
	}, nil, nil, nil)

	title, ok := dom.FindByID("page-title")
	if !ok || title.Text != "Ref And Computed" {
		t.Fatalf("expected reactive page title, got %#v", title)
	}
	if _, ok := dom.FindByID("reactive-ref-name"); !ok {
		t.Fatalf("expected reactive demo input")
	}
	if _, ok := dom.FindByID("reactive-derived-summary"); !ok {
		t.Fatalf("expected reactive derived summary")
	}
}

func TestShowcaseGameTracksCurrentPageInDebugState(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	scene := game.sceneSnapshot()
	if scene.Summary["currentPageID"] == "" {
		t.Fatalf("expected currentPageID in scene summary")
	}

	ui := game.uiSnapshot()
	if ui.Root.Props["currentPageID"] == "" {
		t.Fatalf("expected currentPageID in root props")
	}
	if ui.Root.Props["themePreset"] == "" {
		t.Fatalf("expected themePreset in root props")
	}
}

func TestShowcaseSidebarClickNavigatesPages(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	before := game.currentPageID()
	result := game.debugBridgeLikeCommand("ui_click", map[string]any{
		"node_id": "nav-item-inputs",
	})
	if !result.Success {
		t.Fatalf("expected click command to succeed: %#v", result)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("step failed: %v", err)
		}
	}
	after := game.currentPageID()
	if before == after {
		t.Fatalf("expected page navigation after sidebar click")
	}
	if got, want := after, "inputs/input-field"; got != want {
		t.Fatalf("page mismatch: got %q want %q", got, want)
	}
}

func TestShowcaseCaptureUsesCurrentPageLayout(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	result, ok := game.uiDebug.UICapture(ebitendebug.UICaptureRequest{
		Target: "node_id",
		NodeID: "showcase-detail",
	})
	if !ok {
		t.Fatalf("expected capture result")
	}
	if result.Path == "" {
		t.Fatalf("expected artifact path")
	}
}
