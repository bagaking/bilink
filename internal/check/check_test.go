package check

import "testing"

func TestDetectConflicts(t *testing.T) {
	files := []File{{Path: "a.md", ResolveKey: "foo"}, {Path: "b.md", ResolveKey: "foo"}}
	errs, warns := Detect(files)
	if len(errs) != 1 {
		t.Fatalf("expected conflict error")
	}
	if len(warns) != 0 {
		t.Fatalf("expected no warnings")
	}
}
