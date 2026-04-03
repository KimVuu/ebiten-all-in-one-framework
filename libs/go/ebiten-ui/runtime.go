package ebitenui

type InputSnapshot struct {
	PointerX     float64
	PointerY     float64
	PointerDown  bool
	InputBlocked bool
	ScrollX      float64
	ScrollY      float64
	Text         string
	Backspace    bool
	Delete       bool
	Home         bool
	End          bool
	Submit       bool
	Space        bool
	SelectAll    bool
	Shortcut     string
	Tab          bool
	Escape       bool
	ArrowUp      bool
	ArrowDown    bool
	ArrowLeft    bool
	ArrowRight   bool
	Shift        bool
	Control      bool
	Alt          bool
	Meta         bool
}

type EventContext struct {
	Runtime  *Runtime
	DOM      *DOM
	Node     *Node
	Layout   *LayoutNode
	Viewport Viewport
	Input    InputSnapshot
	X        float64
	Y        float64
	LocalX   float64
	LocalY   float64
	Text     string
	ScrollX  float64
	ScrollY  float64
	Shortcut string
}

type EventHandlers struct {
	OnPointerEnter func(EventContext)
	OnPointerLeave func(EventContext)
	OnPointerMove  func(EventContext)
	OnPointerDown  func(EventContext)
	OnPointerHold  func(EventContext)
	OnPointerUp    func(EventContext)
	OnClick        func(EventContext)
	OnFocus        func(EventContext)
	OnBlur         func(EventContext)
	OnShortcut     func(EventContext)
	OnCursorMove   func(EventContext)
	OnTextInput    func(EventContext)
	OnBackspace    func(EventContext)
	OnDelete       func(EventContext)
	OnSelectAll    func(EventContext)
	OnSubmit       func(EventContext)
	OnScroll       func(EventContext)
}

func (handlers EventHandlers) hasAny() bool {
	return handlers.OnPointerEnter != nil ||
		handlers.OnPointerLeave != nil ||
		handlers.OnPointerMove != nil ||
		handlers.OnPointerDown != nil ||
		handlers.OnPointerHold != nil ||
		handlers.OnPointerUp != nil ||
		handlers.OnClick != nil ||
		handlers.OnFocus != nil ||
		handlers.OnBlur != nil ||
		handlers.OnShortcut != nil ||
		handlers.OnCursorMove != nil ||
		handlers.OnTextInput != nil ||
		handlers.OnBackspace != nil ||
		handlers.OnDelete != nil ||
		handlers.OnSelectAll != nil ||
		handlers.OnSubmit != nil ||
		handlers.OnScroll != nil
}

type TextSelection struct {
	Start int
	End   int
}

type Runtime struct {
	layout         *LayoutNode
	viewport       Viewport
	prevInput      InputSnapshot
	hoveredID      string
	pressedID      string
	focusedID      string
	bools          map[string]bool
	numbers        map[string]float64
	textValues     map[string]string
	textCursors    map[string]int
	textSelections map[string]TextSelection
}

func NewRuntime() *Runtime {
	return &Runtime{
		bools:          map[string]bool{},
		numbers:        map[string]float64{},
		textValues:     map[string]string{},
		textCursors:    map[string]int{},
		textSelections: map[string]TextSelection{},
	}
}

