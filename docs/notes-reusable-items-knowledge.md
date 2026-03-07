---
title: Reusable Items - Knowledge
sop:
  - Update this catalog when one note, index, or query pattern becomes worth reusing across tasks.
  - Keep source-of-truth links current and remove duplicate entries.
---

# Reusable Items - Knowledge

This catalog tracks reusable knowledge assets for the repository.

## Canonical Indexes

| Item | Level | When To Use | Source Of Truth |
| --- | --- | --- | --- |
| Shared Guidebook | MUST | Start any task that depends on repository knowledge loading order or topic discovery | `docs/must-guidebook.md` |
| Shared Authority | MUST | Clarify where durable truth lives before changing shared knowledge surfaces | `docs/must-authority.md` |
| Shared SOP | MUST | Find maintenance-route guidance before editing shared knowledge pages | `docs/must-sop.md` |
| Shared Recall Workflow | SHOULD | Use deterministic search and line-level recall instead of memory | `docs/must-recall.md` |

## High-Signal Notes

| Item | Level | When To Use | Source Of Truth |
| --- | --- | --- | --- |
| Bilink PRD | MUST | Check product scope, user stories, and quality gates before changing user-facing behavior | `docs/notes-prd-bilink.md` |
| Bilink Architecture | MUST | Check core design constraints before refactors or behavior changes | `docs/architecture-bilink-design.md` |
| Living Knowledge System | SHOULD | Change `.bagakit/knowledge_conf.toml`, shared `must-*` pages, or AGENTS bootstrap safely | `docs/specs/living-knowledge-system.md` |

## Reusable Query Patterns

| Item | Level | When To Use | Source Of Truth |
| --- | --- | --- | --- |
| `recall search --root . 'Bilink'` | SHOULD | Find the project’s main product and architecture pages quickly | `docs/must-recall.md` |
| `recall search --root . 'living knowledge'` | SHOULD | Find substrate rules, specs, and reusable-item governance pages | `docs/must-recall.md` |
