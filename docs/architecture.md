# Architecture — gdocs-markdown-sync

Summary
-------
`gdocs-markdown-sync` is a small CLI with a modular architecture to keep Google API concerns, markdown transformation, and local track layout separate.

Commands
--------
- `export` / `preview` / `import` / `list`: one Google Doc ID ↔ one Markdown file (legacy body or a single leaf tab).
- `track`: one-way Docs → nested local folder for **tabbed** documents. Pull-only. See `Missives/2026-08-13-spec-track-tabbed-docs.md` and ADR-2026-08-13-track.

`export` and `preview` fail closed on tabbed documents (more than one tab at any depth, or any `childTabs`) and direct the user to `track`. `import` refuses paths under a track root.

Components
----------
- CLI (`cmd/gdocs-markdown-sync`): argument parsing, UX, and wiring. Keep this thin.
- Auth Layer (`internal/google`): OAuth2 and service-account credential flows; authenticated HTTP client.
- Markdown adapter (`internal/markdown`): document model, Markdown conversion, fetch/apply against the Docs API. Fetch is injectable (`Fetcher`) so unit tests never call the network. Production CLI uses the API fetcher (`documents.get` with `includeTabsContent=true`).
- Sync orchestrator (`internal/sync`): `track` (tab-tree → nested paths, `_track.toml`, chmod cycle, dirty/`--force`, first-run guard) and import refusal for tracked trees.
- Test utilities: fakes implementing `markdown.Fetcher`; fixtures for tab trees.

The following were planned in earlier drafts and are **not present**: `internal/google/docs_adapter.go`, `internal/storage`. Fetch/apply live in `internal/markdown`. Track file IO lives in `internal/sync`.

Deployment & Packaging
----------------------
- Build cross-compiled native binaries using `goreleaser` or `xgo` in CI.
- Provide a Docker image for CI and headless usage.

Observability & Reliability
---------------------------
- Structured logging with configurable verbosity.
- Retries with exponential backoff for transient Google API errors.
- Sentinel errors (`ErrTabbedDocument`, `ErrTrackedPath`, `ErrDirtyTrack`, `ErrNotTrackRoot`, `ErrDocIDMismatch`) for stable CLI hints.

Security
--------
- Follow least-privilege: request minimal OAuth scopes.
- Use `GOOGLE_APPLICATION_CREDENTIALS` for service-account-based CI.
