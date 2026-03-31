package uidom_test

import (
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestLayoutSpecFallsBackToStyleFlowLayout(t *testing.T) {
	styleDOM := uidom.New(uidom.Div(uidom.Props{
		ID: "grid",
		Style: uidom.Style{
			Width:     uidom.Px(240),
			Height:    uidom.Px(120),
			Direction: uidom.Column,
			Padding:   uidom.All(12),
			Gap:       8,
		},
	},
		uidom.Div(uidom.Props{ID: "header", Style: uidom.Style{Height: uidom.Px(24)}}),
		uidom.Div(uidom.Props{ID: "body", Style: uidom.Style{Height: uidom.Px(40)}}),
	))

	layoutDOM := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Layout: uidom.LayoutSpec{
			Mode:    uidom.LayoutModeFlowVertical,
			Padding: uidom.All(12),
			Gap:     8,
			Size: uidom.LayoutSize{
				Width:  uidom.Px(240),
				Height: uidom.Px(120),
			},
		},
	},
		uidom.Div(uidom.Props{ID: "header", Style: uidom.Style{Height: uidom.Px(24)}}),
		uidom.Div(uidom.Props{ID: "body", Style: uidom.Style{Height: uidom.Px(40)}}),
	))

	styleLayout := styleDOM.Layout(uidom.Viewport{Width: 240, Height: 120})
	layoutLayout := layoutDOM.Layout(uidom.Viewport{Width: 240, Height: 120})

	for _, id := range []string{"header", "body"} {
		styleNode, ok := styleLayout.FindByID(id)
		if !ok {
			t.Fatalf("expected style node %q", id)
		}
		layoutNode, ok := layoutLayout.FindByID(id)
		if !ok {
			t.Fatalf("expected layout node %q", id)
		}
		if got, want := layoutNode.Frame, styleNode.Frame; got != want {
			t.Fatalf("frame mismatch for %q: got %#v want %#v", id, got, want)
		}
	}
}

func TestAnchoredLayoutComputesRelativeToParent(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Style: uidom.Style{
			Width:  uidom.Px(400),
			Height: uidom.Px(240),
		},
	},
		uidom.Div(uidom.Props{
			ID: "panel",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorCenter,
				Pivot:  uidom.PivotCenter,
				Size: uidom.LayoutSize{
					Width:  uidom.Px(120),
					Height: uidom.Px(40),
				},
			},
		}),
	))

	layout := dom.Layout(uidom.Viewport{Width: 400, Height: 240})
	if got, want := layout.ContentBounds, layout.Frame; got != want {
		t.Fatalf("expected root content bounds to match frame, got %#v want %#v", got, want)
	}
	if layout.Overflow.Any {
		t.Fatalf("expected root overflow to be false")
	}
	panel, ok := layout.FindByID("panel")
	if !ok {
		t.Fatalf("expected panel layout")
	}
	if got, want := panel.Frame, (uidom.Rect{X: 140, Y: 100, Width: 120, Height: 40}); got != want {
		t.Fatalf("anchored frame mismatch: got %#v want %#v", got, want)
	}
	if panel.ParentID != "root" {
		t.Fatalf("expected parent id to be root, got %q", panel.ParentID)
	}
	if panel.ClipRect != panel.Frame {
		t.Fatalf("expected clip rect to match frame, got %#v", panel.ClipRect)
	}
	if panel.ClickableRect != panel.Frame {
		t.Fatalf("expected clickable rect to match frame, got %#v", panel.ClickableRect)
	}
}

