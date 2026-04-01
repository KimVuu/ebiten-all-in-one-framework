package ebitenui

import "math"

func ValidateLayout(root *LayoutNode, viewport Viewport, opts ValidationOptions) LayoutValidationReport {
	report := LayoutValidationReport{}
	if root == nil {
		return report
	}

	viewportRect := Rect{X: 0, Y: 0, Width: viewport.Width, Height: viewport.Height}
	validateLayoutNode(root, nil, viewportRect, opts.SafeArea, &report)
	return report
}

func validateLayoutNode(node *LayoutNode, parent *LayoutNode, viewport Rect, safeArea Insets, report *LayoutValidationReport) {
	if node == nil || node.Node == nil {
		return
	}

	frame := node.Frame
	spec, hasLayout := effectiveLayoutSpec(node.Node)
	skipPositionalValidation := parent != nil && parent.ClipChildren

	if !skipPositionalValidation && !rectWithin(frame, viewport) {
		report.Issues = append(report.Issues, LayoutIssue{
			NodeID:   node.Node.Props.ID,
			Severity: IssueSeverityError,
			Code:     IssueOutOfViewport,
			Message:  "node extends beyond viewport",
			SuggestedConstraintChanges: []ConstraintPatch{
				{Field: "layout.keepInsideParent", Value: true, Note: "keep node visible within viewport"},
				{Field: "layout.offset", Value: Point{X: clampFloat(frame.X, viewport.X, viewport.Width), Y: clampFloat(frame.Y, viewport.Y, viewport.Height)}, Note: "reposition inside viewport"},
			},
		})
	}

	if parent != nil && !skipPositionalValidation && !rectWithin(frame, parent.Frame) {
		report.Issues = append(report.Issues, LayoutIssue{
			NodeID:   node.Node.Props.ID,
			Severity: IssueSeverityError,
			Code:     IssueOutOfParent,
			Message:  "node extends beyond parent bounds",
			SuggestedConstraintChanges: []ConstraintPatch{
				{Field: "layout.keepInsideParent", Value: true, Note: "keep node inside parent bounds"},
				{Field: "layout.anchor", Value: AnchorCenter, Note: "re-anchor inside parent"},
			},
		})
	}

	if safeArea != (Insets{}) && !skipPositionalValidation {
		safeRect := Rect{
			X:      viewport.X + safeArea.Left,
			Y:      viewport.Y + safeArea.Top,
			Width:  maxFloat(0, viewport.Width-safeArea.Horizontal()),
			Height: maxFloat(0, viewport.Height-safeArea.Vertical()),
		}
		if !rectWithin(frame, safeRect) {
			report.Issues = append(report.Issues, LayoutIssue{
				NodeID:   node.Node.Props.ID,
				Severity: IssueSeverityWarning,
				Code:     IssueSafeAreaViolation,
				Message:  "node intersects safe area boundary",
				SuggestedConstraintChanges: []ConstraintPatch{
					{Field: "layout.offset", Value: Point{X: safeRect.X, Y: safeRect.Y}, Note: "nudge node inside safe area"},
				},
			})
		}
	}

	minHitTarget := spec.Constraints.MinHitTarget
	if minHitTarget <= 0 {
		minHitTarget = 44
	}
	if hasLayout || node.Node.Props.Focusable || node.Node.Tag == TagButton {
		if math.Min(frame.Width, frame.Height) < minHitTarget {
			report.Issues = append(report.Issues, LayoutIssue{
				NodeID:   node.Node.Props.ID,
				Severity: IssueSeverityWarning,
				Code:     IssueMinHitTarget,
				Message:  "node hit target is too small",
				SuggestedConstraintChanges: []ConstraintPatch{
					{Field: "layout.minSize", Value: LayoutSize{Width: Px(minHitTarget), Height: Px(minHitTarget)}, Note: "increase interactive hit target"},
				},
			})
		}
	}

	if node.Node.Tag == TagText || node.Node.Tag == TagTextBlock {
		overflow := textOverflowDetected(node)
		if overflow {
			report.Issues = append(report.Issues, LayoutIssue{
				NodeID:   node.Node.Props.ID,
				Severity: IssueSeverityWarning,
				Code:     IssueTextOverflow,
				Message:  "text exceeds available bounds",
				SuggestedConstraintChanges: []ConstraintPatch{
					{Field: "layout.size.width", Value: maxFloat(frame.Width, textWidthEstimate(node)), Note: "widen text container"},
				},
			})
		}
	}

	if len(node.Children) > 1 {
		for i := 0; i < len(node.Children); i++ {
			left := node.Children[i]
			if left == nil {
				continue
			}
			for j := i + 1; j < len(node.Children); j++ {
				right := node.Children[j]
				if right == nil {
					continue
				}
				if !rectsOverlap(left.Frame, right.Frame) {
					continue
				}
				leftAllow := left.Node != nil && effectiveOverlapAllowed(left.Node)
				rightAllow := right.Node != nil && effectiveOverlapAllowed(right.Node)
				if leftAllow || rightAllow {
					continue
				}
				report.Issues = append(report.Issues, LayoutIssue{
					NodeID:   right.Node.Props.ID,
					Severity: IssueSeverityWarning,
					Code:     IssueOverlap,
					Message:  "sibling nodes overlap",
					SuggestedConstraintChanges: []ConstraintPatch{
						{Field: "layout.offset", Value: Point{X: right.Frame.X + 8, Y: right.Frame.Y + 8}, Note: "separate overlapping nodes"},
					},
				})
				if left.Node.Props.Layout.ZIndex == right.Node.Props.Layout.ZIndex {
					report.Issues = append(report.Issues, LayoutIssue{
						NodeID:   right.Node.Props.ID,
						Severity: IssueSeverityWarning,
						Code:     IssueZOrderConflict,
						Message:  "sibling nodes overlap with same z-index",
						SuggestedConstraintChanges: []ConstraintPatch{
							{Field: "layout.zIndex", Value: right.Node.Props.Layout.ZIndex + 1, Note: "raise node above overlapping sibling"},
						},
					})
				}
			}
		}
	}

	for _, child := range node.Children {
		validateLayoutNode(child, node, viewport, safeArea, report)
	}
}

