package ebitenuidebug

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	stdimagedraw "image/draw"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

type compactNodeRef struct {
	layout *ebitenui.LayoutNode
	parent *ebitenui.LayoutNode
}

func buildCompactUIOverview(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, runtime *ebitenui.Runtime, input ebitenui.InputSnapshot, queueDepth int) ebitendebug.UIOverviewSnapshot {
	if layout == nil || layout.Node == nil {
		return ebitendebug.UIOverviewSnapshot{}
	}

	refs := flattenCompactLayout(layout)
	topLevel := make([]ebitendebug.UINodeSummarySnapshot, 0, len(layout.Children))
	visibleCount := 0
	for _, ref := range refs {
		if ref.layout != nil && inViewport(ref.layout.Frame, viewport) {
			visibleCount++
		}
	}
	for _, child := range layout.Children {
		topLevel = append(topLevel, compactSummaryForNode(child, layout, report, viewport))
	}

	overview := ebitendebug.UIOverviewSnapshot{
		Viewport:         ebitendebug.UIViewportSnapshot{Width: viewport.Width, Height: viewport.Height, Scale: 1},
		SafeArea:         ebitendebug.UIInsetsSnapshot{},
		RootID:           layout.Node.Props.ID,
		TotalNodeCount:   len(refs),
		VisibleNodeCount: visibleCount,
		InvalidNodeCount: report.InvalidNodeCount,
		TopLevelSections: topLevel,
		IssueSummary:     report.SummarySnapshot,
	}
	if runtime != nil {
		overview.FocusedNodeID = runtime.FocusedID()
		overview.HoveredNodeID = runtime.HoveredID()
	}
	if overview.HoveredNodeID == "" && (input.PointerX != 0 || input.PointerY != 0) {
		if target, ok := resolvePointerTargetAt(layout, input.PointerX, input.PointerY); ok {
			overview.HoveredNodeID = target.ID
		}
	}
	_ = queueDepth
	return overview
}

