package resolve

import "testing"

func TestNormalizeName_ResolveRules(t *testing.T) {
	rules := Rules{
		CaseInsensitive:      true,
		IgnoreExtension:      true,
		SeparatorEquivalents: []string{" ", "-", "_"},
		UnicodeNormalize:     "NFKC",
	}
	got := NormalizeName("Foo_Bar.md", rules)
	if got != "foo-bar" {
		t.Fatalf("expected foo-bar, got %q", got)
	}
}
