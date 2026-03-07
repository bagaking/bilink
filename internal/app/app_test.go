package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunUnknownCommand(t *testing.T) {
	err := Run([]string{"nope"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunRefsMissingArg(t *testing.T) {
	if err := Run([]string{"refs"}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunRenameParsesFlagAfterPositionals(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "new.md")
	refPath := filepath.Join(dir, "ref.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	if err := os.WriteFile(refPath, []byte("See [Old](old.md)"), 0o644); err != nil {
		t.Fatalf("write ref: %v", err)
	}

	if err := Run([]string{"rename", "--root", dir, oldPath, newPath, "--no-move"}); err != nil {
		t.Fatalf("run rename: %v", err)
	}

	if _, err := os.Stat(oldPath); err != nil {
		t.Fatalf("expected old file to remain: %v", err)
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Fatalf("expected new file not to exist, got %v", err)
	}
	data, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatalf("read ref: %v", err)
	}
	if string(data) != "See [Old](new.md)" {
		t.Fatalf("expected link rewrite, got %q", string(data))
	}
}

func TestRunRenameRejectsUnknownFlagAfterPositionals(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "new.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	if err := Run([]string{"rename", "--root", dir, oldPath, newPath, "--bogus"}); err == nil {
		t.Fatalf("expected unknown flag error")
	}
}

func TestRunRenameRejectsValueFlagMissingValueAfterPositionals(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.md")
	newPath := filepath.Join(dir, "new.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	err := Run([]string{"rename", oldPath, newPath, "--root"})
	if err == nil || err.Error() != "flag needs an argument: -root" {
		t.Fatalf("Run(rename old new --root) error = %v, want missing value error", err)
	}
}

func TestRunRenamePreservesDelimiterPositionals(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "-old.md")
	newPath := filepath.Join(dir, "new.md")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	if err := Run([]string{"rename", "--root", dir, "--", oldPath, newPath}); err != nil {
		t.Fatalf("run rename: %v", err)
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old file to move, got %v", err)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected new file: %v", err)
	}
}
