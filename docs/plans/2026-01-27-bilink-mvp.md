---
title: Bilink MVP Implementation Plan
required: false
sop:
  - Read this plan before implementing Bilink MVP tasks.
  - Update this plan when execution steps change.
  - Regenerate must-sop.md after updating this doc.
---
# Bilink MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement Bilink MVP (Go CLI + Watch TUI + TS wrapper) per OpenSpec change `add-bilink-mvp`.

**Architecture:** Go core handles scanning/parsing/indexing/updates; Bubble Tea provides watch TUI; a Node wrapper downloads platform binaries for `npx`. Resolve vs lint rules split parsing and style warnings, and update policy is Balanced.

**Tech Stack:** Go 1.22+, Bubble Tea, fsnotify, TOML, Node (TS wrapper).

---

### Task 1: Initialize Go module and baseline CLI entrypoint

**Files:**
- Create: `go.mod`
- Create: `cmd/bilink/main.go`

**Step 1: Create go.mod**

```go
module github.com/bagaking/bilink

go 1.22
```

**Step 2: Create minimal CLI entrypoint**

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stdout, "bilink: MVP bootstrap")
}
```

**Step 3: Run build to verify baseline**

Run: `go build ./...`
Expected: PASS

**Step 4: Commit**

```bash
git add go.mod cmd/bilink/main.go
git commit -m "chore: initialize go module and CLI entrypoint"
```

---

### Task 2: Config loader and defaults (TDD)

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

**Step 1: Write failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config -v`
Expected: FAIL with "undefined: Load"

**Step 3: Implement minimal config loader**

```go
package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Workspace Workspace `toml:"workspace"`
	Scan      Scan      `toml:"scan"`
	Resolve   Resolve   `toml:"resolveRules"`
	Lint      Lint      `toml:"lintRules"`
	Anchors   Anchors   `toml:"anchors"`
	Update    Update    `toml:"updatePolicy"`
	Index     Index     `toml:"index"`
}

type Workspace struct{ Roots []string `toml:"roots"` }
type Scan struct{ Extensions []string `toml:"extensions"` }
type Resolve struct {
	CaseInsensitive      bool     `toml:"caseInsensitive"`
	IgnoreExtension      bool     `toml:"ignoreExtension"`
	SeparatorEquivalents []string `toml:"separatorEquivalents"`
	UnicodeNormalize     string   `toml:"unicodeNormalize"`
}
type Lint struct {
	RequireExactCase        bool `toml:"requireExactCase"`
	RequireExplicitExtension bool `toml:"requireExplicitExtension"`
	RequireExactSeparators  bool `toml:"requireExactSeparators"`
}
type Anchors struct {
	Style string `toml:"style"`
	Mode  string `toml:"mode"`
}
type Update struct{ Mode string `toml:"mode"` }
type Index struct {
	Path             string `toml:"path"`
	RequiredForWatch bool   `toml:"requiredForWatch"`
}

type ConfigOpts struct {
	Roots      []string
	ConfigPath string
}

func Defaults() Config {
	return Config{
		Workspace: Workspace{Roots: []string{"."}},
		Scan:      Scan{Extensions: []string{".md", ".markdown", ".mdx"}},
		Resolve: Resolve{
			CaseInsensitive:      true,
			IgnoreExtension:      true,
			SeparatorEquivalents: []string{" ", "-", "_"},
			UnicodeNormalize:     "NFKC",
		},
		Lint: Lint{
			RequireExactCase:        true,
			RequireExplicitExtension: true,
			RequireExactSeparators:  true,
		},
		Anchors: Anchors{Style: "github", Mode: "resolve-only"},
		Update:  Update{Mode: "balanced"},
		Index:   Index{Path: ".bilink/index.json", RequiredForWatch: true},
	}
}

func Load(opts ConfigOpts) (Config, error) {
	cfg := Defaults()
	if len(opts.Roots) > 0 {
		cfg.Workspace.Roots = opts.Roots
	}
	if opts.ConfigPath == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(opts.ConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/config -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add config defaults and loader"
```

---

### Task 3: Resolve/lint normalization (TDD)

**Files:**
- Create: `internal/resolve/normalize.go`
- Create: `internal/resolve/normalize_test.go`

**Step 1: Write failing test**

