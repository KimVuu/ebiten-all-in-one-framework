package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
	"golang.org/x/image/font/basicfont"
)

func TestSetTextFaceOverridesAndResetsDefaultFace(t *testing.T) {
	original := ebitenui.TextFace()
	defer ebitenui.SetTextFace(original)

	ebitenui.SetTextFace(basicfont.Face7x13)
	if got := ebitenui.TextFace(); got == nil {
		t.Fatalf("expected text face")
	}

	ebitenui.SetTextFace(nil)
	if got := ebitenui.TextFace(); got == nil {
		t.Fatalf("expected fallback text face after reset")
	}
}
