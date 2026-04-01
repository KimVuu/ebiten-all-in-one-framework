package ebitendebug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Enabled bool
	GameID  string
	Version string
	Addr    string
}

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type FrameSnapshot struct {
	Frame        int     `json:"frame"`
	Tick         int     `json:"tick"`
	FPS          float64 `json:"fps"`
	TPS          float64 `json:"tps"`
	Paused       bool    `json:"paused"`
	DebugEnabled bool    `json:"debugEnabled"`
}

type SceneRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SceneSnapshot struct {
	Current SceneRef       `json:"current"`
	Known   []SceneRef     `json:"known"`
	Summary map[string]any `json:"summary,omitempty"`
}

type EntitySnapshot struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Visible  bool           `json:"visible"`
	Enabled  bool           `json:"enabled"`
	Tags     []string       `json:"tags,omitempty"`
	Position Vector2        `json:"position"`
	Size     Vector2        `json:"size"`
	Props    map[string]any `json:"props,omitempty"`
}

type WorldSnapshot struct {
	Entities []EntitySnapshot `json:"entities"`
}

type CommandRequest struct {
	Name string         `json:"-"`
	Args map[string]any `json:"args"`
}

func (request *CommandRequest) UnmarshalJSON(data []byte) error {
	type rawCommandRequest struct {
		Args map[string]any `json:"args"`
	}

	var raw rawCommandRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	request.Args = raw.Args
	if request.Args == nil {
		request.Args = map[string]any{}
	}

	return nil
}

type CommandResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message,omitempty"`
	Status         string `json:"status,omitempty"`
	ResolvedTarget string `json:"resolvedTarget,omitempty"`
	QueuedFrame    int    `json:"queuedFrame,omitempty"`
	Reason         string `json:"reason,omitempty"`
	Payload        any    `json:"payload,omitempty"`
}

type CommandHandler func(CommandRequest) CommandResult

type commandDescriptor struct {
	Name string `json:"name"`
}

type Bridge struct {
	config             Config
	frameProvider      func() FrameSnapshot
	sceneProvider      func() SceneSnapshot
	worldProvider      func() WorldSnapshot
	uiProvider         func() UISnapshot
	uiOverviewProvider func() UIOverviewSnapshot
	uiQueryProvider    func(UIQueryRequest) UIQueryResult
	uiNodeProvider     func(UINodeInspectRequest) (UINodeDetailSnapshot, bool)
	uiIssuesProvider   func(UIIssueListRequest) UIIssueListSnapshot
	uiCaptureProvider  func(UICaptureRequest) (UICaptureResult, bool)
	uiArtifactProvider func(string) (UIArtifact, bool)

	mu       sync.RWMutex
	commands map[string]CommandHandler
	handler  http.Handler
	server   *http.Server
	listener net.Listener
}

func New(config Config) *Bridge {
	bridge := &Bridge{
		config:   config,
		commands: map[string]CommandHandler{},
	}
	bridge.handler = bridge.buildHandler()
	return bridge
}

func (bridge *Bridge) Handler() http.Handler {
	return bridge.handler
}

func (bridge *Bridge) SetFrameProvider(provider func() FrameSnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.frameProvider = provider
}

func (bridge *Bridge) SetSceneProvider(provider func() SceneSnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.sceneProvider = provider
}

func (bridge *Bridge) SetWorldProvider(provider func() WorldSnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.worldProvider = provider
}

func (bridge *Bridge) SetUIProvider(provider func() UISnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiProvider = provider
}

func (bridge *Bridge) SetUIOverviewProvider(provider func() UIOverviewSnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiOverviewProvider = provider
}

func (bridge *Bridge) SetUIQueryProvider(provider func(UIQueryRequest) UIQueryResult) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiQueryProvider = provider
}

func (bridge *Bridge) SetUINodeProvider(provider func(UINodeInspectRequest) (UINodeDetailSnapshot, bool)) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiNodeProvider = provider
}

func (bridge *Bridge) SetUIIssuesProvider(provider func(UIIssueListRequest) UIIssueListSnapshot) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiIssuesProvider = provider
}

func (bridge *Bridge) SetUICaptureProvider(provider func(UICaptureRequest) (UICaptureResult, bool)) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiCaptureProvider = provider
}

func (bridge *Bridge) SetUIArtifactProvider(provider func(string) (UIArtifact, bool)) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.uiArtifactProvider = provider
}

func (bridge *Bridge) RegisterCommand(name string, handler CommandHandler) {
	bridge.mu.Lock()
	defer bridge.mu.Unlock()
	bridge.commands[name] = handler
}

