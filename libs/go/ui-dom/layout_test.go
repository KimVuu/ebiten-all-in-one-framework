package uidom_test

import (
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestColumnLayoutStacksChildren(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "app",
		Style: uidom.Style{
			Width:     uidom.Px(220),
			Height:    uidom.Px(160),
			Direction: uidom.Column,
			Padding:   uidom.All(12),
			Gap:       10,
		},
	},
		uidom.Div(uidom.Props{
			ID: "header",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Px(30),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "content",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Px(40),
			},
		}),
	)

	dom := uidom.New(root)
	layout := dom.Layout(uidom.Viewport{Width: 220, Height: 160})

	header, ok := layout.FindByID("header")
	if !ok {
		t.Fatalf("expected header node in layout tree")
	}
	if got, want := header.Frame, (uidom.Rect{X: 12, Y: 12, Width: 196, Height: 30}); got != want {
		t.Fatalf("header frame mismatch: got %#v want %#v", got, want)
	}

	content, ok := layout.FindByID("content")
	if !ok {
		t.Fatalf("expected content node in layout tree")
	}
	if got, want := content.Frame, (uidom.Rect{X: 12, Y: 52, Width: 196, Height: 40}); got != want {
		t.Fatalf("content frame mismatch: got %#v want %#v", got, want)
	}
}

func TestRowLayoutDistributesFillWidth(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "row",
		Style: uidom.Style{
			Width:     uidom.Px(300),
			Height:    uidom.Px(80),
			Direction: uidom.Row,
			Padding:   uidom.All(10),
			Gap:       10,
		},
	},
		uidom.Div(uidom.Props{
			ID: "fixed",
			Style: uidom.Style{
				Width:  uidom.Px(40),
				Height: uidom.Fill(),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "grow-a",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Fill(),
			},
		}),
		uidom.Div(uidom.Props{
			ID: "grow-b",
			Style: uidom.Style{
				Width:  uidom.Fill(),
				Height: uidom.Fill(),
			},
		}),
	)

	dom := uidom.New(root)
	layout := dom.Layout(uidom.Viewport{Width: 300, Height: 80})

	growA, ok := layout.FindByID("grow-a")
	if !ok {
		t.Fatalf("expected grow-a node in layout tree")
	}
	if got, want := growA.Frame, (uidom.Rect{X: 60, Y: 10, Width: 110, Height: 60}); got != want {
		t.Fatalf("grow-a frame mismatch: got %#v want %#v", got, want)
	}

	growB, ok := layout.FindByID("grow-b")
	if !ok {
		t.Fatalf("expected grow-b node in layout tree")
	}
	if got, want := growB.Frame, (uidom.Rect{X: 180, Y: 10, Width: 110, Height: 60}); got != want {
		t.Fatalf("grow-b frame mismatch: got %#v want %#v", got, want)
	}
}

func TestTextLayoutUsesIntrinsicSize(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "app",
		Style: uidom.Style{
			Width:   uidom.Px(180),
			Height:  uidom.Px(90),
			Padding: uidom.All(8),
		},
	},
		uidom.Text("Hello, Ebiten UI DOM", uidom.Props{ID: "label"}),
	)

	dom := uidom.New(root)
	layout := dom.Layout(uidom.Viewport{Width: 180, Height: 90})

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
	root := uidom.Div(uidom.Props{ID: "app"},
		uidom.Section(uidom.Props{ID: "profile"},
			uidom.Text("Kim", uidom.Props{ID: "display-name"}),
		),
	)

	dom := uidom.New(root)
	node, ok := dom.FindByID("display-name")
	if !ok {
		t.Fatalf("expected node lookup to succeed")
	}
	if got, want := node.Tag, uidom.TagText; got != want {
		t.Fatalf("tag mismatch: got %q want %q", got, want)
	}
	if got, want := node.Text, "Kim"; got != want {
		t.Fatalf("text mismatch: got %q want %q", got, want)
	}
}

func TestRowIntrinsicLayoutDoesNotOverflowWhenChildUsesFill(t *testing.T) {
	root := uidom.Div(uidom.Props{
		ID: "app",
		Style: uidom.Style{
			Width:     uidom.Fill(),
			Height:    uidom.Fill(),
			Direction: uidom.Column,
			Padding:   uidom.All(10),
		},
	},
		uidom.Header(uidom.Props{
			ID: "header",
			Style: uidom.Style{
				Direction: uidom.Row,
				Padding:   uidom.All(10),
				Gap:       8,
			},
		},
			uidom.Div(uidom.Props{
				ID: "copy",
				Style: uidom.Style{
					Width:     uidom.Fill(),
					Direction: uidom.Column,
				},
			},
				uidom.Text("Title", uidom.Props{ID: "title"}),
			),
			uidom.Span(uidom.Props{
				ID: "badge",
				Style: uidom.Style{
					Padding: uidom.All(6),
				},
			},
				uidom.Text("Badge", uidom.Props{ID: "badge-text"}),
			),
		),
	)

	layout := uidom.New(root).Layout(uidom.Viewport{Width: 200, Height: 100})
	header, ok := layout.FindByID("header")
	if !ok {
		t.Fatalf("expected header layout")
	}
	if header.Frame.X+header.Frame.Width > 190 {
		t.Fatalf("expected header to fit content width, got %#v", header.Frame)
	}
}
