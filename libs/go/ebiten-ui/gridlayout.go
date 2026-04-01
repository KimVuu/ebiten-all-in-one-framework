package ebitenui

type gridPlacement struct {
	child          *Node
	rowStart       int
	colStart       int
	rowSpan        int
	colSpan        int
	measuredWidth  float64
	measuredHeight float64
	justifySelf    LayoutAlignment
	alignSelf      LayoutAlignment
	explicitRow    bool
	explicitColumn bool
}

type gridAxisDistribution struct {
	startOffset float64
	gap         float64
	sizes       []float64
}

func buildGridChildrenAdvanced(layout *LayoutNode, baseX, baseY, contentWidth, contentHeight float64) (float64, float64) {
	if layout == nil || layout.Node == nil {
		return 0, 0
	}

	node := layout.Node
	gridSpec := node.Props.Layout.Grid
	columns := gridSpec.Columns
	if columns <= 0 {
		columns = 2
	}
	gap := gridSpec.Gap
	if gap <= 0 {
		gap = node.Props.Style.Gap
	}
	placements, rowCount := planGridPlacements(node, columns, gridSpec.AutoFlow)
	colWidths, rowHeights := measureGridTracks(node, placements, columns, contentWidth, contentHeight, gap)

	if gridSpec.Rows > 0 && gridSpec.Rows > rowCount {
		for len(rowHeights) < gridSpec.Rows {
			rowHeights = append(rowHeights, 0)
		}
		rowCount = gridSpec.Rows
	}

	horizontal := distributeGridAxis(colWidths, contentWidth, gap, firstGridAlignment(gridSpec.JustifyContent, LayoutAlignmentStretch))
	vertical := distributeGridAxis(rowHeights, contentHeight, gap, firstGridAlignment(gridSpec.AlignContent, LayoutAlignmentStretch))

	colOrigins := accumulateGridOrigins(baseX, horizontal)
	rowOrigins := accumulateGridOrigins(baseY, vertical)

	layout.Children = make([]*LayoutNode, 0, len(placements))
	usedWidth := 0.0
	usedHeight := 0.0

	for _, placement := range placements {
		rowStart := clampInt(placement.rowStart, 1, len(rowOrigins))
		colStart := clampInt(placement.colStart, 1, len(colOrigins))
		rowEnd := minInt(len(rowOrigins), rowStart+placement.rowSpan-1)
		colEnd := minInt(len(colOrigins), colStart+placement.colSpan-1)

		areaX := colOrigins[colStart-1]
		areaY := rowOrigins[rowStart-1]
		areaWidth := spanTrackSize(colWidths, horizontal.gap, colStart, colEnd)
		areaHeight := spanTrackSize(rowHeights, vertical.gap, rowStart, rowEnd)

		childWidth := placement.measuredWidth
		childHeight := placement.measuredHeight
		childWidth, childHeight = resolveLayoutSize(placement.child, areaWidth, areaHeight, childWidth, childHeight)

		justification := placement.justifySelf
		if justification == "" {
			justification = gridSpec.JustifyItems
		}
		alignment := placement.alignSelf
		if alignment == "" {
			alignment = gridSpec.AlignItems
		}

		childX, childY, childWidth, childHeight := alignGridChildWithinArea(areaX, areaY, areaWidth, areaHeight, childWidth, childHeight, justification, alignment)
		if layoutModeFor(placement.child) == LayoutModeAnchored {
			cellParent := &LayoutNode{
				Node: layout.Node,
				ContentBounds: Rect{
					X:      areaX,
					Y:      areaY,
					Width:  areaWidth,
					Height: areaHeight,
				},
			}
			childX, childY = resolveAnchoredChildFrame(cellParent, placement.child, areaX, areaY, areaWidth, areaHeight, childWidth, childHeight)
			childX, childY, childWidth, childHeight = applyLayoutMargin(placement.child, childX, childY, childWidth, childHeight)
		} else {
			childX, childY, childWidth, childHeight = applyLayoutMargin(placement.child, childX, childY, childWidth, childHeight)
		}

		layout.Children = append(layout.Children, buildLayout(placement.child, childX, childY, childWidth, childHeight, layout.Node.Props.ID))
		usedWidth = maxFloat(usedWidth, (childX-baseX)+childWidth)
		usedHeight = maxFloat(usedHeight, (childY-baseY)+childHeight)
	}

	if len(colOrigins) > 0 && len(rowOrigins) > 0 {
		lastX := colOrigins[len(colOrigins)-1]
		lastY := rowOrigins[len(rowOrigins)-1]
		usedWidth = maxFloat(usedWidth, (lastX-baseX)+colWidths[len(colWidths)-1])
		usedHeight = maxFloat(usedHeight, (lastY-baseY)+rowHeights[len(rowHeights)-1])
	}

	return usedWidth, usedHeight
}

