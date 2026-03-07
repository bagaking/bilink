---
title: Living Knowledge System
sop:
  - Read this doc before changing `.bagakit/knowledge_conf.toml`, shared `must-*` pages, or the managed `bagakit-living-knowledge` block in `AGENTS.md`.
  - Refresh the substrate with `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" apply --root .`, `index --root .`, and `doctor --root .` after changing shared knowledge surfaces.
---
# Living Knowledge System

## Purpose

This repository uses `bagakit-living-knowledge` as its shared
filesystem-first knowledge substrate.

The goal is to keep durable project knowledge:
- checked in
- progressively loadable through `must-*` pages
- discoverable through deterministic recall
- separate from research, selector runtime, and evolver memory

## Path Protocol

The local path contract is declared in `.bagakit/knowledge_conf.toml`.

Current values:
- `shared_root = docs`
- `system_root = docs`
- `generated_root = .bagakit/living-knowledge/.generated`
- `researcher_root = .bagakit/researcher`
- `selector_root = .bagakit/skill-selector`
- `evolver_root = .bagakit/evolver`

## Shared Surfaces

Checked-in shared knowledge:
- `docs/`

Bootstrap:
- `AGENTS.md`

System pages:
- `docs/must-guidebook.md`
- `docs/must-authority.md`
- `docs/must-sop.md`
- `docs/must-recall.md`

Reusable-items governance:
- `docs/norms-maintaining-reusable-items.md`
- `docs/notes-reusable-items-knowledge.md`

Generated local helper outputs:
- `.bagakit/living-knowledge/.generated/`

## Local Rules

- `AGENTS.md` is only the bootstrap layer; durable shared knowledge belongs under `docs/`.
- `bagakit-living-knowledge` owns path protocol, guidebook indexing, generated `must-sop.md`, and recall helpers.
- `docs/must-docs-taxonomy.md` remains a local naming guide for ordinary pages, not the shared bootstrap authority.
- `openspec/` remains the source of truth for spec-driven change proposals and capability deltas.

## Maintenance Flow

Recommended maintenance sequence:
1. `apply` when the substrate or managed bootstrap needs refreshing.
2. `paths` to confirm the resolved path contract.
3. `index` after shared page or frontmatter changes.
4. `doctor` to catch non-destructive substrate issues.

Commands:

```bash
export BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR="<path-to-bagakit-living-knowledge-skill>"
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" apply --root .
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" paths --root .
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" doctor --root .
```

## Recall Flow

Use recall helpers before answering from memory:

```bash
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" recall search --root . "<query>"
sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" recall get --root . <path> --from <line> --lines <n>
```

## Compatibility Notes

- Legacy `bagakit-living-docs` wording may still appear in historical notes and logs.
- Shared system-page refresh should use the living-knowledge `index` command, not `node scripts/generate-sop.mjs`.
- Path-local `AGENTS.md` files may narrow execution guidance, but they must not redefine the shared knowledge root.
