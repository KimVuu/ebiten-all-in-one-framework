package ebitenuidebug

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ebitendebug "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestInspectCompactUINodeHonorsChildDepth(t *testing.T) {
	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:     ebitenui.Px(240),
				Height:    ebitenui.Px(240),
				Direction: ebitenui.Column,
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "parent",
				Style: ebitenui.Style{
					Width:     ebitenui.Px(200),
					Height:    ebitenui.Px(200),
					Direction: ebitenui.Column,
				},
			},
				ebitenui.Div(ebitenui.Props{
					ID: "child",
					Style: ebitenui.Style{
						Width:     ebitenui.Px(180),
						Height:    ebitenui.Px(100),
						Direction: ebitenui.Column,
					},
				},
					ebitenui.Text("leaf", ebitenui.Props{ID: "leaf"}),
				),
			),
		),
	)

	viewport := ebitenui.Viewport{Width: 240, Height: 240}
	layout := dom.Layout(viewport)
	report := buildDebugLayoutReport(layout, viewport)

	depthOne, ok := inspectCompactUINode(layout, viewport, report, ebitendebug.UINodeInspectRequest{
		NodeID:          "parent",
		IncludeChildren: true,
		ChildDepth:      1,
	})
	if !ok {
		t.Fatalf("expected node detail")
	}
	if got, want := len(depthOne.Children), 1; got != want {
		t.Fatalf("expected direct children only, got %d", got)
	}
	if got, want := depthOne.Children[0].ID, "child"; got != want {
		t.Fatalf("unexpected child at depth 1: got %q want %q", got, want)
	}

	depthTwo, ok := inspectCompactUINode(layout, viewport, report, ebitendebug.UINodeInspectRequest{
		NodeID:          "parent",
		IncludeChildren: true,
		ChildDepth:      2,
	})
	if !ok {
		t.Fatalf("expected node detail")
	}
	if got, want := len(depthTwo.Children), 2; got != want {
		t.Fatalf("expected descendant summaries up to depth 2, got %d", got)
	}
	if got, want := depthTwo.Children[1].ID, "leaf"; got != want {
		t.Fatalf("unexpected descendant at depth 2: got %q want %q", got, want)
	}
}

func TestCompactVisibilityRespectsClipChain(t *testing.T) {
	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(240),
				Height: ebitenui.Px(200),
			},
			Layout: ebitenui.LayoutSpec{Mode: ebitenui.LayoutModeStack},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "clip-parent",
				Style: ebitenui.Style{
					Width:  ebitenui.Px(120),
					Height: ebitenui.Px(60),
				},
				Layout: ebitenui.LayoutSpec{
					Mode: ebitenui.LayoutModeStack,
					Constraints: ebitenui.LayoutConstraints{
						ClipChildren: true,
					},
				},
			},
				ebitenui.Div(ebitenui.Props{
					ID: "visible-child",
					Layout: ebitenui.LayoutSpec{
						Mode:   ebitenui.LayoutModeAnchored,
						Anchor: ebitenui.AnchorTopLeft,
						Pivot:  ebitenui.PivotTopLeft,
						Offset: ebitenui.Point{X: 0, Y: 10},
						Size:   ebitenui.LayoutSize{Width: ebitenui.Px(40), Height: ebitenui.Px(20)},
					},
				}),
				ebitenui.Div(ebitenui.Props{
					ID: "hidden-child",
					Layout: ebitenui.LayoutSpec{
						Mode:   ebitenui.LayoutModeAnchored,
						Anchor: ebitenui.AnchorTopLeft,
						Pivot:  ebitenui.PivotTopLeft,
						Offset: ebitenui.Point{X: 0, Y: 80},
						Size:   ebitenui.LayoutSize{Width: ebitenui.Px(40), Height: ebitenui.Px(20)},
					},
				}),
			),
		),
	)

	viewport := ebitenui.Viewport{Width: 240, Height: 200}
	layout := dom.Layout(viewport)
	report := buildDebugLayoutReport(layout, viewport)

	overview := buildCompactUIOverview(layout, viewport, report, nil, ebitenui.InputSnapshot{}, 0)
	if got, want := overview.VisibleNodeCount, 3; got != want {
		t.Fatalf("visible node count mismatch: got %d want %d", got, want)
	}

	query := queryCompactUINodes(layout, viewport, report, ebitendebug.UIQueryRequest{
		VisibleOnly: true,
		Limit:       10,
	})
	for _, node := range query.Nodes {
		if node.ID == "hidden-child" {
			t.Fatalf("expected clipped child to be excluded from visible query")
		}
	}

	snapshot := buildDebugUISnapshot(layout, viewport, report, false, nil, ebitenui.InputSnapshot{}, 0)
	if len(snapshot.Root.Children) == 0 || len(snapshot.Root.Children[0].Children) < 2 {
		t.Fatalf("unexpected snapshot structure")
	}
	hidden := snapshot.Root.Children[0].Children[1]
	if hidden.Visible {
		t.Fatalf("expected hidden child snapshot to be invisible")
	}
	if hidden.Computed == nil || hidden.Computed.Visible {
		t.Fatalf("expected hidden child computed visibility to be false")
	}
}

func TestCaptureCompactUIScreenshotDefaultsToWorkingDirectoryScreenshots(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWD)
	}()

	dom := ebitenui.New(
		ebitenui.Div(ebitenui.Props{
			ID: "root",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(40),
				Height: ebitenui.Px(40),
			},
		}),
	)
	viewport := ebitenui.Viewport{Width: 40, Height: 40}
	layout := dom.Layout(viewport)
	report := buildDebugLayoutReport(layout, viewport)

	result, _, ok := captureCompactUIScreenshot("test-game", "", layout, viewport, report, ebitendebug.UICaptureRequest{Target: "viewport"})
	if !ok {
		t.Fatalf("expected capture result")
	}
	wantSuffix := filepath.Join("screenshots", "test-game")
	if got := filepath.Dir(result.Path); !strings.HasSuffix(got, wantSuffix) {
		t.Fatalf("unexpected screenshot dir: got %q want suffix %q", got, wantSuffix)
	}
}
