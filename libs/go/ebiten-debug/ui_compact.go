package ebitendebug

type UINodeSummarySnapshot struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	Role        string `json:"role,omitempty"`
	Slot        string `json:"slot,omitempty"`
	Bounds      Rect   `json:"bounds"`
	Visible     bool   `json:"visible"`
	Enabled     bool   `json:"enabled,omitempty"`
	ChildCount  int    `json:"childCount,omitempty"`
	IssueCount  int    `json:"issueCount,omitempty"`
	Interactive bool   `json:"interactive,omitempty"`
	Scrollable  bool   `json:"scrollable,omitempty"`
	TextPreview string `json:"textPreview,omitempty"`
}

type UIOverviewSnapshot struct {
	Viewport         UIViewportSnapshot      `json:"viewport,omitempty"`
	SafeArea         UIInsetsSnapshot        `json:"safeArea,omitempty"`
	RootID           string                  `json:"rootId,omitempty"`
	FocusedNodeID    string                  `json:"focusedNodeId,omitempty"`
	HoveredNodeID    string                  `json:"hoveredNodeId,omitempty"`
	TotalNodeCount   int                     `json:"totalNodeCount,omitempty"`
	VisibleNodeCount int                     `json:"visibleNodeCount,omitempty"`
	InvalidNodeCount int                     `json:"invalidNodeCount,omitempty"`
	TopLevelSections []UINodeSummarySnapshot `json:"topLevelSections,omitempty"`
	IssueSummary     UIIssueSummarySnapshot  `json:"issueSummary,omitempty"`
}

type UIQueryRequest struct {
	ID              string `json:"id,omitempty"`
	Role            string `json:"role,omitempty"`
	Slot            string `json:"slot,omitempty"`
	Type            string `json:"type,omitempty"`
	TextContains    string `json:"text_contains,omitempty"`
	VisibleOnly     bool   `json:"visible_only,omitempty"`
	InteractiveOnly bool   `json:"interactive_only,omitempty"`
	IssueCode       string `json:"issue_code,omitempty"`
	InViewport      bool   `json:"in_viewport,omitempty"`
	Limit           int    `json:"limit,omitempty"`
	Cursor          string `json:"cursor,omitempty"`
}

type UIQueryResult struct {
	Nodes      []UINodeSummarySnapshot `json:"nodes,omitempty"`
	NextCursor string                  `json:"nextCursor,omitempty"`
	Total      int                     `json:"total,omitempty"`
}

type UINodeInspectRequest struct {
	NodeID          string `json:"node_id,omitempty"`
	IncludeChildren bool   `json:"include_children,omitempty"`
	ChildDepth      int    `json:"child_depth,omitempty"`
	IncludeProps    bool   `json:"include_props,omitempty"`
	IncludeIssues   bool   `json:"include_issues,omitempty"`
}

type UINodeDetailSnapshot struct {
	Summary  UINodeSummarySnapshot   `json:"summary"`
	Semantic *UISemanticSnapshot     `json:"semantic,omitempty"`
	Layout   *UILayoutSnapshot       `json:"layout,omitempty"`
	Computed *UIComputedSnapshot     `json:"computed,omitempty"`
	Issues   []UIIssueSnapshot       `json:"issues,omitempty"`
	Props    map[string]any          `json:"props,omitempty"`
	Children []UINodeSummarySnapshot `json:"children,omitempty"`
}

type UIIssueListRequest struct {
	Severity string `json:"severity,omitempty"`
	Code     string `json:"code,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Cursor   string `json:"cursor,omitempty"`
}

type UIIssueListSnapshot struct {
	IssueSummary UIIssueSummarySnapshot `json:"issueSummary,omitempty"`
	Issues       []UIIssueSnapshot      `json:"issues,omitempty"`
	NextCursor   string                 `json:"nextCursor,omitempty"`
	Total        int                    `json:"total,omitempty"`
}

type UICaptureRequest struct {
	Target      string `json:"target,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	WithOverlay bool   `json:"with_overlay,omitempty"`
	Scale       int    `json:"scale,omitempty"`
	Padding     int    `json:"padding,omitempty"`
	Rect        *Rect  `json:"rect,omitempty"`
}

type UICaptureResult struct {
	ArtifactID     string `json:"artifactId,omitempty"`
	Path           string `json:"path,omitempty"`
	Target         string `json:"target,omitempty"`
	CapturedRect   Rect   `json:"capturedRect"`
	Width          int    `json:"width,omitempty"`
	Height         int    `json:"height,omitempty"`
	Hash           string `json:"hash,omitempty"`
	OverlayEnabled bool   `json:"overlayEnabled,omitempty"`
	ContentType    string `json:"-"`
}

type UIArtifact struct {
	ID          string
	Path        string
	ContentType string
}
