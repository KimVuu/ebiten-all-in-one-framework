package ebitenuidebug

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

type Config struct {
	GameID         string
	ScreenshotsDir string
}

type Callbacks struct {
	CurrentLayout   func() *ebitenui.LayoutNode
	CurrentViewport func() ebitenui.Viewport
	CurrentRuntime  func() *ebitenui.Runtime
	CurrentInput    func() ebitenui.InputSnapshot
	CurrentFrame    func() int
	OverlayEnabled  func() bool
	SetOverlay      func(bool)
}

type Adapter struct {
	config     Config
	callbacks  Callbacks
	debugQueue *debugInputQueue

	mu        sync.RWMutex
	artifacts map[string]ebitendebug.UIArtifact
}

func NewAdapter(config Config, callbacks Callbacks) *Adapter {
	return &Adapter{
		config:     config,
		callbacks:  callbacks,
		debugQueue: newDebugInputQueue(),
		artifacts:  map[string]ebitendebug.UIArtifact{},
	}
}

func (adapter *Adapter) Attach(bridge *ebitendebug.Bridge) {
	if adapter == nil || bridge == nil {
		return
	}
	bridge.SetUIProvider(adapter.UISnapshot)
	bridge.SetUIOverviewProvider(adapter.UIOverview)
	bridge.SetUIQueryProvider(adapter.UIQuery)
	bridge.SetUINodeProvider(adapter.UINodeDetail)
	bridge.SetUIIssuesProvider(adapter.UIIssues)
	bridge.SetUICaptureProvider(adapter.UICapture)
	bridge.SetUIArtifactProvider(adapter.Artifact)
	adapter.registerCommands(bridge)
}

func (adapter *Adapter) QueueLength() int {
	if adapter == nil || adapter.debugQueue == nil {
		return 0
	}
	return adapter.debugQueue.len()
}

func (adapter *Adapter) Artifact(id string) (ebitendebug.UIArtifact, bool) {
	adapter.mu.RLock()
	defer adapter.mu.RUnlock()
	artifact, ok := adapter.artifacts[id]
	return artifact, ok
}

func (adapter *Adapter) DrawOverlay(screen *ebiten.Image, layout *ebitenui.LayoutNode, overlayEnabled bool) {
	if adapter == nil || screen == nil || layout == nil {
		return
	}
	drawDebugOverlay(screen, layout, buildDebugLayoutReport(layout, adapter.currentViewport()), overlayEnabled)
}

func (adapter *Adapter) ApplyQueuedInput(frame int, dom *ebitenui.DOM, runtime *ebitenui.Runtime, layout *ebitenui.LayoutNode, input ebitenui.InputSnapshot) ebitenui.InputSnapshot {
	if adapter == nil || adapter.debugQueue == nil {
		return input
	}
	effects := adapter.debugQueue.drain(frame)
	if len(effects) == 0 {
		return input
	}

	for _, effect := range effects {
		switch effect.kind {
		case debugFrameInputPointerMove, debugFrameInputPointerDown, debugFrameInputPointerUp, debugFrameInputScroll, debugFrameInputText:
			input = mergeDebugInputSnapshot(input, effect.input, effect.kind)
		case debugFrameInputFocusNode:
			if runtime != nil && dom != nil && effect.focusTargetID != "" {
				runtime.FocusNodeByID(dom, effect.focusTargetID, input)
			}
		case debugFrameInputClearFocus:
			if runtime != nil && dom != nil {
				runtime.ClearFocus(dom, input)
			}
		case debugFrameInputOverlayToggle:
			if effect.overlayEnabled != nil && adapter.callbacks.SetOverlay != nil {
				adapter.callbacks.SetOverlay(*effect.overlayEnabled)
			}
		}
	}
	return input
}

func (adapter *Adapter) UISnapshot() ebitendebug.UISnapshot {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UISnapshot{}
	}
	viewport := adapter.currentViewport()
	report := buildDebugLayoutReport(layout, viewport)
	return buildDebugUISnapshot(layout, viewport, report, adapter.overlayEnabled(), adapter.currentRuntime(), adapter.currentInput(), adapter.QueueLength())
}

