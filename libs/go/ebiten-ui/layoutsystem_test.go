package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestPageLayoutBuildsFixedHeaderAndScrollableBody(t *testing.T) {
	dom := ebitenui.New(ebitenui.PageLayout(ebitenui.PageLayoutConfig{
		ID:       "page-root",
		ScrollID: "page-scroll",
		Padding:  12,
		Gap:      8,
		Header: ebitenui.Header(ebitenui.Props{
			ID: "page-header",
			Style: ebitenui.Style{
				Direction: ebitenui.Row,
				Padding:   ebitenui.All(10),
			},
		}, ebitenui.Text("Header", ebitenui.Props{ID: "page-header-label"})),
		Content: ebitenui.Main(ebitenui.Props{
			ID: "page-body",
			Style: ebitenui.Style{
				Direction: ebitenui.Column,
				Gap:       8,
			},
		},
			ebitenui.Section(ebitenui.Props{ID: "section-a", Style: ebitenui.Style{Height: ebitenui.Px(120)}}),
			ebitenui.Section(ebitenui.Props{ID: "section-b", Style: ebitenui.Style{Height: ebitenui.Px(120)}}),
		),
		Footer: ebitenui.Footer(ebitenui.Props{
			ID: "page-footer",
			Style: ebitenui.Style{
				Padding: ebitenui.All(10),
			},
		}, ebitenui.Text("Footer", ebitenui.Props{ID: "page-footer-label"})),
	}))

	layout := dom.Layout(ebitenui.Viewport{Width: 320, Height: 200})
	scroll, ok := layout.FindByID("page-scroll")
	if !ok || scroll.Node.Tag != ebitenui.TagScrollView {
		t.Fatalf("expected page scroll view")
	}

	header, ok := layout.FindByID("page-header")
	if !ok {
		t.Fatalf("expected page header")
	}
	footer, ok := layout.FindByID("page-footer")
	if !ok {
		t.Fatalf("expected page footer")
	}
	if footer.Frame.Y <= header.Frame.Y {
		t.Fatalf("expected footer after header, got header=%#v footer=%#v", header.Frame, footer.Frame)
	}
}

func TestPageLayoutDispatchesScrollOffsetChanges(t *testing.T) {
	var offsets []float64
	dom := ebitenui.New(ebitenui.PageLayout(ebitenui.PageLayoutConfig{
		ID:       "page-root",
		ScrollID: "page-scroll",
		Padding:  12,
		Gap:      8,
		Content: ebitenui.Main(ebitenui.Props{
			ID: "page-body",
			Style: ebitenui.Style{
				Direction: ebitenui.Column,
				Gap:       8,
			},
		},
			ebitenui.Section(ebitenui.Props{ID: "section-a", Style: ebitenui.Style{Height: ebitenui.Px(180)}}),
			ebitenui.Section(ebitenui.Props{ID: "section-b", Style: ebitenui.Style{Height: ebitenui.Px(180)}}),
			ebitenui.Section(ebitenui.Props{ID: "section-c", Style: ebitenui.Style{Height: ebitenui.Px(180)}}),
		),
		OnScrollChange: func(offset float64) {
			offsets = append(offsets, offset)
		},
	}))

	runtime := ebitenui.NewRuntime()
	viewport := ebitenui.Viewport{Width: 320, Height: 200}
	layout := runtime.Update(dom, viewport, ebitenui.InputSnapshot{})
	scroll, ok := layout.FindByID("page-scroll")
	if !ok {
		t.Fatalf("expected page scroll")
	}

	x := scroll.Frame.X + 10
	y := scroll.Frame.Y + 10
	runtime.Update(dom, viewport, ebitenui.InputSnapshot{PointerX: x, PointerY: y, ScrollY: -1})

	if len(offsets) == 0 {
		t.Fatalf("expected scroll callback")
	}
	if offsets[0] <= 0 {
		t.Fatalf("expected positive scroll offset, got %v", offsets[0])
	}
}
