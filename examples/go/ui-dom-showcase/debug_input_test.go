package main

import (
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestDebugInputQueueSchedulesClickAcrossFrames(t *testing.T) {
	queue := newDebugInputQueue()
	queue.queueClick(12, debugResolvedTarget{
		ID: "start-button",
		Frame: uidom.Rect{
			X:      40,
			Y:      60,
			Width:  120,
			Height: 48,
		},
	})

	move := queue.drain(13)
	if len(move) != 1 {
		t.Fatalf("expected one move event, got %d", len(move))
	}
	if got, want := move[0].resolvedTarget, "start-button"; got != want {
		t.Fatalf("resolved target mismatch: got %q want %q", got, want)
	}
	if move[0].input.PointerX != 100 || move[0].input.PointerY != 84 {
		t.Fatalf("unexpected pointer move: %#v", move[0].input)
	}

	down := queue.drain(14)
	if len(down) != 1 || !down[0].input.PointerDown {
		t.Fatalf("expected pointer down event, got %#v", down)
	}

	up := queue.drain(15)
	if len(up) != 1 || up[0].input.PointerDown {
		t.Fatalf("expected pointer up event, got %#v", up)
	}
}

func TestResolveDebugTargetFindsInteractiveAncestor(t *testing.T) {
	dom := uidom.New(
		uidom.Div(uidom.Props{
			ID: "root",
			Style: uidom.Style{
				Width:  uidom.Px(240),
				Height: uidom.Px(120),
			},
		},
			uidom.InteractiveButton(uidom.Props{
				ID: "play-button",
				Style: uidom.Style{
					Width:  uidom.Px(160),
					Height: uidom.Px(48),
				},
			},
				uidom.Text("Play", uidom.Props{ID: "play-label"}),
			),
		),
	)

	layout := dom.Layout(uidom.Viewport{Width: 240, Height: 120})
	target, ok := resolveDebugTarget(layout, "play-label")
	if !ok {
		t.Fatalf("expected interactive target")
	}
	if got, want := target.ID, "play-button"; got != want {
		t.Fatalf("target mismatch: got %q want %q", got, want)
	}
	if target.Frame.Width <= 0 || target.Frame.Height <= 0 {
		t.Fatalf("expected positive frame, got %#v", target.Frame)
	}
}

func TestDebugLayoutReportFlagsInvalidGeometry(t *testing.T) {
	dom := uidom.New(
		uidom.Div(uidom.Props{
			ID: "root",
			Style: uidom.Style{
				Width:  uidom.Px(100),
				Height: uidom.Px(100),
			},
		},
			uidom.Div(uidom.Props{
				ID: "overflow-child",
				Style: uidom.Style{
					Width:  uidom.Px(180),
					Height: uidom.Px(60),
				},
			}),
		),
	)

	layout := dom.Layout(uidom.Viewport{Width: 100, Height: 100})
	report := buildDebugLayoutReport(layout, uidom.Viewport{Width: 100, Height: 100})
	if report.InvalidNodeCount == 0 {
		t.Fatalf("expected invalid node count")
	}
	if !hasIssueCode(report.Issues, "overflow-child", "out_of_viewport") &&
		!hasIssueCode(report.Issues, "overflow-child", "out_of_parent") {
		t.Fatalf("expected overflow issue, got %#v", report.Issues)
	}
}
