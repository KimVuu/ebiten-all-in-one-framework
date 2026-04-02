package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestBuildShowcaseDOMUsesPageBasedLayout(t *testing.T) {
	dom := buildShowcaseDOM()

	for _, id := range []string{
		"showcase-root",
		"showcase-header",
		"showcase-main",
		"showcase-sidebar",
		"showcase-detail",
		"page-title",
		"page-demo",
		"page-usage",
		"page-code-block",
	} {
		if _, ok := dom.FindByID(id); !ok {
			t.Fatalf("expected node %q", id)
		}
	}
}

func TestBuildShowcaseDOMLaysOutSidebarAndDetail(t *testing.T) {
	layout := buildShowcaseDOM().Layout(ebitenui.Viewport{Width: 1280, Height: 720})

	sidebar, ok := layout.FindByID("showcase-sidebar")
	if !ok {
		t.Fatalf("expected sidebar layout")
	}
	detail, ok := layout.FindByID("showcase-detail")
	if !ok {
		t.Fatalf("expected detail layout")
	}
	if sidebar.Frame.Width <= 0 || detail.Frame.Width <= 0 {
		t.Fatalf("expected positive widths: sidebar=%#v detail=%#v", sidebar.Frame, detail.Frame)
	}
	if detail.Frame.X <= sidebar.Frame.X {
		t.Fatalf("expected detail to be placed to the right of sidebar")
	}
}

func TestBuildShowcaseDOMShowsCurrentPageTitleCodeAndDemo(t *testing.T) {
	dom := buildShowcaseDOMWithState(showcaseLayoutState{
		CurrentPageID: "inputs/input-field",
	}, nil, nil, nil)

	title, ok := dom.FindByID("page-title")
	if !ok || title.Text != "InputField" {
		t.Fatalf("expected input page title, got %#v", title)
	}
	code, ok := dom.FindByID("page-code-block")
	if !ok || !strings.Contains(code.Text, "ebitenui.InputField") {
		t.Fatalf("expected input field code example, got %#v", code)
	}
	if _, ok := dom.FindByID("name-input"); !ok {
		t.Fatalf("expected input demo node")
	}
	if _, ok := dom.FindByID("page-live-state"); !ok {
		t.Fatalf("expected live state panel")
	}
}

