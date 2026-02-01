package rename

import "testing"

func TestRewriteWikiLinks(t *testing.T) {
	content := "See [[Foo]] and [[Foo|alias]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar]] and [[Bar|alias]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}

func TestRewriteWikiAlias(t *testing.T) {
	content := "See [[Foo|My Alias]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar|My Alias]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}

func TestRewriteWikiAnchor(t *testing.T) {
	content := "See [[Foo#section]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar#section]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}

func TestRewriteWikiAnchorAlias(t *testing.T) {
	content := "See [[Foo#section|Alias]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar#section|Alias]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}

func TestRewriteMarkdown(t *testing.T) {
	content := "See [B](b.md)"
	out := RewriteMarkdown(content, "b.md", "c.md")
	if out != "See [B](c.md)" {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestRewriteMarkdownAnchor(t *testing.T) {
	content := "See [B](b.md#intro)"
	out := RewriteMarkdown(content, "b.md", "c.md")
	if out != "See [B](c.md#intro)" {
		t.Fatalf("unexpected: %s", out)
	}
}