func rectWithin(child, parent Rect) bool {
	if child.Width < 0 || child.Height < 0 || parent.Width < 0 || parent.Height < 0 {
		return false
	}
	return child.X >= parent.X &&
		child.Y >= parent.Y &&
		child.X+child.Width <= parent.X+parent.Width &&
		child.Y+child.Height <= parent.Y+parent.Height
}

func rectsOverlap(a, b Rect) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

func effectiveOverlapAllowed(node *Node) bool {
	if node == nil {
		return false
	}
	spec, ok := effectiveLayoutSpec(node)
	if !ok {
		return false
	}
	return spec.Constraints.AllowOverlap
}

func textOverflowDetected(layout *LayoutNode) bool {
	if layout == nil || layout.Node == nil {
		return false
	}
	width := layout.Frame.Width
	if width <= 0 {
		return false
	}
	if layout.Node.Tag == TagText {
		estimated, _ := measureText(layout.Node.Text, layout.Node.Props.Style)
		return estimated > width
	}
	if len(layout.TextLines) == 0 {
		return false
	}
	maxWidth := 0.0
	for _, line := range layout.TextLines {
		lineWidth, _ := measureText(line, layout.Node.Props.Style)
		maxWidth = maxFloat(maxWidth, lineWidth)
	}
	return maxWidth > width
}

func textWidthEstimate(layout *LayoutNode) float64 {
	if layout == nil || layout.Node == nil {
		return 0
	}
	if layout.Node.Tag == TagText {
		width, _ := measureText(layout.Node.Text, layout.Node.Props.Style)
		return width
	}
	maxWidth := 0.0
	for _, line := range layout.TextLines {
		lineWidth, _ := measureText(line, layout.Node.Props.Style)
		maxWidth = maxFloat(maxWidth, lineWidth)
	}
	return maxWidth
}
