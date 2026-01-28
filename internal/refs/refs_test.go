package refs

import "testing"

func TestRefsFor(t *testing.T) {
	idx := Index{
		Outbound: map[string][]Link{"a.md": {{Target: "b"}}},
		Inbound:  map[string][]Link{"a.md": {{Target: "a"}}},
	}
	out, in := RefsFor(idx, "a.md")
	if len(out) != 1 || len(in) != 1 {
		t.Fatalf("unexpected refs")
	}
}