func TestValidatorReportsOverflowOverlapAndPatches(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Style: uidom.Style{
			Width:  uidom.Px(200),
			Height: uidom.Px(120),
		},
	},
		uidom.Div(uidom.Props{
			ID: "overflow",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorTopLeft,
				Offset: uidom.Point{X: 170, Y: 12},
				Size: uidom.LayoutSize{
					Width:  uidom.Px(60),
					Height: uidom.Px(24),
				},
				Constraints: uidom.LayoutConstraints{
					KeepInsideParent: true,
				},
			},
		}),
		uidom.Div(uidom.Props{
			ID: "overlap",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorTopLeft,
				Offset: uidom.Point{X: 168, Y: 10},
				Size: uidom.LayoutSize{
					Width:  uidom.Px(60),
					Height: uidom.Px(24),
				},
			},
		}),
		uidom.Div(uidom.Props{
			ID: "tiny",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorTopLeft,
				Offset: uidom.Point{X: 12, Y: 72},
				Size: uidom.LayoutSize{
					Width:  uidom.Px(12),
					Height: uidom.Px(12),
				},
				Constraints: uidom.LayoutConstraints{
					MinHitTarget: 44,
				},
			},
		}),
	))

	layout := dom.Layout(uidom.Viewport{Width: 200, Height: 120})
	report := uidom.ValidateLayout(layout, uidom.Viewport{Width: 200, Height: 120}, uidom.ValidationOptions{})

	if len(report.Issues) == 0 {
		t.Fatalf("expected validation issues")
	}

	foundOverflow := false
	foundOverlap := false
	foundHitTarget := false
	for _, issue := range report.Issues {
		switch issue.Code {
		case uidom.IssueOutOfViewport:
			foundOverflow = true
		case uidom.IssueOverlap:
			foundOverlap = true
		case uidom.IssueMinHitTarget:
			foundHitTarget = true
		}
		if len(issue.SuggestedConstraintChanges) == 0 {
			t.Fatalf("expected suggested constraint change for issue %#v", issue)
		}
	}

	if !foundOverflow {
		t.Fatalf("expected out of viewport issue")
	}
	if !foundOverlap {
		t.Fatalf("expected overlap issue")
	}
	if !foundHitTarget {
		t.Fatalf("expected min hit target issue")
	}
}

func TestNestedAnchoredLayoutKeepsChildInsideAnchoredParent(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Style: uidom.Style{
			Width:  uidom.Px(640),
			Height: uidom.Px(360),
		},
	},
		uidom.Div(uidom.Props{
			ID: "panel",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorCenter,
				Pivot:  uidom.PivotCenter,
				Size: uidom.LayoutSize{
					Width:  uidom.Px(280),
					Height: uidom.Px(180),
				},
				Padding: uidom.All(20),
			},
		},
			uidom.Div(uidom.Props{
				ID: "close-button",
				Layout: uidom.LayoutSpec{
					Mode:   uidom.LayoutModeAnchored,
					Anchor: uidom.AnchorTopRight,
					Pivot:  uidom.PivotTopRight,
					Offset: uidom.Point{X: -8, Y: 8},
					Size: uidom.LayoutSize{
						Width:  uidom.Px(32),
						Height: uidom.Px(32),
					},
					Constraints: uidom.LayoutConstraints{
						KeepInsideParent: true,
					},
				},
			}),
		),
	))

	layout := dom.Layout(uidom.Viewport{Width: 640, Height: 360})
	panel, ok := layout.FindByID("panel")
	if !ok {
		t.Fatalf("expected panel")
	}
	button, ok := layout.FindByID("close-button")
	if !ok {
		t.Fatalf("expected close button")
	}
	if button.Frame.X < panel.ContentBounds.X || button.Frame.Y < panel.ContentBounds.Y {
		t.Fatalf("expected child to stay inside parent content bounds, panel=%#v child=%#v", panel.ContentBounds, button.Frame)
	}
	if button.Frame.X+button.Frame.Width > panel.ContentBounds.X+panel.ContentBounds.Width {
		t.Fatalf("expected child right edge inside parent content bounds, panel=%#v child=%#v", panel.ContentBounds, button.Frame)
	}
}