func queryCompactUINodes(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, request ebitendebug.UIQueryRequest) ebitendebug.UIQueryResult {
	refs := flattenCompactLayout(layout)
	filtered := make([]compactNodeRef, 0, len(refs))
	for _, ref := range refs {
		if matchesCompactQuery(ref, viewport, report, request) {
			filtered = append(filtered, ref)
		}
	}

	start := decodeCursorOffset(request.Cursor)
	if start > len(filtered) {
		start = len(filtered)
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 25
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := ebitendebug.UIQueryResult{
		Nodes: make([]ebitendebug.UINodeSummarySnapshot, 0, end-start),
		Total: len(filtered),
	}
	for _, ref := range filtered[start:end] {
		result.Nodes = append(result.Nodes, compactSummaryForNode(ref.layout, ref.parent, report, viewport))
	}
	if end < len(filtered) {
		result.NextCursor = strconv.Itoa(end)
	}
	return result
}

func inspectCompactUINode(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, request ebitendebug.UINodeInspectRequest) (ebitendebug.UINodeDetailSnapshot, bool) {
	target, ok := resolveDebugTarget(layout, request.NodeID)
	if !ok {
		return ebitendebug.UINodeDetailSnapshot{}, false
	}
	parent := parentLayoutForPath(target.Path)
	detail := ebitendebug.UINodeDetailSnapshot{
		Summary:  compactSummaryForNode(target.Node, parent, report, viewport),
		Semantic: semanticMetadata(target.Node),
		Layout:   layoutMetadata(target.Node),
		Computed: computedMetadata(target.Node, parent),
	}
	if request.IncludeIssues {
		detail.Issues = issuesForNode(report, target.Node.Node.Props.ID)
	}
	if request.IncludeProps {
		detail.Props = compactPropsForNode(target.Node)
	}
	if request.IncludeChildren {
		detail.Children = collectDirectChildSummaries(target.Node, report, viewport)
	}
	return detail, true
}

func listCompactUIIssues(report debugLayoutReport, request ebitendebug.UIIssueListRequest) ebitendebug.UIIssueListSnapshot {
	issues := append([]debugLayoutIssue(nil), report.Issues...)
	sort.SliceStable(issues, func(i, j int) bool {
		left := issueSeverityRank(issues[i].Severity)
		right := issueSeverityRank(issues[j].Severity)
		if left != right {
			return left > right
		}
		if issues[i].NodeID != issues[j].NodeID {
			return issues[i].NodeID < issues[j].NodeID
		}
		return issues[i].Code < issues[j].Code
	})

	filtered := make([]debugLayoutIssue, 0, len(issues))
	for _, issue := range issues {
		if request.Severity != "" && !strings.EqualFold(issue.Severity, request.Severity) {
			continue
		}
		if request.Code != "" && issue.Code != request.Code {
			continue
		}
		if request.NodeID != "" && issue.NodeID != request.NodeID {
			continue
		}
		filtered = append(filtered, issue)
	}

	start := decodeCursorOffset(request.Cursor)
	if start > len(filtered) {
		start = len(filtered)
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := ebitendebug.UIIssueListSnapshot{
		IssueSummary: report.SummarySnapshot,
		Issues:       make([]ebitendebug.UIIssueSnapshot, 0, end-start),
		Total:        len(filtered),
	}
	for _, issue := range filtered[start:end] {
		result.Issues = append(result.Issues, ebitendebug.UIIssueSnapshot{
			NodeID:   issue.NodeID,
			Severity: issue.Severity,
			Code:     issue.Code,
			Message:  issue.Message,
		})
	}
	if end < len(filtered) {
		result.NextCursor = strconv.Itoa(end)
	}
	return result
}

func captureCompactUIScreenshot(gameID string, screenshotsDir string, layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, request ebitendebug.UICaptureRequest) (ebitendebug.UICaptureResult, ebitendebug.UIArtifact, bool) {
	if layout == nil {
		return ebitendebug.UICaptureResult{}, ebitendebug.UIArtifact{}, false
	}
	target := request.Target
	if target == "" {
		target = "viewport"
	}
	if request.Scale <= 0 {
		request.Scale = 1
	}

	captureRect, ok := resolveCaptureRect(layout, viewport, request)
	if !ok {
		return ebitendebug.UICaptureResult{}, ebitendebug.UIArtifact{}, false
	}

	full := renderLayoutToRGBA(layout, viewport, report, request.WithOverlay)
	cropped := cropRGBA(full, captureRect)
	scaled := scaleRGBA(cropped, request.Scale)
	encoded := encodePNGBytes(scaled)
	if len(encoded) == 0 {
		return ebitendebug.UICaptureResult{}, ebitendebug.UIArtifact{}, false
	}

	hashBytes := sha256.Sum256(encoded)
	hashValue := hex.EncodeToString(hashBytes[:])
	artifactID := "capture-" + hashValue[:12]
	if strings.TrimSpace(screenshotsDir) == "" {
		screenshotsDir = filepath.Join(repoRootDir(), "screenshots")
	}
	dir := filepath.Join(screenshotsDir, gameID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ebitendebug.UICaptureResult{}, ebitendebug.UIArtifact{}, false
	}
	path := filepath.Join(dir, artifactID+".png")
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		return ebitendebug.UICaptureResult{}, ebitendebug.UIArtifact{}, false
	}

	result := ebitendebug.UICaptureResult{
		ArtifactID:     artifactID,
		Path:           path,
		Target:         target,
		CapturedRect:   ebitendebug.Rect{X: float64(captureRect.Min.X), Y: float64(captureRect.Min.Y), Width: float64(captureRect.Dx()), Height: float64(captureRect.Dy())},
		Width:          scaled.Bounds().Dx(),
		Height:         scaled.Bounds().Dy(),
		Hash:           hashValue,
		OverlayEnabled: request.WithOverlay,
	}
	return result, ebitendebug.UIArtifact{ID: artifactID, Path: path, ContentType: "image/png"}, true
}

func flattenCompactLayout(root *ebitenui.LayoutNode) []compactNodeRef {
	if root == nil {
		return nil
	}
	result := []compactNodeRef{}
	var walk func(node *ebitenui.LayoutNode, parent *ebitenui.LayoutNode)
	walk = func(node *ebitenui.LayoutNode, parent *ebitenui.LayoutNode) {
		if node == nil || node.Node == nil {
			return
		}
		result = append(result, compactNodeRef{layout: node, parent: parent})
		for _, child := range node.Children {
			walk(child, node)
		}
	}
	walk(root, nil)
	return result
}

func matchesCompactQuery(ref compactNodeRef, viewport ebitenui.Viewport, report debugLayoutReport, request ebitendebug.UIQueryRequest) bool {
	summary := compactSummaryForNode(ref.layout, ref.parent, report, viewport)
	if request.ID != "" && summary.ID != request.ID {
		return false
	}
	if request.Role != "" && !strings.EqualFold(summary.Role, request.Role) {
		return false
	}
	if request.Slot != "" && !strings.EqualFold(summary.Slot, request.Slot) {
		return false
	}
	if request.Type != "" && !strings.EqualFold(summary.Type, request.Type) {
		return false
	}
	if request.TextContains != "" && !strings.Contains(strings.ToLower(summary.TextPreview), strings.ToLower(request.TextContains)) {
		return false
	}
	if request.VisibleOnly && !summary.Visible {
		return false
	}
	if request.InteractiveOnly && !summary.Interactive {
		return false
	}
	if request.IssueCode != "" && !hasIssueCode(report.Issues, summary.ID, request.IssueCode) {
		return false
	}
	if request.InViewport && !inViewport(ref.layout.Frame, viewport) {
		return false
	}
	return true
}

func compactSummaryForNode(layout *ebitenui.LayoutNode, parent *ebitenui.LayoutNode, report debugLayoutReport, viewport ebitenui.Viewport) ebitendebug.UINodeSummarySnapshot {
	if layout == nil || layout.Node == nil {
		return ebitendebug.UINodeSummarySnapshot{}
	}
	textPreview := strings.TrimSpace(layout.Node.Text)
	if layout.Node.Tag == ebitenui.TagTextBlock && len(layout.TextLines) > 0 {
		textPreview = strings.Join(layout.TextLines, " ")
	}
	textPreview = compactTextPreview(textPreview, 64)
	_ = parent

	return ebitendebug.UINodeSummarySnapshot{
		ID:          layout.Node.Props.ID,
		Type:        string(layout.Node.Tag),
		Role:        nodeRole(layout),
		Slot:        nodeSlot(layout),
		Bounds:      rectToDebug(layout.Frame),
		Visible:     inViewport(layout.Frame, viewport),
		Enabled:     !layout.Node.Props.State.Disabled,
		ChildCount:  len(layout.Children),
		IssueCount:  len(report.issuesByNode[layout.Node.Props.ID]),
		Interactive: isInteractiveLayoutNode(layout),
		Scrollable:  isScrollLayoutNode(layout),
		TextPreview: textPreview,
	}
}

func compactPropsForNode(layout *ebitenui.LayoutNode) map[string]any {
	if layout == nil || layout.Node == nil {
		return nil
	}
	return map[string]any{
		"tag":        string(layout.Node.Tag),
		"state":      interactionStateMetadata(layout),
		"clickable":  isInteractiveLayoutNode(layout),
		"scrollable": isScrollLayoutNode(layout),
	}
}

func collectDirectChildSummaries(layout *ebitenui.LayoutNode, report debugLayoutReport, viewport ebitenui.Viewport) []ebitendebug.UINodeSummarySnapshot {
	if layout == nil {
		return nil
	}
	children := make([]ebitendebug.UINodeSummarySnapshot, 0, len(layout.Children))
	for _, child := range layout.Children {
		children = append(children, compactSummaryForNode(child, layout, report, viewport))
	}
	return children
}

func issueSeverityRank(severity string) int {
	switch strings.ToLower(severity) {
	case "error":
		return 3
	case "warning":
		return 2
	default:
		return 1
	}
}

func decodeCursorOffset(cursor string) int {
	if strings.TrimSpace(cursor) == "" {
		return 0
	}
	value, err := strconv.Atoi(cursor)
	if err != nil || value < 0 {
		return 0
	}
	return value
}

func inViewport(rect ebitenui.Rect, viewport ebitenui.Viewport) bool {
	view := ebitenui.Rect{X: 0, Y: 0, Width: viewport.Width, Height: viewport.Height}
	return intersects(view, rect)
}

func compactTextPreview(value string, limit int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.Join(strings.Fields(value), " ")
	if len(value) <= limit {
		return value
	}
	if limit <= 3 {
		return value[:limit]
	}
	return value[:limit-3] + "..."
}

func renderLayoutToRGBA(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, report debugLayoutReport, withOverlay bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(viewport.Width), int(viewport.Height)))
	stdimagedraw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 16, G: 18, B: 24, A: 255}}, image.Point{}, stdimagedraw.Src)
	drawCompactNode(img, layout)
	if withOverlay {
		drawCompactOverlay(img, layout, report)
	}
	return img
}

