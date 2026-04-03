package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestGameDebugSurfaceAndScreenFlow(t *testing.T) {
	game := newGame(GameConfig{
		DebugEnabled: true,
		Seed:         7,
	})
	game.width = DefaultWindowWidth
	game.height = DefaultWindowHeight

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("initial step failed: %v", err)
	}
	if _, ok := game.dom.FindByID("party-selection-screen"); !ok {
		t.Fatalf("expected party selection screen")
	}
	if _, ok := game.dom.FindByID("party-option-human-warrior"); !ok {
		t.Fatalf("expected warrior party button")
	}
	if _, ok := game.dom.FindByID("party-selection-grid"); !ok {
		t.Fatalf("expected party selection grid")
	}
	if containsTag(game.dom.Root, ebitenui.TagScrollView) {
		t.Fatalf("expected no scroll view nodes on party selection screen")
	}
	assertNoLayoutOverflow(t, game)

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
	var overview map[string]any
	if err := json.NewDecoder(overviewResponse.Body).Decode(&overview); err != nil {
		t.Fatalf("decode overview failed: %v", err)
	}
	if overview["rootId"] == "" {
		t.Fatalf("expected rootId in overview payload")
	}

	captureBody := bytes.NewBufferString(`{"target":"node_id","node_id":"party-selection-screen","with_overlay":true}`)
	captureResponse, err := client.Post("http://"+game.debugBridge.Address()+"/debug/ui/capture", "application/json", captureBody)
	if err != nil {
		t.Fatalf("capture request failed: %v", err)
	}
	defer captureResponse.Body.Close()
	var capture map[string]any
	if err := json.NewDecoder(captureResponse.Body).Decode(&capture); err != nil {
		t.Fatalf("decode capture failed: %v", err)
	}
	if path, _ := capture["path"].(string); !strings.Contains(path, filepath.Join("screenshots", "dice-rogue")) {
		t.Fatalf("expected screenshot to be stored under screenshots/dice-rogue, got %q", path)
	}

	clickAndStep(t, game, "party-option-human-warrior")
	clickAndStep(t, game, "party-option-human-guard")
	clickAndStep(t, game, "party-option-human-guide")
	clickAndStep(t, game, "start-run-button")

	if got, want := game.currentScreen(), ScreenMap; got != want {
		t.Fatalf("screen mismatch after party start: got %q want %q", got, want)
	}
	if _, ok := game.dom.FindByID("map-screen"); !ok {
		t.Fatalf("expected map screen after starting run")
	}
	if _, ok := game.dom.FindByID("map-node-grid"); !ok {
		t.Fatalf("expected map node grid")
	}
	if containsTag(game.dom.Root, ebitenui.TagScrollView) {
		t.Fatalf("expected no scroll view nodes on map screen")
	}
	assertNoLayoutOverflow(t, game)

	clickAndStep(t, game, "map-node-normal-a")
	if got, want := game.currentScreen(), ScreenCombat; got != want {
		t.Fatalf("screen mismatch after map node selection: got %q want %q", got, want)
	}
	if _, ok := game.dom.FindByID("combat-screen"); !ok {
		t.Fatalf("expected combat screen")
	}
	if _, ok := game.dom.FindByID("resolve-turn-button"); !ok {
		t.Fatalf("expected resolve turn button")
	}
	if _, ok := game.dom.FindByID("available-dice-grid"); !ok {
		t.Fatalf("expected available dice grid")
	}
	if _, ok := game.dom.FindByID("combat-log-grid"); !ok {
		t.Fatalf("expected combat log grid")
	}
	for _, id := range []string{
		"available-dice-scroll",
		"used-dice-scroll",
		"combat-log-scroll",
	} {
		if _, ok := game.dom.FindByID(id); !ok {
			t.Fatalf("expected combat scroll node %q", id)
		}
	}
	if !containsTag(game.dom.Root, ebitenui.TagScrollView) {
		t.Fatalf("expected scroll view nodes on combat screen")
	}
	assertNoLayoutOverflow(t, game)

	ui := game.uiSnapshot()
	if got, want := ui.Root.Props["currentScreen"], string(ScreenCombat); got != want {
		t.Fatalf("ui snapshot currentScreen mismatch: got %#v want %q", got, want)
	}
	for _, key := range []string{
		"partyIDs",
		"partyHP",
		"availableDiceCount",
		"graveyardDiceCount",
		"currentNodeID",
		"turn",
	} {
		if _, ok := ui.Root.Props[key]; !ok {
			t.Fatalf("expected %q in root props", key)
		}
	}
}

