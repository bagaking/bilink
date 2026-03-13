package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/resolve"
)

func TestRunCheck(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err == nil {
		_ = payload
	}
}

func TestRunCheckGroups(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	_ = payload.ErrorGroups
}

func TestRunCheckPayloadSortedDeterministically(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "zeta-one.md"), "Hi")
	writeFile(t, filepath.Join(dir, "zeta_one.md"), "Hi")
	writeFile(t, filepath.Join(dir, "alpha-one.md"), "Hi")
	writeFile(t, filepath.Join(dir, "alpha_one.md"), "Hi")

	payload, err := RunCheck(CheckInput{
		Roots:      []string{dir},
		Extensions: []string{".md"},
		ResolveRules: resolve.Rules{
			SeparatorEquivalents: []string{"_", "-"},
		},
		LintRules: resolve.Rules{
			SeparatorEquivalents: []string{"_", "-"},
		},
	})
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	assertStringSlice(t, payload.Errors, []string{"alpha-one", "zeta-one"})
	assertStringSlice(t, payload.Warnings, []string{"alpha-one", "zeta-one"})
	assertCheckGroups(t, payload.ErrorGroups, []wantCheckGroup{
		{Key: "alpha-one", Paths: []string{
			filepath.Join(dir, "alpha-one.md"),
			filepath.Join(dir, "alpha_one.md"),
		}},
		{Key: "zeta-one", Paths: []string{
			filepath.Join(dir, "zeta-one.md"),
			filepath.Join(dir, "zeta_one.md"),
		}},
	})
	assertCheckGroups(t, payload.WarningGroups, []wantCheckGroup{
		{Key: "alpha-one", Paths: []string{
			filepath.Join(dir, "alpha-one.md"),
			filepath.Join(dir, "alpha_one.md"),
		}},
		{Key: "zeta-one", Paths: []string{
			filepath.Join(dir, "zeta-one.md"),
			filepath.Join(dir, "zeta_one.md"),
		}},
	})
}

type wantCheckGroup struct {
	Key   string
	Paths []string
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func assertStringSlice(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d strings, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("string %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func assertCheckGroups(t *testing.T, got []output.CheckGroup, want []wantCheckGroup) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d groups, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i].Key != want[i].Key {
			t.Fatalf("group %d key: got %q, want %q", i, got[i].Key, want[i].Key)
		}
		assertStringSlice(t, got[i].Paths, want[i].Paths)
	}
}