```go
package resolve

import "testing"

func TestNormalizeName_ResolveRules(t *testing.T) {
	rules := Rules{
		CaseInsensitive: true,
		IgnoreExtension: true,
		SeparatorEquivalents: []string{" ", "-", "_"},
		UnicodeNormalize: "NFKC",
	}
	got := NormalizeName("Foo_Bar.md", rules)
	if got != "foo-bar" {
		t.Fatalf("expected foo-bar, got %q", got)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/resolve -v`
Expected: FAIL with "undefined: Rules"

**Step 3: Implement minimal normalization**

```go
package resolve

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Rules struct {
	CaseInsensitive      bool
	IgnoreExtension      bool
	SeparatorEquivalents []string
	UnicodeNormalize     string
}

func NormalizeName(input string, rules Rules) string {
	name := input
	if rules.IgnoreExtension {
		ext := filepath.Ext(name)
		if ext != "" {
			name = strings.TrimSuffix(name, ext)
		}
	}
	if rules.UnicodeNormalize == "NFKC" {
		name = norm.NFKC.String(name)
	}
	if rules.CaseInsensitive {
		name = strings.ToLower(name)
	}
	for _, sep := range rules.SeparatorEquivalents {
		name = strings.ReplaceAll(name, sep, "-")
	}
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return name
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/resolve -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/resolve/normalize.go internal/resolve/normalize_test.go
git commit -m "feat: add resolve normalization rules"
```

---

### Task 4: Link parsing (TDD)

**Files:**
- Create: `internal/parse/links.go`
- Create: `internal/parse/links_test.go`

**Step 1: Write failing test**

```go
package parse

import "testing"

func TestParseLinks(t *testing.T) {
	content := "See [[Foo|alias]] and [Bar](docs/bar.md#intro) and http://example.com"
	links := ParseLinks(content)
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
	if links[0].Target != "Foo" || links[0].Alias != "alias" {
		t.Fatalf("unexpected wiki link %#v", links[0])
	}
	if links[1].Path != "docs/bar.md" || links[1].Anchor != "intro" {
		t.Fatalf("unexpected md link %#v", links[1])
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/parse -v`
Expected: FAIL with "undefined: ParseLinks"

**Step 3: Implement minimal parser**

