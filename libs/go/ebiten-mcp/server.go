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
	"strconv"
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
	case "get_ui_overview":
		return server.client.get(ctx, "/debug/ui/overview")
	case "query_ui_nodes":
		if params == nil {
			params = map[string]any{}
		}
		return server.client.post(ctx, "/debug/ui/query", params)
	case "inspect_ui_node":
		nodeID, _ := params["node_id"].(string)
		if nodeID == "" {
			return nil, fmt.Errorf("inspect_ui_node requires node_id")
		}
		query := map[string]string{}
		if value, ok := params["include_children"].(bool); ok {
			query["include_children"] = strconv.FormatBool(value)
		}
		if value, ok := intParam(params, "child_depth"); ok {
			query["child_depth"] = strconv.Itoa(value)
		}
		if value, ok := params["include_props"].(bool); ok {
			query["include_props"] = strconv.FormatBool(value)
		}
		if value, ok := params["include_issues"].(bool); ok {
			query["include_issues"] = strconv.FormatBool(value)
		}
		return server.client.getWithQuery(ctx, "/debug/ui/node/"+url.PathEscape(nodeID), query)
	case "list_ui_issues":
		query := map[string]string{}
		for _, key := range []string{"severity", "code", "node_id", "cursor"} {
			if value, ok := params[key].(string); ok && strings.TrimSpace(value) != "" {
				query[key] = value
			}
		}
		if value, ok := intParam(params, "limit"); ok {
			query["limit"] = strconv.Itoa(value)
		}
		return server.client.getWithQuery(ctx, "/debug/ui/issues", query)
	case "capture_ui_screenshot":
		if params == nil {
			params = map[string]any{}
		}
		return server.client.post(ctx, "/debug/ui/capture", params)
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

func (server *Server) StreamableHTTPHandler(opts *mcp.StreamableHTTPOptions) http.Handler {
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server.sdkServer()
	}, opts)
}

func (server *Server) streamableHTTPHandler(opts *mcp.StreamableHTTPOptions) http.Handler {
	return server.StreamableHTTPHandler(opts)
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
		case "query_ui_nodes", "list_ui_issues", "capture_ui_screenshot":
			name := tool.Name
			description := tool.Description
			mcp.AddTool(sdkServer, &mcp.Tool{
				Name:        name,
				Description: description,
			}, func(ctx context.Context, req *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, any, error) {
				payload, err := server.CallTool(ctx, name, input)
				if err != nil {
					return nil, nil, err
				}
				return nil, payload, nil
			})
		case "inspect_ui_node":
			name := tool.Name
			description := tool.Description
			mcp.AddTool(sdkServer, &mcp.Tool{
				Name:        name,
				Description: description,
			}, server.inspectNodeTool)
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

type inspectNodeInput struct {
	NodeID          string `json:"node_id" jsonschema:"node id to inspect"`
	IncludeChildren *bool  `json:"include_children,omitempty" jsonschema:"include direct child summaries"`
	ChildDepth      int    `json:"child_depth,omitempty" jsonschema:"child summary depth"`
	IncludeProps    *bool  `json:"include_props,omitempty" jsonschema:"include props metadata"`
	IncludeIssues   *bool  `json:"include_issues,omitempty" jsonschema:"include issue metadata"`
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

func (server *Server) inspectNodeTool(ctx context.Context, req *mcp.CallToolRequest, input inspectNodeInput) (*mcp.CallToolResult, any, error) {
	params := map[string]any{
		"node_id": input.NodeID,
	}
	if input.IncludeChildren != nil {
		params["include_children"] = *input.IncludeChildren
	}
	if input.ChildDepth > 0 {
		params["child_depth"] = input.ChildDepth
	}
	if input.IncludeProps != nil {
		params["include_props"] = *input.IncludeProps
	}
	if input.IncludeIssues != nil {
		params["include_issues"] = *input.IncludeIssues
	}
	payload, err := server.CallTool(ctx, "inspect_ui_node", params)
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
			Description: "Return full UI tree with semantic/layout/computed/issues/inputState snapshot state. Legacy full dump with high token cost.",
		},
		{
			Name:        "get_ui_overview",
			Description: "Return compact UI overview for low-token design and layout testing.",
		},
		{
			Name:        "query_ui_nodes",
			Description: "Return paginated compact UI node summaries filtered by id, role, slot, type, text, visibility, interactivity, issue code, or viewport.",
			InputSchema: map[string]any{
				"type": "object",
			},
		},
		{
			Name:        "inspect_ui_node",
			Description: "Return compact detail for a single UI node with summary, semantic, layout, computed metadata, issues, and direct child summaries.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"node_id":          map[string]any{"type": "string"},
					"include_children": map[string]any{"type": "boolean"},
					"child_depth":      map[string]any{"type": "integer"},
					"include_props":    map[string]any{"type": "boolean"},
					"include_issues":   map[string]any{"type": "boolean"},
				},
				"required": []string{"node_id"},
			},
		},
		{
			Name:        "list_ui_issues",
			Description: "Return paginated flat UI issue rows for low-token layout validation.",
			InputSchema: map[string]any{
				"type": "object",
			},
		},
		{
			Name:        "capture_ui_screenshot",
			Description: "Capture a UI screenshot artifact and return metadata plus absolute artifact path without inlining image bytes.",
			InputSchema: map[string]any{
				"type": "object",
			},
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

func (client *BridgeClient) getWithQuery(ctx context.Context, path string, query map[string]string) (any, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	values := request.URL.Query()
	for key, value := range query {
		if strings.TrimSpace(value) == "" {
			continue
		}
		values.Set(key, value)
	}
	request.URL.RawQuery = values.Encode()
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

func intParam(params map[string]any, key string) (int, bool) {
	if params == nil {
		return 0, false
	}
	switch value := params[key].(type) {
	case int:
		return value, true
	case int32:
		return int(value), true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}
