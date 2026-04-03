package ui

import "testing"

func TestNewRendererLoadsEmbeddedFont(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}
	if renderer == nil {
		t.Fatalf("expected renderer")
	}
	if renderer.face == nil {
		t.Fatalf("expected font face")
	}
	if got := renderer.face.Metrics().Height.Ceil(); got <= 0 {
		t.Fatalf("expected positive font height, got %d", got)
	}
}
