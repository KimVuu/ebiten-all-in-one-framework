package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebitendebug"
	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
	"github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom/ebitenrenderer"
)

type game struct {
	mu sync.RWMutex

	width        int
	height       int
	frame        int
	tick         int
	debugEnabled bool

	renderer       *ebitenrenderer.Renderer
	runtime        *uidom.Runtime
	dom            *uidom.DOM
	lastInput      uidom.InputSnapshot
	pageScroll     float64
	overlayEnabled bool
	debugQueue     *debugInputQueue
	debugBridge    *ebitendebug.Bridge
	artifacts      map[string]ebitendebug.UIArtifact
}

func newGame(debugMode bool) *game {
	return &game{
		width:          1280,
		height:         720,
		debugEnabled:   debugMode,
		overlayEnabled: debugMode,
		renderer:       ebitenrenderer.New(),
		runtime:        uidom.NewRuntime(),
		debugQueue:     newDebugInputQueue(),
		artifacts:      map[string]ebitendebug.UIArtifact{},
	}
}

func (g *game) Update() error {
	return g.step(g.collectInput())
}

func (g *game) step(input uidom.InputSnapshot) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.tick++
	g.frame++
	frame := g.frame
	width := maxInt(g.width, 1280)
	height := maxInt(g.height, 720)
	pageScroll := g.pageScroll

	viewport := uidom.Viewport{
		Width:  float64(width),
		Height: float64(height),
	}

	dom := buildShowcaseDOMWithState(showcaseLayoutState{PageScroll: pageScroll}, nil, g.runtime)
	layout := dom.Layout(viewport)
	input = g.applyQueuedDebugEffects(frame, dom, layout, input)
	input = g.applyHostKeyboardInput(dom, layout, input)

	nextPageScroll := pageScroll
	dom = buildShowcaseDOMWithState(showcaseLayoutState{PageScroll: pageScroll}, func(offset float64) {
		nextPageScroll = offset
	}, g.runtime)
	g.runtime.Update(dom, viewport, input)

	if nextPageScroll != pageScroll {
		pageScroll = nextPageScroll
		dom = buildShowcaseDOMWithState(showcaseLayoutState{PageScroll: pageScroll}, func(offset float64) {
			nextPageScroll = offset
		}, g.runtime)
		g.runtime.Update(dom, viewport, input)
	}

	g.pageScroll = pageScroll
	g.dom = dom
	g.lastInput = input
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	size := screen.Bounds().Size()

	g.mu.Lock()
	g.width = size.X
	g.height = size.Y
	dom := g.dom
	pageScroll := g.pageScroll
	overlayEnabled := g.overlayEnabled
	g.mu.Unlock()

	if dom == nil {
		dom = buildShowcaseDOMWithState(showcaseLayoutState{PageScroll: pageScroll}, nil, g.runtime)
	}

	viewport := uidom.Viewport{
		Width:  float64(size.X),
		Height: float64(size.Y),
	}
	layout := g.renderer.Draw(screen, dom, viewport)
	report := buildDebugLayoutReport(layout, viewport)
	drawDebugOverlay(screen, layout, report, overlayEnabled)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func (g *game) startDebugBridge(addr string) error {
	bridge := ebitendebug.New(ebitendebug.Config{
		Enabled: true,
		Addr:    addr,
		GameID:  "ui-dom-showcase",
		Version: "v1",
	})
	bridge.SetFrameProvider(g.frameSnapshot)
	bridge.SetSceneProvider(g.sceneSnapshot)
	bridge.SetWorldProvider(g.worldSnapshot)
	bridge.SetUIProvider(g.uiSnapshot)
	bridge.SetUIOverviewProvider(g.uiOverview)
	bridge.SetUIQueryProvider(g.uiQuery)
	bridge.SetUINodeProvider(g.uiNodeDetail)
	bridge.SetUIIssuesProvider(g.uiIssues)
	bridge.SetUICaptureProvider(g.uiCapture)
	bridge.SetUIArtifactProvider(g.uiArtifact)
	g.registerDebugCommands(bridge)
	if err := bridge.Start(); err != nil {
		return err
	}

	g.mu.Lock()
	g.debugBridge = bridge
	g.mu.Unlock()
	return nil
}

func (g *game) stopDebugBridge() error {
	g.mu.Lock()
	bridge := g.debugBridge
	g.debugBridge = nil
	g.mu.Unlock()

	if bridge == nil {
		return nil
	}
	return bridge.Close(context.Background())
}

