package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanAndIndex(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("See [[b]]"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	idx, err := ScanAndIndex([]string{dir}, []string{".md"})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.Outbound[filepath.Join(dir, "a.md")]) != 1 {
		t.Fatalf("expected outbound")
	}
}