func measureGridAdvanced(node *Node, availableWidth, availableHeight float64) (float64, float64) {
	if node == nil {
		return 0, 0
	}

	gridSpec := node.Props.Layout.Grid
	columns := gridSpec.Columns
	if columns <= 0 {
		columns = 2
	}
	gap := gridSpec.Gap
	if gap <= 0 {
		gap = node.Props.Style.Gap
	}
	contentWidth := maxFloat(0, availableWidth)
	contentHeight := maxFloat(0, availableHeight)

	placements, rowCount := planGridPlacements(node, columns, gridSpec.AutoFlow)
	colWidths, rowHeights := measureGridTracks(node, placements, columns, contentWidth, contentHeight, gap)
	if gridSpec.Rows > 0 && gridSpec.Rows > rowCount {
		for len(rowHeights) < gridSpec.Rows {
			rowHeights = append(rowHeights, 0)
		}
		rowCount = gridSpec.Rows
	}

	horizontal := distributeGridAxis(colWidths, contentWidth, gap, firstGridAlignment(gridSpec.JustifyContent, LayoutAlignmentStretch))
	vertical := distributeGridAxis(rowHeights, contentHeight, gap, firstGridAlignment(gridSpec.AlignContent, LayoutAlignmentStretch))

	totalWidth := sumGridAxis(horizontal.sizes, horizontal.gap)
	totalHeight := sumGridAxis(vertical.sizes, vertical.gap)
	padding := effectivePadding(node)
	width := resolveCrossAxisSize(node.Props.Style.Width, availableWidth, totalWidth+padding.Horizontal())
	height := resolveCrossAxisSize(node.Props.Style.Height, availableHeight, totalHeight+padding.Vertical())
	return resolveLayoutSize(node, availableWidth, availableHeight, width, height)
}

func planGridPlacements(node *Node, columns int, autoFlow LayoutAutoFlow) ([]gridPlacement, int) {
	placements := make([]gridPlacement, 0, len(node.Children))
	if columns <= 0 {
		columns = 2
	}

	occupied := make([][]bool, 0)
	nextRow := 1
	nextColumn := 1

	for _, child := range node.Children {
		spec, _ := effectiveLayoutSpec(child)
		placement := gridPlacementFromSpec(child, spec)
		if placement.colSpan > columns {
			placement.colSpan = columns
		}
		if placement.rowSpan < 1 {
			placement.rowSpan = 1
		}
		if placement.colSpan < 1 {
			placement.colSpan = 1
		}

		switch {
		case placement.explicitRow && placement.explicitColumn:
			placement.rowStart, placement.colStart = fitExplicitGridPlacement(occupied, columns, placement.rowStart, placement.colStart, placement.rowSpan, placement.colSpan, autoFlow, nextRow, nextColumn)
		case placement.explicitRow:
			placement.colStart = fitExplicitRow(occupied, columns, placement.rowStart, placement.colSpan)
			if placement.colStart == 0 {
				placement.rowStart, placement.colStart = findNextGridPlacement(occupied, columns, placement.rowSpan, placement.colSpan, autoFlow, nextRow, nextColumn)
			}
		case placement.explicitColumn:
			placement.rowStart = fitExplicitColumn(occupied, columns, placement.colStart, placement.rowSpan)
			if placement.rowStart == 0 {
				placement.rowStart, placement.colStart = findNextGridPlacement(occupied, columns, placement.rowSpan, placement.colSpan, autoFlow, nextRow, nextColumn)
			}
		default:
			placement.rowStart, placement.colStart = findNextGridPlacement(occupied, columns, placement.rowSpan, placement.colSpan, autoFlow, nextRow, nextColumn)
		}

		if placement.rowStart <= 0 {
			placement.rowStart = 1
		}
		if placement.colStart <= 0 {
			placement.colStart = 1
		}
		markGridPlacement(&occupied, placement.rowStart, placement.colStart, placement.rowSpan, placement.colSpan, columns)
		nextRow = placement.rowStart
		nextColumn = placement.colStart
		placements = append(placements, placement)
	}

	rowCount := len(occupied)
	if rowCount == 0 {
		rowCount = 1
	}
	return placements, rowCount
}