func (g *game) frameSnapshot() ebitendebug.FrameSnapshot {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return ebitendebug.FrameSnapshot{
		Frame:        g.frame,
		Tick:         g.tick,
		FPS:          ebiten.ActualFPS(),
		TPS:          ebiten.ActualTPS(),
		Paused:       false,
		DebugEnabled: g.debugEnabled,
	}
}

func (g *game) sceneSnapshot() ebitendebug.SceneSnapshot {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return ebitendebug.SceneSnapshot{
		Current: ebitendebug.SceneRef{
			ID:   "ui-dom-showcase",
			Name: "UI DOM Showcase",
		},
		Known: []ebitendebug.SceneRef{
			{ID: "ui-dom-showcase", Name: "UI DOM Showcase"},
		},
		Summary: map[string]any{
			"viewportWidth":  g.width,
			"viewportHeight": g.height,
			"frame":          g.frame,
			"tick":           g.tick,
		},
	}
}

func (g *game) worldSnapshot() ebitendebug.WorldSnapshot {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.WorldSnapshot{}
	}

	sectionIDs := []string{
		"overview-section",
		"button-section",
		"foundation-section",
		"components-section",
		"prefabs-section",
	}

	entities := make([]ebitendebug.EntitySnapshot, 0, len(sectionIDs))
	for _, id := range sectionIDs {
		node, ok := layout.FindByID(id)
		if !ok {
			continue
		}
		entities = append(entities, ebitendebug.EntitySnapshot{
			ID:      id,
			Type:    "section",
			Visible: true,
			Enabled: true,
			Tags:    []string{"showcase", "ui"},
			Position: ebitendebug.Vector2{
				X: node.Frame.X,
				Y: node.Frame.Y,
			},
			Size: ebitendebug.Vector2{
				X: node.Frame.Width,
				Y: node.Frame.Height,
			},
			Props: map[string]any{
				"title": strings.TrimSuffix(id, "-section"),
			},
		})
	}

	return ebitendebug.WorldSnapshot{Entities: entities}
}

func (g *game) uiSnapshot() ebitendebug.UISnapshot {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UISnapshot{}
	}

	g.mu.RLock()
	width := g.width
	height := g.height
	pageScroll := g.pageScroll
	overlayEnabled := g.overlayEnabled
	lastInput := g.lastInput
	g.mu.RUnlock()

	viewport := uidom.Viewport{
		Width:  float64(maxInt(width, 1280)),
		Height: float64(maxInt(height, 720)),
	}
	report := buildDebugLayoutReport(layout, viewport)
	snapshot := buildDebugUISnapshot(layout, viewport, report, overlayEnabled, g.runtime, lastInput, g.debugQueue.len())
	if snapshot.Root.Props == nil {
		snapshot.Root.Props = map[string]any{}
	}
	snapshot.Root.Props["pageScroll"] = pageScroll
	return snapshot
}

func (g *game) uiOverview() ebitendebug.UIOverviewSnapshot {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UIOverviewSnapshot{}
	}
	return buildCompactUIOverview(layout, g.currentViewport(), buildDebugLayoutReport(layout, g.currentViewport()), g.runtime, g.lastInput, g.debugQueue.len())
}

func (g *game) uiQuery(request ebitendebug.UIQueryRequest) ebitendebug.UIQueryResult {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UIQueryResult{}
	}
	viewport := g.currentViewport()
	return queryCompactUINodes(layout, viewport, buildDebugLayoutReport(layout, viewport), request)
}

func (g *game) uiNodeDetail(request ebitendebug.UINodeInspectRequest) (ebitendebug.UINodeDetailSnapshot, bool) {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UINodeDetailSnapshot{}, false
	}
	viewport := g.currentViewport()
	return inspectCompactUINode(layout, viewport, buildDebugLayoutReport(layout, viewport), request)
}

func (g *game) uiIssues(request ebitendebug.UIIssueListRequest) ebitendebug.UIIssueListSnapshot {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UIIssueListSnapshot{}
	}
	viewport := g.currentViewport()
	return listCompactUIIssues(buildDebugLayoutReport(layout, viewport), request)
}