func (adapter *Adapter) UIOverview() ebitendebug.UIOverviewSnapshot {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UIOverviewSnapshot{}
	}
	viewport := adapter.currentViewport()
	return buildCompactUIOverview(layout, viewport, buildDebugLayoutReport(layout, viewport), adapter.currentRuntime(), adapter.currentInput(), adapter.QueueLength())
}

func (adapter *Adapter) UIQuery(request ebitendebug.UIQueryRequest) ebitendebug.UIQueryResult {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UIQueryResult{}
	}
	viewport := adapter.currentViewport()
	return queryCompactUINodes(layout, viewport, buildDebugLayoutReport(layout, viewport), request)
}

func (adapter *Adapter) UINodeDetail(request ebitendebug.UINodeInspectRequest) (ebitendebug.UINodeDetailSnapshot, bool) {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UINodeDetailSnapshot{}, false
	}
	viewport := adapter.currentViewport()
	return inspectCompactUINode(layout, viewport, buildDebugLayoutReport(layout, viewport), request)
}

func (adapter *Adapter) UIIssues(request ebitendebug.UIIssueListRequest) ebitendebug.UIIssueListSnapshot {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UIIssueListSnapshot{}
	}
	viewport := adapter.currentViewport()
	return listCompactUIIssues(buildDebugLayoutReport(layout, viewport), request)
}

func (adapter *Adapter) UICapture(request ebitendebug.UICaptureRequest) (ebitendebug.UICaptureResult, bool) {
	layout := adapter.currentLayout()
	if layout == nil {
		return ebitendebug.UICaptureResult{}, false
	}
	viewport := adapter.currentViewport()
	result, artifact, ok := captureCompactUIScreenshot(adapter.config.GameID, adapter.config.ScreenshotsDir, layout, viewport, buildDebugLayoutReport(layout, viewport), request)
	if !ok {
		return ebitendebug.UICaptureResult{}, false
	}
	adapter.mu.Lock()
	adapter.artifacts[artifact.ID] = artifact
	adapter.mu.Unlock()
	return result, true
}

