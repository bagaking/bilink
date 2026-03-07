# Feature Proposal: f-222nu8eev

## Why
- Bilink's core commands currently infer link identity in different places, so Markdown relative paths, inbound source reporting, unresolved links, and ambiguity handling drift across `refs`, `check`, `rename`, and `watch`.
- The next architectural step is to define one canonical edge model before expanding rename, watch, or agent-context features.

## Goal
- Define and plan Bilink's canonical link edge model, then use it to make refs/check report source-aware resolved and unresolved link relationships.

## Scope
- In scope:
  - OpenSpec change for canonical edge schema and refs/check behavior.
  - Edge fields for source path, syntax, serialized target, canonical target, anchor/subpath, confidence, candidates, and source context.
  - Markdown relative path resolution from the source file directory.
  - `refs` output that can answer who references a target and where.
  - `check` output for unresolved and ambiguous explicit links.
- Out of scope:
  - Two-phase rename engine.
  - Durable watch state machine.
  - Agent-native context/MCP layer.
  - Block refs and automatic anchor rewrites.

## Impact
- Code paths: `internal/parse`, `internal/index`, `internal/refs`, `internal/check`, `internal/service`, `internal/output`, and CLI JSON contracts.
- Tests: resolver/index unit tests plus black-box service fixtures for relative Markdown links, duplicate basenames, missing targets, and refs/check JSON shape.
- Rollout notes: start as `proposal_only`; implementation should not begin until the OpenSpec change is written and approved.