func TestShowcaseSnapshotsIncludeCurrentPageMetadata(t *testing.T) {
	game := newGame(false)

	frame := game.frameSnapshot()
	if frame.DebugEnabled {
		t.Fatalf("expected debug disabled frame snapshot")
	}

	scene := game.sceneSnapshot()
	if got, want := scene.Current.ID, "ebiten-ui-showcase"; got != want {
		t.Fatalf("scene mismatch: got %q want %q", got, want)
	}
	if scene.Summary["currentPageID"] == "" {
		t.Fatalf("expected current page in scene summary")
	}

	world := game.worldSnapshot()
	if len(world.Entities) == 0 {
		t.Fatalf("expected world entities")
	}

	ui := game.uiSnapshot()
	if got, want := ui.Root.ID, "showcase-root"; got != want {
		t.Fatalf("ui root mismatch: got %q want %q", got, want)
	}
	if ui.Root.Props["currentPageID"] == "" {
		t.Fatalf("expected currentPageID in ui snapshot props")
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

func TestShowcaseCompactUIEndpointsReturnSmallPayloads(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if err := game.step(ebitenui.InputSnapshot{PointerX: 24, PointerY: 24}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	client := http.Client{Timeout: 2 * time.Second}

	overviewResponse, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui/overview")
	if err != nil {
		t.Fatalf("overview request failed: %v", err)
	}
	defer overviewResponse.Body.Close()
	var overviewPayload map[string]any
	if err := json.NewDecoder(overviewResponse.Body).Decode(&overviewPayload); err != nil {
		t.Fatalf("decode overview failed: %v", err)
	}
	if overviewPayload["rootId"] == "" {
		t.Fatalf("expected rootId in overview payload")
	}
	overviewEncoded, _ := json.Marshal(overviewPayload)
	if bytes.Contains(overviewEncoded, []byte(`"children"`)) {
		t.Fatalf("overview payload should not include full tree")
	}
	if len(overviewEncoded) > 4096 {
		t.Fatalf("expected compact overview payload, got %d bytes", len(overviewEncoded))
	}

	queryBody := bytes.NewBufferString(`{"visible_only":true,"limit":10}`)
	queryResponse, err := client.Post("http://"+game.debugBridge.Address()+"/debug/ui/query", "application/json", queryBody)
	if err != nil {
		t.Fatalf("query request failed: %v", err)
	}
	defer queryResponse.Body.Close()
	var queryPayload map[string]any
	if err := json.NewDecoder(queryResponse.Body).Decode(&queryPayload); err != nil {
		t.Fatalf("decode query failed: %v", err)
	}
	nodes, ok := queryPayload["nodes"].([]any)
	if !ok || len(nodes) == 0 {
		t.Fatalf("expected queried nodes, got %#v", queryPayload["nodes"])
	}
	queryEncoded, _ := json.Marshal(queryPayload)
	if len(queryEncoded) > 6144 {
		t.Fatalf("expected compact query payload, got %d bytes", len(queryEncoded))
	}

	nodeResponse, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui/node/showcase-detail")
	if err != nil {
		t.Fatalf("node request failed: %v", err)
	}
	defer nodeResponse.Body.Close()
	var nodePayload map[string]any
	if err := json.NewDecoder(nodeResponse.Body).Decode(&nodePayload); err != nil {
		t.Fatalf("decode node failed: %v", err)
	}
	if _, ok := nodePayload["summary"].(map[string]any); !ok {
		t.Fatalf("expected summary in node payload")
	}
	if _, ok := nodePayload["children"].([]any); !ok {
		t.Fatalf("expected children summaries in node payload")
	}

	issuesResponse, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui/issues?limit=10")
	if err != nil {
		t.Fatalf("issues request failed: %v", err)
	}
	defer issuesResponse.Body.Close()
	var issuesPayload map[string]any
	if err := json.NewDecoder(issuesResponse.Body).Decode(&issuesPayload); err != nil {
		t.Fatalf("decode issues failed: %v", err)
	}
	if _, ok := issuesPayload["issues"].([]any); !ok {
		t.Fatalf("expected issues array in issues payload")
	}
}

func TestShowcaseCaptureEndpointReturnsArtifactMetadataAndFile(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720
	if err := game.step(ebitenui.InputSnapshot{PointerX: 24, PointerY: 24}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.stopDebugBridge()
	}()

	client := http.Client{Timeout: 2 * time.Second}
	captureBody := bytes.NewBufferString(`{"target":"node_id","node_id":"showcase-detail","with_overlay":true}`)
	captureResponse, err := client.Post("http://"+game.debugBridge.Address()+"/debug/ui/capture", "application/json", captureBody)
	if err != nil {
		t.Fatalf("capture request failed: %v", err)
	}
	defer captureResponse.Body.Close()

	var capturePayload map[string]any
	if err := json.NewDecoder(captureResponse.Body).Decode(&capturePayload); err != nil {
		t.Fatalf("decode capture failed: %v", err)
	}
	artifactID, _ := capturePayload["artifactId"].(string)
	if artifactID == "" {
		t.Fatalf("expected artifactId in capture payload")
	}
	if path, _ := capturePayload["path"].(string); path == "" {
		t.Fatalf("expected artifact path in capture payload")
	} else if !strings.Contains(path, filepath.Join("screenshots", "ebiten-ui-showcase")) {
		t.Fatalf("expected artifact path to be stored under screenshots/ebiten-ui-showcase, got %q", path)
	}
	captureEncoded, _ := json.Marshal(capturePayload)
	if bytes.Contains(captureEncoded, []byte(`iVBOR`)) {
		t.Fatalf("capture response should not inline image bytes")
	}

	artifactResponse, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui/artifacts/" + artifactID)
	if err != nil {
		t.Fatalf("artifact request failed: %v", err)
	}
	defer artifactResponse.Body.Close()
	if got, want := artifactResponse.Header.Get("Content-Type"), "image/png"; got != want {
		t.Fatalf("artifact content type mismatch: got %q want %q", got, want)
	}
	header := make([]byte, 8)
	if _, err := io.ReadFull(artifactResponse.Body, header); err != nil {
		t.Fatalf("read artifact header failed: %v", err)
	}
	if !bytes.HasPrefix(header, []byte{0x89, 'P', 'N', 'G'}) {
		t.Fatalf("expected png header, got %x", header)
	}
}

func TestShowcaseGameAppliesWheelScrollToSidebar(t *testing.T) {
	game := newGame(false)
	game.width = 1280
	game.height = 720

	initial := game.sidebarScroll
	if err := game.step(ebitenui.InputSnapshot{PointerX: 120, PointerY: 220, ScrollY: -1}); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	if game.sidebarScroll <= initial {
		t.Fatalf("expected sidebar scroll to increase, got initial=%v current=%v", initial, game.sidebarScroll)
	}

	layout := game.currentLayout()
	scroll, ok := layout.FindByID("showcase-sidebar-scroll")
	if !ok {
		t.Fatalf("expected showcase sidebar scroll layout")
	}
	if got, want := scroll.Node.Props.Scroll.OffsetY, game.sidebarScroll; got != want {
		t.Fatalf("scroll offset mismatch: got %v want %v", got, want)
	}
}

func TestShowcaseUISnapshotIncludesExpandedMetadata(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(ebitenui.InputSnapshot{PointerX: 24, PointerY: 24}); err != nil {
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
	if got, want := ui.Root.Props["themePreset"], "default"; got != want {
		t.Fatalf("theme preset mismatch: got %#v want %q", got, want)
	}
}

func TestShowcaseDebugOverlayDefaultsOff(t *testing.T) {
	game := newGame(true)
	if game.overlayEnabled {
		t.Fatalf("expected debug overlay to default off")
	}

	ui := game.uiSnapshot()
	if got, want := ui.Root.Props["overlay"], false; got != want {
		t.Fatalf("overlay prop mismatch: got %#v want %#v", got, want)
	}
}

func TestShowcaseThemePresetButtonsSwitchThemeState(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if got, want := game.currentState().ThemePreset, "default"; got != want {
		t.Fatalf("expected default preset, got %q", got)
	}

	result := game.debugBridgeLikeCommand("ui_click", map[string]any{
		"node_id": "theme-preset-forest",
	})
	if !result.Success {
		t.Fatalf("expected theme preset click to succeed: %#v", result)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("theme step failed: %v", err)
		}
	}

	if got, want := game.currentState().ThemePreset, "forest"; got != want {
		t.Fatalf("expected forest preset, got %q", got)
	}
	ui := game.uiSnapshot()
	if got, want := ui.Root.Props["themePreset"], "forest"; got != want {
		t.Fatalf("expected ui snapshot theme preset %q, got %#v", want, got)
	}
}

func TestShowcaseLiveStatePanelTracksBindingValues(t *testing.T) {
	game := newGame(true)
	game.width = 1280
	game.height = 720

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("initial step failed: %v", err)
	}
	if result := game.debugBridgeLikeCommand("ui_click", map[string]any{
		"node_id": "nav-item-inputs",
	}); !result.Success {
		t.Fatalf("ui_click failed: %#v", result)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("navigation step failed: %v", err)
		}
	}
	if result := game.debugBridgeLikeCommand("ui_type_text", map[string]any{
		"node_id": "name-input",
		"text":    "A",
	}); !result.Success {
		t.Fatalf("ui_type_text failed: %#v", result)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("text step failed: %v", err)
	}

	node, ok := game.dom.FindByID("live-state-name-value")
	if !ok {
		t.Fatalf("expected live state value node")
	}
	if got, want := node.Text, "KimA"; got != want {
		t.Fatalf("expected live state text %q, got %q", want, got)
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
		"node_id": "showcase-sidebar-scroll",
		"delta_y": -1.0,
	})
	if !scrollResult.Success {
		t.Fatalf("expected queued scroll command, got %#v", scrollResult)
	}
	initialScroll := game.sidebarScroll
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if game.sidebarScroll <= initialScroll {
		t.Fatalf("expected sidebar scroll after debug command, got initial=%v current=%v", initialScroll, game.sidebarScroll)
	}

	clickResult := game.debugBridge.InvokeCommand("ui_click", map[string]any{
		"node_id": "nav-item-inputs",
	})
	if !clickResult.Success {
		t.Fatalf("expected queued click command, got %#v", clickResult)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("navigation step failed: %v", err)
		}
	}

	typeResult := game.debugBridge.InvokeCommand("ui_type_text", map[string]any{
		"node_id": "name-input",
		"text":    "A",
	})
	if !typeResult.Success {
		t.Fatalf("expected queued type command, got %#v", typeResult)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("text step failed: %v", err)
	}
	if got, want := game.runtime.FocusedID(), "name-input"; got != want {
		t.Fatalf("focus mismatch: got %q want %q", got, want)
	}
	if got := game.runtime.TextValueOrDefault("name-input", "Kim"); got != "KimA" {
		t.Fatalf("expected runtime text value to update, got %q", got)
	}
	if got := game.bindings.NameInput.Get(); got != "KimA" {
		t.Fatalf("expected bound input value to update, got %q", got)
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

	if result := game.debugBridge.InvokeCommand("ui_click", map[string]any{
		"node_id": "nav-item-inputs",
	}); !result.Success {
		t.Fatalf("ui_click failed: %s", result.Message)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("navigation step failed: %v", err)
		}
	}

	if result := game.debugBridge.InvokeCommand("ui_type_text", map[string]any{
		"node_id": "name-input",
		"text":    "Z",
	}); !result.Success {
		t.Fatalf("ui_type_text failed: %s", result.Message)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("typing step failed: %v", err)
	}

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "home",
	}); !result.Success {
		t.Fatalf("ui_key_event home failed: %s", result.Message)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("home focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("home key step failed: %v", err)
	}

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "delete",
	}); !result.Success {
		t.Fatalf("ui_key_event delete failed: %s", result.Message)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("delete focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("delete key step failed: %v", err)
	}

	if got, want := game.runtime.TextValueOrDefault("name-input", "Kim"), "imZ"; got != want {
		t.Fatalf("expected home/delete editing result, got %q want %q", got, want)
	}
	if got, want := game.bindings.NameInput.Get(), "imZ"; got != want {
		t.Fatalf("expected binding editing result, got %q want %q", got, want)
	}
}