func (adapter *Adapter) registerCommands(bridge *ebitendebug.Bridge) {
	bridge.RegisterCommand("validate_ui_layout", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		layout := adapter.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		viewport := adapter.currentViewport()
		report := buildDebugLayoutReport(layout, viewport)
		return ebitendebug.CommandResult{
			Success: true,
			Payload: map[string]any{
				"viewport":         report.Viewport,
				"issues":           report.Issues,
				"issueSummary":     report.IssueSummary,
				"invalidNodeCount": report.InvalidNodeCount,
			},
		}
	})

	bridge.RegisterCommand("inspect_ui_node", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		nodeID := commandNodeID(request.Args)
		if nodeID == "" {
			return ebitendebug.CommandResult{Success: false, Message: "inspect_ui_node requires node_id"}
		}
		detail, ok := adapter.UINodeDetail(ebitendebug.UINodeInspectRequest{
			NodeID:          nodeID,
			IncludeChildren: true,
			ChildDepth:      1,
			IncludeIssues:   true,
		})
		if !ok {
			return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unknown node: %s", nodeID)}
		}
		return ebitendebug.CommandResult{Success: true, Payload: map[string]any{"node": detail}}
	})

	bridge.RegisterCommand("suggest_ui_constraint_fixes", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		layout := adapter.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		viewport := adapter.currentViewport()
		report := buildDebugLayoutReport(layout, viewport)
		return ebitendebug.CommandResult{
			Success: true,
			Payload: map[string]any{
				"issues":       report.Issues,
				"issueSummary": report.IssueSummary,
			},
		}
	})

	bridge.RegisterCommand("set_ui_debug_overlay", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		enabled, ok := commandBool(request.Args, "enabled")
		if !ok {
			return ebitendebug.CommandResult{Success: false, Message: "set_ui_debug_overlay requires enabled"}
		}
		if adapter.callbacks.SetOverlay != nil {
			adapter.callbacks.SetOverlay(enabled)
		}
		return ebitendebug.CommandResult{Success: true, Payload: map[string]any{"overlayEnabled": enabled}}
	})

	registerQueueCommand := func(name string, build func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error)) {
		bridge.RegisterCommand(name, func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
			layout := adapter.currentLayout()
			if layout == nil {
				return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
			}
			resolved, effects, err := build(adapter.currentFrame(), layout, request.Args)
			if err != nil {
				return ebitendebug.CommandResult{Success: false, Message: err.Error()}
			}
			return debugQueueResult(name, resolved, effects)
		})
	}

	registerQueueCommand("ui_click", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		target, resolved, err := resolvePointerTarget(layout, args)
		if err != nil {
			return "", nil, err
		}
		return resolved, adapter.debugQueue.queueClick(frame, target), nil
	})
	registerQueueCommand("ui_pointer_move", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		target, resolved, err := resolvePointerTarget(layout, args)
		if err != nil {
			return "", nil, err
		}
		effect := adapter.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerMove,
			input:          ebitenui.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5},
			resolvedTarget: resolved,
		})
		return resolved, []debugFrameInput{effect}, nil
	})
	registerQueueCommand("ui_pointer_down", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		target, resolved, err := resolvePointerTarget(layout, args)
		if err != nil {
			return "", nil, err
		}
		effect := adapter.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerDown,
			input:          ebitenui.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5, PointerDown: true},
			resolvedTarget: resolved,
		})
		return resolved, []debugFrameInput{effect}, nil
	})
	registerQueueCommand("ui_pointer_up", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		target, resolved, err := resolvePointerTarget(layout, args)
		if err != nil {
			return "", nil, err
		}
		effect := adapter.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerUp,
			input:          ebitenui.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5},
			resolvedTarget: resolved,
		})
		return resolved, []debugFrameInput{effect}, nil
	})
	registerQueueCommand("ui_scroll", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		deltaX, _ := commandFloat(args, "delta_x")
		deltaY, _ := commandFloat(args, "delta_y")
		target, resolved, err := resolveScrollTargetFromArgs(layout, args)
		if err != nil {
			return "", nil, err
		}
		return resolved, []debugFrameInput{adapter.debugQueue.queueScroll(frame, target, deltaX, deltaY)}, nil
	})
	registerQueueCommand("ui_focus_node", func(frame int, layout *ebitenui.LayoutNode, args map[string]any) (string, []debugFrameInput, error) {
		nodeID := commandNodeID(args)
		if nodeID == "" {
			return "", nil, fmt.Errorf("ui_focus_node requires node_id")
		}
		target, ok := resolveDebugTarget(layout, nodeID)
		if !ok {
			return "", nil, fmt.Errorf("unknown node: %s", nodeID)
		}
		return target.ID, []debugFrameInput{adapter.debugQueue.queueFocus(frame, target.ID)}, nil
	})

	bridge.RegisterCommand("ui_type_text", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		value, ok := commandString(request.Args, "text")
		if !ok || value == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_type_text requires text"}
		}
		layout := adapter.currentLayout()
		runtime := adapter.currentRuntime()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		nodeID := commandNodeID(request.Args)
		targetID := fallbackString(nodeID, focusedID(runtime))
		startFrame := adapter.currentFrame()
		queuedFrames := make([]int, 0, len(value)+1)
		if nodeID != "" {
			target, ok := resolveDebugTarget(layout, nodeID)
			if !ok {
				return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unknown node: %s", nodeID)}
			}
			targetID = target.ID
			adapter.debugQueue.queueFocus(startFrame, targetID)
			queuedFrames = append(queuedFrames, startFrame+1)
			startFrame++
		} else if focusedID(runtime) == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_type_text requires a focused node or node_id"}
		}
		for _, effect := range adapter.debugQueue.queueText(startFrame, targetID, value) {
			queuedFrames = append(queuedFrames, effect.frame)
		}
		return ebitendebug.CommandResult{
			Success:        true,
			Status:         "queued",
			ResolvedTarget: targetID,
			QueuedFrame:    queuedFrames[len(queuedFrames)-1],
			Payload:        map[string]any{"status": "queued", "resolvedTarget": targetID, "queuedFrames": queuedFrames},
		}
	})

	bridge.RegisterCommand("ui_key_event", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		keyName, ok := commandString(request.Args, "key")
		if !ok || keyName == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_key_event requires key"}
		}
		layout := adapter.currentLayout()
		runtime := adapter.currentRuntime()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		startFrame := adapter.currentFrame()
		shift, _ := commandBool(request.Args, "shift")
		control, _ := commandBool(request.Args, "control")
		if !control {
			control, _ = commandBool(request.Args, "ctrl")
		}
		alt, _ := commandBool(request.Args, "alt")
		meta, _ := commandBool(request.Args, "meta")
		nodeID := commandNodeID(request.Args)
		targetID := fallbackString(nodeID, focusedID(runtime))
		if targetID == "" {
			targetID = firstFocusableID(layout)
		}
		queuedFrames := make([]int, 0, 2)
		queueFocus := func() {
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := adapter.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
		}
		switch normalizeKeyName(keyName) {
		case "enter":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{Submit: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "backspace":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{Backspace: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "delete":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{Delete: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "home":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{Home: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "end":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{End: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "space":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{Space: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "selectall":
			queueFocus()
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: ebitenui.InputSnapshot{SelectAll: true, Shift: shift, Control: control, Alt: alt, Meta: meta}, resolvedTarget: targetID}).frame)
		case "tab":
			nextID := nextFocusableID(layout, targetID, shift)
			if nextID == "" {
				return ebitendebug.CommandResult{Success: false, Message: "no focusable node available"}
			}
			queuedFrames = append(queuedFrames, adapter.debugQueue.queueFocus(startFrame, nextID).frame)
			targetID = nextID
		case "escape":
			queuedFrames = append(queuedFrames, adapter.debugQueue.queueClearFocus(startFrame).frame)
			targetID = ""
		case "arrowup", "arrowdown", "arrowleft", "arrowright":
			queueFocus()
			snapshot, ok := shortcutAwareKeyEventSnapshot(normalizeKeyName(keyName), shift, control, alt, meta)
			if !ok {
				snapshot = ebitenui.InputSnapshot{}
			}
			queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: snapshot, resolvedTarget: targetID}).frame)
		default:
			if len([]rune(keyName)) == 1 {
				queueFocus()
				if control || meta {
					snapshot, ok := shortcutAwareKeyEventSnapshot(strings.ToLower(keyName), shift, control, alt, meta)
					if !ok {
						return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unsupported shortcut: %s", keyName)}
					}
					queuedFrames = append(queuedFrames, adapter.debugQueue.enqueue(debugFrameInput{frame: startFrame + 1, kind: debugFrameInputText, input: snapshot, resolvedTarget: targetID}).frame)
				} else {
					for _, effect := range adapter.debugQueue.queueText(startFrame, targetID, keyName) {
						queuedFrames = append(queuedFrames, effect.frame)
					}
				}
			} else {
				return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unsupported key: %s", keyName)}
			}
		}

		return ebitendebug.CommandResult{
			Success:        true,
			Status:         "queued",
			ResolvedTarget: targetID,
			QueuedFrame:    queuedFrames[len(queuedFrames)-1],
			Payload:        map[string]any{"status": "queued", "resolvedTarget": targetID, "queuedFrames": queuedFrames},
		}
	})

	bridge.RegisterCommand("ui_clear_input_queue", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		adapter.debugQueue.clear()
		return ebitendebug.CommandResult{Success: true, Status: "cleared", Payload: map[string]any{"status": "cleared"}}
	})
}

