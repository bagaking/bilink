package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(ConfigOpts{Roots: []string{dir}})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(cfg.Scan.Extensions) == 0 {
		t.Fatalf("expected default extensions")
	}
}

func TestLoadConfig_DefaultPathOverrides(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".bilink")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "settings.toml"), []byte(`
[scan]
extensions = [".md"]
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(ConfigOpts{Roots: []string{dir}})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(cfg.Scan.Extensions) != 1 || cfg.Scan.Extensions[0] != ".md" {
		t.Fatalf("expected default path override extensions, got %#v", cfg.Scan.Extensions)
	}
}

func TestLoadConfig_ExplicitMissingErrors(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(ConfigOpts{ConfigPath: filepath.Join(dir, "missing.toml")})
	if err == nil {
		t.Fatalf("expected missing explicit config to error")
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
