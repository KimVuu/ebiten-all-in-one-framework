package ebitenui

type SemanticSpec struct {
	Screen  string
	Element string
	Role    string
	Slot    string
}

type Point struct {
	X float64
	Y float64
}

type LayoutMode string

const (
	LayoutModeFlowVertical   LayoutMode = "flow-vertical"
	LayoutModeFlowHorizontal LayoutMode = "flow-horizontal"
	LayoutModeAnchored       LayoutMode = "anchored"
	LayoutModeGrid           LayoutMode = "grid"
	LayoutModeStack          LayoutMode = "stack"
)

type LayoutAlignment string

const (
	LayoutAlignmentStart        LayoutAlignment = "start"
	LayoutAlignmentCenter       LayoutAlignment = "center"
	LayoutAlignmentEnd          LayoutAlignment = "end"
	LayoutAlignmentStretch      LayoutAlignment = "stretch"
	LayoutAlignmentSpaceBetween LayoutAlignment = "space-between"
	LayoutAlignmentSpaceAround  LayoutAlignment = "space-around"
	LayoutAlignmentSpaceEvenly  LayoutAlignment = "space-evenly"
)

const (
	LayoutAlignStart        = LayoutAlignmentStart
	LayoutAlignCenter       = LayoutAlignmentCenter
	LayoutAlignEnd          = LayoutAlignmentEnd
	LayoutAlignStretch      = LayoutAlignmentStretch
	LayoutAlignSpaceBetween = LayoutAlignmentSpaceBetween
	LayoutAlignSpaceAround  = LayoutAlignmentSpaceAround
	LayoutAlignSpaceEvenly  = LayoutAlignmentSpaceEvenly
)

type LayoutAutoFlow string

const (
	LayoutAutoFlowRow    LayoutAutoFlow = "row"
	LayoutAutoFlowColumn LayoutAutoFlow = "column"
)

type Anchor string

const (
	AnchorTopLeft     Anchor = "top-left"
	AnchorTop         Anchor = "top"
	AnchorTopRight    Anchor = "top-right"
	AnchorLeft        Anchor = "left"
	AnchorCenter      Anchor = "center"
	AnchorRight       Anchor = "right"
	AnchorBottomLeft  Anchor = "bottom-left"
	AnchorBottom      Anchor = "bottom"
	AnchorBottomRight Anchor = "bottom-right"
)

type Pivot string

const (
	PivotTopLeft     Pivot = "top-left"
	PivotTop         Pivot = "top"
	PivotTopRight    Pivot = "top-right"
	PivotLeft        Pivot = "left"
	PivotCenter      Pivot = "center"
	PivotRight       Pivot = "right"
	PivotBottomLeft  Pivot = "bottom-left"
	PivotBottom      Pivot = "bottom"
	PivotBottomRight Pivot = "bottom-right"
)

type LayoutSize struct {
	Width  Length
	Height Length
}

type LayoutConstraints struct {
	KeepInsideParent bool
	AllowOverlap     bool
	ClipChildren     bool
	MinHitTarget     float64
}

type LayoutGrid struct {
	Columns int
	Rows    int
	Gap     float64

	AutoFlow       LayoutAutoFlow
	JustifyContent LayoutAlignment
	AlignContent   LayoutAlignment
	JustifyItems   LayoutAlignment
	AlignItems     LayoutAlignment
	JustifySelf    LayoutAlignment
	AlignSelf      LayoutAlignment

	ColumnStart int
	ColumnSpan  int
	RowStart    int
	RowSpan     int
}

type LayoutSpec struct {
	Mode        LayoutMode
	ParentID    string
	Anchor      Anchor
	Pivot       Pivot
	Offset      Point
	Size        LayoutSize
	MinSize     LayoutSize
	MaxSize     LayoutSize
	Margin      Insets
	Padding     Insets
	Gap         float64
	ZIndex      int
	Constraints LayoutConstraints
	Grid        LayoutGrid

	GridColumn     int
	GridColumnSpan int
	GridRow        int
	GridRowSpan    int
}

type LayoutOverflow struct {
	Horizontal bool
	Vertical   bool
	Any        bool
}

type LayoutIssueSeverity string

const (
	IssueSeverityInfo    LayoutIssueSeverity = "info"
	IssueSeverityWarning LayoutIssueSeverity = "warning"
	IssueSeverityError   LayoutIssueSeverity = "error"
)

type LayoutIssueCode string

const (
	IssueOutOfViewport     LayoutIssueCode = "out_of_viewport"
	IssueOutOfParent       LayoutIssueCode = "out_of_parent"
	IssueOverlap           LayoutIssueCode = "overlap"
	IssueMinHitTarget      LayoutIssueCode = "min_hit_target"
	IssueTextOverflow      LayoutIssueCode = "text_overflow"
	IssueSafeAreaViolation LayoutIssueCode = "safe_area_violation"
	IssueZOrderConflict    LayoutIssueCode = "z_order_conflict"
)

type ConstraintPatch struct {
	Field string
	Value any
	Note  string
}

type LayoutIssue struct {
	NodeID                     string
	Severity                   LayoutIssueSeverity
	Code                       LayoutIssueCode
	Message                    string
	SuggestedConstraintChanges []ConstraintPatch
}

type LayoutValidationReport struct {
	Issues []LayoutIssue
}

type ValidationOptions struct {
	SafeArea Insets
}
