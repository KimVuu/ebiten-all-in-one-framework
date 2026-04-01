package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestColumnLayoutStacksChildren(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "app",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(220),
			Height:    ebitenui.Px(160),
			Direction: ebitenui.Column,
			Padding:   ebitenui.All(12),
			Gap:       10,
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "header",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Px(30),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "content",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Px(40),
			},
		}),
	)

	dom := ebitenui.New(root)
	layout := dom.Layout(ebitenui.Viewport{Width: 220, Height: 160})

	header, ok := layout.FindByID("header")
	if !ok {
		t.Fatalf("expected header node in layout tree")
	}
	if got, want := header.Frame, (ebitenui.Rect{X: 12, Y: 12, Width: 196, Height: 30}); got != want {
		t.Fatalf("header frame mismatch: got %#v want %#v", got, want)
	}

	content, ok := layout.FindByID("content")
	if !ok {
		t.Fatalf("expected content node in layout tree")
	}
	if got, want := content.Frame, (ebitenui.Rect{X: 12, Y: 52, Width: 196, Height: 40}); got != want {
		t.Fatalf("content frame mismatch: got %#v want %#v", got, want)
	}
}

func TestRowLayoutDistributesFillWidth(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "row",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(300),
			Height:    ebitenui.Px(80),
			Direction: ebitenui.Row,
			Padding:   ebitenui.All(10),
			Gap:       10,
		},
	},
		ebitenui.Div(ebitenui.Props{
			ID: "fixed",
			Style: ebitenui.Style{
				Width:  ebitenui.Px(40),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "grow-a",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		}),
		ebitenui.Div(ebitenui.Props{
			ID: "grow-b",
			Style: ebitenui.Style{
				Width:  ebitenui.Fill(),
				Height: ebitenui.Fill(),
			},
		}),
	)

	dom := ebitenui.New(root)
	layout := dom.Layout(ebitenui.Viewport{Width: 300, Height: 80})

	growA, ok := layout.FindByID("grow-a")
	if !ok {
		t.Fatalf("expected grow-a node in layout tree")
	}
	if got, want := growA.Frame, (ebitenui.Rect{X: 60, Y: 10, Width: 110, Height: 60}); got != want {
		t.Fatalf("grow-a frame mismatch: got %#v want %#v", got, want)
	}

	growB, ok := layout.FindByID("grow-b")
	if !ok {
		t.Fatalf("expected grow-b node in layout tree")
	}
	if got, want := growB.Frame, (ebitenui.Rect{X: 180, Y: 10, Width: 110, Height: 60}); got != want {
		t.Fatalf("grow-b frame mismatch: got %#v want %#v", got, want)
	}
}

func TestTextLayoutUsesIntrinsicSize(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "app",
		Style: ebitenui.Style{
			Width:   ebitenui.Px(180),
			Height:  ebitenui.Px(90),
			Padding: ebitenui.All(8),
		},
	},
		ebitenui.Text("Hello, Ebiten Ebiten UI", ebitenui.Props{ID: "label"}),
	)

	dom := ebitenui.New(root)
	layout := dom.Layout(ebitenui.Viewport{Width: 180, Height: 90})

	label, ok := layout.FindByID("label")
	if !ok {
		t.Fatalf("expected label node in layout tree")
	}
	if label.Frame.Width <= 0 {
		t.Fatalf("expected intrinsic text width, got %v", label.Frame.Width)
	}
	if label.Frame.Height <= 0 {
		t.Fatalf("expected intrinsic text height, got %v", label.Frame.Height)
	}
}

func TestFindByIDReturnsDOMNode(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{ID: "app"},
		ebitenui.Section(ebitenui.Props{ID: "profile"},
			ebitenui.Text("Kim", ebitenui.Props{ID: "display-name"}),
		),
	)

	dom := ebitenui.New(root)
	node, ok := dom.FindByID("display-name")
	if !ok {
		t.Fatalf("expected node lookup to succeed")
	}
	if got, want := node.Tag, ebitenui.TagText; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if got, want := node.Text, "Kim"; got != want {
		t.Fatalf("text mismatch: got %q want %q", got, want)
	}
}

func TestRowIntrinsicLayoutDoesNotOverflowWhenChildUsesFill(t *testing.T) {
	root := ebitenui.Div(ebitenui.Props{
		ID: "app",
		Style: ebitenui.Style{
			Width:     ebitenui.Fill(),
			Height:    ebitenui.Fill(),
			Direction: ebitenui.Column,
			Padding:   ebitenui.All(10),
		},
	},
		ebitenui.Header(ebitenui.Props{
			ID: "header",
			Style: ebitenui.Style{
				Direction: ebitenui.Row,
				Padding:   ebitenui.All(10),
				Gap:       8,
			},
		},
			ebitenui.Div(ebitenui.Props{
				ID: "copy",
				Style: ebitenui.Style{
					Width:     ebitenui.Fill(),
					Direction: ebitenui.Column,
				},
			},
				ebitenui.Text("Title", ebitenui.Props{ID: "title"}),
			),
			ebitenui.Span(ebitenui.Props{
				ID: "badge",
				Style: ebitenui.Style{
					Padding: ebitenui.All(6),
				},
			},
				ebitenui.Text("Badge", ebitenui.Props{ID: "badge-text"}),
			),
		),
	)

	layout := ebitenui.New(root).Layout(ebitenui.Viewport{Width: 200, Height: 100})
	header, ok := layout.FindByID("header")
	if !ok {
		t.Fatalf("expected header layout")
	}
	if header.Frame.X+header.Frame.Width > 190 {
		t.Fatalf("expected header to fit content width, got %#v", header.Frame)
	}
}