func (adapter *Adapter) currentLayout() *ebitenui.LayoutNode {
	if adapter == nil || adapter.callbacks.CurrentLayout == nil {
		return nil
	}
	return adapter.callbacks.CurrentLayout()
}

func (adapter *Adapter) currentViewport() ebitenui.Viewport {
	if adapter == nil || adapter.callbacks.CurrentViewport == nil {
		return ebitenui.Viewport{}
	}
	return adapter.callbacks.CurrentViewport()
}

func (adapter *Adapter) currentRuntime() *ebitenui.Runtime {
	if adapter == nil || adapter.callbacks.CurrentRuntime == nil {
		return nil
	}
	return adapter.callbacks.CurrentRuntime()
}

func (adapter *Adapter) currentInput() ebitenui.InputSnapshot {
	if adapter == nil || adapter.callbacks.CurrentInput == nil {
		return ebitenui.InputSnapshot{}
	}
	return adapter.callbacks.CurrentInput()
}

func (adapter *Adapter) currentFrame() int {
	if adapter == nil || adapter.callbacks.CurrentFrame == nil {
		return 0
	}
	return adapter.callbacks.CurrentFrame()
}

func (adapter *Adapter) overlayEnabled() bool {
	if adapter == nil || adapter.callbacks.OverlayEnabled == nil {
		return false
	}
	return adapter.callbacks.OverlayEnabled()
}

