package index

import (
	"encoding/json"
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

func TestBuildIndexStableAcrossInputOrder(t *testing.T) {
	files := []FileInput{
		{Path: "z.md", Content: "See [[same]]"},
		{Path: "a.md", Content: "See [[same]]"},
		{Path: "dir/same.md", Content: "Hi"},
		{Path: "same.md", Content: "Hi"},
	}
	reversed := []FileInput{files[3], files[2], files[1], files[0]}

	forward := mustJSON(t, Build(files))
	backward := mustJSON(t, Build(reversed))
	if forward != backward {
		t.Fatalf("expected stable index output across input order\nforward:  %s\nbackward: %s", forward, backward)
	}
}

func TestBuildIndexDuplicateBasenameUsesLexicographicPathWinner(t *testing.T) {
	files := []FileInput{
		{Path: "roots/z/same.md", Content: "Hi"},
		{Path: "roots/a/same.md", Content: "Hi"},
		{Path: "source.md", Content: "See [[same]]"},
	}

	idx := Build(files)
	if got := len(idx.Inbound["roots/a/same.md"]); got != 1 {
		t.Fatalf("expected lexicographic path winner to receive inbound link, got %d", got)
	}
	if got := len(idx.Inbound["roots/z/same.md"]); got != 0 {
		t.Fatalf("expected later lexicographic duplicate basename to receive no inbound links, got %d", got)
	}
}

func mustJSON(t *testing.T, idx Index) string {
	t.Helper()
	data, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("marshal index: %v", err)
	}
	return string(data)
}
