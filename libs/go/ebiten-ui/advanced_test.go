package ebitenui_test

import (
	"image/color"
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestImageUsesIntrinsicSize(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "app",
		Style: ebitenui.Style{
			Width:   ebitenui.Px(100),
			Height:  ebitenui.Px(80),
			Padding: ebitenui.All(8),
		},
	},
		ebitenui.Image(ebitenui.Props{
			ID:    "avatar",
			Image: ebitenui.SolidImage(24, 18, color.RGBA{R: 80, G: 120, B: 180, A: 255}),
		}),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 100, Height: 80})
	avatar, ok := layout.FindByID("avatar")
	if !ok {
		t.Fatalf("expected avatar layout node")
	}
	if got, want := avatar.Frame, (ebitenui.Rect{X: 8, Y: 8, Width: 24, Height: 18}); got != want {
		t.Fatalf("avatar frame mismatch: got %#v want %#v", got, want)
	}
}

func TestTextBlockWrapsWithinConfiguredWidth(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "app",
		Style: ebitenui.Style{
			Width:   ebitenui.Px(120),
			Height:  ebitenui.Px(120),
			Padding: ebitenui.All(8),
		},
	},
		ebitenui.TextBlock("A long text block that should wrap across multiple lines.", ebitenui.Props{
			ID: "copy",
			Style: ebitenui.Style{
				Width: ebitenui.Px(72),
			},
		}),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 120, Height: 120})
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
	root := ebitenui.Div(ebitenui.Props{
		ID: "row",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(220),
			Height:    ebitenui.Px(50),
			Direction: ebitenui.Row,
			Padding:   ebitenui.All(10),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "left",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(20),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Spacer(ebitenui.Props{
			ID: "spacer",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "right",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(30),
				Height: ebitenui.Fill(),
			},
		}),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 220, Height: 50})
	spacer, ok := layout.FindByID("spacer")
	if !ok {
		t.Fatalf("expected spacer layout node")
	}
	if got, want := spacer.Frame, (ebitenui.Rect{X: 30, Y: 10, Width: 150, Height: 30}); got != want {
		t.Fatalf("spacer frame mismatch: got %#v want %#v", got, want)
	}
}

func TestStackOverlaysChildrenInsideContentBox(t *testing.T) {
	root := ebitenui.Stack(ebitenui.Props{
		ID: "stack",
		Style: ebitenui.Style{
			Width:   ebitenui.Px(120),
			Height:  ebitenui.Px(80),
			Padding: ebitenui.All(4),
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "background",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "badge",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(20),
				Height: ebitenui.Px(12),
			},
		}),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 120, Height: 80})
	background, ok := layout.FindByID("background")
	if !ok {
		t.Fatalf("expected background node")
	}
	if got, want := background.Frame, (ebitenui.Rect{X: 4, Y: 4, Width: 112, Height: 72}); got != want {
		t.Fatalf("background frame mismatch: got %#v want %#v", got, want)
	}

	badge, ok := layout.FindByID("badge")
	if !ok {
		t.Fatalf("expected badge node")
	}
	if got, want := badge.Frame, (ebitenui.Rect{X: 4, Y: 4, Width: 20, Height: 12}); got != want {
		t.Fatalf("badge frame mismatch: got %#v want %#v", got, want)
	}
}

func TestScrollViewOffsetsChildrenAndTracksContentSize(t *testing.T) {
	root := ebitenui.ScrollView(ebitenui.Props{
		ID: "scroll",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(120),
			Height:    ebitenui.Px(60),
			Direction: ebitenui.Column,
			Padding:   ebitenui.All(4),
			Gap:       4,
		},
		Scroll: ebitenui.ScrollState{
			OffsetY: 18,
		},
	},
		ebitenui.Div(ebitenui.Props{ID: "item-a", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
		ebitenui.Div(ebitenui.Props{ID: "item-b", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
		ebitenui.Div(ebitenui.Props{ID: "item-c", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 120, Height: 60})
	itemA, ok := layout.FindByID("item-a")
	if !ok {
		t.Fatalf("expected first scroll child")
	}
	if got, want := itemA.Frame, (ebitenui.Rect{X: 4, Y: -14, Width: 112, Height: 20}); got != want {
		t.Fatalf("scroll offset mismatch: got %#v want %#v", got, want)
	}
	if got, want := layout.ContentHeight, float64(68); got != want {
		t.Fatalf("content height mismatch: got %v want %v", got, want)
	}
}

func TestInteractiveButtonCarriesInteractionState(t *testing.T) {
	node := ebitenui.InteractiveButton(ebitenui.Props{
		ID: "cta",
		State: ebitenui.InteractionState{
			Hovered:  true,
			Focused:  true,
			Selected: true,
		},
	},
		ebitenui.Text("Continue", ebitenui.Props{ID: "cta-label"}),
	)

	if got, want := node.Tag, ebitenui.TagButton; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if !node.Props.State.Hovered || !node.Props.State.Focused || !node.Props.State.Selected {
		t.Fatalf("expected interaction state to be preserved, got %#v", node.Props.State)
	}
}
