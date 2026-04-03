package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	ebitenuidebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui-debug"
	renderer "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui/renderer"
)

type game struct {
	mu sync.RWMutex

	width        int
	height       int
	frame        int
	tick         int
	debugEnabled bool

	renderer       *renderer.Renderer
	runtime        *ebitenui.Runtime
	dom            *ebitenui.DOM
	lastInput      ebitenui.InputSnapshot
	sidebarScroll  float64
	detailScroll   float64
	themePreset    string
	fontPreset     string
	overlayEnabled bool

	registry    ShowcasePageRegistry
	router      *ebitenui.PageRouter
	bindings    *showcaseBindings
	uiDebug     *ebitenuidebug.Adapter
	debugBridge *ebitendebug.Bridge
}

type showcaseBindingSnapshot struct {
	NameInput      string
	Resolution     string
	ResolutionOpen bool
	Bio            string
	Hardcore       bool
	MusicVolume    float64
	ButtonPhase    string
	ButtonDowns    int
	ButtonHolds    int
	ButtonUps      int
	ButtonClicks   int
}

func newGame(debugMode bool) *game {
	registry := buildShowcasePageRegistry()
	router := ebitenui.NewPageRouter(ebitenui.PageRouterConfig{
		Routes:        registry.Routes,
		InitialPageID: "overview",
	})

	g := &game{
		width:          1280,
		height:         720,
		debugEnabled:   debugMode,
		overlayEnabled: false,
		themePreset:    "default",
		fontPreset:     "default",
		registry:       registry,
		router:         router,
		bindings:       newShowcaseBindings(),
	}
	g.renderer = renderer.New()
	g.runtime = ebitenui.NewRuntime()
	g.uiDebug = ebitenuidebug.NewAdapter(ebitenuidebug.Config{
		GameID:         "ebiten-ui-showcase",
		ScreenshotsDir: showcaseScreenshotsDir(),
	}, ebitenuidebug.Callbacks{
		CurrentLayout:   g.currentLayout,
		CurrentViewport: g.currentViewport,
		CurrentRuntime: func() *ebitenui.Runtime {
			return g.runtime
		},
		CurrentInput: func() ebitenui.InputSnapshot {
			g.mu.RLock()
			defer g.mu.RUnlock()
			return g.lastInput
		},
		CurrentFrame:   g.currentFrame,
		OverlayEnabled: func() bool { g.mu.RLock(); defer g.mu.RUnlock(); return g.overlayEnabled },
		SetOverlay: func(enabled bool) {
			g.mu.Lock()
			g.overlayEnabled = enabled
			g.mu.Unlock()
		},
	})
	return g
}

func showcaseScreenshotsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join(".", "screenshots")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "..", "screenshots"))
}

func (g *game) Update() error {
	return g.step(g.collectInput())
}

