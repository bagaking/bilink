## Context
Bilink MVP adds a local CLI for bidirectional links across a directory forest, including parsing, index, ref queries, rename updates, conflict checks, and watch TUI. See `docs/notes-prd-bilink.md` and `docs/architecture-bilink-design.md` for product and architecture details.

## Goals / Non-Goals
- Goals: Strong/weak link model, resolve rules vs lint rules, optional index persistence, interactive watch TUI, npx wrapper.
- Non-Goals: Block IDs (`^blockid`), non-Markdown defaults, anchor auto-rewrite, remote/cloud sync.

## Decisions
- Decision: Go core for parsing/index/update + Bubble Tea TUI
- Decision: TS wrapper for npx distribution
- Decision: Balanced update policy (high/medium auto-update; low warn)
- Decision: Default anchor style GitHub, resolve-only (no auto-rewrite)

## Risks / Trade-offs
- Ambiguity can block automation → default to error + interactive choice
- Index persistence adds state → required for watch, optional for static analysis
- Watch TUI adds dependency → Bubble Tea chosen for mature components

## Migration Plan
- New project; no migrations required

## Open Questions
- None (defaults set in `docs/architecture-bilink-design.md`)
