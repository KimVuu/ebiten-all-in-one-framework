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

func TestBuildShowcaseDOMIncludesAllSupportedTags(t *testing.T) {
	dom := buildShowcaseDOM()

	required := map[ebitenui.Tag]bool{
		ebitenui.TagDiv:        false,
		ebitenui.TagHeader:     false,
		ebitenui.TagMain:       false,
		ebitenui.TagSection:    false,
		ebitenui.TagFooter:     false,
		ebitenui.TagButton:     false,
		ebitenui.TagSpan:       false,
		ebitenui.TagText:       false,
		ebitenui.TagImage:      false,
		ebitenui.TagTextBlock:  false,
		ebitenui.TagSpacer:     false,
		ebitenui.TagStack:      false,
		ebitenui.TagScrollView: false,
	}

	var walk func(*ebitenui.Node)
	walk = func(node *ebitenui.Node) {
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
	layout := dom.Layout(ebitenui.Viewport{Width: 1440, Height: 2200})

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
	layout := dom.Layout(ebitenui.Viewport{Width: 1280, Height: 720})

	scroll, ok := layout.FindByID("showcase-scroll")
	if !ok || scroll.Node.Tag != ebitenui.TagScrollView {
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
	layout := buildShowcaseDOM().Layout(ebitenui.Viewport{Width: 1280, Height: 720})

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
	if got, want := scene.Current.ID, "ebiten-ui-showcase"; got != want {
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

	nodeResponse, err := client.Get("http://" + game.debugBridge.Address() + "/debug/ui/node/name-input")
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
	nodeEncoded, _ := json.Marshal(nodePayload)
	if len(nodeEncoded) > 8192 {
		t.Fatalf("expected compact node payload, got %d bytes", len(nodeEncoded))
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
	captureBody := bytes.NewBufferString(`{"target":"node_id","node_id":"name-input","with_overlay":true}`)
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

func TestShowcaseGameAppliesWheelScrollToPageLayout(t *testing.T) {
	game := newGame(false)
	game.width = 1280
	game.height = 720

	initial := game.pageScroll
	if err := game.step(ebitenui.InputSnapshot{PointerX: 120, PointerY: 160, ScrollY: -1}); err != nil {
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
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
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

	if result := game.debugBridge.InvokeCommand("ui_key_event", map[string]any{
		"node_id": "name-input",
		"key":     "a",
		"control": true,
	}); !result.Success {
		t.Fatalf("ui_key_event ctrl+a failed: %s", result.Message)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("ctrl+a focus step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
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

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("step failed: %v", err)
	}
	if game.pageScroll <= 0 {
		t.Fatalf("expected page scroll to move, got %v", game.pageScroll)
	}
}
