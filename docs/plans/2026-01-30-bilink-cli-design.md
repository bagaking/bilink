---
title: Bilink CLI 完整实现设计
required: false
sop:
  - Read this doc before implementing service layer or CLI wiring changes.
  - Update this doc when command orchestration or output formats change.
  - Refresh shared system pages after updating this doc.
---
# Bilink CLI 完整实现设计

## Context
当前实现包含解析、索引、重命名等核心模块，但 CLI 只提供命令分发骨架。需要新增 `internal/service` 作为应用层，把配置、扫描、索引、规则与输出串起来，满足 PRD 与 OpenSpec 对 `refs/rename/check/watch/--json` 的要求，并补全 npx 自动下载。

## Goals / Non-Goals
- Goals:
  - 引入 `internal/service` 统一编排命令行为与数据流。
  - CLI 变薄：只做参数解析与调用 service。
  - 文本/JSON 输出由 `internal/output` 统一格式化。
  - npx wrapper 从 GitHub Releases 自动下载二进制（命名 `bilink-<os>-<arch>`）。
- Non-Goals:
  - 重做底层解析/索引/归一化逻辑。
  - 自动改写 anchor；仅解析跳转。

## Architecture
- **service 层**：按用例拆分（RefsService/CheckService/RenameService/WatchService），共享 `LoadConfig/ScanAndIndex/Resolve/FormatOutput`。
- **output 层**：统一输出结构（text + json）。
- **CLI 层**：解析 flags → 构造 `ServiceConfig` → 调用 service。
- **watch 层**：service 产出事件流，TUI 订阅并展示。

## Data Flow (每个命令)
1) Load config（--config 或默认 `.bilink/settings.toml`）
2) Scan roots → Read file contents → Parse links/anchors
3) Build index → Apply resolve/lint rules
4) Command-specific action（refs/check/rename/watch）
5) Output formatter（text/json）

## Error Handling & UX
- 强关联且唯一命中才自动改写；歧义默认 error，`--interactive` 进入选择。
- `check`：resolve 冲突 error，lint 违规 warning。
- `watch`：缺少 `.bilink/index.json` 直接报错并提示初始化。

## Testing
- service 层新增集成测试覆盖 refs/check/rename JSON 输出与歧义行为。
- CLI 仅测参数解析与错误码。
- 质量门禁：`make lint`、`make test`、SOP 生成检查。

## Release & npx
- Releases 资产命名：`bilink-<os>-<arch>`，Windows 追加 `.exe`。
- npx 包自动下载到 `~/.cache/bilink` 并转发 CLI 参数。
