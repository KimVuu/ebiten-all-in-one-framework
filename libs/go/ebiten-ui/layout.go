package ebitenui

import (
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Viewport struct {
	Width  float64
	Height float64
}

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type LayoutNode struct {
	Node          *Node
	ParentID      string
	Frame         Rect
	Children      []*LayoutNode
	TextLines     []string
	ContentBounds Rect
	ClickableRect Rect
	ClipRect      Rect
	ContentWidth  float64
	ContentHeight float64
	Overflow      LayoutOverflow
	ClipChildren  bool
}

func (n *LayoutNode) FindByID(id string) (*LayoutNode, bool) {
	if n == nil {
		return nil, false
	}
	if n.Node != nil && n.Node.Props.ID == id {
		return n, true
	}
	for _, child := range n.Children {
		if found, ok := child.FindByID(id); ok {
			return found, true
		}
	}
	return nil, false
}

func (d *DOM) Layout(viewport Viewport) *LayoutNode {
	if d == nil || d.Root == nil {
		return nil
	}

	width, height := measureNode(d.Root, viewport.Width, viewport.Height)
	if width == 0 {
		width = viewport.Width
	}
	if height == 0 {
		height = viewport.Height
	}

	return buildLayout(d.Root, 0, 0, width, height, "")
}

func buildLayout(node *Node, x, y, width, height float64, parentID string) *LayoutNode {
	layout := &LayoutNode{
		Node:     node,
		ParentID: parentID,
		Frame:    Rect{X: x, Y: y, Width: width, Height: height},
	}
	if node == nil {
		return layout
	}

	switch node.Tag {
	case TagText:
		layout.TextLines = []string{node.Text}
		layout.ClickableRect = layout.Frame
		layout.ClipRect = layout.Frame
		layout.ContentBounds = layout.Frame
		return layout
	case TagTextBlock:
		layout.TextLines = wrapTextLines(node.Text, node.Props.Style, width)
		layout.ClickableRect = layout.Frame
		layout.ClipRect = layout.Frame
		layout.ContentBounds = layout.Frame
		return layout
	case TagImage, TagSpacer:
		layout.ClickableRect = layout.Frame
		layout.ClipRect = layout.Frame
		layout.ContentBounds = layout.Frame
		return layout
	}
	if len(node.Children) == 0 {
		layout.ClickableRect = layout.Frame
		layout.ClipRect = layout.Frame
		layout.ContentBounds = layout.Frame
		return layout
	}

	layoutSpec, hasLayout := effectiveLayoutSpec(node)
	style := node.Props.Style
	padding := style.Padding
	if hasLayout {
		padding = mergeInsets(layoutSpec.Padding, style.Padding)
	}
	gap := style.Gap
	if hasLayout && layoutSpec.Gap != 0 {
		gap = layoutSpec.Gap
	}
	contentX := x + padding.Left
	contentY := y + padding.Top
	contentWidth := maxFloat(0, width-padding.Horizontal())
	contentHeight := maxFloat(0, height-padding.Vertical())
	layout.ContentBounds = Rect{X: contentX, Y: contentY, Width: contentWidth, Height: contentHeight}
	layout.ClickableRect = layout.Frame
	layout.ClipRect = layout.Frame
	if node.Tag == TagScrollView || (hasLayout && layoutSpec.Constraints.ClipChildren) {
		layout.ClipRect = layout.ContentBounds
		layout.ClipChildren = true
	}

	mode := layoutModeFor(node)
	switch mode {
	case LayoutModeStack:
		layout.ContentWidth, layout.ContentHeight = buildStackChildren(layout, contentX, contentY, contentWidth, contentHeight)
	case LayoutModeGrid:
		layout.ContentWidth, layout.ContentHeight = buildGridChildren(layout, contentX, contentY, contentWidth, contentHeight)
	case LayoutModeFlowHorizontal, LayoutModeFlowVertical:
		offsetX := 0.0
		offsetY := 0.0
		if node.Tag == TagScrollView {
			offsetX = node.Props.Scroll.OffsetX
			offsetY = node.Props.Scroll.OffsetY
		}
		layout.ContentWidth, layout.ContentHeight = buildFlowChildren(layout, contentX, contentY, contentWidth, contentHeight, offsetX, offsetY, gap)
	default:
		layout.ContentWidth, layout.ContentHeight = buildFlowChildren(layout, contentX, contentY, contentWidth, contentHeight, 0, 0, gap)
	}

	layout.Overflow.Horizontal = layout.ContentWidth > contentWidth || layout.Frame.Width > width
	layout.Overflow.Vertical = layout.ContentHeight > contentHeight || layout.Frame.Height > height
	layout.Overflow.Any = layout.Overflow.Horizontal || layout.Overflow.Vertical
	if layout.ClickableRect == (Rect{}) {
		layout.ClickableRect = layout.Frame
	}
	if layout.ClipRect == (Rect{}) {
		layout.ClipRect = layout.Frame
	}

	return layout
}

func buildStackChildren(layout *LayoutNode, baseX, baseY, contentWidth, contentHeight float64) (float64, float64) {
	layout.Children = make([]*LayoutNode, 0, len(layout.Node.Children))
	maxWidth := 0.0
	maxHeight := 0.0

	for _, child := range layout.Node.Children {
		childWidth, childHeight := measureNode(child, contentWidth, contentHeight)
		childWidth, childHeight = resolveLayoutSize(child, contentWidth, contentHeight, childWidth, childHeight)

		childX, childY := baseX, baseY
		if layoutModeFor(child) == LayoutModeAnchored {
			childX, childY = resolveAnchoredChildFrame(layout, child, baseX, baseY, contentWidth, contentHeight, childWidth, childHeight)
		}
		childX, childY, childWidth, childHeight = applyLayoutMargin(child, childX, childY, childWidth, childHeight)

		layout.Children = append(layout.Children, buildLayout(child, childX, childY, childWidth, childHeight, layout.Node.Props.ID))
		maxWidth = maxFloat(maxWidth, (childX-baseX)+childWidth)
		maxHeight = maxFloat(maxHeight, (childY-baseY)+childHeight)
	}

	return maxWidth, maxHeight
}

func buildFlowChildren(layout *LayoutNode, baseX, baseY, contentWidth, contentHeight, offsetX, offsetY, gap float64) (float64, float64) {
	direction := layoutDirectionFor(layout.Node)
	gapTotal := gap * float64(maxInt(0, len(layout.Node.Children)-1))
	mainAvailable := contentHeight
	if direction == Row {
		mainAvailable = contentWidth
	}

	fixedMain := 0.0
	fillCount := 0
	measured := make([]Rect, len(layout.Node.Children))

	for i, child := range layout.Node.Children {
		childWidth, childHeight := measureNode(child, contentWidth, contentHeight)
		childWidth, childHeight = resolveLayoutSize(child, contentWidth, contentHeight, childWidth, childHeight)
		measured[i] = Rect{Width: childWidth, Height: childHeight}

		mainLength := child.Props.Style.Height
		if direction == Row {
			mainLength = child.Props.Style.Width
		}
		if layoutSpec, ok := effectiveLayoutSpec(child); ok && layoutSpec.Size.Width.Kind != LengthAuto && direction == Row {
			mainLength = layoutSpec.Size.Width
		}
		if layoutSpec, ok := effectiveLayoutSpec(child); ok && layoutSpec.Size.Height.Kind != LengthAuto && direction == Column {
			mainLength = layoutSpec.Size.Height
		}
		if mainLength.Kind == LengthFill {
			fillCount++
			continue
		}
		if direction == Row {
			fixedMain += childWidth
		} else {
			fixedMain += childHeight
		}
	}

	remaining := maxFloat(0, mainAvailable-gapTotal-fixedMain)
	fillSize := 0.0
	if fillCount > 0 {
		fillSize = remaining / float64(fillCount)
	}

	layout.Children = make([]*LayoutNode, 0, len(layout.Node.Children))
	localX := 0.0
	localY := 0.0
	usedWidth := 0.0
	usedHeight := 0.0

	for i, child := range layout.Node.Children {
		childWidth := resolveCrossAxisSize(child.Props.Style.Width, contentWidth, measured[i].Width)
		childHeight := resolveCrossAxisSize(child.Props.Style.Height, contentHeight, measured[i].Height)

		if childWidth == 0 && measured[i].Width > 0 {
			childWidth = measured[i].Width
		}
		if childHeight == 0 && measured[i].Height > 0 {
			childHeight = measured[i].Height
		}

		childMode := layoutModeFor(child)
		if direction == Row {
			childWidth = resolveMainAxisSize(child.Props.Style.Width, contentWidth, measured[i].Width, fillSize)
			if child.Props.Style.Height.Kind == LengthFill {
				childHeight = contentHeight
			}
			if childMode == LayoutModeAnchored {
				childX, childY := resolveAnchoredChildFrame(layout, child, baseX, baseY, contentWidth, contentHeight, childWidth, childHeight)
				childX, childY, childWidth, childHeight = applyLayoutMargin(child, childX, childY, childWidth, childHeight)
				layout.Children = append(layout.Children, buildLayout(child, childX-offsetX, childY-offsetY, childWidth, childHeight, layout.Node.Props.ID))
				usedWidth = maxFloat(usedWidth, (childX-baseX)+childWidth)
				usedHeight = maxFloat(usedHeight, (childY-baseY)+childHeight)
				continue
			}
			childX := baseX + localX - offsetX
			childY := baseY - offsetY
			childX, childY, childWidth, childHeight = applyLayoutMargin(child, childX, childY, childWidth, childHeight)
			layout.Children = append(layout.Children, buildLayout(child, childX, childY, childWidth, childHeight, layout.Node.Props.ID))
			localX += childWidth + gap
			usedWidth = maxFloat(usedWidth, localX-gap)
			usedHeight = maxFloat(usedHeight, childHeight)
			continue
		}

		childHeight = resolveMainAxisSize(child.Props.Style.Height, contentHeight, measured[i].Height, fillSize)
		if child.Props.Style.Width.Kind == LengthFill {
			childWidth = contentWidth
		}
		if childMode == LayoutModeAnchored {
			childX, childY := resolveAnchoredChildFrame(layout, child, baseX, baseY, contentWidth, contentHeight, childWidth, childHeight)
			childX, childY, childWidth, childHeight = applyLayoutMargin(child, childX, childY, childWidth, childHeight)
			layout.Children = append(layout.Children, buildLayout(child, childX-offsetX, childY-offsetY, childWidth, childHeight, layout.Node.Props.ID))
			usedWidth = maxFloat(usedWidth, (childX-baseX)+childWidth)
			usedHeight = maxFloat(usedHeight, (childY-baseY)+childHeight)
			continue
		}
		childX := baseX - offsetX
		childY := baseY + localY - offsetY
		childX, childY, childWidth, childHeight = applyLayoutMargin(child, childX, childY, childWidth, childHeight)
		layout.Children = append(layout.Children, buildLayout(child, childX, childY, childWidth, childHeight, layout.Node.Props.ID))
		localY += childHeight + gap
		usedWidth = maxFloat(usedWidth, childWidth)
		usedHeight = maxFloat(usedHeight, localY-gap)
	}

	return usedWidth, usedHeight
}

func buildGridChildren(layout *LayoutNode, baseX, baseY, contentWidth, contentHeight float64) (float64, float64) {
	return buildGridChildrenAdvanced(layout, baseX, baseY, contentWidth, contentHeight)
}

func resolveAnchoredChildFrame(parent *LayoutNode, child *Node, baseX, baseY, contentWidth, contentHeight, childWidth, childHeight float64) (float64, float64) {
	spec, _ := effectiveLayoutSpec(child)
	anchorX, anchorY := anchoredPoint(parent, spec.Anchor, baseX, baseY, contentWidth, contentHeight)
	pivotX, pivotY := pivotOffset(spec.Pivot, childWidth, childHeight)
	x := anchorX + spec.Offset.X - pivotX
	y := anchorY + spec.Offset.Y - pivotY
	if spec.Constraints.KeepInsideParent && parent != nil {
		bounds := parent.ContentBounds
		if bounds == (Rect{}) {
			bounds = Rect{X: baseX, Y: baseY, Width: contentWidth, Height: contentHeight}
		}
		x = clampFloat(x, bounds.X, bounds.X+bounds.Width-childWidth)
		y = clampFloat(y, bounds.Y, bounds.Y+bounds.Height-childHeight)
	}
	return x, y
}

func applyLayoutMargin(node *Node, x, y, width, height float64) (float64, float64, float64, float64) {
	spec, ok := effectiveLayoutSpec(node)
	if !ok {
		return x, y, width, height
	}
	x += spec.Margin.Left
	y += spec.Margin.Top
	return x, y, width, height
}

func anchoredPoint(parent *LayoutNode, anchor Anchor, baseX, baseY, contentWidth, contentHeight float64) (float64, float64) {
	bounds := parent.ContentBounds
	if bounds == (Rect{}) {
		bounds = Rect{X: baseX, Y: baseY, Width: contentWidth, Height: contentHeight}
	}
	switch anchor {
	case AnchorTop:
		return bounds.X + bounds.Width*0.5, bounds.Y
	case AnchorTopRight:
		return bounds.X + bounds.Width, bounds.Y
	case AnchorLeft:
		return bounds.X, bounds.Y + bounds.Height*0.5
	case AnchorCenter, "":
		return bounds.X + bounds.Width*0.5, bounds.Y + bounds.Height*0.5
	case AnchorRight:
		return bounds.X + bounds.Width, bounds.Y + bounds.Height*0.5
	case AnchorBottomLeft:
		return bounds.X, bounds.Y + bounds.Height
	case AnchorBottom:
		return bounds.X + bounds.Width*0.5, bounds.Y + bounds.Height
	case AnchorBottomRight:
		return bounds.X + bounds.Width, bounds.Y + bounds.Height
	default:
		return bounds.X, bounds.Y
	}
}

func pivotOffset(pivot Pivot, width, height float64) (float64, float64) {
	switch pivot {
	case PivotTop:
		return width * 0.5, 0
	case PivotTopRight:
		return width, 0
	case PivotLeft:
		return 0, height * 0.5
	case PivotCenter, "":
		return width * 0.5, height * 0.5
	case PivotRight:
		return width, height * 0.5
	case PivotBottomLeft:
		return 0, height
	case PivotBottom:
		return width * 0.5, height
	case PivotBottomRight:
		return width, height
	default:
		return 0, 0
	}
}

func layoutDirectionFor(node *Node) Direction {
	spec, ok := effectiveLayoutSpec(node)
	if ok {
		switch spec.Mode {
		case LayoutModeFlowHorizontal:
			return Row
		case LayoutModeFlowVertical, LayoutModeAnchored, LayoutModeGrid, LayoutModeStack:
			return Column
		}
	}
	if node != nil && node.Tag == TagStack {
		return Column
	}
	if node != nil && node.Props.Style.Direction == Row {
		return Row
	}
	return Column
}

func layoutModeFor(node *Node) LayoutMode {
	if node == nil {
		return LayoutModeFlowVertical
	}
	if spec, ok := effectiveLayoutSpec(node); ok && spec.Mode != "" {
		return spec.Mode
	}
	if node.Tag == TagStack {
		return LayoutModeStack
	}
	if node.Tag == TagScrollView {
		return LayoutModeFlowVertical
	}
	if node.Props.Style.Direction == Row {
		return LayoutModeFlowHorizontal
	}
	return LayoutModeFlowVertical
}

func effectiveLayoutSpec(node *Node) (LayoutSpec, bool) {
	if node == nil {
		return LayoutSpec{}, false
	}
	spec := node.Props.Layout
	if spec == (LayoutSpec{}) {
		return LayoutSpec{}, false
	}
	return spec, true
}

func resolveLayoutSize(node *Node, availableWidth, availableHeight, measuredWidth, measuredHeight float64) (float64, float64) {
	spec, ok := effectiveLayoutSpec(node)
	if !ok {
		return measuredWidth, measuredHeight
	}

	width := measuredWidth
	height := measuredHeight
	if spec.Size.Width.Kind != LengthAuto {
		width = resolveLayoutLength(spec.Size.Width, availableWidth, measuredWidth)
	}
	if spec.Size.Height.Kind != LengthAuto {
		height = resolveLayoutLength(spec.Size.Height, availableHeight, measuredHeight)
	}
	if spec.MinSize.Width.Kind != LengthAuto {
		width = maxFloat(width, resolveLayoutLength(spec.MinSize.Width, availableWidth, measuredWidth))
	}
	if spec.MinSize.Height.Kind != LengthAuto {
		height = maxFloat(height, resolveLayoutLength(spec.MinSize.Height, availableHeight, measuredHeight))
	}
	if spec.MaxSize.Width.Kind != LengthAuto {
		width = minFloat(width, resolveLayoutLength(spec.MaxSize.Width, availableWidth, measuredWidth))
	}
	if spec.MaxSize.Height.Kind != LengthAuto {
		height = minFloat(height, resolveLayoutLength(spec.MaxSize.Height, availableHeight, measuredHeight))
	}
	return width, height
}

func resolveLayoutLength(length Length, available, intrinsic float64) float64 {
	switch length.Kind {
	case LengthPx:
		return length.Value
	case LengthFill:
		return available
	default:
		return intrinsic
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func mergeInsets(primary, fallback Insets) Insets {
	if primary == (Insets{}) {
		return fallback
	}
	result := fallback
	if primary.Top != 0 {
		result.Top = primary.Top
	}
	if primary.Right != 0 {
		result.Right = primary.Right
	}
	if primary.Bottom != 0 {
		result.Bottom = primary.Bottom
	}
	if primary.Left != 0 {
		result.Left = primary.Left
	}
	return result
}

func measureNode(node *Node, availableWidth, availableHeight float64) (float64, float64) {
	if node == nil {
		return 0, 0
	}

	style := node.Props.Style
	gap := effectiveGap(node)
	padding := effectivePadding(node)
	switch node.Tag {
	case TagText:
		textWidth, textHeight := measureText(node.Text, style)
		width := resolveCrossAxisSize(style.Width, availableWidth, textWidth)
		height := resolveCrossAxisSize(style.Height, availableHeight, textHeight)
		return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
	case TagTextBlock:
		textWidth, textHeight := measureTextBlock(node.Text, style, availableWidth)
		width := resolveCrossAxisSize(style.Width, availableWidth, textWidth)
		height := resolveCrossAxisSize(style.Height, availableHeight, textHeight)
		return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
	case TagImage:
		imageWidth, imageHeight := node.Props.Image.intrinsicSize()
		width := resolveCrossAxisSize(style.Width, availableWidth, imageWidth)
		height := resolveCrossAxisSize(style.Height, availableHeight, imageHeight)
		return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
	}

	if len(node.Children) == 0 {
		width := resolveCrossAxisSize(style.Width, availableWidth, 0)
		height := resolveCrossAxisSize(style.Height, availableHeight, 0)
		return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
	}

	switch layoutModeFor(node) {
	case LayoutModeStack:
		return measureStack(node, availableWidth, availableHeight)
	case LayoutModeGrid:
		return measureGrid(node, availableWidth, availableHeight)
	}

	direction := style.directionOrDefault()
	contentWidth := maxFloat(0, availableWidth-padding.Horizontal())
	contentHeight := maxFloat(0, availableHeight-padding.Vertical())

	mainTotal := 0.0
	crossMax := 0.0
	gapTotal := 0.0
	visibleChildren := 0

	for _, child := range node.Children {
		childWidth, childHeight := measureNode(child, contentWidth, contentHeight)
		mainLength := childHeight
		crossLength := childWidth
		mainStyleLength := child.Props.Style.Height

		if direction == Row {
			mainLength = childWidth
			crossLength = childHeight
			mainStyleLength = child.Props.Style.Width
		}

		if visibleChildren > 0 {
			gapTotal += gap
		}
		visibleChildren++

		if mainStyleLength.Kind != LengthFill {
			mainTotal += mainLength
		}
		crossMax = maxFloat(crossMax, crossLength)
	}

	intrinsicWidth := padding.Horizontal()
	intrinsicHeight := padding.Vertical()

	if direction == Row {
		intrinsicWidth += mainTotal + gapTotal
		intrinsicHeight += crossMax
	} else {
		intrinsicWidth += crossMax
		intrinsicHeight += mainTotal + gapTotal
	}

	width := resolveCrossAxisSize(style.Width, availableWidth, intrinsicWidth)
	height := resolveCrossAxisSize(style.Height, availableHeight, intrinsicHeight)
	return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
}

func measureStack(node *Node, availableWidth, availableHeight float64) (float64, float64) {
	style := node.Props.Style
	padding := effectivePadding(node)
	contentWidth := maxFloat(0, availableWidth-padding.Horizontal())
	contentHeight := maxFloat(0, availableHeight-padding.Vertical())
	maxWidth := 0.0
	maxHeight := 0.0

	for _, child := range node.Children {
		childWidth, childHeight := measureNode(child, contentWidth, contentHeight)
		maxWidth = maxFloat(maxWidth, childWidth)
		maxHeight = maxFloat(maxHeight, childHeight)
	}

	width := resolveCrossAxisSize(style.Width, availableWidth, maxWidth+padding.Horizontal())
	height := resolveCrossAxisSize(style.Height, availableHeight, maxHeight+padding.Vertical())
	return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
}

func measureGrid(node *Node, availableWidth, availableHeight float64) (float64, float64) {
	return measureGridAdvanced(node, availableWidth, availableHeight)
}

func effectivePadding(node *Node) Insets {
	if node == nil {
		return Insets{}
	}
	padding := node.Props.Style.Padding
	if spec, ok := effectiveLayoutSpec(node); ok {
		padding = mergeInsets(spec.Padding, padding)
	}
	return padding
}

func effectiveGap(node *Node) float64 {
	if node == nil {
		return 0
	}
	gap := node.Props.Style.Gap
	if spec, ok := effectiveLayoutSpec(node); ok && spec.Gap != 0 {
		gap = spec.Gap
	}
	return gap
}

func measureText(value string, style Style) (float64, float64) {
	face := defaultFontFace()
	advance := font.MeasureString(face, value)
	metrics := face.Metrics()

	width := float64(advance.Ceil())
	height := float64(metrics.Height.Ceil())

	if width == 0 {
		width = 1
	}
	if height == 0 {
		height = 1
	}
	return width, height
}

func measureTextBlock(value string, style Style, availableWidth float64) (float64, float64) {
	lines := wrapTextLines(value, style, widthConstraint(style, availableWidth))
	if len(lines) == 0 {
		lines = []string{""}
	}

	width := 0.0
	for _, line := range lines {
		lineWidth, _ := measureText(line, style)
		width = maxFloat(width, lineWidth)
	}

	lineHeight := textLineHeight(style)
	if constrained := widthConstraint(style, availableWidth); constrained > 0 {
		width = constrained
	}

	return width, lineHeight * float64(len(lines))
}

func defaultFontFace() font.Face {
	return basicfont.Face7x13
}

func textLineHeight(style Style) float64 {
	if style.LineHeight > 0 {
		return style.LineHeight
	}
	return float64(defaultFontFace().Metrics().Height.Ceil())
}

func widthConstraint(style Style, availableWidth float64) float64 {
	switch style.Width.Kind {
	case LengthPx:
		return style.Width.Value
	case LengthFill:
		return availableWidth
	default:
		return 0
	}
}

func wrapTextLines(value string, style Style, availableWidth float64) []string {
	if value == "" {
		return []string{""}
	}

	normalized := strings.ReplaceAll(value, "\r\n", "\n")
	if availableWidth <= 0 {
		return strings.Split(normalized, "\n")
	}

	result := make([]string, 0)
	for _, rawLine := range strings.Split(normalized, "\n") {
		result = append(result, wrapLine(rawLine, availableWidth, style)...)
	}
	if len(result) == 0 {
		return []string{""}
	}
	return result
}

func wrapLine(line string, maxWidth float64, style Style) []string {
	if line == "" {
		return []string{""}
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0)
	current := words[0]
	for _, word := range words[1:] {
		candidate := current + " " + word
		candidateWidth, _ := measureText(candidate, style)
		if candidateWidth <= maxWidth {
			current = candidate
			continue
		}
		lines = append(lines, current)
		current = word
	}
	lines = append(lines, current)
	return lines
}

func resolveMainAxisSize(length Length, available, intrinsic, fillSize float64) float64 {
	switch length.Kind {
	case LengthPx:
		return length.Value
	case LengthFill:
		return fillSize
	default:
		return intrinsic
	}
}

func resolveCrossAxisSize(length Length, available, intrinsic float64) float64 {
	switch length.Kind {
	case LengthPx:
		return length.Value
	case LengthFill:
		return available
	default:
		return intrinsic
	}
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
