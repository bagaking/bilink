# SOP

This page is generated from optional frontmatter in shared knowledge pages.
Do not hand-edit it directly.

## Update Rules
- Refresh this page with `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .`.
- Add or update route guidance in page frontmatter under `sop:` and optional `directives:`.
- Keep items short, concrete, and scoped to the source page.

## Sources
- shared root: `docs`
- system root: `docs`

## SOP Items

### Bilink 设计文档（核心理念与架构）
Source: `docs/architecture-bilink-design.md`
- Read this doc before making architectural or behavioral changes.
- Update this doc when core design principles or data flow change.
- Refresh shared system pages after updating this doc.

### Maintaining Reusable Items
Source: `docs/norms-maintaining-reusable-items.md`
- When one pattern, checklist, or example becomes a stable project default, add or update the matching reusable-items catalog entry in the same change.
- When reusable-items route guidance changes, refresh `docs/must-sop.md` by running `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .`.

### Bilink PRD（MVP + Watch TUI）
Source: `docs/notes-prd-bilink.md`
- Read this doc before changing product scope or user-facing behavior.
- Update this doc when requirements or goals change.
- Refresh shared system pages after updating this doc.

### Reusable Items - Knowledge
Source: `docs/notes-reusable-items-knowledge.md`
- Update this catalog when one note, index, or query pattern becomes worth reusing across tasks.
- Keep source-of-truth links current and remove duplicate entries.

### Bilink MVP Implementation Plan
Source: `docs/plans/2026-01-27-bilink-mvp.md`
- Read this plan before implementing Bilink MVP tasks.
- Update this plan when execution steps change.
- Refresh shared system pages after updating this doc.

### Bilink CLI 完整实现设计
Source: `docs/plans/2026-01-30-bilink-cli-design.md`
- Read this doc before implementing service layer or CLI wiring changes.
- Update this doc when command orchestration or output formats change.
- Refresh shared system pages after updating this doc.

### Bilink CLI 完整实现计划
Source: `docs/plans/2026-01-30-bilink-cli-implementation.md`
- Read this plan before implementing Bilink CLI completion tasks.
- Update this plan when steps or files change.
- Refresh shared system pages after updating this doc.

### Bilink Check 明细输出 + Watch 赛博朋克 TUI
Source: `docs/plans/2026-02-01-bilink-check-watch-ui.md`
- Read this plan before implementing check detail output or watch TUI styling.
- Update this plan when steps or files change.
- Refresh shared system pages after updating this doc.

### Living Knowledge System
Source: `docs/specs/living-knowledge-system.md`
- Read this doc before changing `.bagakit/knowledge_conf.toml`, shared `must-*` pages, or the managed `bagakit-living-knowledge` block in `AGENTS.md`.
- Refresh the substrate with `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" apply --root .`, `index --root .`, and `doctor --root .` after changing shared knowledge surfaces.
