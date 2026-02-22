# MANIFEST — markdown-sync

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
- Export: convert a Google Doc into one or more Markdown files with assets saved locally.
- Import: push a Markdown file to a Google Doc (create or update content).
- Preview: render a Google Doc as Markdown without saving.
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

Technical Requirements
----------------------
- Language: Go (minimum 1.21).
- Module: `github.com/caracolazuldev/markdown-sync` (placeholder).
- CLI commands: `export`, `import`, `preview`, `list` (see docs/cli-spec.md).
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

Roadmap and Next Work
----------------------
v0.1: Manual export/import/preview + auth flows + mapping for core features.

Contact
-------
Owner: repository maintainer (TBD)
