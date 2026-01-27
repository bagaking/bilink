package index

import (
	"testing"

	"github.com/bagaking/bilink/internal/parse"
)

func TestBuildIndex(t *testing.T) {
	files := []FileInput{
		{Path: "a.md", Content: "See [[b]]"},
		{Path: "b.md", Content: "Hi"},
	}
	idx := Build(files)
	if len(idx.Outbound["a.md"]) != 1 {
		t.Fatalf("expected outbound link")
	}
	if len(idx.Inbound["b.md"]) != 1 {
		t.Fatalf("expected inbound link")
	}
	if idx.Outbound["a.md"][0].Target != "b" {
		t.Fatalf("unexpected link %#v", idx.Outbound["a.md"][0])
	}
	_ = parse.Link{}
}
