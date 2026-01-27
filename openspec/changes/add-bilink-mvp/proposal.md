# Change: Add Bilink MVP (CLI + Watch TUI + npx wrapper)

## Why
Provide a robust local bidirectional-link tool that can scan Markdown, query refs, update links on rename, and watch changes with an interactive TUI.

## What Changes
- Add Bilink CLI with `refs`, `rename`, `check`, `watch` and `--json` output
- Add link parsing, resolution rules, and update policy (strong vs weak)
- Add optional index persistence for watch mode
- Add Bubble Tea-based Watch TUI
- Add TS wrapper for `npx` distribution

## Impact
- Affected specs: `bilink-cli`
- Affected code: new Go CLI, parsing/index/update modules, TUI, TS wrapper, configuration loader
