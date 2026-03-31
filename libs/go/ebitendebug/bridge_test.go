package ebitendebug

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDisabledBridgeUsesNoopHandler(t *testing.T) {
	bridge := New(Config{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	bridge.Handler().ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusServiceUnavailable; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}
}

func TestBridgeExposesRegisteredSnapshots(t *testing.T) {
	bridge := New(Config{
		Enabled: true,
		GameID:  "debug-bridge",
		Version: "v1.0.0",
	})
	bridge.SetFrameProvider(func() FrameSnapshot {
		return FrameSnapshot{
			Frame:        12,
			Tick:         24,
			FPS:          60,
			TPS:          60,
			Paused:       false,
			DebugEnabled: true,
		}
	})
	bridge.SetSceneProvider(func() SceneSnapshot {
		return SceneSnapshot{
			Current: SceneRef{ID: "menu", Name: "Menu"},
			Known: []SceneRef{
				{ID: "menu", Name: "Menu"},
				{ID: "battle", Name: "Battle"},
			},
			Summary: map[string]any{
				"selection": "start",
			},
		}
	})
	bridge.SetWorldProvider(func() WorldSnapshot {
		return WorldSnapshot{
			Entities: []EntitySnapshot{
				{
					ID:      "hero",
					Type:    "player",
					Tags:    []string{"party", "selected"},
					Visible: true,
					Enabled: true,
					Position: Vector2{
						X: 32,
						Y: 48,
					},
					Size: Vector2{
						X: 16,
						Y: 24,
					},
					Props: map[string]any{
						"hp": 80,
					},
				},
			},
		}
	})
	bridge.SetUIProvider(func() UISnapshot {
		return UISnapshot{
			Width:  960,
			Height: 540,
			Viewport: UIViewportSnapshot{
				Width:  960,
				Height: 540,
				Scale:  1,
			},
			SafeArea: UISafeAreaSnapshot{
				Top:    12,
				Right:  12,
				Bottom: 12,
				Left:   12,
			},
			IssueSummary: UIIssueSummarySnapshot{
				Total:        1,
				InvalidNodes: 1,
				Errors:       1,
				Warnings:     0,
				Info:         0,
			},
			InputState: UIInputSnapshot{
				FocusedNodeID: "start-button",
				HoveredNodeID: "start-button",
				Pointer: &UIPointerSnapshot{
					X:    120,
					Y:    40,
					Down: false,
				},
			},
			Root: UINodeSnapshot{
				ID:       "screen",
				Type:     "screen",
				Visible:  true,
				ParentID: "",
				Semantic: &UISemanticSnapshot{
					Screen:  "main_menu",
					Element: "root",
					Role:    "screen",
					Slot:    "root",
				},
				Layout: &UILayoutSnapshot{
					Mode:     LayoutModeStack,
					Anchor:   "center",
					Pivot:    "center",
					ParentID: "",
					Offset:   UIPositionSnapshot{X: 0, Y: 0},
					Size:     UISizeSnapshot{Width: 960, Height: 540},
					Constraints: []UIConstraintSnapshot{
						{Field: "keep_inside_parent", Op: "set", Value: true},
					},
				},
				Computed: &UIComputedSnapshot{
					Bounds: Rect{
						X:      0,
						Y:      0,
						Width:  960,
						Height: 540,
					},
					ContentBounds: &Rect{
						X:      0,
						Y:      0,
						Width:  960,
						Height: 540,
					},
					ClickableRect: &Rect{
						X:      0,
						Y:      0,
						Width:  960,
						Height: 540,
					},
					ClipRect: &Rect{
						X:      0,
						Y:      0,
						Width:  960,
						Height: 540,
					},
					Visible: true,
					Overflow: &UIOverflowSnapshot{
						Top:    false,
						Right:  false,
						Bottom: false,
						Left:   false,
					},
				},
				Issues: []UIIssueSnapshot{
					{
						NodeID:   "start-button",
						Severity: "warning",
						Code:     "min_hit_target",
						Message:  "button too small",
						SuggestedConstraintChanges: []UIConstraintSnapshot{
							{Field: "minWidth", Op: "set", Value: 180},
						},
					},
				},
				Children: []UINodeSnapshot{
					{
						ID:      "title",
						Type:    "text",
						Text:    "Debug Bridge Example",
						Visible: true,
						Bounds: Rect{
							X:      12,
							Y:      16,
							Width:  220,
							Height: 16,
						},
					},
				},
			},
		}
	})
	bridge.RegisterCommand("pause.toggle", func(request CommandRequest) CommandResult {
		if request.Name != "pause.toggle" {
			t.Fatalf("command name mismatch: got %q", request.Name)
		}
		return CommandResult{
			Success: true,
			Message: "paused",
		}
	})

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "health", path: "/health", want: `"gameId":"debug-bridge"`},
		{name: "frame", path: "/debug/frame", want: `"frame":12`},
		{name: "scene", path: "/debug/scene", want: `"current":{"id":"menu","name":"Menu"}`},
		{name: "world", path: "/debug/world", want: `"id":"hero"`},
		{name: "ui", path: "/debug/ui", want: `"semantic":{"screen":"main_menu"`},
		{name: "commands", path: "/debug/commands", want: `"pause.toggle"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			bridge.Handler().ServeHTTP(recorder, request)

			if got, want := recorder.Code, http.StatusOK; got != want {
				t.Fatalf("status mismatch: got %d want %d", got, want)
			}
			if !bytes.Contains(recorder.Body.Bytes(), []byte(test.want)) {
				t.Fatalf("expected response to contain %q, got %s", test.want, recorder.Body.String())
			}
		})
	}
}

func TestBridgeReturnsNotFoundForUnknownCommand(t *testing.T) {
	bridge := New(Config{Enabled: true})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/debug/commands/missing", bytes.NewBufferString(`{"args":{"enabled":true}}`))
	bridge.Handler().ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusNotFound; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}
}

func TestBridgeExecutesCommands(t *testing.T) {
	bridge := New(Config{Enabled: true})
	bridge.RegisterCommand("scene.switch", func(request CommandRequest) CommandResult {
		if got, want := request.Args["scene"], "battle"; got != want {
			t.Fatalf("arg mismatch: got %v want %v", got, want)
		}
		return CommandResult{
			Success:        true,
			Message:        "switched",
			Status:         "completed",
			ResolvedTarget: "scene:battle",
			QueuedFrame:    42,
			Payload: map[string]any{
				"scene": "battle",
			},
		}
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/debug/commands/scene.switch", bytes.NewBufferString(`{"args":{"scene":"battle"}}`))
	bridge.Handler().ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"success":true`)) {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"completed"`)) {
		t.Fatalf("expected status in response, got %s", recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"resolvedTarget":"scene:battle"`)) {
		t.Fatalf("expected resolvedTarget in response, got %s", recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"scene":"battle"`)) {
		t.Fatalf("expected payload in response, got %s", recorder.Body.String())
	}
}

func TestCommandResultEncodesInputMetadata(t *testing.T) {
	payload, err := json.Marshal(CommandResult{
		Success:        true,
		Status:         "queued",
		ResolvedTarget: "start-button",
		QueuedFrame:    18,
		Reason:         "queued for next frame",
	})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	for _, want := range []string{
		`"status":"queued"`,
		`"resolvedTarget":"start-button"`,
		`"queuedFrame":18`,
		`"reason":"queued for next frame"`,
	} {
		if !bytes.Contains(payload, []byte(want)) {
			t.Fatalf("expected payload to contain %q, got %s", want, string(payload))
		}
	}
}

func TestBridgeStartAndCloseLifecycle(t *testing.T) {
	bridge := New(Config{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	})

	if err := bridge.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if bridge.Address() == "" {
		t.Fatalf("expected bound address")
	}

	response, err := http.Get("http://" + bridge.Address() + "/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer response.Body.Close()

	if got, want := response.StatusCode, http.StatusOK; got != want {
		t.Fatalf("status mismatch: got %d want %d", got, want)
	}

	if err := bridge.Close(context.Background()); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}

func TestBridgeRejectsNonLoopbackAddresses(t *testing.T) {
	bridge := New(Config{
		Enabled: true,
		Addr:    "0.0.0.0:3456",
	})

	if err := bridge.Start(); err == nil {
		t.Fatalf("expected start to reject non-loopback address")
	}
}

func TestDecodeCommandRequestWithoutArgs(t *testing.T) {
	var request CommandRequest
	if err := json.Unmarshal([]byte(`{}`), &request); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if request.Args == nil {
		t.Fatalf("expected args map to be initialized")
	}
}
