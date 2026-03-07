# Task Plan: Bilink Frontier Research and Next Plan

## Source Snapshot

This file preserves the content that previously lived at the repository root.

- Original path: `task_plan.md`
- Migration target: `feature-tracker/artifacts`

## Goal
围绕 Bilink 的项目目标，建立一套本地可复用的研究证据库，汇总前沿实践与代表性文章，抽取可借鉴模式，识别优化空间，并输出下一轮实现计划。

## Current Phase
Phase 5

## Phases

### Phase 1: Scope and Workspace Setup
- [x] Confirm research goal and repository constraints
- [x] Read relevant skills and Bagakit researcher rules
- [x] Decide topic taxonomy and research axes
- **Status:** complete

### Phase 2: Evidence Collection
- [x] Create topic workspace under the researcher runtime
- [x] Search and collect representative sources across product / architecture / AI workflow domains
- [x] Preserve source cards for selected materials
- **Status:** complete

### Phase 3: Summarization
- [x] Write reusable per-source summaries
- [x] Refresh the topic index with reading order
- [x] Consolidate cross-source patterns
- **Status:** complete

### Phase 4: Gap Analysis
- [x] Compare Bilink current design/implementation with research findings
- [x] Identify product and architecture optimization opportunities
- [x] Rank opportunities by impact and dependency
- **Status:** complete

### Phase 5: Next Implementation Plan
- [x] Draft a concrete follow-up implementation plan
- [x] Align plan with current repo constraints and docs
- [x] Deliver summary to user with references
- **Status:** complete

## Key Questions
1. 对 Bilink 最有代表性的外部参照系是什么？
2. 哪些实践对“本地文件 + 双链 + rename 安全 + agent 友好”最关键？
3. 当前实现和前沿实践相比，最大的缺口在哪里？

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| Use Bagakit researcher layout directly under the configured researcher root | `knowledge_conf.toml` already configures this root |
| Split research into 3 axes: bidirectional-link products, local-first/indexing architecture, AI-friendly documentation/agent workflows | Matches Bilink’s product goals and design claims |
| Use subagents for parallel online research but keep file writing local | User explicitly requested a subagent team; local integration avoids write conflicts |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| `scripts/bagakit-researcher.sh` missing from repo | 1 | Fall back to manual workspace creation following the skill spec |
| Prior review planning files absent in current worktree | 1 | Recreate research-specific planning files for this turn |

## Notes
- Prefer primary sources and official docs/blogs
- Keep research evidence local and reusable
- Re-read this plan before final recommendations