func TestGameOnlyButtonsKeepHoverAndFocusState(t *testing.T) {
	game := newGame(GameConfig{Seed: 7})
	game.width = DefaultWindowWidth
	game.height = DefaultWindowHeight

	if err := game.step(ebitenui.InputSnapshot{}); err != nil {
		t.Fatalf("initial step failed: %v", err)
	}
	clickAndStep(t, game, "party-option-human-warrior")
	clickAndStep(t, game, "party-option-human-guard")
	clickAndStep(t, game, "party-option-human-guide")
	clickAndStep(t, game, "start-run-button")
	clickAndStep(t, game, "map-node-normal-a")

	layout := game.currentLayout()
	scrollLayout, ok := layout.FindByID("available-dice-scroll")
	if !ok {
		t.Fatalf("expected available dice scroll layout")
	}
	firstWrap, ok := layout.FindByID("available-dice-grid-wrap-0")
	if !ok {
		t.Fatalf("expected first available die wrapper")
	}
	secondWrap, ok := layout.FindByID("available-dice-grid-wrap-1")
	if !ok {
		t.Fatalf("expected second available die wrapper")
	}
	gapY := firstWrap.Frame.Y + firstWrap.Frame.Height + ((secondWrap.Frame.Y - (firstWrap.Frame.Y + firstWrap.Frame.Height)) / 2)
	gapX := scrollLayout.Frame.X + (scrollLayout.Frame.Width / 2)

	if err := game.step(ebitenui.InputSnapshot{PointerX: gapX, PointerY: gapY}); err != nil {
		t.Fatalf("hover gap step failed: %v", err)
	}
	scrollNode, ok := game.dom.FindByID("available-dice-scroll")
	if !ok {
		t.Fatalf("expected available dice scroll node")
	}
	if scrollNode.Props.State.Hovered || scrollNode.Props.State.Pressed || scrollNode.Props.State.Focused {
		t.Fatalf("expected scroll background to stay visually inactive, got %#v", scrollNode.Props.State)
	}

	if err := game.step(ebitenui.InputSnapshot{PointerX: gapX, PointerY: gapY, PointerDown: true}); err != nil {
		t.Fatalf("background press step failed: %v", err)
	}
	if err := game.step(ebitenui.InputSnapshot{PointerX: gapX, PointerY: gapY}); err != nil {
		t.Fatalf("background release step failed: %v", err)
	}
	if got := game.runtime.FocusedID(); got != "" {
		t.Fatalf("expected background click to leave no focus, got %q", got)
	}

	game.mu.RLock()
	dieID := game.run.CurrentCombat.AvailableDice[0].ID
	game.mu.RUnlock()

	layout = game.currentLayout()
	dieLayout, ok := layout.FindByID("available-die-" + dieID)
	if !ok {
		t.Fatalf("expected available die button layout")
	}
	dieX := dieLayout.Frame.X + (dieLayout.Frame.Width / 2)
	dieY := dieLayout.Frame.Y + (dieLayout.Frame.Height / 2)

	if err := game.step(ebitenui.InputSnapshot{PointerX: dieX, PointerY: dieY}); err != nil {
		t.Fatalf("button hover step failed: %v", err)
	}
	dieNode, ok := game.dom.FindByID("available-die-" + dieID)
	if !ok {
		t.Fatalf("expected available die button node")
	}
	if !dieNode.Props.State.Hovered {
		t.Fatalf("expected button hover state to remain active")
	}
	scrollNode, ok = game.dom.FindByID("available-dice-scroll")
	if !ok {
		t.Fatalf("expected available dice scroll node after button hover")
	}
	if scrollNode.Props.State.Hovered || scrollNode.Props.State.Pressed || scrollNode.Props.State.Focused {
		t.Fatalf("expected scroll node to stay passive while button is hovered, got %#v", scrollNode.Props.State)
	}
}

func clickAndStep(t *testing.T, game *Game, nodeID string) {
	t.Helper()
	result := game.debugBridgeLikeCommand("ui_click", map[string]any{"node_id": nodeID})
	if !result.Success {
		t.Fatalf("ui_click(%q) failed: %#v", nodeID, result)
	}
	for i := 0; i < 3; i++ {
		if err := game.step(ebitenui.InputSnapshot{}); err != nil {
			t.Fatalf("step %d after click failed: %v", i, err)
		}
	}
}

func containsTag(node *ebitenui.Node, tag ebitenui.Tag) bool {
	if node == nil {
		return false
	}
	if node.Tag == tag {
		return true
	}
	for _, child := range node.Children {
		if containsTag(child, tag) {
			return true
		}
	}
	return false
}

func assertNoLayoutOverflow(t *testing.T, game *Game) {
	t.Helper()

	layout := game.currentLayout()
	if layout == nil {
		t.Fatalf("expected current layout")
	}
	report := ebitenui.ValidateLayout(layout, game.currentViewport(), ebitenui.ValidationOptions{})
	for _, issue := range report.Issues {
		if isWithinScrollView(layout, issue.NodeID) {
			continue
		}
		switch issue.Code {
		case ebitenui.IssueOutOfViewport, ebitenui.IssueOutOfParent, ebitenui.IssueTextOverflow:
			t.Fatalf("unexpected layout issue: node=%s code=%s message=%s", issue.NodeID, issue.Code, issue.Message)
		}
	}
	assertNodeOverflowFree(t, layout, false)
}

func assertNodeOverflowFree(t *testing.T, layout *ebitenui.LayoutNode, withinScrollView bool) {
	t.Helper()
	if layout == nil {
		return
	}
	if layout.Node != nil && layout.Node.Tag == ebitenui.TagScrollView {
		withinScrollView = true
	}
	if layout.Overflow.Any && !withinScrollView {
		nodeID := ""
		if layout.Node != nil {
			nodeID = layout.Node.Props.ID
		}
		t.Fatalf("unexpected overflow on node %q: %#v", nodeID, layout.Overflow)
	}
	for _, child := range layout.Children {
		assertNodeOverflowFree(t, child, withinScrollView)
	}
}

func isWithinScrollView(layout *ebitenui.LayoutNode, nodeID string) bool {
	if layout == nil || nodeID == "" {
		return false
	}
	node, ok := layout.FindByID(nodeID)
	if !ok {
		return false
	}
	for node != nil {
		if node.Node != nil && node.Node.Tag == ebitenui.TagScrollView {
			return true
		}
		if node.ParentID == "" {
			break
		}
		parent, ok := layout.FindByID(node.ParentID)
		if !ok {
			break
		}
		node = parent
	}
	return false
}
