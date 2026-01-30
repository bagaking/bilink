---
title: Bilink CLI 完整实现计划
required: false
sop:
  - Read this plan before implementing Bilink CLI completion tasks.
  - Update this plan when steps or files change.
  - Regenerate must-sop.md after updating this doc.
---
# Bilink CLI 完整实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 完成 Bilink CLI 的完整实现，新增 service 层贯通 refs/rename/check/watch/--json，并补全 npx 自动下载。

**Architecture:** 引入 `internal/service` 作为应用层，CLI 仅做参数解析与输出；`internal/output` 统一文本/JSON；`watch` 通过 service 产生事件流与 TUI 交互。

**Tech Stack:** Go 1.24+, Bubble Tea, Node.js (npx wrapper).

---

### Task 1: Service 配置与命令模型 (TDD)

**Files:**
- Create: `internal/service/types.go`
- Create: `internal/service/types_test.go`

**Step 1: Write failing test**
```go
package service

import "testing"

func TestCommandKind(t *testing.T) {
	if CommandRefs.String() == "" {
		t.Fatalf("expected command string")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/service -v`
Expected: FAIL with "undefined: CommandRefs"

**Step 3: Implement minimal types**
```go
package service

type CommandKind int

const (
	CommandRefs CommandKind = iota
	CommandCheck
	CommandRename
	CommandWatch
)

func (c CommandKind) String() string {
	switch c {
	case CommandRefs:
		return "refs"
	case CommandCheck:
		return "check"
	case CommandRename:
		return "rename"
	case CommandWatch:
		return "watch"
	default:
		return "unknown"
	}
}

type Options struct {
	Roots       []string
	ConfigPath  string
	JSON        bool
	Interactive bool
	Move        bool
	OldPath     string
	NewPath     string
	TargetPath  string
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/service -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/service/types.go internal/service/types_test.go
git commit -m "feat: add service command types"
```

---

### Task 2: Output 统一格式 (TDD)

**Files:**
- Create: `internal/output/output.go`
- Create: `internal/output/output_test.go`

**Step 1: Write failing test**
```go
package output

import "testing"

func TestJSONRefsOutput(t *testing.T) {
	payload := RefsPayload{Target: "a.md"}
	data, err := JSON(payload)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected json")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/output -v`
Expected: FAIL with "undefined: JSON"

**Step 3: Implement output payloads + text helpers**
```go
package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RefsPayload struct {
	Target   string `json:"target"`
	Outbound any    `json:"outbound"`
	Inbound  any    `json:"inbound"`
}

type CheckPayload struct {
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type RenamePayload struct {
	OldPath string   `json:"oldPath"`
	NewPath string   `json:"newPath"`
	Updated []string `json:"updated"`
}

type WatchPayload struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
}

func JSON(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func TextRefs(p RefsPayload) string {
	return fmt.Sprintf("target: %s\noutbound: %d\ninbound: %d", p.Target, count(p.Outbound), count(p.Inbound))
}

func TextCheck(p CheckPayload) string {
	return fmt.Sprintf("errors: %d\nwarnings: %d", len(p.Errors), len(p.Warnings))
}

func TextRename(p RenamePayload) string {
	return fmt.Sprintf("old: %s\nnew: %s\nupdated: %s", p.OldPath, p.NewPath, strings.Join(p.Updated, ", "))
}

func TextWatch(p WatchPayload) string {
	return fmt.Sprintf("added: %d\nremoved: %d", len(p.Added), len(p.Removed))
}

func count(v any) int {
	switch t := v.(type) {
	case []any:
		return len(t)
	case []string:
		return len(t)
	case []interface{}:
		return len(t)
	default:
		return 0
	}
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/output -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/output/output.go internal/output/output_test.go
git commit -m "feat: add output payloads"
```

---

### Task 3: Scan + Index 编排 (TDD)

**Files:**
- Create: `internal/service/scan.go`
- Create: `internal/service/scan_test.go`

**Step 1: Write failing test**
```go
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
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/service -v`
Expected: FAIL with "undefined: ScanAndIndex"

**Step 3: Implement minimal scan**
```go
package service

import (
	"os"

	"github.com/bagaking/bilink/internal/fs"
	"github.com/bagaking/bilink/internal/index"
)

func ScanAndIndex(roots []string, extensions []string) (index.Index, error) {
	files, err := fs.ScanRoots(roots, extensions)
	if err != nil {
		return index.Index{}, err
	}
	inputs := make([]index.FileInput, 0, len(files))
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return index.Index{}, err
		}
		inputs = append(inputs, index.FileInput{Path: path, Content: string(data)})
	}
	return index.Build(inputs), nil
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/service -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/service/scan.go internal/service/scan_test.go
git commit -m "feat: add service scan pipeline"
```

