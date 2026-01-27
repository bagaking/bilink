package parse

import (
	"regexp"
	"strings"
)

type Link struct {
	Kind   string
	Target string
	Alias  string
	Path   string
	Anchor string
}

var wikiRe = regexp.MustCompile(`\[\[([^\]|#]+)(?:#([^\]|]+))?(?:\|([^\]]+))?\]\]`)
var mdRe = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

func ParseLinks(content string) []Link {
	var links []Link
	for _, m := range wikiRe.FindAllStringSubmatch(content, -1) {
		link := Link{Kind: "wiki", Target: strings.TrimSpace(m[1])}
		if m[2] != "" {
			link.Anchor = strings.TrimSpace(m[2])
		}
		if m[3] != "" {
			link.Alias = strings.TrimSpace(m[3])
		}
		links = append(links, link)
	}
	for _, m := range mdRe.FindAllStringSubmatch(content, -1) {
		target := strings.TrimSpace(m[2])
		if isExternal(target) {
			continue
		}
		path, anchor := splitAnchor(target)
		links = append(links, Link{Kind: "md", Path: path, Anchor: anchor})
	}
	return links
}

func splitAnchor(path string) (string, string) {
	parts := strings.SplitN(path, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return path, ""
}

func isExternal(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:")
}