func gridPlacementFromSpec(child *Node, spec LayoutSpec) gridPlacement {
	placement := gridPlacement{
		child:          child,
		rowStart:       spec.Grid.RowStart,
		colStart:       spec.Grid.ColumnStart,
		rowSpan:        maxInt(1, spec.Grid.RowSpan),
		colSpan:        maxInt(1, spec.Grid.ColumnSpan),
		justifySelf:    spec.Grid.JustifySelf,
		alignSelf:      spec.Grid.AlignSelf,
		explicitRow:    spec.Grid.RowStart > 0,
		explicitColumn: spec.Grid.ColumnStart > 0,
	}
	if spec.GridRow > 0 {
		placement.rowStart = spec.GridRow
		placement.explicitRow = true
	}
	if spec.GridColumn > 0 {
		placement.colStart = spec.GridColumn
		placement.explicitColumn = true
	}
	if spec.GridRowSpan > 0 {
		placement.rowSpan = spec.GridRowSpan
	}
	if spec.GridColumnSpan > 0 {
		placement.colSpan = spec.GridColumnSpan
	}
	return placement
}

func measureGridTracks(node *Node, placements []gridPlacement, columns int, contentWidth, contentHeight, gap float64) ([]float64, []float64) {
	if columns <= 0 {
		columns = 2
	}
	colWidths := make([]float64, columns)
	rowHeights := make([]float64, 0)
	cellWidth := maxFloat(0, (contentWidth-gap*float64(maxInt(0, columns-1)))/float64(columns))
	for i := range colWidths {
		colWidths[i] = cellWidth
	}

	for i := range placements {
		placement := &placements[i]
		spanWidth := cellWidth*float64(placement.colSpan) + gap*float64(maxInt(0, placement.colSpan-1))
		spanHeight := contentHeight
		if placement.rowSpan > 0 {
			spanHeight = contentHeight
		}
		childWidth, childHeight := measureNode(placement.child, spanWidth, spanHeight)
		childWidth, childHeight = resolveLayoutSize(placement.child, spanWidth, spanHeight, childWidth, childHeight)
		placement.measuredWidth = childWidth
		placement.measuredHeight = childHeight

		shareWidth := childWidth / float64(maxInt(1, placement.colSpan))
		for col := placement.colStart; col < placement.colStart+placement.colSpan; col++ {
			if col < 1 || col > columns {
				continue
			}
			colWidths[col-1] = maxFloat(colWidths[col-1], shareWidth)
		}

		rowEnd := placement.rowStart + placement.rowSpan - 1
		for len(rowHeights) < rowEnd {
			rowHeights = append(rowHeights, 0)
		}
		shareHeight := childHeight / float64(maxInt(1, placement.rowSpan))
		for row := placement.rowStart; row <= rowEnd; row++ {
			rowHeights[row-1] = maxFloat(rowHeights[row-1], shareHeight)
		}
	}

	for i := range rowHeights {
		if rowHeights[i] == 0 {
			rowHeights[i] = 0
		}
	}

	return colWidths, rowHeights
}

