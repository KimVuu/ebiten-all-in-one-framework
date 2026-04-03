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
	game.width = 1280
	game.height = 720

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
	if containsTag(game.dom.Root, ebitenui.TagScrollView) {
		t.Fatalf("expected no scroll view nodes on combat screen")
	}

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
