package parse

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

func ExtractAnchors(content string) map[string]string {
	anchors := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	inFence := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(line, "#") {
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			slug := slugifyGitHub(title, anchors)
			anchors[slug] = title
		}
	}
	return anchors
}

func slugifyGitHub(title string, existing map[string]string) string {
	lower := strings.ToLower(title)
	var b strings.Builder
	prevDash := false
	for _, r := range lower {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if r == ' ' || r == '-' {
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if _, ok := existing[slug]; !ok {
		return slug
	}
	for i := 1; ; i++ {
		candidate := slug + "-" + strconv.Itoa(i)
		if _, ok := existing[candidate]; !ok {
			return candidate
		}
	}
}
