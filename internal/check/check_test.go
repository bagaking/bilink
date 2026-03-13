package check

import "testing"

func TestDetectGroups(t *testing.T) {
	files := []File{
		{Path: "a.md", ResolveKey: "foo", LintKey: "foo"},
		{Path: "b.md", ResolveKey: "foo", LintKey: "bar"},
		{Path: "c.md", ResolveKey: "baz", LintKey: "bar"},
	}
	errs, warns := Detect(files)
	if len(errs) != 1 || errs[0].Key != "foo" || len(errs[0].Paths) != 2 {
		t.Fatalf("unexpected errors")
	}
	if len(warns) != 1 || warns[0].Key != "bar" || len(warns[0].Paths) != 2 {
		t.Fatalf("unexpected warnings")
	}
}

func TestDetectGroupsSortedDeterministically(t *testing.T) {
	files := []File{
		{Path: "z.md", ResolveKey: "zeta", LintKey: "zeta"},
		{Path: "b.md", ResolveKey: "alpha", LintKey: "alpha"},
		{Path: "a.md", ResolveKey: "alpha", LintKey: "alpha"},
		{Path: "y.md", ResolveKey: "zeta", LintKey: "zeta"},
	}

	errs, warns := Detect(files)
	assertGroups(t, errs, []Group{
		{Key: "alpha", Paths: []string{"a.md", "b.md"}},
		{Key: "zeta", Paths: []string{"y.md", "z.md"}},
	})
	assertGroups(t, warns, []Group{
		{Key: "alpha", Paths: []string{"a.md", "b.md"}},
		{Key: "zeta", Paths: []string{"y.md", "z.md"}},
	})
}

func assertGroups(t *testing.T, got []Group, want []Group) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d groups, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i].Key != want[i].Key {
			t.Fatalf("group %d key: got %q, want %q", i, got[i].Key, want[i].Key)
		}
		if len(got[i].Paths) != len(want[i].Paths) {
			t.Fatalf("group %d paths length: got %d, want %d", i, len(got[i].Paths), len(want[i].Paths))
		}
		for j := range want[i].Paths {
			if got[i].Paths[j] != want[i].Paths[j] {
				t.Fatalf("group %d path %d: got %q, want %q", i, j, got[i].Paths[j], want[i].Paths[j])
			}
		}
	}
}
