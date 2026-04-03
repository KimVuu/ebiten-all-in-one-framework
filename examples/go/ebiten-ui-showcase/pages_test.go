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
	if _, ok := registry.Pages["inputs/button-events"]; !ok {
		t.Fatalf("expected button-events page")
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
		"font-preset-default",
		"font-preset-neo-dunggeunmo",
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

func TestShowcaseButtonEventsPageUpdatesLifecycleBindings(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if !game.router.Navigate("inputs/button-events") {
		t.Fatalf("expected button-events page route")
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	layout := game.currentLayout()
	button, ok := layout.FindByID("button-events-demo-button")
	if !ok {
		t.Fatalf("expected button events demo button")
	}
	x := button.Frame.X + button.Frame.Width*0.5
	y := button.Frame.Y + button.Frame.Height*0.5

	if err := game.step(ebitenui.InputSnapshot{PointerX: x, PointerY: y}); err != nil {
		t.Fatalf("hover step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true}); err != nil {
		t.Fatalf("down step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true}); err != nil {
		t.Fatalf("hold step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{PointerX: x, PointerY: y}); err != nil {
		t.Fatalf("up step failed: %v", err)
	}

	if got, want := game.bindings.ButtonDowns.Get(), 1; got != want {
		t.Fatalf("button down count mismatch: got %d want %d", got, want)
	}
	if got := game.bindings.ButtonHolds.Get(); got < 1 {
		t.Fatalf("expected hold count, got %d", got)
	}
	if got, want := game.bindings.ButtonUps.Get(), 1; got != want {
		t.Fatalf("button up count mismatch: got %d want %d", got, want)
	}
	if got, want := game.bindings.ButtonClicks.Get(), 1; got != want {
		t.Fatalf("button click count mismatch: got %d want %d", got, want)
	}
	if got, want := game.bindings.ButtonPhase.Get(), "click"; got != want {
		t.Fatalf("button phase mismatch: got %q want %q", got, want)
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
	if ui.Root.Props["fontPreset"] == "" {
		t.Fatalf("expected fontPreset in root props")
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
