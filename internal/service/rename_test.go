package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bagaking/bilink/internal/resolve"
)

func TestRunRename(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.md")
	b := filepath.Join(dir, "b.md")
	if err := os.WriteFile(a, []byte("See [[b]] and [B](b.md)"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(b, []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	payload, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      b,
		NewPath:      filepath.Join(dir, "c.md"),
		Move:         true,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if len(payload.Updated) == 0 {
		t.Fatalf("expected updates")
	}
}

func TestRunRenameAmbiguous(t *testing.T) {
	dir := t.TempDir()
	rootA := filepath.Join(dir, "a")
	rootB := filepath.Join(dir, "b")
	if err := os.MkdirAll(rootA, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(rootB, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	oldPath := filepath.Join(rootA, "same.md")
	if err := os.WriteFile(oldPath, []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rootB, "same.md"), []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      filepath.Join(rootA, "new.md"),
		Move:         false,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
		Interactive:  false,
	})
	if err == nil {
		t.Fatalf("expected ambiguity error")
	}
}