func (runtime *Runtime) Update(dom *DOM, viewport Viewport, input InputSnapshot) *LayoutNode {
	if runtime == nil || dom == nil || dom.Root == nil {
		return nil
	}

	runtime.viewport = viewport
	runtime.layout = dom.Layout(viewport)
	clearTransientStates(dom.Root)

	if input.InputBlocked {
		runtime.hoveredID = ""
		runtime.pressedID = ""
		runtime.prevInput = InputSnapshot{InputBlocked: true}
		return runtime.layout
	}

	hovered := findInteractiveHit(runtime.layout, input.PointerX, input.PointerY)
	hoveredID := ""
	if hovered != nil && hovered.Node != nil {
		hoveredID = hovered.Node.Props.ID
	}

	if runtime.hoveredID != hoveredID {
		runtime.dispatchByID(dom, runtime.hoveredID, input, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnPointerLeave != nil {
				handlers.OnPointerLeave(ctx)
			}
		})
		runtime.dispatchLayout(dom, hovered, input, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnPointerEnter != nil {
				handlers.OnPointerEnter(ctx)
			}
		})
	}
	runtime.hoveredID = hoveredID

	if hovered != nil {
		hovered.Node.Props.State.Hovered = true
	}

	pointerMoved := input.PointerX != runtime.prevInput.PointerX || input.PointerY != runtime.prevInput.PointerY
	if pointerMoved {
		if input.PointerDown && runtime.pressedID != "" {
			runtime.dispatchByID(dom, runtime.pressedID, input, func(handlers EventHandlers, ctx EventContext) {
				if handlers.OnPointerMove != nil {
					handlers.OnPointerMove(ctx)
				}
			})
		} else if hovered != nil {
			runtime.dispatchLayout(dom, hovered, input, func(handlers EventHandlers, ctx EventContext) {
				if handlers.OnPointerMove != nil {
					handlers.OnPointerMove(ctx)
				}
			})
		}
	}

	if !runtime.prevInput.PointerDown && input.PointerDown {
		if hovered != nil && hovered.Node != nil {
			runtime.pressedID = hovered.Node.Props.ID
			runtime.setFocus(dom, hovered.Node.Props.ID, input)
			runtime.dispatchLayout(dom, hovered, input, func(handlers EventHandlers, ctx EventContext) {
				if handlers.OnPointerDown != nil {
					handlers.OnPointerDown(ctx)
				}
			})
		} else {
			runtime.setFocus(dom, "", input)
		}
	}

	if runtime.pressedID != "" {
		if pressed, ok := dom.FindByID(runtime.pressedID); ok {
			pressed.Props.State.Pressed = input.PointerDown
		}
	}

	if runtime.prevInput.PointerDown && input.PointerDown && runtime.pressedID != "" {
		runtime.dispatchByID(dom, runtime.pressedID, input, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnPointerHold != nil {
				handlers.OnPointerHold(ctx)
			}
		})
	}

	if runtime.prevInput.PointerDown && !input.PointerDown {
		pressedID := runtime.pressedID
		runtime.dispatchByID(dom, pressedID, input, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnPointerUp != nil {
				handlers.OnPointerUp(ctx)
			}
			if handlers.OnClick != nil && pressedID != "" && pressedID == runtime.hoveredID {
				handlers.OnClick(ctx)
			}
		})
		runtime.pressedID = ""
	}

	input = runtime.handleKeyboardNavigation(dom, input)
	input = runtime.dispatchFocusedText(dom, input)
	runtime.dispatchKeyboardActivation(dom, input)
	runtime.dispatchScroll(dom, hovered, input)
	runtime.applyFocusState(dom)

	runtime.prevInput = input
	return runtime.layout
}

func (runtime *Runtime) Layout() *LayoutNode {
	if runtime == nil {
		return nil
	}
	return runtime.layout
}

func (runtime *Runtime) FocusedID() string {
	if runtime == nil {
		return ""
	}
	return runtime.focusedID
}

func (runtime *Runtime) HoveredID() string {
	if runtime == nil {
		return ""
	}
	return runtime.hoveredID
}

func (runtime *Runtime) FocusNodeByID(dom *DOM, id string, input InputSnapshot) {
	if runtime == nil {
		return
	}
	runtime.setFocus(dom, id, input)
}

func (runtime *Runtime) ClearFocus(dom *DOM, input InputSnapshot) {
	if runtime == nil {
		return
	}
	runtime.setFocus(dom, "", input)
}

func (runtime *Runtime) BoolValueOrDefault(id string, fallback bool) bool {
	if runtime == nil {
		return fallback
	}
	if value, ok := runtime.bools[id]; ok {
		return value
	}
	runtime.bools[id] = fallback
	return fallback
}

func (runtime *Runtime) SetBoolValue(id string, value bool) {
	if runtime == nil {
		return
	}
	runtime.bools[id] = value
}

func (runtime *Runtime) NumberValueOrDefault(id string, fallback float64) float64 {
	if runtime == nil {
		return fallback
	}
	if value, ok := runtime.numbers[id]; ok {
		return value
	}
	runtime.numbers[id] = fallback
	return fallback
}

func (runtime *Runtime) SetNumberValue(id string, value float64) {
	if runtime == nil {
		return
	}
	runtime.numbers[id] = value
}

func (runtime *Runtime) TextValueOrDefault(id string, fallback string) string {
	if runtime == nil {
		return fallback
	}
	if value, ok := runtime.textValues[id]; ok {
		return value
	}
	runtime.textValues[id] = fallback
	return fallback
}

