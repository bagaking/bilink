package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(ConfigOpts{
		Roots: []string{dir},
		ConfigPath: filepath.Join(dir, ".bilink", "settings.toml"),
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(cfg.Scan.Extensions) == 0 {
		t.Fatalf("expected default extensions")
	}
}

func TestLoadConfig_FileOverrides(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "settings.toml")
	if err := os.WriteFile(cfgPath, []byte(`
[scan]
extensions = [".md"]
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(ConfigOpts{ConfigPath: cfgPath})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(cfg.Scan.Extensions) != 1 || cfg.Scan.Extensions[0] != ".md" {
		t.Fatalf("expected override extensions, got %#v", cfg.Scan.Extensions)
	}
}
