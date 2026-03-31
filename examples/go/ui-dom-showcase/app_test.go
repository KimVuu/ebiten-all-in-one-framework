package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestBuildShowcaseDOMIncludesAllSupportedTags(t *testing.T) {
	dom := buildShowcaseDOM()

	required := map[uidom.Tag]bool{
		uidom.TagDiv:        false,
		uidom.TagHeader:     false,
		uidom.TagMain:       false,
		uidom.TagSection:    false,
		uidom.TagFooter:     false,
		uidom.TagButton:     false,
		uidom.TagSpan:       false,
		uidom.TagText:       false,
		uidom.TagImage:      false,
		uidom.TagTextBlock:  false,
		uidom.TagSpacer:     false,
		uidom.TagStack:      false,
		uidom.TagScrollView: false,
	}

	var walk func(*uidom.Node)
	walk = func(node *uidom.Node) {
		if node == nil {
			return
		}
		if _, ok := required[node.Tag]; ok {
			required[node.Tag] = true
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(dom.Root)

	for tag, found := range required {
		if !found {
			t.Fatalf("expected showcase DOM to include tag %q", tag)
		}
	}
}

func TestBuildShowcaseDOMLaysOutKeySections(t *testing.T) {
	dom := buildShowcaseDOM()
	layout := dom.Layout(uidom.Viewport{Width: 1440, Height: 2200})

	ids := []string{
		"showcase-root",
		"showcase-header",
		"showcase-main",
		"overview-section",
		"button-section",
		"foundation-section",
		"components-section",
		"prefabs-section",
		"scroll-preview",
		"showcase-footer",
	}

	for _, id := range ids {
		node, ok := layout.FindByID(id)
		if !ok {
			t.Fatalf("expected layout node %q", id)
		}
		if node.Frame.Width <= 0 || node.Frame.Height <= 0 {
			t.Fatalf("expected positive frame for %q, got %#v", id, node.Frame)
		}
	}
}

func TestBuildShowcaseDOMUsesPageScrollAndKeepsHeaderWithinViewport(t *testing.T) {
	dom := buildShowcaseDOM()
	layout := dom.Layout(uidom.Viewport{Width: 1280, Height: 720})

	scroll, ok := layout.FindByID("showcase-scroll")
	if !ok || scroll.Node.Tag != uidom.TagScrollView {
		t.Fatalf("expected page scroll view")
	}

	header, ok := layout.FindByID("showcase-header")
	if !ok {
		t.Fatalf("expected header layout")
	}
	if header.Frame.X+header.Frame.Width > 1280 {
		t.Fatalf("expected header within viewport, got %#v", header.Frame)
	}

	badge, ok := layout.FindByID("header-badge")
	if !ok {
		t.Fatalf("expected header badge")
	}
	if badge.Frame.X+badge.Frame.Width > 1280 {
		t.Fatalf("expected badge within viewport, got %#v", badge.Frame)
	}
}

func TestBuildShowcaseDOMIncludesRepresentativeComponentsAndPrefabs(t *testing.T) {
	dom := buildShowcaseDOM()

	ids := []string{
		"form-section",
		"layout-section",
		"overlay-section",
		"data-section",
		"status-section",
		"name-input",
		"difficulty-toggle",
		"music-slider",
		"inventory-scrollbar",
		"resolution-dropdown",
		"bio-textarea",
		"mode-radio",
		"party-stepper",
		"exp-progress",
		"content-grid",
		"virtual-items",
		"settings-modal",
		"loot-tooltip",
		"slot-context-menu",
		"tabs-demo",
		"accordion-demo",
		"elite-badge",
		"fire-chip",
		"dialog-demo",
		"hud-demo",
		"inventory-demo",
		"pause-demo",
		"settings-demo",
		"tooltip-demo",
	}

	for _, id := range ids {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected unified showcase node %q", id)
		}
	}
}

func TestBuildShowcaseDOMPlacesFooterAfterPrefabs(t *testing.T) {
	layout := buildShowcaseDOM().Layout(uidom.Viewport{Width: 1280, Height: 720})

	prefabsNode, ok := layout.FindByID("prefabs-section")
	if !ok {
		t.Fatalf("expected prefabs section")
	}
	footer, ok := layout.FindByID("showcase-footer")
	if !ok {
		t.Fatalf("expected footer")
	}
	if footer.Frame.Y < prefabsNode.Frame.Y+prefabsNode.Frame.Height {
		t.Fatalf("expected footer after prefabs section, got prefabs=%#v footer=%#v", prefabsNode.Frame, footer.Frame)
	}
}

func TestShowcaseSnapshotsIncludeRootAndSections(t *testing.T) {
	game := newGame(false)

	frame := game.frameSnapshot()
	if frame.DebugEnabled {
		t.Fatalf("expected debug disabled frame snapshot")
	}

	scene := game.sceneSnapshot()
	if got, want := scene.Current.ID, "ui-dom-showcase"; got != want {
		t.Fatalf("scene mismatch: got %q want %q", got, want)
	}

	world := game.worldSnapshot()
	if len(world.Entities) < 4 {
		t.Fatalf("expected major section entities, got %d", len(world.Entities))
	}

	ui := game.uiSnapshot()
	if got, want := ui.Root.ID, "showcase-root"; got != want {
		t.Fatalf("ui root mismatch: got %q want %q", got, want)
	}
	if len(ui.Root.Children) < 2 {
		t.Fatalf("expected root child nodes, got %d", len(ui.Root.Children))
	}
	if got, want := ui.Root.Children[1].ID, "showcase-scroll"; got != want {
		t.Fatalf("expected page scroll child, got %q want %q", got, want)
	}
}

func TestShowcaseDebugBridgeStartsWhenEnabled(t *testing.T) {
	game := newGame(true)
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	client := http.Client{Timeout: time.Second}
	response, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui")
	if err != nil {
		t.Fatalf("ui request failed: %v", err)
	}
	defer response.Body.Close()

	if got, want := response.StatusCode, http.StatusOK; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}

	if err := game.debugBridge.Close(context.Background()); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	game.debugBridge = nil
}