func (runtime *Runtime) SetTextValue(id string, value string) {
	if runtime == nil {
		return
	}
	runtime.textValues[id] = value
	runtime.clampTextCursor(id, value)
}

func (runtime *Runtime) TextCursorOrDefault(id string, fallback string) int {
	if runtime == nil {
		return len([]rune(fallback))
	}
	if cursor, ok := runtime.textCursors[id]; ok {
		return cursor
	}
	cursor := len([]rune(fallback))
	runtime.textCursors[id] = cursor
	return cursor
}

func (runtime *Runtime) SetTextCursor(id string, cursor int) {
	if runtime == nil {
		return
	}
	runtime.textCursors[id] = cursor
}

func (runtime *Runtime) TextSelectionOrDefault(id string) TextSelection {
	if runtime == nil {
		return TextSelection{}
	}
	if selection, ok := runtime.textSelections[id]; ok {
		return selection
	}
	return TextSelection{}
}

func (runtime *Runtime) SetTextSelection(id string, selection TextSelection) {
	if runtime == nil {
		return
	}
	if selection.Start == selection.End {
		delete(runtime.textSelections, id)
		return
	}
	if selection.Start > selection.End {
		selection.Start, selection.End = selection.End, selection.Start
	}
	runtime.textSelections[id] = selection
}

func (runtime *Runtime) clampTextCursor(id string, value string) {
	if runtime == nil {
		return
	}
	cursor := runtime.textCursors[id]
	length := len([]rune(value))
	if cursor < 0 {
		cursor = 0
	}
	if cursor > length {
		cursor = length
	}
	runtime.textCursors[id] = cursor
	if selection, ok := runtime.textSelections[id]; ok {
		selection.Start = clampInt(selection.Start, 0, length)
		selection.End = clampInt(selection.End, 0, length)
		if selection.Start == selection.End {
			delete(runtime.textSelections, id)
		} else {
			if selection.Start > selection.End {
				selection.Start, selection.End = selection.End, selection.Start
			}
			runtime.textSelections[id] = selection
		}
	}
}

func (runtime *Runtime) textValueAndCursor(id string, fallback string) (string, int, TextSelection) {
	value := runtime.TextValueOrDefault(id, fallback)
	cursor := runtime.TextCursorOrDefault(id, value)
	selection := runtime.TextSelectionOrDefault(id)
	return value, cursor, selection
}

func (runtime *Runtime) setFocus(dom *DOM, id string, input InputSnapshot) {
	if runtime.focusedID == id {
		return
	}
	oldID := runtime.focusedID
	runtime.focusedID = id

	runtime.dispatchByID(dom, oldID, input, func(handlers EventHandlers, ctx EventContext) {
		if handlers.OnBlur != nil {
			handlers.OnBlur(ctx)
		}
	})
	runtime.dispatchByID(dom, id, input, func(handlers EventHandlers, ctx EventContext) {
		if handlers.OnFocus != nil {
			handlers.OnFocus(ctx)
		}
	})
}

func (runtime *Runtime) applyFocusState(dom *DOM) {
	if dom == nil {
		return
	}
	if runtime.focusedID == "" {
		return
	}
	if node, ok := dom.FindByID(runtime.focusedID); ok {
		node.Props.State.Focused = true
	}
}

func (runtime *Runtime) dispatchFocusedText(dom *DOM, input InputSnapshot) InputSnapshot {
	if runtime.focusedID == "" {
		return input
	}
	consumed := false
	runtime.dispatchByID(dom, runtime.focusedID, input, func(handlers EventHandlers, ctx EventContext) {
		if input.Shortcut != "" && handlers.OnShortcut != nil {
			ctx.Shortcut = input.Shortcut
			handlers.OnShortcut(ctx)
			input.Shortcut = ""
			consumed = true
		}
		if input.SelectAll && handlers.OnSelectAll != nil {
			handlers.OnSelectAll(ctx)
			input.SelectAll = false
			consumed = true
		}
		if (input.Home || input.End || input.ArrowLeft || input.ArrowRight || input.ArrowUp || input.ArrowDown) && handlers.OnCursorMove != nil {
			handlers.OnCursorMove(ctx)
			input.Home = false
			input.End = false
			input.ArrowLeft = false
			input.ArrowRight = false
			input.ArrowUp = false
			input.ArrowDown = false
			consumed = true
		}
		if input.Text != "" && handlers.OnTextInput != nil {
			ctx.Text = input.Text
			handlers.OnTextInput(ctx)
			input.Text = ""
			consumed = true
		}
		if input.Backspace && handlers.OnBackspace != nil {
			handlers.OnBackspace(ctx)
			input.Backspace = false
			consumed = true
		}
		if input.Delete && handlers.OnDelete != nil {
			handlers.OnDelete(ctx)
			input.Delete = false
			consumed = true
		}
		if input.Submit && handlers.OnSubmit != nil {
			handlers.OnSubmit(ctx)
			input.Submit = false
			consumed = true
		}
	})
	if consumed {
		return input
	}
	return input
}

