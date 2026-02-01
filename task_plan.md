# Task Plan: Bilink MVP Implementation
<!-- 
  WHAT: This is your roadmap for the entire task. Think of it as your "working memory on disk."
  WHY: After 50+ tool calls, your original goals can get forgotten. This file keeps them fresh.
  WHEN: Create this FIRST, before starting any work. Update after each phase completes.
-->

## Goal
升级 `check` 输出为分组明细 + JSON 结构化，并将 `watch` TUI 升级为赛博朋克风格（含最近变更列表）。

## Current Phase
Phase 9

## Phases

### Phase 1: Requirements & Discovery
- [x] Review docs and OpenSpec instructions
- [x] Confirm existing specs/changes and constraints
- [x] Document findings in findings.md
- **Status:** complete

### Phase 2: Planning & Structure
- [x] Create OpenSpec change proposal (proposal/tasks/design/spec deltas)
- [x] Decide on implementation plan and worktree strategy
- [x] Document decisions with rationale
- **Status:** complete

### Phase 3: Implementation
- [x] Execute plan step by step (TDD)
- [x] Implement Go core, CLI, Watch TUI, TS wrapper
- [x] Test incrementally
- **Status:** complete

### Phase 4: Testing & Verification
- [x] Verify quality gates
- [x] Document test results in progress.md
- [x] Fix any issues found
- **Status:** complete

### Phase 5: Delivery
- [x] Review output files
- [x] Ensure deliverables are complete
- [x] Deliver summary to user
- **Status:** complete

### Phase 6: CLI Completion (Service Layer)
- [x] Design service layer and update docs
- [x] Implement service orchestration for refs/check/rename/watch
- [x] Implement JSON/text outputs
- **Status:** complete

### Phase 7: Testing & Verification
- [x] Verify quality gates
- [x] Document test results in progress.md
- [x] Fix any issues found
- **Status:** complete

### Phase 8: Delivery
- [x] Review output files
- [x] Ensure deliverables are complete
- [x] Deliver summary to user
- **Status:** complete

### Phase 9: Planning (Check Detail + Cyberpunk Watch)
- [x] Draft implementation plan doc
- [x] Update must-sop.md for new plan frontmatter
- **Status:** complete

### Phase 10: Implementation
- [ ] Implement check group details (text + json)
- [ ] Implement watch cyberpunk TUI (layout + change list)
- [ ] Update README
- **Status:** in_progress

### Phase 11: Testing & Verification
- [ ] Run quality gates (make lint/test + SOP check)
- [ ] Document results in progress.md
- **Status:** pending

## Key Questions
1. OpenSpec change-id and capability split (single spec vs multi-capability)
2. Worktree setup feasibility (repo has no commits)
3. Testing strategy and `make test` definition for Go/TS

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| Go core + TS wrapper + Bubble Tea TUI | Performance + npx distribution + UX |
| Balanced update policy (high/medium/low) | Recall/precision trade-off aligned with tooling norms |
| Default anchor style GitHub; resolve-only | Consistent behavior with minimal risk |
| Change id `add-bilink-mvp` with `bilink-cli` capability | Single capability for CLI feature set |
| Worktree path `.worktrees/add-bilink-mvp` | Isolated implementation workspace |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| here-doc parse error near `|` | 1 | Switched to apply_patch to add progress.md |

## Notes
- Update phase status as you progress: pending → in_progress → complete
- Re-read this plan before major decisions (attention manipulation)
- Log ALL errors - they help avoid repetition
