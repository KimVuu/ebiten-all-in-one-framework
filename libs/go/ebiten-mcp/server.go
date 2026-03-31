package ebitenmcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BridgeClient struct {
	baseURL string
	client  *http.Client
}

func NewBridgeClient(baseURL string) *BridgeClient {
	baseURL = strings.TrimSpace(baseURL)
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	return &BridgeClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func newBridgeClient(baseURL string) *BridgeClient {
	return NewBridgeClient(baseURL)
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"inputSchema,omitempty"`
}

type Server struct {
	client *BridgeClient
}

func New(client *BridgeClient) *Server {
	return &Server{client: client}
}

func newMCPServer(client *BridgeClient) *Server {
	return New(client)
}

func (server *Server) CallTool(ctx context.Context, name string, params map[string]any) (any, error) {
	switch name {
	case "game_health":
		return server.client.get(ctx, "/health")
	case "get_frame_state":
		return server.client.get(ctx, "/debug/frame")
	case "get_scene_state":
		return server.client.get(ctx, "/debug/scene")
	case "get_world_state":
		return server.client.get(ctx, "/debug/world")
	case "get_ui_state":
		return server.client.get(ctx, "/debug/ui")
	case "list_commands":
		return server.client.get(ctx, "/debug/commands")
	case "run_command":
		commandName, _ := params["name"].(string)
		if commandName == "" {
			return nil, fmt.Errorf("run_command requires name")
		}
		rawArgs, _ := params["args"].(map[string]any)
		if rawArgs == nil {
			rawArgs = map[string]any{}
		}
		return server.client.post(ctx, "/debug/commands/"+url.PathEscape(commandName), map[string]any{
			"args": rawArgs,
		})
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (server *Server) callTool(ctx context.Context, name string, params map[string]any) (any, error) {
	return server.CallTool(ctx, name, params)
}

func (server *Server) ServeStdio(ctx context.Context, input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var request rpcRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			if err := writeResponse(output, rpcResponse{
				JSONRPC: "2.0",
				Error: &rpcError{
					Code:    -32700,
					Message: "parse error",
				},
			}); err != nil {
				return err
			}
			continue
		}

		response, shouldWrite := server.handleRequest(ctx, request)
		if !shouldWrite {
			continue
		}
		if err := writeResponse(output, response); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (server *Server) serveStdio(ctx context.Context, input io.Reader, output io.Writer) error {
	return server.ServeStdio(ctx, input, output)
}

func (server *Server) handleRequest(ctx context.Context, request rpcRequest) (rpcResponse, bool) {
	switch request.Method {
	case "initialize":
		var params struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(request.Params, &params)

		return rpcResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: initializeResult{
				ProtocolVersion: defaultString(params.ProtocolVersion, "2024-11-05"),
				Capabilities: map[string]any{
					"tools": map[string]any{},
				},
				ServerInfo: map[string]any{
					"name":    "ebiten-mcp",
					"version": "0.1.0",
				},
			},
		}, true
	case "notifications/initialized":
		return rpcResponse{}, false
	case "tools/list":
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]any{
				"tools": server.tools(),
			},
		}, true
	case "tools/call":
		var params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(request.Params, &params); err != nil {
			return rpcResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: &rpcError{
					Code:    -32602,
					Message: "invalid params",
				},
			}, true
		}

		payload, err := server.CallTool(ctx, params.Name, params.Arguments)
		if err != nil {
			return rpcResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: &rpcError{
					Code:    -32000,
					Message: err.Error(),
				},
			}, true
		}

		return rpcResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result:  payload,
		}, true
	default:
		return rpcResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &rpcError{
				Code:    -32601,
				Message: "method not found",
			},
		}, true
	}
}

func (server *Server) tools() []Tool {
	return []Tool{
		{
			Name:        "game_health",
			Description: "Return bridge health and connection state.",
		},
		{
			Name:        "get_frame_state",
			Description: "Return current frame timing and pause state.",
		},
		{
			Name:        "get_scene_state",
			Description: "Return current and known scene state.",
		},
		{
			Name:        "get_world_state",
			Description: "Return entity snapshot state.",
		},
		{
			Name:        "get_ui_state",
			Description: "Return UI tree with semantic/layout/computed/issues/inputState snapshot state.",
		},
		{
			Name:        "list_commands",
			Description: "List registered debug commands.",
		},
		{
			Name:        "run_command",
			Description: "Invoke a registered debug command for validate/inspect/suggest/overlay/input flows such as validate_ui_layout, inspect_ui_node, suggest_ui_constraint_fixes, set_ui_debug_overlay, ui_click, ui_scroll, ui_type_text, or ui_key_event.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"args": map[string]any{
						"type": "object",
					},
				},
				"required": []string{"name"},
			},
		},
	}
}

func (client *BridgeClient) get(ctx context.Context, path string) (any, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return client.do(request)
}

func (client *BridgeClient) post(ctx context.Context, path string, body any) (any, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+path, bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return client.do(request)
}

func (client *BridgeClient) do(request *http.Request) (any, error) {
	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = response.Status
		}
		return nil, fmt.Errorf("bridge request failed: %s %s: %s", request.Method, request.URL.Path, message)
	}

	var payload any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type initializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ServerInfo      map[string]any `json:"serverInfo"`
}

func writeResponse(output io.Writer, response rpcResponse) error {
	encoded, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = output.Write(append(encoded, '\n'))
	return err
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
