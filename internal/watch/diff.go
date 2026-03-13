package watch

import "sort"

type Index struct{ Files []string }

type DiffResult struct {
	Added   []string
	Removed []string
}

func Diff(oldIdx, newIdx Index) DiffResult {
	oldSet := map[string]struct{}{}
	for _, f := range oldIdx.Files {
		oldSet[f] = struct{}{}
	}
	newSet := map[string]struct{}{}
	for _, f := range newIdx.Files {
		newSet[f] = struct{}{}
	}
	var added []string
	var removed []string
	for f := range newSet {
		if _, ok := oldSet[f]; !ok {
			added = append(added, f)
		}
	}
	for f := range oldSet {
		if _, ok := newSet[f]; !ok {
			removed = append(removed, f)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return DiffResult{Added: added, Removed: removed}
}
