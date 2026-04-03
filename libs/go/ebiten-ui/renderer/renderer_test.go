package renderer

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestIntersectImageRectReturnsSharedArea(t *testing.T) {
	got, ok := intersectImageRect(image.Rect(0, 0, 20, 20), image.Rect(10, 8, 30, 16))
	if !ok {
		t.Fatalf("expected intersecting rects")
	}
	if want := image.Rect(10, 8, 20, 16); got != want {
		t.Fatalf("intersection mismatch: got %v want %v", got, want)
	}
}

func TestClippedImageUsesRequestedBounds(t *testing.T) {
	screen := ebiten.NewImage(40, 40)
	clip := image.Rect(10, 12, 26, 30)
	got := clippedImage(screen, clip)
	if got == nil {
		t.Fatalf("expected clipped image")
	}
	if want := clip; got.Bounds() != want {
		t.Fatalf("clipped image bounds mismatch: got %v want %v", got.Bounds(), want)
	}
}

func TestScrollViewLayoutExposesContentClipForRenderer(t *testing.T) {
	dom := ebitenui.New(ebitenui.ScrollView(ebitenui.Props{
		ID: "scroll",
		Style: ebitenui.Style{
			Width:     ebitenui.Px(120),
			Height:    ebitenui.Px(60),
			Direction: ebitenui.Column,
			Padding:   ebitenui.All(4),
			Gap:       4,
		},
		Scroll: ebitenui.ScrollState{OffsetY: 18},
	},
		ebitenui.Div(ebitenui.Props{ID: "item-a", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
		ebitenui.Div(ebitenui.Props{ID: "item-b", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
		ebitenui.Div(ebitenui.Props{ID: "item-c", Style: ebitenui.Style{Width: ebitenui.Fill(), Height: ebitenui.Px(20)}}),
	))

	layout := dom.Layout(ebitenui.Viewport{Width: 120, Height: 60})
	scroll, ok := layout.FindByID("scroll")
	if !ok {
		t.Fatalf("expected scroll layout")
	}
	if !scroll.ClipChildren {
		t.Fatalf("expected scroll view to clip children")
	}
	if want := (ebitenui.Rect{X: 4, Y: 4, Width: 112, Height: 52}); scroll.ClipRect != want {
		t.Fatalf("clip rect mismatch: got %#v want %#v", scroll.ClipRect, want)
	}

	childClip, ok := intersectImageRect(imageRect(layout.Frame), imageRect(scroll.ClipRect))
	if !ok {
		t.Fatalf("expected child clip intersection")
	}
	if want := image.Rect(4, 4, 116, 56); childClip != want {
		t.Fatalf("child clip mismatch: got %v want %v", childClip, want)
	}

	itemA, ok := scroll.FindByID("item-a")
	if !ok {
		t.Fatalf("expected first item")
	}
	if !imageRectsIntersect(childClip, imageRect(itemA.Frame)) {
		t.Fatalf("expected first child to still partially intersect clip rect")
	}
}
