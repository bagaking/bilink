# Project SOP

This SOP is generated from docs frontmatter. Do not edit manually.

## Update Requirements
- When a document with SOP frontmatter changes, regenerate this file with `node scripts/generate-sop.mjs` and commit the result.
- Add new SOP items by updating the `sop` list in the source document frontmatter.
- Keep SOP items small and actionable; use the source document for details.

## SOP Items
### Bilink 设计文档（核心理念与架构）
Source: `docs/architecture-bilink-design.md`
- Read this doc before making architectural or behavioral changes.
- Update this doc when core design principles or data flow change.
- Regenerate must-sop.md after updating this doc.

### Bilink PRD（MVP + Watch TUI）
Source: `docs/notes-prd-bilink.md`
- Read this doc before changing product scope or user-facing behavior.
- Update this doc when requirements or goals change.
- Regenerate must-sop.md after updating this doc.

### Bilink MVP Implementation Plan
Source: `docs/plans/2026-01-27-bilink-mvp.md`
- Read this plan before implementing Bilink MVP tasks.
- Update this plan when execution steps change.
- Regenerate must-sop.md after updating this doc.

### Bilink CLI 完整实现设计
Source: `docs/plans/2026-01-30-bilink-cli-design.md`
- Read this doc before implementing service layer or CLI wiring changes.
- Update this doc when command orchestration or output formats change.
- Regenerate must-sop.md after updating this doc.

### Bilink CLI 完整实现计划
Source: `docs/plans/2026-01-30-bilink-cli-implementation.md`
- Read this plan before implementing Bilink CLI completion tasks.
- Update this plan when steps or files change.
- Regenerate must-sop.md after updating this doc.
