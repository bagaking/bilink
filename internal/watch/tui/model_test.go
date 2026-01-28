package tui

import "testing"

func TestModelInit(t *testing.T) {
	m := NewModel()
	if m.Status == "" {
		t.Fatalf("expected status")
	}
}
