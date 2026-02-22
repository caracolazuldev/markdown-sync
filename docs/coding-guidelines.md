# Coding Guidelines — markdown-sync

Go conventions
--------------
- Target Go 1.21. Use modules (`go.mod`).
- Follow `gofmt`, `go vet`, and `staticcheck` in CI.

Error handling
--------------
- Return wrapped errors with context (use `fmt.Errorf("...: %w", err)`).
- Define sentinel error types where callers must detect specific conditions.

Context
-------
- Accept `context.Context` on all functions that perform IO or network requests.

Dependencies
------------
- Prefer standard library when feasible. Pin third-party libs in `go.mod` and review licenses.

Testing
-------
- Unit test pure functions. Use fixtures in `internal/testdata` for mapping tests.
- Integration tests that require Google APIs must be gated and documented in `docs/testing.md`.
