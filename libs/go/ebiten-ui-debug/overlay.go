package ebitenuidebug

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"golang.org/x/image/font/basicfont"
)

type debugConstraintPatch struct {
	Field  string `json:"field"`
	Value  any    `json:"value"`
	Reason string `json:"reason"`
}

type debugLayoutIssue struct {
	NodeID                    string                 `json:"nodeId"`
	Severity                  string                 `json:"severity"`
	Code                      string                 `json:"code"`
	Message                   string                 `json:"message"`
	SuggestedConstraintChange []debugConstraintPatch `json:"suggestedConstraintChanges,omitempty"`
}

type debugLayoutReport struct {
	Viewport         ebitenui.Viewport
	Issues           []debugLayoutIssue
	IssueSummary     map[string]int
	SummarySnapshot  ebitendebug.UIIssueSummarySnapshot
	InvalidNodeCount int
	issuesByNode     map[string][]debugLayoutIssue
	invalidNodes     map[string]bool
	visibleByNode    map[string]bool
}

func buildDebugLayoutReport(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport) debugLayoutReport {
	report := debugLayoutReport{
		Viewport:      viewport,
		IssueSummary:  map[string]int{},
		issuesByNode:  map[string][]debugLayoutIssue{},
		invalidNodes:  map[string]bool{},
		visibleByNode: map[string]bool{},
	}
	if layout == nil {
		return report
	}

	validation := ebitenui.ValidateLayout(layout, viewport, ebitenui.ValidationOptions{})
	for _, issue := range validation.Issues {
		converted := debugLayoutIssue{
			NodeID:   issue.NodeID,
			Severity: string(issue.Severity),
			Code:     string(issue.Code),
			Message:  issue.Message,
		}
		for _, patch := range issue.SuggestedConstraintChanges {
			converted.SuggestedConstraintChange = append(converted.SuggestedConstraintChange, debugConstraintPatch{
				Field:  patch.Field,
				Value:  patch.Value,
				Reason: patch.Note,
			})
		}
		report.Issues = append(report.Issues, converted)
		report.IssueSummary[converted.Code]++
		report.issuesByNode[converted.NodeID] = append(report.issuesByNode[converted.NodeID], converted)
		if !report.invalidNodes[converted.NodeID] {
			report.invalidNodes[converted.NodeID] = true
			report.InvalidNodeCount++
		}

		switch issue.Severity {
		case ebitenui.IssueSeverityError:
			report.SummarySnapshot.Errors++
		case ebitenui.IssueSeverityWarning:
			report.SummarySnapshot.Warnings++
		default:
			report.SummarySnapshot.Info++
		}
	}
	report.SummarySnapshot.Total = len(report.Issues)
	report.SummarySnapshot.InvalidNodes = report.InvalidNodeCount
	populateVisibility(layout, ebitenui.Rect{X: 0, Y: 0, Width: viewport.Width, Height: viewport.Height}, report.visibleByNode)
	return report
}

func hasIssueCode(issues []debugLayoutIssue, nodeID, code string) bool {
	for _, issue := range issues {
		if issue.NodeID == nodeID && issue.Code == code {
			return true
		}
	}
	return false
}

func buildDebugUISnapshot(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, overlayEnabled bool, runtime *ebitenui.Runtime, input ebitenui.InputSnapshot, queueDepth int) ebitendebug.UISnapshot {
	if layout == nil {
		return ebitendebug.UISnapshot{}
	}

	root := convertLayoutToUISnapshot(layout, nil, report, viewport, overlayEnabled)

	return ebitendebug.UISnapshot{
		Width:            viewport.Width,
		Height:           viewport.Height,
		Viewport:         ebitendebug.UIViewportSnapshot{Width: viewport.Width, Height: viewport.Height, Scale: 1},
		SafeArea:         ebitendebug.UIInsetsSnapshot{},
		IssueSummary:     report.SummarySnapshot,
		InvalidNodeCount: report.InvalidNodeCount,
		InputState:       buildUIInputState(runtime, input, queueDepth),
		Root:             root,
	}
}

