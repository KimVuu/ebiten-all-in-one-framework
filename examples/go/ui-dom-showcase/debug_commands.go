package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kimyechan/ebiten-aio-framework/libs/go/ebitendebug"
	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func (g *game) registerDebugCommands(bridge *ebitendebug.Bridge) {
	bridge.RegisterCommand("validate_ui_layout", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		viewport := g.currentViewport()
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

		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		target, ok := resolveDebugTarget(layout, nodeID)
		if !ok {
			return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unknown node: %s", nodeID)}
		}

		viewport := g.currentViewport()
		report := buildDebugLayoutReport(layout, viewport)
		snapshot := convertLayoutToUISnapshot(target.Node, parentLayoutForPath(target.Path), report, viewport, g.overlayEnabled)
		return ebitendebug.CommandResult{
			Success: true,
			Payload: map[string]any{
				"node": snapshot,
			},
		}
	})

	bridge.RegisterCommand("suggest_ui_constraint_fixes", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		viewport := g.currentViewport()
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

		g.mu.Lock()
		g.overlayEnabled = enabled
		g.mu.Unlock()

		return ebitendebug.CommandResult{
			Success: true,
			Payload: map[string]any{
				"overlayEnabled": enabled,
			},
		}
	})

	bridge.RegisterCommand("ui_click", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		frame := g.currentFrame()
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		target, resolved, err := resolvePointerTarget(layout, request.Args)
		if err != nil {
			return ebitendebug.CommandResult{Success: false, Message: err.Error()}
		}

		queued := g.debugQueue.queueClick(frame, target)
		return debugQueueResult("click", resolved, queued)
	})

	bridge.RegisterCommand("ui_pointer_move", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		frame := g.currentFrame()
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		target, resolved, err := resolvePointerTarget(layout, request.Args)
		if err != nil {
			return ebitendebug.CommandResult{Success: false, Message: err.Error()}
		}

		queued := g.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerMove,
			input:          uidom.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5},
			resolvedTarget: resolved,
		})
		return debugQueueResult("pointer_move", resolved, []debugFrameInput{queued})
	})

	bridge.RegisterCommand("ui_pointer_down", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		frame := g.currentFrame()
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		target, resolved, err := resolvePointerTarget(layout, request.Args)
		if err != nil {
			return ebitendebug.CommandResult{Success: false, Message: err.Error()}
		}

		queued := g.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerDown,
			input:          uidom.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5, PointerDown: true},
			resolvedTarget: resolved,
		})
		return debugQueueResult("pointer_down", resolved, []debugFrameInput{queued})
	})

	bridge.RegisterCommand("ui_pointer_up", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		frame := g.currentFrame()
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		target, resolved, err := resolvePointerTarget(layout, request.Args)
		if err != nil {
			return ebitendebug.CommandResult{Success: false, Message: err.Error()}
		}

		queued := g.debugQueue.enqueue(debugFrameInput{
			frame:          frame + 1,
			kind:           debugFrameInputPointerUp,
			input:          uidom.InputSnapshot{PointerX: target.Frame.X + target.Frame.Width*0.5, PointerY: target.Frame.Y + target.Frame.Height*0.5, PointerDown: false},
			resolvedTarget: resolved,
		})
		return debugQueueResult("pointer_up", resolved, []debugFrameInput{queued})
	})

	bridge.RegisterCommand("ui_scroll", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		frame := g.currentFrame()
		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		deltaX, _ := commandFloat(request.Args, "delta_x")
		deltaY, _ := commandFloat(request.Args, "delta_y")
		target, resolved, err := resolveScrollTargetFromArgs(layout, request.Args)
		if err != nil {
			return ebitendebug.CommandResult{Success: false, Message: err.Error()}
		}

		queued := g.debugQueue.queueScroll(frame, target, deltaX, deltaY)
		return debugQueueResult("scroll", resolved, []debugFrameInput{queued})
	})

	bridge.RegisterCommand("ui_focus_node", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		nodeID := commandNodeID(request.Args)
		if nodeID == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_focus_node requires node_id"}
		}

		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}
		target, ok := resolveDebugTarget(layout, nodeID)
		if !ok {
			return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unknown node: %s", nodeID)}
		}

		queued := g.debugQueue.queueFocus(g.currentFrame(), target.ID)
		return debugQueueResult("focus", target.ID, []debugFrameInput{queued})
	})

	bridge.RegisterCommand("ui_type_text", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		value, ok := commandString(request.Args, "text")
		if !ok || value == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_type_text requires text"}
		}

		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		nodeID := commandNodeID(request.Args)
		targetID := fallbackString(nodeID, g.runtime.FocusedID())
		startFrame := g.currentFrame()
		queue := g.debugQueue
		queuedFrames := make([]int, 0, len(value)+1)

		if nodeID != "" {
			target, ok := resolveDebugTarget(layout, nodeID)
			if !ok {
				return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unknown node: %s", nodeID)}
			}
			targetID = target.ID
			queue.queueFocus(startFrame, targetID)
			queuedFrames = append(queuedFrames, startFrame+1)
			startFrame++
		} else if g.runtime.FocusedID() == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_type_text requires a focused node or node_id"}
		}

		for _, effect := range queue.queueText(startFrame, targetID, value) {
			queuedFrames = append(queuedFrames, effect.frame)
		}

		return ebitendebug.CommandResult{
			Success:        true,
			Status:         "queued",
			ResolvedTarget: targetID,
			QueuedFrame:    queuedFrames[len(queuedFrames)-1],
			Payload: map[string]any{
				"status":         "queued",
				"resolvedTarget": targetID,
				"queuedFrames":   queuedFrames,
			},
		}
	})

	bridge.RegisterCommand("ui_key_event", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		keyName, ok := commandString(request.Args, "key")
		if !ok || keyName == "" {
			return ebitendebug.CommandResult{Success: false, Message: "ui_key_event requires key"}
		}

		layout := g.currentLayout()
		if layout == nil {
			return ebitendebug.CommandResult{Success: false, Message: "ui layout unavailable"}
		}

		startFrame := g.currentFrame()
		shift, _ := commandBool(request.Args, "shift")
		control, _ := commandBool(request.Args, "control")
		if !control {
			control, _ = commandBool(request.Args, "ctrl")
		}
		alt, _ := commandBool(request.Args, "alt")
		meta, _ := commandBool(request.Args, "meta")
		nodeID := commandNodeID(request.Args)
		targetID := fallbackString(nodeID, g.runtime.FocusedID())
		if targetID == "" {
			targetID = firstFocusableID(layout)
		}

		queuedFrames := make([]int, 0, 2)
		switch normalizeKeyName(keyName) {
		case "enter":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{Submit: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "backspace":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{Backspace: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "delete":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{Delete: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "home":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{Home: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "end":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{End: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "space":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{Space: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "selectall":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame:          startFrame + 1,
				kind:           debugFrameInputText,
				input:          uidom.InputSnapshot{SelectAll: true, Shift: shift, Control: control, Alt: alt, Meta: meta},
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		case "tab":
			nextID := nextFocusableID(layout, targetID, shift)
			if nextID == "" {
				return ebitendebug.CommandResult{Success: false, Message: "no focusable node available"}
			}
			effect := g.debugQueue.queueFocus(startFrame, nextID)
			queuedFrames = append(queuedFrames, effect.frame)
			targetID = nextID
		case "escape":
			effect := g.debugQueue.queueClearFocus(startFrame)
			queuedFrames = append(queuedFrames, effect.frame)
			targetID = ""
		case "arrowup", "arrowdown", "arrowleft", "arrowright":
			if nodeID != "" {
				if target, ok := resolveDebugTarget(layout, nodeID); ok {
					focus := g.debugQueue.queueFocus(startFrame, target.ID)
					queuedFrames = append(queuedFrames, focus.frame)
					startFrame++
					targetID = target.ID
				}
			}
			effect := g.debugQueue.enqueue(debugFrameInput{
				frame: startFrame + 1,
				kind:  debugFrameInputText,
				input: func() uidom.InputSnapshot {
					snapshot, ok := shortcutAwareKeyEventSnapshot(normalizeKeyName(keyName), shift, control, alt, meta)
					if !ok {
						return uidom.InputSnapshot{}
					}
					return snapshot
				}(),
				resolvedTarget: targetID,
			})
			queuedFrames = append(queuedFrames, effect.frame)
		default:
			if len([]rune(keyName)) == 1 {
				if nodeID != "" {
					if target, ok := resolveDebugTarget(layout, nodeID); ok {
						focus := g.debugQueue.queueFocus(startFrame, target.ID)
						queuedFrames = append(queuedFrames, focus.frame)
						startFrame++
						targetID = target.ID
					}
				}
				if control || meta {
					snapshot, ok := shortcutAwareKeyEventSnapshot(strings.ToLower(keyName), shift, control, alt, meta)
					if !ok {
						return ebitendebug.CommandResult{Success: false, Message: fmt.Sprintf("unsupported shortcut: %s", keyName)}
					}
					effect := g.debugQueue.enqueue(debugFrameInput{
						frame:          startFrame + 1,
						kind:           debugFrameInputText,
						input:          snapshot,
						resolvedTarget: targetID,
					})
					queuedFrames = append(queuedFrames, effect.frame)
				} else {
					queued := g.debugQueue.queueText(startFrame, targetID, keyName)
					for _, effect := range queued {
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
			Payload: map[string]any{
				"status":         "queued",
				"resolvedTarget": targetID,
				"queuedFrames":   queuedFrames,
			},
		}
	})

	bridge.RegisterCommand("ui_clear_input_queue", func(request ebitendebug.CommandRequest) ebitendebug.CommandResult {
		g.debugQueue.clear()
		return ebitendebug.CommandResult{
			Success: true,
			Status:  "cleared",
			Payload: map[string]any{
				"status": "cleared",
			},
		}
	})
}

func (g *game) currentFrame() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.frame
}

func (g *game) applyQueuedDebugEffects(frame int, dom *uidom.DOM, layout *uidom.LayoutNode, input uidom.InputSnapshot) uidom.InputSnapshot {
	effects := g.debugQueue.drain(frame)
	if len(effects) == 0 {
		return input
	}

	for _, effect := range effects {
		switch effect.kind {
		case debugFrameInputPointerMove, debugFrameInputPointerDown, debugFrameInputPointerUp, debugFrameInputScroll, debugFrameInputText:
			input = mergeDebugInputSnapshot(input, effect.input, effect.kind)
		case debugFrameInputFocusNode:
			if effect.focusTargetID != "" {
				g.runtime.FocusNodeByID(dom, effect.focusTargetID, input)
			}
		case debugFrameInputClearFocus:
			g.runtime.ClearFocus(dom, input)
		case debugFrameInputOverlayToggle:
			if effect.overlayEnabled != nil {
				g.overlayEnabled = *effect.overlayEnabled
			}
		}
	}

	return input
}

func mergeDebugInputSnapshot(base, add uidom.InputSnapshot, kind debugFrameInputKind) uidom.InputSnapshot {
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

func resolvePointerTarget(layout *uidom.LayoutNode, args map[string]any) (debugResolvedTarget, string, error) {
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
		frame := uidom.Rect{X: x, Y: y}
		return debugResolvedTarget{ID: "pointer", Frame: frame}, "pointer", nil
	}

	target, ok := findFirstInteractiveTarget(layout)
	if !ok {
		return debugResolvedTarget{}, "", fmt.Errorf("no interactive target available")
	}
	return target, target.ID, nil
}

func resolveScrollTargetFromArgs(layout *uidom.LayoutNode, args map[string]any) (debugResolvedTarget, string, error) {
	x, hasX := commandFloat(args, "x")
	y, hasY := commandFloat(args, "y")
	if hasX && hasY {
		frame := uidom.Rect{X: x, Y: y}
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

func findFirstInteractiveTarget(layout *uidom.LayoutNode) (debugResolvedTarget, bool) {
	if layout == nil || layout.Node == nil {
		return debugResolvedTarget{}, false
	}
	if isInteractiveLayoutNode(layout) {
		frame := layout.ClickableRect
		if frame == (uidom.Rect{}) {
			frame = layout.Frame
		}
		return debugResolvedTarget{ID: layout.Node.Props.ID, Frame: frame, Node: layout, Path: []*uidom.LayoutNode{layout}}, true
	}
	for _, child := range layout.Children {
		if target, ok := findFirstInteractiveTarget(child); ok {
			return target, true
		}
	}
	return debugResolvedTarget{}, false
}

func parentLayoutForPath(path []*uidom.LayoutNode) *uidom.LayoutNode {
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
		default:
			return false, false
		}
	default:
		return false, false
	}
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
	default:
		return 0, false
	}
}

func normalizeKeyName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func shortcutAwareKeyEventSnapshot(key string, shift, control, alt, meta bool) (uidom.InputSnapshot, bool) {
	snapshot := uidom.InputSnapshot{
		Shift:   shift,
		Control: control,
		Alt:     alt,
		Meta:    meta,
	}

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
	case "selectall":
		snapshot.SelectAll = true
	case "a":
		snapshot.SelectAll = true
	case "b":
		snapshot.ArrowLeft = true
	case "d":
		snapshot.Delete = true
	case "e":
		snapshot.End = true
	case "f":
		snapshot.ArrowRight = true
	case "h":
		snapshot.Backspace = true
	case "n":
		snapshot.ArrowDown = true
	case "p":
		snapshot.ArrowUp = true
	case "w":
		snapshot.Backspace = true
	default:
		return uidom.InputSnapshot{}, false
	}

	return snapshot, true
}
