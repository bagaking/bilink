package app

import (
	"encoding/json"
	"io"
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

func TestRunRefsLoadsDefaultConfigFromRoot(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".bilink")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "settings.toml"), []byte(`
[workspace]
roots = ["."]

[scan]
extensions = [".txt"]
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	source := filepath.Join(dir, "source.txt")
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(source, []byte("See [[target]]"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(target, []byte("target"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}

	var runErr error
	stdout := captureStdout(t, func() {
		runErr = Run([]string{"refs", "--root", dir, "--json", target})
	})
	if runErr != nil {
		t.Fatalf("run refs: %v", runErr)
	}

	var payload struct {
		Inbound []any `json:"inbound"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("unmarshal refs output: %v\noutput: %s", err, stdout)
	}
	if len(payload.Inbound) != 1 {
		t.Fatalf("expected default root config to include .txt refs, got %#v", payload.Inbound)
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

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = writePipe

	fn()

	os.Stdout = oldStdout
	if err := writePipe.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	out, err := io.ReadAll(readPipe)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := readPipe.Close(); err != nil {
		t.Fatalf("close read pipe: %v", err)
	}
	os.Stdout = oldStdout
	return string(out)
}
