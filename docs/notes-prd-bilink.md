---
title: Bilink PRD（MVP + Watch TUI）
required: true
sop:
  - Read this doc before changing product scope or user-facing behavior.
  - Update this doc when requirements or goals change.
  - Regenerate must-sop.md after updating this doc.
---
# PRD: Bilink 双向链接工具（MVP + Watch TUI）

## Overview
Bilink 是一个面向本地文件系统的双向链接工具，支持在一组目录“森林”中扫描并解析 Markdown 内的 `[[...]]` 与 `[text](path)` 链接，提供引用查询、重命名更新、冲突检测与 Watch 交互界面。系统区分强关联与弱关联：强关联可自动更新，弱关联仅提示或交互确认。`#anchor` 作为补充协议，仅用于跳转解析，发生变更时提醒但不自动改写。

## Goals
- 在任意目录集合中快速找到出链/反链。
- 支持 `[[...]]` 与 `[text](path)` 的强/弱关联识别与更新。
- 提供 `rename` 自动更新引用，歧义时交互选择。
- 提供 `check` 进行全局唯一性与冲突检测（error/warn）。
- 提供 `watch` TUI：进度/状态、交互选择、配置查询、酷炫 ask 动画。
- 可选 `.bilink/index.json` 索引；无索引时也可作为静态分析工具使用。
- 配置文件默认读取 `.bilink/settings.toml`，可用 `--config` 覆盖。

## Quality Gates
These commands must pass for every user story:
- `make lint`
- `make test`
- `node scripts/generate-sop.mjs && git diff --exit-code docs/must-sop.md`

Manual TUI verification is not required.

## User Stories

### US-001: 扫描与索引
**Description:** 作为用户，我希望在给定目录集合中扫描 Markdown 链接并建立索引，以便后续查询与更新。

**Acceptance Criteria:**
- [ ] 解析 `[[...]]` 与 `[text](path)` 两类链接。
- [ ] 支持配置文件扩展名白名单（默认 `.md/.markdown/.mdx`）。
- [ ] 提供 `--json` 输出，返回解析到的链接与分类信息。

### US-002: 引用查询（refs）
**Description:** 作为用户，我希望查询任意文件的出链/反链，以便快速定位引用关系。

**Acceptance Criteria:**
- [ ] `refs` 支持输入单文件路径并输出出链/反链列表。
- [ ] 输出标明强/弱关联与歧义候选。
- [ ] `--json` 输出包含可解析的结构化字段。

### US-003: 重命名/移动更新（rename）
**Description:** 作为用户，我希望重命名或移动文件时自动更新相关引用，且歧义时可交互选择。

**Acceptance Criteria:**
- [ ] 支持 `--move/--no-move`，默认移动文件并更新引用。
- [ ] 仅在目标唯一且强关联时自动改写。
- [ ] 歧义时默认报错；`--interactive` 提供选择并继续更新。

### US-004: 冲突检测（check）
**Description:** 作为用户，我希望检测全局唯一性冲突与命名规范问题，以便保持链接稳定。

**Acceptance Criteria:**
- [ ] `check` 基于 resolveRules 发现冲突并以 error 输出。
- [ ] 基于 lintRules 输出 warning，但不阻断。
- [ ] 输出包含冲突候选与建议修复方向。

### US-005: 配置与索引持久化
**Description:** 作为用户，我希望通过配置文件与可选索引来控制规则与性能。

**Acceptance Criteria:**
- [ ] CLI 支持 `--config`；默认读取 `.bilink/settings.toml`（不存在则使用默认值）。
- [ ] `watch` 模式要求 `.bilink/index.json` 可用；静态模式不强制。
- [ ] 支持初始化目录（创建 `.bilink/` 与可选配置模板）。

### US-006: Watch TUI
**Description:** 作为用户，我希望在 watch 模式下获得可交互界面以处理变更与冲突。

**Acceptance Criteria:**
- [ ] TUI 展示索引进度、变更数量、错误/告警。
- [ ] 支持交互处理歧义与确认改写。
- [ ] 支持查询当前配置与规则状态。
- [ ] 包含“ask”动画与明显的交互反馈。

### US-007: TS 包装器（npx）
**Description:** 作为用户，我希望通过 `npx` 直接运行 bilink，无需手动安装 Go 二进制。

**Acceptance Criteria:**
- [ ] NPM 包安装时自动下载对应平台二进制。
- [ ] CLI 参数原样转发到 Go 程序。

## Functional Requirements
- FR-1: 系统必须支持 `[[...]]` 与 `[text](path)` 的解析。
- FR-2: 强关联必须在无歧义时自动更新；弱关联仅提示或交互确认。
- FR-3: `#anchor` 仅用于跳转解析；发生变更时只告警。
- FR-4: 配置支持 resolveRules 与 lintRules（分别控制解析与提示）。
- FR-5: CLI 提供 `refs/rename/check/watch` 与 `--json` 输出。
- FR-6: 允许传入多个目录作为“森林”并合并索引。
- FR-7: 默认仅处理 Markdown 扩展名，支持用户自定义。
- FR-8: 默认锚点规则为 GitHub，并允许配置为 Obsidian。
- FR-9: 默认更新策略为 Balanced（高/中置信度自动改写，低置信度仅提示）。

## Non-Goals (Out of Scope)
- `^blockid`/块引用语义。
- 非 Markdown 扩展名默认支持。
- `#anchor` 自动改写。
- 远程/云端同步。

## Technical Considerations
- Go 实现核心；Bubble Tea 作为 TUI 框架。
- 解析与改写分层：解析追求召回，写入追求精确。
- `.bilink/index.json` 为可选持久化快照；watch 依赖索引 diff。
- `resolveRules` 用于解析；`lintRules` 用于规范提示。
- 默认锚点规则为 GitHub（可配置为 Obsidian），仅用于跳转解析。
- Balanced 更新策略建议：高置信度自动改写，中置信度自动改写并 warning，低置信度仅提示。

## Success Metrics
- 重命名更新无误改案例（precision 优先）。
- `check` 能稳定发现冲突与歧义。
- watch 模式在大目录下保持稳定且响应迅速。

## Open Questions
- 暂无。
