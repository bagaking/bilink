---
title: Bilink 设计文档（核心理念与架构）
required: true
sop:
  - Read this doc before making architectural or behavioral changes.
  - Update this doc when core design principles or data flow change.
  - Regenerate must-sop.md after updating this doc.
---
# Bilink 设计文档（核心理念与架构）

## Context
Bilink 面向本地文件系统的双向链接与引用维护，覆盖任意目录集合（森林），支持 `[[...]]` 与 `[text](path)` 解析、引用查询、重命名更新、冲突检测与 watch 交互界面。

## Goals / Non-Goals
- Goals: 强/弱关联分层；强关联自动改写、弱关联告警；`#anchor` 仅用于跳转解析；可选索引快照；高可解释性输出。
- Non-Goals: 块引用 `^blockid`；非 Markdown 扩展默认支持；`#anchor` 自动改写；远程/云端同步。

## 核心设计理念（原汁原味）
> 这里的设计哲学是我们认为名字后面的 # 各种插件都视为是补充协议, 我们现在就识别出来, 确实方便做一些扩展插件.
>
> 但是这不影响 block first 的设计. 也就是说, 牵引性是避免复杂文件, 尽量保持简单小文件为主.
>
> 这样其实也更加 AI Friendly, 变相要求用户或 Agent 写入时已经有按逻辑拆分和用名字建立索引的行为, 而不是把索引和摘要全都嵌入在文件正文里
>
> “block first, 一个文件就是一个 block”.

## Architecture Overview
系统由三层组成：
1) **解析与索引层**：扫描目录集合，解析 Markdown 链接与锚点，构建 in-memory 索引；可选落盘 `.bilink/index.json`。
2) **规则与更新层**：`resolveRules` 控制解析范围与匹配；`lintRules` 负责规范提示；更新策略仅在强关联且无歧义时自动改写。
3) **交互与分发层**：CLI 子命令（refs/rename/check/watch）与 `--json` 输出；watch 使用 Bubble Tea 提供交互式 TUI；TS 包装器支持 `npx` 使用。

## Link Model
- **强关联**：`[[file]]`、`[[file|alias]]`、`[text](path)` 等显式文件引用，目标唯一可自动改写。
- **弱关联**：通过 title/aliases/heading 推断或需要更宽松等价化的匹配，仅告警或交互确认。
- **锚点**：仅用于跳转解析；发生变更时提醒，不自动改写。

## Data Flow
1) 扫描目录集合 → 过滤扩展名（默认 `.md/.markdown/.mdx`）。
2) 解析文件 → 抽取链接与锚点 → 归一化路径 → 生成索引。
3) 命令执行（refs/rename/check/watch）→ 应用规则 → 输出结果或执行更新。

## Error Handling & UX
- 歧义默认报错并列出候选；`--interactive` 允许选择。
- 解析尽量高召回；写入操作高精度，避免误改。
- 提供 `--json` 结构化输出，方便脚本与 CI 集成。

## Configuration
- 默认读取 `.bilink/settings.toml`，可用 `--config` 覆盖。
- 允许定义 resolveRules、lintRules 与扩展名白名单。
- 无配置文件时使用内置默认值，不自动创建。

### Default Settings (TOML)
```
[workspace]
roots = ["."]

[scan]
extensions = [".md", ".markdown", ".mdx"]

[resolveRules]
caseInsensitive = true
ignoreExtension = true
separatorEquivalents = [" ", "-", "_"]
unicodeNormalize = "NFKC"

[lintRules]
requireExactCase = true
requireExplicitExtension = true
requireExactSeparators = true

[anchors]
style = "github"
mode = "resolve-only"

[updatePolicy]
mode = "balanced"

[index]
path = ".bilink/index.json"
requiredForWatch = true
```

## Confidence Tiers (Balanced)
- **High**: 强关联且目标唯一，只使用“核心归一化”（路径清理、大小写折叠、缺省扩展名补全）。
- **Medium**: 仅在应用等价规则后才唯一命中（空格/下划线/短横线等价、Unicode 归一化）。
- **Low**: 由内容推断（title/aliases/heading）或存在歧义。

Balanced 策略：High 自动改写；Medium 自动改写并 warning；Low 仅提示。

## Indexing Strategy
- 默认全量扫描构建临时索引；`watch` 模式依赖 `.bilink/index.json` 快照以对比变更。
- 无快照时作为静态分析工具使用。

## Testing & Verification
- 质量门禁：`make lint`、`make test`、`node scripts/generate-sop.mjs && git diff --exit-code docs/must-sop.md`。

## Future Considerations
- 可配置锚点规则（GitHub/Obsidian）。
- 远程/云端同步作为未来扩展能力。
- 远程同步若引入：优先同步 `.bilink/index.json` 元数据快照，再通过本地全量扫描与冲突提示进行修复。
