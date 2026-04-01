package ebitenmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolHandlersCallExpectedBridgeEndpoints(t *testing.T) {
	var paths []string
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		paths = append(paths, request.Method+" "+request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/health":
			_, _ = writer.Write([]byte(`{"gameId":"debug-bridge","version":"v1","connected":true}`))
		case "/debug/frame":
			_, _ = writer.Write([]byte(`{"frame":12,"tick":30,"fps":60,"tps":60,"paused":false,"debugEnabled":true}`))
		case "/debug/scene":
			_, _ = writer.Write([]byte(`{"current":{"id":"menu","name":"Menu"},"known":[{"id":"menu","name":"Menu"}],"summary":{"selection":"start"}}`))
		case "/debug/world":
			_, _ = writer.Write([]byte(`{"entities":[{"id":"hero","type":"player","visible":true,"enabled":true,"tags":["party"],"position":{"x":1,"y":2},"size":{"x":3,"y":4},"props":{"hp":90}}]}`))
		case "/debug/ui":
			_, _ = writer.Write([]byte(`{"width":960,"height":540,"viewport":{"width":960,"height":540},"safeArea":{"top":12,"right":12,"bottom":12,"left":12},"issueSummary":{"total":1,"errors":1,"warnings":0,"info":0,"invalidNodes":1},"inputState":{"focusedNodeId":"start-button","hoveredNodeId":"start-button"},"root":{"id":"screen","type":"screen","visible":true,"semantic":{"screen":"main_menu","element":"root","role":"screen","slot":"root"},"layout":{"mode":"stack","anchor":"center","pivot":"center","offset":{"x":0,"y":0},"size":{"width":960,"height":540},"constraints":[{"field":"keep_inside_parent","op":"set","value":true}]},"computed":{"bounds":{"x":0,"y":0,"width":960,"height":540},"visible":true},"issues":[{"nodeId":"start-button","severity":"warning","code":"min_hit_target","message":"button too small","suggestedConstraintChanges":[{"field":"minWidth","op":"set","value":180}]}],"children":[{"id":"title","type":"text","text":"Debug Bridge Example","visible":true,"bounds":{"x":12,"y":16,"width":220,"height":16}}]}}`))
		case "/debug/ui/overview":
			_, _ = writer.Write([]byte(`{"rootId":"screen","totalNodeCount":42,"visibleNodeCount":12,"invalidNodeCount":1,"topLevelSections":[{"id":"showcase-header","type":"header","role":"header","bounds":{"x":0,"y":0,"width":960,"height":88}}]}`))
		case "/debug/ui/query":
			_, _ = writer.Write([]byte(`{"nodes":[{"id":"name-input","type":"input","role":"input","textPreview":"Kim","interactive":true,"scrollable":false,"bounds":{"x":40,"y":80,"width":220,"height":40}}],"nextCursor":"1","total":1}`))
		case "/debug/ui/node/name-input":
			_, _ = writer.Write([]byte(`{"summary":{"id":"name-input","type":"input","role":"input","textPreview":"Kim","interactive":true,"scrollable":false,"bounds":{"x":40,"y":80,"width":220,"height":40}},"children":[{"id":"name-input-label","type":"text","role":"text","bounds":{"x":40,"y":56,"width":100,"height":16}}]}`))
		case "/debug/ui/issues":
			_, _ = writer.Write([]byte(`{"issueSummary":{"total":2,"errors":1,"warnings":1,"info":0,"invalidNodes":2},"issues":[{"nodeId":"hero-title","severity":"error","code":"out_of_parent","message":"node extends beyond parent bounds"}],"nextCursor":"1","total":2}`))
		case "/debug/ui/capture":
			_, _ = writer.Write([]byte(`{"artifactId":"artifact-1","path":"/repo/screenshots/ebiten-ui-showcase/artifact-1.png","width":1280,"height":720,"hash":"abc123","overlayEnabled":false,"target":"viewport","capturedRect":{"x":0,"y":0,"width":1280,"height":720}}`))
		case "/debug/commands":
			_, _ = writer.Write([]byte(`{"commands":[{"name":"pause.toggle"},{"name":"scene.switch"}]}`))
		case "/debug/commands/scene.switch":
			_, _ = writer.Write([]byte(`{"success":true,"message":"ok","payload":{"scene":"battle"}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer bridge.Close()

	client := NewBridgeClient(bridge.URL)
	server := New(client)

	cases := []struct {
		name   string
		params map[string]any
		want   string
	}{
		{name: "game_health", want: `"gameId":"debug-bridge"`},
		{name: "get_frame_state", want: `"frame":12`},
		{name: "get_scene_state", want: `"selection":"start"`},
		{name: "get_world_state", want: `"id":"hero"`},
		{name: "get_ui_state", want: `"text":"Debug Bridge Example"`},
		{name: "get_ui_overview", want: `"rootId":"screen"`},
		{name: "query_ui_nodes", params: map[string]any{"visible_only": true, "limit": 10}, want: `"id":"name-input"`},
		{name: "inspect_ui_node", params: map[string]any{"node_id": "name-input"}, want: `"id":"name-input"`},
		{name: "list_ui_issues", params: map[string]any{"limit": 10}, want: `"nodeId":"hero-title"`},
		{name: "capture_ui_screenshot", params: map[string]any{"target": "viewport"}, want: `"artifactId":"artifact-1"`},
		{name: "list_commands", want: `"pause.toggle"`},
		{name: "run_command", params: map[string]any{"name": "scene.switch", "args": map[string]any{"scene": "battle"}}, want: `"scene":"battle"`},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			payload, err := server.CallTool(context.Background(), test.name, test.params)
			if err != nil {
				t.Fatalf("tool call failed: %v", err)
			}
			encoded, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if !bytes.Contains(encoded, []byte(test.want)) {
				t.Fatalf("expected response to contain %q, got %s", test.want, string(encoded))
			}
		})
	}

	got := strings.Join(paths, ",")
	for _, want := range []string{
		"GET /health",
		"GET /debug/frame",
		"GET /debug/scene",
		"GET /debug/world",
		"GET /debug/ui",
		"GET /debug/ui/overview",
		"POST /debug/ui/query",
		"GET /debug/ui/node/name-input",
		"GET /debug/ui/issues",
		"POST /debug/ui/capture",
		"GET /debug/commands",
		"POST /debug/commands/scene.switch",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected path list to contain %q, got %s", want, got)
		}
	}
}

func TestToolDefinitionsDescribeUiInspectAndInputSurface(t *testing.T) {
	server := New(NewBridgeClient("http://127.0.0.1:9999"))
	tools := server.tools()

	var uiToolDesc, runCommandDesc string
	for _, tool := range tools {
		switch tool.Name {
		case "get_ui_state":
			uiToolDesc = tool.Description
		case "capture_ui_screenshot":
			uiToolDesc += " " + tool.Description
		case "run_command":
			runCommandDesc = tool.Description
		}
	}

	if !strings.Contains(uiToolDesc, "semantic/layout/computed/issues/inputState") {
		t.Fatalf("expected ui tool description to mention expanded snapshot, got %q", uiToolDesc)
	}
	if !strings.Contains(uiToolDesc, "artifact") {
		t.Fatalf("expected compact capture tool description to mention artifact output, got %q", uiToolDesc)
	}
	if !strings.Contains(runCommandDesc, "validate/inspect/suggest/overlay/input") {
		t.Fatalf("expected run_command description to mention ui commands, got %q", runCommandDesc)
	}
}

func TestToolHandlersReturnBridgeErrors(t *testing.T) {
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, `{"error":"boom"}`, http.StatusBadGateway)
	}))
	defer bridge.Close()

	server := New(NewBridgeClient(bridge.URL))
	if _, err := server.CallTool(context.Background(), "get_world_state", nil); err == nil {
		t.Fatalf("expected bridge error")
	}
}

func TestToolHandlersReturnDecodeErrors(t *testing.T) {
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"broken"`))
	}))
	defer bridge.Close()

	server := New(NewBridgeClient(bridge.URL))
	if _, err := server.CallTool(context.Background(), "game_health", nil); err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestUnknownToolReturnsError(t *testing.T) {
	server := New(NewBridgeClient("http://127.0.0.1:9999"))
	if _, err := server.CallTool(context.Background(), "missing_tool", nil); err == nil {
		t.Fatalf("expected unknown tool error")
	}
}

