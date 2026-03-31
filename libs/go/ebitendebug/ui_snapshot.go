package ebitendebug

type LayoutMode string

const (
	LayoutModeFlowVertical   LayoutMode = "flow-vertical"
	LayoutModeFlowHorizontal LayoutMode = "flow-horizontal"
	LayoutModeAnchored       LayoutMode = "anchored"
	LayoutModeGrid           LayoutMode = "grid"
	LayoutModeStack          LayoutMode = "stack"
)

type UIAnchor string

type UIPivot string

type UIPositionSnapshot struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type UISizeSnapshot struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type UIInsetsSnapshot struct {
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
}

type UISafeAreaSnapshot = UIInsetsSnapshot

type UIGridSnapshot struct {
	Columns int                `json:"columns,omitempty"`
	Rows    int                `json:"rows,omitempty"`
	Gap     UIPositionSnapshot `json:"gap,omitempty"`
}

type UIConstraintSnapshot struct {
	Field string `json:"field"`
	Op    string `json:"op,omitempty"`
	Value any    `json:"value,omitempty"`
}

type UIConstraintPatchSnapshot = UIConstraintSnapshot

type UISemanticSnapshot struct {
	Screen  string `json:"screen,omitempty"`
	Element string `json:"element,omitempty"`
	Role    string `json:"role,omitempty"`
	Slot    string `json:"slot,omitempty"`
}

type UILayoutSnapshot struct {
	Mode        LayoutMode             `json:"mode,omitempty"`
	ParentID    string                 `json:"parentId,omitempty"`
	Anchor      UIAnchor               `json:"anchor,omitempty"`
	Pivot       UIPivot                `json:"pivot,omitempty"`
	Offset      UIPositionSnapshot     `json:"offset,omitempty"`
	Size        UISizeSnapshot         `json:"size,omitempty"`
	MinSize     *UISizeSnapshot        `json:"minSize,omitempty"`
	MaxSize     *UISizeSnapshot        `json:"maxSize,omitempty"`
	Margin      *UIInsetsSnapshot      `json:"margin,omitempty"`
	Padding     *UIInsetsSnapshot      `json:"padding,omitempty"`
	ZIndex      int                    `json:"zIndex,omitempty"`
	Constraints []UIConstraintSnapshot `json:"constraints,omitempty"`
	Grid        *UIGridSnapshot        `json:"grid,omitempty"`
}

type UIOverflowSnapshot struct {
	Top    bool `json:"top,omitempty"`
	Right  bool `json:"right,omitempty"`
	Bottom bool `json:"bottom,omitempty"`
	Left   bool `json:"left,omitempty"`
}

type UIComputedSnapshot struct {
	Bounds        Rect                `json:"bounds"`
	ParentBounds  *Rect               `json:"parentBounds,omitempty"`
	ContentBounds *Rect               `json:"contentBounds,omitempty"`
	ClickableRect *Rect               `json:"clickableRect,omitempty"`
	ClipRect      *Rect               `json:"clipRect,omitempty"`
	Visible       bool                `json:"visible"`
	Overflow      *UIOverflowSnapshot `json:"overflow,omitempty"`
}

type UIIssueSnapshot struct {
	NodeID                     string                 `json:"nodeId,omitempty"`
	Severity                   string                 `json:"severity"`
	Code                       string                 `json:"code"`
	Message                    string                 `json:"message"`
	SuggestedConstraintChanges []UIConstraintSnapshot `json:"suggestedConstraintChanges,omitempty"`
}

type UIIssueSummarySnapshot struct {
	Total        int `json:"total"`
	Errors       int `json:"errors"`
	Warnings     int `json:"warnings"`
	Info         int `json:"info"`
	InvalidNodes int `json:"invalidNodes"`
}

type UIViewportSnapshot struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Scale  float64 `json:"scale,omitempty"`
}

type UIPointerSnapshot struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Down bool    `json:"down,omitempty"`
}

type UIScrollSnapshot struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Source string  `json:"source,omitempty"`
}

type UIKeyboardSnapshot struct {
	Text        string   `json:"text,omitempty"`
	KeysDown    []string `json:"keysDown,omitempty"`
	KeysPressed []string `json:"keysPressed,omitempty"`
	Modifiers   []string `json:"modifiers,omitempty"`
	QueueDepth  int      `json:"queueDepth,omitempty"`
}

type UIInputSnapshot struct {
	FocusedNodeID string              `json:"focusedNodeId,omitempty"`
	HoveredNodeID string              `json:"hoveredNodeId,omitempty"`
	Pointer       *UIPointerSnapshot  `json:"pointer,omitempty"`
	Scroll        *UIScrollSnapshot   `json:"scroll,omitempty"`
	Keyboard      *UIKeyboardSnapshot `json:"keyboard,omitempty"`
}

type UINodeSnapshot struct {
	ID       string              `json:"id,omitempty"`
	Type     string              `json:"type"`
	Text     string              `json:"text,omitempty"`
	Visible  bool                `json:"visible"`
	Enabled  bool                `json:"enabled,omitempty"`
	ParentID string              `json:"parentId,omitempty"`
	Semantic *UISemanticSnapshot `json:"semantic,omitempty"`
	Layout   *UILayoutSnapshot   `json:"layout,omitempty"`
	Computed *UIComputedSnapshot `json:"computed,omitempty"`
	Bounds   Rect                `json:"bounds"`
	Issues   []UIIssueSnapshot   `json:"issues,omitempty"`
	Props    map[string]any      `json:"props,omitempty"`
	Children []UINodeSnapshot    `json:"children,omitempty"`
}

type UISnapshot struct {
	Width            float64                `json:"width"`
	Height           float64                `json:"height"`
	Viewport         UIViewportSnapshot     `json:"viewport,omitempty"`
	SafeArea         UIInsetsSnapshot       `json:"safeArea,omitempty"`
	IssueSummary     UIIssueSummarySnapshot `json:"issueSummary,omitempty"`
	InvalidNodeCount int                    `json:"invalidNodeCount,omitempty"`
	InputState       UIInputSnapshot        `json:"inputState,omitempty"`
	Root             UINodeSnapshot         `json:"root"`
}