func focusedID(runtime *ebitenui.Runtime) string {
	if runtime == nil {
		return ""
	}
	return runtime.FocusedID()
}

func mergeDebugInputSnapshot(base, add ebitenui.InputSnapshot, kind debugFrameInputKind) ebitenui.InputSnapshot {
	base.InputBlocked = false
	if kind == debugFrameInputPointerMove || kind == debugFrameInputPointerDown || kind == debugFrameInputPointerUp || kind == debugFrameInputScroll {
		base.PointerX = add.PointerX
		base.PointerY = add.PointerY
	}
	if kind == debugFrameInputPointerDown {
		base.PointerDown = true
	}
	if kind == debugFrameInputPointerUp {
		base.PointerDown = false
	}
	base.ScrollX += add.ScrollX
	base.ScrollY += add.ScrollY
	if add.Text != "" {
		base.Text += add.Text
	}
	base.Backspace = base.Backspace || add.Backspace
	base.Delete = base.Delete || add.Delete
	base.Home = base.Home || add.Home
	base.End = base.End || add.End
	base.Submit = base.Submit || add.Submit
	base.Space = base.Space || add.Space
	base.SelectAll = base.SelectAll || add.SelectAll
	base.Tab = base.Tab || add.Tab
	base.Escape = base.Escape || add.Escape
	base.ArrowUp = base.ArrowUp || add.ArrowUp
	base.ArrowDown = base.ArrowDown || add.ArrowDown
	base.ArrowLeft = base.ArrowLeft || add.ArrowLeft
	base.ArrowRight = base.ArrowRight || add.ArrowRight
	base.Shift = base.Shift || add.Shift
	base.Control = base.Control || add.Control
	base.Alt = base.Alt || add.Alt
	base.Meta = base.Meta || add.Meta
	return base
}

