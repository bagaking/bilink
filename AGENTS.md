<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

<!-- BAGAKIT:LIVEDOCS:START -->
Legacy note: this repository now uses `bagakit-living-knowledge` as the shared
knowledge substrate.
Use the managed `BAGAKIT:LIVING-KNOWLEDGE` block below for bootstrap and refresh
shared system pages with:
- `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .`
At the end of every response, include a single line:
- `[[Bagakit.LivingDoc]] <short note about which shared knowledge rule was followed>`
<!-- BAGAKIT:LIVEDOCS:END -->

## Commit Rule

For commits made under this migration effort, use `git-commit-nt` instead of
calling `git commit` directly.

Usage:

```bash
git-commit-nt <commit-message> [git-commit-args...]
```

Do not bypass this wrapper unless the user explicitly overrides the rule.
<!-- BAGAKIT:LIVING-KNOWLEDGE:START -->
This is a managed block for `bagakit-living-knowledge`. Do not hand-edit the
managed region directly; refresh it through the skill operator instead.

Resolve the installed skill dir before using the operator directly:

- `export BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR="<path-to-bagakit-living-knowledge-skill>"`

Boot layer:

- Read the resolved `must-guidebook.md` before relying on memory.
- If a task needs shared knowledge rules, read `must-authority.md`.
- If a task needs maintenance-route guidance or shared directives, read `must-sop.md`.
- If a task needs prior decisions or facts, follow `must-recall.md`.
- `AGENTS.md` is only the bootstrap layer; the shared checked-in knowledge root
  is configured in `.bagakit/knowledge_conf.toml`.

Recall discipline:

- Search first:
  - `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" recall search --root . '<query>'`
- Then inspect only the needed lines:
  - `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" recall get --root . <path> --from <line> --lines <n>`
- Prefer quoting only needed lines over paraphrasing from memory.

Substrate discipline:

- Shared knowledge belongs under the configured shared root.
- Research runtime belongs to `bagakit-researcher`.
- Task-level composition/runtime belongs to `bagakit-skill-selector`.
- Repository evolution memory belongs to `bagakit-skill-evolver`.
- `living-knowledge` owns path protocol, normalization, indexing, and recall.
- `living-knowledge` also owns generated `must-sop.md` and reusable-items
  governance inside the shared knowledge root.

Inspection helpers:

- View the resolved path protocol:
  - `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" paths --root .`
- Refresh the guidebook and helper map:
  - `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" index --root .`
- Run non-destructive diagnostics:
  - `sh "$BAGAKIT_LIVING_KNOWLEDGE_SKILL_DIR/scripts/bagakit-living-knowledge.sh" doctor --root .`

If the surrounding workflow explicitly asks for `living-knowledge` task
reporting, the response footer may use:

- `[[BAGAKIT]]`
- `- LivingKnowledge: Surface=<updated shared surfaces or none>; Evidence=<commands/checks>; Next=<one deterministic next action>`
<!-- BAGAKIT:LIVING-KNOWLEDGE:END -->
