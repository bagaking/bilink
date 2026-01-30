# Findings & Decisions

## Requirements
- CLI supports `refs`, `rename`, `check`, `watch` with `--json` output.
- Parse `[[...]]` and `[text](path)` links in Markdown files.
- Strong vs weak association; auto-update only for strong + unique targets.
- `#anchor` resolved for navigation only; changes warn, no auto-rewrite.
- Default extensions: `.md`, `.markdown`, `.mdx`; configurable.
- Config: optional `.bilink/settings.toml`, override via `--config`.
- Optional `.bilink/index.json` for watch mode diffing.
- Go core, Bubble Tea TUI, TS wrapper for npx.
- Quality gates: `make lint`, `make test`, SOP regeneration check.

## Research Findings
- OpenSpec required for new capabilities; must create proposal and get approval before implementation.
- Worktree creation may be blocked if repo has no commits.
- Must read and follow docs/must-sop.md; SOP regeneration required after frontmatter changes.
- `openspec/project.md` is a template with no project-specific details.
- `openspec list --specs` shows no existing specs.
- `openspec list` shows no active changes; `openspec/changes` only has archive/.
- `openspec validate add-bilink-mvp --strict` passed.
- No `.worktrees/` or `worktrees/` directory present; `CLAUDE.md` missing.
- Git HEAD missing (no commits yet); worktree creation will require an initial commit.
- Worktree created at `.worktrees/add-bilink-mvp`; implementation plan saved at `docs/plans/2026-01-27-bilink-mvp.md`.
- SOP regenerated to include implementation plan.

## Technical Decisions
| Decision | Rationale |
|----------|-----------|
| Go core + TS wrapper + Bubble Tea TUI | Performance + distribution + UX |
| Balanced update policy | Precision/recall balance aligned with tooling norms |
| Default anchor style GitHub | Widely compatible; minimal risk |
| Change id `add-bilink-mvp` with `bilink-cli` capability | Single spec for CLI scope |
| Keep go.mod at 1.24.0 | golangci-lint requires go mod tidy; tidy sets go version to 1.24.0 with current toolchain |
| GitHub Releases asset naming | `bilink-<os>-<arch>` (Windows adds `.exe`) |
| CLI completion approach | Add `internal/service` layer (Option B) to orchestrate commands |

## Issues Encountered
| Issue | Resolution |
|-------|------------|
| planning-with-files templates needed | Created task_plan.md and seeded findings/progress |
| make lint failed due to module needing tidy | Ran go mod tidy and retained go.mod at 1.24.0 |

## Resources
- `docs/notes-prd-bilink.md`
- `docs/architecture-bilink-design.md`
- `docs/must-guidebook.md`
- `docs/must-docs-taxonomy.md`
- `openspec/AGENTS.md`

## Visual/Browser Findings
- None.
