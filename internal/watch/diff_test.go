package watch

import "testing"

func TestDiff(t *testing.T) {
	oldIdx := Index{Files: []string{"a.md"}}
	newIdx := Index{Files: []string{"a.md", "b.md"}}
	diff := Diff(oldIdx, newIdx)
	if len(diff.Added) != 1 || diff.Added[0] != "b.md" {
		t.Fatalf("unexpected diff")
	}
}
