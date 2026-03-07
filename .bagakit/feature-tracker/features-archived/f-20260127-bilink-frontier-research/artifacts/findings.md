# Findings & Decisions

## Source Snapshot

This file preserves the content that previously lived at the repository root.

- Original path: `findings.md`
- Migration target: `feature-tracker/artifacts`

## Requirements
- Create a durable local-first research workspace for Bilink.
- Find frontier practices and representative high-quality sources.
- For AI-related work, explicitly collect Anthropic-written materials where they matter.
- Save source refs and write summaries in the project directory.
- Compare the findings against Bilink and identify optimization space.
- End with a proposed next implementation plan.

## Research Axes
- Bidirectional linking products and semantics
- Local-first / indexing / watch / rename integrity architecture
- AI-friendly documentation and agent-oriented project structure

## Research Findings
- The repository configured a local researcher root.
- Bagakit researcher requires per-topic `originals/`, `summaries/`, and `index.md`.
- The repository does not contain `scripts/bagakit-researcher.sh`, so the workspace must be managed manually.
- User explicitly requested subagent delegation; parallel research is appropriate because the three axes are independent.
- Research evidence was saved under the research topic workspace.
- The strongest recurring pattern across sources is: `canonical target model` first, `plan/apply refactor` second, `durable watch cursor/snapshot` third.
- The most important AI-facing conclusion is that Bilink should evolve toward an agent context and evidence substrate, not remain only a link maintenance CLI.
- The recommended next implementation sequence is: OpenSpec change -> canonical edge model -> refs/check upgrade -> two-phase rename -> real watch state -> agent-native extension layer.

## Technical Decisions
| Decision | Rationale |
|----------|-----------|
| Topic workspace for Bilink system design | Covers the user’s requested frontier-oriented design research |
| Preserve source cards first, then write summaries, then synthesize | Matches bagakit-researcher evidence workflow |
| Use docs/PRD/architecture as the comparison baseline for gap analysis | Project intent is already documented there |
| Keep research artifacts in the researcher runtime instead of `docs/` | Matches configured researcher root and avoids mixing evidence workspace with user-facing project docs |

## Issues Encountered
| Issue | Resolution |
|-------|------------|
| Missing helper script for bagakit-researcher | Use the skill’s documented directory contract directly |

## Resources
- `.bagakit/knowledge_conf.toml`
- `docs/notes-prd-bilink.md`
- `docs/architecture-bilink-design.md`
- `bagakit-researcher/SKILL.md`
- `bagakit-researcher/references/research-workspace-spec.md`
- Bilink system design research topic index
- Bilink system design research implementation plan

## Visual/Browser Findings
- None yet.