func (g *game) step(input ebitenui.InputSnapshot) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.tick++
	g.frame++
	frame := g.frame
	viewport := ebitenui.Viewport{
		Width:  float64(maxInt(g.width, 1280)),
		Height: float64(maxInt(g.height, 720)),
	}

	state := showcaseLayoutState{
		CurrentPageID: g.router.CurrentPageID(),
		ThemePreset:   g.themePreset,
		FontPreset:    g.fontPreset,
		SidebarScroll: g.sidebarScroll,
		DetailScroll:  g.detailScroll,
	}

	dom := buildShowcaseDOMWithState(state, nil, g.runtime, g.bindings)
	layout := dom.Layout(viewport)
	bindingBefore := g.bindingSnapshot()
	input = g.uiDebug.ApplyQueuedInput(frame, dom, g.runtime, layout, input)
	input = g.applyHostKeyboardInput(dom, layout, input)

	nextState := state
	callbacks := &showcaseCallbacks{
		OnNavigate: func(pageID string) {
			nextState.CurrentPageID = pageID
		},
		OnThemePresetChange: func(themePreset string) {
			nextState.ThemePreset = themePreset
		},
		OnFontPresetChange: func(fontPreset string) {
			nextState.FontPreset = fontPreset
		},
		OnSidebarScrollChange: func(offset float64) {
			nextState.SidebarScroll = offset
		},
		OnDetailScrollChange: func(offset float64) {
			nextState.DetailScroll = offset
		},
	}

	dom = buildShowcaseDOMWithState(state, callbacks, g.runtime, g.bindings)
	g.runtime.Update(dom, viewport, input)

	if nextState.CurrentPageID == "" {
		nextState.CurrentPageID = state.CurrentPageID
	}
	if !g.router.Navigate(nextState.CurrentPageID) {
		nextState.CurrentPageID = state.CurrentPageID
	}
	nextState.CurrentPageID = g.router.CurrentPageID()

	pageChanged := nextState.CurrentPageID != state.CurrentPageID
	if pageChanged {
		nextState.DetailScroll = 0
	}

	bindingsChanged := g.bindingSnapshot() != bindingBefore
	scrollChanged := nextState.SidebarScroll != state.SidebarScroll || nextState.DetailScroll != state.DetailScroll
	if pageChanged || scrollChanged || bindingsChanged {
		dom = buildShowcaseDOMWithState(nextState, callbacks, g.runtime, g.bindings)
		g.runtime.Update(dom, viewport, stabilizeShowcaseInput(input))
	}

	g.sidebarScroll = nextState.SidebarScroll
	g.detailScroll = nextState.DetailScroll
	g.themePreset = initialShowcaseThemePreset(nextState.ThemePreset)
	g.fontPreset = initialShowcaseFontPreset(nextState.FontPreset)
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
	overlayEnabled := g.overlayEnabled
	g.mu.Unlock()

	if dom == nil {
		dom = buildShowcaseDOMWithState(g.currentState(), nil, g.runtime, g.bindings)
	}

	viewport := ebitenui.Viewport{
		Width:  float64(size.X),
		Height: float64(size.Y),
	}
	layout := g.renderer.Draw(screen, dom, viewport)
	g.uiDebug.DrawOverlay(screen, layout, overlayEnabled)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func (g *game) startDebugBridge(addr string) error {
	bridge := ebitendebug.New(ebitendebug.Config{
		Enabled: true,
		Addr:    addr,
		GameID:  "ebiten-ui-showcase",
		Version: "v1",
	})
	bridge.SetFrameProvider(g.frameSnapshot)
	bridge.SetSceneProvider(g.sceneSnapshot)
	bridge.SetWorldProvider(g.worldSnapshot)
	g.uiDebug.Attach(bridge)
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

func (g *game) debugBridgeLikeCommand(name string, args map[string]any) ebitendebug.CommandResult {
	g.mu.RLock()
	bridge := g.debugBridge
	g.mu.RUnlock()
	if bridge != nil {
		return bridge.InvokeCommand(name, args)
	}

	bridge = ebitendebug.New(ebitendebug.Config{Enabled: true, GameID: "ebiten-ui-showcase", Version: "v1"})
	g.uiDebug.Attach(bridge)
	return bridge.InvokeCommand(name, args)
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
			ID:   "ebiten-ui-showcase",
			Name: "Ebiten UI Showcase",
		},
		Known: []ebitendebug.SceneRef{
			{ID: "ebiten-ui-showcase", Name: "Ebiten UI Showcase"},
		},
		Summary: map[string]any{
			"viewportWidth":  g.width,
			"viewportHeight": g.height,
			"frame":          g.frame,
			"tick":           g.tick,
			"currentPageID":  g.router.CurrentPageID(),
			"themePreset":    g.themePreset,
			"fontPreset":     g.fontPreset,
		},
	}
}

func (g *game) worldSnapshot() ebitendebug.WorldSnapshot {
	layout := g.currentLayout()
	if layout == nil {
		return ebitendebug.WorldSnapshot{}
	}

	currentPageID := g.currentPageID()
	entityIDs := []struct {
		ID   string
		Type string
	}{
		{ID: "showcase-sidebar", Type: "navigation"},
		{ID: "showcase-detail", Type: "content"},
		{ID: "nav-item-" + sanitizeID(currentPageID), Type: "nav-item"},
		{ID: "page-demo", Type: "page-demo"},
	}

	entities := make([]ebitendebug.EntitySnapshot, 0, len(entityIDs))
	for _, entry := range entityIDs {
		node, ok := layout.FindByID(entry.ID)
		if !ok {
			continue
		}
		entities = append(entities, ebitendebug.EntitySnapshot{
			ID:      entry.ID,
			Type:    entry.Type,
			Visible: true,
			Enabled: true,
			Tags:    []string{"showcase", "ui", currentPageID},
			Position: ebitendebug.Vector2{
				X: node.Frame.X,
				Y: node.Frame.Y,
			},
			Size: ebitendebug.Vector2{
				X: node.Frame.Width,
				Y: node.Frame.Height,
			},
			Props: map[string]any{
				"currentPageID": currentPageID,
			},
		})
	}

	return ebitendebug.WorldSnapshot{Entities: entities}
}