func firstGridAlignment(primary, fallback LayoutAlignment) LayoutAlignment {
	if primary != "" {
		return primary
	}
	return fallback
}

func distributeGridAxis(sizes []float64, available, gap float64, align LayoutAlignment) gridAxisDistribution {
	result := gridAxisDistribution{
		sizes: append([]float64{}, sizes...),
		gap:   gap,
	}
	if len(result.sizes) == 0 {
		return result
	}

	used := sumGridAxis(result.sizes, gap)
	extra := maxFloat(0, available-used)
	switch align {
	case LayoutAlignmentEnd:
		result.startOffset = extra
	case LayoutAlignmentCenter:
		result.startOffset = extra / 2
	case LayoutAlignmentStretch:
		share := extra / float64(len(result.sizes))
		for i := range result.sizes {
			result.sizes[i] += share
		}
	case LayoutAlignmentSpaceBetween:
		if len(result.sizes) > 1 {
			result.gap = gap + extra/float64(len(result.sizes)-1)
		} else {
			result.startOffset = extra / 2
		}
	case LayoutAlignmentSpaceAround:
		if len(result.sizes) > 0 {
			result.gap = gap + extra/float64(len(result.sizes))
			result.startOffset = result.gap / 2
		}
	case LayoutAlignmentSpaceEvenly:
		result.gap = gap + extra/float64(len(result.sizes)+1)
		result.startOffset = result.gap
	}
	return result
}

func accumulateGridOrigins(origin float64, axis gridAxisDistribution) []float64 {
	origins := make([]float64, len(axis.sizes))
	position := origin + axis.startOffset
	for i, size := range axis.sizes {
		origins[i] = position
		position += size
		if i < len(axis.sizes)-1 {
			position += axis.gap
		}
	}
	return origins
}

func spanTrackSize(sizes []float64, gap float64, start, end int) float64 {
	if start < 1 || end < start {
		return 0
	}
	start--
	end--
	if start >= len(sizes) {
		return 0
	}
	if end >= len(sizes) {
		end = len(sizes) - 1
	}
	total := 0.0
	for i := start; i <= end; i++ {
		total += sizes[i]
		if i < end {
			total += gap
		}
	}
	return total
}

func alignGridChildWithinArea(areaX, areaY, areaWidth, areaHeight, childWidth, childHeight float64, justify, align LayoutAlignment) (float64, float64, float64, float64) {
	childX := areaX
	childY := areaY

	if justify == "" {
		justify = LayoutAlignmentStretch
	}
	if align == "" {
		align = LayoutAlignmentStretch
	}

	switch justify {
	case LayoutAlignmentCenter:
		childX = areaX + maxFloat(0, (areaWidth-childWidth)/2)
	case LayoutAlignmentEnd:
		childX = areaX + maxFloat(0, areaWidth-childWidth)
	case LayoutAlignmentStretch:
		childWidth = areaWidth
	}

	switch align {
	case LayoutAlignmentCenter:
		childY = areaY + maxFloat(0, (areaHeight-childHeight)/2)
	case LayoutAlignmentEnd:
		childY = areaY + maxFloat(0, areaHeight-childHeight)
	case LayoutAlignmentStretch:
		childHeight = areaHeight
	}

	return childX, childY, childWidth, childHeight
}

func fitExplicitGridPlacement(occupied [][]bool, columns, rowStart, colStart, rowSpan, colSpan int, autoFlow LayoutAutoFlow, nextRow, nextColumn int) (int, int) {
	if canGridPlace(occupied, rowStart, colStart, rowSpan, colSpan, columns) {
		return rowStart, colStart
	}
	if autoFlow == LayoutAutoFlowColumn {
		if row := fitExplicitColumn(occupied, columns, colStart, rowSpan); row > 0 {
			return row, colStart
		}
	}
	if col := fitExplicitRow(occupied, columns, rowStart, colSpan); col > 0 {
		return rowStart, col
	}
	return findNextGridPlacement(occupied, columns, rowSpan, colSpan, autoFlow, nextRow, nextColumn)
}