func (runtime *Runtime) handleKeyboardNavigation(dom *DOM, input InputSnapshot) InputSnapshot {
	if runtime.layout == nil {
		return input
	}
	if input.Escape {
		runtime.setFocus(dom, "", input)
		input.Escape = false
	}
	if input.Tab {
		nextID := nextFocusableLayoutID(runtime.layout, runtime.focusedID, input.Shift)
		if nextID != "" {
			runtime.setFocus(dom, nextID, input)
		}
		input.Tab = false
	}
	return input
}

func (runtime *Runtime) dispatchKeyboardActivation(dom *DOM, input InputSnapshot) {
	if runtime.focusedID == "" {
		return
	}
	layout, ok := runtime.layout.FindByID(runtime.focusedID)
	if !ok || layout == nil || layout.Node == nil {
		return
	}
	handlers := layout.Node.Props.Handlers
	if handlers.OnSubmit != nil {
		return
	}
	if (input.Submit || input.Space) && (handlers.OnClick != nil || layout.Node.Tag == TagButton) {
		runtime.dispatchLayout(dom, layout, input, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnClick != nil {
				handlers.OnClick(ctx)
			}
		})
	}
}

func (runtime *Runtime) dispatchScroll(dom *DOM, hovered *LayoutNode, input InputSnapshot) {
	scrollInput := input
	if scrollInput.ScrollX == 0 && scrollInput.ScrollY == 0 {
		if input.ArrowUp {
			scrollInput.ScrollY = -1
		}
		if input.ArrowDown {
			scrollInput.ScrollY = 1
		}
		if input.ArrowLeft {
			scrollInput.ScrollX = -1
		}
		if input.ArrowRight {
			scrollInput.ScrollX = 1
		}
		if (scrollInput.ScrollX != 0 || scrollInput.ScrollY != 0) && runtime.focusedID != "" {
			if focused, ok := runtime.layout.FindByID(runtime.focusedID); ok && focused != nil {
				if target := nearestScrollTarget(runtime.layout, focused); target != nil {
					scrollInput.PointerX = target.Frame.X + target.Frame.Width*0.5
					scrollInput.PointerY = target.Frame.Y + target.Frame.Height*0.5
					hovered = target
				}
			}
		}
	}
	if scrollInput.ScrollX == 0 && scrollInput.ScrollY == 0 {
		return
	}
	if target := findScrollTarget(runtime.layout, scrollInput.PointerX, scrollInput.PointerY); target != nil {
		runtime.dispatchLayout(dom, target, scrollInput, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnScroll != nil {
				handlers.OnScroll(ctx)
			}
		})
		return
	}
	if hovered != nil {
		runtime.dispatchLayout(dom, hovered, scrollInput, func(handlers EventHandlers, ctx EventContext) {
			if handlers.OnScroll != nil {
				handlers.OnScroll(ctx)
			}
		})
		return
	}
	runtime.dispatchByID(dom, runtime.focusedID, scrollInput, func(handlers EventHandlers, ctx EventContext) {
		if handlers.OnScroll != nil {
			handlers.OnScroll(ctx)
		}
	})
}

func (runtime *Runtime) dispatchByID(dom *DOM, id string, input InputSnapshot, fn func(EventHandlers, EventContext)) {
	if id == "" || runtime.layout == nil {
		return
	}
	layout, ok := runtime.layout.FindByID(id)
	if !ok || layout.Node == nil {
		return
	}
	runtime.dispatchLayout(dom, layout, input, fn)
}

