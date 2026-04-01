package ebitenuidebug

import (
	"sort"
	"strings"
	"sync"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

type debugResolvedTarget struct {
	ID    string
	Frame ebitenui.Rect
	Node  *ebitenui.LayoutNode
	Path  []*ebitenui.LayoutNode
}

func (target debugResolvedTarget) center() (float64, float64) {
	return target.Frame.X + target.Frame.Width*0.5, target.Frame.Y + target.Frame.Height*0.5
}

type debugFrameInputKind string

const (
	debugFrameInputPointerMove   debugFrameInputKind = "pointer-move"
	debugFrameInputPointerDown   debugFrameInputKind = "pointer-down"
	debugFrameInputPointerUp     debugFrameInputKind = "pointer-up"
	debugFrameInputScroll        debugFrameInputKind = "scroll"
	debugFrameInputText          debugFrameInputKind = "text"
	debugFrameInputFocusNode     debugFrameInputKind = "focus-node"
	debugFrameInputClearFocus    debugFrameInputKind = "clear-focus"
	debugFrameInputOverlayToggle debugFrameInputKind = "overlay-toggle"
)

type debugFrameInput struct {
	frame int
	seq   int
	kind  debugFrameInputKind

	input ebitenui.InputSnapshot

	resolvedTarget string
	focusTargetID  string
	clearFocus     bool
	overlayEnabled *bool
}

type debugInputQueue struct {
	mu     sync.Mutex
	nextID int
	items  []debugFrameInput
}

func newDebugInputQueue() *debugInputQueue {
	return &debugInputQueue{}
}

func (queue *debugInputQueue) enqueue(effect debugFrameInput) debugFrameInput {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	queue.nextID++
	effect.seq = queue.nextID
	queue.items = append(queue.items, effect)
	return effect
}

func (queue *debugInputQueue) drain(frame int) []debugFrameInput {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if len(queue.items) == 0 {
		return nil
	}

	due := make([]debugFrameInput, 0, len(queue.items))
	remaining := queue.items[:0]
	for _, item := range queue.items {
		if item.frame <= frame {
			due = append(due, item)
			continue
		}
		remaining = append(remaining, item)
	}
	queue.items = remaining

	sort.SliceStable(due, func(i, j int) bool {
		if due[i].frame != due[j].frame {
			return due[i].frame < due[j].frame
		}
		return due[i].seq < due[j].seq
	})
	return due
}

func (queue *debugInputQueue) clear() {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.items = nil
}

func (queue *debugInputQueue) len() int {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	return len(queue.items)
}

func (queue *debugInputQueue) queueClick(startFrame int, target debugResolvedTarget) []debugFrameInput {
	x, y := target.center()
	return []debugFrameInput{
		queue.enqueue(debugFrameInput{
			frame:          startFrame + 1,
			kind:           debugFrameInputPointerMove,
			input:          ebitenui.InputSnapshot{PointerX: x, PointerY: y},
			resolvedTarget: target.ID,
		}),
		queue.enqueue(debugFrameInput{
			frame:          startFrame + 2,
			kind:           debugFrameInputPointerDown,
			input:          ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: true},
			resolvedTarget: target.ID,
		}),
		queue.enqueue(debugFrameInput{
			frame:          startFrame + 3,
			kind:           debugFrameInputPointerUp,
			input:          ebitenui.InputSnapshot{PointerX: x, PointerY: y, PointerDown: false},
			resolvedTarget: target.ID,
		}),
	}
}

func (queue *debugInputQueue) queueScroll(startFrame int, target debugResolvedTarget, deltaX, deltaY float64) debugFrameInput {
	x, y := target.center()
	return queue.enqueue(debugFrameInput{
		frame:          startFrame + 1,
		kind:           debugFrameInputScroll,
		input:          ebitenui.InputSnapshot{PointerX: x, PointerY: y, ScrollX: deltaX, ScrollY: deltaY},
		resolvedTarget: target.ID,
	})
}

func (queue *debugInputQueue) queueText(startFrame int, targetID string, value string) []debugFrameInput {
	if value == "" {
		return nil
	}

	runes := []rune(value)
	effects := make([]debugFrameInput, 0, len(runes))
	for index, r := range runes {
		effects = append(effects, queue.enqueue(debugFrameInput{
			frame:          startFrame + index + 1,
			kind:           debugFrameInputText,
			input:          ebitenui.InputSnapshot{Text: string(r)},
			resolvedTarget: targetID,
		}))
	}
	return effects
}

func (queue *debugInputQueue) queueFocus(startFrame int, targetID string) debugFrameInput {
	return queue.enqueue(debugFrameInput{
		frame:          startFrame + 1,
		kind:           debugFrameInputFocusNode,
		focusTargetID:  targetID,
		resolvedTarget: targetID,
	})
}

func (queue *debugInputQueue) queueClearFocus(startFrame int) debugFrameInput {
	return queue.enqueue(debugFrameInput{
		frame:      startFrame + 1,
		kind:       debugFrameInputClearFocus,
		clearFocus: true,
	})
}

func (queue *debugInputQueue) queueOverlayToggle(startFrame int, enabled bool) debugFrameInput {
	return queue.enqueue(debugFrameInput{
		frame: startFrame + 1,
		kind:  debugFrameInputOverlayToggle,
		overlayEnabled: func(value bool) *bool {
			return &value
		}(enabled),
	})
}

func resolveExactTarget(layout *ebitenui.LayoutNode, nodeID string) (debugResolvedTarget, bool) {
	path := findLayoutPath(layout, nodeID)
	if len(path) == 0 {
		return debugResolvedTarget{}, false
	}

	leaf := path[len(path)-1]
	frame := leaf.Frame
	if leaf.ClickableRect != (ebitenui.Rect{}) {
		frame = leaf.ClickableRect
	}
	return debugResolvedTarget{
		ID:    leaf.Node.Props.ID,
		Frame: frame,
		Node:  leaf,
		Path:  path,
	}, true
}