func TestGridLayoutAutoMeasuresRowsFromNestedContent(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Style: uidom.Style{
			Width:  uidom.Px(480),
			Height: uidom.Px(320),
		},
	},
		uidom.Div(uidom.Props{
			ID: "grid",
			Layout: uidom.LayoutSpec{
				Mode: uidom.LayoutModeGrid,
				Grid: uidom.LayoutGrid{
					Columns: 2,
					Gap:     12,
				},
				Size: uidom.LayoutSize{
					Width: uidom.Fill(),
				},
			},
			Style: uidom.Style{
				Padding: uidom.All(12),
			},
		},
			uidom.Div(uidom.Props{
				ID: "card-a",
				Style: uidom.Style{
					Direction: uidom.Column,
					Gap:       8,
					Padding:   uidom.All(8),
				},
			},
				uidom.TextBlock("A card with wrapped text content that should increase row height.", uidom.Props{
					ID: "card-a-copy",
					Style: uidom.Style{
						Width: uidom.Fill(),
					},
				}),
			),
			uidom.Div(uidom.Props{
				ID: "card-b",
				Style: uidom.Style{
					Height: uidom.Px(40),
				},
			}),
			uidom.Div(uidom.Props{
				ID: "card-c",
				Layout: uidom.LayoutSpec{
					Mode:   uidom.LayoutModeAnchored,
					Anchor: uidom.AnchorCenter,
					Pivot:  uidom.PivotCenter,
					Size: uidom.LayoutSize{
						Width:  uidom.Px(80),
						Height: uidom.Px(28),
					},
				},
				Style: uidom.Style{
					Height: uidom.Px(80),
				},
			}),
		),
	))

	layout := dom.Layout(uidom.Viewport{Width: 480, Height: 320})
	grid, ok := layout.FindByID("grid")
	if !ok {
		t.Fatalf("expected grid")
	}
	cardA, ok := layout.FindByID("card-a")
	if !ok {
		t.Fatalf("expected card-a")
	}
	cardC, ok := layout.FindByID("card-c")
	if !ok {
		t.Fatalf("expected card-c")
	}
	if grid.ContentHeight <= cardA.Frame.Height {
		t.Fatalf("expected grid content height to account for multiple rows, grid=%#v cardA=%#v", grid.ContentHeight, cardA.Frame)
	}
	if cardC.Frame.Y <= cardA.Frame.Y {
		t.Fatalf("expected second row item below first row item, cardA=%#v cardC=%#v", cardA.Frame, cardC.Frame)
	}
}

func TestGridLayoutSupportsSpanAndAutoPlacement(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "root",
		Style: uidom.Style{
			Width:  uidom.Px(360),
			Height: uidom.Px(240),
		},
	},
		uidom.Div(uidom.Props{
			ID: "grid",
			Layout: uidom.LayoutSpec{
				Mode: uidom.LayoutModeGrid,
				Grid: uidom.LayoutGrid{
					Columns: 3,
					Gap:     12,
				},
				Size: uidom.LayoutSize{
					Width: uidom.Fill(),
				},
				Padding: uidom.All(12),
			},
		},
			uidom.Div(uidom.Props{
				ID: "hero",
				Layout: uidom.LayoutSpec{
					Grid: uidom.LayoutGrid{
						ColumnStart: 1,
						ColumnSpan:  2,
					},
				},
				Style: uidom.Style{
					Width:  uidom.Px(220),
					Height: uidom.Px(48),
				},
			}),
			uidom.Div(uidom.Props{
				ID: "aside",
				Style: uidom.Style{
					Height: uidom.Px(48),
				},
			}),
			uidom.Div(uidom.Props{
				ID: "content",
				Style: uidom.Style{
					Height: uidom.Px(40),
				},
			}),
		),
	))

	layout := dom.Layout(uidom.Viewport{Width: 360, Height: 240})
	hero, ok := layout.FindByID("hero")
	if !ok {
		t.Fatalf("expected hero")
	}
	aside, ok := layout.FindByID("aside")
	if !ok {
		t.Fatalf("expected aside")
	}
	content, ok := layout.FindByID("content")
	if !ok {
		t.Fatalf("expected content")
	}

	if hero.Frame.Width < 220 {
		t.Fatalf("expected spanned item to stay wider than a single track, got %#v", hero.Frame)
	}
	if aside.Frame.X <= hero.Frame.X+hero.Frame.Width-1 {
		t.Fatalf("expected aside to auto-place after hero span, hero=%#v aside=%#v", hero.Frame, aside.Frame)
	}
	if content.Frame.Y <= hero.Frame.Y {
		t.Fatalf("expected content to wrap to next row after first row is filled, hero=%#v content=%#v", hero.Frame, content.Frame)
	}
}

