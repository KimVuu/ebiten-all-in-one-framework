package ebitenmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	session, err := server.sdkServer().Connect(ctx, &mcp.IOTransport{
		Reader: io.NopCloser(input),
		Writer: nopWriteCloser{Writer: output},
	}, nil)
	if err != nil {
		return err
	}
	err = session.Wait()
	if isBenignStdioClose(err) {
		return nil
	}
	return err
}

func (server *Server) serveStdio(ctx context.Context, input io.Reader, output io.Writer) error {
	return server.ServeStdio(ctx, input, output)
}

func (server *Server) sdkServer() *mcp.Server {
	sdkServer := mcp.NewServer(&mcp.Implementation{
		Name:    "ebiten-mcp",
		Version: "0.1.0",
	}, nil)

	for _, tool := range server.tools() {
		switch tool.Name {
		case "run_command":
			mcp.AddTool(sdkServer, &mcp.Tool{
				Name:        tool.Name,
				Description: tool.Description,
			}, server.runCommandTool)
		default:
			name := tool.Name
			description := tool.Description
			mcp.AddTool(sdkServer, &mcp.Tool{
				Name:        name,
				Description: description,
			}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
				payload, err := server.CallTool(ctx, name, nil)
				if err != nil {
					return nil, nil, err
				}
				return nil, payload, nil
			})
		}
	}

	return sdkServer
}

type runCommandInput struct {
	Name string         `json:"name" jsonschema:"registered debug command name"`
	Args map[string]any `json:"args,omitempty" jsonschema:"command arguments"`
}

func (server *Server) runCommandTool(ctx context.Context, req *mcp.CallToolRequest, input runCommandInput) (*mcp.CallToolResult, any, error) {
	payload, err := server.CallTool(ctx, "run_command", map[string]any{
		"name": input.Name,
		"args": input.Args,
	})
	if err != nil {
		return nil, nil, err
	}
	return nil, payload, nil
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

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func isBenignStdioClose(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	message := err.Error()
	return strings.Contains(message, "EOF") || strings.Contains(message, "server is closing")
}