func fitExplicitRow(occupied [][]bool, columns, row, colSpan int) int {
	for col := 1; col <= columns-colSpan+1; col++ {
		if canGridPlace(occupied, row, col, 1, colSpan, columns) {
			return col
		}
	}
	return 0
}

func fitExplicitColumn(occupied [][]bool, columns, column, rowSpan int) int {
	row := 1
	for {
		if canGridPlace(occupied, row, column, rowSpan, 1, columns) {
			return row
		}
		row++
	}
}

func findNextGridPlacement(occupied [][]bool, columns, rowSpan, colSpan int, autoFlow LayoutAutoFlow, nextRow, nextColumn int) (int, int) {
	if autoFlow == LayoutAutoFlowColumn {
		return findNextGridPlacementColumnMajor(occupied, columns, rowSpan, colSpan, nextRow, nextColumn)
	}
	return findNextGridPlacementRowMajor(occupied, columns, rowSpan, colSpan, nextRow, nextColumn)
}

func findNextGridPlacementRowMajor(occupied [][]bool, columns, rowSpan, colSpan, startRow, startColumn int) (int, int) {
	if colSpan > columns {
		colSpan = columns
	}
	row := maxInt(1, startRow)
	for {
		ensureGridRows(&occupied, row+rowSpan-1, columns)
		start := 1
		if row == startRow {
			start = maxInt(1, startColumn)
		}
		for col := start; col <= columns-colSpan+1; col++ {
			if canGridPlace(occupied, row, col, rowSpan, colSpan, columns) {
				return row, col
			}
		}
		row++
	}
}

func findNextGridPlacementColumnMajor(occupied [][]bool, columns, rowSpan, colSpan, startRow, startColumn int) (int, int) {
	if colSpan > columns {
		colSpan = columns
	}
	startColumn = maxInt(1, startColumn)
	for col := startColumn; col <= columns-colSpan+1; col++ {
		row := maxInt(1, startRow)
		for {
			ensureGridRows(&occupied, row+rowSpan-1, columns)
			if canGridPlace(occupied, row, col, rowSpan, colSpan, columns) {
				return row, col
			}
			row++
		}
	}
	return findNextGridPlacementRowMajor(occupied, columns, rowSpan, colSpan, startRow, 1)
}

func canGridPlace(occupied [][]bool, rowStart, colStart, rowSpan, colSpan, columns int) bool {
	if rowStart < 1 || colStart < 1 {
		return false
	}
	if colStart+colSpan-1 > columns {
		return false
	}
	for row := rowStart; row < rowStart+rowSpan; row++ {
		if row > len(occupied) {
			continue
		}
		for col := colStart; col < colStart+colSpan; col++ {
			if col <= 0 || col >= len(occupied[row-1]) {
				return false
			}
			if occupied[row-1][col] {
				return false
			}
		}
	}
	return true
}

func markGridPlacement(occupied *[][]bool, rowStart, colStart, rowSpan, colSpan, columns int) {
	ensureGridRows(occupied, rowStart+rowSpan-1, columns)
	for row := rowStart; row < rowStart+rowSpan; row++ {
		for col := colStart; col < colStart+colSpan; col++ {
			(*occupied)[row-1][col] = true
		}
	}
}

func ensureGridRows(occupied *[][]bool, rows, columns int) {
	for len(*occupied) < rows {
		row := make([]bool, columns+1)
		*occupied = append(*occupied, row)
	}
}

func sumGridAxis(sizes []float64, gap float64) float64 {
	if len(sizes) == 0 {
		return 0
	}
	total := 0.0
	for i, size := range sizes {
		total += size
		if i < len(sizes)-1 {
			total += gap
		}
	}
	return total
}
