# MANIFEST — gdocs-markdown-sync

Version: 0.1.0
Status: draft

Purpose
-------
Provide a command-line tool that bridges Google Docs and Markdown to enable collaborative editing in Docs and reproducible, version-controlled artifacts in Markdown.

Primary Personas
-----------------
- Content authors: write in Google Docs, expect clean Markdown exports.
- Engineers: keep Markdown-driven repos, import or preview Docs content.
- CI/Automation: use headless service accounts to run imports/exports in pipelines.

Key Use Cases
-------------
- Export: convert a single-tab (or legacy) Google Doc into one Markdown file.
- Track: one-way pull of a tabbed Google Doc into a nested folder of Markdown files (`_track.toml`, pull-only).
- Import: push a Markdown file to a Google Doc (create or update content). Refuses tracked trees.
- Preview: render a Google Doc as Markdown without saving (fails on tabbed Docs).
- List: enumerate accessible Docs and their metadata.

Business Requirements
---------------------
- Manual-first: operations are explicit CLI commands.
- Dual auth support: interactive OAuth2 for users and service-account for CI.
- Fidelity: preserve structure (headings, paragraphs, code blocks, images, tables, links) for >=90% of typical documents in the sample suite.
- Deterministic mapping: round-trip conversions must be reproducible for supported features.

Non-functional Requirements
---------------------------
- Cross-platform binaries (Linux, macOS, Windows) built from Go 1.21.
- Reasonable performance for documents up to ~10k words; network-friendly API usage with retries and backoff.
- Secure handling of credentials: never store unencrypted keys in repo; follow Google best-practices.
- MIT license for repository.

- Technical Requirements
- ----------------------
- Language: Go (minimum 1.21).
- Module: `github.com/caracolazuldev/gdocs-markdown-sync` (placeholder).
- CLI commands: `export`, `import`, `preview`, `list`, `track` (see docs/cli-spec.md).
- Auth: support both OAuth2 user-consent and service account flows (see docs/auth.md).
- Mapping rules: documented in docs/mapping.md (headings, paragraphs, code blocks, images, tables, links, frontmatter).
- Tests: unit tests for mapping and adapters; gated integration tests against Google Docs (service account or test user), described in docs/testing.md.

Acceptance Criteria
--------------------
- Unit test coverage for mapping logic > 80%.
- Exported Markdown for sample docs matches approved fixtures in docs/examples/ with structural parity in 90% of cases.
- CLI builds on Linux/macos/windows and runs basic `list` and `preview` flows with interactive OAuth.

Supporting Documents (drafts)
-----------------------------
- docs/cli-spec.md — concrete CLI surface and examples.
- docs/architecture.md — component responsibilities and deployment notes.
- docs/auth.md — OAuth2 and service-account flows, scopes.
- docs/mapping.md — deterministic mapping rules and limitations.
- docs/coding-guidelines.md — Go-specific coding conventions and rules for contributors.
- docs/testing.md — testing strategy and CI gating for integration tests.
 - docs/approach.md — guidance and records for implementation/algorithmic approaches and when to use them.
 - docs/decisions.md — recorded architecture and business-logic decisions and approach entries.
 - agents/references/design-patterns.md — project-relevant design patterns and practical notes.
 - docs/testing.md — testing strategy and CI gating for integration tests.
 - agents/ — AI coding agent resources (skills, workflows, references). See agents/index.yaml for discovery.
 - AGENTS.md — cross-tool agent instructions (devcontainer-only toolchain).

Roadmap and Next Work
----------------------
v0.1: Manual export/import/preview + auth flows + mapping for core features.
v0.2: One-way `track` for tabbed Docs (ADR-2026-08-13-track). Bidirectional tab sync is declined (see Missives/2026-08-13-response-tabbed-doc-directory-sync.md). Comments/suggestions and watchers remain later.

CI and Automation
-----------------
Continuous integration, automated release workflows, and gated integration tests are deferred until the core CLI functionality and mapping rules stabilize. For now, contributors should use the `Makefile` targets for local build, lint, and test workflows. CI will be added later and tracked in the project plan (task: "CI setup (deferred)").

Contact
-------
Owner: repository maintainer (TBD)

Decisions
---------
- See `docs/decisions.md` for recorded architecture and business-logic decisions; PRs updating architecture or business logic must add or update a decision record.
- ADR-2026-08-18-devcontainer-toolchain — agents use `.devcontainer/Dockerfile`; do not install Go on the host.
- ADR-2026-08-13-track — one-way `track`; `export` fails on tabbed Docs.
- APP-2026-08-13-tab-paths — tab tree to nested folder paths.
- Product spec: `Missives/2026-08-13-spec-track-tabbed-docs.md` (references the decline missive; does not amend it).