func convertLayoutToUISnapshot(layout *ebitenui.LayoutNode, parent *ebitenui.LayoutNode, report debugLayoutReport, viewport ebitenui.Viewport, overlayEnabled bool) ebitendebug.UINodeSnapshot {
	if layout == nil || layout.Node == nil {
		return ebitendebug.UINodeSnapshot{}
	}

	children := make([]ebitendebug.UINodeSnapshot, 0, len(layout.Children))
	for _, child := range layout.Children {
		children = append(children, convertLayoutToUISnapshot(child, layout, report, viewport, overlayEnabled))
	}

	textValue := layout.Node.Text
	if layout.Node.Tag == ebitenui.TagTextBlock && len(layout.TextLines) > 0 {
		textValue = strings.Join(layout.TextLines, "\n")
	}

	props := map[string]any{
		"tag":         string(layout.Node.Tag),
		"overlay":     overlayEnabled,
		"state":       interactionStateMetadata(layout),
		"clickable":   isInteractiveLayoutNode(layout),
		"scrollable":  isScrollLayoutNode(layout),
		"debugTarget": layout.Node.Props.ID,
	}

	bounds := rectToDebug(layout.Frame)
	return ebitendebug.UINodeSnapshot{
		ID:       layout.Node.Props.ID,
		Type:     string(layout.Node.Tag),
		Text:     textValue,
		Visible:  report.visibleByNode[layout.Node.Props.ID],
		Enabled:  !layout.Node.Props.State.Disabled,
		ParentID: layout.ParentID,
		Semantic: semanticMetadata(layout),
		Layout:   layoutMetadata(layout),
		Computed: computedMetadata(layout, parent, report.visibleByNode[layout.Node.Props.ID]),
		Bounds:   bounds,
		Issues:   issuesForNode(report, layout.Node.Props.ID),
		Props:    props,
		Children: children,
	}
}

func issuesForNode(report debugLayoutReport, nodeID string) []ebitendebug.UIIssueSnapshot {
	issues := report.issuesByNode[nodeID]
	if len(issues) == 0 {
		return nil
	}
	result := make([]ebitendebug.UIIssueSnapshot, 0, len(issues))
	for _, issue := range issues {
		patches := make([]ebitendebug.UIConstraintSnapshot, 0, len(issue.SuggestedConstraintChange))
		for _, patch := range issue.SuggestedConstraintChange {
			patches = append(patches, ebitendebug.UIConstraintSnapshot{
				Field: patch.Field,
				Op:    "set",
				Value: patch.Value,
			})
		}
		result = append(result, ebitendebug.UIIssueSnapshot{
			NodeID:                     issue.NodeID,
			Severity:                   issue.Severity,
			Code:                       issue.Code,
			Message:                    issue.Message,
			SuggestedConstraintChanges: patches,
		})
	}
	return result
}

func semanticMetadata(layout *ebitenui.LayoutNode) *ebitendebug.UISemanticSnapshot {
	if layout == nil || layout.Node == nil {
		return nil
	}
	spec := layout.Node.Props.Semantic
	if spec.Screen == "" {
		spec.Screen = "ebiten-ui-showcase"
	}
	if spec.Element == "" {
		spec.Element = layout.Node.Props.ID
	}
	if spec.Role == "" {
		spec.Role = nodeRole(layout)
	}
	if spec.Slot == "" {
		spec.Slot = nodeSlot(layout)
	}
	return &ebitendebug.UISemanticSnapshot{
		Screen:  spec.Screen,
		Element: spec.Element,
		Role:    spec.Role,
		Slot:    spec.Slot,
	}
}

func interactionStateMetadata(layout *ebitenui.LayoutNode) map[string]any {
	if layout == nil || layout.Node == nil {
		return nil
	}
	return map[string]any{
		"hovered":  layout.Node.Props.State.Hovered,
		"focused":  layout.Node.Props.State.Focused,
		"pressed":  layout.Node.Props.State.Pressed,
		"selected": layout.Node.Props.State.Selected,
		"disabled": layout.Node.Props.State.Disabled,
	}
}