func TestGridLayoutAppliesJustifyAndAlignWithinCell(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "grid",
		Style: uidom.Style{
			Width:  uidom.Px(280),
			Height: uidom.Px(180),
		},
	},
		uidom.Div(uidom.Props{
			ID: "grid",
			Layout: uidom.LayoutSpec{
				Mode: uidom.LayoutModeGrid,
				Grid: uidom.LayoutGrid{
					Columns:      2,
					Gap:          8,
					JustifyItems: uidom.LayoutAlignmentCenter,
					AlignItems:   uidom.LayoutAlignmentCenter,
				},
				Size: uidom.LayoutSize{
					Width: uidom.Fill(),
				},
				Padding: uidom.All(12),
			},
		},
			uidom.Div(uidom.Props{
				ID: "chip",
				Style: uidom.Style{
					Width:  uidom.Px(40),
					Height: uidom.Px(20),
				},
			}),
			uidom.Div(uidom.Props{
				ID: "tall",
				Style: uidom.Style{
					Width:  uidom.Px(32),
					Height: uidom.Px(60),
				},
			}),
		),
	))

	layout := dom.Layout(uidom.Viewport{Width: 280, Height: 180})
	grid, ok := layout.FindByID("grid")
	if !ok {
		t.Fatalf("expected grid")
	}
	chip, ok := layout.FindByID("chip")
	if !ok {
		t.Fatalf("expected chip")
	}

	if chip.Frame.X <= grid.ContentBounds.X {
		t.Fatalf("expected centered x inside the cell, got %#v", chip.Frame)
	}
	if chip.Frame.Y <= grid.ContentBounds.Y {
		t.Fatalf("expected centered y inside the row, got %#v", chip.Frame)
	}
}

func TestGridLayoutSupportsSpanAutoPlacementAndNestedAnchors(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "grid",
		Layout: uidom.LayoutSpec{
			Mode: uidom.LayoutModeGrid,
			Grid: uidom.LayoutGrid{
				Columns: 3,
				Gap:     10,
			},
			Padding: uidom.All(10),
			Size: uidom.LayoutSize{
				Width:  uidom.Px(360),
				Height: uidom.Px(220),
			},
		},
	},
		uidom.Div(uidom.Props{
			ID: "hero",
			Layout: uidom.LayoutSpec{
				Grid: uidom.LayoutGrid{
					ColumnStart: 1,
					ColumnSpan:  2,
					RowStart:    1,
					RowSpan:     1,
					JustifySelf: uidom.LayoutAlignmentCenter,
					AlignSelf:   uidom.LayoutAlignmentEnd,
				},
			},
			Style: uidom.Style{
				Width:  uidom.Px(240),
				Height: uidom.Px(20),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "auto-a",
			Style: uidom.Style{
				Width:  uidom.Px(50),
				Height: uidom.Px(30),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "auto-b",
			Style: uidom.Style{
				Width:  uidom.Px(50),
				Height: uidom.Px(30),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "auto-c",
			Layout: uidom.LayoutSpec{
				Mode:   uidom.LayoutModeAnchored,
				Anchor: uidom.AnchorBottomRight,
				Pivot:  uidom.PivotBottomRight,
				Offset: uidom.Point{X: -4, Y: -4},
				Size: uidom.LayoutSize{
					Width:  uidom.Px(24),
					Height: uidom.Px(16),
				},
				Constraints: uidom.LayoutConstraints{
					KeepInsideParent: true,
				},
			},
			Style: uidom.Style{
				Width:  uidom.Px(40),
				Height: uidom.Px(28),
			},
		}),
	))

	layout := dom.Layout(uidom.Viewport{Width: 360, Height: 220})
	hero, ok := layout.FindByID("hero")
	if !ok {
		t.Fatalf("expected hero")
	}
	autoA, ok := layout.FindByID("auto-a")
	if !ok {
		t.Fatalf("expected auto-a")
	}
	autoB, ok := layout.FindByID("auto-b")
	if !ok {
		t.Fatalf("expected auto-b")
	}
	autoC, ok := layout.FindByID("auto-c")
	if !ok {
		t.Fatalf("expected auto-c")
	}

	if hero.Frame.X >= autoA.Frame.X {
		t.Fatalf("expected auto-a to resolve after the hero span, hero=%#v autoA=%#v", hero.Frame, autoA.Frame)
	}
	if hero.Frame.Y+hero.Frame.Height != autoA.Frame.Y+autoA.Frame.Height {
		t.Fatalf("expected first-row items to share a row baseline, hero=%#v autoA=%#v", hero.Frame, autoA.Frame)
	}
	if autoB.Frame.Y <= hero.Frame.Y {
		t.Fatalf("expected auto-b to auto-place on the next row, hero=%#v autoB=%#v", hero.Frame, autoB.Frame)
	}
	if autoC.Frame.X < layout.ContentBounds.X || autoC.Frame.Y < layout.ContentBounds.Y {
		t.Fatalf("expected anchored child to stay within grid content bounds, autoC=%#v content=%#v", autoC.Frame, layout.ContentBounds)
	}
	if autoC.Frame.X+autoC.Frame.Width > layout.ContentBounds.X+layout.ContentBounds.Width || autoC.Frame.Y+autoC.Frame.Height > layout.ContentBounds.Y+layout.ContentBounds.Height {
		t.Fatalf("expected anchored child to stay within grid content bounds, autoC=%#v content=%#v", autoC.Frame, layout.ContentBounds)
	}
}