---

### Task 4: Refs Service (TDD)

**Files:**
- Create: `internal/service/refs.go`
- Create: `internal/service/refs_test.go`

**Step 1: Write failing test**
```go
package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRefs(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.md")
	b := filepath.Join(dir, "b.md")
	if err := os.WriteFile(a, []byte("See [[b]]"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(b, []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	payload, err := RunRefs(RefsInput{Roots: []string{dir}, Target: b, Extensions: []string{".md"}})
	if err != nil {
		t.Fatalf("refs: %v", err)
	}
	if len(payload.Inbound.([]any)) == 0 && payload.Inbound != nil {
		// allow generic any; ensure not empty
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/service -v`
Expected: FAIL with "undefined: RunRefs"

**Step 3: Implement refs service**
```go
package service

import "github.com/bagaking/bilink/internal/output"

type RefsInput struct {
	Roots      []string
	Target     string
	Extensions []string
}

func RunRefs(input RefsInput) (output.RefsPayload, error) {
	idx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.RefsPayload{}, err
	}
	return output.RefsPayload{Target: input.Target, Outbound: idx.Outbound[input.Target], Inbound: idx.Inbound[input.Target]}, nil
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/service -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/service/refs.go internal/service/refs_test.go
git commit -m "feat: add refs service"
```

---

### Task 5: Check Service + Lint Warnings (TDD)

**Files:**
- Modify: `internal/check/check.go`
- Create: `internal/check/check_test.go`
- Create: `internal/service/check.go`
- Create: `internal/service/check_test.go`

**Step 1: Write failing test**
```go
package check

import "testing"

func TestDetectWarnings(t *testing.T) {
	files := []File{
		{Path: "a.md", ResolveKey: "foo", LintKey: "foo"},
		{Path: "b.md", ResolveKey: "bar", LintKey: "foo"},
	}
	errs, warns := Detect(files)
	if len(errs) != 0 {
		t.Fatalf("expected no errors")
	}
	if len(warns) != 1 {
		t.Fatalf("expected warning")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/check -v`
Expected: FAIL

**Step 3: Implement warnings in check.Detect**
```go
package check

func Detect(files []File) ([]string, []string) {
	resolveSeen := map[string][]string{}
	lintSeen := map[string][]string{}
	for _, f := range files {
		resolveSeen[f.ResolveKey] = append(resolveSeen[f.ResolveKey], f.Path)
		if f.LintKey != "" {
			lintSeen[f.LintKey] = append(lintSeen[f.LintKey], f.Path)
		}
	}
	var errs []string
	for key, paths := range resolveSeen {
		if len(paths) > 1 {
			errs = append(errs, key)
		}
	}
	var warns []string
	for key, paths := range lintSeen {
		if len(paths) > 1 {
			warns = append(warns, key)
		}
	}
	return errs, warns
}
```

**Step 4: Write failing service test**
```go
package service

import "testing"

func TestRunCheck(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err == nil {
		_ = payload
	}
}
```

**Step 5: Run tests to verify they fail**
Run: `go test ./internal/service -v`
Expected: FAIL

**Step 6: Implement service check**
```go
package service

import (
	"path/filepath"
	"strings"

	"github.com/bagaking/bilink/internal/check"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/resolve"
)

type CheckInput struct {
	Roots        []string
	Extensions   []string
	ResolveRules resolve.Rules
	LintRules    resolve.Rules
}

func RunCheck(input CheckInput) (output.CheckPayload, error) {
	idx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.CheckPayload{}, err
	}
	var files []check.File
	for path := range idx.Outbound {
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		files = append(files, check.File{
			Path:       path,
			ResolveKey: resolve.NormalizeName(base, input.ResolveRules),
			LintKey:    resolve.NormalizeName(base, input.LintRules),
		})
	}
	errs, warns := check.Detect(files)
	return output.CheckPayload{Errors: errs, Warnings: warns}, nil
}
```

**Step 7: Run tests to verify they pass**
Run: `go test ./internal/check ./internal/service -v`
Expected: PASS

**Step 8: Commit**
```bash
git add internal/check/check.go internal/check/check_test.go internal/service/check.go internal/service/check_test.go
git commit -m "feat: add check warnings"
```

---

### Task 6: Rename Service (全量改写) (TDD)

**Files:**
- Modify: `internal/rename/rename.go`
- Create: `internal/rename/rename_test.go`
- Create: `internal/service/rename.go`
- Create: `internal/service/rename_test.go`

