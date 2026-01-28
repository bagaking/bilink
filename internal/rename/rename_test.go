package rename

import "testing"

func TestRewriteWikiLinks(t *testing.T) {
	content := "See [[Foo]] and [[Foo|alias]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar]] and [[Bar|alias]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}