func TestShowcaseGameAppliesWheelScrollToPageLayout(t *testing.T) {
	game := newGame(false)
	game.width = 1280
	game.height = 720

	initial := game.pageScroll
	if err := game.step(uidom.InputSnapshot{PointerX: 120, PointerY: 160, ScrollY: -1}); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	if game.pageScroll <= initial {
		t.Fatalf("expected page scroll to increase, got initial=%v current=%v", initial, game.pageScroll)
	}

	layout := game.currentLayout()
	scroll, ok := layout.FindByID("showcase-scroll")
	if !ok {
		t.Fatalf("expected showcase scroll layout")
	}
	if got, want := scroll.Node.Props.Scroll.OffsetY, game.pageScroll; got != want {
		t.Fatalf("scroll offset mismatch: got %v want %v", got, want)
	}
}

func TestShowcaseUISnapshotIncludesExpandedMetadata(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(uidom.InputSnapshot{PointerX: 24, PointerY: 24}); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	ui := game.uiSnapshot()
	if ui.Root.Semantic == nil {
		t.Fatalf("expected semantic metadata on root")
	}
	if ui.Root.Layout == nil {
		t.Fatalf("expected layout metadata on root")
	}
	if ui.Root.Computed == nil {
		t.Fatalf("expected computed metadata on root")
	}
	if ui.InputState.Pointer == nil {
		t.Fatalf("expected pointer input state")
	}
	if ui.Viewport.Width != 1280 || ui.Viewport.Height != 720 {
		t.Fatalf("unexpected viewport snapshot: %#v", ui.Viewport)
	}
}

func TestDebugBridgeCommandsQueueScrollAndTextInput(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	scrollResult := game.debugBridge.InvokeCommand("ui_scroll", map[string]any{
		"node_id": "showcase-scroll",
		"delta_y": -1.0,
	})
	if !scrollResult.Success {
		t.Fatalf("expected queued scroll command, got %#v", scrollResult)
	}
	initialScroll := game.pageScroll
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if game.pageScroll <= initialScroll {
		t.Fatalf("expected page scroll after debug command, got initial=%v current=%v", initialScroll, game.pageScroll)
	}

	typeResult := game.debugBridge.InvokeCommand("ui_type_text", map[string]any{
		"node_id": "name-input",
		"text":    "A",
	})
	if !typeResult.Success {
		t.Fatalf("expected queued type command, got %#v", typeResult)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("focus step failed: %v", err)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("text step failed: %v", err)
	}
	if got, want := game.runtime.FocusedID(), "name-input"; got != want {
		t.Fatalf("focus mismatch: got %q want %q", got, want)
	}
	if got := game.runtime.TextValueOrDefault("name-input", "Kim"); got != "KimA" {
		t.Fatalf("expected runtime text value to update, got %q", got)
	}
}

func TestDebugBridgeKeyEventRoutesShortcutsAndEditing(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	if result := game.debugBridge.InvokeCommand("ui_type_text", map[string]any{
		"node_id": "name-input",
		"text":    "Z",
	}); !result.Success {
		t.Fatalf("ui_type_text failed: %s", result.Message)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("focus step failed: %v", err)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("typing step failed: %v", err)
	}

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "home",
	}); !result.Success {
		t.Fatalf("ui_key_event home failed: %s", result.Message)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("home focus step failed: %v", err)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("home key step failed: %v", err)
	}

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "delete",
	}); !result.Success {
		t.Fatalf("ui_key_event delete failed: %s", result.Message)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("delete focus step failed: %v", err)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("delete key step failed: %v", err)
	}

	if got, want := game.runtime.TextValueOrDefault("name-input", "Kim"), "imZ"; got != want {
		t.Fatalf("expected home/delete editing result, got %q want %q", got, want)
	}

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "a",
		"control": true,
	}); !result.Success {
		t.Fatalf("ui_key_event ctrl+a failed: %s", result.Message)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("ctrl+a focus step failed: %v", err)
	}
	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("ctrl+a key step failed: %v", err)
	}
}

func TestShowcaseDebugCommandsQueueScrollAndToggleOverlay(t *testing.T) {
	game := newGame(true)
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	result := game.debugBridge.InvokeCommand("set_ui_debug_overlay", map[string]any{
		"enabled": false,
	})
	if !result.Success {
		t.Fatalf("set_ui_debug_overlay failed: %s", result.Message)
	}
	if game.overlayEnabled {
		t.Fatalf("expected overlay to be disabled")
	}

	result = game.debugBridge.InvokeCommand("ui_scroll", map[string]any{
		"node_id": "showcase-scroll",
		"delta_y": -1,
	})
	if !result.Success {
		t.Fatalf("ui_scroll failed: %s", result.Message)
	}

	if err := game.step(uidom.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if game.pageScroll <= 0 {
		t.Fatalf("expected page scroll to move, got %v", game.pageScroll)
	}
}
