package ui

import "testing"

func TestLoadTextFace(t *testing.T) {
	face, err := LoadTextFace()
	if err != nil {
		t.Fatalf("LoadTextFace() failed: %v", err)
	}
	if face == nil {
		t.Fatalf("expected font face")
	}
	if got := face.Metrics().Height.Ceil(); got <= 0 {
		t.Fatalf("expected positive font height, got %d", got)
	}
}
