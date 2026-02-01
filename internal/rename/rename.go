package rename

import (
	"regexp"
	"strings"
)

var wikiRe = regexp.MustCompile(`\[\[([^\]|#]+)(?:#([^\]|]+))?(?:\|([^\]]+))?\]\]`)
var mdTargetRe = regexp.MustCompile(`\]\(([^\)]+)\)`)

func RewriteWiki(content, from, to string) string {
	return wikiRe.ReplaceAllStringFunc(content, func(match string) string {
		parts := wikiRe.FindStringSubmatch(match)
		if len(parts) == 0 {
			return match
		}
		if strings.TrimSpace(parts[1]) != from {
			return match
		}
		var b strings.Builder
		b.WriteString("[[")
		b.WriteString(to)
		if parts[2] != "" {
			b.WriteString("#")
			b.WriteString(parts[2])
		}
		if parts[3] != "" {
			b.WriteString("|")
			b.WriteString(parts[3])
		}
		b.WriteString("]]")
		return b.String()
	})
}

func RewriteMarkdown(content, from, to string) string {
	return mdTargetRe.ReplaceAllStringFunc(content, func(match string) string {
		if len(match) < 3 {
			return match
		}
		target := match[2 : len(match)-1]
		if target == from {
			return "](" + to + ")"
		}
		if strings.HasPrefix(target, from+"#") {
			suffix := strings.TrimPrefix(target, from)
			return "](" + to + suffix + ")"
		}
		return match
	})
}
