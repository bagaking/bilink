# Guidebook

This guidebook is a reading map. Keep it stable and index-style; do not duplicate volatile details.

## How to Use
- First-time: read in order.
- Afterward: jump to the relevant section.
- When in doubt: follow system docs first.

## Fast Path
1) System docs
- `docs/must-docs-taxonomy.md`
- `docs/must-sop.md`

2) Project intent and constraints
- `openspec/project.md`

3) Current changes or proposals
- `openspec/changes/`
- `openspec list` (CLI)

4) Build/run entrypoints
- `Makefile` (lint/build)
- `node scripts/generate-sop.mjs` (SOP regeneration)

## Deep Dives
- Add domain-specific docs by category (norms, guidelines, notes).
- Use taxonomy suffixes and keep ordering consistent.
- `docs/notes-prd-bilink.md`
- `docs/architecture-bilink-design.md`

## Docs Maintenance
- This guidebook must reference `docs/must-docs-taxonomy.md`.
- When docs are added/renamed, update this guidebook.
- Do not move system docs without updating `AGENTS.md`.

## Response Footer
- Every task response must end with `[[Bagakit.LivingDoc]] ...`.