func (g *game) uiCapture(request ebitendebug.UICaptureRequest) (ebitendebug.UICaptureResult, bool) {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.UICaptureResult{}, false
	}
	result, artifact, ok := captureCompactUIScreenshot("ui-dom-showcase", layout, g.currentViewport(), buildDebugLayoutReport(layout, g.currentViewport()), request)
	if !ok {
		return ebitendebug.UICaptureResult{}, false
	}
	g.mu.Lock()
	g.artifacts[artifact.ID] = artifact
	g.mu.Unlock()
	return result, true
}

func (g *game) uiArtifact(id string) (ebitendebug.UIArtifact, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	artifact, ok := g.artifacts[id]
	return artifact, ok
}

func (g *game) currentViewport() uidom.Viewport {
	g.mu.RLock()
	width := g.width
	height := g.height
	g.mu.RUnlock()
	return uidom.Viewport{
		Width:  float64(maxInt(width, 1280)),
		Height: float64(maxInt(height, 720)),
	}
}

func (g *game) currentLayout() *uidom.LayoutNode {
	g.mu.RLock()
	dom := g.dom
	width := g.width
	height := g.height
	pageScroll := g.pageScroll
	runtimeLayout := g.runtime.Layout()
	g.mu.RUnlock()

	if width <= 0 {
		width = 1280
	}
	if height <= 0 {
		height = 720
	}

	if dom == nil {
		dom = buildShowcaseDOMWithState(showcaseLayoutState{PageScroll: pageScroll}, nil, g.runtime)
	}
	if runtimeLayout != nil {
		return runtimeLayout
	}

	return dom.Layout(uidom.Viewport{
		Width:  float64(width),
		Height: float64(height),
	})
}

func (g *game) collectInput() uidom.InputSnapshot {
	pointerX, pointerY := ebiten.CursorPosition()
	scrollX, scrollY := ebiten.Wheel()
	textInput := ebiten.AppendInputChars(nil)

	input := uidom.InputSnapshot{
		PointerX:    float64(pointerX),
		PointerY:    float64(pointerY),
		PointerDown: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		ScrollX:     scrollX,
		ScrollY:     scrollY,
	}
	if len(textInput) > 0 {
		input.Text = string(textInput)
	}
	input.Backspace = inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
	input.Delete = inpututil.IsKeyJustPressed(ebiten.KeyDelete)
	input.Home = inpututil.IsKeyJustPressed(ebiten.KeyHome)
	input.End = inpututil.IsKeyJustPressed(ebiten.KeyEnd)
	input.Submit = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	input.Space = inpututil.IsKeyJustPressed(ebiten.KeySpace)
	input.Tab = inpututil.IsKeyJustPressed(ebiten.KeyTab)
	input.Escape = inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	input.ArrowUp = inpututil.IsKeyJustPressed(ebiten.KeyArrowUp)
	input.ArrowDown = inpututil.IsKeyJustPressed(ebiten.KeyArrowDown)
	input.ArrowLeft = inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft)
	input.ArrowRight = inpututil.IsKeyJustPressed(ebiten.KeyArrowRight)
	input.Shift = ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	input.Control = ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)
	input.Alt = ebiten.IsKeyPressed(ebiten.KeyAltLeft) || ebiten.IsKeyPressed(ebiten.KeyAltRight)
	input.Meta = ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight)
	input.SelectAll = (input.Control || input.Meta) && inpututil.IsKeyJustPressed(ebiten.KeyA)
	if input.Control || input.Meta {
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			input.ArrowLeft = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			input.Delete = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			input.End = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			input.ArrowRight = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyH) {
			input.Backspace = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			input.ArrowDown = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			input.ArrowUp = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			input.Backspace = true
		}
	}
	return input
}

func (g *game) applyHostKeyboardInput(dom *uidom.DOM, layout *uidom.LayoutNode, input uidom.InputSnapshot) uidom.InputSnapshot {
	return input
}

func maxInt(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func resolveShowcaseDebugAddr() string {
	if value := strings.TrimSpace(ebitenDebugAddrFromEnv()); value != "" {
		return value
	}
	return showcaseDebugDefaultAddr
}

func ebitenDebugAddrFromEnv() string {
	return strings.TrimSpace(getenv("EBITEN_DEBUG_ADDR"))
}

func debugModeEnabled() bool {
	value := strings.TrimSpace(strings.ToLower(getenv("EBITEN_DEBUG_MODE")))
	switch value {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func debugBridgeBanner(addr string) string {
	if addr == "" {
		return "debug bridge disabled"
	}
	return fmt.Sprintf("debug bridge: http://%s", addr)
}

const showcaseDebugDefaultAddr = "127.0.0.1:47832"
