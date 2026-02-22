# Testing — markdown-sync

Unit tests
----------
- Keep unit tests fast and deterministic. Mock Google adapters and only test mapping logic and local storage behavior.

Integration tests
-----------------
- Integration tests exercise Google Docs API and require credentials. They must be gated behind CI secrets and marked `// +build integration` or run via `go test -tags=integration`.
- Provide instructions for creating test documents and service accounts in this doc.

CI
--
- CI runs unit tests for every PR. Integration tests run on a separate pipeline with secrets provided.

Local dev
---------
- For local integration testing, set `GOOGLE_APPLICATION_CREDENTIALS` to a service account JSON and use `go test -tags=integration ./...`.
