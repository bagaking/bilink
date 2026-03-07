# Feature Proposal: f-20260127-bilink-frontier-research

## Why
- Replace ad hoc root-level planning files with tracked feature-tracker state.
- Preserve the completed Bilink frontier research and next-plan work as durable planning truth tied to the repository.

## Goal
- Build reusable Bilink research evidence, summarize frontier practices, identify optimization space, and produce the next implementation plan.

## Scope
- In scope:
- Local-first Bilink system design research workspace
- Source cards, reusable summaries, topic index, and implementation plan
- Migration of root `task_plan.md`, `findings.md`, and `progress.md` into feature-tracker
- Out of scope:
- Direct product-code implementation
- OpenSpec execution for the next phase

## Impact
- Code paths: none directly; planning truth and researcher assets only
- Tests: tracker validation and doctor checks
- Rollout notes: root-level planning scratch files are removed after migration; feature-tracker and researcher become the durable planning surfaces
