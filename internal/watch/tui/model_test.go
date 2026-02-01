package tui

import (
	"strings"
	"testing"

	"github.com/bagaking/bilink/internal/output"
)

func TestModelViewCyberpunk(t *testing.T) {
	m := NewModel(output.WatchPayload{Added: []string{"a.md"}, Removed: []string{"b.md"}}, "config")
	view := m.View()
	if !strings.Contains(view, "BILINK WATCH") || !strings.Contains(view, "ASK") {
		t.Fatalf("expected title and ask")
	}
	if !strings.Contains(view, "+ a.md") || !strings.Contains(view, "- b.md") {
		t.Fatalf("expected change list")
	}
}
