package tui

import (
	"strings"
	"testing"

	"github.com/bagaking/bilink/internal/output"
)

func TestModelViewIncludesAsk(t *testing.T) {
	m := NewModel(output.WatchPayload{Added: []string{"a.md"}}, "config")
	view := m.View()
	if !strings.Contains(view, "ASK") {
		t.Fatalf("expected ask prompt")
	}
}
