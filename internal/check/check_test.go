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
