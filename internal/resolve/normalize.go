package resolve

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Rules struct {
	CaseInsensitive      bool
	IgnoreExtension      bool
	SeparatorEquivalents []string
	UnicodeNormalize     string
}

func NormalizeName(input string, rules Rules) string {
	name := input
	if rules.IgnoreExtension {
		ext := filepath.Ext(name)
		if ext != "" {
			name = strings.TrimSuffix(name, ext)
		}
	}
	if rules.UnicodeNormalize == "NFKC" {
		name = norm.NFKC.String(name)
	}
	if rules.CaseInsensitive {
		name = strings.ToLower(name)
	}
	for _, sep := range rules.SeparatorEquivalents {
		name = strings.ReplaceAll(name, sep, "-")
	}
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return name
}
