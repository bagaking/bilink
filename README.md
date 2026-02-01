# Bilink

Bilink 是一个面向本地文件系统的双向链接工具，支持在目录“森林”中扫描 Markdown 链接，提供引用查询、冲突检查、重命名更新与 watch TUI。

## 功能（当前实现）
- 解析 `[[...]]` 与 `[text](path)` 两类链接（忽略 `http/https/mailto` 外链）。
- `refs` 查看出链/反链；`check` 输出冲突（error）与 lint（warning）明细分组。
- `rename` 支持改写引用并可移动文件，歧义时默认报错。
- `watch` 赛博朋克风格 TUI（ASK 动画 + 配置面板 + 最近变更列表），支持 `--json` 结构化输出。
- 默认扫描 `.md/.markdown/.mdx`，可通过配置扩展。

> 说明：`--interactive` 当前仅跳过歧义拦截（不提供交互选择器）。

## 安装

### Go
```bash
go build -o bilink ./cmd/bilink
```

### npx
```bash
npx bilink refs path/to/file.md
```
要求 GitHub Releases 产物命名为：`bilink-<os>-<arch>`（Windows 追加 `.exe`）。

## 快速使用
```bash
./bilink refs path/to/file.md
./bilink check
./bilink rename old.md new.md
./bilink watch
```
常用参数：
- `--root <dir>`：扫描根目录（单根）
- `--config <path>`：指定配置文件
- `--json`：输出 JSON（check 包含 errorGroups/warningGroups）
- `--interactive`：歧义时允许继续（当前无交互选择器）
- `--no-move`：只改引用不移动文件

## check 输出说明
- 文本输出会按 key 分组列出冲突文件路径。
- JSON 输出包含：
  - `errors` / `warnings`（key 列表）
  - `errorGroups` / `warningGroups`（key + paths 明细）

## 配置
默认读取 `.bilink/settings.toml`（不存在则使用内置默认值）。示例：
```toml
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

## watch 模式
`watch` 需要 `.bilink/index.json`：
```bash
mkdir -p .bilink
cat <<JSON > .bilink/index.json
{"outbound":{},"inbound":{}}
JSON
```
TUI 快捷键：
- `y/Enter` 确认
- `c` 显示配置摘要
- `q` 退出

## 开发
```bash
make lint
make test
```
