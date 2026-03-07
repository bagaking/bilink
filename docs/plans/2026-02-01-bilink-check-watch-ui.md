---
title: Bilink Check 明细输出 + Watch 赛博朋克 TUI
required: false
sop:
  - Read this plan before implementing check detail output or watch TUI styling.
  - Update this plan when steps or files change.
  - Refresh shared system pages after updating this doc.
---
# Bilink Check 明细输出 + Watch 赛博朋克 TUI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `check` 增加分组明细与 JSON 结构化输出，同时升级 `watch` TUI 为赛博朋克风格并支持最近变更列表展示。

**Architecture:** `check` 明细通过 `internal/check` 输出分组结果，`internal/service` 组装 payload，`internal/output` 统一文本/JSON。`watch` 使用 Bubble Tea + lipgloss 进行布局与配色，保留 ASK 动画与配置切换。

**Tech Stack:** Go 1.24+, Bubble Tea, lipgloss.

---

### Task 1: Check 分组检测 (TDD)

**Files:**
- Modify: `internal/check/check.go`
- Modify: `internal/check/check_test.go`

**Step 1: Write failing test**
```go
package check

import "testing"

func TestDetectGroups(t *testing.T) {
	files := []File{
		{Path: "a.md", ResolveKey: "foo", LintKey: "foo"},
		{Path: "b.md", ResolveKey: "foo", LintKey: "bar"},
		{Path: "c.md", ResolveKey: "baz", LintKey: "bar"},
	}
	errs, warns := Detect(files)
	if len(errs) != 1 || errs[0].Key != "foo" || len(errs[0].Paths) != 2 {
		t.Fatalf("unexpected errors")
	}
	if len(warns) != 1 || warns[0].Key != "bar" || len(warns[0].Paths) != 2 {
		t.Fatalf("unexpected warnings")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/check -run TestDetectGroups -v`
Expected: FAIL

**Step 3: Implement group detection**
```go
package check

type Group struct {
	Key   string
	Paths []string
}

func Detect(files []File) ([]Group, []Group) {
	resolveSeen := map[string][]string{}
	lintSeen := map[string][]string{}
	for _, f := range files {
		resolveSeen[f.ResolveKey] = append(resolveSeen[f.ResolveKey], f.Path)
		if f.LintKey != "" {
			lintSeen[f.LintKey] = append(lintSeen[f.LintKey], f.Path)
		}
	}
	var errs []Group
	for key, paths := range resolveSeen {
		if len(paths) > 1 {
			errs = append(errs, Group{Key: key, Paths: paths})
		}
	}
	var warns []Group
	for key, paths := range lintSeen {
		if len(paths) > 1 {
			warns = append(warns, Group{Key: key, Paths: paths})
		}
	}
	return errs, warns
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/check -run TestDetectGroups -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/check/check.go internal/check/check_test.go
git commit -m "feat: add check group detection"
```

---

### Task 2: Check 输出结构化 (TDD)

**Files:**
- Modify: `internal/output/output.go`
- Modify: `internal/output/output_test.go`

**Step 1: Write failing test**
```go
package output

import "testing"

func TestTextCheckIncludesGroups(t *testing.T) {
	payload := CheckPayload{
		Errors: []string{"foo"},
		ErrorGroups: []CheckGroup{{Key: "foo", Paths: []string{"a.md", "b.md"}}},
	}
	text := TextCheck(payload)
	if text == "" || !strings.Contains(text, "foo") {
		t.Fatalf("expected grouped output")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/output -run TestTextCheckIncludesGroups -v`
Expected: FAIL