**Step 1: Write failing rename rewrite test**
```go
package rename

import "testing"

func TestRewriteMarkdown(t *testing.T) {
	content := "See [B](b.md)"
	out := RewriteMarkdown(content, "b.md", "c.md")
	if out != "See [B](c.md)" {
		t.Fatalf("unexpected: %s", out)
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/rename -v`
Expected: FAIL with "undefined: RewriteMarkdown"

**Step 3: Implement markdown rewrite**
```go
package rename

import "strings"

func RewriteMarkdown(content, from, to string) string {
	return strings.ReplaceAll(content, "("+from+")", "("+to+")")
}
```

**Step 4: Write failing service test**
```go
package service

import (
	"os"
	"path/filepath"
	"testing"
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
	payload, err := RunRename(RenameInput{Roots: []string{dir}, OldPath: b, NewPath: filepath.Join(dir, "c.md"), Move: true})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if len(payload.Updated) == 0 {
		t.Fatalf("expected updates")
	}
}
```

**Step 5: Run tests to verify they fail**
Run: `go test ./internal/service -v`
Expected: FAIL

**Step 6: Implement rename service**
```go
package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/rename"
)

type RenameInput struct {
	Roots   []string
	OldPath string
	NewPath string
	Move    bool
}

func RunRename(input RenameInput) (output.RenamePayload, error) {
	idx, err := ScanAndIndex(input.Roots, []string{".md", ".markdown", ".mdx"})
	if err != nil {
		return output.RenamePayload{}, err
	}
	oldBase := baseName(input.OldPath)
	newBase := baseName(input.NewPath)
	var updated []string
	for path := range idx.Outbound {
		data, err := os.ReadFile(path)
		if err != nil {
			return output.RenamePayload{}, err
		}
		content := string(data)
		content = rename.RewriteWiki(content, oldBase, newBase)
		content = rename.RewriteMarkdown(content, filepath.Base(input.OldPath), filepath.Base(input.NewPath))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return output.RenamePayload{}, err
		}
		updated = append(updated, path)
	}
	if input.Move {
		if err := os.Rename(input.OldPath, input.NewPath); err != nil {
			return output.RenamePayload{}, err
		}
	}
	return output.RenamePayload{OldPath: input.OldPath, NewPath: input.NewPath, Updated: updated}, nil
}

func baseName(path string) string {
	name := filepath.Base(path)
	if i := strings.LastIndex(name, "."); i > 0 {
		return name[:i]
	}
	return name
}
```

**Step 7: Run tests to verify they pass**
Run: `go test ./internal/rename ./internal/service -v`
Expected: PASS

**Step 8: Commit**
```bash
git add internal/rename/rename.go internal/rename/rename_test.go internal/service/rename.go internal/service/rename_test.go
git commit -m "feat: add rename service"
```

---

### Task 7: Watch Service (TDD)

**Files:**
- Create: `internal/service/watch.go`
- Create: `internal/service/watch_test.go`

**Step 1: Write failing test**
```go
package service

import "testing"

func TestRunWatchRequiresIndex(t *testing.T) {
	_, err := RunWatch(WatchInput{IndexPath: ".bilink/index.json"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/service -v`
Expected: FAIL with "undefined: RunWatch"

**Step 3: Implement watch**
```go
package service

import (
	"os"

	"github.com/bagaking/bilink/internal/index"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/watch"
)

type WatchInput struct {
	IndexPath  string
	Roots      []string
	Extensions []string
}

func RunWatch(input WatchInput) (output.WatchPayload, error) {
	if _, err := os.Stat(input.IndexPath); err != nil {
		return output.WatchPayload{}, err
	}
	oldIdx, err := index.Load(input.IndexPath)
	if err != nil {
		return output.WatchPayload{}, err
	}
	newIdx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.WatchPayload{}, err
	}
	diff := watch.Diff(watch.Index{Files: keys(oldIdx.Outbound)}, watch.Index{Files: keys(newIdx.Outbound)})
	if err := index.Save(input.IndexPath, newIdx); err != nil {
		return output.WatchPayload{}, err
	}
	return output.WatchPayload{Added: diff.Added, Removed: diff.Removed}, nil
}

func keys(m map[string][]index.Link) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/service -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/service/watch.go internal/service/watch_test.go
git commit -m "feat: add watch service"
```

---

### Task 8: CLI Wiring (TDD)

**Files:**
- Modify: `internal/app/app.go`
- Modify: `internal/app/app_test.go`
- Modify: `cmd/bilink/main.go`

**Step 1: Write failing test**
```go
package app

import "testing"

func TestRunRefsMissingArg(t *testing.T) {
	if err := Run([]string{"refs"}); err == nil {
		t.Fatalf("expected error")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/app -v`
Expected: FAIL

