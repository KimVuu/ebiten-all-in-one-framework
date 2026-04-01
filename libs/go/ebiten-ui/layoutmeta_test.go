package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestLayoutSpecFallsBackToStyleFlowLayout(t *testing.T) {
	styleDOM := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "grid",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(240),
			Height:    ebitenui.Px(120),
			Direction: ebitenui.Column,
			Padding:   ebitenui.All(12),
			Gap:       8,
		},
	},
		ebitenui.Div(ebitenui.Props{ID: "header", Style: ebitenui.Style{Height: ebitenui.Px(24)}}),
		ebitenui.Div(ebitenui.Props{ID: "body", Style: ebitenui.Style{Height: ebitenui.Px(40)}}),
	))

	layoutDOM := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Layout: ebitenui.LayoutSpec{
			Mode:    ebitenui.LayoutModeFlowVertical,
			Padding: ebitenui.All(12),
			Gap:     8,
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Px(240),
				Height: ebitenui.Px(120),
			},
		},
	},
		ebitenui.Div(ebitenui.Props{ID: "header", Style: ebitenui.Style{Height: ebitenui.Px(24)}}),
		ebitenui.Div(ebitenui.Props{ID: "body", Style: ebitenui.Style{Height: ebitenui.Px(40)}}),
	))

	styleLayout := styleDOM.Layout(ebitenui.Viewport{Width: 240, Height: 120})
	layoutLayout := layoutDOM.Layout(ebitenui.Viewport{Width: 240, Height: 120})

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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(400),
			Height: ebitenui.Px(240),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "panel",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorCenter,
				Pivot:  ebitenui.PivotCenter,
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(120),
					Height: ebitenui.Px(40),
				},
			},
		}),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 400, Height: 240})
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
	if got, want := panel.Frame, (ebitenui.Rect{X: 140, Y: 100, Width: 120, Height: 40}); got != want {
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(200),
			Height: ebitenui.Px(120),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "overflow",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorTopLeft,
				Offset: ebitenui.Point{X: 170, Y: 12},
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(60),
					Height: ebitenui.Px(24),
				},
				Constraints: ebitenui.LayoutConstraints{
					KeepInsideParent: true,
				},
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "overlap",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorTopLeft,
				Offset: ebitenui.Point{X: 168, Y: 10},
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(60),
					Height: ebitenui.Px(24),
				},
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "tiny",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorTopLeft,
				Offset: ebitenui.Point{X: 12, Y: 72},
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(12),
					Height: ebitenui.Px(12),
				},
				Constraints: ebitenui.LayoutConstraints{
					MinHitTarget: 44,
				},
			},
		}),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 200, Height: 120})
	report := ebitenui.ValidateLayout(layout, ebitenui.Viewport{Width: 200, Height: 120}, ebitenui.ValidationOptions{})

	if len(report.Issues) == 0 {
		t.Fatalf("expected validation issues")
	}

	foundOverflow := false
	foundOverlap := false
	foundHitTarget := false
	for _, issue := range report.Issues {
		switch issue.Code {
		case ebitenui.IssueOutOfViewport:
			foundOverflow = true
		case ebitenui.IssueOverlap:
			foundOverlap = true
		case ebitenui.IssueMinHitTarget:
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(640),
			Height: ebitenui.Px(360),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "panel",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorCenter,
				Pivot:  ebitenui.PivotCenter,
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(280),
					Height: ebitenui.Px(180),
				},
				Padding: ebitenui.All(20),
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "close-button",
				Layout: ebitenui.LayoutSpec{
					Mode:   ebitenui.LayoutModeAnchored,
					Anchor: ebitenui.AnchorTopRight,
					Pivot:  ebitenui.PivotTopRight,
					Offset: ebitenui.Point{X: -8, Y: 8},
					Size: ebitenui.LayoutSize{
						Width:  ebitenui.Px(32),
						Height: ebitenui.Px(32),
					},
					Constraints: ebitenui.LayoutConstraints{
						KeepInsideParent: true,
					},
				},
			}),
		),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 640, Height: 360})
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(480),
			Height: ebitenui.Px(320),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "grid",
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeGrid,
				Grid: ebitenui.LayoutGrid{
					Columns: 2,
					Gap:     12,
				},
				Size: ebitenui.LayoutSize{
					Width: ebitenui.Fill(),
				},
			},
			Style: ebitenui.Style{
				Padding: ebitenui.All(12),
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "card-a",
				Style: ebitenui.Style{
					Direction: ebitenui.Column,
					Gap:       8,
					Padding:   ebitenui.All(8),
				},
			},
				ebitenui.TextBlock("A card with wrapped text content that should increase row height.", ebitenui.Props{
					ID: "card-a-copy",
					Style: ebitenui.Style{
						Width: ebitenui.Fill(),
					},
				}),
			),
			ebitenui.Div(ebitenui.Props{
				ID: "card-b",
				Style: ebitenui.Style{
					Height: ebitenui.Px(40),
				},
			}),
			ebitenui.Div(ebitenui.Props{
				ID: "card-c",
				Layout: ebitenui.LayoutSpec{
					Mode:   ebitenui.LayoutModeAnchored,
					Anchor: ebitenui.AnchorCenter,
					Pivot:  ebitenui.PivotCenter,
					Size: ebitenui.LayoutSize{
						Width:  ebitenui.Px(80),
						Height: ebitenui.Px(28),
					},
				},
				Style: ebitenui.Style{
					Height: ebitenui.Px(80),
				},
			}),
		),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 480, Height: 320})
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "root",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(360),
			Height: ebitenui.Px(240),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "grid",
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeGrid,
				Grid: ebitenui.LayoutGrid{
					Columns: 3,
					Gap:     12,
				},
				Size: ebitenui.LayoutSize{
					Width: ebitenui.Fill(),
				},
				Padding: ebitenui.All(12),
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "hero",
				Layout: ebitenui.LayoutSpec{
					Grid: ebitenui.LayoutGrid{
						ColumnStart: 1,
						ColumnSpan:  2,
					},
				},
				Style: ebitenui.Style{
					Width:  ebitenui.Px(220),
					Height: ebitenui.Px(48),
				},
			}),
			ebitenui.Div(ebitenui.Props{
				ID: "aside",
				Style: ebitenui.Style{
					Height: ebitenui.Px(48),
				},
			}),
			ebitenui.Div(ebitenui.Props{
				ID: "content",
				Style: ebitenui.Style{
					Height: ebitenui.Px(40),
				},
			}),
		),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 360, Height: 240})
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "grid",
		Style: ebitenui.Style{
			Width:  ebitenui.Px(280),
			Height: ebitenui.Px(180),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "grid",
			Layout: ebitenui.LayoutSpec{
				Mode: ebitenui.LayoutModeGrid,
				Grid: ebitenui.LayoutGrid{
					Columns:      2,
					Gap:          8,
					JustifyItems: ebitenui.LayoutAlignmentCenter,
					AlignItems:   ebitenui.LayoutAlignmentCenter,
				},
				Size: ebitenui.LayoutSize{
					Width: ebitenui.Fill(),
				},
				Padding: ebitenui.All(12),
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "chip",
				Style: ebitenui.Style{
					Width:  ebitenui.Px(40),
					Height: ebitenui.Px(20),
				},
			}),
			ebitenui.Div(ebitenui.Props{
				ID: "tall",
				Style: ebitenui.Style{
					Width:  ebitenui.Px(32),
					Height: ebitenui.Px(60),
				},
			}),
		),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 280, Height: 180})
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "grid",
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeGrid,
			Grid: ebitenui.LayoutGrid{
				Columns: 3,
				Gap:     10,
			},
			Padding: ebitenui.All(10),
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Px(360),
				Height: ebitenui.Px(220),
			},
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "hero",
			Layout: ebitenui.LayoutSpec{
				Grid: ebitenui.LayoutGrid{
					ColumnStart: 1,
					ColumnSpan:  2,
					RowStart:    1,
					RowSpan:     1,
					JustifySelf: ebitenui.LayoutAlignmentCenter,
					AlignSelf:   ebitenui.LayoutAlignmentEnd,
				},
			},
			Style: ebitenui.Style{
				Width:  ebitenui.Px(240),
				Height: ebitenui.Px(20),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "auto-a",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(50),
				Height: ebitenui.Px(30),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "auto-b",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(50),
				Height: ebitenui.Px(30),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "auto-c",
			Layout: ebitenui.LayoutSpec{
				Mode:   ebitenui.LayoutModeAnchored,
				Anchor: ebitenui.AnchorBottomRight,
				Pivot:  ebitenui.PivotBottomRight,
				Offset: ebitenui.Point{X: -4, Y: -4},
				Size: ebitenui.LayoutSize{
					Width:  ebitenui.Px(24),
					Height: ebitenui.Px(16),
				},
				Constraints: ebitenui.LayoutConstraints{
					KeepInsideParent: true,
				},
			},
			Style: ebitenui.Style{
				Width:  ebitenui.Px(40),
				Height: ebitenui.Px(28),
			},
		}),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 360, Height: 220})
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
	dom := ebitenui.New(ebitenui.Div(ebitenui.Props{
		ID: "grid",
		Layout: ebitenui.LayoutSpec{
			Mode: ebitenui.LayoutModeGrid,
			Grid: ebitenui.LayoutGrid{
				Columns:        2,
				Gap:            10,
				JustifyContent: ebitenui.LayoutAlignmentCenter,
				AlignContent:   ebitenui.LayoutAlignmentEnd,
				JustifyItems:   ebitenui.LayoutAlignmentCenter,
				AlignItems:     ebitenui.LayoutAlignmentEnd,
			},
			Padding: ebitenui.All(10),
			Size: ebitenui.LayoutSize{
				Width:  ebitenui.Px(320),
				Height: ebitenui.Px(220),
			},
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "left",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(40),
				Height: ebitenui.Px(20),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "right",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(50),
				Height: ebitenui.Px(30),
			},
		}),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 320, Height: 220})
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