**Step 3: Implement payload + grouped output**
```go
package output

type CheckGroup struct {
	Key   string   `json:"key"`
	Paths []string `json:"paths"`
}

type CheckPayload struct {
	Errors        []string     `json:"errors"`
	Warnings      []string     `json:"warnings"`
	ErrorGroups   []CheckGroup `json:"errorGroups"`
	WarningGroups []CheckGroup `json:"warningGroups"`
}

func TextCheck(p CheckPayload) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("errors: %d\nwarnings: %d\n", len(p.Errors), len(p.Warnings)))
	if len(p.ErrorGroups) > 0 {
		b.WriteString("\nerror groups:\n")
		for _, g := range p.ErrorGroups {
			b.WriteString(fmt.Sprintf("- %s\n", g.Key))
			for _, path := range g.Paths {
				b.WriteString(fmt.Sprintf("  - %s\n", path))
			}
		}
	}
	if len(p.WarningGroups) > 0 {
		b.WriteString("\nwarning groups:\n")
		for _, g := range p.WarningGroups {
			b.WriteString(fmt.Sprintf("- %s\n", g.Key))
			for _, path := range g.Paths {
				b.WriteString(fmt.Sprintf("  - %s\n", path))
			}
		}
	}
	return strings.TrimRight(b.String(), "\n")
}
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/output -run TestTextCheckIncludesGroups -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/output/output.go internal/output/output_test.go
git commit -m "feat: add check grouped output"
```

---

### Task 3: Service Check 组装分组 (TDD)

**Files:**
- Modify: `internal/service/check.go`
- Modify: `internal/service/check_test.go`

**Step 1: Write failing test**
```go
package service

import "testing"

func TestRunCheckGroups(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	_ = payload.ErrorGroups
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/service -run TestRunCheckGroups -v`
Expected: FAIL

**Step 3: Implement group mapping**
```go
errs, warns := check.Detect(files)
return output.CheckPayload{
	Errors: keys(errs),
	Warnings: keys(warns),
	ErrorGroups: toGroups(errs),
	WarningGroups: toGroups(warns),
}, nil
```

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/service -run TestRunCheckGroups -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/service/check.go internal/service/check_test.go
git commit -m "feat: add check groups to payload"
```

---

### Task 4: Watch TUI 赛博朋克布局 (TDD)

**Files:**
- Modify: `internal/watch/tui/model.go`
- Modify: `internal/watch/tui/model_test.go`
- Modify: `go.mod`, `go.sum`

**Step 1: Write failing test**
```go
package tui

import (
	"strings"
	"testing"

	"github.com/bagaking/bilink/internal/output"
)

func TestModelViewCyberpunk(t *testing.T) {
	m := NewModel(output.WatchPayload{Added: []string{"a.md"}, Removed: []string{"b.md"}}, "config")
	view := m.View()
	if !strings.Contains(view, "BILINK WATCH") || !strings.Contains(view, "ASK") {
		t.Fatalf("expected title and ask")
	}
	if !strings.Contains(view, "+ a.md") || !strings.Contains(view, "- b.md") {
		t.Fatalf("expected change list")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./internal/watch/tui -run TestModelViewCyberpunk -v`
Expected: FAIL

**Step 3: Implement layout with lipgloss**
- 引入 `github.com/charmbracelet/lipgloss`
- 标题区：ASCII 标题 + 霓虹渐变色
- 统计区：added/removed/errors/warnings
- 事件区：最近 N 条变更 `+ file` / `- file`
- 底栏：快捷键提示 + ASK 动画

**Step 4: Run tests to verify they pass**
Run: `go test ./internal/watch/tui -run TestModelViewCyberpunk -v`
Expected: PASS

**Step 5: Commit**
```bash
git add internal/watch/tui/model.go internal/watch/tui/model_test.go go.mod go.sum
git commit -m "feat: style watch tui"
```

---

### Task 5: README 更新

**Files:**
- Modify: `README.md`

**Step 1: Update README**
- 说明 `check` 明细输出格式
- 说明 `watch` TUI 新布局与变更列表

**Step 2: Commit**
```bash
git add README.md
git commit -m "docs: update readme for check/watch"
```

---

### Task 6: Verification

Run:
- `make lint`
- `make test`
- `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root . && git diff --exit-code docs/must-guidebook.md docs/must-sop.md`

Expected: PASS
