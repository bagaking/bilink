package parse

import "testing"

func TestExtractAnchors(t *testing.T) {
	content := "# Overview\n\n## Intro\n```\n# Not a header\n```\n"
	anchors := ExtractAnchors(content)
	if anchors["overview"] == "" || anchors["intro"] == "" {
		t.Fatalf("expected headings, got %#v", anchors)
	}
}
