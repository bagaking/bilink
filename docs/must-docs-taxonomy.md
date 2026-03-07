# Docs Taxonomy Guidelines

## Scope
This document defines local naming and frontmatter conventions for ordinary
pages under `docs/`.

`bagakit-living-knowledge` owns the shared substrate and system pages:
- `must-guidebook.md`
- `must-authority.md`
- `must-sop.md`
- `must-recall.md`

This page is supplementary. It guides how repository-authored knowledge is
named and grouped inside the shared root, but it is not the bootstrap layer.

## Categories (Non-System Docs)
- Norms: `norms-*.md`
- Architecture: `architecture-*.md`
- Guidelines: `guidelines-*.md`
- Notes: `notes-*.md`
- Plans: `plans/YYYY-MM-DD-<topic>.md`
- Specs: `specs/*.md`
- Runbook (optional): `runbook-*.md`
- Manual test (optional): `manual-test-*.md`

## System Docs
System pages are generated or maintained by the shared knowledge substrate and
use the `must-` prefix.

Required system docs:
- `must-guidebook.md`: reading map and doc index
- `must-authority.md`: authority and path protocol boundaries
- `must-sop.md`: generated SOP output (do not hand-edit)
- `must-recall.md`: recall workflow over shared knowledge

## Naming Rules (Non-System Docs)
Use lowercase kebab-case filenames.

- Top-level pages should keep the category first: `<type>-<topic>.md`.
- Plans should stay in `docs/plans/` and keep the date prefix.
- Stable subsystem contracts should live in `docs/specs/`.

Examples:
- `norms-dev.md`
- `notes-voice-input.md`
- `guidelines-ui-components.md`
- `plans/2026-04-20-docs-reorg.md`
- `specs/living-knowledge-system.md`

## Frontmatter Templates
Non-system docs that need maintenance-route guidance should include frontmatter
with `title` and `sop`. The legacy `required` field may still appear in older
pages, but new pages do not need to introduce it unless another workflow relies
on it.

### Norms
```
---
title: <Norms Title>
sop:
  - Read this doc before touching the affected area.
  - Update this doc when rules or constraints change.
  - Refresh shared system pages after updating this doc.
---
```

### Architecture
```
---
title: <Architecture Title>
sop:
  - Read this doc before architecture changes or refactors.
  - Update this doc when components or boundaries change.
  - Refresh shared system pages after updating this doc.
---
```

### Guidelines
```
---
title: <Guidelines Title>
sop:
  - Follow this doc when implementing or reviewing related UI/UX or engineering patterns.
  - Update this doc when conventions change.
  - Refresh shared system pages after updating this doc.
---
```

### Notes
```
---
title: <Notes Title>
sop:
  - Read this doc when working on the related feature area.
  - Update this doc when implementation details or gotchas change.
  - Refresh shared system pages after updating this doc.
---
```

### Plans
```
---
title: <Plan Title>
sop:
  - Read this plan before executing the covered work.
  - Update this plan when execution steps or verification routes change.
  - Refresh shared system pages after updating this doc.
---
```

### Specs
```
---
title: <Spec Title>
sop:
  - Read this spec before changing the owned subsystem boundary.
  - Update this spec when shared contracts or maintenance routes change.
  - Refresh shared system pages after updating this doc.
---
```

### Runbook (optional)
```
---
title: <Runbook Title>
sop:
  - Read this doc before operational changes or incident response.
  - Update this doc when procedures or tooling change.
  - Refresh shared system pages after updating this doc.
---
```

### Manual Test (optional)
```
---
title: <Manual Test Title>
sop:
  - Run this checklist before release or when related features change.
  - Update this doc when verification steps change.
  - Refresh shared system pages after updating this doc.
---
```

## System Page Refresh
Refresh shared system pages after frontmatter changes with:
- `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .`
