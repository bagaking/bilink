# Feature Proposal: f-20260127-add-bilink-mvp

## Why
- Migrate the historical OpenSpec change `add-bilink-mvp` into bagakit-feature-tracker so the repository has one explicit feature-tracking SSOT.
- Preserve the original Bilink MVP planning truth without deleting or rewriting the existing `openspec/` history.

## Goal
- Deliver the Bilink MVP CLI, watch TUI, and npx wrapper as originally planned in OpenSpec.

## Scope
- In scope:
- Bilink CLI with `refs`, `rename`, `check`, `watch`, and `--json` output
- Link parsing, resolution rules, update policy, and optional index persistence
- Bubble Tea-based watch TUI
- TS wrapper for `npx`
- Out of scope:
- Block IDs (`^blockid`)
- Non-Markdown defaults
- Anchor auto-rewrite
- Remote/cloud sync

## Impact
- Code paths: `cmd/bilink`, `internal/*`, `packages/bilink-npx`, `docs/*`
- Tests: Go package tests, wrapper tests, quality-gate commands
- Rollout notes: Historical implementation already exists in the repository; this feature records that state in feature-tracker while keeping the OpenSpec source files in place.
