package service

import (
	"os"

	"github.com/bagaking/bilink/internal/index"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/watch"
)

type WatchInput struct {
	IndexPath  string
	Roots      []string
	Extensions []string
}

func RunWatch(input WatchInput) (output.WatchPayload, error) {
	if _, err := os.Stat(input.IndexPath); err != nil {
		return output.WatchPayload{}, err
	}
	oldIdx, err := index.Load(input.IndexPath)
	if err != nil {
		return output.WatchPayload{}, err
	}
	newIdx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.WatchPayload{}, err
	}
	diff := watch.Diff(watch.Index{Files: keys(oldIdx.Outbound)}, watch.Index{Files: keys(newIdx.Outbound)})
	if err := index.Save(input.IndexPath, newIdx); err != nil {
		return output.WatchPayload{}, err
	}
	return output.WatchPayload{Added: diff.Added, Removed: diff.Removed}, nil
}

func keys(m map[string][]index.Link) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
