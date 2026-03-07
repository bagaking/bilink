# Feature Tasks: f-20260127-add-bilink-mvp

JSON SSOT: `tasks.json`

## Task Checklist
- [x] T-001 Initialize Go module and baseline CLI entrypoint
- [x] T-002 Add config loader for `.bilink/settings.toml` and defaults
- [x] T-003 Implement Markdown parser for `[[...]]` and `[text](path)`
- [x] T-004 Build in-memory index (files, outbound refs, inbound refs)
- [x] T-005 Implement `refs` command (text + `--json` output)
- [x] T-006 Implement resolve/lint rules and conflict detection for `check`
- [x] T-007 Implement rename/update engine (strong links only, ambiguity handling)
- [x] T-008 Implement watch mode with Bubble Tea TUI and index persistence
- [x] T-009 Add TS wrapper package that downloads Go binaries for `npx`
- [x] T-010 Add tests per feature (TDD) and update docs

## Status Legend
- todo
- in_progress
- done
- blocked