func TestGridLayoutAppliesContentAndItemAlignment(t *testing.T) {
	dom := uidom.New(uidom.Div(uidom.Props{
		ID: "grid",
		Layout: uidom.LayoutSpec{
			Mode: uidom.LayoutModeGrid,
			Grid: uidom.LayoutGrid{
				Columns:        2,
				Gap:            10,
				JustifyContent: uidom.LayoutAlignmentCenter,
				AlignContent:   uidom.LayoutAlignmentEnd,
				JustifyItems:   uidom.LayoutAlignmentCenter,
				AlignItems:     uidom.LayoutAlignmentEnd,
			},
			Padding: uidom.All(10),
			Size: uidom.LayoutSize{
				Width:  uidom.Px(320),
				Height: uidom.Px(220),
			},
		},
	},
		uidom.Div(uidom.Props{
			ID: "left",
			Style: uidom.Style{
				Width:  uidom.Px(40),
				Height: uidom.Px(20),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "right",
			Style: uidom.Style{
				Width:  uidom.Px(50),
				Height: uidom.Px(30),
			},
		}),
	))

	layout := dom.Layout(uidom.Viewport{Width: 320, Height: 220})
	grid, ok := layout.FindByID("grid")
	if !ok {
		t.Fatalf("expected grid")
	}
	left, ok := layout.FindByID("left")
	if !ok {
		t.Fatalf("expected left item")
	}
	right, ok := layout.FindByID("right")
	if !ok {
		t.Fatalf("expected right item")
	}

	cellWidth := (grid.ContentBounds.Width - 10) / 2
	expectedX := grid.ContentBounds.X + (cellWidth-left.Frame.Width)/2
	if left.Frame.X != expectedX {
		t.Fatalf("expected centered grid start x=%.0f, got %#v", expectedX, left.Frame)
	}
	if left.Frame.Y != 190 {
		t.Fatalf("expected bottom-aligned left item y=190, got %#v", left.Frame)
	}
	if right.Frame.X <= left.Frame.X {
		t.Fatalf("expected right item to follow left item, left=%#v right=%#v", left.Frame, right.Frame)
	}
	if right.Frame.Y != 180 {
		t.Fatalf("expected taller item to align to row bottom, got %#v", right.Frame)
	}
}
