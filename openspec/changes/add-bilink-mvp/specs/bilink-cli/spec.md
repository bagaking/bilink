## ADDED Requirements

### Requirement: Workspace Roots and File Filtering
The system SHALL scan a user-provided set of directory roots (a forest) and only process files matching configured extensions.

#### Scenario: Scan a directory forest
- **WHEN** the user runs any command that requires scanning
- **THEN** the system indexes only files under the specified roots with allowed extensions

### Requirement: Link Parsing and Classification
The system SHALL parse `[[...]]` and `[text](path)` links from Markdown files and classify strong vs weak associations.

#### Scenario: Parse wiki and markdown links
- **WHEN** a Markdown file contains `[[foo]]` and `[bar](docs/bar.md)`
- **THEN** both links are parsed and classified for downstream commands

### Requirement: Reference Queries
The system SHALL provide an explicit command to list outbound and inbound references for a file.

#### Scenario: Query refs
- **WHEN** the user runs `refs` for a file
- **THEN** the system outputs outbound and inbound references for that file

### Requirement: Rename Updates with Ambiguity Handling
The system SHALL update references when renaming or moving a file, and SHALL block on ambiguity unless the user resolves it interactively.

#### Scenario: Rename with unique target
- **WHEN** the user renames a file and all strong references resolve uniquely
- **THEN** references are updated automatically

#### Scenario: Rename with ambiguous target
- **WHEN** a rename would result in ambiguous references
- **THEN** the system reports an error unless `--interactive` is used to resolve choices

### Requirement: Conflict Checks and Lint Warnings
The system SHALL detect resolve-rule conflicts as errors and lint-rule violations as warnings.

#### Scenario: Report resolve conflicts and lint warnings
- **WHEN** the user runs `check`
- **THEN** resolve-rule conflicts are reported as errors and lint-rule violations as warnings

### Requirement: Watch Mode with TUI
The system SHALL provide a watch mode with an interactive TUI that reports changes and supports user choices.

#### Scenario: Watch mode interaction
- **WHEN** the user runs `watch`
- **THEN** the TUI displays progress, warnings, and interactive prompts for ambiguity

### Requirement: Configuration and Defaults
The system SHALL load configuration from `.bilink/settings.toml` when present and allow override via `--config`.

#### Scenario: Config override
- **WHEN** the user supplies `--config path/to/settings.toml`
- **THEN** the specified config is used instead of the default path

### Requirement: Structured Output
The system SHALL provide a `--json` option for machine-readable output in supported commands.

#### Scenario: JSON output
- **WHEN** the user passes `--json` to a supported command
- **THEN** the output is emitted as JSON

### Requirement: Index Persistence for Watch
The system SHALL use `.bilink/index.json` as persistent storage for watch-mode diffing.

#### Scenario: Watch requires index
- **WHEN** the user runs `watch`
- **THEN** `.bilink/index.json` is used to compare previous and current state
