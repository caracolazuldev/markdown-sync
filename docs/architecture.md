# Architecture — markdown-sync

Summary
-------
`markdown-sync` is a small CLI with a modular architecture to keep Google API concerns, markdown transformation, and storage separate.

Components
----------
- CLI (`cmd/markdown-sync`): argument parsing, UX, and orchestration.
- Auth Layer (`internal/google`): handle OAuth2 and service-account credential flows and provide an authenticated HTTP client.
- Google Docs Adapter (`internal/google/docs_adapter.go`): thin wrapper around the Google Docs REST API; handles fetch/patch operations and retries.
- Markdown Engine (`internal/markdown`): deterministic conversion rules between Docs structured content and Markdown text (+asset extraction).
- Sync Orchestrator (`internal/sync`): high-level operations for export/import/preview and conflict policies.
- Storage (`internal/storage`): local file IO, asset caching, and path mapping.
- Test utilities (`internal/testutil`): fixtures, mocks, and small in-memory fakes for unit tests.

Deployment & Packaging
----------------------
- Build cross-compiled native binaries using `goreleaser` or `xgo` in CI.
- Provide a Docker image for CI and headless usage.

Observability & Reliability
---------------------------
- Structured logging with configurable verbosity.
- Retries with exponential backoff for transient Google API errors.
- Clear error types to enable programmatic handling by higher-level scripts.

Security
--------
- Follow least-privilege: request minimal OAuth scopes.
- Use `GOOGLE_APPLICATION_CREDENTIALS` for service-account-based CI.
