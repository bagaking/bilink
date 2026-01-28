package index

import (
	"path/filepath"
	"strings"

	"github.com/bagaking/bilink/internal/parse"
)

type FileInput struct {
	Path    string
	Content string
}

type Link = parse.Link

type Index struct {
	Outbound map[string][]Link
	Inbound  map[string][]Link
}

func Build(files []FileInput) Index {
	idx := Index{Outbound: map[string][]parse.Link{}, Inbound: map[string][]parse.Link{}}
	nameToPath := map[string]string{}
	for _, f := range files {
		base := strings.TrimSuffix(filepath.Base(f.Path), filepath.Ext(f.Path))
		if _, ok := nameToPath[base]; !ok {
			nameToPath[base] = f.Path
		}
	}
	for _, f := range files {
		links := parse.ParseLinks(f.Content)
		idx.Outbound[f.Path] = links
		for _, link := range links {
			target := link.Target
			if link.Kind == "md" {
				target = link.Path
			} else if link.Kind == "wiki" {
				if resolved, ok := nameToPath[link.Target]; ok {
					target = resolved
				}
			}
			idx.Inbound[target] = append(idx.Inbound[target], link)
		}
	}
	return idx
}