func nodeRole(layout *ebitenui.LayoutNode) string {
	if layout == nil || layout.Node == nil {
		return ""
	}
	switch layout.Node.Tag {
	case ebitenui.TagHeader:
		return "header"
	case ebitenui.TagMain:
		return "main"
	case ebitenui.TagSection:
		return "section"
	case ebitenui.TagFooter:
		return "footer"
	case ebitenui.TagButton:
		return "button"
	case ebitenui.TagScrollView:
		return "scroll"
	case ebitenui.TagImage:
		return "image"
	case ebitenui.TagText, ebitenui.TagTextBlock:
		return "text"
	case ebitenui.TagSpacer:
		return "spacer"
	case ebitenui.TagStack:
		return "stack"
	default:
		return string(layout.Node.Tag)
	}
}

func nodeSlot(layout *ebitenui.LayoutNode) string {
	if layout == nil || layout.Node == nil {
		return ""
	}
	id := layout.Node.Props.ID
	switch {
	case id == "showcase-header":
		return "hero"
	case id == "showcase-scroll":
		return "page"
	case strings.Contains(id, "input") || strings.Contains(id, "textarea"):
		return "input"
	case strings.Contains(id, "button"):
		return "action"
	case strings.Contains(id, "dialog"):
		return "modal"
	case strings.Contains(id, "overlay"):
		return "overlay"
	default:
		return ""
	}
}

func layoutMetadata(layout *ebitenui.LayoutNode) *ebitendebug.UILayoutSnapshot {
	if layout == nil || layout.Node == nil {
		return nil
	}

	spec := layout.Node.Props.Layout
	mode := ebitendebug.LayoutMode(derivedLayoutMode(layout.Node))
	result := &ebitendebug.UILayoutSnapshot{
		Mode:     mode,
		ParentID: layout.ParentID,
		Anchor:   ebitendebug.UIAnchor(spec.Anchor),
		Pivot:    ebitendebug.UIPivot(spec.Pivot),
		Offset:   ebitendebug.UIPositionSnapshot{X: spec.Offset.X, Y: spec.Offset.Y},
		Size: ebitendebug.UISizeSnapshot{
			Width:  resolvedLayoutLength(spec.Size.Width, layout.Frame.Width),
			Height: resolvedLayoutLength(spec.Size.Height, layout.Frame.Height),
		},
		ZIndex: spec.ZIndex,
	}
	if result.Anchor == "" {
		result.Anchor = ebitendebug.UIAnchor(ebitenui.AnchorTopLeft)
	}
	if result.Pivot == "" {
		result.Pivot = ebitendebug.UIPivot(ebitenui.PivotTopLeft)
	}
	if spec.MinSize != (ebitenui.LayoutSize{}) {
		result.MinSize = &ebitendebug.UISizeSnapshot{
			Width:  resolvedLayoutLength(spec.MinSize.Width, 0),
			Height: resolvedLayoutLength(spec.MinSize.Height, 0),
		}
	}
	if spec.MaxSize != (ebitenui.LayoutSize{}) {
		result.MaxSize = &ebitendebug.UISizeSnapshot{
			Width:  resolvedLayoutLength(spec.MaxSize.Width, layout.Frame.Width),
			Height: resolvedLayoutLength(spec.MaxSize.Height, layout.Frame.Height),
		}
	}
	margin := spec.Margin
	if margin != (ebitenui.Insets{}) {
		result.Margin = &ebitendebug.UIInsetsSnapshot{Top: margin.Top, Right: margin.Right, Bottom: margin.Bottom, Left: margin.Left}
	}
	padding := spec.Padding
	if padding == (ebitenui.Insets{}) {
		padding = layout.Node.Props.Style.Padding
	}
	if padding != (ebitenui.Insets{}) {
		result.Padding = &ebitendebug.UIInsetsSnapshot{Top: padding.Top, Right: padding.Right, Bottom: padding.Bottom, Left: padding.Left}
	}
	constraints := make([]ebitendebug.UIConstraintSnapshot, 0, 4)
	if spec.Constraints.KeepInsideParent {
		constraints = append(constraints, ebitendebug.UIConstraintSnapshot{Field: "keepInsideParent", Op: "set", Value: true})
	}
	if spec.Constraints.AllowOverlap {
		constraints = append(constraints, ebitendebug.UIConstraintSnapshot{Field: "allowOverlap", Op: "set", Value: true})
	}
	if spec.Constraints.ClipChildren {
		constraints = append(constraints, ebitendebug.UIConstraintSnapshot{Field: "clipChildren", Op: "set", Value: true})
	}
	if spec.Constraints.MinHitTarget > 0 {
		constraints = append(constraints, ebitendebug.UIConstraintSnapshot{Field: "minHitTarget", Op: "set", Value: spec.Constraints.MinHitTarget})
	}
	if len(constraints) > 0 {
		result.Constraints = constraints
	}
	if spec.Grid != (ebitenui.LayoutGrid{}) {
		result.Grid = &ebitendebug.UIGridSnapshot{
			Columns: spec.Grid.Columns,
			Gap:     ebitendebug.UIPositionSnapshot{X: spec.Grid.Gap, Y: spec.Grid.Gap},
		}
	}
	return result
}

