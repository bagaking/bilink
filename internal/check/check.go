package check

import "sort"

type File struct {
	Path       string
	ResolveKey string
	LintKey    string
}

type Group struct {
	Key   string
	Paths []string
}

func Detect(files []File) ([]Group, []Group) {
	resolveSeen := map[string][]string{}
	lintSeen := map[string][]string{}
	for _, f := range files {
		resolveSeen[f.ResolveKey] = append(resolveSeen[f.ResolveKey], f.Path)
		if f.LintKey != "" {
			lintSeen[f.LintKey] = append(lintSeen[f.LintKey], f.Path)
		}
	}
	var errs []Group
	for key, paths := range resolveSeen {
		if len(paths) > 1 {
			errs = append(errs, Group{Key: key, Paths: sortedStrings(paths)})
		}
	}
	var warns []Group
	for key, paths := range lintSeen {
		if len(paths) > 1 {
			warns = append(warns, Group{Key: key, Paths: sortedStrings(paths)})
		}
	}
	sortGroups(errs)
	sortGroups(warns)
	return errs, warns
}

func sortedStrings(values []string) []string {
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	return sorted
}

func sortGroups(groups []Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})
}