func fallbackString(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func debugQueueResult(action string, resolvedTarget string, effects []debugFrameInput) ebitendebug.CommandResult {
	queuedFrames := make([]int, 0, len(effects))
	for _, effect := range effects {
		queuedFrames = append(queuedFrames, effect.frame)
	}
	return ebitendebug.CommandResult{
		Success:        true,
		Status:         "queued",
		ResolvedTarget: resolvedTarget,
		QueuedFrame:    queuedFrames[len(queuedFrames)-1],
		Payload: map[string]any{
			"action":         action,
			"resolvedTarget": resolvedTarget,
			"queuedFrames":   queuedFrames,
			"status":         "queued",
		},
	}
}

func resolvePointerTarget(layout *ebitenui.LayoutNode, args map[string]any) (debugResolvedTarget, string, error) {
	nodeID := commandNodeID(args)
	if nodeID != "" {
		target, ok := resolveDebugTarget(layout, nodeID)
		if !ok {
			return debugResolvedTarget{}, "", fmt.Errorf("unknown node: %s", nodeID)
		}
		return target, target.ID, nil
	}

	x, hasX := commandFloat(args, "x")
	y, hasY := commandFloat(args, "y")
	if hasX && hasY {
		frame := ebitenui.Rect{X: x, Y: y}
		return debugResolvedTarget{ID: "pointer", Frame: frame}, "pointer", nil
	}

	target, ok := findFirstInteractiveTarget(layout)
	if !ok {
		return debugResolvedTarget{}, "", fmt.Errorf("no interactive target available")
	}
	return target, target.ID, nil
}

func resolveScrollTargetFromArgs(layout *ebitenui.LayoutNode, args map[string]any) (debugResolvedTarget, string, error) {
	x, hasX := commandFloat(args, "x")
	y, hasY := commandFloat(args, "y")
	if hasX && hasY {
		frame := ebitenui.Rect{X: x, Y: y}
		return debugResolvedTarget{ID: "pointer", Frame: frame}, "pointer", nil
	}

	nodeID := commandNodeID(args)
	if nodeID == "" {
		target, ok := findFirstScrollTarget(layout)
		if !ok {
			return debugResolvedTarget{}, "", fmt.Errorf("no scroll target available")
		}
		return target, target.ID, nil
	}

	target, ok := resolveScrollTarget(layout, nodeID)
	if !ok {
		return debugResolvedTarget{}, "", fmt.Errorf("unknown node: %s", nodeID)
	}
	return target, target.ID, nil
}

func findFirstInteractiveTarget(layout *ebitenui.LayoutNode) (debugResolvedTarget, bool) {
	if layout == nil || layout.Node == nil {
		return debugResolvedTarget{}, false
	}
	if isInteractiveLayoutNode(layout) {
		frame := layout.ClickableRect
		if frame == (ebitenui.Rect{}) {
			frame = layout.Frame
		}
		return debugResolvedTarget{ID: layout.Node.Props.ID, Frame: frame, Node: layout, Path: []*ebitenui.LayoutNode{layout}}, true
	}
	for _, child := range layout.Children {
		if target, ok := findFirstInteractiveTarget(child); ok {
			return target, true
		}
	}
	return debugResolvedTarget{}, false
}

func parentLayoutForPath(path []*ebitenui.LayoutNode) *ebitenui.LayoutNode {
	if len(path) < 2 {
		return nil
	}
	return path[len(path)-2]
}

func commandNodeID(args map[string]any) string {
	if value, ok := commandString(args, "node_id"); ok {
		return value
	}
	if value, ok := commandString(args, "target_id"); ok {
		return value
	}
	if value, ok := commandString(args, "id"); ok {
		return value
	}
	return ""
}

func commandString(args map[string]any, key string) (string, bool) {
	if args == nil {
		return "", false
	}
	raw, ok := args[key]
	if !ok {
		return "", false
	}
	value, ok := raw.(string)
	if !ok {
		return "", false
	}
	value = strings.TrimSpace(value)
	return value, value != ""
}

func commandBool(args map[string]any, key string) (bool, bool) {
	if args == nil {
		return false, false
	}
	raw, ok := args[key]
	if !ok {
		return false, false
	}
	switch value := raw.(type) {
	case bool:
		return value, true
	case string:
		value = strings.TrimSpace(strings.ToLower(value))
		switch value {
		case "1", "true", "yes", "on":
			return true, true
		case "0", "false", "no", "off":
			return false, true
		}
	}
	return false, false
}

func commandFloat(args map[string]any, key string) (float64, bool) {
	if args == nil {
		return 0, false
	}
	raw, ok := args[key]
	if !ok {
		return 0, false
	}
	switch value := raw.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	case int32:
		return float64(value), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	}
	return 0, false
}

func normalizeKeyName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func shortcutAwareKeyEventSnapshot(key string, shift, control, alt, meta bool) (ebitenui.InputSnapshot, bool) {
	snapshot := ebitenui.InputSnapshot{Shift: shift, Control: control, Alt: alt, Meta: meta}
	switch key {
	case "arrowup":
		snapshot.ArrowUp = true
	case "arrowdown":
		snapshot.ArrowDown = true
	case "arrowleft":
		snapshot.ArrowLeft = true
	case "arrowright":
		snapshot.ArrowRight = true
	case "enter":
		snapshot.Submit = true
	case "backspace":
		snapshot.Backspace = true
	case "delete":
		snapshot.Delete = true
	case "home":
		snapshot.Home = true
	case "end":
		snapshot.End = true
	case "space":
		snapshot.Space = true
	case "tab":
		snapshot.Tab = true
	case "escape":
		snapshot.Escape = true
	case "a":
		if control || meta {
			snapshot.SelectAll = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "b":
		if control || meta {
			snapshot.ArrowLeft = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "d":
		if control || meta {
			snapshot.Delete = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "e":
		if control || meta {
			snapshot.End = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "f":
		if control || meta {
			snapshot.ArrowRight = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "h":
		if control || meta {
			snapshot.Backspace = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "n":
		if control || meta {
			snapshot.ArrowDown = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "p":
		if control || meta {
			snapshot.ArrowUp = true
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	case "w":
		if control || meta {
			snapshot.Backspace = true
			snapshot.Control = true
			snapshot.Meta = meta
		} else {
			return ebitenui.InputSnapshot{}, false
		}
	default:
		return ebitenui.InputSnapshot{}, false
	}
	return snapshot, true
}
