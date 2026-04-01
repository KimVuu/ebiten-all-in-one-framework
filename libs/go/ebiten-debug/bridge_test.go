package ebitendebug

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestBridgeExposesCompactUIEndpoints(t *testing.T) {
	bridge := New(Config{Enabled: true})
	bridge.SetUIOverviewProvider(func() UIOverviewSnapshot {
		return UIOverviewSnapshot{
			Viewport:         UIViewportSnapshot{Width: 1280, Height: 720, Scale: 1},
			RootID:           "showcase-root",
			FocusedNodeID:    "name-input",
			HoveredNodeID:    "name-input",
			TotalNodeCount:   320,
			VisibleNodeCount: 55,
			InvalidNodeCount: 8,
			IssueSummary:     UIIssueSummarySnapshot{Total: 8, Errors: 5, Warnings: 3, InvalidNodes: 8},
			TopLevelSections: []UINodeSummarySnapshot{
				{ID: "showcase-header", Type: "header", Role: "header", Bounds: Rect{Width: 1280, Height: 96}},
			},
		}
	})
	bridge.SetUIQueryProvider(func(request UIQueryRequest) UIQueryResult {
		return UIQueryResult{
			Nodes: []UINodeSummarySnapshot{
				{ID: "name-input", Type: "input", Role: "input", TextPreview: "Kim", Interactive: true},
			},
			NextCursor: "1",
			Total:      1,
		}
	})
	bridge.SetUINodeProvider(func(request UINodeInspectRequest) (UINodeDetailSnapshot, bool) {
		if request.NodeID != "name-input" {
			return UINodeDetailSnapshot{}, false
		}
		return UINodeDetailSnapshot{
			Summary: UINodeSummarySnapshot{
				ID:          "name-input",
				Type:        "input",
				Role:        "input",
				TextPreview: "Kim",
				Interactive: true,
			},
			Semantic: &UISemanticSnapshot{Screen: "showcase", Element: "name-input", Role: "input", Slot: "input"},
			Children: []UINodeSummarySnapshot{
				{ID: "name-input-label", Type: "text", Role: "text"},
			},
		}, true
	})
	bridge.SetUIIssuesProvider(func(request UIIssueListRequest) UIIssueListSnapshot {
		return UIIssueListSnapshot{
			IssueSummary: UIIssueSummarySnapshot{Total: 2, Errors: 1, Warnings: 1, InvalidNodes: 2},
			Issues: []UIIssueSnapshot{
				{NodeID: "hero-title", Severity: "error", Code: "out_of_parent", Message: "node extends beyond parent bounds"},
			},
			NextCursor: "1",
			Total:      2,
		}
	})

	dir := t.TempDir()
	path := filepath.Join(dir, "capture.png")
	if err := os.WriteFile(path, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, 0o644); err != nil {
		t.Fatalf("write artifact failed: %v", err)
	}
	bridge.SetUICaptureProvider(func(request UICaptureRequest) (UICaptureResult, bool) {
		target := request.Target
		if target == "" {
			target = "viewport"
		}
		return UICaptureResult{
			ArtifactID:     "artifact-1",
			Path:           path,
			Width:          1280,
			Height:         720,
			Hash:           "abc123",
			OverlayEnabled: request.WithOverlay,
			Target:         target,
			CapturedRect:   Rect{X: 0, Y: 0, Width: 1280, Height: 720},
			ContentType:    "image/png",
		}, true
	})
	bridge.SetUIArtifactProvider(func(id string) (UIArtifact, bool) {
		if id != "artifact-1" {
			return UIArtifact{}, false
		}
		return UIArtifact{
			ID:          id,
			Path:        path,
			ContentType: "image/png",
		}, true
	})

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   string
	}{
		{name: "overview", method: http.MethodGet, path: "/debug/ui/overview", want: `"rootId":"showcase-root"`},
		{name: "query", method: http.MethodPost, path: "/debug/ui/query", body: `{"visible_only":true,"limit":10}`, want: `"nodes":[{"id":"name-input"`},
		{name: "node", method: http.MethodGet, path: "/debug/ui/node/name-input", want: `"summary":{"id":"name-input"`},
		{name: "issues", method: http.MethodGet, path: "/debug/ui/issues?limit=10", want: `"issues":[{"nodeId":"hero-title"`},
		{name: "capture", method: http.MethodPost, path: "/debug/ui/capture", body: `{"target":"viewport","with_overlay":true}`, want: `"artifactId":"artifact-1"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.path, bytes.NewBufferString(test.body))
			if test.body != "" {
				request.Header.Set("Content-Type", "application/json")
			}
			bridge.Handler().ServeHTTP(recorder, request)
			if got, want := recorder.Code, http.StatusOK; got != want {
				t.Fatalf("status mismatch: got %d want %d", got, want)
			}
			if !bytes.Contains(recorder.Body.Bytes(), []byte(test.want)) {
				t.Fatalf("expected response to contain %q, got %s", test.want, recorder.Body.String())
			}
		})
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/ui/artifacts/artifact-1", nil)
	bridge.Handler().ServeHTTP(recorder, request)
	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Fatalf("artifact status mismatch: got %d want %d", got, want)
	}
	if got, want := recorder.Header().Get("Content-Type"), "image/png"; got != want {
		t.Fatalf("artifact content type mismatch: got %q want %q", got, want)
	}
	if !bytes.HasPrefix(recorder.Body.Bytes(), []byte{0x89, 'P', 'N', 'G'}) {
		t.Fatalf("expected png bytes, got %x", recorder.Body.Bytes())
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
