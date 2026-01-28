package rename

import "regexp"

var wikiRe = regexp.MustCompile(`\[\[([^\]|#]+)([^\]]*)\]\]`)

func RewriteWiki(content, from, to string) string {
	return wikiRe.ReplaceAllStringFunc(content, func(match string) string {
		if match == "[["+from+"]]" {
			return "[[" + to + "]]"
		}
		if match == "[["+from+"|alias]]" {
			return "[[" + to + "|alias]]"
		}
		return match
	})
}