func drawCompactNode(dst *image.RGBA, layout *ebitenui.LayoutNode) {
	if layout == nil || layout.Node == nil {
		return
	}
	frame := image.Rect(
		int(layout.Frame.X),
		int(layout.Frame.Y),
		int(layout.Frame.X+layout.Frame.Width),
		int(layout.Frame.Y+layout.Frame.Height),
	)
	style := layout.Node.Props.Style
	if style.BackgroundColor != nil {
		fillRect(dst, frame, style.BackgroundColor)
	}
	if style.BorderColor != nil && style.BorderWidth > 0 {
		strokeRect(dst, frame, style.BorderColor, int(style.BorderWidth))
	}
	drawCompactInteractionState(dst, frame, layout)

	switch layout.Node.Tag {
	case ebitenui.TagText:
		drawCompactText(dst, frame, []string{layout.Node.Text}, style.Color)
	case ebitenui.TagTextBlock:
		drawCompactText(dst, frame, layout.TextLines, style.Color)
	case ebitenui.TagImage:
		fill := layout.Node.Props.Image.Fill
		if fill == nil {
			fill = color.RGBA{R: 72, G: 82, B: 96, A: 255}
		}
		fillRect(dst, frame, fill)
	}

	for _, child := range layout.Children {
		drawCompactNode(dst, child)
	}
}