func (bridge *Bridge) InvokeCommand(name string, args map[string]any) CommandResult {
	bridge.mu.RLock()
	handler, ok := bridge.commands[name]
	bridge.mu.RUnlock()
	if !ok {
		return CommandResult{
			Success: false,
			Message: fmt.Sprintf("unknown command: %s", name),
		}
	}

	if args == nil {
		args = map[string]any{}
	}

	return handler(CommandRequest{
		Name: name,
		Args: args,
	})
}

func (bridge *Bridge) Start() error {
	if !bridge.config.Enabled {
		return nil
	}
	if bridge.listener != nil {
		return nil
	}

	addr := bridge.config.Addr
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	if err := validateLoopbackAddress(addr); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	bridge.listener = listener
	bridge.server = &http.Server{Handler: bridge.Handler()}

	go func() {
		_ = bridge.server.Serve(listener)
	}()

	return nil
}

func (bridge *Bridge) Address() string {
	if bridge.listener == nil {
		return ""
	}
	return bridge.listener.Addr().String()
}

func (bridge *Bridge) Close(ctx context.Context) error {
	if bridge.server == nil {
		return nil
	}

	err := bridge.server.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (bridge *Bridge) buildHandler() http.Handler {
	if !bridge.config.Enabled {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writeJSON(writer, http.StatusServiceUnavailable, map[string]any{
				"error": "debug bridge disabled",
			})
		})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{
			"gameId":    defaultString(bridge.config.GameID, "ebiten-app"),
			"version":   defaultString(bridge.config.Version, "dev"),
			"connected": true,
		})
	})
	mux.HandleFunc("/debug/frame", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.frameSnapshot())
	})
	mux.HandleFunc("/debug/scene", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.sceneSnapshot())
	})
	mux.HandleFunc("/debug/world", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.worldSnapshot())
	})
	mux.HandleFunc("/debug/ui", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.uiSnapshot())
	})
	mux.HandleFunc("/debug/ui/overview", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.uiOverview())
	})
	mux.HandleFunc("/debug/ui/query", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		var query UIQueryRequest
		if err := json.NewDecoder(request.Body).Decode(&query); err != nil && !errors.Is(err, io.EOF) {
			writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "invalid query request"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.uiQuery(query))
	})
	mux.HandleFunc("/debug/ui/node/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		nodeID := strings.TrimPrefix(request.URL.Path, "/debug/ui/node/")
		if nodeID == "" {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "missing node id"})
			return
		}
		detail, ok := bridge.uiNode(UINodeInspectRequest{
			NodeID:          nodeID,
			IncludeChildren: defaultBoolQuery(request, "include_children", true),
			ChildDepth:      defaultIntQuery(request, "child_depth", 1),
			IncludeProps:    boolQuery(request, "include_props"),
			IncludeIssues:   defaultBoolQuery(request, "include_issues", true),
		})
		if !ok {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "unknown node"})
			return
		}
		writeJSON(writer, http.StatusOK, detail)
	})
	mux.HandleFunc("/debug/ui/issues", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, bridge.uiIssues(UIIssueListRequest{
			Severity: request.URL.Query().Get("severity"),
			Code:     request.URL.Query().Get("code"),
			NodeID:   request.URL.Query().Get("node_id"),
			Limit:    defaultIntQuery(request, "limit", 50),
			Cursor:   request.URL.Query().Get("cursor"),
		}))
	})
	mux.HandleFunc("/debug/ui/capture", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		var capture UICaptureRequest
		if err := json.NewDecoder(request.Body).Decode(&capture); err != nil && !errors.Is(err, io.EOF) {
			writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "invalid capture request"})
			return
		}
		if capture.Scale == 0 {
			capture.Scale = 1
		}
		result, ok := bridge.uiCapture(capture)
		if !ok {
			writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "capture unavailable"})
			return
		}
		writeJSON(writer, http.StatusOK, result)
	})
	mux.HandleFunc("/debug/ui/artifacts/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		artifactID := strings.TrimPrefix(request.URL.Path, "/debug/ui/artifacts/")
		if artifactID == "" {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "missing artifact id"})
			return
		}
		artifact, ok := bridge.uiArtifact(artifactID)
		if !ok {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "unknown artifact"})
			return
		}
		if artifact.ContentType != "" {
			writer.Header().Set("Content-Type", artifact.ContentType)
		}
		http.ServeFile(writer, request, artifact.Path)
	})
	mux.HandleFunc("/debug/commands", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{
			"commands": bridge.commandList(),
		})
	})
	mux.HandleFunc("/debug/commands/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSON(writer, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}

		name := strings.TrimPrefix(request.URL.Path, "/debug/commands/")
		if name == "" {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "missing command"})
			return
		}

		var commandRequest CommandRequest
		if err := json.NewDecoder(request.Body).Decode(&commandRequest); err != nil {
			writeJSON(writer, http.StatusBadRequest, map[string]any{"error": "invalid command request"})
			return
		}
		commandRequest.Name = name

		bridge.mu.RLock()
		handler, ok := bridge.commands[name]
		bridge.mu.RUnlock()
		if !ok {
			writeJSON(writer, http.StatusNotFound, map[string]any{"error": "unknown command"})
			return
		}

		writeJSON(writer, http.StatusOK, handler(commandRequest))
	})

	return mux
}

