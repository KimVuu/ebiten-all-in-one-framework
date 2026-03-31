package uidom_test

import (
	"image/color"
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestImageUsesIntrinsicSize(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "app",
		Style: uidom.Style{
			Width:   uidom.Px(100),
			Height:  uidom.Px(80),
			Padding: uidom.All(8),
		},
	},
		uidom.Image(uidom.Props{
			ID:    "avatar",
			Image: uidom.SolidImage(24, 18, color.RGBA{R: 80, G: 120, B: 180, A: 255}),
		}),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 100, Height: 80})
	avatar, ok := layout.FindByID("avatar")
	if !ok {
		t.Fatalf("expected avatar layout node")
	}
	if got, want := avatar.Frame, (uidom.Rect{X: 8, Y: 8, Width: 24, Height: 18}); got != want {
		t.Fatalf("avatar frame mismatch: got %#v want %#v", got, want)
	}
}

func TestTextBlockWrapsWithinConfiguredWidth(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "app",
		Style: uidom.Style{
			Width:   uidom.Px(120),
			Height:  uidom.Px(120),
			Padding: uidom.All(8),
		},
	},
		uidom.TextBlock("A long text block that should wrap across multiple lines.", uidom.Props{
			ID: "copy",
			Style: uidom.Style{
				Width: uidom.Px(72),
			},
		}),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 120, Height: 120})
	copyNode, ok := layout.FindByID("copy")
	if !ok {
		t.Fatalf("expected text block layout node")
	}
	if copyNode.Frame.Width != 72 {
		t.Fatalf("expected wrapped text width 72, got %v", copyNode.Frame.Width)
	}
	if copyNode.Frame.Height <= 13 {
		t.Fatalf("expected wrapped text height to span multiple lines, got %v", copyNode.Frame.Height)
	}
}

func TestSpacerConsumesRemainingSpace(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "row",
		Style: uidom.Style{
			Width:     uidom.Px(220),
			Height:    uidom.Px(50),
			Direction: uidom.Row,
			Padding:   uidom.All(10),
		},
	},
		uidom.Div(uidom.Props{
			ID: "left",
			Style: uidom.Style{
				Width:  uidom.Px(20),
				Height: uidom.Fill(),
			},
		}),
		uidom.Spacer(uidom.Props{
			ID: "spacer",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Fill(),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "right",
			Style: uidom.Style{
				Width:  uidom.Px(30),
				Height: uidom.Fill(),
			},
		}),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 220, Height: 50})
	spacer, ok := layout.FindByID("spacer")
	if !ok {
		t.Fatalf("expected spacer layout node")
	}
	if got, want := spacer.Frame, (uidom.Rect{X: 30, Y: 10, Width: 150, Height: 30}); got != want {
		t.Fatalf("spacer frame mismatch: got %#v want %#v", got, want)
	}
}

func TestStackOverlaysChildrenInsideContentBox(t *testing.T) {
	root := uidom.Stack(uidom.Props{
		ID: "stack",
		Style: uidom.Style{
			Width:   uidom.Px(120),
			Height:  uidom.Px(80),
			Padding: uidom.All(4),
		},
	},
		uidom.Div(uidom.Props{
			ID: "background",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Fill(),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "badge",
			Style: uidom.Style{
				Width:  uidom.Px(20),
				Height: uidom.Px(12),
			},
		}),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 120, Height: 80})
	background, ok := layout.FindByID("background")
	if !ok {
		t.Fatalf("expected background node")
	}
	if got, want := background.Frame, (uidom.Rect{X: 4, Y: 4, Width: 112, Height: 72}); got != want {
		t.Fatalf("background frame mismatch: got %#v want %#v", got, want)
	}

	badge, ok := layout.FindByID("badge")
	if !ok {
		t.Fatalf("expected badge node")
	}
	if got, want := badge.Frame, (uidom.Rect{X: 4, Y: 4, Width: 20, Height: 12}); got != want {
		t.Fatalf("badge frame mismatch: got %#v want %#v", got, want)
	}
}

func TestScrollViewOffsetsChildrenAndTracksContentSize(t *testing.T) {
	root := uidom.ScrollView(uidom.Props{
		ID: "scroll",
		Style: uidom.Style{
			Width:     uidom.Px(120),
			Height:    uidom.Px(60),
			Direction: uidom.Column,
			Padding:   uidom.All(4),
			Gap:       4,
		},
		Scroll: uidom.ScrollState{
			OffsetY: 18,
		},
	},
		uidom.Div(uidom.Props{ID: "item-a", Style: uidom.Style{Width: uidom.Fill(), Height: uidom.Px(20)}}),
		uidom.Div(uidom.Props{ID: "item-b", Style: uidom.Style{Width: uidom.Fill(), Height: uidom.Px(20)}}),
		uidom.Div(uidom.Props{ID: "item-c", Style: uidom.Style{Width: uidom.Fill(), Height: uidom.Px(20)}}),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 120, Height: 60})
	itemA, ok := layout.FindByID("item-a")
	if !ok {
		t.Fatalf("expected first scroll child")
	}
	if got, want := itemA.Frame, (uidom.Rect{X: 4, Y: -14, Width: 112, Height: 20}); got != want {
		t.Fatalf("scroll offset mismatch: got %#v want %#v", got, want)
	}
	if got, want := layout.ContentHeight, float64(68); got != want {
		t.Fatalf("content height mismatch: got %v want %v", got, want)
	}
}

func TestInteractiveButtonCarriesInteractionState(t *testing.T) {
	node := uidom.InteractiveButton(uidom.Props{
		ID: "cta",
		State: uidom.InteractionState{
			Hovered:  true,
			Focused:  true,
			Selected: true,
		},
	},
		uidom.Text("Continue", uidom.Props{ID: "cta-label"}),
	)

	if got, want := node.Tag, uidom.TagButton; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if !node.Props.State.Hovered || !node.Props.State.Focused || !node.Props.State.Selected {
		t.Fatalf("expected interaction state to be preserved, got %#v", node.Props.State)
	}
}