func TestSDKServerListsAndCallsToolsOverInMemoryTransport(t *testing.T) {
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/debug/commands/ui_click":
			_, _ = writer.Write([]byte(`{"success":true,"status":"queued","resolvedTarget":"start-button","queuedFrame":42}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer bridge.Close()

	server := New(NewBridgeClient(bridge.URL))
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.sdkServer().Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect failed: %v", err)
	}
	defer serverSession.Close()

	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}
	if len(tools.Tools) == 0 {
		t.Fatalf("expected tools from sdk server")
	}

	result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "run_command",
		Arguments: map[string]any{"name": "ui_click", "args": map[string]any{"node_id": "start-button"}},
	})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if !bytes.Contains(encoded, []byte(`"resolvedTarget":"start-button"`)) {
		t.Fatalf("expected resolved target in sdk result, got %s", string(encoded))
	}
}

func TestStreamableHTTPHandlerServesTools(t *testing.T) {
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/health":
			_, _ = writer.Write([]byte(`{"gameId":"debug-bridge","version":"v1","connected":true}`))
		case "/debug/commands/ui_click":
			_, _ = writer.Write([]byte(`{"success":true,"status":"queued","resolvedTarget":"start-button","queuedFrame":42}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer bridge.Close()

	server := New(NewBridgeClient(bridge.URL))
	handler := server.StreamableHTTPHandler(&mcp.StreamableHTTPOptions{
		JSONResponse:               true,
		DisableLocalhostProtection: true,
	})
	httpServer := httptest.NewServer(handler)
	defer httpServer.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint: httpServer.URL,
	}, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}
	if len(tools.Tools) == 0 {
		t.Fatalf("expected tools from streamable http server")
	}

	health, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "game_health"})
	if err != nil {
		t.Fatalf("game_health call failed: %v", err)
	}
	healthJSON, err := json.Marshal(health)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if !bytes.Contains(healthJSON, []byte(`"gameId":"debug-bridge"`)) {
		t.Fatalf("expected game health payload, got %s", string(healthJSON))
	}

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "run_command",
		Arguments: map[string]any{"name": "ui_click", "args": map[string]any{"node_id": "start-button"}},
	})
	if err != nil {
		t.Fatalf("run_command call failed: %v", err)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if !bytes.Contains(encoded, []byte(`"resolvedTarget":"start-button"`)) {
		t.Fatalf("expected resolved target in streamable result, got %s", string(encoded))
	}
}
