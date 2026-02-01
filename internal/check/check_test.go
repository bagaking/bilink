package check

import "testing"

func TestDetectWarnings(t *testing.T) {
	files := []File{
		{Path: "a.md", ResolveKey: "foo", LintKey: "foo"},
		{Path: "b.md", ResolveKey: "bar", LintKey: "foo"},
	}
	errs, warns := Detect(files)
	if len(errs) != 0 {
		t.Fatalf("expected no errors")
	}
	if len(warns) != 1 {
		t.Fatalf("expected warning")
	}
}
