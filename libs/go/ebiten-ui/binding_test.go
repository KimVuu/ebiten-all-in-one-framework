package ebitenui_test

import (
	"testing"

	ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"
)

func TestRefGetAndSet(t *testing.T) {
	ref := ebitenui.NewRef("Kim")
	if got, want := ref.Get(), "Kim"; got != want {
		t.Fatalf("ref get mismatch: got %q want %q", got, want)
	}

	ref.Set("Yechan")
	if got, want := ref.Get(), "Yechan"; got != want {
		t.Fatalf("ref set mismatch: got %q want %q", got, want)
	}
}

func TestComputedReflectsLatestSourceValue(t *testing.T) {
	score := ebitenui.NewRef(12)
	double := ebitenui.NewComputed(func() int {
		return score.Get() * 2
	})

	if got, want := double.Get(), 24; got != want {
		t.Fatalf("computed mismatch: got %d want %d", got, want)
	}

	score.Set(18)
	if got, want := double.Get(), 36; got != want {
		t.Fatalf("computed updated mismatch: got %d want %d", got, want)
	}
}
