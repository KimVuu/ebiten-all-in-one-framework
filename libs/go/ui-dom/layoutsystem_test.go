package uidom_test

import (
	"testing"

	uidom "github.com/kimyechan/ebiten-aio-framework/libs/go/ui-dom"
)

func TestPageLayoutBuildsFixedHeaderAndScrollableBody(t *testing.T) {
	dom := uidom.New(uidom.PageLayout(uidom.PageLayoutConfig{
		ID:       "page-root",
		ScrollID: "page-scroll",
		Padding:  12,
		Gap:      8,
		Header: uidom.Header(uidom.Props{
			ID: "page-header",
			Style: uidom.Style{
				Direction: uidom.Row,
				Padding:   uidom.All(10),
			},
		}, uidom.Text("Header", uidom.Props{ID: "page-header-label"})),
		Content: uidom.Main(uidom.Props{
			ID: "page-body",
			Style: uidom.Style{
				Direction: uidom.Column,
				Gap:       8,
			},
		},
			uidom.Section(uidom.Props{ID: "section-a", Style: uidom.Style{Height: uidom.Px(120)}}),
			uidom.Section(uidom.Props{ID: "section-b", Style: uidom.Style{Height: uidom.Px(120)}}),
		),
		Footer: uidom.Footer(uidom.Props{
			ID: "page-footer",
			Style: uidom.Style{
				Padding: uidom.All(10),
			},
		}, uidom.Text("Footer", uidom.Props{ID: "page-footer-label"})),
	}))

	layout := dom.Layout(uidom.Viewport{Width: 320, Height: 200})
	scroll, ok := layout.FindByID("page-scroll")
	if !ok || scroll.Node.Tag != uidom.TagScrollView {
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
	dom := uidom.New(uidom.PageLayout(uidom.PageLayoutConfig{
		ID:       "page-root",
		ScrollID: "page-scroll",
		Padding:  12,
		Gap:      8,
		Content: uidom.Main(uidom.Props{
			ID: "page-body",
			Style: uidom.Style{
				Direction: uidom.Column,
				Gap:       8,
			},
		},
			uidom.Section(uidom.Props{ID: "section-a", Style: uidom.Style{Height: uidom.Px(180)}}),
			uidom.Section(uidom.Props{ID: "section-b", Style: uidom.Style{Height: uidom.Px(180)}}),
			uidom.Section(uidom.Props{ID: "section-c", Style: uidom.Style{Height: uidom.Px(180)}}),
		),
		OnScrollChange: func(offset float64) {
			offsets = append(offsets, offset)
		},
	}))

	runtime := uidom.NewRuntime()
	viewport := uidom.Viewport{Width: 320, Height: 200}
	layout := runtime.Update(dom, viewport, uidom.InputSnapshot{})
	scroll, ok := layout.FindByID("page-scroll")
	if !ok {
		t.Fatalf("expected page scroll")
	}

	x := scroll.Frame.X + 10
	y := scroll.Frame.Y + 10
	runtime.Update(dom, viewport, uidom.InputSnapshot{PointerX: x, PointerY: y, ScrollY: -1})

	if len(offsets) == 0 {
		t.Fatalf("expected scroll callback")
	}
	if offsets[0] <= 0 {
		t.Fatalf("expected positive scroll offset, got %v", offsets[0])
	}
}