func computedMetadata(layout *ebitenui.LayoutNode, parent *ebitenui.LayoutNode, visible bool) *ebitendebug.UIComputedSnapshot {
	if layout == nil {
		return nil
	}
	computed := &ebitendebug.UIComputedSnapshot{
		Bounds:  rectToDebug(layout.Frame),
		Visible: visible,
	}
	if parent != nil {
		parentBounds := rectToDebug(parent.Frame)
		computed.ParentBounds = &parentBounds
	}
	if layout.ContentBounds != (ebitenui.Rect{}) {
		contentBounds := rectToDebug(layout.ContentBounds)
		computed.ContentBounds = &contentBounds
	}
	if layout.ClickableRect != (ebitenui.Rect{}) {
		clickableRect := rectToDebug(layout.ClickableRect)
		computed.ClickableRect = &clickableRect
	}
	if layout.ClipRect != (ebitenui.Rect{}) {
		clipRect := rectToDebug(layout.ClipRect)
		computed.ClipRect = &clipRect
	}
	if layout.Overflow.Any {
		computed.Overflow = &ebitendebug.UIOverflowSnapshot{
			Top:    layout.Overflow.Vertical,
			Bottom: layout.Overflow.Vertical,
			Left:   layout.Overflow.Horizontal,
			Right:  layout.Overflow.Horizontal,
		}
	}
	return computed
}

func buildUIInputState(runtime *ebitenui.Runtime, input ebitenui.InputSnapshot, queueDepth int) ebitendebug.UIInputSnapshot {
	state := ebitendebug.UIInputSnapshot{
		Pointer: &ebitendebug.UIPointerSnapshot{
			X:    input.PointerX,
			Y:    input.PointerY,
			Down: input.PointerDown,
		},
	}
	if runtime != nil {
		state.FocusedNodeID = runtime.FocusedID()
		state.HoveredNodeID = runtime.HoveredID()
	}
	if input.ScrollX != 0 || input.ScrollY != 0 {
		state.Scroll = &ebitendebug.UIScrollSnapshot{
			X:      input.ScrollX,
			Y:      input.ScrollY,
			Source: "runtime",
		}
	}
	keysPressed := make([]string, 0, 8)
	if input.SelectAll {
		keysPressed = append(keysPressed, "SelectAll")
	}
	if input.Backspace {
		keysPressed = append(keysPressed, "Backspace")
	}
	if input.Delete {
		keysPressed = append(keysPressed, "Delete")
	}
	if input.Home {
		keysPressed = append(keysPressed, "Home")
	}
	if input.End {
		keysPressed = append(keysPressed, "End")
	}
	if input.Submit {
		keysPressed = append(keysPressed, "Enter")
	}
	if input.Space {
		keysPressed = append(keysPressed, "Space")
	}
	if input.Tab {
		keysPressed = append(keysPressed, "Tab")
	}
	if input.Escape {
		keysPressed = append(keysPressed, "Escape")
	}
	if input.ArrowUp {
		keysPressed = append(keysPressed, "ArrowUp")
	}
	if input.ArrowDown {
		keysPressed = append(keysPressed, "ArrowDown")
	}
	if input.ArrowLeft {
		keysPressed = append(keysPressed, "ArrowLeft")
	}
	if input.ArrowRight {
		keysPressed = append(keysPressed, "ArrowRight")
	}
	modifiers := make([]string, 0, 4)
	if input.Shift {
		modifiers = append(modifiers, "Shift")
	}
	if input.Control {
		modifiers = append(modifiers, "Control")
	}
	if input.Alt {
		modifiers = append(modifiers, "Alt")
	}
	if input.Meta {
		modifiers = append(modifiers, "Meta")
	}
	if input.Text != "" || len(keysPressed) > 0 || len(modifiers) > 0 || queueDepth > 0 {
		state.Keyboard = &ebitendebug.UIKeyboardSnapshot{
			Text:        input.Text,
			KeysPressed: keysPressed,
			Modifiers:   modifiers,
			QueueDepth:  queueDepth,
		}
	}
	return state
}