func (bridge *Bridge) frameSnapshot() FrameSnapshot {
	bridge.mu.RLock()
	provider := bridge.frameProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return FrameSnapshot{}
	}
	return provider()
}

func (bridge *Bridge) sceneSnapshot() SceneSnapshot {
	bridge.mu.RLock()
	provider := bridge.sceneProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return SceneSnapshot{
			Summary: map[string]any{},
		}
	}
	snapshot := provider()
	if snapshot.Summary == nil {
		snapshot.Summary = map[string]any{}
	}
	return snapshot
}

func (bridge *Bridge) worldSnapshot() WorldSnapshot {
	bridge.mu.RLock()
	provider := bridge.worldProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return WorldSnapshot{}
	}
	return provider()
}

func (bridge *Bridge) uiSnapshot() UISnapshot {
	bridge.mu.RLock()
	provider := bridge.uiProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UISnapshot{}
	}
	return provider()
}

func (bridge *Bridge) commandList() []commandDescriptor {
	bridge.mu.RLock()
	names := make([]string, 0, len(bridge.commands))
	for name := range bridge.commands {
		names = append(names, name)
	}
	bridge.mu.RUnlock()

	sort.Strings(names)

	commands := make([]commandDescriptor, 0, len(names))
	for _, name := range names {
		commands = append(commands, commandDescriptor{Name: name})
	}
	return commands
}

func (bridge *Bridge) uiOverview() UIOverviewSnapshot {
	bridge.mu.RLock()
	provider := bridge.uiOverviewProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UIOverviewSnapshot{}
	}
	return provider()
}

func (bridge *Bridge) uiQuery(request UIQueryRequest) UIQueryResult {
	bridge.mu.RLock()
	provider := bridge.uiQueryProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UIQueryResult{}
	}
	if request.Limit <= 0 {
		request.Limit = 25
	}
	return provider(request)
}

func (bridge *Bridge) uiNode(request UINodeInspectRequest) (UINodeDetailSnapshot, bool) {
	bridge.mu.RLock()
	provider := bridge.uiNodeProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UINodeDetailSnapshot{}, false
	}
	if request.ChildDepth <= 0 {
		request.ChildDepth = 1
	}
	return provider(request)
}

func (bridge *Bridge) uiIssues(request UIIssueListRequest) UIIssueListSnapshot {
	bridge.mu.RLock()
	provider := bridge.uiIssuesProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UIIssueListSnapshot{}
	}
	if request.Limit <= 0 {
		request.Limit = 50
	}
	return provider(request)
}

func (bridge *Bridge) uiCapture(request UICaptureRequest) (UICaptureResult, bool) {
	bridge.mu.RLock()
	provider := bridge.uiCaptureProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UICaptureResult{}, false
	}
	if request.Scale <= 0 {
		request.Scale = 1
	}
	return provider(request)
}

func (bridge *Bridge) uiArtifact(id string) (UIArtifact, bool) {
	bridge.mu.RLock()
	provider := bridge.uiArtifactProvider
	bridge.mu.RUnlock()
	if provider == nil {
		return UIArtifact{}, false
	}
	return provider(id)
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func validateLoopbackAddress(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host == "" || host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("debug bridge must bind to loopback, got %q", addr)
	}
	return nil
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func boolQuery(request *http.Request, name string) bool {
	value := strings.TrimSpace(request.URL.Query().Get(name))
	if value == "" {
		return false
	}
	switch strings.ToLower(value) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func defaultBoolQuery(request *http.Request, name string, fallback bool) bool {
	value := strings.TrimSpace(request.URL.Query().Get(name))
	if value == "" {
		return fallback
	}
	return boolQuery(request, name)
}

func defaultIntQuery(request *http.Request, name string, fallback int) int {
	value := strings.TrimSpace(request.URL.Query().Get(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
