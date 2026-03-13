package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bagaking/bilink/internal/index"
)

func TestRunWatchRequiresIndex(t *testing.T) {
	_, err := RunWatch(WatchInput{IndexPath: ".bilink/index.json"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunWatchPayloadSortedDeterministically(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, ".bilink", "index.json")
	oldRemovedZ := filepath.Join(dir, "z-removed.md")
	oldRemovedB := filepath.Join(dir, "b-removed.md")
	stable := filepath.Join(dir, "stable.md")
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		t.Fatalf("mkdir index dir: %v", err)
	}
	if err := index.Save(indexPath, index.Index{
		Outbound: map[string][]index.Link{
			oldRemovedZ: nil,
			oldRemovedB: nil,
			stable:      nil,
		},
	}); err != nil {
		t.Fatalf("save index: %v", err)
	}

	writeFile(t, stable, "Hi")
	writeFile(t, filepath.Join(dir, "y-added.md"), "Hi")
	writeFile(t, filepath.Join(dir, "a-added.md"), "Hi")

	payload, err := RunWatch(WatchInput{
		IndexPath:  indexPath,
		Roots:      []string{dir},
		Extensions: []string{".md"},
	})
	if err != nil {
		t.Fatalf("watch: %v", err)
	}

	assertStringSlice(t, payload.Added, []string{
		filepath.Join(dir, "a-added.md"),
		filepath.Join(dir, "y-added.md"),
	})
	assertStringSlice(t, payload.Removed, []string{
		oldRemovedB,
		oldRemovedZ,
	})
}