func rectToDebug(rect ebitenui.Rect) ebitendebug.Rect {
	return ebitendebug.Rect{
		X:      rect.X,
		Y:      rect.Y,
		Width:  rect.Width,
		Height: rect.Height,
	}
}

func derivedLayoutMode(node *ebitenui.Node) ebitenui.LayoutMode {
	if node == nil {
		return ebitenui.LayoutModeFlowVertical
	}
	if node.Props.Layout.Mode != "" {
		return node.Props.Layout.Mode
	}
	switch {
	case node.Tag == ebitenui.TagStack:
		return ebitenui.LayoutModeStack
	case node.Tag == ebitenui.TagScrollView:
		return ebitenui.LayoutModeFlowVertical
	case node.Props.Style.Direction == ebitenui.Row:
		return ebitenui.LayoutModeFlowHorizontal
	default:
		return ebitenui.LayoutModeFlowVertical
	}
}

func resolvedLayoutLength(length ebitenui.Length, fallback float64) float64 {
	switch length.Kind {
	case ebitenui.LengthPx:
		return length.Value
	case ebitenui.LengthFill:
		return fallback
	default:
		return fallback
	}
}

func drawDebugOverlay(screen *ebiten.Image, layout *ebitenui.LayoutNode, report debugLayoutReport, overlayEnabled bool) {
	if !overlayEnabled || screen == nil || layout == nil {
		return
	}

	drawDebugOverlayNode(screen, layout, report)
}

func drawDebugOverlayNode(screen *ebiten.Image, layout *ebitenui.LayoutNode, report debugLayoutReport) {
	if layout == nil || layout.Node == nil {
		return
	}

	frame := layout.Frame
	issueCount := len(report.issuesByNode[layout.Node.Props.ID])

	strokeColor := color.RGBA{R: 108, G: 140, B: 179, A: 160}
	if issueCount > 0 {
		strokeColor = color.RGBA{R: 255, G: 84, B: 84, A: 220}
	}
	if layout.Node.Props.State.Focused {
		strokeColor = color.RGBA{R: 72, G: 214, B: 160, A: 220}
	}
	if layout.Node.Props.State.Hovered {
		strokeColor = color.RGBA{R: 246, G: 196, B: 62, A: 220}
	}
	if isScrollLayoutNode(layout) {
		strokeColor = color.RGBA{R: 100, G: 162, B: 255, A: 180}
	}

	vector.StrokeRect(screen, float32(frame.X), float32(frame.Y), float32(frame.Width), float32(frame.Height), 1.5, strokeColor, false)

	anchorColor := color.RGBA{R: 120, G: 192, B: 255, A: 220}
	pivotColor := color.RGBA{R: 255, G: 120, B: 192, A: 220}
	drawMarker(screen, frame.X, frame.Y, anchorColor)
	drawMarker(screen, frame.X+frame.Width*0.5, frame.Y+frame.Height*0.5, pivotColor)

	if issueCount > 0 || layout.Node.Props.State.Focused || layout.Node.Props.State.Hovered || isInteractiveLayoutNode(layout) || isScrollLayoutNode(layout) {
		label := layout.Node.Props.ID
		if label == "" {
			label = string(layout.Node.Tag)
		}
		if issueCount > 0 {
			label = fmt.Sprintf("%s [%s]", label, report.issuesByNode[layout.Node.Props.ID][0].Code)
		}
		drawDebugLabel(screen, frame.X+2, frame.Y+2, label, strokeColor)
	}

	for _, child := range layout.Children {
		drawDebugOverlayNode(screen, child, report)
	}
}