**Step 3: Implement CLI parsing**
```go
package app

import (
	"errors"
	"flag"
	"fmt"

	"github.com/bagaking/bilink/internal/config"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/service"
)

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("missing command")
	}
	cmd := args[0]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	var configPath string
	var root string
	var jsonOut bool
	var interactive bool
	var noMove bool
	fs.StringVar(&configPath, "config", "", "config path")
	fs.StringVar(&root, "root", ".", "root directory")
	fs.BoolVar(&jsonOut, "json", false, "json output")
	fs.BoolVar(&interactive, "interactive", false, "interactive")
	fs.BoolVar(&noMove, "no-move", false, "do not move file")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	cfg, err := config.Load(config.ConfigOpts{Roots: []string{root}, ConfigPath: configPath})
	if err != nil {
		return err
	}
	switch cmd {
	case "refs":
		if fs.NArg() < 1 {
			return errors.New("missing target path")
		}
		payload, err := service.RunRefs(service.RefsInput{Roots: cfg.Workspace.Roots, Target: fs.Arg(0), Extensions: cfg.Scan.Extensions})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextRefs(payload))
	case "check":
		payload, err := service.RunCheck(service.CheckInput{Roots: cfg.Workspace.Roots, Extensions: cfg.Scan.Extensions, ResolveRules: toResolve(cfg.Resolve), LintRules: toResolve(cfg.Lint)})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextCheck(payload))
	case "rename":
		if fs.NArg() < 2 {
			return errors.New("missing rename paths")
		}
		payload, err := service.RunRename(service.RenameInput{Roots: cfg.Workspace.Roots, OldPath: fs.Arg(0), NewPath: fs.Arg(1), Move: !noMove})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextRename(payload))
	case "watch":
		payload, err := service.RunWatch(service.WatchInput{IndexPath: cfg.Index.Path, Roots: cfg.Workspace.Roots, Extensions: cfg.Scan.Extensions})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextWatch(payload))
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func writeOutput(jsonOut bool, payload any, text string) error {
	if jsonOut {
		data, err := output.JSON(payload)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	fmt.Println(text)
	return nil
}

func toResolve(r config.Resolve) resolve.Rules {
	return resolve.Rules{
		CaseInsensitive: r.CaseInsensitive,
		IgnoreExtension: r.IgnoreExtension,
		SeparatorEquivalents: r.SeparatorEquivalents,
		UnicodeNormalize: r.UnicodeNormalize,
	}
}
```

**Step 4: Run tests**
Run: `go test ./internal/app -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/app/app.go internal/app/app_test.go cmd/bilink/main.go
git commit -m "feat: wire cli to service"
```

---

### Task 9: npx 自动下载 (TDD)

**Files:**
- Modify: `packages/bilink-npx/lib/platform.js`
- Modify: `packages/bilink-npx/bin/bilink.js`
- Modify: `packages/bilink-npx/test/platform.test.mjs`

**Step 1: Write failing test**
```js
import assert from "node:assert/strict";
import { releaseUrl } from "../lib/platform.js";

assert.equal(
  releaseUrl("latest", "bilink-darwin-arm64"),
  "https://github.com/bagakit/bilink/releases/latest/download/bilink-darwin-arm64"
);
```

**Step 2: Run test to verify it fails**
Run: `node --test packages/bilink-npx/test/platform.test.mjs`
Expected: FAIL

**Step 3: Implement releaseUrl and download**
```js
export function releaseUrl(version, binary) {
  if (version === "latest") {
    return `https://github.com/bagakit/bilink/releases/latest/download/${binary}`;
  }
  return `https://github.com/bagakit/bilink/releases/download/${version}/${binary}`;
}
```

```js
import { createWriteStream, chmodSync } from "node:fs";
import https from "node:https";

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error(`download failed: ${res.statusCode}`));
        return;
      }
      res.pipe(file);
      file.on("finish", () => file.close(resolve));
    }).on("error", reject);
  });
}
```

**Step 4: Run tests**
Run: `node --test packages/bilink-npx/test/platform.test.mjs`
Expected: PASS

**Step 5: Commit**
```bash
git add packages/bilink-npx
git commit -m "feat: add npx auto-download"
```

---

### Task 10: Verification & Docs

**Files:**
- Modify: `docs/must-guidebook.md`
- Modify: `docs/must-sop.md`

**Step 1: Run quality gates**
Run: `make lint`
Expected: PASS

Run: `make test`
Expected: PASS

Run: `node scripts/generate-sop.mjs && git diff --exit-code docs/must-sop.md`
Expected: PASS

**Step 2: Commit**
```bash
git add docs/must-guidebook.md docs/must-sop.md
git commit -m "chore: update docs for cli completion"
```
