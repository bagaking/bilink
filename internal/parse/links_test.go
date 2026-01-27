package parse

import "testing"

func TestParseLinks(t *testing.T) {
	content := "See [[Foo|alias]] and [Bar](docs/bar.md#intro) and http://example.com"
	links := ParseLinks(content)
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
	if links[0].Target != "Foo" || links[0].Alias != "alias" {
		t.Fatalf("unexpected wiki link %#v", links[0])
	}
	if links[1].Path != "docs/bar.md" || links[1].Anchor != "intro" {
		t.Fatalf("unexpected md link %#v", links[1])
	}
}
