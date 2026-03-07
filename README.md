# Bilink

Bilink 是一个面向本地文件系统的 Markdown 双向链接工具。它的目标不是做笔记应用，而是把一组目录里的 Markdown 文件变成一个可扫描、可查询、可校验、可维护的链接系统。

当前仓库主要包含：
- Go CLI 核心实现
- `watch` TUI
- 一个 Node 包装器原型，位于 `packages/bilink-npx/`
- 一套基于 `bagakit-living-knowledge` 的共享文档基座

## 项目目标

Bilink 想解决 4 类问题：
- 扫描 Markdown 目录森林，识别 `[[...]]` 与 `[text](path)` 两类链接，建立出链和反链。
- 在文件重命名或移动时，自动更新足够确定的强关联引用，尽量避免误改。
- 对命名冲突、解析歧义和规范问题做全局检查。
- 用 `watch` 界面把索引变化、状态和待确认信息展示出来。

设计取向很明确：
- 解析阶段追求高召回。
- 改写阶段追求高精度。
- `#anchor` 只做跳转解析，不自动改写。
- “block first，一个文件就是一个 block”。

更完整的项目目标和约束见：
- [docs/notes-prd-bilink.md](docs/notes-prd-bilink.md)
- [docs/architecture-bilink-design.md](docs/architecture-bilink-design.md)

## 当前状态

当前实现已经具备基础命令和主要文档，但还不是一个完全收口的稳定版。仓库内可以直接开发和验证的部分是：
- `refs`
- `check`
- `rename`
- `watch`
- JSON 输出
- 基础 TUI

当前行为里有几个需要明确知道的边界：
- `--interactive` 目前只会跳过歧义拦截，不提供交互选择器。
- `watch` 当前更接近“一次 index diff + TUI 展示”，不是长期驻留的文件系统监听器。
- 配置能力已经有，但当前使用时建议显式传 `--config`，不要假设默认配置路径一定会被读取。
- Node 包装器代码已经在仓库里，但它依赖匹配的 GitHub Releases 资产，不应把它当成已验证发布链路。

## 快速开始

最直接的使用方式是本地构建 Go CLI：

```bash
go build -o bilink ./cmd/bilink
```

准备一个最小配置：

```bash
mkdir -p .bilink
cat > .bilink/settings.toml <<'EOF'
[workspace]
roots = ["."]

[scan]
extensions = [".md", ".markdown", ".mdx"]
EOF
```

然后显式传配置运行：

```bash
./bilink refs --root . --config .bilink/settings.toml path/to/file.md
./bilink check --root . --config .bilink/settings.toml
./bilink rename --root . --config .bilink/settings.toml old.md new.md
./bilink watch --root . --config .bilink/settings.toml
```

常用参数：
- `--root <dir>`：扫描根目录
- `--config <path>`：显式指定配置文件
- `--json`：输出结构化 JSON
- `--interactive`：歧义时允许继续，当前不带选择器
- `--no-move`：只改引用，不移动文件

## 命令说明

`refs`
- 输入一个目标文件路径。
- 当前文本输出偏摘要；JSON 输出更适合脚本消费。

`check`
- 输出 resolve 级别的冲突错误和 lint 级别的 warning。
- 文本输出按 key 分组列出路径。
- JSON 输出包含 `errors`、`warnings`、`errorGroups`、`warningGroups`。

`rename`
- 支持改写引用并移动文件。
- 默认要求目标唯一；歧义时直接报错。

`watch`
- 依赖 `.bilink/index.json` 作为前后状态对比基线。
- TUI 快捷键：
  - `y` / `Enter`：确认
  - `c`：显示配置摘要
  - `q`：退出

初始化 `watch` 所需索引文件：

```bash
mkdir -p .bilink
cat > .bilink/index.json <<'EOF'
{"outbound":{},"inbound":{}}
EOF
```

## 配置

一个完整配置示例如下：

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

## 仓库文档入口

这个仓库现在使用 `bagakit-living-knowledge` 做共享知识基座。进入仓库后，优先读这些页面：
- `docs/must-guidebook.md`
- `docs/must-authority.md`
- `docs/must-sop.md`
- `docs/must-recall.md`

如果你要改文档基座本身，再看：
- `docs/specs/living-knowledge-system.md`

## 开发与验证

代码质量门禁：

```bash
make lint
make test
```

共享文档系统页刷新：

```bash
export BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR="<path-to-bagakit-living-knowledge-skill>"
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .
git diff --exit-code docs/must-guidebook.md docs/must-sop.md
```