func drawCompactText(dst *image.RGBA, frame image.Rectangle, lines []string, textColor color.Color) {
	if len(lines) == 0 {
		return
	}
	if textColor == nil {
		textColor = color.White
	}
	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(textColor),
		Face: basicfont.Face7x13,
	}
	ascent := basicfont.Face7x13.Metrics().Ascent.Ceil()
	lineHeight := basicfont.Face7x13.Metrics().Height.Ceil()
	for index, line := range lines {
		drawer.Dot = fixed.P(frame.Min.X, frame.Min.Y+ascent+index*lineHeight)
		drawer.DrawString(line)
	}
}

func drawCompactInteractionState(dst *image.RGBA, frame image.Rectangle, layout *ebitenui.LayoutNode) {
	state := layout.Node.Props.State
	if state.Hovered {
		fillRect(dst, frame, color.RGBA{R: 255, G: 255, B: 255, A: 20})
	}
	if state.Pressed {
		fillRect(dst, frame, color.RGBA{R: 0, G: 0, B: 0, A: 40})
	}
	if state.Focused || state.Selected {
		strokeRect(dst, frame, color.RGBA{R: 255, G: 214, B: 92, A: 255}, 2)
	}
	if state.Disabled {
		fillRect(dst, frame, color.RGBA{R: 20, G: 20, B: 24, A: 96})
	}
}

func drawCompactOverlay(dst *image.RGBA, layout *ebitenui.LayoutNode, report debugLayoutReport) {
	for _, ref := range flattenCompactLayout(layout) {
		if ref.layout == nil || ref.layout.Node == nil {
			continue
		}
		frame := image.Rect(
			int(ref.layout.Frame.X),
			int(ref.layout.Frame.Y),
			int(ref.layout.Frame.X+ref.layout.Frame.Width),
			int(ref.layout.Frame.Y+ref.layout.Frame.Height),
		)
		if len(report.issuesByNode[ref.layout.Node.Props.ID]) > 0 {
			strokeRect(dst, frame, color.RGBA{R: 255, G: 87, B: 87, A: 255}, 2)
			continue
		}
		if ref.layout.Node.Props.State.Focused || ref.layout.Node.Props.State.Hovered {
			strokeRect(dst, frame, color.RGBA{R: 255, G: 214, B: 92, A: 255}, 1)
		}
	}
}

