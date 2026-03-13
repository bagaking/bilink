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

func TestDiffSortsAddedAndRemoved(t *testing.T) {
	oldIdx := Index{Files: []string{"z.md", "b.md", "stable.md"}}
	newIdx := Index{Files: []string{"stable.md", "y.md", "a.md"}}

	diff := Diff(oldIdx, newIdx)
	assertStrings(t, diff.Added, []string{"a.md", "y.md"})
	assertStrings(t, diff.Removed, []string{"b.md", "z.md"})
}

func assertStrings(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d strings, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("string %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