```go
package parse

import (
	"regexp"
	"strings"
)

type Link struct {
	Kind   string // "wiki" or "md"
	Target string // for wiki
	Alias  string
	Path   string // for md
	Anchor string
}

var wikiRe = regexp.MustCompile(`\\[\\[([^\\]|#]+)(?:#([^\\]|]+))?(?:\\|([^\\]]+))?\\]\\]`)
var mdRe = regexp.MustCompile(`\\[([^\\]]+)\\]\\(([^\\)]+)\\)`)

func ParseLinks(content string) []Link {
	var links []Link
	for _, m := range wikiRe.FindAllStringSubmatch(content, -1) {
		link := Link{Kind: "wiki", Target: strings.TrimSpace(m[1])}
		if m[2] != "" {
			link.Anchor = strings.TrimSpace(m[2])
		}
		if m[3] != "" {
			link.Alias = strings.TrimSpace(m[3])
		}
		links = append(links, link)
	}
	for _, m := range mdRe.FindAllStringSubmatch(content, -1) {
		target := strings.TrimSpace(m[2])
		if isExternal(target) {
			continue
		}
		path, anchor := splitAnchor(target)
		links = append(links, Link{Kind: "md", Path: path, Anchor: anchor})
	}
	return links
}

func splitAnchor(path string) (string, string) {
	parts := strings.SplitN(path, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return path, ""
}

func isExternal(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:")
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/parse -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/parse/links.go internal/parse/links_test.go
git commit -m "feat: add markdown link parser"
```

---

### Task 5: Anchor extraction (TDD)

**Files:**
- Create: `internal/parse/anchors.go`
- Create: `internal/parse/anchors_test.go`

**Step 1: Write failing test**

```go
package parse

import "testing"

func TestExtractAnchors(t *testing.T) {
	content := "# Overview\\n\\n## Intro\\n```\\n# Not a header\\n```\\n"
	anchors := ExtractAnchors(content)
	if anchors["overview"] == "" || anchors["intro"] == "" {
		t.Fatalf("expected headings, got %#v", anchors)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/parse -v`
Expected: FAIL with "undefined: ExtractAnchors"

**Step 3: Implement minimal anchor extraction**

```go
package parse

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

func ExtractAnchors(content string) map[string]string {
	anchors := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	inFence := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(line, "#") {
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			slug := slugifyGitHub(title, anchors)
			anchors[slug] = title
		}
	}
	return anchors
}

func slugifyGitHub(title string, existing map[string]string) string {
	lower := strings.ToLower(title)
	var b strings.Builder
	prevDash := false
	for _, r := range lower {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if r == ' ' || r == '-' {
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if _, ok := existing[slug]; !ok {
		return slug
	}
	for i := 1; ; i++ {
		candidate := slug + "-" + strconv.Itoa(i)
		if _, ok := existing[candidate]; !ok {
			return candidate
		}
	}
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/parse -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/parse/anchors.go internal/parse/anchors_test.go
git commit -m "feat: add heading anchor extraction"
```

---

### Task 6: Scan and index (TDD)

**Files:**
- Create: `internal/fs/scan.go`
- Create: `internal/fs/scan_test.go`
- Create: `internal/index/index.go`
- Create: `internal/index/index_test.go`

**Step 1: Write failing tests**

```go
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
```

```go
package index

import (
	"testing"

	"github.com/bagaking/bilink/internal/parse"
)

func TestBuildIndex(t *testing.T) {
	files := []FileInput{
		{Path: "a.md", Content: "See [[b]]"},
		{Path: "b.md", Content: "Hi"},
	}
	idx := Build(files)
	if len(idx.Outbound["a.md"]) != 1 {
		t.Fatalf("expected outbound link")
	}
	if len(idx.Inbound["b.md"]) != 1 {
		t.Fatalf("expected inbound link")
	}
	if idx.Outbound["a.md"][0].Target != "b" {
		t.Fatalf("unexpected link %#v", idx.Outbound["a.md"][0])
	}
	_ = parse.Link{}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/fs ./internal/index -v`
Expected: FAIL with "undefined: ScanRoots / Build"

**Step 3: Implement scan + index**

```go
package fs

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func ScanRoots(roots []string, extensions []string) ([]string, error) {
	exts := map[string]struct{}{}
	for _, ext := range extensions {
		exts[strings.ToLower(ext)] = struct{}{}
	}
	var files []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			if _, ok := exts[ext]; ok {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
```

```go
package index

import "github.com/bagaking/bilink/internal/parse"

type FileInput struct {
	Path    string
	Content string
}

type Index struct {
	Outbound map[string][]parse.Link
	Inbound  map[string][]parse.Link
}

func Build(files []FileInput) Index {
	idx := Index{Outbound: map[string][]parse.Link{}, Inbound: map[string][]parse.Link{}}
	for _, f := range files {
		links := parse.ParseLinks(f.Content)
		idx.Outbound[f.Path] = links
		for _, link := range links {
			target := link.Target
			if link.Kind == "md" {
				target = link.Path
			}
			idx.Inbound[target] = append(idx.Inbound[target], link)
		}
	}
	return idx
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/fs ./internal/index -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/fs/scan.go internal/fs/scan_test.go internal/index/index.go internal/index/index_test.go
git commit -m "feat: add scanning and index build"
```

---

### Task 7: Refs command core (TDD)

**Files:**
- Create: `internal/refs/refs.go`
- Create: `internal/refs/refs_test.go`

**Step 1: Write failing test**

```go
package refs

import "testing"

func TestRefsFor(t *testing.T) {
	idx := Index{
		Outbound: map[string][]Link{"a.md": {{Target: "b"}}},
		Inbound:  map[string][]Link{"a.md": {{Target: "a"}}},
	}
	out, in := RefsFor(idx, "a.md")
	if len(out) != 1 || len(in) != 1 {
		t.Fatalf("unexpected refs")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/refs -v`
Expected: FAIL with "undefined: Index"

**Step 3: Implement minimal**

```go
package refs

type Link struct{ Target string }

type Index struct {
	Outbound map[string][]Link
	Inbound  map[string][]Link
}

func RefsFor(idx Index, path string) ([]Link, []Link) {
	return idx.Outbound[path], idx.Inbound[path]
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/refs -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/refs/refs.go internal/refs/refs_test.go
git commit -m "feat: add refs query core"
```

---

### Task 8: Conflict checks (TDD)

**Files:**
- Create: `internal/check/check.go`
- Create: `internal/check/check_test.go`

**Step 1: Write failing test**

```go
package check

import "testing"

func TestDetectConflicts(t *testing.T) {
	files := []File{{Path: "a.md", ResolveKey: "foo"}, {Path: "b.md", ResolveKey: "foo"}}
	errs, warns := Detect(files)
	if len(errs) != 1 {
		t.Fatalf("expected conflict error")
	}
	if len(warns) != 0 {
		t.Fatalf("expected no warnings")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/check -v`
Expected: FAIL with "undefined: Detect"

**Step 3: Implement minimal**

```go
package check

type File struct {
	Path       string
	ResolveKey string
	LintKey    string
}

func Detect(files []File) ([]string, []string) {
	seen := map[string][]string{}
	for _, f := range files {
		seen[f.ResolveKey] = append(seen[f.ResolveKey], f.Path)
	}
	var errs []string
	for key, paths := range seen {
		if len(paths) > 1 {
			errs = append(errs, key)
		}
	}
	return errs, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/check -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/check/check.go internal/check/check_test.go
git commit -m "feat: add conflict detection core"
```

---

### Task 9: Rename/update engine (TDD)

**Files:**
- Create: `internal/rename/rename.go`
- Create: `internal/rename/rename_test.go`

**Step 1: Write failing test**

```go
package rename

import "testing"

func TestRewriteWikiLinks(t *testing.T) {
	content := "See [[Foo]] and [[Foo|alias]]"
	out := RewriteWiki(content, "Foo", "Bar")
	if out != "See [[Bar]] and [[Bar|alias]]" {
		t.Fatalf("unexpected rewrite: %s", out)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/rename -v`
Expected: FAIL with "undefined: RewriteWiki"

**Step 3: Implement minimal**

```go
package rename

import "regexp"

var wikiRe = regexp.MustCompile(`\\[\\[([^\\]|#]+)([^\\]]*)\\]\\]`)

func RewriteWiki(content, from, to string) string {
	return wikiRe.ReplaceAllStringFunc(content, func(match string) string {
		if match == "[["+from+"]]" {
			return "[["+to+"]]"
		}
		if match == "[["+from+"|alias]]" {
			return "[["+to+"|alias]]"
		}
		return match
	})
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/rename -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/rename/rename.go internal/rename/rename_test.go
git commit -m "feat: add wiki link rewrite core"
```

---

### Task 10: Index persistence (TDD)

**Files:**
- Create: `internal/index/store.go`
- Create: `internal/index/store_test.go`

**Step 1: Write failing test**

```go
package index

import (
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/index.json"
	idx := Index{Outbound: map[string][]Link{"a.md": {{Target: "b"}}}}
	if err := Save(path, idx); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Outbound["a.md"]) != 1 {
		t.Fatalf("unexpected index")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/index -v`
Expected: FAIL with "undefined: Save"

**Step 3: Implement minimal**

```go
package index

import (
	"encoding/json"
	"os"
)

func Save(path string, idx Index) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func Load(path string) (Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Index{}, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return Index{}, err
	}
	return idx, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/index -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/index/store.go internal/index/store_test.go
git commit -m "feat: add index persistence"
```

---

### Task 11: Watch diff + TUI skeleton (TDD)

**Files:**
- Create: `internal/watch/diff.go`
- Create: `internal/watch/diff_test.go`
- Create: `internal/watch/tui/model.go`
- Create: `internal/watch/tui/model_test.go`

**Step 1: Write failing tests**

```go
package watch

import "testing"

func TestDiff(t *testing.T) {
	oldIdx := Index{Files: []string{"a.md"}}
	newIdx := Index{Files: []string{"a.md", "b.md"}}
	diff := Diff(oldIdx, newIdx)
	if len(diff.Added) != 1 || diff.Added[0] != "b.md" {
		t.Fatalf("unexpected diff")
	}
}
```

```go
package tui

import "testing"

func TestModelInit(t *testing.T) {
	m := NewModel()
	if m.Status == "" {
		t.Fatalf("expected status")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/watch ./internal/watch/tui -v`
Expected: FAIL with "undefined: Diff / NewModel"

**Step 3: Implement minimal**

```go
package watch

type Index struct{ Files []string }

type DiffResult struct {
	Added   []string
	Removed []string
}

func Diff(oldIdx, newIdx Index) DiffResult {
	oldSet := map[string]struct{}{}
	for _, f := range oldIdx.Files {
		oldSet[f] = struct{}{}
	}
	newSet := map[string]struct{}{}
	for _, f := range newIdx.Files {
		newSet[f] = struct{}{}
	}
	var added []string
	var removed []string
	for f := range newSet {
		if _, ok := oldSet[f]; !ok {
			added = append(added, f)
		}
	}
	for f := range oldSet {
		if _, ok := newSet[f]; !ok {
			removed = append(removed, f)
		}
	}
	return DiffResult{Added: added, Removed: removed}
}
```

```go
package tui

type Model struct {
	Status string
}

func NewModel() Model {
	return Model{Status: "starting"}
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/watch ./internal/watch/tui -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/watch/diff.go internal/watch/diff_test.go internal/watch/tui/model.go internal/watch/tui/model_test.go
git commit -m "feat: add watch diff and TUI skeleton"
```

---

### Task 12: CLI wiring and JSON output (TDD)

**Files:**
- Create: `internal/app/app.go`
- Create: `internal/app/app_test.go`
- Modify: `cmd/bilink/main.go`

**Step 1: Write failing test**

```go
package app

import "testing"

func TestRunUnknownCommand(t *testing.T) {
	err := Run([]string{"nope"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/app -v`
Expected: FAIL with "undefined: Run"

**Step 3: Implement minimal app runner**

```go
package app

import "fmt"

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}
	switch args[0] {
	case "refs", "rename", "check", "watch":
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
```

**Step 4: Update main to call app.Run**

```go
package main

import (
	"fmt"
	"os"

	"github.com/bagaking/bilink/internal/app"
)

func main() {
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

**Step 5: Run tests**

Run: `go test ./internal/app -v`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/app/app.go internal/app/app_test.go cmd/bilink/main.go
git commit -m "feat: add CLI dispatcher"
```

---

### Task 13: TS wrapper (TDD)

**Files:**
- Create: `packages/bilink-npx/package.json`
- Create: `packages/bilink-npx/bin/bilink.js`
- Create: `packages/bilink-npx/lib/platform.js`
- Create: `packages/bilink-npx/test/platform.test.mjs`
- Modify: `Makefile` (optional: add node test in make test)

**Step 1: Write failing test**

```js
import assert from "node:assert/strict";
import { resolveTarget } from "../lib/platform.js";

assert.equal(resolveTarget("darwin", "arm64").binary, "bilink-darwin-arm64");
assert.equal(resolveTarget("win32", "x64").binary, "bilink-windows-amd64.exe");
```

**Step 2: Run test to verify it fails**

Run: `node --test packages/bilink-npx/test/platform.test.mjs`
Expected: FAIL with "Cannot find module"

**Step 3: Implement platform resolver**

```js
export function resolveTarget(platform, arch) {
  if (platform === "darwin" && arch === "arm64") return { binary: "bilink-darwin-arm64" };
  if (platform === "darwin" && arch === "x64") return { binary: "bilink-darwin-amd64" };
  if (platform === "linux" && arch === "x64") return { binary: "bilink-linux-amd64" };
  if (platform === "linux" && arch === "arm64") return { binary: "bilink-linux-arm64" };
  if (platform === "win32" && arch === "x64") return { binary: "bilink-windows-amd64.exe" };
  throw new Error(`unsupported platform: ${platform}-${arch}`);
}
```

**Step 4: Run tests**

Run: `node --test packages/bilink-npx/test/platform.test.mjs`
Expected: PASS

**Step 5: Implement bin stub**

```js
#!/usr/bin/env node
import { resolveTarget } from "../lib/platform.js";
import { spawnSync } from "node:child_process";
import { existsSync, mkdirSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";

const { binary } = resolveTarget(process.platform, process.arch);
const cacheDir = path.join(homedir(), ".cache", "bilink");
const binPath = path.join(cacheDir, binary);

if (!existsSync(cacheDir)) mkdirSync(cacheDir, { recursive: true });
if (!existsSync(binPath)) {
  console.error("bilink binary missing; please download releases first");
  process.exit(1);
}
const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
process.exit(result.status ?? 1);
```

**Step 6: Commit**

```bash
git add packages/bilink-npx
git commit -m "feat: add npx wrapper skeleton"
```

---

### Task 14: Update quality gate (optional)

**Files:**
- Modify: `Makefile`

**Step 1: Add node tests to make test (optional)**

```makefile
test:
	@if [ -f go.mod ]; then \
		$(GO) test ./...; \
	else \
		echo "go.mod not found; skipping go test"; \
	fi
	@if [ -f packages/bilink-npx/test/platform.test.mjs ]; then \
		node --test packages/bilink-npx/test/platform.test.mjs; \
	fi
```

**Step 2: Commit**

```bash
git add Makefile
git commit -m "chore: add node tests to make test"
```