func drawMarker(screen *ebiten.Image, x, y float64, fill color.Color) {
	vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, fill, false)
}

func drawDebugLabel(screen *ebiten.Image, x, y float64, label string, fill color.Color) {
	if label == "" {
		return
	}
	text.Draw(screen, label, basicfont.Face7x13, int(x), int(y)+basicfont.Face7x13.Metrics().Ascent.Ceil(), fill)
}

func outOfViewport(frame ebitenui.Rect, viewport ebitenui.Viewport) bool {
	if viewport.Width <= 0 || viewport.Height <= 0 {
		return false
	}
	return frame.X < 0 || frame.Y < 0 || frame.X+frame.Width > viewport.Width || frame.Y+frame.Height > viewport.Height
}

func outOfParent(frame, parent ebitenui.Rect) bool {
	return frame.X < parent.X || frame.Y < parent.Y || frame.X+frame.Width > parent.X+parent.Width || frame.Y+frame.Height > parent.Y+parent.Height
}

func outsideSafeArea(frame ebitenui.Rect, viewport ebitenui.Viewport, inset float64) bool {
	if viewport.Width <= 0 || viewport.Height <= 0 {
		return false
	}
	left := inset
	top := inset
	right := viewport.Width - inset
	bottom := viewport.Height - inset
	return frame.X < left || frame.Y < top || frame.X+frame.Width > right || frame.Y+frame.Height > bottom
}

func textOverflows(layout *ebitenui.LayoutNode) bool {
	if layout == nil || layout.Node == nil {
		return false
	}

	lines := layout.TextLines
	if len(lines) == 0 && layout.Node.Text != "" {
		lines = []string{layout.Node.Text}
	}
	if len(lines) == 0 {
		return false
	}

	maxWidth := 0
	for _, line := range lines {
		width := text.BoundString(basicfont.Face7x13, line).Dx()
		if width > maxWidth {
			maxWidth = width
		}
	}
	return float64(maxWidth) > layout.Frame.Width
}

func intersects(a, b ebitenui.Rect) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X && a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}

func intersectionRect(a, b ebitenui.Rect) ebitenui.Rect {
	minX := maxFloat(a.X, b.X)
	minY := maxFloat(a.Y, b.Y)
	maxX := minFloat(a.X+a.Width, b.X+b.Width)
	maxY := minFloat(a.Y+a.Height, b.Y+b.Height)
	if maxX <= minX || maxY <= minY {
		return ebitenui.Rect{}
	}
	return ebitenui.Rect{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

func populateVisibility(layout *ebitenui.LayoutNode, inheritedClip ebitenui.Rect, visibleByNode map[string]bool) {
	if layout == nil || layout.Node == nil {
		return
	}

	nodeVisible := intersects(layout.Frame, inheritedClip)
	visibleByNode[layout.Node.Props.ID] = nodeVisible

	childClip := inheritedClip
	if layout.ClipChildren {
		clip := layout.ClipRect
		if clip == (ebitenui.Rect{}) {
			clip = layout.Frame
		}
		childClip = intersectionRect(childClip, clip)
	}

	for _, child := range layout.Children {
		populateVisibility(child, childClip, visibleByNode)
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