func resolveCaptureRect(layout *ebitenui.LayoutNode, viewport ebitenui.Viewport, request ebitendebug.UICaptureRequest) (image.Rectangle, bool) {
	bounds := image.Rect(0, 0, int(viewport.Width), int(viewport.Height))
	switch request.Target {
	case "", "viewport":
		return applyPadding(bounds, request.Padding, bounds), true
	case "node_id":
		target, ok := resolveDebugTarget(layout, request.NodeID)
		if !ok {
			return image.Rectangle{}, false
		}
		rect := target.Frame
		if target.Node != nil {
			if target.Node.ClickableRect != (ebitenui.Rect{}) {
				rect = target.Node.ClickableRect
			} else {
				rect = target.Node.Frame
			}
		}
		return applyPadding(image.Rect(int(rect.X), int(rect.Y), int(rect.X+rect.Width), int(rect.Y+rect.Height)), request.Padding, bounds), true
	case "rect":
		if request.Rect == nil {
			return image.Rectangle{}, false
		}
		rect := image.Rect(int(request.Rect.X), int(request.Rect.Y), int(request.Rect.X+request.Rect.Width), int(request.Rect.Y+request.Rect.Height))
		return applyPadding(rect, request.Padding, bounds), true
	default:
		return image.Rectangle{}, false
	}
}

func cropRGBA(src *image.RGBA, rect image.Rectangle) *image.RGBA {
	rect = rect.Intersect(src.Bounds())
	if rect.Empty() {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	stdimagedraw.Draw(dst, dst.Bounds(), src, rect.Min, stdimagedraw.Src)
	return dst
}

func scaleRGBA(src *image.RGBA, scale int) *image.RGBA {
	if scale <= 1 {
		return src
	}
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx()*scale, src.Bounds().Dy()*scale))
	xdraw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), stdimagedraw.Src, nil)
	return dst
}

func applyPadding(rect image.Rectangle, padding int, bounds image.Rectangle) image.Rectangle {
	if padding > 0 {
		rect = image.Rect(rect.Min.X-padding, rect.Min.Y-padding, rect.Max.X+padding, rect.Max.Y+padding)
	}
	rect = rect.Intersect(bounds)
	if rect.Empty() {
		return image.Rect(0, 0, 1, 1)
	}
	return rect
}

func encodePNGBytes(img image.Image) []byte {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil
	}
	return buffer.Bytes()
}

func fillRect(dst *image.RGBA, rect image.Rectangle, fill color.Color) {
	if rect.Empty() {
		return
	}
	stdimagedraw.Draw(dst, rect, &image.Uniform{C: fill}, image.Point{}, stdimagedraw.Over)
}

func strokeRect(dst *image.RGBA, rect image.Rectangle, stroke color.Color, width int) {
	if rect.Empty() || width <= 0 {
		return
	}
	fillRect(dst, image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+width), stroke)
	fillRect(dst, image.Rect(rect.Min.X, rect.Max.Y-width, rect.Max.X, rect.Max.Y), stroke)
	fillRect(dst, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+width, rect.Max.Y), stroke)
	fillRect(dst, image.Rect(rect.Max.X-width, rect.Min.Y, rect.Max.X, rect.Max.Y), stroke)
}

func resolvePointerTargetAt(layout *ebitenui.LayoutNode, x, y float64) (debugResolvedTarget, bool) {
	refs := flattenCompactLayout(layout)
	for i := len(refs) - 1; i >= 0; i-- {
		ref := refs[i]
		if ref.layout == nil || ref.layout.Node == nil {
			continue
		}
		rect := ref.layout.ClickableRect
		if rect == (ebitenui.Rect{}) {
			rect = ref.layout.Frame
		}
		if !containsPoint(rect, x, y) {
			continue
		}
		path := findLayoutPath(layout, ref.layout.Node.Props.ID)
		if len(path) == 0 {
			continue
		}
		return debugResolvedTarget{
			ID:    ref.layout.Node.Props.ID,
			Frame: rect,
			Node:  ref.layout,
			Path:  path,
		}, true
	}
	return debugResolvedTarget{}, false
}

func containsPoint(rect ebitenui.Rect, x, y float64) bool {
	return x >= rect.X && x <= rect.X+rect.Width && y >= rect.Y && y <= rect.Y+rect.Height
}

func repoRootDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}