func (g *game) uiSnapshot() ebitendebug.UISnapshot {
	snapshot := g.uiDebug.UISnapshot()
	if snapshot.Root.Props == nil {
		snapshot.Root.Props = map[string]any{}
	}
	g.mu.RLock()
	snapshot.Root.Props["currentPageID"] = g.router.CurrentPageID()
	snapshot.Root.Props["sidebarScroll"] = g.sidebarScroll
	snapshot.Root.Props["detailScroll"] = g.detailScroll
	snapshot.Root.Props["themePreset"] = g.themePreset
	snapshot.Root.Props["fontPreset"] = g.fontPreset
	g.mu.RUnlock()
	return snapshot
}

func (g *game) currentViewport() ebitenui.Viewport {
	g.mu.RLock()
	width := g.width
	height := g.height
	g.mu.RUnlock()
	return ebitenui.Viewport{
		Width:  float64(maxInt(width, 1280)),
		Height: float64(maxInt(height, 720)),
	}
}

func (g *game) currentLayout() *ebitenui.LayoutNode {
	g.mu.RLock()
	dom := g.dom
	width := g.width
	height := g.height
	runtimeLayout := g.runtime.Layout()
	state := g.currentStateLocked()
	g.mu.RUnlock()

	if width <= 0 {
		width = 1280
	}
	if height <= 0 {
		height = 720
	}

	if dom == nil {
		dom = buildShowcaseDOMWithState(state, nil, g.runtime, g.bindings)
	}
	if runtimeLayout != nil {
		return runtimeLayout
	}

	return dom.Layout(ebitenui.Viewport{
		Width:  float64(width),
		Height: float64(height),
	})
}

func (g *game) currentPageID() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.router.CurrentPageID()
}

func (g *game) currentState() showcaseLayoutState {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.currentStateLocked()
}

func (g *game) currentStateLocked() showcaseLayoutState {
	return showcaseLayoutState{
		CurrentPageID: g.router.CurrentPageID(),
		ThemePreset:   g.themePreset,
		FontPreset:    g.fontPreset,
		SidebarScroll: g.sidebarScroll,
		DetailScroll:  g.detailScroll,
	}
}

func (g *game) bindingSnapshot() showcaseBindingSnapshot {
	if g == nil || g.bindings == nil {
		return showcaseBindingSnapshot{}
	}
	return showcaseBindingSnapshot{
		NameInput:      g.bindings.NameInput.Get(),
		Resolution:     g.bindings.Resolution.Get(),
		ResolutionOpen: g.bindings.ResolutionOpen.Get(),
		Bio:            g.bindings.Bio.Get(),
		Hardcore:       g.bindings.Hardcore.Get(),
		MusicVolume:    g.bindings.MusicVolume.Get(),
		ButtonPhase:    g.bindings.ButtonPhase.Get(),
		ButtonDowns:    g.bindings.ButtonDowns.Get(),
		ButtonHolds:    g.bindings.ButtonHolds.Get(),
		ButtonUps:      g.bindings.ButtonUps.Get(),
		ButtonClicks:   g.bindings.ButtonClicks.Get(),
	}
}

func (g *game) collectInput() ebitenui.InputSnapshot {
	pointerX, pointerY := ebiten.CursorPosition()
	scrollX, scrollY := ebiten.Wheel()
	textInput := ebiten.AppendInputChars(nil)

	input := ebitenui.InputSnapshot{
		PointerX:     float64(pointerX),
		PointerY:     float64(pointerY),
		PointerDown:  ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		InputBlocked: !ebiten.IsFocused(),
		ScrollX:      scrollX,
		ScrollY:      scrollY,
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

func (g *game) applyHostKeyboardInput(dom *ebitenui.DOM, layout *ebitenui.LayoutNode, input ebitenui.InputSnapshot) ebitenui.InputSnapshot {
	return input
}

func (g *game) currentFrame() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.frame
}

func stabilizeShowcaseInput(input ebitenui.InputSnapshot) ebitenui.InputSnapshot {
	input.ScrollX = 0
	input.ScrollY = 0
	input.Text = ""
	input.Backspace = false
	input.Delete = false
	input.Home = false
	input.End = false
	input.Submit = false
	input.Space = false
	input.SelectAll = false
	input.Shortcut = ""
	input.Tab = false
	input.Escape = false
	input.ArrowUp = false
	input.ArrowDown = false
	input.ArrowLeft = false
	input.ArrowRight = false
	input.Shift = false
	input.Control = false
	input.Alt = false
	input.Meta = false
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