func (runtime *Runtime) dispatchLayout(dom *DOM, layout *LayoutNode, input InputSnapshot, fn func(EventHandlers, EventContext)) {
	if layout == nil || layout.Node == nil {
		return
	}
	handlers := layout.Node.Props.Handlers
	if !handlers.hasAny() {
		return
	}
	ctx := EventContext{
		Runtime:  runtime,
		DOM:      dom,
		Node:     layout.Node,
		Layout:   layout,
		Viewport: runtime.viewport,
		Input:    input,
		X:        input.PointerX,
		Y:        input.PointerY,
		LocalX:   input.PointerX - layout.Frame.X,
		LocalY:   input.PointerY - layout.Frame.Y,
		Text:     input.Text,
		ScrollX:  input.ScrollX,
		ScrollY:  input.ScrollY,
		Shortcut: input.Shortcut,
	}
	fn(handlers, ctx)
}

func clearTransientStates(node *Node) {
	if node == nil {
		return
	}
	state := node.Props.State
	state.Hovered = false
	state.Pressed = false
	state.Focused = false
	node.Props.State = state
	for _, child := range node.Children {
		clearTransientStates(child)
	}
}

func findInteractiveHit(layout *LayoutNode, x, y float64) *LayoutNode {
	hitPath := findHitPath(layout, x, y)
	for _, node := range hitPath {
		if isInteractiveNode(node) {
			return node
		}
	}
	return nil
}

func findScrollTarget(layout *LayoutNode, x, y float64) *LayoutNode {
	hitPath := findHitPath(layout, x, y)
	for _, node := range hitPath {
		if node == nil || node.Node == nil {
			continue
		}
		if node.Node.Props.Handlers.OnScroll != nil {
			return node
		}
	}
	return nil
}

func nearestScrollTarget(root *LayoutNode, target *LayoutNode) *LayoutNode {
	if root == nil || target == nil || target.Node == nil {
		return nil
	}
	path := findPathByID(root, target.Node.Props.ID)
	for _, node := range path {
		if node == nil || node.Node == nil {
			continue
		}
		if node.Node.Tag == TagScrollView || node.Node.Props.Handlers.OnScroll != nil {
			return node
		}
	}
	return nil
}

func findPathByID(layout *LayoutNode, id string) []*LayoutNode {
	if layout == nil || layout.Node == nil {
		return nil
	}
	if layout.Node.Props.ID == id {
		return []*LayoutNode{layout}
	}
	for _, child := range layout.Children {
		if path := findPathByID(child, id); len(path) > 0 {
			return append(path, layout)
		}
	}
	return nil
}

func findHitPath(layout *LayoutNode, x, y float64) []*LayoutNode {
	if layout == nil || layout.Node == nil {
		return nil
	}
	if layout.ClipChildren && !containsPoint(layout.Frame, x, y) {
		return nil
	}

	for i := len(layout.Children) - 1; i >= 0; i-- {
		if path := findHitPath(layout.Children[i], x, y); len(path) > 0 {
			return append(path, layout)
		}
	}

	if containsPoint(layout.Frame, x, y) {
		return []*LayoutNode{layout}
	}
	return nil
}

func containsPoint(frame Rect, x, y float64) bool {
	return x >= frame.X && x <= frame.X+frame.Width && y >= frame.Y && y <= frame.Y+frame.Height
}

func isInteractiveNode(layout *LayoutNode) bool {
	if layout == nil || layout.Node == nil {
		return false
	}
	if layout.Node.Props.Focusable {
		return true
	}
	if layout.Node.Tag == TagButton {
		return true
	}
	return layout.Node.Props.Handlers.hasAny()
}

func nextFocusableLayoutID(layout *LayoutNode, currentID string, reverse bool) string {
	ids := collectFocusableLayoutIDs(layout)
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

func collectFocusableLayoutIDs(layout *LayoutNode) []string {
	ids := make([]string, 0)
	seen := map[string]bool{}
	var walk func(*LayoutNode)
	walk = func(node *LayoutNode) {
		if node == nil || node.Node == nil {
			return
		}
		if (node.Node.Props.Focusable || node.Node.Tag == TagButton || node.Node.Props.Handlers.hasAny()) && node.Node.Props.ID != "" && !seen[node.Node.Props.ID] {
			seen[node.Node.Props.ID] = true
			ids = append(ids, node.Node.Props.ID)
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(layout)
	return ids
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