func resolveDebugTarget(layout *ebitenui.LayoutNode, nodeID string) (debugResolvedTarget, bool) {
	path := findLayoutPath(layout, nodeID)
	if len(path) == 0 {
		return debugResolvedTarget{}, false
	}

	for index := len(path) - 1; index >= 0; index-- {
		node := path[index]
		if isInteractiveLayoutNode(node) {
			frame := node.ClickableRect
			if frame == (ebitenui.Rect{}) {
				frame = node.Frame
			}
			return debugResolvedTarget{
				ID:    node.Node.Props.ID,
				Frame: frame,
				Node:  node,
				Path:  path[:index+1],
			}, true
		}
	}

	leaf := path[len(path)-1]
	return debugResolvedTarget{
		ID:    leaf.Node.Props.ID,
		Frame: leaf.Frame,
		Node:  leaf,
		Path:  path,
	}, true
}

func resolveScrollTarget(layout *ebitenui.LayoutNode, nodeID string) (debugResolvedTarget, bool) {
	if nodeID == "" {
		return findFirstScrollTarget(layout)
	}

	path := findLayoutPath(layout, nodeID)
	if len(path) == 0 {
		return debugResolvedTarget{}, false
	}

	for index := len(path) - 1; index >= 0; index-- {
		node := path[index]
		if isScrollLayoutNode(node) {
			frame := node.ClipRect
			if frame == (ebitenui.Rect{}) {
				frame = node.Frame
			}
			return debugResolvedTarget{
				ID:    node.Node.Props.ID,
				Frame: frame,
				Node:  node,
				Path:  path[:index+1],
			}, true
		}
	}

	return resolveDebugTarget(layout, nodeID)
}

func findFirstScrollTarget(layout *ebitenui.LayoutNode) (debugResolvedTarget, bool) {
	if layout == nil || layout.Node == nil {
		return debugResolvedTarget{}, false
	}
	if isScrollLayoutNode(layout) {
		frame := layout.ClipRect
		if frame == (ebitenui.Rect{}) {
			frame = layout.Frame
		}
		return debugResolvedTarget{
			ID:    layout.Node.Props.ID,
			Frame: frame,
			Node:  layout,
			Path:  []*ebitenui.LayoutNode{layout},
		}, true
	}
	for _, child := range layout.Children {
		if target, ok := findFirstScrollTarget(child); ok {
			return target, true
		}
	}
	return debugResolvedTarget{}, false
}

func findLayoutPath(layout *ebitenui.LayoutNode, nodeID string) []*ebitenui.LayoutNode {
	if layout == nil || layout.Node == nil {
		return nil
	}
	if layout.Node.Props.ID == nodeID {
		return []*ebitenui.LayoutNode{layout}
	}
	for _, child := range layout.Children {
		if path := findLayoutPath(child, nodeID); len(path) > 0 {
			return append([]*ebitenui.LayoutNode{layout}, path...)
		}
	}
	return nil
}

func isInteractiveLayoutNode(layout *ebitenui.LayoutNode) bool {
	if layout == nil || layout.Node == nil {
		return false
	}
	if layout.Node.Tag == ebitenui.TagButton || layout.Node.Props.Focusable {
		return true
	}
	handlers := layout.Node.Props.Handlers
	return handlers.OnClick != nil ||
		handlers.OnFocus != nil ||
		handlers.OnSubmit != nil ||
		handlers.OnTextInput != nil ||
		handlers.OnBackspace != nil ||
		handlers.OnScroll != nil
}

func isScrollLayoutNode(layout *ebitenui.LayoutNode) bool {
	if layout == nil || layout.Node == nil {
		return false
	}
	return layout.Node.Tag == ebitenui.TagScrollView || layout.Node.Props.Handlers.OnScroll != nil
}

func collectFocusableIDs(layout *ebitenui.LayoutNode) []string {
	ids := make([]string, 0)
	visited := map[string]bool{}

	var walk func(*ebitenui.LayoutNode)
	walk = func(node *ebitenui.LayoutNode) {
		if node == nil || node.Node == nil {
			return
		}
		if isInteractiveLayoutNode(node) || node.Node.Props.Focusable {
			if node.Node.Props.ID != "" && !visited[node.Node.Props.ID] {
				visited[node.Node.Props.ID] = true
				ids = append(ids, node.Node.Props.ID)
			}
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(layout)
	return ids
}

func nextFocusableID(layout *ebitenui.LayoutNode, currentID string, reverse bool) string {
	ids := collectFocusableIDs(layout)
	if len(ids) == 0 {
		return ""
	}
	if currentID == "" {
		if reverse {
			return ids[len(ids)-1]
		}
		return ids[0]
	}

	index := -1
	for i, id := range ids {
		if id == currentID {
			index = i
			break
		}
	}
	if index == -1 {
		if reverse {
			return ids[len(ids)-1]
		}
		return ids[0]
	}

	if reverse {
		index--
		if index < 0 {
			index = len(ids) - 1
		}
		return ids[index]
	}

	index++
	if index >= len(ids) {
		index = 0
	}
	return ids[index]
}

func firstFocusableID(layout *ebitenui.LayoutNode) string {
	ids := collectFocusableIDs(layout)
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}

func targetIDs(effects []debugFrameInput) []string {
	ids := make([]string, 0, len(effects))
	for _, effect := range effects {
		if strings.TrimSpace(effect.resolvedTarget) != "" {
			ids = append(ids, effect.resolvedTarget)
		}
	}
	return ids
}
