package ebitenmcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		case "run_command":
			runCommandDesc = tool.Description
		}
	}

	if !strings.Contains(uiToolDesc, "semantic/layout/computed/issues/inputState") {
		t.Fatalf("expected ui tool description to mention expanded snapshot, got %q", uiToolDesc)
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

func TestServeStdioHandlesInitializeListAndCall(t *testing.T) {
	bridge := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/health":
			_, _ = writer.Write([]byte(`{"gameId":"debug-bridge","version":"v1","connected":true}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer bridge.Close()

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"game_health","arguments":{}}}`,
	}, "\n") + "\n"

	output := &bytes.Buffer{}
	server := New(NewBridgeClient(bridge.URL))
	if err := server.ServeStdio(context.Background(), strings.NewReader(input), output); err != nil {
		t.Fatalf("serveStdio failed: %v", err)
	}

	lines := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(output.Bytes()))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines) < 3 {
		t.Fatalf("expected at least three responses, got %d", len(lines))
	}

	if !strings.Contains(lines[0], `"result":{"protocolVersion"`) {
		t.Fatalf("expected initialize response, got %s", lines[0])
	}
	if !strings.Contains(lines[1], `"tools":[`) {
		t.Fatalf("expected tools/list response, got %s", lines[1])
	}
	if !strings.Contains(lines[2], `"gameId":"debug-bridge"`) {
		t.Fatalf("expected game health payload, got %s", lines[2])
	}
}
