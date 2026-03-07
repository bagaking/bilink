# Feature Summary: f-20260127-bilink-frontier-research

- Title: Bilink Frontier Research and Next Plan
- Status: done
- Workspace mode: current_tree
- Related researcher topic: Bilink system design research workspace

## Outcome

This feature captured the completed research/planning loop for Bilink:

- representative product-semantics sources
- local-first and watch/index architecture sources
- Anthropic and agent-workflow sources
- a consolidated gap analysis
- a phased next implementation plan

## Main Conclusions

- Bilink should first become correct around canonical edges, not broader in syntax surface.
- Safe rename should move to a `plan -> preview -> verify -> apply` model.
- Watch should become a durable state machine with cursor/snapshot integrity rather than a one-shot diff.
- The long-term differentiated direction is an agent context and evidence substrate, not only a link-maintenance CLI.

## Next

Use the researcher topic index and `implementation-plan.md` as the source of truth for the next OpenSpec change.
