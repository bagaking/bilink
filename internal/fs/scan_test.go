package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanRoots(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("no"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	files, err := ScanRoots([]string{dir}, []string{".md"})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}
