package service

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestRunRenameOnlyWritesChangedFilesAndPreservesMode(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "new.md")
	changedPath := filepath.Join(dir, "changed.md")
	unchangedPath := filepath.Join(dir, "unchanged.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	if err := os.WriteFile(changedPath, []byte("See [Old](old.md)"), 0o640); err != nil {
		t.Fatalf("write changed: %v", err)
	}
	if err := os.WriteFile(unchangedPath, []byte("No references here"), 0o600); err != nil {
		t.Fatalf("write unchanged: %v", err)
	}

	payload, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         false,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if !reflect.DeepEqual(payload.Updated, []string{changedPath}) {
		t.Fatalf("expected only changed file in updated, got %#v", payload.Updated)
	}
	for path, wantMode := range map[string]os.FileMode{
		changedPath:   0o640,
		unchangedPath: 0o600,
		oldPath:       0o644,
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != wantMode {
			t.Fatalf("expected %s mode %v, got %v", path, wantMode, got)
		}
	}
}

func TestRunRenameRewritesRelativeMarkdownTargetsFromSourceDir(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "sub")
	notesDir := filepath.Join(dir, "notes")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.MkdirAll(notesDir, 0o755); err != nil {
		t.Fatalf("mkdir notes: %v", err)
	}
	oldPath := filepath.Join(subDir, "b.md")
	newPath := filepath.Join(subDir, "c.md")
	rootRef := filepath.Join(dir, "root.md")
	nestedRef := filepath.Join(notesDir, "nested.md")
	unrelated := filepath.Join(dir, "unrelated.md")
	if err := os.WriteFile(oldPath, []byte("b"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	if err := os.WriteFile(rootRef, []byte("See [B](sub/b.md#x)"), 0o644); err != nil {
		t.Fatalf("write root ref: %v", err)
	}
	if err := os.WriteFile(nestedRef, []byte("See [B](../sub/b.md)"), 0o644); err != nil {
		t.Fatalf("write nested ref: %v", err)
	}
	if err := os.WriteFile(unrelated, []byte("See [Other](b.md)"), 0o644); err != nil {
		t.Fatalf("write unrelated: %v", err)
	}

	payload, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         false,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if !sameStrings(payload.Updated, []string{rootRef, nestedRef}) {
		t.Fatalf("expected root and nested refs updated, got %#v", payload.Updated)
	}
	assertFileContent(t, rootRef, "See [B](sub/c.md#x)")
	assertFileContent(t, nestedRef, "See [B](../sub/c.md)")
	assertFileContent(t, unrelated, "See [Other](b.md)")
}

func TestRunRenameDoesNotRewriteReferencesWhenMoveFails(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "missing", "new.md")
	refPath := filepath.Join(dir, "ref.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	if err := os.WriteFile(refPath, []byte("See [Old](old.md)"), 0o644); err != nil {
		t.Fatalf("write ref: %v", err)
	}

	_, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         true,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err == nil {
		t.Fatalf("expected rename to fail")
	}
	assertFileContent(t, refPath, "See [Old](old.md)")
	if _, err := os.Stat(oldPath); err != nil {
		t.Fatalf("expected old path to remain: %v", err)
	}
}

func TestRunRenameDoesNotMoveWhenPreparingWritesFails(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "missing", "new.md")
	if err := os.WriteFile(oldPath, []byte("Self [[old]]"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	_, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         true,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err == nil {
		t.Fatalf("expected rename to fail")
	}
	assertFileContent(t, oldPath, "Self [[old]]")
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Fatalf("expected new path not to exist, got %v", err)
	}
}

func TestRunRenameRewritesMovedFileWithoutRecreatingOldPath(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "new.md")
	if err := os.WriteFile(oldPath, []byte("Self [[old]] and [Old](old.md)"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	_, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         true,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old path not to be recreated, got %v", err)
	}
	assertFileContent(t, newPath, "Self [[new]] and [Old](new.md)")
}

func TestRunRenameRewritesMovedFileMarkdownFromNewDirectory(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(subDir, "new.md")
	if err := os.WriteFile(oldPath, []byte("Self [Old](old.md)"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	_, err := RunRename(RenameInput{
		Roots:        []string{dir},
		OldPath:      oldPath,
		NewPath:      newPath,
		Move:         true,
		Extensions:   []string{".md"},
		ResolveRules: resolve.Rules{CaseInsensitive: true, IgnoreExtension: true},
	})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	assertFileContent(t, newPath, "Self [Old](new.md)")
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

func assertFileContent(t *testing.T, path string, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("expected %s content %q, got %q", path, want, string(data))
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	counts := map[string]int{}
	for _, item := range a {
		counts[item]++
	}
	for _, item := range b {
		counts[item]--
		if counts[item] < 0 {
			return false
		}
	}
	return true
}
