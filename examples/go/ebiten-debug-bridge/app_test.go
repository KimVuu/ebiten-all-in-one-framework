package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
)

func TestBuildGameStateIncludesScenesAndEntities(t *testing.T) {
	game := newGame(false)

	frame := game.frameSnapshot()
	if !frame.DebugEnabled {
		t.Fatalf("expected debug flag in frame snapshot")
	}

	scene := game.sceneSnapshot()
	if got, want := scene.Current.ID, "menu"; got != want {
		t.Fatalf("scene mismatch: got %q want %q", got, want)
	}

	world := game.worldSnapshot()
	if len(world.Entities) < 2 {
		t.Fatalf("expected multiple entities, got %d", len(world.Entities))
	}
	if world.Entities[0].ID == "" {
		t.Fatalf("expected entity id")
	}

	ui := game.uiSnapshot()
	if got, want := ui.Root.ID, "screen-root"; got != want {
		t.Fatalf("ui root mismatch: got %q want %q", got, want)
	}
	if len(ui.Root.Children) < 4 {
		t.Fatalf("expected ui child nodes, got %d", len(ui.Root.Children))
	}
	if got, want := ui.Root.Children[0].Text, "Debug Bridge Example"; got != want {
		t.Fatalf("ui title mismatch: got %q want %q", got, want)
	}
}

func TestDebugCommandsMutateGameState(t *testing.T) {
	game := newGame(false)
	bridge := ebitendebug.New(ebitendebug.Config{Enabled: true})
	game.registerDebugProviders(bridge)

	result := bridge.InvokeCommand("pause.toggle", map[string]any{})
	if !result.Success || !game.paused {
		t.Fatalf("expected pause toggle to pause game, got %#v", result)
	}

	result = bridge.InvokeCommand("scene.switch", map[string]any{"scene": "battle"})
	if !result.Success || game.currentScene != "battle" {
		t.Fatalf("expected scene switch, got %#v current=%q", result, game.currentScene)
	}

	result = bridge.InvokeCommand("entity.visibility.toggle", map[string]any{"entity": "npc-guide"})
	if !result.Success {
		t.Fatalf("expected entity visibility toggle, got %#v", result)
	}
	entity := game.entityByID("npc-guide")
	if entity == nil || entity.Visible {
		t.Fatalf("expected npc-guide visibility to toggle off")
	}
}

func TestDebugBridgeStartsWhenEnabled(t *testing.T) {
	game := newGame(true)
	if err := game.startDebugBridge("127.0.0.1:0"); err != nil {
		t.Fatalf("startDebugBridge failed: %v", err)
	}
	defer func() {
		_ = game.debugBridge.Close(context.Background())
	}()

	client := http.Client{Timeout: time.Second}
	response, err := client.Get("http://" + game.debugBridge.Address() + "/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer response.Body.Close()

	if got, want := response.StatusCode, http.StatusOK; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}
}
